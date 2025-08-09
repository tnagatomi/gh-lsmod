package size

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tnagatomi/gh-lsmod/model"
)
// CalculatePackageSize calculates the size of a package
func CalculatePackageSize(pkg *model.Package) (int64, error) {
	// Get GOMODCACHE or GOPATH
	goModCache := os.Getenv("GOMODCACHE")
	if goModCache == "" {
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return 0, fmt.Errorf("failed to get user home directory: %w", err)
			}
			goPath = filepath.Join(home, "go")
		}
		goModCache = filepath.Join(goPath, "pkg", "mod")
	}

	// Construct package path
	pkgPath := pkg.Path
	if pkg.Version != "" {
		pkgPath = fmt.Sprintf("%s@%s", pkgPath, pkg.Version)
	}

	// Replace / with OS path separator
	pkgPath = strings.ReplaceAll(pkgPath, "/", string(os.PathSeparator))
	pkgPath = filepath.Join(goModCache, pkgPath)

	// Check if directory exists
	_, err := os.Stat(pkgPath)
	if err != nil {
		return 0, fmt.Errorf("failed to stat package directory: %w", err)
	}

	var size int64
	err = filepath.WalkDir(pkgPath, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		
		return nil
	})
	
	return size, err
}
