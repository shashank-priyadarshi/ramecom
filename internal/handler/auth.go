package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rajabhishekmaurya/ecom/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (a *AuthHandler) Register(e *echo.Echo) {
	e.POST("/login", a.login)
}

func (a *AuthHandler) login(c echo.Context) error {
	var req service.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	resp, err := a.service.Login(
		c.Request().Context(),
		&req,
	)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}
