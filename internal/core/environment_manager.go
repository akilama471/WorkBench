package core

import (
	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
)

type EnvironmentManager struct {
	paths       *filesystem.Paths
	initializer *filesystem.Initializer
	log         *logger.Logger
}

func NewEnvironmentManager(paths *filesystem.Paths, log *logger.Logger) *EnvironmentManager {
	return &EnvironmentManager{
		paths:       paths,
		initializer: filesystem.NewInitializer(paths),
		log:         log,
	}
}

func (em *EnvironmentManager) Initialize() error {
	em.log.Info(logger.CategoryApplication, "initializing environment", "root", em.paths.Root())
	return em.initializer.Initialize()
}

func (em *EnvironmentManager) Root() string {
	return em.paths.Root()
}

func (em *EnvironmentManager) Paths() *filesystem.Paths {
	return em.paths
}
