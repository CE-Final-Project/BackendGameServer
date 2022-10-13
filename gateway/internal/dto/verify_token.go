package dto

import (
	authService "github.com/ce-final-project/backend_game_server/authentication/proto"
)

type VerifyTokenResponse struct {
	Valid     bool   `json:"valid"`
	AccountID uint64 `json:"account_id"`
	PlayerID  string `json:"player_id"`
}

func VerifyTokenResponseFromGRPC(verifyToken *authService.VerifyTokenRes) (*VerifyTokenResponse, error) {
	return &VerifyTokenResponse{
		Valid:     verifyToken.GetValid(),
		AccountID: verifyToken.GetAccountID(),
		PlayerID:  verifyToken.GetPlayerID(),
	}, nil
}
