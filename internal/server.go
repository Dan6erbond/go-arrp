package internal

import (
	"github.com/gin-gonic/gin"
	"openwt.com/go-arrp/internal/cadence"
	"openwt.com/go-arrp/internal/config"
	"openwt.com/go-arrp/internal/controllers"
	"openwt.com/go-arrp/internal/router"
	"openwt.com/go-arrp/internal/services"
)

func NewGin() *gin.Engine {
	r := gin.Default()
	return r
}

func MakeServer(c ...*config.AppConfig) *gin.Engine {
	appConfig := config.ProvideConfig(c...)
	var cadenceClient cadence.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	r := gin.Default()

	jobsService := services.NewJobsService(&cadenceClient, appConfig.Logger)
	jobsController := controllers.NewJobsController(jobsService)

	router.RegisterRoutes(r, jobsController)

	return r
}
