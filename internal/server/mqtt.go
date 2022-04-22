package server

import (
	"context"
	"fmt"
	"net"

	"github.com/eclipse/paho.golang/paho"
)

func (s *Server) mqttConnect() error {

	url := fmt.Sprintf("%s:%s", s.config.MQTTHost, s.config.MQTTPort)
	s.logger.Infof("try to connect to mqtt broker: %s", url)

	conn, err := net.Dial("tcp", url)
	if err != nil {
		return fmt.Errorf("mqtt connect error: %v. Error connect to tcp: %s", err, url)
	}

	pc := paho.NewClient(paho.ClientConfig{
		Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
			s.logger.Infof("MQTT received: %s", string(m.Payload))
			if s.mqtt != nil {
				if _, err := s.mqtt.Publish(context.Background(), &paho.Publish{
					QoS:     1,
					Topic:   "go-test-server/mqtt/response",
					Payload: []byte("Awesome! Received:\n" + string(m.Payload)),
				}); err != nil {
					s.logger.Error(err)
				}
			}
		}),
		Conn: conn,
	})

	opts := &paho.Connect{
		KeepAlive:  30,
		CleanStart: true,
		ClientID:   "go-test-server",
		Username:   s.config.MQTTUser,
		Password:   []byte(s.config.MQTTPass),
	}

	if s.config.MQTTUser != "" {
		opts.UsernameFlag = true
	}

	if s.config.MQTTPass != "" {
		opts.PasswordFlag = true
	}

	res, err := pc.Connect(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("mqtt connect error: %v", err)
	}

	if res.ReasonCode != 0 {
		return fmt.Errorf("failed to connect with reason: %d - %s", res.ReasonCode, res.Properties.ReasonString)
	}

	s.logger.Info("Connected to MQTT Broker successfully")

	s.mqtt = pc

	go s.mqttSubscribe()

	return nil
}

func (s *Server) mqttSubscribe() {
	if _, err := s.mqtt.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			"go-test-app/mqtt": {QoS: 0, NoLocal: false},
		},
	}); err != nil {
		s.logger.Fatal(err)
	}
}
