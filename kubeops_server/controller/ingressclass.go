package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	networkv1 "k8s.io/api/networking/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

var IngressClass ingressClass

type ingressClass struct {
}

// GetIngressClassList list 获取列表信息
func (ic *ingressClass) GetIngressClassList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	err := c.ShouldBind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	//从header获取uuid，并转为int
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.IngressClass.GetIngressClass(params.FilterName, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": data,
	})

}

// GetIngressClassDetail 获取详情信息
func (ic *ingressClass) GetIngressClassDetail(c *gin.Context) {
	param := c.Query("Name")
	//从header获取uuid，并转为int
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.IngressClass.GetIngressClassDetail(param, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": data,
	})

}

// DelIngressClass 删除ingressClass
func (ic *ingressClass) DelIngressClass(c *gin.Context) {
	param := c.Query("Name")
	//从header获取uuid，并转为int
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.IngressClass.DelIngressClass(param, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// UpdateIngressClass 更新ingressClass
func (ic *ingressClass) UpdateIngressClass(c *gin.Context) {
	param := new(struct {
		Data *networkv1.IngressClass `json:"data"`
	})
	err := c.ShouldBind(&param)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.IngressClass.UpdateIngressClass(param.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "更新失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新成功",
		"data": nil,
	})

}