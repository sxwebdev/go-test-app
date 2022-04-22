package server

import (
	"fmt"
	"math/rand"

	"github.com/nats-io/nats.go"
	"github.com/tkcrm/modules/cfg"
)

func (s *Server) natsConnect() error {

	url := cfg.GetNATSURL(s.config.NATSHost, s.config.NATSPort)
	s.logger.Infof("try to connect to nats: %s", url)

	nc, err := nats.Connect(
		url,
		cfg.GetNATSOpts(
			s.logger,
			s.config.APPName,
			s.config.NATSUser,
			s.config.NATSPass,
			s.config.NATSToken,
		)...,
	)
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
