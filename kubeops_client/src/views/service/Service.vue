<template>
    <el-row>
        <el-col :span="4">
            <div class="header-grid-content">
                <el-select v-model="namespace" placeholder="Select" @visible-change="getnsselect()" @change="getSevices()"
                    clearable>
                    <el-option v-for="item in nslist" :key="item.namespace" :label="item.label" :value="item.namespace" />
                </el-select>
            </div>
        </el-col>
        <el-col :span="16">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="Please input" class="input-with-select" clearable>
                    <template #prepend>
                        <el-button icon="Search" @click="getSevices()" />
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="4">
            <div class="header-grid-content" style="text-align: center;">
                <el-button type="info" @click="Refresh()" round>
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
        <el-table :data="serviceItem" :header-cell-style="{ background: '#e6e7e9' }" size="small">
            <el-table-column label="Name" prop="name" width="200" align="center">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        {{ scope.row.name
                        }}
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="Namespace" prop="namespace" width="150" align="center" />
            <el-table-column label="labels" width="280" align="center">
                <template #default="scope">
                    <el-tag type="info" size="small" v-for="(v, k) in scope.row.labels " :key="k">{{ k }}:{{ v
                    }}<br></el-tag>
                </template>
            </el-table-column>
            <el-table-column label="Type" prop="type" width="100" align="center" />
            <el-table-column label="ClusterIP" prop="cluster_ip" width="100" align="center" />
            <el-table-column label="ExternalIP" width="100" align="center">
                <template #default="scope">
                    <div v-if="scope.row.external_ip == undefined">-</div>
                    <div v-if="scope.row.external_ip != undefined" v-for="(v, k) in scope.row.external_ip " :key="k">{{ v
                    }}</div>
                </template>
            </el-table-column>
            <el-table-column label="Ports" width="130" align="center">
                <template #default="scope">
                    <div v-for="(v, k) in scope.row.port " :key="k">
                        <span v-if="v.nodePort == undefined">{{ v.port }}/{{ v.protocol }}</span>
                        <span v-if="v.nodePort != undefined">{{ v.port }} : {{ v.nodePort }} / {{ v.protocol
                        }}</span>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="Age" prop="age" width="50" align="center" />
            <el-table-column label="Operations" align="center">
                <template #default="scope">
                    <el-dropdown>
                        <el-button type="primary">
                            Operate<el-icon class="el-icon--right"><arrow-down /></el-icon>
                        </el-button>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item icon="Edit"
                                    @click="dialogFormVisible = true, handleEdit(scope.row.namespace, scope.row.name)">编辑</el-dropdown-item>
                                <el-dropdown-item icon="Delete"
                                    @click="messageboxOperate(scope.row, 'delete')">删除</el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </template>
            </el-table-column>
        </el-table>
        <el-dialog v-model="dialogFormVisible" title="实例详情/json" center>
            <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
                <el-tab-pane label="Yaml" name="yaml"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                        style="min-height: 350px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
                <el-tab-pane label="Json" name="json"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                        style="min-height: 350px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
            </el-tabs>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary"
                        @click="dialogFormVisible = false, handleUpdate(detailnamespace)">更新</el-button>
                    <el-button @click="dialogFormVisible = false">
                        取消
                    </el-button>
                </span>
            </template>
        </el-dialog>
        <el-row>
            <el-col :span="1.5">
                <div class="demo-pagination-block">
                    <el-text type="info">Total:{{ total }}</el-text>
                </div>
            </el-col>
            <el-col :span="22.5">
                <div class="demo-pagination-block">
                    <el-pagination v-model:current-page="page" v-model:page-size="limit" :page-sizes="page_size"
                        :pager-count="5" small background layout="sizes, prev, pager, next" :total="total"
                        @size-change="handleSizeChange" @current-change="handleCurrentChange" />
                </div>
            </el-col>
        </el-row>
    </div>

    <el-dialog v-model="dialogcreatens" title="创建资源" center>
        <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="25%" status-icon>
            <el-form-item label="名称" prop="name">
                <el-input v-model="ruleForm.name"></el-input>
            </el-form-item>
            <el-form-item label="名称空间" prop="namespace">
                <el-select v-model="ruleForm.namespace" placeholder="请选择名称空间" @visible-change="getnsselect()">
                    <el-option v-for="item in nslist" :key="item.namespace" :label="item.label" :value="item.namespace"
                        v-show="item.label != 'All'" />
                </el-select>
            </el-form-item>
            <el-form-item label="标签" prop="labels">
                <el-input v-model="ruleForm.labels"></el-input>
            </el-form-item>
            <el-form-item label="类型" prop="type">
                <el-select v-model="ruleForm.type" placeholder="请选择服务类型">
                    <el-option v-for="item in typelist" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>


            <el-form-item label="容器名称" prop="port_name">
                <el-input v-model="ruleForm.service_ports.port_name"></el-input>
            </el-form-item>
            <el-form-item label="服务端口" prop="ports">
                <el-input v-model.number="ruleForm.service_ports.port"></el-input>
            </el-form-item>
            <el-form-item label="协议" prop="resource">
                <el-radio-group v-model="ruleForm.service_ports.protocol">
                    <el-radio label="TCP"></el-radio>
                    <el-radio label="UDP"></el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="暴露端口">
                <el-input type="number" v-model.number="ruleForm.service_ports.target_port"></el-input>
            </el-form-item>
            <el-form-item label="访问端口" v-show="ruleForm.type == 'NodePort' || ruleForm.type == 'LoadBalancer'">
                <el-input type="number" v-model.number="ruleForm.service_ports.node_port"></el-input>
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
import { ElFormItem } from 'element-plus'
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
            detailnamespace: '',
            dialogFormVisible: false,
            dialogcreatens: false,
            serviceItem: [],
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
            typelist: [
                {
                    value: 'ClusterIP',
                    label: 'ClusterIP'
                },
                {
                    value: 'NodePort',
                    label: 'NodePort'
                },
                {
                    value: 'LoadBalancer',
                    label: 'LoadBalancer'
                }
            ],
            ruleForm: {
                name: null,
                namespace: null,
                labels: null,
                type: null,
                service_ports: {
                    port_name: null,
                    port: null,
                    protocol: null,
                    target_port: null,
                    node_port: null
                },
            },
            rules: {
                name: [
                    { required: true, message: '请输入资源名称', trigger: 'change' },
                    { min: 3, max: 10, message: '长度在 3 到 10 个字符', trigger: 'change' }
                ],
                namespace: [
                    { required: true, message: '请选择名称空间', trigger: 'change' }
                ],
                resource: [
                    { message: '请选择一种协议', trigger: 'change' }
                ],
                type: [
                    { required: true, message: '请选择一种服务类型', trigger: 'change' }
                ],
            }
        }
    },
    created() {
        this.getSevices()
    },
    methods: {
        getSevices() {
            this.$ajax({
                method: 'get',
                url: '/svc/list',
                params: {
                    filter_name: this.filter_name,
                    namespace: this.namespace,
                    limit: this.limit,
                    page: this.page,
                }
            }).then((res) => {
                this.total = res.data.total
                this.serviceItem = res.data.item
            }).catch(function (res) {
                console.log(res);
            })
        },
        createService() {
            this.$ajax.post(
                '/svc/create',
                {
                    data: this.ruleForm
                },
            ).then((res) => {
                this.$message({
                    message: res.msg,
                    type: 'success'
                });
            }).catch((res) => {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            })
            this.Refresh()
        },
        Refresh() {
            setTimeout(() => {
                this.reload()
            }, 2000)
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
            this.page = 1
            this.getSevices()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getSevices()
        },
        handleEdit(namespace, name) {
            this.$ajax({
                method: 'get',
                url: '/svc/detail',
                params: {
                    service_name: name,
                    namespace: namespace
                }
            }).then((res) => {
                this.activeName = 'json'
                this.aceConfig.lang = "json"
                this.detailnamespace = res.data.detail.metadata.namespace
                this.content = JSON.stringify(res.data.detail, null, 2)
            }).catch((res) => {
                console.log(res.data);
            })
        },
        handleDelete(namespace, name) {
            this.$ajax({
                method: 'delete',
                url: '/svc/delete',
                params: {
                    namespace: namespace,
                    service_name: name
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
                    message: res.reason,
                    type: 'error'
                });
            })
            this.Refresh()
        },
        handleUpdate(namespace) {
            let data = this.content
            if (this.aceConfig.lang == 'yaml') {
                data = JSON.stringify(yaml.load(data), null, 2);
            }
            this.$ajax.put(
                '/svc/update',
                {
                    namespace: namespace,
                    data: JSON.parse(data)
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
                    message: res.reason,
                    type: 'error'
                });
            })
            this.Refresh()
        },
        messageboxOperate(row, name) {
            this.$confirm(`是否${name}实例${row.name}`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.handleDelete(row.namespace, row.name)
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
                    if (this.ruleForm.labels != null) {
                        this.ruleForm.labels = JSON.parse(this.ruleForm.labels)
                    }
                    this.dialogcreatens = false
                    console.log(this.ruleForm)
                    this.createService()
                } else {
                    return false;
                }
            });
        },
        resetForm(formName) {
            this.$refs[formName].resetFields();
        },
        addService() {
            this.ruleForm.service_ports.push({
                port_name: null,
                port: null,
                protocol: null,
                target_port: null,
                node_port: null
            })
        },
        handle(row) {
            this.$router.push({
                path: '/services/services_detail',
                name: 'services 详情',
                query: {
                    namespace: row.namespace,
                    service_name: row.name
                }
            })
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

.el-tabs--border-card>.el-tabs__content {
    padding: 0px;
}


.el-form-item__content .el-input {
    width: 400px;
}

.el-dialog {
    .el-select {
        .el-input {
            width: 180px;
        }
    }
}
</style>