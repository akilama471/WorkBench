package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akilama471/WorkBench/internal/config"
	"github.com/akilama471/WorkBench/internal/core"
	"github.com/akilama471/WorkBench/internal/database"
	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
	pkg "github.com/akilama471/WorkBench/internal/package"
	"github.com/akilama471/WorkBench/internal/process"
	"github.com/akilama471/WorkBench/internal/project"
	"github.com/akilama471/WorkBench/internal/service/apache"
	"github.com/akilama471/WorkBench/internal/service/mariadb"
	"github.com/akilama471/WorkBench/pkg/events"
)

type Application struct {
	log     *logger.Logger
	paths   *filesystem.Paths
	events  *events.Bus
	db      *database.Database
	process process.Manager
	config  *config.Manager

	ServiceManager *core.ServiceManager
	RuntimeManager *core.RuntimeManager
	Environment    *core.EnvironmentManager
	ProjectManager *project.Manager
	PackageManager *pkg.Manager
	ConfigManager  *config.Manager
	ConfigWriter   *config.Writer
	ConfigLoader   *config.Loader
}

func New(rootDir string) (*Application, error) {
	log := logger.New(logger.LevelInfo, os.Stderr)

	paths := filesystem.NewPaths(rootDir)

	eventBus := events.NewBus()

	processMgr := process.NewManager()

	configMgr := config.NewManager()
	configWriter := config.NewWriter(configMgr)
	configLoader := config.NewLoader(configMgr)

	envManager := core.NewEnvironmentManager(paths, log)

	serviceManager := core.NewServiceManager()

	apacheSvc := apache.NewService(paths, processMgr, log)
	mariadbSvc := mariadb.NewService(paths, processMgr, log)

	if err := serviceManager.Register(apacheSvc); err != nil {
		log.Warn(logger.CategoryApplication, "failed to register Apache service", "error", err)
	}
	if err := serviceManager.Register(mariadbSvc); err != nil {
		log.Warn(logger.CategoryApplication, "failed to register MariaDB service", "error", err)
	}

	runtimeManager := core.NewRuntimeManager(paths, log)

	projectManager := project.NewManager()
	packageManager := pkg.NewManager(paths, log)

	app := &Application{
		log:            log,
		paths:          paths,
		events:         eventBus,
		process:        processMgr,
		config:         configMgr,
		ServiceManager: serviceManager,
		RuntimeManager: runtimeManager,
		Environment:    envManager,
		ProjectManager: projectManager,
		PackageManager: packageManager,
		ConfigManager:  configMgr,
		ConfigWriter:   configWriter,
		ConfigLoader:   configLoader,
	}

	return app, nil
}

func (a *Application) InitializeEnvironment() error {
	a.log.Info(logger.CategoryApplication, "initializing environment", "root", a.paths.Root())
	return a.Environment.Initialize()
}

func (a *Application) OpenDatabase() error {
	dbPath := a.paths.Database()
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := database.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Initialize(); err != nil {
		db.Close()
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	a.db = db
	a.log.Info(logger.CategoryDatabase, "database opened", "path", dbPath)
	return nil
}

func (a *Application) Database() *database.Database {
	return a.db
}

func (a *Application) Paths() *filesystem.Paths {
	return a.paths
}

func (a *Application) Logger() *logger.Logger {
	return a.log
}

func (a *Application) Startup() error {
	if err := a.InitializeEnvironment(); err != nil {
		return fmt.Errorf("failed to initialize environment: %w", err)
	}

	if err := a.OpenDatabase(); err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	a.log.Info(logger.CategoryApplication, "application started")
	return nil
}

func (a *Application) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func (a *Application) StartService(id string) error {
	a.log.Info(logger.CategoryService, "starting service", "id", id)

	if id == "apache" {
		if phpVersion, err := a.CurrentPHPVersion(); err == nil && phpVersion != "" && phpVersion != "none" {
			if apacheSvc, err := a.ServiceManager.GetService("apache"); err == nil {
				if configurable, ok := apacheSvc.(interface{ ConfigurePHP(string) error }); ok {
					phpPath := filepath.Dir(a.RuntimeManager.PHP().PHPBinaryPath(phpVersion))
					if confErr := configurable.ConfigurePHP(phpPath); confErr != nil {
						a.log.Warn(logger.CategoryService, "failed to auto-configure Apache for PHP before starting", "error", confErr)
					}
				}
			}
		}
	}

	err := a.ServiceManager.Start(id)
	if err == nil {
		a.events.Publish(events.Event{Type: events.ServiceStarted, Payload: id})
	} else {
		a.events.Publish(events.Event{Type: events.ServiceError, Payload: map[string]string{"service": id, "error": err.Error()}})
	}
	return err
}

func (a *Application) StopService(id string) error {
	a.log.Info(logger.CategoryService, "stopping service", "id", id)
	err := a.ServiceManager.Stop(id)
	if err == nil {
		a.events.Publish(events.Event{Type: events.ServiceStopped, Payload: id})
	} else {
		a.events.Publish(events.Event{Type: events.ServiceError, Payload: map[string]string{"service": id, "error": err.Error()}})
	}
	return err
}

func (a *Application) RestartService(id string) error {
	a.log.Info(logger.CategoryService, "restarting service", "id", id)
	
	if err := a.StopService(id); err != nil && err.Error() != "service is not running" && !strings.Contains(err.Error(), "service is not running") {
		// We ignore "not running" errors because we just want to ensure it's stopped before starting
		a.log.Warn(logger.CategoryService, "warning during stop phase of restart", "error", err)
	}

	err := a.StartService(id)
	if err == nil {
		a.events.Publish(events.Event{Type: events.ServiceStatusChanged, Payload: id})
	}
	return err
}

func (a *Application) ServiceStatus(id string) (string, error) {
	status, err := a.ServiceManager.Status(id)
	if err != nil {
		return "", err
	}
	return status.String(), nil
}

func (a *Application) ListPHPVersions() ([]string, error) {
	versions, err := a.RuntimeManager.ListPHPVersions()
	if err != nil {
		return nil, err
	}

	var result []string
	for _, v := range versions {
		result = append(result, v.Version)
	}
	return result, nil
}

func (a *Application) CurrentPHPVersion() (string, error) {
	return a.RuntimeManager.CurrentPHPVersion()
}

func (a *Application) SwitchPHPVersion(version string) error {
	a.log.Info(logger.CategoryRuntime, "switching PHP version", "target", version)

	err := a.RuntimeManager.SwitchPHPVersion(version)
	if err != nil {
		a.events.Publish(events.Event{Type: events.ServiceError, Payload: map[string]string{"component": "php", "error": err.Error()}})
		return err
	}

	if apacheSvc, err := a.ServiceManager.GetService("apache"); err == nil {
		if configurable, ok := apacheSvc.(interface{ ConfigurePHP(string) error }); ok {
			phpPath := filepath.Dir(a.RuntimeManager.PHP().PHPBinaryPath(version))
			if confErr := configurable.ConfigurePHP(phpPath); confErr != nil {
				a.log.Warn(logger.CategoryRuntime, "Failed to configure Apache for PHP", "error", confErr)
				return fmt.Errorf("failed to configure Apache for PHP: %w", confErr)
			}
		}
	}

	a.events.Publish(events.Event{Type: events.PHPVersionChanged, Payload: version})

	if a.isServiceRunning("apache") {
		a.log.Info(logger.CategoryRuntime, "restarting Apache for PHP change", "php", version)
		if restartErr := a.RestartService("apache"); restartErr != nil {
			a.log.Warn(logger.CategoryRuntime, "Apache restart after PHP switch failed", "error", restartErr)
		}
	}

	return nil
}

func (a *Application) isServiceRunning(id string) bool {
	status, err := a.ServiceManager.Status(id)
	if err != nil {
		return false
	}
	return status.String() == "Running"
}

func (a *Application) GetEnvironmentStatus() map[string]string {
	summary := make(map[string]string)

	for _, svc := range a.ServiceManager.AllServices() {
		summary[svc.ID()] = svc.Status().String()
	}

	phpVersion, err := a.CurrentPHPVersion()
	if err != nil || phpVersion == "" || phpVersion == "none" {
		summary["php"] = "none"
	} else {
		summary["php"] = phpVersion
	}

	return summary
}

func (a *Application) Events() *events.Bus {
	return a.events
}
