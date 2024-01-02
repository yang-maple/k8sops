
<template>
    <div class="home-header">
        <el-row>
            <el-col :span="24" style="text-align: left">
                <el-text>集群信息</el-text>
            </el-col>
        </el-row>
        <el-row>
            <el-col :span="6">
                <el-statistic title="集群总数" :value="1111" />
            </el-col>
            <el-col :span="6">
                <el-space direction="vertical">
                    <el-text
                        style="font-size: 12px;font-weight: 400; color:#606266;line-height: 20px;margin-bottom: -4px;">当前集群</el-text>
                    <el-text style="font-size: 20px; font-weight: 400; color:#303133;">{{ clusterName }}</el-text>
                </el-space>
            </el-col>
            <el-col :span="6">
                <el-statistic :value="222">
                    <template #title>
                        <div style="display: inline-flex; align-items: center">
                            私有云集群
                            <el-icon style="margin-left: 4px" :size="12">
                                <Male />
                            </el-icon>
                        </div>
                    </template>
                </el-statistic>
            </el-col>
            <el-col :span="6">
                <el-statistic title="阿里云集群" :value="333" />
            </el-col>
        </el-row>
    </div>
    <div class="home-nodes">
        <el-row>
            <el-col :span="24" style="text-align: left">
                <el-text>节点信息</el-text>
            </el-col>
        </el-row>
        <el-row>
            <el-col :span="16">
                <el-table :data="nodeItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
                    :default-sort="{ prop: 'date', order: 'descending' }">
                    <el-table-column label="名称" min-width="40" align="center">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <!-- <el-link style="margin-left: 10px" type="primary" :underline="false"
                                    @click="handle(scope.row)">{{
                                        scope.row.name
                                    }}</el-link> -->
                                <span style="margin-left: 10px">{{ scope.row.name }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column label="状态" align="center" width="100">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <span slot="reference" v-if="scope.row.status == 'Ready'">
                                    <el-tooltip placement="bottom" effect="light"><template #content>Ready</template>
                                        <i class="dotClass" style="background-color: springgreen"></i></el-tooltip>
                                </span>
                                <span slot="reference" v-if="scope.row.status != 'Ready'">
                                    <el-tooltip placement="bottom" effect="light"><template #content> NotReady </template>
                                        <i class="dotClass" style="background-color: red"></i></el-tooltip>
                                </span>
                                <span style="margin-left: 10px">{{ scope.row.status }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="cpu" align="center" label="CPU(核)" width="100" />
                    <el-table-column prop="memory" align="center" label="可分配内存" width="100">
                    </el-table-column>
                    <el-table-column prop="pods" align="center" label="pod 数量" width="70">
                    </el-table-column>
                    <el-table-column prop="create_time" align="center" label="创建时间" />

                </el-table>
            </el-col>
            <el-col :span="8">
                <div id="node-chart" class="node-chart"></div>
            </el-col>
        </el-row>
    </div>
    <div class="home-nodes">
        <el-row>
            <el-col :span="6">
            </el-col>
        </el-row>
    </div>
</template>

<script>
import { Male } from '@element-plus/icons-vue'
import * as echarts from 'echarts';
export default {
    data() {
        return {
            clusterName: "aaa",
            name: null,
            nodeItem: [],
            data: [],
        }
    },
    computed: {
    },
    created() {
        this.initdata()
    },
    mounted() {
        this.initEchart()
    },
    methods: {
        initEchart() {
            console.log(typeof this.data)
            var myChart = echarts.init(document.getElementById('node-chart'));
            myChart.setOption({
                title: {
                    subtext: 'Fake Data',
                    left: 'center'
                },
                tooltip: {
                    trigger: 'item'
                },
                series: [
                    {
                        name: 'Access From',
                        type: 'pie',
                        radius: '50%',
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
        },
        initdata() {
            this.$ajax({
                method: 'get',
                url: '/homepage/getInfo',
            }).then((res) => {
                this.nodeItem = res.data.node_info
                this.name = res.data.node_info[0].name
                for (var i = 0; i < res.data.node_info.length; i++) {
                    this.data.push({
                        value: res.data.node_info[i].pods,
                        name: res.data.node_info[i].name
                    })
                }
                console.log(this.data)
            }).catch((res) => {
                console.log(res);
            })


        }
    }
}
</script>

<style>
.home-header {
    padding-top: 5px;
    padding-left: 10px;
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
    height: 200px;
}
</style>