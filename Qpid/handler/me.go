package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
)

type tag struct {
	Label    optional.Option[string] `json:"label"`
	Affinity string                  `json:"affinity"`
	Strength float64                 `json:"strength"`
}

type meResponse struct {
	Username     string                  `json:"username"`
	IconFileID   optional.Option[string] `json:"iconFileId"`
	Major        optional.Option[string] `json:"major"`
	Affiliations []string                `json:"affiliations"`
	Hometown     optional.Option[string] `json:"hometown"`
	Tags         map[string]tag          `json:"tags"`
	Bio          optional.Option[string] `json:"bio"`
}

// GET /api/me
func (h *handler) getMe(c echo.Context) error {
	username := middleware.GetUsername(c)
	if username == nil {
		return unauthorized(c)
	}

	user, err := h.repository.FindByUsername(*username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to load user"})
	}
	if user == nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "user not found"})
	}

	return c.JSON(http.StatusOK, toMeResponse(*user))
}

func toMeResponse(user domain.User) *meResponse { //FindByUsernameで取得したデータをuserをjsonにして返す用
	tags := make(map[string]tag, len(user.Tags))
	for name, userTag := range user.Tags {
		tags[name] = tag{
			Label:    userTag.Label,
			Affinity: string(userTag.Affinity),
			Strength: userTag.Strength,
		}
	}

	return &meResponse{
		Username:     user.Username,
		IconFileID:   user.IconFileID,
		Major:        user.Major,
		Affiliations: user.Affiliations,
		Hometown:     user.Hometown,
		Tags:         tags,
		Bio:          user.Bio,
	}
}

type meRequest struct {
	IconFileID   string         `json:"iconFileId"`
	Major        string         `json:"major"`
	Affiliations []string       `json:"affiliations"`
	Hometown     string         `json:"hometown"`
	Tags         map[string]tag `json:"tags"`
	Bio          string         `json:"bio"`
}

// PUT /api/me
func (h *handler) updateMe(c echo.Context) error {
	loginUserRetriever := middleware.GetLoginUserRetriever(c)

	if !loginUserRetriever.IsUserLoggedIn() {
		return unauthorized(c)
	}

	username, err := loginUserRetriever.GetLoginUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to get login user"})
	}

	data := &meRequest{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
	}

	tags := make(map[string]domain.Tag, len(data.Tags))
	for name, t := range data.Tags {
		aff := domain.TagAffinity(t.Affinity)
		switch aff {
		case domain.TagAffinityPositive, domain.TagAffinityNeutral, domain.TagAffinityNegative:
		default:
			return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid tag affinity"})
		}
		tags[name] = domain.Tag{
			Label:    t.Label,
			Affinity: aff,
			Strength: t.Strength,
		}
	}
	userData := domain.User{
		Username:     username,
		IconFileID:   data.IconFileID,
		Major:        data.Major,
		Affiliations: data.Affiliations,
		Hometown:     data.Hometown,
		Tags:         tags,
		Bio:          data.Bio,
	}

	if err := h.repository.UpdateProfile(username, userData); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to update user"})
	}

	return c.JSON(http.StatusOK, toMeResponse(userData))
}

type userSummaryResponse struct {
	Username string `json:"username"`
}

func toUserSummaryResponses(users []domain.UserSummary) ([]userSummaryResponse, error) {
	result := make([]userSummaryResponse, len(users))
	for i, user := range users {
		if user.Username == "" {
			return nil, errors.New("invalid user summary")
		}
		result[i] = userSummaryResponse{Username: user.Username}
	}

	return result, nil
}

// GET /api/me/likes
func (h *handler) listMyLikes(c echo.Context) error {
	username := middleware.GetUsername(c)
	if username == nil {
		return unauthorized(c)
	}

	users, err := h.repository.ListLikedUsers(*username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to list liked users"})
	}

	result, err := toUserSummaryResponses(users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to validate liked users"})
	}
	return c.JSON(http.StatusOK, result)
}

type userActionRequest struct {
	Username string `json:"username"`
}

// POST /api/me/likes
func (h *handler) likeUser(c echo.Context) error {
	username := middleware.GetUsername(c)
	if username == nil {
		return unauthorized(c)
	}

	toUser := &userActionRequest{}
	if err := c.Bind(toUser); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
	}

	if toUser.Username == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "username is required"})
	}
	if toUser.Username == *username {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "cannot like yourself"})
	}

	isExist, err := h.repository.Exists(toUser.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to check user existence"})
	}
	if !isExist {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "user does not exist"})
	}

	if err = h.repository.Like(*username, toUser.Username); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to like user"})
	}

	return c.NoContent(http.StatusNoContent)
}

// GET /api/me/liked-by
func (h *handler) listUsersWhoLikedMe(c echo.Context) error {
	loginUserRetriever := middleware.GetLoginUserRetriever(c)

	if !loginUserRetriever.IsUserLoggedIn() {
		return unauthorized(c)
	}

	username, err := loginUserRetriever.GetLoginUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to get login user"})
	}

	toUser := &userActionRequest{}
	if err := c.Bind(toUser); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
	}

	if toUser.Username == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "username is required"})
	}
	if toUser.Username == username {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "cannot like yourself"})
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
	loginUserRetriever := middleware.GetLoginUserRetriever(c)

	if !loginUserRetriever.IsUserLoggedIn() {
		return unauthorized(c)
	}

	username, err := loginUserRetriever.GetLoginUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to get login user"})
	}

	toUser := &userActionRequest{}
	if err := c.Bind(toUser); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
	}

	if toUser.Username == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "username is required"})
	}
	if toUser.Username == username {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "cannot nope yourself"})
	}

	isExist, err := h.repository.Exists(toUser.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to check user existence"})
	}
	if !isExist {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "user does not exist"})
	}

	if err = h.repository.Nope(username, toUser.Username); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to nope user"})
	}

	return c.NoContent(http.StatusNoContent)
}
