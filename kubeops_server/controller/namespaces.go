package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type namespace struct{}

var Namespace namespace

func (n *namespace) GetNsList(c *gin.Context) {
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Namespace.GetNsList(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "获取namespaces失败",
			"reason": errors.New(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

func (n *namespace) GetNsDetail(c *gin.Context) {
	params := new(struct {
		NamespaceName string `form:"namespace_name"`
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
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Namespace.GetNsDetail(params.NamespaceName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "获取namespaces失败",
			"reason": errors.New(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

func (n *namespace) DelNs(c *gin.Context) {
	params := new(struct {
		NamespaceName string `form:"namespace_name"`
	})
	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "绑定参数失败",
			"reason": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Namespace.DelNs(params.NamespaceName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "名称空间" + params.NamespaceName + "删除失败",
			"reason": errors.New(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "名称空间" + params.NamespaceName + "删除成功",
	})
}

func (n *namespace) CreateNs(c *gin.Context) {
	params := new(struct {
		NamespaceName string `json:"namespace_name"`
	})
	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "绑定参数失败",
			"reason": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Namespace.CreateNs(params.NamespaceName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    params.NamespaceName + "创建失败",
			"reason": errors.New(err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "创建成功",
	})
}

func (n *namespace) UpdateNs(c *gin.Context) {
	params := new(struct {
		Data *corev1.Namespace `json:"data"`
	})
	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("绑定参数失败" + err.Error()),
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Namespace.UpdateNs(params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "更新失败",
			"reason": errors.New(err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})

}
