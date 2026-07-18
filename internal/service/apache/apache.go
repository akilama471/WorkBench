package apache

import (
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

	if s.status != service.StatusRunning {
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
	if err := s.Stop(); err != nil && err != service.ErrServiceNotRunning {
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
	}
	return s.status
}

func (s *Service) resolveBinDir() string {
	dirs, _ := os.ReadDir(s.paths.Bin())
	for _, d := range dirs {
		if d.IsDir() && strings.HasPrefix(d.Name(), "apache") {
			return filepath.Join(s.paths.Bin(), d.Name())
		}
	}

	for _, d := range dirs {
		if d.IsDir() && d.Name() == "apache" {
			subDirs, _ := os.ReadDir(filepath.Join(s.paths.Bin(), "apache"))
			for _, sd := range subDirs {
				if sd.IsDir() {
					return filepath.Join(s.paths.Bin(), "apache", sd.Name())
				}
			}
		}
	}

	return filepath.Join(s.paths.Bin(), "apache")
}

func (s *Service) resolveExecutable() string {
	binDir := s.resolveBinDir()
	httpd := filepath.Join(binDir, "httpd.exe")
	if _, err := os.Stat(httpd); err == nil {
		return httpd
	}
	httpd = filepath.Join(binDir, "httpd")
	if _, err := os.Stat(httpd); err == nil {
		return httpd
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

func (s *Service) ErrorLogPath() string {
	return filepath.Join(s.paths.ApacheLogs(), "error.log")
}

func (s *Service) AccessLogPath() string {
	return filepath.Join(s.paths.ApacheLogs(), "access.log")
}
