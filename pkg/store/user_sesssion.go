package store

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type UserSession struct {
	UserInfo UserInfo
	Token    oauth2.Token
	ExpireAt time.Time
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
