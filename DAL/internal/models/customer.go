package models

import (
	"fmt"
)

type Customer struct {
	Id       int    `json:"id" bson:"id" gorm:"primaryKey"`
	Login    string `json:"login" bson:"login" gorm:"unique;not null"`
	Password string `json:"password" bson:"password" gorm:"not null"`
	Role     string `json:"role" bson:"role" gorm:"default:'user'"`
	Surname  string `json:"surname" bson:"surname" gorm:"column:surname;not null"`
	Name     string `json:"name" bson:"name" gorm:"column:name;not null"`
}

func (c Customer) String() string {
	return fmt.Sprintf("Customer ID: %d, Login: %s, Name: %s %s, Role: %s",
		c.Id, c.Login, c.Name, c.Surname, c.Role)
}
