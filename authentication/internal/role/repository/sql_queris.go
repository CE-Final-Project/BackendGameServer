package repository

var (
	getRoleByIDQuery = `SELECT r.id, r.name, r.created_at, r.created_by, r.updated_at, r.updated_by FROM role r WHERE r.id = $1`

	getRoleByNameQuery = `SELECT r.id, r.name, r.created_at, r.created_by, r.updated_at, r.updated_by FROM role r WHERE r.name = $1`

	searchRole = `SELECT count(*) over() as total, r.id, r.name, r.created_at, r.created_by, r.updated_at, r.updated_by FROM role r WHERE r.name ILIKE $1 ORDER BY $2 LIMIT $3 OFFSET $4`
	// TODO: Add sql query string
	insertRole = ``
)
