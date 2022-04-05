package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

type SSHKeysService interface {
	GenerateSSHKey(ctx context.Context, genRequest SSHKeyGenerationRequest) (SSHKey, error)
}

type sshKeysService struct {
	client *client
}

func (s *sshKeysService) GenerateSSHKey(ctx context.Context, genRequest SSHKeyGenerationRequest) (SSHKey, error) {
	var res SSHKey
	body, err := json.Marshal(genRequest)
	if err != nil {
		return res, errors.Wrap(err, "Fail to marshall ssh generation request")
	}
	request, err := s.client.NewRequest(ctx, http.MethodPost, path.Join("/ssh_keys"), bytes.NewReader(body))
	if err != nil {
		return res, errors.Wrap(err, "failed to create http request")
	}
	request.Header.Add("Accept", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return res, errors.Wrap(err, "failed to send http request to generate ssh key")
	}

	err = ReadResponse(response, &res)
	return res, err
}
