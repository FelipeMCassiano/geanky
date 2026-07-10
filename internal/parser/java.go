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
	Name string `json:"name"`
}

type ClassJava struct {
	Name    string   `json:"name"`
	Fields  []Field  `json:"fields"`
	Methods []Method `json:"methods"`
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
		type: (type_identifier)  @field_name)

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

		for _, capture := range match.Captures {
			tag := query.CaptureNames()[capture.Index]
			texto := capture.Node.Utf8Text(content)

			switch tag {
			case "class_name":
				classData.Name = texto
			case "method_name":
				newMethod := Method{Name: texto}
				classData.Methods = append(classData.Methods, newMethod)
			case "field_name":
				newField := Field{Name: texto}
				classData.Fields = append(classData.Fields, newField)

			}
		}
	}

	fmt.Printf("Estrutura preenchida com sucesso!\n\n")
	fmt.Printf("Classe: %s\n", classData.Name)
	fmt.Printf("Quantidade de métodos: %d\n", len(classData.Methods))
	for _, f := range classData.Fields {
		fmt.Printf(" - Propriedade: %s\n", f.Name)
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
