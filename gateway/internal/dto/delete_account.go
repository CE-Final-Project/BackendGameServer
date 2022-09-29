package dto

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type DeleteAccountById struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
}

type DeleteAccountByIdResponse struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	DeletedAt time.Time `json:"deleted_at"`
}
