package api

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetClient(t *testing.T) {

	tests := []struct {
		name        string
		cc          Configuration
		wantErr     bool
		errContains string
	}{
		//		{"simpleConfig", Configuration{}, false, ""},
		{"simpleTLSConfig", Configuration{
			APIURL:        "my.com",
			SSLEnabled:    true,
			SkipTLSVerify: true,
			KeyFile:       "testdata/comp.key",
			CertFile:      "testdata/comp.pem",
			CAFile:        "testdata/ca.pem",
			CAPath:        "testdata",
		}, false, ""},
		{"badApiURL", Configuration{APIURL: "thisIsBad!!!", SSLEnabled: true}, true, "Malformed API URL"},
		{"badApiURLPort", Configuration{APIURL: "thisIsBad:very", SSLEnabled: true}, true, "Malformed API URL"},
		{"badCerts", Configuration{APIURL: "my.com", SSLEnabled: true, KeyFile: "missing", CertFile: "missing"}, true, "Failed to load TLS certificates"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient, err := GetClient(tt.cc)
			if tt.wantErr != (err != nil) {
				t.Errorf("GetClient() err= %v, wantError %v", err, tt.wantErr)
			}
			if tt.wantErr == true && tt.errContains != "" {
				assert.ErrorContains(t, err, tt.errContains)
			} else {
				assert.Assert(t, httpClient != nil)
			}
		})
	}
}
