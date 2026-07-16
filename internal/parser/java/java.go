package java

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	java_lang "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

func AnalyzeDirectory(rootDir string, outputDir string) {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Erro ao criar diretório de saída: %v", err)
	}

	var allClasses []ClassJava

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".java") {
			fmt.Printf("🔍 Analisando: %s\n", path)
			classData := Analyze(path)

			if classData.Name != "" {
				allClasses = append(allClasses, classData) // Salva no array global

				outFileName := fmt.Sprintf("%s.md", classData.Name)
				outFilePath := filepath.Join(outputDir, outFileName)
				GenerateMarkdown(classData, outFilePath)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Erro ao varrer diretórios: %v", err)
	}

	// <-- GERA O MAPA GLOBAL AO FINAL DA VARREDURA
	if len(allClasses) > 0 {
		globalOutPath := filepath.Join(outputDir, "00_Architecture_Overview.md")
		GenerateGlobalArchitecture(allClasses, globalOutPath)
	}

	fmt.Println("🚀 Varredura concluída com sucesso!")
}
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

	// Removido o GenerateMarkdown daqui!
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
