
<!-- <template>
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
      }).catch((res) => {
        console.log(res);
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
</script> -->

<template>
  <!-- <el-upload class="upload-demo" ref="upload" action="http://127.0.0.1:8090/v1/api/upload/uploadFile"
    :on-preview="handlePreview" :on-remove="handleRemove" :file-list="upload" :auto-upload="false" :headers="config">
    <el-button slot="trigger" size="small" type="primary">选取文件</el-button>
    <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">上传到服务器</el-button>
    <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div>
  </el-upload> -->
  <el-upload ref="upload" action="http://127.0.0.1:8090/v1/api/upload/uploadFile" :auto-upload="false" :headers="config"
    :name="file" multiple :file-list="upload">
    <template #trigger>
      <el-button type="primary">select file</el-button>
    </template>

    <el-button class="ml-3" type="success" @click="submitUpload">
      upload to server
    </el-button>

    <template #tip>
      <div class="el-upload__tip">
        jpg/png files with a size less than 500kb
      </div>
    </template>
  </el-upload>
</template>

<script>
export default {
  data() {
    return {
      token: "",
      upload: []
    };
  },
  computed: {
    config() {
      return {
        Authorization: this.token,
      };
    },
  },
  mounted() {
    this.token = localStorage.getItem('user_token')
  },
  methods: {
    submitUpload() {
      this.$refs.upload.submit();
    },
    handleRemove(file, fileList) {
      console.log(file, fileList);
    },
    handlePreview(file) {
      console.log(file);
    }
  }
}
</script>