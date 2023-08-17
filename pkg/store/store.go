package store

import (
	"errors"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Store interface {
	SaveSession(gc *gin.Context, userInfo *UserInfo, token *oauth2.Token) error
	LoadSession(gc *gin.Context) (*UserSession, error)
	ClearSession(gc *gin.Context) (bool, error)
}

// A class to store user sessions
type store struct {
	sessionDuration time.Duration
	userSessions    map[string]*UserSession
}

// Instanciate a store
//
// sessionMaxSeconds is the maximum length of a user session, in seconds
func NewStore(sessionMaxSeconds int64) Store {
	sessionDuration := time.Duration(sessionMaxSeconds) * time.Second
	return &store{
		sessionDuration: sessionDuration,
		userSessions:    map[string]*UserSession{},
	}
}

// Save user session
//
// No argument may be nil
func (s *store) SaveSession(gc *gin.Context, userInfo *UserInfo, token *oauth2.Token) error {
	if gc == nil {
		return errors.New("empty Gin context")
	}
	if userInfo == nil {
		return errors.New("empty user info")
	}
	if token == nil {
		return errors.New("empty access token")
	}

	// Save user ID in store
	session := sessions.Default(gc)
	if session == nil {
		return errors.New("no store found in Gin context")
	}
	session.Set(gin.AuthUserKey, userInfo.Sub)
	err := session.Save()
	if err != nil {
		return err
	}

	// Store session in map
	userSession := NewUserSession(*userInfo, *token, s.sessionDuration)
	s.userSessions[userInfo.Sub] = userSession
	return nil
}

// Load session for the current user
//
// Return nil if the user is not logged in
func (s *store) LoadSession(gc *gin.Context) (*UserSession, error) {
	if gc == nil {
		return nil, errors.New("empty Gin context")
	}

	// Get user ID from store
	session := sessions.Default(gc)
	if session == nil {
		return nil, errors.New("no store found in Gin context")
	}
	sub := session.Get(gin.AuthUserKey)
	if sub == nil {
		// User not logged in
		return nil, nil
	}

	// Get session for current user
	userSession := s.userSessions[sub.(string)]
	// This session may be nil if the server has been restarted
	// In this case, the user will need to log in again
	return userSession, nil
}

// Delete the session for the current user
//
// Return false if the user is not logged in
func (s *store) ClearSession(gc *gin.Context) (bool, error) {
	if gc == nil {
		return false, errors.New("empty Gin context")
	}

	// Get user ID from store
	session := sessions.Default(gc)
	if session == nil {
		return false, errors.New("no store found in Gin context")
	}
	subPtr := session.Get(gin.AuthUserKey)
	if subPtr == nil {
		// User not logged in
		return false, nil
	}

	sub := subPtr.(string)
	delete(s.userSessions, sub)

	session.Delete(gin.AuthUserKey)
	err := session.Save()
	if err != nil {
		return false, err
	}

	return true, nil
}
