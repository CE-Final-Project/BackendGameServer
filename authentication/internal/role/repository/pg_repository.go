package repository

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type postgresRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func (p *postgresRepo) GetRoleByID(ctx context.Context, roleID uint64) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "repository.GetRoleByID")
	defer span.Finish()

}

func (p *postgresRepo) GetRoleByName(ctx context.Context, roleName string) (*models.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postgresRepo) SearchRole(ctx context.Context, search string, pagination *utils.Pagination) (*models.RolesList, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postgresRepo) InsertRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postgresRepo) UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postgresRepo) DeleteRoleByID(ctx context.Context, roleID uint64) error {
	//TODO implement me
	panic("implement me")
}

func NewRoleRepository(db *sqlx.DB, log logger.Logger) RoleRepository {
	return &postgresRepo{
		db:  db,
		log: log,
	}
}
