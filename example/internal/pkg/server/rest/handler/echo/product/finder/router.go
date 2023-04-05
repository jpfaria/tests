package finder

import "github.com/labstack/echo/v4"

func Router(router *echo.Echo, handler *Handler) {
	router.GET("/product", handler.Handle)
}
