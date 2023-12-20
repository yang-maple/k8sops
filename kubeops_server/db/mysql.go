package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
	"kubeops/config"
	"kubeops/model"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

func Init() {
	//判断是否已经初始化 如果是直接 return
	if isInit {
		return
	}
	// 组装数据
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPasswd,
		config.DbHost,
		config.DbPort,
		config.DbName)
	GORM, err = gorm.Open(config.DbType, dsn)
	if err != nil {
		panic("数据库连接失败" + err.Error())
	}
	logger.Info("数据库初始化成功")
	GORM.LogMode(config.LogMode)

	//打开连接池
	//最大空闲连接数
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	//最大连接池数量
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	// 连接池存活时间
	GORM.DB().SetConnMaxLifetime(config.MaxLifeTime)
	isInit = true
	GORM.AutoMigrate(&model.Workflow{})
	GORM.AutoMigrate(&model.User{})
}

// Close 关闭数据库连接
func Close() error {
	return GORM.Close()
}
