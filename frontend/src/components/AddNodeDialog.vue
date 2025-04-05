<script setup>
import { ref } from 'vue';
import { AddNode, GetAllS3NodesInfo } from '../../wailsjs/go/main/S3Manager';
import { LogDebug } from '../../wailsjs/runtime/runtime';

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['close', 'node-added']);

// 表单数据
const form = ref({
  nodeName: '',
  endpoint: '',
  accessKey: '',
  secretKey: '',
  region: ''
});

// 表单验证状态
const formErrors = ref({
  nodeName: '',
  endpoint: '',
  accessKey: '',
  secretKey: '',
  region: ''
});

// 提交状态
const isSubmitting = ref(false);

// 验证表单
function validateForm() {
  let isValid = true;
  
  // 重置错误信息
  formErrors.value = {
    nodeName: '',
    endpoint: '',
    accessKey: '',
    secretKey: '',
    region: ''
  };

  // 验证节点名称
  if (!form.value.nodeName.trim()) {
    formErrors.value.nodeName = '节点名称不能为空';
    isValid = false;
  }

  // 验证端点
  if (!form.value.endpoint.trim()) {
    formErrors.value.endpoint = '端点不能为空';
    isValid = false;
  } else if (!form.value.endpoint.startsWith('http://') && !form.value.endpoint.startsWith('https://')) {
    formErrors.value.endpoint = '端点必须以http://或https://开头';
    isValid = false;
  }

  // 验证访问密钥
  if (!form.value.accessKey.trim()) {
    formErrors.value.accessKey = '访问密钥不能为空';
    isValid = false;
  }

  // 验证密钥
  if (!form.value.secretKey.trim()) {
    formErrors.value.secretKey = '密钥不能为空';
    isValid = false;
  }

  // 验证区域
  if (!form.value.region.trim()) {
    formErrors.value.region = '区域不能为空';
    isValid = false;
  }

  return isValid;
}

// 提交表单
async function submitForm() {
  if (!validateForm()) {
    return;
  }

  isSubmitting.value = true;
  try {
    const result = await AddNode(
      form.value.nodeName,
      form.value.endpoint,
      form.value.accessKey,
      form.value.secretKey,
      form.value.region
    );

    if (result) {
      LogDebug('节点添加成功');
      // 重置表单
      form.value = {
        nodeName: '',
        endpoint: '',
        accessKey: '',
        secretKey: '',
        region: ''
      };
      // 通知父组件节点已添加
      emit('node-added');
      // 关闭对话框
      emit('close');
    } else {
      LogDebug('节点添加失败');
    }
  } catch (error) {
    LogDebug('添加节点出错: ' + error);
  } finally {
    isSubmitting.value = false;
  }
}

// 关闭对话框
function closeDialog() {
  // 重置表单
  form.value = {
    nodeName: '',
    endpoint: '',
    accessKey: '',
    secretKey: '',
    region: ''
  };
  // 重置错误信息
  formErrors.value = {
    nodeName: '',
    endpoint: '',
    accessKey: '',
    secretKey: '',
    region: ''
  };
  // 通知父组件关闭对话框
  emit('close');
}
</script>

<template>
  <div v-if="show" class="dialog-overlay" @click="closeDialog">
    <div class="dialog" @click.stop>
      <div class="dialog-header">
        <h2>添加节点</h2>
        <button class="close-btn" @click="closeDialog">×</button>
      </div>
      <div class="dialog-body">
        <form @submit.prevent="submitForm">
          <div class="form-group">
            <label for="nodeName">节点名称</label>
            <input 
              type="text" 
              id="nodeName" 
              v-model="form.nodeName" 
              :class="{ 'error-input': formErrors.nodeName }"
            >
            <div class="error-message" v-if="formErrors.nodeName">{{ formErrors.nodeName }}</div>
          </div>
          
          <div class="form-group">
            <label for="endpoint">端点</label>
            <input 
              type="text" 
              id="endpoint" 
              v-model="form.endpoint" 
              placeholder="https://example.com"
              :class="{ 'error-input': formErrors.endpoint }"
            >
            <div class="error-message" v-if="formErrors.endpoint">{{ formErrors.endpoint }}</div>
          </div>
          
          <div class="form-group">
            <label for="region">区域</label>
            <input 
              type="text" 
              id="region" 
              v-model="form.region" 
              placeholder="us-east-1"
              :class="{ 'error-input': formErrors.region }"
            >
            <div class="error-message" v-if="formErrors.region">{{ formErrors.region }}</div>
          </div>
          
          <div class="form-group">
            <label for="accessKey">访问密钥</label>
            <input 
              type="text" 
              id="accessKey" 
              v-model="form.accessKey" 
              :class="{ 'error-input': formErrors.accessKey }"
            >
            <div class="error-message" v-if="formErrors.accessKey">{{ formErrors.accessKey }}</div>
          </div>
          
          <div class="form-group">
            <label for="secretKey">密钥</label>
            <input 
              type="password" 
              id="secretKey" 
              v-model="form.secretKey" 
              :class="{ 'error-input': formErrors.secretKey }"
            >
            <div class="error-message" v-if="formErrors.secretKey">{{ formErrors.secretKey }}</div>
          </div>
          
          <div class="form-actions">
            <button type="button" class="cancel-btn" @click="closeDialog">取消</button>
            <button type="submit" class="submit-btn" :disabled="isSubmitting">
              {{ isSubmitting ? '添加中...' : '添加' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.dialog {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  width: 90%;
  max-width: 500px;
  overflow: hidden;
  animation: dialog-fade-in 0.3s ease;
}

@keyframes dialog-fade-in {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.dialog-header h2 {
  margin: 0;
  font-size: 1.25rem;
  color: #343a40;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6c757d;
  padding: 0;
  line-height: 1;
}

.close-btn:hover {
  color: #343a40;
}

.dialog-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #495057;
}

input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.15s ease-in-out;
}

input:focus {
  border-color: #80bdff;
  outline: 0;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.error-input {
  border-color: #dc3545;
}

.error-input:focus {
  box-shadow: 0 0 0 0.2rem rgba(220, 53, 69, 0.25);
}

.error-message {
  color: #dc3545;
  font-size: 0.875rem;
  margin-top: 5px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 30px;
}

.cancel-btn, .submit-btn {
  padding: 10px 20px;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-btn {
  background-color: #f8f9fa;
  border: 1px solid #ced4da;
  color: #495057;
}

.cancel-btn:hover {
  background-color: #e9ecef;
}

.submit-btn {
  background-color: #007bff;
  border: 1px solid #007bff;
  color: white;
}

.submit-btn:hover {
  background-color: #0069d9;
  border-color: #0062cc;
}

.submit-btn:disabled {
  background-color: #6c757d;
  border-color: #6c757d;
  cursor: not-allowed;
  opacity: 0.65;
}
</style>