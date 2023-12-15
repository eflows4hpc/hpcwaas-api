package store

import (
	"errors"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Store interface {
	CreateSession(gc *gin.Context, userInfo *UserInfo, accessToken string) error
	GetSession(gc *gin.Context, accessToken string) (*UserSession, error)
	DeleteSession(gc *gin.Context, accessToken string) error
}

// A class to store user sessions
type store struct {
	sessionDuration time.Duration
	userSessions    map[string]UserSession
}

// Instanciate a store
//
// sessionMaxSeconds is the maximum length of a user session, in seconds
func NewStore(sessionMaxSeconds int64) Store {
	sessionDuration := time.Duration(sessionMaxSeconds) * time.Second
	return &store{
		sessionDuration: sessionDuration,
		userSessions:    make(map[string]UserSession, 0),
	}
}

// Create and store a new user session
func (s *store) CreateSession(gc *gin.Context, userInfo *UserInfo, accessToken string) error {
	if gc == nil {
		return errors.New("empty Gin context")
	}
	if userInfo == nil {
		return errors.New("empty user info")
	}
	if accessToken == "" {
		return errors.New("empty access token")
	}

	// Save user ID in store
	session := sessions.Default(gc)
	session.Set(gin.AuthUserKey, userInfo.Sub)
	err := session.Save()
	if err != nil {
		return err
	}

	// Store session in map
	userSession := NewUserSession(*userInfo, s.sessionDuration)
	s.userSessions[accessToken] = *userSession
	return nil
}

// Get the session associated to the specified token
//
// Return nil if the token is not present in the store
func (s *store) GetSession(gc *gin.Context, accessToken string) (*UserSession, error) {
	if gc == nil {
		return nil, errors.New("empty Gin context")
	}
	if accessToken == "" {
		return nil, errors.New("empty access token")
	}

	// Get session for specified token
	userSession, present := s.userSessions[accessToken]
	if !present {
		return nil, nil
	}

	return &userSession, nil
}

// Delete the session associated to the specified token
func (s *store) DeleteSession(gc *gin.Context, accessToken string) error {
	if gc == nil {
		return errors.New("empty Gin context")
	}
	if accessToken == "" {
		return errors.New("empty access token")
	}

	delete(s.userSessions, accessToken)
	return nil
}
