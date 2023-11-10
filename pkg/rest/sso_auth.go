package rest

import (
	"log"
	"net/http"

	"github.com/eflows4hpc/hpcwaas-api/pkg/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
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
		userSession, err := s.store.LoadSession(gc)
		if err != nil {
			writeError(gc, newInternalServerError(err))
			return
		}

		if userSession == nil || userSession.IsExpired() {
			// User is not logged in, we redirect to authorize endpoint
			url := oauthConf.AuthCodeURL(s.Config.Auth.State)
			gc.Redirect(http.StatusTemporaryRedirect, url)
			return
		}

		if userSession.IsTokenExpired() {
			userSession.RefreshToken(s.Config.Auth.OAuth2)
		}
		gc.Next()
	}
}
