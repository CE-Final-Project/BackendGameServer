package v1

import (
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/account/commands"
	"github.com/ce-final-project/backend_game_server/gateway/internal/account/queries"
	"github.com/ce-final-project/backend_game_server/gateway/internal/account/service"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/gateway/internal/middlewares"
	httpErrors "github.com/ce-final-project/backend_game_server/pkg/http_errors"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

type accountHandlers struct {
	group *echo.Group
	log   logger.Logger
	mw    middlewares.MiddlewareManager
	cfg   *config.Config
	acs   *service.AccountService
	v     *validator.Validate
}

func NewAccountHandlers(group *echo.Group, log logger.Logger, mw middlewares.MiddlewareManager, cfg *config.Config, acs *service.AccountService, v *validator.Validate) *accountHandlers {
	return &accountHandlers{
		group: group,
		log:   log,
		mw:    mw,
		cfg:   cfg,
		acs:   acs,
		v:     v,
	}
}

// SearchAccount
// @Tags Account
// @Summary Search Account
// @Description Search Account with Username
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        text   query      string  true  "Username"
// @Success 200 {object} dto.AccountsListResponseDto
// @Router /account/search [get]
func (a *accountHandlers) SearchAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "accountHandlers.SearchAccount")
		defer span.Finish()

		text := c.QueryParam("text")
		if text == "" {
			a.log.WarnMsg("text query param is empty", httpErrors.BadRequest)
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest, a.cfg.HTTP.DebugErrorsResponse)
		}
		pagination := utils.NewPaginationQuery(10, 0)
		query := queries.NewSearchAccountQuery(text, pagination)

		result, err := a.acs.Queries.SearchAccount.Handle(ctx, query)
		if err != nil {
			a.log.WarnMsg("Queries.SearchAccount.Handle", err)
			return httpErrors.ErrorCtxResponse(c, errors.Cause(err), a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, result)
	}
}

// GetAccountById
// @Tags Account
// @Summary Account information by ID
// @Description Get all information Account by ID
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Account ID"
// @Success 200 {object} dto.AccountResponseDto
// @Router /account/{id} [get]
func (a *accountHandlers) GetAccountById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "accountHandlers.GetAccountById")
		defer span.Finish()

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

		query := queries.NewGetAccountByIdQuery(accountUUID)

		result, err := a.acs.Queries.GetAccountById.Handle(ctx, query)
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
// @Param        request body dto.UpdateAccount true "Update body request"
// @Success 200 {object} dto.UpdateAccountResponse
// @Router /account [put]
func (a *accountHandlers) UpdateAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "accountHandlers.UpdateAccount")
		defer span.Finish()

		updateDto := &dto.UpdateAccount{}
		if err := c.Bind(updateDto); err != nil {
			a.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		userID := c.Request().Header.Get("User-ID")

		if updateDto.AccountID.String() != userID {
			a.log.WarnMsg("userId and accountId not match", httpErrors.WrongCredentials)
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials, a.cfg.HTTP.DebugErrorsResponse)
		}

		command := commands.NewUpdateAccountCommand(updateDto)

		err := a.acs.Commands.UpdateAccount.Handle(ctx, command)
		if err != nil {
			a.log.WarnMsg("UpdateAccount", errors.Cause(err))
			return httpErrors.ErrorCtxResponse(c, errors.Cause(err), a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, &dto.UpdateAccountResponse{
			AccountID: updateDto.AccountID,
			UpdatedAt: time.Now(),
		})
	}
}

// ChangePassword
// @Tags Account
// @Summary Change password account
// @Description Update password account
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        request body dto.ChangePassword true "ChangePassword body request"
// @Success 200 {object} dto.UpdateAccountResponse
// @Router /account/password [put]
func (a *accountHandlers) ChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "accountHandlers.ChangePassword")
		defer span.Finish()

		changePwdDto := &dto.ChangePassword{}
		if err := c.Bind(changePwdDto); err != nil {
			a.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, a.cfg.HTTP.DebugErrorsResponse)
		}

		userID := c.Request().Header.Get("User-ID")

		if changePwdDto.AccountID.String() != userID {
			a.log.WarnMsg("userId and accountId not match", httpErrors.Unauthorized)
			return httpErrors.ErrorCtxResponse(c, httpErrors.Unauthorized, a.cfg.HTTP.DebugErrorsResponse)
		}

		command := commands.NewChangePasswordCommand(changePwdDto)

		err := a.acs.Commands.ChangePassword.Handle(ctx, command)
		if err != nil {
			a.log.WarnMsg("ChangePassword", errors.Cause(err))
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials, a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, &dto.UpdateAccountResponse{
			AccountID: changePwdDto.AccountID,
			UpdatedAt: time.Now(),
		})
	}
}

// DeleteAccount
// @Tags Account
// @Summary Delete account
// @Description Delete account by id
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Account ID"
// @Success 200 {object} dto.DeleteAccountByIdResponse
// @Router /account/{id} [delete]
func (a *accountHandlers) DeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "accountHandlers.DeleteAccountById")
		defer span.Finish()

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
		command := commands.NewDeleteAccountCommand(&dto.DeleteAccountById{AccountID: accountUUID})

		err = a.acs.Commands.DeleteAccount.Handle(ctx, command)
		if err != nil {
			a.log.WarnMsg("DeleteAccountById", errors.Cause(err))
			return httpErrors.ErrorCtxResponse(c, errors.Cause(err), a.cfg.HTTP.DebugErrorsResponse)
		}
		return c.JSON(http.StatusOK, &dto.DeleteAccountByIdResponse{
			AccountID: accountUUID,
			DeletedAt: time.Now(),
		})
	}
}
