package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type BanAccountByIdCmdHandler interface {
	Handle(ctx context.Context, command *BanAccountByIdCommand) error
}

type banAccountByIdHandler struct {
	log          logger.Logger
	cfg          *config.Config
	postgresRepo repository.Repository
	redisRepo    repository.CacheRepository
}

func NewBanAccountByIdCmdHandler(log logger.Logger, cfg *config.Config, postgresRepo repository.Repository, redisRepo repository.CacheRepository) *banAccountByIdHandler {
	return &banAccountByIdHandler{log: log, cfg: cfg, postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (u *banAccountByIdHandler) Handle(ctx context.Context, command *BanAccountByIdCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "banAccountByIdHandler.Handle")
	defer span.Finish()

	account := &models.Account{
		AccountID: command.AccountID,
		IsBan:     command.IsBan,
		UpdatedAt: command.UpdatedAt,
	}

	updated, err := u.postgresRepo.UpdateAccount(ctx, account)
	if err != nil {
		return err
	}

	u.redisRepo.DelAccount(ctx, updated.AccountID.String())
	u.redisRepo.PutAccount(ctx, updated.AccountID.String(), updated)
	return nil
}
