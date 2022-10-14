package dto

type LoginAccount struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginAccountResponse struct {
	Account Account `json:"account"`
	Token   string  `json:"token" validate:"required"`
}
