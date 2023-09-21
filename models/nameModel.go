package models

import (
	"time"
)

type User struct {
	UID                      uint   `json:"id" gorm:"primaryKey:type:uuid"`
	Name                     string `json:"name" validate:"required"`
	Surname                  string `json:"surname" validate:"required"`
	Patronymic               string `json:"patronymic"`
	Age, Gender, Nationality string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

func validateUser(u *User) {

}
