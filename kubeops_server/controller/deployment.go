package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	"kubeops/model"
	"kubeops/service"
	"net/http"
)

var Deployment deployment

type deployment struct {
}

// GetDeploylist 获取deployment列表
func (d *deployment) GetDeploylist(c *gin.Context) {

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
	data, err := service.Deployment.GetDeploymentList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取podList失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})
}

// ModifyDeployReplicas 修改副本数
func (d *deployment) ModifyDeployReplicas(c *gin.Context) {
	params := new(struct {
		DeployName string `json:"deploy_name"`
		Namespace  string `json:"namespace"`
		Replicas   int    `json:"replicas"`
	})

	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		return
	}
	replicas := model.Int32Ptr(int32(params.Replicas))
	err = service.Deployment.ModifyDeployReplicas(params.Namespace, params.DeployName, replicas)
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

// RestartDeploy 重启 deployment
func (d *deployment) RestartDeploy(c *gin.Context) {
	params := new(struct {
		DeployName string `json:"deploy_name"`
		Namespace  string `json:"namespace"`
	})

	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		return
	}
	err = service.Deployment.RestartDeploy(params.Namespace, params.DeployName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": errors.New("重启失败" + err.Error()),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "重启成功",
		})
	}
}

// GetDeployDetail 获取 deploy 详情
func (d *deployment) GetDeployDetail(c *gin.Context) {
	params := new(struct {
		DeployName string `form:"deploy_name"`
		Namespace  string `form:"namespace"`
	})

	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("绑定参数失败" + err.Error()),
		})
		return
	}
	deploy, err := service.Deployment.GetDeployDetail(params.Namespace, params.DeployName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("获取deploy 失败" + err.Error()),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "获取成功",
			"data": deploy,
		})
	}
}

// CreateDeploy 创建 deploy 实例
func (d *deployment) CreateDeploy(c *gin.Context) {
	params := new(struct {
		Data *service.DeploymentCreate `json:"data"`
	})
	fmt.Println("接受参数")
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("绑定参数失败" + err.Error()),
		})
		return
	}
	err = service.Deployment.CreateDeploy(params.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "创建失败",
			"reason": errors.New(err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "创建成功",
	})

}

// DelDeploy 删除 deploy 实例
func (d *deployment) DelDeploy(c *gin.Context) {
	params := new(struct {
		DeployName string `form:"deploy_name"`
		Namespace  string `form:"namespace"`
	})
	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("绑定参数失败" + err.Error()),
		})
		return
	}
	err = service.Deployment.DelDeploy(params.Namespace, params.DeployName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    "删除失败",
			"reason": errors.New(err.Error()),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除成功",
		})
	}

}

// UpdateDeploy 更新 deploy 实例
func (d *deployment) UpdateDeploy(c *gin.Context) {

	params := new(struct {
		Namespace string             `json:"namespace"`
		Data      *appsv1.Deployment `json:"data"`
	})

	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("绑定参数失败" + err.Error()),
		})
		return
	}
	err = service.Deployment.UpdateDeploy(params.Namespace, params.Data)
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

// GetDeployPer 获取 每个namespace 下的deployment 实例
func (d *deployment) GetDeployPer(c *gin.Context) {

	dps, err := service.Deployment.GetDeployPer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errors.New("获取deploy 失败" + err.Error()),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "获取成功",
			"data": dps,
		})
	}
}

// RolloutDeploy 回滚 deploy 实例
func (d *deployment) RolloutDeploy(c *gin.Context) {
	params := new(struct {
		DeployName string `form:"deploy_name"`
		Namespace  string `form:"namespace"`
	})

	err := c.Bind(&params)
	if err != nil {
		logger.Info("绑定参数失败" + err.Error())
		return
	}
	err = service.Deployment.RolloutDeploy(params.Namespace, params.DeployName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": errors.New("回滚失败" + err.Error()),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "回滚成功",
		})
	}
}
