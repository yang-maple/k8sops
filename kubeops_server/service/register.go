package service

import (
	"fmt"
	"kubeops/dao"
)

var Register register

type register struct{}

type RegisterInfo struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	VerifyCode string `json:"verify_code"`
}

func (r *register) RegisterUser(info *RegisterInfo) (err error) {
	err = dao.Register.RegisterUser(info.Username, info.Password, info.Email)
	if err != nil {
		return err
	}
	return nil

}

func (r *register) EmailVerify() {

}

func (r *register) SendEmail(email string) (err error) {
	fmt.Println(email)
	return err
}
