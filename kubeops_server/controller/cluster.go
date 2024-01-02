package controller

import (
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"net/http"
	"strconv"
)

var Cluster cluster

type cluster struct{}

// 默认配置文件存放路径
const configfileurl = "./static/configfile/"

// Create 创建集群
func (clt *cluster) Create(c *gin.Context) {
	//绑定参数信息
	params := new(struct {
		ClusterName string `form:"cluster_name"`
		ClusterType string `form:"cluster_type"`
	})
	_ = c.Bind(&params)
	//捕获文件信息
	file, _ := c.FormFile("config_file")
	//重新dir url路径
	uuid := c.Request.Header.Get("Uuid")
	dir := configfileurl + uuid
	//获取用户id
	userid, _ := strconv.Atoi(uuid)
	//创建集群
	err := service.Cluster.Create(dir, params.ClusterName, params.ClusterType, userid, file, c)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "集群创建失败" + err.Error(),
		})
		return
	}
	//创建集群
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "集群" + params.ClusterName + "创建成功",
	})
}

// List 获取集群列表
func (clt *cluster) List(c *gin.Context) {
	userid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Cluster.List(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "获取集群列表失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取集群列表成功",
		"data": data,
	})
}

// Change 更换集群信息
func (clt *cluster) Change(c *gin.Context) {
	userid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	param := new(struct {
		ClusterName string `json:"cluster_name"`
	})
	_ = c.ShouldBindJSON(&param)
	err := service.Cluster.Change(param.ClusterName, userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "应用失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "已切换集群  " + param.ClusterName,
	})
}

// Delete 删除集群
func (clt *cluster) Delete(c *gin.Context) {
	userid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	param := new(struct {
		ClusterName string `form:"cluster_name"`
	})
	_ = c.Bind(&param)
	err := service.Cluster.Delete(param.ClusterName, userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "删除失败" + err.Error(),
		})
		return
	}
	//删除成功返回
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "集群" + param.ClusterName + "删除成功",
	})
}
