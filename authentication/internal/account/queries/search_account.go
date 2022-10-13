package queries

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type SearchAccountHandler interface {
	Handle(ctx context.Context, query *SearchAccountQuery) (*models.AccountsList, error)
}

type searchAccountHandler struct {
	log         logger.Logger
	cfg         *config.Config
	accountRepo repository.AccountRepository
	cacheRepo   repository.CacheRepository
}

func NewSearchAccountQueryHandler(log logger.Logger, cfg *config.Config, accountRepo repository.AccountRepository, cacheRepo repository.CacheRepository) SearchAccountHandler {
	return &searchAccountHandler{log: log, cfg: cfg, accountRepo: accountRepo, cacheRepo: cacheRepo}
}

func (s *searchAccountHandler) Handle(ctx context.Context, query *SearchAccountQuery) (*models.AccountsList, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "searchAccountHandler.Handle")
	defer span.Finish()

	return s.accountRepo.SearchAccount(ctx, query.Text, query.Pagination)
}
