package model

import (
	"time"
)

type User struct {
	Id        uint `gorm:"primaryKey"`
	Username  string
	Password  string
	Email     string
	Status    int
	CreatedAt time.Time
}

func (User) TableName() string {
	return "user"
}
