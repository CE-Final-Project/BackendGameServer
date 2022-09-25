package delivery

import "github.com/labstack/echo/v4"

type HttpDelivery interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
}
