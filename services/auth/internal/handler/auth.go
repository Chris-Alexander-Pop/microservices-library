package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Register(e *echo.Echo) {
	e.POST("/auth/signup", h.Signup)
	e.POST("/auth/login", h.Login)
}

func (h *Handler) Signup(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"msg": "signed up"})
}

func (h *Handler) Login(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"msg": "logged in", "token": "fake-jwt"})
}
