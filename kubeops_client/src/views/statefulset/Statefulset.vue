<template>
    <el-row>
        <el-col :span="24">
            <div class="grid-content bg-purple">
                <el-col :span="24">
                    <el-row>
                        <el-col :span="2">
                            <div class="grid-content3 ep-bg-purple">名称空间</div>
                        </el-col>
                        <el-col :span="4">
                            <div class="grid-content2 ep-bg-purple">
                                <el-select v-model="namespace" filterable placeholder="Select"
                                    @visible-change="getnsselect()" @change="getStatefulset()" clearable>
                                    <el-option v-for="item in nslist" :key="item.namespace" :label="item.label"
                                        :value="item.namespace" style="width:100%" />
                                </el-select>
                            </div>
                        </el-col>
                        <el-col :span="6">
                        </el-col>
                        <el-col :span="6"></el-col>
                        <el-col :span="6">
                            <div class="grid-content1 ep-bg-purple">
                                <el-input v-model="filter_name" placeholder="Please input" class="input-with-select"
                                    clearable>
                                    <template #prepend>
                                        <el-button icon="Search" @click="getStatefulset()" />
                                    </template>
                                </el-input>
                            </div>
                        </el-col>
                    </el-row>
                </el-col>
                <el-table :data="statefulItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%"
                    size="small">
                    <el-table-column label="Name" width="200">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <el-icon>
                                    <timer />
                                </el-icon>
                                <span style="margin-left: 10px">{{ scope.row.name }}</span>
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
                                        <el-dropdown-item icon="Delete"
                                            @click="messageboxOperate(scope.row, 'delete')">删除</el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown>
                        </template>
                    </el-table-column>
                </el-table>
                <el-row>
                    <el-col :span="6">
                        <div class="grid-content ep-bg-purple" />
                    </el-col>
                    <el-col :span="6">
                        <div class="grid-content ep-bg-purple-light" />
                    </el-col>
                    <el-col :span="4">
                        <div class="grid-content ep-bg-purple" />
                    </el-col>
                    <el-col :span="8">
                        <div class="demo-pagination-block">
                            <el-pagination v-model:current-page="page" v-model:page-size="limit" :page-sizes="page_size"
                                small background layout="sizes, prev, pager, next" :total="total" :pager-count="5"
                                @size-change="handleSizeChange" @current-change="handleCurrentChange" />
                        </div>
                    </el-col>
                </el-row>
            </div>
        </el-col>
    </el-row>
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
</template>

<script scoped>
import { ElFormItem, ElLoading } from 'element-plus'
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
            detailnamespace: null,
            dialogFormVisible: false,
            dialogcreatens: false,
            statefulItem: [],
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
                labels: {
                    "app": "demo"
                },
                container: [
                    {
                        container_name: '',
                        image: '',
                        cpu: '0',
                        memory: '0',
                        container_port: [
                            {
                                port_name: '',
                                container_port: 0,
                                protocol: '',
                            }
                        ]
                    }
                ],
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
        this.getStatefulset()
    },
    methods: {
        getStatefulset() {
            this.$ajax({
                method: 'get',
                url: '/stateful/list',
                params: {
                    filter_name: this.filter_name,
                    namespace: this.namespace,
                    limit: this.limit,
                    page: this.page,
                }
            }).then((res) => {
                this.total = res.data.total
                this.statefulItem = res.data.item
            }).catch(function (res) {
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
            }, 2000)
            this.reload()
        },
        getnsselect() {
            if (this.nslist == "") {
                this.nslist.push({ 'namespace': '', 'label': "All" })
                this.$ajax({
                    method: 'get',
                    url: '/namespaces/list',
                }).then((res) => {
                    res.data.item.forEach(v => {
                        this.nslist.push({ 'namespace': v.name, 'label': v.name })
                    })
                }).catch((res) => {
                    console.log(res.data)
                })
            }
        },
        handleSizeChange(limit) {
            this.limit = limit
            this.getStatefulset()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getStatefulset()
        },
        handleEdit(namespace, name) {
            this.$ajax({
                method: 'get',
                url: '/stateful/detail',
                params: {
                    stateful_set_name: name,
                    namespace: namespace
                }
            }).then((res) => {
                this.activeName = 'json'
                this.aceConfig.lang = "json"
                this.detailnamespace = res.data.metadata.namespace
                this.content = JSON.stringify(res.data, null, 2)
            }).catch(function (res) {
                console.log(res.data);
            })
        },
        handleReplica(namespace, name, replicas) {
            this.$ajax.put(
                '/stateful/modify',
                {
                    stateful_set_name: name,
                    namespace: namespace,
                    replicas: replicas
                },
            ).then((res) => {
                this.Loading("Modifying")
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
        handleDelete(namespace, name) {
            this.Loading("Deleting······")
            this.$ajax({
                method: 'delete',
                url: '/stateful/delete',
                params: {
                    namespace: namespace,
                    stateful_set_name: name
                }
            }
            ).then((res) => {
                this.Loading("Deleting")
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'warning'
                });
            }).catch(function (res) {
                console.log(res);
            })
            this.reload()

        },
        handleUpdate(namespace) {
            let data = this.content
            if (this.aceConfig.lang == 'yaml') {
                data = JSON.stringify(yaml.load(data), null, 2);
            }
            this.$ajax.put(
                '/stateful/update',
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
                this.handleDelete(row.namespaces, row.name)
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
        addContainer() {
            this.ruleForm.container.push({
                container_name: '',
                image: '',
                cpu: '0',
                memory: '0',
                container_port: [
                    {
                        port_name: '',
                        container_port: 0,
                        protocol: '',
                    }
                ]
            })
        },
        addContainerPort(index) {
            this.ruleForm.container[index].container_port.push({
                port_name: '',
                container_port: 0,
                protocol: '',
            })
        },
        removeContainer(item) {
            var index = this.ruleForm.container.indexOf(item)
            if (index !== -1) {
                this.ruleForm.container.splice(index, 1)
            }
        },
        removeContainerPort(portindex, portitem) {
            var index = this.ruleForm.container[portindex].container_port.indexOf(portitem)
            if (index !== -1) {
                this.ruleForm.container[portindex].container_port.splice(index, 1)
            }
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
        }
    }
}
</script>

<style>
.el-row {
    margin-bottom: 10px;

    &:last-child {
        margin-bottom: 0;
    }
}

.el-col {
    border-radius: 4px;
}

.bg-purple {
    background: #f0f2f5;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
}

.grid-content {
    border-radius: 4px;
    min-height: 10px;
    padding: 10px;
}

.el-row:last-child {
    margin-bottom: 0;
}

.grid-content1 {

    border-radius: 4px;
}

.grid-content3 {
    position: fixed;
    margin-top: 5px;
    margin-left: 20px;
}

.grid-content2 {

    border-radius: 4px;
}

.el-table,
.el-table__expanded-cell {
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

.el-dialog--center .el-dialog__body {
    min-height: 400px;
}

.demo-pagination-block+.demo-pagination-block {
    margin-top: 10px;
}

.el-form-item__content .el-input {
    width: 400px;
}

.el-tabs--border-card>.el-tabs__content {
    padding: 0px;
}
</style>