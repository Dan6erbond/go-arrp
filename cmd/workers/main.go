package main

import (
	"openwt.com/go-arrp/internal/cadence"
	"openwt.com/go-arrp/internal/config"
	"openwt.com/go-arrp/internal/tasks"
)

func main() {
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadence.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)
	cadence.StartWorkers(&cadenceClient, tasks.TaskListName)
	select {}
}
