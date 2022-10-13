package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"strconv"
)

type DeleteAccountCmdHandler interface {
	Handle(ctx context.Context, command *DeleteAccountCommand) error
}

type deleteAccountCmdHandler struct {
	log         logger.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
	cacheRepo   repository.CacheRepository
}

func NewDeleteAccountCmdHandler(log logger.Logger, cfg *config.Config, accountRepo repository.AccountRepository, cacheRepo repository.CacheRepository) DeleteAccountCmdHandler {
	return &deleteAccountCmdHandler{log: log, cfg: cfg, accountRepo: accountRepo, cacheRepo: cacheRepo}
}

func (c *deleteAccountCmdHandler) Handle(ctx context.Context, command *DeleteAccountCommand) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "deleteAccountCmdHandler.Handle")
	defer span.Finish()

	if err := c.accountRepo.DeleteAccountByID(ctx, command.ID); err != nil {
		return err
	}

	c.cacheRepo.DelAccount(ctx, strconv.FormatUint(command.ID, 10))
	return nil
}
