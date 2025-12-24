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
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="360" fixed="right">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="handleRun(row)" link>
              <el-icon><VideoPlay /></el-icon>
              运行
            </el-button>
            <el-button type="primary" size="small" @click="handleViewLogs(row)" link>
              <el-icon><Document /></el-icon>
              日志
            </el-button>
            <el-button type="info" size="small" @click="handlePublish(row)" link>
              <el-icon><Upload /></el-icon>
              发布
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
              <span class="value">{{ formatTime(app.created_at) }}</span>
            </div>
          </div>
          <div class="card-actions">
            <el-button type="success" size="small" @click="handleRun(app)">运行</el-button>
            <el-button type="primary" size="small" @click="handleViewLogs(app)">日志</el-button>
            <el-button type="info" size="small" @click="handlePublish(app)">发布</el-button>
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
        
        <!-- 应用配置 -->
        <el-divider>应用配置（可选）</el-divider>
        <el-alert 
          title="配置说明" 
          type="info" 
          :closable="false"
          style="margin-bottom: 16px;"
        >
          应用配置用于设置应用运行时的默认参数。例如：API密钥、URL地址等。运行应用时这些配置会自动传入。
        </el-alert>
        
        <div v-for="(config, index) in formConfigs" :key="index" class="config-item">
          <el-row :gutter="10">
            <el-col :span="6">
              <el-input v-model="config.key" placeholder="参数名" size="small" />
            </el-col>
            <el-col :span="8">
              <el-input v-model="config.value" placeholder="参数值" size="small" />
            </el-col>
            <el-col :span="5">
              <el-select v-model="config.type" placeholder="类型" size="small" style="width: 100%">
                <el-option label="字符串" value="string" />
                <el-option label="数字" value="number" />
                <el-option label="布尔" value="boolean" />
                <el-option label="JSON" value="json" />
              </el-select>
            </el-col>
            <el-col :span="3">
              <el-checkbox v-model="config.encrypted" size="small">加密</el-checkbox>
            </el-col>
            <el-col :span="2">
              <el-button type="danger" size="small" @click="removeConfig(index)" :icon="Delete" circle />
            </el-col>
          </el-row>
        </div>
        <el-button type="primary" size="small" @click="addConfig" plain style="width: 100%;">
          <el-icon><Plus /></el-icon>
          添加配置参数
        </el-button>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 日志对话框 -->
    <el-dialog v-model="logDialogVisible" title="应用执行日志" :width="isMobile ? '95%' : '900px'" :fullscreen="isMobile">
      <div class="log-container" v-loading="logLoading">
        <div v-if="logs.length === 0" class="empty-logs">暂无日志</div>
        <div v-else class="log-list">
          <div v-for="log in logs" :key="log.id" class="log-item" :class="log.status">
            <div class="log-header">
              <el-tag :type="log.status === 'success' ? 'success' : 'danger'" size="small">
                {{ log.status === 'success' ? '成功' : '失败' }}
              </el-tag>
              <span class="log-time">{{ formatTime(log.started_at) }}</span>
              <span class="log-duration">耗时: {{ log.duration }}ms</span>
            </div>
            <div v-if="log.output" class="log-content">
              <div class="log-label">输出：</div>
              <pre class="log-output">{{ formatOutput(log.output) }}</pre>
            </div>
            <div v-if="log.error" class="log-content">
              <div class="log-label">错误：</div>
              <pre class="log-error">{{ log.error }}</pre>
            </div>
          </div>
        </div>
      </div>
      <el-pagination
        v-model:current-page="logPage"
        :page-size="logPageSize"
        :total="logTotal"
        layout="total, prev, pager, next"
        @current-change="loadLogs"
        style="margin-top: 16px; text-align: center;"
      />
    </el-dialog>

    <!-- 发布到应用商店对话框 -->
    <el-dialog v-model="publishDialogVisible" title="发布到应用商店" :width="isMobile ? '95%' : '600px'" :fullscreen="isMobile">
      <el-form :model="publishForm" :rules="publishRules" ref="publishFormRef" label-width="100px">
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="publishForm.display_name" placeholder="请输入应用显示名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="publishForm.description" type="textarea" :rows="3" placeholder="请输入应用描述" />
        </el-form-item>
        <el-form-item label="版本" prop="version">
          <el-input v-model="publishForm.version" placeholder="例如: 1.0.0" />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="publishForm.category" placeholder="例如: 工具、数据处理等" />
        </el-form-item>
        <el-form-item label="作者">
          <el-input v-model="publishForm.author" placeholder="请输入作者名称" />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="publishForm.tags" placeholder="多个标签用逗号分隔" />
        </el-form-item>
        <el-form-item label="图标URL">
          <el-input v-model="publishForm.icon" placeholder="请输入图标URL" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="publishDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handlePublishSubmit" :loading="publishing">发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { getApps, createApp, updateApp, deleteApp, runApp, getAppLogs, publishToStore, type App, type AppLog, type PublishToStoreRequest, type AppConfig } from '../api/app'

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

const publishDialogVisible = ref(false)
const publishing = ref(false)
const publishFormRef = ref<FormInstance>()
const currentPublishAppId = ref<number>()
const publishForm = reactive<PublishToStoreRequest>({
  display_name: '',
  description: '',
  version: '1.0.0',
  category: '',
  author: '',
  tags: '',
  icon: ''
})

const form = reactive<App>({
  name: '',
  description: '',
  code: '',
  language: 'go',
  enabled: true
})

const formConfigs = ref<AppConfig[]>([])

const addConfig = () => {
  formConfigs.value.push({
    key: '',
    value: '',
    type: 'string',
    encrypted: false
  })
}

const removeConfig = (index: number) => {
  formConfigs.value.splice(index, 1)
}

const rules: FormRules = {
  name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入应用代码', trigger: 'blur' }]
}

const publishRules: FormRules = {
  display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  description: [{ required: true, message: '请输入描述', trigger: 'blur' }],
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }]
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
  formConfigs.value = []
  dialogVisible.value = true
}

const handleEdit = async (row: App) => {
  dialogTitle.value = '编辑应用'
  currentAppId.value = row.id
  
  // 获取应用详情和配置
  try {
    const res: any = await fetch(`/api/v1/apps/${row.id}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    }).then(r => r.json())
    
    if (res.code === 0) {
      Object.assign(form, res.data.app)
      formConfigs.value = res.data.configs || []
    }
  } catch (error) {
    Object.assign(form, row)
    formConfigs.value = []
  }
  
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const data = {
          ...form,
          configs: formConfigs.value.filter(c => c.key)
        }
        
        if (currentAppId.value) {
          await updateApp(currentAppId.value, data)
          ElMessage.success('更新成功')
        } else {
          await createApp(data)
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
    
    const result: any = await runApp(row.id!)
    
    // 显示运行结果
    const status = result.status || 'unknown'
    let output = result.output || '无输出'
    const error = result.error || ''
    const duration = result.duration || 0
    
    // 尝试格式化输出为 JSON
    try {
      if (typeof output === 'string') {
        const parsed = JSON.parse(output)
        output = JSON.stringify(parsed, null, 2)
      } else if (typeof output === 'object') {
        output = JSON.stringify(output, null, 2)
      }
    } catch (e) {
      // 格式化失败，使用原始输出
    }
    
    if (status === 'success') {
      ElMessageBox.alert(
        `<div style="max-height: 400px; overflow-y: auto;">
          <p><strong>状态：</strong><span style="color: #67c23a;">成功</span></p>
          <p><strong>耗时：</strong>${duration}ms</p>
          <p><strong>输出：</strong></p>
          <pre style="background: #f5f7fa; padding: 12px; border-radius: 4px; white-space: pre-wrap; word-break: break-all; font-family: 'Courier New', monospace;">${output}</pre>
        </div>`,
        '运行结果',
        {
          dangerouslyUseHTMLString: true,
          confirmButtonText: '确定'
        }
      )
    } else {
      ElMessageBox.alert(
        `<div style="max-height: 400px; overflow-y: auto;">
          <p><strong>状态：</strong><span style="color: #f56c6c;">失败</span></p>
          <p><strong>耗时：</strong>${duration}ms</p>
          <p><strong>错误：</strong></p>
          <pre style="background: #fef0f0; padding: 12px; border-radius: 4px; color: #f56c6c; white-space: pre-wrap; word-break: break-all; font-family: 'Courier New', monospace;">${error}</pre>
        </div>`,
        '运行结果',
        {
          dangerouslyUseHTMLString: true,
          confirmButtonText: '确定',
          type: 'error'
        }
      )
    }
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

const formatOutput = (output: string) => {
  if (!output) return ''
  
  try {
    // 尝试解析为 JSON 并美化
    const parsed = JSON.parse(output)
    return JSON.stringify(parsed, null, 2)
  } catch (e) {
    // 解析失败，返回原始内容
    return output
  }
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

const handlePublish = (row: App) => {
  currentPublishAppId.value = row.id
  publishForm.display_name = row.name
  publishForm.description = row.description
  publishForm.version = '1.0.0'
  publishForm.category = ''
  publishForm.author = ''
  publishForm.tags = ''
  publishForm.icon = ''
  publishDialogVisible.value = true
}

const handlePublishSubmit = async () => {
  if (!publishFormRef.value || !currentPublishAppId.value) return
  
  await publishFormRef.value.validate(async (valid) => {
    if (valid) {
      publishing.value = true
      try {
        await publishToStore(currentPublishAppId.value!, publishForm)
        ElMessage.success('发布成功！应用已添加到应用商店')
        publishDialogVisible.value = false
      } catch (error: any) {
        const errorMsg = error.response?.data?.message || error.response?.data?.error || error.message || '发布失败'
        ElMessage.error(errorMsg)
        console.error('发布错误:', error)
      } finally {
        publishing.value = false
      }
    }
  })
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

.config-item {
  margin-bottom: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.log-container {
  min-height: 400px;
  max-height: 600px;
  overflow-y: auto;
}

.empty-logs {
  text-align: center;
  padding: 60px 0;
  color: #909399;
  font-size: 14px;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.log-item {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 12px;
  background: #fff;
  transition: all 0.3s;
}

.log-item:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.log-item.success {
  border-left: 3px solid #67c23a;
}

.log-item.error {
  border-left: 3px solid #f56c6c;
}

.log-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f5f7fa;
}

.log-time {
  color: #909399;
  font-size: 13px;
}

.log-duration {
  color: #606266;
  font-size: 13px;
  margin-left: auto;
}

.log-content {
  margin-top: 8px;
}

.log-label {
  color: #606266;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 4px;
}

.log-output, .log-error {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}

.log-error {
  background: #fef0f0;
  color: #f56c6c;
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
