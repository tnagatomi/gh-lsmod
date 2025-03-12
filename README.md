# gh-go-mod-browser

`gh-go-mod-browser` is a gh extension which allow you to browse a project go.mod's direct dependent packages.

![demo](https://github.com/user-attachments/assets/ceb8fcf0-c6e1-487a-bf4c-1fcfe37b076e)

## Installation

Install as a [gh](https://cli.github.com/) extension ([GitHub CLI extensions](https://cli.github.com/manual/gh_extension))

```console
gh extension install tnagatomi/gh-go-mod-browser
```

## Usage

Navigate to a directory containing a `go.mod` and run:

```console
gh go-mod-browser
```

## Features

- Browse direct dependencies of your project's go.mod
- Open GitHub repository in browser for GitHub-hosted packages
- Open pkg.go.dev page in browser
- Add/remove stars to GitHub repositories
