package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *authsHandlers) MapRoutes() {

	a.group.POST("/login", a.Login())
	a.group.POST("/register", a.Register())
	a.group.PUT("/account", a.UpdateAccount(), a.mw.AuthorizationMiddleware())
	a.group.PUT("/account/password", a.ChangePassword(), a.mw.AuthorizationMiddleware())
	a.group.GET("/account/:id", a.GetAccount(), a.mw.AuthorizationMiddleware())
	a.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
