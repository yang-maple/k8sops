<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content bg-purple-light">
        <el-collapse v-model="activeNames">
          <el-collapse-item title="Metadata" name="1">
            <div>
              <el-descriptions direction="vertical" :column="4">
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
              <el-descriptions-item>
                <el-descriptions direction="vertical">
                  <el-descriptions-item label="状 态">{{ status.phase }}</el-descriptions-item>
                </el-descriptions>
              </el-descriptions-item>
            </el-descriptions>
          </el-collapse-item>
          <el-collapse-item title="Resource Quotas" name="4">
            <el-descriptions direction="vertical">
              <el-descriptions-item>There is nothing to display here</el-descriptions-item>
            </el-descriptions>
          </el-collapse-item>
          <el-collapse-item title="Resource Limits" name="5">
            <el-descriptions direction="vertical">
              <el-descriptions-item>There is nothing to display here</el-descriptions-item>
            </el-descriptions>
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
      namespace_name: '',
      metadata: {},
      spec: {},
      status: {},
      age: '',
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
      let routerParams = this.$route.query.namespace_name
      // 将数据放在当前组件的数据内
      this.namespace_name = routerParams
      // 获取Namespace
      this.getNamespacedetail()
    },
    //获取Namespace详情的方法
    getNamespacedetail() {
      //获取单个Namespace 详情
      this.$ajax({
        method: 'get',
        url: '/namespaces/detail',
        params: {
          namespace_name: this.namespace_name
        }
      }).then((res) => {
        this.metadata = res.data.detail.metadata
        this.spec = res.data.detail.spec
        this.status = res.data.detail.status
        this.age = res.data.age
      }).catch((res) => {

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
  min-height: 36px;
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