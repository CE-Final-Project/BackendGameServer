package dto

import uuid "github.com/satori/go.uuid"

type LoginResponseDto struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	PlayerID  string    `json:"player_id" validate:"required,max=11"`
	Token     string    `json:"token" validate:"required"`
}
