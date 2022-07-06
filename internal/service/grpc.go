package service

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/sxwebdev/go-test-app/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func (s *Service) newGRPCServer() error {

	s.grpc = grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(s.grpc, health.NewServer())
	pb.RegisterHelloServiceServer(s.grpc, s)

	reflection.Register(s.grpc)

	listen, err := net.Listen("tcp", s.config.GrpcDSN)
	if err != nil {
		return err
	}

	s.logger.Infof("GRPC server start successfully on %s", s.config.GrpcDSN)

	if err := s.grpc.Serve(listen); err != nil {
		return err
	}

	return nil
}

func (s *Service) Say(ctx context.Context, req *pb.SayRequest) (*pb.SayResponse, error) {

	resp := &pb.SayResponse{
		Message: fmt.Sprintf("Received message: %s", req.GetMessage()),
		Time:    time.Now().String(),
		Rand:    uint32(rand.Int()),
	}

	return resp, nil
}
