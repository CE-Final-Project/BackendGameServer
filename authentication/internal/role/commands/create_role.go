package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type CreateRoleCmdHandler interface {
	Handle(ctx context.Context, command *CreateRoleCommand) (*models.Role, error)
}

type createRoleHandler struct {
	log      logger.Logger
	cfg      *config.Config
	roleRepo repository.RoleRepository
}

func (r *createRoleHandler) Handle(ctx context.Context, command *CreateRoleCommand) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "createRoleHandler.Handle")
	defer span.Finish()

	role := &models.Role{
		Name:      command.Name,
		CreatedBy: command.CreatedBy,
		UpdatedBy: command.UpdatedBy,
	}

	var err error
	role, err = r.roleRepo.InsertRole(ctx, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func NewCreateRoleCmdHandler(log logger.Logger, cfg *config.Config, roleRepo repository.RoleRepository) CreateRoleCmdHandler {
	return &createRoleHandler{
		log:      log,
		cfg:      cfg,
		roleRepo: roleRepo,
	}
}
