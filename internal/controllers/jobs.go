package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openwt.com/go-arrp/internal/dto"
	"openwt.com/go-arrp/internal/services"
)

type JobsController struct {
	jobsService *services.JobsService
}

func NewJobsController(jobsService *services.JobsService) *JobsController {
	return &JobsController{jobsService: jobsService}
}

func (c *JobsController) TriggerHelloWorld(ctx *gin.Context) {
	var request dto.TriggerHelloWorldRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.jobsService.TriggerHelloWorld(request.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}

func (c *JobsController) SignalHelloWorld(ctx *gin.Context) {
	workflowID := ctx.Param("workflowID")
	var request dto.SignalHelloWorldRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.jobsService.SignalHelloWorld(workflowID, request.Age)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

func (c *JobsController) GetHelloWorldStatus(ctx *gin.Context) {
	workflowID := ctx.Param("workflowID")
	response, err := c.jobsService.GetHelloWorldStatus(workflowID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
