package handler

import "github.com/labstack/echo/v4"

// GET /api/users/:id
func (h *handler) getUser(c echo.Context) error {
	return notFound(c)
}
