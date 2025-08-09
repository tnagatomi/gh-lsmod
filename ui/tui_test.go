package ui

import (
	"testing"

	"github.com/tnagatomi/gh-lsmod/model"
)

// MockGitHubClient is a mock implementation of the github.GitHubClient interface
type MockGitHubClient struct {
	starredRepos     map[string]bool
	starCallCount    int
	unstarCallCount  int
	starAllCallCount int
}

// NewMockGitHubClient creates a new mock GitHub client
func NewMockGitHubClient() *MockGitHubClient {
	return &MockGitHubClient{
		starredRepos: make(map[string]bool),
	}
}

// CheckStarredStatus mocks checking if repositories are starred
func (m *MockGitHubClient) CheckStarredStatus(packages []*model.Package) error {
	for _, pkg := range packages {
		if pkg.IsGitHub {
			pkg.IsStarred = m.starredRepos[pkg.Path]
		}
	}
	return nil
}

// StarRepository mocks starring a repository
func (m *MockGitHubClient) StarRepository(pkg *model.Package) error {
	m.starCallCount++
	m.starredRepos[pkg.Path] = true
	pkg.IsStarred = true
	return nil
}

// UnstarRepository mocks unstarring a repository
func (m *MockGitHubClient) UnstarRepository(pkg *model.Package) error {
	m.unstarCallCount++
	delete(m.starredRepos, pkg.Path)
	pkg.IsStarred = false
	return nil
}

// StarAllUnstarred mocks starring all unstarred repositories
func (m *MockGitHubClient) StarAllUnstarred(packages []*model.Package) (int, error) {
	m.starAllCallCount++
	count := 0
	for _, pkg := range packages {
		if pkg.IsGitHub && !pkg.IsStarred {
			m.starredRepos[pkg.Path] = true
			pkg.IsStarred = true
			count++
		}
	}
	return count, nil
}

func TestAppView(t *testing.T) {
	// Create test packages
	pkg := model.NewPackage("github.com/charmbracelet/bubbles", "v0.20.0")
	packages := []*model.Package{pkg}

	// Create mock GitHub client
	mockClient := NewMockGitHubClient()

	// Create app
	app := NewApp(packages, mockClient)

	// Test list view
	app.state = StateList
	listView := app.View()
	if listView == "" {
		t.Errorf("Expected list view to be non-empty")
	}

	// Test dialog view
	app.state = StateDialog
	app.dialog = NewDialog("Test Title", "Test Message")
	dialogView := app.View()
	if dialogView == "" {
		t.Errorf("Expected dialog view to be non-empty")
	}
}
