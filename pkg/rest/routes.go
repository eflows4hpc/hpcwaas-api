package rest

func (s *Server) setupRoutes() {
	s.router.GET("/workflows", s.getWorkflows)
	s.router.POST("/workflows/:workflow_name", s.triggerWorkflow)
	s.router.GET("/executions/:execution_id", s.getExecution)
	s.router.DELETE("/executions/:execution_id", s.cancelExecution)
	s.router.POST("/users/:user_name/ssh_key", s.createKey)
}
