<template>
  <transition name="toast-fade">
    <div v-if="visible" class="toast-container" :class="type">
      <div class="toast-content">
        <span class="toast-icon" v-if="type === 'success'">✓</span>
        <span class="toast-icon" v-else-if="type === 'error'">✗</span>
        <span class="toast-icon" v-else-if="type === 'info'">ℹ</span>
        <span class="toast-message">{{ message }}</span>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';

const props = defineProps({
  message: {
    type: String,
    required: true
  },
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['success', 'error', 'info'].includes(value)
  },
  duration: {
    type: Number,
    default: 3000
  },
  onClose: {
    type: Function,
    default: () => {}
  }
});

const visible = ref(false);
let timer = null;

onMounted(() => {
  visible.value = true;
  
  if (props.duration > 0) {
    timer = setTimeout(() => {
      close();
    }, props.duration);
  }
});

onBeforeUnmount(() => {
  if (timer) {
    clearTimeout(timer);
  }
});

function close() {
  visible.value = false;
  setTimeout(() => {
    props.onClose();
  }, 300); // 等待过渡动画完成
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  padding: 12px 20px;
  border-radius: 4px;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.16);
  display: flex;
  align-items: center;
  min-width: 250px;
  max-width: 80%;
}

.toast-content {
  display: flex;
  align-items: center;
}

.toast-icon {
  margin-right: 10px;
  font-size: 18px;
  font-weight: bold;
}

.toast-message {
  font-size: 14px;
  word-break: break-word;
}

.success {
  background-color: #f0fff4;
  border-left: 4px solid #27ae60;
  color: #27ae60;
}

.error {
  background-color: #fff5f5;
  border-left: 4px solid #e53e3e;
  color: #e53e3e;
}

.info {
  background-color: #ebf8ff;
  border-left: 4px solid #3498db;
  color: #3498db;
}

.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.toast-fade-enter-from,
.toast-fade-leave-to {
  opacity: 0;
  transform: translate(-50%, -20px);
}
</style>