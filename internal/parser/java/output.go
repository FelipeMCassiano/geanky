package java

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

// getDependencyCalls varre os métodos de uma classe e mapeia: NomeDaDependencia -> "metodo(arg1, arg2)"
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

				// Extrai os argumentos formatados para exibir no diagrama
				var args []string
				for _, arg := range e.Args {
					argStr := formatExpression(arg)
					// Troca aspas duplas por simples para não quebrar a label do Mermaid
					argStr = strings.ReplaceAll(argStr, "\"", "'")
					args = append(args, argStr)
				}

				// Monta a assinatura completa: metodo(arg1, arg2)
				methodSignature := fmt.Sprintf("%s(%s)", e.Accessed.Identifier.Name, strings.Join(args, ", "))
				deps[typeName][methodSignature] = true
			}

			for _, arg := range e.Args {
				traverse(arg)
			}
		case ReturnNode:
			traverse(e.Value)
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

// formatModifiers transforma a slice de modificadores em string

const docTemplate = `
# 📄 Technical Specification: {{bt}}{{.Name}}{{bt}}

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
{{if or (eq .TypeName "String") (eq .TypeName "int") (eq .TypeName "boolean") (eq .TypeName "double") (eq .TypeName "long") (eq .TypeName "float")}}
- {{bt}}{{formatModifiers .Modifiers}}{{bt}} **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{else}}
- {{bt}}{{formatModifiers .Modifiers}}{{bt}} **{{.Declarator}}** ([{{.TypeName}}]({{.TypeName}}.md)) 🔗
{{end}}{{end}}{{end}}

**Available Methods:**
{{if not .Methods}}> *No methods defined.*
{{else}}{{range .Methods}}- **{{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.TypeName}} {{$p.Declarator}}{{end}})** ➞ returns {{bt}}{{.ReturnType}}{{bt}}
{{end}}{{end}}

---

## 2. Architecture & Data Flow Diagram
Visual representation of how data enters the class, internal state, and external dependencies.

{{bt}}{{bt}}{{bt}}mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#2b3137,stroke:#fff,stroke-width:2px,color:#fff;
    classDef stateNode fill:#f4f6f8,stroke:#d0d7de,color:#24292f;
    classDef extNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    Caller(("Caller"))
    ThisClass["{{.Name}}"]:::classNode

    %% Method Calls
    {{range .Methods}}
    Caller -- "Calls {{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.TypeName}} {{$p.Declarator}}{{end}})" --> ThisClass
    ThisClass -. "Returns {{.ReturnType}}" .-> Caller
    {{end}}

    %% State vs External Dependencies
    {{range .Fields}}
    {{if or (eq .TypeName "String") (eq .TypeName "int") (eq .TypeName "boolean") (eq .TypeName "double") (eq .TypeName "long") (eq .TypeName "float")}}
    ThisClass -- "Maintains State" --- State_{{.Declarator}}(["{{.TypeName}} {{.Declarator}}"]):::stateNode
    {{else}}
    ThisClass -- "Depends on" ---> Dep_{{.Declarator}}["{{.TypeName}}"]:::extNode
    {{end}}
    {{end}}
{{bt}}{{bt}}{{bt}}

---

## 3. Deep Dive (Constructors & Methods)
Expand the sections below to read the exact pseudo-code and business rules.

{{if .Constructors}}
### 🛠️ Constructors
{{range .Constructors}}
<details>
<summary><b>{{.Name}}</b>({{range $i, $p := .Parameters}}{{if $i}}, {{end}}<i>{{$p.TypeName}}</i> {{$p.Declarator}}{{end}}) (Click to expand)</summary>

**Parameters:**
{{if not .Parameters}}> *None.*
{{else}}{{range .Parameters}}
- **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{end}}{{end}}

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

> **Signature:** {{bt}}{{formatModifiers .Modifiers}}{{.ReturnType}} {{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.TypeName}} {{$p.Declarator}}{{end}}){{bt}}

**Parameters:**
{{if not .Parameters}}> *None.*
{{else}}{{range .Parameters}}
- **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{end}}{{end}}

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
    %% Styling
    classDef classNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    %% Nodes Creation Grouped by Package
    {{range $pkgName, $pkgClasses := .GroupedClasses}}
    subgraph {{$pkgName}}
        {{range $pkgClasses}}
        {{.Name}}["{{.Name}}"]:::classNode
        {{end}}
    end
    {{end}}

    %% Relationships / Dependencies
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
		"bt":               func() string { return "`" },
		"formatModifiers":  formatModifiers,
		"formatExpression": formatExpression,
		"extractClassName": extractClassName,
		"isProjectClass":   isProjectClass,
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

func extractClassName(importPath string) string {
	parts := strings.Split(importPath, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return importPath
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

	// Agrupa as classes pelo nome completo do pacote obtido na AST
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
