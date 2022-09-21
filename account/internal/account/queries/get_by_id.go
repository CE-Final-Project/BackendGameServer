package queries

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type GetAccountByIdHandler interface {
	Handle(ctx context.Context, query *GetAccountByIdQuery) (*models.Account, error)
}

type getAccountByIdHandler struct {
	log          logger.Logger
	cfg          *config.Config
	postgresRepo repository.Repository
	redisRepo    repository.CacheRepository
}

func NewGetAccountByIdHandler(log logger.Logger, cfg *config.Config, postgresRepo repository.Repository, redisRepo repository.CacheRepository) *getAccountByIdHandler {
	return &getAccountByIdHandler{log: log, cfg: cfg, postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (q *getAccountByIdHandler) Handle(ctx context.Context, query *GetAccountByIdQuery) (*models.Account, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getAccountByIdHandler.Handle")
	defer span.Finish()

	if product, err := q.redisRepo.GetAccount(ctx, query.AccountID.String()); err == nil && product != nil {
		return product, nil
	}

	product, err := q.postgresRepo.GetAccountById(ctx, query.AccountID)
	if err != nil {
		return nil, err
	}

	q.redisRepo.PutAccount(ctx, product.AccountID, product)
	return product, nil
}
