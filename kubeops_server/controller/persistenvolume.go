package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
)

type persistenvolume struct{}

var Persistentvolume persistenvolume

// GetPersistentVolumeList 获取PersistentVolume清单
func (p *persistenvolume) GetPersistentVolumeList(c *gin.Context) {
	data, err := service.Persistenvolume.GetPvList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取PersistentVolume失败",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	data, err := service.Persistenvolume.GetPvDetail(params.PersistentVolumeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取PersistentVolume失败",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	err = service.Persistenvolume.DelPv(params.PersistentVolumeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "删除 PersistentVolume 失败",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "绑定参数失败",
			"data": nil,
		})
		return
	}
	err = service.Persistenvolume.CreatePv(createPv.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建PersistentVolume失败",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("绑定参数失败" + err.Error()),
		})
		return
	}
	err = service.Persistenvolume.UpdatePv(params.Data)
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
