package rest

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/eflows4hpc/hpcwaas-api/pkg/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	unityEndpoint oauth2.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://zam10045.zam.kfa-juelich.de:7000/oauth2-as/oauth2-authz",
		TokenURL: "https://zam10045.zam.kfa-juelich.de:7000/oauth2/token",
	}
	userInfoEndpoint string         = "https://zam10045.zam.kfa-juelich.de:7000/oauth2/userinfo"
	oauthConf        *oauth2.Config = &oauth2.Config{
		ClientID:     "580b8e3e-b4f8-444a-8f04-841a1dd3453b",
		ClientSecret: "b41b1c24-58de-487d-bfc2-e3892ecd2f45",
		Endpoint:     unityEndpoint,
		Scopes:       []string{"profile", "email", "eflows"},
		RedirectURL:  "http://localhost:9090/auth/authorize",
	}
	sessionDuration time.Duration = time.Hour * 24
)

// getRandomState returns a number of random bytes, encoded in base64
func getRandomState(length int) string {
	randomBytes := util.SecureRandomBytes(length)
	return base64.StdEncoding.EncodeToString(randomBytes)
}

func (s *Server) initSsoConf() {
	s.Config.Auth.State = getRandomState(64)
	s.Config.Auth.OAuth2 = oauthConf
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
