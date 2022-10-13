package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type DeleteRoleCmdHandler interface {
	Handle(ctx context.Context, command *DeleteRoleCommand) error
}

type deleteRoleHandler struct {
	log      logger.Logger
	cfg      *config.Config
	roleRepo repository.RoleRepository
}

func (r *deleteRoleHandler) Handle(ctx context.Context, command *DeleteRoleCommand) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "deleteRoleHandler.Handle")
	defer span.Finish()

	if err := r.roleRepo.DeleteRoleByID(ctx, command.ID); err != nil {
		return err
	}

	return nil
}

func NewDeleteRoleCmdHandler(log logger.Logger, cfg *config.Config, roleRepo repository.RoleRepository) DeleteRoleCmdHandler {
	return &deleteRoleHandler{
		log:      log,
		cfg:      cfg,
		roleRepo: roleRepo,
	}
}
