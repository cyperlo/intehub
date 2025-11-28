<template>
  <div class="dashboard" v-loading="loading">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #ecf5ff; color: #409eff;">
              <el-icon :size="40"><Setting /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.totalConfigs }}</div>
              <div class="stat-label">推送配置</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #f0f9ff; color: #67c23a;">
              <el-icon :size="40"><CircleCheck /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.successCount }}</div>
              <div class="stat-label">成功推送</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #fef0f0; color: #f56c6c;">
              <el-icon :size="40"><CircleClose /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.failureCount }}</div>
              <div class="stat-label">失败推送</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #fdf6ec; color: #e6a23c;">
              <el-icon :size="40"><Timer /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.avgDuration }}ms</div>
              <div class="stat-label">平均耗时</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card class="recent-history" shadow="hover" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>最近推送记录</span>
          <el-button type="primary" link @click="goToHistory">查看全部</el-button>
        </div>
      </template>
      
      <el-table :data="recentHistory" style="width: 100%">
        <el-table-column prop="config_name" label="配置名称" width="180" />
        <el-table-column prop="url" label="推送URL" show-overflow-tooltip />
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column prop="created_at" label="推送时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getPushConfigs, getPushHistory } from '../api/push'
import type { PushHistory } from '../types'

const router = useRouter()
const loading = ref(true)

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
  loadData()
})
</script>

<style scoped>
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
}

.stat-content {
  flex: 1;
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

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
