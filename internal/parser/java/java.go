package java

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	java_lang "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

const maxGoroutines = 8

func AnalyzeDirectory(rootDir string, outputDir string) {

	var allClasses []ClassJava
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, maxGoroutines)

	var filePaths []string
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".java") {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Erro ao varrer diretórios: %v", err)
	}

	for _, path := range filePaths {
		wg.Add(1)
		sem <- struct{}{}

		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()

			fmt.Printf("🔍 Analisando: %s\n", p)
			classData := Analyze(p)

			if classData.Name != "" {
				outFileName := fmt.Sprintf("%s.md", classData.Name)
				outFilePath := filepath.Join(outputDir, outFileName)
				GenerateMarkdown(classData, outFilePath)

				mu.Lock()
				allClasses = append(allClasses, classData)
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()

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

	return classData
}
