<script setup>
import Top from "./components/Top.vue";
import { RouterView } from "vue-router";
import ToastManager from "./components/ToastManager.vue";
import { toast } from "./components/toast.js";

// 将toast方法挂载到window对象，使其可以全局访问
window.toast = toast;
</script>

<template>
  <div class="layout">
    <!-- 导航栏 -->
    <div class="navbar"  >
      <Top />
    </div>

    <!-- 主体区域 -->
    <div class="main">
      <RouterView />
    </div>
    
    <!-- Toast通知管理器 -->
    <ToastManager />
  </div>
</template>

<style>
/* 布局样式 */
.layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden; /* 防止外层出现滚动条 */
}

/* 导航栏样式 */
.navbar {
  position: fixed; /* 固定在顶部 */
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000; /* 确保在最上层 */
  background-color: #fff; /* 添加背景色防止内容透视 */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 可选：添加阴影 */
  text-align: center;
}

/* 主体区域样式 */
.main {
  flex: 1;
  padding: 20px;
  background-color: #f4f4f4;
  margin-top: 32px; /* 预留导航栏高度，需要根据实际Top组件高度调整 */
  overflow-y: auto; /* 主体区域内容超出时显示滚动条 */
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 8px; /* 滚动条宽度 */
  height: 8px; /* 滚动条高度 */
}

::-webkit-scrollbar-track {
  background: #f1f1f1; /* 滚动条轨道背景色 */
  border-radius: 10px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1; /* 滚动条滑块颜色 */
  border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8; /* 鼠标悬停时滑块颜色 */
}
</style>
