package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) logout(gc *gin.Context) {
	// To log out, we just remove user info from the store
	loggedIn, err := s.store.ClearSession(gc)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}

	if !loggedIn {
		gc.String(http.StatusOK, "You are not logged in")
	} else {
		gc.String(http.StatusOK, "Logout successsful")
	}
}
