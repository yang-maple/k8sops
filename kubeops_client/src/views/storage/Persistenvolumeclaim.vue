<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-pvc" aria-hidden="true">
                            <use xlink:href="#icon-cunchujuan"></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">持久卷声明</span>
                        <br>
                        <span style="font-size: 12px;color: #79879c!important">持久卷声明（Persistent Volume Claim,
                            PVC）是Kubernetes中用于定义和管理应用程序对其它持久卷（Persistent Volume, PV）的请求和绑定的组件。</span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>
    <el-row>
        <el-col :span="4">
            <div class="header-grid-content">
                <el-select v-model="namespace" placeholder="命名空间（默认ALL）" @visible-change="getnsselect()"
                    @change="getPersistentVolumeClaims()">
                    <el-option-group v-for="group in nslist" :key="group.label" :label="group.label">
                        <el-option v-for="item in group.options" :key="item.namespace" :label="item.label"
                            :value="item.namespace" />
                    </el-option-group>
                </el-select>
            </div>
        </el-col>
        <el-col :span="16">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="请输入搜索资源的名称" clearable @clear="getPersistentVolumeClaims()"
                    @keyup.enter="getPersistentVolumeClaims()">
                    <template #prepend>
                        <el-button icon="Search" @click="getPersistentVolumeClaims()" />
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="4">
            <div class="header-grid-content" style="text-align: right;">
                <el-button type="info" @click="getPersistentVolumeClaims()" round>
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
        <el-table :data="claimsItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
            empty-text="抱歉，暂无数据">
            <el-table-column label="名称" width="300">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <el-icon>
                            <Coin />
                        </el-icon>
                        <span style="margin-left: 3px">{{ scope.row.name }}</span>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="容量" prop="claim" width="100" align="center" />
            <el-table-column label="持久卷" prop="volume" width="300">
                <template #default="scope">{{ scope.row.volume ? scope.row.volume : '---' }}</template>
            </el-table-column>
            <el-table-column label="状态" prop="status" width="100">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <span slot="reference" v-if="scope.row.status == 'Bound'">
                            <el-tooltip placement="bottom" effect="light"><template #content>Bound</template>
                                <i class="dotClass" style="background-color: springgreen"></i></el-tooltip>
                        </span>
                        <span slot="reference" v-if="scope.row.status != 'Bound'">
                            <el-tooltip placement="bottom" effect="light"><template #content> Available </template>
                                <i class="dotClass" style="background-color: red"></i></el-tooltip>
                        </span>
                        <el-text size="small" style="margin-left: 10px">{{
                            scope.row.status
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="访问模式" width="150">
                <template #default="scope">
                    <div v-for="(v, k) in scope.row.access_modes " :key="k">{{ v }}<br></div>
                </template>
            </el-table-column>
            <el-table-column label="存储类型" prop="storage_class">
                <template #default="scope">{{ scope.row.storage_class ? scope.row.storage_class : '---' }}</template>
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
                <div class="demo-pagination-block">
                    <el-pagination v-model:current-page="page" v-model:page-size="limit" :page-sizes="page_size" small
                        background layout="sizes, prev, pager, next" :total="total" :pager-count="5"
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
            <el-form-item label="命名空间" prop="namespace">
                <el-select v-model="ruleForm.namespace" placeholder="Select" @visible-change="getnsselect()">
                    <el-option-group v-for="group in  nslist " :key="group.label">
                        <el-option v-for=" item  in  group.options " :key="item.namespace" :value="item.namespace"
                            v-show="item.label != 'ALL'" />
                    </el-option-group>
                </el-select>
            </el-form-item>
            <el-form-item label="标签" prop="labels">
                <el-input v-model="ruleForm.labels"></el-input>
            </el-form-item>
            <el-form-item label="访问模式">
                <el-checkbox-group v-model="ruleForm.access_modes">
                    <el-checkbox label="ReadWriteOnce" name="access_modes" />
                    <el-checkbox label="ReadOnlyMany" name="access_modes" />
                    <el-checkbox label="ReadWriteMany" name="access_modes" />
                    <el-checkbox label="ReadWriteOncePod" name="access_modes" />
                </el-checkbox-group>
            </el-form-item>
            <el-form-item label="存储容量">
                <div>
                    <el-input v-model="ruleForm.storage" placeholder="Please input" style="width: fit-content;">
                        <template #append>
                            <el-select v-model="storage_type">
                                <el-option label="Gi" value="Gi" />
                                <el-option label="Mi" value="Mi" />
                            </el-select>
                        </template>
                    </el-input>
                </div>
            </el-form-item>
            <el-form-item label="存储类型">
                <el-select v-model="ruleForm.storage_classname" placeholder="请选择存储类型">
                    <el-option v-for="item in storageclasslist" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
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
            claimsItem: [],
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
                access_modes: [],
                labels: null,
                storage: null,
                storage_classname: null,
            },
            rules: {
                name: [
                    { required: true, message: '请输入资源名称', trigger: 'blur' },
                    { min: 3, max: 10, message: '长度在 3 到 10 个字符', trigger: 'blur' }
                ],
                namespace: [
                    { required: true, message: '请选择命名空间', trigger: 'change' }
                ],
                resource: [
                    { message: '请选择一种协议', trigger: 'change' }
                ],

            }
        }
    },
    created() {
        this.getPersistentVolumeClaims()
    },
    methods: {
        getPersistentVolumeClaims() {
            this.$ajax({
                method: 'get',
                url: '/pvc/list',
                params: {
                    filter_name: this.filter_name,
                    namespace: this.namespace,
                    limit: this.limit,
                    page: this.page,
                }
            }).then((res) => {
                this.total = res.data.total
                this.claimsItem = res.data.item
            }).catch((res) => {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
                console.log(res);
            })
        },
        createPersistentVolumeClaims() {
            this.$ajax.post(
                '/pvc/create',
                {
                    data: this.ruleForm
                },
            ).then((res) => {
                this.dialogcreatens = false
                this.$message({
                    message: res.msg,
                    type: 'success'
                });
                this.reload()
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
            this.getPersistentVolumeClaims()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getPersistentVolumeClaims()
        },
        handleEdit(namespace, name) {
            this.$ajax({
                method: 'get',
                url: '/pvc/detail',
                params: {
                    persistent_volume_claim_name: name,
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
                url: '/pvc/delete',
                params: {
                    namespace: namespace,
                    persistent_volume_claim_name: name
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
                '/pvc/update',
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
                console.log(res)
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
                    console.log(this.ruleForm)
                    if (this.ruleForm.labels != null) {
                        this.ruleForm.labels = JSON.parse(this.ruleForm.labels)
                    }
                    this.createPersistentVolumeClaims()
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

.el-dialog {
    .el-select {
        .el-input {
            width: 180px;
        }
    }
}

.el-form-item__content .el-input {
    width: 400px;
}

.dotClass {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    display: block;
    margin-left: 10px;
}

.icon-pvc {
    width: 2.5em;
    height: 2.5em;
    vertical-align: -0.7em;
    fill: currentColor;
    overflow: hidden;
}
</style>