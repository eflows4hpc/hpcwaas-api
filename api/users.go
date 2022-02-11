package api

import (
	"context"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

type UsersService interface {
	GenerateSSHKey(ctx context.Context, userName string) (SSHKey, error)
}

type usersService struct {
	client *client
}

func (s *usersService) GenerateSSHKey(ctx context.Context, userName string) (SSHKey, error) {
	var res SSHKey
	request, err := s.client.NewRequest(ctx, http.MethodPost, path.Join("/users", userName, "ssh_key"), nil)
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
