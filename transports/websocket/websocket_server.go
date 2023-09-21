package websocket

import (
	"bux-wallet/domain"
	"bux-wallet/logging"
	"bux-wallet/transports/http/auth"
	router "bux-wallet/transports/http/endpoints/routes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
)

// Server websocket server controller.
type Server interface {
	Start() error
	Shutdown() error
	SetupEntrypoint(*gin.Engine)
	GetNode() *centrifuge.Node
	GetSocket(userId string) *Socket
	GetSockets() map[string]*Socket
}

type server struct {
	node     *centrifuge.Node
	log      logging.Logger
	sockets  map[string]*Socket
	services *domain.Services
	db       *sql.DB
}

// NewServer creates new websocket server.
func NewServer(lf logging.LoggerFactory, services *domain.Services, db *sql.DB) (Server, error) {
	node, err := newNode(lf)
	if err != nil {
		return nil, err
	}
	s := &server{
		node:     node,
		log:      lf.NewLogger("Websocket"),
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
	return nil
}

// Shutdown stoping a server.
func (s *server) Shutdown() error {
	return s.ShutdownWithContext(context.Background())
}

// ShutdownWithContext stoping a server in a provided context.
func (s *server) ShutdownWithContext(ctx context.Context) error {
	s.log.Infof("Shutting down a websocket server")
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

	r := engine.Group("/connection/websocket", apiMiddlewares...)

	config := centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	r.Use(auth.GinContextToContextMiddleware())
	r.GET("", gin.WrapH(auth.WsAuthMiddleware(centrifuge.NewWebsocketHandler(s.GetNode(), config))))
}

func newNode(lf logging.LoggerFactory) (*centrifuge.Node, error) {
	lh := newLogHandler(lf)
	return centrifuge.New(centrifuge.Config{
		Name:       "bux-wallet",
		LogLevel:   lh.Level(),
		LogHandler: lh.Log,
	})
}

type connectData struct {
	Email string `json:"email"`
}

func (s *server) setupNode() {
	s.node.OnConnecting(func(ctx context.Context, event centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
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

		client.OnRefresh(func(e centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			cb(centrifuge.RefreshReply{
				ExpireAt: time.Now().Unix() + 10,
			}, nil)
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			cb(centrifuge.SubscribeReply{
				Options: centrifuge.SubscribeOptions{
					EmitPresence: true,
				},
			}, nil)
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			delete(s.sockets, client.ID())
		})
	})

}

func (s *server) GetNode() *centrifuge.Node {
	return s.node
}

func (s *server) GetSocket(userId string) *Socket {
	return s.sockets[userId]
}

func (s *server) GetSockets() map[string]*Socket {
	return s.sockets
}
