<template>
  <div class="node-info-container">
    <div class="node-header">
      <div class="title-row">
        <button class="back-button" @click="goBack">
          <span class="back-icon">←</span> 返回
        </button>
        <h1 class="node-title">{{ nodeName }} 节点信息</h1>
      </div>
      <div class="node-details">
        <div class="detail-item">
          <span class="detail-label">端点:</span>
          <span class="detail-value">{{ endpoint }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">区域:</span>
          <span class="detail-value">{{ region }}</span>
        </div>
      </div>
    </div>

    <!-- 显示桶对象列表 -->
    <BucketObjects 
      v-if="showingBucketObjects && selectedBucket" 
      :bucketName="selectedBucket.name"
      :endpoint="endpoint"
      :region="region"
      :accessKey="accessKey"
      :secretKey="secretKey"
      @back="backToBucketList"
    />

    <!-- 显示桶列表 -->
    <template v-else>
      <div class="loading-container" v-if="loading">
        <div class="loading-spinner"></div>
        <p>正在加载桶信息...</p>
      </div>

      <div class="error-container" v-else-if="error">
        <p class="error-message">{{ error }}</p>
        <button class="retry-button" @click="fetchBucketInfo">重试</button>
      </div>

      <div class="buckets-container" v-else>
        <h2 class="section-title">桶列表 ({{ bucketInfo.length > 0 ? bucketInfo[0].buckets.length : 0 }})</h2>
        
        <div class="no-buckets" v-if="bucketInfo.length === 0 || bucketInfo[0].buckets.length === 0">
          <p>该节点下暂无桶信息</p>
        </div>
        
        <div class="bucket-list" v-else>
        <div class="bucket-card" v-for="bucket in bucketInfo[0].buckets" :key="bucket.name" @click="showBucketObjects(bucket)">
          <div class="bucket-header">
            <h3 class="bucket-name">{{ bucket.name }}</h3>
            <span class="bucket-region">{{ bucket.region }}</span>
          </div>
          
          <div class="bucket-details">
            <div class="bucket-detail">
              <span class="detail-label">创建时间:</span>
              <span class="detail-value">{{ formatDate(bucket.creationDate) }}</span>
            </div>
            <div class="bucket-detail">
              <span class="detail-label">已用空间:</span>
              <span class="detail-value">{{ formatSize(bucket.usedSpace) }}</span>
            </div>
            <div class="bucket-detail">
              <span class="detail-label">对象总数:</span>
              <span class="detail-value">{{ bucket.totalObjects }}</span>
            </div>
            
            <div class="bucket-features">
              <div class="feature" :class="{ 'feature-enabled': bucket.versioningEnabled }">
                版本控制: {{ bucket.versioningEnabled ? '已启用' : '未启用' }}
              </div>
              <div class="feature" :class="{ 'feature-enabled': bucket.publicAccessBlocked }">
                公共访问: {{ bucket.publicAccessBlocked ? '已阻止' : '未阻止' }}
              </div>
              <div class="feature" :class="{ 'feature-enabled': bucket.hasPolicy }">
                桶策略: {{ bucket.hasPolicy ? '已设置' : '未设置' }}
              </div>
              <div class="feature" :class="{ 'feature-enabled': bucket.encryptionEnabled }">
                加密: {{ bucket.encryptionEnabled ? bucket.encryptionType : '未启用' }}
              </div>
              <div class="feature" :class="{ 'feature-enabled': bucket.hasLifecycleRules }">
                生命周期规则: {{ bucket.hasLifecycleRules ? bucket.lifecycleRulesCount + '条' : '无' }}
              </div>
              <div class="feature" :class="{ 'feature-enabled': bucket.websiteEnabled }">
                静态网站: {{ bucket.websiteEnabled ? '已启用' : '未启用' }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { GetNodeBucketInfo } from '../../wailsjs/go/main/S3Manager';
import { LogDebug } from '../../wailsjs/runtime/runtime';
import BucketObjects from '../components/BucketObjects.vue';

const route = useRoute();
const router = useRouter();

// 返回上一级页面
function goBack() {
  router.push('/');
}

// 从路由参数中获取节点信息
const nodeName = ref(route.query.nodeName || '');
const endpoint = ref(route.query.endpoint || '');
const region = ref(route.query.region || '');
const accessKey = ref(route.query.accessKey || '');
const secretKey = ref(route.query.secretKey || '');

// 状态变量
const bucketInfo = ref([]);
const loading = ref(true);
const error = ref('');
const showingBucketObjects = ref(false);
const selectedBucket = ref(null);

// 在组件挂载时获取桶信息
onMounted(() => {
  fetchBucketInfo();
});

// 获取节点桶信息
async function fetchBucketInfo() {
  if (!endpoint.value || !accessKey.value || !secretKey.value || !region.value) {
    error.value = '节点信息不完整，无法获取桶信息';
    loading.value = false;
    return;
  }
  
  loading.value = true;
  error.value = '';
  
  try {
    LogDebug(`获取节点 ${nodeName.value} 的桶信息`);
    const result = await GetNodeBucketInfo(
      endpoint.value,
      accessKey.value,
      secretKey.value,
      region.value
    );
    
    bucketInfo.value = result || [];
    LogDebug(`获取到 ${bucketInfo.value.length > 0 ? bucketInfo.value[0].buckets.length : 0} 个桶信息`);
  } catch (err) {
    error.value = `获取桶信息失败: ${err.message || err}`;
    LogDebug(`获取桶信息失败: ${err}`);
  } finally {
    loading.value = false;
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

// 显示桶对象列表
function showBucketObjects(bucket) {
  selectedBucket.value = bucket;
  showingBucketObjects.value = true;
  LogDebug(`显示桶 ${bucket.name} 的对象列表`);
}

// 返回桶列表
function backToBucketList() {
  showingBucketObjects.value = false;
  selectedBucket.value = null;
}
</script>

<style scoped>
.node-info-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.node-header {
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e0e0e0;
}

.title-row {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
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

.node-title {
  font-size: 24px;
  color: #333;
  margin: 0;
}

.node-details {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.detail-item {
  display: flex;
  align-items: center;
}

.detail-label {
  font-weight: 600;
  margin-right: 8px;
  color: #666;
}

.detail-value {
  color: #333;
}

.section-title {
  font-size: 20px;
  margin-bottom: 20px;
  color: #333;
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

.no-buckets {
  text-align: center;
  padding: 30px;
  background-color: #f8f9fa;
  border-radius: 8px;
  color: #666;
}

.bucket-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.bucket-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
  padding: 20px;
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.bucket-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  background-color: #f8f9fa;
}

.bucket-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.bucket-name {
  font-size: 18px;
  margin: 0;
  color: #333;
}

.bucket-region {
  font-size: 14px;
  color: #666;
  background-color: #f0f0f0;
  padding: 3px 8px;
  border-radius: 4px;
}

.bucket-details {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.bucket-detail {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.bucket-features {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.feature {
  font-size: 13px;
  padding: 5px 10px;
  border-radius: 4px;
  background-color: #f5f5f5;
  color: #666;
}

.feature-enabled {
  background-color: #e6f7ff;
  color: #1890ff;
}
</style>