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
> Structural mapping, dependency injection, and initialization rules.

---

## 1. Overview and Internal State
This section details the state properties (fields) stored in the class instance.

{{if not .Fields}}> *No state properties defined for this class.*
{{else}}| Modifiers | Property Name | Variable Type |
| :--- | :--- | :--- |
{{range .Fields}}| {{bt}}{{formatModifiers .Modifiers}}{{bt}} | **{{bt}}{{.Declarator}}{{bt}}** | {{bt}}{{.TypeName}}{{bt}} |
{{end}}{{end}}

---

## 2. Constructors and Dependency Injection
{{if not .Constructors}}> *No custom constructor detected. The class uses the default empty constructor.*
{{end}}{{range .Constructors}}
### 🛠️ Initializer: {{bt}}{{.Name}}{{bt}}
Instantiates and prepares the class for use.

**Input Parameters:**
{{if not .Parameters}}> *This constructor takes no parameters.*
{{else}}| Parameter | Type |
| :--- | :--- |
{{range .Parameters}}| **{{bt}}{{.Declarator}}{{bt}}** | {{bt}}{{.TypeName}}{{bt}} |
{{end}}{{end}}

**Execution Flow (Constructor Rules):**
{{if not .Body.Statements}}> *Empty execution body.*
{{else}}{{range .Body.Statements}}{{range .Expressions}}1. {{bt}}{{formatExpression .}}{{bt}}
{{end}}{{end}}{{end}}
{{end}}

---

## 3. Behavior and Business Rules (Methods)
Below are all the class methods detailed, including their signatures and the literal step-by-step of what they do internally.

{{if not .Methods}}
> ℹ️ **Architecture Warning:** No methods were detected for this class.
{{else}}{{range .Methods}}
### ⚙️ Function: {{bt}}{{.Name}}{{bt}}
* **Return:** {{bt}}{{.ReturnType}}{{bt}}
* **Visibility/Modifiers:** {{bt}}{{formatModifiers .Modifiers}}{{bt}}

**Received Parameters:**
{{if not .Parameters}}> *Method without parameters.*
{{else}}| Parameter Name | Type |
| :--- | :--- |
{{range .Parameters}}| **{{bt}}{{.Declarator}}{{bt}}** | {{bt}}{{.TypeName}}{{bt}} |
{{end}}{{end}}

**Internal Execution Flow:**
{{if not .Body.Statements}}> *Method body empty or not implemented.*
{{else}}{{range .Body.Statements}}{{range .Expressions}}1. {{bt}}{{formatExpression .}}{{bt}}
{{end}}{{end}}{{end}}
---
{{end}}{{end}}
`

func GenerateMarkdown(classData ClassJava, outputFilename string) {
	tmpl, err := template.New("classDoc").Funcs(template.FuncMap{
		"formatExpression": formatExpression,
		"formatModifiers":  formatModifiers,
		"bt":               func() string { return "`" },
	}).Parse(docTemplate)

	if err != nil {
		log.Fatalf("Error creating template: %v", err)
	}

	file, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, classData)
	if err != nil {
		log.Fatalf("Error generating documentation: %v", err)
	}

	fmt.Printf("✅ Documentation successfully generated at: %s\n", outputFilename)
}
