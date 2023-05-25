package httpserver

import (
	"bux-wallet/config"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// GinEngineOpt represents functions to configure server engine.
type GinEngineOpt func(*gin.Engine)

// HttpServer represents server http.
type HttpServer struct {
	httpServer *http.Server
	handler    *gin.Engine
}

// NewHttpServer creates and returns HttpServer instance.
func NewHttpServer(port int) *HttpServer {
	handler := gin.Default()

	return &HttpServer{
		httpServer: &http.Server{
			Addr:         ":" + fmt.Sprint(port),
			Handler:      handler,
			ReadTimeout:  time.Duration(viper.GetInt(config.EnvHttpServerReadTimeout)) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt(config.EnvHttpServerWriteTimeout)) * time.Second,
		},
		handler: handler,
	}
}

// ApplyConfiguration it's entrypoint to configure a gin engine used by a server.
func (s *HttpServer) ApplyConfiguration(opts ...GinEngineOpt) {
	for _, config := range opts {
		config(s.handler)
	}
}

// Start is used to start http server.
func (s *HttpServer) Start() error {
	return s.httpServer.ListenAndServe()
}

// ShutdownWithContext is used to stop http server using provided context.
func (s *HttpServer) ShutdownWithContext(ctx context.Context) error {
	fmt.Println("HTTP Server Shutdown")
	return s.httpServer.Shutdown(ctx)
}

// Shutdown is used to stop http server.
func (s *HttpServer) Shutdown() error {
	return s.ShutdownWithContext(context.Background())
}
