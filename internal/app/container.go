package app

type Container struct {
	App *Application
}

func NewContainer(app *Application) *Container {
	return &Container{App: app}
}
