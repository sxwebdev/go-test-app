package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/eclipse/paho.golang/paho"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/sxwebdev/go-test-app/internal/config"
	"github.com/sxwebdev/go-test-app/internal/listenerserver"
	"github.com/sxwebdev/go-test-app/internal/store"
	"github.com/sxwebdev/go-test-app/pb"
	"github.com/tkcrm/modules/logger"
)

type Server struct {
	config *config.Config
	logger logger.Logger
	store  store.Store

	fiber *fiber.App
	nats  *nats.Conn
	mqtt  *paho.Client

	mx              sync.RWMutex
	listenerServers map[string]*listenerserver.ListenerServer

	grpcClient pb.HelloServiceClient
}

func Start(l logger.Logger) error {
	s := &Server{
		logger:          l,
		listenerServers: make(map[string]*listenerserver.ListenerServer),
	}

	// Read configuration and envirioments
	config := config.New()
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration params validation error: %v", err)
	}
	s.config = config

	// Connect to database
	if err := s.newDB(); err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}

	// Start listeners servers
	sigCh := make(chan os.Signal, 1)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for server_name, port := range config.TCPServers {
		go func(server_name string, port uint16) {
			ls, err := listenerserver.New(l, server_name, port)
			if err != nil {
				s.logger.Errorf("start server %s error: %+v", server_name, err)
				sigCh <- os.Interrupt
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

	// Conect to external service throw GRPC
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

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	cancel()

	// Stop listener servers
	for _, ls := range s.listenerServers {
		if err := ls.StopServer(); err != nil {
			s.logger.Errorf("%+v", err)
		}
	}

	wg.Wait()
	os.Exit(0)

	return nil
}
