<template>
  <router-view v-if="isRouterAlive"></router-view>
</template>

<script>
const debounce = (fn, delay) => {
  let timer = null;

  return function () {
    let context = this;

    let args = arguments;

    clearTimeout(timer);

    timer = setTimeout(function () {
      fn.apply(context, args);
    }, delay);
  };
};
const _ResizeObserver = window.ResizeObserver;

window.ResizeObserver = class ResizeObserver extends _ResizeObserver {
  constructor(callback) {
    callback = debounce(callback, 16);

    super(callback);
  }
};
export default {
  name: 'App',
  provide() {
    return {
      reload: this.reload
    }
  },

  data() {
    return {
      isRouterAlive: true
    }
  },

  methods: {
    reload() {
      this.isRouterAlive = false;
      this.$nextTick(function () {
        this.isRouterAlive = true;
      });
    }
  },
}
</script>

<style>
.html,
body {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}

#nprogress .bar {
  background: #2186c0 !important;
}

/* 滚动条颜色样式 */
::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 0;
}

/* 滚动条整体部分 */
::-webkit-scrollbar {
  -webkit-appearance: none;
  width: 6px;
  height: 6px;
}

/*滚动条里面的小方块，能向上向下移动（或往左往右移动，取决于是垂直滚动条还是水平滚动条） */
::-webkit-scrollbar-thumb {
  cursor: pointer;
  border-radius: 5px;
  background: rgba(0, 0, 0, 0.15);
  transition: color 0.2s ease;
}
</style>
