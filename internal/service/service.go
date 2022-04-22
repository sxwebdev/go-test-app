package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sxwebdev/go-test-app/internal/config"
	"github.com/sxwebdev/go-test-app/pb"
	"github.com/tkcrm/modules/logger"
	"google.golang.org/grpc"
)

type Service struct {
	config *config.Config
	logger logger.Logger
	grpc   *grpc.Server

	pb.UnimplementedHelloServiceServer
}

func Start(l logger.Logger) error {
	s := &Service{
		logger: l,
	}

	// Read configuration and envirioments
	config := config.New()
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration params validation error: %v", err)
	}
	s.config = config

	// Start listeners servers
	sigCh := make(chan os.Signal, 1)

	// Start GRPC server
	go func() {
		if err := s.newGRPCServer(); err != nil {
			s.logger.Fatal(err)
		}
	}()

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	os.Exit(0)
	return nil
}
