package dto

import uuid "github.com/satori/go.uuid"

type LoginAccount struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginAccountResponse struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	PlayerID  string    `json:"player_id" validate:"required"`
	Token     string    `json:"token" validate:"required"`
}
