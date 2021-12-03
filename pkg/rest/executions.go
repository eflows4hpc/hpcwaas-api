package rest

import (
	"net/http"

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

func (s *Server) cancelExecution(gc *gin.Context) {
	executionID := gc.Param("execution_id")
	err := s.a4cManager.CancelExecution(gc.Request.Context(), executionID)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}
	gc.Status(http.StatusAccepted)
}
