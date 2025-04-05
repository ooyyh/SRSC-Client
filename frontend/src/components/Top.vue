<template>
  <div class="navbar" style="--wails-draggable:drag">
    <!-- Logo 和标题 -->
    <div class="logo-container">
      <div class="logo">
        <img src="../assets/logo.png" alt="Logo" />
      </div>
      <span class="title">SRSC-Client</span>
    </div>

    <!-- 控制按钮 -->
    <div class="controls">
      <!-- 新增设置按钮 -->
      <button @click="openSettingsDialog" class="control-button settings-button">
        <img src="../assets/setting.png" alt="Settings" />
      </button>
      <button @click="minimize" class="control-button">
        <img src="../assets/minimize.png" alt="Minimize" />
      </button>
      <button @click="maximize" class="control-button">
        <img src="../assets/maximize.png" alt="Maximize" />
      </button>
      <button @click="close" class="control-button">
        <img src="../assets/close.png" alt="Close" />
      </button>
    </div>
  </div>

  <!-- 设置弹窗 -->
  <div v-if="showSettingsDialog" class="settings-dialog-overlay">
    <div class="settings-dialog">
      <h3>设置 Endpoint</h3>
      <div class="form-group">
        <label for="endpoint">Endpoint URL:</label>
        <input type="text" id="endpoint" v-model="endpointUrl" placeholder="例如: http://localhost:8080" />
      </div>
      <div class="dialog-buttons">
        <button @click="saveSettings" class="dialog-button save">保存</button>
        <button @click="closeSettingsDialog" class="dialog-button cancel">取消</button>
      </div>
    </div>
  </div>
  </template>
  
<script setup>
import { ref } from 'vue';
import { WindowFullscreen, WindowUnfullscreen, WindowMinimise, Quit } from '../../wailsjs/runtime'; // 导入 Wails Runtime 方法, 添加 Quit

const showSettingsDialog = ref(false);
const endpointUrl = ref(''); // 可以设置一个默认值，例如从 localStorage 读取
const isFullscreen = ref(false); // 添加状态来跟踪全屏

const openSettingsDialog = () => {
  // 可以在打开时从 localStorage 或其他地方加载当前设置
  // endpointUrl.value = localStorage.getItem('endpointUrl') || 'http://localhost:8080';
  showSettingsDialog.value = true;
};

const closeSettingsDialog = () => {
  showSettingsDialog.value = false;
};

const saveSettings = () => {
  console.log("保存 Endpoint:", endpointUrl.value);
  SetEndpoint(endpointUrl.value).then(() => {
    console.log("Endpoint 设置成功:", endpointUrl.value);
  }).catch((error) => {
    console.error("设置 Endpoint 失败:", error);
  });
  // 在这里添加保存设置的逻辑，例如保存到 localStorage
  // localStorage.setItem('endpointUrl', endpointUrl.value);
  // alert(`Endpoint 已保存: ${endpointUrl.value}`); // 简单提示
  closeSettingsDialog();
};

const minimize = () => {
  WindowMinimise(); // 调用 Wails 最小化函数
};

const maximize = () => {
  if (isFullscreen.value) {
    WindowUnfullscreen(); // 如果已全屏，则退出全屏
  } else {
    WindowFullscreen(); // 否则，进入全屏
  }
  isFullscreen.value = !isFullscreen.value; // 切换状态
};

const close = () => {
  Quit(); // 调用 Wails 关闭函数
};
</script>
  
  <style scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* Logo 和标题容器 */
.logo-container {
  display: flex;
  align-items: center;
  gap: 10px; /* Logo 和标题之间的间距 */
}

.logo img {
  height: 28px; /* 稍微调整 logo 大小 */
}

.title {
  font-size: 14px; /* 稍微调整标题大小 */
  font-weight: bold;
  color: #333; /* 标题颜色 */
}

.controls {
  display: flex;
  align-items: center;
  gap: 5px; /* 控制按钮之间的间距 */
  -webkit-app-region: no-drag; /* 禁止拖动按钮区域 */
}

.control-button {
  background: none;
  border: none;
  padding: 5px;
  cursor: pointer;
  transition: background-color 0.2s, opacity 0.2s;
  border-radius: 4px; /* 轻微圆角 */
  display: flex; /* 让图片居中 */
  align-items: center;
  justify-content: center;
}

.control-button img {
  width: 16px;
  height: 16px;
  display: block; /* 移除图片下方可能的空隙 */
}

.control-button:hover {
  background-color: rgba(0, 0, 0, 0.1); /* 悬停背景 */
}

.control-button:active {
  background-color: rgba(0, 0, 0, 0.2); /* 点击背景 */
}

/* 特别为设置按钮添加一些右边距 */
.settings-button {
  margin-right: 10px;
}

/* 设置弹窗样式 */
.settings-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5); /* 半透明遮罩 */
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000; /* 确保在顶层 */
  -webkit-app-region: no-drag; /* 弹窗区域不可拖动 */
}

.settings-dialog {
  background-color: #fff;
  padding: 25px 30px;
  border-radius: 8px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 400px; /* 限制最大宽度 */
}

.settings-dialog h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  text-align: center;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  color: #555;
  font-size: 0.9em;
}

.form-group input[type="text"] {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box; /* 防止 padding 撑大元素 */
  font-size: 1em;
}

.dialog-buttons {
  display: flex;
  justify-content: flex-end; /* 按钮靠右 */
  gap: 10px; /* 按钮间距 */
  margin-top: 25px;
}

.dialog-button {
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9em;
  transition: background-color 0.2s, color 0.2s;
}

.dialog-button.save {
  background-color: #007bff;
  color: white;
}
.dialog-button.save:hover {
  background-color: #0056b3;
}

.dialog-button.cancel {
  background-color: #f0f0f0;
  color: #333;
  border: 1px solid #ccc;
}
.dialog-button.cancel:hover {
  background-color: #e0e0e0;
}
</style>