package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type ChangeRoleNameCmdHandler interface {
	Handle(ctx context.Context, command *ChangeRoleNameCommand) (*models.Role, error)
}

type changeRoleNameHandler struct {
	log      logger.Logger
	cfg      *config.Config
	roleRepo repository.RoleRepository
}

func (r *changeRoleNameHandler) Handle(ctx context.Context, command *ChangeRoleNameCommand) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "changeRoleNameHandler.Handle")
	defer span.Finish()

	role := &models.Role{
		ID:        command.ID,
		Name:      command.Name,
		UpdatedBy: command.UpdatedBy,
	}

	role, err := r.roleRepo.UpdateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func NewChangeRoleNameCmdHandler(log logger.Logger, cfg *config.Config, roleRepo repository.RoleRepository) ChangeRoleNameCmdHandler {
	return &changeRoleNameHandler{
		log:      log,
		cfg:      cfg,
		roleRepo: roleRepo,
	}
}
