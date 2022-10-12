package dto

type ChangePasswordAccountDTO struct {
	ID          uint64 `json:"id" validate:"required,numeric"`
	OldPassword string `json:"old_password" validate:"required,min=8,max=100"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}
