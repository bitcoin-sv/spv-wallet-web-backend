package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type contextKey int

var ginContextKey contextKey

// GinContextToContextMiddleware add gin context to context for centrifuge.
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinContextFromContext recover the gin context from the context.Context.
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}
	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

// WsAuthMiddleware is used to authenticate websocket connections.
func WsAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		gc, err := GinContextFromContext(ctx)
		if err != nil {
			return
		}
		s := sessions.Default(gc)
		userId := s.Get(SessionUserId)
		if userId == "" {
			return
		}
		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
			UserID: strconv.Itoa(userId.(int)),
		})
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}
