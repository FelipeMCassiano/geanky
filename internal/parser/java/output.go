package java

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

// extractClassName pega um import completo (ex: java.util.List) e retorna só a classe (List)
func extractClassName(importPath string) string {
	parts := strings.Split(importPath, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return importPath
}

// formatModifiers transforma a slice de modificadores em string
// func formatModifiers(modifiers []Modifier) string {
// 	var sb strings.Builder
// 	for _, m := range modifiers {
// 		sb.WriteString(m.Modifier)
// 		sb.WriteString(" ")
// 	}
// 	return sb.String()
// }

// exprToText é uma função auxiliar para formatar expressões em texto legível
func exprToText(expr Expression) string {
	if expr == nil {
		return ""
	}
	defer func() { recover() }() // Evita panic em tipos inesperados
	switch e := expr.(type) {
	case Identifier:
		return e.Name
	case Literal:
		return e.Value
	case Access:
		obj := exprToText(e.Object)
		if obj != "" {
			return obj + "." + e.Identifier.Name
		}
		return e.Identifier.Name
	case MethodInvocation:
		var args []string
		for _, a := range e.Args {
			args = append(args, exprToText(a))
		}
		return fmt.Sprintf("%s(%s)", exprToText(e.Accessed), strings.Join(args, ", "))
	case Binary:
		return fmt.Sprintf("%s %s %s", exprToText(e.Left), e.Operator, exprToText(e.Right))
	case Assignment:
		return fmt.Sprintf("%s = %s", exprToText(e.Left), exprToText(e.Right))
	case ReturnNode:
		return exprToText(e.Value)
	case Variable:
		if e.Value != nil {
			return fmt.Sprintf("%s %s = %s", e.TypeName, e.Declarator, exprToText(e.Value))
		}
		return fmt.Sprintf("%s %s", e.TypeName, e.Declarator)
	case ThrowNode:
		return fmt.Sprintf("throw %s", exprToText(e.Value))
	default:
		return "..."
	}
}

// // formatExpression formata o node para uma string pseudo-code clean na lista de Step-by-Step
// func formatExpression(expr Expression) string {
// 	if expr == nil {
// 		return ""
// 	}
// 	switch e := expr.(type) {
// 	case IfNode:
// 		return fmt.Sprintf("If %s then execute block", exprToText(e.Condition))
// 	case ForNode:
// 		return fmt.Sprintf("Loop (for %s)", exprToText(e.Condition))
// 	case EnhancedForNode:
// 		return fmt.Sprintf("Loop (for each %s in %s)", e.Name, exprToText(e.Value))
// 	case WhileNode:
// 		return fmt.Sprintf("Loop (while %s)", exprToText(e.Condition))
// 	case TryNode:
// 		return "Try executing block"
// 	case ThrowNode:
// 		return fmt.Sprintf("Throw exception: %s", exprToText(e.Value))
// 	case BreakNode:
// 		return "Break loop"
// 	case Variable:
// 		return fmt.Sprintf("Declare variable: %s", exprToText(e))
// 	case ReturnNode:
// 		return fmt.Sprintf("Return: %s", exprToText(e.Value))
// 	default:
// 		return exprToText(e)
// 	}
// }

// getDependencyCalls varre os métodos de uma classe e mapeia dependências
func getDependencyCalls(c ClassJava) map[string]string {
	fieldsMap := make(map[string]string)
	for _, f := range c.Fields {
		fieldsMap[f.Declarator] = f.TypeName
	}

	deps := make(map[string]map[string]bool)

	var traverse func(expr Expression)
	traverse = func(expr Expression) {
		if expr == nil {
			return
		}
		switch e := expr.(type) {
		case Assignment:
			traverse(e.Left)
			traverse(e.Right)
		case Variable:
			traverse(e.Value) // Extrai chamadas dentro de var = obj.metodo()
		case Binary:
			traverse(e.Left)
			traverse(e.Right)
		case IfNode:
			traverse(e.Condition)
			for _, s := range e.Consequence.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			if e.Alternative != nil {
				for _, s := range e.Alternative.Statements {
					for _, ex := range s.Expressions {
						traverse(ex)
					}
				}
			}
		case ForNode:
			traverse(e.Init)
			traverse(e.Condition)
			traverse(e.Update)
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
		case EnhancedForNode:
			traverse(e.Value)
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
		case WhileNode:
			traverse(e.Condition)
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
		case TryNode:
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			for _, c := range e.Catches {
				for _, s := range c.Body.Statements {
					for _, ex := range s.Expressions {
						traverse(ex)
					}
				}
			}
		case ThrowNode:
			traverse(e.Value)
		case ReturnNode:
			traverse(e.Value)
		case MethodInvocation:
			targetObj := ""
			if e.Accessed.Object != nil {
				switch obj := e.Accessed.Object.(type) {
				case Access:
					targetObj = obj.Identifier.Name
				case Identifier:
					targetObj = obj.Name
				}
			}

			if typeName, exists := fieldsMap[targetObj]; exists {
				if deps[typeName] == nil {
					deps[typeName] = make(map[string]bool)
				}

				var args []string
				for _, arg := range e.Args {
					switch a := arg.(type) {
					case Identifier:
						args = append(args, a.Name)
					default:
						args = append(args, "...")
					}
				}

				methodSignature := fmt.Sprintf("%s(%s)", e.Accessed.Identifier.Name, strings.Join(args, ", "))
				deps[typeName][methodSignature] = true
			}

			for _, arg := range e.Args {
				traverse(arg)
			}
		}
	}

	for _, m := range c.Methods {
		for _, s := range m.Body.Statements {
			for _, e := range s.Expressions {
				traverse(e)
			}
		}
	}

	result := make(map[string]string)
	for typeName, methods := range deps {
		var methodList []string
		for m := range methods {
			methodList = append(methodList, m)
		}
		result[typeName] = strings.Join(methodList, "<br>")
	}
	return result
}

// generateSequenceDiagram constrói dinamicamente um diagrama de sequência Mermaid
func generateSequenceDiagram(m Executable) string {
	var sb strings.Builder
	participants := make(map[string]bool)

	sb.WriteString("sequenceDiagram\n")
	sb.WriteString("    actor Caller\n")
	sb.WriteString("    participant ThisClass\n")

	ensureParticipant := func(name string) {
		if name == "" || name == "this" || name == "ThisClass" || name == "Caller" {
			return
		}
		if !participants[name] {
			participants[name] = true
			sb.WriteString(fmt.Sprintf("    participant %s\n", name))
		}
	}

	cleanForMermaid := func(s string) string {
		s = strings.ReplaceAll(s, "\n", " ")
		s = strings.ReplaceAll(s, "\r", "")
		s = strings.ReplaceAll(s, "\"", "'")
		if len(s) > 60 {
			s = s[:57] + "..."
		}
		return s
	}

	var resolveTarget func(expr Expression) string
	resolveTarget = func(expr Expression) string {
		if expr == nil {
			return ""
		}
		switch o := expr.(type) {
		case Identifier:
			if o.Name == "this" {
				return ""
			}
			return o.Name
		case Access:
			if inner := resolveTarget(o.Object); inner != "" {
				return inner
			}
			return o.Identifier.Name
		default:
			return ""
		}
	}

	var traverse func(expr Expression)
	traverse = func(expr Expression) {
		if expr == nil {
			return
		}
		defer func() { recover() }()

		switch e := expr.(type) {
		case Assignment:
			traverse(e.Right)
		case Variable:
			traverse(e.Value)
		case Binary:
			traverse(e.Left)
			traverse(e.Right)
		case IfNode:
			cond := cleanForMermaid(exprToText(e.Condition))
			if cond == "" {
				cond = "condition"
			}
			sb.WriteString(fmt.Sprintf("    alt %s\n", cond))
			for _, s := range e.Consequence.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			if e.Alternative != nil {
				sb.WriteString("    else\n")
				for _, s := range e.Alternative.Statements {
					for _, ex := range s.Expressions {
						traverse(ex)
					}
				}
			}
			sb.WriteString("    end\n")
		case ForNode:
			cond := cleanForMermaid(exprToText(e.Condition))
			sb.WriteString(fmt.Sprintf("    loop for %s\n", cond))
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			sb.WriteString("    end\n")
		case EnhancedForNode:
			val := cleanForMermaid(exprToText(e.Value))
			sb.WriteString(fmt.Sprintf("    loop for each %s in %s\n", e.Name, val))
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			sb.WriteString("    end\n")
		case WhileNode:
			cond := cleanForMermaid(exprToText(e.Condition))
			sb.WriteString(fmt.Sprintf("    loop while %s\n", cond))
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			sb.WriteString("    end\n")
		case TryNode:
			sb.WriteString("    alt try\n")
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			for _, c := range e.Catches {
				sb.WriteString(fmt.Sprintf("    else catch %s\n", cleanForMermaid(c.Parameter)))
				for _, s := range c.Body.Statements {
					for _, ex := range s.Expressions {
						traverse(ex)
					}
				}
			}
			sb.WriteString("    end\n")
		case ThrowNode:
			val := cleanForMermaid(exprToText(e.Value))
			sb.WriteString(fmt.Sprintf("    ThisClass-->>Caller: throw %s\n", val))
		case MethodInvocation:
			target := resolveTarget(e.Accessed.Object)
			if target != "" {
				ensureParticipant(target)
				var args []string
				for _, a := range e.Args {
					args = append(args, exprToText(a))
				}
				callArgs := cleanForMermaid(strings.Join(args, ", "))
				sb.WriteString(fmt.Sprintf("    ThisClass->>%s: %s(%s)\n", target, e.Accessed.Identifier.Name, callArgs))
			}
			for _, a := range e.Args {
				traverse(a)
			}
		case ReturnNode:
			val := cleanForMermaid(exprToText(e.Value))
			sb.WriteString(fmt.Sprintf("    ThisClass-->>Caller: return %s\n", val))
		}
	}

	var params []string
	for _, p := range m.Parameters {
		params = append(params, p.Declarator)
	}

	initCall := cleanForMermaid(strings.Join(params, ", "))
	sb.WriteString(fmt.Sprintf("\n    Caller->>ThisClass: %s(%s)\n", m.Name, initCall))

	for _, s := range m.Body.Statements {
		for _, e := range s.Expressions {
			traverse(e)
		}
	}

	return sb.String()
}

const docTemplate = `
{{range .Annotations}}> **{{.}}**
{{end}}# 📄 Technical Specification: {{bt}}{{.Name}}{{bt}}

{{if .Package.Name}}> **Package:** {{.Package.Name}}
{{end}}{{if .Imports}}> **Dependencies (Imports):**
{{range .Imports}}> - {{if isProjectClass .}}[{{.}}]({{extractClassName .}}.md) 🔗{{else}}{{.}}{{end}}
{{end}}{{end}}> **Automatically generated documentation** by the Geanky tool.

---

## 1. Quick Summary (API & State)
A high-level overview of the class, its internal state, and available methods.

**Internal State & Dependencies:**
{{if not .Fields}}> *No state properties defined.*
{{else}}{{range .Fields}}
{{if isProjectClass .TypeName}}
- {{range .Annotations}}{{bt}}{{.}}{{bt}} {{end}}{{bt}}{{formatModifiers .Modifiers}}{{bt}} **{{.Declarator}}** ([{{.TypeName}}]({{.TypeName}}.md)) 🔗
{{else}}
- {{range .Annotations}}{{bt}}{{.}}{{bt}} {{end}}{{bt}}{{formatModifiers .Modifiers}}{{bt}} **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{end}}{{end}}{{end}}

**Available Methods:**
{{if not .Methods}}> *No methods defined.*
{{else}}{{range .Methods}}- **{{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.TypeName}} {{$p.Declarator}}{{end}})** ➞ returns {{bt}}{{.ReturnType}}{{bt}}{{if .Throws}} (throws {{.Throws}}){{end}}
{{end}}{{end}}

---

## 2. Class Dependencies & State
Visual representation of the internal state and external dependencies this class maintains.

{{bt}}{{bt}}{{bt}}mermaid
flowchart LR
    classDef classNode fill:#2b3137,stroke:#fff,stroke-width:2px,color:#fff;
    classDef stateNode fill:#f4f6f8,stroke:#d0d7de,color:#24292f;
    classDef extNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    ThisClass["{{.Name}}"]:::classNode

    {{range .Fields}}
    {{if or (eq .TypeName "String") (eq .TypeName "int") (eq .TypeName "boolean") (eq .TypeName "double") (eq .TypeName "long") (eq .TypeName "float")}}
    ThisClass -- "Maintains State" --- State_{{.Declarator}}(["{{.TypeName}}<br>{{.Declarator}}"]):::stateNode
    {{else}}
    ThisClass -- "Depends on" ---> Dep_{{.Declarator}}["{{.TypeName}}"]:::extNode
    {{end}}
    {{end}}
{{bt}}{{bt}}{{bt}}

---

## 3. Deep Dive (Constructors & Methods)

{{if .Constructors}}
### 🛠️ Constructors
{{range .Constructors}}
<details>
<summary><b>{{.Name}}</b>({{range $i, $p := .Parameters}}{{if $i}}, {{end}}<i>{{$p.TypeName}}</i> {{$p.Declarator}}{{end}}) (Click to expand)</summary>

> **Signature:**
{{range .Annotations}}> {{bt}}{{.}}{{bt}}
{{end}}> {{bt}}{{formatModifiers .Modifiers}}{{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{range $p.Annotations}}{{.}} {{end}}{{$p.TypeName}} {{$p.Declarator}}{{end}}){{if .Throws}} throws {{.Throws}}{{end}}{{bt}}

**Sequence Diagram:**
{{bt}}{{bt}}{{bt}}mermaid
{{generateSequenceDiagram .}}
{{bt}}{{bt}}{{bt}}

**Step-by-Step Logic:**
{{if not .Body.Statements}}> *Empty body.*
{{else}}
{{range .Body.Statements}}{{range .Expressions}}
1. {{formatExpression .}}
{{end}}{{end}}
{{end}}
</details>
{{end}}
{{end}}

{{if .Methods}}
### ⚙️ Methods
{{range .Methods}}
<details>
<summary><b>{{.Name}}</b>({{range $i, $p := .Parameters}}{{if $i}}, {{end}}<i>{{$p.TypeName}}</i> {{$p.Declarator}}{{end}}) ➞ {{bt}}{{.ReturnType}}{{bt}} (Click to expand)</summary>

> **Signature:**
{{range .Annotations}}> {{bt}}{{.}}{{bt}}
{{end}}> {{bt}}{{formatModifiers .Modifiers}}{{.ReturnType}} {{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{range $p.Annotations}}{{.}} {{end}}{{$p.TypeName}} {{$p.Declarator}}{{end}}){{if .Throws}} throws {{.Throws}}{{end}}{{bt}}

**Sequence Diagram:**
{{bt}}{{bt}}{{bt}}mermaid
{{generateSequenceDiagram .}}
{{bt}}{{bt}}{{bt}}

**Step-by-Step Logic:**
{{if not .Body.Statements}}> *Empty body.*
{{else}}
{{range .Body.Statements}}{{range .Expressions}}
1. {{formatExpression .}}
{{end}}{{end}}
{{end}}
</details>
{{end}}
{{end}}
`

const globalDocTemplate = `
# 🌍 Global Architecture Diagram

> Visão geral de alto nível mostrando as dependências entre todas as classes analisadas e seus respectivos pacotes.

{{bt}}{{bt}}{{bt}}mermaid
flowchart LR
    classDef classNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    {{range $pkgName, $pkgClasses := .GroupedClasses}}
    subgraph {{$pkgName}}
        {{range $pkgClasses}}
        {{.Name}}["{{.Name}}"]:::classNode
        {{end}}
    end
    {{end}}

    {{range .AllClasses}}
    {{$className := .Name}}
    {{$callsMap := getDependencyCalls .}}
    
    {{range .Fields}}
    {{if not (or (eq .TypeName "String") (eq .TypeName "int") (eq .TypeName "boolean") (eq .TypeName "double") (eq .TypeName "long") (eq .TypeName "float"))}}
        {{$usedMethods := index $callsMap .TypeName}}
        {{if $usedMethods}}
            {{$className}} -->|"Calls:<br><b>{{$usedMethods}}</b>"| {{.TypeName}}
        {{else}}
            {{$className}} -->|"Depends on"| {{.TypeName}}
        {{end}}
    {{end}}
    {{end}}
    {{end}}
{{bt}}{{bt}}{{bt}}
`

func GenerateMarkdown(classData ClassJava, allClasses []ClassJava, outputFilename string) {
	isProjectClass := func(importPath string) bool {
		className := extractClassName(importPath)
		for _, c := range allClasses {
			if c.Name == className {
				return true
			}
		}
		return false
	}

	tmpl, err := template.New("classDoc").Funcs(template.FuncMap{
		"bt":                      func() string { return "`" },
		"formatModifiers":         formatModifiers,
		"formatExpression":        formatExpression,
		"extractClassName":        extractClassName,
		"isProjectClass":          isProjectClass,
		"generateSequenceDiagram": generateSequenceDiagram,
	}).Parse(docTemplate)

	if err != nil {
		log.Fatalf("Erro ao criar template: %v", err)
	}

	file, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, classData)
	if err != nil {
		log.Fatalf("Erro ao executar template: %v", err)
	}
}

func GenerateGlobalArchitecture(classes []ClassJava, outputFilename string) {
	tmpl, err := template.New("globalDoc").Funcs(template.FuncMap{
		"bt":                 func() string { return "`" },
		"getDependencyCalls": getDependencyCalls,
	}).Parse(globalDocTemplate)

	if err != nil {
		log.Fatalf("Erro ao criar template global: %v", err)
	}

	file, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo global: %v", err)
	}
	defer file.Close()

	groupedClasses := make(map[string][]ClassJava)
	for _, c := range classes {
		pkgName := "Default Package"
		if c.Package.Name != "" {
			pkgName = c.Package.Name
		}
		groupedClasses[pkgName] = append(groupedClasses[pkgName], c)
	}

	templateData := struct {
		AllClasses     []ClassJava
		GroupedClasses map[string][]ClassJava
	}{
		AllClasses:     classes,
		GroupedClasses: groupedClasses,
	}

	err = tmpl.Execute(file, templateData)
	if err != nil {
		log.Fatalf("Erro ao gerar documentação global: %v", err)
	}

	fmt.Printf("🗺️ Diagrama de Arquitetura Global gerado em: %s\n", outputFilename)
}
