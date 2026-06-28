package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
)

type suggestionResponse struct {
	Username   string  `json:"username"`
	Similarity float64 `json:"similarity"`
}

func toSuggestionResponses(suggestions []domain.Suggestion) ([]suggestionResponse, error) {
	result := make([]suggestionResponse, len(suggestions))
	for i, suggestion := range suggestions {
		if suggestion.Username == "" {
			return nil, errors.New("invalid suggestion")
		}
		if !(suggestion.Similarity >= 0 && suggestion.Similarity <= 1) {
			return nil, errors.New("invalid suggestion similarity")
		}
		result[i] = suggestionResponse{
			Username:   suggestion.Username,
			Similarity: suggestion.Similarity,
		}
	}

	return result, nil
}

// GET /api/suggestions
func (h *handler) listSuggestions(c *echo.Context) error {
	username := middleware.GetUsername(c)
	if username == nil {
		return unauthorized(c)
	}

	suggestions, err := h.repository.ListSuggestions(*username, 20)
	if err != nil {
		c.Logger().Error(
			"failed to list suggestions",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to list suggestions"})
	}

	result, err := toSuggestionResponses(suggestions)
	if err != nil {
		c.Logger().Error(
			"failed to convert suggestions",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to validate suggestions"})
	}

	return c.JSON(http.StatusOK, result)
}
