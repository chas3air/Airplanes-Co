package models

type Card struct {
	Number string `json:"number" xml:"number"` // Номер карты
	//ExpirationDate time.Time `json:"expiration_date" xml:"expiration_date"` // Дата истечения срока
	//CVV     string `json:"cvv" xml:"cvv"`         // Код безопасности карты
	Balance int `json:"balance" xml:"balance"` // Баланс на карте
}
