<template>
    <div class="login">
        <div class="login__particles"></div>
        <div class="loginPart">
            <h2>找回密码</h2>
            <el-form ref="user" :rules="rules" :model="user" status-icon label-width="100px" class="demo-ruleForm"
                style="transform:translate(-30px);">
                <el-form-item label="电子邮箱：" prop="email">
                    <el-input v-model="user.email" placeholder="请输入注册账号的邮箱地址" clearable />
                </el-form-item>
                <el-form-item label="验证码：" prop="verifycode">
                    <el-input v-model="user.verifycode" maxlength="6" placeholder="请输入验证码">
                        <template #append>
                            <el-button :disabled="button.disabled" @click="sendemail()">
                                <template #default>
                                    <span v-if="button.disabled == true"><el-icon class="is-loading">
                                            <Loading />
                                        </el-icon></span>
                                    <span>{{ button.text }}</span>
                                </template>
                            </el-button>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item label="新密码：" prop="password">
                    <el-input v-model="user.password" type="password" placeholder="请输入新密码" maxlength="20" show-password
                        clearable />
                </el-form-item>
                <el-form-item label="确认密码：" prop="checkpassword">
                    <el-input v-model="user.checkpassword" type="password" placeholder="请再次输入密码" maxlength="20"
                        show-password clearable />
                </el-form-item>
                <el-button class="btn" type="primary" @click="handleSubmit('user')">修改密码</el-button>
            </el-form>
        </div>

    </div>
</template>


<script>
import { ElLoading } from 'element-plus'
export default {
    data() {
        var validatePass2 = (rule, value, callback) => {
            if (value === '') {
                callback(new Error('请再次输入密码'));
            } else if (value !== this.user.password) {
                callback(new Error('两次输入密码不一致!'));
            } else {
                callback();
            }
        };
        return {
            user: {
                email: '',
                verifycode: '',
                password: '',
                checkpassword: '',
            },
            button: {
                text: '获取验证码',
                disabled: false,
                duration: 90,
                timer: null
            },
            rules: {
                email: [
                    { required: true, message: '邮箱不能为空', trigger: 'blur' },
                    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
                ],
                verifycode: [
                    { required: true, message: '验证码不能为空', trigger: 'blur' },
                ],
                password: [
                    { required: true, message: '密码不能为空', trigger: 'blur' },
                    { min: 6, message: '长度至少为6个字符', trigger: 'blur' }
                ],
                checkpassword: [
                    { required: true, message: '密码不能为空', trigger: 'blur' },
                    { validator: validatePass2, }
                ],
            }

        }
    },
    unmounted() {
        clearInterval(this.button.timer);
    },
    methods: {
        sendemail() {
            //清除计时器
            this.button.timer && clearInterval(this.button.timer)
            if (this.user.email != "") {
                this.$ajax.post(
                    '/user/findPassword/email',
                    {
                        email: this.user.email,
                    },
                ).then((res) => {



                    this.notify("success", "发送成功", res.msg)
                    this.button.timer = setInterval(() => {
                        // 倒计时期间按钮不能点击
                        this.button.disabled = true
                        const tmp = this.button.duration--
                        this.button.text = `剩余${tmp}秒`
                        if (tmp <= 0) {
                            // 清除掉定时器
                            clearInterval(this.button.timer)
                            this.button.duration = 90
                            this.button.text = '请重新获取'
                            // 设置按钮可以单击
                            this.button.disabled = false
                        }
                    }, 1000)

                }).catch((res) => {
                    this.notify("error", "发送失败", res.msg)
                    console.log(res);
                })
                return
            }
            this.notify("error", "发送失败", "邮箱地址不正确或为空")
        },
        handleSubmit(form) {
            this.$refs[form].validate((valid) => {
                if (valid) {
                    this.$ajax.post(
                        '/user/findPassword',
                        {
                            password: this.user.password,
                            email: this.user.email,
                            verify_code: this.user.verifycode,
                        },
                    ).then((res) => {
                        this.notify("success", "重置成功", res.msg)
                        this.$router.push("/login")
                        console.log(res)
                    }).catch((res) => {
                        this.notify("error", "重置失败", res.msg)
                        console.log(res);
                    })
                } else {
                    return false;
                }
            });
        },
        notify(status, info, msg) {
            this.$notify({
                title: info,
                message: msg,
                type: status
            });
        },
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
    background-image: url('@/assets/img/00001.jpg');
    opacity: 0.9;
    position: fixed;
    pointer-events: none;
}

.loginPart {
    position: fixed;
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

.verifyCodeImg {
    background-color: #fff;
    cursor: pointer;
}
</style>