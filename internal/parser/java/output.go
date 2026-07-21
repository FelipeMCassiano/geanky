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
			traverse(e.Value)
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
	var header strings.Builder
	var steps strings.Builder
	participants := make(map[string]bool)

	header.WriteString("sequenceDiagram\n")
	header.WriteString("    actor Caller\n")
	header.WriteString("    participant ThisClass\n")

	ensureParticipant := func(name string) {
		if name == "" || name == "this" || name == "ThisClass" || name == "Caller" {
			return
		}
		cleanName := strings.TrimSpace(name)
		cleanName = strings.Split(cleanName, " ")[0]
		cleanName = strings.Split(cleanName, ".")[0]
		cleanName = strings.Split(cleanName, "(")[0]

		// Remove colchetes e símbolos de genéricos que quebram a engine do Mermaid
		cleanName = strings.ReplaceAll(cleanName, "<", "")
		cleanName = strings.ReplaceAll(cleanName, ">", "")
		cleanName = strings.ReplaceAll(cleanName, "[", "")
		cleanName = strings.ReplaceAll(cleanName, "]", "")

		if cleanName == "" {
			return
		}

		if !participants[cleanName] {
			participants[cleanName] = true
			header.WriteString(fmt.Sprintf("    participant %s\n", cleanName))
		}
	}

	cleanForMermaid := func(s string) string {
		s = strings.ReplaceAll(s, "\n", " ")
		s = strings.ReplaceAll(s, "\r", "")
		s = strings.ReplaceAll(s, "\t", " ")

		s = strings.ReplaceAll(s, "\"", "'")

		s = strings.ReplaceAll(s, "&lt;", "<")
		s = strings.ReplaceAll(s, "&gt;", ">")
		s = strings.ReplaceAll(s, "&amp;", "&")

		s = strings.Join(strings.Fields(s), " ")

		// Limita o tamanho para não estourar o diagrama visualmente
		if len(s) > 60 {
			s = s[:57] + "..."
		}
		return strings.TrimSpace(s)
	}

	var resolveTarget func(expr Expression) string
	resolveTarget = func(expr Expression) string {
		if expr == nil {
			return ""
		}
		switch o := expr.(type) {
		case Identifier:
			if o.Name == "this" || o.Name == "super" {
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

	var exprToString func(expr Expression) string
	exprToString = func(expr Expression) string {
		if expr == nil {
			return ""
		}
		defer func() { recover() }()
		switch e := expr.(type) {
		case Identifier:
			return e.Name
		case Literal:
			return e.Value
		case Access:
			obj := exprToString(e.Object)
			if obj != "" {
				return obj + "." + e.Identifier.Name
			}
			return e.Identifier.Name
		case MethodInvocation:
			var args []string
			for _, a := range e.Args {
				args = append(args, exprToString(a))
			}
			return fmt.Sprintf("%s(%s)", exprToString(e.Accessed), strings.Join(args, ", "))
		case Binary:
			return fmt.Sprintf("%s %s %s", exprToString(e.Left), e.Operator, exprToString(e.Right))
		case Assignment:
			return exprToString(e.Right)
		case Variable:
			return exprToString(e.Value)
		case ThrowNode:
			return exprToString(e.Value)
		case ReturnNode:
			return exprToString(e.Value)
		default:
			return "..."
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
			cond := cleanForMermaid(exprToString(e.Condition))
			if cond == "" {
				cond = "condition"
			}
			steps.WriteString(fmt.Sprintf("    alt %s\n", cond))
			for _, s := range e.Consequence.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			if e.Alternative != nil {
				steps.WriteString("    else\n")
				for _, s := range e.Alternative.Statements {
					for _, ex := range s.Expressions {
						traverse(ex)
					}
				}
			}
			steps.WriteString("    end\n")
		case ForNode:
			cond := cleanForMermaid(exprToString(e.Condition))
			if cond == "" {
				cond = "true"
			}
			steps.WriteString(fmt.Sprintf("    loop for %s\n", cond))
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			steps.WriteString("    end\n")
		case EnhancedForNode:
			val := cleanForMermaid(exprToString(e.Value))
			steps.WriteString(fmt.Sprintf("    loop for each %s in %s\n", e.Name, val))
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			steps.WriteString("    end\n")
		case WhileNode:
			cond := cleanForMermaid(exprToString(e.Condition))
			if cond == "" {
				cond = "true"
			}
			steps.WriteString(fmt.Sprintf("    loop while %s\n", cond))
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			steps.WriteString("    end\n")
		case TryNode:
			steps.WriteString("    alt try\n")
			for _, s := range e.Body.Statements {
				for _, ex := range s.Expressions {
					traverse(ex)
				}
			}
			for _, c := range e.Catches {
				steps.WriteString(fmt.Sprintf("    else catch %s\n", cleanForMermaid(c.Parameter)))
				for _, s := range c.Body.Statements {
					for _, ex := range s.Expressions {
						traverse(ex)
					}
				}
			}
			steps.WriteString("    end\n")
		case ThrowNode:
			val := cleanForMermaid(exprToString(e.Value))
			steps.WriteString(fmt.Sprintf("    ThisClass-->>Caller: throw %s\n", val))
		case BreakNode:
			steps.WriteString("    Note right of ThisClass: break loop\n")
		case MethodInvocation:
			target := resolveTarget(e.Accessed.Object)
			if target == "" {
				target = "ThisClass"
			} else {
				target = strings.Split(target, " ")[0]
				target = strings.Split(target, ".")[0]
				target = strings.Split(target, "(")[0]

				target = strings.ReplaceAll(target, "<", "")
				target = strings.ReplaceAll(target, ">", "")
				target = strings.ReplaceAll(target, "[", "")
				target = strings.ReplaceAll(target, "]", "")
			}
			ensureParticipant(target)

			var args []string
			for _, a := range e.Args {
				args = append(args, exprToString(a))
			}

			callArgs := cleanForMermaid(strings.Join(args, ", "))

			// AQUI É A MUDANÇA: Exibe objeto.metodo(...) se não for da própria classe
			methodLabel := e.Accessed.Identifier.Name
			if target != "ThisClass" {
				methodLabel = target + "." + methodLabel
			}

			steps.WriteString(fmt.Sprintf(
				"    ThisClass->>%s: %s(%s)\n",
				target,
				methodLabel,
				callArgs,
			))

			for _, a := range e.Args {
				traverse(a)
			}
		case ReturnNode:
			val := cleanForMermaid(exprToString(e.Value))
			steps.WriteString(fmt.Sprintf("    ThisClass-->>Caller: return %s\n", val))
		default:
			// Nó não mapeado: ignorado silenciosamente
		}
	}

	var params []string
	for _, p := range m.Parameters {
		params = append(params, p.Declarator)
	}

	initCall := cleanForMermaid(strings.Join(params, ", "))
	steps.WriteString(fmt.Sprintf("\n    Caller->>ThisClass: %s(%s)\n", m.Name, initCall))

	for _, s := range m.Body.Statements {
		for _, e := range s.Expressions {
			traverse(e)
		}
	}

	header.WriteString(steps.String())
	return header.String()
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
		"formatModifiers":         formatModifiers,  // Certifique-se que você tem essa func no seu build
		"formatExpression":        formatExpression, // Certifique-se que você tem essa func no seu build
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
			pkgName = c.Package.Scope + c.Package.Name
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
