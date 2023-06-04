package entity

import "time"

type Token struct {
	TwitchID              int64     `json:"twitchId"`
	TwitchAuthCode        string    `json:"twitchAuthCode"`
	TwitchBearer          string    `json:"twitchBearer"`
	TwitchBearerExpiresAt time.Time `json:"twitchBearerExpiresAt"`
	TwitchRefreshToken    string    `json:"twitchRefreshToken"`
	InviteKey             string    `json:"inviteKey"`
}
