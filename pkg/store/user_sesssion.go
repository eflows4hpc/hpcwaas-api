package store

import (
	"time"
)

type UserSession struct {
	UserInfo  UserInfo  `json:"user_info"`
	StartedAt time.Time `json:"started_at"`
	ExpireAt  time.Time `json:"expire_at"`
}

func NewUserSession(userInfo UserInfo, validity time.Duration) *UserSession {
	startTime := time.Now()
	userSession := UserSession{
		UserInfo:  userInfo,
		StartedAt: startTime,
		ExpireAt:  startTime.Add(validity),
	}
	return &userSession
}

func (us *UserSession) IsExpired() bool {
	now := time.Now()
	return now.After(us.ExpireAt)
}
