package handler

import (
	"cmp"
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/traP-jp/hackathon26spring_05/Qpid/env"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
	"github.com/traP-jp/hackathon26spring_05/Qpid/infrastructure"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository/mock"
	"github.com/traPtitech/go-traq"
)

type handler struct {
	env        env.Env
	repository repository.Repository
	sessions   sessions.Store
	traq       traqClientWithContext
	traqUsers  sync.Map // username → traQ user
}

type traqClientWithContext struct {
	client  *traq.APIClient
	context context.Context
}

func Serve() {
	e := echo.New()

	db, err := infrastructure.NewDB()
	if err != nil {
		e.Logger.Error("Failed to connect to the database", slog.Any("error", err))
		return
	}

	cfg := env.GetEnv()

	var repo repository.Repository
	if cfg.IsProduction() {
		repo = infrastructure.NewRepository(db)
	} else {
		repo = mock.NewMockRepository()
	}

	h := &handler{
		env:        cfg,
		repository: repo,
		sessions: sessions.NewCookieStore([]byte(
			cmp.Or(os.Getenv("SESSION_SECRET"), "secret"),
		)),
		traq: traqClientWithContext{},
	}

	if h.env.TraqAccessToken != "" {
		h.traq = traqClientWithContext{
			client:  newTraqAPIClient(h.env.TraqHost),
			context: context.WithValue(context.Background(), traq.ContextAccessToken, h.env.TraqAccessToken),
		}
	} else {
		e.Logger.Warn("traQ access token is not set. traQ API will not be available.")
	}

	h.mapRoutes(e)
	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Error occurred", slog.Any("error", err))
	}
}

func (h *handler) mapRoutes(e *echo.Echo) {
	api := e.Group("/api",
		echoMiddleware.RequestLogger(),
		echoMiddleware.Recover(),
		middleware.UsernameExtractorMiddleware(&h.env))
	{
		api.POST("/signup", h.signup)

		// 認証が必要な API 群
		authenticated := api.Group("", middleware.AuthenticationMiddleware(h.repository, h.getUserUUID))
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
