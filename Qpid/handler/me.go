package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /api/me
func (h *handler) getMe(c echo.Context) error {
	return unauthorized(c)
}

// PUT /api/me
func (h *handler) updateMe(c echo.Context) error {
	return unauthorized(c)
}

type userSummaryResponse struct {
	Username string `json:"username"`
}

// GET /api/me/likes
func (h *handler) listMyLikes(c echo.Context) error {
	if !h.loginUserRetriever.IsUserLoggedIn() {
		return unauthorized(c)
	}

	username, err := h.loginUserRetriever.GetLoginUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to get login user"})
	}

	users, err := h.repository.ListLikedUsers(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to list liked users"})
	}

	result := make([]userSummaryResponse, len(users))
	for i, user := range users {
		result[i] = userSummaryResponse{Username: user.Username}
	}

	return c.JSON(http.StatusOK, result)
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
