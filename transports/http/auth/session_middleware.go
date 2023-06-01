package auth

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	sessionToken  = "token"
	sessionUserId = "userId"
)

// SessionMiddleware middleware that is retrieving auth token from cookie.
type SessionMiddleware struct{}

// NewSessionMiddleware create Session middleware that is retrieving auth token from cookie.
func NewSessionMiddleware() *SessionMiddleware {
	return &SessionMiddleware{}
}

// ApplyToApi is a middleware which checks if the request has a valid token.
func (h *SessionMiddleware) ApplyToApi(c *gin.Context) {
	session := sessions.Default(c)

	// Try to retrieve session token.
	token := session.Get(sessionToken)
	if token == nil {
		c.AbortWithStatusJSON(401, errors.New("missing token in cookie"))
		return
	}

	// Try to retrieve session user id.
	userId := session.Get(sessionUserId)
	if userId == nil {
		c.AbortWithStatusJSON(401, errors.New("missing user id in cookie"))
		return
	}

	// Set token and user id in gin context.
	c.Set("token", token)
	c.Set("userId", userId)
}
