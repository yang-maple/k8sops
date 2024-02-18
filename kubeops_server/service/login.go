package service

import (
	"fmt"
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
		utils.Logger.Error("User identity verification failed")
		return 0, "", "", err
	}
	utils.Logger.Info("User identity verification succeeded")
	//初始化k8s集群
	if dir != "" {
		fmt.Println(dir, uuid)
		K8s.ConfigDir[uuid] = &dir
		//err = K8s.Init(dir, uuid)
		err = K8s.Init(uuid)
		if err != nil {
			utils.Logger.Error("Failed to initialize the k8s cluster:" + err.Error())
			return 0, "", "", err
		}
	}
	token, err = utils.CreateJwtToken(uuid, info.Username, utils.UserExpireDuration, utils.UserSecret)
	if err != nil {
		utils.Logger.Error("CreateJwtToken failed :" + err.Error())
		return 0, "", "", err
	}
	return uuid, token, name, nil
}
