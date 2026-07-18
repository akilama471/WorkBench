package core

import (
	"fmt"

	svcmgr "github.com/akilama471/WorkBench/internal/service"
)

type ServiceManager struct {
	manager *svcmgr.Manager
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		manager: svcmgr.NewManager(),
	}
}

func (sm *ServiceManager) Register(svc svcmgr.Service) error {
	return sm.manager.Register(svc)
}

func (sm *ServiceManager) Start(id string) error {
	return sm.manager.Start(id)
}

func (sm *ServiceManager) Stop(id string) error {
	return sm.manager.Stop(id)
}

func (sm *ServiceManager) Restart(id string) error {
	return sm.manager.Restart(id)
}

func (sm *ServiceManager) Status(id string) (svcmgr.Status, error) {
	return sm.manager.Status(id)
}

func (sm *ServiceManager) GetService(id string) (svcmgr.Service, error) {
	return sm.manager.Get(id)
}

func (sm *ServiceManager) AllServices() []svcmgr.Service {
	return sm.manager.All()
}

func (sm *ServiceManager) StatusSummary() map[string]string {
	summary := make(map[string]string)
	for _, svc := range sm.manager.All() {
		summary[svc.ID()] = fmt.Sprintf("%s", svc.Status())
	}
	return summary
}
