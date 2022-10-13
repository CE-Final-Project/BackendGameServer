package queries

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type GetRoleByIDQueryHandler interface {
	Handle(ctx context.Context, query *GetRoleByIDQuery) (*models.Role, error)
}

type getRoleByIDHandler struct {
	log      logger.Logger
	cfg      *config.Config
	roleRepo repository.RoleRepository
}

func (r *getRoleByIDHandler) Handle(ctx context.Context, query *GetRoleByIDQuery) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "getRoleByIDHandler.Handle")
	defer span.Finish()

	return r.roleRepo.GetRoleByID(ctx, query.ID)
}

func NewGetRoleByIDQueryHandler(log logger.Logger, cfg *config.Config, roleRepo repository.RoleRepository) GetRoleByIDQueryHandler {
	return &getRoleByIDHandler{
		log:      log,
		cfg:      cfg,
		roleRepo: roleRepo,
	}
}
