package mysql

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func (s *Service) ID() string   { return "mysql" }
func (s *Service) Name() string { return "MySQL" }

func (s *Service) IsInstalled() bool {
	binDir := s.resolveBinDir()
	info, err := os.Stat(binDir)
	return err == nil && info.IsDir()
}

func (s *Service) Start() error {
	s.log.Info(logger.CategoryService, "starting MySQL")

	if !s.IsInstalled() {
		return fmt.Errorf("%s: %w", s.ID(), service.ErrServiceNotInstalled)
	}

	if s.status == service.StatusRunning {
		return fmt.Errorf("%s: %w", s.ID(), service.ErrServiceAlreadyRunning)
	}

	s.status = service.StatusStarting

	if err := s.ensureDataDir(); err != nil {
		s.status = service.StatusError
		return fmt.Errorf("failed to prepare MySQL data directory: %w", err)
	}

	mysqld := s.resolveExecutable()
	if mysqld == "" {
		s.status = service.StatusError
		return fmt.Errorf("MySQL mysqld executable not found")
	}

	dataDir := s.paths.MySQLData()
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
		return fmt.Errorf("failed to start MySQL: %w", err)
	}

	s.pid = proc.PID
	s.status = service.StatusRunning
	s.log.Info(logger.CategoryService, "MySQL started", "pid", s.pid)
	return nil
}

func (s *Service) Stop() error {
	s.log.Info(logger.CategoryService, "stopping MySQL")

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
		return fmt.Errorf("failed to stop MySQL: %w", err)
	}

	s.status = service.StatusStopped
	s.pid = 0
	s.log.Info(logger.CategoryService, "MySQL stopped")
	return nil
}

func (s *Service) Restart() error {
	if err := s.Stop(); err != nil && !errors.Is(err, service.ErrServiceNotRunning) {
		return fmt.Errorf("failed to restart MySQL (stop phase): %w", err)
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
	pidPath := filepath.Join(s.paths.MySQLData(), "mysql.pid")
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
	mysqlBase := filepath.Join(s.paths.Bin(), "mysql")
	
	activeFile := filepath.Join(s.paths.Active(), "mysql")
	if data, err := os.ReadFile(activeFile); err == nil {
		activeVersion := strings.TrimSpace(string(data))
		targetDir := filepath.Join(mysqlBase, activeVersion)
		if info, err := os.Stat(targetDir); err == nil && info.IsDir() {
			return targetDir
		}
	}

	subDirs, err := os.ReadDir(mysqlBase)
	if err == nil {
		for _, sd := range subDirs {
			if sd.IsDir() {
				return filepath.Join(mysqlBase, sd.Name())
			}
		}
		return mysqlBase
	}
	return mysqlBase
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
	conf := filepath.Join(s.paths.MySQLConfig(), "my.ini")
	if _, err := os.Stat(conf); err == nil {
		return conf
	}
	conf = filepath.Join(s.paths.MySQLConfig(), "my.cnf")
	if _, err := os.Stat(conf); err == nil {
		return conf
	}
	return ""
}

func (s *Service) ensureDataDir() error {
	dataDir := s.paths.MySQLData()
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}

	// Check if the system tables exist
	mysqlDir := filepath.Join(dataDir, "mysql")
	if info, err := os.Stat(mysqlDir); err == nil && info.IsDir() {
		return nil
	}

	s.log.Info(logger.CategoryService, "initializing MySQL data directory", "dir", dataDir)

	mysqldExe := s.resolveExecutable()
	if mysqldExe == "" {
		return fmt.Errorf("mysqld executable not found")
	}

	cmd := exec.Command(mysqldExe, "--initialize-insecure", "--datadir="+dataDir)
	cmd.Dir = filepath.Dir(mysqldExe)
	if err := cmd.Run(); err != nil {
		// Fallback to older mysql_install_db for older versions
		installDbExe := filepath.Join(filepath.Dir(mysqldExe), "mysql_install_db.exe")
		if _, err2 := os.Stat(installDbExe); err2 == nil {
			cmd2 := exec.Command(installDbExe, "--datadir="+dataDir)
			cmd2.Dir = filepath.Dir(mysqldExe)
			if err3 := cmd2.Run(); err3 != nil {
				return fmt.Errorf("failed to run mysql_install_db: %w", err3)
			}
			return nil
		}
		return fmt.Errorf("failed to initialize mysql data dir: %w", err)
	}

	return nil
}

func (s *Service) ErrorLogPath() string {
	return filepath.Join(s.paths.MySQLLogs(), "error.log")
}

func (s *Service) DataDirectory() string {
	return s.paths.MySQLData()
}
