package collection

import (
	"errors"
	"os/exec"
)

const (
	PackageManagerRPM    = "rpm"
	PackageManagerDebian = "dpkg"
)

// ErrNoPackageManager is returned when we could not detect one.
var ErrNoPackageManager = errors.New("could not detect a supported package manager")

var FoundPackageManager string

func DetectPackageManager() string {
	if FoundPackageManager != "" {
		return FoundPackageManager
	}

	priority := []string{PackageManagerDebian, PackageManagerRPM}

	for _, name := range priority {
		if _, err := exec.LookPath(name); err == nil {
			FoundPackageManager = name
			return name
		}
	}

	return ""
}

func ListInstalledPackagesRaw(pattern string) ([]byte, error) {
	switch DetectPackageManager() {
	case PackageManagerRPM:
		return LoadCommandOutput("rpm", "-qa", pattern)
	case PackageManagerDebian:
		return LoadCommandOutput("dpkg-query", "-f", "${Package} ${Version} ${Architecture} ${Status}\\n", "-W", pattern)
	default:
		return []byte{}, ErrNoPackageManager
	}
}
