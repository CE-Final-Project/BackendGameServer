package queries

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"strconv"
)

type GetAccountByEmailHandler interface {
	Handle(ctx context.Context, query *GetAccountByEmailQuery) (*models.Account, error)
}

type getAccountByEmailHandler struct {
	log         logger.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
	cacheRepo   repository.CacheRepository
}

func NewGetAccountByEmailHandler(log logger.Logger, cfg *config.Config, accountRepo repository.AccountRepository, cacheRepo repository.CacheRepository) GetAccountByEmailHandler {
	return &getAccountByEmailHandler{log: log, cfg: cfg, accountRepo: accountRepo, cacheRepo: cacheRepo}
}

func (q *getAccountByEmailHandler) Handle(ctx context.Context, query *GetAccountByEmailQuery) (*models.Account, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "getAccountByEmailHandler.Handle")
	defer span.Finish()

	if account, err := q.cacheRepo.GetAccountByKeyReference(ctx, query.Email); err == nil && account != nil {
		return account, nil
	}

	account, err := q.accountRepo.GetAccountByEmail(ctx, query.Email)
	if err != nil {
		return nil, err
	}

	q.cacheRepo.PutKeyReference(ctx, account.Email, strconv.FormatUint(account.ID, 10))
	return account, nil
}
