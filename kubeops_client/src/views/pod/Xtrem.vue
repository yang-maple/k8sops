<template>
    <el-row>
        <el-col :span="24">
            <div class="grid-content bg-purple">
                <el-row>
                    <el-col :span="1">
                        <div class="grid-content3">shell</div>
                    </el-col>
                    <el-col :span="3">
                        <div class="grid-content ">
                            <el-select v-model="container_name" filterable placeholder="Select Container"
                                @change="switchContainer()">
                                <el-option v-for="(value, index) in containerItem" :key="index" :label="value"
                                    :value="value" style="width:100%" />
                            </el-select>
                        </div>
                    </el-col>
                    <el-col :span="1" class="grid-content2">
                        In
                    </el-col>
                    <el-col :span="5" class="grid-content ">
                        <el-input v-model="pod_name" disabled placeholder="Pod Name" />
                    </el-col>
                    <el-col :span="7">
                    </el-col>
                    <el-col :span="6">
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
import { AttachAddon } from 'xterm-addon-attach'
import "xterm/css/xterm.css"
export default {
    data() {
        return {
            term: "", // 保存terminal实例
            socketUri: 'ws://127.0.0.1:8081/ws?namespace=test&pod_name=ops-64786f5c8f-np5q7&container_name=qqq',
            rows: 100,
            cols: 135,
            data: '',
            submitdata: '',
            socket: '',
            jsdata: {},
            containerItem: [],
            pod_name: '',
            container_name: '',
            namespace: '',
        }
    },
    created() {
    },
    mounted() {
        this.getParams()
    },
    methods: {
        initXterm() {
            let _this = this
            let term = new Terminal({
                rendererType: "canvas", //渲染类型
                rows: _this.rows, //行数
                cols: _this.cols, // 不指定行数，自动回车后光标从下一行开始
                convertEol: true, //启用时，光标将设置为下一行的开头
                scrollback: 50, //终端中的回滚量
                disableStdin: false, //是否应禁用输入
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
            // canvas背景全屏
            const fitAddon = new FitAddon()
            term.loadAddon(fitAddon)
            fitAddon.fit()
            window.addEventListener("resize", resizeScreen)
            term.resize(_this.cols, _this.rows)
            function resizeScreen() {
                try { // 窗口大小改变时，触发xterm的resize方法使自适应
                    let order = JSON.stringify({
                        "operation": "resize",
                        "data": '',
                        "rows": _this.rows,
                        "cols": _this.cols,
                    })
                    _this.socket.send(order)
                } catch (e) {
                    console.log("e", e.message)
                }
            }
            _this.term = term
            if (this.container_name != "") {
                _this.runFakeTerminal(false)
            }
        },
        runFakeTerminal(init) {
            let _this = this
            let term = _this.term
            term._initialized = init
            if (term._initialized) return
            // 初始化
            term._initialized = true
            this.socket = new WebSocket(`ws://127.0.0.1:8081/ws?namespace=${this.namespace}&pod_name=${this.pod_name}&container_name=${this.container_name}`);

            // 换行并输入起始符
            term.prompt = _ => {
                if (this.jsdata.data != undefined) {
                    term.write(`\r\n${this.jsdata.data}\x1b[0m`)
                }
            }
            term.prompt()
            term.writeln("Welcome to \x1b[1;32m墨天轮\x1b[0m.")
            term.writeln('This is Web Terminal of Modb; Good Good Study, Day Day Up.')
            term.prompt()
            term.focus();

            // 添加事件监听器，支持输入方法
            term.onKey(e => {
                const printable = !e.domEvent.altKey && !e.domEvent.altGraphKey && !e.domEvent.ctrlKey && !e.domEvent.metaKey
                if (e.domEvent.keyCode === 13) {
                    term.prompt()
                } else if (e.domEvent.keyCode === 8) { // back 删除的情况
                    if (term._core.buffer.x > 2) {
                        term.write('\b \b')
                    }
                } else if (printable) {
                    term.write(e.key)
                }
                if (e.key != '\r') {
                    this.data = this.data + e.key
                }
                if (e.key == '\r') {
                    this.submitdata = this.data
                    this.data = ''
                }
            })
            term.onData(key => {  // 粘贴的情况
                if (key.length > 1) {
                    term.write(key)
                    this.data = this.data + key
                }
                if (this.submitdata != "\r" && this.submitdata != "") {
                    let order = JSON.stringify({
                        "operation": "stdin",
                        "data": this.submitdata,
                        "rows": 0,
                        "cols": 0
                    })
                    console.log(order)
                    _this.socket.send(order)
                }
            })
            const messageHandler = (event) => {
                // Handle the received data here
                this.jsdata = JSON.parse(event.data)
                if (this.jsdata.data != (this.submitdata + '\r\n')) {
                    this.submitdata = ""
                    _this.term.write(`\r\n${this.jsdata.data}\x1b[0m`)
                    // _this.term.write("\r\n\x1b[33m$\x1b[0m ")
                }
            };
            let resize = JSON.stringify({
                "operation": "resize",
                "data": '',
                "rows": _this.rows * 100,
                "cols": _this.cols,
            })
            _this.socket.removeEventListener('message', messageHandler, true);
            _this.socket.addEventListener('message', messageHandler);
            term.onSelectionChange(() => {
                _this.socket.send(resize)
                term.focus();
            })
            term.onLineFeed(() => {
            })
            term.onBinary((data) => {
                console.log('Received binary data:', data);
            })
        },
        switchContainer() {
            let _this = this
            this.jsdata = {}
            this.term.writeln(`\r\n`)
            this.term.clear()
            _this.runFakeTerminal(false)
        },
        //获取路由参数的方法
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