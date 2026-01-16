package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port         string        `env:"PORT" env-default:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"10s"`
}

type Server struct {
	echo *echo.Echo
	cfg  Config
	log  *slog.Logger
}

func New(cfg Config, log *slog.Logger) *Server {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			status := c.Response().Status
			log.Info("request",
				"method", c.Request().Method,
				"uri", c.Request().RequestURI,
				"status", status,
				"latency", time.Since(start),
			)
			return err
		}
	})

	return &Server{echo: e, cfg: cfg, log: log}
}

func (s *Server) Start() error {
	s.log.Info("starting http server", "port", s.cfg.Port)
	return s.echo.Start(":" + s.cfg.Port)
}

func (s *Server) Echo() *echo.Echo {
	return s.echo
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
