package dto

type FriendInvite struct {
	PlayerID       string `json:"player_id" validate:"required,max=15"`
	FriendPlayerID string `json:"friend_player_id" validate:"required,max=15"`
}
