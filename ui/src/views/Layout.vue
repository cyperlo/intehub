<template>
  <el-container class="layout-container">
    <!-- 移动端顶部栏 -->
    <el-header class="mobile-header" v-if="isMobile">
      <div class="mobile-header-content">
        <el-icon @click="drawerVisible = true" :size="24" class="menu-icon">
          <Menu />
        </el-icon>
        <div class="logo">InteHub</div>
        <el-dropdown>
          <el-icon :size="24"><User /></el-icon>
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
    </el-header>

    <!-- 移动端抽屉菜单 -->
    <el-drawer v-model="drawerVisible" direction="ltr" :size="250" v-if="isMobile">
      <template #header>
        <div class="drawer-logo">InteHub</div>
      </template>
      <el-menu
        :default-active="activeMenu"
        router
        @select="drawerVisible = false"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-sub-menu index="/apps-menu">
          <template #title>
            <el-icon><Box /></el-icon>
            <span>应用平台</span>
          </template>
          <el-menu-item index="/apps">应用管理</el-menu-item>
          <el-menu-item index="/appstore">应用商店</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/push-menu">
          <template #title>
            <el-icon><Promotion /></el-icon>
            <span>推送平台</span>
          </template>
          <el-menu-item index="/push-configs">集成配置</el-menu-item>
          <el-menu-item index="/push-history">集成历史</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/schedule-tasks">
          <el-icon><Clock /></el-icon>
          <span>定时任务</span>
        </el-menu-item>
        <el-sub-menu index="/system" v-if="authStore.user?.role === 'admin'">
          <template #title>
            <el-icon><Tools /></el-icon>
            <span>系统设置</span>
          </template>
          <el-menu-item index="/system/users">用户管理</el-menu-item>
          <el-menu-item index="/system/fields">字段定义</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-drawer>

    <!-- 桌面端侧边栏 -->
    <el-aside width="200px" v-if="!isMobile">
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
        <el-sub-menu index="/apps-menu">
          <template #title>
            <el-icon><Box /></el-icon>
            <span>应用平台</span>
          </template>
          <el-menu-item index="/apps">应用管理</el-menu-item>
          <el-menu-item index="/appstore">应用商店</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/push-menu">
          <template #title>
            <el-icon><Promotion /></el-icon>
            <span>推送平台</span>
          </template>
          <el-menu-item index="/push-configs">集成配置</el-menu-item>
          <el-menu-item index="/push-history">集成历史</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/schedule-tasks">
          <el-icon><Clock /></el-icon>
          <span>定时任务</span>
        </el-menu-item>
        <el-sub-menu index="/system" v-if="authStore.user?.role === 'admin'">
          <template #title>
            <el-icon><Tools /></el-icon>
            <span>系统设置</span>
          </template>
          <el-menu-item index="/system/users">用户管理</el-menu-item>
          <el-menu-item index="/system/fields">字段定义</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <el-container>
      <!-- 桌面端顶部栏 -->
      <el-header v-if="!isMobile">
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
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import { logout } from '../api/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const drawerVisible = ref(false)
const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value <= 768)
const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta.title || '')

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

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

.el-sub-menu {
  background-color: #304156;
}

.el-sub-menu :deep(.el-sub-menu__title) {
  color: #bfcbd9;
}

.el-sub-menu :deep(.el-sub-menu__title):hover {
  background-color: #263445 !important;
  color: #409eff !important;
}

.el-sub-menu :deep(.el-menu-item) {
  background-color: #1f2d3d !important;
  color: #bfcbd9;
  padding-left: 50px !important;
}

.el-sub-menu :deep(.el-menu-item):hover,
.el-sub-menu :deep(.el-menu-item.is-active) {
  background-color: #001528 !important;
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

/* 移动端样式 */
.mobile-header {
  background-color: #304156;
  border-bottom: none;
  padding: 0 16px;
  height: 56px !important;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.mobile-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.mobile-header-content .logo {
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  background: none;
  height: auto;
}

.mobile-header-content .el-icon {
  color: #fff;
}

.menu-icon {
  color: #fff;
  cursor: pointer;
}

.drawer-logo {
  font-size: 20px;
  font-weight: bold;
  color: #304156;
}

@media (max-width: 768px) {
  .layout-container {
    flex-direction: column;
  }
  
  .el-main {
    padding: 12px;
    padding-top: 68px;
  }
}
</style>
