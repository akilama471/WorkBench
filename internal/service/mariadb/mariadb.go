package mariadb

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

func (s *Service) ID() string   { return "mariadb" }
func (s *Service) Name() string { return "MariaDB" }

func (s *Service) IsInstalled() bool {
	binDir := s.resolveBinDir()
	info, err := os.Stat(binDir)
	return err == nil && info.IsDir()
}

func (s *Service) Start() error {
	s.log.Info(logger.CategoryService, "starting MariaDB")

	if !s.IsInstalled() {
		return fmt.Errorf("%s: %w", s.ID(), service.ErrServiceNotInstalled)
	}

	if s.status == service.StatusRunning {
		return fmt.Errorf("%s: %w", s.ID(), service.ErrServiceAlreadyRunning)
	}

	s.status = service.StatusStarting

	if err := s.ensureDataDir(); err != nil {
		s.status = service.StatusError
		return fmt.Errorf("failed to prepare MariaDB data directory: %w", err)
	}

	mysqld := s.resolveExecutable()
	if mysqld == "" {
		s.status = service.StatusError
		return fmt.Errorf("MariaDB mysqld executable not found")
	}

	dataDir := s.paths.MariaDBData()
	config := s.resolveConfig()

	args := []string{
		"--datadir=" + dataDir,
		"--port=3306",
	}

	if config != "" {
		args = append(args, "--defaults-file="+config)
	}

	procConfig := process.StartConfig{
		Executable: mysqld,
		Args:       args,
		Directory:  s.resolveBinDir(),
	}

	proc, err := s.process.Start(procConfig)
	if err != nil {
		s.status = service.StatusError
		return fmt.Errorf("failed to start MariaDB: %w", err)
	}

	s.pid = proc.PID
	s.status = service.StatusRunning
	s.log.Info(logger.CategoryService, "MariaDB started", "pid", s.pid)
	return nil
}

func (s *Service) Stop() error {
	s.log.Info(logger.CategoryService, "stopping MariaDB")

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
		return fmt.Errorf("failed to stop MariaDB: %w", err)
	}

	s.status = service.StatusStopped
	s.pid = 0
	s.log.Info(logger.CategoryService, "MariaDB stopped")
	return nil
}

func (s *Service) Restart() error {
	if err := s.Stop(); err != nil && !errors.Is(err, service.ErrServiceNotRunning) {
		return fmt.Errorf("failed to restart MariaDB (stop phase): %w", err)
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
	pidPath := filepath.Join(s.paths.MariaDBData(), "mariadb.pid")
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
		if d.IsDir() && (d.Name() == "mariadb" || d.Name() == "mysql") {
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
		if d.IsDir() && (strings.HasPrefix(name, "mariadb") || strings.HasPrefix(name, "mysql")) && name != "mariadb" && name != "mysql" {
			return filepath.Join(s.paths.Bin(), name)
		}
	}

	return filepath.Join(s.paths.Bin(), "mariadb")
}

func (s *Service) resolveExecutable() string {
	binDir := s.resolveBinDir()

	candidates := []string{
		filepath.Join(binDir, "mysqld.exe"),
		filepath.Join(binDir, "bin", "mysqld.exe"),
		filepath.Join(binDir, "mysqld"),
		filepath.Join(binDir, "bin", "mysqld"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func (s *Service) resolveConfig() string {
	conf := filepath.Join(s.paths.MariaDBConfig(), "my.ini")
	if _, err := os.Stat(conf); err == nil {
		return conf
	}
	conf = filepath.Join(s.paths.MariaDBConfig(), "my.cnf")
	if _, err := os.Stat(conf); err == nil {
		return conf
	}
	return ""
}

func (s *Service) ensureDataDir() error {
	dataDir := s.paths.MariaDBData()
	info, err := os.Stat(dataDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(dataDir, 0o755)
}

func (s *Service) ErrorLogPath() string {
	return filepath.Join(s.paths.MariaDBLogs(), "error.log")
}

func (s *Service) DataDirectory() string {
	return s.paths.MariaDBData()
}
