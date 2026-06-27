package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
)

type signupRequest struct {
	Agreed bool `json:"agreed"`
}

// POST /api/signup
func (h *handler) signup(c echo.Context) error {
	req := &signupRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			&errorResponse{"Invalid request body"},
		)
	}
	if !req.Agreed {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			&errorResponse{"You must agree to terms"},
		)
	}

	username := middleware.GetUsername(c)
	if username != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			&errorResponse{"Already signed up"},
		)
	}

	precomputed, err := h.repository.FindPrecomputedProfileByUsername(*username)
	if err != nil {
		c.Logger().Errorf("failed to find precomputed profile of user %s: %v", username, err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			&errorResponse{"Failed to find precomputed profile"},
		)
	}

	newUser := mergeUserProfile(
		domain.User{Username: *username},
		precomputed,
	)
	err = h.repository.CreateUser(newUser)
	if err != nil {
		c.Logger().Errorf("failed to create user %s: %v", username, err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			&errorResponse{"Failed to create user"},
		)
	}
	return c.JSON(http.StatusOK, toUserResponse(newUser))
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
