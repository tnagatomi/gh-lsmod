package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/browser"
	"github.com/tnagatomi/gh-go-mod-browser/github"
	"github.com/tnagatomi/gh-go-mod-browser/model"
)

// State represents the current state of the TUI
type State int

const (
	StateList State = iota
	StateDialog
)

// Layout constants
const (
	DetailViewHeight = 10
	HelpViewHeight   = 3
	MinListHeight    = 5
)

// App represents the TUI application
type App struct {
	packages     []*model.Package
	list         *PackageList
	details      *PackageDetails
	state        State
	githubClient *github.Client
	dialog       *Dialog
	err          error
	width        int
	height       int
}

// NewApp creates a new TUI application
func NewApp(packages []*model.Package, githubClient *github.Client) *App {
	list := NewPackageList(packages)
	details := NewPackageDetails()

	if len(packages) > 0 {
		details.SetPackage(packages[0])
	}

	return &App{
		packages:     packages,
		list:         list,
		details:      details,
		state:        StateList,
		githubClient: githubClient,
		width:        80,
		height:       24,
	}
}

// Init initializes the TUI application
func (a *App) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.updateComponentSizes()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.list.keyMap.Quit):
			return a, tea.Quit
		}
	}

	switch a.state {
	case StateList:
		return a.updateList(msg)
	case StateDialog:
		return a.updateDialog(msg)
	}

	return a, cmd
}

// updateComponentSizes updates the sizes of all components
func (a *App) updateComponentSizes() {
	// Reserve space for details and help message
	listHeight := a.height - DetailViewHeight - HelpViewHeight
	if listHeight < MinListHeight {
		listHeight = MinListHeight // Ensure minimum list height
	}
	a.list.SetSize(a.width, listHeight)
	a.details.SetSize(a.width, DetailViewHeight)
}

// updateList handles user input in the list view
func (a *App) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.list.keyMap.OpenGitHub):
			// Open GitHub repository in browser
			pkg := a.list.SelectedPackage()
			if pkg != nil && pkg.IsGitHub {
				url := pkg.GitHubURL()
				if url != "" {
					_ = browser.OpenURL(url)
				}
			}

		case key.Matches(msg, a.list.keyMap.OpenPkgGoDev):
			// Open pkg.go.dev page in browser
			pkg := a.list.SelectedPackage()
			if pkg != nil {
				url := pkg.PkgGoDevURL()
				_ = browser.OpenURL(url)
			}

		case key.Matches(msg, a.list.keyMap.ToggleStar):
			// Toggle star status
			pkg := a.list.SelectedPackage()
			if pkg != nil && pkg.IsGitHub {
				if pkg.IsStarred {
					_ = a.githubClient.UnstarRepository(pkg)
				} else {
					_ = a.githubClient.StarRepository(pkg)
				}
			}

		case key.Matches(msg, a.list.keyMap.StarAll):
			// Show confirmation dialog for starring all unstarred repositories
			unstarredCount := 0
			for _, pkg := range a.packages {
				if pkg.IsGitHub && !pkg.IsStarred {
					unstarredCount++
				}
			}
			if unstarredCount > 0 {
				a.dialog = NewDialog(
					"Star all unstarred GitHub repositories?",
					"This will add stars to "+string(rune(unstarredCount))+" repositories.",
				)
				a.state = StateDialog
				return a, nil
			}
		}
	}

	// Update the list component
	_, cmd = a.list.Update(msg)

	// Update the details component with the selected package
	pkg := a.list.SelectedPackage()
	if pkg != nil {
		a.details.SetPackage(pkg)
	}

	return a, cmd
}

// updateDialog handles user input in the dialog view
func (a *App) updateDialog(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.dialog.keyMap.Confirm):
			// Confirm dialog
			_, _ = a.githubClient.StarAllUnstarred(a.packages)
			a.state = StateList
			a.dialog = nil

		case key.Matches(msg, a.dialog.keyMap.Cancel):
			// Cancel dialog
			a.state = StateList
			a.dialog = nil
		}
	}
	return a, nil
}

// View renders the TUI
func (a *App) View() string {
	switch a.state {
	case StateList:
		return a.list.View() + "\n" + a.details.View()
	case StateDialog:
		return a.dialog.View()
	}
	return ""
}

// Run runs the TUI application
func Run(packages []*model.Package, githubClient *github.Client) error {
	app := NewApp(packages, githubClient)
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
