package rest

import (
	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/a4c"
	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/vault"
)

const DefaultListenAddress = "0.0.0.0:9090"

type Config struct {
	ListenAddress string       `mapstructure:"listen_address"`
	AlienConfig   a4c.Config   `mapstructure:"alien_config"`
	VaultConfig   vault.Config `mapstructure:"vault_config"`
}
