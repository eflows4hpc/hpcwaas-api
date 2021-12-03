package api

import (
	"net/http"
)

// Configuration holds the api client configuration
type Configuration struct {
	APIURL        string `mapstructure:"api_url"`
	SSLEnabled    bool   `mapstructure:"ssl_enabled"`
	SkipTLSVerify bool   `mapstructure:"skip_tls_verify"`
	KeyFile       string `mapstructure:"key_file"`
	CertFile      string `mapstructure:"cert_file"`
	CAFile        string `mapstructure:"ca_file"`
	CAPath        string `mapstructure:"ca_path"`

	HttpClient *http.Client
}
