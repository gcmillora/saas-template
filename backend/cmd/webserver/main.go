package main

import (
	"os"
	"saas-template/config"
	"saas-template/internal/webserver"
)

func main() {
	command := ""
	if len(os.Args[1:]) > 0 {
		command = os.Args[1]
	}

	app := config.NewApp()
	ws := webserver.NewWebserver(app)

	if command == "routes:list" {
		ws.PrintRoutes()
	} else {
		ws.Start()
	}
}
