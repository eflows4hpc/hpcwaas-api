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

// SSHKey is the response structure of a CreateKey operation for a given user
//
// The response contains only the public key, private key is never disclosed
type SSHKey struct {
	PublicKey string `json:"public_key"`
}
