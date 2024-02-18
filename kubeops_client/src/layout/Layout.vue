<template>
    <div class="common-layout">
        <!-- container 布局-->
        <el-container>
            <!--侧边导航栏-->
            <el-aside :width="asideWidth">
                <!--固定控制标头组件-->
                <div class="aside-logo" v-bind:translate="false">
                    <el-image class="logo-image" :src="logo" />
                    <span :class="[isCollapse ? 'is-collapse' : 'logo-name']">Kubernetes</span>
                </div>
                <!--菜单导航栏-->
                <el-menu class="el-menu-vertical-demo" router :default-active="$route.path" :collapse="isCollapse"
                    text-color="#bfcbd9" active-text-color="#20a0ff" background-color="#131b27">
                    <!-- routers 就是路由规则的 route -->
                    <!-- 第一种情况 children 只有一个的路由 -->
                    <!-- 第一种情况 children 只有一个的路由 -->
                    <template v-for="menu in routers" :key="menu">
                        <el-menu-item v-if="menu.children && menu.children.length == 1" :index="menu.children[0].path">
                            <!-- 处理图标和菜单栏 -->
                            <el-icon>
                                <svg class="icon" aria-hidden="true">
                                    <use :xlink:href="menu.children[0].icon"></use>
                                </svg>
                            </el-icon>
                            <template #title>
                                {{ menu.children[0].name }}
                            </template>
                        </el-menu-item>
                        <!-- 第二种情况 children 大于一个的路由 -->
                        <el-sub-menu class="aside-submenu" v-if="menu.children && menu.children.length > 1"
                            :index="menu.path">
                            <!-- 处理父菜单栏 -->
                            <template #title>
                                <el-icon>
                                    <svg class="icon-layui" aria-hidden="true">
                                        <use :xlink:href="menu.children[0].icon"></use>
                                    </svg>
                                </el-icon>
                                <span>{{ menu.name }}</span>
                            </template>
                            <!-- 处理子菜单栏 -->
                            <el-menu-item class="aside-children" v-for="child in menu.children" :key="child"
                                :index="child.path">
                                <template #title>
                                    <span>{{ child.name }}</span>
                                </template>
                                <!-- 处理三级子菜单栏 -->
                                <!-- <el-menu-item v-for="third_child in child.children" :key="third_child"
                                    :index="third_child.path" class="aside-children">
                                </el-menu-item> -->
                            </el-menu-item>
                        </el-sub-menu>
                    </template>
                </el-menu>
            </el-aside>
            <el-container>
                <!--头部导航栏-->
                <el-header>
                    <el-row :gutter="20">
                        <!-- 折叠按钮 -->
                        <el-col :span="1">
                            <div class="header-collapse" @click="collapse()">
                                <el-icon>
                                    <component :is="isCollapse?'Expand':'Fold'" />
                                </el-icon>
                            </div>
                        </el-col>
                        <!--面包屑组件-->
                        <el-col :span="8">
                            <div class="header-breadcrumb">
                                <el-breadcrumb separator="/">
                                    <!-- 最外层 固定 -->
                                    <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
                                    <!-- 循环内层 name -->
                                    <template v-for="(matched, m) in this.$route.matched" :key="m">
                                        <el-breadcrumb-item v-if="matched.name != undefined">
                                            {{ matched.name }}
                                        </el-breadcrumb-item>
                                    </template>
                                </el-breadcrumb>
                            </div>
                        </el-col>
                        <!--通过yaml 创建资源-->
                        <el-col :span="10">
                            <div style="text-align: right;">
                                <el-button type="info" @click="checkcluster()">
                                    YAML创建资源
                                </el-button>
                            </div>
                        </el-col>
                        <!--用户信息-->
                        <el-col :span="5">
                            <div style="text-align: right;">
                                <el-dropdown>
                                    <div class="header-username" @mouseover="isOver = true" @mouseleave="isOver = false"
                                        style="text-align: right;">
                                        <el-image class="avator-image" :src="avator" />
                                        <span style="font-size: 20px;">{{ username }}</span>
                                        <el-icon size="16" style="padding-left: 5px; padding-bottom: 5px;">
                                            <component :is="isOver?'ArrowDown':'ArrowLeft'" />
                                        </el-icon>
                                    </div>
                                    <!-- 下拉框内容 -->
                                    <template #dropdown>
                                        <el-dropdown-menu>
                                            <el-dropdown-item @click="dialogpassword = true">修改密码</el-dropdown-item>
                                            <el-dropdown-item @click="logout">退出</el-dropdown-item>
                                        </el-dropdown-menu>
                                    </template>
                                </el-dropdown>
                            </div>
                        </el-col>
                    </el-row>
                </el-header>
                <!--Main导航栏-->
                <el-main>
                    <div v-if="cluster != null && this.$route.path != '/clusterinfo'"><router-view></router-view></div>
                    <div v-if="this.$route.path == '/clusterinfo'"><router-view></router-view></div>
                    <div v-if="cluster == null && this.$route.path != '/clusterinfo'"><el-empty description="还没有添加任何集群哦~"
                            :image="empty" :image-size="400">
                            <el-button type="primary" link @click="addCluster">去添加集群</el-button>
                        </el-empty>
                    </div>
                </el-main>
            </el-container>
        </el-container>
    </div>
    <el-dialog v-model="dialogyaml" title="YAML创建资源" center>
        <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
            <el-tab-pane label="YAML" name="yaml"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                    style="min-height: 400px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
            <el-tab-pane label="YAML文件" name="yamlfile">
                <el-upload ref="upload" class="upload-demo" drag action="http://127.0.0.1:8090/v1/api/upload/uploadFile"
                    multiple :auto-upload="false" :name="yamlfile" :file-list="file_list" :headers="config"
                    :on-success="sumbitsuccess" :on-error="submiterror">
                    <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                    <div class="el-upload__text">
                        Drop file here or <em>click to upload</em>
                    </div>
                </el-upload>
            </el-tab-pane>
        </el-tabs>
        <template #footer>
            <el-text type="info" size="small" v-if="activeName == 'yaml'">仅支持单资源创建，多资源创建请使用文件上传</el-text><br>
            <span class="dialog-footer">
                <el-button type="primary" v-if="activeName == 'yaml'" @click="createyaml()">创建</el-button>
                <el-button type="primary" v-if="activeName == 'yamlfile'" @click="submitUpload">上传并创建</el-button>
                <el-button @click="handleCancel">
                    取消
                </el-button>
            </span>
        </template>
    </el-dialog>
    <el-dialog v-model="dialogpassword" title="重置密码" center style="width: 400px;">
        <el-form :model="passwordForm" :rules="rules" ref="passwordForm" label-width="100px">
            <el-form-item label="旧密码" prop="oldpasswd">
                <el-input v-model="passwordForm.oldpasswd"></el-input>
            </el-form-item>
            <el-form-item label="新密码" prop="newpasswd">
                <el-input type="password" v-model="passwordForm.newpasswd" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="确认新密码" prop="checknewpasswd">
                <el-input type="password" v-model="passwordForm.checknewpasswd" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="submitForm('passwordForm')">提交</el-button>
                <el-button @click="resetForm('passwordForm')">重置</el-button>
            </el-form-item>
        </el-form>
    </el-dialog>
</template>

<script>
import { useRouter } from 'vue-router';
import { VAceEditor } from 'vue3-ace-editor';
import { ref } from 'vue';
import '../ace.config.js';
import jsyaml from 'js-yaml'
export default {
    components: {
        VAceEditor
    },
    setup() {
        const yamlfile = ref('yamlfile')
        return { yamlfile }
    },
    computed: {
        username() {
            let username = localStorage.getItem('user_name')
            return username ? username : '未知'
        },
        config() {
            let token = localStorage.getItem('user_token')
            let uuid = localStorage.getItem('user_id')
            return {
                Authorization: token,
                Uuid: uuid
            };
        },
        cluster() {
            let cluster = localStorage.getItem('user_cluster')
            return cluster ? cluster : null
        }
    },
    data() {
        var checkpass = (rule, value, callback) => {
            if (value === '') {
                callback(new Error('请再次输入密码'));
            } else if (value !== this.passwordForm.newpasswd) {
                callback(new Error('两次输入密码不一致!'));
            } else {
                callback();
            }
        };
        return {
            logo: require('@/assets/img/logo.png'),
            avator: require('@/assets/img/user.png'),
            empty: require('@/assets/img/empty.png'),
            asideWidth: "220px",
            isCollapse: false,
            isOver: false,
            routers: [],
            dialogyaml: false,
            dialogpassword: false,
            aceConfig: {
                lang: 'yaml',
                theme: "cloud9_day",
                options: {
                    showPrintMargin: false,
                }
            },
            content: '',
            file_list: [],
            activeName: 'yaml',
            passwordForm: {
                oldpasswd: '',
                newpasswd: '',
                checknewpasswd: '',
            },
            rules: {
                oldpasswd: [
                    { required: true, message: '请输入旧密码', trigger: 'blur' }
                ],
                newpasswd: [
                    { required: true, message: '请输入新密码', trigger: 'blur' },
                    { min: 6, max: 12, message: '长度在 6 到 12 个字符', trigger: 'blur' }
                ],
                checknewpasswd: [
                    { required: true, validator: checkpass, trigger: 'blur' }
                ]
            }
        }
    },
    methods: {
        collapse() {
            if (this.isCollapse == false) {
                this.isCollapse = true
                this.asideWidth = "60px"
            } else {
                this.isCollapse = false
                this.asideWidth = "220px"
            }
        },
        createyaml() {
            try {
                if (this.content == '') {
                    this.$message({
                        message: '内容不能为空',
                        type: 'error'
                    });
                    return
                }
                let data = jsyaml.load(this.content);
                this.$ajax.post(
                    '/upload/uploadYaml',
                    {
                        yamlContent: JSON.stringify(jsyaml.load(this.content), null, 2),
                    },
                ).then((res) => {
                    this.$message({
                        message: res.msg,
                        type: 'success'
                    });
                    this.dialogyaml = false
                    this.content = ''
                    console.log(res)
                }).catch((res) => {
                    this.$message({
                        message: res.msg,
                        type: 'error'
                    });
                    console.log(res);
                })
            } catch (error) {
                this.$message({
                    message: 'yaml 格式无效',
                    type: 'error'
                });
                console.error(error);
            }
        },
        submitUpload() {
            this.$refs.upload.submit();
        },
        sumbitsuccess(response, uploadFile, uploadFiles) {
            this.$message({
                message: response.msg,
                type: 'success'
            });
            this.dialogyaml = false
            console.log(response)
        },
        submiterror(error, uploadFile, uploadFiles) {
            this.$message({
                message: JSON.parse(error.message).msg,
                type: 'error'
            });
            console.log(JSON.parse(error.message))
        },
        handleCancel() {
            this.dialogyaml = false
            this.content = ''
        },
        handleClick(tab) {
            console.log(tab)
        },
        logout() {
            this.$auth.delUserAuth()
            this.$router.push('/login')
        },
        addCluster() {
            this.$router.push('/clusterinfo')
        },
        checkcluster() {
            let clusters = localStorage.getItem('user_cluster')
            if (clusters == "") {
                this.dialogyaml = false;
                this.$message({
                    message: "请先添加集群",
                    type: 'error'
                });
            } else {
                this.dialogyaml = true;

            }

        },
        submitForm(formName) {
            console.log(typeof this.passwordForm.oldpasswd)
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    console.log(valid);
                    this.$ajax.post(
                        '/user/resetPassword',
                        {
                            id: parseInt(localStorage.getItem('user_id')),
                            old_passwd: this.passwordForm.oldpasswd,
                            new_passwd: this.passwordForm.newpasswd,
                        },
                    ).then((res) => {
                        this.$message({
                            message: res.msg,
                            type: 'success'
                        });
                        this.dialogpassword = false
                        console.log(res)
                    }).catch((res) => {
                        this.$message({
                            message: res.msg,
                            type: 'error'
                        });
                    })
                } else {
                    this.$message({
                        message: "请填写完整的数据",
                        type: 'error'
                    });
                    return false;
                }
            });
        },
        resetForm(formName) {
            this.$refs[formName].resetFields();
        }

    },
    beforeMount() {
        this.routers = useRouter().options.routes
    },
}
</script>

<style>
.common-layout {
    height: 100vh;
}

.el-container {
    height: 100%;
}



/* 
    边框样式，隐藏X轴滚轮
*/
.el-aside {
    transition: all .5s;
    background-color: #131b27;
    overflow-x: hidden !important;
}

/* 
    边框样式，隐藏Y轴滚轮
*/
.el-aside::-webkit-scrollbar {
    display: none;
}

/* 
   logo样式
*/
.aside-logo {
    height: 60px;
}

/* 
   logo 图片样式
*/
.logo-image {
    width: 40px;
    height: 40px;
    top: 12px;
    padding-left: 12px;

}

/* 
   logo 名称样式
*/
.logo-name {
    color: white;
    font-size: 20px;
    font-weight: bold;
    padding: 10px;
    transition: none !important;
}

/* 
    隐藏logo 名称
*/
.is-collapse {
    display: none;
}

/* 
   取消菜单栏的边框
*/
.el-menu {
    border-right: solid 0px !important;
}

/* 
   菜单栏折叠样式
*/
.el-menu-vertical-demo:not(.el-menu--collapse) {
    width: 220px;
    height: 100%;
}

/* 
   鼠标悬浮菜单栏样式
*/
.el-menu .el-menu-item:hover {
    background-color: rgba(3, 60, 174, 0.81);
}

/* 
   鼠标选中菜单栏样式
*/
.el-menu .el-menu-item.is-active {
    background-color: #263448;
}

/* 
   main 页面样式
*/
.el-main {
    width: 100%;
    height: 100%;
    padding: 5px !important;
    padding-top: 10px !important;
}

/* 
   固定头部样式
*/
.el-header {
    z-index: 1200;
    line-height: 60px;
    font-size: 24px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
}

/* 
  点击图标折叠
*/
.header-collapse {
    .el-icon {
        cursor: pointer;
    }

}

/* 
  面包屑样式
*/
.header-breadcrumb {
    padding-top: 0.85em;
}

/* 
  el-dialog 中 yaml内边距样式
*/
.el-tabs--border-card>.el-tabs__content {
    padding: 0px !important;
    padding-top: 5px !important;
}


.el-dialog--center {
    .el-dialog__body {
        padding-top: 0px !important;
        padding-bottom: 10px !important;
    }
}

.el-dialog__footer {
    padding: 0px !important;
    padding-bottom: 10px !important;
}




.header-username {
    line-height: 40px;
    cursor: pointer;
}

:focus-visible {
    outline: -webkit-focus-ring-color auto 0px;
}

/*用户图片信息*/
.avator-image {
    top: 12px;
    width: 40px;
    height: 40px;
    /*直角变圆角*/
    border-radius: 50%;
    margin-right: 8px;
}

.icon-layui {
    width: 3em;
    height: 3em;
    vertical-align: -0.15em;
    fill: currentColor;
    overflow: hidden;

}
</style>