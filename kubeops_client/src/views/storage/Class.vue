<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-stclass" aria-hidden="true">
                            <use xlink:href="#icon-cunchuleixingpeizhi1
"></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">存储类</span>
                        <br>
                        <span
                            style="font-size: 12px;color: #79879c!important">存储类（StorageClass）支持动态卷供应，使管理员能够按需创建新的卷。</span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>
    <el-row>
        <el-col :span="8">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="请输入搜索资源的名称" clearable @clear="getstorageClass()"
                    @keyup.enter="getstorageClass()">
                    <template #prepend>
                        <el-button icon="Search" @click="getstorageClass()"></el-button>
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="10">
            <div class="header-grid-content" style="min-height: 32px;"></div>
        </el-col>
        <el-col :span="6">
            <div class="header-grid-content" style="text-align: right;">
                <el-button type="info" @click="getstorageClass()" round>
                    <el-icon>
                        <Refresh />
                    </el-icon>
                </el-button>
            </div>
        </el-col>
    </el-row>
    <div class="table-bg-purple">
        <el-table :data="storageClassItem" :header-cell-style="{ background: '#e6e7e9' }" size="small" empty-text="抱歉，暂无数据">
            <el-table-column label="名称" prop="name" width="200">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <svg class="icon-table-storage" aria-hidden="true">
                            <use xlink:href="#icon-cunchuleixingpeizhi1"></use>
                        </svg>
                        <el-text style="margin-left:3px" size="small">{{
                            scope.row.name
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="提供者" prop="provisioner" width="280" />
            <el-table-column label="参数">
                <template #default="scope">
                    <div v-for="(v, k, index) in scope.row.parameters " :key="k">
                        <div class="content">
                            <el-tag type="info" style="margin-left: 5px;" effect="plain" round>
                                {{ k }}:{{ v }}<br></el-tag>
                        </div>
                    </div>
                    <div v-if="scope.row.parameters == null"> --- </div>
                </template>
            </el-table-column>
            <el-table-column label="回收策略" prop="reclaim_policy" width="150" />
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
                                <el-dropdown-item icon="Edit" @click="handleEdit(scope.row.name)">编辑</el-dropdown-item>
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
            storageClassItem: [],
            dialogFormVisible: false,
            filter_name: '',
            total: 0,
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
        this.getstorageClass()
    },
    methods: {
        handleEdit(name) {
            this.dialogFormVisible = true
            this.$ajax({
                method: 'get',
                url: '/storageClass/detail',
                params: {
                    Name: name
                }
            }).then((res) => {
                this.activeName = 'json'
                this.aceConfig.lang = "json"
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
                '/storageClass/update',
                {
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
        handleDelete(row) {
            this.$ajax({
                method: 'delete',
                url: '/storageClass/delete',
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
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
                console.log(res);
            })
            this.reload()
        },
        getstorageClass() {
            this.$ajax({
                method: 'get',
                url: '/storageClass/list',
                params: {
                    filter_name: this.filter_name,
                    limit: this.limit,
                    page: this.page,
                }
            }).then((res) => {
                this.total = res.data.total
                this.storageClassItem = res.data.item
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
        handleSizeChange(limit) {
            this.limit = limit
            this.page = 1
            this.getstorageClass()
        },
        handleCurrentChange(page) {
            this.page = page
            this.getstorageClass()
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

.icon-stclass {
    width: 2.5em;
    height: 2.5em;
    vertical-align: -0.7em;
    fill: currentColor;
    overflow: hidden;
}

.icon-table-storage {
    width: 1.5em;
    height: 1.5em;
    vertical-align: -0.5em;
    fill: currentColor;
    overflow: hidden;
}
</style>