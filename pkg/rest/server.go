package rest

import (
	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/a4c"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Config *Config

	router     *gin.Engine
	a4cManager a4c.Manager
}

func (s *Server) StartServer() error {
	var err error

	s.a4cManager, err = a4c.GetManager(s.Config.AlienConfig)
	if err != nil {
		return err
	}

	s.router = gin.Default()
	s.setupRoutes()

	return s.router.Run(s.Config.ListenAddress)
}
