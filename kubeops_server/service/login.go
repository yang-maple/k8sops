package service

import (
	"errors"
	"kubeops/dao"
	"kubeops/utils"
)

var Login login

type login struct {
}

type LoginInfo struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	CaptchaId       string `json:"captchaId"`
	CaptchaSolution string `json:"captchaSolution"`
}

// VerifyUserInfo 验证用户身份信息
func (l *login) VerifyUserInfo(info *LoginInfo) (token string, err error) {
	err = dao.Login.VerifyUser(info.Username, info.Password)
	if err != nil {
		return "", errors.New(err.Error())
	}
	token, err = utils.CreateJwtToken(info.Username, info.Password, utils.UserExpireDuration, utils.UserSecret)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return token, nil
}
