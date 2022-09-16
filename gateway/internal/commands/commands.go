package commands

import "github.com/ce-final-project/backend_game_server/gateway/internal/dto"

type AuthCommands struct {
	RegisterAccount RegisterAccountCmdHandler
	UpdateAccount   UpdateAccountCmdHandler
}

func NewAuthCommands(registerAccount RegisterAccountCmdHandler, updateAccount UpdateAccountCmdHandler) *AuthCommands {
	return &AuthCommands{
		RegisterAccount: registerAccount,
		UpdateAccount:   updateAccount,
	}
}

type RegisterAccountCommand struct {
	RegisterDto *dto.RegisterAccountDto
}

func NewRegisterAccountCommand(registerDto *dto.RegisterAccountDto) *RegisterAccountCommand {
	return &RegisterAccountCommand{registerDto}
}

type UpdateAccountCommand struct {
	UpdateDto *dto.UpdateAccountDto
}

func NewUpdateAccountCommand(updateDto *dto.UpdateAccountDto) *UpdateAccountCommand {
	return &UpdateAccountCommand{updateDto}
}
