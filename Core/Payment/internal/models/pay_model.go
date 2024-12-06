package models

type PayInfo struct {
	Cost     int  `json:"cost"`
	CardInfo Card `json:"card_info"`
}
