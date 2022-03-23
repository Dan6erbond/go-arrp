package dto

type TriggerHelloWorldRequest struct {
	Name string `json:"name"`
}

type TriggerHelloWorldResponse struct {
	Message    string `json:"message"`
	WorkflowID string `json:"workflowId"`
}

type SignalHelloWorldRequest struct {
	Age int `json:"age"`
}

type SignalHelloWorldResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type HelloWorldStatusResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Result  string `json:"result"`
}
