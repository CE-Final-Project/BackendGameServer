package dto

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type BanAccount struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
}

type BanAccountResponse struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}
