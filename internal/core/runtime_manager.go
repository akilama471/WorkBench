package core

import (
	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
	"github.com/akilama471/WorkBench/internal/runtime"
	"github.com/akilama471/WorkBench/internal/runtime/php"
)

type RuntimeManager struct {
	paths   *filesystem.Paths
	log     *logger.Logger
	manager *runtime.Manager
	php     *php.Manager
}

func NewRuntimeManager(paths *filesystem.Paths, log *logger.Logger) *RuntimeManager {
	rm := runtime.NewManager()
	phpMgr := php.NewManager(paths, log)

	return &RuntimeManager{
		paths:   paths,
		log:     log,
		manager: rm,
		php:     phpMgr,
	}
}

func (rm *RuntimeManager) Register(r runtime.Runtime) error {
	return rm.manager.Register(r)
}

func (rm *RuntimeManager) Get(id string) (runtime.Runtime, error) {
	return rm.manager.Get(id)
}

func (rm *RuntimeManager) All() []runtime.Runtime {
	return rm.manager.All()
}

func (rm *RuntimeManager) PHP() *php.Manager {
	return rm.php
}

func (rm *RuntimeManager) ListPHPVersions() ([]runtime.Version, error) {
	return rm.php.ListVersions()
}

func (rm *RuntimeManager) CurrentPHPVersion() (string, error) {
	return rm.php.CurrentVersion()
}

func (rm *RuntimeManager) SwitchPHPVersion(version string) error {
	return rm.php.SwitchVersion(version)
}
