package service

import (
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type k8s struct {
	Clientset *kubernetes.Clientset
}

var K8s k8s

// 初始化方法
func (k *k8s) Init() {
	//读取配置文件
	conf, err := clientcmd.BuildConfigFromFlags("", "./config/config")
	if err != nil {
		panic("获取 配置文件失败" + err.Error())
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("创建 clientset 失败" + err.Error())
	} else {
		logger.Info("k8s 实例初始化成功")
	}
	k.Clientset = clientset

}
