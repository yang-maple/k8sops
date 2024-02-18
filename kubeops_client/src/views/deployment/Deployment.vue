<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-deploy" aria-hidden="true">
                            <use xlink:href="#icon-deployment-copy"></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">无状态副本集</span>
                        <br>
                        <span style="font-size: 12px;color: #79879c!important">
                            无状态副本集（Deployment）是指在 Kubernetes 集群中运行的不保存任何状态或数据软件应用程序，它们将所有需要保存的状态存储在外部系统中，并且在重启或迁移时可以很快恢复。
                        </span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>
    <el-row>
        <el-col :span="4">
            <div class="header-grid-content">
                <el-select v-model="namespace" filterable placeholder="命名空间（默认ALL）" @change="getDeployment()" clearable>
                    <el-option-group v-for="group in nslist" :key="group.label" :label="group.label">
                        <el-option v-for="item in group.options" :key="item.namespace" :label="item.label"
                            :value="item.namespace" />
                    </el-option-group>
                </el-select>
            </div>
        </el-col>
        <el-col :span="16">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="请输入搜索资源的名称" clearable @clear="getDeployment()"
                    @keyup.enter="getDeployment()">
                    <template #prepend>
                        <el-button icon="Search" @click="getDeployment()"></el-button>
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="4">
            <div class="header-grid-content" style="text-align: right;">
                <el-button type="info" @click="getDeployment()" round>
                    <el-icon>
                        <Refresh />
                    </el-icon>
                </el-button>
                <el-button type="info" @click="dialogcreatens = true" round>
                    创建
                </el-button>
            </div>
        </el-col>
    </el-row>
    <div class="table-bg-purple">
        <el-table :data="deploymentItem" :header-cell-style="{ background: '#e6e7e9' }" size="small" empty-text="抱歉，暂无数据">
            <el-table-column label="名称" width="200">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <svg class="icon-table-deploy" aria-hidden="true">
                            <use xlink:href="#icon-deployment-copy"></use>
                        </svg>
                        <el-text style="margin-left:3px" size="small">{{
                            scope.row.name
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="命名空间" prop="namespaces" width="150" />
            <el-table-column label="标签">
                <template #default="scope">
                    <div v-for="(v, k, index) in scope.row.labels " :key="k">
                        <div v-if="index < maxitem[scope.$index]">
                            <el-tag type="info" style="margin-left: 5px;" size="small" effect="plain" round>
                                {{ k }}:{{ v }}</el-tag>
                        </div>
                    </div>
                    <div v-if="scope.row.labels == null">---</div>
                    <div v-if="scope.row.labels != null && Object.keys(scope.row.labels).length > 3"><el-button size="small"
                            type="primary" link @click="showLabels(scope.$index)">{{
                                maxitem[scope.$index] == 3 ?
                                '展开' : '收起'
                            }}</el-button></div>
                </template>
            </el-table-column>
            <el-table-column label="容器组数量" prop="pods" width="100" align="center" />
            <el-table-column label="状态" prop="status" width="100">
                <template #default="scope">
                    <div v-if="scope.row.status == 'Running'" style="color: #00a2ae">
                        <el-tag type="success" size="small" effect="plain" round>
                            Running
                        </el-tag>
                    </div>
                    <div v-if="scope.row.status != 'Running'" style="color: #f56c6c">
                        <el-tag type="danger" size="small" effect="plain" round>
                            {{ scope.row.status }}
                        </el-tag>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="镜像" width="200">
                <template #default="scope">
                    <div v-for="(v, k) in scope.row.image " :key="k">{{ v }}<br></div>
                </template>
            </el-table-column>
            <el-table-column label="创建时间" prop="age" width="140" />
            <el-table-column label="操作" width="70" align="center">
                <template #default="scope">
                    <el-dropdown trigger="click">
                        <el-button type="primary" link>
                            <el-icon :size="24">
                                <Operation />
                            </el-icon>
                        </el-button>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item icon="Edit"
                                    @click="dialogFormVisible = true, handleEdit(scope.row.namespaces, scope.row.name)">编辑</el-dropdown-item>
                                <el-dropdown-item icon="DocumentAdd" @click="messageboxReplica(scope.row)">
                                    副本
                                </el-dropdown-item>
                                <el-dropdown-item icon="Refresh"
                                    @click="messageboxOperate(scope.row, '重启')">重启</el-dropdown-item>
                                <el-dropdown-item icon="Delete"
                                    @click="messageboxOperate(scope.row, '删除')">删除</el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </template>
            </el-table-column>
        </el-table>
        <el-row>
            <el-col :span="1.5">
                <div class="demo-pagination-block">
                    <el-text type="info">Total:{{ total }}</el-text>
                </div>
            </el-col>
            <el-col :span="22.5">
                <div class="demo-pagination-block" style="padding-bottom: 5px;">
                    <el-pagination v-model:current-page="page" v-model:page-size="limit" :page-sizes="page_size" small
                        background layout="sizes, prev, pager, next" :total="total" :pager-count="5"
                        @size-change="handleSizeChange" @current-change="handleCurrentChange" />
                </div>
            </el-col>
        </el-row>
    </div>


    <el-dialog v-model="dialogFormVisible" title="实例详情" center>
        <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
            <el-tab-pane label="Yaml" name="yaml"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                    style="min-height: 350px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
            <el-tab-pane label="Json" name="json"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                    style="min-height: 350px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
        </el-tabs>
        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="handleUpdate(detailnamespace)">更新</el-button>
                <el-button @click="dialogFormVisible = false">
                    取消
                </el-button>
            </span>
        </template>
    </el-dialog>
    <el-dialog v-model="dialogcreatens" title="创建资源" center>
        <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="25%">
            <el-form-item label="名称" prop="name">
                <el-input v-model="ruleForm.name"></el-input>
            </el-form-item>
            <el-form-item label="命名空间" prop="namespace">
                <el-select v-model="ruleForm.namespace" placeholder="请选择命名空间">
                    <el-option v-for="item in nslist[1].options" :key="item.namespace" :label="item.label"
                        :value="item.namespace" />
                </el-select>
            </el-form-item>
            <el-form-item label="标签" prop="labels">
                <el-input v-model="ruleForm.labels" placeholder="`请输入如下格式 {'a':'b'}`"></el-input>
            </el-form-item>
            <el-form-item label="副本" prop="replicas">
                <el-input-number v-model="ruleForm.replicas" :min="1" :max="99" />
            </el-form-item>
            <el-form-item label="容器名称" prop="container.container_name">
                <el-input v-model="ruleForm.container.container_name"></el-input>
            </el-form-item>
            <el-form-item label="镜像" prop="container.image">
                <el-input v-model="ruleForm.container.image"></el-input>
            </el-form-item>
            <!-- <el-form-item label="Cpu需求">
                <el-input v-model="ruleForm.container.cpu" placeholder="单位为 m"></el-input>m
            </el-form-item>
            <el-form-item label="Memory需求">
                <el-input v-model="ruleForm.container.memory" placeholder="单位为 Mi"></el-input>Mi
            </el-form-item> -->
            <el-form-item label="端口名称" prop="container.container_port.port_name">
                <el-input v-model="ruleForm.container.container_port.port_name"></el-input>
            </el-form-item>
            <el-form-item label="端口" prop="container.container_port.container_port">
                <el-input v-model.number="ruleForm.container.container_port.container_port"></el-input>
            </el-form-item>
            <el-form-item label="协议" prop="container.container_port.protocol">
                <el-radio-group v-model="ruleForm.container.container_port.protocol">
                    <el-radio label="TCP"></el-radio>
                    <el-radio label="UDP"></el-radio>
                </el-radio-group>
            </el-form-item>

            <el-form-item label="健康检查" prop="delivery">
                <el-switch v-model="ruleForm.health_check"></el-switch>
            </el-form-item>
            <el-form-item label="路径" v-show="ruleForm.health_check == true"
                :prop="ruleForm.health_check ? 'health_path' : 'nocheck'">
                <el-input v-model="ruleForm.health_path"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="submitForm('ruleForm')">立即创建</el-button>
                <el-button @click="resetForm('ruleForm')">立即重置</el-button>
                <el-button type="danger" @click="dialogcreatens = false, resetForm('ruleForm')">
                    取消
                </el-button>
            </el-form-item>
        </el-form>
    </el-dialog>
</template>

<script scoped>
import { VAceEditor } from 'vue3-ace-editor';
import '../../ace.config.js';
import yaml from 'js-yaml';
export default {
    inject: ['reload'],
    components: {
        VAceEditor
    },
    data() {
        return {
            maxitem: [],
            detailnamespace: null,
            dialogFormVisible: false,
            dialogcreatens: false,
            deploymentItem: [],
            filter_name: '',
            namespace: '',
            limit: 10,
            page: 1,
            content: '',
            aceConfig: {
                lang: 'json',
                theme: "cloud9_day",
                options: {
                    showPrintMargin: false,
                }
            },
            activeName: 'json',
            nslist: [],
            total: 0,
            page_size: [1, 10, 20, 50, 100],
            ruleForm: {
                name: '',
                namespace: '',
                replicas: 1,
                labels: null,
                container:
                {
                    container_name: '',
                    image: '',
                    cpu: '0m',
                    memory: '0Mi',
                    container_port:
                    {
                        port_name: '',
                        container_port: null,
                        protocol: '',
                    }
                },
                health_check: false,
                health_path: '',
            },
            rules: {
                name: [
                    { required: true, message: '请输入资源名称', trigger: 'blur' },
                ],
                namespace: [
                    { required: true, message: '请选择命名空间', trigger: 'change' }
                ],
                labels: [
                    { required: true, message: '请输入标签', trigger: 'blur' },
                ],
                replica: [
                    { required: true, message: '请输入副本数量', trigger: 'change' }
                ],
                'container.container_name': [
                    { required: true, message: '请输入容器名称', trigger: 'blur' },
                ],
                'container.image': [
                    { required: true, message: '请输入镜像名称', trigger: 'blur' },
                ],
                'container.container_port.container_port': [
                    { required: true, message: '请输入端口', trigger: 'change' },
                    { type: 'number', min: 1, max: 65535, message: '端口范围为1-65535', trigger: 'change' },
                ],
                'container.container_port.protocol': [
                    { required: true, message: '请选择一种协议', trigger: 'change' },
                ],
                health_path: [
                    { required: true, message: '请输入健康检查路径', trigger: 'blur' },
                ],
                nocheck: [],
            }
        }
    },
    created() {
        this.getDeployment()
        this.getnsselect()
    },
    methods: {
        getDeployment() {
            this.$ajax({
                method: 'get',
                url: '/deploy/list',
                params: {
                    filter_name: this.filter_name,
                    namespace: this.namespace,
                    limit: this.limit,
                    page: this.page,
                }
            }).then((res) => {
                this.total = res.data.total
                this.deploymentItem = res.data.item
                for (let i = 0; i < res.data.item.length; i++) {
                    this.maxitem.push(3)
                }
            }).catch((res) => {
                console.log(res);
            })
        },
        createDeployment() {
            this.$ajax.post(
                '/deploy/create',
                {
                    data: this.ruleForm
                },
            ).then((res) => {
                this.$message({
                    message: res.msg,
                    type: 'success'
                });
                this.dialogcreatens = false
                this.Refresh()
            }).catch((res) => {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            })
        },
        Refresh() {
            setTimeout(() => {
                this.reload()
            }, 1000)
        },
        getnsselect() {
            if (this.nslist == "") {
                this.nslist.push({
                    label: '',
                    options: [
                        {
                            namespace: '',
                            label: '全部空间',
                        },
                    ],
                })
                this.$ajax({
                    method: 'get',
                    url: '/namespaces/list',
                }).then((res) => {
                    this.nslist.push({
                        label: '',
                        options: [],
                    })
                    res.data.item.forEach(v => {
                        this.nslist[1].options.push({ 'namespace': v.name, 'label': v.name })
                    })
                }).catch((res) => {
                    this.$message({
                        showClose: true,
                        message: "获取名称空间失败",
                        type: 'error'
                    });
                    console.log(res.data)
                })
            }
        },
        handleSizeChange(limit) {
            this.limit = limit
            this.page = 1
            this.getDeployment()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getDeployment()
        },
        handleEdit(namespace, name) {
            this.$ajax({
                method: 'get',
                url: '/deploy/detail',
                params: {
                    deploy_name: name,
                    namespace: namespace
                }
            }).then((res) => {
                this.activeName = 'json'
                this.aceConfig.lang = "json"
                this.detailnamespace = res.data.metadata.namespace
                this.content = JSON.stringify(res.data, null, 2)
            }).catch((res) => {
                console.log(res.data);
            })
        },
        handleReplica(namespace, name, replicas) {
            this.$ajax.put(
                '/deploy/modify',
                {
                    deploy_name: name,
                    namespace: namespace,
                    replicas: replicas
                },
            ).then((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'success'
                });
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
            })
            this.Refresh()
        },
        handleRestart(namespace, name) {
            this.$ajax.post(
                '/deploy/restart',
                {
                    deploy_name: name,
                    namespace: namespace,
                },
            ).then((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'success'
                });
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
                console.log(res);
            })
            this.Refresh()

        },
        handleDelete(namespace, name) {
            this.$ajax({
                method: 'delete',
                url: '/deploy/delete',
                params: {
                    namespace: namespace,
                    deploy_name: name
                }
            }
            ).then((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'warning'
                });

            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
                console.log(res);
            })
            this.Refresh()
        },
        handleUpdate(namespace) {
            let data = this.content
            let datajson = {}
            try {
                if (this.aceConfig.lang == 'yaml') {
                    data = JSON.stringify(yaml.load(data), null, 2);
                }
                datajson = JSON.parse(data)
            } catch (e) {
                this.$message({
                    showClose: true,
                    message: '格式错误,请检查格式',
                    type: 'error'
                });
                return
            }
            this.$ajax.put(
                '/deploy/update',
                {
                    namespace: namespace,
                    data: datajson
                },
            ).then((res) => {
                this.dialogFormVisible = false
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'success'
                });
                this.Refresh()
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
                console.log(res);
            })
        },
        messageboxReplica(row) {
            this.$prompt('修改的副本数为', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                inputType: "number",
                inputPattern: /^([1-9]|[1-9][0-9]|0?[0-9])$/,
                inputErrorMessage: '副本的数量范围 0-99'
            }).then(({ value }) => {
                this.handleReplica(row.namespaces, row.name, parseInt(value))
                console.log(value)
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '取消输入'
                });
            });
        },
        messageboxOperate(row, name) {
            this.$confirm(`是否${name}实例${row.name}`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                if (name == "重启") {
                    this.handleRestart(row.namespaces, row.name)
                } else if (name == "删除") {
                    this.handleDelete(row.namespaces, row.name)
                }

            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                });
            });
        },
        submitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    try {
                        this.ruleForm.labels = JSON.parse(this.ruleForm.labels);
                        this.createDeployment()
                    } catch (error) {
                        this.$message({
                            showClose: true,
                            message: '请填写正确格式，例如: {"a":"b"}',
                            type: 'error'
                        });
                    }
                } else {
                    return false;
                }
            });
        },
        resetForm(formName) {
            this.$refs[formName].resetFields();
        },
        yamlFormat() {
            if (this.aceConfig.lang == "yaml") {
                return
            }
            this.aceConfig.lang = "yaml"
            this.content = yaml.dump(JSON.parse(this.content))
        },
        jsonFormat() {
            if (this.aceConfig.lang == "json") {
                return
            }
            this.aceConfig.lang = "json"
            this.content = JSON.stringify(yaml.load(this.content), null, 2);
        },
        handleClick(tab) {
            if (tab.props.name == "yaml") {
                this.yamlFormat()
                return
            }
            this.jsonFormat()
        },
        showLabels(index) {
            if (this.maxitem[index] == 3) {
                this.maxitem[index] = 99
            } else {
                this.maxitem[index] = 3
            }
        }
    }
}
</script>

<style>
.el-row {
    margin-bottom: 0px;
}

.el-row:last-child {
    margin-bottom: 0;
}

.el-col {
    border-radius: 4px;
}

.header-grid-content {
    padding-top: 5px;
    padding-left: 10px;
    padding-right: 10px;
    padding-bottom: 5px;
    border-radius: 4px;
    background: #f0f2f5;
}


.table-bg-purple {
    padding-right: 2px;
    padding-left: 2px;
    border-radius: 4px;
    background: #f0f2f5;
}

.el-table,
.el-table__expanded-cell {
    border-radius: 4px;
    background-color: transparent;
}

.el-table th,
.el-table tr {
    background-color: transparent;
}

.el-button--text {
    margin-right: 15px;
}




.dialog-footer button:first-child {
    margin-right: 10px;
}

.demo-pagination-block {
    margin-top: 5px;
    margin-bottom: 5px;
    padding-left: 10px;
}

.el-form-item__content .el-input {
    width: 400px;
}

.el-tabs--border-card>.el-tabs__content {
    padding: 0px;
}

.dotClass {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    display: block;
    margin-left: 10px;
}

.icon-deploy {
    width: 2.5em;
    height: 2.5em;
    vertical-align: -0.7em;
    fill: currentColor;
    overflow: hidden;
}

.icon-table-deploy {
    width: 1.5em;
    height: 1.5em;
    vertical-align: -0.5em;
    fill: currentColor;
    overflow: hidden;
}
</style>