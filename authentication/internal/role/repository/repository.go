package repository

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
)

type RoleRepository interface {
	GetRoleByID(ctx context.Context, roleID uint64) (*models.Role, error)
	GetRoleByName(ctx context.Context, roleName string) (*models.Role, error)
	SearchRole(ctx context.Context, search string, pagination *utils.Pagination) (*models.RolesList, error)

	InsertRole(ctx context.Context, role *models.Role) (*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error)
	DeleteRoleByID(ctx context.Context, roleID uint64) error
}
