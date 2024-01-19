package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

var Pod pod

type pod struct {
}

// GetPods 获取pod list
func (p *pod) GetPods(c *gin.Context) {
	//定义匿名结构体 处理入参的请求
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
			"msg": "绑定参数失败",
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取plist失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})

}

// GetPodsDetail 获取pod 详情
func (p *pod) GetPodsDetail(c *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		PodName   string `form:"pod_name"`
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
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取成功",
		"data": data,
	})
}

// DeletePod 删除 pod
func (p *pod) DeletePod(c *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		PodName   string `form:"pod_name"`
	})
	err := c.ShouldBind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	if params.Namespace == "" {
		params.Namespace = "default"
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Pod.DelPod(params.PodName, params.Namespace, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "删除失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除成功",
		})
	}
}

// UpdatePod 更新 pod
func (p *pod) UpdatePod(c *gin.Context) {
	params := new(struct {
		Namespace string     `json:"namespace"`
		Data      corev1.Pod `json:"data"`
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
	err = service.Pod.UpdatePod(&params.Data, params.Namespace, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "更新失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新成功",
		})
	}
}

// GetContainerList 获取容器列表
func (p *pod) GetContainerList(c *gin.Context) {
	params := new(struct {
		Namespace string `form:"namespace"`
		PodName   string `form:"pod_name"`
	})

	err := c.ShouldBind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	containers, err := service.Pod.GetContainer(params.PodName, params.Namespace, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取容器信息失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "获取容器信息成功",
			"item": containers,
		})
	}
}

// GetContainerLog 获取容器日志
func (p *pod) GetContainerLog(c *gin.Context) {

	params := new(struct {
		Namespace     string `form:"namespace"`
		PodName       string `form:"pod_name"`
		ContainerName string `form:"container_name"`
	})
	err := c.ShouldBind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	log, err := service.Pod.GetContainerLog(params.PodName, params.ContainerName, params.Namespace, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取日志失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "获取日志成功",
			"data": log,
		})
	}

}

// GetPodNumber  获取每个namespace下的pod 数量
func (p *pod) GetPodNumber(c *gin.Context) {
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	total, err := service.Pod.CountPod(uuid)
	fmt.Println("total:", total)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "获取数据成功",
			"data": total,
		})
	}
}
