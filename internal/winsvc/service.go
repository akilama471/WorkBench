package winsvc

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/akilama471/WorkBench/internal/app"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type workbenchService struct {
	app      *app.Application
	services []string
	delayMs  int
}

func (m *workbenchService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	// Wait for delay
	if m.delayMs > 0 {
		time.Sleep(time.Duration(m.delayMs) * time.Millisecond)
	}

	for _, s := range m.services {
		if s != "" {
			m.app.StartService(s)
		}
	}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			for _, s := range m.services {
				if s != "" {
					m.app.StopService(s)
				}
			}
			break loop
		default:
			log.Printf("unexpected control request #%d", c)
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func RunAsService(application *app.Application, services []string, delayMs int) error {
	return svc.Run("WorkBenchAutoStart", &workbenchService{
		app:      application,
		services: services,
		delayMs:  delayMs,
	})
}

// ElevateAndInstall spawns a new elevated process to install the service
func ElevateAndInstall(delayMs int, services []string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	args := fmt.Sprintf("service install %d %s", delayMs, strings.Join(services, " "))

	verb := "runas"
	cwd, _ := os.Getwd()

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exePath)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, windows.SW_NORMAL)
	if err != nil {
		return fmt.Errorf("failed to prompt for elevation: %v", err)
	}
	return nil
}

// ElevateAndRemove spawns a new elevated process to remove the service
func ElevateAndRemove() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	verb := "runas"
	cwd, _ := os.Getwd()

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exePath)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString("service remove")

	err = windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, windows.SW_NORMAL)
	if err != nil {
		return fmt.Errorf("failed to prompt for elevation: %v", err)
	}
	return nil
}

func InstallService(delayMs int, services []string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService("WorkBenchAutoStart")
	if err == nil {
		s.Close()
		// Remove existing to recreate
		RemoveService()
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	args := []string{"service", "run", fmt.Sprintf("%d", delayMs)}
	args = append(args, services...)

	s, err = m.CreateService("WorkBenchAutoStart", exePath, mgr.Config{
		DisplayName: "WorkBench AutoStart Service",
		Description: "Automatically starts selected WorkBench services",
		StartType:   mgr.StartAutomatic,
	}, args...)
	if err != nil {
		return err
	}
	defer s.Close()

	return nil
}

func RemoveService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService("WorkBenchAutoStart")
	if err != nil {
		return fmt.Errorf("service not found")
	}
	defer s.Close()

	// Stop it if running
	st, _ := s.Query()
	if st.State != svc.Stopped {
		s.Control(svc.Stop)
	}

	return s.Delete()
}

func IsServiceInstalled() bool {
	m, err := mgr.Connect()
	if err != nil {
		return false
	}
	defer m.Disconnect()

	s, err := m.OpenService("WorkBenchAutoStart")
	if err != nil {
		return false
	}
	s.Close()
	return true
}

func GetServiceConfig() (int, []string, error) {
	m, err := mgr.Connect()
	if err != nil {
		return 0, nil, err
	}
	defer m.Disconnect()

	s, err := m.OpenService("WorkBenchAutoStart")
	if err != nil {
		return 0, nil, err
	}
	defer s.Close()

	cfg, err := s.Config()
	if err != nil {
		return 0, nil, err
	}

	// BinaryPathName contains the exe and args
	// e.g. C:\path\to\workbench.exe service run 10000 apache mariadb
	// We can parse it by getting the arguments after "run"
	parts := strings.Split(cfg.BinaryPathName, " ")

	runIdx := -1
	for i, p := range parts {
		if p == "run" {
			runIdx = i
			break
		}
	}

	if runIdx == -1 || runIdx+1 >= len(parts) {
		return 0, nil, fmt.Errorf("invalid config")
	}

	delayMs := 0
	fmt.Sscanf(parts[runIdx+1], "%d", &delayMs)

	var services []string
	if runIdx+2 < len(parts) {
		services = parts[runIdx+2:]
	}

	return delayMs, services, nil
}
