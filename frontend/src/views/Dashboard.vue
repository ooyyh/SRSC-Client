<template>
  <div class="dashboard">
    <!-- 添加节点按钮 -->
    <div class="add-node-container">
      <button class="add-node-btn" @click="showAddNodeDialog = true">
        <span class="plus-icon">+</span>
        <span>添加节点</span>
      </button>
    </div>
    
    <!-- 节点信息卡片 -->
    <div class="cards-container">
      <NodeCard 
        v-for="node in nodesInfo" 
        :key="node.NodeName"
        :title="node.NodeName" 
        :endpoint="node.EndPoint" 
        :region="node.Region"
        :accessKey="node.AccessKey"
        :secretKey="node.SecretKey"
        @view-node="handleViewNode"
      />
      <div v-if="nodesInfo.length === 0" class="no-data">暂无节点信息</div>
    </div>
    
    <!-- 添加节点对话框 -->
    <AddNodeDialog 
      :show="showAddNodeDialog" 
      @close="showAddNodeDialog = false"
      @node-added="handleNodeAdded"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import NodeCard from "../components/NodeCard.vue";
import AddNodeDialog from "../components/AddNodeDialog.vue";
import { GetAllS3NodesInfo } from "../../wailsjs/go/main/S3Manager";
import { LogDebug } from "../../wailsjs/runtime/runtime";

// 响应式数据
const nodesInfo = ref([]);
const isLoading = ref(false);
const showAddNodeDialog = ref(false);

// 处理查看节点按钮点击
function handleViewNode(nodeName) {
  LogDebug("查看节点: " + nodeName);
  // 节点详情页的跳转逻辑已在NodeCard组件中实现
}

// 处理节点添加成功
function handleNodeAdded() {
  LogDebug("节点添加成功，刷新节点列表");
  // 重新加载节点数据
  loadData();
}

// 加载所有数据
async function loadData() {
  isLoading.value = true;
  try {
    // 获取节点信息
    const nodes = await GetAllS3NodesInfo();
    console.log(nodes);
    nodesInfo.value = nodes || [];
    LogDebug("获取到节点信息: " + JSON.stringify(nodesInfo.value));
  } catch (error) {
    LogDebug("加载数据出错: " + error);
  } finally {
    isLoading.value = false;
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadData();
});
</script>

<style scoped>
.dashboard {
  padding: 30px;
  min-height: 100vh;
  background: linear-gradient(135deg, #f8f9fa, #e9ecef);
}

/* 添加节点按钮容器 */
.add-node-container {
  display: flex;
  justify-content: flex-end;
  max-width: 1600px;
  margin: 0 auto 20px;
  padding: 0 20px;
}

/* 添加节点按钮样式 */
.add-node-btn {
  display: flex;
  align-items: center;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 20px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.add-node-btn:hover {
  background-color: #0069d9;
}

.plus-icon {
  font-size: 1.2rem;
  margin-right: 8px;
  font-weight: bold;
}

.cards-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 70px;
  margin: 15px auto;
  max-width: 1600px;
  padding: 20px;
}

.no-data {
  grid-column: 1 / -1;
  text-align: center;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  color: #888;
  font-style: italic;
}

@media (max-width: 768px) {
  .cards-container {
    grid-template-columns: 1fr;
  }
  
  .add-node-container {
    justify-content: center;
  }
}
</style>