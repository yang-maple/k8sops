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

type services struct{}

var Services services

// GetServiceList 获取 services 列表
func (s *services) GetServiceList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
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
	data, err := service.Services.GetSvcList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "获取Services失败",
			"reason": errors.New(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// GetServiceDetail 获取 services 详情
func (s *services) GetServiceDetail(c *gin.Context) {
	params := new(struct {
		ServiceName string `form:"service_name"`
		Namespace   string `form:"namespace"`
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
	detail, err := service.Services.GetSvcDetail(params.Namespace, params.ServiceName, uuid)
	if err != nil {
		logger.Info("获取Services 详情失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取Services 详情失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": detail,
	})
}

// DelServices  删除 services 资源
func (s *services) DelServices(c *gin.Context) {
	params := new(struct {
		ServiceName string `form:"service_name"`
		Namespace   string `form:"namespace"`
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
	err = service.Services.DelSvc(params.Namespace, params.ServiceName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "删除 Services 失败",
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除 Services 成功",
	})
}

// CreateService 创建 Services 资源
func (s *services) CreateService(c *gin.Context) {
	createSvc := new(struct {
		Data *service.CreateService `json:"data"`
	})
	err := c.ShouldBindJSON(&createSvc)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Services.CreateSvc(createSvc.Data, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "创建 Services 失败",
			"reason": errors.New(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "创建 Services 成功",
	})
}

// UpdateService 更新 Services 资源
func (s *services) UpdateService(c *gin.Context) {

	params := new(struct {
		Namespace string          `json:"namespace"`
		Data      *corev1.Service `json:"data"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("绑定参数失败" + err.Error()),
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Services.UpdateSvc(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "更新失败",
			"reason": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})

}
