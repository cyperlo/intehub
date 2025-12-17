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
        <el-table-column prop="last_run_at" label="上次运行" width="160" />
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
              <span class="value">{{ task.last_run_at }}</span>
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
          <el-select v-model="form.task_type">
            <el-option label="推送任务" value="push" />
          </el-select>
        </el-form-item>
        <el-form-item label="关联配置" prop="config_id">
          <el-select v-model="form.config_id" placeholder="选择推送配置">
            <el-option 
              v-for="config in pushConfigs" 
              :key="config.id" 
              :label="config.name" 
              :value="config.id" 
            />
          </el-select>
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

const windowWidth = ref(window.innerWidth)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const loading = ref(false)
const tasks = ref<ScheduleTask[]>([])
const pushConfigs = ref<any[]>([])
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

const form = reactive<ScheduleTask>({
  name: '',
  description: '',
  cron_expr: '',
  task_type: 'push',
  config_id: 0,
  enabled: false
})

const isMobile = computed(() => windowWidth.value <= 768)

const rules: FormRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  cron_expr: [{ required: true, message: '请输入Cron表达式', trigger: 'blur' }],
  task_type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  config_id: [{ required: true, message: '请选择关联配置', trigger: 'change' }]
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
    console.error('加载推送配置失败:', error)
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增任务'
  currentTaskId.value = undefined
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: ScheduleTask) => {
  dialogTitle.value = '编辑任务'
  currentTaskId.value = row.id
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (currentTaskId.value) {
          await updateScheduleTask(currentTaskId.value, form)
          ElMessage.success('更新成功')
        } else {
          await createScheduleTask(form)
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
    enabled: false
  })
  formRef.value?.clearValidate()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadTasks()
  loadPushConfigs()
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
</style>
