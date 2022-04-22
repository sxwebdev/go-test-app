package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func (s *Server) newApiServer() {

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		CaseSensitive:         true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			loggerExtendedFields := []interface{}{"status_code", code}
			errText := fmt.Sprintf("%s %s %s - %s",
				ctx.Method(),
				ctx.Path(),
				err.Error(),
				ctx.IP())

			switch {
			case code >= 500:
				s.logger.With(loggerExtendedFields...).Error(errText)
			case code >= 400:
				s.logger.With(loggerExtendedFields...).Warn(errText)
			}

			return ctx.Status(code).JSON(newError(code, err.Error()))
		},
	})

	app.Use(etag.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,PUT,POST,DELETE,OPTIONS",
		ExposeHeaders: "Content-Type,Authorization,Accept",
	}))
	app.Use(requestid.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"test": "value",
		})
	})

	s.fiber = app
}

func (s *Server) ApiStart() error {
	go func() {
		time.Sleep(time.Millisecond * 50)
		s.logger.Infof("api server start successfully on port %s", s.config.APIServerPort)
	}()
	return s.fiber.Listen(fmt.Sprintf(":%s", s.config.APIServerPort))
}
