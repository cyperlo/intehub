<template>
  <div class="dashboard" v-loading="loading">
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card">
            <div class="stat-icon" style="background: #ecf5ff; color: #409eff;">
              <el-icon :size="isMobile ? 30 : 40"><Box /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.totalApps }}</div>
              <div class="stat-label">应用总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card">
            <div class="stat-icon" style="background: #f0f9ff; color: #67c23a;">
              <el-icon :size="isMobile ? 30 : 40"><ShoppingCart /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.totalTemplates }}</div>
              <div class="stat-label">应用商店模板</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card">
            <div class="stat-icon" style="background: #f0f9ff; color: #67c23a;">
              <el-icon :size="isMobile ? 30 : 40"><CircleCheck /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.successRuns }}</div>
              <div class="stat-label">成功运行</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card">
            <div class="stat-icon" style="background: #fef0f0; color: #f56c6c;">
              <el-icon :size="isMobile ? 30 : 40"><CircleClose /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.failedRuns }}</div>
              <div class="stat-label">失败运行</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card class="recent-history" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>最近运行记录</span>
          <el-button type="primary" link @click="goToApps">查看全部</el-button>
        </div>
      </template>
      
      <!-- 桌面端表格 -->
      <el-table :data="recentLogs" class="desktop-table">
        <el-table-column prop="app_name" label="应用名称" min-width="150" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="output" label="输出" show-overflow-tooltip min-width="200" />
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column prop="started_at" label="运行时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.started_at) }}
          </template>
        </el-table-column>
      </el-table>

      <!-- 移动端卡片列表 -->
      <div class="mobile-list">
        <el-card v-for="item in recentLogs" :key="item.id" class="history-card" shadow="hover">
          <div class="history-header">
            <span class="config-name">{{ item.app_name }}</span>
            <el-tag :type="item.status === 'success' ? 'success' : 'danger'" size="small">
              {{ item.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </div>
          <div class="history-info">
            <div class="info-row">
              <span class="label">输出：</span>
              <span class="value">{{ item.output || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="label">耗时：</span>
              <span class="value">{{ item.duration }}ms</span>
            </div>
            <div class="info-row">
              <span class="label">时间：</span>
              <span class="value">{{ formatTime(item.started_at) }}</span>
            </div>
          </div>
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Box, ShoppingCart, CircleCheck, CircleClose } from '@element-plus/icons-vue'
import { getApps, getAppLogs } from '../api/app'
import { getTemplates } from '../api/appstore'

const router = useRouter()
const loading = ref(true)
const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value <= 768)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const stats = ref({
  totalApps: 0,
  totalTemplates: 0,
  successRuns: 0,
  failedRuns: 0
})

const recentLogs = ref<any[]>([])

const loadData = async () => {
  loading.value = true
  try {
    // 加载应用数量
    const apps = await getApps()
    stats.value.totalApps = apps.length
    
    // 加载模板数量
    const templates = await getTemplates()
    stats.value.totalTemplates = templates.length
    
    // 加载所有日志用于统计
    const allLogs = await getAppLogs({ page: 1, page_size: 1000 })
    const allLogsList = allLogs.list || []
    
    // 计算统计数据（基于所有日志）
    stats.value.successRuns = allLogsList.filter(l => l.status === 'success').length
    stats.value.failedRuns = allLogsList.filter(l => l.status === 'error').length
    
    // 只显示最近10条
    recentLogs.value = allLogsList.slice(0, 10)
  } catch (error: any) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const formatTime = (time: string) => {
  if (!time) return '-'
  try {
    const date = new Date(time)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hour = String(date.getHours()).padStart(2, '0')
    const minute = String(date.getMinutes()).padStart(2, '0')
    const second = String(date.getSeconds()).padStart(2, '0')
    return `${year}/${month}/${day} ${hour}:${minute}:${second}`
  } catch (e) {
    return time
  }
}

const goToApps = () => {
  router.push('/apps')
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadData()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.dashboard {
  width: 100%;
}

.stats-row {
  margin-bottom: 24px;
}

.stat-card-wrapper {
  margin-bottom: 16px;
}

:deep(.el-card) {
  border-radius: 16px !important;
  transition: all 0.3s ease;
}

:deep(.el-card:hover) {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.1) !important;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 72px;
  height: 72px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #303133 0%, #606266 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 6px;
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  font-weight: 500;
}

.recent-history {
  margin-top: 24px;
}

:deep(.recent-history .el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f2f5;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-weight: 600;
  font-size: 16px;
  color: #303133;
}

:deep(.el-button--primary) {
  font-weight: 500;
}

.desktop-table {
  width: 100%;
}

:deep(.el-table th) {
  background: #fafafa !important;
  font-weight: 600;
}

.mobile-list {
  display: none;
}

.history-card {
  margin-bottom: 12px;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f2f5;
}

.config-name {
  font-weight: 600;
  font-size: 15px;
  color: #303133;
}

.history-info {
  font-size: 14px;
}

.info-row {
  display: flex;
  margin-bottom: 8px;
}

.info-row .label {
  color: #909399;
  min-width: 50px;
  flex-shrink: 0;
}

.info-row .value {
  color: #606266;
  flex: 1;
  word-break: break-all;
}

@media (max-width: 768px) {
  .stats-row {
    margin-bottom: 16px;
  }

  .stat-card {
    gap: 14px;
  }

  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
  }

  .stat-value {
    font-size: 24px;
  }

  .stat-label {
    font-size: 13px;
  }

  .desktop-table {
    display: none;
  }

  .mobile-list {
    display: block;
  }

  .recent-history {
    margin-top: 16px;
  }
}

@media (max-width: 480px) {
  .stat-value {
    font-size: 22px;
  }

  .stat-icon {
    width: 48px;
    height: 48px;
  }
}
</style>
