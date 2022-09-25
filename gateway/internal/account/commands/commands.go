package commands

import "github.com/ce-final-project/backend_game_server/gateway/internal/dto"

type AccountCommands struct {
	UpdateAccount  UpdateAccountCmdHandler
	ChangePassword ChangePasswordCmdHandler
	BanAccount     BanAccountCmdHandler
}

func NewAccountCommands(updateAccount UpdateAccountCmdHandler, changePassword ChangePasswordCmdHandler, banAccount BanAccountCmdHandler) *AccountCommands {
	return &AccountCommands{
		UpdateAccount:  updateAccount,
		ChangePassword: changePassword,
		BanAccount:     banAccount,
	}
}

type UpdateAccountCommand struct {
	UpdateDto *dto.UpdateAccount
}

func NewUpdateAccountCommand(updateDto *dto.UpdateAccount) *UpdateAccountCommand {
	return &UpdateAccountCommand{updateDto}
}

type ChangePasswordCommand struct {
	ChangePasswordDto *dto.ChangePassword
}

func NewChangePasswordCommand(changePasswordDto *dto.ChangePassword) *ChangePasswordCommand {
	return &ChangePasswordCommand{ChangePasswordDto: changePasswordDto}
}

type BanAccountCommand struct {
	BanAccountDto *dto.BanAccount
}

func NewBanAccountCommand(banAccountDto *dto.BanAccount) *BanAccountCommand {
	return &BanAccountCommand{BanAccountDto: banAccountDto}
}
