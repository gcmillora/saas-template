package console

import "saas-template/config"

type Console struct {
	app *config.App
}

func NewConsole(app *config.App) *Console {
	return &Console{
		app: app,
	}
}

	