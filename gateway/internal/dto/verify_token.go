package dto

import uuid "github.com/satori/go.uuid"

type VerifyTokenResponseDto struct {
	Valid     bool      `json:"valid"`
	AccountID uuid.UUID `json:"account_id"`
	PlayerID  string    `json:"player_id"`
}
