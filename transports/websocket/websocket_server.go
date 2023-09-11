package websocket

import (
	"bux-wallet/domain"
	"bux-wallet/domain/websockets"
	"bux-wallet/logging"
	"context"
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
}

type server struct {
	node    *centrifuge.Node
	log     logging.Logger
	sockets map[string]*websockets.Socket
}

// NewServer creates new websocket server.
func NewServer(lf logging.LoggerFactory, services *domain.Services) (Server, error) {
	node, err := newNode(lf)
	if err != nil {
		return nil, err
	}
	s := &server{
		node:    node,
		log:     lf.NewLogger("Websocket"),
		sockets: services.TransactionsService.Websockets,
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
	config := centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	engine.GET("/connection/websocket", gin.WrapH(centrifuge.NewWebsocketHandler(s.node, config)))
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
		s.sockets[client.UserID()] = &websockets.Socket{
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
