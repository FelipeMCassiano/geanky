package java

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

const docTemplate = `{{- define "executable" -}}
### {{formatModifiers .Modifiers}}{{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.TypeName}} {{$p.Declarator}}{{end}})
 
**📥 Contrato**
 
{{if .Parameters -}}
| Parâmetro | Tipo | Modificadores |
|:---|:---|:---|
{{range .Parameters}}| <code>{{.Declarator}}</code> | <code>{{.TypeName}}</code> | {{range .Modifiers}}<code>{{.Modifier}}</code> {{end}} |
{{end -}}
{{else -}}
> Nenhum parâmetro de entrada.
{{end}}
{{if .ReturnType}}**↩️ Retorna:** <code>{{.ReturnType}}</code>{{else}}**🚫 Retorna:** nada (<code>void</code>){{end}}
 
**⚙️ Comportamento**
 
{{if .Body.Statements -}}
{{range .Body.Statements}}{{range .Expressions}}- [x] {{formatExpression .}}
{{end}}{{end -}}
{{else -}}
> Nenhuma lógica interna definida neste bloco.
{{end}}
<br>
 
---
 
{{end -}}
# 📦 {{range .Modifiers}}<code>{{.Modifier}}</code> {{end}}{{.Name}}
 
> Documentação gerada automaticamente a partir da análise estática do código-fonte.
 
<p>
{{range .Modifiers}}<img src="https://img.shields.io/badge/-{{.Modifier}}-6f42c1?style=flat-square" alt="{{.Modifier}}" /> {{end}}<img src="https://img.shields.io/badge/fields-{{len .Fields}}-informational?style=flat-square" /> <img src="https://img.shields.io/badge/constructors-{{len .Constructors}}-informational?style=flat-square" /> <img src="https://img.shields.io/badge/methods-{{len .Methods}}-informational?style=flat-square" />
</p>
 
---
 
## 🧩 Atributos
 
{{if .Fields -}}
| Modificadores | Tipo | Nome |
|:---|:---|:---|
{{range .Fields}}| {{range .Modifiers}}<code>{{.Modifier}}</code> {{end}} | <code>{{.TypeName}}</code> | <code>{{.Declarator}}</code> |
{{end -}}
{{else -}}
> Esta classe não possui atributos declarados.
{{end}}
 
---
 
{{if .Constructors -}}
## 🏗️ Construtores
 
{{range .Constructors}}{{template "executable" .}}
{{end}}
{{end -}}
{{if .Methods -}}
## ⚙️ Métodos
 
{{range .Methods}}{{template "executable" .}}
{{end}}
{{end -}}
`

func GenerateMarkdown(classData ClassJava, outputFilename string) {
	tmpl, err := template.New("classDoc").Funcs(template.FuncMap{
		"formatExpression": formatExpression,
		// Adicionamos o formatModifiers para exibir 'public', 'private' corretamente
		"formatModifiers": formatModifiers,
		"bt":              func() string { return "`" },
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
