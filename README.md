# k8sops
gin-vue3 realize k8s-dashboard

问题：

1. 终端连接还需要调试 前后端交互不够灵活
2. 退出时清除token 和 客户端连接
3. 跨域访问控制

实现

1. 初始化需要在config.go 中配置相关内容
