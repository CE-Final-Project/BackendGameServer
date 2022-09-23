package repository

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	DeleteAccount(ctx context.Context, uuid uuid.UUID) error

	GetAccountById(ctx context.Context, uuid uuid.UUID) (*models.Account, error)
	GetAccountByEmailOrUsername(ctx context.Context, email string, username string) (*models.Account, error)
	Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.AccountsList, error)
}

type CacheRepository interface {
	PutAccount(ctx context.Context, key string, account *models.Account)
	PutKeyReference(ctx context.Context, key string, targetKey string)
	GetAccount(ctx context.Context, key string) (*models.Account, error)
	GetAccountReference(ctx context.Context, key string) (*models.Account, error)
	DelAccount(ctx context.Context, key string)
	DelAllAccounts(ctx context.Context)
}
