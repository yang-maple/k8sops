<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content bg-purple-light">
        <el-collapse v-model="activeNames">
          <el-collapse-item title="Metadata" name="1">
            <div>
              <el-descriptions direction="vertical">
                <el-descriptions-item label="名 称">{{ metadata.name }}</el-descriptions-item>
                <el-descriptions-item label="创建时间">{{ metadata.creationTimestamp }}</el-descriptions-item>
                <el-descriptions-item label="Age">{{ age }}</el-descriptions-item>
              </el-descriptions>
              <el-descriptions direction="vertical">
                <el-descriptions-item label="标 签">
                  <el-tag type="info" size="small" v-for="(v, k) in metadata.labels" :key="k"
                    style="margin-left: 5px;">&nbsp;&nbsp;{{ k }}:{{
                      v }}&nbsp;&nbsp;</el-tag>
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-collapse-item>
          <el-collapse-item title="Resource information" name="2">
            <el-descriptions direction="vertical">
              <el-descriptions-item label="ADDRESS">
                <el-tag type="info" size="small" style="margin-left: 5px;" v-for="(v, k) in status.addresses " :key="k">{{
                  v.type }}:{{ v.address
  }}&nbsp;&nbsp;</el-tag>
              </el-descriptions-item>
            </el-descriptions>
            <el-descriptions direction="vertical">
              <el-descriptions-item label="Taints">
                <el-tag size="small" type="info" style="margin-left: 5px;" v-for="(v, k) in spec.taints " :key="k">{{
                  v.key }}:{{ v.effect
  }}&nbsp;&nbsp;</el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </el-collapse-item>
          <el-collapse-item title="System information" name="3">
            <el-descriptions direction="vertical" :column="4">
              <el-descriptions-item label="Machine ID">{{ status.nodeInfo.machineID }}</el-descriptions-item>
              <el-descriptions-item label="System UUID">{{ status.nodeInfo.systemUUID }}</el-descriptions-item>
              <el-descriptions-item label="Boot ID">{{ status.nodeInfo.bootID }}</el-descriptions-item>
              <el-descriptions-item label="Kernel version">{{ status.nodeInfo.kernelVersion }}</el-descriptions-item>
            </el-descriptions>
            <el-descriptions direction="vertical" :column="5">
              <el-descriptions-item label="OS Image">{{ status.nodeInfo.osImage }}</el-descriptions-item>
              <el-descriptions-item label="Container runtime version">{{ status.nodeInfo.containerRuntimeVersion
              }}</el-descriptions-item>
              <el-descriptions-item label="kubelet version">{{ status.nodeInfo.kubeletVersion }}</el-descriptions-item>
              <el-descriptions-item label="kube-proxy version">{{ status.nodeInfo.kubeProxyVersion
              }}</el-descriptions-item>
              <el-descriptions-item label="Architecture">{{ status.nodeInfo.architecture }}</el-descriptions-item>
              <el-descriptions-item label="Operating system">{{ status.nodeInfo.operatingSystem
              }}</el-descriptions-item>
              <el-descriptions-item label="CPU capacity">{{ status.allocatable.cpu }}</el-descriptions-item>
              <el-descriptions-item label="Memory capacity">{{ memory }}</el-descriptions-item>
              <el-descriptions-item label="Pods capacity">{{ status.allocatable.pods }}</el-descriptions-item>
            </el-descriptions>
          </el-collapse-item>
          <el-collapse-item title="Allocation" name="4">
          </el-collapse-item>
          <el-collapse-item title="Pods" name="5">
            <el-table :data="pods" style="width: 100%" size="small" :header-cell-style="{ background: '#e6e7e9' }">
              <el-table-column prop="Name" label="名称" width="200" />
              <el-table-column prop="Image" label="镜像" width="350" />
              <el-table-column label="标签" width="300">
                <template #default="scope">
                  <el-tag type="info" size="small" v-for="(v, k) in scope.row.Labels " :key="k">{{ k }}:{{ v
                  }}<br></el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="Status" label="状态" width="120" />
              <el-table-column prop="Restart" label="重启" width="120" />
              <el-table-column prop="PodAge" label="创建时间" />
            </el-table>
          </el-collapse-item>
          <el-collapse-item title="Events" name="6">
          </el-collapse-item>
        </el-collapse>
      </div>
    </el-col>
  </el-row>
</template>

<script>
// import * as echarts from 'echarts';
export default {
  name: '',
  data() {
    return {
      nodenames: '',
      node: {},
      metadata: {},
      spec: {},
      status: {
        nodeInfo: {},
        allocatable: {},
      },
      pods: [],
      age: '',
      activeNames: ['1'],
      memory: '',
      result: {
        total: 1,
        cpu: 1
      },
      value: '',
      nodelist: [],

    }
  },
  created() {
    this.getParams()
    this.getNodedetail()
  },
  mounted() {
    // this.getchart(this.result);//定义一个接口方法在methods中调用 
  },
  methods: {
    //获取路由参数的方法
    getParams() {
      // 取到路由带过来的参数 
      let routerParams = this.$route.query.nodename
      // 将数据放在当前组件的数据内
      this.nodenames = routerParams
      // 获取nodenamelist
      this.$ajax.get("/node/list").then((res) => {
        res.data.node_name.forEach(v => {
          this.nodelist.push({ 'value': v.metadata.name, 'label': v.metadata.name })
        })
      }).catch((res) => {
        console.log(res)
      })
    },
    //获取node详情的方法
    getNodedetail() {
      //获取单个node 详情
      this.$ajax({
        method: 'get',
        url: '/node/detail',
        params: {
          node_name: this.nodenames
        }
      }).then((res) => {
        this.metadata = res.data.detail.metadata
        this.spec = res.data.detail.spec
        this.status = res.data.detail.status
        this.age = res.data.age
        this.memory = res.data.memory_allocator
        this.pods = res.data.pods
      }).catch((res) => {

      })
    },
    //select 选中跳转的方法
    changeSelect(val) {
      this.nodenames = val
      this.getNodedetail()
    },
    //图表方法
    // getchart(result) {
    //   //获取id并初始化图表
    //   var cpu_request = echarts.init(document.getElementById('cpu_request'));
    //   var cpu_limit = echarts.init(document.getElementById('cpu_limit'));
    //   var memory_request = echarts.init(document.getElementById('memory_request'));
    //   var memory_limit = echarts.init(document.getElementById('memory_limit'));
    //   var pods = echarts.init(document.getElementById('pods'));
    //   //配置项
    //   var option = {
    //     color: ['#5fb878', '#eeeeee'],
    //     title: [{
    //       text: 'Pie',
    //       left: '85%',
    //     }, {
    //       subtext: 'alignTo',
    //       left: '50%',
    //       top: '85%',
    //       textAlign: 'center'
    //     }
    //     ],
    //     grid: {
    //       x2: 5,
    //       y2: 10
    //     },
    //     series: [
    //       {
    //         name: '访问来源',
    //         type: 'pie',
    //         radius: ['60%', '65%'],
    //         avoidLabelOverlap: false,
    //         hoverAnimation: false,   //关闭放大动画
    //         selectedOffset: 0,     //选中块的偏移量
    //         label: {
    //           show: false,
    //           position: 'center',
    //           formatter: '{d}%\nrequest'
    //         },
    //         emphasis: {
    //           label: {
    //             show: false,
    //             fontSize: '20',
    //             fontWeight: 'bold'
    //           }
    //         },
    //         labelLine: {
    //           show: false
    //         },
    //         data: (function () {
    //           var res = []
    //           res[0] = {
    //             value: result.total,
    //             name: '直接访问',
    //             label: {
    //               show: true,
    //               fontSize: '20',
    //               fontWeight: 'bold'
    //             }
    //           },
    //             res[1] = {
    //               value: result.cpu,
    //               name: '邮箱营销'
    //             }
    //           return res
    //         })()
    //       }
    //     ]
    //   };
    //   cpu_request.setOption(option)
    //   cpu_limit.setOption(option)
    //   memory_request.setOption(option)
    //   memory_limit.setOption(option)
    //   pods.setOption(option)
    // },
  }
}
</script>

<style>
.el-row {
  margin-bottom: 10px;

  &:last-child {
    margin-bottom: 0;
  }
}

.el-col {
  border-radius: 4px;
}

.bg-purple-dark {
  background: #99a9bf;
}

.bg-purple {
  background: #d3dce6;
}

.bg-purple-light {
  background: #e5e9f2;
}

.grid-content {

  border-radius: 4px;
  min-height: 36px;
}

.row-bg {
  padding: 10px 0;
  background-color: #f9fafc;
}

.el-collapse-item__header {

  background-color: #f0f2f5 !important;
}

/* .el-descriptions__body {
  background-color: #f0f2f5 !important;
} */

/* .el-collapse-item__content {
  background-color: #e5e9f2 !important;
} */

.el-table__body {
  background-color: #e5e9f2
}
</style>