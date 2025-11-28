<template>
  <el-container class="layout-container">
    <el-aside width="200px">
      <div class="logo">InteHub</div>
      <el-menu
        :default-active="activeMenu"
        router
        class="el-menu-vertical"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/field-schemas">
          <el-icon><Grid /></el-icon>
          <span>字段定义</span>
        </el-menu-item>
        <el-menu-item index="/push-configs">
          <el-icon><Setting /></el-icon>
          <span>推送配置</span>
        </el-menu-item>
        <el-menu-item index="/push-history">
          <el-icon><List /></el-icon>
          <span>推送历史</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header>
        <div class="header-content">
          <h3>{{ currentTitle }}</h3>
          <div class="user-info">
            <el-dropdown>
              <span class="el-dropdown-link">
                <el-icon><User /></el-icon>
                {{ authStore.user?.nickname || authStore.user?.username }}
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleLogout">
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-header>
      
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import { logout } from '../api/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta.title || '')

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await logout()
    authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch (error) {
    // 用户取消操作
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.el-aside {
  background-color: #304156;
  color: #fff;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
  color: #fff;
  background-color: #282f3d;
}

.el-menu {
  border: none;
  background-color: #304156;
}

.el-menu-item {
  color: #bfcbd9;
}

.el-menu-item:hover,
.el-menu-item.is-active {
  background-color: #263445 !important;
  color: #409eff !important;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  padding: 0 20px;
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h3 {
  margin: 0;
  color: #303133;
}

.user-info {
  display: flex;
  align-items: center;
}

.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
}

.el-dropdown-link:hover {
  color: #409eff;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
