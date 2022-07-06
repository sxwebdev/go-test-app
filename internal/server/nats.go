package server

import (
	"fmt"
	"math/rand"

	"github.com/nats-io/nats.go"
	"github.com/tkcrm/modules/broker/natsconn"
)

func (s *Server) natsConnect() error {

	nc, err := natsconn.New(s.logger, s.config.Nats, s.config.AppName)
	if err != nil {
		return err
	}

	s.nats = nc

	go s.natsSubscribe()

	return nil
}

func (s *Server) natsSubscribe() {
	s.nats.Subscribe("go-test-app/test-nats", func(m *nats.Msg) {

		res := fmt.Sprintf(
			"Hello from server! Received: %s. %d",
			string(m.Data),
			rand.Int(),
		)

		s.nats.Publish(m.Reply, []byte(res))
		s.logger.Debugf("NATS received a message: %s", string(m.Data))
	})
}
