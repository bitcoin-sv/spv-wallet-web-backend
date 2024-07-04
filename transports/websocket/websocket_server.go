package websocket

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/auth"
	router "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/routes"
	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Server websocket server controller.
type Server interface {
	Start() error
	Shutdown() error
	SetupEntrypoint(*gin.Engine)
	GetNode() *centrifuge.Node
	GetSocket(userID string) *Socket
	GetSockets() map[string]*Socket
}

type server struct {
	node     *centrifuge.Node
	log      *zerolog.Logger
	sockets  map[string]*Socket
	services *domain.Services
	db       *sql.DB
}

// NewServer creates new websocket server.
func NewServer(log *zerolog.Logger, services *domain.Services, db *sql.DB) (Server, error) {
	websocketLogger := log.With().Str("service", "websocket").Logger()

	node, err := newNode(&websocketLogger)
	if err != nil {
		return nil, err
	}
	s := &server{
		node:     node,
		log:      &websocketLogger,
		sockets:  make(map[string]*Socket),
		services: services,
		db:       db,
	}
	return s, nil
}

// Start starts a server.
func (s *server) Start() error {
	s.setupNode()
	if err := s.node.Run(); err != nil {
		return fmt.Errorf("cannot start websocket server: %w", err)
	}
	s.log.Debug().Msgf("Websocket server started")
	return nil
}

// Shutdown stoping a server.
func (s *server) Shutdown() error {
	return s.ShutdownWithContext(context.Background())
}

// ShutdownWithContext stoping a server in a provided context.
func (s *server) ShutdownWithContext(ctx context.Context) error {
	s.log.Info().Msgf("Shutting down a websocket server")
	if err := s.node.Shutdown(ctx); err != nil {
		return fmt.Errorf("cannot stop websocket server: %w", err)
	}
	return nil
}

// SetupEntrypoint setup gin to init websocket connection.
func (s *server) SetupEntrypoint(engine *gin.Engine) {
	apiMiddlewares := router.ToHandlers(
		auth.NewSessionMiddleware(s.db, engine),
		auth.NewAuthMiddleware(s.services),
	)
	r := engine.Group("/api/websocket", apiMiddlewares...)

	config := centrifuge.WebsocketConfig{
		CheckOrigin: func(_ *http.Request) bool { return true },
	}

	r.Use(auth.GinContextToContextMiddleware())
	r.GET("", gin.WrapH(auth.WsAuthMiddleware(centrifuge.NewWebsocketHandler(s.GetNode(), config))))
}

func newNode(l *zerolog.Logger) (*centrifuge.Node, error) {
	lh := newLogHandler(l)
	return centrifuge.New(centrifuge.Config{
		Name:       "spv-wallet",
		LogLevel:   lh.Level(),
		LogHandler: lh.Log,
	})
}

type connectData struct {
	Email string `json:"email"`
}

func (s *server) setupNode() {
	s.node.OnConnecting(func(ctx context.Context, _ centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		cred, ok := centrifuge.GetCredentials(ctx)
		if !ok {
			return centrifuge.ConnectReply{}, centrifuge.DisconnectServerError
		}
		data, _ := json.Marshal(connectData{
			Email: cred.UserID,
		})
		return centrifuge.ConnectReply{
			Data: data,
		}, nil
	})

	s.node.OnConnect(func(client *centrifuge.Client) {
		s.sockets[client.UserID()] = &Socket{
			Client: client,
			Log:    s.log,
		}

		client.OnRefresh(func(_ centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			cb(centrifuge.RefreshReply{
				ExpireAt: time.Now().Unix() + 10,
			}, nil)
		})

		client.OnSubscribe(func(_ centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			cb(centrifuge.SubscribeReply{
				Options: centrifuge.SubscribeOptions{
					EmitPresence: true,
				},
			}, nil)
		})

		client.OnDisconnect(func(_ centrifuge.DisconnectEvent) {
			delete(s.sockets, client.ID())
		})
	})

}

func (s *server) GetNode() *centrifuge.Node {
	return s.node
}

func (s *server) GetSocket(userID string) *Socket {
	userSocket := s.sockets[userID]
	if userSocket == nil {
		userSocket = &Socket{Log: s.log}
	}
	return userSocket
}

func (s *server) GetSockets() map[string]*Socket {
	return s.sockets
}
