package router

import (
	"github.com/gin-gonic/gin"
	"openwt.com/go-arrp/internal/controllers"
)

func RegisterRoutes(r *gin.Engine, jobsController *controllers.JobsController) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/jobs/hello-world", jobsController.TriggerHelloWorld)
		v1.POST("/jobs/hello-world/:workflowID/age", jobsController.SignalHelloWorld)
		v1.GET("/jobs/hello-world/:workflowID/status", jobsController.GetHelloWorldStatus)
	}
}
