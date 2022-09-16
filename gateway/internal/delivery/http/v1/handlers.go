package v1

import (
	"context"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/gateway/internal/middlewares"
	"github.com/ce-final-project/backend_game_server/gateway/internal/service"
	httpErrors "github.com/ce-final-project/backend_game_server/pkg/http_errors"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type authsHandlers struct {
	group *echo.Group
	log   logger.Logger
	mw    middlewares.MiddlewareManager
	cfg   *config.Config
	acs   service.AccountService
	as    service.AuthService
	v     *validator.Validate
}

func NewAuthsHandlers(group *echo.Group, log logger.Logger, mw middlewares.MiddlewareManager, cfg *config.Config, acs service.AccountService, as service.AuthService, v *validator.Validate) *authsHandlers {
	return &authsHandlers{
		group: group,
		log:   log,
		mw:    mw,
		cfg:   cfg,
		acs:   acs,
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
// @Param        request body dto.RegisterAccountDto true "Register body request"
// @Success 201 {object} dto.RegisterAccountResponseDto
// @Router /register [post]
func (a *authsHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		registerDto := &dto.RegisterAccountDto{}
		if err := c.Bind(registerDto); err != nil {
			a.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		if err := a.v.StructCtx(ctx, registerDto); err != nil {
			a.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}
		result, err := a.as.Register(ctx, registerDto)
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
// @Param        request body dto.LoginAccountDto true "Login body request"
// @Success 200 {object} dto.LoginAccountResponseDto
// @Router /login [post]
func (a *authsHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		LoginDto := &dto.LoginAccountDto{}
		if err := c.Bind(LoginDto); err != nil {
			a.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		if err := a.v.StructCtx(ctx, LoginDto); err != nil {
			a.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}
		result, err := a.as.Login(ctx, LoginDto)
		if err != nil {
			a.log.WarnMsg("LoginAccount", httpErrors.WrongCredentials)
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials, a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, result)
	}
}

// GetAccount
// @Tags Account
// @Summary Account information by ID
// @Description Get all information Account by ID
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Account ID"
// @Success 200 {object} dto.AccountResponseDto
// @Router /account/{id} [get]
func (a *authsHandlers) GetAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		accountID := c.Param("id")

		userID := c.Request().Header.Get("User-ID")

		if accountID != userID {
			a.log.WarnMsg("userId and accountId not match", httpErrors.Unauthorized)
			return httpErrors.ErrorCtxResponse(c, httpErrors.Unauthorized, a.cfg.HTTP.DebugErrorsResponse)
		}

		accountUUID, err := uuid.FromString(accountID)
		if err != nil {
			a.log.WarnMsg("uuid.FromString", errors.Cause(err))
			return httpErrors.ErrorCtxResponse(c, errors.Cause(err), a.cfg.HTTP.DebugErrorsResponse)
		}
		result, err := a.acs.GetAccountById(ctx, accountUUID)
		if err != nil {
			a.log.WarnMsg("GetAccountById", errors.Cause(err))
			return httpErrors.ErrorCtxResponse(c, errors.Cause(err), a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, result)
	}
}

// UpdateAccount
// @Tags Account
// @Summary Update account
// @Description Update account: username or email
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        request body dto.UpdateAccountDto true "Update body request"
// @Success 200 {object} dto.UpdateAccountResponseDto
// @Router /account [put]
func (a *authsHandlers) UpdateAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		updateDto := &dto.UpdateAccountDto{}
		if err := c.Bind(updateDto); err != nil {
			a.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		userID := c.Request().Header.Get("User-ID")

		if updateDto.AccountID.String() != userID {
			a.log.WarnMsg("userId and accountId not match", httpErrors.WrongCredentials)
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials, a.cfg.HTTP.DebugErrorsResponse)
		}

		result, err := a.acs.UpdateAccount(ctx, updateDto)
		if err != nil {
			a.log.WarnMsg("UpdateAccount", errors.Cause(err))
			return httpErrors.ErrorCtxResponse(c, errors.Cause(err), a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, result)
	}
}

// ChangePassword
// @Tags Account
// @Summary Change password account
// @Description Update password account
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        request body dto.ChangePasswordDto true "ChangePassword body request"
// @Success 200 {object} dto.UpdateAccountResponseDto
// @Router /account/password [put]
func (a *authsHandlers) ChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		changePwdDto := &dto.ChangePasswordDto{}
		if err := c.Bind(changePwdDto); err != nil {
			a.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		userID := c.Request().Header.Get("User-ID")

		if changePwdDto.AccountID.String() != userID {
			a.log.WarnMsg("userId and accountId not match", httpErrors.Unauthorized)
			return httpErrors.ErrorCtxResponse(c, httpErrors.Unauthorized, a.cfg.HTTP.DebugErrorsResponse)
		}

		result, err := a.acs.ChangePassword(ctx, changePwdDto)
		if err != nil {
			a.log.WarnMsg("ChangePassword", errors.Cause(err))
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials, a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, result)
	}
}
