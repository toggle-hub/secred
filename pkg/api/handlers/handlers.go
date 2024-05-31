package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HttpError struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Healthy bool `json:"healthy"`
}

func HealthZHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, HealthResponse{true})
	}
}
