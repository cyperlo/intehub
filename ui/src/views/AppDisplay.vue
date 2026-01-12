<template>
  <div class="app-display">
    <el-card v-loading="loading || running">
      <template #header>
        <div class="header">
          <h2>{{ app?.name || '应用展示' }}</h2>
          <el-button type="primary" @click="handleRun" :loading="running">
            <el-icon><Refresh /></el-icon>
            重新运行
          </el-button>
        </div>
      </template>

      <div v-if="!result && !running" class="empty-state">
        <el-empty description="应用运行中..." />
      </div>

      <div v-else-if="result" class="result-container">
        <!-- 视频展示 -->
        <div v-if="mediaType === 'video'" class="media-container">
          <video 
            :src="mediaUrl" 
            controls 
            autoplay
            class="media-video"
            @error="handleMediaError"
          >
            您的浏览器不支持视频播放
          </video>
        </div>

        <!-- 图片展示 -->
        <div v-else-if="mediaType === 'image'" class="media-container">
          <img 
            :src="mediaUrl" 
            class="media-image"
            @error="handleMediaError"
          />
        </div>

        <!-- 链接展示 -->
        <div v-else-if="mediaType === 'link'" class="link-container">
          <el-result icon="success" title="应用执行成功">
            <template #extra>
              <el-button type="primary" @click="openLink">
                <el-icon><Link /></el-icon>
                打开链接
              </el-button>
            </template>
            <template #sub-title>
              <div class="link-info">
                <p>链接地址：</p>
                <el-input v-model="mediaUrl" readonly>
                  <template #append>
                    <el-button @click="copyLink">
                      <el-icon><CopyDocument /></el-icon>
                    </el-button>
                  </template>
                </el-input>
              </div>
            </template>
          </el-result>
        </div>

        <!-- 原始输出 -->
        <div v-else class="output-container">
          <el-alert 
            :type="result.status === 'success' ? 'success' : 'error'"
            :title="result.status === 'success' ? '执行成功' : '执行失败'"
            :closable="false"
            style="margin-bottom: 16px;"
          >
            <template #default>
              <p>耗时：{{ result.duration }}ms</p>
            </template>
          </el-alert>
          
          <div v-if="result.status === 'success'">
            <h4>输出结果：</h4>
            <pre class="output-content">{{ formatOutput(result.output) }}</pre>
          </div>
          
          <div v-else>
            <h4>错误信息：</h4>
            <pre class="error-content">{{ result.error }}</pre>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Refresh, Link, CopyDocument } from '@element-plus/icons-vue'
import { getApp, runApp, type App } from '../api/app'

const route = useRoute()
const appId = computed(() => Number(route.params.id))

const loading = ref(false)
const running = ref(false)
const app = ref<App>()
const result = ref<any>(null)
const mediaUrl = ref('')
const mediaType = ref<'video' | 'image' | 'link' | 'text'>('text')

const loadApp = async () => {
  loading.value = true
  try {
    app.value = await getApp(appId.value)
  } catch (error) {
    ElMessage.error('加载应用失败')
  } finally {
    loading.value = false
  }
}

const handleRun = async () => {
  running.value = true
  try {
    const res: any = await runApp(appId.value)
    result.value = res
    
    if (res.status === 'success') {
      parseOutput(res.output)
      ElMessage.success('应用执行成功')
    } else {
      ElMessage.error('应用执行失败')
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '执行失败')
  } finally {
    running.value = false
  }
}

const parseOutput = (output: any) => {
  try {
    let url = ''
    
    // 尝试解析 JSON
    if (typeof output === 'string') {
      try {
        const parsed = JSON.parse(output)
        // 查找 URL 字段
        const urlFields = ['url', 'video_url', 'image_url', 'link', 'video', 'image', 'file_url']
        for (const field of urlFields) {
          if (parsed[field]) {
            url = parsed[field]
            break
          }
        }
      } catch (e) {
        // 可能是纯 URL 字符串
        url = output.trim()
      }
    } else if (typeof output === 'object' && output !== null) {
      const urlFields = ['url', 'video_url', 'image_url', 'link', 'video', 'image', 'file_url']
      for (const field of urlFields) {
        if (output[field]) {
          url = output[field]
          break
        }
      }
    }
    
    if (!url || !(url.startsWith('http://') || url.startsWith('https://'))) {
      mediaType.value = 'text'
      return
    }
    
    mediaUrl.value = url
    
    // 判断媒体类型
    const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov', '.avi', '.mkv', '.m3u8']
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg', '.bmp']
    
    const urlLower = url.toLowerCase()
    
    if (videoExtensions.some(ext => urlLower.includes(ext)) || urlLower.includes('video')) {
      mediaType.value = 'video'
    } else if (imageExtensions.some(ext => urlLower.includes(ext)) || urlLower.includes('image')) {
      mediaType.value = 'image'
    } else {
      mediaType.value = 'link'
    }
  } catch (e) {
    mediaType.value = 'text'
  }
}

const formatOutput = (output: any) => {
  if (!output) return '无输出'
  
  try {
    if (typeof output === 'string') {
      const parsed = JSON.parse(output)
      return JSON.stringify(parsed, null, 2)
    } else if (typeof output === 'object') {
      return JSON.stringify(output, null, 2)
    }
  } catch (e) {
    // 格式化失败
  }
  
  return output
}

const handleMediaError = () => {
  ElMessage.error('媒体加载失败')
  mediaType.value = 'text'
}

const openLink = () => {
  window.open(mediaUrl.value, '_blank')
}

const copyLink = () => {
  navigator.clipboard.writeText(mediaUrl.value)
  ElMessage.success('链接已复制')
}

const formatTime = (time: string | undefined) => {
  if (!time) return '-'
  try {
    const date = new Date(time)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (e) {
    return time
  }
}

onMounted(async () => {
  await loadApp()
  // 页面加载后自动运行应用
  handleRun()
})
</script>

<style scoped>
.app-display {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 500;
}

.empty-state {
  padding: 60px 0;
}

.result-container {
  min-height: 400px;
}

.media-container {
  display: flex;
  justify-content: center;
  align-items: center;
  background: #000;
  border-radius: 8px;
  overflow: hidden;
  min-height: 500px;
}

.media-video {
  max-width: 100%;
  max-height: 80vh;
  width: 100%;
}

.media-image {
  max-width: 100%;
  max-height: 80vh;
  object-fit: contain;
}

.link-container {
  padding: 40px 0;
}

.link-info {
  max-width: 600px;
  margin: 20px auto 0;
  text-align: left;
}

.link-info p {
  margin-bottom: 8px;
  color: #606266;
  font-size: 14px;
}

.output-container h4 {
  margin: 16px 0 8px;
  color: #606266;
  font-size: 14px;
  font-weight: 500;
}

.output-content,
.error-content {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
  max-height: 600px;
  overflow-y: auto;
}

.error-content {
  background: #fef0f0;
  color: #f56c6c;
}

@media (max-width: 768px) {
  .app-display {
    padding: 12px;
  }
  
  .header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
  
  .media-container {
    min-height: 300px;
  }
}
</style>
