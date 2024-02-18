<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-ing" aria-hidden="true">
                            <use xlink:href="#icon-ingress"></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">应用路由</span>
                        <br>
                        <span style="font-size: 12px;color: #79879c!important">应用路由（Ingress）提供一种聚合服务的方式，您可以通过一个外部可访问的 IP
                            地址将集群的内部服务暴露给外部。</span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>
    <el-row>
        <el-col :span="4">
            <div class="header-grid-content">
                <el-select v-model="namespace" filterable placeholder="命名空间（默认ALL）" @change="getIngress()" clearable>
                    <el-option-group v-for="group in nslist" :key="group.label" :label="group.label">
                        <el-option v-for="item in group.options" :key="item.namespace" :label="item.label"
                            :value="item.namespace" />
                    </el-option-group>
                </el-select>
            </div>
        </el-col>
        <el-col :span="16">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="请输入搜索资源的名称" clearable @clear="getIngress()"
                    @keyup.enter="getIngress()">
                    <template #prepend>
                        <el-button icon="Search" @click="getIngress()" />
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="4">
            <div class="header-grid-content" style="text-align: right;">
                <el-button type="info" @click="getIngress()" round>
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
        <el-table :data="ingressItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
            empty-text="抱歉，暂无数据">
            <el-table-column label="名称" width="300">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <svg class="icon-table-ing" aria-hidden="true">
                            <use xlink:href="#icon-ingress"></use>
                        </svg>
                        <el-text style="margin-left:3px" size="small">{{
                            scope.row.name
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="命名空间" prop="namespace" width="100" />
            <el-table-column label="端点" width="200">
                <template #default="scope">
                    <div v-for="(v, k) in scope.row.endpoint " :key="k">{{ v.ip ? v.ip : '---' }}<br></div>
                </template>
            </el-table-column>
            <el-table-column label="路由">
                <template #default="scope">
                    <div v-for="(v, k) in scope.row.host " :key="k">{{ v ? v : '---' }}<br></div>
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
                                    @click="dialogFormVisible = true, handleEdit(scope.row.namespace, scope.row.name)">编辑</el-dropdown-item>
                                <el-dropdown-item icon="Delete"
                                    @click="messageboxOperate(scope.row, '删除')">删除</el-dropdown-item>
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
                    <el-button type="primary" @click="handleUpdate(detailnamespace)">更新</el-button>
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
                    <el-pagination v-model:current-page="page" v-model:page-size="limit" :page-sizes="page_size" small
                        background layout="sizes, prev, pager, next" :total="total" :pager-count="5"
                        @size-change="handleSizeChange" @current-change="handleCurrentChange" />
                </div>
            </el-col>
        </el-row>
    </div>

    <el-dialog v-model="dialogcreatens" title="创建资源" center>
        <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="25%">
            <el-form-item label="资源名称" prop="name">
                <el-input v-model="ruleForm.name"></el-input>
            </el-form-item>
            <el-form-item label="命名空间" prop="namespace">
                <el-select v-model="ruleForm.namespace" placeholder="请选择命名空间" fit-input-width>
                    <el-option v-for="item in nslist[1].options" :key="item.namespace" :label="item.label"
                        :value="item.namespace" />
                </el-select>
            </el-form-item>
            <el-form-item label="资源标签" prop="labels">
                <el-input v-model="ruleForm.labels" placeholder="`请输入如下格式 {'a':'b'}`"></el-input>
            </el-form-item>
            <el-form-item label="控制器" prop="ingress_class_name">
                <el-select v-model="ruleForm.ingress_class_name" placeholder="请选择路由控制器">
                    <el-option v-for="item in ingressclasslist" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
            <el-form-item label="访问域名" prop="port_name">
                <el-input v-model="ruleForm.rules[0].host"></el-input>
            </el-form-item>
            <el-form-item label="访问路径" prop="rules[0].http_ingress_rule_values[0].path">
                <el-input v-model="ruleForm.rules[0].http_ingress_rule_values[0].path"></el-input>
            </el-form-item>
            <el-form-item label="路径类型" prop="rules[0].http_ingress_rule_values[0].path_type">
                <el-select v-model="ruleForm.rules[0].http_ingress_rule_values[0].path_type" placeholder="请选择路径类型">
                    <el-option v-for="item in path_type" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
            <el-form-item label="服务名称" prop="rules[0].http_ingress_rule_values[0].service_name">
                <el-input v-model="ruleForm.rules[0].http_ingress_rule_values[0].service_name"></el-input>
            </el-form-item>
            <el-form-item label="服务端口" prop="rules[0].http_ingress_rule_values[0].service_port">
                <el-input v-model.number="ruleForm.rules[0].http_ingress_rule_values[0].service_port"></el-input>
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
            path_type: [
                {
                    "value": 'ImplementationSpecific',
                    "label": 'ImplementationSpecific'
                },
                {
                    "value": 'Exact',
                    "label": 'Exact'
                },
                {
                    "value": 'Prefix',
                    "label": 'Prefix'
                },
            ],
            ruleForm: {
                name: null,
                namespace: null,
                labels: null,
                ingress_class_name: null,
                rules: [
                    {
                        host: null,
                        http_ingress_rule_values: [
                            {
                                path: null,
                                path_type: null,
                                service_name: null,
                                service_port: null,
                            }
                        ]
                    }
                ],
            },
            rules: {
                name: [
                    { required: true, message: '请输入资源名称', trigger: 'change' },
                ],
                namespace: [
                    { required: true, message: '请选择命名空间', trigger: 'change' }
                ],
                ingress_class_name: [
                    { required: true, message: '请选择路由类型', trigger: 'change' }
                ],
                resource: [
                    { message: '请选择一种协议', trigger: 'change' }
                ],
                type: [
                    { required: true, message: '请选择一种服务类型', trigger: 'change' }
                ],
                'rules[0].http_ingress_rule_values[0].path': [
                    { required: true, message: '请输入服务访问路径', trigger: 'change' }
                ],
                'rules[0].http_ingress_rule_values[0].path_type': [
                    { required: true, message: '请选择路径类型', trigger: 'change' }
                ],
                'rules[0].http_ingress_rule_values[0].service_name': [
                    { required: true, message: '请输入服务名称', trigger: 'change' }
                ],
                'rules[0].http_ingress_rule_values[0].service_port': [
                    { required: true, message: '请输入服务访问端口', trigger: 'change' },
                    { type: 'number', min: 1, max: 65535, message: '请输入1-65535数字' }
                ]
            }
        }
    },
    created() {
        this.getIngress()
        this.getnsselect()
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

            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
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
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'success'
                });
                this.dialogcreatens = false
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
            })
            this.reload()
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
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
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
            })
            this.reload()
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
                '/ing/update',
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
                this.reload()
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
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
                    try {
                        if (this.ruleForm.labels != null) {
                            this.ruleForm.labels = JSON.parse(this.ruleForm.labels)
                        }
                        this.createIngress()
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

.icon-ing {
    width: 2.5em;
    height: 2.5em;
    vertical-align: -0.8em;
    fill: currentColor;
    overflow: hidden;
}

.icon-table-ing {
    width: 1.5em;
    height: 1.5em;
    vertical-align: -0.5em;
    fill: currentColor;
    overflow: hidden;
}
</style>