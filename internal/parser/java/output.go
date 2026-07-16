package java

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

const docTemplate = `
# 📄 Technical Specification: {{bt}}{{.Name}}{{bt}}

> **Automatically generated documentation** by the Geanky tool.

---

## 1. Quick Summary (API & State)
A high-level overview of the class, its internal state, and available methods.

**Internal State & Dependencies:**
{{if not .Fields}}> *No state properties defined.*
{{else}}{{range .Fields}}- {{bt}}{{formatModifiers .Modifiers}}{{bt}} **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{end}}{{end}}

**Available Methods:**
{{if not .Methods}}> *No methods defined.*
{{else}}{{range .Methods}}- **{{.Name}}()** ➞ returns {{bt}}{{.ReturnType}}{{bt}}
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
    Caller -- "Calls {{.Name}}()" --> ThisClass
    ThisClass -. "Returns {{.ReturnType}}" .-> Caller
    {{end}}

    %% State vs External Dependencies
    {{range .Fields}}
    {{if or (eq .TypeName "String") (eq .TypeName "int") (eq .TypeName "boolean") (eq .TypeName "double")}}
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
<summary><b>{{.Name}}</b> (Click to expand)</summary>

**Parameters:**
{{if not .Parameters}}> *None.*
{{else}}{{range .Parameters}}- **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{end}}{{end}}

**Step-by-Step Logic:**
{{if not .Body.Statements}}> *Empty body.*
{{else}}{{range .Body.Statements}}{{range .Expressions}}1. {{bt}}{{formatExpression .}}{{bt}}
{{end}}{{end}}{{end}}
</details>
{{end}}
{{end}}

{{if .Methods}}
### ⚙️ Methods
{{range .Methods}}
<details>
<summary><b>{{.Name}}</b> ➞ {{bt}}{{.ReturnType}}{{bt}} (Click to expand)</summary>

> **Signature:** {{bt}}{{formatModifiers .Modifiers}}{{.ReturnType}} {{.Name}}(){{bt}}

**Parameters:**
{{if not .Parameters}}> *None.*
{{else}}{{range .Parameters}}- **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{end}}{{end}}

**Step-by-Step Logic:**
{{if not .Body.Statements}}> *Empty body.*
{{else}}{{range .Body.Statements}}{{range .Expressions}}1. {{bt}}{{formatExpression .}}{{bt}}
{{end}}{{end}}{{end}}
</details>
{{end}}
{{end}}
`

// GenerateMarkdown compila o template com os dados mapeados da AST
func GenerateMarkdown(classData ClassJava, outputFilename string) {
	tmpl, err := template.New("classDoc").Funcs(template.FuncMap{
		"formatExpression": formatExpression,
		"formatModifiers":  formatModifiers,
		"bt":               func() string { return "`" }, // Hack para imprimir backticks no template
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
		log.Fatalf("Erro ao gerar documentação: %v", err)
	}

	fmt.Printf("✅ Documentação gerada com sucesso em: %s\n", outputFilename)
}
