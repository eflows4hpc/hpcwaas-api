package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/goware/urlx"
	"github.com/hashicorp/go-rootcerts"
	"github.com/pkg/errors"
)

// HTTPClient represents an HTTP client
type HTTPClient interface {
	NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
}

// GetClient returns a HTTP Client
func GetClient(cc Configuration) (HTTPClient, error) {
	apiURL := cc.APIURL
	apiURL = strings.TrimRight(apiURL, "/")
	caFile := cc.CAFile
	caPath := cc.CAPath
	certFile := cc.CertFile
	keyFile := cc.KeyFile

	httpClient := cc.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 1 * time.Minute,
		}
	}
	client := &client{
		baseURL: "http://" + apiURL,
		Client:  httpClient,
	}
	if cc.SSLEnabled || cc.CAFile != "" || cc.CAPath != "" || (certFile != "" && keyFile != "") {
		url, err := urlx.Parse(apiURL)
		if err != nil {
			return nil, errors.Wrap(err, "Malformed API URL")
		}
		apiHost, _, err := urlx.SplitHostPort(url)
		if err != nil {
			return nil, errors.Wrap(err, "Malformed API URL")
		}

		tlsConfig := &tls.Config{ServerName: apiHost}
		if certFile != "" && keyFile != "" {
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to load TLS certificates")
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
		if caFile != "" || caPath != "" {
			cfg := &rootcerts.Config{
				CAFile: caFile,
				CAPath: caPath,
			}
			rootcerts.ConfigureTLS(tlsConfig, cfg)
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		if cc.SkipTLSVerify {
			tlsConfig.InsecureSkipVerify = true
			fmt.Println("Warning : usage of skip_tls_verify is not recommended for production and may expose to MITM attack")
		}

		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client.baseURL = "https://" + apiURL
		httpClient.Transport = tr
		client.Client = httpClient

	}

	return client, nil
}

// client is the HTTP client structure
type client struct {
	*http.Client
	baseURL string
}

// NewRequest returns a new HTTP request
func (c *client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
}

// ReadResponse is an helper function that allow to fully read and close a response body and
// unmarshal its json content into a provided data structure.
// If response status code is greater or equal to 400 it automatically parse an error response and
// returns it as a non-nil error.
func ReadResponse(response *http.Response, data interface{}) error {
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "Cannot read the response from Prov")
	}
	if response.StatusCode >= 400 {
		return handleErrors(response.StatusCode, responseBody)
	}
	// Don't try to retrieve data if StatusNoContent or StatusNotModified
	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusNotModified && data != nil {
		err = json.Unmarshal(responseBody, &data)
	}
	return errors.Wrap(err, "Unable to unmarshal content of the Prov response")
}

func ReadTextResponse(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot read response")
	}
	if response.StatusCode >= 400 {
		return nil, handleErrors(response.StatusCode, responseBody)
	}
	return responseBody, nil
}

func handleErrors(statusCode int, responseBody []byte) error {
	errs := new(Errors)
	err := json.Unmarshal(responseBody, errs)
	if err != nil {
		// No error message can be read : return error with http status
		return errors.New(fmt.Sprintf("Code: %d - Status: %s", statusCode, http.StatusText(statusCode)))
	}
	return errors.WithStack(errs)
}
