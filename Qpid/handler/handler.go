package handler

import "github.com/labstack/echo/v4"

type handler struct{}

func Serve() {
	e := echo.New()

	h := &handler{}

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
