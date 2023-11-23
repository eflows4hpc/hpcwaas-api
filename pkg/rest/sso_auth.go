package rest

import (
	"encoding/base64"
	"log"
	"regexp"

	"github.com/eflows4hpc/hpcwaas-api/pkg/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	bearerPattern = regexp.MustCompile("^Bearer *([^ ]+) *$")
)

// getRandomState returns a number of random bytes, encoded in base64
func getRandomState(length int) string {
	return util.SecureRandomSecret(length)
}

func (s *Server) initSsoConf() {
	auth := s.Config.Auth
	s.Config.Auth.OAuth2 = &oauth2.Config{
		ClientID:     auth.ClientID,
		ClientSecret: auth.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  auth.AuthURL,
			TokenURL: auth.TokenURL,
		},
		Scopes:      auth.Scopes,
		RedirectURL: auth.RedirectURL,
	}
	s.Config.Auth.State = getRandomState(64)
}

func (s *Server) ssoAuth(oauthConf *oauth2.Config) gin.HandlerFunc {
	if oauthConf == nil {
		log.Fatal("Empty oauth2 config")
	}

	return func(gc *gin.Context) {
		authorization := gc.Request.Header.Get("Authorization")
		if authorization == "" {
			writeError(gc, newUnauthorizedRequest(gc, "Authorization Required"))
			return
		}
		if !bearerPattern.MatchString(authorization) {
			writeError(gc, newUnauthorizedRequest(gc, "Invalid authorization format"))
			return
		}
		base64AccessToken := bearerPattern.FindStringSubmatch(authorization)[1]
		bytesAccessToken, err := base64.StdEncoding.DecodeString(base64AccessToken)
		if err != nil {
			writeError(gc, newUnauthorizedRequest(gc, "Invalid authorization token"))
			return
		}
		accessToken := string(bytesAccessToken)

		userSession, err := s.store.GetSession(gc, accessToken)
		if err != nil || userSession == nil || userSession.IsExpired() {
			writeError(gc, newUnauthorizedRequest(gc, "You are not logged in or your session has expired"))
			return
		}

		gc.Next()
	}
}
