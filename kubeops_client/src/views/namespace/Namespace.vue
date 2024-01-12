<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-ns" aria-hidden="true">
                            <use xlink:href="#icon-namespace1
 "></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">命名空间</span>
                        <br>
                        <span
                            style="font-size: 12px;color: #79879c!important">命名空间（Namespace）提供了一种逻辑分隔的方式，可以将集群中的资源划分到不同的命名空间中，以实现不同用户或应用之间的资源隔离和管理。
                        </span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>

    <el-row>
        <el-col :span="8">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="Please input" clearable @clear="getNamespaces()"
                    @keyup.enter="getNamespaces()">
                    <template #prepend>
                        <el-button icon="Search" @click="getNamespaces()"></el-button>
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="10">
            <div class="header-grid-content" style="min-height: 32px;"></div>
        </el-col>
        <el-col :span="6">
            <div class="header-grid-content" style="text-align: center;">
                <el-button type="info" @click="getNamespaces()" round>
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
        <el-table :data="namespacesItem" :header-cell-style="{ background: '#e6e7e9' }" size="small">
            <el-table-column label="名称" width="200">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <el-text style="margin-left: 10px" size="small">{{
                            scope.row.name
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="标签">
                <template #default="scope">
                    <div v-for="(v, k, index) in scope.row.labels " :key="k">
                        <div v-if="index < maxitem[scope.$index]">
                            <el-tag type="info" style="margin-left: 5px;" size="small" effect="plain" round>
                                {{ k }}:{{ v }}</el-tag>
                        </div>
                    </div>
                    <div v-if="scope.row.labels == null" align="center">---</div>
                    <div v-if="scope.row.labels != null"><el-button size="small" type="primary" link
                            @click="showLabels(scope.$index)">{{
                                maxitem[scope.$index] == 3 ?
                                '展开' : '收起'
                            }}</el-button></div>
                </template>
            </el-table-column>
            <el-table-column label="状态" prop="status" width="300">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <span slot="reference" v-if="scope.row.status == 'Active'">
                            <el-tooltip placement="top" effect="light"><template #content>Active</template>
                                <i class="dotClass" style="background-color: springgreen"></i></el-tooltip>
                        </span>
                        <span slot="reference" v-if="scope.row.status != 'Active'">
                            <el-tooltip placement="bottom" effect="light"><template #content> UnActive </template>
                                <i class="dotClass" style="background-color: red"></i></el-tooltip>
                        </span>
                        <el-link style="margin-left: 10px" type="primary" :underline="false" @click="handle(scope.row)">{{
                            scope.row.status
                        }}</el-link>
                    </div>
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
    <el-dialog v-model="dialogcreatens" title="创建命名空间" center>
        <el-form :model="form" label-width="25%">
            <el-form-item label="名称">
                <el-input v-model="form.newnamespaces" autocomplete="off" style="width: 80%;" />
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="dialogcreatens = false, createNamespace()">确认创建</el-button>
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
            maxitem: [],
            namespacesItem: [],
            dialogFormVisible: false,
            dialogcreatens: false,
            form: {
                name: '',
                region: '',
                date1: '',
                date2: '',
                delivery: false,
                type: [],
                resource: '',
                desc: '',
                newnamespaces: '',
            },
            filter_name: '',
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
        this.getNamespaces()
    },
    methods: {
        handleEdit(name) {
            this.$ajax({
                method: 'get',
                url: '/namespaces/detail',
                params: {
                    namespace_name: name
                }
            }).then((res) => {
                this.activeName = 'json'
                this.aceConfig.lang = "json"
                this.content = JSON.stringify(res.data.detail, null, 2)
            }).catch((res) => {
                console.log(res.data);
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
                '/namespaces/update',
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
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.reason,
                    type: 'error'
                });
            })
            this.Refresh()
        },
        handleDelete(row) {
            this.$ajax({
                method: 'delete',
                url: '/namespaces/delete',
                params: {
                    namespace_name: row.name
                }
            }
            ).then((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'warning'
                });
            }).catch((res) => {
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
        getNamespaces() {
            this.$ajax({
                method: 'get',
                url: '/namespaces/list',
                params: {
                    page: this.page,
                    limit: this.limit,
                    filter_name: this.filter_name
                }
            }).then((res) => {
                this.total = res.data.total
                this.namespacesItem = res.data.item
                for (let i = 0; i < res.data.item.length; i++) {
                    this.maxitem.push(3)
                }
            }).catch((res) => {
                console.log(res);
            })
        },
        createNamespace() {
            this.$ajax.post(
                '/namespaces/create',
                {
                    namespace_name: this.form.newnamespaces
                },
            ).then((res) => {
                this.$message({
                    message: res.msg,
                    type: 'success'
                });
            }).catch((res) => {
                console.log(res);
            })
            this.Refresh()
        },
        Refresh() {
            setTimeout(() => {
                this.reload()
            }, 1000)
        },
        handleSizeChange(limit) {
            this.limit = limit
            this.page = 1
            this.getNamespaces()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getNamespaces()
        },
        handle(row) {
            this.$router.push({
                path: '/namespaces/namespaces_detail',
                name: 'namespaces详情',
                query: {
                    namespace_name: row.name
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
        },
        showLabels(index) {
            if (this.maxitem[index] == 3) {
                this.maxitem[index] = 99
            } else {
                this.maxitem[index] = 3
            }
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

.icon-ns {
    width: 2.5em;
    height: 2.5em;
    vertical-align: -0.7em;
    fill: currentColor;
    overflow: hidden;
}
</style>