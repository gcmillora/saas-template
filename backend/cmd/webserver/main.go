package main

import (
	"adobo/config"
	"adobo/internal/webserver"
	"os"
)

func main() {
	command := ""
	if len(os.Args[1:]) > 0 {
		command = os.Args[1]
	}

	a:= config.NewApp()
	ws := webserver.NewWebserver(a)

	if command == "routes:list"{
		ws.PrintRoutes()
	} else {
		ws.Start()
	}
}