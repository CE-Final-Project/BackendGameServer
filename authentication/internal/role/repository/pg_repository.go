package repository

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type postgresRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func (p *postgresRepo) GetRoleByID(ctx context.Context, roleID uint64) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "repository.GetRoleByID")
	defer span.Finish()

	var role models.Role
	if err := p.db.QueryRowContext(ctx, getRoleByIDQuery, &roleID).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.CreatedBy,
		&role.UpdatedAt,
		&role.UpdatedBy); err != nil {
		return nil, errors.Wrap(err, "repository.GetRoleByID")
	}

	return &role, nil
}

func (p *postgresRepo) GetRoleByName(ctx context.Context, roleName string) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "repository.GetRoleByName")
	defer span.Finish()

	var role models.Role
	if err := p.db.QueryRowContext(ctx, getRoleByNameQuery, &roleName).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.CreatedBy,
		&role.UpdatedAt,
		&role.UpdatedBy); err != nil {
		return nil, errors.Wrap(err, "repository.GetRoleByName")
	}

	return &role, nil
}

func (p *postgresRepo) SearchRole(ctx context.Context, search string, pagination *utils.Pagination) (*models.RolesList, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "repository.SearchRole")
	defer span.Finish()

	rows, err := p.db.QueryContext(ctx, searchRoleQuery, "%"+search+"%", pagination.GetOrderBy(), pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, errors.Wrap(err, "repository.SearchRole")
	}

	var total int64
	roles := make([]*models.Role, 0, pagination.GetSize())

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(
			&total,
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.CreatedBy,
			&role.UpdatedAt,
			&role.UpdatedBy,
		); err != nil {
			return nil, errors.Wrap(err, "Scan Search Role")
		}
		roles = append(roles, &role)
	}

	return models.NewRoleListWithPagination(roles, total, pagination), nil
}

func (p *postgresRepo) InsertRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "repository.InsertRole")
	defer span.Finish()

	var created models.Role

	if err := p.db.QueryRowContext(ctx, insertRoleQuery, role.Name, role.CreatedBy, role.UpdatedBy).Scan(
		&created.ID,
		&created.Name,
		&created.CreatedAt,
		&created.CreatedBy,
		&created.UpdatedAt,
		&created.UpdatedBy,
	); err != nil {
		return nil, errors.Wrap(err, "repository.InsertRole")
	}

	return &created, nil
}

func (p *postgresRepo) UpdateRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "repository.UpdateRole")
	defer span.Finish()

	var updated models.Role

	if err := p.db.QueryRowContext(ctx, updateRoleQuery, role.Name, role.UpdatedBy).Scan(
		&updated.ID,
		&updated.Name,
		&updated.CreatedAt,
		&updated.CreatedBy,
		&updated.UpdatedAt,
		&updated.UpdatedBy,
	); err != nil {
		return nil, errors.Wrap(err, "repository.UpdateRole")
	}

	return &updated, nil
}

func (p *postgresRepo) DeleteRoleByID(ctx context.Context, roleID uint64) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "repository.DeleteRoleByID")
	defer span.Finish()

	row := p.db.QueryRowContext(ctx, deleteRoleQuery, roleID)
	if row.Err() != nil {
		return errors.Wrap(row.Err(), "repository.DeleteRoleByID")
	}

	return nil
}

func NewRoleRepository(db *sqlx.DB, log logger.Logger) RoleRepository {
	return &postgresRepo{
		db:  db,
		log: log,
	}
}
