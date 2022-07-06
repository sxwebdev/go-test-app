package server

import (
	"github.com/sxwebdev/go-test-app/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Server) grpcConnect() error {

	conn, err := grpc.Dial(
		s.config.GrpcDSN,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	s.grpcClient = pb.NewHelloServiceClient(conn)

	go s.grpcTestRequest()

	return nil
}

func (s *Server) grpcTestRequest() {

	resp, err := s.grpcClient.Say(context.Background(), &pb.SayRequest{
		Message: "Hello from server",
	})
	if err != nil {
		s.logger.Errorf("grpc send response error: %v", err)
		return
	}

	s.logger.Infof("grpc response: %v", resp.GetMessage())
}
