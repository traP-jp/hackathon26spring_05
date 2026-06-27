package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository/mock"
)

type suggestionHandlerTestRepository struct {
	*mock.MockRepository
	listSuggestions func(username string, limit int) ([]domain.Suggestion, error)
}

func (r *suggestionHandlerTestRepository) ListSuggestions(username string, limit int) ([]domain.Suggestion, error) {
	return r.listSuggestions(username, limit)
}

type testLoginUserRetriever struct {
	username string
}

func (r *testLoginUserRetriever) IsUserLoggedIn() bool {
	return r.username != ""
}

func (r *testLoginUserRetriever) GetLoginUser() (string, error) {
	return r.username, nil
}

func TestListSuggestions(t *testing.T) {
	t.Parallel()

	repo := &suggestionHandlerTestRepository{
		MockRepository: mock.NewMockRepository(),
		listSuggestions: func(username string, limit int) ([]domain.Suggestion, error) {
			if username != "test-user" {
				t.Fatalf("unexpected username: %s", username)
			}
			if limit != defaultSuggestionLimit {
				t.Fatalf("unexpected limit: %d", limit)
			}

			return []domain.Suggestion{
				{Username: "suggested-user", Similarity: 0.5},
			}, nil
		},
	}
	h := &handler{repository: repo}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/suggestions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("loginUserRetriever", &testLoginUserRetriever{username: "test-user"})

	if err := h.listSuggestions(c); err != nil {
		t.Fatalf("listSuggestions returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", rec.Code, http.StatusOK)
	}

	var response []suggestionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("unexpected response length: %d", len(response))
	}
	if response[0].Username != "suggested-user" {
		t.Fatalf("unexpected username: %s", response[0].Username)
	}
	if response[0].Similarity != 0.5 {
		t.Fatalf("unexpected similarity: %v", response[0].Similarity)
	}
}

func TestListSuggestionsReturnsValidationErrorForInvalidSuggestion(t *testing.T) {
	t.Parallel()

	repo := &suggestionHandlerTestRepository{
		MockRepository: mock.NewMockRepository(),
		listSuggestions: func(username string, limit int) ([]domain.Suggestion, error) {
			return []domain.Suggestion{
				{Username: "", Similarity: 0.5},
			}, nil
		},
	}
	h := &handler{repository: repo}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/suggestions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("loginUserRetriever", &testLoginUserRetriever{username: "test-user"})

	if err := h.listSuggestions(c); err != nil {
		t.Fatalf("listSuggestions returned error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: got %d want %d", rec.Code, http.StatusInternalServerError)
	}
}
