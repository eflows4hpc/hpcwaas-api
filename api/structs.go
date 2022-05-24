package api

// Workflows is the response structure of a GetWorkflows operation
type Workflows struct {
	Workflows []Workflow `json:"workflows"`
}

type Workflow struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ApplicationID   string `json:"application_id"`
	EnvironmentID   string `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
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

// SSHKeyGenerationRequest is the request structure of a CreateKey operation
type SSHKeyGenerationRequest struct {
	MetaData map[string]interface{} `json:"metadata,omitempty"`
}

// SSHKey is the response structure of a CreateKey operation for a given user
//
// The response contains only the public key, private key is never disclosed
type SSHKey struct {
	ID        string `json:"id"`
	PublicKey string `json:"public_key"`
}
