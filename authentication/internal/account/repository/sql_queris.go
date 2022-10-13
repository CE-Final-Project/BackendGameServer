package repository

const (
	insertAccountQuery = `INSERT INTO accounts ( player_id, username, email, password_hashed, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, now(), now()) RETURNING id, player_id, username, email, password_hashed, created_at, updated_at`

	updateAccountQuery = `UPDATE accounts a SET 
                      username=COALESCE(NULLIF($1, ''), username), 
                      email=COALESCE(NULLIF($2, ''), email), 
                      password_hashed=COALESCE(NULLIF($3, ''), password_hashed),
                      updated_at = now()
                      WHERE id=$4
                      RETURNING id, player_id, username, email, password_hashed, created_at, updated_at`

	searchAccountQuery = `SELECT count(*) over() as total, a.id, a.player_id, a.username, a.email, a.created_at, a.updated_at
	FROM accounts a WHERE a.username ILIKE $1 OR a.player_id ILIKE $1 ORDER BY $2 LIMIT $3 OFFSET $4`

	getAccountByIdQuery = `SELECT a.id, a.player_id, a.username, a.email, a.password_hashed, a.created_at, a.updated_at
	FROM accounts a WHERE a.id = $1`

	getAccountByUsernameQuery = `SELECT a.id, a.player_id, a.username, a.email, a.password_hashed, a.created_at, a.updated_at
	FROM accounts a WHERE a.username = $1`

	getAccountByEmailQuery = `SELECT a.id, a.player_id, a.username, a.email, a.password_hashed, a.created_at, a.updated_at
	FROM accounts a WHERE a.email = $1`

	deleteAccountByIdQuery = `DELETE FROM accounts WHERE id = $1`
)
