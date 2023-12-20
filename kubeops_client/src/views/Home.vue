
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
  <div class="header-collapse" @click="onCollapse">
    <!---->
    <el-icon>
      <component :is="isCollapse?'Expand':'Fold'" />
    </el-icon>
  </div>
  <el-menu default-active="2" class="el-menu-vertical-demo" :collapse="isCollapse" @open="handleOpen"
    @close="handleClose">
    <el-sub-menu index="1">
      <template #title>
        <el-icon>
          <location />
        </el-icon>
        <span>Navigator One</span>
      </template>
      <el-menu-item-group>
        <template #title><span>Group One</span></template>
        <el-menu-item index="1-1">item one</el-menu-item>
        <el-menu-item index="1-2">item two</el-menu-item>
      </el-menu-item-group>
      <el-menu-item-group title="Group Two">
        <el-menu-item index="1-3">item three</el-menu-item>
      </el-menu-item-group>
      <el-sub-menu index="1-4">
        <template #title><span>item four</span></template>
        <el-menu-item index="1-4-1">item one</el-menu-item>
      </el-sub-menu>
    </el-sub-menu>
    <el-menu-item index="2">
      <el-icon><icon-menu /></el-icon>
      <template #title>Navigator Two</template>
    </el-menu-item>
    <el-menu-item index="3" disabled>
      <el-icon>
        <document />
      </el-icon>
      <template #title>Navigator Three</template>
    </el-menu-item>
    <el-menu-item index="4">
      <el-icon>
        <setting />
      </el-icon>
      <template #title>Navigator Four</template>
    </el-menu-item>
  </el-menu>
</template>

<script  >
export default {
  data() {
    return {
      isCollapse: true
    };
  },
  methods: {
    handleOpen(key, keyPath) {
      console.log(key, keyPath);
    },
    handleClose(key, keyPath) {
      console.log(key, keyPath);
    },
    onCollapse() {
      if (this.isCollapse == true) {
        this.isCollapse = false
        return
      }
      this.isCollapse = true
    }
  }
}
</script>

<style>
.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 200px;
  min-height: 400px;
}
</style>