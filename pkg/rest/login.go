package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) login(gc *gin.Context) {
	url := s.Config.Auth.OAuth2.AuthCodeURL(s.Config.Auth.State)
	gc.Redirect(http.StatusSeeOther, url)
}
