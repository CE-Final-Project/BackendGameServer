package queries

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
)

type SearchAccountHandler interface {
	Handle(ctx context.Context, query *SearchAccountQuery) (*models.AccountsList, error)
}

type searchAccountHandler struct {
	log          logger.Logger
	cfg          *config.Config
	postgresRepo repository.Repository
	redisRepo    repository.CacheRepository
}

func NewSearchAccountHandler(log logger.Logger, cfg *config.Config, postgresRepo repository.Repository, redisRepo repository.CacheRepository) *searchAccountHandler {
	return &searchAccountHandler{log: log, cfg: cfg, postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (s *searchAccountHandler) Handle(ctx context.Context, query *SearchAccountQuery) (*models.AccountsList, error) {
	return s.postgresRepo.Search(ctx, query.Text, query.Pagination)
}
