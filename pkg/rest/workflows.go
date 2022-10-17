package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/eflows4hpc/hpcwaas-api/pkg/ctxauth"
	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/a4c"
)

func (s *Server) getWorkflows(gc *gin.Context) {
	ctx := gc.Request.Context()
	if auth, ok := gc.Get(gin.AuthUserKey); ok {
		log.Printf("authenticated user %+v", auth)
		ctx = ctxauth.WithCurrentUser(ctx, auth.(AuthAccount).Username)
	}

	workflows, err := s.a4cManager.GetWorkflows(ctx)
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

	ctx := gc.Request.Context()
	if auth, ok := gc.Get(gin.AuthUserKey); ok {
		log.Printf("authenticated user %+v", auth)
		ctx = ctxauth.WithCurrentUser(ctx, auth.(AuthAccount).Username)
	}

	execID, err := s.a4cManager.TriggerWorkflow(ctx, wfName, inputsReq.Inputs)
	if err != nil {
		if a4c.IsUnauthorizedError(err) {
			writeError(gc, newForbiddenRequest(err.Error()))
			return
		}
		writeError(gc, newInternalServerError(err))
		return
	}
	gc.Header("Location", fmt.Sprintf("/executions/%s", execID))
	gc.Status(http.StatusCreated)
}
