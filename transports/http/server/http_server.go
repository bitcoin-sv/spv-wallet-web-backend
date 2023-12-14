package httpserver

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"

	"bux-wallet/config"
	"bux-wallet/transports/http/endpoints/api/cors"
	"bux-wallet/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// GinEngineOpt represents functions to configure server engine.
type GinEngineOpt func(*gin.Engine)

// HttpServer represents server http.
type HttpServer struct {
	httpServer *http.Server
	handler    *gin.Engine
	log        *zerolog.Logger
}

// NewHttpServer creates and returns HttpServer instance.
func NewHttpServer(port int, log *zerolog.Logger) *HttpServer {
	httpLogger := log.With().Str("service", "http-server").Logger()

	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(debugWriter(&httpLogger)), gin.Recovery())
	engine.Use(cors.CorsMiddleware())

	return &HttpServer{
		httpServer: &http.Server{
			Addr:         ":" + fmt.Sprint(port),
			Handler:      engine,
			ReadTimeout:  time.Duration(viper.GetInt(config.EnvHttpServerReadTimeout)) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt(config.EnvHttpServerWriteTimeout)) * time.Second,
		},
		handler: engine,
		log:     &httpLogger,
	}
}

func debugWriter(logger *zerolog.Logger) io.Writer {
	w := func(p []byte) (n int, err error) {
		logger.Debug().Msg(string(p))
		return len(p), err
	}
	return util.WriterFunc(w)
}

// ApplyConfiguration it's entrypoint to configure a gin engine used by a server.
func (s *HttpServer) ApplyConfiguration(opts ...GinEngineOpt) {
	for _, opt := range opts {
		opt(s.handler)
	}
}

// Start is used to start http server.
func (s *HttpServer) Start() error {
	s.log.Info().Msgf("Starting server on address %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// ShutdownWithContext is used to stop http server using provided context.
func (s *HttpServer) ShutdownWithContext(ctx context.Context) error {
	s.log.Info().Msg("HTTP Server Shutdown")
	return s.httpServer.Shutdown(ctx)
}

// Shutdown is used to stop http server.
func (s *HttpServer) Shutdown() error {
	return s.ShutdownWithContext(context.Background())
}

// Logger return http server logger.
func (s *HttpServer) Logger() *zerolog.Logger {
	return s.log
}
