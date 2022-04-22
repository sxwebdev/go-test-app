package main

import (
	"os"

	"github.com/sakirsensoy/genv"
	"github.com/sakirsensoy/genv/dotenv"
	"github.com/sxwebdev/go-test-app/internal/server"
	"github.com/sxwebdev/go-test-app/internal/service"
	"github.com/tkcrm/modules/logger"
	"github.com/tkcrm/modules/utils"
)

func main() {

	if utils.GetDefaultString(os.Getenv("ENV"), "dev") == "dev" {
		dotenv.Load()
	}

	l := logger.DefaultLogger(
		utils.GetDefaultString(genv.Key("LOG_LEVEL").String(), "info"),
		utils.GetDefaultString(genv.Key("APP_MS_NAME").String(), "undefined_app_ms_name"),
	)

	appType := os.Getenv("APP_TYPE")

	if appType == "" {
		l.Fatal("undefined APP_TYPE")
	}

	switch appType {
	case "server":
		if err := server.Start(l); err != nil {
			l.Errorf("INIT APP ERROR: %v", err)
		}
	case "service":
		if err := service.Start(l); err != nil {
			l.Errorf("INIT APP ERROR: %v", err)
		}
	default:
		l.Fatalf("unavailable APP_TYPE: %s", appType)
	}

}
