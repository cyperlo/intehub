<template>
  <div class="push-configs">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>推送配置管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增配置
          </el-button>
        </div>
      </template>
      
      <el-table :data="configs" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="url" label="URL" show-overflow-tooltip />
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" @change="handleToggleStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="handleConfigFields(row)">
              <el-icon><Grid /></el-icon>
              字段
            </el-button>
            <el-button type="primary" size="small" @click="handleTest(row)">
              <el-icon><Promotion /></el-icon>
              测试
            </el-button>
            <el-button type="warning" size="small" @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="resetForm"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="配置名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入配置名称" />
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
        
        <el-form-item label="推送URL" prop="url">
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
        
        <el-form-item label="推送模板" prop="template">
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
    <el-dialog v-model="fieldDialogVisible" title="配置字段" width="600px">
      <div style="margin-bottom: 16px;">
        <el-alert type="info" :closable="false">
          选择此配置使用的字段，推送时将根据这些字段生成表单
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
      />
      
      <template #footer>
        <el-button @click="fieldDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveFields" :loading="savingFields">保存</el-button>
      </template>
    </el-dialog>
    
    <!-- 测试推送对话框 -->
    <el-dialog v-model="testDialogVisible" title="测试推送" width="600px">
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
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getPushConfigs, createPushConfig, updatePushConfig, deletePushConfig, sendPush } from '../api/push'
import { getFieldSchemas, getConfigFields, updateConfigFields } from '../api/field'
import type { PushConfig, FieldSchema } from '../types'

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
  url: [{ required: true, message: '请输入推送URL', trigger: 'blur' }],
  method: [{ required: true, message: '请选择请求方法', trigger: 'change' }],
  template: [{ required: true, message: '请输入推送模板', trigger: 'blur' }]
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
  loadConfigs()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
