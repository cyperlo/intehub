<template>
  <div class="push-history">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>推送历史记录</span>
        </div>
      </template>
      
      <el-table :data="history" style="width: 100%" v-loading="loading">
        <el-table-column prop="config_name" label="配置名称" width="150" />
        <el-table-column prop="url" label="推送URL" show-overflow-tooltip />
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status_code" label="状态码" width="100" />
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column prop="created_at" label="推送时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadHistory"
          @current-change="loadHistory"
        />
      </div>
    </el-card>
    
    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="推送详情" width="700px">
      <el-descriptions :column="2" border v-if="currentDetail">
        <el-descriptions-item label="配置名称">
          {{ currentDetail.config_name }}
        </el-descriptions-item>
        <el-descriptions-item label="推送URL">
          {{ currentDetail.url }}
        </el-descriptions-item>
        <el-descriptions-item label="请求方法">
          {{ currentDetail.method }}
        </el-descriptions-item>
        <el-descriptions-item label="状态码">
          <el-tag :type="currentDetail.success ? 'success' : 'danger'">
            {{ currentDetail.status_code }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="耗时">
          {{ currentDetail.duration }} ms
        </el-descriptions-item>
        <el-descriptions-item label="推送时间">
          {{ formatTime(currentDetail.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="推送内容" :span="2">
          <el-input
            v-model="currentDetail.content"
            type="textarea"
            :rows="5"
            readonly
          />
        </el-descriptions-item>
        <el-descriptions-item label="响应内容" :span="2">
          <el-input
            v-model="currentDetail.response"
            type="textarea"
            :rows="5"
            readonly
          />
        </el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="currentDetail.error">
          <el-text type="danger">{{ currentDetail.error }}</el-text>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPushHistory } from '../api/push'
import type { PushHistory } from '../types'

const loading = ref(false)
const history = ref<PushHistory[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const detailVisible = ref(false)
const currentDetail = ref<PushHistory | null>(null)

const loadHistory = async () => {
  loading.value = true
  try {
    const res = await getPushHistory({
      page: currentPage.value,
      page_size: pageSize.value
    })
    history.value = res.data
    total.value = res.total
  } catch (error) {
    console.error('加载历史记录失败:', error)
  } finally {
    loading.value = false
  }
}

const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN')
}

const handleViewDetail = (row: PushHistory) => {
  currentDetail.value = row
  detailVisible.value = true
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
