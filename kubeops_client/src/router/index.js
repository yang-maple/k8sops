import { createRouter, createWebHistory } from "vue-router";
//导入进度条组件
import NProgress from "nprogress";
import 'nprogress/nprogress.css'
import layout from '@/layout/Layout.vue'

// 定义路由规划
const routes = [
    {
        path: '/home',
        component: layout,
        meta: { title: '概要', requireAuth: true },
        children: [
            {
                path: '/home',
                icon: '#icon-s-promotion',
                component: () => import('@/views/homepage/Homepage.vue'),
                name: '概要',
            }
        ]
    },
    {
        path: '/clusterinfo',
        component: layout,
        meta: { title: '集群信息', requireAuth: true },
        children: [
            {
                path: '/clusterinfo',
                icon: '#icon-jiqunguanli2',
                component: () => import('@/views/clustermanager/Cluster.vue'),
                name: '集群信息',
            }
        ]
    },
    {
        path: '/cluster',
        component: layout,
        name: '集群资源',
        icon: "Menu",
        meta: { title: '集群资源', requireAuth: true },
        children: [
            {
                path: '/node',
                icon: '#icon-jiqunguanli1',
                component: () => import('@/layout/index.vue'),
                name: '节点',
                meta: { title: '节点', requireAuth: true },
                children: [
                    {
                        path: '/node',
                        component: () => import('@/views/nodes/Node.vue'),
                    },
                    {
                        path: '/node/node_detail',
                        name: '节点详情',
                        icon: 'LocationInformation',
                        component: () => import('@/views/nodes/Nodedetail.vue'),
                        meta: { title: '节点详情', requireAuth: true },
                    }
                ]
            },
            {
                path: '/namespaces',
                icon: 'LocationInformation',
                component: () => import('@/layout/index.vue'),
                name: '命名空间',
                meta: { title: '命名空间', requireAuth: true },
                children: [
                    {
                        path: '/namespaces',
                        component: () => import('@/views/namespace/Namespace.vue'),
                    },
                    {
                        path: '/namespaces/namespaces_detail',
                        name: '命名空间详情',
                        icon: 'LocationInformation',
                        component: () => import('@/views/namespace/Namespacedetail.vue'),
                        meta: { title: '命名空间详情', requireAuth: true },
                    }
                ]
            },
        ]
    },
    {
        path: '/workload',
        component: layout,
        name: '工作负载',
        icon: "Menu",
        meta: { title: '工作负载', requireAuth: true },
        children: [
            {
                path: '/deployment',
                name: "无状态",
                component: () => import('@/views/deployment/Deployment.vue'),
                icon: '#icon-gongzuofuzai1',
                meta: { title: '无状态', requireAuth: true },
            },
            {
                path: '/statefulset',
                name: "有状态",
                component: () => import('@/views/statefulset/Statefulset.vue'),
                meta: { title: '有状态', requireAuth: true },
            },
            {
                path: '/daemonset',
                name: "守护进程",
                component: () => import('@/views/daemonset/Daemonset.vue'),
                meta: { title: '守护进程', requireAuth: true },
            },
            {
                path: '/pod',
                name: '容器组',
                component: () => import('@/layout/index.vue'),
                meta: { title: '容器组', requireAuth: true },
                children: [
                    {
                        path: '/pod',
                        component: () => import('@/views/pod/Pod.vue'),
                    },
                    {
                        path: '/pod/shell',
                        name: 'shell',
                        component: () => import('@/views/pod/Xtrem.vue'),
                        meta: {
                            title: 'shell',
                            requireAuth: false
                        }
                    },
                    {
                        path: '/pod/log',
                        name: 'logs',
                        component: () => import('@/views/pod/Xlogs.vue'),
                        meta: {
                            title: 'logs',
                            requireAuth: false
                        }
                    }
                ],
            },
        ]
    },
    {
        path: '/sever',
        component: layout,
        name: '网络资源',
        meta: { title: '网络资源', requireAuth: true },
        children: [
            {
                path: '/services',
                icon: '#icon-gongzuofuzai',
                component: () => import('@/layout/index.vue'),
                name: '服务',
                meta: { title: '服务', requireAuth: true },
                children: [
                    {
                        path: '/services',
                        component: () => import('@/views/service/Service.vue'),
                    },
                    {
                        path: '/services/services_detail',
                        name: '服务详情',
                        component: () => import('@/views/service/Servicedetail.vue'),
                        meta: { title: '服务详情', requireAuth: true },
                    }
                ]
            },
            {
                path: '/ingress',
                icon: 'LocationInformation',
                component: () => import('@/views/service/Ingress.vue'),
                name: '应用路由',
                meta: { title: '应用路由', requireAuth: true },
            }, {
                path: '/storageclass',
                icon: 'LocationInformation',
                component: () => import('@/views/service/Class.vue'),
                name: '应用路由类',
                meta: { title: '应用路由类', requireAuth: true },
            }
        ]
    },
    {
        path: '/configure',
        component: layout,
        name: '配置与存储',
        icon: "Menu",
        meta: { title: '配置与存储', requireAuth: true },
        children: [
            {
                path: '/configure/persistenvolumeclaim',
                icon: '#icon-yidongyunkongzhitaiicon06',
                component: () => import('@/views/storage/Persistenvolumeclaim.vue'),
                name: '持久卷声明',
                meta: { title: '持久卷声明', requireAuth: true },
            },
            {
                path: '/persistentvolume',
                icon: 'LocationInformation',
                component: () => import('@/layout/index.vue'),
                name: '持久卷',
                meta: { title: '持久卷', requireAuth: true },
                children: [
                    {
                        path: '/persistentvolume',
                        component: () => import('@/views/PersistentVolume/PersistentVolume.vue'),
                    },
                    {
                        path: '/persistentvolume/persistentvolume_detail',
                        name: '持久卷详情',
                        icon: 'LocationInformation',
                        component: () => import('@/views/PersistentVolume/PersistentVolumedetail.vue'),
                        meta: { title: '持久卷详情', requireAuth: true },
                    }
                ]
            },
            {
                path: '/configure/storageclass',
                icon: 'LocationInformation',
                component: () => import('@/views/storage/Class.vue'),
                name: '存储类',
                meta: { title: '存储类', requireAuth: true },
            },
            {
                path: '/configure/configmap',
                icon: 'LocationInformation',
                component: () => import('@/views/storage/Configmap.vue'),
                name: '配置字典',
                meta: { title: '配置字典', requireAuth: true },
            },
            {
                path: '/configure/secret',
                icon: 'LocationInformation',
                component: () => import('@/views/storage/Secret.vue'),
                name: '保密字典',
                meta: { title: '保密字典', requireAuth: true },
            },
        ]
    },
    {
        path: '/404',
        name: '404',
        component: () => import('@/views/common/404.vue'),
        meta: {
            title: "404",
            requireAuth: false
        }
    },
    {
        path: '/:pathMatch(.*)',
        redirect: '/404',
    },
    {
        path: '/',
        redirect: '/home',
    },
    {
        path: '/403',
        name: '403',
        component: () => import('@/views/common/403.vue'),
        meta: {
            title: "403",
            requireAuth: false
        }
    },
    {
        path: '/login',
        name: '登陆',
        component: () => import('@/views/user/Login.vue'),
        meta: {
            title: "登陆",
            requireAuth: false
        }
    },
    {
        path: '/register',
        name: '注册',
        component: () => import('@/views/user/Register.vue'),
        meta: {
            title: "注册",
            requireAuth: false
        }
    },
    {
        path: '/verify_identity',
        name: '重置密码',
        component: () => import('@/views/user/Verify_identity.vue'),
        meta: {
            title: "重置密码",
            requireAuth: false
        }
    },
]

// 创建路由实例
const router = createRouter({
    history: createWebHistory(),
    routes
})


// 进度条配置
//**要递增进度条，只需使用.inc()
NProgress.inc(0.2) //这将获取当前状态值并添加0.2直到状态为0.994
NProgress.configure({
    easing: 'ease',// 动画方式
    speed: 600,      // 递增进度条的速度
    showSpinner: false //是否显示加载ico
})
//路由拦截
router.beforeEach((to, from, next) => {
    //启动进度条
    NProgress.start()
    //设置头部
    if (to.meta.title) {
        document.title = to.meta.title
    } else {
        document.title = "Kubernetes"
    }
    //to 满足条件放过 login register resetpasswd
    if (to.path == "/login" || to.path == "/register" || to.path == "/verify_identity") {
        next()
    } else {
        //不满足条件 做token 校验
        const user_token = localStorage.getItem("user_token")
        if (user_token == "" || user_token == null || user_token == '') {
            //重定向到登录页面
            next("/login")
        } else {
            //放行
            next()
        }
    }
})

router.afterEach(() => {
    //关闭进度条
    NProgress.done()
})


//导出路由实例，被其他js引用
export default router