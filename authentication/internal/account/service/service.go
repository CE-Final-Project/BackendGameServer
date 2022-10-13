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

	createAccountCmdHandler := commands.NewCreateAccountCmdHandler(log, cfg, accountRepo, cacheRepo)
	deleteAccountCmdHandler := commands.NewDeleteAccountCmdHandler(log, cfg, accountRepo, cacheRepo)
	changePasswordCmdHandler := commands.NewChangePasswordCmdHandler(log, cfg, accountRepo, cacheRepo)

	getAccountByIdQueryHandler := queries.NewGetAccountByIdQueryHandler(log, cfg, accountRepo, cacheRepo)
	getAccountByUsernameQueryHandler := queries.NewGetAccountByUsernameQueryHandler(log, cfg, accountRepo, cacheRepo)
	getAccountByEmailQueryHandler := queries.NewGetAccountByEmailQueryHandler(log, cfg, accountRepo, cacheRepo)
	searchAccountQueryHandler := queries.NewSearchAccountQueryHandler(log, cfg, accountRepo, cacheRepo)

	accountCommands := commands.NewAccountCommands(createAccountCmdHandler, changePasswordCmdHandler, deleteAccountCmdHandler)
	accountQueries := queries.NewAccountQueries(getAccountByIdQueryHandler, getAccountByUsernameQueryHandler, getAccountByEmailQueryHandler, searchAccountQueryHandler)

	return &AccountService{Commands: accountCommands, Queries: accountQueries}
}
