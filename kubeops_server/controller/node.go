package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubeops/service"
	"net/http"
	"strconv"
)

type node struct{}

var Node node

func (n *node) GetNodeList(c *gin.Context) {
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Node.GetNodeList(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Node.GetNodeDetail(params.NodeName, uuid)
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

// SetNodeSchedule 设置节点状态
func (n *node) SetNodeSchedule(c *gin.Context) {
	//获取节点名称以及状态
	params := new(struct {
		NodeName string `json:"node_name"`
		Status   bool   `json:"status"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Node.SetNodeSchedule(params.NodeName, params.Status, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "设置失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": nil,
	})
}
