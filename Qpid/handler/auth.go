package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
)

type signupRequest struct {
	Agreed bool `json:"agreed"`
}

// POST /api/signup
func (h *handler) signup(c *echo.Context) error {
	req := &signupRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if !req.Agreed {
		return echo.NewHTTPError(http.StatusBadRequest, "You must agree to terms")
	}

	username := middleware.GetUsername(c)
	if username == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "You must be logged in to traQ to sign up")
	}

	exists, err := h.repository.IsUserExists(*username)
	if err != nil {
		c.Logger().Error(
			"failed to check if user exists",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if user exists")
	}
	if exists {
		return echo.NewHTTPError(http.StatusConflict, "User already exists")
	}

	precomputed, err := h.repository.FindPrecomputedProfileByUsername(*username)
	if err != nil {
		c.Logger().Error(
			"failed to find precomputed profile",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to find precomputed profile")
	}

	newUser := mergeUserProfile(
		domain.User{Username: *username},
		precomputed,
	)
	err = h.repository.CreateUser(newUser)
	if err != nil {
		c.Logger().Error(
			"failed to create user",
			slog.String("username", *username),
			slog.Any("error", err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	uuid, err := h.getUserUUID(*username)
	if err != nil {
		c.Logger().Warn("failed to get traQ UUID, affiliations will be empty",
			slog.String("username", *username),
			slog.Any("error", err),
		)
	}
	if uuid != "" {
		newUser.Affiliations = h.fetchAffiliations(c, uuid)
	} else {
		newUser.Affiliations = []domain.UserAffiliation{}
	}

	return c.JSON(http.StatusCreated, toUserResponse(newUser))
}

func mergeUserProfile(base domain.User, override *domain.UserOverride) domain.User {
	if override == nil {
		return base
	}
	if override.Major.IsSome() {
		if major := override.Major.Unwrap(); major != "" {
			base.Major = optional.Some(major)
		}
	}
	if override.Hometown.IsSome() {
		if hometown := override.Hometown.Unwrap(); hometown != "" {
			base.Hometown = optional.Some(hometown)
		}
	}
	if override.Bio.IsSome() {
		if bio := override.Bio.Unwrap(); bio != "" {
			base.Bio = optional.Some(bio)
		}
	}
	if len(override.Affiliations) > 0 {
		if base.Affiliations == nil {
			base.Affiliations = []domain.UserAffiliation{}
		}
		base.Affiliations = distinctValues(append(base.Affiliations, override.Affiliations...))
	}
	if len(override.Tags) > 0 {
		if base.Tags == nil {
			base.Tags = []string{}
		}
		base.Tags = distinctValues(append(base.Tags, override.Tags...))
	}
	if len(override.Technologies) > 0 {
		base.Technologies = distinctValues(append(base.Technologies, override.Technologies...))
	}
	if override.FavoriteTopic.IsSome() {
		base.FavoriteTopic = override.FavoriteTopic
	}
	if override.DislikedTopic.IsSome() {
		base.DislikedTopic = override.DislikedTopic
	}
	return base
}

func distinctValues[T comparable](items []T) []T {
	unique := make([]T, 0, len(items))
	seen := make(map[T]struct{})
	for _, item := range items {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			unique = append(unique, item)
		}
	}
	return unique
}
