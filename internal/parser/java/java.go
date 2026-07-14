package java

import (
	"fmt"
	"log"
	"os"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	java_lang "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

func Analyze(filePath string) ClassJava {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	parser := tree_sitter.NewParser()
	defer parser.Close()
	lang := tree_sitter.NewLanguage(java_lang.Language())
	parser.SetLanguage(lang)

	tree := parser.Parse(content, nil)
	defer tree.Close()

	root := tree.RootNode()

	queryPattern := JavaQueries()

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
			if handler, ok := handlers[tag]; ok {
				handler(&capture.Node, content, &classData)
			}
		}

	}
	GenerateMarkdown(classData, "doccccc.md")

	return classData

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
