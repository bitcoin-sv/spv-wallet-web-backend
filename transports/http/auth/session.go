package auth

import (
	"bux-wallet/domain/users"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UpdateSession updates session with accessKeyId and userId.
func UpdateSession(c *gin.Context, user *users.AuthenticatedUser) error {
	session := sessions.Default(c)
	session.Set(SessionAccessKeyId, user.AccessKey.Id)
	session.Set(SessionAccessKey, user.AccessKey.Key)
	session.Set(SessionUserId, user.User.Id)
	session.Set(SessionUserPaymail, user.User.Paymail)
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
