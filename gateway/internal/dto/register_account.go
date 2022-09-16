package dto

import uuid "github.com/satori/go.uuid"

type RegisterAccountDto struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterAccountResponseDto struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	PlayerID  string    `json:"player_id" validate:"required"`
	Token     string    `json:"token" validate:"required"`
}
