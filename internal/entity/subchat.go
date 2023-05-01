package entity

type Subchat struct {
	Dates
	OwnerTwitchID     int64 `json:"ownerTwitchId"`
	SubchatTelegramID int64 `json:"subchatTelegramId"`
	Disabled          bool  `json:"disabled"`
}
