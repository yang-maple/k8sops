<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content bg-purple-light">
        <el-collapse v-model="activeNames">
          <el-collapse-item title="Metadata" name="1">
            <div>
              <el-descriptions direction="vertical" :column="4">
                <el-descriptions-item label="名   称">{{ metadata.name }}</el-descriptions-item>
                <el-descriptions-item label="名称空间">{{ metadata.namespace }}</el-descriptions-item>
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
              <el-descriptions-item>
                <el-descriptions direction="vertical" :column="5">
                  <el-descriptions-item label="类  别">{{ spec.type }}</el-descriptions-item>
                  <el-descriptions-item label="集群 IP">{{ spec.clusterIP }}</el-descriptions-item>
                  <el-descriptions-item label="分发策略">{{ spec.sessionAffinity }}</el-descriptions-item>
                </el-descriptions>
                <el-descriptions direction="vertical">
                  <el-descriptions-item label="标 签">
                    <el-tag type="info" size="small" v-for="(v, k) in spec.selector" :key="k"
                      style="margin-left: 5px;">&nbsp;&nbsp;{{ k }}:{{
                        v }}&nbsp;&nbsp;</el-tag>
                  </el-descriptions-item>
                </el-descriptions>
              </el-descriptions-item>
            </el-descriptions>
          </el-collapse-item>
          <el-collapse-item title="Source" name="4">
            <el-row>
              <el-col :span="8">
                <el-descriptions direction="vertical">
                  <el-descriptions-item label="Type"></el-descriptions-item>
                </el-descriptions>
              </el-col>
              <el-col :span="16">
                <div class="grid-content1" />
              </el-col>
            </el-row>
          </el-collapse-item>
          <el-collapse-item title="Capacity" name="5">
            <el-row>
              <el-col :span="4">
                <el-descriptions direction="vertical">
                </el-descriptions>
              </el-col>
              <el-col :span="20">
                <div class="grid-content1" />
              </el-col>
            </el-row>
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
      service_name: null,
      namespace: null,
      metadata: {},
      spec: {},
      status: {},
      age: null,
      source: {},
      activeNames: ['1'],
    }
  },
  created() {
    this.getParams()
  },
  mounted() {
    // this.getchart(this.result);//定义一个接口方法在methods中调用 
  },
  methods: {
    //获取路由参数的方法
    getParams() {
      // 取到路由带过来的参数 
      this.service_name = this.$route.query.service_name
      this.namespace = this.$route.query.namespace
      // 将数据放在当前组件的数据内
      // 获取Namespace
      this.getServicedetail()
    },
    //获取Namespace详情的方法
    getServicedetail() {
      //获取单个Namespace 详情
      this.$ajax({
        method: 'get',
        url: '/svc/detail',
        params: {
          service_name: this.service_name,
          namespace: this.namespace
        }
      }).then((res) => {
        this.metadata = res.data.detail.metadata
        this.spec = res.data.detail.spec
        this.status = res.data.detail.status
        this.age = res.data.age
        console.log(res.data)
      }).catch((res) => {
        console.log(res)
      })
    },
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
  min-height: 10px;
}

.grid-content1 {
  min-height: 65px;
  background: #f0f2f5;
}

.row-bg {
  padding: 10px 0;
  background-color: #f9fafc;
}

.el-collapse-item__header {

  background-color: #d3dce6 !important;
}

.el-descriptions__body {
  background-color: #f0f2f5 !important;
}

/* .el-collapse-item__content {
  background-color: #e5e9f2 !important;
} */

.el-table__body {
  background-color: #e5e9f2
}
</style>