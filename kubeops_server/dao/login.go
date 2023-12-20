package dao

import (
	"errors"
	"kubeops/db"
	"kubeops/model"
	"kubeops/utils"
)

var Login login

type login struct {
}

func (l *login) VerifyUser(username, password string) (err error) {
	var userInfo model.User
	tx := db.GORM.Where("username = ?", username).Where("password = ?", utils.PasswordToMd5(password)).First(&userInfo)
	if tx.Error != nil && tx.Error.Error() != "record not found" || tx.RowsAffected == 0 {
		return errors.New("用户名或密码错误")
	}
	return nil
}
