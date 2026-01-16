package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	v1 := e.Group("/api/v1")
	v1.GET("/hello", HandleHello)
}

func HandleHello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello World",
	})
}
