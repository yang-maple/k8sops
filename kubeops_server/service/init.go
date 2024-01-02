package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type k8s struct {
	Clientset map[int]*kubernetes.Clientset
}

var K8s = &k8s{
	Clientset: make(map[int]*kubernetes.Clientset),
}

// Init 初始化方法
func (k *k8s) Init(dir string, uuid int) error {
	// 是否初始化 - 是(关闭连接，重新连接) - 否(直接连接)
	//关闭连接存在问题
	//读取配置文件
	conf, err := clientcmd.BuildConfigFromFlags("", dir)
	if err != nil {
		logger.Info("获取 配置文件失败" + err.Error())
		return errors.New("获取 配置文件失败" + err.Error())
	}
	clientset, err := kubernetes.NewForConfig(conf)

	if err != nil {
		logger.Info("创建 client set 失败" + err.Error())
		return errors.New("创建 client set 失败" + err.Error())
	}
	logger.Info("k8s 初始化成功")
	k.Clientset[uuid] = clientset
	return nil
}
