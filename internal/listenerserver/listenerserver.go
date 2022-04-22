package listenerserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tkcrm/modules/logger"
)

type ListenerServer struct {
	logger     logger.Logger
	serverName string
	listener   net.Listener
}

func New(l logger.Logger, serverName string, port uint16) (*ListenerServer, error) {

	loggerExtendedFields := []interface{}{"protocol_type", serverName}
	l = l.With(loggerExtendedFields...)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	l.Infof("server %s started", serverName)

	return &ListenerServer{
		logger:     l,
		serverName: serverName,
		listener:   listener,
	}, nil
}

func (s *ListenerServer) StopServer() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *ListenerServer) AcceptConnections(ctx context.Context) error {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				s.logger.Errorf("Error accepting connection %v", err)
				continue
			}
		}

		s.logger.Debugf("Accepted connection from %v", connection.RemoteAddr())

		go func(connection net.Conn) {
			if err := s.handleConnection(ctx, connection); err != nil {
				s.logger.Errorf("handleConnection error: %v", err)
			}
		}(connection)
	}
}

func (s *ListenerServer) handleConnection(ctx context.Context, conn net.Conn) error {

	defer func() {
		conn.Close()
	}()

	res := fmt.Sprintf(
		"Hi %s. You are connected to %s, %s",
		conn.RemoteAddr().String(),
		s.serverName,
		time.Now(),
	)
	_, err := conn.Write([]byte(res))
	return err
}
