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
                                    @visible-change="getnsselect()" @change="getIngress()" clearable>
                                    <el-option v-for="item in nslist" :key="item.namespace" :label="item.label"
                                        :value="item.namespace" style="width:100%" />
                                </el-select>
                            </div>
                        </el-col>
                        <el-col :span="6">
                            <div class="grid-content1 ep-bg-purple">
                                <el-button type="primary" @click="dialogcreatens = true" style="margin-left: 18px;">
                                    创建 Ingress 资源
                                </el-button>
                            </div>
                        </el-col>
                        <el-col :span="6"></el-col>
                        <el-col :span="6">
                            <div class="grid-content1 ep-bg-purple">
                                <el-input v-model="filter_name" placeholder="Please input" class="input-with-select"
                                    clearable>
                                    <template #prepend>
                                        <el-button icon="Search" @click="getIngress()" />
                                    </template>
                                </el-input>
                            </div>
                        </el-col>
                    </el-row>
                </el-col>
                <el-table :data="ingressItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%"
                    size="small">
                    <el-table-column label="Name" width="300" align="center">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <el-icon>
                                    <timer />
                                </el-icon>
                                <span style="margin-left: 10px">{{ scope.row.name }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column label="labels" width="300" align="center">
                        <template #default="scope">
                            <el-tag type="info" size="small" v-for="(v, k) in scope.row.labels " :key="k">{{ k }}:{{ v
                            }}<br></el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="Endpoints" width="200" align="center">
                        <template #default="scope">
                            <div v-for="(v, k) in scope.row.endpoint " :key="k">{{ v.ip }}<br></div>
                        </template>
                    </el-table-column>
                    <el-table-column label="Hosts" width="200" align="center">
                        <template #default="scope">
                            <div v-for="(v, k) in scope.row.host " :key="k">{{ v }}<br></div>
                        </template>
                    </el-table-column>
                    <el-table-column label="Age" prop="age" width="70" align="center" />
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
                                style="min-height: 350px" :theme="aceConfig.theme"
                                :options="aceConfig.options" /></el-tab-pane>
                        <el-tab-pane label="Json" name="json"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                                style="min-height: 350px" :theme="aceConfig.theme"
                                :options="aceConfig.options" /></el-tab-pane>
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

            <el-dialog v-model="dialogcreatens" title="创建资源" center>
                <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="100px" class="demo-ruleForm"
                    status-icon>
                    <el-form-item label="名称" prop="name" style="width: 80%;">
                        <el-input v-model="ruleForm.name"></el-input>
                    </el-form-item>
                    <el-form-item label="名称空间" prop="namespace" style="width: 80%;">
                        <el-select v-model="ruleForm.namespace" placeholder="请选择名称空间" @visible-change="getnsselect()">
                            <el-option v-for="item in nslist" :key="item.namespace" :label="item.label"
                                :value="item.namespace" v-show="item.label != 'All'" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="标签" prop="labels" style="width: 80%;">
                        <el-input v-model="ruleForm.labels"></el-input>
                    </el-form-item>
                    <el-form-item label="Class" prop="Class" style="width: 80%;">
                        <el-select v-model="ruleForm.ingress_class_name" placeholder="请选择Class">
                            <el-option v-for="item in ingressclasslist" :key="item.value" :label="item.label"
                                :value="item.value" />
                        </el-select>
                    </el-form-item>

                    <el-form-item v-for="(rules, index) in ruleForm.rules" :label="'Rules'" :key="rules.key">
                        <el-form-item style="width: 100%;" v-show="index == 0">
                            <el-tooltip placement="right"><template #content> 新增Rules </template>
                                <el-button @click="addRules" type="success" circle icon="Plus" size="small"></el-button>
                            </el-tooltip>
                            <el-tooltip placement="right"><template #content> 删除Rules </template>
                                <el-button @click.prevent="removeRules(rules)" type="danger" icon="Minus" size="small"
                                    circle v-show="isremove" />
                            </el-tooltip>
                        </el-form-item>
                        <el-form-item label="Host" prop="port_name" style="width: 77%;">
                            <el-input v-model="rules.host"></el-input>
                        </el-form-item>
                        <el-form-item v-for="(values, index) in rules.http_ingress_rule_values" :label="'HttpValue'"
                            :key="values.key">
                            <el-form-item label="Path" style="width: 73%;">
                                <el-input v-model="values.path"></el-input>
                            </el-form-item>
                            <el-form-item label="svcname" style="width: 73%;">
                                <el-input v-model="values.service_name"></el-input>
                            </el-form-item>
                            <el-form-item label="svcport" style="width: 73%;">
                                <el-input v-model.number="values.service_port"></el-input>
                            </el-form-item>
                        </el-form-item>

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

        </el-col>
    </el-row>
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
            ingressItem: [],
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
            isremove: false,
            isvalueremove: false,
            total: 0,
            page_size: [1, 10, 20, 50, 100],
            ingressclasslist: [
                {
                    "value": 'nginx',
                    "label": 'nginx'
                },
                {
                    "value": 'nginx-example',
                    "label": 'nginx-example'
                }
            ],
            ruleForm: {
                name: null,
                namespace: null,
                labels: {},
                ingress_class_name: null,
                rules: [
                    {
                        host: null,
                        http_ingress_rule_values: [
                            {
                                path: null,
                                service_name: null,
                                service_port: null
                            }
                        ]
                    }
                ],
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
        this.getIngress()
    },
    methods: {
        getIngress() {
            this.$ajax({
                method: 'get',
                url: '/ing/list',
                params: {
                    filter_name: this.filter_name,
                    namespace: this.namespace,
                    limit: this.limit,
                    page: this.page,
                }
            }).then((res) => {
                this.total = res.data.total
                this.ingressItem = res.data.item
            }).catch(function (res) {
                console.log(res);
            })
        },
        createIngress() {
            this.$ajax.post(
                '/ing/create',
                {
                    data: this.ruleForm
                },
            ).then((res) => {
                this.Loading("Creating")
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
            this.getIngress()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getIngress()
        },
        handleEdit(namespace, name) {
            this.$ajax({
                method: 'get',
                url: '/ing/detail',
                params: {
                    ingress_name: name,
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
        handleDelete(namespace, name) {
            this.$ajax({
                method: 'delete',
                url: '/ing/delete',
                params: {
                    namespace: namespace,
                    ingress_name: name,
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
                this.$message({
                    showClose: true,
                    message: res.reason,
                    type: 'error'
                });
            })

        },
        handleUpdate(namespace) {
            let data = this.content
            if (this.aceConfig.lang == 'yaml') {
                data = JSON.stringify(yaml.load(data), null, 2);
            }
            this.$ajax.put(
                '/ing/update',
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
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.reason,
                    type: 'error'
                });
            })
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
                    this.dialogcreatens = false
                    console.log(this.ruleForm)
                    this.createIngress()
                } else {
                    return false;
                }
            });
        },
        resetForm(formName) {
            this.$refs[formName].resetFields();
        },
        addRules() {
            this.ruleForm.rules.push({
                host: null,
                http_ingress_rule_values: [
                    {
                        path: null,
                        service_name: null,
                        service_port: null
                    }
                ]
            })
            this.isremove = true
        },
        removeRules(item) {
            var arr = new Array();
            var index = this.ruleForm.rules.indexOf(item)
            if (index !== -1) {
                this.ruleForm.rules.splice(index, 1)
            }
            arr = this.ruleForm.rules
            if (arr.length == 1) {
                this.isremove = false
            }


        },
        addValue(index) {
            this.ruleForm.rules[index].http_ingress_rule_values.push({
                path: null,
                service_name: null,
                service_port: null
            })
            this.isvalueremove = true
        },
        removeValue(item, i) {
            var arr = new Array();
            var index = this.ruleForm.rules[i].http_ingress_rule_values.indexOf(item)
            if (index !== -1) {
                this.ruleForm.rules[i].http_ingress_rule_values.splice(index, 1)
            }
            arr = this.ruleForm.rules[i].http_ingress_rule_values
            if (arr.length == 1) {
                this.isvalueremove = false
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
    padding: 10px;
    border-radius: 4px;
    min-height: 10px;
}

.el-row:last-child {
    margin-bottom: 0;
}

.grid-content1 {

    border-radius: 4px;
}

.grid-content2 {

    border-radius: 4px;
}

.grid-content3 {
    position: fixed;
    margin-top: 5px;
    margin-left: 20px;
}

.el-table,
.el-table__expanded-cell {
    margin-top: 5px;
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

.el-dialog--center {
    min-height: 400px;
    min-width: 200px;
}

.el-dialog__body {
    min-height: 400px;
}

.demo-pagination-block+.demo-pagination-block {
    float: right;
    margin-right: 10px;
    margin-top: 10px;
}

.grid-content3 .el-select {
    width: 100px;
}

.el-tabs--border-card>.el-tabs__content {
    padding: 0px;
}

/* .el-form-item__content .el-input {
    width: 500px;
} */
</style>