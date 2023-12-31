package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/wonderivan/logger"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/config"
	"kubeops/model"
)

var Pod pod

type pod struct {
}
type podNp struct {
	Name   string
	Number int
}

//类型转换的两个方法
// core1.pod --> Data cell

func (p *pod) toCells(pods []corev1.Pod) []DataCell {
	// 定义并初始化 []DataCell 数组
	cells := make([]DataCell, len(pods))
	//数据转换，将pods类型转化为data-cells 类型
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

// Data cell --> core1.pod
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {

	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//先断言判断是否为 pod cell 然后进行转换
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

type PodResp struct {
	Total int       `json:"total"`
	Item  []podInfo `json:"item"`
}

type podInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Image     []string          `json:"image"`
	Labels    map[string]string `json:"labels"`
	Node      string            `json:"node"`
	Status    string            `json:"status"`
	Restart   int32             `json:"restart"`
	Cpu       string            `json:"cpu"`
	Memory    string            `json:"memory"`
	Age       string            `json:"age"`
}

type CountPodReady struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
}

// GetPods 定义 get pod 方法获取pod 列表 支持过滤排序和分页
func (p *pod) GetPods(FilterName, NameSpaces string, Limit, Page int, uuid int) (PR *PodResp, err error) {
	//获取post-list
	postList, err := K8s.Clientset[uuid].CoreV1().Pods(NameSpaces).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取post-list 失败" + err.Error())
	}
	//实例化data select对象
	selectData := &dataselector{
		GenericDataList: p.toCells(postList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{FilterName},
			Paginate: &PaginateQuery{
				limit: Limit,
				page:  Page,
			},
		},
	}
	//先过滤 后排序
	filtered := selectData.Filter()
	total := len(filtered.GenericDataList)
	//分页
	dataPage := filtered.Sort().Pagination()
	pods := p.fromCells(dataPage.GenericDataList)
	items := make([]podInfo, 0, len(pods))
	for _, v := range pods {
		items = append(items, podInfo{
			Name:      v.Name,
			Namespace: v.Namespace,
			Labels:    v.Labels,
			Image:     model.GetImage(v.Spec),
			Node:      v.Spec.NodeName,
			Status:    GetPodStatus(v.Status),
			Restart:   model.GetRestart(v.Status),
			Cpu:       model.GetResourcesRequests("Cpu", v.Spec),
			Memory:    model.GetResourcesRequests("Memory", v.Spec),
			Age:       model.GetAge(v.CreationTimestamp.Unix()),
		})
	}

	return &PodResp{
		Total: total,
		Item:  items,
	}, err
}

// GetPodDetail 获取pod 详情
func (p *pod) GetPodDetail(Name, Namespace string, uuid int) (pod *corev1.Pod, err error) {
	podDetail, err := K8s.Clientset[uuid].CoreV1().Pods(Namespace).Get(context.TODO(), Name, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取pod 详情失败" + err.Error())
		return nil, errors.New("获取pod 详情失败")
	}
	podDetail.Kind = "Pod"
	podDetail.APIVersion = "v1"
	return podDetail, nil
}

// DelPod 删除 pod
func (p *pod) DelPod(Name, Namespace string, uuid int) error {
	err := K8s.Clientset[uuid].CoreV1().Pods(Namespace).Delete(context.TODO(), Name, metav1.DeleteOptions{})
	if err != nil {
		logger.Info("删除失败" + err.Error())
		return errors.New("删除失败" + err.Error())
	}
	return nil
}

// UpdatePod  更新 pod
func (p *pod) UpdatePod(pod *corev1.Pod, Namespaces string, uuid int) error {
	_, err := K8s.Clientset[uuid].CoreV1().Pods(Namespaces).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("更新失败" + err.Error())
		return errors.New("更新失败" + err.Error())
	}
	return nil
}

// GetContainer 获取pod 的容器列表
func (p *pod) GetContainer(podName, Namespaces string, uuid int) (containers []string, err error) {
	pod, err := p.GetPodDetail(podName, Namespaces, uuid)
	if err != nil {
		logger.Info("获取pod 失败" + err.Error())
		return nil, errors.New("获取pod 失败")
	}
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

// GetContainerLog 获取pod 内容器日志
func (p *pod) GetContainerLog(podName, containerName, namespaces string, uuid int) (log string, err error) {
	tailLines := int64(config.Logtaillimit)
	//设置日志的配置，容器名，获取的内容的配置
	logOpt := &corev1.PodLogOptions{
		Container: containerName,
		Follow:    false,
		TailLines: &tailLines,
	}

	//获取一个request 实例
	request := K8s.Clientset[uuid].CoreV1().Pods(namespaces).GetLogs(podName, logOpt)

	//发起stream连接，到底response.body
	podLog, err := request.Stream(context.TODO())
	if err != nil {
		logger.Info("获取pod-log 失败" + err.Error())
		return "", errors.New("获取pod-log 失败")
	}
	//延时关闭
	defer func(podLog io.ReadCloser) {
		err := podLog.Close()
		if err != nil {

		}
	}(podLog)
	//将response.body 写入缓冲区，用于转化成string 类型
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLog)
	if err != nil {
		logger.Info("写入缓冲失败 失败" + err.Error())
		return "", errors.New("写入缓冲失败 失败")
	}
	return buf.String(), nil
}

// CountPod  获取每个namespace 下pod的数量
func (p *pod) CountPod(uuid int) (total map[string]int, err error) {
	podCountPerNamespace := make(map[string]int)
	list, err := K8s.Clientset[uuid].CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取pod 失败" + err.Error())
		return nil, errors.New("获取pod 失败")
	}
	for _, v := range list.Items {
		podCountPerNamespace[v.Namespace]++
	}
	return podCountPerNamespace, nil
}

func (p *pod) PodCount(namespace string, uuid int) (*CountPodReady, error) {
	podList, err := K8s.Clientset[uuid].CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取pod 失败" + err.Error())
		return nil, errors.New("获取pod 失败")
	}
	count := 0

	for _, v := range podList.Items {
		if v.Status.Phase == "Running" {
			count++
		}
	}
	return &CountPodReady{
		Total:    len(podList.Items),
		Ready:    count,
		NotReady: len(podList.Items) - count,
	}, nil
}

// GetPodStatus 获取pod状态的方法
func GetPodStatus(status corev1.PodStatus) string {
	for _, v := range status.ContainerStatuses {
		if v.Ready {
			return "Running"
		}
		if v.State.Waiting != nil {
			return v.State.Waiting.Reason
		}
		if v.State.Terminated != nil {
			return v.State.Terminated.Reason
		}
		if v.LastTerminationState.Terminated != nil {
			return v.LastTerminationState.Terminated.Reason
		}
		if v.LastTerminationState.Waiting != nil {
			return v.LastTerminationState.Waiting.Reason
		}
	}
	//if status.Phase == "Pending" {
	//	return "Pending"
	//}
	return "UnKnow"
}
