/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	java "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

// sexpCmd represents the sexp command
var sexpCmd = &cobra.Command{
	Use:   "sexp <arquivo.java>",
	Short: "Imprime a Árvore Sintática (AST) no formato S-expression indentado",
	Long: `Este comando lê um arquivo fonte em Java, faz o parse utilizando o Tree-sitter 
e imprime toda a estrutura da árvore em formato S-expression (LISP-like),
adicionando quebras de linha e indentação para facilitar a leitura humana e 
ajudar na criação de queries.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Valida se o usuário passou o caminho do arquivo
		if len(args) == 0 {
			log.Fatal("Por favor, forneça o caminho do arquivo Java. Exemplo: go run main.go sexp ./Arquivo.java")
		}

		// Lê o conteúdo do arquivo
		content, err := os.ReadFile(args[0])
		if err != nil {
			log.Fatalf("Erro ao ler o arquivo: %v", err)
		}

		// Configura e inicializa o Parser do Tree-sitter para Java
		parser := tree_sitter.NewParser()
		defer parser.Close()

		lang := tree_sitter.NewLanguage(java.Language())
		parser.SetLanguage(lang)

		// Faz o parse do conteúdo e extrai a árvore
		tree := parser.Parse(content, nil)
		defer tree.Close()

		root := tree.RootNode()

		// Chama a nossa função de impressão formatada
		fmt.Println("--- S-Expression Indentada ---")
		printSexpIndented(root, "", 0)
	},
}

func init() {
	rootCmd.AddCommand(sexpCmd)
}

// printSexpIndented imprime a AST no formato S-expression com quebras de linha e indentação
func printSexpIndented(node *tree_sitter.Node, fieldName string, depth int) {
	if node == nil {
		return
	}

	// Cria o recuo com base na profundidade
	indent := strings.Repeat("  ", depth)

	// Prepara o prefixo do campo, se existir (ex: "name: ", "body: ")
	fieldPrefix := ""
	if fieldName != "" {
		fieldPrefix = fieldName + ": "
	}

	namedChildCount := 0
	childCount := node.ChildCount()
	for i := uint(0); i < childCount; i++ {
		if node.Child(i).IsNamed() {
			namedChildCount++
		}
	}

	// Se for uma "folha" (sem filhos nomeados), imprime tudo na mesma linha
	if namedChildCount == 0 {
		fmt.Printf("%s%s(%s)\n", indent, fieldPrefix, node.Kind())
		return
	}

	// Se tiver filhos, abre o parêntese e quebra a linha
	fmt.Printf("%s%s(%s\n", indent, fieldPrefix, node.Kind())

	// Chama a recursão para os filhos
	for i := uint(0); i < childCount; i++ {
		child := node.Child(i)
		// Ignora pontuações (chaves, ponto-e-vírgula) para manter a sintaxe de S-expression pura
		if child.IsNamed() {
			childFieldName := node.FieldNameForChild(uint32(i))
			printSexpIndented(child, childFieldName, depth+1)
		}
	}

	// Fecha o parêntese na mesma linha de recuo de quando abriu
	fmt.Printf("%s)\n", indent)
}
