package queries

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type GetAccountByUsernameHandler interface {
	Handle(ctx context.Context, query *GetAccountByUsernameQuery) (*models.Account, error)
}

type getAccountByUsernameHandler struct {
	log          logger.Logger
	cfg          *config.Config
	postgresRepo repository.Repository
	redisRepo    repository.CacheRepository
}

func NewGetAccountByUsernameHandler(log logger.Logger, cfg *config.Config, postgresRepo repository.Repository, redisRepo repository.CacheRepository) *getAccountByUsernameHandler {
	return &getAccountByUsernameHandler{log: log, cfg: cfg, postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (q *getAccountByUsernameHandler) Handle(ctx context.Context, query *GetAccountByUsernameQuery) (*models.Account, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getAccountByUsernameHandler.Handle")
	defer span.Finish()

	if account, err := q.redisRepo.GetAccountReference(ctx, query.Username); err == nil && account != nil {
		return account, nil
	}

	account, err := q.postgresRepo.GetAccountByEmailOrUsername(ctx, "", query.Username)
	if err != nil {
		return nil, err
	}

	q.redisRepo.PutKeyReference(ctx, account.Username, account.AccountID.String())
	return account, nil
}
