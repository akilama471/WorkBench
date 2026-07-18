package filesystem

import (
	"fmt"
	"os"
)

type Initializer struct {
	paths *Paths
}

func NewInitializer(paths *Paths) *Initializer {
	return &Initializer{paths: paths}
}

func (i *Initializer) Initialize() error {
	dirs := []string{
		i.paths.Bin(),
		i.paths.Active(),
		i.paths.WWW(),
		i.paths.Data(),
		i.paths.Etc(),
		i.paths.Logs(),
		i.paths.Backup(),
		i.paths.Packages(),
		i.paths.Cache(),
		i.paths.CacheDownloads(),
		i.paths.CacheExtract(),
		i.paths.ApacheConfig(),
		i.paths.PHPConfig(),
		i.paths.MariaDBConfig(),
		i.paths.MariaDBData(),
		i.paths.ApacheLogs(),
		i.paths.PHPLogs(),
		i.paths.MariaDBLogs(),
	}

	for _, dir := range dirs {
		if err := ensureDir(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func ensureDir(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return nil
		}
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}
	return os.MkdirAll(path, 0o755)
}
