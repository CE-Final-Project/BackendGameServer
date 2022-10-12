package service

import (
	"github.com/ce-final-project/backend_game_server/authentication/config"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/commands"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/queries"
	"github.com/ce-final-project/backend_game_server/authentication/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
)

type AccountService struct {
	Commands *commands.AccountCommands
	Queries  *queries.AccountQueries
}

func NewAccountService(
	log logger.Logger,
	cfg *config.Config,
	accountRepo repository.AccountRepository,
	cacheRepo repository.CacheRepository,
) *AccountService {

	createAccountHandler := commands.NewCreateAccountHandler(log, cfg, accountRepo, cacheRepo)
	deleteAccountCmdHandler := commands.NewDeleteAccountCmdHandler(log, cfg, accountRepo, cacheRepo)
	changePasswordCmdHandler := commands.NewChangePasswordCmdHandler(log, cfg, accountRepo, cacheRepo)

	getAccountByIdHandler := queries.NewGetAccountByIdHandler(log, cfg, accountRepo, cacheRepo)
	getAccountByUsernameHandler := queries.NewGetAccountByUsernameHandler(log, cfg, accountRepo, cacheRepo)
	getAccountByEmailHandler := queries.NewGetAccountByEmailHandler(log, cfg, accountRepo, cacheRepo)
	searchAccountHandler := queries.NewSearchAccountHandler(log, cfg, accountRepo, cacheRepo)

	accountCommands := commands.NewAccountCommands(createAccountHandler, changePasswordCmdHandler, deleteAccountCmdHandler)
	accountQueries := queries.NewAccountQueries(getAccountByIdHandler, getAccountByUsernameHandler, getAccountByEmailHandler, searchAccountHandler)

	return &AccountService{Commands: accountCommands, Queries: accountQueries}
}
