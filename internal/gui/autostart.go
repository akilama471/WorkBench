package gui

import (
	"fmt"

	"github.com/akilama471/WorkBench/internal/winsvc"
)

func EnableAutoStart(delaySeconds string, services []string) error {
	delayMs := 10000
	if delaySeconds != "" {
		delaySec := 10
		fmt.Sscanf(delaySeconds, "%d", &delaySec)
		delayMs = delaySec * 1000
	}

	return winsvc.ElevateAndInstall(delayMs, services)
}

func DisableAutoStart() error {
	return winsvc.ElevateAndRemove()
}

func IsAutoStartEnabled() bool {
	return winsvc.IsServiceInstalled()
}
