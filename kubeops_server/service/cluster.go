package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"kubeops/dao"
	"kubeops/model"
	"mime/multipart"
	"os"
	"path"
)

var Cluster cluster

type cluster struct{}
type clusterList struct {
	Item  []clusterInfo `json:"item"`
	Total int           `json:"total"`
}

type clusterInfo struct {
	ClusterName string `json:"cluster_name"`
	FileName    string `json:"file_name"`
	Type        string `json:"type"`
	Status      bool   `json:"status"`
	CreateTime  string `json:"create_time"`
}

// Create 创建集群
func (clt *cluster) Create(dir, clusterName, clusterType string, uuid int, file *multipart.FileHeader, c *gin.Context) error {
	//判断目录是不是存在,如果存在则不创建
	if !model.DirExists(dir) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			logger.Info("创建目录失败" + err.Error())
			return err
		}
	}
	//保留原始文件名
	oldFilename := file.Filename
	//重新命名文件
	file.Filename = clusterName + "_" + clusterType + ".conf"
	//保存文件
	dst := path.Join(dir, file.Filename)
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		logger.Info("保存文件失败" + err.Error())
		return err
	}
	//数据库存储
	err = dao.Cluster.AddCluster(clusterName, oldFilename, clusterType, dir, uuid)
	if err != nil {
		return err
	}
	return nil
}

// List 获取集群列表
func (clt *cluster) List(uuid int) (item *clusterList, err error) {
	//获取用户集群列表
	value, err := dao.Cluster.ListCluster(uuid)
	if err != nil {
		return nil, err
	}
	clusters := make([]clusterInfo, 0, len(*value))
	for _, v := range *value {
		cluster := clusterInfo{
			ClusterName: v.ClusterName,
			FileName:    v.FileName,
			Type:        v.Type,
			Status:      v.Status,
			CreateTime:  fmt.Sprintf("%v", v.CreatedAt),
		}
		clusters = append(clusters, cluster)
	}
	//返回集群列表
	return &clusterList{
		Item:  clusters,
		Total: len(clusters),
	}, nil
}

// Change 更换集群信息
func (clt *cluster) Change(clusterName string, uuid int) error {
	//找到dir - 初始化 - 更改
	// 查询dir
	dir, err := dao.Cluster.GetClusterDir(clusterName, uuid)
	if err != nil {
		return err
	}
	//更新k8s client set
	err = K8s.Init(dir, uuid)
	if err != nil {
		return err
	}
	//更改sql 中集群的状态
	err = dao.Cluster.ChangeCluster(clusterName, uuid)
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除集群
func (clt *cluster) Delete(name string, userid int) error {
	//获取 dir -- 删除本地 -- 删除库
	//获取集群路径
	dir, err := dao.Cluster.GetClusterDir(name, userid)
	if err != nil {
		return err
	}
	if dir != "" {
		//删除本地文件
		fmt.Println(dir)
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}
	err = dao.Cluster.DeleteCluster(name, userid)
	if err != nil {
		return err
	}
	return nil
}
