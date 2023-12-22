package controller

import (
	"github.com/gin-gonic/gin"
)

// 定义 router 结构体
type router struct{}

// Router 实例化结构体
var Router router

// InitApiRouter 初始化路由规则
func (router *router) InitApiRouter(r *gin.Engine) {
	v1 := r.Group("/v1/api/pod")
	{
		v1.GET("/list", Pod.GetPods)
		v1.DELETE("/delete", Pod.DeletePod)
		v1.PUT("/update", Pod.UpdatePod)
		v1.GET("/container/list", Pod.GetContainerList)
		v1.GET("/detail", Pod.GetPodsDetail)
		v1.GET("/container/log", Pod.GetContainerLog)
		v1.GET("/ns", Pod.GetPodNumber)
	}
	v2 := r.Group("/v1/api/deploy")
	{
		v2.GET("/list", Deployment.GetDeploylist)
		v2.PUT("/modify", Deployment.ModifyDeployReplicas)
		v2.POST("/restart", Deployment.RestartDeploy)
		v2.GET("/detail", Deployment.GetDeployDetail)
		v2.POST("/create", Deployment.CreateDeploy)
		v2.DELETE("/delete", Deployment.DelDeploy)
		v2.PUT("/update", Deployment.UpdateDeploy)
		v2.GET("/ns", Deployment.GetDeployPer)
		v2.POST("/rollout", Deployment.RolloutDeploy)
	}
	v3 := r.Group("/v1/api/daemon")
	{
		v3.GET("/list", DaemonSet.GetDaemonList)
		v3.GET("/detail", DaemonSet.GetDaemonDetail)
		v3.DELETE("/delete", DaemonSet.DelDaemon)
		v3.PUT("/update", DaemonSet.UpdateDaemon)
	}
	v4 := r.Group("/v1/api/stateful")
	{
		v4.GET("/list", StatefulSet.GetStatefulSetList)
		v4.GET("/detail", StatefulSet.GetStatefulDetail)
		v4.DELETE("/delete", StatefulSet.DelStatefulSet)
		v4.PUT("/update", StatefulSet.UpDataStatefulSet)
		v4.PUT("/modify", StatefulSet.ModifyStatefulReplicas)
	}
	v5 := r.Group("/v1/api/namespaces")
	{
		v5.GET("/list", Namespace.GetNsList)
		v5.GET("/detail", Namespace.GetNsDetail)
		v5.DELETE("/delete", Namespace.DelNs)
		v5.POST("/create", Namespace.CreateNs)
		v5.PUT("/update", Namespace.UpdateNs)
	}
	v6 := r.Group("/v1/api/node")
	{
		v6.GET("/list", Node.GetNodeList)
		v6.GET("/detail", Node.GetNodeDeTal)
	}
	v7 := r.Group("/v1/api/pv")
	{
		v7.GET("/list", Persistentvolume.GetPersistentVolumeList)
		v7.GET("/detail", Persistentvolume.GetPersistentVolumeDetail)
		v7.DELETE("/delete", Persistentvolume.DelPersistentVolume)
		v7.POST("/create", Persistentvolume.CreatePersistentVolume)
		v7.PUT("/update", Persistentvolume.UpdatePersistentVolume)
	}
	v8 := r.Group("/v1/api/pvc")
	{
		v8.GET("/list", PersistentVolumeClaim.GetPersistentVolumeClaimList)
		v8.GET("/detail", PersistentVolumeClaim.GetPersistentVolumeClaimDetail)
		v8.DELETE("/delete", PersistentVolumeClaim.DelPersistentVolumeClaim)
		v8.POST("/create", PersistentVolumeClaim.CreatePersistentVolumeClaim)
		v8.PUT("/update", PersistentVolumeClaim.UpdatePersistentVolumeClaim)
	}
	v9 := r.Group("/v1/api/svc")
	{
		v9.GET("/list", Services.GetServiceList)
		v9.GET("/detail", Services.GetServiceDetail)
		v9.DELETE("/delete", Services.DelServices)
		v9.POST("/create", Services.CreateService)
		v9.PUT("/update", Services.UpdateService)
	}
	v10 := r.Group("/v1/api/ing")
	{
		v10.GET("/list", Ingress.GetIngressList)
		v10.GET("/detail", Ingress.GetIngressDetail)
		v10.DELETE("/delete", Ingress.DelIngress)
		v10.POST("/create", Ingress.CreateIngress)
		v10.PUT("/update", Ingress.UpdateIngress)
	}
	v11 := r.Group("/v1/api/cm")
	{
		v11.GET("/list", Configmaps.GetConfigmapList)
		v11.GET("/detail", Configmaps.GetConfigmapDetail)
		v11.DELETE("/delete", Configmaps.DelConfigmap)
		v11.POST("/create", Configmaps.CreateConfigmap)
		v11.PUT("/update", Configmaps.UpdateConfigmap)
	}
	v12 := r.Group("/v1/api/secret")
	{
		v12.GET("/list", Secrets.GetSecretList)
		v12.GET("/detail", Secrets.GetSecretDetail)
		v12.DELETE("/delete", Secrets.DelSecret)
		v12.POST("/create", Secrets.CreateSecret)
		v12.PUT("/update", Secrets.UpdateSecret)
	}
	v13 := r.Group("/v1/api/workflow")
	{
		v13.GET("/list", Workflow.GetWorkflowList)
		v13.GET("/detail", Workflow.GetWorkflowDetail)
		v13.DELETE("/delete", Workflow.DeleteWorkflow)
		v13.POST("/create", Workflow.CreateWorkflow)
	}
	v14 := r.Group("/v1/api/user")
	{
		v14.GET("/getCaptcha", Login.CaptchaImage)
		v14.POST("/login", Login.VerifyInfo)
		v14.POST("/register", Register.RegisterUser)
		v14.POST("/register/email", Register.SendEmail)
		v14.POST("/resetPassword/email", ResetPassword.verifyIdentity)
		v14.POST("/resetPassword", ResetPassword.reset)
	}
	v15 := r.Group("/v1/api/upload")
	{
		v15.POST("/uploadFile", Upload.uploadFile)
	}
}
