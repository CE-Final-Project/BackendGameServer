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

type GetAccountByIdHandler interface {
	Handle(ctx context.Context, query *GetAccountByIdQuery) (*models.Account, error)
}

type getAccountByIdHandler struct {
	log         logger.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
	cacheRepo   repository.CacheRepository
}

func NewGetAccountByIdQueryHandler(log logger.Logger, cfg *config.Config, accountRepo repository.AccountRepository, cacheRepo repository.CacheRepository) GetAccountByIdHandler {
	return &getAccountByIdHandler{log: log, cfg: cfg, accountRepo: accountRepo, cacheRepo: cacheRepo}
}

func (q *getAccountByIdHandler) Handle(ctx context.Context, query *GetAccountByIdQuery) (*models.Account, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "getAccountByIdHandler.Handle")
	defer span.Finish()

	if account, err := q.cacheRepo.GetAccount(ctx, strconv.FormatUint(query.ID, 10)); err == nil && account != nil {
		return account, nil
	}

	account, err := q.accountRepo.GetAccountByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	q.cacheRepo.PutAccount(ctx, strconv.FormatUint(query.ID, 10), account)
	return account, nil
}
