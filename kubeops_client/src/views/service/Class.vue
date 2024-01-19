<template>
    <el-row style="padding-bottom: 10px;">
        <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
                <span>
                    <div>
                        <svg class="icon-ingclass" aria-hidden="true">
                            <use xlink:href="#icon-ingress_nginx"></use>
                        </svg>
                        <span
                            style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;">应用路由类</span>
                        <br>
                        <span
                            style="font-size: 12px;color: #79879c!important">应用路由类（IngressClass）一种用于控制Ingress资源的访问权限的机制</span>
                    </div>
                </span>
            </el-card>
        </el-col>
    </el-row>
    <el-row>
        <el-col :span="8">
            <div class="header-grid-content">
                <el-input v-model="filter_name" placeholder="请输入搜索资源的名称" clearable @clear="getIngress()"
                    @keyup.enter="getIngress()">
                    <template #prepend>
                        <el-button icon="Search" @click="getIngress()"></el-button>
                    </template>
                </el-input>
            </div>
        </el-col>
        <el-col :span="10">
            <div class="header-grid-content" style="min-height: 32px;"></div>
        </el-col>
        <el-col :span="6">
            <div class="header-grid-content" style="text-align: right;">
                <el-button type="info" @click="getIngress()" round>
                    <el-icon>
                        <Refresh />
                    </el-icon>
                </el-button>
            </div>
        </el-col>
    </el-row>
    <div class="table-bg-purple">
        <el-table :data="ingressClassItem" :header-cell-style="{ background: '#e6e7e9' }" style="width: 100%" size="small"
            empty-text="抱歉，暂无数据">
            <el-table-column label="名称" prop="name" width="200">
                <template #default="scope">
                    <div style="display: flex; align-items: center">
                        <svg class="icon-table-ingclass" aria-hidden="true">
                            <use xlink:href="#icon-ingress_nginx"></use>
                        </svg>
                        <el-text style="margin-left:3px" size="small">{{
                            scope.row.name
                        }}</el-text>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="控制器">
                <template #default="scope">{{ scope.row.controller ? scope.row.controller : '---' }}</template>
            </el-table-column>
            <el-table-column label="创建时间" prop="age" width="360"></el-table-column>
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

    <el-dialog v-model="dialogFormVisible" title="实例详情" center>
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


<script>
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
            ingressClassItem: [],
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
            activeName: 'json',
        }
    },
    created() {
        this.getIngress()
    },
    methods: {
        getIngress() {
            this.$ajax({
                method: 'get',
                url: '/ingressClass/list',
                params: {
                    filter_name: this.filter_name,
                    limit: this.limit,
                    page: this.page,
                }
            }).then((res) => {
                this.total = res.data.total
                this.ingressClassItem = res.data.item
            }).catch((res) => {
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
                console.log(res.data);
            })
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
        handleEdit(name) {
            this.$ajax({
                method: 'get',
                url: '/ingressClass/detail',
                params: {
                    Name: name,
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
                console.log(res);
            })
        },
        messageboxOperate(row, name) {
            this.$confirm(`是否${name}实例${row.name}`, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.handleDelete(row.name)
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                });
            });
        },
        handleDelete(name) {
            this.$ajax({
                method: 'delete',
                url: '/ingressClass/delete',
                params: {
                    name: name
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
                '/ingressClass/update',
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
                console.log(res);
                this.$message({
                    showClose: true,
                    message: res.msg,
                    type: 'error'
                });
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


.demo-pagination-block+.demo-pagination-block {
    float: right;
    margin-right: 10px;
    margin-top: 10px;
}

.icon-ingclass {
    width: 3em;
    height: 3em;
    vertical-align: -1em;
    fill: currentColor;
    overflow: hidden;
}

.icon-table-ingclass {
    width: 1.5em;
    height: 1.5em;
    vertical-align: -0.5em;
    fill: currentColor;
    overflow: hidden;
}
</style>