package entity

type User struct {
	Dates
	TwitchID          int64  `json:"twitchId"`
	TwitchUsername    string `json:"twitchUsername"`
	TelegramID        int64  `json:"telegramId"`
	TelegramUsername  string `json:"telegramUsername"`
	SubchatTelegramID int64  `json:"subchatTelegramId"`
}
