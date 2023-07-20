package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eflows4hpc/hpcwaas-api/api"
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

	// Save user info and token
	err = s.store.SaveSession(gc, userInfo, token)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}

	msg := fmt.Sprintf("\tLog in successful\nWelcome %s %s\nYou can now use HPCWaaS", userInfo.FirstName, userInfo.Surname)
	gc.String(http.StatusOK, msg)
}

func (s *Server) getUserInfo(ctx context.Context, accessToken string) (*api.UserInfo, error) {
	var res api.UserInfo
	url := userInfoEndpoint
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
