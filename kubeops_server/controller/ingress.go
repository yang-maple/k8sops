package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	networkv1 "k8s.io/api/networking/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type ingress struct{}

var Ingress ingress

// GetIngressList 获取 Ingress 列表
func (i *ingress) GetIngressList(c *gin.Context) {
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
	data, err := service.Ingress.GetIngList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取Ingress失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// GetIngressDetail  获取 Ingress 详情
func (i *ingress) GetIngressDetail(c *gin.Context) {
	params := new(struct {
		IngressName string `form:"ingress_name"`
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
	detail, err := service.Ingress.GetIngDetail(params.Namespace, params.IngressName, uuid)
	if err != nil {
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

// DelIngress   删除 Ingress 资源
func (i *ingress) DelIngress(c *gin.Context) {
	params := new(struct {
		IngressName string `form:"ingress_name"`
		Namespace   string `form:"namespace"`
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
	err = service.Ingress.DelIng(params.Namespace, params.IngressName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "删除 Services 失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除 Services 成功",
	})
}

// CreateIngress  创建 Ingress 资源
func (i *ingress) CreateIngress(c *gin.Context) {
	createIngress := new(struct {
		Data *service.CreateIngress `json:"data"`
	})
	err := c.ShouldBindJSON(&createIngress)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Ingress.CreateIng(createIngress.Data, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建 Services 失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "创建 Services 成功",
	})
}

// UpdateIngress  更新 Ingress 资源
func (i *ingress) UpdateIngress(c *gin.Context) {

	params := new(struct {
		Namespace string             `json:"namespace"`
		Data      *networkv1.Ingress `json:"data"`
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
	err = service.Ingress.UpdateIng(params.Namespace, params.Data, uuid)
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
