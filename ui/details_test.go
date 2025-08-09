package ui

import (
	"strings"
	"testing"

	"github.com/tnagatomi/gh-lsmod/model"
)

func TestPackageDetailsSetSize(t *testing.T) {
	details := NewPackageDetails()

	details.SetSize(100, 20)

	if details.width != 100 {
		t.Errorf("Expected width to be 100, got %d", details.width)
	}

	if details.height != 20 {
		t.Errorf("Expected height to be 20, got %d", details.height)
	}
}

func TestPackageDetailsView(t *testing.T) {
	tests := []struct {
		name     string
		pkg      *model.Package
		contains []string
	}{
		{
			name: "No package selected",
			pkg:  nil,
			contains: []string{
				"No package selected",
			},
		},
		{
			name: "GitHub package",
			pkg:  func() *model.Package {
				pkg := model.NewPackage("github.com/charmbracelet/bubbles", " v0.20.0")
				pkg.Size = 1024 * 1024 // 1MB
				return pkg
			}(),
			contains: []string{
				"Name: github.com/charmbracelet/bubbles",
				"Version:  v0.20.0",
				"Size: 1.00 MB",
				"GitHub: https://github.com/charmbracelet/bubbles",
				"pkg.go.dev: https://pkg.go.dev/github.com/charmbracelet/bubbles",
			},
		},
		{
			name: "GitHub package with unknown size",
			pkg:  model.NewPackage("golang.org/x/mod", "v1.0.0"),
			contains: []string{
				"Name: golang.org/x/mod",
				"Version: v1.0.0",
				"Size: unknown",
				"pkg.go.dev: https://pkg.go.dev/golang.org/x/mod",
			},
		},
		{
			name: "Non-GitHub package",
			pkg:  func() *model.Package {
				pkg := model.NewPackage("golang.org/x/mod", "v1.0.0")
				pkg.Size = 1024 // 1KB
				return pkg
			}(),
			contains: []string{
				"Name: golang.org/x/mod",
				"Version: v1.0.0",
				"Size: 1.00 KB",
				"pkg.go.dev: https://pkg.go.dev/golang.org/x/mod",
			},
		},
}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			details := NewPackageDetails()
			details.SetPackage(tt.pkg)

			result := details.View()

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected view to contain %q, but it didn't.\nGot: %s", expected, result)
				}
			}

			// For GitHub package, ensure it contains GitHub URL
			if tt.pkg != nil && tt.pkg.IsGitHub {
				if !strings.Contains(result, "GitHub:") {
					t.Errorf("Expected view to contain GitHub URL for GitHub package, but it didn't.\nGot: %s", result)
				}
			}

			// For non-GitHub package, ensure it doesn't contain GitHub URL
			if tt.pkg != nil && !tt.pkg.IsGitHub {
				if strings.Contains(result, "GitHub:") {
					t.Errorf("Expected view not to contain GitHub URL for non-GitHub package, but it did.\nGot: %s", result)
				}
			}
		})
	}
}
