package rest

import (
	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/a4c"
	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/vault"
	"golang.org/x/oauth2"
)

const DefaultListenAddress = "0.0.0.0:9090"

type Config struct {
	ListenAddress string       `mapstructure:"listen_address"`
	AlienConfig   a4c.Config   `mapstructure:"alien_config"`
	VaultConfig   vault.Config `mapstructure:"vault_config"`
	Auth          AuthConfig   `mapstructure:"auth"`
}

type AuthConfig struct {
	AuthType  string           `mapstructure:"auth_type,omitempty"`
	BasicAuth *BasicAuthConfig `mapstructure:"basic_auth,omitempty"`
	OAuth     *oauth2.Config
}

type BasicAuthConfig struct {
	Accounts []AuthAccount `mapstructure:"accounts"`
}

type AuthAccount struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
