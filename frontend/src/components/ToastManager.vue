<template>
  <div class="toast-manager">
    <Toast
      v-for="(toast, index) in toasts"
      :key="index"
      :message="toast.message"
      :type="toast.type"
      :duration="toast.duration"
      :onClose="() => removeToast(index)"
    />
  </div>
</template>

<script setup>
import { provide } from 'vue';
import Toast from './Toast.vue';
import { toast, toasts } from './toast.js';

// 从toast.js中导入removeToast函数的替代实现
function removeToast(index) {
  if (index >= 0 && index < toasts.value.length) {
    toasts.value.splice(index, 1);
  }
}

// 提供toast方法给其他组件使用
provide('toast', toast);
</script>

<!-- 使用普通script标签导出toast方法 -->
<script>
import { toast } from './toast.js';
export { toast };
</script>

<style scoped>
.toast-manager {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 0;
  pointer-events: none;
  z-index: 9999;
}
</style>