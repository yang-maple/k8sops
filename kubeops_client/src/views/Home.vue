
<template>
  <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
    <el-tab-pane label="Yaml" name="yaml"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
        style="min-height: 300px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
    <el-tab-pane label="Json" name="json"><v-ace-editor v-model:value="content" :lang="aceConfig.lang"
        style="min-height: 300px" :theme="aceConfig.theme" :options="aceConfig.options" /></el-tab-pane>
  </el-tabs>
</template>


<script>
import { VAceEditor } from 'vue3-ace-editor';
import '../ace.config.js';
import yaml from 'js-yaml';
export default {
  data() {
    return {
      content: '',
      aceConfig: {
        lang: 'json',
        theme: "cloud9_day",
        options: {
          showPrintMargin: false,
        }
      },
      activeName: 'json'
    };
  },
  created() {
    this.editorInit()
  },
  methods: {
    editorInit() {
      this.$ajax({
        method: 'get',
        url: '/pod/detail',
        params: {
          pod_name: "web-jenkins-1",
          namespace: "test"
        }
      }).then((res) => {
        this.content = JSON.stringify(res.data, null, 2)
      }).catch(function (res) {
        console.log(res.data);
      })
    },
    yamlFormat() {
      this.aceConfig.lang = "yaml"
      this.content = yaml.dump(JSON.parse(this.content))
    },
    jsonFormat() {
      this.aceConfig.lang = "json"
      this.content = JSON.stringify(yaml.load(this.content), null, 2);
    },
    handleClick(tab) {
      console.log(tab.props.name);
      if (tab.props.name == "yaml") {
        this.yamlFormat()
        return
      }
      this.jsonFormat()
    }

  },
  components: {
    VAceEditor,
  },
};
</script>