<template>
    <div class="login">
        <div class="login__particles" />
        <div class="loginPart">
            <h2>登录</h2>
            <el-form ref="user" :rules="rules" :model="user" label-width="100px" class="demo-ruleForm"
                style="transform:translate(-30px);">
                <el-form-item label="账号：" prop="username">
                    <el-input v-model="user.username" placeholder="请输入账号" maxlength="20" />
                </el-form-item>
                <el-form-item label="密码：" prop="password">
                    <el-input v-model="user.password" type="password" placeholder="请输入密码" maxlength="20" show-password />
                </el-form-item>
                <el-form-item label="验证码：" prop="verifyCode">
                    <el-input style="width: 140px; padding-right: 10px;" v-model="user.verifyCode" placeholder="请输入验证码"
                        maxlength="4" />
                    <img id="image" class="verifyCodeImg" :src="imgUrl" @click="GetCaptcha()">
                </el-form-item>
                <el-button class="btn" type="primary" @click="handleSubmit('user')">登录</el-button>
                <div style="text-align: right;transform:translate(0,30px);">
                    <el-link style="padding-right: 60px;" type="warning" href="http://localhost:7070/register"
                        :underline="false">没有账号？去注册</el-link>
                    <el-link type="warning" href="http://localhost:7070/verify_identity"
                        :underline="false">忘记密码？去找回</el-link>
                </div>
            </el-form>
        </div>

    </div>
</template>


<script>
export default {
    data() {
        return {
            user: {
                username: '',
                password: '',
                verifyCode: ''
            },
            captchaId: null,
            imgUrl: null,
            rules: {
                username: [
                    { required: true, message: '用户名不能为空', trigger: 'blur' },
                ],
                password: [
                    { required: true, message: '密码不能为空', trigger: 'blur' },
                    { min: 6, message: '长度至少为6个字符', trigger: 'blur' }
                ],
                verifyCode: [
                    { required: true, message: '验证码不能为空', trigger: 'blur' },
                ],
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
                url: '/user/getCaptcha',
            }).then((res) => {
                this.captchaId = res.captchaId
                this.imgUrl = res.imageUrl
            }).catch(function (res) {
                console.log(res.data);
            })
        },
        handleSubmit(form) {
            this.$refs[form].validate((valid) => {
                if (valid) {
                    this.$ajax.post(
                        '/user/login',
                        {
                            username: this.user.username,
                            password: this.user.password,
                            captcha_id: this.captchaId,
                            verify_value: this.user.verifyCode
                        },
                    ).then((res) => {
                        this.notify("success", "登录成功", res.msg)
                        this.$auth.setUserAuth(res.username, res.token, res.uid)
                        this.$auth.setUserCluster(res.cluster_name)
                        this.$router.push("/home")
                    }).catch((res) => {
                        this.notify("error", "登录失败", res.msg)
                        this.GetCaptcha();
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

.verifyCodeImg {
    background-color: #fff;
    cursor: pointer;
}
</style>