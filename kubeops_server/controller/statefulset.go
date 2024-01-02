package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	"kubeops/model"
	"kubeops/service"
	"net/http"
	"strconv"
)

var StatefulSet statefulSet

type statefulSet struct {
}

// GetStatefulSetList 获取list
func (s *statefulSet) GetStatefulSetList(c *gin.Context) {

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
	data, err := service.StatefulSet.GetStatefulList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		logger.Info("获取 StatefulSet 失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取StatefulSet失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// GetStatefulDetail 获取Stateful 详情
func (s *statefulSet) GetStatefulDetail(c *gin.Context) {
	params := new(struct {
		StatefulSetName string `form:"stateful_set_name"`
		Namespace       string `form:"namespace"`
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
	Ds, err := service.StatefulSet.GetStatefulDetail(params.Namespace, params.StatefulSetName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("获取Stateful set 失败" + err.Error()),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "获取成功",
			"data": Ds,
		})
	}
}

// DelStatefulSet 删除 实例
func (s *statefulSet) DelStatefulSet(c *gin.Context) {
	params := new(struct {
		StatefulSetName string `form:"stateful_set_name"`
		Namespace       string `form:"namespace"`
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
	err = service.StatefulSet.DelStateful(params.Namespace, params.StatefulSetName, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("删除失败" + err.Error()),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除成功",
		})
	}

}

// UpDataStatefulSet 更新实例
func (s *statefulSet) UpDataStatefulSet(c *gin.Context) {
	params := new(struct {
		Namespace string              `json:"namespace"`
		Data      *appsv1.StatefulSet `json:"data"`
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
	err = service.StatefulSet.UpdateDelStateful(params.Namespace, params.Data, uuid)
	if err != nil {
		logger.Info("更新失败" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新失败",
			"err": errors.New(err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})

}

// ModifyStatefulReplicas  修改副本数
func (s *statefulSet) ModifyStatefulReplicas(c *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"stateful_set_name"`
		Namespace       string `json:"namespace"`
		Replicas        int    `json:"replicas"`
	})

	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		return
	}
	replicas := model.Int32Ptr(int32(params.Replicas))
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err = service.StatefulSet.ModifyStatefulReplicas(params.Namespace, params.StatefulSetName, replicas, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": errors.New("更新副本数失败" + err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新副本数成功",
	})

}
