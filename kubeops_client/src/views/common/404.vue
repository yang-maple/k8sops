<template>
    <div class="main-body-div">
        <el-row>
            <!--image-->
            <el-col :span="24">
            <img class="main-body-img" src="../../assets/logo.png">
            </el-col>
            <!--describe-->
            <el-col>
                <div>
                    <p class="status-code">404</p>
                    <p class="status-describe">你所访问的页面不存在.....</p>
                </div>
            </el-col>
        </el-row>
    </div>
</template>

<script>
export default{

}
</script>

<style>
/*image 属性*/
.main-body-img {
    margin-top: 150px;
}
/*整体布局 属性*/
.main-body-div {
    text-align: center;
    height: 100%;
    width: 100%;
}

/*状态码*/
.status-code{
    margin-top: 20px;
    margin-bottom: 10px;
    font-size: 95px;
    font-weight: bold;
    color:rgb(54, 95, 230) 
}
/*描述*/
.status-describe{
    color: rgb(145, 143, 143);
}
</style>