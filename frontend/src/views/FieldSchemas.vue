<template>
  <div class="field-schemas">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>字段定义管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增字段
          </el-button>
        </div>
      </template>
      
      <el-table :data="fields" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="字段名称" width="150" />
        <el-table-column prop="key" label="字段Key" width="150" />
        <el-table-column label="字段类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeColor(row.type)">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column label="必填" width="80" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.required" color="#67c23a"><CircleCheck /></el-icon>
            <el-icon v-else color="#909399"><CircleClose /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="placeholder" label="占位符" width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
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
        <el-form-item label="字段名称" prop="name">
          <el-input v-model="form.name" placeholder="例如：标题" />
        </el-form-item>
        
        <el-form-item label="字段Key" prop="key">
          <el-input v-model="form.key" placeholder="例如：title（用于模板变量）" />
          <template #extra>
            <span style="color: #909399; font-size: 12px;">用于模板中的变量名，如 {{title}}</span>
          </template>
        </el-form-item>
        
        <el-form-item label="字段类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择">
            <el-option label="单行文本" value="text" />
            <el-option label="多行文本" value="textarea" />
            <el-option label="数字" value="number" />
            <el-option label="日期" value="date" />
            <el-option label="下拉选择" value="select" />
            <el-option label="URL" value="url" />
            <el-option label="邮箱" value="email" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" placeholder="字段说明" />
        </el-form-item>
        
        <el-form-item label="占位符">
          <el-input v-model="form.placeholder" placeholder="输入提示文字" />
        </el-form-item>
        
        <el-form-item label="默认值">
          <el-input v-model="form.default_value" placeholder="默认值" />
        </el-form-item>
        
        <el-form-item label="选项" v-if="form.type === 'select'">
          <el-input 
            v-model="form.options" 
            type="textarea" 
            :rows="3" 
            placeholder='["选项1", "选项2", "选项3"]'
          />
          <template #extra>
            <span style="color: #909399; font-size: 12px;">JSON 数组格式</span>
          </template>
        </el-form-item>
        
        <el-form-item label="必填">
          <el-switch v-model="form.required" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getFieldSchemas, createFieldSchema, updateFieldSchema, deleteFieldSchema } from '../api/field'
import type { FieldSchema, FieldType } from '../types'

const loading = ref(false)
const fields = ref<FieldSchema[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('新增字段')
const submitting = ref(false)
const formRef = ref<FormInstance>()
const currentFieldId = ref<number>()

const form = reactive<FieldSchema>({
  name: '',
  key: '',
  type: 'text',
  description: '',
  required: false,
  default_value: '',
  options: '',
  placeholder: '',
  validation: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入字段名称', trigger: 'blur' }],
  key: [
    { required: true, message: '请输入字段Key', trigger: 'blur' },
    { pattern: /^[a-zA-Z_][a-zA-Z0-9_]*$/, message: '只能包含字母、数字和下划线，且不能以数字开头', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择字段类型', trigger: 'change' }]
}

const getTypeLabel = (type: FieldType) => {
  const labels: Record<FieldType, string> = {
    text: '单行文本',
    textarea: '多行文本',
    number: '数字',
    date: '日期',
    select: '下拉选择',
    url: 'URL',
    email: '邮箱'
  }
  return labels[type] || type
}

const getTypeColor = (type: FieldType) => {
  const colors: Record<FieldType, string> = {
    text: '',
    textarea: 'success',
    number: 'warning',
    date: 'danger',
    select: 'info',
    url: 'primary',
    email: 'primary'
  }
  return colors[type] || ''
}

const loadFields = async () => {
  loading.value = true
  try {
    fields.value = await getFieldSchemas()
  } catch (error) {
    console.error('加载字段失败:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增字段'
  currentFieldId.value = undefined
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: FieldSchema) => {
  dialogTitle.value = '编辑字段'
  currentFieldId.value = row.id
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (currentFieldId.value) {
          await updateFieldSchema(currentFieldId.value, form)
          ElMessage.success('更新成功')
        } else {
          await createFieldSchema(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadFields()
      } catch (error) {
        console.error('提交失败:', error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDelete = async (row: FieldSchema) => {
  try {
    await ElMessageBox.confirm('确定要删除这个字段吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteFieldSchema(row.id!)
    ElMessage.success('删除成功')
    loadFields()
  } catch (error) {
    // 用户取消操作
  }
}

const resetForm = () => {
  Object.assign(form, {
    name: '',
    key: '',
    type: 'text',
    description: '',
    required: false,
    default_value: '',
    options: '',
    placeholder: '',
    validation: ''
  })
  formRef.value?.clearValidate()
}

onMounted(() => {
  loadFields()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
