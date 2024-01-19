<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-cluster" aria-hidden="true">
                            <use xlink:href="#icon-jiqun"></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">集群config配置</span>
                        <br>
                        <span style="font-size: 12px;color: #79879c!important">集在 Kubernetes 集群中，config 资源用于存储集群的配置信息，包括 API
                            Server 的地址、认证信息、授权信息等。通过 config 资源，可以方便地连接和管理 Kubernetes 集群。
                        </span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>
    <el-row class="cluster-header">
        <el-col :span="6">
            <el-statistic title="集群总数" :value="clusteritem.length" />
        </el-col>
        <el-col :span="6">
            <el-statistic :value="k8stotal.length">
                <template #title>
                    <div style="display: inline-flex; align-items: center">
                        私有云集群
                        <el-icon style="margin-left: 4px" :size="12">
                            <Male />
                        </el-icon>
                    </div>
                </template>
            </el-statistic>
        </el-col>
        <el-col :span="6">
            <el-statistic title="阿里云集群" :value="acktotal.length" />
        </el-col>
        <el-col :span="6">
            <el-statistic title="Feedback number" :value="562">
                <template #suffix>
                    <el-icon style="vertical-align: -0.125em">
                        <ChatLineRound />
                    </el-icon>
                </template>
            </el-statistic>
        </el-col>
    </el-row>
    <el-row :gutter="10">
        <el-col v-for="(value, index) in clusteritem" :key="value" :span="6" class="cluster-card-content">
            <el-card :body-style="{ padding: '0px' }">
                <img :src="k8s" v-if="value.type == 'k8s'" class="image">
                <img :src="ack" v-if="value.type == 'ack'" class="image">
                <div style="padding: 14px">
                    <span>{{ value.cluster_name }}</span>
                    <div class="bottom">
                        <time class="time">{{ value.create_time }}</time>
                        <div class="status">
                            <el-tag :type="value.status == false ? 'danger' : 'success'"
                                v-if="value.status == false">未应用</el-tag>
                            <el-tag :type="value.status == true ? 'success' : 'danger'"
                                v-if="value.status == true">已应用</el-tag>
                        </div>
                    </div>
                    <el-button text class="button" v-if="value.status == false"
                        @click="change(value.cluster_name)">应用</el-button>
                    <el-button text class="button" disabled v-if="value.status == true">应用中</el-button>
                    <el-button text class="button" type="danger" @click="remove(value.cluster_name)"
                        v-if="value.status == false">删除</el-button>

                </div>
            </el-card>
        </el-col>
        <el-col :span="6" style="padding-top: 10px;">
            <div class="config-add" @click="dialogcreatens = true">
                <img v-if="imageUrl" :src="imageUrl" class="avatar" />
                <el-icon v-else class="avatar-uploader-icon">
                    <Plus />
                </el-icon>
            </div>
        </el-col>
    </el-row>
    <el-dialog v-model="dialogcreatens" title="新增集群" center>
        <el-form :model="ruleForm" label-width="25%">
            <el-form-item label="集群名称" prop="cluster_name">
                <el-input v-model="ruleForm.cluster_name" autocomplete="off" style="width: 80%;" />
            </el-form-item>
            <el-form-item label="集群类型" prop="cluster_type">
                <el-select v-model="ruleForm.cluster_type" placeholder="请选择集群类型">
                    <el-option label="私有云K8S" value="k8s" />
                    <el-option label="阿里云ACK" value="ack" />
                </el-select>
            </el-form-item>
            <el-form-item label="集群config" prop="cluster_file">
                <el-upload ref="upload" action="http://127.0.0.1:8090/v1/api/cluster/create" :auto-upload="false"
                    :headers="config" :name="file" :file-list="upload" :limit="1" style="width:300px" :data="ruleForm"
                    :on-success="sumbitsuccess" :on-error="submiterror" :on-change="onChange">
                    <template #trigger>
                        <el-button type="primary">select file<el-icon class="el-icon--right">
                                <Upload />
                            </el-icon></el-button>
                    </template>
                    <template #tip>
                        <div class="el-upload__tip">
                            集群信息一般位于 /root/.kube/config
                        </div>
                    </template>
                </el-upload>
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="submitUpload()"
                    style="margin-right:50px;padding-left:20px;padding-right:20px">新增</el-button>
                <el-button @click="dialogcreatens = false" style="margin-right:50px;padding-left:20px;padding-right:20px">
                    取消
                </el-button>
            </span>
        </template>
    </el-dialog>
</template>

<script>
import { ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus'
import { ChatLineRound, Male } from '@element-plus/icons-vue'
export default {
    inject: ['reload'],
    setup() {
        const file = ref('config_file')
        const imageUrl = ref('')
        return { file, imageUrl }
    },
    data() {
        return {
            k8s: require('@/assets/img/k8s.jpg'),
            ack: require('@/assets/img/ack1.png'),
            currentDate: new Date(),
            upload: [],
            uploadfile: [],
            disabled: false,
            dialogcreatens: false,
            ruleForm: {
                cluster_name: '',
                cluster_type: '',
            },
            name: '',
            clusteritem: [],
            k8stotal: [],
            acktotal: [],
        };
    },
    computed: {
        config() {
            return {
                Authorization: localStorage.getItem('user_token'),
                Uuid: localStorage.getItem('user_id')
            };
        },
    },
    created() {
        this.clusterlist()
    },
    methods: {
        submitUpload() {
            //判断是否上传文件
            if (!this.ruleForm.cluster_name || !this.ruleForm.cluster_type) {
                this.$message({
                    message: '请填写完整信息',
                    type: 'error'
                });
                return;
            }
            if (this.uploadfile.length < 1) {
                this.$message({
                    message: '请上传文件',
                    type: 'error'
                });
                return;
            }
            this.$refs.upload.submit();
        },
        sumbitsuccess(response, uploadFile, uploadFiles) {
            this.notify("success", "新增集群成功", response.msg)
            this.$refs.upload.clearFiles()
            this.ruleForm.cluster_name = null
            this.ruleForm.cluster_type = null
            this.dialogcreatens = false
            this.reload()
            console.log(response)
        },
        submiterror(error, uploadFile, uploadFiles) {
            this.$message({
                message: JSON.parse(error.message).msg,
                type: 'error'
            });
            console.log(JSON.parse(error.message))
        },
        clusterlist() {
            this.$ajax({
                method: 'get',
                url: '/cluster/list',
                params: {
                    name: this.name
                }
            }).then((res) => {
                this.clusteritem = res.data.item
                this.Counter()
            }).catch((res) => {
                console.log(res);
            })
        },
        change(clustername) {
            this.$ajax.post(
                '/cluster/change',
                {
                    cluster_name: clustername
                },
            ).then((res) => {
                this.notify("success", "成功", res.msg)
                this.$auth.setUserCluster(clustername)
                this.reload()
            }).catch((res) => {
                this.notify("error", "失败", res.msg)
                console.log(res);
            })
        },
        onChange(file, fileList) {
            this.uploadfile = fileList
        },
        remove(cluster_name) {
            ElMessageBox.confirm(
                `您确定要删除集群 ${cluster_name} 吗?`,
                '提示',
                {
                    confirmButtonText: '确认',
                    cancelButtonText: '取消',
                    type: 'warning',
                }
            ).then(() => {
                console.log(cluster_name)
                this.$ajax(
                    {
                        method: 'delete',
                        url: '/cluster/delete',
                        params: {
                            cluster_name: cluster_name
                        }
                    }
                )
                    .then((res) => {
                        this.notify("success", "成功", res.msg)
                        this.reload()
                        console.log(res)
                    }).catch((res) => {
                        this.notify("error", "失败", res.msg)
                        console.log(res);
                    })
            }).catch(() => {
                ElMessage({
                    type: 'info',
                    message: '删除取消',
                })
            })
        },
        notify(status, info, msg) {
            this.$notify({
                title: info,
                message: msg,
                type: status
            });
        },
        Counter() {
            this.acktotal = [];
            this.k8stotal = [];
            for (let i = 0; i < this.clusteritem.length; i++) {
                if (this.clusteritem[i].type == 'k8s') {
                    this.k8stotal.push(this.clusteritem[i])
                } else if (this.clusteritem[i].type == 'ack') {
                    this.acktotal.push(this.clusteritem[i])

                }
            }
        }
    }
}
</script>

<style>
.cluster-header {
    padding-top: 5px;
    padding-left: 10px;
    padding-right: 10px;
    padding-bottom: 5px;
    border-radius: 4px;
    background: #f0f2f5;
    text-align: center;
}

.cluster-header.el-icon-plus {
    font-size: 16px;
    color: #409eff;
    margin-right: 10px;
}


.time {
    font-size: 12px;
    color: #999;
}

.bottom {
    margin-top: 13px;
    line-height: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.button {
    padding: 0;
    min-height: auto;
}

.image {
    width: 100%;
    display: block;

}

.el-row {
    margin-bottom: 0px;
}

.el-row:last-child {
    margin-bottom: 0;
}

.el-col {
    border-radius: 4px;
}

.cluster-card-content {
    position: relative;
    border-radius: 4px;
    padding-top: 10px;
}


.config-add {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);
    width: 178px;
    height: 200px;
}

.config-add {
    border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 178px;
    height: 178px;
    text-align: center;
}

.icon-cluster {
    width: 3em;
    height: 3em;
    vertical-align: -1em;
    fill: currentColor;
    overflow: hidden;
}
</style>
