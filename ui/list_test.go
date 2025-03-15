package ui

import (
	"testing"

	"github.com/tnagatomi/gh-go-mod-browser/model"
)

func TestPackageItemFilterValue(t *testing.T) {
	pkg := model.NewPackage("github.com/charmbracelet/bubbles", "v0.20.0")
	item := PackageItem{pkg: pkg}

	if got := item.FilterValue(); got != pkg.Path {
		t.Errorf("FilterValue() = %v, want %v", got, pkg.Path)
	}
}

func TestPackageItemTitle(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		version  string
		isStarred bool
		expected string
	}{
		{
			name:     "GitHub package starred",
			path:     "github.com/charmbracelet/bubbles",
			version:  "v0.20.0",
			isStarred: true,
			expected: "★ github.com/charmbracelet/bubbles",
		},
		{
			name:     "GitHub package not starred",
			path:     "github.com/charmbracelet/bubbles",
			version:  "v0.20.0",
			isStarred: false,
			expected: "☆ github.com/charmbracelet/bubbles",
		},
		{
			name:     "Non-GitHub package",
			path:     "golang.org/x/mod",
			version:  "v0.8.0",
			isStarred: false,
			expected: "  golang.org/x/mod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := model.NewPackage(tt.path, tt.version)
			pkg.IsStarred = tt.isStarred
			item := PackageItem{pkg: pkg}

			if got := item.Title(); got != tt.expected {
				t.Errorf("Title() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPackageItemDescription(t *testing.T) {
	tests := []struct {
		name     string
		pkg      *model.Package
		expected string
	}{
		{
			name: "GitHub package",
			pkg: func() *model.Package {
				pkg := model.NewPackage("github.com/charmbracelet/bubbles", "v1.0.0")
				pkg.Size = 1024 * 1024 // 1MB
				return pkg
			}(),
			expected: "[pkg.go] [GitHub] (1.00 MB)",
		},
		{
			name: "GitHub package with unknown size",
			pkg:  model.NewPackage("github.com/charmbracelet/bubbles", "v1.0.0"),
			expected: "[pkg.go] [GitHub] (unknown)",
		},
		{
			name: "Non-GitHub package",
			pkg: func() *model.Package {
				pkg := model.NewPackage("golang.org/x/mod", "v1.0.0")
				pkg.Size = 1024 // 1KB
				return pkg
			}(),
			expected: "[pkg.go] (1.00 KB)",
		},
		{
			name: "Non-GitHub package with unknown size",
			pkg:  model.NewPackage("golang.org/x/mod", "v1.0.0"),
			expected: "[pkg.go] (unknown)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := PackageItem{pkg: tt.pkg}

			if got := item.Description(); got != tt.expected {
				t.Errorf("Description() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSelectedPackage(t *testing.T) {
	// Create test packages
	pkg1 := model.NewPackage("github.com/charmbracelet/bubbles", "v0.20.0")
	pkg2 := model.NewPackage("golang.org/x/mod", "v0.8.0")
	packages := []*model.Package{pkg1, pkg2}

	// Create package list
	list := NewPackageList(packages)

	// Test selected package (default is first item)
	if got := list.SelectedPackage(); got != pkg1 {
		t.Errorf("SelectedPackage() = %v, want %v", got, pkg1)
	}

	// Manually set the index to the second item
	list.list.Select(1)

	// Test selected package after selection change
	if got := list.SelectedPackage(); got != pkg2 {
		t.Errorf("SelectedPackage() = %v, want %v", got, pkg2)
	}

	// Test with invalid index
	list.list.Select(-1)
	if got := list.SelectedPackage(); got != nil {
		t.Errorf("SelectedPackage() with invalid index = %v, want nil", got)
	}

	// Test with out of bounds index
	list.list.Select(len(packages))
	if got := list.SelectedPackage(); got != nil {
		t.Errorf("SelectedPackage() with out of bounds index = %v, want nil", got)
	}
}
