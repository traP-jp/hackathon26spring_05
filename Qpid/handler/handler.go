package handler

import "github.com/labstack/echo/v4"

type handler struct{}

func Serve() {
	e := echo.New()

	h := &handler{}

	h.mapRoutes(e)
	e.Start(":8080")
}

func (h *handler) mapRoutes(e *echo.Echo) {
	e.GET("/", h.indexHandler)
}
