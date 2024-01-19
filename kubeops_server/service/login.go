package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"kubeops/dao"
	"kubeops/utils"
)

var Login login

type login struct {
}

type LoginInfo struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaId   string `json:"captcha_id"`
	VerifyValue string `json:"verify_value"`
}

// VerifyUserInfo 验证用户身份信息
func (l *login) VerifyUserInfo(info *LoginInfo) (uid int, token, name string, err error) {
	uuid, dir, name, err := dao.Login.VerifyUser(info.Username, info.Password)
	if err != nil {
		logger.Info("用户名或密码错误" + err.Error())
		return 0, "", "", errors.New(err.Error())
	}
	//初始化k8s集群
	if dir != "" {
		fmt.Println(dir, uuid)
		K8s.ConfigDir[uuid] = &dir
		//err = K8s.Init(dir, uuid)
		err = K8s.Init(uuid)
		if err != nil {
			logger.Info("初始化k8s集群失败" + err.Error())
			return 0, "", "", errors.New(err.Error())
		}
	}
	token, err = utils.CreateJwtToken(uuid, info.Username, utils.UserExpireDuration, utils.UserSecret)
	if err != nil {
		logger.Info("生成token失败" + err.Error())
		return 0, "", "", errors.New(err.Error())
	}

	return uuid, token, name, nil
}
