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

type persistenvolume struct{}

var Persistentvolume persistenvolume

// GetPersistentVolumeList 获取PersistentVolume清单
func (p *persistenvolume) GetPersistentVolumeList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Persistenvolume.GetPvList(params.FilterName, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取PersistentVolume失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// GetPersistentVolumeDetail 获取 PersistentVolume 详情
func (p *persistenvolume) GetPersistentVolumeDetail(c *gin.Context) {
	params := new(struct {
		PersistentVolumeName string `form:"persistent_volume_name"`
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
	data, err := service.Persistenvolume.GetPvDetail(params.PersistentVolumeName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取PersistentVolume失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// DelPersistentVolume 删除 PersistentVolume 资源
func (p *persistenvolume) DelPersistentVolume(c *gin.Context) {
	params := new(struct {
		PersistentVolumeName string `form:"persistent_volume_name"`
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
	err = service.Persistenvolume.DelPv(params.PersistentVolumeName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "删除 PersistentVolume 失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// CreatePersistentVolume 创建 PersistentVolume 资源
func (p *persistenvolume) CreatePersistentVolume(c *gin.Context) {
	createPv := new(struct {
		Data *service.CreatePVConfig `json:"data"`
	})
	err := c.ShouldBindJSON(&createPv)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Persistenvolume.CreatePv(createPv.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "创建PersistentVolume失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "创建成功",
	})
}

// UpdatePersistentVolume 更新PersistentVolume 资源
func (p *persistenvolume) UpdatePersistentVolume(c *gin.Context) {
	params := new(struct {
		Data *corev1.PersistentVolume `json:"data"`
	})
	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": errors.New("绑定参数失败" + err.Error()),
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Persistenvolume.UpdatePv(params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "更新失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})

}
