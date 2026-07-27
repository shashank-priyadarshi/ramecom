package server

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/rajabhishekmaurya/ecom/internal/middleware"
	"github.com/rajabhishekmaurya/ecom/libs/config"
)

type register interface {
	Register(*echo.Echo)
}

type Server struct {
	cfg *config.Config
	e   *echo.Echo
}

func New(cfg *config.Config, r register) (*Server, error) {
	e := echo.New()

	e.Use(middleware.RequestID)
	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Recover())
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	r.Register(e)

	return &Server{
		cfg: cfg,
		e:   e,
	}, nil
}

func (s *Server) Start() error {
	log.Printf("%s started on :%s", s.cfg.Srv.Name, s.cfg.Srv.Port)
	return s.e.Start(":" + s.cfg.Srv.Port)
}

func (s *Server) Shutdown() error {
	return s.e.Shutdown(context.Background())
}
