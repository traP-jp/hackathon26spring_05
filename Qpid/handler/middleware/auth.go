package middleware

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/hackathon26spring_05/Qpid/env"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository"
)

func UsernameExtractorMiddleware(env *env.Env) echo.MiddlewareFunc {
	if env == nil {
		panic("env is nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if env.IsProduction() {
				username := c.Request().Header.Get("X-Forwarded-User")
				if username != "" {
					setUsername(c, username)
				}
				return next(c)
			}
			// ローカル環境では常にログイン済みユーザーとして扱う
			setUsername(c, "test-user")
			return next(c)
		}
	}
}

func AuthenticationMiddleware(
	repo repository.Repository,
	getUserUUID func(username string) (string, error),
) echo.MiddlewareFunc {
	if repo == nil {
		panic("repo is nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			username := GetUsername(c)
			if username == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			exists, err := repo.IsUserExists(*username)
			if err != nil {
				c.Logger().Error(
					"failed to check if user exists",
					slog.String("username", *username),
					slog.Any("error", err),
				)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check if user exists")
			}
			if !exists {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			uuid, err := getUserUUID(*username)
			if err != nil {
				c.Logger().Warn("failed to get traQ UUID",
					slog.String("username", *username),
					slog.Any("error", err),
				)
			} else if uuid != "" {
				setUserUUID(c, uuid)
			}

			return next(c)
		}
	}
}

const (
	usernameKey = "Username"
	uuidKey     = "UserUUID"
)

func setUsername(c *echo.Context, username string) {
	c.Set(usernameKey, username)
}

// ログイン中のユーザー名を取得する. ログインしていない場合は nil.
func GetUsername(c *echo.Context) *string {
	username, ok := c.Get(usernameKey).(string)
	if !ok {
		return nil
	}
	return &username
}

func setUserUUID(c *echo.Context, uuid string) {
	c.Set(uuidKey, uuid)
}

// ログイン中のユーザーの traQ UUID を取得する. 未取得の場合は nil.
func GetUserUUID(c *echo.Context) *string {
	uuid, ok := c.Get(uuidKey).(string)
	if !ok || uuid == "" {
		return nil
	}
	return &uuid
}
