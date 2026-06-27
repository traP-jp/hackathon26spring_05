package handler

import (
	"cmp"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/traP-jp/hackathon26spring_05/Qpid/env"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
	"github.com/traP-jp/hackathon26spring_05/Qpid/infrastructure"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository/mock"
)

type handler struct {
	env        env.Env
	repository repository.Repository
	sessions   sessions.Store
}

func Serve() {
	e := echo.New()

	_, err := infrastructure.NewDB()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	repo := mock.NewMockRepository()

	h := &handler{
		env:        env.GetEnv(),
		repository: repo,
		sessions: sessions.NewCookieStore([]byte(
			cmp.Or(os.Getenv("SESSION_SECRET"), "secret"),
		)),
	}

	h.mapRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func (h *handler) mapRoutes(e *echo.Echo) {
	api := e.Group("/api",
		echoMiddleware.RequestLogger(),
		echoMiddleware.Recover(),
		middleware.UsernameExtractorMiddleware(&h.env))
	{
		api.POST("/signup", h.signup)

		// 認証が必要な API 群
		authenticated := api.Group("", middleware.AuthenticationMiddleware(h.repository))
		{
			me := authenticated.Group("/me")
			{
				me.GET("", h.getMe)
				me.PUT("", h.updateMe)
				me.GET("/likes", h.listMyLikes)
				me.POST("/likes", h.likeUser)
				me.GET("/liked-by", h.listUsersWhoLikedMe)
				me.POST("/nopes", h.nopeUser)
				me.PATCH("/icon", h.updateMyIcon)
				me.DELETE("/icon", h.deleteMyIcon)
			}
			users := authenticated.Group("/users")
			{
				users.GET("/:id", h.getUser)
				users.GET("/:id/icon", h.getUserIcon)
			}
			authenticated.GET("/suggestions", h.listSuggestions)
		}
	}
}
