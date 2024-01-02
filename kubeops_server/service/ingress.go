package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/model"
)

type ingress struct{}

var Ingress ingress

type IngResp struct {
	Total int           `json:"total"`
	Item  []ingressInfo `json:"item"`
}

type ingressInfo struct {
	Name      string                                 `json:"name"`
	Namespace string                                 `json:"namespace"`
	Labels    map[string]string                      `json:"labels"`
	Endpoint  []networkv1.IngressLoadBalancerIngress `json:"endpoint"`
	Host      []string                               `json:"host"`
	Age       string                                 `json:"age"`
}

type CreateIngress struct {
	Name             string            `json:"name"`
	Namespace        string            `json:"namespace"`
	Labels           map[string]string `json:"labels"`
	IngressClassName string            `json:"ingress_class_name"`
	Rules            []Rule            `json:"rules"`
}
type Rule struct {
	Host                  string     `json:"host"`
	HTTPIngressRuleValues []HTTPRule `json:"http_ingress_rule_values"`
}
type HTTPRule struct {
	Path        string `json:"path"`
	ServiceName string `json:"service_name"`
	ServicePort int32  `json:"service_port"`
}

func (i *ingress) toCells(ing []networkv1.Ingress) []DataCell {
	cells := make([]DataCell, len(ing))
	for i := range ing {
		cells[i] = ingressCell(ing[i])
	}
	return cells
}

func (i *ingress) fromCells(cells []DataCell) []networkv1.Ingress {
	svc := make([]networkv1.Ingress, len(cells))
	for i := range cells {
		svc[i] = networkv1.Ingress(cells[i].(ingressCell))
	}
	return svc
}

// GetIngList 获取 ingress 资源列表
func (i *ingress) GetIngList(ingName, Namespace string, Limit, Page int, uuid int) (DP *IngResp, err error) {
	//获取deployment 的所有清单列表
	ingList, err := K8s.Clientset[uuid].NetworkingV1().Ingresses(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取 ingress 失败" + err.Error())
	}

	//组装数据
	selectData := &dataselector{
		GenericDataList: i.toCells(ingList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{ingName},
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
	ing := i.fromCells(dataPage.GenericDataList)
	item := make([]ingressInfo, 0, total)
	for _, v := range ing {
		host := make([]string, 0)
		for _, value := range v.Spec.Rules {
			host = append(host, value.Host)
		}
		item = append(item, ingressInfo{
			Name:      v.Name,
			Namespace: v.Namespace,
			Labels:    v.Labels,
			Endpoint:  v.Status.LoadBalancer.Ingress,
			Host:      host,
			Age:       model.GetAge(v.CreationTimestamp.Unix()),
		})
	}
	return &IngResp{
		Total: total,
		Item:  item,
	}, err
}

// GetIngDetail 获取 ingress 资源详情
func (i *ingress) GetIngDetail(Namespace, ingName string, uuid int) (detail *networkv1.Ingress, err error) {
	//获取deploy
	detail, err = K8s.Clientset[uuid].NetworkingV1().Ingresses(Namespace).Get(context.TODO(), ingName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取ingress 详情失败" + err.Error())
		return nil, errors.New("获取ingress 详情失败" + err.Error())
	}
	detail.Kind = "Ingress"
	detail.APIVersion = "networking.k8s.io/v1"
	return detail, nil
}

// CreateIng 创建 ingress 资源
func (i *ingress) CreateIng(data *CreateIngress, uuid int) (err error) {
	pathType := networkv1.PathTypeImplementationSpecific
	createIng := &networkv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Spec: networkv1.IngressSpec{
			IngressClassName: &data.IngressClassName,
			TLS:              nil,
		},
		Status: networkv1.IngressStatus{},
	}
	for i := range data.Rules {
		specRole := networkv1.IngressRule{
			Host: data.Rules[i].Host,
		}
		for j := range data.Rules[i].HTTPIngressRuleValues {
			specRolesHttpPath := &networkv1.HTTPIngressRuleValue{Paths: []networkv1.HTTPIngressPath{
				{
					Path:     data.Rules[i].HTTPIngressRuleValues[j].Path,
					PathType: &pathType,
					Backend: networkv1.IngressBackend{
						Service: &networkv1.IngressServiceBackend{
							Name: data.Rules[i].HTTPIngressRuleValues[j].ServiceName,
							Port: networkv1.ServiceBackendPort{
								Name:   "",
								Number: data.Rules[i].HTTPIngressRuleValues[j].ServicePort,
							},
						},
					},
				},
			}}
			specRole.IngressRuleValue.HTTP = specRolesHttpPath
		}
		createIng.Spec.Rules = append(createIng.Spec.Rules, specRole)
	}
	_, err = K8s.Clientset[uuid].NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), createIng, metav1.CreateOptions{})
	if err != nil {
		logger.Info("创建 service 失败" + err.Error())
		return errors.New("创建 service 失败" + err.Error())
	}
	return nil

}

// DelIng 删除ingress 资源
func (i *ingress) DelIng(Namespace, ingName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].NetworkingV1().Ingresses(Namespace).Delete(context.TODO(), ingName, metav1.DeleteOptions{})
	if err != nil {
		logger.Info("删除ingress失败" + err.Error())
		return errors.New("删除ingress失败" + err.Error())
	}
	return nil
}

// UpdateIng 更新ingress 资源详情
func (i *ingress) UpdateIng(Namespace string, ing *networkv1.Ingress, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].NetworkingV1().Ingresses(Namespace).Update(context.TODO(), ing, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("ingress 更新失败" + err.Error())
		return errors.New("ingress 更新失败" + err.Error())
	}
	return nil
}
