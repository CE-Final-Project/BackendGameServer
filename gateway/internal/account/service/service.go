package service

import (
	accountService "github.com/ce-final-project/backend_game_server/account/proto"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/account/commands"
	"github.com/ce-final-project/backend_game_server/gateway/internal/account/queries"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
)

type AccountService struct {
	Commands *commands.AccountCommands
	Queries  *queries.AccountQueries
}

func NewAccountService(
	log logger.Logger,
	cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	accService accountService.AccountServiceClient,
) *AccountService {

	updateAccountCmdHandler := commands.NewUpdateAccountHandler(log, cfg, kafkaProducer)
	changePasswordCmdHandler := commands.NewChangePasswordHandler(log, cfg, kafkaProducer)
	banAccountCmdHandler := commands.NewBanAccountHandler(log, cfg, kafkaProducer)

	getAccountByIdQueryHandler := queries.NewGetAccountByIdHandler(log, cfg, accService)
	searchAccountQueryHandler := queries.NewSearchAccountHandler(log, cfg, accService)

	accountCommands := commands.NewAccountCommands(updateAccountCmdHandler, changePasswordCmdHandler, banAccountCmdHandler)
	accountQueries := queries.NewAccountQueries(getAccountByIdQueryHandler, searchAccountQueryHandler)

	return &AccountService{
		Commands: accountCommands,
		Queries:  accountQueries,
	}
}
