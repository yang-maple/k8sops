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

func (l *login) VerifyUser(username, password string) (uid int, dir, clusterName string, err error) {
	var userInfo model.User
	var cluster model.ClusterInfo
	tx := db.GORM.Where("username = ?", username).Where("password = ?", utils.PasswordToMd5(password)).First(&userInfo)
	if tx.Error != nil && tx.Error.Error() != "record not found" || tx.RowsAffected == 0 {
		return 0, "", "", errors.New("用户名或密码错误")
	}
	if userInfo.Status == 0 {
		return 0, "", "", errors.New("用户已被禁用")
	}
	//判断用户是否有集群
	tx = db.GORM.Where("user_id =? And status = ?", userInfo.Id, 1).First(&cluster)
	if tx.Error != nil && tx.Error.Error() != "record not found" || tx.RowsAffected == 0 {
		return int(userInfo.Id), "", "", nil
	}
	return int(userInfo.Id), cluster.Dir + "/" + cluster.ClusterName + "_" + cluster.Type + ".conf", cluster.ClusterName, nil
}
