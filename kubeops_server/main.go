package main

import (
	"github.com/gin-gonic/gin"
	"kubeops/config"
	"kubeops/controller"
	"kubeops/db"
	"kubeops/middle"
	"kubeops/service"
	"net/http"
)

func main() {
	//初始化k8s
	service.K8s.Init()
	//初始化db
	db.Init()
	//初始化redis
	db.InitRedis()
	//初始化 gin
	r := gin.Default()
	//使用全局中间件
	r.Use(middle.CrosHandler()) //跨域中间件
	r.Use(middle.JwtAuth())     //jwt 验证中间件
	//初始化路由
	controller.Router.InitApiRouter(r)
	//启动ws 服务并监听 ws 端口
	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		_ = http.ListenAndServe(":8081", nil)
	}()
	// 启动gin 服务
	_ = r.Run(config.ListenAddr) // 监听并在 0.0.0.0:9090 上启动服务
	// 关闭数据库连接
	_ = db.Close()
	// 关闭redis连接
	_ = db.CloseRedis()
}
