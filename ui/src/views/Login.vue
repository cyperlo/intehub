<template>
  <div class="login-wrapper">
    <!-- 动态背景 -->
    <div class="bg">
      <div class="gradient-bg">
        <div class="gradient-blob blob1"></div>
        <div class="gradient-blob blob2"></div>
        <div class="gradient-blob blob3"></div>
      </div>
    </div>

    <div class="login-container">
      <div class="login-card">
        <!-- Logo 区域 -->
        <div class="login-header">
          <div class="logo-wrapper">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <h1 class="logo-text">InteHub</h1>
          </div>
          <p class="logo-subtitle">智能集成平台</p>
        </div>

        <!-- 表单区域 -->
        <el-form 
          :model="loginForm" 
          :rules="rules" 
          ref="formRef" 
          @keyup.enter="handleLogin"
          class="login-form"
        >
          <el-form-item prop="username">
            <div class="input-wrapper">
              <el-icon class="input-icon"><User /></el-icon>
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                size="large"
                class="custom-input"
              />
            </div>
          </el-form-item>
          
          <el-form-item prop="password">
            <div class="input-wrapper">
              <el-icon class="input-icon"><Lock /></el-icon>
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                size="large"
                show-password
                class="custom-input"
              />
            </div>
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              size="large"
              class="login-btn"
              :loading="loading"
              @click="handleLogin"
            >
              <span v-if="!loading">登 录</span>
              <span v-else>登录中...</span>
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 底部版权 -->
      <div class="copyright">
        <p>Powered by InteHub</p>
      </div>
    </div>

    <!-- 备案信息 -->
    <div class="footer">
      <a href="https://beian.miit.gov.cn/" target="_blank">京ICP备XXXXXXXX号-1</a>
      <span class="separator">|</span>
      <a href="http://www.beian.gov.cn/" target="_blank">京公网安备XXXXXXXXXXXXXXXX号</a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { login } from '../api/auth'
import { useAuthStore } from '../stores/auth'
import type { LoginForm } from '../types'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const formRef = ref<FormInstance>()

const loginForm = reactive<LoginForm>({
  username: '',
  password: ''
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await login(loginForm)
        
        if (!res || !res.token) {
          throw new Error('登录响应格式错误')
        }
        
        authStore.setToken(res.token)
        authStore.setUser(res.user_info)
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error: any) {
        console.error('登录失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

/* 动态渐变背景 */
.bg {
  position: absolute;
  inset: 0;
  z-index: 0;
}

.gradient-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
}

.gradient-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.6;
  animation: float 20s ease-in-out infinite;
}

.blob1 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  top: -10%;
  left: -10%;
  animation-delay: 0s;
}

.blob2 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  top: 60%;
  right: -5%;
  animation-delay: -5s;
}

.blob3 {
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  bottom: -20%;
  left: 30%;
  animation-delay: -10s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(30px, -30px) rotate(5deg);
  }
  50% {
    transform: translate(0, -50px) rotate(0deg);
  }
  75% {
    transform: translate(-30px, -30px) rotate(-5deg);
  }
}

/* 登录容器 */
.login-container {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.login-card {
  width: 420px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 48px 40px;
  box-shadow: 
    0 25px 50px -12px rgba(0, 0, 0, 0.25),
    0 0 0 1px rgba(255, 255, 255, 0.1);
}

/* Logo 区域 */
.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 8px;
}

.logo-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.logo-icon svg {
  width: 28px;
  height: 28px;
}

.logo-text {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.logo-subtitle {
  color: #6b7280;
  font-size: 14px;
  margin: 0;
  letter-spacing: 4px;
}

/* 表单区域 */
.login-form {
  margin-top: 0;
}

.input-wrapper {
  position: relative;
  width: 100%;
}

.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 1;
  color: #9ca3af;
  font-size: 18px;
}

:deep(.custom-input .el-input__wrapper) {
  padding: 4px 16px 4px 48px;
  border-radius: 12px;
  background: #f9fafb;
  border: 2px solid transparent;
  box-shadow: none !important;
  transition: all 0.3s ease;
}

:deep(.custom-input .el-input__wrapper:hover) {
  background: #f3f4f6;
}

:deep(.custom-input .el-input__wrapper.is-focus) {
  background: #fff;
  border-color: #667eea;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1) !important;
}

:deep(.custom-input .el-input__inner) {
  height: 48px;
  font-size: 15px;
}

:deep(.custom-input .el-input__inner::placeholder) {
  color: #9ca3af;
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  height: 52px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  border: none !important;
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 28px rgba(102, 126, 234, 0.5);
}

.login-btn:active {
  transform: translateY(0);
}

/* 底部版权 */
.copyright {
  margin-top: 24px;
  text-align: center;
}

.copyright p {
  color: rgba(255, 255, 255, 0.5);
  font-size: 12px;
  margin: 0;
}

/* 备案信息 */
.footer {
  position: fixed;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 12px;
}

.footer a {
  color: rgba(255, 255, 255, 0.4);
  font-size: 12px;
  text-decoration: none;
  transition: color 0.3s ease;
  white-space: nowrap;
}

.footer a:hover {
  color: rgba(255, 255, 255, 0.7);
}

.separator {
  color: rgba(255, 255, 255, 0.3);
  font-size: 12px;
}

/* 响应式 */
@media (max-width: 480px) {
  .login-card {
    width: 100%;
    max-width: 360px;
    padding: 32px 24px;
    border-radius: 20px;
  }
  
  .logo-icon {
    width: 40px;
    height: 40px;
    border-radius: 12px;
  }
  
  .logo-icon svg {
    width: 24px;
    height: 24px;
  }
  
  .logo-text {
    font-size: 28px;
  }
}
</style>
