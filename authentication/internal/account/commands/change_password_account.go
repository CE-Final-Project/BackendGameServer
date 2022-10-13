package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strconv"
)

type ChangePasswordCmdHandler interface {
	Handle(ctx context.Context, command *ChangePasswordCommand) error
}

type changePasswordHandler struct {
	log         logger.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
	cacheRepo   repository.CacheRepository
}

func NewChangePasswordCmdHandler(log logger.Logger, cfg *config.Config, accountRepo repository.AccountRepository, cacheRepo repository.CacheRepository) ChangePasswordCmdHandler {
	return &changePasswordHandler{log: log, cfg: cfg, accountRepo: accountRepo, cacheRepo: cacheRepo}
}

func (u *changePasswordHandler) Handle(ctx context.Context, command *ChangePasswordCommand) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "changePasswordHandler.Handle")
	defer span.Finish()

	var account *models.Account
	account, _ = u.cacheRepo.GetAccount(ctx, strconv.FormatUint(command.ID, 10))
	if account == nil {
		var err error
		account, err = u.accountRepo.GetAccountByID(ctx, command.ID)
		if err != nil {
			return err
		}
	}

	if !utils.CheckPasswordHash(command.OldPassword, account.PasswordHashed) {
		return errors.New("Invalid Old Password!")
	}

	passwordHashed, err := utils.HashPassword(command.NewPassword)
	if err != nil {
		return err
	}

	account = &models.Account{
		ID:             command.ID,
		Email:          account.Email,
		Username:       account.Username,
		PasswordHashed: passwordHashed,
	}

	updated, err := u.accountRepo.UpdateAccount(ctx, account)
	if err != nil {
		return err
	}

	u.cacheRepo.DelAccount(ctx, strconv.FormatUint(updated.ID, 10))
	u.cacheRepo.PutAccount(ctx, strconv.FormatUint(updated.ID, 10), updated)
	return nil
}
