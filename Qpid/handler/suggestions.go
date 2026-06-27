package handler

import "github.com/labstack/echo/v4"

// GET /api/suggestions
func (h *handler) listSuggestions(c echo.Context) error {
	return unauthorized(c)
}
