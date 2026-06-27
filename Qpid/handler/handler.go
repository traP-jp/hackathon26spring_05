package handler

import (
	"cmp"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository/mock"
)

type handler struct {
	repository repository.Repository
	sessions   sessions.Store
}

func Serve() {
	e := echo.New()

	// _, err := infrastructure.NewDB()
	// if err != nil {
	// 	e.Logger.Fatal(err)
	// 	return
	// }

	repo := mock.NewMockRepository()

	h := &handler{
		repository: repo,
		sessions: sessions.NewCookieStore([]byte(
			cmp.Or(os.Getenv("SESSION_SECRET"), "secret"),
		)),
	}

	h.mapRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func (h *handler) mapRoutes(e *echo.Echo) {
	api := e.Group("/api")
	{
		login := api.Group("/login")
		{
			login.GET("", h.startLogin)
			login.POST("/callback", h.loginCallback)
		}
		api.POST("/signup", h.signup)
		api.POST("/logout", h.logout)
		me := api.Group("/me")
		{
			me.GET("", h.getMe)
			me.PUT("", h.updateMe)
			me.GET("/likes", h.listMyLikes)
			me.POST("/likes", h.likeUser)
			me.GET("/liked-by", h.listUsersWhoLikedMe)
			me.POST("/nopes", h.nopeUser)
		}
		users := api.Group("/users")
		{
			users.GET("/:id", h.getUser)
		}
		api.GET("/suggestions", h.listSuggestions)
	}
}
