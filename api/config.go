package api

import (
	"net/http"
)

const DefaultAPIAddress = "https://127.0.0.1:9090"

// Configuration holds the api client configuration
type Configuration struct {
	APIURL        string `mapstructure:"api_url"`
	SkipTLSVerify bool   `mapstructure:"skip_tls_verify"`
	KeyFile       string `mapstructure:"key_file"`
	CertFile      string `mapstructure:"cert_file"`
	CAFile        string `mapstructure:"ca_file"`
	CAPath        string `mapstructure:"ca_path"`
	User          string `mapstructure:"user"`
	AccessToken   string `mapstructure:"access_token"`

	HttpClient *http.Client
}
