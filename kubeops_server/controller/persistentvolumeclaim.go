package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type persistentVolumeClaim struct{}

var PersistentVolumeClaim persistentVolumeClaim

// GetPersistentVolumeClaimList 获取 PersistentVolumeClaim 的列表
func (p *persistentVolumeClaim) GetPersistentVolumeClaimList(c *gin.Context) {
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
	data, err := service.Claim.GetPVClaimList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取PersistentVolumeClaim失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// GetPersistentVolumeClaimDetail 获取 PersistentVolumeClaim 详情
func (p *persistentVolumeClaim) GetPersistentVolumeClaimDetail(c *gin.Context) {
	params := new(struct {
		PersistentVolumeClaimName string `form:"persistent_volume_claim_name"`
		Namespace                 string `form:"namespace"`
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
	detail, err := service.Claim.GetPVClaimDetail(params.Namespace, params.PersistentVolumeClaimName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取PersistentVolumeClaim 详情失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": detail,
	})
}

// CreatePersistentVolumeClaim 创建PersistentVolumeClaim 资源
func (p *persistentVolumeClaim) CreatePersistentVolumeClaim(c *gin.Context) {
	createPvc := new(struct {
		Data *service.CreateClaim `json:"data"`
	})
	err := c.ShouldBindJSON(&createPvc)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "绑定参数失败",
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Claim.CreatePVClaim(createPvc.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "创建PersistentVolumeClaim失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "创建PersistentVolumeClaim成功",
	})
}

// DelPersistentVolumeClaim 删除 PersistentVolumeClaim 资源
func (p *persistentVolumeClaim) DelPersistentVolumeClaim(c *gin.Context) {
	params := new(struct {
		PersistentVolumeClaimName string `form:"persistent_volume_claim_name"`
		Namespace                 string `form:"namespace"`
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
	err = service.Claim.DelPVClaim(params.Namespace, params.PersistentVolumeClaimName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "删除 PersistentVolumeClaim 失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除 PersistentVolumeClaim 成功",
	})
}

// UpdatePersistentVolumeClaim 更新 PersistentVolumeClaim 资源
func (p *persistentVolumeClaim) UpdatePersistentVolumeClaim(c *gin.Context) {
	params := new(struct {
		Data      *corev1.PersistentVolumeClaim `json:"data"`
		Namespace string                        `json:"namespace"`
	})
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "绑定参数失败",
		})
		return
	}
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.Claim.UpdatePVClaim(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "更新 PersistentVolumeClaim 失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新 PersistentVolumeClaim 成功",
	})
}
