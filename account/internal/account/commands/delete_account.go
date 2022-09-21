package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type DeleteAccountCmdHandler interface {
	Handle(ctx context.Context, command *DeleteAccountCommand) error
}

type deleteAccountCmdHandler struct {
	log          logger.Logger
	cfg          *config.Config
	postgresRepo repository.Repository
	redisRepo    repository.CacheRepository
}

func NewDeleteAccountCmdHandler(log logger.Logger, cfg *config.Config, postgresRepo repository.Repository, redisRepo repository.CacheRepository) *deleteAccountCmdHandler {
	return &deleteAccountCmdHandler{log: log, cfg: cfg, postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (c *deleteAccountCmdHandler) Handle(ctx context.Context, command *DeleteAccountCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteAccountCmdHandler.Handle")
	defer span.Finish()

	if err := c.postgresRepo.DeleteAccount(ctx, command.AccountID); err != nil {
		return err
	}

	c.redisRepo.DelAccount(ctx, command.AccountID.String())
	return nil
}
