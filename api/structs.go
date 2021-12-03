package api

// Workflows is the response structure of a GetWorkflows operation
type Workflows struct {
	Workflows []string `json:"workflows"`
}

// WorkflowInputs is the payload send when triggering a workflow execution
type WorkflowInputs struct {
	Inputs map[string]interface{} `json:"inputs"`
}

// Execution is the response structure of a GetExecution operation
type Execution struct {
	ID      string                 `json:"id"`
	Status  string                 `json:"status"`
	Outputs map[string]interface{} `json:"outputs,omitempty"`
}
