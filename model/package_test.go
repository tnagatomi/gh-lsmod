package model

import (
	"testing"
)

func TestGitHubRepoPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Non-GitHub repository",
			path:     "golang.org/x/mod",
			expected: "",
		},
		{
			name:     "GitHub repository with version suffix",
			path:     "github.com/cli/go-gh/v2",
			expected: "cli/go-gh",
		},
		{
			name:     "GitHub repository without version suffix",
			path:     "github.com/charmbracelet/bubbles",
			expected: "charmbracelet/bubbles",
		},
		{
			name:     "GitHub repository with deep path",
			path:     "github.com/charmbracelet/bubbles/list/item",
			expected: "charmbracelet/bubbles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := NewPackage(tt.path, "v1.0.0")
			if got := pkg.GitHubRepoPath(); got != tt.expected {
				t.Errorf("GitHubRepoPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGitHubURL(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Non-GitHub repository",
			path:     "golang.org/x/mod",
			expected: "",
		},
		{
			name:     "GitHub repository with version suffix",
			path:     "github.com/cli/go-gh/v2",
			expected: "https://github.com/cli/go-gh",
		},
		{
			name:     "GitHub repository without version suffix",
			path:     "github.com/charmbracelet/bubbles",
			expected: "https://github.com/charmbracelet/bubbles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := NewPackage(tt.path, "v1.0.0")
			if got := pkg.GitHubURL(); got != tt.expected {
				t.Errorf("GitHubURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPkgGoDevURL(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Non-GitHub repository",
			path:     "golang.org/x/mod",
			expected: "https://pkg.go.dev/golang.org/x/mod",
		},
		{
			name:     "GitHub repository with version suffix",
			path:     "github.com/cli/go-gh/v2",
			expected: "https://pkg.go.dev/github.com/cli/go-gh/v2",
		},
		{
			name:     "GitHub repository without version suffix",
			path:     "github.com/charmbracelet/bubbles",
			expected: "https://pkg.go.dev/github.com/charmbracelet/bubbles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := NewPackage(tt.path, "v1.0.0")
			if got := pkg.PkgGoDevURL(); got != tt.expected {
				t.Errorf("PkgGoDevURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStarSymbol(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		isStarred bool
		expected  string
	}{
		{
			name:      "Non-GitHub repository",
			path:      "golang.org/x/mod",
			isStarred: false,
			expected:  "",
		},
		{
			name:      "GitHub repository starred",
			path:      "github.com/charmbracelet/bubbles",
			isStarred: true,
			expected:  "★",
		},
		{
			name:      "GitHub repository not starred",
			path:      "github.com/charmbracelet/bubbles",
			isStarred: false,
			expected:  "☆",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := NewPackage(tt.path, "v1.0.0")
			pkg.IsStarred = tt.isStarred
			if got := pkg.StarSymbol(); got != tt.expected {
				t.Errorf("StarSymbol() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormattedSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "Zero size",
			size:     0,
			expected: "unknown",
		},
		{
			name:     "Bytes",
			size:     500,
			expected: "500.00 B",
		},
		{
			name:     "Kilobytes",
			size:     1500,
			expected: "1.46 KB",
		},
		{
			name:     "Megabytes",
			size:     1500000,
			expected: "1.43 MB",
		},
		{
			name:     "Gigabytes",
			size:     1500000000,
			expected: "1.40 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := NewPackage("test/package", "v1.0.0")
			pkg.Size = tt.size
			if got := pkg.FormattedSize(); got != tt.expected {
				t.Errorf("FormattedSize() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		isStarred bool
		expected  string
	}{
		{
			name:      "Non-GitHub repository",
			path:      "golang.org/x/mod",
			isStarred: false,
			expected:  "  golang.org/x/mod",
		},
		{
			name:      "GitHub repository starred",
			path:      "github.com/charmbracelet/bubbles",
			isStarred: true,
			expected:  "★ github.com/charmbracelet/bubbles",
		},
		{
			name:      "GitHub repository not starred",
			path:      "github.com/charmbracelet/bubbles",
			isStarred: false,
			expected:  "☆ github.com/charmbracelet/bubbles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := NewPackage(tt.path, "v1.0.0")
			pkg.IsStarred = tt.isStarred
			if got := pkg.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
