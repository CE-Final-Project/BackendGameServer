package commands_test

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/commands"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestHandler_Handle(t *testing.T) {
	cfg, err := config.InitConfig("../../../config")
	if err != nil {
		log.Fatal(err)
	}
	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("TestAccountService")

	accounts := []*models.Account{
		&models.Account{
			AccountID:      uuid.NewV4(),
			PlayerID:       "test01",
			Username:       "kenta01",
			Email:          "kenta01@gmail.com",
			PasswordHashed: "passwordHashed",
			IsBan:          false,
			UpdatedAt:      time.Now(),
			CreatedAt:      time.Now(),
		},
		&models.Account{
			AccountID:      uuid.NewV4(),
			PlayerID:       "test02",
			Username:       "kenta02",
			Email:          "kenta02@gmail.com",
			PasswordHashed: "passwordHashed",
			IsBan:          false,
			UpdatedAt:      time.Now(),
			CreatedAt:      time.Now(),
		},
		&models.Account{
			AccountID:      uuid.NewV4(),
			PlayerID:       "test03",
			Username:       "kenta03",
			Email:          "kenta03@gmail.com",
			PasswordHashed: "passwordHashed",
			IsBan:          false,
			UpdatedAt:      time.Now(),
			CreatedAt:      time.Now(),
		},
		&models.Account{
			AccountID:      uuid.NewV4(),
			PlayerID:       "test04",
			Username:       "kenta04",
			Email:          "kenta04@gmail.com",
			PasswordHashed: "passwordHashed",
			IsBan:          false,
			UpdatedAt:      time.Now(),
			CreatedAt:      time.Now(),
		},
		&models.Account{
			AccountID:      uuid.NewV4(),
			PlayerID:       "test05",
			Username:       "kenta05",
			Email:          "kenta05@gmail.com",
			PasswordHashed: "passwordHashed",
			IsBan:          false,
			UpdatedAt:      time.Now(),
			CreatedAt:      time.Now(),
		},
	}

	accountList := &models.AccountsList{
		TotalCount: 5,
		TotalPages: 1,
		Page:       0,
		Size:       10,
		HasMore:    false,
		Accounts:   accounts,
	}

	accountRepo := repository.NewAccountRepositoryMock()
	for _, acc := range accounts {
		accountRepo.On("CreateAccount", context.Background(), acc).Return(acc, nil)
		acc.Username = acc.Username + "Updated"
		accountRepo.On("UpdateAccount", context.Background(), acc).Return(acc, nil)
		accountRepo.On("DeleteAccount", context.Background(), acc.AccountID).Return(nil)
		accountRepo.On("GetAccountById", context.Background(), acc.AccountID).Return(acc, nil)
		accountRepo.On("GetAccountByEmailOrUsername", context.Background(), "", acc.Username).Return(acc, nil)
	}
	accountRepo.On("Search", context.Background(), "kenta", utils.NewPaginationQuery(10, 0)).Return(accountList, nil)

	cacheRepo := repository.NewAccountCacheRepositoryMock()
	for _, acc := range accounts {
		cacheRepo.On("GetAccountReference", context.Background(), "auth:account"+acc.Username, acc).Return(acc, nil)
		cacheRepo.On("GetAccount", context.Background(), "auth:account"+acc.AccountID.String(), acc).Return(acc, nil)
	}

	t.Run("Create Account Handler", func(t *testing.T) {
		createAccountHandler := commands.NewCreateAccountHandler(appLogger, cfg, accountRepo, cacheRepo)
		err := createAccountHandler.Handle(context.Background(), commands.NewCreateAccountCommand(accounts[0].AccountID, accounts[0].PlayerID, accounts[0].Username, accounts[0].Email, accounts[0].PasswordHashed, accounts[0].IsBan, accounts[0].CreatedAt, accounts[0].UpdatedAt))
		assert.ErrorIs(t, err, nil)
	})
}
