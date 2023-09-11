package websocket

import (
	"bux-wallet/transports/http/auth"
	"bux-wallet/transports/websocket"
	"github.com/centrifugal/centrifuge"
	"net/http"

	router "bux-wallet/transports/http/endpoints/routes"

	"github.com/gin-gonic/gin"
)

// NewHandler creates new endpoint handler.
func NewHandler(ws websocket.Server) router.ApiEndpoints {
	config := centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Register api endpoints which are authorized by session token.
	apiEndpoints := router.ApiEndpointsFunc(func(router *gin.RouterGroup) {
		router.Use(auth.GinContextToContextMiddleware())
		router.GET("/connection/websocket", gin.WrapH(auth.WsAuthMiddleware(centrifuge.NewWebsocketHandler(ws.GetNode(), config))))
	})

	return apiEndpoints
}
