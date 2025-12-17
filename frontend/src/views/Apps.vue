<template>
  <div class="apps">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="header-title">应用管理</span>
          <el-button type="primary" @click="handleAdd" size="default">
            <el-icon><Plus /></el-icon>
            <span class="btn-text">新增应用</span>
          </el-button>
        </div>
      </template>

      <!-- 桌面端表格 -->
      <el-table :data="apps" v-loading="loading" class="desktop-table">
        <el-table-column prop="name" label="应用名称" min-width="150" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip min-width="200" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" @change="handleToggleStatus(row)" size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="handleRun(row)" link>
              <el-icon><VideoPlay /></el-icon>
              运行
            </el-button>
            <el-button type="primary" size="small" @click="handleViewLogs(row)" link>
              <el-icon><Document /></el-icon>
              日志
            </el-button>
            <el-button type="warning" size="small" @click="handleEdit(row)" link>
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)" link>
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 移动端卡片列表 -->
      <div class="mobile-list" v-loading="loading">
        <el-card v-for="app in apps" :key="app.id" class="app-card" shadow="hover">
          <div class="app-header">
            <div class="app-name-row">
              <h4>{{ app.name }}</h4>
              <el-switch v-model="app.enabled" @change="handleToggleStatus(app)" size="small" />
            </div>
          </div>
          <div class="app-info">
            <div class="info-row" v-if="app.description">
              <span class="label">描述：</span>
              <span class="value">{{ app.description }}</span>
            </div>
            <div class="info-row">
              <span class="label">创建时间：</span>
              <span class="value">{{ app.created_at }}</span>
            </div>
          </div>
          <div class="card-actions">
            <el-button type="success" size="small" @click="handleRun(app)">运行</el-button>
            <el-button type="primary" size="small" @click="handleViewLogs(app)">日志</el-button>
            <el-button type="warning" size="small" @click="handleEdit(app)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(app)">删除</el-button>
          </div>
        </el-card>
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" :width="isMobile ? '95%' : '800px'" :fullscreen="isMobile">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="应用名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入应用名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="请输入应用描述" />
        </el-form-item>
        <el-form-item label="应用代码" prop="code">
          <el-input 
            v-model="form.code" 
            type="textarea" 
            :rows="15" 
            placeholder="请输入 Go 代码，例如：&#10;package main&#10;&#10;import &quot;fmt&quot;&#10;&#10;func main() {&#10;    fmt.Println(&quot;Hello, World!&quot;)&#10;}"
            style="font-family: 'Courier New', monospace;"
          />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 日志对话框 -->
    <el-dialog v-model="logDialogVisible" title="应用执行日志" :width="isMobile ? '95%' : '900px'" :fullscreen="isMobile">
      <el-table :data="logs" v-loading="logLoading">
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="output" label="输出" show-overflow-tooltip />
        <el-table-column prop="error" label="错误" show-overflow-tooltip />
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column prop="started_at" label="开始时间" width="180" />
      </el-table>
      <el-pagination
        v-model:current-page="logPage"
        :page-size="logPageSize"
        :total="logTotal"
        layout="total, prev, pager, next"
        @current-change="loadLogs"
        style="margin-top: 20px"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getApps, createApp, updateApp, deleteApp, runApp, getAppLogs, type App, type AppLog } from '../api/app'

const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value <= 768)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const loading = ref(false)
const apps = ref<App[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('新增应用')
const submitting = ref(false)
const formRef = ref<FormInstance>()
const currentAppId = ref<number>()

const logDialogVisible = ref(false)
const logLoading = ref(false)
const logs = ref<AppLog[]>([])
const logPage = ref(1)
const logPageSize = ref(20)
const logTotal = ref(0)
const currentLogAppId = ref<number>()

const form = reactive<App>({
  name: '',
  description: '',
  code: '',
  language: 'go',
  enabled: true
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入应用代码', trigger: 'blur' }]
}

const loadApps = async () => {
  loading.value = true
  try {
    apps.value = await getApps()
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增应用'
  currentAppId.value = undefined
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: App) => {
  dialogTitle.value = '编辑应用'
  currentAppId.value = row.id
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (currentAppId.value) {
          await updateApp(currentAppId.value, form)
          ElMessage.success('更新成功')
        } else {
          await createApp(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadApps()
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleToggleStatus = async (row: App) => {
  try {
    await updateApp(row.id!, row)
    ElMessage.success(row.enabled ? '已启用' : '已禁用')
  } catch (error) {
    row.enabled = !row.enabled
    ElMessage.error('操作失败')
  }
}

const handleRun = async (row: App) => {
  try {
    await ElMessageBox.confirm('确定要运行这个应用吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    const result = await runApp(row.id!)
    ElMessage.success('运行成功')
    console.log('运行结果:', result)
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '运行失败')
    }
  }
}

const handleDelete = async (row: App) => {
  try {
    await ElMessageBox.confirm('确定要删除这个应用吗?', '提示', {
      type: 'warning'
    })
    await deleteApp(row.id!)
    ElMessage.success('删除成功')
    loadApps()
  } catch (error) {
    // 用户取消
  }
}

const handleViewLogs = (row: App) => {
  currentLogAppId.value = row.id
  logPage.value = 1
  loadLogs()
  logDialogVisible.value = true
}

const loadLogs = async () => {
  if (!currentLogAppId.value) return
  
  logLoading.value = true
  try {
    const res = await getAppLogs({
      app_id: currentLogAppId.value,
      page: logPage.value,
      page_size: logPageSize.value
    })
    logs.value = res.list
    logTotal.value = res.total
  } finally {
    logLoading.value = false
  }
}

const resetForm = () => {
  Object.assign(form, {
    id: undefined,
    name: '',
    description: '',
    code: '',
    language: 'go',
    enabled: true
  })
  formRef.value?.clearValidate()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadApps()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.apps {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 16px;
  font-weight: 500;
}

.desktop-table {
  width: 100%;
}

.mobile-list {
  display: none;
}

.app-card {
  margin-bottom: 12px;
}

.app-header {
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.app-name-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-name-row h4 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.app-info {
  margin-bottom: 12px;
  font-size: 14px;
}

.info-row {
  display: flex;
  margin-bottom: 8px;
}

.info-row .label {
  color: #909399;
  min-width: 80px;
  flex-shrink: 0;
}

.info-row .value {
  color: #303133;
  flex: 1;
  word-break: break-all;
}

.card-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  padding-top: 8px;
  border-top: 1px solid #ebeef5;
}

@media (max-width: 768px) {
  .desktop-table {
    display: none;
  }
  
  .mobile-list {
    display: block;
  }
  
  .btn-text {
    margin-left: 4px;
  }
}
</style>
