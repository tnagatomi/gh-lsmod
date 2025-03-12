package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tnagatomi/gh-go-mod-browser/model"
)

// PackageDetails represents the details view of a package
type PackageDetails struct {
	pkg    *model.Package
	width  int
	height int
	styles DetailsStyles
}

// DetailsStyles contains the styles for the details view
type DetailsStyles struct {
	Title       lipgloss.Style
	Label       lipgloss.Style
	Value       lipgloss.Style
	Border      lipgloss.Style
	EmptyBorder lipgloss.Style
}

// DefaultDetailsStyles returns the default styles for the details view
func DefaultDetailsStyles() DetailsStyles {
	return DetailsStyles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			MarginLeft(2),
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true),
		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		Border: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2),
		EmptyBorder: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Foreground(lipgloss.Color("240")).
			Italic(true),
	}
}

// NewPackageDetails creates a new package details view
func NewPackageDetails() *PackageDetails {
	return &PackageDetails{
		pkg:    nil,
		width:  80,
		height: 10,
		styles: DefaultDetailsStyles(),
	}
}

// SetPackage sets the package to display
func (d *PackageDetails) SetPackage(pkg *model.Package) {
	d.pkg = pkg
}

// SetSize sets the size of the details view
func (d *PackageDetails) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// Init initializes the details view
func (d *PackageDetails) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the details view
func (d *PackageDetails) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return d, nil
}

// View renders the details view
func (d *PackageDetails) View() string {
	if d.pkg == nil {
		return d.styles.EmptyBorder.Width(d.width - 4).Render("No package selected")
	}

	// Create the details content
	var content string

	// Add the package name
	content += d.styles.Label.Render("Name: ") + d.styles.Value.Render(d.pkg.Path) + "\n\n"

	// Add the package version
	content += d.styles.Label.Render("Version: ") + d.styles.Value.Render(d.pkg.Version) + "\n\n"

	// Add the GitHub URL if it's a GitHub repository
	if d.pkg.IsGitHub {
		content += d.styles.Label.Render("GitHub: ") + d.styles.Value.Render(d.pkg.GitHubURL()) + "\n\n"
	}

	// Add the pkg.go.dev URL
	content += d.styles.Label.Render("pkg.go.dev: ") + d.styles.Value.Render(d.pkg.PkgGoDevURL())

	// Apply border to the content
	return d.styles.Border.Width(d.width - 4).Render(content)
}