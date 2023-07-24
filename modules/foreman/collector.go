package foreman

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"os"
	"path/filepath"
)

const ModuleName = "foreman"

var files = []string{
	"/etc/foreman",
	"/etc/foreman-installer",
	"/etc/foreman-proxy",
}

var detailedFiles = []string{
	"/var/log/foreman",
	"/var/log/foreman-installer",
	"/var/log/foreman-proxy",
}

func detect() bool {
	_, err := os.Stat("/etc/foreman")
	return err == nil
}

var obfuscaters = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password)\s*:\s*(.*)`, "yml"),
	obfuscate.NewFile(`(?i)(?:ENCRYPTION_KEY)\s*=\s*(.*)`, "rb"),
}

func Collect(c *collection.Collection) {
	if !detect() {
		c.Log.Info("Could not find Foreman")
		return
	}

	c.Log.Info("Collection Foreman information")

	c.RegisterObfuscators(obfuscaters...)

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"),
		"foreman",
		"foreman-installer",
		"foreman-proxy",
	)

	c.AddServiceStatusRaw(filepath.Join(ModuleName, "service.txt"), "foreman")

	if collection.DetectServiceManager() == "systemd" {
		c.AddCommandOutput(filepath.Join(ModuleName, "systemd-foreman.service"), "systemctl", "cat", "foreman.service")
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFiles(ModuleName, file)
		}
	}
}
