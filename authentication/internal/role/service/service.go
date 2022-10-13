package service

import (
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/commands"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/queries"
	"github.com/ce-final-project/backend_game_server/authentication/internal/role/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
)

type RoleService struct {
	Commands *commands.RoleCommands
	Queries  *queries.RoleQueries
}

func NewRoleService(log logger.Logger, cfg *config.Config, roleRepo repository.RoleRepository) *RoleService {
	createRoleCmdHandler := commands.NewCreateRoleCmdHandler(log, cfg, roleRepo)
	changeRoleNameCmdHandler := commands.NewChangeRoleNameCmdHandler(log, cfg, roleRepo)
	deleteRoleCmdHandler := commands.NewDeleteRoleCmdHandler(log, cfg, roleRepo)

	searchRoleQueryHandler := queries.NewSearchRoleQueryHandler(log, cfg, roleRepo)
	getRoleByIDQueryHandler := queries.NewGetRoleByIDQueryHandler(log, cfg, roleRepo)
	getRoleByNameQueryHandler := queries.NewGetRoleByNameQueryHandler(log, cfg, roleRepo)

	roleCommands := commands.NewRoleCommands(createRoleCmdHandler, changeRoleNameCmdHandler, deleteRoleCmdHandler)
	roleQueries := queries.NewRoleQueries(searchRoleQueryHandler, getRoleByIDQueryHandler, getRoleByNameQueryHandler)

	return &RoleService{
		Commands: roleCommands,
		Queries:  roleQueries,
	}
}
