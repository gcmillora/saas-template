package console

import "adobo/config"

type Console struct {
	app *config.App
}

func NewConsole(app *config.App) *Console {
	return &Console{
		app: app,
	}
}

	