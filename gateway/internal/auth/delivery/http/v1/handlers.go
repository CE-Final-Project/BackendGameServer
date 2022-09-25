package v1

import (
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/auth/commands"
	"github.com/ce-final-project/backend_game_server/gateway/internal/auth/service"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/gateway/internal/middlewares"
	httpErrors "github.com/ce-final-project/backend_game_server/pkg/http_errors"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type authsHandlers struct {
	group *echo.Group
	log   logger.Logger
	mw    middlewares.MiddlewareManager
	cfg   *config.Config
	as    *service.AuthService
	v     *validator.Validate
}

func NewAuthsHandlers(group *echo.Group, log logger.Logger, mw middlewares.MiddlewareManager, cfg *config.Config, as *service.AuthService, v *validator.Validate) *authsHandlers {
	return &authsHandlers{
		group: group,
		log:   log,
		mw:    mw,
		cfg:   cfg,
		as:    as,
		v:     v,
	}
}

// Register
// @Tags Authentication
// @Summary Register Account
// @Description Create new user account
// @Accept json
// @Produce json
// @Param        request body dto.RegisterAccount true "Register body request"
// @Success 201 {object} dto.RegisterAccountResponse
// @Router /register [post]
func (a *authsHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "authHandlers.Register")
		defer span.Finish()
		registerDto := &dto.RegisterAccount{}
		if err := c.Bind(registerDto); err != nil {
			a.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		if err := a.v.StructCtx(ctx, registerDto); err != nil {
			a.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		command := commands.NewRegisterAccountCommand(registerDto)

		result, err := a.as.Commands.RegisterAccount.Handle(ctx, command)
		if err != nil {
			a.log.WarnMsg("RegisterAccount", httpErrors.BadRequest)
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest, a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusCreated, result)
	}
}

// Login
// @Tags Authentication
// @Summary Login
// @Description Login with Username and Password
// @Accept json
// @Produce json
// @Param        request body dto.LoginAccount true "Login body request"
// @Success 200 {object} dto.LoginAccountResponse
// @Router /login [post]
func (a *authsHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "authHandlers.Login")
		defer span.Finish()
		LoginDto := &dto.LoginAccount{}
		if err := c.Bind(LoginDto); err != nil {
			a.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		if err := a.v.StructCtx(ctx, LoginDto); err != nil {
			a.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}
		command := commands.NewLoginAccountCommand(LoginDto)
		result, err := a.as.Commands.LoginAccount.Handle(ctx, command)
		if err != nil {
			a.log.WarnMsg("LoginAccount", httpErrors.WrongCredentials)
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials, a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, result)
	}
}
