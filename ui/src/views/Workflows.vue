<template>
  <div class="appflow-container">
    <el-card>
      <template #header>
        <div class="header">
          <span>应用流管理</span>
          <el-button type="primary" @click="showCreateDialog">创建应用流</el-button>
        </div>
      </template>

      <el-table :data="workflows" v-loading="loading">
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editWorkflow(row)">编辑</el-button>
            <el-button type="success" link @click="runWorkflow(row)">运行</el-button>
            <el-button type="info" link @click="viewLogs(row)">日志</el-button>
            <el-button type="danger" link @click="deleteWorkflow(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑对话框 - 可视化画布 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑应用流' : '创建应用流'"
      fullscreen
      :close-on-click-modal="false"
      class="flow-dialog"
      destroy-on-close
    >
      <div class="flow-editor">
        <div class="flow-header">
          <el-form :model="form" inline>
            <el-form-item label="名称">
              <el-input v-model="form.name" placeholder="请输入应用流名称" style="width: 250px" />
            </el-form-item>
            <el-form-item label="描述">
              <el-input v-model="form.description" placeholder="请输入描述" style="width: 400px" />
            </el-form-item>
            <el-form-item label="状态">
              <el-switch v-model="form.enabled" active-text="启用" inactive-text="禁用" />
            </el-form-item>
            <el-form-item style="margin-left: auto">
              <el-button @click="dialogVisible = false">取消</el-button>
              <el-button type="primary" @click="saveWorkflow">保存应用流</el-button>
            </el-form-item>
          </el-form>
        </div>
        
        <div class="flow-content">
          <!-- 左侧应用列表 -->
          <div class="app-sidebar">
            <div class="sidebar-header">
              <el-icon><Box /></el-icon>
              <span>可用应用</span>
            </div>
            <div class="sidebar-search">
              <el-input
                v-model="appSearch"
                placeholder="搜索应用"
                size="small"
                clearable
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </div>
            <div class="app-list">
              <div
                v-for="app in filteredApps"
                :key="app.id"
                class="app-item"
                draggable="true"
                @dragstart="onDragStart($event, app)"
              >
                <div class="app-item-icon">
                  <el-icon><Box /></el-icon>
                </div>
                <div class="app-item-content">
                  <div class="app-item-name">{{ app.name }}</div>
                  <div class="app-item-desc">{{ app.description || '暂无描述' }}</div>
                </div>
              </div>
              <el-empty v-if="filteredApps.length === 0" description="暂无应用" :image-size="80" />
            </div>
            <div class="sidebar-tip">
              <el-alert
                title="拖拽应用到右侧画布创建节点，连接节点定义执行顺序"
                type="info"
                :closable="false"
              />
            </div>
            <div class="sidebar-special">
              <div class="sidebar-header">
                <el-icon><Operation /></el-icon>
                <span>特殊节点</span>
              </div>
              <div
                class="special-node start-node"
                draggable="true"
                @dragstart="onDragStartSpecial($event, 'start')"
              >
                <el-icon><VideoPlay /></el-icon>
                <span>开始节点</span>
              </div>
              <div
                class="special-node end-node"
                draggable="true"
                @dragstart="onDragStartSpecial($event, 'end')"
              >
                <el-icon><CircleCheck /></el-icon>
                <span>结束节点</span>
              </div>
              <div
                class="special-node condition-node"
                draggable="true"
                @dragstart="onDragStartSpecial($event, 'condition')"
              >
                <el-icon><QuestionFilled /></el-icon>
                <span>判断节点</span>
              </div>
              <div
                class="special-node parallel-node"
                draggable="true"
                @dragstart="onDragStartSpecial($event, 'parallel')"
              >
                <el-icon><Share /></el-icon>
                <span>并行节点</span>
              </div>
            </div>
          </div>

          <!-- 右侧画布 -->
          <div class="canvas-area" @drop="onDrop" @dragover.prevent @contextmenu.prevent="onCanvasContextMenu">
            <div class="canvas-toolbar">
              <el-button-group>
                <el-tooltip content="撤销" placement="bottom">
                  <el-button size="small" @click="undo" :disabled="!canUndo">
                    <el-icon><RefreshLeft /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="重做" placement="bottom">
                  <el-button size="small" @click="redo" :disabled="!canRedo">
                    <el-icon><RefreshRight /></el-icon>
                  </el-button>
                </el-tooltip>
              </el-button-group>
              
              <el-button-group style="margin-left: 10px">
                <el-tooltip content="放大" placement="bottom">
                  <el-button size="small" @click="zoomIn">
                    <el-icon><ZoomIn /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="缩小" placement="bottom">
                  <el-button size="small" @click="zoomOut">
                    <el-icon><ZoomOut /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="适应画布" placement="bottom">
                  <el-button size="small" @click="fitView">
                    <el-icon><FullScreen /></el-icon>
                  </el-button>
                </el-tooltip>
              </el-button-group>
              
              <div style="margin-left: 10px; display: inline-flex; align-items: center; gap: 8px">
                <span style="font-size: 12px; color: #606266">连线类型：</span>
                <el-select v-model="edgeType" size="small" style="width: 120px">
                  <el-option
                    v-for="item in edgeTypes"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </div>
              
              <el-button-group style="margin-left: 10px">
                <el-tooltip content="删除选中连线" placement="bottom">
                  <el-button size="small" @click="deleteSelectedEdges" type="warning">
                    <el-icon><Connection /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="自动布局" placement="bottom">
                  <el-button size="small" @click="autoLayout">
                    <el-icon><Grid /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="清空画布" placement="bottom">
                  <el-button size="small" @click="clearCanvas" type="danger">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-tooltip>
              </el-button-group>
            </div>
            <VueFlow
              v-model:nodes="flowNodes"
              v-model:edges="flowEdges"
              @nodes-change="onNodesChange"
              @edges-change="onEdgesChange"
              @connect="onConnect"
              class="vue-flow"
              :default-zoom="0.8"
              :min-zoom="0.1"
              :max-zoom="2"
            >
              <Background pattern-color="#aaa" :gap="16" />
              <Controls />
              
              <template #node-custom="{ data }">
                <Handle type="target" :position="Position.Left" />
                <div class="custom-node" @contextmenu.prevent.stop="onNodeContextMenu($event, data)">
                  <div class="node-header">
                    <el-icon class="node-icon"><Box /></el-icon>
                    <span class="node-title">{{ data.label }}</span>
                    <el-icon class="delete-icon" @click="deleteNode(data.id)"><Close /></el-icon>
                  </div>
                  <div class="node-body">
                    <el-button size="small" type="primary" link @click="configNode(data)">
                      <el-icon><Setting /></el-icon>
                      配置参数
                    </el-button>
                  </div>
                </div>
                <Handle type="source" :position="Position.Right" />
              </template>

              <template #node-start="{ data }">
                <div class="special-node-canvas start" @contextmenu.prevent.stop="onNodeContextMenu($event, data)">
                  <el-icon><VideoPlay /></el-icon>
                  <span>开始</span>
                  <el-icon class="delete-icon" @click="deleteNode(data.id)"><Close /></el-icon>
                </div>
                <Handle type="source" :position="Position.Right" />
              </template>

              <template #node-end="{ data }">
                <Handle type="target" :position="Position.Left" />
                <div class="special-node-canvas end" @contextmenu.prevent.stop="onNodeContextMenu($event, data)">
                  <el-icon><CircleCheck /></el-icon>
                  <span>结束</span>
                  <el-icon class="delete-icon" @click="deleteNode(data.id)"><Close /></el-icon>
                </div>
              </template>

              <template #node-condition="{ data }">
                <Handle type="target" :position="Position.Left" />
                <div class="special-node-canvas condition" @contextmenu.prevent.stop="onNodeContextMenu($event, data)">
                  <el-icon><QuestionFilled /></el-icon>
                  <span>{{ data.label || '判断' }}</span>
                  <el-icon class="delete-icon" @click="deleteNode(data.id)"><Close /></el-icon>
                </div>
                <Handle type="source" :position="Position.Right" id="true" style="top: 30%" />
                <Handle type="source" :position="Position.Right" id="false" style="top: 70%" />
                <el-button size="small" type="primary" link @click="configCondition(data)" class="config-btn">
                  <el-icon><Setting /></el-icon>
                </el-button>
              </template>

              <template #node-parallel="{ data }">
                <Handle type="target" :position="Position.Left" />
                <div class="special-node-canvas parallel" @contextmenu.prevent.stop="onNodeContextMenu($event, data)">
                  <el-icon><Share /></el-icon>
                  <span>并行</span>
                  <el-icon class="delete-icon" @click="deleteNode(data.id)"><Close /></el-icon>
                </div>
                <Handle type="source" :position="Position.Right" />
              </template>
            </VueFlow>
            <div v-if="flowNodes.length === 0" class="canvas-empty">
              <el-empty description="从左侧拖拽应用到此处开始创建应用流" :image-size="120" />
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 节点配置对话框 -->
    <el-dialog v-model="nodeConfigVisible" title="配置节点参数" width="600px">
      <el-form label-width="100px">
        <el-form-item label="应用名称">
          <el-input :value="currentNodeData?.label" disabled />
        </el-form-item>
        <el-form-item label="配置参数">
          <el-input
            v-model="nodeConfigStr"
            type="textarea"
            :rows="8"
            placeholder='JSON格式，如: {"key": "value"}'
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="nodeConfigVisible = false">取消</el-button>
        <el-button type="primary" @click="saveNodeConfig">保存</el-button>
      </template>
    </el-dialog>

    <!-- 判断节点配置对话框 -->
    <el-dialog v-model="conditionConfigVisible" title="配置判断条件" width="600px">
      <el-form label-width="100px">
        <el-form-item label="条件表达式">
          <el-input
            v-model="conditionExpression"
            type="textarea"
            :rows="4"
            placeholder='例如: data.status == "success" 或 data.value > 100'
          />
          <div style="color: #909399; font-size: 12px; margin-top: 4px;">
            支持JavaScript表达式，data为上一节点的输出数据
          </div>
        </el-form-item>
        <el-form-item label="True分支">
          <el-input v-model="conditionTrueLabel" placeholder="满足条件时的标签" />
        </el-form-item>
        <el-form-item label="False分支">
          <el-input v-model="conditionFalseLabel" placeholder="不满足条件时的标签" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="conditionConfigVisible = false">取消</el-button>
        <el-button type="primary" @click="saveConditionConfig">保存</el-button>
      </template>
    </el-dialog>

    <!-- 运行对话框 -->
    <el-dialog v-model="runDialogVisible" title="运行工作流" width="600px">
      <el-form label-width="100px">
        <el-form-item label="输入参数">
          <el-input
            v-model="runInput"
            type="textarea"
            :rows="6"
            placeholder='JSON格式，如: {"key": "value"}'
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="runDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="executeWorkflow" :loading="running">运行</el-button>
      </template>
    </el-dialog>

    <!-- 运行结果对话框 -->
    <el-dialog v-model="resultDialogVisible" title="运行结果" width="800px">
      <div class="result-container">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="状态">
            <el-tag :type="runResult.status === 'success' ? 'success' : 'danger'">
              {{ runResult.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="耗时">{{ runResult.duration }}ms</el-descriptions-item>
        </el-descriptions>
        <div class="result-section">
          <h4>输出结果</h4>
          <pre class="result-output">{{ formatOutput(runResult.output) }}</pre>
        </div>
        <div v-if="runResult.error" class="result-section">
          <h4>错误信息</h4>
          <pre class="result-error">{{ runResult.error }}</pre>
        </div>
      </div>
    </el-dialog>

    <!-- 日志对话框 -->
    <el-dialog v-model="logsDialogVisible" title="执行日志" width="90%" destroy-on-close>
      <el-table :data="logs" v-loading="logsLoading" style="width: 100%" border>
        <el-table-column prop="name" label="应用流" min-width="200" />
        <el-table-column label="状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时(ms)" width="120" align="center" />
        <el-table-column prop="started_at" label="开始时间" width="200">
          <template #default="{ row }">
            {{ formatTime(row.started_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="viewLogDetail(row)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!logsLoading && logs.length === 0" description="暂无执行日志" />
    </el-dialog>

    <!-- 日志详情对话框 -->
    <el-dialog v-model="logDetailVisible" title="执行日志详情" width="90%" class="log-detail-dialog">
      <div v-if="currentLog" class="log-detail">
        <el-card class="summary-card" shadow="never">
          <el-descriptions :column="4" border>
            <el-descriptions-item label="应用流名称">
              <el-tag type="info">{{ currentLog.name }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="执行状态">
              <el-tag :type="currentLog.status === 'success' ? 'success' : 'danger'" effect="dark">
                {{ currentLog.status === 'success' ? '✓ 成功' : '✗ 失败' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="总耗时">
              <el-tag type="warning">{{ currentLog.duration }}ms</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="开始时间">
              {{ formatTime(currentLog.started_at) }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card class="input-card" shadow="never">
          <template #header>
            <div class="card-header-title">
              <el-icon><Upload /></el-icon>
              <span>输入参数</span>
            </div>
          </template>
          <pre class="log-content">{{ formatOutput(currentLog.input) }}</pre>
        </el-card>

        <el-card class="nodes-card" shadow="never">
          <template #header>
            <div class="card-header-title">
              <el-icon><Connection /></el-icon>
              <span>节点执行流程（共 {{ parseNodeLogs(currentLog.node_logs).length }} 个节点）</span>
            </div>
          </template>
          
          <el-timeline>
            <el-timeline-item
              v-for="(nodeLog, index) in parseNodeLogs(currentLog.node_logs)"
              :key="index"
              :type="nodeLog.status === 'success' ? 'success' : 'danger'"
              :icon="nodeLog.status === 'success' ? 'CircleCheck' : 'CircleClose'"
              :size="'large'"
            >
              <el-card class="node-log-card" :class="nodeLog.status">
                <template #header>
                  <div class="node-log-header">
                    <div class="node-info">
                      <el-tag :type="getNodeTypeTag(nodeLog)" size="small">
                        {{ getNodeTypeName(nodeLog) }}
                      </el-tag>
                      <span class="node-name">{{ getNodeDisplayName(nodeLog) }}</span>
                    </div>
                    <div class="node-status">
                      <el-tag :type="nodeLog.status === 'success' ? 'success' : 'danger'" size="small">
                        {{ nodeLog.status === 'success' ? '成功' : '失败' }}
                      </el-tag>
                      <span class="node-duration">{{ nodeLog.duration }}ms</span>
                    </div>
                  </div>
                </template>
                
                <el-collapse>
                  <el-collapse-item name="1">
                    <template #title>
                      <div class="collapse-title">
                        <el-icon><View /></el-icon>
                        <span>查看输入输出详情</span>
                      </div>
                    </template>
                    <div class="node-detail">
                      <el-row :gutter="16">
                        <el-col :xs="24" :sm="12">
                          <div class="detail-section input-section">
                            <div class="section-title">
                              <el-icon><Download /></el-icon>
                              <span>输入数据</span>
                              <el-tag size="small" type="info">Input</el-tag>
                            </div>
                            <div class="detail-content-wrapper">
                              <pre class="detail-content">{{ formatJSON(nodeLog.input) }}</pre>
                            </div>
                          </div>
                        </el-col>
                        <el-col :xs="24" :sm="12">
                          <div class="detail-section output-section">
                            <div class="section-title">
                              <el-icon><Upload /></el-icon>
                              <span>输出数据</span>
                              <el-tag size="small" type="success">Output</el-tag>
                            </div>
                            <div class="detail-content-wrapper">
                              <pre class="detail-content">{{ formatJSON(nodeLog.output) }}</pre>
                            </div>
                          </div>
                        </el-col>
                      </el-row>
                      
                      <div v-if="nodeLog.error" class="detail-section error-section">
                        <el-divider />
                        <div class="section-title">
                          <el-icon><Warning /></el-icon>
                          <span>错误信息</span>
                          <el-tag size="small" type="danger">Error</el-tag>
                        </div>
                        <div class="error-content-wrapper">
                          <pre class="error-content">{{ nodeLog.error }}</pre>
                        </div>
                      </div>
                    </div>
                  </el-collapse-item>
                </el-collapse>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <el-card class="output-card" shadow="never">
          <template #header>
            <div class="card-header-title">
              <el-icon><Download /></el-icon>
              <span>最终输出</span>
            </div>
          </template>
          <pre class="log-content">{{ formatOutput(currentLog.output) }}</pre>
        </el-card>

        <el-card v-if="currentLog.error" class="error-card" shadow="never">
          <template #header>
            <div class="card-header-title">
              <el-icon><Warning /></el-icon>
              <span>错误信息</span>
            </div>
          </template>
          <pre class="log-error">{{ currentLog.error }}</pre>
        </el-card>
      </div>
    </el-dialog>

    <!-- 右键菜单 -->
    <teleport to="body">
      <div
        v-show="contextMenuVisible"
        class="context-menu"
        :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
        @click="contextMenuVisible = false"
      >
        <template v-if="contextMenuType === 'node'">
          <div class="menu-item" @click="configNode(contextMenuNode)">
            <el-icon><Setting /></el-icon>
            <span>配置参数</span>
          </div>
          <div class="menu-item" @click="copyNode">
            <el-icon><DocumentCopy /></el-icon>
            <span>复制节点</span>
          </div>
          <div class="menu-divider"></div>
          <div class="menu-item danger" @click="deleteNode(contextMenuNode?.id)">
            <el-icon><Delete /></el-icon>
            <span>删除节点</span>
          </div>
        </template>
        <template v-else>
          <div class="menu-item" @click="pasteNode">
            <el-icon><DocumentCopy /></el-icon>
            <span>粘贴节点</span>
          </div>
          <div class="menu-item" @click="autoLayout">
            <el-icon><Grid /></el-icon>
            <span>自动布局</span>
          </div>
          <div class="menu-item" @click="fitView">
            <el-icon><FullScreen /></el-icon>
            <span>适应画布</span>
          </div>
          <div class="menu-divider"></div>
          <div class="menu-item" @click="clearCanvas">
            <el-icon><Delete /></el-icon>
            <span>清空画布</span>
          </div>
        </template>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { VueFlow, Handle, Position, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { Box, Close, Search, Delete, Grid, Setting, Operation, VideoPlay, CircleCheck, RefreshLeft, RefreshRight, ZoomIn, ZoomOut, FullScreen, Connection, DocumentCopy, Upload, Download, Warning, View, QuestionFilled, Share } from '@element-plus/icons-vue'
import {
  getWorkflows,
  createWorkflow,
  updateWorkflow,
  deleteWorkflow as deleteWorkflowAPI,
  runWorkflow as runWorkflowAPI,
  getWorkflowLogs,
  type Workflow,
  type WorkflowNode,
  type WorkflowLog
} from '../api/workflow'
import { getApps, type App } from '../api/app'
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'

const loading = ref(false)
const workflows = ref<Workflow[]>([])
const apps = ref<App[]>([])
const appSearch = ref('')

const dialogVisible = ref(false)
const isEdit = ref(false)
const form = ref<Workflow>({
  name: '',
  description: '',
  nodes: '[]',
  enabled: true
})

// Vue Flow 相关
const flowNodes = ref<any[]>([])
const flowEdges = ref<any[]>([])
let nodeIdCounter = 0

// 连线类型
const edgeType = ref('default')
const edgeTypes = [
  { label: '直线', value: 'default' },
  { label: '平滑曲线', value: 'smoothstep' },
  { label: '阶梯线', value: 'step' },
  { label: '贝塞尔曲线', value: 'bezier' }
]

// 历史记录
const history = ref<Array<{ nodes: any[], edges: any[] }>>([])
const historyIndex = ref(-1)
const canUndo = computed(() => historyIndex.value > 0)
const canRedo = computed(() => historyIndex.value < history.value.length - 1)

// Vue Flow 实例
const { zoomIn: vueFlowZoomIn, zoomOut: vueFlowZoomOut, fitView: vueFlowFitView } = useVueFlow()

// 节点配置
const nodeConfigVisible = ref(false)
const currentNodeData = ref<any>(null)
const nodeConfigStr = ref('{}')

// 判断节点配置
const conditionConfigVisible = ref(false)
const conditionExpression = ref('')
const conditionTrueLabel = ref('True')
const conditionFalseLabel = ref('False')
const currentConditionNode = ref<any>(null)

// 右键菜单
const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const contextMenuType = ref<'node' | 'canvas'>('canvas')
const contextMenuNode = ref<any>(null)

const runDialogVisible = ref(false)
const runInput = ref('{}')
const currentWorkflowId = ref(0)
const running = ref(false)

const resultDialogVisible = ref(false)
const runResult = ref<any>({})

const logsDialogVisible = ref(false)
const logsLoading = ref(false)
const logs = ref<WorkflowLog[]>([])

const logDetailVisible = ref(false)
const currentLog = ref<WorkflowLog | null>(null)

const filteredApps = computed(() => {
  if (!appSearch.value) return apps.value
  return apps.value.filter(app => 
    app.name.toLowerCase().includes(appSearch.value.toLowerCase())
  )
})

const loadWorkflows = async () => {
  loading.value = true
  try {
    workflows.value = await getWorkflows()
  } catch (error: any) {
    ElMessage.error('加载应用流失败')
  } finally {
    loading.value = false
  }
}

const loadApps = async () => {
  try {
    apps.value = await getApps()
  } catch (error) {
    ElMessage.error('加载应用列表失败')
  }
}

const showCreateDialog = () => {
  isEdit.value = false
  form.value = {
    name: '',
    description: '',
    nodes: '[]',
    enabled: true
  }
  flowNodes.value = []
  flowEdges.value = []
  nodeIdCounter = 0
  history.value = []
  historyIndex.value = -1
  dialogVisible.value = true
  saveHistory()
}

const editWorkflow = (workflow: Workflow) => {
  isEdit.value = true
  form.value = { ...workflow }
  
  try {
    const parsedNodes = JSON.parse(workflow.nodes)
    // 转换为 Vue Flow 格式
    flowNodes.value = parsedNodes.map((node: WorkflowNode, index: number) => {
      // 根据节点类型设置不同的属性
      if (node.type === 'start') {
        return {
          id: node.id,
          type: 'start',
          position: node.position || { x: 100, y: 100 + index * 150 },
          data: {
            id: node.id,
            label: '开始',
            nodeType: 'start'
          }
        }
      } else if (node.type === 'end') {
        return {
          id: node.id,
          type: 'end',
          position: node.position || { x: 100, y: 100 + index * 150 },
          data: {
            id: node.id,
            label: '结束',
            nodeType: 'end'
          }
        }
      } else {
        return {
          id: node.id,
          type: 'custom',
          position: node.position || { x: 100, y: 100 + index * 150 },
          data: {
            id: node.id,
            label: apps.value.find(a => a.id === node.app_id)?.name || '未知应用',
            app_id: node.app_id,
            config: node.config || {}
          }
        }
      }
    })
    
    // 转换边
    flowEdges.value = []
    parsedNodes.forEach((node: WorkflowNode) => {
      if (node.next && node.next.length > 0) {
        node.next.forEach(nextId => {
          flowEdges.value.push({
            id: `e${node.id}-${nextId}`,
            source: node.id,
            target: nextId,
            type: 'default'
          })
        })
      }
    })
    
    nodeIdCounter = parsedNodes.length
  } catch (e) {
    flowNodes.value = []
    flowEdges.value = []
  }
  
  dialogVisible.value = true
}

// 拖拽相关
const onDragStart = (event: DragEvent, app: App) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('application/json', JSON.stringify({ type: 'app', data: app }))
  }
}

const onDragStartSpecial = (event: DragEvent, nodeType: string) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('application/json', JSON.stringify({ type: 'special', nodeType }))
  }
}

const onDrop = (event: DragEvent) => {
  event.preventDefault()
  const dataStr = event.dataTransfer?.getData('application/json')
  if (!dataStr) return
  
  const dragData = JSON.parse(dataStr)
  const canvasArea = event.currentTarget as HTMLElement
  const rect = canvasArea.getBoundingClientRect()
  
  const nodeId = `node_${++nodeIdCounter}`
  const position = {
    x: event.clientX - rect.left - 100,
    y: event.clientY - rect.top - 50
  }
  
  if (dragData.type === 'app') {
    const app = dragData.data as App
    const newNode = {
      id: nodeId,
      type: 'custom',
      position,
      data: {
        id: nodeId,
        label: app.name,
        app_id: app.id,
        config: {}
      }
    }
    flowNodes.value.push(newNode)
  } else if (dragData.type === 'special') {
    let label = '未知'
    if (dragData.nodeType === 'start') label = '开始'
    else if (dragData.nodeType === 'end') label = '结束'
    else if (dragData.nodeType === 'condition') label = '判断'
    else if (dragData.nodeType === 'parallel') label = '并行'
    
    const newNode = {
      id: nodeId,
      type: dragData.nodeType,
      position,
      data: {
        id: nodeId,
        label,
        nodeType: dragData.nodeType,
        config: {}
      }
    }
    flowNodes.value.push(newNode)
  }
  
  saveHistory()
}

const onNodesChange = () => {
  // 处理节点变化
}

const onEdgesChange = () => {
  // 处理边变化
}

const onConnect = (params: any) => {
  flowEdges.value.push({
    id: `e${params.source}-${params.target}`,
    source: params.source,
    target: params.target,
    type: edgeType.value
  })
  saveHistory()
}

const deleteNode = (nodeId: string) => {
  flowNodes.value = flowNodes.value.filter(n => n.id !== nodeId)
  flowEdges.value = flowEdges.value.filter(e => e.source !== nodeId && e.target !== nodeId)
  saveHistory()
}

// 删除选中的连线
const deleteSelectedEdges = () => {
  const selectedEdges = flowEdges.value.filter((e: any) => e.selected)
  if (selectedEdges.length === 0) {
    ElMessage.warning('请先选中要删除的连线')
    return
  }
  
  ElMessageBox.confirm(`确定要删除选中的 ${selectedEdges.length} 条连线吗？`, '提示', {
    type: 'warning'
  }).then(() => {
    flowEdges.value = flowEdges.value.filter((e: any) => !e.selected)
    saveHistory()
    ElMessage.success('已删除')
  }).catch(() => {})
}

const configNode = (data: any) => {
  currentNodeData.value = data
  nodeConfigStr.value = JSON.stringify(data.config || {}, null, 2)
  nodeConfigVisible.value = true
}

const saveNodeConfig = () => {
  try {
    const config = JSON.parse(nodeConfigStr.value)
    const node = flowNodes.value.find(n => n.data.id === currentNodeData.value.id)
    if (node) {
      node.data.config = config
    }
    nodeConfigVisible.value = false
    ElMessage.success('配置已保存')
  } catch (e) {
    ElMessage.warning('配置格式错误')
  }
}

const configCondition = (data: any) => {
  currentConditionNode.value = data
  conditionExpression.value = data.config?.expression || ''
  conditionTrueLabel.value = data.config?.trueLabel || 'True'
  conditionFalseLabel.value = data.config?.falseLabel || 'False'
  conditionConfigVisible.value = true
}

const saveConditionConfig = () => {
  const node = flowNodes.value.find(n => n.data.id === currentConditionNode.value.id)
  if (node) {
    node.data.config = {
      expression: conditionExpression.value,
      trueLabel: conditionTrueLabel.value,
      falseLabel: conditionFalseLabel.value
    }
    node.data.label = `判断: ${conditionExpression.value.substring(0, 20)}${conditionExpression.value.length > 20 ? '...' : ''}`
  }
  conditionConfigVisible.value = false
  ElMessage.success('判断条件已保存')
}

const clearCanvas = () => {
  ElMessageBox.confirm('确定要清空画布吗？', '提示', {
    type: 'warning'
  }).then(() => {
    flowNodes.value = []
    flowEdges.value = []
    saveHistory()
    ElMessage.success('已清空')
  }).catch(() => {})
}

const autoLayout = () => {
  // 简单的自动布局：垂直排列
  flowNodes.value.forEach((node, index) => {
    node.position = {
      x: 100,
      y: 100 + index * 150
    }
  })
  ElMessage.success('布局已调整')
}

// 保存历史记录
const saveHistory = () => {
  // 删除当前位置之后的历史
  history.value = history.value.slice(0, historyIndex.value + 1)
  
  // 添加新的历史记录
  history.value.push({
    nodes: JSON.parse(JSON.stringify(flowNodes.value)),
    edges: JSON.parse(JSON.stringify(flowEdges.value))
  })
  
  historyIndex.value = history.value.length - 1
  
  // 限制历史记录数量
  if (history.value.length > 50) {
    history.value.shift()
    historyIndex.value--
  }
}

// 撤销
const undo = () => {
  if (canUndo.value) {
    historyIndex.value--
    const state = history.value[historyIndex.value]
    flowNodes.value = JSON.parse(JSON.stringify(state.nodes))
    flowEdges.value = JSON.parse(JSON.stringify(state.edges))
  }
}

// 重做
const redo = () => {
  if (canRedo.value) {
    historyIndex.value++
    const state = history.value[historyIndex.value]
    flowNodes.value = JSON.parse(JSON.stringify(state.nodes))
    flowEdges.value = JSON.parse(JSON.stringify(state.edges))
  }
}

// 放大
const zoomIn = () => {
  vueFlowZoomIn()
}

// 缩小
const zoomOut = () => {
  vueFlowZoomOut()
}

// 适应画布
const fitView = () => {
  vueFlowFitView()
}

// 右键菜单
const onNodeContextMenu = (event: MouseEvent, data: any) => {
  contextMenuType.value = 'node'
  contextMenuNode.value = data
  contextMenuPosition.value = { x: event.clientX, y: event.clientY }
  contextMenuVisible.value = true
}

const onCanvasContextMenu = (event: MouseEvent) => {
  contextMenuType.value = 'canvas'
  contextMenuPosition.value = { x: event.clientX, y: event.clientY }
  contextMenuVisible.value = true
}

// 复制节点
let copiedNode: any = null
const copyNode = () => {
  if (contextMenuNode.value) {
    const node = flowNodes.value.find(n => n.data.id === contextMenuNode.value.id)
    if (node) {
      copiedNode = JSON.parse(JSON.stringify(node))
      ElMessage.success('已复制节点')
    }
  }
}

// 粘贴节点
const pasteNode = () => {
  if (!copiedNode) {
    ElMessage.warning('没有可粘贴的节点')
    return
  }
  
  const nodeId = `node_${++nodeIdCounter}`
  const newNode = {
    ...copiedNode,
    id: nodeId,
    position: {
      x: copiedNode.position.x + 50,
      y: copiedNode.position.y + 50
    },
    data: {
      ...copiedNode.data,
      id: nodeId
    }
  }
  
  flowNodes.value.push(newNode)
  saveHistory()
  ElMessage.success('已粘贴节点')
}

const saveWorkflow = async () => {
  if (!form.value.name) {
    ElMessage.warning('请输入应用流名称')
    return
  }

  if (flowNodes.value.length === 0) {
    ElMessage.warning('请至少添加一个节点')
    return
  }

  // 检查是否有开始节点
  const hasStart = flowNodes.value.some(n => n.type === 'start')
  if (!hasStart) {
    ElMessage.warning('请添加开始节点')
    return
  }

  // 构建节点数据
  const nodes = flowNodes.value.map(node => {
    // 找到该节点的所有出边
    const outEdges = flowEdges.value.filter(e => e.source === node.id)
    const next = outEdges.map(e => e.target)
    
    return {
      id: node.id,
      type: node.type === 'start' || node.type === 'end' ? node.type : 'app',
      app_id: node.data.app_id || 0,
      config: node.data.config || {},
      position: node.position,
      next
    }
  })

  form.value.nodes = JSON.stringify(nodes)

  try {
    if (isEdit.value && form.value.id) {
      await updateWorkflow(form.value.id, form.value)
      ElMessage.success('更新成功')
    } else {
      await createWorkflow(form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadWorkflows()
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  }
}

const runWorkflow = (workflow: Workflow) => {
  currentWorkflowId.value = workflow.id!
  runInput.value = '{}'
  runDialogVisible.value = true
}

const executeWorkflow = async () => {
  let input = {}
  try {
    input = JSON.parse(runInput.value)
  } catch (e) {
    ElMessage.warning('输入参数格式错误')
    return
  }

  running.value = true
  try {
    const result = await runWorkflowAPI(currentWorkflowId.value, input)
    runResult.value = result
    runDialogVisible.value = false
    resultDialogVisible.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '运行失败')
  } finally {
    running.value = false
  }
}

const deleteWorkflow = async (workflow: Workflow) => {
  try {
    await ElMessageBox.confirm('确定要删除该工作流吗？', '提示', {
      type: 'warning'
    })
    await deleteWorkflowAPI(workflow.id!)
    ElMessage.success('删除成功')
    loadWorkflows()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const viewLogs = async (workflow: Workflow) => {
  logsDialogVisible.value = true
  logsLoading.value = true
  try {
    const result = await getWorkflowLogs({ workflow_id: workflow.id, page: 1, page_size: 50 })
    logs.value = result.list || []
  } catch (error) {
    ElMessage.error('加载日志失败')
  } finally {
    logsLoading.value = false
  }
}

const viewLogDetail = (log: WorkflowLog) => {
  currentLog.value = log
  logDetailVisible.value = true
}

const parseNodeLogs = (nodeLogsStr: string) => {
  try {
    return JSON.parse(nodeLogsStr || '[]')
  } catch (e) {
    return []
  }
}

const formatTime = (time: string) => {
  if (!time) return '-'
  try {
    const date = new Date(time)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hour = String(date.getHours()).padStart(2, '0')
    const minute = String(date.getMinutes()).padStart(2, '0')
    const second = String(date.getSeconds()).padStart(2, '0')
    return `${year}/${month}/${day} ${hour}:${minute}:${second}`
  } catch (e) {
    return time
  }
}

const formatOutput = (output: string) => {
  if (!output) return '-'
  try {
    const parsed = JSON.parse(output)
    return JSON.stringify(parsed, null, 2)
  } catch (e) {
    return output
  }
}

const formatJSON = (data: any) => {
  if (!data) return '{}'
  if (typeof data === 'string') {
    try {
      const parsed = JSON.parse(data)
      return JSON.stringify(parsed, null, 2)
    } catch (e) {
      return data
    }
  }
  return JSON.stringify(data, null, 2)
}

const getNodeTypeName = (nodeLog: any) => {
  if (nodeLog.app_name === '开始') return '开始节点'
  if (nodeLog.app_name === '结束') return '结束节点'
  return '应用节点'
}

const getNodeTypeTag = (nodeLog: any) => {
  if (nodeLog.app_name === '开始') return 'success'
  if (nodeLog.app_name === '结束') return 'info'
  return 'primary'
}

const getNodeDisplayName = (nodeLog: any) => {
  // 如果有 app_name 就显示 app_name
  if (nodeLog.app_name) {
    return nodeLog.app_name
  }
  
  // 尝试从当前工作流的节点配置中查找
  if (currentLog.value) {
    try {
      const workflow = workflows.value.find(w => w.id === currentLog.value?.workflow_id)
      if (workflow) {
        const nodes = JSON.parse(workflow.nodes)
        const node = nodes.find((n: any) => n.id === nodeLog.node_id)
        if (node) {
          // 如果是特殊节点
          if (node.type === 'start') return '开始'
          if (node.type === 'end') return '结束'
          // 如果是应用节点，查找应用名称
          if (node.app_id) {
            const app = apps.value.find(a => a.id === node.app_id)
            if (app) return app.name
          }
        }
      }
    } catch (e) {
      // 解析失败，继续使用默认逻辑
    }
  }
  
  // 否则根据 node_id 判断类型
  if (nodeLog.node_id?.includes('start')) {
    return '开始'
  }
  if (nodeLog.node_id?.includes('end')) {
    return '结束'
  }
  // 最后才显示 node_id
  return nodeLog.node_id || '未知节点'
}

onMounted(() => {
  loadWorkflows()
  loadApps()
  
  document.addEventListener('click', () => {
    contextMenuVisible.value = false
  })
})
</script>

<style>
/* 全局样式 - 不使用 scoped */
.flow-dialog.el-dialog.is-fullscreen {
  display: flex;
  flex-direction: column;
}

.flow-dialog.el-dialog.is-fullscreen .el-dialog__header {
  flex-shrink: 0;
}

.flow-dialog.el-dialog.is-fullscreen .el-dialog__body {
  flex: 1;
  overflow: hidden;
  padding: 0;
  display: flex;
  flex-direction: column;
}
</style>

<style scoped>
.appflow-container {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.flow-dialog :deep(.el-dialog__header) {
  border-bottom: 1px solid #ebeef5;
  padding: 15px 20px;
  margin: 0;
}

.flow-dialog :deep(.el-dialog__footer) {
  display: none;
}

.flow-editor {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.flow-header {
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
  background: #f5f7fa;
  flex-shrink: 0;
}

.flow-header :deep(.el-form) {
  display: flex;
  align-items: center;
  margin: 0;
}

.flow-header :deep(.el-form-item) {
  margin-bottom: 0;
}

.flow-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-height: 0;
}

.app-sidebar {
  width: 300px;
  border-right: 1px solid #ebeef5;
  background: #fff;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 15px 20px;
  font-weight: 600;
  font-size: 16px;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #303133;
}

.sidebar-search {
  padding: 15px 20px 10px;
}

.app-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px 15px;
}

.app-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  margin-bottom: 10px;
  background: #f5f7fa;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: move;
  transition: all 0.3s;
}

.app-item:hover {
  background: #ecf5ff;
  border-color: #409eff;
  transform: translateX(5px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.app-item-icon {
  font-size: 24px;
  color: #409eff;
  flex-shrink: 0;
}

.app-item-content {
  flex: 1;
  min-width: 0;
}

.app-item-name {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.app-item-desc {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sidebar-tip {
  padding: 15px;
  border-top: 1px solid #ebeef5;
}

.sidebar-special {
  border-top: 1px solid #ebeef5;
  padding: 15px;
}

.special-node {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: move;
  transition: all 0.3s;
  font-weight: 500;
}

.special-node.start-node {
  background: #f0f9ff;
  border: 2px solid #67c23a;
  color: #67c23a;
}

.special-node.start-node:hover {
  background: #e1f3d8;
  transform: translateX(5px);
  box-shadow: 0 2px 8px rgba(103, 194, 58, 0.3);
}

.special-node.end-node {
  background: #fef0f0;
  border: 2px solid #f56c6c;
  color: #f56c6c;
}

.special-node.end-node:hover {
  background: #fde2e2;
  transform: translateX(5px);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.3);
}

.special-node.condition-node {
  background: #fff7e6;
  border: 2px solid #e6a23c;
  color: #e6a23c;
}

.special-node.condition-node:hover {
  background: #faecd8;
  transform: translateX(5px);
  box-shadow: 0 2px 8px rgba(230, 162, 60, 0.3);
}

.special-node.parallel-node {
  background: #f0f9ff;
  border: 2px solid #409eff;
  color: #409eff;
}

.special-node.parallel-node:hover {
  background: #ecf5ff;
  transform: translateX(5px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.canvas-area {
  flex: 1;
  position: relative;
  background: #fafafa;
  min-width: 0;
  min-height: 0;
}

.canvas-toolbar {
  position: absolute;
  top: 15px;
  left: 15px;
  z-index: 10;
}

.vue-flow {
  width: 100%;
  height: 100%;
}

/* 优化连线样式 */
.vue-flow :deep(.vue-flow__edge-path) {
  stroke: #409eff;
  stroke-width: 2;
}

.vue-flow :deep(.vue-flow__edge.selected .vue-flow__edge-path) {
  stroke: #67c23a;
  stroke-width: 3;
}

.vue-flow :deep(.vue-flow__edge:hover .vue-flow__edge-path) {
  stroke: #67c23a;
  stroke-width: 3;
}

.vue-flow :deep(.vue-flow__edge-textwrapper) {
  pointer-events: all;
}

.vue-flow :deep(.vue-flow__edge-text) {
  fill: #409eff;
  font-size: 12px;
}

/* 连接点样式 */
.vue-flow :deep(.vue-flow__handle) {
  width: 10px;
  height: 10px;
  background: #409eff;
  border: 2px solid white;
}

.vue-flow :deep(.vue-flow__handle:hover) {
  background: #67c23a;
  transform: scale(1.3);
}

.canvas-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  pointer-events: none;
}

.custom-node {
  background: white;
  border: 2px solid #409eff;
  border-radius: 10px;
  padding: 12px 15px;
  min-width: 200px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s;
}

.custom-node:hover {
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.3);
  transform: translateY(-2px);
}

.custom-node .node-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
  color: #303133;
}

.custom-node .node-icon {
  font-size: 18px;
  color: #409eff;
}

.custom-node .node-title {
  flex: 1;
}

.custom-node .delete-icon {
  cursor: pointer;
  color: #909399;
  font-size: 16px;
  transition: all 0.3s;
}

.custom-node .delete-icon:hover {
  color: #f56c6c;
  transform: scale(1.2);
}

.custom-node .node-body {
  text-align: center;
}

.special-node-canvas {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 50px;
  font-weight: 600;
  font-size: 14px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s;
  position: relative;
}

.special-node-canvas:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
  transform: translateY(-2px);
}

.special-node-canvas.start {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  color: white;
  border: 2px solid #67c23a;
}

.special-node-canvas.end {
  background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%);
  color: white;
  border: 2px solid #f56c6c;
}

.special-node-canvas.condition {
  background: linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%);
  color: white;
  border: 2px solid #e6a23c;
  position: relative;
}

.special-node-canvas.parallel {
  background: linear-gradient(135deg, #409eff 0%, #79bbff 100%);
  color: white;
  border: 2px solid #409eff;
}

.special-node-canvas .config-btn {
  position: absolute;
  bottom: -30px;
  left: 50%;
  transform: translateX(-50%);
  color: white;
}

.special-node-canvas .delete-icon {
  position: absolute;
  top: -8px;
  right: -8px;
  background: white;
  border-radius: 50%;
  padding: 2px;
  cursor: pointer;
  color: #909399;
  font-size: 16px;
  transition: all 0.3s;
}

.special-node-canvas .delete-icon:hover {
  color: #f56c6c;
  transform: scale(1.2);
}

.result-container,
.log-detail {
  padding: 12px 0;
}

.result-section,
.log-section {
  margin-top: 20px;
}

.result-section h4,
.log-section h4 {
  margin-bottom: 8px;
  color: #303133;
}

.result-output,
.log-content {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
  max-height: 400px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
}

.result-error,
.log-error {
  background: #fef0f0;
  color: #f56c6c;
  padding: 12px;
  border-radius: 4px;
  max-height: 200px;
  overflow: auto;
}

.node-log {
  margin-bottom: 12px;
}

.node-log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-weight: 500;
}

.node-log-content {
  font-size: 14px;
}

.node-log-content p {
  margin: 8px 0;
}

.node-log-content pre {
  background: #f5f7fa;
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
  max-height: 200px;
  overflow: auto;
}

/* 右键菜单 */
.context-menu {
  position: fixed;
  background: white;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 4px 0;
  z-index: 9999;
  min-width: 150px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
  color: #606266;
}

.menu-item:hover {
  background: #f5f7fa;
  color: #409eff;
}

.menu-item.danger {
  color: #f56c6c;
}

.menu-item.danger:hover {
  background: #fef0f0;
}

.menu-divider {
  height: 1px;
  background: #ebeef5;
  margin: 4px 0;
}

/* 日志详情样式优化 */
.log-detail-dialog :deep(.el-dialog__body) {
  padding: 20px;
  background: #f5f7fa;
}

.log-detail {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

.summary-card,
.input-card,
.nodes-card,
.output-card,
.error-card {
  border-radius: 8px;
  width: 100%;
}

.nodes-card :deep(.el-timeline) {
  padding-left: 0;
  width: 100%;
}

.nodes-card :deep(.el-timeline-item) {
  width: 100%;
}

.nodes-card :deep(.el-timeline-item__wrapper) {
  padding-left: 28px;
  width: calc(100% - 28px);
  box-sizing: border-box;
}

.nodes-card :deep(.el-timeline-item__content) {
  width: 100%;
  box-sizing: border-box;
}

.nodes-card :deep(.el-timeline-item__tail) {
  left: 4px;
}

.nodes-card :deep(.el-timeline-item__node) {
  left: 0;
}

.card-header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 16px;
  color: #303133;
}

.node-log-card {
  margin-bottom: 0;
  width: 100%;
  box-sizing: border-box;
}

.node-log-card.success {
  border-left: 4px solid #67c23a;
}

.node-log-card.error {
  border-left: 4px solid #f56c6c;
}

.node-log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.node-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.node-name {
  font-weight: 600;
  font-size: 15px;
  color: #303133;
}

.node-status {
  display: flex;
  align-items: center;
  gap: 12px;
}

.node-duration {
  font-size: 13px;
  color: #909399;
}

.node-detail {
  padding: 12px 0;
}

.collapse-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
}

.detail-section {
  margin-bottom: 16px;
}

.input-section,
.output-section {
  margin-bottom: 0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  font-size: 14px;
  color: #606266;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 2px solid #e4e7ed;
}

.input-section .section-title {
  border-bottom-color: #409eff;
}

.output-section .section-title {
  border-bottom-color: #67c23a;
}

.error-section .section-title {
  border-bottom-color: #f56c6c;
}

.detail-content-wrapper {
  background: #fafafa;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  overflow: hidden;
}

.detail-content {
  background: #f5f7fa;
  padding: 14px;
  font-size: 12px;
  line-height: 1.8;
  max-height: 350px;
  overflow: auto;
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  color: #2c3e50;
}

.error-section .section-title {
  color: #f56c6c;
}

.error-content-wrapper {
  background: #fef0f0;
  border: 1px solid #fde2e2;
  border-radius: 6px;
  overflow: hidden;
}

.error-content {
  background: #fef0f0;
  color: #f56c6c;
  padding: 14px;
  font-size: 12px;
  line-height: 1.8;
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
}

.log-content {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
  font-size: 13px;
  line-height: 1.6;
  max-height: 400px;
  overflow: auto;
  border: 1px solid #e4e7ed;
}

.log-error {
  background: #fef0f0;
  color: #f56c6c;
  padding: 16px;
  border-radius: 4px;
  font-size: 13px;
  line-height: 1.6;
  border: 1px solid #fde2e2;
}
</style>
