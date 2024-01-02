package service

var Homepage homepage

type homepage struct {
}

type homepageInfo struct {
	ClusterInfo []clusterInfo `json:"cluster_info"`
	NodeInfo    []nodeInfo    `json:"node_info"`
	DeployInfo  []deployinfo  `json:"deploy_info"`
	PodInfo     []podInfo     `json:"pod_info"`
}

//type homepageNode struct {
//}
//type homepageDeploy struct {
//}
//type homepagePod struct {
//}

// GetHomepage 获取首页信息
func (h *homepage) GetHomepage(uuid int) (info *homepageInfo, err error) {
	cInfo, err := getClusterInfo(uuid)
	if err != nil {
		return nil, err
	}
	nInfo, err := getNodeInfo(uuid)
	if err != nil {
		return nil, err
	}
	dInfo, err := getDeployInfo(uuid)
	if err != nil {
		return nil, err
	}
	pInfo, err := getPodInfo(uuid)
	if err != nil {
		return nil, err
	}
	return &homepageInfo{
		ClusterInfo: cInfo,
		NodeInfo:    nInfo,
		DeployInfo:  dInfo,
		PodInfo:     pInfo,
	}, nil
}

// 集群信息的处理方法
func getClusterInfo(uuid int) (info []clusterInfo, err error) {
	list, err := Cluster.List(uuid)
	if err != nil {
		return nil, err
	}
	return list.Item, nil
}

// node信息的处理方法
func getNodeInfo(uuid int) (info []nodeInfo, err error) {
	list, err := Node.GetNodeList(uuid)
	if err != nil {
		return nil, err
	}
	return list.Item, nil
}

// pod 信息的处理方法
func getPodInfo(uuid int) (info []podInfo, err error) {
	list, err := Pod.GetPods("", "", 0, 0, uuid)
	if err != nil {
		return nil, err
	}
	return list.Item, nil
}

// deployment 信息的处理方法
func getDeployInfo(uuid int) (info []deployinfo, err error) {
	list, err := Deployment.GetDeploymentList("", "", 0, 0, uuid)
	if err != nil {
		return nil, err
	}
	return list.Item, nil
}
