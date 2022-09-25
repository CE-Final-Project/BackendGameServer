package commands

import "github.com/ce-final-project/backend_game_server/gateway/internal/dto"

type AuthCommands struct {
	RegisterAccount RegisterAccountCmdHandler
	LoginAccount    LoginAccountCmdHandler
}

func NewAuthCommands(registerAccount RegisterAccountCmdHandler, loginAccount LoginAccountCmdHandler) *AuthCommands {
	return &AuthCommands{
		RegisterAccount: registerAccount,
		LoginAccount:    loginAccount,
	}
}

type RegisterAccountCommand struct {
	RegisterDto *dto.RegisterAccount
}

func NewRegisterAccountCommand(registerDto *dto.RegisterAccount) *RegisterAccountCommand {
	return &RegisterAccountCommand{RegisterDto: registerDto}
}

type LoginAccountCommand struct {
	LoginDto *dto.LoginAccount
}

func NewLoginAccountCommand(loginDto *dto.LoginAccount) *LoginAccountCommand {
	return &LoginAccountCommand{LoginDto: loginDto}
}
