package repository

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type postgresRepository struct {
	log logger.Logger
	db  *sqlx.DB
}

func NewAccountRepository(log logger.Logger, db *sqlx.DB) (*postgresRepository, error) {
	_, err := db.ExecContext(context.Background(), initAllTable)
	if err != nil {
		return nil, err
	}
	return &postgresRepository{
		log: log,
		db:  db,
	}, nil
}
func (a *postgresRepository) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	var created models.Account
	if err := a.db.QueryRowContext(ctx, createAccountQuery, &account.AccountID, &account.PlayerID, &account.Username, &account.Email, &account.PasswordHashed, &account.IsBan).Scan(
		&created.AccountID,
		&created.PlayerID,
		&created.Username,
		&created.Email,
		&created.PasswordHashed,
		&created.IsBan,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow")
	}

	return &created, nil
}

func (a *postgresRepository) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	var prod models.Account
	if err := a.db.QueryRowContext(
		ctx,
		updateAccountQuery,
		&account.Username,
		&account.Email,
		&account.PasswordHashed,
		&account.IsBan,
		&account.AccountID,
	).Scan(&prod.AccountID, &prod.PlayerID, &prod.Username, &prod.Email, &prod.PasswordHashed, &prod.IsBan, &prod.CreatedAt, &prod.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &prod, nil
}

func (a *postgresRepository) DeleteAccount(ctx context.Context, uuid uuid.UUID) error {
	_, err := a.db.ExecContext(ctx, deleteAccountByIdQuery, uuid)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}

	return nil
}
func (a *postgresRepository) GetAccountById(ctx context.Context, uuid uuid.UUID) (*models.Account, error) {
	var account models.Account
	if err := a.db.QueryRowContext(ctx, getAccountByIdQuery, uuid).Scan(
		&account.AccountID,
		&account.PlayerID,
		&account.Username,
		&account.Email,
		&account.PasswordHashed,
		&account.IsBan,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &account, nil
}

func (a *postgresRepository) GetAccountByEmailOrUsername(ctx context.Context, email string, username string) (*models.Account, error) {
	var account models.Account
	if username != "" {
		if err := a.db.QueryRowContext(ctx, getAccountByUsernameQuery, username).Scan(
			&account.AccountID,
			&account.PlayerID,
			&account.Username,
			&account.Email,
			&account.PasswordHashed,
			&account.IsBan,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "Scan GetAccountByUsername")
		}
	} else if email != "" {
		if err := a.db.QueryRowContext(ctx, getAccountByEmailQuery, email).Scan(
			&account.AccountID,
			&account.PlayerID,
			&account.Username,
			&account.Email,
			&account.PasswordHashed,
			&account.IsBan,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "Scan GetAccountByEmail")
		}
	}
	return &account, nil
}

func (a *postgresRepository) Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.AccountsList, error) {
	return nil, nil
}
