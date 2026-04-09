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
          <el-menu-item index="/workflows">应用流</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/push-menu">
          <template #title>
            <el-icon><Promotion /></el-icon>
            <span>推送平台</span>
          </template>
          <el-menu-item index="/push-configs">推送配置</el-menu-item>
          <el-menu-item index="/push-history">推送历史</el-menu-item>
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
          <el-menu-item index="/workflows">应用流</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/push-menu">
          <template #title>
            <el-icon><Promotion /></el-icon>
            <span>推送平台</span>
          </template>
          <el-menu-item index="/push-configs">推送配置</el-menu-item>
          <el-menu-item index="/push-history">推送历史</el-menu-item>
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
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  color: #fff;
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.15);
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  background: rgba(255, 255, 255, 0.05);
  letter-spacing: 2px;
  position: relative;
}

.logo::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 20px;
  right: 20px;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
}

.el-menu {
  border: none;
  background: transparent;
  padding: 8px 0;
}

.el-menu-item {
  color: rgba(255, 255, 255, 0.7);
  margin: 4px 8px;
  border-radius: 10px;
  height: 48px;
  display: flex;
  align-items: center;
}

.el-menu-item:hover,
.el-menu-item.is-active {
  background: rgba(64, 158, 255, 0.15) !important;
  color: #fff !important;
}

.el-menu-item.is-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: linear-gradient(180deg, #409eff, #66b1ff);
  border-radius: 0 2px 2px 0;
}

.el-menu-item .el-icon {
  margin-right: 10px;
  font-size: 18px;
}

.el-sub-menu {
  margin: 4px 8px;
}

.el-sub-menu :deep(.el-sub-menu__title) {
  color: rgba(255, 255, 255, 0.7);
  margin: 4px 0;
  border-radius: 10px;
  height: 48px;
  display: flex;
  align-items: center;
}

.el-sub-menu :deep(.el-sub-menu__title:hover) {
  background: rgba(64, 158, 255, 0.15) !important;
  color: #fff !important;
}

.el-sub-menu :deep(.el-sub-menu__title) .el-icon {
  margin-right: 10px;
  font-size: 18px;
}

.el-sub-menu :deep(.el-sub-menu__arrow) {
  color: rgba(255, 255, 255, 0.5);
}

.el-sub-menu :deep(.el-menu) {
  background: transparent !important;
}

.el-sub-menu :deep(.el-menu-item) {
  background: transparent !important;
  color: rgba(255, 255, 255, 0.6);
  padding-left: 48px !important;
  height: 40px;
  margin: 2px 8px;
  border-radius: 8px;
}

.el-sub-menu :deep(.el-menu-item):hover,
.el-sub-menu :deep(.el-menu-item.is-active) {
  background: rgba(64, 158, 255, 0.1) !important;
  color: #fff !important;
}

.el-header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  padding: 0 24px;
  height: 64px !important;
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
  font-weight: 600;
  font-size: 18px;
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
  padding: 8px 16px;
  border-radius: 8px;
  transition: all 0.3s;
}

.el-dropdown-link:hover {
  background: #f5f7fa;
  color: #409eff;
}

.el-main {
  background: linear-gradient(180deg, #f5f7fa 0%, #eef1f5 100%);
  padding: 24px;
}

/* 移动端样式 */
.mobile-header {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border-bottom: none;
  padding: 0 16px;
  height: 56px !important;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
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
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
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
