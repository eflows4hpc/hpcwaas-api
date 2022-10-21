package a4c

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/alien4cloud/alien4cloud-go-client/v3/alien4cloud"
	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/eflows4hpc/hpcwaas-api/pkg/ctxauth"
)

type Manager interface {
	GetWorkflows(ctx context.Context) ([]api.Workflow, error)
	TriggerWorkflow(ctx context.Context, id string, inputs map[string]interface{}) (string, error)

	GetExecution(ctx context.Context, id string) (alien4cloud.Execution, error)
	GetExecutionLogs(ctx context.Context, id string, fromIndex, size int, levels ...string) ([]alien4cloud.Log, int, error)
	CancelExecution(ctx context.Context, id string) error
}

var unauthorizedError = errors.New("not authorized")

func IsUnauthorizedError(e error) bool {
	return errors.Is(e, unauthorizedError)
}

const DefaultAddress = "http://127.0.0.1:8088"
const DefaultUser = "admin"
const DefaultPassword = "admin"

type Config struct {
	Address    string `mapstructure:"address"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	CaFile     string `mapstructure:"ca_file"`
	SkipSecure bool   `mapstructure:"skip_secure"`
}

type manager struct {
	client  alien4cloud.Client
	baseURL string
}

var l sync.Mutex
var managerCache map[Config]Manager = make(map[Config]Manager)

func GetManager(c Config) (
	Manager, error) {
	l.Lock()
	defer l.Unlock()
	if m, ok := managerCache[c]; ok {
		return m, nil
	}

	m := &manager{}
	m.baseURL = c.Address
	var err error
	m.client, err = alien4cloud.NewClient(c.Address, c.User, c.Password, c.CaFile, c.SkipSecure)
	if err != nil {
		return nil, err
	}
	managerCache[c] = m

	return m, nil
}

func checkUserInAuthorizedTag(currentUsername string, tagValue string) bool {
	log.Printf("Checking if %q is in authorized users %q", currentUsername, tagValue)
	for _, user := range strings.Split(tagValue, ",") {
		if user == currentUsername {
			return true
		}
	}
	return false
}

func (m *manager) isUserInAuthorizedTag(ctx context.Context, currentUsername, applicationID string) (bool, error) {
	application, err := m.client.ApplicationService().GetApplicationByID(ctx, applicationID)
	if err != nil {
		return false, err
	}
	for _, tag := range application.Tags {
		if tag.Key == "hpcwaas-authorized-users" {
			return checkUserInAuthorizedTag(currentUsername, tag.Value), nil
		}
	}
	// If hpcwaas-authorized-users not in tags every user is authorized
	return true, nil

}

func (m *manager) GetWorkflows(ctx context.Context) ([]api.Workflow, error) {
	appsSearchReq := alien4cloud.SearchRequest{
		Size: 1000000,
		Filters: map[string][]string{
			"tags.name": {"hpcwaas-workflows"},
		}}
	// TODO(loicalbertin): add pagination
	apps, _, err := m.client.ApplicationService().SearchApplications(ctx, appsSearchReq)
	if err != nil {
		return nil, err
	}

	// If no auth and no tags on app this is ok
	// if no auth and tag then request is not authorized
	currentUsername, _ := ctxauth.GetCurrentUser(ctx)

	var ids []api.Workflow
	for _, app := range apps {
		var declaredWF []string
		// by default if not specified all users are authorized
		userAuthorized := true
		for _, tag := range app.Tags {
			if tag.Key == "hpcwaas-workflows" {
				declaredWF = strings.Split(tag.Value, ",")
			}
			if tag.Key == "hpcwaas-authorized-users" {
				userAuthorized = checkUserInAuthorizedTag(currentUsername, tag.Value)
			}
		}
		if !userAuthorized {
			continue
		}
		envs, _, err := m.client.ApplicationService().SearchEnvironments(ctx, app.ID, alien4cloud.SearchRequest{
			Size: 100000,
			// Filters seems to not apply here
			// Filters: map[string][]string{"status": {"DEPLOYED"}},
		})
		if err != nil {
			return nil, err
		}
		for _, env := range envs {
			if env.Status != "DEPLOYED" {
				continue
			}
			for _, wf := range declaredWF {
				ids = append(ids, api.Workflow{
					ID:              strings.Join([]string{app.ID, env.ID, wf}, "@"),
					ApplicationID:   app.ID,
					EnvironmentID:   env.ID,
					EnvironmentName: env.Name,
					Name:            wf,
				})
			}
		}
	}
	return ids, nil
}

func (m *manager) TriggerWorkflow(ctx context.Context, id string, inputs map[string]interface{}) (string, error) {
	parts := strings.Split(id, "@")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid workflow id %q", id)
	}
	appID := parts[0]
	envID := parts[1]
	wfID := parts[2]

	// If no auth and no tags on app this is ok
	// if no auth and tag then request is not authorized
	currentUsername, _ := ctxauth.GetCurrentUser(ctx)
	authorized, err := m.isUserInAuthorizedTag(ctx, currentUsername, appID)
	if err != nil {
		return "", err
	}
	if !authorized {
		return "", fmt.Errorf("user %q is not authorized to trigger workflow %q: %w", currentUsername, id, unauthorizedError)
	}

	return m.client.DeploymentService().RunWorkflowAsyncWithParameters(ctx, appID, envID, wfID, inputs, func(exec *alien4cloud.Execution, err error) {})
}

func (m *manager) GetExecution(ctx context.Context, id string) (alien4cloud.Execution, error) {
	return m.client.DeploymentService().GetExecutionByID(ctx, id)
}

func (m *manager) GetExecutionLogs(ctx context.Context, id string, fromIndex, size int, levels ...string) ([]alien4cloud.Log, int, error) {
	filters := alien4cloud.LogFilter{
		ExecutionID: []string{id},
		Level:       levels,
	}
	return m.getLogsOfExecution(ctx, id, filters, fromIndex, size)
}

func (m *manager) CancelExecution(ctx context.Context, id string) error {

	execution, err := m.client.DeploymentService().GetExecutionByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get execution %q: %w", id, err)
	}

	deployment, err := m.client.DeploymentService().GetDeployment(ctx, execution.DeploymentID)
	if err != nil {
		return fmt.Errorf("failed to get deployment %q linked to execution %q: %w", execution.DeploymentID, id, err)
	}

	return m.client.DeploymentService().CancelExecution(ctx, deployment.EnvironmentID, id)
}
