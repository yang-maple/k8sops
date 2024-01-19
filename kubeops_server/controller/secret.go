package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type secret struct{}

var Secrets secret

// GetSecretList   获取 Secret 列表
func (s *secret) GetSecretList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
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
	data, err := service.Secrets.GetSecretList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取Secret失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// GetSecretDetail    获取 Secret 详情
func (s *secret) GetSecretDetail(c *gin.Context) {
	params := new(struct {
		SecretName string `form:"secret_name"`
		Namespace  string `form:"namespace"`
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
	detail, err := service.Secrets.GetSecretDetail(params.Namespace, params.SecretName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取Secret 详情失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": detail,
	})
}

// DelSecret    删除 Secret 资源
func (s *secret) DelSecret(c *gin.Context) {
	params := new(struct {
		SecretName string `form:"secret_name"`
		Namespace  string `form:"namespace"`
	})
	err := c.Bind(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Secrets.DelSecret(params.Namespace, params.SecretName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "删除 Secret 失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除 Secret 成功",
	})
}

// CreateSecret   创建 Secret 资源
func (s *secret) CreateSecret(c *gin.Context) {
	params := new(struct {
		Data *service.CreateSecret `json:"data"`
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
	err = service.Secrets.CreateSecret(params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "创建 Secret 失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "创建 Services 成功",
	})
}

// UpdateSecret   更新 Secret 资源
func (s *secret) UpdateSecret(c *gin.Context) {
	params := new(struct {
		Namespace string         `json:"namespace"`
		Data      *corev1.Secret `json:"data"`
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
	err = service.Secrets.UpdateSecret(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "更新 Secret 失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新 Secret 成功",
	})
}
