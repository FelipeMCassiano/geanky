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
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Erro ao criar diretório de saída: %v", err)
	}

	var allClasses []ClassJava
	var mu sync.Mutex // Mutex para proteger o append no allClasses concorrente
	var wg sync.WaitGroup

	// Canal limitador (Semaphore) para controlar o máximo de concorrência
	sem := make(chan struct{}, maxGoroutines)

	// Coleta todos os caminhos de arquivos .java primeiro
	var filePaths []string
	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
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

	// Dispara a análise concorrente de cada arquivo
	for _, path := range filePaths {
		wg.Add(1)
		sem <- struct{}{} // Bloqueia se atingir o limite de maxGoroutines

		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }() // Libera espaço no canal ao terminar

			fmt.Printf("🔍 Analisando: %s\n", p)
			classData := Analyze(p)

			if classData.Name != "" {
				// Escreve o Markdown individual de forma assíncrona
				outFileName := fmt.Sprintf("%s.md", classData.Name)
				outFilePath := filepath.Join(outputDir, outFileName)
				GenerateMarkdown(classData, outFilePath)

				// Lock para evitar Race Condition ao salvar os dados no array global
				mu.Lock()
				allClasses = append(allClasses, classData)
				mu.Unlock()
			}
		}(path)
	}

	// Aguarda todas as goroutines finalizarem
	wg.Wait()

	if len(allClasses) > 0 {
		globalOutPath := filepath.Join(outputDir, "00_Architecture_Overview.md")
		GenerateGlobalArchitecture(allClasses, globalOutPath)
	}

	fmt.Println("🚀 Varredura concorrente concluída com sucesso!")
}

func Analyze(filePath string) ClassJava {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// IMPORTANTE: Cada Goroutine precisa criar seu próprio Parser e instanciar sua própria Language.
	// O Tree-Sitter NÃO é thread-safe se compartilhado no mesmo objeto!
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
