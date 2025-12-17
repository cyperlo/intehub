<template>
  <div class="system-logs">
    <div class="filter-bar">
      <el-select v-model="filters.level" placeholder="日志级别" clearable size="default" class="filter-item">
        <el-option label="信息" value="info" />
        <el-option label="警告" value="warning" />
        <el-option label="错误" value="error" />
      </el-select>
      <el-input v-model="filters.module" placeholder="模块" clearable size="default" class="filter-item" />
      <el-button type="primary" @click="loadLogs" size="default">查询</el-button>
    </div>

    <!-- 桌面端表格 -->
    <el-table :data="logs" v-loading="loading" class="desktop-table">
      <el-table-column label="级别" width="100">
        <template #default="{ row }">
          <el-tag :type="getLevelType(row.level)" size="small">{{ row.level }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="module" label="模块" width="120" />
      <el-table-column prop="action" label="操作" width="150" />
      <el-table-column prop="content" label="内容" show-overflow-tooltip min-width="200" />
      <el-table-column prop="ip" label="IP" width="130" />
      <el-table-column prop="created_at" label="时间" width="180" />
    </el-table>

    <!-- 移动端卡片列表 -->
    <div class="mobile-list" v-loading="loading">
      <el-card v-for="log in logs" :key="log.id" class="log-card" shadow="hover">
        <div class="log-header">
          <el-tag :type="getLevelType(log.level)" size="small">{{ log.level }}</el-tag>
          <span class="log-time">{{ log.created_at }}</span>
        </div>
        <div class="log-info">
          <div class="info-row">
            <span class="label">模块：</span>
            <span class="value">{{ log.module }}</span>
          </div>
          <div class="info-row">
            <span class="label">操作：</span>
            <span class="value">{{ log.action }}</span>
          </div>
          <div class="info-row">
            <span class="label">内容：</span>
            <span class="value">{{ log.content }}</span>
          </div>
          <div class="info-row">
            <span class="label">IP：</span>
            <span class="value">{{ log.ip }}</span>
          </div>
        </div>
      </el-card>
    </div>

    <el-pagination
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      :layout="isMobile ? 'prev, pager, next' : 'total, prev, pager, next'"
      @current-change="loadLogs"
      class="pagination"
      small
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, computed } from 'vue'
import { getSystemLogs, type SystemLog } from '../../api/system'

const loading = ref(false)
const logs = ref<SystemLog[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const filters = reactive({
  level: '',
  module: ''
})

const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value <= 768)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const getLevelType = (level: string) => {
  const types: Record<string, any> = {
    info: 'info',
    warning: 'warning',
    error: 'danger'
  }
  return types[level] || 'info'
}

const loadLogs = async () => {
  loading.value = true
  try {
    const res = await getSystemLogs({
      page: page.value,
      page_size: pageSize.value,
      ...filters
    })
    logs.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('加载日志失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadLogs()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.system-logs {
  width: 100%;
}

.filter-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.filter-item {
  flex: 1;
  min-width: 120px;
  max-width: 200px;
}

.desktop-table {
  width: 100%;
}

.mobile-list {
  display: none;
}

.log-card {
  margin-bottom: 12px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.log-time {
  font-size: 12px;
  color: #909399;
}

.log-info {
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

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .desktop-table {
    display: none;
  }
  
  .mobile-list {
    display: block;
  }
  
  .filter-bar {
    gap: 8px;
  }
  
  .filter-item {
    flex: 1;
    min-width: 100px;
    max-width: none;
  }
}
</style>
