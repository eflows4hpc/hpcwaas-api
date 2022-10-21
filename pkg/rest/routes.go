package rest

import "github.com/gin-gonic/gin"

func (s *Server) setupRoutes() {

	rootGrp := s.router.Group("/")
	s.setupAuth(rootGrp)
	{
		rootGrp.GET("/workflows", s.getWorkflows)
		rootGrp.POST("/workflows/:workflow_name", s.triggerWorkflow)
		rootGrp.GET("/executions/:execution_id", s.getExecution)
		rootGrp.GET("/executions/:execution_id/logs", s.getExecutionLogs)
		rootGrp.DELETE("/executions/:execution_id", s.cancelExecution)
		rootGrp.POST("/ssh_keys", s.createKey)
	}
}

func (s *Server) setupAuth(rootGrp *gin.RouterGroup) {
	if s.Config.Auth.BasicAuth != nil {
		rootGrp.Use(basicAuth(s.Config.Auth.BasicAuth))
	}
}
