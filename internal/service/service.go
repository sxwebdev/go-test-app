package service

import (
	"context"

	"github.com/sxwebdev/go-test-app/internal/config"
	"github.com/sxwebdev/go-test-app/pb"
	"github.com/tkcrm/modules/logger"
	"google.golang.org/grpc"
)

type Service struct {
	logger logger.Logger
	config *config.Config
	grpc   *grpc.Server

	pb.UnimplementedHelloServiceServer
}

func Start(ctx context.Context, logger logger.Logger, config *config.Config) error {
	s := &Service{
		logger: logger,
		config: config,
	}

	// Start GRPC server
	go func() {
		if err := s.newGRPCServer(); err != nil {
			s.logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	return nil
}
