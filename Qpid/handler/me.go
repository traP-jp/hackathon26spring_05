package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
)

// GET /api/me
func (h *handler) getMe(c *echo.Context) error {
	username := middleware.GetUsername(c)
	if username == nil {
		return unauthorized(c)
	}

	user, err := h.repository.FindUserByUsername(*username)
	if err != nil {
		c.Logger().Error(
			"failed to find user",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to load user"})
	}
	if user == nil {
		c.Logger().Error(
			"authenticated user was not found",
			slog.String("username", *username),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "user not found"})
	}

	if uuid := middleware.GetUserUUID(c); uuid != nil {
		user.Affiliations = h.fetchAffiliations(c, *uuid)
	} else {
		user.Affiliations = []domain.UserAffiliation{}
	}

	return c.JSON(http.StatusOK, toUserResponse(*user))
}

type updateMeRequest struct {
	HasIcon       bool                                 `json:"hasIcon"`
	Major         optional.Option[string]              `json:"major"`
	Hometown      optional.Option[string]              `json:"hometown"`
	Tags          []string                             `json:"tags"`
	Technologies  []string                             `json:"technologies"`
	Bio           optional.Option[string]              `json:"bio"`
	FavoriteTopic optional.Option[updateTopicAndValue] `json:"favoriteTopic"`
	DislikedTopic optional.Option[updateTopicAndValue] `json:"dislikedTopic"`
}

type updateTopicAndValue struct {
	Topic string `json:"topic"`
	Value string `json:"value"`
}

func toDomainTopicAndValue(updateTopicAndValue optional.Option[updateTopicAndValue]) optional.Option[domain.TopicAndValue] {
	if updateTopicAndValue.IsNone() {
		return optional.None[domain.TopicAndValue]()
	}
	updateTopicAndValueValue := updateTopicAndValue.Unwrap()
	return optional.Some(domain.TopicAndValue{
		Topic: updateTopicAndValueValue.Topic,
		Value: updateTopicAndValueValue.Value,
	})
}

// PUT /api/me
func (h *handler) updateMe(c *echo.Context) error {
	username := middleware.GetUsername(c)
	if username == nil {
		return unauthorized(c)
	}

	data := &updateMeRequest{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
	}

	userData := domain.User{
		Username:      *username,
		HasIcon:       data.HasIcon,
		Major:         data.Major,
		Hometown:      data.Hometown,
		Tags:          data.Tags,
		Technologies:  data.Technologies,
		Bio:           data.Bio,
		FavoriteTopic: toDomainTopicAndValue(data.FavoriteTopic),
		DislikedTopic: toDomainTopicAndValue(data.DislikedTopic),
	}

	if err := h.repository.UpdateUser(*username, userData); err != nil {
		c.Logger().Error(
			"failed to update user",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to update user"})
	}

	if uuid := middleware.GetUserUUID(c); uuid != nil {
		userData.Affiliations = h.fetchAffiliations(c, *uuid)
	} else {
		userData.Affiliations = []domain.UserAffiliation{}
	}

	return c.JSON(http.StatusOK, toUserResponse(userData))
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
func (h *handler) listMyLikes(c *echo.Context) error {
	username := middleware.GetUsername(c)
	if username == nil {
		return unauthorized(c)
	}

	users, err := h.repository.ListLikedUsers(*username)
	if err != nil {
		c.Logger().Error(
			"failed to list liked users",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to list liked users"})
	}

	result, err := toUserSummaryResponses(users)
	if err != nil {
		c.Logger().Error(
			"failed to convert liked users",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to validate liked users"})
	}
	return c.JSON(http.StatusOK, result)
}

type userActionRequest struct {
	Username string `json:"username"`
}

// POST /api/me/likes
func (h *handler) likeUser(c *echo.Context) error {
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

	isExist, err := h.repository.IsUserExists(toUser.Username)
	if err != nil {
		c.Logger().Error(
			"failed to check if user exists before like",
			slog.String("username", *username),
			slog.String("targetUsername", toUser.Username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to check user existence"})
	}
	if !isExist {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "user does not exist"})
	}

	isActionExist, err := h.repository.IsActionExists(*username, toUser.Username)
	if err != nil {
		c.Logger().Error(
			"failed to check if action exists before like",
			slog.String("username", *username),
			slog.String("targetUsername", toUser.Username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to check action existence"})
	}
	if isActionExist {
		return c.JSON(http.StatusConflict, errorResponse{Message: "action already exists"})
	}

	if err = h.repository.LikeUser(*username, toUser.Username); err != nil {
		c.Logger().Error(
			"failed to like user",
			slog.String("username", *username),
			slog.String("targetUsername", toUser.Username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to like user"})
	}

	return c.NoContent(http.StatusNoContent)
}

// GET /api/me/liked-by
func (h *handler) listUsersWhoLikedMe(c *echo.Context) error {
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

	users, err := h.repository.ListUsersWhoLiked(*username)
	if err != nil {
		c.Logger().Error(
			"failed to list users who liked user",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to list users who liked me"})
	}

	result, err := toUserSummaryResponses(users)
	if err != nil {
		c.Logger().Error(
			"failed to convert users who liked user",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to validate liked-by users"})
	}

	return c.JSON(http.StatusOK, result)
}

// POST /api/me/nopes
func (h *handler) nopeUser(c *echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "cannot nope yourself"})
	}

	isExist, err := h.repository.IsUserExists(toUser.Username)
	if err != nil {
		c.Logger().Error(
			"failed to check if user exists before nope",
			slog.String("username", *username),
			slog.String("targetUsername", toUser.Username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to check user existence"})
	}
	if !isExist {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "user does not exist"})
	}

	isActionExist, err := h.repository.IsActionExists(*username, toUser.Username)
	if err != nil {
		c.Logger().Error(
			"failed to check if action exists before nope",
			slog.String("username", *username),
			slog.String("targetUsername", toUser.Username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to check action existence"})
	}
	if isActionExist {
		return c.JSON(http.StatusConflict, errorResponse{Message: "action already exists"})
	}

	if err = h.repository.NopeUser(*username, toUser.Username); err != nil {
		c.Logger().Error(
			"failed to nope user",
			slog.String("username", *username),
			slog.String("targetUsername", toUser.Username),
			slog.Any("error", err),
		)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to nope user"})
	}

	return c.NoContent(http.StatusNoContent)
}
