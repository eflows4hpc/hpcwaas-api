package vault

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

type Manager interface {
	StoreKV(path string, data map[string]interface{}) error
}

const DefaultAddress = "http://127.0.0.1:8200"

type Config struct {
	Address         string `mapstructure:"address"`
	RoleID          string `mapstructure:"role_id"`
	SecretID        string `mapstructure:"secret_id"`
	IsSecretWrapped bool   `mapstructure:"is_secret_wrapped"`
}

type manager struct {
	client *api.Client
}

var l sync.Mutex
var managerCache map[Config]managerAndRenewer = make(map[Config]managerAndRenewer)

type managerAndRenewer struct {
	manager Manager
	renewer *api.LifetimeWatcher
}

func CloseRenewers() {

	l.Lock()
	defer l.Unlock()

	for _, m := range managerCache {
		m.renewer.Stop()
	}
}

func GetManager(ctx context.Context, config Config) (Manager, error) {
	l.Lock()
	defer l.Unlock()
	if m, ok := managerCache[config]; ok {
		return m.manager, nil
	}

	m := &manager{}
	var err error

	vConfig := api.DefaultConfig()
	vConfig.Address = config.Address

	m.client, err = api.NewClient(vConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	secretID := &approle.SecretID{
		FromString: config.SecretID,
	}

	opts := []approle.LoginOption{}
	if config.IsSecretWrapped {
		log.Println("[INFO] Vault SecretID is wrapped")
		opts = append(opts, approle.WithWrappingToken())
	}

	appRoleAuth, err := approle.NewAppRoleAuth(config.RoleID, secretID, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize AppRole auth method: %w", err)
	}

	authInfo, err := m.client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	watcher, err := m.client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
		Secret:        authInfo,
		RenewBehavior: api.RenewBehaviorErrorOnErrors,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize token renewer: %w", err)
	}

	go watcher.Start()

	managerCache[config] = managerAndRenewer{manager: m, renewer: watcher}

	return m, nil
}

func (m *manager) StoreKV(path string, data map[string]interface{}) error {
	_, err := m.client.Logical().Write(path, data)
	return err

}
