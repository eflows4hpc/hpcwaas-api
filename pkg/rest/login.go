package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) login(gc *gin.Context) {
	// TODO: randomize state here
	url := s.Config.Auth.OAuth.AuthCodeURL(state)
	gc.Redirect(http.StatusSeeOther, url)
}
