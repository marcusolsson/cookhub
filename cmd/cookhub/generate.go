package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcusolsson/cookhub/generator"
	"github.com/spf13/cobra"
)

var (
	inputDir  string
	outputDir string
	title     string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a static site from CookLang recipes",
	Long:  `Generate a static HTML website from a folder of CookLang (.cook) recipe files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputDir == "" {
			return fmt.Errorf("input directory is required")
		}

		// Validate input directory exists
		info, err := os.Stat(inputDir)
		if err != nil {
			return fmt.Errorf("input directory error: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("input path is not a directory: %s", inputDir)
		}

		// Default title to input directory name
		if title == "" {
			title = filepath.Base(inputDir)
		}

		cfg := generator.Config{
			InputDir:  inputDir,
			OutputDir: outputDir,
			Title:     title,
		}

		fmt.Printf("Generating site from %s to %s...\n", inputDir, outputDir)

		if err := generator.Generate(cfg); err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}

		fmt.Println("Site generated successfully!")
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&inputDir, "input", "i", "", "Input directory containing .cook files (required)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "public", "Output directory for generated HTML")
	generateCmd.Flags().StringVarP(&title, "title", "t", "", "Site title (defaults to input directory name)")
	generateCmd.MarkFlagRequired("input")
}
