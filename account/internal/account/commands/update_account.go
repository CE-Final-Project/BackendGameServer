package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type UpdateAccountCmdHandler interface {
	Handle(ctx context.Context, command *UpdateAccountCommand) error
}

type updateAccountHandler struct {
	log          logger.Logger
	cfg          *config.Config
	postgresRepo repository.Repository
	redisRepo    repository.CacheRepository
}

func NewUpdateAccountCmdHandler(log logger.Logger, cfg *config.Config, postgresRepo repository.Repository, redisRepo repository.CacheRepository) *updateAccountHandler {
	return &updateAccountHandler{log: log, cfg: cfg, postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (u *updateAccountHandler) Handle(ctx context.Context, command *UpdateAccountCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createAccountHandler.Handle")
	defer span.Finish()

	account := &models.Account{
		AccountID: command.AccountID,
		Username:  command.Username,
		Email:     command.Email,
		UpdatedAt: command.UpdatedAt,
	}

	updated, err := u.postgresRepo.UpdateAccount(ctx, account)
	if err != nil {
		return err
	}

	u.redisRepo.DelAccount(ctx, updated.AccountID)
	u.redisRepo.PutAccount(ctx, updated.AccountID, updated)
	return nil
}
