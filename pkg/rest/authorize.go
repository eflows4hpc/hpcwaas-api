package rest

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/eflows4hpc/hpcwaas-api/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) authorize(gc *gin.Context) {
	// Check state
	requestState := gc.Request.FormValue("state")
	if requestState != s.Config.Auth.State {
		writeError(gc, newBadRequestMessage("request state doesn't match session state"))
		return
	}

	// Exchange code
	authorizationCode := gc.Request.FormValue("code")
	token, err := s.Config.Auth.OAuth2.Exchange(context.Background(), authorizationCode)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}

	// Get user info from endpoint
	userInfo, err := s.getUserInfo(gc, token.AccessToken)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}

	// Start a new session with user info and token
	err = s.store.CreateSession(gc, userInfo, token.AccessToken)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}

	encodedToken := base64.StdEncoding.EncodeToString([]byte(token.AccessToken))
	msg := fmt.Sprintf(`	Log in successful

Welcome %s %s
You can now use HPCWaaS

For using the CLI, please use the following token:
	%s
`, userInfo.FirstName, userInfo.Surname, encodedToken)

	gc.String(http.StatusOK, msg)
}

func (s *Server) getUserInfo(ctx context.Context, accessToken string) (*store.UserInfo, error) {
	var res store.UserInfo
	url := s.Config.Auth.UserInfoURL
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send http request to get user info")
	}

	err = api.ReadResponse(response, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
