package service

import (
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/account/commands"
	"github.com/ce-final-project/backend_game_server/account/internal/account/queries"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
)

type AccountService struct {
	Commands *commands.AccountCommands
	Queries  *queries.AccountQueries
}

func NewAccountService(
	log logger.Logger,
	cfg *config.Config,
	postgresRepo repository.Repository,
	redisRepo repository.CacheRepository,
) *AccountService {

	createAccountHandler := commands.NewCreateAccountHandler(log, cfg, postgresRepo, redisRepo)
	deleteAccountCmdHandler := commands.NewDeleteAccountCmdHandler(log, cfg, postgresRepo, redisRepo)
	updateAccountCmdHandler := commands.NewUpdateAccountCmdHandler(log, cfg, postgresRepo, redisRepo)
	changePasswordCmdHandler := commands.NewChangePasswordCmdHandler(log, cfg, postgresRepo, redisRepo)
	banAccountByIdCmdHandler := commands.NewBanAccountByIdCmdHandler(log, cfg, postgresRepo, redisRepo)

	getAccountByIdHandler := queries.NewGetAccountByIdHandler(log, cfg, postgresRepo, redisRepo)
	searchAccountHandler := queries.NewSearchAccountHandler(log, cfg, postgresRepo, redisRepo)

	accountCommands := commands.NewAccountCommands(createAccountHandler, updateAccountCmdHandler, changePasswordCmdHandler, banAccountByIdCmdHandler, deleteAccountCmdHandler)
	accountQueries := queries.NewAccountQueries(getAccountByIdHandler, searchAccountHandler)

	return &AccountService{Commands: accountCommands, Queries: accountQueries}
}
