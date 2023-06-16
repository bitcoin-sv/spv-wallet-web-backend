package auth

import (
	"bux-wallet/domain/users"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UpdateSession updates session with accessKeyId and userId.
func UpdateSession(c *gin.Context, accessKey users.AccessKey, userId int) error {
	session := sessions.Default(c)
	session.Set(SessionAccessKeyId, accessKey.Id)
	session.Set(SessionAccessKey, accessKey.Key)
	session.Set(SessionUserId, userId)
	err := session.Save()
	if err != nil {
		return err
	}
	c.Header("Access-Control-Allow-Credentials", "true")
	return nil
}

// TerminateSession terminates current (default) session.
func TerminateSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()

	err := session.Save()
	if err != nil {
		return err
	}

	return nil
}
