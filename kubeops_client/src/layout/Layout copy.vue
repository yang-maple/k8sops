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
                                    <el-icon>
                                        <component :is="child.icon" />
                                    </el-icon>
                                    {{ child.name }}
                                </template>
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
                        <el-col :span="10">
                            <div class="header-breadcrumb">
                                <el-breadcrumb separator="/">
                                    <!-- 最外层 固定 -->
                                    <el-breadcrumb-item :to="{ path: '/' }">工作台</el-breadcrumb-item>
                                    <!-- 循环内层 name -->
                                    <template v-for="(matched, m) in this.$route.matched" :key="m">
                                        <el-breadcrumb-item :to="{ path: this.$route.path }" v-if="matched.name != undefined">
                                            {{ matched.name }}
                                        </el-breadcrumb-item>
                                    </template>
                                </el-breadcrumb>
                            </div>
                        </el-col>
                        <!-- 用户信息 -->
                        <el-col :span="13">
                            <div class="header-user">
                                <el-dropdown>
                                    <div class="header-dropdown">
                                        <el-image class="avator-image" :src="avator" />
                                        <span class="header-username"> {{ username }}</span>
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
                <el-main class="main">
                    <router-view></router-view>
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
</template>

<script>
import { useRouter } from 'vue-router';

export default {
    data() {
        return {
            logo: require('@/assets/img/logo.png'),
            avator: require('@/assets/img/user.png'),
            asideWidth: '220px',
            isCollapse: false,
            routers: [],
        }
    },
    computed: {
        username() {
            let username = localStorage.getItem('username')
            return username ? username : '未知'
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
            localStorage.removeItem("username")
            localStorage.removeItem("token")
            this.$router.push('/login')
        }
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

.aside-menu {
    background-color: #131b27;
    border-right-width: 0;
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

.aside-menu-item:not(.el-menu--collapse) {
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
    padding-left: 25px !important;
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

.main {
    padding: 10px;
}

.footer {
    z-index: 1200;
    color: rgb(187, 184, 184);
    font-size: 14px;
    text-align: center;
    line-height: 60px;
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


.el-sub-menu .el-sub-menu__icon-arrow:not(.el-menu--collapse) {
    display: none;
}</style>