package parser

import (
	"fmt"
	"log"
	"os"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	java "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

type Method struct {
	Name string `json:"name"`
}
type Field struct {
	TypeName   string `json:"typeName"`
	Declarator string `json:"declarator"`
}

type ClassJava struct {
	Name         string        `json:"name"`
	Constructors []Constructor `json:"constructors"`
	Fields       []Field       `json:"fields"`
	Methods      []Method      `json:"methods"`
}

type Constructor struct {
	Name       string  `json:"name"`
	Parameters []Field `json:"parameters"`
}

func Analyze(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	parser := tree_sitter.NewParser()
	defer parser.Close()
	lang := tree_sitter.NewLanguage(java.Language())
	parser.SetLanguage(lang)

	tree := parser.Parse(content, nil)
	defer tree.Close()

	root := tree.RootNode()

	queryPattern := `
	(class_declaration
		name: (identifier) @class_name)
	(field_declaration
		type: (type_identifier)  @field_type
		declarator: (variable_declarator name: (identifier) @declarator_name)
		)
	(constructor_declaration
		name: (identifier) @constructor_name
		parameters: (
			formal_parameters (
				formal_parameter type: (type_identifier) @constructor_parameter_type
				name : (identifier) @constructor_parameter_name
				)
			)
		)


	(method_declaration
		name: (identifier) @method_name)
	`

	query, _ := tree_sitter.NewQuery(lang, (queryPattern))

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(query, root, content)

	var classData ClassJava

	for {
		match := matches.Next()
		if match == nil {
			break
		}
		var newField Field
		isFieldMatch := false

		var newConstructor Constructor
		isContructor := false

		for _, capture := range match.Captures {
			tag := query.CaptureNames()[capture.Index]
			text := capture.Node.Utf8Text(content)

			switch tag {
			case "class_name":
				classData.Name = text
			case "method_name":
				newMethod := Method{Name: text}
				classData.Methods = append(classData.Methods, newMethod)
			case "field_type":
				newField.TypeName = text
				isFieldMatch = true
			case "declarator_name":
				newField.Declarator = text
				isFieldMatch = true
			case "constructor_name":
				newConstructor.Name = text
				isContructor = true
			case "constructor_parameter_type":
				newField.TypeName = text
				isContructor = true
			case "constructor_parameter_name":
				newField.Declarator = text
				isContructor = true

			}

		}
		if isContructor {
			newConstructor.Parameters = append(newConstructor.Parameters, newField)

			classData.Constructors = append(classData.Constructors, newConstructor)
		}

		if isFieldMatch {
			classData.Fields = append(classData.Fields, newField)
		}
	}

	fmt.Printf("Estrutura preenchida com sucesso!\n\n")
	fmt.Printf("Classe: %s\n", classData.Name)

	fmt.Printf("Quantidade de propriedades: %d\n", len(classData.Fields))
	fmt.Printf("Quantidade de construtores: %d\n", len(classData.Constructors))

	fmt.Printf("Quantidade de métodos: %d\n", len(classData.Methods))
	for _, f := range classData.Fields {

		fmt.Printf(" - Propriedades: \n")
		fmt.Printf("  - Tipo: %s Propriedade: %s\n ", f.Declarator, f.TypeName)
	}
	for _, c := range classData.Constructors {
		fmt.Printf(" - Construtor: %s\n", c.Name)
		for _, p := range c.Parameters {
			fmt.Printf("    - Tipo: %s Parameter: %s\n", p.TypeName, p.Declarator)

		}
	}

	for _, m := range classData.Methods {
		fmt.Printf("  - Método: %s\n", m.Name)

	}

}

func printReadableTree(node *tree_sitter.Node, sourceCode []byte, depth int) {
	if node == nil {
		return

	}

	indent := strings.Repeat(" ", depth)

	nodeType := node.Kind()

	var textSample string
	if node.ChildCount() == 0 {
		textSample = fmt.Sprintf(" -> \"%s\"", node.Utf8Text(sourceCode))
	}

	fmt.Printf("%s[%s]%s\n", indent, nodeType, textSample)

	childCount := node.ChildCount()
	for i := uint(0); i < uint(childCount); i++ {
		child := node.Child(i)
		printReadableTree(child, sourceCode, depth+1)
	}

}
