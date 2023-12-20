package service

import (
	"errors"
	"kubeops/dao"
	"kubeops/utils"
)

var Register register

type register struct{}

type emailVerity struct {
	Email string `json:"email"`
	Value string `json:"value"`
}

type RegisterInfo struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	VerifyCode string `json:"verify_code"`
}

func (r *register) RegisterUser(info *RegisterInfo) (err error) {
	//从缓存中获取验证码
	value, err := dao.RdbValue.GetValue(info.Email + "_rgs")
	if err != nil {
		return err
	}
	//验证 验证码
	if value != info.VerifyCode {
		return errors.New("验证码错误")
	}
	//注册用户 并检测用户是否已存在
	err = dao.Register.RegisterUser(info.Username, info.Password, info.Email)
	if err != nil {
		return err
	}
	//创建成功后删除缓存中的验证码
	_ = dao.RdbValue.DelValue(info.Email + "_rgs")
	return nil

}

// SendEmail 发送邮件
func (r *register) SendEmail(email string) (err error) {
	err = dao.Register.SendEmail(email)
	if err != nil {
		return err
	}
	// 生成随机数
	verity := emailVerity{
		Email: email + "_rgs",
		Value: utils.GenRandNum(6),
	}
	// 将随机数缓存至redis
	err = dao.RdbValue.SetValue(verity.Email, verity.Value, 90)
	if err != nil {
		return err
	}
	// 随机数验证码发送邮件
	err = utils.Emails(email, utils.FormatEmailBody("view/verfityemail.html", verity), "KubeOps注册")
	if err != nil {
		return err
	}
	return nil
}
