//go:build cgo

package pages

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/akilama471/WorkBench/internal/service"
	"github.com/akilama471/WorkBench/ui/components"

	apppkg "github.com/akilama471/WorkBench/internal/app"
)

func BuildServicesPage(app *apppkg.Application) fyne.CanvasObject {
	title := widget.NewLabel("Services")
	title.TextStyle = fyne.TextStyle{Bold: true}

	services := []struct {
		id   string
		name string
	}{
		{"apache", "Apache"},
		{"mariadb", "MariaDB"},
	}

	var cards []fyne.CanvasObject
	cards = append(cards, title)

	for _, svc := range services {
		svcID := svc.id
		svcName := svc.name

		statusLabel := widget.NewLabel(fmt.Sprintf("%s: Unknown", svcName))

		startBtn := widget.NewButton("Start", func() {
			_ = app.StartService(svcID)
			refreshServiceStatus(app, svcID, statusLabel)
		})
		stopBtn := widget.NewButton("Stop", func() {
			_ = app.StopService(svcID)
			refreshServiceStatus(app, svcID, statusLabel)
		})
		restartBtn := widget.NewButton("Restart", func() {
			_ = app.RestartService(svcID)
			refreshServiceStatus(app, svcID, statusLabel)
		})

		card := components.ServiceCard(svcName, service.StatusUnknown, nil, nil, nil)
		_ = card

		svcStatus, err := app.ServiceStatus(svcID)
		if err != nil {
			svcStatus = "Unknown"
		}
		statusLabel.SetText(fmt.Sprintf("%s: %s", svcName, svcStatus))

		cards = append(cards, statusLabel, container.NewHBox(startBtn, stopBtn, restartBtn))
	}

	return container.NewVBox(cards...)
}

func refreshServiceStatus(app *apppkg.Application, svcID string, label *widget.Label) {
	svcStatus, err := app.ServiceStatus(svcID)
	if err != nil {
		svcStatus = "Unknown"
	}
	svc, err := app.ServiceManager.GetService(svcID)
	name := svcID
	if err == nil {
		name = svc.Name()
	}
	label.SetText(fmt.Sprintf("%s: %s", name, svcStatus))
}
