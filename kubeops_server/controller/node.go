package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubeops/service"
	"net/http"
)

type node struct{}

var Node node

func (n *node) GetNodeList(c *gin.Context) {
	data, err := service.Node.GetNodelist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

func (n *node) GetNodeDeTal(c *gin.Context) {
	params := new(struct {
		NodeName string `form:"node_name"`
	})
	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	data, err := service.Node.GetNodeDetail(params.NodeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": data,
	})
}
