<template>
  <!-- 用来放echarts图表的容器，一定要有宽高 -->
  <div style="height: 250px;width: 250px; float: left; overflow: hidden;">
    <div id="main" style="height: 250px;width: 250px;"></div>
  </div>
  <div style="height: 250px;width: 250px; float: left; overflow: hidden;">
    <div id="cpu" style="height: 250px;width: 250px;"></div>
  </div>
</template>


<script>
import * as echarts from 'echarts';
export default {
  data() {
    return {
      result: {
        total: 1,
        cpu: 1
      }
    }
  },
  mounted() {
    this.getchart(this.result);//定义一个接口方法在methods中调用 
  },
  methods: {

    //图表方法
    getchart(result) {
      //获取id并初始化图表
      var chart = echarts.init(document.getElementById('main'));
      var cpu = echarts.init(document.getElementById('cpu'));
      //配置项
      var option = {
        color: ['#5fb878', '#eeeeee'],
        title: [{
          text: 'Pie',
          left: '85%',
        }, {
          subtext: 'alignTo',
          left: '50%',
          top: '85%',
          textAlign: 'center'
        }
        ],

        series: [

          {
            name: '访问来源',
            type: 'pie',
            radius: ['60%', '70%'],
            avoidLabelOverlap: false,
            hoverAnimation: false,   //关闭放大动画
            selectedOffset: 0,     //选中块的偏移量
            label: {
              show: false,
              position: 'center',
              formatter: '{d}%\nrequest'
            },
            emphasis: {
              label: {
                show: false,
                fontSize: '20',
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: false
            },
            data: (function () {
              var res = []
              res[0] = {
                value: result.total,
                name: '直接访问',
                label: {
                  show: true,
                  fontSize: '20',
                  fontWeight: 'bold'
                }
              },
                res[1] = {
                  value: result.cpu,
                  name: '邮箱营销'
                }
              return res
            })()
          }
        ]
      };

      chart.setOption(option)
      cpu.setOption(option)
    },
  },
}
</script>
<style scoped>
.generalStyle {
  width: 100%;
  height: 90%;
}
</style>