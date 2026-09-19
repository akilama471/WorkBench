package pkg

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
)

func createTestZip(t *testing.T, zipPath string, files map[string]string) {
	t.Helper()
	outFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create test zip: %v", err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			t.Fatalf("failed to add file %s to zip: %v", name, err)
		}
		if _, err := f.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write file %s to zip: %v", name, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}
}

func TestInstallLocalZip(t *testing.T) {
	root := t.TempDir()
	paths := filesystem.NewPaths(root)
	log := logger.New(logger.LevelDebug, os.Stdout)
	mgr := NewManager(paths, log)

	zipDir := t.TempDir()

	// Test PHP ZIP
	phpZip := filepath.Join(zipDir, "php-8.3.30-Win32-vs16-x64.zip")
	createTestZip(t, phpZip, map[string]string{
		"php.exe":             "dummy php binary",
		"php.ini-development": "; php ini dev",
	})

	version, err := mgr.InstallLocalZip("php", phpZip)
	if err != nil {
		t.Fatalf("InstallLocalZip php failed: %v", err)
	}
	if version != "8.3.30" {
		t.Errorf("expected version 8.3.30, got %s", version)
	}

	targetPHPDir := paths.PHPBin("8.3.30")
	if _, err := os.Stat(filepath.Join(targetPHPDir, "php.exe")); err != nil {
		t.Errorf("php.exe should exist in target directory %s", targetPHPDir)
	}
	if _, err := os.Stat(filepath.Join(targetPHPDir, "php.ini")); err != nil {
		t.Errorf("php.ini should be auto-created in target directory %s", targetPHPDir)
	}

	// Test Apache ZIP
	apacheZip := filepath.Join(zipDir, "httpd-2.4.63-win64-VS17.zip")
	createTestZip(t, apacheZip, map[string]string{
		"Apache24/bin/httpd.exe":     "dummy httpd binary",
		"Apache24/conf/httpd.conf":  "# dummy httpd conf",
	})

	versionApache, err := mgr.InstallLocalZip("apache", apacheZip)
	if err != nil {
		t.Fatalf("InstallLocalZip apache failed: %v", err)
	}
	if versionApache != "2.4.63" {
		t.Errorf("expected version 2.4.63, got %s", versionApache)
	}

	targetApacheDir := paths.ApacheBin("2.4.63")
	if _, err := os.Stat(filepath.Join(targetApacheDir, "bin", "httpd.exe")); err != nil {
		t.Errorf("httpd.exe should exist in target directory %s", targetApacheDir)
	}

	// Test MariaDB ZIP
	mariadbZip := filepath.Join(zipDir, "mariadb-11.7.2-winx64.zip")
	createTestZip(t, mariadbZip, map[string]string{
		"mariadb-11.7.2-winx64/bin/mariadbd.exe": "dummy mariadbd binary",
	})

	versionMaria, err := mgr.InstallLocalZip("mariadb", mariadbZip)
	if err != nil {
		t.Fatalf("InstallLocalZip mariadb failed: %v", err)
	}
	if versionMaria != "11.7.2" {
		t.Errorf("expected version 11.7.2, got %s", versionMaria)
	}

	targetMariaDir := paths.MariaDBBin("11.7.2")
	if _, err := os.Stat(filepath.Join(targetMariaDir, "bin", "mariadbd.exe")); err != nil {
		t.Errorf("mariadbd.exe should exist in target directory %s", targetMariaDir)
	}

	// Test MySQL ZIP
	mysqlZip := filepath.Join(zipDir, "mysql-8.0.36-winx64.zip")
	createTestZip(t, mysqlZip, map[string]string{
		"mysql-8.0.36-winx64/bin/mysqld.exe": "dummy mysqld binary",
	})

	versionMySQL, err := mgr.InstallLocalZip("mysql", mysqlZip)
	if err != nil {
		t.Fatalf("InstallLocalZip mysql failed: %v", err)
	}
	if versionMySQL != "8.0.36" {
		t.Errorf("expected version 8.0.36, got %s", versionMySQL)
	}

	targetMySQLDir := paths.MySQLBin("8.0.36")
	if _, err := os.Stat(filepath.Join(targetMySQLDir, "bin", "mysqld.exe")); err != nil {
		t.Errorf("mysqld.exe should exist in target directory %s", targetMySQLDir)
	}
}

