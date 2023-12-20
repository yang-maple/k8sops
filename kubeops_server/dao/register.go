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
	//判断用户是否存在
	tx := db.GORM.Where("username = ?", username).First(&userInfo)
	if tx.Error != nil && tx.Error.Error() != "record not found" || tx.RowsAffected != 0 {
		return errors.New("用户名已存在")
	}
	userInfo.Username = username
	userInfo.Password = utils.PasswordToMd5(password)
	userInfo.Email = email
	userInfo.Status = 1
	//注册用户
	tx = db.GORM.Create(&userInfo)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// SendEmail 验证邮箱是否已被注册
func (r *register) SendEmail(email string) (err error) {
	var userInfo model.User
	tx := db.GORM.Where("email = ?", email).First(&userInfo)
	if tx.RowsAffected != 0 {
		return errors.New("邮箱已被注册")
	}
	return nil
}
