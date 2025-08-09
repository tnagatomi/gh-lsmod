# gh-lsmod

`gh-lsmod` is a gh extension which allow you to browse a project go.mod's direct dependent packages.

![demo](https://github.com/user-attachments/assets/8535d39f-95b4-464c-99d9-c985e177bde8)

## Installation

Install as a [gh](https://cli.github.com/) extension ([GitHub CLI extensions](https://cli.github.com/manual/gh_extension))

```console
gh extension install tnagatomi/gh-lsmod
```

## Usage

Navigate to a directory containing a `go.mod` and run:

```console
gh lsmod
```

## Features

- Browse direct dependencies of your project's go.mod
- Open GitHub repository in browser for GitHub-hosted packages
- Open pkg.go.dev page in browser
- Add/remove stars to GitHub repositories
