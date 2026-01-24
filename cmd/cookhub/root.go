package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cookhub",
	Short: "A static site generator for CookLang recipes",
	Long:  `CookHub generates a static HTML website from a folder of CookLang (.cook) recipe files.`,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
