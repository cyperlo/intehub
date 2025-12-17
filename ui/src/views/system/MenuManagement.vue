<template>
  <div class="menu-management">
    <div class="header-actions">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        <span class="btn-text">新增菜单</span>
      </el-button>
    </div>

    <!-- 桌面端表格 -->
    <el-table :data="menus" v-loading="loading" class="desktop-table">
      <el-table-column prop="name" label="菜单名称" min-width="120" />
      <el-table-column prop="path" label="路径" min-width="150" />
      <el-table-column prop="icon" label="图标" width="100" />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column label="可见" width="80">
        <template #default="{ row }">
          <el-tag :type="row.visible ? 'success' : 'info'" size="small">
            {{ row.visible ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button type="warning" size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 移动端卡片列表 -->
    <div class="mobile-list" v-loading="loading">
      <el-card v-for="menu in menus" :key="menu.id" class="menu-card" shadow="hover">
        <div class="menu-info">
          <div class="info-row">
            <span class="label">菜单名称：</span>
            <span class="value">{{ menu.name }}</span>
          </div>
          <div class="info-row">
            <span class="label">路径：</span>
            <span class="value">{{ menu.path }}</span>
          </div>
          <div class="info-row">
            <span class="label">图标：</span>
            <span class="value">{{ menu.icon }}</span>
          </div>
          <div class="info-row">
            <span class="label">排序：</span>
            <span class="value">{{ menu.sort }}</span>
          </div>
          <div class="info-row">
            <span class="label">可见：</span>
            <el-tag :type="menu.visible ? 'success' : 'info'" size="small">
              {{ menu.visible ? '是' : '否' }}
            </el-tag>
          </div>
        </div>
        <div class="card-actions">
          <el-button type="warning" size="small" @click="handleEdit(menu)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDelete(menu)">删除</el-button>
        </div>
      </el-card>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" :width="isMobile ? '95%' : '500px'" :fullscreen="isMobile">
      <el-form :model="form" ref="formRef" label-width="80px">
        <el-form-item label="菜单名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="路径">
          <el-input v-model="form.path" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="可见">
          <el-switch v-model="form.visible" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMenus, createMenu, updateMenu, deleteMenu, type Menu } from '../../api/system'

const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value <= 768)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

const loading = ref(false)
const menus = ref<Menu[]>([])
const dialogVisible = ref(false)
const dialogTitle = ref('新增菜单')
const currentMenuId = ref<number>()

const form = reactive<Menu>({
  name: '',
  path: '',
  icon: '',
  parent_id: 0,
  sort: 0,
  visible: true,
  roles: '["admin","user"]'
})

const loadMenus = async () => {
  loading.value = true
  try {
    menus.value = await getMenus()
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增菜单'
  currentMenuId.value = undefined
  Object.assign(form, {
    name: '',
    path: '',
    icon: '',
    parent_id: 0,
    sort: 0,
    visible: true,
    roles: '["admin","user"]'
  })
  dialogVisible.value = true
}

const handleEdit = (row: Menu) => {
  dialogTitle.value = '编辑菜单'
  currentMenuId.value = row.id
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    if (currentMenuId.value) {
      await updateMenu(currentMenuId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createMenu(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadMenus()
  } catch (error) {
    console.error('提交失败:', error)
  }
}

const handleDelete = async (row: Menu) => {
  try {
    await ElMessageBox.confirm('确定要删除这个菜单吗?', '提示', {
      type: 'warning'
    })
    await deleteMenu(row.id!)
    ElMessage.success('删除成功')
    loadMenus()
  } catch (error) {
    // 用户取消
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadMenus()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.menu-management {
  width: 100%;
}

.header-actions {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

.desktop-table {
  width: 100%;
}

.mobile-list {
  display: none;
}

.menu-card {
  margin-bottom: 12px;
}

.menu-info {
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
  word-break: break-all;
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
}
</style>
