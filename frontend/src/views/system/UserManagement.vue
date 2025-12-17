<template>
  <div class="user-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="header-title">用户管理</span>
          <el-button type="primary" @click="handleAdd" size="default">
            <el-icon><Plus /></el-icon>
            <span class="btn-text">新增用户</span>
          </el-button>
        </div>
      </template>

      <!-- 桌面端表格 -->
      <el-table :data="users" v-loading="loading" class="desktop-table">
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : ''" size="small">
              {{ row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" min-width="160" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
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
        <el-card v-for="user in users" :key="user.id" class="user-card" shadow="hover">
          <div class="user-info">
            <div class="info-row">
              <span class="label">用户名：</span>
              <span class="value">{{ user.username }}</span>
            </div>
            <div class="info-row">
              <span class="label">昵称：</span>
              <span class="value">{{ user.nickname }}</span>
            </div>
            <div class="info-row">
              <span class="label">角色：</span>
              <el-tag :type="user.role === 'admin' ? 'danger' : ''" size="small">
                {{ user.role === 'admin' ? '管理员' : '普通用户' }}
              </el-tag>
            </div>
            <div class="info-row">
              <span class="label">创建时间：</span>
              <span class="value">{{ user.created_at }}</span>
            </div>
          </div>
          <div class="card-actions">
            <el-button type="warning" size="small" @click="handleEdit(user)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(user)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
        </el-card>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" :width="isMobile ? '95%' : '500px'" :fullscreen="isMobile">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="!!currentUserId" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" :placeholder="currentUserId ? '不修改请留空' : '请输入密码'" />
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
import { ref, onMounted, onUnmounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getUsers, createUser, updateUser, deleteUser, type User } from '../../api/system'

const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value <= 768)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const loading = ref(false)
const users = ref<User[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('新增用户')
const submitting = ref(false)
const formRef = ref<FormInstance>()
const currentUserId = ref<number>()

const form = reactive<User>({
  username: '',
  nickname: '',
  role: 'user',
  password: ''
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  password: [
    { 
      validator: (rule, value, callback) => {
        if (!currentUserId.value && !value) {
          callback(new Error('请输入密码'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
}

const loadUsers = async () => {
  loading.value = true
  try {
    users.value = await getUsers()
  } catch (error) {
    console.error('加载用户失败:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增用户'
  currentUserId.value = undefined
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  dialogTitle.value = '编辑用户'
  currentUserId.value = row.id
  Object.assign(form, { ...row, password: '' })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (currentUserId.value) {
          const updateData: any = {
            nickname: form.nickname,
            role: form.role
          }
          if (form.password) {
            updateData.password = form.password
          }
          await updateUser(currentUserId.value, updateData)
          ElMessage.success('更新成功')
        } else {
          await createUser(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadUsers()
      } catch (error) {
        console.error('提交失败:', error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDelete = async (row: User) => {
  try {
    await ElMessageBox.confirm('确定要删除这个用户吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteUser(row.id!)
    ElMessage.success('删除成功')
    loadUsers()
  } catch (error) {
    // 用户取消操作
  }
}

const resetForm = () => {
  Object.assign(form, {
    username: '',
    nickname: '',
    role: 'user',
    password: ''
  })
  formRef.value?.clearValidate()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadUsers()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.user-management {
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

.user-card {
  margin-bottom: 12px;
}

.user-info {
  margin-bottom: 12px;
}

.info-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
}

.info-row .label {
  color: #909399;
  min-width: 80px;
}

.info-row .value {
  color: #303133;
  flex: 1;
}

.card-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
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
}
</style>
