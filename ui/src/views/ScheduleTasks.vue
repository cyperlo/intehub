<template>
  <div class="schedule-tasks">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="header-title">定时任务管理</span>
          <el-button type="primary" @click="handleAdd" size="default">
            <el-icon><Plus /></el-icon>
            <span class="btn-text">新增任务</span>
          </el-button>
        </div>
      </template>

      <!-- 桌面端表格 -->
      <el-table :data="tasks" v-loading="loading" class="desktop-table">
        <el-table-column prop="name" label="任务名称" min-width="120" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip min-width="150" />
        <el-table-column prop="cron_expr" label="Cron表达式" width="130" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="上次运行" width="180">
          <template #default="{ row }">
            {{ formatTime(row.last_run_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button 
              :type="row.enabled ? 'warning' : 'success'" 
              size="small" 
              @click="handleToggle(row)"
            >
              {{ row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button type="primary" size="small" @click="handleViewLogs(row)">日志</el-button>
            <el-button type="warning" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 移动端卡片列表 -->
      <div class="mobile-list" v-loading="loading">
        <el-card v-for="task in tasks" :key="task.id" class="task-card" shadow="hover">
          <div class="task-header">
            <h4>{{ task.name }}</h4>
            <el-tag :type="task.enabled ? 'success' : 'info'" size="small">
              {{ task.enabled ? '启用' : '禁用' }}
            </el-tag>
          </div>
          <div class="task-info">
            <div class="info-row" v-if="task.description">
              <span class="label">描述：</span>
              <span class="value">{{ task.description }}</span>
            </div>
            <div class="info-row">
              <span class="label">Cron：</span>
              <span class="value">{{ task.cron_expr }}</span>
            </div>
            <div class="info-row" v-if="task.last_run_at">
              <span class="label">上次运行：</span>
              <span class="value">{{ formatTime(task.last_run_at) }}</span>
            </div>
          </div>
          <div class="card-actions">
            <el-button 
              :type="task.enabled ? 'warning' : 'success'" 
              size="small" 
              @click="handleToggle(task)"
            >
              {{ task.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button type="primary" size="small" @click="handleViewLogs(task)">日志</el-button>
            <el-button type="warning" size="small" @click="handleEdit(task)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(task)">删除</el-button>
          </div>
        </el-card>
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" :width="isMobile ? '95%' : '600px'" :fullscreen="isMobile">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="Cron表达式" prop="cron_expr">
          <el-input v-model="form.cron_expr" placeholder="例如: 0 */5 * * * (每5分钟)" />
          <div style="color: #909399; font-size: 12px; margin-top: 4px;">
            格式: 秒 分 时 日 月 周，例如 "0 0 * * * *" 表示每小时执行
          </div>
        </el-form-item>
        <el-form-item label="任务类型" prop="task_type">
          <el-select v-model="form.task_type" @change="handleTaskTypeChange">
            <el-option label="集成任务" value="push" />
            <el-option label="应用任务" value="app" />
            <el-option label="工作流任务" value="workflow" />
          </el-select>
        </el-form-item>
        <el-form-item label="关联配置" prop="config_id" v-if="form.task_type === 'push'">
          <el-select v-model="form.config_id" placeholder="选择集成配置" @change="handleConfigChange">
            <el-option 
              v-for="config in pushConfigs" 
              :key="config.id" 
              :label="config.name" 
              :value="config.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关联应用" prop="app_id" v-if="form.task_type === 'app'">
          <el-select v-model="form.app_id" placeholder="选择应用">
            <el-option 
              v-for="app in appList" 
              :key="app.id" 
              :label="app.name" 
              :value="app.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关联工作流" prop="workflow_id" v-if="form.task_type === 'workflow'">
          <el-select v-model="form.workflow_id" placeholder="选择工作流">
            <el-option 
              v-for="workflow in workflowList" 
              :key="workflow.id" 
              :label="workflow.name" 
              :value="workflow.id" 
            />
          </el-select>
        </el-form-item>
        
        <!-- 字段数据配置 -->
        <el-divider v-if="configFields.length > 0">字段数据配置</el-divider>
        <div v-if="configFields.length > 0" class="field-data-section">
          <el-form-item 
            v-for="field in configFields" 
            :key="field.id"
            :label="field.name"
            :required="field.required"
          >
            <el-input 
              v-if="field.type === 'text' || field.type === 'url' || field.type === 'email'"
              v-model="fieldDataForm[field.key]"
              :placeholder="field.placeholder || field.description"
            />
            <el-input 
              v-else-if="field.type === 'textarea'"
              v-model="fieldDataForm[field.key]"
              type="textarea"
              :rows="3"
              :placeholder="field.placeholder || field.description"
            />
            <el-input-number 
              v-else-if="field.type === 'number'"
              v-model="fieldDataForm[field.key]"
              style="width: 100%"
            />
            <el-date-picker 
              v-else-if="field.type === 'date'"
              v-model="fieldDataForm[field.key]"
              type="date"
              style="width: 100%"
            />
            <el-select 
              v-else-if="field.type === 'select'"
              v-model="fieldDataForm[field.key]"
              style="width: 100%"
            >
              <el-option 
                v-for="(opt, idx) in parseOptions(field.options)" 
                :key="idx"
                :label="opt"
                :value="opt"
              />
            </el-select>
          </el-form-item>
        </div>
        
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
    <el-dialog v-model="logDialogVisible" title="任务执行日志" :width="isMobile ? '95%' : '900px'" :fullscreen="isMobile">
      <el-table :data="logs" v-loading="logLoading">
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="信息" show-overflow-tooltip />
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column label="开始时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.started_at) }}
          </template>
        </el-table-column>
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
import { 
  getScheduleTasks, 
  createScheduleTask, 
  updateScheduleTask, 
  deleteScheduleTask,
  toggleScheduleTask,
  getScheduleLogs,
  type ScheduleTask,
  type ScheduleLog
} from '../api/schedule'
import { getPushConfigs } from '../api/push'
import { getWorkflows, type Workflow } from '../api/workflow'

const windowWidth = ref(window.innerWidth)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const loading = ref(false)
const tasks = ref<ScheduleTask[]>([])
const pushConfigs = ref<any[]>([])
const workflowList = ref<Workflow[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('新增任务')
const submitting = ref(false)
const formRef = ref<FormInstance>()
const currentTaskId = ref<number>()

const logDialogVisible = ref(false)
const logLoading = ref(false)
const logs = ref<ScheduleLog[]>([])
const logPage = ref(1)
const logPageSize = ref(20)
const logTotal = ref(0)
const currentLogTaskId = ref<number>()

// 字段相关
const configFields = ref<any[]>([])
const fieldDataForm = reactive<Record<string, any>>({})

// 应用列表
const appList = ref<any[]>([])

const form = reactive<ScheduleTask>({
  name: '',
  description: '',
  cron_expr: '',
  task_type: 'push',
  config_id: 0,
  app_id: 0,
  workflow_id: 0,
  enabled: false
})

const isMobile = computed(() => windowWidth.value <= 768)

const rules: FormRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  cron_expr: [{ required: true, message: '请输入Cron表达式', trigger: 'blur' }],
  task_type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  config_id: [
    { 
      validator: (rule, value, callback) => {
        if (form.task_type === 'push' && !value) {
          callback(new Error('请选择集成配置'))
        } else {
          callback()
        }
      }, 
      trigger: 'change' 
    }
  ],
  app_id: [
    { 
      validator: (rule, value, callback) => {
        if (form.task_type === 'app' && !value) {
          callback(new Error('请选择应用'))
        } else {
          callback()
        }
      }, 
      trigger: 'change' 
    }
  ],
  workflow_id: [
    { 
      validator: (rule, value, callback) => {
        if (form.task_type === 'workflow' && !value) {
          callback(new Error('请选择工作流'))
        } else {
          callback()
        }
      }, 
      trigger: 'change' 
    }
  ]
}

const formatTime = (time: string | null | undefined) => {
  if (!time) return '-'
  try {
    return new Date(time).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
  } catch {
    return time
  }
}

const loadTasks = async () => {
  loading.value = true
  try {
    tasks.value = await getScheduleTasks()
  } finally {
    loading.value = false
  }
}

const loadPushConfigs = async () => {
  try {
    pushConfigs.value = await getPushConfigs()
  } catch (error) {
    console.error('加载集成配置失败:', error)
  }
}

const loadWorkflows = async () => {
  try {
    workflowList.value = await getWorkflows()
  } catch (error) {
    console.error('加载工作流失败:', error)
  }
}

const loadApps = async () => {
  try {
    const { getApps } = await import('../api/app')
    appList.value = await getApps()
  } catch (error) {
    console.error('加载应用失败:', error)
  }
}

const handleTaskTypeChange = () => {
  // 切换任务类型时清空相关字段
  form.config_id = 0
  form.app_id = 0
  configFields.value = []
  Object.keys(fieldDataForm).forEach(key => delete fieldDataForm[key])
}

const handleConfigChange = async (configId: number) => {
  if (!configId) {
    configFields.value = []
    return
  }
  
  try {
    const { getConfigFields } = await import('../api/field')
    configFields.value = await getConfigFields(configId)
    
    // 清空字段数据
    Object.keys(fieldDataForm).forEach(key => delete fieldDataForm[key])
    
    // 设置默认值
    configFields.value.forEach((field: any) => {
      if (field.default_value) {
        fieldDataForm[field.key] = field.default_value
      }
    })
  } catch (error) {
    console.error('加载字段失败:', error)
  }
}

const parseOptions = (optionsStr: string) => {
  try {
    return JSON.parse(optionsStr || '[]')
  } catch {
    return []
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增任务'
  currentTaskId.value = undefined
  resetForm()
  configFields.value = []
  Object.keys(fieldDataForm).forEach(key => delete fieldDataForm[key])
  dialogVisible.value = true
}

const handleEdit = async (row: ScheduleTask) => {
  dialogTitle.value = '编辑任务'
  currentTaskId.value = row.id
  Object.assign(form, row)
  
  // 加载配置的字段
  if (row.config_id) {
    await handleConfigChange(row.config_id)
    
    // 加载已保存的字段数据
    if (row.field_data) {
      try {
        const savedData = JSON.parse(row.field_data)
        Object.assign(fieldDataForm, savedData)
      } catch (error) {
        console.error('解析字段数据失败:', error)
      }
    }
  }
  
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        // 将字段数据转换为 JSON 字符串
        const taskData = {
          ...form,
          field_data: JSON.stringify(fieldDataForm)
        }
        
        if (currentTaskId.value) {
          await updateScheduleTask(currentTaskId.value, taskData)
          ElMessage.success('更新成功')
        } else {
          await createScheduleTask(taskData)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadTasks()
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleToggle = async (row: ScheduleTask) => {
  try {
    await toggleScheduleTask(row.id!)
    ElMessage.success(row.enabled ? '已禁用' : '已启用')
    loadTasks()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: ScheduleTask) => {
  try {
    await ElMessageBox.confirm('确定要删除这个任务吗?', '提示', {
      type: 'warning'
    })
    await deleteScheduleTask(row.id!)
    ElMessage.success('删除成功')
    loadTasks()
  } catch (error) {
    // 用户取消
  }
}

const handleViewLogs = (row: ScheduleTask) => {
  currentLogTaskId.value = row.id
  logPage.value = 1
  loadLogs()
  logDialogVisible.value = true
}

const loadLogs = async () => {
  if (!currentLogTaskId.value) return
  
  logLoading.value = true
  try {
    const res = await getScheduleLogs({
      task_id: currentLogTaskId.value,
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
    cron_expr: '',
    task_type: 'push',
    config_id: 0,
    app_id: 0,
    enabled: false
  })
  configFields.value = []
  Object.keys(fieldDataForm).forEach(key => delete fieldDataForm[key])
  formRef.value?.clearValidate()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadTasks()
  loadPushConfigs()
  loadWorkflows()
  loadApps()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.schedule-tasks {
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

.task-card {
  margin-bottom: 12px;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.task-header h4 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.task-info {
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
  
  .card-header {
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .header-title {
    width: 100%;
    margin-bottom: 8px;
  }
}

.field-data-section {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
  margin-bottom: 16px;
}
</style>
