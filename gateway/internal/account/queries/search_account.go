package queries

import (
	"context"
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type SearchAccountHandler interface {
	Handle(ctx context.Context, query *SearchAccountQuery) (*dto.AccountsListResponseDto, error)
}

type searchAccountHandler struct {
	log      logger.Logger
	cfg      *config.Config
	asClient accountService.AccountServiceClient
}

func NewSearchAccountHandler(log logger.Logger, cfg *config.Config, asClient accountService.AccountServiceClient) *searchAccountHandler {
	return &searchAccountHandler{log: log, cfg: cfg, asClient: asClient}
}

func (s *searchAccountHandler) Handle(ctx context.Context, query *SearchAccountQuery) (*dto.AccountsListResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "searchAccountHandler.Handle")
	defer span.Finish()

	res, err := s.asClient.SearchAccount(ctx, &accountService.SearchReq{
		Search: query.Text,
		Page:   int64(query.Pagination.GetPage()),
		Size:   int64(query.Pagination.GetSize()),
	})
	if err != nil {
		return nil, err
	}

	return dto.AccountsListResponseFromGrpc(res), nil
}
