<template>
  <div class="dashboard" v-loading="loading">
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card">
            <div class="stat-icon" style="background: #ecf5ff; color: #409eff;">
              <el-icon :size="isMobile ? 30 : 40"><Setting /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.totalConfigs }}</div>
              <div class="stat-label">推送配置</div>
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
              <div class="stat-value">{{ stats.successCount }}</div>
              <div class="stat-label">成功推送</div>
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
              <div class="stat-value">{{ stats.failureCount }}</div>
              <div class="stat-label">失败推送</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card">
            <div class="stat-icon" style="background: #fdf6ec; color: #e6a23c;">
              <el-icon :size="isMobile ? 30 : 40"><Timer /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.avgDuration }}ms</div>
              <div class="stat-label">平均耗时</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card class="recent-history" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>最近推送记录</span>
          <el-button type="primary" link @click="goToHistory">查看全部</el-button>
        </div>
      </template>
      
      <!-- 桌面端表格 -->
      <el-table :data="recentHistory" class="desktop-table">
        <el-table-column prop="config_name" label="配置名称" min-width="150" />
        <el-table-column prop="url" label="推送URL" show-overflow-tooltip min-width="200" />
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'" size="small">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column prop="created_at" label="推送时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <!-- 移动端卡片列表 -->
      <div class="mobile-list">
        <el-card v-for="item in recentHistory" :key="item.id" class="history-card" shadow="hover">
          <div class="history-header">
            <span class="config-name">{{ item.config_name }}</span>
            <el-tag :type="item.success ? 'success' : 'danger'" size="small">
              {{ item.success ? '成功' : '失败' }}
            </el-tag>
          </div>
          <div class="history-info">
            <div class="info-row">
              <span class="label">URL：</span>
              <span class="value">{{ item.url }}</span>
            </div>
            <div class="info-row">
              <span class="label">方法：</span>
              <span class="value">{{ item.method }}</span>
            </div>
            <div class="info-row">
              <span class="label">耗时：</span>
              <span class="value">{{ item.duration }}ms</span>
            </div>
            <div class="info-row">
              <span class="label">时间：</span>
              <span class="value">{{ formatTime(item.created_at) }}</span>
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
import { getPushConfigs, getPushHistory } from '../api/push'
import type { PushHistory } from '../types'

const router = useRouter()
const loading = ref(true)
const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value <= 768)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const stats = ref({
  totalConfigs: 0,
  successCount: 0,
  failureCount: 0,
  avgDuration: 0
})

const recentHistory = ref<PushHistory[]>([])

const loadData = async () => {
  loading.value = true
  try {
    // 加载配置数量
    const configs = await getPushConfigs()
    stats.value.totalConfigs = configs.length
    
    // 加载历史记录
    const history = await getPushHistory({ page: 1, page_size: 10 })
    recentHistory.value = history.data
    
    // 计算统计数据
    if (history.data.length > 0) {
      stats.value.successCount = history.data.filter(h => h.success).length
      stats.value.failureCount = history.data.filter(h => !h.success).length
      const totalDuration = history.data.reduce((sum, h) => sum + h.duration, 0)
      stats.value.avgDuration = Math.round(totalDuration / history.data.length)
    }
  } catch (error: any) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载数据失败：' + (error?.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN')
}

const goToHistory = () => {
  router.push('/push-history')
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
  margin-bottom: 20px;
}

.stat-card-wrapper {
  margin-bottom: 12px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.recent-history {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.desktop-table {
  width: 100%;
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
  border-bottom: 1px solid #ebeef5;
}

.config-name {
  font-weight: 500;
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
  color: #303133;
  flex: 1;
  word-break: break-all;
}

@media (max-width: 768px) {
  .stats-row {
    margin-bottom: 12px;
  }

  .stat-card {
    gap: 12px;
  }

  .stat-icon {
    width: 60px;
    height: 60px;
  }

  .stat-value {
    font-size: 22px;
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
    margin-top: 12px;
  }
}

@media (max-width: 480px) {
  .stat-value {
    font-size: 20px;
  }

  .stat-icon {
    width: 50px;
    height: 50px;
  }
}
</style>
