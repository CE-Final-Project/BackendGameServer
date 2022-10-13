package repository

var (
	getRoleByIDQuery = `SELECT r.id, r.name, r.created_at, r.created_by, r.updated_at, r.updated_by FROM roles r WHERE r.id = $1`

	getRoleByNameQuery = `SELECT r.id, r.name, r.created_at, r.created_by, r.updated_at, r.updated_by FROM roles r WHERE r.name = $1`

	searchRoleQuery = `SELECT count(*) over() as total, r.id, r.name, r.created_at, r.created_by, r.updated_at, r.updated_by FROM roles r WHERE r.name ILIKE $1 ORDER BY $2 LIMIT $3 OFFSET $4`

	insertRoleQuery = `INSERT INTO roles (name, created_at, created_by, updated_at, updated_by) VALUES ($1, now(), $2, now(), $3) RETURNING id, name, created_at, created_by, updated_at, updated_by`

	updateRoleQuery = `UPDATE roles r SET
				name=COALESCE(NULLIF($1, ''), name),
				updated_at = now(),
              	updated_by = COALESCE(NULLIF($2, 0), updated_by)
              	WHERE r.id = $3
              	RETURNING r.id, r.name, r.created_at, r.created_by, r.updated_at, r.updated_by `

	deleteRoleQuery = `DELETE FROM roles WHERE id = $1`
)
