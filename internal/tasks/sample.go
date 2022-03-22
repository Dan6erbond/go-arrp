package tasks

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	TaskListName = "helloWorldGroup"
	SignalName   = "helloWorldSignal"
)

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(Workflow)
	activity.Register(helloworldActivity)
}

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
}

func helloworldActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity started")
	return "Hello " + name + "! How old are you!", nil
}

func Workflow(ctx workflow.Context, name string) (string, error) {
	currentState := "started"
	err := workflow.SetQueryHandler(ctx, "status", func() (string, error) {
		return currentState, nil
	})
	if err != nil {
		currentState = "failed to register query handler"
		return "", err
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("helloworld workflow started")

	currentState = "activity started"
	var activityResult string
	err = workflow.ExecuteActivity(ctx, helloworldActivity, name).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return "", err
	}

	// After saying hello, the workflow will wait for you to inform it of your age!
	currentState = "waiting signal"
	signalName := SignalName
	selector := workflow.NewSelector(ctx)
	var ageResult int

	for {
		signalChan := workflow.GetSignalChannel(ctx, signalName)
		selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
			c.Receive(ctx, &ageResult)
			workflow.GetLogger(ctx).Info("Received age results from signal!", zap.String("signal", signalName), zap.Int("value", ageResult))
		})
		workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
		// Wait for signal
		selector.Select(ctx)

		// We can check the age and return an appropriate response
		if ageResult > 0 && ageResult < 150 {
			logger.Info("Workflow completed.", zap.String("NameResult", activityResult), zap.Int("AgeResult", ageResult))

			currentState = "completed"
			return fmt.Sprintf("Hello "+name+"! You are %v years old!", ageResult), nil
		} else {
			currentState = "completed"
			return "You can't be that old!", nil
		}
	}
}
