const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  devServer: {
    host: '0.0.0.0',//监听地址
    port: 7070,//监听端口
    open: true, //启动后自动打开页面
    client: {
      overlay: false,
    }

  },
  transpileDependencies: true,

  //关闭语法检测
  lintOnSave: false
})
