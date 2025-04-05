<template>
  <div class="bucket-objects-container">
    <div class="objects-header">
      <div class="title-row">
        <button class="back-button" @click="$emit('back')">
          <span class="back-icon">←</span> 返回桶列表
        </button>
        <h2 class="objects-title">{{ bucketName }} 对象列表</h2>
      </div>
      <div class="actions-row">
        <label class="upload-button" :class="{ 'uploading': uploading }">
          <span class="upload-icon">+</span> 
          <span v-if="!uploading">上传对象</span>
          <span v-else>上传中...</span>
          <input class="file-input" @click="uploadFile" :disabled="uploading" />
        </label>
        <div class="upload-status" v-if="uploadStatus">
          <span :class="{'status-error': uploadError, 'status-success': !uploadError}">
            {{ uploadStatus }}
          </span>
        </div>
      </div>
    </div>

    <div class="loading-container" v-if="loading">
      <div class="loading-spinner"></div>
      <p>正在加载对象列表...</p>
    </div>

    <div class="error-container" v-else-if="error">
      <p class="error-message">{{ error }}</p>
      <button class="retry-button" @click="fetchObjects">重试</button>
    </div>

    <div class="objects-container" v-else>
      <div class="no-objects" v-if="objects.length === 0">
        <p>该桶中暂无对象</p>
      </div>
      
      <div class="objects-list" v-else>
        <table class="objects-table">
          <thead>
            <tr>
              <th>名称</th>
              <th>大小</th>
              <th>最后修改时间</th>
              <th>存储类型</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="object in objects" :key="object.key">
              <td class="object-name">{{ object.key }}</td>
              <td>{{ formatSize(object.size) }}</td>
              <td>{{ formatDate(object.lastModified) }}</td>
              <td>{{ object.storageClass }}</td>
              <td class="actions-cell">
                <button class="action-button info-button" @click="showObjectInfo(object)">
                  详情
                </button>
                <button class="action-button download-button" @click="downloadObject(object)">
                  下载
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 对象详情弹窗 -->
    <div class="modal" v-if="showObjectDetails">
      <div class="modal-content">
        <div class="modal-header">
          <h3>对象详情</h3>
          <button class="close-button" @click="showObjectDetails = false">×</button>
        </div>
        <div class="modal-body" v-if="selectedObject">
          <div class="detail-row">
            <span class="detail-label">键名:</span>
            <span class="detail-value">{{ selectedObject.key }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">大小:</span>
            <span class="detail-value">{{ formatSize(selectedObject.size) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">最后修改时间:</span>
            <span class="detail-value">{{ formatDate(selectedObject.lastModified) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">ETag:</span>
            <span class="detail-value">{{ selectedObject.etag }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">内容类型:</span>
            <span class="detail-value">{{ selectedObject.contentType }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">存储类型:</span>
            <span class="detail-value">{{ selectedObject.storageClass }}</span>
          </div>
          <div class="detail-row" v-if="selectedObject.versionId">
            <span class="detail-label">版本ID:</span>
            <span class="detail-value">{{ selectedObject.versionId }}</span>
          </div>
          <div class="metadata-section" v-if="Object.keys(selectedObject.metadata || {}).length > 0">
            <h4>元数据</h4>
            <div class="detail-row" v-for="(value, key) in selectedObject.metadata" :key="key">
              <span class="detail-label">{{ key }}:</span>
              <span class="detail-value">{{ value }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 移除上传对象弹窗 -->
  </div>
</template>

<script setup>
import { ref, onMounted, defineProps, defineEmits } from 'vue';
import { ListObjects, GetObjectInfo, DownloadObject, UploadObject } from '../../wailsjs/go/main/S3Manager';
import { LogDebug } from '../../wailsjs/runtime/runtime';

const props = defineProps({
  bucketName: {
    type: String,
    required: true
  },
  endpoint: {
    type: String,
    required: true
  },
  region: {
    type: String,
    required: true
  },
  accessKey: {
    type: String,
    required: true
  },
  secretKey: {
    type: String,
    required: true
  }
});

const emit = defineEmits(['back']);

// 状态变量
const objects = ref([]);
const loading = ref(true);
const error = ref('');
const showObjectDetails = ref(false);
const selectedObject = ref(null);
const selectedFile = ref(null);
const uploading = ref(false);
const uploadStatus = ref('');
const uploadError = ref(false);

// 在组件挂载时获取对象列表
onMounted(() => {
  fetchObjects();
});

// 获取桶中的对象列表
async function fetchObjects() {
  if (!props.endpoint || !props.accessKey || !props.secretKey || !props.region || !props.bucketName) {
    error.value = '信息不完整，无法获取对象列表';
    loading.value = false;
    return;
  }
  
  loading.value = true;
  error.value = '';
  
  try {
    LogDebug(`获取桶 ${props.bucketName} 的对象列表`);
    // 调用后端API获取对象列表
    const result = await ListObjects(
      props.endpoint,
      props.accessKey,
      props.secretKey,
      props.region,
      props.bucketName
    );
    
    objects.value = result || [];
    LogDebug(`获取到 ${objects.value.length} 个对象`);
  } catch (err) {
    error.value = `获取对象列表失败: ${err.message || err}`;
    LogDebug(`获取对象列表失败: ${err}`);
  } finally {
    loading.value = false;
  }
}

// 显示对象详情
async function showObjectInfo(object) {
  try {
    // 获取对象的详细信息
    const detailedInfo = await GetObjectInfo(
      props.endpoint,
      props.accessKey,
      props.secretKey,
      props.region,
      props.bucketName,
      object.key
    );
    
    selectedObject.value = detailedInfo;
    showObjectDetails.value = true;
    LogDebug(`获取对象详情成功: ${object.key}`);
  } catch (err) {
    LogDebug(`获取对象详情失败: ${err.message || err}`);
    // 使用Toast通知替代alert
    window.toast.error(`获取对象详情失败: ${err.message || err}`);
  }
}

// 下载对象
async function downloadObject(object) {
  try {
    LogDebug(`开始下载对象: ${object.key}`);
    await DownloadObject(
      props.endpoint,
      props.accessKey,
      props.secretKey,
      props.region,
      props.bucketName,
      object.key,
      '' // 下载路径，空字符串表示使用默认路径
    );
    LogDebug(`对象下载成功: ${object.key}`);
    // 使用Toast通知替代alert
    window.toast.success(`对象 ${object.key} 下载成功`);
  } catch (err) {
    LogDebug(`对象下载失败: ${err.message || err}`);
    // 使用Toast通知替代alert
    window.toast.error(`下载失败: ${err.message || err}`);
  }
}

// 处理文件选择
// function handleFileSelect(event) {
//     console.log('文件选择事件触发');
//     uploadFile();
// //   }
// }

// 上传文件
async function uploadFile() {

  
  uploading.value = true;
  uploadStatus.value = '正在上传...';
  uploadError.value = false;
  
  try {
    LogDebug(`开始上传文件到桶: ${props.bucketName}`);
    // 注意：由于后端的UploadObject方法会自动打开文件选择对话框
    // 这里我们需要先让用户选择文件，然后再调用后端方法
    // 但是由于我们已经在前端选择了文件，所以这里会导致用户看到两次文件选择框
    // 调用后端上传方法
    const etag = await UploadObject(
      props.endpoint,
      props.accessKey,
      props.secretKey,
      props.region,
      props.bucketName
    );
    
    LogDebug(`文件上传成功，ETag: ${etag}`);
    uploadStatus.value = '上传成功!';
    // 刷新对象列表
    fetchObjects();
    // 延迟清除状态
    setTimeout(() => {
      uploadStatus.value = '';
      selectedFile.value = null;
    }, 3000);
  } catch (err) {
    LogDebug(`文件上传失败: ${err.message || err}`);
    uploadStatus.value = `上传失败: ${err.message || err}`;
    uploadError.value = true;
  } finally {
    uploading.value = false;
  }
}

// 格式化日期
function formatDate(dateStr) {
  if (!dateStr) return '未知';
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
}

// 格式化文件大小
function formatSize(bytes) {
  if (bytes === 0) return '0 B';
  if (!bytes) return '未知';
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + units[i];
}
</script>

<style scoped>
.bucket-objects-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.objects-header {
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e0e0e0;
}

.title-row {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  justify-content: space-between;
}

.back-button {
  display: flex;
  align-items: center;
  background-color: #3498db;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-right: 15px;
  transition: background-color 0.3s;
}

.back-button:hover {
  background-color: #2980b9;
}

.back-icon {
  margin-right: 5px;
  font-size: 16px;
}

.objects-title {
  font-size: 22px;
  color: #333;
  margin: 0;
  flex-grow: 1;
}

.actions-row {
  display: flex;
  justify-content: flex-end;
}

.upload-button {
  display: flex;
  align-items: center;
  background-color: #27ae60;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
  position: relative;
}

.upload-button:hover:not(.uploading) {
  background-color: #219653;
}

.upload-button.uploading {
  background-color: #7f8c8d;
  cursor: wait;
}

.file-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.upload-status {
  margin-left: 10px;
  font-size: 14px;
}

.status-error {
  color: #e53e3e;
}

.status-success {
  color: #27ae60;
}

.upload-icon {
  margin-right: 5px;
  font-size: 16px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.loading-spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top: 4px solid #3498db;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-container {
  text-align: center;
  padding: 30px;
  background-color: #fff5f5;
  border-radius: 8px;
  margin: 20px 0;
}

.error-message {
  color: #e53e3e;
  margin-bottom: 15px;
}

.retry-button {
  background-color: #3498db;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.retry-button:hover {
  background-color: #2980b9;
}

.no-objects {
  text-align: center;
  padding: 30px;
  background-color: #f8f9fa;
  border-radius: 8px;
  color: #666;
}

.objects-table {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
}

.objects-table th,
.objects-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

.objects-table th {
  background-color: #f8f9fa;
  font-weight: 600;
  color: #333;
}

.objects-table tr:last-child td {
  border-bottom: none;
}

.objects-table tr:hover {
  background-color: #f5f5f5;
}

.object-name {
  font-weight: 500;
  color: #2c3e50;
  word-break: break-all;
}

.actions-cell {
  white-space: nowrap;
}

.action-button {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  margin-right: 5px;
  transition: background-color 0.3s;
}

.info-button {
  background-color: #f0f0f0;
  color: #333;
}

.info-button:hover {
  background-color: #e0e0e0;
}

.download-button {
  background-color: #3498db;
  color: white;
}

.download-button:hover {
  background-color: #2980b9;
}

/* 弹窗样式 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #e0e0e0;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.close-button:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.detail-row {
  display: flex;
  margin-bottom: 10px;
}

.detail-label {
  font-weight: 600;
  width: 150px;
  color: #666;
}

.detail-value {
  flex-grow: 1;
  word-break: break-all;
}

.metadata-section {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.metadata-section h4 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #333;
}

/* 上传表单样式 */
.upload-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.form-group label {
  font-weight: 600;
  color: #333;
}

.form-group input[type="text"] {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.selected-file {
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  font-size: 14px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

.upload-submit-button {
  background-color: #27ae60;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.upload-submit-button:hover:not(:disabled) {
  background-color: #219653;
}

.upload-submit-button:disabled {
  background-color: #a0a0a0;
  cursor: not-allowed;
}

.cancel-button {
  background-color: #f0f0f0;
  color: #333;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.cancel-button:hover:not(:disabled) {
  background-color: #e0e0e0;
}

.cancel-button:disabled {
  color: #a0a0a0;
  cursor: not-allowed;
}

.upload-status {
  padding: 10px;
  border-radius: 4px;
  font-size: 14px;
}

.status-error {
  color: #e53e3e;
  background-color: #fff5f5;
}

.status-success {
  color: #27ae60;
  background-color: #f0fff4;
}
</style>