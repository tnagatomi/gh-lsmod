package parser

import (
	"os"
	"path/filepath"

	"github.com/tnagatomi/gh-go-mod-browser/model"
	"golang.org/x/mod/modfile"
)

// GoModParser parses go.mod files and extracts direct dependencies
type GoModParser struct {
	filePath string
}

// NewGoModParser creates a new GoModParser instance
func NewGoModParser(filePath string) *GoModParser {
	return &GoModParser{
		filePath: filePath,
	}
}

// NewParserForCurrentDirectory creates a new GoModParser for the go.mod file in the current directory
func NewParserForCurrentDirectory() (*GoModParser, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return NewGoModParser(filepath.Join(cwd, "go.mod")), nil
}

// Parse parses the go.mod file and returns a list of direct dependencies
func (p *GoModParser) Parse() ([]*model.Package, error) {
	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return nil, err
	}

	file, err := modfile.Parse(p.filePath, data, nil)
	if err != nil {
		return nil, err
	}

	var packages []*model.Package

	// Add direct requires
	for _, req := range file.Require {
		if !req.Indirect {
			pkg := model.NewPackage(req.Mod.Path, req.Mod.Version)
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}
