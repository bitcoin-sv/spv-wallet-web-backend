package auth

import (
	"bux-wallet/domain/users"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UpdateSession updates session with accessKeyId and userId.
func UpdateSession(c *gin.Context, accessKey users.AccessKey, userId int) error {
	fmt.Println("UpdateSession")
	fmt.Println(accessKey.Id)
	fmt.Println(accessKey.Key)
	fmt.Println(userId)

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
