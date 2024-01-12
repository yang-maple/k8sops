# k8sops
gin-vue3 realize k8s-dashboard

1. bug 终端连接还需要调试
2. node 资源使用量

```go
cpulimit := make(map[string]float64)
	cpurequest := make(map[string]float64)
	memlimit := make(map[string]int64)
	memrequest := make(map[string]int64)
	start := time.Now().Unix()
	for _, pod := range pods.Items {
		if pod.Spec.NodeName != "" {
			for _, container := range pod.Spec.Containers {
				if _, ok := cpulimit[pod.Spec.NodeName]; ok {
					cpulimit[pod.Spec.NodeName] += float62(container.Resources.Limits.Cpu().AsApproximateFloat64())
					cpurequest[pod.Spec.NodeName] += float62(container.Resources.Requests.Cpu().AsApproximateFloat64())
					memlimit[pod.Spec.NodeName] += container.Resources.Limits.Memory().Value()
					memrequest[pod.Spec.NodeName] += container.Resources.Requests.Memory().Value()
					continue
				} else {
					cpulimit[pod.Spec.NodeName] += float62(container.Resources.Limits.Cpu().AsApproximateFloat64())
					cpurequest[pod.Spec.NodeName] += float62(container.Resources.Requests.Cpu().AsApproximateFloat64())
					memlimit[pod.Spec.NodeName] += container.Resources.Limits.Memory().Value()
					memrequest[pod.Spec.NodeName] += container.Resources.Requests.Memory().Value()
				}
			}
		}
	}

Math.trunc(this.node_resource.memory_used * 100) / 10
```

