package dto

import (
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
	uuid "github.com/satori/go.uuid"
)

type VerifyTokenResponse struct {
	Valid     bool      `json:"valid"`
	AccountID uuid.UUID `json:"account_id"`
	PlayerID  string    `json:"player_id"`
}

func VerifyTokenResponseFromGRPC(verifyToken *authService.VerifyTokenRes) (*VerifyTokenResponse, error) {
	accountUUID, err := uuid.FromString(verifyToken.GetAccountID())
	if err != nil {
		return nil, err
	}
	return &VerifyTokenResponse{
		Valid:     verifyToken.GetValid(),
		AccountID: accountUUID,
		PlayerID:  verifyToken.GetPlayerID(),
	}, nil
}
