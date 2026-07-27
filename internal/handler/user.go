package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rajabhishekmaurya/ecom/internal/model"
	"github.com/rajabhishekmaurya/ecom/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u *UserHandler) Register(e *echo.Echo) {
	e.POST("/users", u.create)
}

func (u *UserHandler) create(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := u.service.Register(c.Request().Context(), &user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	user.Password = "*****"

	return c.JSON(http.StatusCreated, user)
}
