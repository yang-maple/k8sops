package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/model"
)

type namespace struct{}

var Namespace namespace

type NsList struct {
	Total int              `json:"total"`
	Item  []namespacesInfo `json:"item"`
}

type namespacesInfo struct {
	Name   string                `json:"name"`
	Labels map[string]string     `json:"labels"`
	Status corev1.NamespacePhase `json:"status"`
	Age    string                `json:"age"`
}

type NamespaceDetail struct {
	Detail *corev1.Namespace `json:"detail"`
	Age    string            `json:"age"`
}

// GetNsList  列表
func (n *namespace) GetNsList() (namespacesList *NsList, err error) {
	nsList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取namespaces-list 失败" + err.Error())
		return nil, errors.New("获取namespaces-list 失败" + err.Error())
	}
	list := make([]namespacesInfo, 0, len(nsList.Items))
	for _, v := range nsList.Items {
		list = append(list, namespacesInfo{
			Name:   v.Name,
			Labels: v.Labels,
			Status: v.Status.Phase,
			Age:    model.GetAge(v.CreationTimestamp.Unix()),
		})
	}
	return &NsList{
		Total: len(nsList.Items),
		Item:  list,
	}, err
}

// GetNsDetail 获取namespace 详情
func (n *namespace) GetNsDetail(NsName string) (*NamespaceDetail, error) {
	//获取deploy
	details, err := K8s.Clientset.CoreV1().Namespaces().Get(context.TODO(), NsName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取Namespace 详情失败" + err.Error())
		return nil, errors.New("获取Namespace 详情失败" + err.Error())
	}
	details.Kind = "Namespace"
	details.APIVersion = "v1"
	return &NamespaceDetail{
		Detail: details,
		Age:    model.GetAge(details.CreationTimestamp.Unix()),
	}, nil
}

// DelNs 删除
func (n *namespace) DelNs(NsName string) (err error) {
	err = K8s.Clientset.CoreV1().Namespaces().Delete(context.TODO(), NsName, metav1.DeleteOptions{})
	if err != nil {
		logger.Info("删除实例失败" + err.Error())
		return errors.New("删除实例失败" + err.Error())
	}
	return nil
}

// CreateNs 创建
func (n *namespace) CreateNs(NsName string) (err error) {
	namespaceConfig := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: NsName,
		},
	}
	_, err = K8s.Clientset.CoreV1().Namespaces().Create(context.TODO(), namespaceConfig, metav1.CreateOptions{})
	if err != nil {
		logger.Info("namespace 创建失败" + err.Error())
		return errors.New("namespace 创建失败" + err.Error())
	}

	return nil
}

// UpdateNs  更新
func (n *namespace) UpdateNs(deploy *corev1.Namespace) (err error) {
	_, err = K8s.Clientset.CoreV1().Namespaces().Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("NameSpaces 更新失败" + err.Error())
		return errors.New("NameSpaces 更新失败" + err.Error())
	}
	return nil
}
