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
                <el-descriptions direction="vertical" :column="5">
                  <el-descriptions-item label="状   态">{{ status.phase }}</el-descriptions-item>
                  <el-descriptions-item label="claimRef"><el-link type="primary" :underline="false">{{
                    spec.claimRef.namespace }}\{{ spec.claimRef.name
  }}</el-link></el-descriptions-item>
                  <el-descriptions-item label="ReclaimPolicy">{{ spec.persistentVolumeReclaimPolicy
                  }}</el-descriptions-item>
                  <el-descriptions-item label="storageClassNames">{{ spec.storageClassName }}</el-descriptions-item>
                  <el-descriptions-item label="accessModes">
                    <el-tag type="info" size="small" v-for="(v, k) in spec.accessModes" :key="k"
                      style="margin-left: 5px;">{{
                        v }}</el-tag>
                  </el-descriptions-item>
                </el-descriptions>

              </el-descriptions-item>
            </el-descriptions>
          </el-collapse-item>
          <el-collapse-item title="Source" name="4">
            <el-row>
              <el-col :span="8">
                <el-descriptions direction="vertical">
                  <el-descriptions-item label="Type">{{ source.type }}</el-descriptions-item>
                  <el-descriptions-item label="Server" v-if="source.type == 'NFS'">{{ source.server
                  }}</el-descriptions-item>
                  <el-descriptions-item label="Path">{{ source.path }}</el-descriptions-item>
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
                  <el-descriptions-item label="Resource name">storage</el-descriptions-item>
                  <el-descriptions-item label="Quantity">{{ spec.capacity.storage }}</el-descriptions-item>
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
      persistent_volume_name: '',
      metadata: {},
      spec: {
        capacity: {},
        nfs: null,
        accessModes: [],
        claimRef: {},
        hostPath: null,
      },
      status: {},
      age: '',
      source: {
        type: null,
        server: null,
        path: null
      },
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
      let routerParams = this.$route.query.persistent_volume_name
      // 将数据放在当前组件的数据内
      this.persistent_volume_name = routerParams
      // 获取Namespace
      this.getNamespacedetail()
    },
    //获取Namespace详情的方法
    getNamespacedetail() {
      //获取单个Namespace 详情
      this.$ajax({
        method: 'get',
        url: '/pv/detail',
        params: {
          persistent_volume_name: this.persistent_volume_name
        }
      }).then((res) => {
        this.metadata = res.data.detail.metadata
        this.spec = res.data.detail.spec
        this.status = res.data.detail.status
        this.age = res.data.age
        if (this.spec.nfs == undefined) {
          console.log(111)
          this.source.type = 'HostPath'
          this.source.path = this.spec.hostPath.path
        } else {
          console.log(222)
          this.source.type = 'NFS'
          this.source.server = this.spec.nfs.server
          this.source.path = this.spec.nfs.path
        }
        console.log(this.source.type)
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