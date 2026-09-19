package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
)

var versionRegex = regexp.MustCompile(`\b(\d+\.\d+(?:\.\d+)?)\b`)

// InstallLocalZip extracts a local ZIP archive for a service (apache, mariadb, php)
// and installs it to application-root/bin/<service>/<version>/.
func (m *Manager) InstallLocalZip(serviceType, zipPath string) (string, error) {
	svc := strings.ToLower(strings.TrimSpace(serviceType))
	if svc != "apache" && svc != "mariadb" && svc != "mysql" && svc != "php" && svc != "httpd" {
		return "", fmt.Errorf("unsupported service type '%s': expected apache, mariadb, mysql, or php", serviceType)
	}
	if svc == "httpd" {
		svc = "apache"
	}

	zipAbs, err := filepath.Abs(zipPath)
	if err != nil {
		return "", fmt.Errorf("invalid zip path: %w", err)
	}

	info, err := os.Stat(zipAbs)
	if err != nil || info.IsDir() {
		return "", fmt.Errorf("zip file not found: %s", zipPath)
	}

	tempExtractDir := filepath.Join(m.paths.CacheExtract(), fmt.Sprintf("%s-install-%d", svc, time.Now().UnixNano()))
	defer os.RemoveAll(tempExtractDir)

	m.log.Info(logger.CategoryPackage, "extracting zip archive", "service", svc, "zip", zipAbs)
	if err := filesystem.ExtractZip(zipAbs, tempExtractDir); err != nil {
		return "", fmt.Errorf("failed to extract zip archive: %w", err)
	}

	sourceDir := unwrapSingleRootFolder(tempExtractDir)
	version := detectServiceVersion(svc, zipAbs, sourceDir)

	var targetDir string
	switch svc {
	case "apache":
		targetDir = m.paths.ApacheBin(version)
	case "mariadb":
		targetDir = m.paths.MariaDBBin(version)
	case "mysql":
		targetDir = m.paths.MySQLBin(version)
	case "php":
		targetDir = m.paths.PHPBin(version)
	}

	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		return "", fmt.Errorf("failed to create target bin directory: %w", err)
	}

	if _, err := os.Stat(targetDir); err == nil {
		m.log.Info(logger.CategoryPackage, "removing existing installation at target directory", "path", targetDir)
		os.RemoveAll(targetDir)
	}

	m.log.Info(logger.CategoryPackage, "installing package to bin directory", "service", svc, "version", version, "target", targetDir)
	if err := filesystem.MoveOrCopyDir(sourceDir, targetDir); err != nil {
		return "", fmt.Errorf("failed to move extracted files to target directory %s: %w", targetDir, err)
	}

	if err := m.postInstallSetup(svc, targetDir); err != nil {
		m.log.Warn(logger.CategoryPackage, "post-installation setup warning", "service", svc, "error", err)
	}

	m.log.Info(logger.CategoryPackage, "package installed successfully", "service", svc, "version", version)
	return version, nil
}

func unwrapSingleRootFolder(dir string) string {
	entries, err := os.ReadDir(dir)
	if err == nil && len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(dir, entries[0].Name())
	}
	return dir
}

func detectServiceVersion(svc, zipPath, sourceDir string) string {
	zipName := filepath.Base(zipPath)
	if match := versionRegex.FindString(zipName); match != "" {
		return match
	}

	sourceName := filepath.Base(sourceDir)
	if match := versionRegex.FindString(sourceName); match != "" {
		return match
	}

	switch svc {
	case "apache":
		return "2.4.63"
	case "mariadb":
		return "11.7.2"
	case "mysql":
		return "8.0.36"
	case "php":
		return "8.3.30"
	default:
		return "1.0.0"
	}
}

func (m *Manager) postInstallSetup(svc, targetDir string) error {
	switch svc {
	case "php":
		phpIni := filepath.Join(targetDir, "php.ini")
		if _, err := os.Stat(phpIni); os.IsNotExist(err) {
			devIni := filepath.Join(targetDir, "php.ini-development")
			prodIni := filepath.Join(targetDir, "php.ini-production")
			if _, err := os.Stat(devIni); err == nil {
				data, _ := os.ReadFile(devIni)
				os.WriteFile(phpIni, data, 0o644)
			} else if _, err := os.Stat(prodIni); err == nil {
				data, _ := os.ReadFile(prodIni)
				os.WriteFile(phpIni, data, 0o644)
			}
		}
	case "apache":
		os.MkdirAll(filepath.Join(targetDir, "logs"), 0o755)
	case "mariadb", "mysql":
		os.MkdirAll(filepath.Join(targetDir, "data"), 0o755)
	}
	return nil
}

