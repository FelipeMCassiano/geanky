package java

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

const docTemplate = `
# 📄 Especificação Técnica: {{bt}}{{.Name}}{{bt}}

> **Documentação gerada automaticamente** pela ferramenta Geanky.
> Mapeamento estrutural, injeção de dependências e regras de inicialização.

---

## 1. Visão Geral e Estado Interno
Esta seção detalha as propriedades de estado (fields) armazenadas na instância de {{.Name}}.

{{if not .Fields}}> *Nenhuma propriedade de estado definida para esta classe.*
{{else}}| Acesso | Nome da Propriedade | Tipo da Variável | Descrição / Uso |
| :--- | :--- | :--- | :--- |
{{range .Fields}}| {{bt}}private{{bt}} | **{{bt}}{{.Declarator}}{{bt}}** | {{bt}}{{.TypeName}}{{bt}} | Armazena o estado interno da propriedade. |
{{end}}{{end}}
---

## 2. Construtores e Injeção de Dependência
{{if not .Constructors}}> *Nenhum construtor customizado detectado. A classe utiliza o construtor padrão vazio.*
{{end}}{{range .Constructors}}
### 🛠️ Inicializador: {{.Name}}
Instancia e prepara a classe para uso, exigindo a injeção dos seguintes parâmetros:

**Assinatura de Entrada:**
{{if not .Parameters}}> *Este construtor não recebe parâmetros de entrada.*
{{else}}| Parâmetro | Tipo | Finalidade / Origem |
| :--- | :--- | :--- | :--- |
{{range .Parameters}}| **{{bt}}{{.Declarator}}{{bt}}** | {{bt}}{{.TypeName}}{{bt}} | Valor injetado durante a instanciação. |
{{end}}{{end}}

**Fluxo de Execução (Regras do Construtor):**
{{if not .Body.Statements}}> *Corpo de execução vazio (Nenhuma instrução executada).*
{{else}}{{range .Body.Statements}}{{range .Expressions}}1. Executa a seguinte instrução lógica: {{bt}}{{formatExpression .}}{{bt}}
{{end}}{{end}}{{end}}
{{end}}
---

## 3. Comportamento e Regras de Negócio (Métodos)
{{if not .Methods}}
> ℹ️ **Aviso de Arquitetura:** Nenhum método ou API pública foi detectado para esta classe no escopo atual.
{{else}}{{range .Methods}}### ⚙️ Função: {{.Name}}()
*(A implementação interna deste método precisa ser mapeada)*
{{end}}{{end}}
`

func GenerateMarkdown(classData ClassJava, outputFilename string) {
	tmpl, err := template.New("classDoc").Funcs(template.FuncMap{
		"formatExpression": formatExpression,
		"bt":               func() string { return "`" }, // Adicionamos isso aqui!
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
