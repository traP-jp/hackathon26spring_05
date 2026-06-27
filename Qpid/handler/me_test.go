package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
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

type stubListUsersWhoLikedRepository struct {
	*mock.MockRepository
	users []domain.UserSummary
}

func (r stubListUsersWhoLikedRepository) ListUsersWhoLiked(username string) ([]domain.UserSummary, error) {
	return r.users, nil
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

	tests := []struct {
		name  string
		users []domain.UserSummary
		want  []userSummaryResponse
	}{
		{
			name:  "single user",
			users: []domain.UserSummary{{Username: "liked-by-user"}},
			want:  []userSummaryResponse{{Username: "liked-by-user"}},
		},
		{
			name:  "empty users",
			users: []domain.UserSummary{},
			want:  []userSummaryResponse{},
		},
		{
			name: "multiple users",
			users: []domain.UserSummary{
				{Username: "liked-by-user-1"},
				{Username: "liked-by-user-2"},
			},
			want: []userSummaryResponse{
				{Username: "liked-by-user-1"},
				{Username: "liked-by-user-2"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/me/liked-by", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &handler{
				repository: stubListUsersWhoLikedRepository{
					MockRepository: mock.NewMockRepository(),
					users:          tt.users,
				},
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("response body = %+v, want %+v", got, tt.want)
			}
		})
	}
}
