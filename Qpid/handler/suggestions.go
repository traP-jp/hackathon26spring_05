package handler

import (
	"errors"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
)

const defaultSuggestionLimit = 10

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
		if math.IsNaN(suggestion.Similarity) || suggestion.Similarity < 0 || suggestion.Similarity > 1 {
			return nil, errors.New("invalid suggestion")
		}

		result[i] = suggestionResponse{
			Username:   suggestion.Username,
			Similarity: suggestion.Similarity,
		}
	}

	return result, nil
}

// GET /api/suggestions
func (h *handler) listSuggestions(c echo.Context) error {
	loginUserRetriever := middleware.GetLoginUserRetriever(c)

	if !loginUserRetriever.IsUserLoggedIn() {
		return unauthorized(c)
	}

	username, err := loginUserRetriever.GetLoginUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to get login user"})
	}

	suggestions, err := h.repository.ListSuggestions(username, defaultSuggestionLimit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to list suggestions"})
	}

	result, err := toSuggestionResponses(suggestions)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to validate suggestions"})
	}

	return c.JSON(http.StatusOK, result)
}
