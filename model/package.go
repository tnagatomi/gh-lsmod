package model

import (
	"fmt"
	"strings"
)

// Package represents a Go module dependency
type Package struct {
	Path      string // Import path
	Version   string // Version
	IsGitHub  bool   // Whether it's a GitHub repository
	IsStarred bool   // Whether it's starred by the user
	Size      int64  // Size in bytes
}

// NewPackage creates a new Package instance
func NewPackage(path, version string) *Package {
	return &Package{
		Path:      path,
		Version:   version,
		IsGitHub:  strings.HasPrefix(path, "github.com/"),
		IsStarred: false,
		Size:      0,
	}
}

// GitHubRepoPath returns the GitHub repository path (owner/repo)
// Returns empty string if not a GitHub repository
func (p *Package) GitHubRepoPath() string {
	if !p.IsGitHub {
		return ""
	}

	// Remove github.com/ prefix
	repoPath := strings.TrimPrefix(p.Path, "github.com/")

	// Handle version suffix (e.g., github.com/cli/go-gh/v2)
	parts := strings.Split(repoPath, "/")
	if len(parts) >= 2 {
		// Check if the last part starts with 'v' followed by a number
		lastPart := parts[len(parts)-1]
		if len(lastPart) > 1 && lastPart[0] == 'v' && '0' <= lastPart[1] && lastPart[1] <= '9' {
			// Remove the version suffix
			return strings.Join(parts[:len(parts)-1], "/")
		}
	}

	// Handle case with no version suffix
	if len(parts) >= 2 {
		return strings.Join(parts[:2], "/")
	}

	return repoPath
}

// GitHubURL returns the GitHub repository URL
// Returns empty string if not a GitHub repository
func (p *Package) GitHubURL() string {
	if !p.IsGitHub {
		return ""
	}
	return fmt.Sprintf("https://github.com/%s", p.GitHubRepoPath())
}

// PkgGoDevURL returns the pkg.go.dev URL for the package
func (p *Package) PkgGoDevURL() string {
	return fmt.Sprintf("https://pkg.go.dev/%s", p.Path)
}

// StarSymbol returns the star symbol based on the starred status
// Returns empty string if not a GitHub repository
func (p *Package) StarSymbol() string {
	if !p.IsGitHub {
		return ""
	}
	if p.IsStarred {
		return "★"
	}
	return "☆"
}

// FormattedSize returns the size in a human-readable format
func (p *Package) FormattedSize() string {
	if p.Size == 0 {
		return "unknown"
	}

	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
	)

	var (
		value float64
		unit  string
	)

	switch {
	case p.Size >= int64(GB):
		value = float64(p.Size) / GB
		unit = "GB"
	case p.Size >= int64(MB):
		value = float64(p.Size) / MB
		unit = "MB"
	case p.Size >= int64(KB):
		value = float64(p.Size) / KB
		unit = "KB"
	default:
		value = float64(p.Size)
		unit = "B"
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

// String returns a string representation of the package
func (p *Package) String() string {
	symbol := p.StarSymbol()
	if symbol == "" {
		return fmt.Sprintf("  %s", p.Path) // Add space for alignment
	}
	return fmt.Sprintf("%s %s", symbol, p.Path)
}
