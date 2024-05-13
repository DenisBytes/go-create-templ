/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-templ-create",
	Short: "create-next-app@latest but for go",
	Long: `
	The opinionated version of "npx create-next-app", 
	where you will have HTMX, Templ, Air, Supabase, Bun (ORM), Bunjs, and a framework of choice. 
	The login boilerplate code will be setup for you.
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
