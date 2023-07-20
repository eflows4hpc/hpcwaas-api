package rest

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes() {
	s.setupStore()
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

	authGrp := s.router.Group("/auth")
	{
		authGrp.GET("/login", s.login)
		authGrp.GET("/authorize", s.authorize)
		authGrp.GET("/logout", s.logout)
	}
}

func (s *Server) setupAuth(group *gin.RouterGroup) {
	auth := s.Config.Auth
	switch auth.AuthType {
	case "basic":
		log.Println("Using basic authentication")
		group.Use(basicAuth(auth.BasicAuth))
	case "sso":
		log.Println("Using SSO authentication")
		s.initSsoConf()
		group.Use(s.ssoAuth(s.Config.Auth.OAuth))
	default:
		log.Printf("Invalid authentication type*: '%s'", auth.AuthType)
	}
}

func (s *Server) setupStore() {
	store := cookie.NewStore(storeSecret)
	session := sessions.Sessions(sessionName, store)
	s.router.Use(session)
}
