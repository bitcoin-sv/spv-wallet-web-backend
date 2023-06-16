package httpserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"bux-wallet/config"
	"bux-wallet/logging"
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
	log        logging.Logger
}

// NewHttpServer creates and returns HttpServer instance.
func NewHttpServer(port int, lf logging.LoggerFactory) *HttpServer {
	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(debugWriter(lf.NewLogger("gin"))), gin.Recovery())
	engine.Use(cors.CorsMiddleware())

	return &HttpServer{
		httpServer: &http.Server{
			Addr:         ":" + fmt.Sprint(port),
			Handler:      engine,
			ReadTimeout:  time.Duration(viper.GetInt(config.EnvHttpServerReadTimeout)) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt(config.EnvHttpServerWriteTimeout)) * time.Second,
		},
		handler: engine,
		log:     lf.NewLogger("http"),
	}
}

func debugWriter(logger logging.Logger) io.Writer {
	w := func(p []byte) (n int, err error) {
		logger.Debug(string(p))
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
	s.log.Infof("Starting server on address %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// ShutdownWithContext is used to stop http server using provided context.
func (s *HttpServer) ShutdownWithContext(ctx context.Context) error {
	s.log.Info("HTTP Server Shutdown")
	return s.httpServer.Shutdown(ctx)
}

// Shutdown is used to stop http server.
func (s *HttpServer) Shutdown() error {
	return s.ShutdownWithContext(context.Background())
}
