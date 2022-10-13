package queries

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type GetRoleByNameQueryHandler interface {
	Handle(ctx context.Context, query *GetRoleByNameQuery) (*models.Role, error)
}

type getRoleByNameHandler struct {
	log      logger.Logger
	cfg      *config.Config
	roleRepo repository.RoleRepository
}

func (r *getRoleByNameHandler) Handle(ctx context.Context, query *GetRoleByNameQuery) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "getRoleByIDHandler.Handle")
	defer span.Finish()

	return r.roleRepo.GetRoleByName(ctx, query.Name)
}

func NewGetRoleByNameQueryHandler(log logger.Logger, cfg *config.Config, roleRepo repository.RoleRepository) GetRoleByNameQueryHandler {
	return &getRoleByNameHandler{
		log:      log,
		cfg:      cfg,
		roleRepo: roleRepo,
	}
}
