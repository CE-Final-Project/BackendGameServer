package queries

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type SearchRoleQueryHandler interface {
	Handle(ctx context.Context, query *SearchRoleQuery) (*models.RolesList, error)
}

type searchRoleHandler struct {
	log      logger.Logger
	cfg      *config.Config
	roleRepo repository.RoleRepository
}

func (r *searchRoleHandler) Handle(ctx context.Context, query *SearchRoleQuery) (*models.RolesList, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "searchRoleHandler.Handle")
	defer span.Finish()

	// default text query
	if query.Text == "" {
		query.Text = "%"
	}

	return r.roleRepo.SearchRole(ctx, query.Text, query.Pagination)
}

func NewSearchRoleQueryHandler(log logger.Logger, cfg *config.Config, roleRepo repository.RoleRepository) SearchRoleQueryHandler {
	return &searchRoleHandler{
		log:      log,
		cfg:      cfg,
		roleRepo: roleRepo,
	}
}
