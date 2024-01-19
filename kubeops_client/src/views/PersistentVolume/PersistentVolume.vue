<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-pv" aria-hidden="true">
                            <use xlink:href="#icon-cunchujuan1"></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">持久卷</span>
                        <br>
                        <span style="font-size: 12px;color: #79879c!important">持久卷（Persistent Volume,
                            PV）是Kubernetes中用于定义和管理集群中持久性存储资源的组件。</span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>

    <el-row>
        <el-col :span="8">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="请输入搜索资源的名称" clearable @clear="getPersistentVolume()"
                    @keyup.enter="getPersistentVolume()">
                    <template #prepend>
                        <el-button icon="Search"></el-button>
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="10">
            <div class="header-grid-content" style="min-height: 32px;"></div>
        </el-col>
        <el-col :span="6">
            <div class="header-grid-content" style="text-align: right;">
                <el-button type="info" @click="getPersistentVolume()" round>
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
        <el-table :data="persistentvolumeItem" :header-cell-style="{ background: '#e6e7e9' }" size="small"
            empty-text="抱歉，暂无数据">
            <el-table-column label="名称" width="300">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <svg class="icon-table-pv" aria-hidden="true">
                            <use xlink:href="#icon-cunchujuan1"></use>
                        </svg>
                        <el-text style="margin-left:3px" size="small">{{
                            scope.row.name
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="容量" prop="capacity.storage" width="100" align="center" />
            <el-table-column label="访问模式" prop="access_mode" width="120" />
            <el-table-column label="回收策略" prop="reclaim_policy" width="100" align="center" />
            <el-table-column label="状态" prop="status.phase" width="100">
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
                        <el-text size="small" style="margin-left: 10px">{{
                            scope.row.status.phase
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="绑定声明" prop="claim">
                <template #default="scope">{{ scope.row.claim ? scope.row.claim : '---' }}</template>
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
                                    @click="dialogFormVisible = true, handleEdit(scope.row.name)">编辑</el-dropdown-item>
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
    <el-dialog v-model="dialogFormVisible" title="实例详情/json" center>
        <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
            <el-tab-pane label="Yaml" name="yaml"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                    style="min-height: 350px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
            <el-tab-pane label="Json" name="json"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
                    style="min-height: 350px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
        </el-tabs>
        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="handleUpdate()">更新</el-button>
                <el-button @click="dialogFormVisible = false">
                    取消
                </el-button>
            </span>
        </template>
    </el-dialog>
    <el-dialog v-model="dialogcreatens" title="创建 PersistentVolume 资源" center>
        <el-form :model="form" label-width="25%">
            <el-form-item label="名称">
                <el-input v-model="form.name" autocomplete="off" />
            </el-form-item>
            <el-form-item label="标签">
                <el-input v-model="form.labels" placeholder="`请输入如下格式 {'a':'b'}`"></el-input>
            </el-form-item>
            <el-form-item label="容量">
                <el-input v-model="form.storage" style="width: fit-content;">
                    <template #append>
                        <el-select v-model="storage_type" placeholder="Select">
                            <el-option label="Gi" value="Gi" />
                            <el-option label="Mi" value="Mi" />
                            <el-option label="Tel" value="3" />
                        </el-select>
                    </template>
                </el-input>
            </el-form-item>
            <el-form-item label="访问模式">
                <el-checkbox-group v-model="form.access_mode">
                    <el-checkbox label="ReadWriteOnce" name="access_mode" />
                    <el-checkbox label="ReadOnlyMany" name="access_mode" />
                    <el-checkbox label="ReadWriteMany" name="access_mode" />
                    <el-checkbox label="ReadWriteOncePod" name="access_mode" />
                </el-checkbox-group>
            </el-form-item>
            <el-form-item label="卷模式">
                <el-radio-group v-model="form.volume_mode">
                    <el-radio label="Filesystem" />
                    <el-radio label="Block" />
                </el-radio-group>
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
            showEl: false,
            select: '',
            persistentvolumeItem: [],
            dialogFormVisible: false,
            dialogcreatens: false,
            storage_type: 'Gi',
            form: {
                name: '',
                labels: null,
                storage: '',
                type: '',
                path: '',
                server: '',
                access_mode: [],
                volume_mode: '',
            },
            formLabelWidth: '140px',
            total: 0,
            page_size: [1, 10, 20, 50, 100],
            limit: 10,
            page: 1,
            filter_name: '',
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
        handleEdit(name) {
            this.$ajax({
                method: 'get',
                url: '/pv/detail',
                params: {
                    persistent_volume_name: name
                }
            }).then((res) => {
                this.activeName = 'json'
                this.aceConfig.lang = "json"
                this.content = JSON.stringify(res.data.detail, null, 2)
                console.log(res.data);
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
                console.log(res.data);
            })
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
        messageboxOperate(row, name) {
            this.$confirm(`是否${name}实例${row.name}`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.handleDelete(row)
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                });
            });
        },
        getPersistentVolume() {
            this.$ajax({
                method: 'get',
                url: '/pv/list',
                params: {
                    page: this.page,
                    limit: this.limit,
                    filter_name: this.filter_name
                }
            }).then((res) => {
                this.total = res.data.total
                this.persistentvolumeItem = res.data.item
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
                console.log(res);
            })
        },
        createPersistentVolume() {
            if (this.form.labels != null) {
                this.form.labels = JSON.parse(this.form.labels)
            }
            this.form.storage = this.form.storage + this.storage_type
            console.log(this.form)
            this.$ajax.post(
                '/pv/create',
                {
                    data: this.form
                },
            ).then((res) => {
                this.$message({
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
                console.log(res);
            })


        },
        handleUpdate() {
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
                '/pv/update',
                {
                    data: datajson
                },
            ).then((res) => {
                this.dialogFormVisible = false,
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
                console.log(res);
            })
        },
        Refresh() {
            setTimeout(() => {
                this.reload()
            }, 1000)

        },
        changeshow() {
            this.showEl = true
        },
        changedisshow() {
            this.showEl = false
        },
        handleSizeChange(limit) {
            this.limit = limit
            this.page = 1
            this.getPersistentVolume()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getPersistentVolume()
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


.demo-pagination-block {
    margin-top: 5px;
    margin-bottom: 5px;
    padding-left: 10px;
}

.el-button--text {
    margin-right: 15px;
}

.dialog-footer button:first-child {
    margin-right: 10px;
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

.icon-pv {
    width: 2.5em;
    height: 2.5em;
    vertical-align: -0.7em;
    fill: currentColor;
    overflow: hidden;
}

.icon-table-pv {
    width: 1.5em;
    height: 1.5em;
    vertical-align: -0.5em;
    fill: currentColor;
    overflow: hidden;
}
</style>