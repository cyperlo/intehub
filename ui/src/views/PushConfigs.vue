<template>
  <div class="push-configs">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="header-title">推送配置管理</span>
          <el-button type="primary" @click="handleAdd" size="default">
            <el-icon><Plus /></el-icon>
            <span class="btn-text">新增配置</span>
          </el-button>
        </div>
      </template>
      
      <!-- 桌面端表格 -->
      <el-table :data="configs" v-loading="loading" class="desktop-table">
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip min-width="150" />
        <el-table-column prop="url" label="URL" show-overflow-tooltip min-width="200" />
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" @change="handleToggleStatus(row)" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="handleConfigFields(row)" link>
              <el-icon><Grid /></el-icon>
              字段
            </el-button>
            <el-button type="info" size="small" @click="handleTest(row)" link>
              <el-icon><Promotion /></el-icon>
              测试
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
        <el-card v-for="config in configs" :key="config.id" class="config-card" shadow="hover">
          <div class="config-header">
            <div class="config-name-row">
              <h4>{{ config.name }}</h4>
              <el-switch v-model="config.enabled" @change="handleToggleStatus(config)" size="small" />
            </div>
            <div class="config-method">{{ config.method }}</div>
          </div>
          <div class="config-info">
            <div class="info-row" v-if="config.description">
              <span class="label">描述：</span>
              <span class="value">{{ config.description }}</span>
            </div>
            <div class="info-row">
              <span class="label">URL：</span>
              <span class="value">{{ config.url }}</span>
            </div>
          </div>
          <div class="card-actions">
            <el-button type="success" size="small" @click="handleConfigFields(config)">
              <el-icon><Grid /></el-icon>
              字段
            </el-button>
            <el-button type="info" size="small" @click="handleTest(config)">
              <el-icon><Promotion /></el-icon>
              测试
            </el-button>
            <el-button type="warning" size="small" @click="handleEdit(config)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(config)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
        </el-card>
      </div>
    </el-card>
    
    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      :width="isMobile ? '95%' : '600px'"
      :fullscreen="isMobile"
      @close="resetForm"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="配置名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入配置名称" />
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
        
        <el-form-item label="集成URL" prop="url">
          <el-input v-model="form.url" placeholder="https://example.com/webhook" />
        </el-form-item>
        
        <el-form-item label="请求方法" prop="method">
          <el-select v-model="form.method" placeholder="请选择">
            <el-option label="POST" value="POST" />
            <el-option label="GET" value="GET" />
            <el-option label="PUT" value="PUT" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="内容类型" prop="content_type">
          <el-select v-model="form.content_type" placeholder="请选择">
            <el-option label="application/json" value="application/json" />
            <el-option label="application/x-www-form-urlencoded" value="application/x-www-form-urlencoded" />
            <el-option label="text/plain" value="text/plain" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="请求头" prop="headers">
          <el-input 
            v-model="form.headers" 
            type="textarea" 
            :rows="3" 
            placeholder='{"Authorization": "Bearer token"}'
          />
        </el-form-item>
        
        <el-form-item label="请求模板" prop="template">
          <el-input 
            v-model="form.template" 
            type="textarea" 
            :rows="5" 
            placeholder='{"title": "{{title}}", "content": "{{content}}"}'
          />
        </el-form-item>
        
        <el-form-item label="启用状态">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
    
    <!-- 字段配置对话框 -->
    <el-dialog v-model="fieldDialogVisible" title="配置字段" :width="isMobile ? '95%' : '800px'" :fullscreen="isMobile">
      <div style="margin-bottom: 16px;">
        <el-alert type="info" :closable="false">
          选择此配置使用的字段，执行时将根据这些字段生成表单
        </el-alert>
      </div>
      
      <el-transfer
        v-model="selectedFieldIds"
        :data="allFields"
        :titles="['可用字段', '已选字段']"
        :props="{
          key: 'id',
          label: 'name'
        }"
        filterable
        filter-placeholder="搜索字段"
        class="custom-transfer"
      />
      
      <template #footer>
        <el-button @click="fieldDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveFields" :loading="savingFields">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- 测试推送对话框 -->
    <el-dialog v-model="testDialogVisible" title="测试推送" :width="isMobile ? '95%' : '600px'" :fullscreen="isMobile">
      <div v-if="testFields.length > 0">
        <el-form :model="testFormData" label-width="100px">
          <el-form-item 
            v-for="field in testFields" 
            :key="field.id"
            :label="field.name"
            :required="field.required"
          >
            <el-input 
              v-if="field.type === 'text' || field.type === 'url' || field.type === 'email'"
              v-model="testFormData[field.key]"
              :placeholder="field.placeholder || field.description"
            />
            <el-input 
              v-else-if="field.type === 'textarea'"
              v-model="testFormData[field.key]"
              type="textarea"
              :rows="3"
              :placeholder="field.placeholder || field.description"
            />
            <el-input-number 
              v-else-if="field.type === 'number'"
              v-model="testFormData[field.key]"
              style="width: 100%"
            />
            <el-date-picker 
              v-else-if="field.type === 'date'"
              v-model="testFormData[field.key]"
              type="date"
              style="width: 100%"
            />
            <el-select 
              v-else-if="field.type === 'select'"
              v-model="testFormData[field.key]"
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
        </el-form>
      </div>
      <div v-else>
        <el-empty description="请先配置字段" />
      </div>
      
      <template #footer>
        <el-button @click="testDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleTestSubmit" :loading="testing" :disabled="testFields.length === 0">
          发送测试
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getPushConfigs, createPushConfig, updatePushConfig, deletePushConfig, sendPush } from '../api/push'
import { getFieldSchemas, getConfigFields, updateConfigFields } from '../api/field'
import type { PushConfig, FieldSchema } from '../types'

const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value <= 768)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const loading = ref(false)
const configs = ref<PushConfig[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('新增配置')
const submitting = ref(false)
const testDialogVisible = ref(false)
const testing = ref(false)
const formRef = ref<FormInstance>()
const currentConfigId = ref<number>()
const testConfigId = ref<number>()

// 字段相关
const fieldDialogVisible = ref(false)
const allFields = ref<FieldSchema[]>([])
const selectedFieldIds = ref<number[]>([])
const currentFieldConfigId = ref<number>()
const savingFields = ref(false)
const testFields = ref<FieldSchema[]>([])
const testFormData = reactive<Record<string, any>>({})

const form = reactive<PushConfig>({
  name: '',
  description: '',
  url: '',
  method: 'POST',
  headers: '',
  content_type: 'application/json',
  template: '',
  enabled: true
})


const rules: FormRules = {
  name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入集成URL', trigger: 'blur' }],
  method: [{ required: true, message: '请选择请求方法', trigger: 'change' }],
  template: [{ required: true, message: '请输入请求模板', trigger: 'blur' }]
}

const loadConfigs = async () => {
  loading.value = true
  try {
    configs.value = await getPushConfigs()
  } catch (error) {
    console.error('加载配置失败:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增配置'
  currentConfigId.value = undefined
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: PushConfig) => {
  dialogTitle.value = '编辑配置'
  currentConfigId.value = row.id
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (currentConfigId.value) {
          await updatePushConfig(currentConfigId.value, form)
          ElMessage.success('更新成功')
        } else {
          await createPushConfig(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadConfigs()
      } catch (error) {
        console.error('提交失败:', error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleToggleStatus = async (row: PushConfig) => {
  try {
    await updatePushConfig(row.id!, row)
    ElMessage.success(row.enabled ? '已启用' : '已禁用')
  } catch (error) {
    row.enabled = !row.enabled
    console.error('更新状态失败:', error)
  }
}

const handleDelete = async (row: PushConfig) => {
  try {
    await ElMessageBox.confirm('确定要删除这个配置吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deletePushConfig(row.id!)
    ElMessage.success('删除成功')
    loadConfigs()
  } catch (error) {
    // 用户取消操作
  }
}

const handleConfigFields = async (row: PushConfig) => {
  currentFieldConfigId.value = row.id
  
  // 加载所有字段
  try {
    allFields.value = await getFieldSchemas()
    // 加载当前配置的字段
    const configFields = await getConfigFields(row.id!)
    selectedFieldIds.value = configFields.map(f => f.id!)
    fieldDialogVisible.value = true
  } catch (error) {
    console.error('加载字段失败:', error)
  }
}

const handleSaveFields = async () => {
  savingFields.value = true
  try {
    await updateConfigFields(currentFieldConfigId.value!, selectedFieldIds.value)
    ElMessage.success('保存成功')
    fieldDialogVisible.value = false
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    savingFields.value = false
  }
}

const handleTest = async (row: PushConfig) => {
  testConfigId.value = row.id
  
  // 加载配置的字段
  try {
    testFields.value = await getConfigFields(row.id!)
    // 重置表单数据
    Object.keys(testFormData).forEach(key => delete testFormData[key])
    // 设置默认值
    testFields.value.forEach(field => {
      if (field.default_value) {
        testFormData[field.key] = field.default_value
      }
    })
    testDialogVisible.value = true
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

const handleTestSubmit = async () => {
  testing.value = true
  try {
    await sendPush({ config_id: testConfigId.value!, data: testFormData })
    ElMessage.success('推送成功')
    testDialogVisible.value = false
  } catch (error: any) {
    console.error('推送失败:', error)
  } finally {
    testing.value = false
  }
}

const resetForm = () => {
  Object.assign(form, {
    id: undefined,
    name: '',
    description: '',
    url: '',
    method: 'POST',
    headers: '',
    content_type: 'application/json',
    template: '',
    enabled: true
  })
  formRef.value?.clearValidate()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadConfigs()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.push-configs {
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

.config-card {
  margin-bottom: 12px;
}

.config-header {
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.config-name-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.config-name-row h4 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.config-method {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
}

.config-info {
  margin-bottom: 12px;
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

.card-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  padding-top: 8px;
  border-top: 1px solid #ebeef5;
}

/* 穿梭框样式修复 */
:deep(.el-transfer) {
  display: flex;
  flex-direction: row;
  align-items: flex-start;
}

:deep(.el-transfer__buttons) {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 16px;
}

:deep(.el-transfer__button) {
  margin: 4px 0;
}

:deep(.el-transfer-panel) {
  flex: 1;
  max-width: 300px;
}

:deep(.el-transfer-panel__body) {
  height: 300px;
}

:deep(.el-transfer-panel__list) {
  height: 100%;
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
}
</style>
