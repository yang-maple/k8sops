
<template>
    <div class="home-header">
        <el-row style="padding-bottom: 5px">
            <el-col :span="24" style="text-align: left">
                <el-text style="font-size: 20px; font-weight: 400; color:#303133;">集群信息</el-text>
            </el-col>
        </el-row>
        <el-row :gutter="10">
            <el-col :span="6" style="padding-left: 10px;">

                <el-statistic :value="clusternum">
                    <template #title>
                        <div style="display: inline-flex; align-items: center">
                            集群总数
                        </div>
                    </template>
                </el-statistic>
            </el-col>
            <el-col :span="6">
                <el-statistic title="当前集群" :value="null">
                    <template #prefix>
                        {{ clusterName }}
                    </template>
                </el-statistic>
            </el-col>
            <el-col :span="6">

                <el-statistic :value="k8snum">
                    <template #title>
                        <div style="display: inline-flex; align-items: center">
                            私有云集群
                        </div>
                    </template>
                </el-statistic>

            </el-col>
            <el-col :span="6" style="padding-right: 10px;">

                <el-statistic :value="clusternum - k8snum">
                    <template #title>
                        <div style="display: inline-flex; align-items: center">
                            公有云集群
                        </div>
                    </template>
                </el-statistic>

            </el-col>
        </el-row>
    </div>
    <div class="home-nodes">
        <el-row>
            <el-col :span="24" style="text-align: left">
                <el-text style="font-size: 20px; font-weight: 400; color:#303133;">节点信息</el-text>
            </el-col>
        </el-row>
        <el-row style="padding-top: 5px;">
            <el-col :span="24" style="text-align: center;">
                <el-table :data="nodeItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
                    :default-sort="{ prop: 'date', order: 'descending' }">
                    <el-table-column label="名称" prop="name" min-width="40" align="center">
                    </el-table-column>
                    <el-table-column label="状态" align="center" width="120">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <span slot="reference"
                                    v-if="scope.row.status == 'Ready' && scope.row.unschedulable == false">
                                    <el-tooltip placement="top" effect="light"><template #content>Ready</template>
                                        <i class="dotClass" style="background-color: springgreen"></i></el-tooltip>
                                </span>
                                <span slot="reference"
                                    v-if="scope.row.status == 'Ready' && scope.row.unschedulable == true">
                                    <el-tooltip placement="top" effect="light"><template #content>Unschedulable</template>
                                        <i class="dotClass" style="background-color: orange"></i></el-tooltip>
                                </span>
                                <span slot="reference" v-if="scope.row.status != 'Ready'">
                                    <el-tooltip placement="top" effect="light"><template #content> NotReady </template>
                                        <i class="dotClass" style="background-color: red"></i></el-tooltip>
                                </span>
                                <span>
                                    <el-text v-if="scope.row.status == 'Ready' && scope.row.unschedulable == false"
                                        style="margin-left: 5px;color: green;">
                                        运行中
                                    </el-text>
                                    <el-text v-if="scope.row.status == 'Ready' && scope.row.unschedulable == true"
                                        style="margin-left: 5px;color: orangered;">
                                        不可调度
                                    </el-text>
                                    <el-text v-if="scope.row.status != 'Ready'" style="margin-left: 5px;color: red;">
                                        异常
                                    </el-text>
                                </span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="cpu" align="center" label="已使用CPU" width="120">
                        <template #default="scope">
                            <div>
                                <span>{{ node_resource["cpu_request"][scope.row.name] }}m</span>
                                <span style="margin-left: 3px;">({{ toNumber(node_resource["cpu_request"][scope.row.name] *
                                    100 / scope.row.cpu_total)
                                }}%)</span>
                            </div>
                        </template>
                    </el-table-column>

                    <el-table-column prop="memory" align="center" label="已使用内存" width="120">
                        <template #default="scope">
                            <div>
                                <span>{{ node_resource["memory_request"][scope.row.name] }}Mi</span>
                                <span style="margin-left: 3px;">({{ toNumber(node_resource["memory_request"][scope.row.name]
                                    *
                                    100 / scope.row.memory_total)
                                }}%)</span>
                            </div>
                        </template>
                    </el-table-column>

                    <el-table-column prop="pods" align="center" label="pod 数量" width="70">
                    </el-table-column>
                    <el-table-column prop="create_time" align="center" label="创建时间" />
                </el-table>
            </el-col>
        </el-row>
        <el-row>
            <el-col :span="8">
                <div id="node-chart" class="node-chart"></div>
            </el-col>
            <el-col :span="8">
                <div id="memory-chart" class="memory-chart"></div>
            </el-col>
            <el-col :span="8">
                <div id="cpu-chart" class="cpu-chart"></div>
            </el-col>
        </el-row>
    </div>
    <div class="home-deploy">
        <el-row>
            <el-col :span="24" style="text-align: left">
                <el-text style="font-size: 20px; font-weight: 400; color:#303133;">工作负载</el-text>
            </el-col>
        </el-row>
        <el-row :gutter="10" style="margin-top: 5px;">
            <el-col :span="6" style="padding-left: 10px;">
                <div class="statistic-card">
                    <el-statistic :value="deployInfo.total">
                        <template #title>
                            <div style="display: inline-flex; align-items: center">
                                无状态服务
                                <el-tooltip effect="dark" content="Deployment" placement="top">
                                    <el-icon style="margin-left: 4px" :size="12">
                                        <Warning />
                                    </el-icon>
                                </el-tooltip>
                            </div>
                        </template>
                    </el-statistic>
                    <div class="statistic-footer">
                        <div class="footer-item">
                            <span slot="reference">
                                <i class="cardClass" style="background-color: springgreen"></i>
                            </span>
                            <span style="color: green;">正常 {{ deployInfo.ready }}</span>
                            <span slot="reference">
                                <i class="cardClass" style="background-color: #FF4500"></i>
                            </span>
                            <span style="color: red;">异常 {{ deployInfo.not_ready }}</span>
                        </div>
                    </div>
                </div>
            </el-col>
            <el-col :span="6">
                <div class="statistic-card">
                    <el-statistic :value="statefulInfo.total">
                        <template #title>
                            <div style="display: inline-flex; align-items: center">
                                有状态服务
                                <el-tooltip effect="dark" content="StatefulSet" placement="top">
                                    <el-icon style="margin-left: 4px" :size="12">
                                        <Warning />
                                    </el-icon>
                                </el-tooltip>
                            </div>
                        </template>
                    </el-statistic>
                    <div class="statistic-footer">
                        <div class="footer-item">
                            <span slot="reference">
                                <i class="cardClass" style="background-color: springgreen"></i>
                            </span>
                            <span style="color: green;">正常 {{ statefulInfo.ready }}</span>
                            <span slot="reference">
                                <i class="cardClass" style="background-color: #FF4500"></i>
                            </span>
                            <span style="color: red;">异常 {{ statefulInfo.not_ready }}</span>
                        </div>
                    </div>
                </div>
            </el-col>
            <el-col :span="6">
                <div class="statistic-card">
                    <el-statistic :value="daemonsetInfo.total">
                        <template #title>
                            <div style="display: inline-flex; align-items: center">
                                守护进程
                                <el-tooltip effect="dark" content="DaemonSet" placement="top">
                                    <el-icon style="margin-left: 4px" :size="12">
                                        <Warning />
                                    </el-icon>
                                </el-tooltip>
                            </div>
                        </template>
                    </el-statistic>
                    <div class="statistic-footer">
                        <div class="footer-item">
                            <span slot="reference">
                                <i class="cardClass" style="background-color: springgreen"></i>
                            </span>
                            <span style="color: green;">正常 {{ daemonsetInfo.ready }}</span>
                            <span slot="reference">
                                <i class="cardClass" style="background-color: #FF4500"></i>
                            </span>
                            <span style="color: red;">异常 {{ daemonsetInfo.not_ready }}</span>
                        </div>
                    </div>
                </div>
            </el-col>
            <el-col :span="6" style="padding-right: 10px;">
                <div class="statistic-card">
                    <el-statistic :value="podInfo.total">
                        <template #title>
                            <div style="display: inline-flex; align-items: center">
                                容器组
                                <el-tooltip effect="dark" content="Pod" placement="top">
                                    <el-icon style="margin-left: 4px" :size="12">
                                        <Warning />
                                    </el-icon>
                                </el-tooltip>
                            </div>
                        </template>
                    </el-statistic>
                    <div class="statistic-footer">
                        <div class="footer-item">
                            <span slot="reference">
                                <i class="cardClass" style="background-color: springgreen"></i>
                            </span>
                            <span style="color: green;">正常 {{ podInfo.ready }}</span>
                            <span slot="reference">
                                <i class="cardClass" style="background-color: #FF4500"></i>
                            </span>
                            <span style="color: red;">异常 {{ podInfo.not_ready }}</span>
                        </div>

                    </div>
                </div>
            </el-col>
        </el-row>
        <el-row style="margin-top: 20px;">
            <el-col :span="24">
                <div id="namespace-chart" class="namespace-chart"></div>
            </el-col>
        </el-row>
    </div>
    <div class="home-event">
        <el-row>
            <el-col :span="24" style="text-align: left">
                <el-text style="font-size: 20px; font-weight: 400; color:#303133;">集群事件</el-text>
            </el-col>
        </el-row>
        <el-row>
            <el-table :data="eventItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
                :default-sort="{ prop: 'date', order: 'descending' }" height="320">
                <el-table-column label="命名空间" prop="namespace" width="120" align="center"></el-table-column>
                <el-table-column label="类型" prop="type" width="120" align="center">
                    <template #default="scope">
                        <div style="display: flex; align-items: center">
                            <span slot="reference" v-if="scope.row.type == 'Warning'">
                                <i class="dotClass" style="background-color:#FFA500"></i>
                            </span>
                            <span slot="reference" v-if="scope.row.type.includes('Failed')">
                                <i class="dotClass" style="background-color:#FF4500"></i>
                            </span>
                            <span slot="reference" v-if="scope.row.type == 'Error'">
                                <i class="dotClass" style="background-color:red"></i>
                            </span>
                            <span style="margin-left: 3px">{{ scope.row.type }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="对象" prop="object" min-width="40" align="center"></el-table-column>
                <el-table-column label="原因" prop="reason" width="120" align="center"></el-table-column>
                <el-table-column label="事件" prop="message" min-width="40" align="center"></el-table-column>
                <el-table-column label="时间" prop="last_time" width="150" align="center"></el-table-column>

            </el-table>
        </el-row>
    </div>
</template>

<script>
import * as echarts from 'echarts';
export default {
    data() {
        return {
            clusterName: null,
            name: null,
            nodeItem: [],
            namespace: [],
            pod_total: [],
            deploy_total: [],
            state_total: [],
            job_total: [],
            daemon_total: [],
            eventItem: [],
            clusternum: null,
            deployInfo: {},
            statefulInfo: {},
            daemonsetInfo: {},
            podInfo: {},
            jobnum: null,
            k8snum: null,
            data: [],
            node_resource: {},
            node_resource_used: {
                cpu_used: 0,
                memory_used: 0,
                cpu_total: 0,
                memory_total: 0,
            },
            // num: 12,
        }
    },
    computed: {
    },
    created() {
        this.initdata()
    },
    mounted() {

    },
    methods: {
        initdata() {
            this.$ajax({
                method: 'get',
                url: '/homepage/getInfo',
            }).then((res) => {
                //集群资源使用率
                this.node_resource = res.data.node_resource
                const resourceMap = new Map(Object.entries(res.data.node_resource.cpu_request))
                resourceMap.forEach((value, key) => {
                    this.node_resource_used.cpu_used += value
                })
                const resourceMap2 = new Map(Object.entries(res.data.node_resource.memory_request))
                resourceMap2.forEach((value, key) => {
                    this.node_resource_used.memory_used += value
                })
                //集群节点信息
                this.nodeItem = res.data.node_info
                //集群信息
                this.clusternum = res.data.cluster_info.length
                this.clusterName = localStorage.getItem('user_cluster')
                //工作负载信息
                this.deployInfo = res.data.deploy_info
                //状态fulset信息
                this.statefulInfo = res.data.stateful_info
                //daemonset信息
                this.daemonsetInfo = res.data.daemon_info
                //pod信息
                this.podInfo = res.data.pod_info
                //echart pod 数量占比 数据统计
                for (var i = 0; i < res.data.node_info.length; i++) {
                    this.node_resource_used.cpu_total += res.data.node_info[i].cpu_total
                    this.node_resource_used.memory_total += res.data.node_info[i].memory_total
                    this.data.push({
                        value: res.data.node_info[i].pods,
                        name: res.data.node_info[i].name
                    })
                }
                //私有云数据统计
                for (var i = 0; i < res.data.cluster_info.length; i++) {
                    this.k8snum = 0
                    if (res.data.cluster_info[i].type == 'k8s') {
                        this.k8snum++
                    }
                }
                //echart namespace资源占比统计
                //统计namespace 名称以及数量
                // 统计pod deploy job daemonset数量
                const map = new Map(Object.entries(res.data.pod_total))
                map.forEach((value, key) => {
                    this.namespace.push(key)
                    this.pod_total.push(value.pod_total)
                    this.deploy_total.push(value.deploy_total)
                    this.state_total.push(value.stateful_total)
                    this.job_total.push(value.job_total)
                    this.daemon_total.push(value.daemon_total)
                })
                //统计event
                this.eventItem = res.data.event_total.item
                this.initEchart()
            }).catch((res) => {
                console.log(res);
            })

        },
        initEchart() {
            var NodeChart = echarts.init(document.getElementById('node-chart'));
            NodeChart.setOption({
                title: {
                    //text: '节点 Pod 数量',
                    subtext: '节点 Pod 数量',
                    left: 'center'
                },
                tooltip: {
                    trigger: 'item'
                },
                series: [
                    {
                        name: 'Pod 数量占比',
                        type: 'pie',
                        radius: '70%',
                        //给data 动态赋值
                        data: this.data,
                        emphasis: {
                            itemStyle: {
                                shadowBlur: 10,
                                shadowOffsetX: 0,
                                shadowColor: 'rgba(0, 0, 0, 0.5)'
                            }
                        }
                    }
                ]
            })
            var MemoryChart = echarts.init(document.getElementById('memory-chart'));
            MemoryChart.setOption({
                title: {
                    //text: '节点 Pod 数量',
                    subtext: '内存使用情况',
                    left: 'center'
                },
                tooltip: {
                    trigger: 'item'
                },
                series: [
                    {

                        type: 'pie',
                        radius: '70%',
                        label: {
                            show: true,
                            formatter(param) {
                                // correct the percentage
                                return param.name + ' (' + param.percent + '%)';
                            }
                        },
                        //给data 动态赋值
                        data: [{
                            value: this.node_resource_used.memory_used,
                            name: '内存总使用量(Gi)'
                        }, {
                            value: this.node_resource_used.memory_total - this.node_resource_used.memory_used,
                            name: '内存总剩余量(Gi)'
                        },
                        ],
                        emphasis: {
                            itemStyle: {
                                shadowBlur: 10,
                                shadowOffsetX: 0,
                                shadowColor: 'rgba(0, 0, 0, 0.5)'
                            },
                        }
                    }
                ]
            })
            var CpuChart = echarts.init(document.getElementById('cpu-chart'));
            CpuChart.setOption({
                title: {
                    //text: '节点 Pod 数量',
                    subtext: 'CPU使用情况',
                    left: 'center'
                },
                tooltip: {
                    trigger: 'item'
                },
                series: [
                    {
                        type: 'pie',
                        radius: '70%',
                        label: {
                            show: true,
                            formatter(param) {
                                // correct the percentage
                                return param.name + ' (' + param.percent + '%)';
                            }
                        },
                        //给data 动态赋值
                        data: [{
                            value: this.node_resource_used.cpu_used,
                            name: 'cpu总使用量(m)'
                        }, {
                            value: this.node_resource_used.cpu_total - this.node_resource_used.cpu_used,
                            name: 'cpu总剩余量(m)'
                        }],
                        emphasis: {
                            itemStyle: {
                                shadowBlur: 10,
                                shadowOffsetX: 0,
                                shadowColor: 'rgba(0, 0, 0, 0.5)'
                            }
                        }
                    }
                ]
            })
            var NamespaceChart = echarts.init(document.getElementById('namespace-chart'));
            NamespaceChart.setOption({
                title: {
                    subtext: '工作负载资源分布',
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'shadow'
                    }
                },
                legend: {},
                grid: {
                    left: '3%',
                    right: '4%',
                    bottom: '3%',
                    containLabel: true
                },
                xAxis: {
                    type: 'value',
                    boundaryGap: [0, 0.01]
                },
                yAxis: {
                    type: 'category',
                    data: this.namespace
                },
                series: [
                    {
                        name: 'Pod',
                        type: 'bar',
                        data: this.pod_total
                    },
                    {
                        name: 'Deployment',
                        type: 'bar',
                        data: this.deploy_total
                    },
                    {
                        name: 'StatefulSet',
                        type: 'bar',
                        data: this.state_total
                    },
                    {
                        name: 'DaemonSet',
                        type: 'bar',
                        data: this.daemon_total
                    },
                    {
                        name: 'Job',
                        type: 'bar',
                        data: this.job_total
                    }
                ]
            })
        },
        toNumber(val) {
            return Math.floor(val)
        },
    }
}
</script>

<style>
.home-header {
    padding-top: 5px;
    padding-right: 10px;
    padding-bottom: 5px;
    border-radius: 4px;
    background: #f0f2f5;
    text-align: center;
    margin-bottom: 20px;

}

.home-nodes {

    padding-top: 5px;
    background: #f0f2f5;
    margin-bottom: 20px;
}

.home-deploy {
    padding-top: 5px;
    background: #f0f2f5;
    text-align: center;
    margin-bottom: 20px;
    min-height: 155px;

    .el-statistic {
        --el-statistic-content-font-size: 40px;
    }
}

.home-content {
    padding-top: 5px;
    background: #f0f2f5;
    text-align: center;
}

.el-table,
.el-table__expanded-cell {
    background-color: transparent;
}

.el-table th,
.el-table tr {
    background-color: transparent;
}

.dotClass {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    display: block;
    margin-left: 10px;
}

.cardClass {
    width: 25px;
    height: 18px;
    border-radius: 0%;
    display: block;
    margin-left: 30px;
    margin-right: 10px;
    margin-top: 3px;
}

.percentage-value {
    display: block;
    margin-top: 10px;
    font-size: 28px;
}

.percentage-label {
    display: block;
    margin-top: 10px;
    font-size: 12px;
}

.node-chart {
    width: 100%;
    height: 320px;
}

.memory-chart {
    width: 100%;
    height: 320px;
}

.cpu-chart {
    width: 100%;
    height: 320px;
}

.namespace-chart {
    width: 100%;
    min-height: 450px;
}

.home-event {
    padding-top: 5px;
    background: #f0f2f5;
    text-align: center;
    margin-bottom: 20px;
    height: 360px;
}


.el-statistic {
    --el-statistic-content-font-size: 28px;
}

.statistic-card {
    background-color: #fff;
    text-align: center;
    padding: 20px;
}

.statistic-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    font-size: 18px;
    color: var(--el-text-color-regular);
    margin-top: 16px;
}

.statistic-footer .footer-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
</style>