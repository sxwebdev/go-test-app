package server

import (
	"context"
	"sync"

	"github.com/eclipse/paho.golang/paho"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/sxwebdev/go-test-app/internal/config"
	"github.com/sxwebdev/go-test-app/internal/listenerserver"
	"github.com/sxwebdev/go-test-app/internal/store"
	"github.com/sxwebdev/go-test-app/pb"
	"github.com/tkcrm/modules/db/bunconn"
	"github.com/tkcrm/modules/logger"
)

type Server struct {
	logger logger.Logger
	config *config.Config
	db     *bunconn.BunConn
	store  store.Store

	fiber *fiber.App
	nats  *nats.Conn
	mqtt  *paho.Client

	mx              sync.RWMutex
	listenerServers map[string]*listenerserver.ListenerServer

	grpcClient pb.HelloServiceClient
}

func New(ctx context.Context, logger logger.Logger, config *config.Config) (*Server, error) {
	s := &Server{
		logger:          logger,
		config:          config,
		listenerServers: make(map[string]*listenerserver.ListenerServer),
	}

	// Connect to database
	conn, err := bunconn.New(config.DB, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing db")
	}
	s.db = conn
	s.store = s.db

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	defer func() {
		s.db.Close()
		s.logger.Info("server stopped")
	}()

	// Start listeners servers
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for server_name, port := range s.config.TCPServers {
		go func(server_name string, port uint16) {
			ls, err := listenerserver.New(s.logger, server_name, port)
			if err != nil {
				s.logger.Errorf("start server %s error: %+v", server_name, err)
				cancel()
			}
			go func() {
				defer wg.Done()
				wg.Add(1)

				if err := ls.AcceptConnections(ctx); err != nil {
					s.logger.Errorf("accept connection error: %+v", err)
				}
			}()
			s.mx.Lock()
			s.listenerServers[server_name] = ls
			s.mx.Unlock()
		}(server_name, port)
	}

	// Start API server
	go func() {
		s.newApiServer()
		if err := s.ApiStart(); err != nil {
			s.logger.Fatal("api server error:", err.Error())
		}
	}()

	// Conect to external service through GRPC
	go func() {
		if err := s.grpcConnect(); err != nil {
			s.logger.Fatal(err)
		}
	}()

	// Connect to NATS
	go func() {
		if err := s.natsConnect(); err != nil {
			s.logger.Fatal(err)
		}
	}()

	// Connect to MQTT
	go func() {
		if err := s.mqttConnect(); err != nil {
			s.logger.Fatal(err)
		}
	}()

	s.logger.Info("server successfully started")

	<-ctx.Done()

	// Stop listener servers
	for _, ls := range s.listenerServers {
		if err := ls.StopServer(); err != nil {
			s.logger.Errorf("%+v", err)
		}
	}

	return nil
}
