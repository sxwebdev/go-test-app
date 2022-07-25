package config

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/tkcrm/modules/broker/mqttconn"
	"github.com/tkcrm/modules/broker/natsconn"
	"github.com/tkcrm/modules/db/bunconn"
	"github.com/tkcrm/modules/logger"
)

type Config struct {
	AppName    string `default:"go-test-app"`
	LogLevel   string `default:"info"`
	Env        string
	AppType    string
	DB         bunconn.Config
	Nats       natsconn.Config
	Mqtt       mqttconn.Config
	ApiDSN     string
	GrpcDSN    string
	TCPServers map[string]uint16
}

func (c *Config) Validate() error {

	// PostgreSQL
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("DB_%s", err)
	}

	// Nats
	if err := c.Nats.Validate(); err != nil {
		return fmt.Errorf("NATS_%s", err)
	}

	// Mqtt
	if err := c.Mqtt.Validate(); err != nil {
		return fmt.Errorf("MQTT_%s", err)
	}

	return validation.ValidateStruct(
		c,
		validation.Field(&c.AppName, validation.Required),
		validation.Field(&c.Env, validation.Required, validation.In("dev", "stage", "prod")),
		validation.Field(&c.LogLevel, validation.Required, validation.In(logger.GetAllLevels()...)),
		validation.Field(&c.AppType, validation.Required, validation.In("server", "service")),
		validation.Field(&c.ApiDSN, validation.Required),
		validation.Field(&c.GrpcDSN, validation.Required),
	)
}
