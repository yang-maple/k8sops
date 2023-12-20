package service

import (
	"errors"
	"kubeops/dao"
	"kubeops/utils"
)

var ResetPassword resetPassword

type resetPassword struct{}

// VerifyIdentity 验证用户身份
func (rs *resetPassword) VerifyIdentity(email string) (err error) {
	// 判断用户是否存在
	if err := dao.ResetPassword.VerifyIdentity(email); err != nil {
		return err
	}
	//组装数据，生成随机数
	verity := emailVerity{
		Email: email + "_rsp",
		Value: utils.GenRandNum(6),
	}
	// 将随机数缓存至redis
	err = dao.RdbValue.SetValue(verity.Email, verity.Value, 90)
	if err != nil {
		return err
	}
	// 随机数验证码发送邮件
	err = utils.Emails(email, utils.FormatEmailBody("view/verfityidentify.html", verity), "KubeOps找回密码")
	if err != nil {
		return err
	}
	return nil
}

// Reset  重置密码
func (rs *resetPassword) Reset(email, verifyCode, password string) (err error) {
	//读取缓存中的验证码
	value, err := dao.RdbValue.GetValue(email + "_rsp")
	if err != nil {
		return err
	}
	if value != verifyCode {
		return errors.New("验证码错误")
	}
	//重置密码
	if err := dao.ResetPassword.Reset(email, password); err != nil {
		return err
	}
	//重置成功后删除缓存中的验证码
	_ = dao.RdbValue.DelValue(email + "_rsp")
	return nil

}
