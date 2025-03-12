package main

import (
	"fmt"
	"os"

	"github.com/tnagatomi/gh-go-mod-browser/github"
	"github.com/tnagatomi/gh-go-mod-browser/parser"
	"github.com/tnagatomi/gh-go-mod-browser/ui"
)

func main() {
	// Create a parser for the go.mod file in the current directory
	gomodParser, err := parser.NewParserForCurrentDirectory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Extract direct dependencies
	packages, err := gomodParser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(packages) == 0 {
		fmt.Println("No direct dependencies found in go.mod file.")
		os.Exit(0)
	}

	// Initialize GitHub client
	githubClient, err := github.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check starred status
	err = githubClient.CheckStarredStatus(packages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Run TUI application
	err = ui.Run(packages, githubClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
