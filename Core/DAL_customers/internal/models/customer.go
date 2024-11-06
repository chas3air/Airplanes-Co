package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Customer struct {
	Id       uuid.UUID `json:"id" bson:"id"`
	Login    string    `json:"login" bson:"login"`
	Password string    `json:"password" bson:"password"`
	Role     string    `json:"role" bson:"role"`
	Surname  string    `json:"surname" bson:"surname"`
	Name     string    `json:"name" bson:"name"`
}

func (c Customer) String() string {
	return fmt.Sprintf("Customer ID: %s, Login: %s, Name: %s %s, Role: %s",
		c.Id.String(), c.Login, c.Name, c.Surname, c.Role)
}
