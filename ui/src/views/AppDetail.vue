<template>
  <div class="app-detail">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button @click="goBack" text>
              <el-icon><ArrowLeft /></el-icon>
              返回
            </el-button>
            <span class="header-title">{{ app?.name || '应用详情' }}</span>
          </div>
          <div class="header-actions">
            <el-tag :type="app?.enabled ? 'success' : 'info'" size="large">
              {{ app?.enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </div>
        </div>
      </template>

      <div v-loading="loading">
        <!-- 基本信息 -->
        <div class="section">
          <h3 class="section-title">
            <el-icon><InfoFilled /></el-icon>
            基本信息
          </h3>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">应用名称</span>
              <span class="value">{{ app?.name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">描述</span>
              <span class="value">{{ app?.description || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">展示类型</span>
              <span class="value">
                <el-tag v-if="app?.display_type === 'none'" type="info" size="small">无展示（后台执行）</el-tag>
                <el-tag v-else-if="app?.display_type === 'page'" type="primary" size="small">独立页面展示</el-tag>
                <el-tag v-else type="success" size="small">弹窗展示</el-tag>
              </span>
            </div>
            <div class="info-item">
              <span class="label">创建时间</span>
              <span class="value">{{ formatTime(app?.created_at) }}</span>
            </div>
            <div class="info-item">
              <span class="label">更新时间</span>
              <span class="value">{{ formatTime(app?.updated_at) }}</span>
            </div>
          </div>
        </div>

        <!-- 配置参数 -->
        <div class="section" v-if="configs.length > 0">
          <h3 class="section-title">
            <el-icon><Setting /></el-icon>
            配置参数
          </h3>
          <el-table :data="configs" border stripe>
            <el-table-column prop="key" label="参数名" min-width="150" />
            <el-table-column prop="type" label="类型" width="120">
              <template #default="{ row }">
                <el-tag v-if="row.type === 'string'" size="small">字符串</el-tag>
                <el-tag v-else-if="row.type === 'number'" type="warning" size="small">数字</el-tag>
                <el-tag v-else type="info" size="small">布尔值</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="value" label="参数值" min-width="200">
              <template #default="{ row }">
                <span v-if="row.encrypted" class="encrypted-value">
                  <el-icon><Lock /></el-icon>
                  加密存储
                </span>
                <span v-else>{{ row.value || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="encrypted" label="加密" width="80" align="center">
              <template #default="{ row }">
                <el-icon v-if="row.encrypted" color="#67c23a"><Lock /></el-icon>
                <el-icon v-else color="#909399"><Unlock /></el-icon>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="section empty-section" v-else>
          <h3 class="section-title">
            <el-icon><Setting /></el-icon>
            配置参数
          </h3>
          <el-empty description="暂无配置参数" :image-size="60" />
        </div>

        <!-- 应用代码 -->
        <div class="section">
          <h3 class="section-title">
            <el-icon><Document /></el-icon>
            应用代码
          </h3>
          <div class="code-block">
            <div class="code-content">{{ app?.code || '// 暂无代码' }}</div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, InfoFilled, Setting, Document, Lock, Unlock } from '@element-plus/icons-vue'
import type { App, AppConfig } from '../api/app'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const app = ref<App>()
const configs = ref<AppConfig[]>([])

const appId = Number(route.params.id)

const loadAppDetail = async () => {
  loading.value = true
  try {
    const res: any = await fetch(`/api/v1/apps/${appId}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(r => r.json())

    if (res.code === 0) {
      app.value = res.data.app
      configs.value = res.data.configs || []
    } else {
      ElMessage.error(res.message || '获取应用详情失败')
    }
  } catch (error) {
    ElMessage.error('获取应用详情失败')
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

const formatTime = (time: string | undefined) => {
  if (!time) return '-'
  try {
    const date = new Date(time)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (e) {
    return time
  }
}

onMounted(() => {
  loadAppDetail()
})
</script>

<style scoped>
.app-detail {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.section {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.section:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 20px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid #409eff;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-item .label {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
}

.info-item .value {
  font-size: 14px;
  color: #303133;
}

.encrypted-value {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #67c23a;
}

.code-block {
  background: #1e1e1e;
  border-radius: 8px;
  overflow: hidden;
}

.code-content {
  padding: 16px;
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #d4d4d4;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 500px;
  overflow-y: auto;
}

.empty-section {
  min-height: 100px;
}

::deep(.el-card) {
  border-radius: 16px !important;
}

::deep(.el-table) {
  border-radius: 8px;
}
</style>
