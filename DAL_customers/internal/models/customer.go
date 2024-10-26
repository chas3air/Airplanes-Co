package models

import (
	"fmt"

	"github.com/google/uuid"
)

// Customer представляет структуру для клиентов.
type Customer struct {
	Id       uuid.UUID `json:"id" bson:"id"` // Использование UUID
	Login    string    `json:"login" bson:"login"`
	Password string    `json:"password" bson:"password"`
	Role     string    `json:"role" bson:"role"`
	Surname  string    `json:"surname" bson:"surname"`
	Name     string    `json:"name" bson:"name"`
}

// String возвращает строковое представление клиента.
func (c Customer) String() string {
	return fmt.Sprintf("Customer ID: %s, Login: %s, Name: %s %s, Role: %s",
		c.Id.String(), c.Login, c.Name, c.Surname, c.Role)
}
