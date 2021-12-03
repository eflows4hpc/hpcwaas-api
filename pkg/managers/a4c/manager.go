package a4c

import (
	"context"
	"strings"
	"sync"

	"github.com/alien4cloud/alien4cloud-go-client/v3/alien4cloud"
	"github.com/pkg/errors"
)

type Manager interface {
	GetWorkflows(ctx context.Context) ([]string, error)
	TriggerWorkflow(ctx context.Context, id string, inputs map[string]interface{}) (string, error)

	GetExecution(ctx context.Context, id string) (alien4cloud.Execution, error)
	CancelExecution(ctx context.Context, id string) error
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
	client alien4cloud.Client
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
	var err error
	m.client, err = alien4cloud.NewClient(c.Address, c.User, c.Password, c.CaFile, c.SkipSecure)
	if err != nil {
		return nil, err
	}
	managerCache[c] = m

	return m, nil
}

func (m *manager) GetWorkflows(ctx context.Context) ([]string, error) {
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

	var ids []string
	for _, app := range apps {
		var declaredWF []string
		for _, tag := range app.Tags {
			if tag.Key == "hpcwaas-workflows" {
				declaredWF = strings.Split(tag.Value, ",")
			}
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
				ids = append(ids, strings.Join([]string{app.ID, env.ID, wf}, "@"))
			}
		}
	}
	return ids, nil
}

func (m *manager) TriggerWorkflow(ctx context.Context, id string, inputs map[string]interface{}) (string, error) {
	parts := strings.Split(id, "@")
	appID := parts[0]
	envID := parts[1]
	wfID := parts[2]

	return m.client.DeploymentService().RunWorkflowAsyncWithParameters(ctx, appID, envID, wfID, inputs, func(exec *alien4cloud.Execution, err error) {})
}

func (m *manager) GetExecution(ctx context.Context, id string) (alien4cloud.Execution, error) {
	return m.client.DeploymentService().GetExecutionByID(ctx, id)
}

func (m *manager) CancelExecution(ctx context.Context, id string) error {

	execution, err := m.client.DeploymentService().GetExecutionByID(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to get execution %q", id)
	}

	deployment, err := m.client.DeploymentService().GetDeployment(ctx, execution.DeploymentID)
	if err != nil {
		return errors.Wrapf(err, "failed to get deployment %q linked to execution %q", execution.DeploymentID, id)
	}

	return m.client.DeploymentService().CancelExecution(ctx, deployment.EnvironmentID, id)
}
