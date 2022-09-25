package delivery

import "github.com/labstack/echo/v4"

type HttpDelivery interface {
	UpdateAccount() echo.HandlerFunc
	ChangePassword() echo.HandlerFunc

	GetAccountById() echo.HandlerFunc
	SearchAccount() echo.HandlerFunc
}
