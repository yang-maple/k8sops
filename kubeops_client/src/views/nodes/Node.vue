<template>
  <el-row>
    <el-col :span="24">
      <div class="grid-content bg-purple">
        <el-table :data="tableData" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
          :default-sort="{ prop: 'date', order: 'descending' }">
          <el-table-column label="名称" width="150" align="center">
            <template #default="scope">
              <div style="display: flex; align-items: center">
                <span slot="reference" v-if="scope.row.status.conditions[4].type == 'Ready'">
                  <el-tooltip placement="bottom" effect="light"><template #content>Ready</template>
                    <i class="dotClass" style="background-color: springgreen"></i></el-tooltip>
                </span>
                <span slot="reference" v-if="scope.row.status.conditions[4].type != 'Ready'">
                  <el-tooltip placement="bottom" effect="light"><template #content> NotReady </template>
                    <i class="dotClass" style="background-color: red"></i></el-tooltip>
                </span>
                <el-link style="margin-left: 10px" type="primary" :underline="false" @click="handle(scope.row)">{{
                  scope.row.metadata.name
                }}</el-link>
                <!-- <span style="margin-left: 10px">{{ scope.row.metadata.name }}</span> -->
              </div>
            </template>
          </el-table-column>
          <el-table-column label="标签" width="350" align="center">
            <template #default="scope">
              <el-tag type="info" size="small" v-for="(v, k) in scope.row.metadata.labels " :key="k">{{ k }}:{{ v
              }}<br></el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status.conditions[4].type" label="状态" align="center" width="100" />
          <el-table-column prop="status.allocatable.cpu" align="center" label="CPU(核)" width="100" />
          <el-table-column align="center" label="可分配内存" width="100">
            <template #default="scope">
              <span> {{ node.memory[scope.row.metadata.name] }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" label="pod 数量" width="100">
            <template #default="scope">
              <span> {{ node.pod[scope.row.metadata.name] }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" align="center" label="创建时间" width="200" />
          <el-table-column prop="status.nodeInfo.kubeletVersion" label="版本" />
        </el-table>
      </div>
    </el-col>
  </el-row>
</template>

<script scoped>
export default {
  data() {
    return {
      tableData: [],
      node: {
        pod: {},
        memory: {},
      },
      msg: 'test message',
    }
  },
  created() {
    this.GetNodelistaxios()
  },
  methods: {
    formatter(row, column) {
      return row.address;
    },
    // tableRowClassName({ row, rowIndex }) {
    //   if (rowIndex === 0) {
    //     return 'warning-row';
    //   } else if (rowIndex === 2) {
    //     return 'success-row';
    //   }
    //   return '';
    // },
    Getnodes(name) {
      this.$ajax({
        method: 'get',
        url: '/node/detail',
        params: {
          node_name: name
        }
      }).then((res) => {
        this.node.pod[name] = res.data.total
        this.node.memory[name] = res.data.memory_allocator
      }).catch(function (res) {
        console.log(res.data);
      })
    },

    GetNodelistaxios() {
      this.$ajax.get("/node/list").then((res) => {
        this.tableData = res.data.node_name
        this.tableData.forEach(v => {
          this.Getnodes(v.metadata.name)
        })
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

.el-row {
  margin-bottom: 100px;

  &:last-child {
    margin-bottom: 0;
  }
}

.el-col {
  border-radius: 4px;
}

.bg-purple {
  background: #f0f2f5;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
}

.grid-content {
  border-radius: 4px;
  min-height: 10px;
  padding: 10px;
}

.el-row:last-child {
  margin-bottom: 0;
}

.grid-content1 {
  padding-top: 10px;
  margin-right: 10px;
  border-radius: 4px;
}

.grid-content2 {
  border-radius: 4px;
}

.grid-content3 {
  border-radius: 4px;
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

.el-dialog--center {
  min-height: 400px;
  min-width: 200px;
}

.el-dialog__body {
  min-height: 400px;
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