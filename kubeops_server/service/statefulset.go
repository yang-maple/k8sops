package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/model"
)

// 定义空结构体
type statefulSet struct{}

// StatefulSet  全局变量供外部调用
var StatefulSet statefulSet

// StsResp 定义返回的结构体
type StsResp struct {
	Total int            `json:"total"`
	Item  []statefulInfo `json:"item"`
}

type statefulInfo struct {
	Name       string            `json:"name"`
	Namespaces string            `json:"namespaces"`
	Image      []string          `json:"image"`
	Labels     map[string]string `json:"labels"`
	Pods       string            `json:"pods"`
	Age        string            `json:"age"`
	Status     string            `json:"status"`
}

// toCells 将 statefulCell 转换为 dataCell
func (s *statefulSet) toCells(statefulCell []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(statefulCell))
	for i := range cells {
		cells[i] = statefulsetCell(statefulCell[i])
	}
	return cells
}

// 将dataCell 转换为 statefulCell
func (s *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulCell := make([]appsv1.StatefulSet, len(cells))
	for i := range statefulCell {
		statefulCell[i] = appsv1.StatefulSet(cells[i].(statefulsetCell))
	}
	return statefulCell
}

// GetStatefulList  列表
func (s *statefulSet) GetStatefulList(StsName, Namespace string, Limit, Page int, uuid int) (stsResp *StsResp, err error) {
	statefulList, err := K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取 statefulList 失败" + err.Error())
		return nil, errors.New("获取 statefulList 失败" + err.Error())
	}
	selectData := &dataselector{
		GenericDataList: s.toCells(statefulList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{StsName},
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
	states := s.fromCells(dataPage.GenericDataList)
	item := make([]statefulInfo, 0, total)
	for _, v := range states {
		images := make([]string, 0, len(v.Spec.Template.Spec.Containers))
		for _, im := range v.Spec.Template.Spec.Containers {
			images = append(images, im.Image)
		}
		pods, status := model.GetStatus(v.Status.Replicas, v.Status.ReadyReplicas)
		item = append(item, statefulInfo{
			Name:       v.Name,
			Namespaces: v.Namespace,
			Image:      images,
			Labels:     v.Labels,
			Pods:       pods,
			Age:        model.GetAge(v.CreationTimestamp.Unix()),
			Status:     status,
		})
	}
	return &StsResp{
		Total: total,
		Item:  item,
	}, err

}

// GetStatefulDetail 详情
func (s *statefulSet) GetStatefulDetail(Namespace, StsName string, uuid int) (detail *appsv1.StatefulSet, err error) {
	//获取deploy
	detail, err = K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Get(context.TODO(), StsName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取 stateful set 详情失败" + err.Error())
		return nil, errors.New("获取 stateful set 详情失败" + err.Error())
	}
	detail.Kind = "StatefulSet"
	detail.APIVersion = "apps/v1"
	return detail, nil
}

// DelStateful 删除
func (s *statefulSet) DelStateful(Namespace, StsName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Delete(context.TODO(), StsName, metav1.DeleteOptions{})
	if err != nil {
		logger.Info("删除实例失败" + err.Error())
		return errors.New("删除实例失败" + err.Error())
	}
	return nil
}

// UpdateDelStateful 更新
func (s *statefulSet) UpdateDelStateful(Namespace string, Sts *appsv1.StatefulSet, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Update(context.TODO(), Sts, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("StatefulSets 更新失败" + err.Error())
		return errors.New("StatefulSets 更新失败" + err.Error())
	}

	return nil
}

// ModifyStatefulReplicas 修改deployment 副本数
func (s *statefulSet) ModifyStatefulReplicas(Namespace, StsName string, Replicas *int32, uuid int) (err error) {
	stateful, err := K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Get(context.TODO(), StsName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取 Stateful 数据失败" + err.Error())
		return errors.New("Stateful 详情失败" + err.Error())
	}
	stateful.Spec.Replicas = Replicas
	_, err = K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Update(context.TODO(), stateful, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("更新副本数失败" + err.Error())
		return errors.New("更新副本数失败" + err.Error())
	}
	return nil
}
