<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-configmap" aria-hidden="true">
                            <use xlink:href="#icon-configmap"></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">配置字典</span>
                        <br>
                        <span
                            style="font-size: 12px;color: #79879c!important">配置字典（ConfigMap）常用于存储工作负载所需的配置信息，许多应用程序会从配置文件、命令行参数或环境变量中读取配置信息。</span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>
    <el-row>
        <el-col :span="4">
            <div class="header-grid-content">
                <el-select v-model="namespace" placeholder="命名空间（默认ALL）" @change="getConfigmap()">
                    <el-option-group v-for="group in nslist" :key="group.label" :label="group.label">
                        <el-option v-for="item in group.options" :key="item.namespace" :label="item.label"
                            :value="item.namespace" />
                    </el-option-group>
                </el-select>
            </div>
        </el-col>
        <el-col :span="16">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="请输入搜索资源的名称" clearable @clear="getConfigmap()"
                    @keyup.enter="getConfigmap()">
                    <template #prepend>
                        <el-button icon="Search" @click="getConfigmap()" />
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="4">
            <div class="header-grid-content" style="text-align: right;">
                <el-button type="info" @click="getConfigmap()" round>
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
        <el-table :data="configmapItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
            empty-text="抱歉，暂无数据">
            <el-table-column label="名称" width="400">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <svg class="icon-table-cm" aria-hidden="true">
                            <use xlink:href="#icon-configmap"></use>
                        </svg>
                        <el-text style="margin-left:3px" size="small">{{
                            scope.row.name
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="命名空间" prop="namespace" width="200" />
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
            <el-table-column label="修改" prop="modify" width="50" align="center">
                <template #default="scope">
                    <span slot="reference" v-if="scope.row.modify != true">
                        <el-tooltip placement="top"><template #content> 支持修改 </template>
                            <i class="dotClass" style="background-color: springgreen"></i></el-tooltip>
                    </span>
                    <span slot="reference" v-if="scope.row.modify == true">
                        <el-tooltip placement="top"><template #content> 禁止修改 </template>
                            <i class="dotClass" style="background-color: red"></i></el-tooltip>
                    </span>
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

    <el-dialog v-model="dialogcreatens" title="创建资源" center class="create_config_dialog">
        <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="15%">
            <el-form-item label="资源名称" prop="name">
                <el-input v-model="ruleForm.name"></el-input>
            </el-form-item>
            <el-form-item label="命名空间" prop="namespace">
                <el-select v-model="ruleForm.namespace" placeholder=" ">
                    <el-option v-for="item in nslist[1].options" :key="item.namespace" :label="item.label"
                        :value="item.namespace" />
                </el-select>
            </el-form-item>
            <el-form-item label="数据字典">
                <div class="data" v-for="(value, key) in dataConfig" :key="key">
                    <el-row :gutter="2" style="margin-bottom: 10px">
                        <el-col :span="1">
                            <el-text>键</el-text>
                        </el-col>
                        <el-col :span="8">
                            <el-input v-model="value.key" :disabled="value.status" />
                        </el-col>
                        <el-col :span="1">
                            <el-text>值</el-text>
                        </el-col>
                        <el-col :span="14">
                            <el-input type="textarea" autosize :disabled="value.status" v-model="value.value"
                                style="width: 120%" />
                        </el-col>
                        <el-col :span="2">
                            <el-button type="primary" size="small" @click="value.status = !value.status">
                                <template #default="scope">
                                    <el-icon :size="12">
                                        <component :is="value.status?'Edit':'Check'" />
                                    </el-icon>
                                </template>
                            </el-button>
                        </el-col>
                        <el-col :span="2" style="margin-left: 5px">
                            <el-button size="small" type="danger" @click="dataConfig.splice(key, 1)">
                                <template #default="scope">
                                    <el-icon :size="12">
                                        <Delete />
                                    </el-icon>
                                </template>
                            </el-button>
                        </el-col>
                    </el-row>
                </div>
                <el-button type="info" style="width: 90%;" @click="addData">添加</el-button>
            </el-form-item>
            <el-form-item label="禁止修改">
                <el-switch v-model="ruleForm.modify" />
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
            maxitem: [],
            detailnamespace: null,
            dialogFormVisible: false,
            dialogcreatens: false,
            configmapItem: [],
            dataKey: 3,
            dataConfig: [],
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
            selectns: [],
            total: 0,
            storage_type: 'Gi',
            page_size: [1, 10, 20, 50, 100],
            storageclasslist: [
                {
                    value: 'standard',
                    label: 'standard'
                },
                {
                    value: 'nfs-client',
                    label: 'nfs-client'
                },
            ],
            ruleForm: {
                name: null,
                namespace: '',
                data: null,
                modify: false,
            },
            rules: {
                name: [
                    { required: true, message: '请输入资源名称', trigger: 'blur' },
                ],
                namespace: [
                    { required: true, message: '请选择命名空间', trigger: 'change' }
                ],
            }
        }
    },
    created() {
        this.getConfigmap()
        this.getnsselect()
    },
    methods: {
        getConfigmap() {
            this.$ajax({
                method: 'get',
                url: '/cm/list',
                params: {
                    filter_name: this.filter_name,
                    namespace: this.namespace,
                    limit: this.limit,
                    page: this.page,
                }
            }).then((res) => {
                this.total = res.data.total
                this.configmapItem = res.data.item
                for (let i = 0; i < res.data.item.length; i++) {
                    this.maxitem.push(3)
                }
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
                console.log(res);
            })
        },
        createConfigmap() {
            this.$ajax.post(
                '/cm/create',
                {
                    data: this.ruleForm
                },
            ).then((res) => {
                this.dialogcreatens = false
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
            this.getConfigmap()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getConfigmap()
        },
        handleEdit(namespace, name) {
            this.$ajax({
                method: 'get',
                url: '/cm/detail',
                params: {
                    configmap_name: name,
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
                url: '/cm/delete',
                params: {
                    namespace: namespace,
                    configmap_name: name
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
                '/cm/update',
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
                    if (this.dataConfig.length != 0) {
                        let datas = {};
                        this.dataConfig.forEach(v => {
                            datas[v.key] = v.value
                        })
                        this.ruleForm.data = datas
                    }
                    console.log(this.ruleForm)
                    this.createConfigmap()
                } else {
                    return false;
                }
            });
        },
        resetForm(formName) {
            this.$refs[formName].resetFields();
            this.dataConfig = [];
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
        },
        addData() {
            this.dataConfig.push({
                key: '',
                value: '',
                status: false,
            })
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


.dotClass {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    display: block;
    margin-left: 10px;
}

.el-tabs--border-card>.el-tabs__content {
    padding: 0px;
}

.create_config_dialog {
    .el-form-item__content .el-input {
        width: 550px;
    }

    .data {
        .el-input {
            width: 100%;
        }
    }
}


.icon-configmap {
    width: 2.5em;
    height: 2.5em;
    color: #8a8a8a;
    vertical-align: -0.7em;
    fill: currentColor;
    overflow: hidden;
}

.icon-table-cm {
    width: 1.5em;
    height: 1.5em;
    color: #8a8a8a;
    vertical-align: -0.5em;
    fill: currentColor;
    overflow: hidden;
}
</style>