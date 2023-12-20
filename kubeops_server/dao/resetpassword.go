package dao

import (
	"errors"
	"kubeops/db"
	"kubeops/model"
	"kubeops/utils"
)

var ResetPassword resetPassword

type resetPassword struct {
}

// VerifyIdentity 验证身份,是否有该用户的注册邮箱
func (rs *resetPassword) VerifyIdentity(email string) (err error) {
	var userInfo model.User
	tx := db.GORM.Where("email = ?", email).First(&userInfo)
	if tx.RowsAffected == 0 && tx.Error.Error() == "record not found" {
		return errors.New("邮箱未被注册")
	}
	return nil

}

// Reset 用户重置密码
func (rs *resetPassword) Reset(email, password string) (err error) {
	var userInfo model.User
	tx := db.GORM.Model(&userInfo).Where("email = ?", email).Update("password", utils.PasswordToMd5(password))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
