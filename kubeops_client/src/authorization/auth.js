const user_name = "user_name"
const user_token = "user_token"
const uid = "user_id"
const user_cluster = "user_cluster"


class Auth {
    //初始化
    constructor() {
        this.user_name = localStorage.getItem(user_name) ? localStorage.getItem(user_name) : ""
        this.user_token = localStorage.getItem(user_token) ? localStorage.getItem(user_token) : ""
        this.user_token = localStorage.getItem(uid) ? localStorage.getItem(uid) : 0
        this.user_cluster = localStorage.getItem(user_cluster) ? localStorage.getItem(user_cluster) : ""
    }

    //缓存token 用户信息
    setUserAuth(username, token, userid, usercluster) {
        this.user_name = username
        this.user_token = token
        this.uid = userid
        localStorage.setItem(user_name, username)
        localStorage.setItem(user_token, token)
        localStorage.setItem(uid, userid)
    }
    //缓存集群信息
    setUserCluster(usercluster) {
        this.user_cluster = usercluster
        localStorage.setItem(user_cluster, usercluster)
    }
    // 清除token 用户信息
    delUserAuth() {
        this.user_name = null
        this.user_token = null
        this.uid = null
        this.user_cluster = null
        localStorage.removeItem(user_cluster)
        localStorage.removeItem(user_name)
        localStorage.removeItem(user_token)
        localStorage.removeItem(uid)
    }

}

//创建单例。不需要每次进行new
const auth = new Auth()
export default auth