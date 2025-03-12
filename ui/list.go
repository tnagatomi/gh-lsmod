package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tnagatomi/gh-go-mod-browser/model"
)

// PackageItem represents a package in the list
type PackageItem struct {
	pkg *model.Package
}

// FilterValue returns the value to filter on
func (i PackageItem) FilterValue() string {
	return i.pkg.Path
}

// Title returns the title of the item
func (i PackageItem) Title() string {
	return i.pkg.String()
}

// Description returns the description of the item
func (i PackageItem) Description() string {
	desc := ""
	if i.pkg.IsGitHub {
		desc += "[GitHub] "
	}
	desc += "[pkg.go]"
	return desc
}

// PackageList represents the list of packages
type PackageList struct {
	list     list.Model
	packages []*model.Package
	keyMap   PackageListKeyMap
	help     help.Model
	width    int
	height   int
}

// PackageListKeyMap defines the key bindings for the package list
type PackageListKeyMap struct {
	OpenGitHub   key.Binding
	OpenPkgGoDev key.Binding
	ToggleStar   key.Binding
	StarAll      key.Binding
	Quit         key.Binding
}

// DefaultPackageListKeyMap returns the default key bindings for the package list
func DefaultPackageListKeyMap() PackageListKeyMap {
	return PackageListKeyMap{
		OpenGitHub: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "open GitHub"),
		),
		OpenPkgGoDev: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "open pkg.go.dev"),
		),
		ToggleStar: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "star/unstar"),
		),
		StarAll: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "star all"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k PackageListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.OpenGitHub, k.OpenPkgGoDev, k.ToggleStar, k.StarAll, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k PackageListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.OpenGitHub, k.OpenPkgGoDev},
		{k.ToggleStar, k.StarAll},
		{k.Quit},
	}
}

// NewPackageList creates a new package list
func NewPackageList(packages []*model.Package) *PackageList {
	// Create list items
	items := make([]list.Item, len(packages))
	for i, pkg := range packages {
		items[i] = PackageItem{pkg: pkg}
	}

	// Create list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Go Module Browser"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Bold(true).
		MarginLeft(2)

	keyMap := DefaultPackageListKeyMap()
	helpModel := help.New()

	return &PackageList{
		list:     l,
		packages: packages,
		keyMap:   keyMap,
		help:     helpModel,
	}
}

// Init initializes the package list
func (l *PackageList) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the package list
func (l *PackageList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = msg.Width
		l.height = msg.Height
		l.help.Width = msg.Width
		l.SetSize(msg.Width, msg.Height-4) // Reserve space for help
	}

	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

// View renders the package list
func (l *PackageList) View() string {
	// Render the list view
	listView := l.list.View()

	// Render the help view
	helpView := l.help.View(l.keyMap)

	// Calculate appropriate spacing
	height := 2 // Default spacing

	return listView + strings.Repeat("\n", height) + helpView
}

// SelectedPackage returns the currently selected package
func (l *PackageList) SelectedPackage() *model.Package {
	idx := l.list.Index()
	if idx < 0 || idx >= len(l.packages) {
		return nil
	}
	return l.packages[idx]
}

// SetSize sets the size of the list
func (l *PackageList) SetSize(width, height int) {
	l.list.SetSize(width, height)
}
