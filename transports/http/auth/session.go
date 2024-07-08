package auth

import (
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// UpdateSession updates session with accessKeyId and userId.
func UpdateSession(c *gin.Context, authUser *users.AuthenticatedUser) error {
	session := sessions.Default(c)
	session.Set(SessionAccessKeyID, authUser.AccessKey.ID)
	session.Set(SessionAccessKey, authUser.AccessKey.Key)
	session.Set(SessionUserID, authUser.User.ID)
	session.Set(SessionUserPaymail, authUser.User.Paymail)
	session.Set(SessionXPriv, authUser.Xpriv)
	err := session.Save()
	if err != nil {
		return errors.Wrap(err, "internal error")
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
		return errors.Wrap(err, "internal error")
	}

	return nil
}
