package models

import (
	"fmt"
)

type Customer struct {
	Id       int    `json:"id" bson:"id"`
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`
	Surname  string `json:"surname" bson:"surname"`
	Name     string `json:"name" bson:"name"`
}

func (c Customer) String() string {
	return fmt.Sprintf("Customer ID: %d, Login: %s, Name: %s %s, Role: %s",
		c.Id, c.Login, c.Name, c.Surname, c.Role)
}
