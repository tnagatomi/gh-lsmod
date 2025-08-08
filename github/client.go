package github

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/tnagatomi/gh-go-mod-browser/model"
)

// GitHubClient defines the interface for GitHub operations
type GitHubClient interface {
	CheckStarredStatus(packages []*model.Package) error
	StarRepository(pkg *model.Package) error
	UnstarRepository(pkg *model.Package) error
	StarAllUnstarred(packages []*model.Package) (int, error)
}

// Client handles GitHub API operations
type Client struct {
	restClient *api.RESTClient
}

// NewClient creates a new GitHub client
func NewClient() (*Client, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub REST client: %w", err)
	}

	return &Client{
		restClient: restClient,
	}, nil
}

// CheckStarredStatus checks if the repositories are starred by the authenticated user
func (c *Client) CheckStarredStatus(packages []*model.Package) error {
	for _, pkg := range packages {
		if !pkg.IsGitHub {
			continue
		}

		repoPath := pkg.GitHubRepoPath()
		if repoPath == "" {
			continue
		}

		// GitHub API returns 204 if starred, 404 if not starred
		statusCode, err := c.checkStarred(repoPath)
		if err != nil {
			return fmt.Errorf("failed to check star status for %s: %w", repoPath, err)
		}

		pkg.IsStarred = (statusCode == 204)
	}

	return nil
}

// StarRepository stars a repository
func (c *Client) StarRepository(pkg *model.Package) error {
	if !pkg.IsGitHub {
		return fmt.Errorf("not a GitHub repository: %s", pkg.Path)
	}

	repoPath := pkg.GitHubRepoPath()
	if repoPath == "" {
		return fmt.Errorf("invalid GitHub repository path: %s", pkg.Path)
	}

	err := c.putStar(repoPath)
	if err != nil {
		return fmt.Errorf("failed to star repository %s: %w", repoPath, err)
	}

	pkg.IsStarred = true
	return nil
}

// UnstarRepository unstars a repository
func (c *Client) UnstarRepository(pkg *model.Package) error {
	if !pkg.IsGitHub {
		return fmt.Errorf("not a GitHub repository: %s", pkg.Path)
	}

	repoPath := pkg.GitHubRepoPath()
	if repoPath == "" {
		return fmt.Errorf("invalid GitHub repository path: %s", pkg.Path)
	}

	err := c.deleteStar(repoPath)
	if err != nil {
		return fmt.Errorf("failed to unstar repository %s: %w", repoPath, err)
	}

	pkg.IsStarred = false
	return nil
}

// StarAllUnstarred stars all unstarred GitHub repositories
func (c *Client) StarAllUnstarred(packages []*model.Package) (int, error) {
	count := 0
	for _, pkg := range packages {
		if pkg.IsGitHub && !pkg.IsStarred {
			err := c.StarRepository(pkg)
			if err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}

// checkStarred checks if a repository is starred by the authenticated user
func (c *Client) checkStarred(repoPath string) (int, error) {
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid repository path: %s", repoPath)
	}

	owner, repo := parts[0], parts[1]
	path := fmt.Sprintf("user/starred/%s/%s", owner, repo)

	resp, err := c.restClient.Request("GET", path, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return 404, nil
		}

		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return resp.StatusCode, nil
}

// putStar stars a repository
func (c *Client) putStar(repoPath string) error {
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository path: %s", repoPath)
	}

	owner, repo := parts[0], parts[1]
	path := fmt.Sprintf("user/starred/%s/%s", owner, repo)

	resp, err := c.restClient.Request("PUT", path, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return nil
}

// deleteStar unstars a repository
func (c *Client) deleteStar(repoPath string) error {
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository path: %s", repoPath)
	}

	owner, repo := parts[0], parts[1]
	path := fmt.Sprintf("user/starred/%s/%s", owner, repo)

	resp, err := c.restClient.Request("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return nil
}
