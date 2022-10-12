package repository

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
)

type AccountRepository interface {
	SearchAccount(ctx context.Context, search string, pagination *utils.Pagination) (*models.AccountsList, error)
	GetAccountByID(ctx context.Context, accountID uint64) (*models.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (*models.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)

	InsertAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	DeleteAccountByID(ctx context.Context, accountID uint64) error
}

type CacheRepository interface {
	PutAccount(ctx context.Context, key string, account *models.Account)
	PutKeyReference(ctx context.Context, key string, targetKey string)
	GetAccount(ctx context.Context, key string) (*models.Account, error)
	GetAccountByKeyReference(ctx context.Context, key string) (*models.Account, error)
	DelAccount(ctx context.Context, key string)
	DelAllAccounts(ctx context.Context)
}
