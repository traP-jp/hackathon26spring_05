package handler

import "github.com/labstack/echo/v4"

// GET /api/me
func (h *handler) getMe(c echo.Context) error {
	return unauthorized(c)
}

// PUT /api/me
func (h *handler) updateMe(c echo.Context) error {
	return unauthorized(c)
}

// GET /api/me/likes
func (h *handler) listMyLikes(c echo.Context) error {
	return unauthorized(c)
}

// POST /api/me/likes
func (h *handler) likeUser(c echo.Context) error {
	return unauthorized(c)
}

// GET /api/me/liked-by
func (h *handler) listUsersWhoLikedMe(c echo.Context) error {
	return unauthorized(c)
}

// POST /api/me/nopes
func (h *handler) nopeUser(c echo.Context) error {
	return unauthorized(c)
}
