package commands

import uuid "github.com/satori/go.uuid"

type AuthCommands struct {
	Register RegisterCmdHandler
}

func NewAuthCommands(register RegisterCmdHandler) *AuthCommands {
	return &AuthCommands{
		Register: register,
	}
}

type RegisterCommand struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	PlayerID  string    `json:"player_id,omitempty" validate:"required,max=11"`
	Username  string    `json:"username,omitempty" validate:"required"`
	Email     string    `json:"email,omitempty" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required"`
}

func NewRegisterCommand(username, email, password string) *RegisterCommand {
	return &RegisterCommand{
		Username: username,
		Email:    email,
		Password: password,
	}
}
