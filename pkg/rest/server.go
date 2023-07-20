package rest

import (
	"context"

	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/a4c"
	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/vault"
	"github.com/eflows4hpc/hpcwaas-api/pkg/store"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Config *Config

	router       *gin.Engine
	a4cManager   a4c.Manager
	vaultManager vault.Manager
	store        store.Store
}

func (s *Server) StartServer() error {
	var err error

	s.a4cManager, err = a4c.GetManager(s.Config.AlienConfig)
	if err != nil {
		return err
	}
	s.vaultManager, err = vault.GetManager(context.Background(), s.Config.VaultConfig)
	if err != nil {
		return err
	}
	defer vault.CloseRenewers()
	s.store = store.NewStore(sessionDuration)

	s.router = gin.Default()
	s.setupRoutes()

	return s.router.Run(s.Config.ListenAddress)
}
