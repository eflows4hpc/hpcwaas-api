package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eflows4hpc/hpcwaas-api/api"
)

func (s *Server) getWorkflows(gc *gin.Context) {
	currentUsername := ""
	if auth, ok := gc.Get(gin.AuthUserKey); ok {
		currentUsername = auth.(AuthAccount).Username
	}

	workflows, err := s.a4cManager.GetWorkflows(gc.Request.Context(), currentUsername)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}

	gc.JSON(http.StatusOK, api.Workflows{Workflows: workflows})

}

func (s *Server) triggerWorkflow(gc *gin.Context) {
	inputsReq := new(api.WorkflowInputs)
	wfName := gc.Param("workflow_name")

	err := gc.ShouldBindJSON(inputsReq)
	if err != nil {
		writeError(gc, newBadRequestError(err))
		return
	}

	execID, err := s.a4cManager.TriggerWorkflow(gc.Request.Context(), wfName, inputsReq.Inputs)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}
	gc.Header("Location", fmt.Sprintf("/executions/%s", execID))
	gc.Status(http.StatusCreated)
}
