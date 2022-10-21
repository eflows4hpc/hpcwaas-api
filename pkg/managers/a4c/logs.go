package a4c

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/alien4cloud/alien4cloud-go-client/v3/alien4cloud"
)

// logsSearchRequest is the representation of a request to search logs of an application in the A4C catalog
type logsSearchRequest struct {
	From              int                   `json:"from"`
	Size              int                   `json:"size,omitempty"`
	Query             string                `json:"query,omitempty"`
	Filters           alien4cloud.LogFilter `json:"filters"`
	SortConfiguration struct {
		Ascending bool   `json:"ascending"`
		SortBy    string `json:"sortBy"`
	} `json:"sortConfiguration"`
}

func (m *manager) getTotalLogs(ctx context.Context, executionID string, filters alien4cloud.LogFilter) (int, error) {
	logsFilter := logsSearchRequest{
		From:    0,
		Size:    1,
		Filters: filters,
	}

	body, err := json.Marshal(logsFilter)

	if err != nil {
		return 0, fmt.Errorf("unable to marshal log filters in order to get the number of logs available for this deployment: %w", err)
	}

	request, err := m.client.NewRequest(ctx,
		"POST",
		"/rest/latest/deployment/logs/search",
		bytes.NewReader(body),
	)
	if err != nil {
		return 0, fmt.Errorf("cannot create a request to get number of logs for execution '%s': %w", executionID, err)
	}
	var res struct {
		Data struct {
			TotalResults int `json:"totalResults"`
		} `json:"data"`
	}

	response, err := m.client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("cannot send a request to get number of logs for execution '%s': %w", executionID, err)
	}
	err = alien4cloud.ReadA4CResponse(response, &res)
	if err != nil {
		return 0, fmt.Errorf("cannot get number of logs for execution '%s': %w", executionID, err)
	}
	return res.Data.TotalResults, nil

}

func (m *manager) getLogsOfExecution(ctx context.Context, executionID string, filters alien4cloud.LogFilter, fromIndex, size int) ([]alien4cloud.Log, int, error) {
	// The first step allow us to get the number of logs available. We will re-use the TotalResults parameters in order to generate the second request.
	totalResults, err := m.getTotalLogs(ctx, executionID, filters)
	if err != nil {
		return nil, totalResults, err
	}
	if size < 0 {
		size = totalResults
	}

	logsFilter := logsSearchRequest{
		From:    fromIndex,
		Size:    size,
		Filters: filters,
		SortConfiguration: struct {
			Ascending bool   `json:"ascending"`
			SortBy    string `json:"sortBy"`
		}{Ascending: true, SortBy: "timestamp"},
	}

	body, err := json.Marshal(logsFilter)
	if err != nil {
		return nil, totalResults, fmt.Errorf("unable to marshal log filters to get logs for the deployment: %w", err)
	}

	request, err := m.client.NewRequest(ctx,
		"POST",
		"/rest/latest/deployment/logs/search",
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, totalResults, fmt.Errorf("cannot create a request to get logs for execution '%s': %w", executionID, err)
	}
	response, err := m.client.Do(request)
	if err != nil {
		return nil, totalResults, fmt.Errorf("cannot send a request to get logs for execution '%s': %w", executionID, err)
	}
	var res struct {
		Data struct {
			Data []alien4cloud.Log `json:"data"`
		} `json:"data"`
	}
	err = alien4cloud.ReadA4CResponse(response, &res)
	if err != nil {
		return nil, totalResults, fmt.Errorf("cannot get logs for execution '%s': %w", executionID, err)
	}

	return res.Data.Data, totalResults, nil
}
