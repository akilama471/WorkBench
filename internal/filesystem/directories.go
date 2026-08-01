package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
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

	indexPath := filepath.Join(i.paths.WWW(), "index.php")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		defaultIndex := `<?php
$php_version = phpversion();
$server_software = $_SERVER['SERVER_SOFTWARE'] ?? 'Unknown';
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WorkBench</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f0f2f5; color: #333; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
        .container { background: white; padding: 40px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; max-width: 600px; }
        h1 { color: #2c3e50; margin-bottom: 10px; }
        p { color: #7f8c8d; line-height: 1.6; }
        .info { background: #ecf0f1; padding: 15px; border-radius: 5px; margin-top: 20px; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to WorkBench</h1>
        <p>Your local development environment is up and running successfully!</p>
        <p>To get started, replace this file with your own project files in the <strong>www</strong> directory.</p>
        <div class="info">
            <strong>Server:</strong> <?php echo htmlspecialchars($server_software); ?><br>
            <strong>PHP Version:</strong> <?php echo htmlspecialchars($php_version); ?>
        </div>
    </div>
</body>
</html>`
		os.WriteFile(indexPath, []byte(defaultIndex), 0o644)
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
