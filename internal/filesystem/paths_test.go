package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPaths(t *testing.T) {
	root := t.TempDir()
	paths := NewPaths(root)

	tests := []struct {
		method   string
		expected string
	}{
		{"Root", root},
		{"Bin", filepath.Join(root, "bin")},
		{"Active", filepath.Join(root, "active")},
		{"WWW", filepath.Join(root, "www")},
		{"Data", filepath.Join(root, "data")},
		{"Etc", filepath.Join(root, "etc")},
		{"Logs", filepath.Join(root, "logs")},
		{"Backup", filepath.Join(root, "backup")},
		{"Packages", filepath.Join(root, "packages")},
		{"Cache", filepath.Join(root, "cache")},
		{"Database", filepath.Join(root, "devbox.db")},
		{"ApacheBin", filepath.Join(root, "bin", "apache")},
		{"PHPBin", filepath.Join(root, "bin", "php")},
		{"MariaDBBin", filepath.Join(root, "bin", "mariadb")},
		{"ActivePHP", filepath.Join(root, "active", "php")},
		{"MariaDBData", filepath.Join(root, "data", "mariadb")},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			var got string
			switch tt.method {
			case "Root":
				got = paths.Root()
			case "Bin":
				got = paths.Bin()
			case "Active":
				got = paths.Active()
			case "WWW":
				got = paths.WWW()
			case "Data":
				got = paths.Data()
			case "Etc":
				got = paths.Etc()
			case "Logs":
				got = paths.Logs()
			case "Backup":
				got = paths.Backup()
			case "Packages":
				got = paths.Packages()
			case "Cache":
				got = paths.Cache()
			case "Database":
				got = paths.Database()
			case "ApacheBin":
				got = paths.ApacheBin("")
			case "PHPBin":
				got = paths.PHPBin("")
			case "MariaDBBin":
				got = paths.MariaDBBin("")
			case "ActivePHP":
				got = paths.ActivePHP()
			case "MariaDBData":
				got = paths.MariaDBData()
			}
			if got != tt.expected {
				t.Errorf("paths.%s() = %q, want %q", tt.method, got, tt.expected)
			}
		})
	}
}

func TestPathsVersioned(t *testing.T) {
	root := t.TempDir()
	paths := NewPaths(root)

	if got := paths.ApacheBin("2.4.50"); got != filepath.Join(root, "bin", "apache", "2.4.50") {
		t.Errorf("ApacheBin(\"2.4.50\") = %q", got)
	}

	if got := paths.PHPBin("8.3.0"); got != filepath.Join(root, "bin", "php", "8.3.0") {
		t.Errorf("PHPBin(\"8.3.0\") = %q", got)
	}

	if got := paths.MariaDBBin("11.4"); got != filepath.Join(root, "bin", "mariadb", "11.4") {
		t.Errorf("MariaDBBin(\"11.4\") = %q", got)
	}
}

func TestInitializer(t *testing.T) {
	root := t.TempDir()
	paths := NewPaths(root)
	init := NewInitializer(paths)

	if err := init.Initialize(); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	requiredDirs := []string{
		"bin", "active", "www", "data", "etc", "logs",
		"backup", "packages", "cache",
	}

	for _, dir := range requiredDirs {
		path := filepath.Join(root, dir)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("directory %s not created: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", dir)
		}
	}
}

func TestInitializerIdempotent(t *testing.T) {
	root := t.TempDir()
	paths := NewPaths(root)
	init := NewInitializer(paths)

	if err := init.Initialize(); err != nil {
		t.Fatalf("first Initialize() failed: %v", err)
	}

	testFile := filepath.Join(root, "www", "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if err := init.Initialize(); err != nil {
		t.Fatalf("second Initialize() failed: %v", err)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("test file missing after second init: %v", err)
	}
	if string(data) != "test" {
		t.Errorf("test file content changed: %q", string(data))
	}
}
