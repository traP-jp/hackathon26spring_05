package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

type tag struct {
	Label    string  `json:"label"`
	Affinity string  `json:"affinity"`
	Strength float64 `json:"strength"`
}

type meResponse struct {
	Username     string         `json:"username"`
	IconFileID   string         `json:"iconFileId"`
	Major        string         `json:"major"`
	Affiliations []string       `json:"affiliations"`
	Hometown     string         `json:"hometown"`
	Tags         map[string]tag `json:"tags"`
	Bio          string         `json:"bio"`
}

// GET /api/me
func (h *handler) getMe(c echo.Context) error {
	if !h.loginUserRetriever.IsUserLoggedIn() {
		return unauthorized(c)
	}

	username, err := h.loginUserRetriever.GetLoginUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to get login user"})
	}

	user, err := h.repository.FindByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to load user"})
	}
	if user == nil {
		return notFound(c)
	}

	return c.JSON(http.StatusOK, toMeResponse(*user))
}

func toMeResponse(user domain.User) meResponse { //FindByUsernameで取得したデータをuserをjsonにして返す用
	tags := make(map[string]tag, len(user.Tags))
	for name, userTag := range user.Tags {
		tags[name] = tag{
			Label:    userTag.Label,
			Affinity: string(userTag.Affinity),
			Strength: userTag.Strength,
		}
	}

	return meResponse{
		Username:     user.Username,
		IconFileID:   user.IconFileID,
		Major:        user.Major,
		Affiliations: user.Affiliations,
		Hometown:     user.Hometown,
		Tags:         tags,
		Bio:          user.Bio,
	}
}

// PUT /api/me
func (h *handler) updateMe(c echo.Context) error {
	return unauthorized(c)
}

type userSummaryResponse struct {
	Username string `json:"username"`
}

func toUserSummaryResponses(users []domain.UserSummary) ([]userSummaryResponse, error) {
	result := make([]userSummaryResponse, len(users))
	for i, user := range users {
		if user.Username == "" {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "invalid user summary")
		}
		result[i] = userSummaryResponse{Username: user.Username}
	}

	return result, nil
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

	result, err := toUserSummaryResponses(users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to validate liked users"})
	}

	return c.JSON(http.StatusOK, result)
}

// POST /api/me/likes
func (h *handler) likeUser(c echo.Context) error {
	return unauthorized(c)
}

// GET /api/me/liked-by
func (h *handler) listUsersWhoLikedMe(c echo.Context) error {
	if !h.loginUserRetriever.IsUserLoggedIn() {
		return unauthorized(c)
	}

	username, err := h.loginUserRetriever.GetLoginUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to get login user"})
	}

	users, err := h.repository.ListUsersWhoLiked(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to list users who liked me"})
	}

	result, err := toUserSummaryResponses(users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to validate liked-by users"})
	}

	return c.JSON(http.StatusOK, result)
}

// POST /api/me/nopes
func (h *handler) nopeUser(c echo.Context) error {
	return unauthorized(c)
}
