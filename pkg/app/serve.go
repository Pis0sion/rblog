package app

import "log"

type App struct {
	Name        string
	BaseName    string
	Description string
	Version     string
	RunFunc     RunFunc
}

type Options func(app *App)

type RunFunc func(string) error

func WithDescription(description string) Options {
	return func(app *App) {
		app.Description = description
	}
}

func WithVersion(version string) Options {
	return func(app *App) {
		app.Version = version
	}
}

func WithRunFunc(run RunFunc) Options {
	return func(app *App) {
		app.RunFunc = run
	}
}

func NewApplication(name, basename string, opts ...Options) *App {

	a := &App{
		Name:     name,
		BaseName: basename,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *App) Run() {

	if a.RunFunc != nil {
		log.Fatalln(a.RunFunc(a.Name))
	}

	return
}
