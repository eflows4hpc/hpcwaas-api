package api

import (
	"context"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

type ExecutionsService interface {
	Status(ctx context.Context, executionID string) (Execution, error)
	Cancel(ctx context.Context, executionID string) error
}

type executionsService struct {
	client *client
}

func (s *executionsService) Status(ctx context.Context, executionID string) (Execution, error) {
	var res Execution
	request, err := s.client.NewRequest(ctx, http.MethodGet, path.Join("/executions", executionID), nil)
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

func (s *executionsService) Cancel(ctx context.Context, executionID string) error {

	request, err := s.client.NewRequest(ctx, http.MethodDelete, path.Join("/executions", executionID), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create http request")
	}
	request.Header.Add("Accept", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to send http request to list workflows")
	}
	return ReadResponse(response, nil)

}
