package httpserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/api/cors"
	"github.com/bitcoin-sv/spv-wallet-web-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// GinEngineOpt represents functions to configure server engine.
type GinEngineOpt func(*gin.Engine)

// HTTPServer represents server http.
type HTTPServer struct {
	httpServer *http.Server
	handler    *gin.Engine
	log        *zerolog.Logger
}

// NewHTTPServer creates and returns HTTPServer instance.
func NewHTTPServer(port int, log *zerolog.Logger) *HTTPServer {
	httpLogger := log.With().Str("service", "http-server").Logger()

	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(debugWriter(&httpLogger)), gin.Recovery())
	engine.Use(cors.Middleware())

	return &HTTPServer{
		httpServer: &http.Server{
			Addr:         ":" + fmt.Sprint(port),
			Handler:      engine,
			ReadTimeout:  time.Duration(viper.GetInt(config.EnvHTTPServerReadTimeout)) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt(config.EnvHTTPServerWriteTimeout)) * time.Second,
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
func (s *HTTPServer) ApplyConfiguration(opts ...GinEngineOpt) {
	for _, opt := range opts {
		opt(s.handler)
	}
}

// Start is used to start http server.
func (s *HTTPServer) Start() error {
	s.log.Info().Msgf("Starting server on address %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe() //nolint:wrapcheck // no need for wrapping as we will stop everything on any error here.
}

// ShutdownWithContext is used to stop http server using provided context.
func (s *HTTPServer) ShutdownWithContext(ctx context.Context) error {
	s.log.Info().Msg("HTTP Server Shutdown")
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck // no need for wrapping as it is done at the shutting down flow.
}

// Shutdown is used to stop http server.
func (s *HTTPServer) Shutdown() error {
	return s.ShutdownWithContext(context.Background())
}

// Logger return http server logger.
func (s *HTTPServer) Logger() *zerolog.Logger {
	return s.log
}
