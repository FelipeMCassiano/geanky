package java

import (
	"fmt"
	"log"
	"os"
	"strings"
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
{{else}}{{range .Fields}}
{{if or (eq .TypeName "String") (eq .TypeName "int") (eq .TypeName "boolean") (eq .TypeName "double") (eq .TypeName "long") (eq .TypeName "float")}}
- {{bt}}{{formatModifiers .Modifiers}}{{bt}} **{{.Declarator}}** ({{bt}}{{.TypeName}}{{bt}})
{{else}}
- {{bt}}{{formatModifiers .Modifiers}}{{bt}} **{{.Declarator}}** ([{{.TypeName}}]({{.TypeName}}.md)) 🔗
{{end}}{{end}}{{end}}

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
    classDef extNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff,cursor:pointer;
    
    Caller(("Caller"))
    ThisClass["{{.Name}}"]:::classNode

    %% Method Calls
    {{range .Methods}}
    Caller -- "Calls {{.Name}}()" --> ThisClass
    ThisClass -. "Returns {{.ReturnType}}" .-> Caller
    {{end}}

    %% State vs External Dependencies
    {{range .Fields}}
    {{if or (eq .TypeName "String") (eq .TypeName "int") (eq .TypeName "boolean") (eq .TypeName "double") (eq .TypeName "long") (eq .TypeName "float")}}
    ThisClass -- "Maintains State" --- State_{{.Declarator}}(["{{.TypeName}} {{.Declarator}}"]):::stateNode
    {{else}}
    ThisClass -- "Depends on" ---> Dep_{{.Declarator}}["{{.TypeName}}"]:::extNode
    click Dep_{{.Declarator}} "{{.TypeName}}.md" "Ir para a documentação de {{.TypeName}}"
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
<summary><b>{{.Name}}</b> ➞ {{bt}}{{.ReturnType}}{{bt}} (Click to expand)</summary>

> **Signature:** {{bt}}{{formatModifiers .Modifiers}}{{.ReturnType}} {{.Name}}(){{bt}}

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

const globalDocTemplate = `
# 🌍 Global Architecture Diagram

> Visão geral de alto nível mostrando as dependências entre todas as classes analisadas.

{{bt}}{{bt}}{{bt}}mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff,cursor:pointer;

    %% Nodes Creation
    {{range .}}
    {{.Name}}["{{.Name}}"]:::classNode
    click {{.Name}} "{{.Name}}.md" "Acessar {{.Name}}"
    {{end}}

    %% Relationships / Dependencies
    {{range .}}
    {{$className := .Name}}
    {{$callsMap := getDependencyCalls .}} {{/* CHAMA A NOSSA FUNCÃO GO! (Agora de forma segura) */}}
    
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

	err = tmpl.Execute(file, classes)
	if err != nil {
		log.Fatalf("Erro ao gerar documentação global: %v", err)
	}

	fmt.Printf("🗺️ Diagrama de Arquitetura Global gerado em: %s\n", outputFilename)
}

// getDependencyCalls varre os métodos de uma classe e mapeia: NomeDaDependencia -> "metodo1(), metodo2()"
func getDependencyCalls(c ClassJava) map[string]string {
	// 1. Mapeia as variáveis da classe (ex: "service" -> "UsuarioService")
	fieldsMap := make(map[string]string)
	for _, f := range c.Fields {
		fieldsMap[f.Declarator] = f.TypeName
	}

	// 2. Prepara um mapa para guardar os métodos únicos chamados por tipo
	deps := make(map[string]map[string]bool)

	// 3. Função recursiva para varrer expressões infinitamente
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
			// Verifica se o objeto dono do método é uma dependência nossa
			targetObj := ""
			if e.Accessed.Object != nil {
				switch obj := e.Accessed.Object.(type) {
				case Access:
					targetObj = obj.Identifier.Name // Pega o "service" de "this.service"
				case Identifier:
					targetObj = obj.Name // Pega o "service" direto
				}
			}

			// Se o objeto for uma dependência mapeada, registra o método chamado!
			if typeName, exists := fieldsMap[targetObj]; exists {
				if deps[typeName] == nil {
					deps[typeName] = make(map[string]bool)
				}
				deps[typeName][e.Accessed.Identifier.Name] = true
			}

			// Continua varrendo os argumentos do método
			for _, arg := range e.Args {
				traverse(arg)
			}
		case ReturnNode:
			traverse(e.Value)
		}
	}

	// 4. Inicia a varredura em todos os métodos da classe
	for _, m := range c.Methods {
		for _, s := range m.Body.Statements {
			for _, e := range s.Expressions {
				traverse(e)
			}
		}
	}

	// 5. Formata a saída (junta os métodos com vírgula)
	result := make(map[string]string)
	for typeName, methods := range deps {
		var methodList []string
		for m := range methods {
			methodList = append(methodList, m+"()")
		}
		result[typeName] = strings.Join(methodList, "<br>") // <br> quebra a linha na seta do Mermaid!
	}
	return result
}
