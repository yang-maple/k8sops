package model

import (
	"time"
)

type User struct {
	Id           uint `gorm:"primaryKey"`
	Username     string
	Password     string
	Email        string
	Status       int
	CreatedAt    time.Time
	ClustersInfo []ClusterInfo `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "user"
}
