<template>
  <div class="appstore-container">
    <!-- 使用说明 -->
    <el-alert 
      title="应用商店使用说明" 
      type="success" 
      :closable="false"
      style="margin-bottom: 16px;"
    >
      <template #default>
        <div style="line-height: 1.8;">
          <p><strong>1. 浏览模板：</strong>查看可用的应用模板，按分类筛选</p>
          <p><strong>2. 安装应用：</strong>点击"安装"按钮，输入应用名称和配置参数</p>
          <p><strong>3. 使用应用：</strong>安装后在"应用管理"页面查看和运行</p>
          <p><strong>4. 发布模板：</strong>在"应用管理"中创建应用后，可点击"发布"分享给其他用户</p>
        </div>
      </template>
    </el-alert>
    
    <el-card class="header-card">
      <div class="header-content">
        <div>
          <h2>应用商店</h2>
          <p class="subtitle">从模板快速创建应用</p>
        </div>
        <el-button v-if="isAdmin" type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          创建模板
        </el-button>
      </div>
    </el-card>

    <!-- 分类筛选 -->
    <el-card class="filter-card">
      <el-radio-group v-model="selectedCategory" @change="loadTemplates">
        <el-radio-button label="">全部</el-radio-button>
        <el-radio-button label="工具">工具</el-radio-button>
        <el-radio-button label="数据处理">数据处理</el-radio-button>
        <el-radio-button label="通知">通知</el-radio-button>
        <el-radio-button label="示例">示例</el-radio-button>
      </el-radio-group>
    </el-card>

    <!-- 模板列表 -->
    <div v-loading="loading" class="templates-grid">
      <el-card
        v-for="template in templates"
        :key="template.id"
        class="template-card"
        shadow="hover"
      >
        <div class="template-header">
          <div class="template-icon">
            <el-icon v-if="!template.icon" :size="40"><Box /></el-icon>
            <img v-else :src="template.icon" alt="" />
          </div>
          <div class="template-info">
            <h3>{{ template.display_name }}</h3>
            <div class="template-meta">
              <el-tag size="small" type="info">{{ template.category }}</el-tag>
              <span class="version">v{{ template.version }}</span>
            </div>
          </div>
        </div>

        <p class="template-description">{{ template.description }}</p>

        <div class="template-tags">
          <el-tag
            v-for="tag in template.tags?.split(',')"
            :key="tag"
            size="small"
            effect="plain"
          >
            {{ tag }}
          </el-tag>
        </div>

        <div class="template-footer">
          <span class="author">by {{ template.author }}</span>
          <div class="actions">
            <el-button size="small" @click="viewTemplate(template)">
              查看
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="showInstallDialog(template)"
            >
              安装
            </el-button>
            <el-dropdown v-if="isAdmin" trigger="click">
              <el-button size="small" :icon="More" circle />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="editTemplate(template)">
                    编辑
                  </el-dropdown-item>
                  <el-dropdown-item @click="deleteTemplateConfirm(template)">
                    删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 创建/编辑模板对话框 -->
    <el-dialog
      v-model="templateDialogVisible"
      :title="isEdit ? '编辑模板' : '创建模板'"
      width="800px"
    >
      <el-form :model="templateForm" label-width="100px">
        <el-form-item label="模板名称" required>
          <el-input v-model="templateForm.name" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="显示名称" required>
          <el-input v-model="templateForm.display_name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="templateForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="图标URL">
          <el-input v-model="templateForm.icon" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="templateForm.category">
            <el-option label="工具" value="工具" />
            <el-option label="数据处理" value="数据处理" />
            <el-option label="通知" value="通知" />
            <el-option label="示例" value="示例" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本">
          <el-input v-model="templateForm.version" />
        </el-form-item>
        <el-form-item label="作者">
          <el-input v-model="templateForm.author" />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="templateForm.tags" placeholder="用逗号分隔" />
        </el-form-item>
        <el-form-item label="代码" required>
          <el-input
            v-model="templateForm.code"
            type="textarea"
            :rows="12"
            placeholder="package goapp&#10;&#10;func Run(input map[string]any) (map[string]any, error) {&#10;  // 你的代码&#10;  return map[string]any{}, nil&#10;}"
          />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="templateForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="templateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTemplate">保存</el-button>
      </template>
    </el-dialog>

    <!-- 查看模板对话框 -->
    <el-dialog v-model="viewDialogVisible" title="模板详情" width="800px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="模板名称">
          {{ currentTemplate?.name }}
        </el-descriptions-item>
        <el-descriptions-item label="显示名称">
          {{ currentTemplate?.display_name }}
        </el-descriptions-item>
        <el-descriptions-item label="分类">
          {{ currentTemplate?.category }}
        </el-descriptions-item>
        <el-descriptions-item label="版本">
          {{ currentTemplate?.version }}
        </el-descriptions-item>
        <el-descriptions-item label="作者">
          {{ currentTemplate?.author }}
        </el-descriptions-item>
        <el-descriptions-item label="标签">
          {{ currentTemplate?.tags }}
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">
          {{ currentTemplate?.description }}
        </el-descriptions-item>
      </el-descriptions>
      <div style="margin-top: 20px">
        <h4>代码</h4>
        <pre class="code-block">{{ currentTemplate?.code }}</pre>
      </div>
    </el-dialog>

    <!-- 安装应用对话框 -->
    <el-dialog v-model="installDialogVisible" title="安装应用" width="600px">
      <el-form :model="installForm" label-width="100px">
        <el-form-item label="应用名称" required>
          <el-input v-model="installForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="installForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-divider>配置项</el-divider>
        <el-form-item
          v-for="(_, key) in installForm.configs"
          :key="key"
          :label="key"
        >
          <el-input v-model="installForm.configs[key]" />
        </el-form-item>
        <el-button size="small" @click="addConfig">添加配置</el-button>
      </el-form>
      <template #footer>
        <el-button @click="installDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="installApp">安装</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Box, More } from '@element-plus/icons-vue'
import {
  getTemplates,
  createTemplate,
  updateTemplate,
  deleteTemplate,
  installFromTemplate,
  type AppTemplate
} from '../api/appstore'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const isAdmin = computed(() => authStore.user?.role === 'admin')

const loading = ref(false)
const templates = ref<AppTemplate[]>([])
const selectedCategory = ref('')

const templateDialogVisible = ref(false)
const viewDialogVisible = ref(false)
const installDialogVisible = ref(false)

const isEdit = ref(false)
const currentTemplate = ref<AppTemplate | null>(null)

const templateForm = ref<Partial<AppTemplate>>({
  name: '',
  display_name: '',
  description: '',
  icon: '',
  code: '',
  language: 'go',
  category: '工具',
  version: '1.0.0',
  author: '',
  tags: '',
  config_schema: '',
  enabled: true
})

const installForm = ref({
  template_id: 0,
  name: '',
  description: '',
  user_id: 0,
  configs: {} as Record<string, string>
})

const loadTemplates = async () => {
  loading.value = true
  try {
    const res = await getTemplates(selectedCategory.value)
    // 后端直接返回数组
    templates.value = res || []
    console.log('加载模板数量:', templates.value.length)
  } catch (error) {
    console.error('加载模板失败:', error)
    ElMessage.error('加载模板失败')
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  isEdit.value = false
  templateForm.value = {
    name: '',
    display_name: '',
    description: '',
    icon: '',
    code: '',
    language: 'go',
    category: '工具',
    version: '1.0.0',
    author: authStore.user?.nickname || '',
    tags: '',
    config_schema: '',
    enabled: true
  }
  templateDialogVisible.value = true
}

const editTemplate = (template: AppTemplate) => {
  isEdit.value = true
  currentTemplate.value = template
  templateForm.value = { ...template }
  templateDialogVisible.value = true
}

const saveTemplate = async () => {
  try {
    if (isEdit.value && currentTemplate.value?.id) {
      await updateTemplate(currentTemplate.value.id, templateForm.value)
      ElMessage.success('更新成功')
    } else {
      await createTemplate(templateForm.value)
      ElMessage.success('创建成功')
    }
    templateDialogVisible.value = false
    loadTemplates()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const deleteTemplateConfirm = (template: AppTemplate) => {
  ElMessageBox.confirm('确定要删除这个模板吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteTemplate(template.id!)
      ElMessage.success('删除成功')
      loadTemplates()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

const viewTemplate = (template: AppTemplate) => {
  currentTemplate.value = template
  viewDialogVisible.value = true
}

const showInstallDialog = (template: AppTemplate) => {
  currentTemplate.value = template
  // 兼容大小写 ID
  const userId = authStore.user?.id || (authStore.user as any)?.ID
  if (!userId) {
    ElMessage.error('请先登录')
    return
  }
  installForm.value = {
    template_id: template.id!,
    name: template.display_name,
    description: template.description,
    user_id: userId,
    configs: {}
  }
  installDialogVisible.value = true
}

const addConfig = () => {
  ElMessageBox.prompt('请输入配置项名称', '添加配置', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(({ value }) => {
    if (value) {
      installForm.value.configs[value] = ''
    }
  })
}

const installApp = async () => {
  try {
    await installFromTemplate(installForm.value)
    ElMessage.success('安装成功！已添加到应用管理')
    installDialogVisible.value = false
    
    // 提示用户去应用管理查看
    ElMessageBox.confirm('应用已安装成功，是否前往应用管理查看？', '提示', {
      confirmButtonText: '前往',
      cancelButtonText: '留在此页',
      type: 'success'
    }).then(() => {
      router.push('/apps')
    }).catch(() => {
      // 用户选择留在当前页
    })
  } catch (error) {
    ElMessage.error('安装失败')
  }
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.appstore-container {
  padding: 20px;
}

.header-card {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h2 {
  margin: 0 0 5px 0;
}

.subtitle {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.filter-card {
  margin-bottom: 20px;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.template-card {
  transition: transform 0.2s;
}

.template-card:hover {
  transform: translateY(-5px);
}

.template-header {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.template-icon {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-radius: 8px;
}

.template-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.template-info {
  flex: 1;
}

.template-info h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
}

.template-meta {
  display: flex;
  gap: 10px;
  align-items: center;
}

.version {
  font-size: 12px;
  color: #909399;
}

.template-description {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 15px;
  min-height: 44px;
}

.template-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 15px;
  min-height: 24px;
}

.template-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
}

.author {
  font-size: 12px;
  color: #909399;
}

.actions {
  display: flex;
  gap: 8px;
}

.code-block {
  background: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 13px;
  line-height: 1.5;
}
</style>
