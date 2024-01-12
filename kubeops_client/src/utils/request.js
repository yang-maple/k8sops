import axios from 'axios'
import { ElLoading } from 'element-plus'
let loading

function startLoading() {
    loading = ElLoading.service({
        lock: true,
        text: '努力加载中······',
        background: 'rgba(0, 0, 0, 0.7)'
    })
}

function endLoading() {
    loading.close()
}

let needLoadingRequestCount = 0

export function showFullScreenLoading() {
    if (needLoadingRequestCount === 0) {
        startLoading()
    }
    needLoadingRequestCount++
}
export function tryHideFullScreenLoading() {
    if (needLoadingRequestCount <= 0) return
    needLoadingRequestCount--
    if (needLoadingRequestCount === 0) {
        endLoading()
    }
}



//新建axios 对象
const httpClient = axios.create({
    validateStatus(status) {
        return status >= 200 && status < 504 //设置默认的合法状态码，不合法状态码不进行接收
    },
    timeout: 10000, //超时10s
    baseURL: "http://localhost:8090/v1/api/"
})

httpClient.defaults.retry = 3 // 请求重试次数
httpClient.defaults.retryDelay = 1000 // 重试时间间隔
httpClient.defaults.shouldRetry = false // 是否允许重试

// 添加请求拦截器
httpClient.interceptors.request.use(
    config => {
        //添加header 
        config.headers['Content-Type'] = 'application/json'
        config.headers['Accept-Language'] = 'zh-CN'
        config.headers['Authorization'] = localStorage.getItem("user_token") //可以全局设置接口请求 header 中带token
        config.headers['Uuid'] = localStorage.getItem("user_id")
        // 处理post 请求
        if (config.method === 'post') {
            if (!config.data) { //没有参数时，config.data 为null
                config.data = {}
            }
        }
        if (config.url == "/user/getCaptcha" || config.url == "/user/login") {
            return config
        }
        showFullScreenLoading()
        return config
    },
    err => {
        //Promise.reject() 方法返回一个带有拒绝原因的promise 对象，在 console 中显示报错
        return Promise.reject(err)
    }
)

// 添加响应拦截器
httpClient.interceptors.response.use(
    response => {
        tryHideFullScreenLoading()
        //处理状态码
        switch (response.status) {
            case 401:
                window.location.href = '/login'
                alert(response.data.msg)
                break;
            case 200:
                return response.data
            default:
                return Promise.reject(response.data)
        }
        // if (response.status != 200) {
        //     return Promise.reject(response.data)
        // } else {
        //     return response.data
        // }
    },
    err => {
        return Promise.reject(err)
    }
)

//导出

export default httpClient