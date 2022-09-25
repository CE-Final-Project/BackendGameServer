package dto

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type UpdateAccount struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty" validate:"email"`
}

type UpdateAccountResponse struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}
