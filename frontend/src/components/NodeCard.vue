<script setup>
import { computed, inject } from "vue";
import { useRouter } from "vue-router";
import { GetNodeBucketInfo } from "../../wailsjs/go/main/S3Manager";

const router = useRouter();

const props = defineProps({
  title: {
    type: String,
    default: "节点名称",
  },
  endpoint: {
    type: String,
    default: "",
  },
  region: {
    type: String,
    default: "",
  },
  accessKey: {
    type: String,
    default: "",
  },
  secretKey: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["view-node"]);
function viewNode() {
  // 发出事件通知父组件
  emit("view-node", props.title);
  
  // 跳转到nodeInfo路由，并传递节点信息作为路由参数
  console.log("viewNode");
  router.push({
    path: "/nodeInfo",
    query: {
      nodeName: props.title,
      endpoint: props.endpoint,
      region: props.region,
      accessKey: props.accessKey,
      secretKey: props.secretKey
    }
  });
}

// 对SecretKey进行掩码处理，中间部分用***代替
function maskSecretKey(key) {
  if (!key || key.length < 8) return key;

  const firstPart = key.substring(0, 4);
  const lastPart = key.substring(key.length - 4);
  return `${firstPart}***${lastPart}`;
}
</script>

<template>
  <div class="node-card">
    <div class="card-inner">
      <div class="card-header">
        <div class="logo">
          <img src="../assets/node.png" alt="Node Logo" />
        </div>
        <h3 class="card-title">{{ title }}</h3>
      </div>
      <div class="card-content">
        <div class="info-item">
          <span class="info-label">端点:</span>
          <span class="info-value">{{ endpoint }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">区域:</span>
          <span class="info-value">{{ region }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">AccessKey:</span>
          <span class="info-value">{{ accessKey }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">SecretKey:</span>
          <span class="info-value">{{ maskSecretKey(secretKey) }}</span>
        </div>
      </div>
      <div class="card-actions">
        <button class="view-btn" @click="viewNode()">查看节点</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.node-card {
  perspective: 1000px;
  height: 100%;
}

.card-inner {
  background: linear-gradient(135deg, #ffffff, #f8f9fa);
  border-radius: 16px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.08), 0 6px 6px rgba(0, 0, 0, 0.05);
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
  transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  border: 1px solid rgba(255, 255, 255, 0.8);
  position: relative;
  overflow: hidden;
}

.card-inner::before {
  content: "";
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(
    circle,
    rgba(255, 255, 255, 0.8) 0%,
    rgba(255, 255, 255, 0) 70%
  );
  opacity: 0;
  transform: scale(0.5);
  transition: opacity 0.5s, transform 0.7s;
  pointer-events: none;
}

.node-card:hover .card-inner {
  transform: translateY(-10px) rotateX(5deg);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.12), 0 8px 8px rgba(0, 0, 0, 0.06);
  border-color: rgba(255, 255, 255, 1);
}

.node-card:hover .card-inner::before {
  opacity: 0.8;
  transform: scale(1);
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  position: relative;
}

.logo {
  margin-right: 12px;
  background: linear-gradient(135deg, #e6f7ff, #b3e0ff);
  border-radius: 10px;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
}

.logo img {
  width: 28px;
  height: 28px;
  object-fit: contain;
  filter: drop-shadow(0 2px 3px rgba(0, 0, 0, 0.1));
}

.card-title {
  margin: 0;
  font-size: 18px;
  color: #2c3e50;
  font-weight: 600;
  letter-spacing: 0.5px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.card-content {
  color: #5a6a7e;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 6px 0;
  border-bottom: 1px dashed rgba(0, 0, 0, 0.06);
}

.info-label {
  font-weight: 500;
  color: #7f8c8d;
  min-width: 90px;
  font-size: 14px;
}

.info-value {
  font-size: 14px;
  color: #34495e;
  font-weight: 500;
  word-break: break-all;
}

.card-actions {
  margin-top: 15px;
  text-align: right;
}

.view-btn {
  background: linear-gradient(135deg, #4caf50, #45a049);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  letter-spacing: 0.5px;
}

.view-btn:hover {
  background: linear-gradient(135deg, #45a049, #3d8b3d);
  transform: translateY(-2px);
  box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
}

.view-btn:active {
  transform: translateY(1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
