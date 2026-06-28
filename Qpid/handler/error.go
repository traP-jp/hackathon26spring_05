package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type errorResponse struct {
	Message string `json:"message"`
}

func notImplemented(c *echo.Context) error {
	return c.JSON(http.StatusBadRequest, errorResponse{Message: "not implemented"})
}

func unauthorized(c *echo.Context) error {
	return c.JSON(http.StatusUnauthorized, errorResponse{Message: "authentication is required"})
}

func notFound(c *echo.Context) error {
	return c.JSON(http.StatusNotFound, errorResponse{Message: "not found"})
}
