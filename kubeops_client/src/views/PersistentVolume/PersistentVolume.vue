<template>
    <el-row>
        <el-col :span="24">
            <div class="grid-content bg-purple">
                <el-table :data="persistentvolumeItem" :header-cell-style="{ background: '#e6e7e9' }" size="small">
                    <el-table-column label="Name" min-width="370">
                        <template #default="scope">
                            <div style="display: flex; align-items: center">
                                <span slot="reference" v-if="scope.row.status.phase == 'Bound'">
                                    <el-tooltip placement="bottom" effect="light"><template #content>Bound</template>
                                        <i class="dotClass" style="background-color: springgreen"></i></el-tooltip>
                                </span>
                                <span slot="reference" v-if="scope.row.status.phase != 'Bound'">
                                    <el-tooltip placement="bottom" effect="light"><template #content> Available </template>
                                        <i class="dotClass" style="background-color: red"></i></el-tooltip>
                                </span>
                                <el-link style="margin-left: 10px" type="primary" :underline="false"
                                    @click="handle(scope.row)">{{
                                        scope.row.name
                                    }}</el-link>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column label="Lables" width="100" align="center">
                        <template #default="scope">
                            <el-tag type="info" style="margin-left: 5px;" size="small" v-for="(v, k) in scope.row.labels "
                                :key="k">{{ k }}:{{ v
                                }}<br></el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="CAPACITY" prop="capacity.storage" width="100" align="center" />
                    <el-table-column label="ACCESS MODES" prop="access_mode" width="150" />
                    <el-table-column label="RECLAIM POLICY" prop="reclaim_policy" width="150" align="center" />
                    <el-table-column label="STATUS" prop="status.phase" width="120" />
                    <el-table-column label="CLAIM" prop="claim" width="220" />
                    <el-table-column label="STORAGECLASS" prop="storage_class" width="220" />
                    <el-table-column label="REASON" prop="reason" width="100" />
                    <el-table-column label="AGE" prop="age" width="80" />
                    <el-table-column fixed="right" label="Operations" width="150">
                        <template #default="scope">
                            <el-button type="primary" icon="Edit" plain
                                @click="dialogFormVisible = true, handleEdit(scope.$index, scope.row)" />
                            <el-popconfirm width="230" confirm-button-text="OK" cancel-button-text="No, Thanks"
                                icon="InfoFilled" icon-color="#626AEF" title="Are you sure to delete this?"
                                @confirm="handleDelete(scope.row)"><template #reference>
                                    <el-button type="danger" icon="Delete" plain />
                                </template>
                            </el-popconfirm>
                        </template>
                    </el-table-column>
                </el-table>
                <el-row>
                    <el-col :span="6">
                        <div class="grid-content ep-bg-purple">
                            <el-button type="primary" @click="dialogcreatens = true">
                                创建资源
                            </el-button>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <div class="grid-content ep-bg-purple-light" />
                    </el-col>
                    <el-col :span="6">
                        <div class="grid-content ep-bg-purple" />
                    </el-col>
                    <el-col :span="6">
                        <div class="grid-content3">
                            <el-pagination v-model:current-page="page" v-model:page-size="limit" :page-sizes="page_size"
                                small background layout="sizes, prev, pager, next" :total="total"
                                @size-change="handleSizeChange" @current-change="handleCurrentChange" />
                        </div>
                    </el-col>
                </el-row>
            </div>
            <el-dialog v-model="dialogFormVisible" title="实例详情/json" center>
                <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
                    <el-tab-pane label="Yaml" name="yaml"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                            style="min-height: 350px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
                    <el-tab-pane label="Json" name="json"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                            style="min-height: 350px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
                </el-tabs>
                <template #footer>
                    <span class="dialog-footer">
                        <el-button type="primary" @click="dialogFormVisible = false, handleUpdate()">更新</el-button>
                        <el-button @click="dialogFormVisible = false">
                            取消
                        </el-button>
                    </span>
                </template>
            </el-dialog>
            <el-dialog v-model="dialogcreatens" title="创建 PersistentVolume 资源" center>
                <el-form :model="form" :label-width="formLabelWidth">
                    <el-form-item label="名称">
                        <el-input v-model="form.name" autocomplete="off" />
                    </el-form-item>
                    <el-form-item label="容量">
                        <div class="mt-4">
                            <el-input v-model="form.storage" placeholder="Please input" class="input-with-select">
                                <template #append>
                                    <el-select v-model="form.storage_type" placeholder="Select" style="width: 120px">
                                        <el-option label="Gi" value="Gi" />
                                        <el-option label="Mi" value="Mi" />
                                        <el-option label="Tel" value="3" />
                                    </el-select>
                                </template>
                            </el-input>
                        </div>
                    </el-form-item>
                    <el-form-item label="类型">
                        <el-select v-model="form.type" placeholder="type">
                            <el-option label="NFS" value="NFS" @click="changeshow()" />
                            <el-option label="HostPATH" value="HostPATH" @click="changedisshow()" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="服务地址" v-model="showEl" v-if="showEl == true">
                        <el-input v-model="form.server" autocomplete="off" />
                    </el-form-item>
                    <el-form-item label="路径">
                        <el-input v-model="form.path" autocomplete="off" />
                    </el-form-item>
                </el-form>
                <template #footer>
                    <span class="dialog-footer">
                        <el-button type="primary" @click="dialogcreatens = false, createPersistentVolume()">确认创建</el-button>
                        <el-button type="danger" @click="dialogcreatens = false">
                            取消
                        </el-button>
                    </span>
                </template>
            </el-dialog>
        </el-col>
    </el-row>
</template>

<script scoped>
import { ElLoading } from 'element-plus'
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
            showEl: false,
            select: '',
            persistentvolumeItem: [],
            dialogFormVisible: false,
            dialogcreatens: false,
            form: {
                name: '',
                label: [],
                storage: '',
                storage_type: '',
                type: '',
                path: '',
                server: '',
            },
            formLabelWidth: '140px',
            total: 1,
            page_size: [1, 10, 20, 50, 100],
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
            activeName: 'json'
        }
    },
    created() {
        this.getPersistentVolume()
    },
    methods: {
        handleEdit(index, row) {
            this.$ajax({
                method: 'get',
                url: '/pv/detail',
                params: {
                    persistent_volume_name: row.name
                }
            }).then((res) => {
                this.activeName = 'json'
                this.aceConfig.lang = "json"
                this.content = JSON.stringify(res.data.detail, null, 2)
                console.log(res.data);
            }).catch((res) => {
                console.log(res.data);
            })
            console.log(index, row);
        },
        handleDelete(row) {
            this.$ajax({
                method: 'delete',
                url: '/pv/delete',
                params: {
                    persistent_volume_name: row.name
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
        getPersistentVolume() {
            this.$ajax({
                method: 'get',
                url: '/pv/list',
            }).then((res) => {
                this.persistentvolumeItem = res.data.item
            }).catch((res) => {
                console.log(res);
            })
        },
        createPersistentVolume() {
            this.$ajax.post(
                '/pv/create',
                {
                    name: this.form.name,
                    storage: this.form.storage + this.form.storage_type,
                    type: this.form.type,
                    path: this.form.path,
                    server: this.form.server
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
        handleUpdate() {
            this.Loading("Updateing")
            let data = this.content
            if (this.aceConfig.lang == 'yaml') {
                data = JSON.stringify(yaml.load(data), null, 2);
            }
            this.$ajax.put(
                '/pv/update',
                {
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
                    message: res.msg + res.reason,
                    type: 'error'
                });
                console.log(res);
            })
        },
        Loading(msg) {
            const loading = ElLoading.service({
                lock: true,
                text: msg,
                background: 'rgba(0, 0, 0, 0.7)',
            })
            setTimeout(() => {
                loading.close()
            }, 2000)
            this.reload()
        },
        changeshow() {
            this.showEl = true
        },
        changedisshow() {
            this.showEl = false
        },
        handleSizeChange(limit) {
            this.limit = limit
        },
        handleCurrentChange(page) {
            this.page = page
        },
        handle(row) {
            this.$router.push({
                path: '/persistentvolume/persistentvolume_detail',
                name: 'PersistentVolume 详情',
                query: {
                    persistent_volume_name: row.name
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

.el-input {
    width: 80%;
}

.dialog-footer button:first-child {
    margin-right: 10px;
}

.el-dialog--center .el-dialog__body {
    min-height: 200px;
}

.el-select .el-input {
    width: 120px;
}

.grid-content3 .el-select {
    width: 100px;
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
</style>