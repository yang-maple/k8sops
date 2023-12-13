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
                                <el-select v-model="namespace" placeholder="Select" @visible-change="getnsselect()"
                                    @change="getPersistentVolumeClaims()">
                                    <el-option-group v-for="group in nslist" :key="group.label" :label="group.label">
                                        <el-option v-for="item in group.options" :key="item.namespace" :label="item.label"
                                            :value="item.namespace" />
                                    </el-option-group>
                                </el-select>
                            </div>
                        </el-col>
                        <el-col :span="6">
                            <div class="grid-content1 ep-bg-purple">
                                <el-button type="primary" @click="dialogcreatens = true" style="margin-left: 18px;">
                                    创建 PersistentVolumeClaim 资源
                                </el-button>
                            </div>
                        </el-col>
                        <el-col :span="6"></el-col>
                        <el-col :span="6">
                            <div class="grid-content1 ep-bg-purple">
                                <el-input v-model="filter_name" placeholder="Please input" class="input-with-select"
                                    clearable>
                                    <template #prepend>
                                        <el-button icon="Search" @click="getPersistentVolumeClaims()" />
                                    </template>
                                </el-input>
                            </div>
                        </el-col>
                    </el-row>
                </el-col>
                <el-table :data="claimsItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%"
                    size="small">
                    <el-table-column label="Name" width="180" align="center">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <el-icon>
                                    <timer />
                                </el-icon>
                                <span style="margin-left: 10px">{{ scope.row.name }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column label="labels" width="150" align="center">
                        <template #default="scope">
                            <el-tag size="small" v-for="(v, k) in scope.row.labels " :key="k">{{ k }}:{{ v
                            }}<br></el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="Status" prop="status" width="80" align="center" />
                    <el-table-column label="Volume" prop="volume" width="300" align="center" />
                    <el-table-column label="Claim" prop="claim" width="100" align="center" />
                    <el-table-column label="AccessModes" width="110" align="center">
                        <template #default="scope">
                            <div v-for="(v, k) in scope.row.access_modes " :key="k">{{ v }}<br></div>
                        </template>
                    </el-table-column>
                    <el-table-column label="StorageClass" prop="storage_class" width="110" align="center" />
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
                <el-dialog v-model="dialogFormVisible" title="实例详情" center>
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
                    <el-col :span="4">
                        <div class="grid-content ep-bg-purple" />
                    </el-col>
                    <el-col :span="4">
                        <div class="grid-content ep-bg-purple-light" />
                    </el-col>
                    <el-col :span="4">
                        <div class="grid-content ep-bg-purple" />
                    </el-col>
                    <el-col :span="4">
                        <div class="grid-content ep-bg-purple-light" />
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
                        <el-select v-model="ruleForm.namespace" placeholder="Select" @visible-change="getnsselect()">
                            <el-option-group v-for="group in  nslist " :key="group.label">
                                <el-option v-for=" item  in  group.options " :key="item.namespace" :value="item.namespace"
                                    v-show="item.label != 'ALL'" />
                            </el-option-group>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="标签" prop="labels" style="width: 80%;">
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
                    <el-form-item label="存储容量" style="width: 80%;">
                        <div class="mt-4">
                            <el-input v-model="ruleForm.storage" placeholder="Please input" class="input-with-select">
                                <template #append>
                                    <el-select v-model="storage_type" style="width: 120px">
                                        <el-option label="Gi" value="Gi" />
                                        <el-option label="Mi" value="Mi" />
                                    </el-select>
                                </template>
                            </el-input>
                        </div>
                        <!-- <el-input v-model="ruleForm.storage"></el-input> -->
                    </el-form-item>
                    <el-form-item label="存储类型" style="width: 80%;">
                        <el-select v-model="ruleForm.storage_classname" placeholder="请选择存储类型">
                            <el-option v-for="item in storageclasslist" :key="item.value" :label="item.label"
                                :value="item.value" />
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
                labels: {},
                storage: null,
                storage_classname: null,
            },
            rules: {
                name: [
                    { required: true, message: '请输入资源名称', trigger: 'blur' },
                    { min: 3, max: 10, message: '长度在 3 到 10 个字符', trigger: 'blur' }
                ],
                namespace: [
                    { required: true, message: '请选择名称空间', trigger: 'change' }
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
                this.Loading("Deleting")
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

        },
        handleUpdate(namespace) {
            let data = this.content
            if (this.aceConfig.lang == 'yaml') {
                data = JSON.stringify(yaml.load(data), null, 2);
            }
            this.$ajax.put(
                '/pvc/update',
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
                    this.ruleForm.labels = JSON.parse(this.ruleForm.labels);
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

.grid-content2 {
    border-radius: 4px;
}

.grid-content3 {
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

.el-tabs--border-card>.el-tabs__content {
    padding: 0px;
}

/* .el-form-item__content .el-input {
    width: 500px;
} */
</style>