package config

import "time"

const (
	ListenAddr   = "0.0.0.0:8090"
	Logtaillimit = 2000

	// DbType 数据库配置
	DbType   = "mysql"
	DbHost   = "10.1.131.32"
	DbPort   = 3306
	DbName   = "my_db"
	DbUser   = "root"
	DbPasswd = "123456"

	// LogMode 打印mysql. debug 日志
	LogMode = true

	// MaxIdleConns 连接池的配置
	MaxIdleConns = 10
	MaxOpenConns = 100
	MaxLifeTime  = 30 * time.Second
)
