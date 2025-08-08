package size

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tnagatomi/gh-go-mod-browser/model"
)

func TestCalculatePackageSize(t *testing.T) {
	originalGOMODCACHE := os.Getenv("GOMODCACHE")
	originalGOPATH := os.Getenv("GOPATH")
	defer func() {
		os.Setenv("GOMODCACHE", originalGOMODCACHE)
		os.Setenv("GOPATH", originalGOPATH)
	}()

	tempDir, err := os.MkdirTemp("", "gomod-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	os.Setenv("GOMODCACHE", tempDir)
	os.Setenv("GOPATH", "")

	testCases := []struct {
		name          string
		pkg           *model.Package
		setupFunc     func() (int64, error)
		expectedError bool
	}{
		{
			name: "Regular package path",
			pkg:  model.NewPackage("github.com/charmbracelet/bubbles", "v0.20.0"),
			setupFunc: func() (int64, error) {
				pkgDir := filepath.Join(tempDir, "github.com", "charmbracelet", "bubbles@v0.20.0")
				if err := os.MkdirAll(pkgDir, 0755); err != nil {
					return 0, err
				}

				file1 := filepath.Join(pkgDir, "file1.go")
				if err := os.WriteFile(file1, []byte("package bubbles"), 0444); err != nil {
					return 0, err
				}

				file2 := filepath.Join(pkgDir, "file2.go")
				if err := os.WriteFile(file2, []byte("func main() {}"), 0444); err != nil {
					return 0, err
				}

				info1, err := os.Stat(file1)
				if err != nil {
					return 0, err
				}
				info2, err := os.Stat(file2)
				if err != nil {
					return 0, err
				}
				return info1.Size() + info2.Size(), nil
			},
			expectedError: false,
		},
		{
			name: "Package path with version suffix",
			pkg:  model.NewPackage("github.com/cli/go-gh/v2", "v2.11.2"),
			setupFunc: func() (int64, error) {
				pkgDir := filepath.Join(tempDir, "github.com", "cli", "go-gh", "v2@v2.11.2")
				if err := os.MkdirAll(pkgDir, 0755); err != nil {
					return 0, err
				}

				file := filepath.Join(pkgDir, "main.go")
				content := []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, world!\")\n}")
				if err := os.WriteFile(file, content, 0644); err != nil {
					return 0, err
				}

				info, err := os.Stat(file)
				if err != nil {
					return 0, err
				}
				return info.Size(), nil
			},
			expectedError: false,
		},
		{
			name:          "Non-existent package",
			pkg:           model.NewPackage("github.com/nonexistent/package", "v1.0.0"),
			setupFunc:     func() (int64, error) { return 0, nil },
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedSize, err := tc.setupFunc()
			if err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			size, err := CalculatePackageSize(tc.pkg)
			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if size != expectedSize {
					t.Errorf("Expected size %d, got %d", expectedSize, size)
				}
			}
		})
	}
}
