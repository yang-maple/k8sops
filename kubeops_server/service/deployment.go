package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kubeops/model"
	"strconv"
	"time"
)

// 定义一个deployment 结构体
type deployment struct {
}

// Deployment 定义一个 deployment 类型的全局变量 Deployment
var Deployment deployment

// appsv1.deployment --> Datacell
func (d *deployment) tocells(deployments []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(deployments))
	//数据转换，将pods类型转化为datacells 类型
	for i := range deployments {
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}

// 转换成 appsv1.deploymentcell 类型
func (d *deployment) fromcells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

// DeployResp 定义获取清单的列表
type DeployResp struct {
	Total int          `json:"total"`
	Item  []deployinfo `json:"item"`
}

type deployinfo struct {
	Name       string            `json:"name"`
	Namespaces string            `json:"namespaces"`
	Image      []string          `json:"image"`
	Labels     map[string]string `json:"labels"`
	Pods       string            `json:"pods"`
	Age        string            `json:"age"`
	Status     string            `json:"status"`
}

// DeploymentCreate 定义DeploymentCreate 结构体用于创建 deployment需要的参数属性的定义
type DeploymentCreate struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Replicas    int32             `json:"replicas"`
	Labels      map[string]string `json:"labels"`
	Container   container         `json:"container"`
	HealthCheck bool              `json:"health_check"`
	HealthPath  string            `json:"health_path"`
}

type container struct {
	ContainerName string        `json:"container_name"`
	Image         string        `json:"image"`
	Cpu           string        `json:"cpu"`
	Memory        string        `json:"memory"`
	Containerport containerport `json:"container_port"`
}

type containerport struct {
	PortName      string          `json:"port_name"`
	ContainerPort int32           `json:"container_port"`
	Protocol      corev1.Protocol `json:"protocol"`
}

// Deploydp 定义返回的数据类型
type Deploydp struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type CountDeploy struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
}

// GetDeploymentList 获取deployment 列表
func (d *deployment) GetDeploymentList(DeployName, Namespace string, Limit, Page int, uuid int) (DP *DeployResp, err error) {
	//获取deployment 的所有清单列表
	deploylist, err := K8s.Clientset[uuid].AppsV1().Deployments(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取deploylist 失败" + err.Error())
	}
	//组装数据
	selectdata := &dataselector{
		GenericDataList: d.tocells(deploylist.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{DeployName},
			Paginate: &PaginateQuery{
				limit: Limit,
				page:  Page,
			},
		},
	}
	//先过滤 后排序
	filtered := selectdata.Filter()
	total := len(filtered.GenericDataList)
	//分页
	datapage := filtered.Sort().Pagination()
	deploys := d.fromcells(datapage.GenericDataList)
	item := make([]deployinfo, 0, len(deploys))
	for _, v := range deploys {
		images := make([]string, 0, len(v.Spec.Template.Spec.Containers))
		for _, im := range v.Spec.Template.Spec.Containers {
			images = append(images, im.Image)
		}
		pods, status := model.GetStatus(v.Status.Replicas, v.Status.ReadyReplicas)
		item = append(item, deployinfo{
			Name:       v.Name,
			Namespaces: v.Namespace,
			Image:      images,
			Labels:     v.Labels,
			Pods:       pods,
			Age:        model.GetAge(v.CreationTimestamp.Unix()),
			Status:     status,
		})
	}

	return &DeployResp{
		Total: total,
		Item:  item,
	}, err
}

// GetDeployDetail 获取deployment 详情
func (d *deployment) GetDeployDetail(Namespace, DeployName string, uuid int) (detail *appsv1.Deployment, err error) {
	//获取deploy
	detail, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Get(context.TODO(), DeployName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取deployment 详情失败" + err.Error())
		return nil, errors.New("获取deployment 详情失败" + err.Error())
	}
	detail.Kind = "Deployment"
	detail.APIVersion = "apps/v1"
	return detail, nil
}

// ModifyDeployReplicas 修改deployment 副本数
func (d *deployment) ModifyDeployReplicas(Namespace, DeployName string, Replicas *int32, uuid int) (err error) {
	deploy, err := K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Get(context.TODO(), DeployName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取deployment 数据失败" + err.Error())
		return errors.New("获取deployment 详情失败" + err.Error())
	}
	deploy.Spec.Replicas = Replicas
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("更新副本数失败" + err.Error())
		return errors.New("更新副本数失败" + err.Error())
	}
	return nil
}

// CreateDeploy 创建 deployment实例
func (d *deployment) CreateDeploy(data *DeploymentCreate, uuid int) (err error) {

	containers := make([]corev1.Container, 0, 1)
	containersPort := make([]corev1.ContainerPort, 0, 1)
	containersPort = append(containersPort, corev1.ContainerPort{
		Name:          data.Container.Containerport.PortName,
		ContainerPort: data.Container.Containerport.ContainerPort,
		Protocol:      data.Container.Containerport.Protocol,
	})
	containers = append(containers, corev1.Container{
		Name:  data.Container.ContainerName,
		Image: data.Container.Image,
		Ports: containersPort,
		Resources: corev1.ResourceRequirements{
			Limits: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse(data.Container.Cpu),
				corev1.ResourceMemory: resource.MustParse(data.Container.Memory),
			},
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse(data.Container.Cpu),
				corev1.ResourceMemory: resource.MustParse(data.Container.Memory),
			},
		},
	})
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: data.Labels,
				},
				Spec: corev1.PodSpec{
					Volumes:        nil,
					InitContainers: nil,
					Containers:     containers,
				},
			},
		},
	}
	// 判断是否打开健康检测
	if data.HealthCheck {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			//设置第一个容器的readinessprobe
			// 若存在多个容器 使用for 进行定义
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						// intstr.IntOrString 的作用是端口可以定义为整型，也可以定义为字符
						// type=0 表示该结构体实例内的数据为整型，转json 时只需要使用 IntVal 的数据
						// type=1 表示该结构体实例内的数据为字符串，转json 时只需要使用 StrVal 的数据
						Type:   0,
						IntVal: data.Container.Containerport.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(data.Namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		logger.Info("创建deployment 失败" + err.Error())
		return errors.New("创建deployment 失败" + err.Error())
	}
	return nil

}

// DelDeploy 删除 deployment 实例
func (d *deployment) DelDeploy(Namespace, DeployName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Delete(context.TODO(), DeployName, metav1.DeleteOptions{})
	if err != nil {
		logger.Info("删除实例失败" + err.Error())
		return errors.New("删除实例失败" + err.Error())
	}
	return nil
}

// RestartDeploy 重启 deployment 实例
func (d *deployment) RestartDeploy(Namespace, DeployName string, uuid int) (err error) {
	restartTime := fmt.Sprintf(`{"spec":{"template":{"metadata":{"labels":{"restart-time":"%s"}}}}}`, time.Now().Format("2006-01-02_15-04-05"))
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Patch(context.TODO(), DeployName, types.StrategicMergePatchType, []byte(restartTime), metav1.PatchOptions{})
	if err != nil {
		logger.Info("重启deploy 失败" + err.Error())
		return errors.New("重启deploy 失败" + err.Error())
	}
	return nil
}

// UpdateDeploy 更新 deployment 实例
func (d *deployment) UpdateDeploy(Namespace string, deploy *appsv1.Deployment, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("deployment 更新失败" + err.Error())
		return errors.New("deployment 更新失败" + err.Error())
	}

	return nil
}

// GetDeployPer 获取 每个namespace 下的deployment 实例
func (d *deployment) GetDeployPer(uuid int) (dps []Deploydp, err error) {

	namespacelist, err := K8s.Clientset[uuid].CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取namespace 列表失败" + err.Error())
		return nil, errors.New("获取namespace 列表失败")
	}

	for _, v := range namespacelist.Items {
		deploylist, err := K8s.Clientset[uuid].AppsV1().Deployments(v.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Info("获取deployment 列表失败" + err.Error())
			return nil, errors.New("获取namespace 列表失败")
		}
		np := &Deploydp{
			Name:  v.Name,
			Total: len(deploylist.Items),
		}
		dps = append(dps, *np)
	}
	return dps, nil
}

// RolloutDeploy 回滚 deployment
func (d *deployment) RolloutDeploy(Namespace, DeployName string, uuid int) (err error) {
	//获取deploy 信息
	deploy, err := K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Get(context.TODO(), DeployName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取实例失败" + err.Error())
		return errors.New("获取实例失败" + err.Error())
	}
	//获取当前版本号
	revison := deploy.Annotations["deployment.kubernetes.io/revision"]
	//回滚至上一个版本
	version, _ := strconv.Atoi(revison)
	patchData := fmt.Sprintf(`{"metadata":{"annotations":{"deployment.kubernetes.io/revision":"%d"}}}`, version)
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Patch(context.TODO(), DeployName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{})
	if err != nil {
		logger.Info("回滚deploy 失败" + err.Error())
		return errors.New("回滚deploy 失败" + err.Error())
	}
	return nil
}

// DeployCount 统计deployment 实例的状态分布
func (d *deployment) DeployCount(Namespace string, uuid int) (*CountDeploy, error) {
	deployList, err := K8s.Clientset[uuid].AppsV1().Deployments(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取deployment 列表失败" + err.Error())
		return nil, errors.New("获取namespace 列表失败")
	}
	count := 0
	for _, v := range deployList.Items {
		//readyReplicas 与 replicas 相等 就表示所有pod 都处于ready状态
		if v.Status.ReadyReplicas == *v.Spec.Replicas {
			count++
		}
	}
	return &CountDeploy{
		Total:    len(deployList.Items),
		Ready:    count,
		NotReady: len(deployList.Items) - count,
	}, nil
}

/*

type DeploymentCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	ContainerName string            `json:"container_name"`
	Image         string            `json:"image"`
	Labels        map[string]string `json:"labels"`
	Cpu           string            `json:"cpu" `
	Memory        string            `json:"memory" `
	PortName      string            `json:"port_name"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Protocol      corev1.Protocol   `json:"protocol"`
}
deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: data.Labels,
				},
				Spec: corev1.PodSpec{
					Volumes:        nil,
					InitContainers: nil,
					Containers: []corev1.Container{
						{
							Name:  data.ContainerName,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          data.PortName,
									ContainerPort: data.ContainerPort,
									Protocol:      data.Protocol,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceCPU:    resource.MustParse(data.Cpu),
									corev1.ResourceMemory: resource.MustParse(data.Memory),
								},
								Requests: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceCPU:    resource.MustParse(data.Cpu),
									corev1.ResourceMemory: resource.MustParse(data.Memory),
								},
							},
						},
					},
				},
			},
		},
	}
	// 判断是否打开健康检测
	if data.HealthCheck {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			//设置第一个容器的readinessprobe
			// 若存在多个容器 使用for 进行定义
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						// intstr.IntOrString 的作用是端口可以定义为整型，也可以定义为字符
						// type=0 表示该结构体实例内的数据为整型，转json 时只需要使用 IntVal 的数据
						// type=1 表示该结构体实例内的数据为字符串，转json 时只需要使用 StrVal 的数据
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}
*/
