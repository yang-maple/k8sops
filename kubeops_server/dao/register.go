package dao

import (
	"errors"
	"kubeops/db"
	"kubeops/model"
	"kubeops/utils"
)

var Register register

type register struct{}

func (r *register) RegisterUser(username, password, email string) (err error) {
	var userInfo model.User
	tx := db.GORM.Where("username = ?", username).First(&userInfo)
	if tx.Error != nil && tx.Error.Error() != "record not found" || tx.RowsAffected != 0 {
		return errors.New("用户名已存在")
	}
	userInfo.Username = username
	userInfo.Password = utils.PasswordToMd5(password)
	userInfo.Email = email
	userInfo.Status = 1
	tx = db.GORM.Create(&userInfo)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
