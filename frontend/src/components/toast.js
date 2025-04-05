// toast.js - 提供全局toast通知方法
import { ref } from 'vue';

// 创建一个响应式数组来存储toast消息
const toasts = ref([]);

// 添加toast消息到数组
function addToast(message, type = 'info', duration = 3000) {
  const toast = {
    message,
    type,
    duration
  };
  toasts.value.push(toast);
  return toasts.value.length - 1;
}

// 从数组中移除toast消息
function removeToast(index) {
  if (index >= 0 && index < toasts.value.length) {
    toasts.value.splice(index, 1);
  }
}

// 导出toast方法和toasts数组
export const toast = {
  show(message, type = 'info', duration = 3000) {
    return addToast(message, type, duration);
  },
  success(message, duration = 3000) {
    return addToast(message, 'success', duration);
  },
  error(message, duration = 3000) {
    return addToast(message, 'error', duration);
  },
  info(message, duration = 3000) {
    return addToast(message, 'info', duration);
  }
};

export { toasts };