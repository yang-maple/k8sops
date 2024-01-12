package service

//列表
// 获取node 详情
import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/model"
)

type persistenvolume struct{}

var Persistenvolume persistenvolume

type pvList struct {
	Total int      `json:"total"`
	Item  []pvInfo `json:"item"`
}

type pvInfo struct {
	Name          string                               `json:"name"`
	Labels        map[string]string                    `json:"labels"`
	Capacity      corev1.ResourceList                  `json:"capacity"`
	AccessMode    []corev1.PersistentVolumeAccessMode  `json:"access_mode"`
	ReclaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reclaim_policy"`
	Status        corev1.PersistentVolumeStatus        `json:"status"`
	Claim         string                               `json:"claim"`
	StorageClass  string                               `json:"storage_class"`
	Reason        string                               `json:"reason"`
	Age           string                               `json:"age"`
}

type pvDetail struct {
	Detail *corev1.PersistentVolume `json:"detail"`
	Age    string                   `json:"age"`
}

type CreatePVConfig struct {
	Name       string                              `json:"name"`
	Labels     map[string]string                   `json:"labels"`
	Storage    string                              `json:"storage"`
	AccessMode []corev1.PersistentVolumeAccessMode `json:"access_mode"`
	VolumeMode *corev1.PersistentVolumeMode        `json:"volume_mode"`
	Type       string                              `json:"type"`
	Path       string                              `json:"path"`
	Server     string                              `json:"server"`
}

func (p *persistenvolume) toCells(pvs []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(pvs))
	for i := range pvs {
		cells[i] = pvCell(pvs[i])
	}
	return cells
}

func (p *persistenvolume) fromCells(cells []DataCell) []corev1.PersistentVolume {
	pvs := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		pvs[i] = corev1.PersistentVolume(cells[i].(pvCell))
	}
	return pvs
}

// GetPvList PersistentVolume列表
func (p *persistenvolume) GetPvList(pvName string, Limit, Page int, uuid int) (*pvList, error) {
	pvs, err := K8s.Clientset[uuid].CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取PersistentVolume 失败" + err.Error())
		return nil, errors.New("获取PersistentVolume 失败" + err.Error())
	}

	//组装数据
	selectData := &dataselector{
		GenericDataList: p.toCells(pvs.Items),
		DataSelect: &DataSelectQuery{
			Paginate: &PaginateQuery{
				limit: Limit,
				page:  Page,
			},
			Filter: &FilterQuery{
				Name: pvName,
			},
		},
	}
	//筛选数据
	filtered := selectData.Filter()
	total := len(filtered.GenericDataList)
	//排序 并分页
	dataPage := filtered.Sort().Pagination()
	volumes := p.fromCells(dataPage.GenericDataList)
	item := make([]pvInfo, 0, len(pvs.Items))
	for _, pv := range volumes {
		item = append(item, pvInfo{
			Name:          pv.Name,
			Labels:        pv.Labels,
			Capacity:      pv.Spec.Capacity,
			AccessMode:    pv.Spec.AccessModes,
			ReclaimPolicy: pv.Spec.PersistentVolumeReclaimPolicy,
			Status:        pv.Status,
			Claim:         getClaim(pv.Spec.ClaimRef),
			StorageClass:  pv.Spec.StorageClassName,
			Age:           pv.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return &pvList{
		Total: total,
		Item:  item,
	}, nil
}

// GetPvDetail 获取PersistentVolume 详情
func (p *persistenvolume) GetPvDetail(PvName string, uuid int) (*pvDetail, error) {
	//获取deploy
	details, err := K8s.Clientset[uuid].CoreV1().PersistentVolumes().Get(context.TODO(), PvName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取PersistentVolume 详情失败" + err.Error())
		return nil, errors.New("获取PersistentVolume 详情失败" + err.Error())
	}
	details.Kind = "PersistentVolume"
	details.APIVersion = "v1"
	return &pvDetail{
		Detail: details,
		Age:    model.GetAge(details.CreationTimestamp.Unix()),
	}, nil
}

// DelPv 删除 PersistentVolume
func (p *persistenvolume) DelPv(PvName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].CoreV1().PersistentVolumes().Delete(context.TODO(), PvName, metav1.DeleteOptions{})
	if err != nil {
		logger.Info("删除 PersistentVolumes 详情失败" + err.Error())
		return errors.New("删除 PersistentVolumes 详情失败" + err.Error())
	}
	return nil
}

// CreatePv 创建 PersistentVolume
func (p *persistenvolume) CreatePv(data *CreatePVConfig, uuid int) (err error) {
	createPv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(data.Storage),
			},
			AccessModes:                   data.AccessMode,
			ClaimRef:                      nil,
			PersistentVolumeReclaimPolicy: "",
			StorageClassName:              "",
			MountOptions:                  nil,
			VolumeMode:                    data.VolumeMode,
			NodeAffinity:                  nil,
		},
	}
	switch data.Type {
	case "NFS":
		createPv.Spec.NFS = &corev1.NFSVolumeSource{
			Server: data.Server,
			Path:   data.Path,
		}
		createPv.Spec.StorageClassName = "nfs-client"
	case "HostPATH":
		createPv.Spec.HostPath = &corev1.HostPathVolumeSource{
			Path: data.Path,
			Type: nil,
		}
		createPv.Spec.StorageClassName = "standard"
	}
	_, err = K8s.Clientset[uuid].CoreV1().PersistentVolumes().Create(context.TODO(), createPv, metav1.CreateOptions{})
	if err != nil {
		logger.Info("创建 PersistentVolume 失败" + err.Error())
		return errors.New("创建 PersistentVolume 失败" + err.Error())
	}
	return nil
}

// UpdatePv 更新 PersistentVolume
func (p *persistenvolume) UpdatePv(deploy *corev1.PersistentVolume, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].CoreV1().PersistentVolumes().Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("PersistentVolume 更新失败" + err.Error())
		return errors.New("PersistentVolume 更新失败" + err.Error())
	}
	return nil
}

// 获取 claim 参数
func getClaim(claimRef *corev1.ObjectReference) string {
	if claimRef != nil {
		return claimRef.Namespace + "/" + claimRef.Name
	}
	return ""
}
