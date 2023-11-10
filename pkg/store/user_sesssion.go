package store

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"golang.org/x/oauth2"
)

type UserSession struct {
	UserInfo UserInfo     `json:"user_info"`
	Token    oauth2.Token `json:"token"`
	ExpireAt time.Time    `json:"expire_at"`
}

func NewUserSession(userInfo UserInfo, token oauth2.Token, validity time.Duration) *UserSession {
	userSession := UserSession{
		UserInfo: userInfo,
		Token:    token,
		ExpireAt: time.Now().Add(validity),
	}
	return &userSession
}

func (us *UserSession) IsExpired() bool {
	now := time.Now()
	return now.After(us.ExpireAt)
}

func (us *UserSession) IsTokenExpired() bool {
	now := time.Now()
	return now.After(us.Token.Expiry)
}

func (us *UserSession) RefreshToken(conf *oauth2.Config) error {
	refreshedToken, err := conf.TokenSource(context.Background(), &us.Token).Token()
	if err != nil {
		return err
	}
	us.Token = *refreshedToken
	return nil
}

func (us *UserSession) Write(filename string) error {
	bytes, err := json.Marshal(us)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, bytes, 0600)
	return err
}

func ReadUserSession(filename string) (*UserSession, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	userSession := new(UserSession)
	err = json.Unmarshal(bytes, userSession)
	if err != nil {
		return nil, err
	}

	return userSession, nil
}
