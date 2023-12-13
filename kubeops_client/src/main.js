import { createApp } from 'vue'
//导入app.vue 主组件
import App from './App.vue'

//导入路由
import router from "./router";


//导入 axios
import axiosjs from './utils/request.js';
//导入 ElementPlus 以及 Css 文件
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
//导入图标视图
import * as ElIcons from '@element-plus/icons-vue'
//导入particles
import Particles from 'particles.vue3'
import 'default-passive-events'

const app = createApp(App)
app.use(router)
app.use(Particles)
app.mount('#app')

// 将图标引入成全局组件
for (let iconName in ElIcons) {
    app.component(iconName, ElIcons[iconName])
}

// 引用 element-plus
app.use(ElementPlus)
app.config.globalProperties.$ajax = axiosjs