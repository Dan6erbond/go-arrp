package internal_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"openwt.com/go-arrp/internal"
	"openwt.com/go-arrp/internal/config"
	"openwt.com/go-arrp/internal/dto"
)

type ServerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *ServerTestSuite) SetupTest() {
	suite.router = internal.MakeServer(&config.AppConfig{
		Env: "test",
		Cadence: config.CadenceConfig{
			Domain:   "go-arrp",
			Service:  "cadence-frontend",
			HostPort: "127.0.0.1:7933",
		},
	})
}

func (suite *ServerTestSuite) TestHelloWorldTriggerRoute() {
	w := httptest.NewRecorder()
	data := dto.TriggerHelloWorldRequest{
		Name: "John Doe",
	}
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/v1/jobs/hello-world", bytes.NewBuffer(body))
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusAccepted, w.Code)

	var resp dto.TriggerHelloWorldResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal("Started workflow.", resp.Message)
}

func (suite *ServerTestSuite) TestHelloWorldStatusRouteAfterStart() {
	w := httptest.NewRecorder()
	data := dto.TriggerHelloWorldRequest{
		Name: "John Doe",
	}
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/v1/jobs/hello-world", bytes.NewBuffer(body))
	suite.router.ServeHTTP(w, req)
	var workflow dto.TriggerHelloWorldResponse
	json.Unmarshal(w.Body.Bytes(), &workflow)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/jobs/hello-world/%s/status", workflow.WorkflowID), nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var resp dto.HelloWorldStatusResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal("Queried status.", resp.Message)
	suite.Condition(func() bool {
		return resp.Status == "started" || resp.Status == "activity started" || resp.Status == "waiting signal"
	})
}

func (suite *ServerTestSuite) TestSendHelloWorldSignalRoute() {
	w := httptest.NewRecorder()
	triggerData := dto.TriggerHelloWorldRequest{
		Name: "John Doe",
	}
	triggerBody, _ := json.Marshal(triggerData)
	req, _ := http.NewRequest("POST", "/api/v1/jobs/hello-world", bytes.NewBuffer(triggerBody))
	suite.router.ServeHTTP(w, req)
	var workflow dto.TriggerHelloWorldResponse
	json.Unmarshal(w.Body.Bytes(), &workflow)

	w = httptest.NewRecorder()
	data := dto.SignalHelloWorldRequest{
		Age: 20,
	}
	body, _ := json.Marshal(data)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/v1/jobs/hello-world/%s/age", workflow.WorkflowID), bytes.NewBuffer(body))
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusAccepted, w.Code)

	var resp dto.SignalHelloWorldResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.True(resp.Success)
	suite.Equal(fmt.Sprintf("Signaled workflow with age as %d!", data.Age), resp.Message)
}

func (suite *ServerTestSuite) TestHelloWorldStatusRouteAfterSignal() {
	w := httptest.NewRecorder()
	data := dto.TriggerHelloWorldRequest{
		Name: "John Doe",
	}
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/api/v1/jobs/hello-world", bytes.NewBuffer(body))
	suite.router.ServeHTTP(w, req)
	var workflow dto.TriggerHelloWorldResponse
	json.Unmarshal(w.Body.Bytes(), &workflow)

	w = httptest.NewRecorder()
	ageData := dto.SignalHelloWorldRequest{
		Age: 20,
	}
	ageBody, _ := json.Marshal(ageData)
	req, _ = http.NewRequest("POST", fmt.Sprintf("/api/v1/jobs/hello-world/%s/age", workflow.WorkflowID), bytes.NewBuffer(ageBody))
	suite.router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/jobs/hello-world/%s/status", workflow.WorkflowID), nil)
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var resp dto.HelloWorldStatusResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal("Queried status.", resp.Message)
	suite.Equal("completed", resp.Status)
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
