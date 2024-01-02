package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type configmap struct{}

var Configmaps configmap

// GetConfigmapList  获取 Configmap 列表
func (cm *configmap) GetConfigmapList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))

	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	data, err := service.Configmaps.GetCmList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取Configmap失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// GetConfigmapDetail   获取 Configmap 详情
func (cm *configmap) GetConfigmapDetail(c *gin.Context) {
	params := new(struct {
		ConfigmapName string `form:"configmap_name"`
		Namespace     string `form:"namespace"`
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
	detail, err := service.Configmaps.GetCmDetail(params.Namespace, params.ConfigmapName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取Configmap 详情失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": detail,
	})
}

// DelConfigmap    删除 Configmap 资源
func (cm *configmap) DelConfigmap(c *gin.Context) {
	params := new(struct {
		ConfigmapName string `form:"configmap_name"`
		Namespace     string `form:"namespace"`
	})
	err := c.Bind(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Configmaps.DelCm(params.Namespace, params.ConfigmapName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "删除 Configmap 失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除 Configmap 成功",
	})
}

// CreateConfigmap  创建 Configmap 资源
func (cm *configmap) CreateConfigmap(c *gin.Context) {
	params := new(struct {
		Data *service.CreateConfig `json:"data"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Configmaps.CreateCm(params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建 Configmap 失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "创建 Configmap 成功",
	})
}

// UpdateConfigmap 更新Configmap 资源
func (cm *configmap) UpdateConfigmap(c *gin.Context) {
	params := new(struct {
		Namespace string            `json:"namespace"`
		Data      *corev1.ConfigMap `json:"data"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "绑定参数失败",
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Configmaps.UpdateCm(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "更新 Configmap 失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新 Configmap 成功",
	})
}
