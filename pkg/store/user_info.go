package store

// UserInfo is the response structure of a GetUserInfo operation
type UserInfo struct {
	Sub       string `json:"sub"`
	FirstName string `json:"firstname"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
}
