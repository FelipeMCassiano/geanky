package java

import (
	"encoding/json"
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

type ClassJob struct {
	Data   ClassJava
	RelDir string
}

func AnalyzeDirectory(rootDir string, outputDir string) {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Erro ao criar diretório de saída: %v", err)
	}

	var jobs []ClassJob
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, maxGoroutines)
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

	for _, path := range filePaths {
		wg.Add(1)
		sem <- struct{}{}
		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()

			fmt.Printf("🔍 Analisando: %s\n", p)
			classData := Analyze(p)

			if classData.Name != "" {
				// Descobre a pasta original relativa ao rootDir
				relPath, err := filepath.Rel(rootDir, p)
				if err != nil {
					relPath = filepath.Base(p) // Fallback seguro
				}
				relDir := filepath.Dir(relPath)

				mu.Lock()
				jobs = append(jobs, ClassJob{Data: classData, RelDir: relDir})
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()

	var allClasses []ClassJava
	for _, job := range jobs {
		allClasses = append(allClasses, job.Data)
	}

	for _, job := range jobs {
		targetDir := filepath.Join(outputDir, job.RelDir)
		os.MkdirAll(targetDir, os.ModePerm)

		outFileName := fmt.Sprintf("%s.md", job.Data.Name)
		outFilePath := filepath.Join(targetDir, outFileName)

		GenerateMarkdown(job, jobs, outFilePath)
	}

	if len(allClasses) > 0 {
		globalOutPath := filepath.Join(outputDir, "00_Architecture_Overview.md")
		GenerateGlobalArchitecture(allClasses, globalOutPath)
	}

	fmt.Println("🚀 Varredura e Geração concluídas com sucesso!")

	data, err := json.Marshal(allClasses)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(filepath.Join(outputDir, "classes.json"), data, os.ModePerm)
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
