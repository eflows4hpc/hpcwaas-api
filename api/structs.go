package api

import "time"

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

// Log represents the log entry return by the a4c rest api
type Log struct {
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

// ExecutionLogs is the response structure of a GetExecutionLog operation
type ExecutionLogs struct {
	Logs         []Log `json:"logs"`
	TotalResults int   `json:"total_results"`
	From         int   `json:"from"`
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

// LogLevel
type LogLevel uint8

const (
	DEBUG LogLevel = 1 << iota
	INFO
	WARN
	ERROR
)

func SetLogLevels(l LogLevel, levels ...LogLevel) LogLevel {
	for _, al := range levels {
		l |= al
	}
	return l
}

func ClearLogLevels(l LogLevel, levels ...LogLevel) LogLevel {
	for _, al := range levels {
		l &^= al
	}
	return l
}

func HasLogLevel(l, level LogLevel) bool {
	return l&level != 0
}
