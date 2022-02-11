package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

type WorkflowsService interface {
	List(ctx context.Context) (Workflows, error)
	Trigger(ctx context.Context, workflowID string, inputs *WorkflowInputs) (string, error)
}

type workflowsService struct {
	client *client
}

func (s *workflowsService) List(ctx context.Context) (Workflows, error) {
	var res Workflows
	request, err := s.client.NewRequest(ctx, http.MethodGet, "/workflows", nil)
	if err != nil {
		return res, errors.Wrap(err, "failed to create http request")
	}
	request.Header.Add("Accept", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return res, errors.Wrap(err, "failed to send http request to list workflows")
	}

	err = ReadResponse(response, &res)
	return res, err
}

func (s *workflowsService) Trigger(ctx context.Context, workflowID string, inputs *WorkflowInputs) (string, error) {

	body, err := json.Marshal(inputs)
	if err != nil {
		return "", errors.Wrapf(err, "Fail to marshall inputs for workflow: %q", workflowID)
	}

	request, err := s.client.NewRequest(ctx, http.MethodPost, path.Join("/workflows", workflowID), bytes.NewReader(body))
	if err != nil {
		return "", errors.Wrap(err, "failed to create http request")
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "failed to send http request to list workflows")
	}
	err = ReadResponse(response, nil)
	if err != nil {
		return "", err
	}
	return path.Base(response.Header.Get("Location")), nil
}
