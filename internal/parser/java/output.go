package java

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

// O template com Quick Summary, Diagrama Mermaid e Deep Dive Expansível
const docTemplate = `
# 📄 Technical Specification: {{bt}}{{.Name}}{{bt}}

> **Automatically generated documentation** by the Geanky tool.

---

## 1. Quick Summary (API & State)
A high-level overview of the class, its internal state, and available methods.

**Internal State:**
{{if not .Fields}}> *No state properties defined.*
{{else}}{{range .Fields}}- {{bt}}{{formatModifiers .Modifiers}}{{bt}} **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{end}}{{end}}

**Available Methods:**
{{if not .Methods}}> *No methods defined.*
{{else}}{{range .Methods}}- **{{.Name}}()** ➞ returns {{bt}}{{.ReturnType}}{{bt}}
{{end}}{{end}}

---

## 2. Data Flow & Navigation Diagram
Visual representation of how data enters the class, how it is stored, and what is returned to external consumers.

{{bt}}{{bt}}{{bt}}mermaid
flowchart LR
    %% Styling
    classDef mainClass fill:#2b3137,stroke:#fff,stroke-width:2px,color:#fff;
    classDef stateNode fill:#f4f6f8,stroke:#d0d7de,color:#24292f;
    
    Caller((External<br>Caller))
    ThisClass[{{.Name}}]:::mainClass

    %% API Interactions (Inputs and Outputs)
    {{if .Methods}}{{range .Methods}}
    Caller -->|Calls {{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.TypeName}}{{end}})| ThisClass
    ThisClass -.->|Returns {{.ReturnType}}| Caller
    {{end}}{{end}}

    %% Internal State (Storage)
    {{if .Fields}}{{range .Fields}}
    ThisClass ---|Maintains State| State_{{.Declarator}}([{{.TypeName}} {{.Declarator}}]):::stateNode
    {{end}}{{end}}
{{bt}}{{bt}}{{bt}}

---

## 3. Detailed Execution Flow (Deep Dive)
Expand the sections below to read the exact pseudo-code and business rules inside each method.

{{if not .Methods}}> *No methods available to detail.*
{{else}}{{range .Methods}}
<details>
<summary><b>⚙️ Function: {{.Name}}</b> (Click to expand)</summary>

> **Signature:** {{bt}}{{formatModifiers .Modifiers}}{{.ReturnType}} {{.Name}}(){{bt}}

**Parameters:**
{{if not .Parameters}}> *None.*
{{else}}{{range .Parameters}}- **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{end}}{{end}}

**Step-by-Step Logic:**
{{if not .Body.Statements}}> *Method body empty or not implemented.*
{{else}}{{range .Body.Statements}}{{range .Expressions}}1. {{bt}}{{formatExpression .}}{{bt}}
{{end}}{{end}}{{end}}

</details>
{{end}}{{end}}
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
