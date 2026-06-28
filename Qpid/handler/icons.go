package handler

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
)

// PATCH /api/me/icon
func (h *handler) updateMyIcon(c *echo.Context) error {
	username := middleware.GetUsername(c)

	//ユーザー認証
	if username == nil {
		return unauthorized(c)
	}

	file, err := c.FormFile("icon")
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "icon file is required"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "cannot open icon file"})
	}
	defer func() { _ = src.Close() }()

	blob, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to load icon file"})
	}

	mimeType := domain.IconMimeType(http.DetectContentType(blob))
	if !isAllowedIconMimeType(mimeType) {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "unsupported icon mime type"})
	}

	if err := h.repository.SaveIcon(*username, domain.Icon{
		Blob:     blob,
		MimeType: mimeType,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to save icon"})
	}

	return c.NoContent(http.StatusNoContent)
}

// DELETE /api/me/icon
func (h *handler) deleteMyIcon(c *echo.Context) error {
	username := middleware.GetUsername(c)
	if username == nil {
		return unauthorized(c)
	}

	if err := h.repository.DeleteIcon(*username); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to delete icon"})
	}

	return (*c).NoContent(http.StatusNoContent)
}

// GET /api/users/:id/icon
func (h *handler) getUserIcon(c *echo.Context) error {
	username, err := echo.PathParam[string](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid path parameter"})
	}

	icon, err := h.repository.FindIconByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to find icon file"})
	}
	if icon == nil {
		return notFound(c)
	}

	return (*c).Blob(http.StatusOK, string(icon.MimeType), icon.Blob)
}

func isAllowedIconMimeType(mimeType domain.IconMimeType) bool {
	return mimeType == domain.IconMimeTypePNG ||
		mimeType == domain.IconMimeTypeJPEG ||
		mimeType == domain.IconMimeTypeWebp
}
