package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UpdateSession updates session with accessKeyId and userId.
func UpdateSession(c *gin.Context, accessKeyId string, userId int) error {
	session := sessions.Default(c)
	session.Set(sessionToken, accessKeyId)
	session.Set(sessionUserId, userId)
	err := session.Save()
	if err != nil {
		return err
	}
	c.Header("Access-Control-Allow-Credentials", "true")
	return nil
}
