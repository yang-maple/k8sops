const user_name = "user_name"
const user_token = "user_token"


class Auth {
    //初始化
    constructor() {
        this.user_name = localStorage.getItem(user_name) ? localStorage.getItem(user_name) : ""
        this.user_token = localStorage.getItem(user_token) ? localStorage.getItem(user_token) : ""
    }

    //缓存token 用户信息
    setUserAuth(username, token) {
        this.user_name = username
        this.user_token = token
        localStorage.setItem(user_name, username)
        localStorage.setItem(user_token, token)
    }

    // 清除token 用户信息
    delUserAuth() {
        this.user_name = null
        this.user_token = null
        localStorage.removeItem(user_name)
        localStorage.removeItem(user_token)
    }

}

//创建单例。不需要每次进行new
const auth = new Auth()
export default auth