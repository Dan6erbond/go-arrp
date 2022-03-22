package internal

import (
	"github.com/gin-gonic/gin"
	"openwt.com/go-arrp/internal/cadence"
	"openwt.com/go-arrp/internal/config"
	"openwt.com/go-arrp/internal/controllers"
	"openwt.com/go-arrp/internal/services"
)

func MakeServer(c ...*config.AppConfig) *gin.Engine {
	var appConfig config.AppConfig
	appConfig.Setup(c...)
	var cadenceClient cadence.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		jobsService := services.NewJobsService(&cadenceClient, appConfig.Logger)
		jobsController := controllers.NewJobsController(jobsService)
		v1.POST("/jobs/hello-world", jobsController.TriggerHelloWorld)
		v1.POST("/jobs/hello-world/:workflowID/age", jobsController.SignalHelloWorld)
		v1.GET("/jobs/hello-world/:workflowID/status", jobsController.GetHelloWorldStatus)
	}

	return r
}
