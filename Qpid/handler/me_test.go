package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository/mock"
)

type stubLoginUserRetriever struct {
	loggedIn bool
	username string
	err      error
}

func (s stubLoginUserRetriever) IsUserLoggedIn() bool {
	return s.loggedIn
}

func (s stubLoginUserRetriever) GetLoginUser() (string, error) {
	return s.username, s.err
}

func TestListUsersWhoLikedMeUnauthorized(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/me/liked-by", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handler{
		loginUserRetriever: stubLoginUserRetriever{loggedIn: false},
	}

	if err := h.listUsersWhoLikedMe(c); err != nil {
		t.Fatalf("listUsersWhoLikedMe() returned error: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestListUsersWhoLikedMeOK(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/me/liked-by", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handler{
		repository:        mock.NewMockRepository(),
		loginUserRetriever: stubLoginUserRetriever{loggedIn: true, username: "test-user"},
	}

	if err := h.listUsersWhoLikedMe(c); err != nil {
		t.Fatalf("listUsersWhoLikedMe() returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var got []userSummaryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	want := []userSummaryResponse{{Username: "liked-by-user"}}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("response body = %+v, want %+v", got, want)
	}
}
