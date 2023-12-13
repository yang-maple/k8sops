<template>
    <div class="hello">
        <h1 style="text-align: center;">vue-codemirror-json-editor</h1>
        <codemirror v-model="jsonData" :options="editorOptions"></codemirror>
    </div>
</template>

<script>
import { codemirror } from "vue-codemirror";
import "codemirror/mode/javascript/javascript";
import "codemirror/addon/fold/foldcode";
import "codemirror/addon/fold/foldgutter";
import "codemirror/addon/fold/brace-fold";
import "codemirror/lib/codemirror.css";
import "codemirror/addon/fold/foldgutter.css";

export default {
    name: "HelloWorld",
    components: {
        codemirror
    },
    data() {
        return {
            jsonData: JSON.stringify({ a: 1, b: 2 }, null, "  "),
            editorOptions: {
                mode: "application/json",
                lineNumbers: true,
                line: true,
                lineWrapping: true,
                foldGutter: true,
                gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"]
            }
        };
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.vue-codemirror {
    border: 1px solid #ddd;
}
</style>
