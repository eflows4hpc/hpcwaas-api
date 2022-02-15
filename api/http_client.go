package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

	Workflows() WorkflowsService
	Executions() ExecutionsService
	Users() UsersService
}

func setupTLSConfig(client *client, cc Configuration, url *url.URL) error {

	url.Scheme = "https"
	apiHost, _, err := urlx.SplitHostPort(url)
	if err != nil {
		return errors.Wrap(err, "Malformed API URL")
	}

	tlsConfig := &tls.Config{ServerName: apiHost}
	if cc.CertFile != "" && cc.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cc.CertFile, cc.KeyFile)
		if err != nil {
			return errors.Wrap(err, "Failed to load TLS certificates")
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	if cc.CAFile != "" || cc.CAPath != "" {
		cfg := &rootcerts.Config{
			CAFile: cc.CAFile,
			CAPath: cc.CAPath,
		}
		rootcerts.ConfigureTLS(tlsConfig, cfg)
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	if cc.SkipTLSVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client.Client.Transport = tr
	return nil
}

// GetClient returns a HTTP Client
func GetClient(cc Configuration) (HTTPClient, error) {
	apiURL := cc.APIURL
	apiURL = strings.TrimRight(apiURL, "/")

	httpClient := cc.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 1 * time.Minute,
		}
	}

	u, err := urlx.ParseWithDefaultScheme(apiURL, "https")
	if err != nil {
		return nil, errors.Wrap(err, "Malformed API URL")
	}

	client := &client{
		Client: httpClient,
	}

	if strings.ToLower(u.Scheme) == "https" || cc.CAFile != "" || cc.CAPath != "" || (cc.CertFile != "" && cc.KeyFile != "") {
		err := setupTLSConfig(client, cc, u)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to setup TLS config")
		}
	}
	client.baseURL = u.String()

	if cc.User != "" {
		parts := strings.SplitN(cc.User, ":", 2)
		if len(parts) != 2 {
			return nil, errors.New("Invalid user format, expecting 'user:password'")
		}
		client.basicAuthUserPass = url.UserPassword(parts[0], parts[1])
	}

	client.workflows = &workflowsService{client: client}
	client.executions = &executionsService{client: client}
	client.users = &usersService{client: client}
	return client, nil
}

// client is the HTTP client structure
type client struct {
	*http.Client
	baseURL    string
	workflows  WorkflowsService
	executions ExecutionsService
	users      UsersService

	basicAuthUserPass *url.Userinfo
}

func (c *client) Workflows() WorkflowsService {
	return c.workflows
}

func (c *client) Executions() ExecutionsService {
	return c.executions
}
func (c *client) Users() UsersService {
	return c.users
}

// NewRequest returns a new HTTP request
func (c *client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot create request")
	}

	req.URL.User = c.basicAuthUserPass
	return req, nil
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
