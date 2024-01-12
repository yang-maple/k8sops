<template>
  <el-row style="padding-bottom: 10px;">
    <el-col :span="24">
      <el-card shadow="always" style="width: 100%;">
        <span>
          <div>
            <svg class="icon-node" aria-hidden="true">
              <use xlink:href="#icon-jisuanjiedian"></use>
            </svg>
            <span
              style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">集群节点</span>
            <br>
            <span style="font-size: 12px;color: #79879c!important">集群节点（Node）是指在 Kubernetes 集群中运行的 worker 节点。它们是实际运行 Pod
              的物理或虚拟机器。
            </span>
          </div>
        </span>
      </el-card>
    </el-col>
  </el-row>
  <div class="table-bg-purple">
    <el-table :data="nodeItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
      :default-sort="{ prop: 'date', order: 'descending' }">
      <el-table-column label="名称" width="200">
        <template #default="scope">
          <el-row>
            <el-col :span="4"><svg class="table-icon" aria-hidden="true">
                <use xlink:href="#icon-jisuanjiedian"></use>
              </svg></el-col>
            <el-col :span="20">
              <el-text style="margin-left: 5px" size="small" tag="b">{{
                scope.row.name
              }}</el-text>
              <div style="margin-top: -5px;">
                <el-text size="small" style="margin-left: 5px;">{{
                  scope.row.nodeIp
                }}</el-text>
              </div>
            </el-col>
          </el-row>
        </template>
      </el-table-column>
      <el-table-column label="标签">
        <template #default="scope">
          <div v-for="(v, k, index) in scope.row.labels " :key="k">
            <div v-if="index < maxitem[scope.$index]">
              <el-tag type="info" style="margin-left: 5px;" size="small" effect="plain" round>
                {{ k }}:{{ v }}</el-tag>
            </div>
          </div>
          <div v-if="scope.row.labels == null" align="center">---</div>
          <div v-if="scope.row.labels != null"><el-button size="small" type="primary" link
              @click="showLabels(scope.$index)">{{
                maxitem[scope.$index] == 3 ?
                '展开' : '收起'
              }}</el-button></div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="150">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            <span slot="reference" v-if="scope.row.status == 'Ready' && scope.row.unschedulable == false">
              <el-tooltip placement="top" effect="light"><template #content>Ready</template>
                <i class="dotClass" style="background-color: springgreen"></i></el-tooltip>
            </span>
            <span slot="reference" v-if="scope.row.status == 'Ready' && scope.row.unschedulable == true">
              <el-tooltip placement="top" effect="light"><template #content>Unschedulable</template>
                <i class="dotClass" style="background-color: orange"></i></el-tooltip>
            </span>
            <span slot="reference" v-if="scope.row.status != 'Ready'">
              <el-tooltip placement="top" effect="light"><template #content> NotReady </template>
                <i class="dotClass" style="background-color: red"></i></el-tooltip>
            </span>
            <span>
              <el-text v-if="scope.row.status == 'Ready' && scope.row.unschedulable == false"
                style="margin-left: 3px;color: green;">
                运行中
              </el-text>
              <el-text v-if="scope.row.status == 'Ready' && scope.row.unschedulable == true"
                style="margin-left: 3px;color: orangered;">
                不可调度
              </el-text>
              <el-text v-if="scope.row.status != 'Ready'" style="margin-left: 3px;color: red;">
                异常
              </el-text>
              <span v-if="scope.row.taints != null">
                <el-tooltip class="box-item" effect="light" placement="top">
                  <template #content>
                    污点：<br />
                    <span v-for="(v, k) in scope.row.taints" :key="k">
                      {{ v.key }}:{{ v.effect }}<br />
                    </span>
                  </template>
                  <el-button size="small" type="info">
                    <template #default>
                      {{ scope.row.taints.length }}
                    </template>
                  </el-button>
                </el-tooltip>
              </span>
            </span>

          </div>
        </template>
      </el-table-column>
      <el-table-column prop="cpu" label="CPU（m）" width="100">
        <template #default="scope">
          <div>
            <el-tooltip class="box-item" effect="light" placement="top">
              <template #content>
                <span>{{ resource["cpu_request"][scope.row.name] }} m</span>
              </template>
              <span>使用：{{ toNumber(resource["cpu_request"][scope.row.name] * 100 / scope.row.cpu_total) }}% </span>
            </el-tooltip>
          </div>
          <div>
            <el-tooltip class="box-item" effect="light" placement="top">
              <template #content>
                <span>{{ resource["cpu_limit"][scope.row.name] }} m</span>
              </template>
              <span>限制：{{ toNumber(resource["cpu_limit"][scope.row.name] * 100 / scope.row.cpu_total) }}% </span>
            </el-tooltip>
          </div>
          <div>
            <span>总量：{{ scope.row.cpu_total }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="memory" label="内存（Mi）" width="100">
        <template #default="scope">
          <div>
            <el-tooltip class="box-item" effect="light" placement="top">
              <template #content>
                <span> {{ resource["memory_request"][scope.row.name] }} Mi</span>
              </template>
              <span>使用：{{ toNumber(resource["memory_request"][scope.row.name] * 100 / scope.row.memory_total) }}% </span>
            </el-tooltip>
          </div>
          <div>
            <el-tooltip class="box-item" effect="light" placement="top">
              <template #content>
                <span> {{ resource["memory_limit"][scope.row.name] }} Mi</span>
              </template>
              <span>限制：{{ toNumber(resource["memory_limit"][scope.row.name] * 100 / scope.row.memory_total) }}% </span>
            </el-tooltip>
          </div>
          <div>
            <span>总量：{{ scope.row.memory_total }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="pods" align="center" label="pod 数量" width="80">
      </el-table-column>
      <el-table-column prop="create_time" label="创建时间" width="140" />
      <el-table-column prop="kubelet_version" label="版本" width="80" align="center" />
      <el-table-column label="操作" width="50" align="center">
        <template #default="scope">
          <el-dropdown trigger="click">
            <el-button type="primary" link>
              <el-icon :size="24">
                <Operation />
              </el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="setNodeSchedule(scope.row)">
                  <el-icon>
                    <component :is="scope.row.unschedulable == false ? 'VideoPause' : 'VideoPlay'" />
                  </el-icon>
                  {{ scope.row.unschedulable == false ? '停止调度' : '启动调度' }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script scoped>
import { ElMessage } from 'element-plus'
export default {
  data() {
    return {
      maxitem: [],
      nodeItem: [],
      resource: {},
    }
  },
  created() {
    this.GetNodelistaxios()
  },
  methods: {
    GetNodelistaxios() {
      this.$ajax.get("/node/list").then((res) => {
        console.log(res.data);
        this.nodeItem = res.data.Item
        this.resource = res.data.resource
        for (let i = 0; i < res.data.Item.length; i++) {
          this.maxitem.push(3)
        }
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
    showLabels(index) {
      if (this.maxitem[index] == 3) {
        this.maxitem[index] = 99
      } else {
        this.maxitem[index] = 3
      }
    },
    setNodeSchedule(row) {
      console.log(row)
      if (row.unschedulable == false) {
        this.$ajax.post("/node/schedule", {
          node_name: row.name,
          status: true
        }).then((res) => {
          ElMessage({
            message: `节点 ${row.name} 已不可调度`,
            type: 'warning',
          })
          setTimeout(() => {
            this.GetNodelistaxios()
          }, 1000)
        }).catch((res) => {
          console.log(res.data);
        })
        return
      }
      this.$ajax.post("/node/schedule", {
        node_name: row.name,
        status: false
      }).then((res) => {
        ElMessage({
          message: `节点 ${row.name} 已可调度`,
          type: 'success',
        })
        setTimeout(() => {
          this.GetNodelistaxios()
        }, 1000)
      }).catch((res) => {
        console.log(res.data);
      })
    },
    toNumber(val) {
      return Math.floor(val)
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

.icon-node {
  width: 3em;
  height: 3em;
  vertical-align: -1em;
  fill: currentColor;
  overflow: hidden;
}

.table-icon {
  width: 3em;
  height: 3em;
  vertical-align: -1.95em;
  fill: currentColor;
  overflow: hidden;
}
</style>