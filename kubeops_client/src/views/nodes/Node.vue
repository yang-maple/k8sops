<template>
  <div class="table-bg-purple">
    <el-table :data="nodeItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
      :default-sort="{ prop: 'date', order: 'descending' }">
      <el-table-column label="名称" min-width="150" align="center">
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
            <el-link style="margin-left: 10px" type="primary" :underline="false" @click="handle(scope.row)">{{
              scope.row.name
            }}</el-link>
            <!-- <span style="margin-left: 10px">{{ scope.row.metadata.name }}</span> -->
          </div>
        </template>
      </el-table-column>
      <el-table-column label="标签" width="350" align="center">
        <template #default="scope">
          <el-tag type="info" size="small" v-for="(v, k) in scope.row.labels " :key="k">{{ k }}:{{ v
          }}<br></el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" align="center" width="100" />
      <el-table-column prop="cpu" align="center" label="CPU(核)" width="100" />
      <el-table-column prop="memory" align="center" label="可分配内存" width="100">
      </el-table-column>
      <el-table-column prop="pods" align="center" label="pod 数量" width="100">
        <!-- <template #default="scope">
          <span> {{ node.pod[scope.row.metadata.name] }}</span>
        </template> -->
      </el-table-column>
      <el-table-column prop="create_time" align="center" label="创建时间" width="200" />
      <el-table-column prop="kubelet_version" label="版本" align="center" />
    </el-table>
  </div>
</template>

<script scoped>
export default {
  data() {
    return {
      nodeItem: [],
      msg: 'test message',
    }
  },
  created() {
    this.GetNodelistaxios()
  },
  methods: {
    GetNodelistaxios() {
      this.$ajax.get("/node/list").then((res) => {
        this.nodeItem = res.data.Item
        console.log(res.data);
      }).catch(function (res) {
        console.log(res.data);
      })
    },

    handle(row) {
      this.$router.push({
        path: '/cluster/node_detail',
        name: 'node详情',
        query: {
          nodename: row.metadata.name
        }
      })
    },
  }
}
</script>


<style>
/* .el-table__body tr:hover>td {
  background-color: #f0f9eb !important;
} */

.table-bg-purple {
  padding-right: 2px;
  padding-left: 2px;
  border-radius: 4px;
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

.el-button--text {
  margin-right: 15px;
}


.dialog-footer button:first-child {
  margin-right: 10px;
}





.demo-pagination-block+.demo-pagination-block {
  float: right;
  margin-right: 10px;
  margin-top: 10px;
}

.dotClass {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: block;
  margin-left: 10px;
}

/* .el-form-item__content .el-input {
    width: 500px;
} */
</style>