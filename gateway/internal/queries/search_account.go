package queries

import (
	"context"
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/dto"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_rest_api/pkg/tracing"
	"github.com/opentracing/opentracing-go"
)

type SearchAccountHandler interface {
	Handle(ctx context.Context, query *SearchAccountQuery) (*dto.AccountsListResponseDto, error)
}

type searchAccountHandler struct {
	log      logger.Logger
	cfg      *config.Config
	asClient authService.AuthServiceClient
}

func NewSearchAccountHandler(log logger.Logger, cfg *config.Config, asClient authService.AuthServiceClient) *searchAccountHandler {
	return &searchAccountHandler{log: log, cfg: cfg, asClient: asClient}
}

func (s *searchAccountHandler) Handle(ctx context.Context, query *SearchAccountQuery) (*dto.AccountsListResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "searchAccountHandler.Handle")
	defer span.Finish()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := s.asClient.SearchAccount(ctx, &authService.SearchReq{
		Search: query.Text,
		Page:   int64(query.Pagination.GetPage()),
		Size:   int64(query.Pagination.GetSize()),
	})
	if err != nil {
		return nil, err
	}

	return dto.AccountsListResponseFromGrpc(res), nil
}
