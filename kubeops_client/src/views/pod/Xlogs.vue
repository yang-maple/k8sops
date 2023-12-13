<template>
    <el-row>
        <el-col :span="24">
            <div class="grid-content bg-purple">
                <el-row>
                    <el-col :span="1" class="grid-content3">
                        Logs
                    </el-col>
                    <el-col :span="3">
                        <div class="grid-content">
                            <el-select v-model="container_name" filterable placeholder="Select Container"
                                @change="switchContainer()">
                                <el-option v-for="(value, index) in containerItem" :key="index" :label="value"
                                    :value="value" />
                            </el-select>
                        </div>
                    </el-col>
                    <el-col :span="1" class="grid-content2">
                        In
                    </el-col>
                    <el-col :span="5" class="grid-content">
                        <el-input v-model="pod_name" disabled placeholder="Pod Name" />
                    </el-col>
                    <el-col :span="11">
                    </el-col>
                    <el-col :span="2" class="grid-content">
                        <el-select v-model="refreshtimeValue" placeholder="5s" @change="switchRefreshtime()">
                            <template #prefix>
                                <el-icon>
                                    <Timer />
                                </el-icon>
                            </template>
                            <el-option v-for="item in options" :key="item.refreshtimeValue" :label="item.label"
                                :value="item.refreshtimeValue" />
                        </el-select>
                    </el-col>
                </el-row>
                <div id="terminal" ref="terminal"></div>
            </div>
        </el-col>
    </el-row>
</template>
<script>
import { Terminal } from "xterm"
import { FitAddon } from 'xterm-addon-fit'
import { ElLoading } from 'element-plus'
import "xterm/css/xterm.css"
export default {
    data() {
        return {
            term: "", // 保存terminal实例
            rows: 40,
            cols: 100,
            namespace: '',
            pod_name: '',
            container_name: '',
            data: '',
            positionTimer: null,
            initXterms: false,
            containerItem: [],
            refreshtimeValue: null,
            options: [
                {
                    refreshtimeValue: 1000,
                    label: '1s',
                },
                {
                    refreshtimeValue: 5000,
                    label: '5s',
                },
                {
                    refreshtimeValue: 30000,
                    label: '30s',
                },
                {
                    refreshtimeValue: 60000,
                    label: '1min',
                },
                {
                    refreshtimeValue: 300000,
                    label: '5min',
                },
            ]
        }
    },
    mounted() {

        /*
          最初始情况，项目刚打开的时候，这个时候页面是必定没有定时器的，那么逻辑就会走else，这个时候就会注册一个定时器去循环调用相应逻辑代码
          后续有三种情况
             情况一：路由跳转，跳走了，就要清除这个定时器，所以在beforeDestroy钩子中要清除定时器
             情况二：关闭项目，关闭项目了以后，系统就会自动停掉定时器，这个不用管它
             情况三：刷新页面，这个时候vue组件的钩子是不会执行beforeDestroy和destroyed钩子的，所以
                    需要加上if判断一下，若还有定时器的话，就清除掉，所以这个就是mounted钩子的if判断的原因
       */
        if (this.positionTimer) {
            clearInterval(this.positionTimer);
            this.positionTimer = null
        } else {
            this.getParams()
        }
    },
    unmounted() {
        clearInterval(this.positionTimer);
        this.positionTimer = null;
    },
    methods: {
        initXterm() {
            let _this = this
            let term = new Terminal({
                rendererType: "canvas", //渲染类型
                rows: _this.rows, //行数
                cols: _this.cols, // 不指定行数，自动回车后光标从下一行开始
                convertEol: true, //启用时，光标将设置为下一行的开头
                // scrollback: 50, //终端中的回滚量
                disableStdin: true, //是否应禁用输入
                // cursorStyle: "underline", //光标样式
                cursorBlink: true, //光标闪烁
                theme: {
                    foreground: "#ECECEC", //字体
                    background: "#000000", //背景色
                    cursor: "help", //设置光标
                    lineHeight: 20
                }
            })
            // 创建terminal实例
            term.open(this.$refs["terminal"])
            // 换行并输入起始符 $
            term.prompt = _ => {
            }
            term.prompt()
            // canvas背景全屏
            const fitAddon = new FitAddon()
            term.loadAddon(fitAddon)
            fitAddon.fit()

            window.addEventListener("resize", resizeScreen)
            function resizeScreen() {
                try { // 窗口大小改变时，触发xterm的resize方法使自适应
                    fitAddon.fit()
                } catch (e) {
                    console.log("e", e.message)
                }
            }
            _this.term = term
            _this.refreshtimeValue = 5000
            _this.addSetInterval()
            _this.Loading()
        },
        runFakeTerminal() {
            let _this = this
            let term = _this.term
            _this.initXterms = true
            // 初始化
            term._initialized = _this.initXterms
            term.writeln(_this.data)
            term.focus()
            if (term._initialized) {
                term.clear()
            }
        },
        getLog() {
            this.$ajax({
                method: 'get',
                url: '/pod/container/log',
                params: {
                    namespace: this.namespace,
                    pod_name: this.pod_name,
                    container_name: this.container_name
                }
            }).then((res) => {
                this.data = res.data
                this.runFakeTerminal()
            }).catch(function (res) {
                console.log(res);
            })
        },
        addSetInterval() {
            const that = this  // 声明一个变量指向vue实例this,保证作用域一致
            let time = 1000
            if (that.refreshtimeValue != 1000) {
                time = that.refreshtimeValue
            }
            that.positionTimer = setInterval(() => {
                this.getLog()
                console.log(Date().toLocaleLowerCase())
            }, time)
        },
        Loading() {
            const loading = ElLoading.service({
                lock: true,
                text: "努力加载中....",
                background: 'rgba(0, 0, 0, 0.7)',
            })
            let time = 5000
            if (this.refreshtimeValue == 1000) {
                time = 1000
            }
            console.log(time)
            setTimeout(() => {
                loading.close()
            }, time)
        },
        switchContainer() {
            let _this = this
            clearInterval(_this.positionTimer);
            _this.positionTimer = null;
            _this.data = {}
            _this.term.writeln(`\r\n`)
            _this.addSetInterval()
            _this.Loading()
        },
        getParams() {
            // 取到路由带过来的参数,将数据放在当前组件的数据内
            this.namespace = this.$route.query.namespace
            this.pod_name = this.$route.query.pod_name
            this.$ajax({
                method: 'get',
                url: '/pod/container/list',
                params: {
                    namespace: this.namespace,
                    pod_name: this.pod_name
                }
            }).then((res) => {
                this.containerItem = res.item
                this.container_name = res.item[0]
                this.initXterm()
            }).catch(function (res) {
                console.log(res);
            })
            console.log("ing")
        },
        switchRefreshtime() {
            let _this = this
            clearInterval(_this.positionTimer);
            _this.positionTimer = null;
            _this.data = {}
            _this.term.writeln(`\r\n`)
            _this.addSetInterval()
            _this.Loading()
        }
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

.bg-purple {
    background: #f0f2f5;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
}

.grid-content3 {
    margin-left: 5px;
    margin-top: 10px;
}

.grid-content2 {
    margin-left: 30px;
    margin-top: 10px;
}

.grid-content {
    margin-top: 5px;
    border-radius: 4px;
    min-height: 10px;
}

.el-row:last-child {
    margin-bottom: 0;
}
</style>