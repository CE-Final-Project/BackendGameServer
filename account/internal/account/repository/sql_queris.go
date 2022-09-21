package repository

const (
	initUUIDExtension = `
     CREATE EXTENSION IF NOT EXISTS citext;
	 CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	initAccountTable = `
	  CREATE TABLE IF NOT EXISTS accounts (
      account_id UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
      player_id VARCHAR(11) NOT NULL UNIQUE CHECK ( player_id <> '' ),
      username VARCHAR(255) NOT NULL UNIQUE CHECK ( username <> '' ),
      email VARCHAR(320) NOT NULL UNIQUE CHECK ( email ~ '^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$' ),
      password_hashed varchar(255) NOT NULL,
      is_ban BOOLEAN DEFAULT FALSE, 
      created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	initFriendshipTable = `
	CREATE TABLE IF NOT EXISTS friendship(
	    account_id UUID,
	    friend_id  UUID,
	    PRIMARY KEY (account_id, friend_id),
	    CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts(account_id),
	    CONSTRAINT fr_friend FOREIGN KEY (friend_id) REFERENCES accounts(account_id)
	);`

	initAllTable = initUUIDExtension + initAccountTable + initFriendshipTable

	createAccountQuery = `INSERT INTO accounts (account_id, player_id, username, email, password_hashed, is_ban, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, now(), now()) RETURNING account_id, player_id, username, email, password_hashed, is_ban, created_at, updated_at`

	updateAccountQuery = `UPDATE accounts a SET 
                      username=COALESCE(NULLIF($1, ''), username), 
                      email=COALESCE(NULLIF($2, ''), email), 
                      updated_at = now()
                      WHERE account_id=$3
                      RETURNING account_id, player_id, username, email, password_hashed, is_ban, created_at, updated_at`

	searchAccountQuery = `SELECT count(*) over() as total, a.account_id, a.player_id, a.username, a.email, a.password_hashed, a.is_ban, a.created_at, a.updated_at
	FROM accounts a WHERE a.username ILIKE '%' + $1 + '%' OR a.player_id = $1 ORDER BY $2 LIMIT $3 OFFSET $4`

	getAccountByIdQuery = `SELECT a.account_id, a.player_id, a.username, a.email, a.password_hashed, a.is_ban, a.created_at, a.updated_at
	FROM accounts a WHERE a.account_id = $1`

	getAccountByUsernameQuery = `SELECT a.account_id, a.player_id, a.username, a.email, a.password_hashed, a.is_ban, a.created_at, a.updated_at
	FROM accounts a WHERE a.username = $1`

	getAccountByEmailQuery = `SELECT a.account_id, a.player_id, a.username, a.email, a.password_hashed, a.is_ban, a.created_at, a.updated_at
	FROM accounts a WHERE a.email = $1`

	deleteAccountByIdQuery = `DELETE FROM accounts WHERE account_id = $1`
)
