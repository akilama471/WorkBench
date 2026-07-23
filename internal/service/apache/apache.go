package apache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
	"github.com/akilama471/WorkBench/internal/process"
	"github.com/akilama471/WorkBench/internal/service"
)

type Service struct {
	paths   *filesystem.Paths
	process process.Manager
	log     *logger.Logger
	status  service.Status
	pid     int
}

func NewService(paths *filesystem.Paths, proc process.Manager, log *logger.Logger) *Service {
	return &Service{
		paths:   paths,
		process: proc,
		log:     log,
		status:  service.StatusUnknown,
	}
}

func (s *Service) ID() string   { return "apache" }
func (s *Service) Name() string { return "Apache" }

func (s *Service) IsInstalled() bool {
	binDir := s.resolveBinDir()
	info, err := os.Stat(binDir)
	return err == nil && info.IsDir()
}

func (s *Service) Start() error {
	s.log.Info(logger.CategoryService, "starting Apache")

	if !s.IsInstalled() {
		return fmt.Errorf("%s: %w", s.ID(), service.ErrServiceNotInstalled)
	}

	if s.status == service.StatusRunning {
		return fmt.Errorf("%s: %w", s.ID(), service.ErrServiceAlreadyRunning)
	}

	s.status = service.StatusStarting

	httpd := s.resolveExecutable()
	if httpd == "" {
		s.status = service.StatusError
		return fmt.Errorf("Apache httpd executable not found")
	}

	if err := s.ensureConfig(); err != nil {
		s.status = service.StatusError
		return fmt.Errorf("failed to prepare Apache config: %w", err)
	}

	conf := s.resolveConfig()
	if conf == "" {
		s.status = service.StatusError
		return fmt.Errorf("Apache configuration file not found")
	}

	config := process.StartConfig{
		Executable: httpd,
		Args:       []string{"-f", conf},
		Directory:  s.resolveBinDir(),
	}

	proc, err := s.process.Start(config)
	if err != nil {
		s.status = service.StatusError
		return fmt.Errorf("failed to start Apache: %w", err)
	}

	s.pid = proc.PID
	s.status = service.StatusRunning
	s.log.Info(logger.CategoryService, "Apache started", "pid", s.pid)
	return nil
}

func (s *Service) Stop() error {
	s.log.Info(logger.CategoryService, "stopping Apache")

	if s.pid == 0 {
		s.pid = s.readPIDFile()
	}

	if s.pid == 0 || !s.process.IsRunning(s.pid) {
		s.status = service.StatusStopped
		return fmt.Errorf("%s: %w", s.ID(), service.ErrServiceNotRunning)
	}

	s.status = service.StatusStopping

	if err := s.process.Stop(s.pid); err != nil {
		s.status = service.StatusError
		return fmt.Errorf("failed to stop Apache: %w", err)
	}

	s.status = service.StatusStopped
	s.pid = 0
	s.log.Info(logger.CategoryService, "Apache stopped")
	return nil
}

func (s *Service) Restart() error {
	if err := s.Stop(); err != nil && !errors.Is(err, service.ErrServiceNotRunning) {
		return fmt.Errorf("failed to restart Apache (stop phase): %w", err)
	}
	return s.Start()
}

func (s *Service) Status() service.Status {
	if s.status == service.StatusRunning && s.pid > 0 {
		if !s.process.IsRunning(s.pid) {
			s.status = service.StatusStopped
			s.pid = 0
		}
		return s.status
	}

	if pid := s.readPIDFile(); pid > 0 && s.process.IsRunning(pid) {
		s.pid = pid
		s.status = service.StatusRunning
		return s.status
	}

	return s.status
}

func (s *Service) readPIDFile() int {
	pidPath := filepath.Join(s.resolveBinDir(), "logs", "httpd.pid")
	data, err := os.ReadFile(pidPath)
	if err != nil {
		return 0
	}
	var pid int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(data)), "%d", &pid); err != nil {
		return 0
	}
	return pid
}

func (s *Service) resolveBinDir() string {
	dirs, _ := os.ReadDir(s.paths.Bin())

	for _, d := range dirs {
		if d.IsDir() && (d.Name() == "apache" || d.Name() == "httpd") {
			subDirs, _ := os.ReadDir(filepath.Join(s.paths.Bin(), d.Name()))
			for _, sd := range subDirs {
				if sd.IsDir() {
					return filepath.Join(s.paths.Bin(), d.Name(), sd.Name())
				}
			}
		}
	}

	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() && (strings.HasPrefix(name, "apache") || strings.HasPrefix(name, "httpd")) && name != "apache" && name != "httpd" {
			return filepath.Join(s.paths.Bin(), name)
		}
	}

	return filepath.Join(s.paths.Bin(), "apache")
}

func (s *Service) resolveExecutable() string {
	binDir := s.resolveBinDir()

	candidates := []string{
		filepath.Join(binDir, "httpd.exe"),
		filepath.Join(binDir, "bin", "httpd.exe"),
		filepath.Join(binDir, "httpd"),
		filepath.Join(binDir, "bin", "httpd"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func (s *Service) resolveConfig() string {
	conf := filepath.Join(s.paths.ApacheConfig(), "httpd.conf")
	if _, err := os.Stat(conf); err == nil {
		return conf
	}
	return ""
}

func (s *Service) resolveDistributionConfig() string {
	binDir := s.resolveBinDir()
	conf := filepath.Join(binDir, "conf", "httpd.conf")
	if _, err := os.Stat(conf); err == nil {
		return conf
	}
	return ""
}

func (s *Service) ensureConfig() error {
	if s.resolveConfig() != "" {
		return nil
	}

	distConf := s.resolveDistributionConfig()
	if distConf == "" {
		return fmt.Errorf("no httpd.conf found in distribution or etc/apache")
	}

	data, err := os.ReadFile(distConf)
	if err != nil {
		return fmt.Errorf("failed to read distribution httpd.conf: %w", err)
	}

	binDir := s.resolveBinDir()
	wwwDir := s.paths.Root()
	etcDir := s.paths.ApacheConfig()
	logsDir := s.paths.ApacheLogs()

	content := string(data)
	content = replaceConfigValue(content, "Define SRVROOT", fmt.Sprintf(`Define SRVROOT "%s"`, filepath.ToSlash(binDir)))
	content = replaceConfigValue(content, "ServerRoot", fmt.Sprintf(`ServerRoot "%s"`, filepath.ToSlash(binDir)))
	content = replaceConfigValue(content, "DocumentRoot", fmt.Sprintf(`DocumentRoot "%s"`, filepath.ToSlash(filepath.Join(wwwDir, "www"))))
	content = replaceConfigValueDir(content, `<Directory "${SRVROOT}/htdocs">`, fmt.Sprintf(`<Directory "%s">`, filepath.ToSlash(filepath.Join(wwwDir, "www"))))
	content = replaceConfigLine(content, "ErrorLog", fmt.Sprintf(`ErrorLog "%s"`, filepath.ToSlash(filepath.Join(logsDir, "error.log"))))
	content = replaceConfigLine(content, "CustomLog", fmt.Sprintf(`CustomLog "%s" common`, filepath.ToSlash(filepath.Join(logsDir, "access.log"))))

	if !strings.Contains(content, "ServerName localhost") {
		content += "\nServerName localhost\n"
	}

	if err := os.MkdirAll(etcDir, 0o755); err != nil {
		return fmt.Errorf("failed to create etc/apache directory: %w", err)
	}

	if err := os.WriteFile(filepath.Join(etcDir, "httpd.conf"), []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write httpd.conf: %w", err)
	}

	s.log.Info(logger.CategoryService, "Apache config bootstrapped from distribution", "target", filepath.Join(etcDir, "httpd.conf"))
	return nil
}

func (s *Service) ErrorLogPath() string {
	return filepath.Join(s.paths.ApacheLogs(), "error.log")
}

func (s *Service) AccessLogPath() string {
	return filepath.Join(s.paths.ApacheLogs(), "access.log")
}

func (s *Service) ConfigurePHP(phpVersionPath string) error {
	conf := s.resolveConfig()
	if conf == "" {
		return fmt.Errorf("Apache configuration file not found")
	}

	if err := s.updatePHPModule(conf, phpVersionPath); err != nil {
		return fmt.Errorf("failed to update PHP module in Apache config: %w", err)
	}

	s.log.Info(logger.CategoryService, "Apache PHP configured", "php_path", phpVersionPath)
	return nil
}

func (s *Service) updatePHPModule(confPath, phpVersionPath string) error {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return fmt.Errorf("failed to read Apache config: %w", err)
	}

	content := string(data)
	phpModule := filepath.Join(phpVersionPath, "php8apache2_4.dll")
	phpIniDir := phpVersionPath

	lines := strings.Split(content, "\n")
	var result []string
	foundModule := false
	foundIniDir := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "LoadModule") && strings.Contains(trimmed, "php") {
			result = append(result, fmt.Sprintf(`LoadModule php_module "%s"`, filepath.ToSlash(phpModule)))
			foundModule = true
			continue
		}

		if strings.HasPrefix(trimmed, "PHPIniDir") {
			result = append(result, fmt.Sprintf(`PHPIniDir "%s"`, filepath.ToSlash(phpIniDir)))
			foundIniDir = true
			continue
		}

		result = append(result, line)
	}

	if !foundModule {
		result = append(result, fmt.Sprintf(`LoadModule php_module "%s"`, filepath.ToSlash(phpModule)))
	}
	if !foundIniDir {
		result = append(result, fmt.Sprintf(`PHPIniDir "%s"`, filepath.ToSlash(phpIniDir)))
	}

	return os.WriteFile(confPath, []byte(strings.Join(result, "\n")), 0o644)
}

func replaceConfigValue(content, directive, newValue string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, directive+" ") || strings.HasPrefix(trimmed, directive+"\t") {
			lines[i] = newValue
		}
	}
	return strings.Join(lines, "\n")
}

func replaceConfigValueDir(content, search, newValue string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, search) {
			lines[i] = newValue
			break
		}
	}
	return strings.Join(lines, "\n")
}

func replaceConfigLine(content, directive, newValue string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, directive+" ") || strings.HasPrefix(trimmed, directive+"\t") {
			lines[i] = newValue
			break
		}
	}
	return strings.Join(lines, "\n")
}

func removeConfigLines(content, directive string) string {
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) == directive {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
