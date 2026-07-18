package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/akilama471/WorkBench/internal/config"
	"github.com/akilama471/WorkBench/internal/core"
	"github.com/akilama471/WorkBench/internal/database"
	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
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
		ConfigManager:  configMgr,
		ConfigWriter:   configWriter,
		ConfigLoader:   configLoader,
	}

	return app, nil
}

func (a *Application) InitializeEnvironment() error {
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
	return nil
}

func (a *Application) Database() *database.Database {
	return a.db
}

func (a *Application) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func (a *Application) StartService(id string) error {
	err := a.ServiceManager.Start(id)
	if err == nil {
		a.events.Publish(events.Event{Type: events.ServiceStarted, Payload: id})
	} else {
		a.events.Publish(events.Event{Type: events.ServiceError, Payload: map[string]string{"service": id, "error": err.Error()}})
	}
	return err
}

func (a *Application) StopService(id string) error {
	err := a.ServiceManager.Stop(id)
	if err == nil {
		a.events.Publish(events.Event{Type: events.ServiceStopped, Payload: id})
	} else {
		a.events.Publish(events.Event{Type: events.ServiceError, Payload: map[string]string{"service": id, "error": err.Error()}})
	}
	return err
}

func (a *Application) RestartService(id string) error {
	err := a.ServiceManager.Restart(id)
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
	err := a.RuntimeManager.SwitchPHPVersion(version)
	if err == nil {
		a.events.Publish(events.Event{Type: events.PHPVersionChanged, Payload: version})
	}
	return err
}

func (a *Application) Events() *events.Bus {
	return a.events
}
