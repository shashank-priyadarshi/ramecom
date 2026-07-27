package handler

import (
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/rajabhishekmaurya/ecom/libs/config"
)

type GatewayHandler struct {
	cfg *config.Config
}

func NewGatewayHandler(cfg *config.Config) *GatewayHandler {
	return &GatewayHandler{
		cfg: cfg,
	}
}

func (g *GatewayHandler) Register(e *echo.Echo) {
	auth := e.Group("/auth")
	auth.Any("/*", reverseProxy(g.cfg.Service.Auth))

	users := e.Group("/users")
	users.Any("/*", reverseProxy(g.cfg.Service.User))

	orders := e.Group("/orders")
	orders.Any("", reverseProxy(g.cfg.Service.Order))
	orders.Any("/*", reverseProxy(g.cfg.Service.Order))

}

func reverseProxy(target string) echo.HandlerFunc {
	u, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(u)

	return func(c echo.Context) error {
		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
