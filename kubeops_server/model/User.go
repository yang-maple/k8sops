package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (User) TableName() string {
	return "user"
}
