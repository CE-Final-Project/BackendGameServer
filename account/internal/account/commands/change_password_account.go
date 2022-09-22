package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type ChangePasswordCmdHandler interface {
	Handle(ctx context.Context, command *ChangePasswordCommand) error
}

type changePasswordHandler struct {
	log          logger.Logger
	cfg          *config.Config
	postgresRepo repository.Repository
	redisRepo    repository.CacheRepository
}

func NewChangePasswordCmdHandler(log logger.Logger, cfg *config.Config, postgresRepo repository.Repository, redisRepo repository.CacheRepository) *changePasswordHandler {
	return &changePasswordHandler{log: log, cfg: cfg, postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (u *changePasswordHandler) Handle(ctx context.Context, command *ChangePasswordCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "changePasswordHandler.Handle")
	defer span.Finish()

	var account *models.Account
	account, _ = u.redisRepo.GetAccount(ctx, command.AccountID.String())
	if account == nil {
		var err error
		account, err = u.postgresRepo.GetAccountById(ctx, command.AccountID)
		if err != nil {
			return nil
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
		AccountID:      command.AccountID,
		PasswordHashed: passwordHashed,
		UpdatedAt:      command.UpdatedAt,
	}

	updated, err := u.postgresRepo.UpdateAccount(ctx, account)
	if err != nil {
		return err
	}

	u.redisRepo.DelAccount(ctx, updated.AccountID.String())
	u.redisRepo.PutAccount(ctx, updated.AccountID.String(), updated)
	return nil
}
