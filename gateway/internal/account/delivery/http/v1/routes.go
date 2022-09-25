package v1

func (a *accountHandlers) MapRoutes() {

	a.group.GET("/account/search", a.SearchAccount(), a.mw.AuthorizationMiddleware())
	a.group.GET("/account/:id", a.GetAccountById(), a.mw.AuthorizationMiddleware())
	a.group.PUT("/account", a.UpdateAccount(), a.mw.AuthorizationMiddleware())
	a.group.PUT("/account/password", a.ChangePassword(), a.mw.AuthorizationMiddleware())
}
