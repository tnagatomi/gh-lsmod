package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tnagatomi/gh-go-mod-browser/model"
)

func TestParse(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "gomod-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create a test go.mod file
	goModContent := `module github.com/tnagatomi/gh-go-mod-browser

go 1.24.1

require (
	github.com/charmbracelet/bubbles v0.20.0
	github.com/charmbracelet/bubbletea v1.3.4
	golang.org/x/mod v0.24.0
)

require (
	github.com/charmbracelet/x/ansi v0.8.0 // indirect
)
`
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to write test go.mod file: %v", err)
	}

	// Test successful parsing
	parser := NewGoModParser(goModPath)
	packages, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse() returned an error: %v", err)
	}

	// Verify the parsed packages
	expectedPackages := []*model.Package{
		model.NewPackage("github.com/charmbracelet/bubbles", "v0.20.0"),
		model.NewPackage("github.com/charmbracelet/bubbletea", "v1.3.4"),
		model.NewPackage("golang.org/x/mod", "v0.24.0"),
	}

	if len(packages) != len(expectedPackages) {
		t.Errorf("Expected %d packages, got %d", len(expectedPackages), len(packages))
	}

	for i, pkg := range packages {
		if i >= len(expectedPackages) {
			break
		}
		expected := expectedPackages[i]
		if pkg.Path != expected.Path {
			t.Errorf("Package %d: expected path %s, got %s", i, expected.Path, pkg.Path)
		}
		if pkg.Version != expected.Version {
			t.Errorf("Package %d: expected version %s, got %s", i, expected.Version, pkg.Version)
		}
	}

	// Test parsing a non-existent file
	nonExistentParser := NewGoModParser(filepath.Join(tempDir, "non-existent.mod"))
	_, err = nonExistentParser.Parse()
	if err == nil {
		t.Error("Expected an error when parsing a non-existent file, got nil")
	}

	// Test parsing an invalid go.mod file
	invalidGoModPath := filepath.Join(tempDir, "invalid.mod")
	if err := os.WriteFile(invalidGoModPath, []byte("invalid content"), 0644); err != nil {
		t.Fatalf("Failed to write invalid go.mod file: %v", err)
	}

	invalidParser := NewGoModParser(invalidGoModPath)
	_, err = invalidParser.Parse()
	if err == nil {
		t.Error("Expected an error when parsing an invalid go.mod file, got nil")
	}
}
