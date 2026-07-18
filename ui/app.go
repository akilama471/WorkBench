//go:build cgo

package ui

import (
	apppkg "github.com/akilama471/WorkBench/internal/app"
)

type App struct {
	application *apppkg.Application
}

func NewApp(application *apppkg.Application) *App {
	return &App{application: application}
}

func (a *App) Run() {
	a.application.InitializeEnvironment()
}
