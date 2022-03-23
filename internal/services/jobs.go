package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	"openwt.com/go-arrp/internal/cadence"
	"openwt.com/go-arrp/internal/dto"
	"openwt.com/go-arrp/internal/tasks"
)

type JobsService struct {
	cadenceAdapter *cadence.CadenceAdapter
	logger         *zap.Logger
}

func NewJobsService(cadenceAdapter *cadence.CadenceAdapter, logger *zap.Logger) *JobsService {
	return &JobsService{
		cadenceAdapter: cadenceAdapter,
		logger:         logger,
	}
}

func (s *JobsService) TriggerHelloWorld(name string) (*dto.TriggerHelloWorldResponse, error) {
	wo := client.StartWorkflowOptions{
		TaskList:                     tasks.TaskListName,
		ExecutionStartToCloseTimeout: time.Hour * 24,
	}
	execution, err := s.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, tasks.Workflow, name)
	if err != nil {
		s.logger.Error("Failed to start workflow", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Started workflow", zap.String("Workflow ID", execution.ID), zap.String("RunId", execution.RunID))
	return &dto.TriggerHelloWorldResponse{
		WorkflowID: execution.ID,
		Message:    "Started workflow.",
	}, nil
}

func (s *JobsService) SignalHelloWorld(workflowID string, age int) (*dto.SignalHelloWorldResponse, error) {
	err := s.cadenceAdapter.CadenceClient.SignalWorkflow(context.Background(), workflowID, "", tasks.SignalName, age)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Signaled work flow with the following params!", zap.String("WorkflowId", workflowID), zap.Int("Age", age))

	return &dto.SignalHelloWorldResponse{
		Message: fmt.Sprintf("Signaled workflow with age as %d!", age),
		Success: true,
	}, nil
}

func (s *JobsService) GetHelloWorldStatus(workflowID string) (*dto.HelloWorldStatusResponse, error) {
	var status string
	resp, err := s.cadenceAdapter.CadenceClient.QueryWorkflow(context.Background(), workflowID, "", "status")
	if err != nil {
		return nil, err
	}

	err = resp.Get(&status)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Got workflow status!", zap.String("WorkflowId", workflowID), zap.String("Status", status))

	var result string
	if status == "completed" {
		workflow := s.cadenceAdapter.CadenceClient.GetWorkflow(context.Background(), workflowID, "")
		if workflow != nil {
			err = workflow.Get(context.Background(), &result)
		}
		if err != nil {
			s.logger.Error("Failed to get workflow result!", zap.Error(err))
		}
	}

	return &dto.HelloWorldStatusResponse{
		Message: "Queried status.",
		Status:  status,
		Result:  result,
	}, nil
}
