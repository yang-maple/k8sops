<template>
    <div class="login">
        <Particles id="tsparticles" class="login__particles" :options="options" />
        <div class="loginPart">
            <h2>登录</h2>
            <el-form ref="ruleFormRef" :model="user" status-icon :rules="rules" label-width="100px" class="demo-ruleForm"
                style="transform:translate(-30px);">
                <el-form-item label="账号：" prop="account">
                    <el-input v-model="username" placeholder="请输入账号" maxlength="20" clearable />
                </el-form-item>
                <el-form-item label="密码：" prop="password">
                    <el-input v-model="password" type="password" placeholder="请输入密码" maxlength="20" show-password
                        clearable />
                </el-form-item>
                <el-form-item label="验证码：" prop="verifyCode">
                    <el-input style="width: 120px;" v-model="user.verifyCode" placeholder="请输入验证码" maxlength="4"
                        clearable />
                    <img id="image" class="verifyCodeImg" :src="imgUrl" @click="GetCaptcha()">
                </el-form-item>
                <el-button class="btn" type="primary" @click="handleSubmit">登录</el-button>
                <div style="text-align: right;transform:translate(0,30px);">
                    <el-link type="warning">没有账号？去注册</el-link>
                </div>
            </el-form>
        </div>

    </div>
</template>


<script>
export default {
    data() {
        return {
            username: null,
            password: null,
            user: {
                account: '',
                password: '',
                verifyCode: ''
            },
            captchaId: null,
            imgUrl: null,
            options: {
                fpsLimit: 60,
                interactivity: {
                    detectsOn: 'canvas',
                    events: {
                        onClick: { // 开启鼠标点击的效果
                            enable: true,
                            mode: 'push'
                        },
                        onHover: { // 开启鼠标悬浮的效果(线条跟着鼠标移动)
                            enable: true,
                            mode: 'grab'
                        },
                        resize: true
                    },
                    modes: { // 配置动画效果
                        bubble: {
                            distance: 400,
                            duration: 2,
                            opacity: 0.8,
                            size: 40
                        },
                        push: {
                            quantity: 4
                        },
                        grab: {
                            distance: 200,
                            duration: 0.4
                        },
                        attract: { // 鼠标悬浮时，集中于一点，鼠标移开时释放产生涟漪效果
                            distance: 200,
                            duration: 0.4,
                            factor: 5
                        }
                    }
                },
                particles: {
                    color: {
                        value: '#BA55D3' // 粒子点的颜色
                    },
                    links: {
                        color: '#FFBBFF', // 线条颜色
                        distance: 150,//线条距离
                        enable: true,
                        opacity: 0.4, // 不透明度
                        width: 1.2 // 线条宽度
                    },
                    collisions: {
                        enable: true
                    },
                    move: {
                        attract: { enable: false, rotateX: 600, rotateY: 1200 },
                        bounce: false,
                        direction: 'none',
                        enable: true,
                        out_mode: 'out',
                        random: false,
                        speed: 0.5, // 移动速度
                        straight: false
                    },
                    number: {
                        density: {
                            enable: true,
                            value_area: 800
                        },
                        value: 80//粒子数
                    },
                    opacity: {//粒子透明度
                        value: 0.7
                    },
                    shape: {//粒子样式
                        type: 'star'
                    },
                    size: {//粒子大小
                        random: true,
                        value: 3
                    }
                },
                detectRetina: true
            }
        }
    },
    created() {
        this.GetCaptcha()
    },
    methods: {
        GetCaptcha() {
            this.$ajax({
                method: 'get',
                url: '/login/captcha',
            }).then((res) => {
                this.captchaId = res.captchaId
                this.imgUrl = "http://localhost:8090/v1/api/login/captchaId/" + this.captchaId
            }).catch(function (res) {
                console.log(res.data);
            })
        },
        handleSubmit() {
            this.$ajax.post(
                '/login/user/verifyInfo',
                {
                    username: this.username,
                    password: this.password,
                    captchaId: this.captchaId,
                    captchaSolution: this.user.verifyCode
                },
            ).then((res) => {
                console.log(res)
            }).catch((res) => {
                this.GetCaptcha();
                console.log(res);
            })
        }

    }
};
</script>

<style>
.login {
    height: 100%;
    width: 100%;
    overflow: hidden;
}

.login__particles {
    height: 100%;
    width: 100%;
    background-size: cover;
    background-repeat: no-repeat;
    background-image: url('@/assets/00001.jpg');
    opacity: 0.9;
    position: fixed;
    pointer-events: none;
}

.loginPart {
    position: absolute;
    /*定位方式绝对定位absolute*/
    top: 50%;
    left: 80%;
    /*顶和高同时设置50%实现的是同时水平垂直居中效果*/
    transform: translate(-50%, -50%);
    /*实现块元素百分比下居中*/
    width: 450px;
    padding: 50px;
    background: rgba(255, 204, 255, .3);
    /*背景颜色为黑色，透明度为0.8*/
    box-sizing: border-box;
    /*box-sizing设置盒子模型的解析模式为怪异盒模型，
    将border和padding划归到width范围内*/
    box-shadow: 0px 15px 25px rgba(0, 0, 0, .5);
    /*边框阴影  水平阴影0 垂直阴影15px 模糊25px 颜色黑色透明度0.5*/
    border-radius: 15px;
    /*边框圆角，四个角均为15px*/
}

h2 {
    margin: 0 0 30px;
    padding: 0;
    color: #fff;
    text-align: center;
    /*文字居中*/
}

.btn {
    transform: translate(170px);
    width: 80px;
    height: 40px;
    font-size: 15px;
}
</style>