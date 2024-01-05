package service

//列表
// 获取node 详情
import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/model"
	"strconv"
)

type node struct{}

var Node node

type NodeList struct {
	Total     int        `json:"total"`
	Item      []nodeInfo `json:"Item"`
	Resources Resource   `json:"resource"`
}

type Resource struct {
	FreeMemory  float64 `json:"free_memory"`
	TotalMemory float64 `json:"total_memory"`
	FreeCpu     float64 `json:"free_cpu"`
	TotalCpu    float64 `json:"total_cpu"`
}

type nodeInfo struct {
	Name           string            `json:"name"`
	Labels         map[string]string `json:"labels"`
	Status         string            `json:"status"`
	Cpu            string            `json:"cpu"`
	Memory         string            `json:"memory"`
	Pods           int               `json:"pods"`
	CreateTime     string            `json:"create_time"`
	KubeletVersion string            `json:"kubelet_version"`
}

type NodeDetail struct {
	Detail          *corev1.Node `json:"detail"`
	Pods            []poddetail  `json:"pods"`
	MemoryAllocator string       `json:"memory_allocator"`
	AGE             string       `json:"age"`
	Total           int          `json:"total"`
}

type poddetail struct {
	Name    string
	Image   []string
	Labels  map[string]string
	Status  corev1.PodPhase
	Restart int32
	PodAge  string
}

// GetNodeList  列表
func (n *node) GetNodeList(uuid int) (nodeList *NodeList, err error) {
	nodeLists, err := K8s.Clientset[uuid].CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取Node 失败" + err.Error())
		return nil, errors.New("获取Node 失败" + err.Error())
	}
	item := make([]nodeInfo, 0, len(nodeLists.Items))
	//计算memory 使用量
	var freeMemory, totalMemory, freeCpu, totalCpu float64
	for _, v := range nodeLists.Items {
		status := "Ready"
		if v.Status.Conditions[len(v.Status.Conditions)-1].Status != "True" {
			status = "NotReady"
		}
		memory, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", v.Status.Allocatable.Memory().AsApproximateFloat64()/float64(1024*1024*1024)), 64)
		item = append(item, nodeInfo{
			Name:           v.Name,
			Labels:         v.Labels,
			Status:         status,
			Cpu:            strconv.FormatFloat(v.Status.Allocatable.Cpu().AsApproximateFloat64()*1000, 'f', -1, 64) + "m",
			Memory:         strconv.FormatFloat(memory, 'f', -1, 64) + "Gi",
			Pods:           len(*n.GetNodePods(v.Name, uuid)),
			CreateTime:     v.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
			KubeletVersion: v.Status.NodeInfo.KubeletVersion,
		})
		freeMemory += v.Status.Allocatable.Memory().AsApproximateFloat64()
		totalMemory += v.Status.Capacity.Memory().AsApproximateFloat64()
		freeCpu += v.Status.Allocatable.Cpu().AsApproximateFloat64()
		totalCpu += v.Status.Capacity.Cpu().AsApproximateFloat64()
	}
	return &NodeList{
		Total: len(nodeLists.Items),
		Item:  item,
		Resources: Resource{
			FreeMemory:  freeMemory,
			TotalMemory: totalMemory,
			FreeCpu:     freeCpu,
			TotalCpu:    totalCpu,
		},
	}, err
}

// GetNodeDetail  获取namespace 详情
func (n *node) GetNodeDetail(NodeName string, uuid int) (details *NodeDetail, err error) {
	//获取deploy
	detail, err := K8s.Clientset[uuid].CoreV1().Nodes().Get(context.TODO(), NodeName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取Node 详情失败" + err.Error())
		return nil, errors.New("获取Node 详情失败" + err.Error())
	}
	pods := n.GetNodePods(NodeName, uuid)
	memory, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(detail.Status.Allocatable.Name(corev1.ResourceMemory, "").Value())/float64(1024*1024*1024)), 64)
	return &NodeDetail{
		Detail:          detail,
		Pods:            *pods,
		MemoryAllocator: strconv.FormatFloat(memory, 'f', -1, 64) + "Gi",
		AGE:             model.GetAge(detail.CreationTimestamp.Unix()),
		Total:           len(*pods),
	}, nil
}

// GetNodePods 获取node 详情
func (n *node) GetNodePods(NodeName string, uuid int) (detail *[]poddetail) {
	//获取pod list
	podlist, err := K8s.Clientset[uuid].CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取Node 详情失败" + err.Error())
	}
	nodepods := make([]poddetail, 0)
	for _, pods := range podlist.Items {
		if pods.Spec.NodeName == NodeName {
			containers := make([]string, 0)
			for _, container := range pods.Spec.Containers {
				containers = append(containers, container.Image)
			}
			if pods.Status.Phase != corev1.PodCompleted {
				nodepods = append(nodepods, poddetail{
					Name:    pods.Name,
					Image:   containers,
					Labels:  pods.Labels,
					Status:  pods.Status.Phase,
					Restart: pods.Status.ContainerStatuses[0].RestartCount,
					PodAge:  model.GetAge(pods.CreationTimestamp.Unix()),
				})
			}

		}
	}
	return &nodepods
}
