
<template>
  <div class="home-header">
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
      <el-col :span="24">

      </el-col>
    </el-row>
  </div>
  <div class="home-nodes">
    <el-row>
      <el-col :span="6">
        <div id="node-chart" class="node-chart"></div>
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
      nodeItem: [],
      option: {
        title: {
          text: 'Referer of a Website',
          subtext: 'Fake Data',
          left: 'center'
        },
        tooltip: {
          trigger: 'item'
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: 'Access From',
            type: 'pie',
            radius: '50%',
            data: [],
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      }
    }
  },
  computed: {
  },
  created() {
    this.initdata()
  },
  mounted() {
    // this.initEchart()
  },
  methods: {
    initEchart() {
      var myChart = echarts.init(document.getElementById('node-chart'));
      myChart.setOption(this.option)
    },
    initdata() {
      this.$ajax({
        method: 'get',
        url: '/homepage/getInfo',
      }).then((res) => {
        this.nodeItem = res.data.cluster_info
        console.log(res.data)
      }).catch((res) => {
        console.log(res);
      })
      console.log(this.nodeItem)
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
  text-align: center;
  background: #f0f2f5;
}

.node-chart {
  width: 600px;
  height: 400px;

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
</style>