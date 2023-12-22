<template>
    <el-row>
        <el-col :span="4">
            <div class="header-grid-content">
                <el-select v-model="namespace" filterable placeholder="Namespace" @visible-change="getnsselect()"
                    @change="getDeployment()" clearable>
                    <el-option-group v-for="group in nslist" :key="group.label" :label="group.label">
                        <el-option v-for="item in group.options" :key="item.namespace" :label="item.label"
                            :value="item.namespace" />
                    </el-option-group>
                </el-select>
            </div>
        </el-col>
        <el-col :span="16">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="Please input" clearable>
                    <template #prepend>
                        <el-button icon="Search"><span style="font-size: 15px;">
                                搜索
                            </span></el-button>
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="4">
            <div class="header-grid-content" style="text-align: center;">
                <el-button type="info" @click="Loading(`Refresh`)" round>
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
        <el-table :data="deploymentItem" :header-cell-style="{ background: '#e6e7e9' }" size="small">
            <el-table-column label="Name" width="200">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <el-icon>
                            <cpu />
                        </el-icon>
                        <span style="margin-left: 5px">{{ scope.row.name }}</span>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="Namespace" prop="namespaces" width="150" />
            <el-table-column label="Images" width="200" align="center">
                <template #default="scope">
                    <div v-for="(v, k) in scope.row.image " :key="k">{{ v }}<br></div>
                </template>
            </el-table-column>
            <el-table-column label="labels" width="300" align="center">
                <template #default="scope">
                    <el-tag type="info" size="small" style="margin-left: 5px;" v-for="(v, k) in scope.row.labels "
                        :key="k">{{ k }}:{{ v
                        }}<br></el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Status" prop="status" width="100" />
            <el-table-column label="Pods" prop="pods" width="70" />
            <el-table-column label="Age" prop="age" width="50" />
            <el-table-column label="Operations" align="center">
                <template #default="scope">
                    <el-dropdown>
                        <el-button type="primary">
                            Operate<el-icon class="el-icon--right"><arrow-down /></el-icon>
                        </el-button>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item icon="Edit"
                                    @click="dialogFormVisible = true, handleEdit(scope.row.namespaces, scope.row.name)">编辑</el-dropdown-item>
                                <el-dropdown-item icon="DocumentAdd" @click="messageboxReplica(scope.row)">
                                    副本
                                </el-dropdown-item>
                                <el-dropdown-item icon="Refresh"
                                    @click="messageboxOperate(scope.row, 'restart')">重启</el-dropdown-item>
                                <el-dropdown-item icon="Delete"
                                    @click="messageboxOperate(scope.row, 'delete')">删除</el-dropdown-item>
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
                <el-button type="primary" @click="dialogFormVisible = false, handleUpdate(detailnamespace)">更新</el-button>
                <el-button @click="dialogFormVisible = false">
                    取消
                </el-button>
            </span>
        </template>
    </el-dialog>
    <el-dialog v-model="dialogcreatens" title="创建资源" center>
        <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="25%" status-icon>
            <el-form-item label="名称" prop="name">
                <el-input v-model="ruleForm.name"></el-input>
            </el-form-item>
            <el-form-item label="名称空间" prop="namespace">
                <el-select v-model="ruleForm.namespace" placeholder="请选择活动区域" @visible-change="getnsselect()">
                    <el-option-group v-for="group in  nslist " :key="group.label">
                        <el-option v-for=" item  in  group.options " :key="item.namespace" :value="item.namespace"
                            v-show="item.label != 'ALL'" />
                    </el-option-group>
                </el-select>
            </el-form-item>
            <el-form-item label="标签">
                <el-input v-model="ruleForm.labels"></el-input>
            </el-form-item>
            <el-form-item label="副本">
                <el-input-number v-model="ruleForm.replicas" :min="1" :max="99" />
            </el-form-item>
            <el-form-item label="容器名称" prop="container_name">
                <el-input v-model="ruleForm.container.container_name"></el-input>
            </el-form-item>
            <el-form-item label="镜像" prop="image">
                <el-input v-model="ruleForm.container.image"></el-input>
            </el-form-item>
            <el-form-item label="Cpu限制">
                <el-input v-model="ruleForm.container.cpu"></el-input>
            </el-form-item>
            <el-form-item label="Memory限制">
                <el-input v-model="ruleForm.container.memory"></el-input>
            </el-form-item>

            <el-form-item label="端口名称">
                <el-input v-model="ruleForm.container.container_port.port_name"></el-input>
            </el-form-item>
            <el-form-item label="端口">
                <el-input type="number" v-model.number="ruleForm.container.container_port.container_port"></el-input>
            </el-form-item>
            <el-form-item label="协议" prop="resource">
                <el-radio-group v-model="ruleForm.container.container_port.protocol">
                    <el-radio label="TCP"></el-radio>
                    <el-radio label="UDP"></el-radio>
                </el-radio-group>
            </el-form-item>

            <el-form-item label="健康检查" prop="delivery">
                <el-switch v-model="ruleForm.health_check"></el-switch>
            </el-form-item>
            <el-form-item label="路径" v-show="ruleForm.health_check == true">
                <el-input v-model="ruleForm.health_path"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="submitForm('ruleForm')">立即创建</el-button>
                <el-button @click="resetForm('ruleForm')">立即重置</el-button>
                <el-button type="danger" @click="dialogcreatens = false">
                    取消
                </el-button>
            </el-form-item>
        </el-form>
    </el-dialog>
</template>

<script scoped>
import { ElFormItem, ElLoading } from 'element-plus'
import { VAceEditor } from 'vue3-ace-editor';
import '../../ace.config.js';
import yaml from 'js-yaml';
import { SelectProps } from 'element-plus/es/components/select-v2/src/defaults';
export default {
    inject: ['reload'],
    components: {
        VAceEditor
    },
    data() {
        return {
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
                replicas: 0,
                labels: {},
                container:
                {
                    container_name: '',
                    image: '',
                    cpu: '0',
                    memory: '0',
                    container_port:
                    {
                        port_name: '',
                        container_port: 0,
                        protocol: '',
                    }
                }
            },
            rules: {
                name: [
                    { required: true, message: '请输入资源名称', trigger: 'blur' },
                    { min: 3, max: 5, message: '长度在 3 到 5 个字符', trigger: 'blur' }
                ],
                namespace: [
                    { required: true, message: '请选择名称空间', trigger: 'change' }
                ],
                replica: [
                    { required: true, message: '请输入副本数量', trigger: 'change' }
                ],
                resource: [
                    { message: '请选择一种协议', trigger: 'change' }
                ],
            }
        }
    },
    created() {
        this.getDeployment()
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
                console.log(res.data.item)
                this.total = res.data.total
                this.deploymentItem = res.data.item
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
                this.Loading("Creating")
                this.$message({
                    message: res.msg,
                    type: 'success'
                });
                console.log(res)
            }).catch((res) => {
                console.log(res);
            })
        },
        Loading(msg) {
            const loading = ElLoading.service({
                lock: true,
                text: msg + ".....",
                background: 'rgba(0, 0, 0, 0.7)',
            })
            setTimeout(() => {
                loading.close()
                this.reload()
            }, 2000)
        },
        getnsselect() {
            if (this.nslist == "") {
                this.nslist.push({
                    label: '',
                    options: [
                        {
                            namespace: '',
                            label: 'ALL',
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
                    console.log(this.nslist)
                }).catch((res) => {
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
                    message: res.msg + res.reason,
                    type: 'error'
                });
                console.log(res);
            })

        },
        handleRestart(namespace, name) {
            this.Loading("Restarting······")
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
                console.log(res)
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg + res.reason,
                    type: 'error'
                });
                console.log(res);
            })


        },
        handleDelete(namespace, name) {
            this.Loading("Deleting······")
            this.$ajax({
                method: 'delete',
                url: '/deploy/delete',
                params: {
                    namespace: namespace,
                    deploy_name: name
                }
            }
            ).then((res) => {
                this.Loading("Deleting")
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'warning'
                });
            }).catch((res) => {
                console.log(res);
            })
        },
        handleUpdate(namespace) {
            let data = this.content
            if (this.aceConfig.lang == 'yaml') {
                data = JSON.stringify(yaml.load(data), null, 2);
            }
            this.$ajax.put(
                '/deploy/update',
                {
                    namespace: namespace,
                    data: JSON.parse(data)
                },
            ).then((res) => {
                this.Loading("Updateing")
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'success'
                });
                console.log(res)
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg + res.reason,
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
                if (name == "restart") {
                    this.handleRestart(row.namespaces, row.name)
                } else if (name == "delete") {
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
                    this.dialogcreatens = false
                    this.ruleForm.labels = JSON.parse(this.ruleForm.labels);
                    console.log(this.ruleForm)
                    this.createDeployment()
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

.el-dialog {
    .el-select {
        .el-input {
            width: 180px;
        }
    }
}
</style>