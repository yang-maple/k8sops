package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/model"
)

type secret struct{}
type SecretResp struct {
	Total int          `json:"total"`
	Item  []secretInfo `json:"item"`
}

type secretInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Immutable *bool             `json:"immutable"`
	Type      corev1.SecretType `json:"type"`
	Age       string            `json:"age"`
}

type CreateSecret struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Immutable bool              `json:"immutable"`
	Type      string            `json:"type"`
	Data      map[string]string `json:"data"`
}

var Secrets secret

func (s *secret) toCells(secrets []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(secrets))
	for i := range secrets {
		cells[i] = secretCell(secrets[i])
	}
	return cells
}

func (s *secret) fromCells(cells []DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		secrets[i] = corev1.Secret(cells[i].(secretCell))
	}
	return secrets
}

// GetSecretList 列表
func (s *secret) GetSecretList(secretName, Namespace string, Limit, Page int, uuid int) (DP *SecretResp, err error) {
	//获取deployment 的所有清单列表
	secretList, err := K8s.Clientset[uuid].CoreV1().Secrets(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取 secret 失败" + err.Error())
	}

	//组装数据
	selectData := &dataselector{
		GenericDataList: s.toCells(secretList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{secretName},
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
	secrets := s.fromCells(dataPage.GenericDataList)
	item := make([]secretInfo, 0, total)
	for _, v := range secrets {
		item = append(item, secretInfo{
			Name:      v.Name,
			Namespace: v.Namespace,
			Labels:    v.Labels,
			Immutable: v.Immutable,
			Type:      v.Type,
			Age:       model.GetAge(v.CreationTimestamp.Unix()),
		})
	}
	return &SecretResp{
		Total: total,
		Item:  item,
	}, err
}

// GetSecretDetail 详情
func (s *secret) GetSecretDetail(Namespace, secretName string, uuid int) (detail *corev1.Secret, err error) {
	//获取deploy
	detail, err = K8s.Clientset[uuid].CoreV1().Secrets(Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取Secret 详情失败" + err.Error())
		return nil, errors.New("获取Secret 详情失败" + err.Error())
	}
	detail.Kind = "Secret"
	detail.APIVersion = "v1"
	return detail, nil
}

// CreateSecret 创建
func (s *secret) CreateSecret(data *CreateSecret, uuid int) (err error) {

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Immutable:  &data.Immutable,
		Data:       nil,
		StringData: data.Data,
		Type:       corev1.SecretType(data.Type),
	}

	_, err = K8s.Clientset[uuid].CoreV1().Secrets(data.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		logger.Info("创建 secret 失败" + err.Error())
		return errors.New("创建 secret 失败" + err.Error())
	}
	return nil

}

// DelSecret 删除
func (s *secret) DelSecret(Namespace, secretName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].CoreV1().Secrets(Namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		logger.Info("删除Secret失败" + err.Error())
		return errors.New("删除Secret失败" + err.Error())
	}
	return nil
}

// UpdateSecret   更新
func (s *secret) UpdateSecret(Namespace string, Config *corev1.Secret, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].CoreV1().Secrets(Namespace).Update(context.TODO(), Config, metav1.UpdateOptions{})
	if err != nil {
		logger.Info("更新 secret 失败" + err.Error())
		return errors.New("更新 secret 失败" + err.Error())
	}
	return nil
}
