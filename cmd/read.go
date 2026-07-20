/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/FelipeMCassiano/geanky/internal/parser/java"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: read,
}

func init() {

	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	readCmd.Flags().StringP("config", "o", "", "sets the config file for geanky")
}

func read(cmd *cobra.Command, args []string) {
	targetDir := args[0]
	// outputFlag, err := cmd.Flags().GetString("output")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// wd, err := os.Getwd()

	outDir := "./docs"
	// if outputFlag != "" {
	// 	outDir = filepath.Join(wd, outputFlag)
	// }

	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	java.AnalyzeDirectory(targetDir, outDir)
}
