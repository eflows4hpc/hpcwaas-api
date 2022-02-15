package rest

import (
	"encoding/base64"
	"log"

	"github.com/gin-gonic/gin"
)

func basicAuth(basicConf *BasicAuthConfig) gin.HandlerFunc {
	if basicConf == nil {
		return nil
	}
	accounts := make(map[string]AuthAccount)
	for _, value := range basicConf.Accounts {
		key := "Basic " + base64.StdEncoding.EncodeToString([]byte(value.Username+":"+value.Password))
		accounts[key] = value
		log.Printf("Registering authentication for user: %q", value.Username)
	}

	return func(c *gin.Context) {
		// Search user in the slice of allowed credentials
		user, found := accounts[c.Request.Header.Get("Authorization")]

		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			writeError(c, newUnauthorizedRequest(c, "Authorization Required"))
			return
		}

		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(gin.AuthUserKey, user)
	}
}
