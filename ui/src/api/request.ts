import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import router from '../router'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data
    // 统一返回格式: { code: 0, message: "Succeed", data: {...} }
    if (res.code === 0) {
      return res.data
    } else {
      // 处理业务错误，抛出错误让调用方自行处理
      return Promise.reject(new Error(res.message || res.detail || '请求失败'))
    }
  },
  (error) => {
    let message = '操作失败，请稍后重试'
    
    if (error.response) {
      const { status, data } = error.response
      
      if (status === 401) {
        const authStore = useAuthStore()
        authStore.logout()
        router.push('/login')
        message = '登录已过期，请重新登录'
      } else if (status === 403) {
        message = '没有权限进行此操作'
      } else if (status === 404) {
        message = '请求的资源不存在'
      } else if (status >= 500) {
        message = '服务器异常，请稍后重试'
      } else {
        // 从响应中提取后端返回的错误信息
        message = data?.message || data?.detail || data?.error || `请求失败 (${status})`
      }
    } else if (error.request) {
      message = '网络连接失败，请检查网络'
    }
    
    ElMessage.error(message)
    return Promise.reject(new Error(message))
  }
)

export default request
