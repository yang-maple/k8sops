<template>
    <div class="common-layout">
        <!-- container 布局-->
        <el-container style="height: 100vh;">
            <!--侧边导航栏-->
            <el-aside class="aside" :width="asideWidth">
                <el-affix class="aside-affix" :z-index="12000">
                    <div class="aside-logo">
                        <el-image class="logo-image" :src="logo" />
                        <span :class="[isCollapse ? 'is-collapse' : 'logo-name']">Kubernetes</span>
                    </div>
                </el-affix>
                <!-- router 定义vue-router模式 菜单栏的index 跟路由规则的path 绑定-->
                <!-- default-active 默认激活菜单栏，根据打开的path 来找到对应的栏 -->
                <el-menu class="aside-menu" router :default-active="$route.path" :collapse="isCollapse"
                    backgroup-color="#131b27" text-color="#bfcbd9" active-text-color="#20a0ff">
                    <!-- routers 就是路由规则的 route -->
                    <div v-for="menu in routers" :key="menu" class="aside-menu-item">
                        <!-- 第一种情况 children 只有一个的路由 -->
                        <el-menu-item class="aside-menu-item" v-if="menu.children && menu.children.length == 1"
                            :index="menu.children[0].path">
                            <!-- 处理图标和菜单栏 -->
                            <el-icon>
                                <component :is="menu.children[0].icon" />
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
                                    <component :is="menu.icon" />
                                </el-icon>
                                <span :class="[isCollapse ? 'is-collapse' : '']">{{ menu.name }}</span>
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
                    </div>
                </el-menu>
            </el-aside>
            <el-container>
                <!-- 头部信息 -->
                <el-header class="header">
                    <el-row :gutter="20">
                        <!-- 折叠按钮 -->
                        <el-col :span="1">
                            <div class="header-collapse" @click="onCollapse">
                                <!---->
                                <el-icon>
                                    <component :is="isCollapse?'Expand':'Fold'" />
                                </el-icon>
                            </div>
                        </el-col>
                        <!-- 面包屑组件 -->
                        <el-col :span="8">
                            <div class="header-breadcrumb">
                                <el-breadcrumb separator="/">
                                    <!-- 最外层 固定 -->
                                    <el-breadcrumb-item :to="{ path: '/' }">工作台</el-breadcrumb-item>
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
                        <el-col :span="12">
                            <div style="text-align: right;">
                                <el-button type="info" @click="dialogcreatens = true">
                                    YAML创建资源
                                </el-button>
                            </div>
                        </el-col>
                        <!-- 用户信息 -->
                        <el-col :span="3">
                            <div class="header-user">
                                <el-dropdown>
                                    <div class="header-dropdown header-username" @mouseover="isOver = true"
                                        @mouseleave="isOver = false">
                                        <el-image class="avator-image" :src="avator" />
                                        <span style="font-size: 20px;">{{ username }}</span>
                                        <el-icon size="16" style="padding-left: 5px; padding-bottom: 5px;">
                                            <component :is="isOver?'ArrowDown':'ArrowLeft'" />
                                        </el-icon>
                                    </div>
                                    <!-- 下拉框内容 -->
                                    <template #dropdown>
                                        <el-dropdown-menu>
                                            <el-dropdown-item>修改密码</el-dropdown-item>
                                            <el-dropdown-item @click="logout">退出</el-dropdown-item>
                                        </el-dropdown-menu>
                                    </template>
                                </el-dropdown>
                            </div>
                        </el-col>
                    </el-row>
                </el-header>
                <el-main>
                    <div v-if="cluster != null && this.$route.path != '/clusterinfo'"><router-view></router-view></div>
                    <div v-if="this.$route.path == '/clusterinfo'"><router-view></router-view></div>
                    <div v-if="cluster == null"><el-empty description="还没有添加任何集群哦~"></el-empty></div>
                </el-main>
                <el-footer class="footer">
                    <el-icon style="width: 2em;top: 3px;font-size: 18px;">
                        <Place />
                    </el-icon>
                    <a>2022 adoo devops</a>
                </el-footer>
            </el-container>
        </el-container>
    </div>
    <el-dialog v-model="dialogcreatens" title="YAML创建资源" center>
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
                <el-button @click="dialogcreatens = false">
                    取消
                </el-button>
            </span>
        </template>
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
    data() {
        return {
            logo: require('@/assets/img/logo.png'),
            avator: require('@/assets/img/user.png'),
            asideWidth: '220px',
            isCollapse: false,
            isOver: false,
            dialogcreatens: false,
            routers: [],
            aceConfig: {
                lang: 'yaml',
                theme: "cloud9_day",
                options: {
                    showPrintMargin: false,
                }
            },
            validationErrors: [],
            activeName: 'yaml',
            content: '',
            file_list: [],
        }
    },
    computed: {
        username() {
            let username = localStorage.getItem('user_name')
            return username ? username : '未知'
        },
        config() {
            let token = localStorage.getItem('user_token')
            return {
                Authorization: token,
            };
        },
        cluster() {
            let cluster = localStorage.getItem('user_cluster')
            return cluster ? cluster : null
        }
    },
    beforeMount() {
        this.routers = useRouter().options.routes
    },
    methods: {
        onCollapse() {
            if (this.isCollapse) {
                this.asideWidth = '220px'
                this.isCollapse = false
            } else {
                this.asideWidth = '63.2px'
                this.isCollapse = true
            }
        },
        logout() {
            this.$auth.delUserAuth()
            this.$router.push('/login')
        },
        mouseover() {
            this.isOver = true
        },
        handleClick(tab) {
            console.log(tab)
        },
        createyaml() {
            try {
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
                    this.dialogcreatens = false
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
            this.dialogcreatens = false
            console.log(response)
        },
        submiterror(error, uploadFile, uploadFiles) {
            this.$message({
                message: JSON.parse(error.message).msg,
                type: 'error'
            });
            console.log(JSON.parse(error.message))
        },
    }
}

</script>

<style>
.aside {
    transition: all .5s;
    background-color: #131b27;
}

/* 滚轴 */
.aside::-webkir-scrollbar {
    display: none;
}

.aside-affix {
    z-index: 12000;
}

.aside-logo {
    color: white;
    background-color: #131b27;
    height: 60px;
}

.logo-image {
    width: 40px;
    height: 40px;
    top: 12px;
    padding-left: 12px;

}

.logo-name {
    font-size: 20px;
    font-weight: bold;
    padding: 10px;
    padding-left: 20px;
}

.aside-menu:not(.el-menu--collapse) {
    background-color: #131b27;
    border-right-width: 0;
}

.is-collapse {
    display: none;
}

/* 概要菜单栏背景色 */
.aside-menu-item {
    background-color: #131b27;
}

.aside-menu-item.is-active {
    background-color: #263448;
}

/* */
.aside-submenu {
    background-color: #131b27;
}

.aside-children {
    background-color: #131b27;

}

.aside-children.is-active {
    background-color: #263448;
}

.header {
    z-index: 1200;
    line-height: 60px;
    font-size: 24px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
}

/*可点击状态*/
.header-collapse {
    cursor: pointer;
}

.header-breadcrumb {
    padding-top: 0.9em;
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

.header-user {
    text-align: right;
}

.header-dropdown {
    line-height: 60px;
}

.header-username {
    cursor: pointer;
}

.footer {
    z-index: 1200;
    color: rgb(187, 184, 184);
    font-size: 14px;
    text-align: center;
    line-height: 30px;
}


.el-menu .el-menu-item:hover {
    background-color: rgba(3, 60, 174, 0.81);
}

.el-menu .el-sub-menu__title:hover {
    background-color: rgba(3, 60, 174, 0.81);
}

.el-menu .el-sub-menu .el-menu-item:hover {
    background-color: rgba(3, 60, 174, 0.81);
}

.el-main {
    width: 100%;
    height: 100%;
    padding: 5px !important;
    padding-top: 10px !important;
}

.el-footer {
    height: 30px !important
}


.el-tabs--border-card>.el-tabs__content {
    padding: 0px !important;
    padding-top: 5px !important;
}

.el-dialog--center {
    .el-dialog__body {
        padding-top: 10px !important;
        padding-bottom: 10px !important;
    }
}

.icon {
    width: 2em;
    height: 2em;
    vertical-align: -0.15em;
    fill: currentColor;
    overflow: hidden;
    color: #ffffffff;
}
</style>