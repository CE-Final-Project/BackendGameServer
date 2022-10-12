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

type postgresRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

func (p *postgresRepository) SearchAccount(ctx context.Context, search string, pagination *utils.Pagination) (*models.AccountsList, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "postgresRepository.SearchAccount")
	defer span.Finish()

	rows, err := p.db.QueryxContext(ctx, searchAccountQuery, "%"+search+"%", pagination.GetOrderBy(), pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		return nil, errors.Wrap(err, "QueryxContext Search Account")
	}

	var total int64
	accounts := make([]*models.Account, 0, pagination.GetSize())

	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			p.log.WarnMsg("Close rows searchAccount", err)
		}
	}(rows)

	for rows.Next() {
		account := models.Account{}
		if err := rows.Scan(
			&total,
			&account.ID,
			&account.PlayerID,
			&account.Username,
			&account.Email,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "Scan Search Account")
		}
		accounts = append(accounts, &account)
	}
	return models.NewAccountListWithPagination(accounts, total, pagination), nil
}

func (p *postgresRepository) GetAccountByID(ctx context.Context, accountID uint64) (*models.Account, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "postgresRepository.GetAccountByID")
	defer span.Finish()

	var account models.Account
	if err := p.db.QueryRowContext(ctx, getAccountByIdQuery, accountID).Scan(
		&account.ID,
		&account.PlayerID,
		&account.Username,
		&account.Email,
		&account.PasswordHashed,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "repository.GetAccountByID")
	}

	return &account, nil
}

func (p *postgresRepository) GetAccountByUsername(ctx context.Context, username string) (*models.Account, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "postgresRepository.GetAccountByUsername")
	defer span.Finish()

	var account models.Account
	if err := p.db.QueryRowContext(ctx, getAccountByUsernameQuery, username).Scan(
		&account.ID,
		&account.PlayerID,
		&account.Username,
		&account.Email,
		&account.PasswordHashed,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "repository.GetAccountByUsername")
	}

	return &account, nil
}

func (p *postgresRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "postgresRepository.GetAccountByEmail")
	defer span.Finish()

	var account models.Account
	if err := p.db.QueryRowContext(ctx, getAccountByEmailQuery, email).Scan(
		&account.ID,
		&account.PlayerID,
		&account.Username,
		&account.Email,
		&account.PasswordHashed,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "repository.GetAccountByEmail")
	}

	return &account, nil
}

func (p *postgresRepository) InsertAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "postgresRepository.InsertAccount")
	defer span.Finish()

	var created models.Account
	if err := p.db.QueryRowContext(ctx, createAccountQuery, &account.ID, &account.PlayerID, &account.Username, &account.Email, &account.PasswordHashed).Scan(
		&created.ID,
		&created.PlayerID,
		&created.Username,
		&created.Email,
		&created.PasswordHashed,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "repository.InsertAccount")
	}

	return &created, nil
}

func (p *postgresRepository) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "postgresRepository.UpdateAccount")
	defer span.Finish()

	var updated models.Account
	if err := p.db.QueryRowContext(
		ctx,
		updateAccountQuery,
		&account.Username,
		&account.Email,
		&account.PasswordHashed,
		&account.ID,
	).Scan(&updated.ID, &updated.PlayerID, &updated.Username, &updated.Email, &updated.PasswordHashed, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "repository.UpdateAccount")
	}
	return &updated, nil
}

func (p *postgresRepository) DeleteAccountByID(ctx context.Context, accountID uint64) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "postgresRepository.DeleteAccount")
	defer span.Finish()

	row := p.db.QueryRowContext(ctx, deleteAccountByIdQuery, &accountID)
	if row.Err() != nil {
		return errors.Wrap(row.Err(), "repository.UpdateAccount")
	}
	return nil
}

func NewAccountRepository(db *sqlx.DB, log logger.Logger) AccountRepository {
	return &postgresRepository{db: db, log: log}
}
