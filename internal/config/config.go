package config

import (
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sakirsensoy/genv"
	"github.com/sakirsensoy/genv/dotenv"
	"github.com/tkcrm/modules/utils"
)

type Config struct {
	APPName       string
	APPMSName     string
	APPVersion    string
	APPHost       string
	ENV           string
	AppType       string
	DBUser        string
	DBPass        string
	DBHost        string
	DBPort        string
	DBName        string
	APIServerPort string
	NATSHost      string
	NATSPort      string
	NATSUser      string
	NATSPass      string
	NATSToken     string
	MQTTHost      string
	MQTTPort      string
	MQTTUser      string
	MQTTPass      string
	GRPCHost      string
	GRPCPort      string
	TCPServers    map[string]uint16
}

func New() *Config {

	env := utils.GetDefaultString(os.Getenv("ENV"), "dev")

	if env == "dev" {
		dotenv.Load()
	}

	return &Config{
		// App
		APPName:    utils.GetDefaultString(genv.Key("APP_NAME").String(), "UNDEFINED_APP_NAME"),
		APPMSName:  genv.Key("APP_MS_NAME").String(),
		APPVersion: utils.GetDefaultString(genv.Key("APP_VERSION").String(), "UNDEFINED_APP_VESION"),
		APPHost:    genv.Key("APP_HOST").String(),
		ENV:        env,
		AppType:    genv.Key("APP_TYPE").String(),

		// Postgres
		DBUser: genv.Key("DB_USER").String(),
		DBPass: genv.Key("DB_PASSWD").String(),
		DBHost: genv.Key("DB_HOST").String(),
		DBPort: genv.Key("DB_PORT").String(),
		DBName: genv.Key("DB_NAME").String(),

		// API Server
		APIServerPort: genv.Key("API_PORT").String(),

		// NATS
		NATSHost:  genv.Key("NATS_HOST").String(),
		NATSPort:  genv.Key("NATS_PORT").String(),
		NATSUser:  genv.Key("NATS_USER").String(),
		NATSPass:  genv.Key("NATS_PASS").String(),
		NATSToken: genv.Key("NATS_TOKEN").String(),

		// MQTT
		MQTTHost: genv.Key("MQTT_HOST").String(),
		MQTTPort: genv.Key("MQTT_PORT").String(),
		MQTTUser: genv.Key("MQTT_USER").String(),
		MQTTPass: genv.Key("MQTT_PASS").String(),

		// GRPC
		GRPCHost: genv.Key("GRPC_HOST").String(),
		GRPCPort: genv.Key("GRPC_PORT").String(),

		// TCP Servers
		TCPServers: map[string]uint16{
			"server_1": 35100,
			"server_2": 35101,
			"server_3": 35102,
		},
	}
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.APPName, validation.Required),
		validation.Field(&c.APPMSName, validation.Required),
		validation.Field(&c.APPHost, validation.Required),
		validation.Field(&c.ENV, validation.Required, validation.In("dev", "stage", "prod")),
		validation.Field(&c.AppType, validation.Required, validation.In("server", "service")),
		validation.Field(&c.DBHost, validation.Required),
		validation.Field(&c.DBPort, validation.Required),
		validation.Field(&c.DBUser, validation.Required),
		validation.Field(&c.DBPass, validation.Required),
		validation.Field(&c.DBName, validation.Required),
		validation.Field(&c.APIServerPort, validation.Required),
		validation.Field(&c.NATSHost, validation.Required),
		validation.Field(&c.NATSPort, validation.Required),
		validation.Field(&c.MQTTHost, validation.Required),
		validation.Field(&c.MQTTPort, validation.Required),
		validation.Field(&c.GRPCHost, validation.Required),
		validation.Field(&c.GRPCPort, validation.Required),
	)
}
