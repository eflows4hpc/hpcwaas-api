package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alien4cloud/alien4cloud-go-client/v3/alien4cloud"
	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/gin-gonic/gin"
)

func (s *Server) getExecution(gc *gin.Context) {
	executionID := gc.Param("execution_id")
	execution, err := s.a4cManager.GetExecution(gc.Request.Context(), executionID)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}
	execResponse := &api.Execution{
		ID:     execution.ID,
		Status: execution.Status,
	}

	gc.JSON(http.StatusOK, execResponse)
}
func getLevels(levelsInt api.LogLevel) []string {
	if levelsInt == 0 {
		return []string{"INFO", "WARN", "ERROR"}
	}
	var r []string
	if api.HasLogLevel(levelsInt, api.DEBUG) {
		r = append(r, "DEBUG")
	}
	if api.HasLogLevel(levelsInt, api.INFO) {
		r = append(r, "INFO")
	}
	if api.HasLogLevel(levelsInt, api.WARN) {
		r = append(r, "WARN")
	}
	if api.HasLogLevel(levelsInt, api.ERROR) {
		r = append(r, "ERROR")
	}

	return r
}
func (s *Server) getExecutionLogs(gc *gin.Context) {
	executionID := gc.Param("execution_id")
	fromStr := gc.DefaultQuery("from", "0")
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		writeError(gc, newBadRequestError(fmt.Errorf("invalid value for 'from' query parameter %q: %w", fromStr, err)))
		return
	}
	sizeStr := gc.DefaultQuery("size", "-1")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		writeError(gc, newBadRequestError(fmt.Errorf("invalid value for 'from' query parameter %q: %w", sizeStr, err)))
		return
	}
	timeoutStr := gc.DefaultQuery("timeout", "1m")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		writeError(gc, newBadRequestError(fmt.Errorf("invalid value for 'timeout' query parameter %q (golang duration expected): %w", timeoutStr, err)))
		return
	}

	levelsStr := gc.DefaultQuery("levels", "0")
	levelsUInt, err := strconv.ParseUint(levelsStr, 10, 8)
	if err != nil {
		writeError(gc, newBadRequestError(fmt.Errorf("invalid value for 'levels' query parameter %q: %w", levelsStr, err)))
		return
	}

	levels := getLevels(api.LogLevel(levelsUInt))

	ctx, cancel := context.WithTimeout(gc.Request.Context(), timeout)
	defer cancel()
	var executionLogs []alien4cloud.Log
	var totalLogs int

	for {
		executionLogs, totalLogs, err = s.a4cManager.GetExecutionLogs(gc.Request.Context(), executionID, from, size, levels...)
		if err != nil {
			writeError(gc, newInternalServerError(err))
			return
		}
		if len(executionLogs) > 0 {
			break
		}
		select {
		case <-ctx.Done():
			gc.JSON(http.StatusOK, &api.ExecutionLogs{
				Logs:         nil,
				TotalResults: totalLogs,
				From:         from,
			})
			return
		case <-time.After(3 * time.Second):
		}

	}

	execLogsResponse := &api.ExecutionLogs{
		Logs:         make([]api.Log, len(executionLogs)),
		TotalResults: totalLogs,
		From:         from,
	}

	for i := range executionLogs {
		execLogsResponse.Logs[i] = api.Log{
			Level:     executionLogs[i].Level,
			Timestamp: executionLogs[i].Timestamp.Time,
			Content:   executionLogs[i].Content,
		}
	}

	gc.JSON(http.StatusOK, execLogsResponse)
}

func (s *Server) cancelExecution(gc *gin.Context) {
	executionID := gc.Param("execution_id")
	err := s.a4cManager.CancelExecution(gc.Request.Context(), executionID)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}
	gc.Status(http.StatusAccepted)
}
