package api

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type LogsRequestOpts struct {
	FromIndex *int
	Size      *int
	Timeout   time.Duration
	Levels    LogLevel
}

type ExecutionsService interface {
	Status(ctx context.Context, executionID string) (Execution, error)
	Logs(ctx context.Context, executionID string, opts *LogsRequestOpts) (ExecutionLogs, error)
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
		return res, errors.Wrap(err, "failed to send http request to get execution status")
	}

	err = ReadResponse(response, &res)
	return res, err
}

func (s *executionsService) Logs(ctx context.Context, executionID string, opts *LogsRequestOpts) (ExecutionLogs, error) {
	if opts == nil {
		opts = &LogsRequestOpts{}
	}
	var res ExecutionLogs
	u, err := url.Parse(path.Join("/executions", executionID, "logs"))
	if err != nil {
		return ExecutionLogs{}, err
	}
	v := u.Query()
	if opts.FromIndex != nil {
		v.Add("from", strconv.Itoa(*opts.FromIndex))
	}
	if opts.Timeout != 0 {
		v.Add("timeout", opts.Timeout.String())
	}
	if opts.Size != nil {
		v.Add("size", strconv.Itoa(*opts.Size))
	}
	if opts.Levels != 0 {
		v.Add("levels", strconv.Itoa(int(opts.Levels)))
	}
	u.RawQuery = v.Encode()

	request, err := s.client.NewRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return res, errors.Wrap(err, "failed to create http request")
	}
	request.Header.Add("Accept", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return res, errors.Wrap(err, "failed to send http request to get execution logs")
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
		return errors.Wrap(err, "failed to send http request to cancel execution")
	}
	return ReadResponse(response, nil)

}
