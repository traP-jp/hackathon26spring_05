package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/env"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository"
)

func AuthenticationMiddleware(env *env.Env, repo repository.Repository) echo.MiddlewareFunc {
	if env == nil {
		panic("env is nil")
	}
	if repo == nil {
		panic("repo is nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if env.IsProduction() {
				username := c.Request().Header.Get("X-Forwarded-User")
				if username == "" {
					return echo.ErrUnauthorized
				}

				exists, err := repo.Exists(username)
				if err != nil {
					return echo.ErrInternalServerError
				}
				if !exists {
					return echo.ErrUnauthorized
				}

				setLoginUserRetriever(c, &loginUserRetrieverImpl{username: username})
				return next(c)
			}
			// ローカル環境では常にログイン済みユーザーとして扱う
			setLoginUserRetriever(c, &loginUserRetrieverImpl{username: "test-user"})
			return next(c)
		}
	}
}

type loginUserRetrieverImpl struct {
	username string
}

func (r *loginUserRetrieverImpl) IsUserLoggedIn() bool {
	return r.username != ""
}

func (r *loginUserRetrieverImpl) GetLoginUser() (string, error) {
	return r.username, nil
}

func setLoginUserRetriever(ctx echo.Context, retriever domain.LoginUserRetriever) {
	ctx.Set("loginUserRetriever", retriever)
}

func GetLoginUserRetriever(ctx echo.Context) domain.LoginUserRetriever {
	v := ctx.Get("loginUserRetriever")
	if retriever, ok := v.(domain.LoginUserRetriever); ok && retriever != nil {
		return retriever
	}
	// Middleware 未適用などの場合は未ログイン扱いにする
	return &loginUserRetrieverImpl{}
}
