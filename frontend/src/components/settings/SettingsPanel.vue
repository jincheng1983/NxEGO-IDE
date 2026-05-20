<script lang="ts" setup>
import { ref, computed, watch, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { EventsOn } from '../../../wailsjs/runtime/runtime'
import { useTheme, type ThemeMode, type ProjectSettings, type AISettings } from '../../stores/theme'
import { useAIStore, type AIModel, type AIAgent } from '../../stores/ai'

const emit = defineEmits<{
  (e: 'close'): void
}>()

const { currentTheme, settings, setTheme, updateEditorSettings, updateGeneralSettings, updateAISettings, updateCustomColors, resetSettings, defaultProjectSettings, loadProjectSettings, saveProjectSettings } = useTheme()
const {
  allModels, localModels, ollamaModels, cloudModels,
  agents: aiAgents, currentModel, currentAgent,
  selectModel, selectAgent, addAgent, updateAgent, removeAgent,
  addCustomModel, updateCustomModel, removeCustomModel, downloadModel, setModelInstalled,
  updateModelProgress, checkInstalledModels,
} = useAIStore()

const activeTab = ref('appearance')

onMounted(() => {
  checkInstalledModels()
  EventsOn('modelDownloadProgress', (data: any) => {
    if (data && data.modelId) {
      updateModelProgress(data.modelId, data.progress || 0, data.status || 'downloading')
    }
  })
})

const currentProjectId = ref('')
const projectSettings = ref<ProjectSettings>({ ...defaultProjectSettings })

watch(() => settings.general, () => {}, { deep: true })

function initProjectSettings(projectId: string, projectName: string) {
  currentProjectId.value = projectId
  projectSettings.value = loadProjectSettings(projectId)
  if (!projectSettings.value.projectId) {
    projectSettings.value.projectId = projectId
  }
  if (projectName) {
    projectSettings.value.projectName = projectName
  }
}

function handleProjectSettingsChange() {
  if (currentProjectId.value) {
    saveProjectSettings({ ...projectSettings.value })
  }
}

const themeOptions: { label: string; value: ThemeMode; icon: string; desc: string }[] = [
  { label: '浅色主题', value: 'light', icon: '☀️', desc: '明亮清晰的浅色界面，适合白天使用' },
  { label: '深色主题', value: 'dark', icon: '🌙', desc: '护眼舒适的深色界面，适合夜间编码' },
  { label: '自定义主题', value: 'custom', icon: '🎨', desc: '自由定制配色方案，打造专属IDE' },
]

const fontOptions = [
  "'Consolas', 'Courier New', monospace",
  "'Fira Code', 'Consolas', monospace",
  "'JetBrains Mono', 'Consolas', monospace",
  "'Source Code Pro', 'Consolas', monospace",
  "'Microsoft YaHei', 'Consolas', monospace",
]

const fontSizeOptions = [12, 13, 14, 15, 16, 18, 20, 22, 24]
const tabSizeOptions = [2, 4, 6, 8]
const lineHeightOptions = [1.2, 1.4, 1.6, 1.8, 2.0]
const autoSaveOptions = [30, 60, 120, 300, 600]

const handleThemeChange = (theme: ThemeMode) => {
  setTheme(theme)
}

const handleReset = () => {
  if (confirm('确定要恢复所有设置为默认值吗？此操作不可撤销。')) {
    resetSettings()
  }
}

const showCustomColors = computed(() => currentTheme.value === 'custom')

// AI agent/model management state
const showAIAddModel = ref(false)
const showAIAddAgent = ref(false)
const aiNewModel = reactive<Partial<AIModel>>({ provider: 'local', supportsImage: false })
const aiNewAgent = reactive<Partial<AIAgent>>({ icon: 'Cpu' })
const aiEditingAgentId = ref<string | null>(null)
const downloadingModelId = ref<string | null>(null)
const downloadError = ref('')
const aiModelTab = ref('local')

const providerLabel = (provider: string) => {
  switch (provider) {
    case 'local': return '本地'
    case 'ollama': return 'Ollama'
    case 'cloud': return '云端'
    default: return provider
  }
}

const downloadStatusLabel = (status: string | undefined) => {
  switch (status) {
    case 'not_installed': return '未安装'
    case 'downloading': return '下载中...'
    case 'installed': return '已安装'
    case 'error': return '下载失败'
    default: return '未知'
  }
}

const handleDownloadModel = async (modelId: string) => {
  downloadingModelId.value = modelId
  downloadError.value = ''
  try {
    await downloadModel(modelId)
    ElMessage.success('模型下载完成！')
  } catch (e: any) {
    downloadError.value = e.message || '下载失败'
    ElMessage.error(`模型下载失败: ${downloadError.value}`)
  } finally {
    downloadingModelId.value = null
  }
}

const handleAddCustomModel = () => {
  if (!aiNewModel.id || !aiNewModel.name) {
    ElMessage.warning('请填写模型 ID 和名称')
    return
  }
  addCustomModel({
    id: aiNewModel.id,
    name: aiNewModel.name,
    provider: aiNewModel.provider || 'local',
    apiUrl: aiNewModel.apiUrl,
    apiKey: aiNewModel.apiKey,
    supportsImage: aiNewModel.supportsImage || false,
    isActive: aiNewModel.provider === 'cloud' ? true : undefined,
  } as AIModel)
  ElMessage.success('模型添加成功')
  showAIAddModel.value = false
  aiNewModel.id = ''
  aiNewModel.name = ''
  aiNewModel.apiUrl = ''
  aiNewModel.apiKey = ''
  aiNewModel.provider = 'local'
  aiNewModel.supportsImage = false
}

const handleAddAgent = () => {
  if (!aiNewAgent.name) {
    ElMessage.warning('请填写智能体名称')
    return
  }
  if (aiEditingAgentId.value) {
    updateAgent(aiEditingAgentId.value, aiNewAgent)
    ElMessage.success('智能体已更新')
  } else {
    addAgent({
      id: `agent_${Date.now()}`,
      name: aiNewAgent.name || '',
      icon: aiNewAgent.icon || 'Cpu',
      description: aiNewAgent.description || '',
      systemPrompt: aiNewAgent.systemPrompt || '',
    } as AIAgent)
    ElMessage.success('智能体添加成功')
  }
  showAIAddAgent.value = false
  aiEditingAgentId.value = null
  aiNewAgent.name = ''
  aiNewAgent.icon = 'Cpu'
  aiNewAgent.description = ''
  aiNewAgent.systemPrompt = ''
}

const handleEditAgent = (agent: AIAgent) => {
  aiEditingAgentId.value = agent.id
  aiNewAgent.name = agent.name
  aiNewAgent.icon = agent.icon
  aiNewAgent.description = agent.description
  aiNewAgent.systemPrompt = agent.systemPrompt
  showAIAddAgent.value = true
}

const handleRemoveAgent = async (agentId: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个智能体吗？', '确认删除', { type: 'warning' })
    removeAgent(agentId)
    ElMessage.success('智能体已删除')
  } catch {}
}

const handleRemoveModel = async (modelId: string) => {
  try {
    await ElMessageBox.confirm('确定要移除此模型吗？', '确认移除', { type: 'warning' })
    removeCustomModel(modelId)
    ElMessage.success('模型已移除')
  } catch {}
}
</script>

<template>
  <el-dialog
    :model-value="true"
    title="软件设置"
    width="720px"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    @close="emit('close')"
    class="settings-dialog"
  >
    <div class="settings-container">
      <div class="settings-sidebar">
        <div
          v-for="tab in [
            { key: 'appearance', label: '外观', icon: 'Brush' },
            { key: 'editor', label: '编辑器', icon: 'Edit' },
            { key: 'general', label: '通用', icon: 'Setting' },
            { key: 'ai', label: 'AI', icon: 'Cpu' },
            { key: 'project', label: '项目', icon: 'FolderOpened' },
          ]"
          :key="tab.key"
          :class="['settings-tab', { active: activeTab === tab.key }]"
          @click="activeTab = tab.key"
        >
          <el-icon><component :is="tab.icon" /></el-icon>
          <span>{{ tab.label }}</span>
        </div>
      </div>

      <div class="settings-content">
        <!-- ========== 外观设置 ========== -->
        <div v-show="activeTab === 'appearance'" class="settings-section">
          <h3 class="section-title">主题选择</h3>
          <div class="theme-cards">
            <div
              v-for="opt in themeOptions"
              :key="opt.value"
              :class="['theme-card', { selected: currentTheme === opt.value }]"
              @click="handleThemeChange(opt.value)"
            >
              <div class="theme-preview" :class="`theme-preview-${opt.value}`">
                <div class="preview-menu"></div>
                <div class="preview-body">
                  <div class="preview-sidebar"></div>
                  <div class="preview-editor">
                    <div class="preview-line"></div>
                    <div class="preview-line short"></div>
                    <div class="preview-line"></div>
                  </div>
                </div>
              </div>
              <div class="theme-info">
                <span class="theme-icon">{{ opt.icon }}</span>
                <span class="theme-label">{{ opt.label }}</span>
              </div>
              <div class="theme-desc">{{ opt.desc }}</div>
            </div>
          </div>

          <!-- 自定义主题颜色 -->
          <div v-if="showCustomColors" class="custom-colors">
            <h3 class="section-title">自定义配色</h3>
            <div class="color-grid">
              <div class="color-item">
                <label>主背景色</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.bgPrimary"
                    @input="(e: Event) => updateCustomColors({ bgPrimary: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.bgPrimary }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>次背景色</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.bgSecondary"
                    @input="(e: Event) => updateCustomColors({ bgSecondary: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.bgSecondary }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>主文字色</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.textPrimary"
                    @input="(e: Event) => updateCustomColors({ textPrimary: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.textPrimary }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>品牌色</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.colorPrimary"
                    @input="(e: Event) => updateCustomColors({ colorPrimary: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.colorPrimary }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>编辑器背景</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.editorBg"
                    @input="(e: Event) => updateCustomColors({ editorBg: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.editorBg }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>编辑器文字</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.editorText"
                    @input="(e: Event) => updateCustomColors({ editorText: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.editorText }}</span>
                </div>
              </div>
            </div>

            <h4 class="sub-title">语法高亮颜色</h4>
            <div class="color-grid">
              <div class="color-item">
                <label>关键字</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.hlKeyword"
                    @input="(e: Event) => updateCustomColors({ hlKeyword: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.hlKeyword }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>字符串</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.hlString"
                    @input="(e: Event) => updateCustomColors({ hlString: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.hlString }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>注释</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.hlComment"
                    @input="(e: Event) => updateCustomColors({ hlComment: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.hlComment }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>E语言关键字</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.hlElang"
                    @input="(e: Event) => updateCustomColors({ hlElang: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.hlElang }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>函数名</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.hlFunction"
                    @input="(e: Event) => updateCustomColors({ hlFunction: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.hlFunction }}</span>
                </div>
              </div>
              <div class="color-item">
                <label>数字</label>
                <div class="color-input-row">
                  <input
                    type="color"
                    :value="settings.customColors.hlNumber"
                    @input="(e: Event) => updateCustomColors({ hlNumber: (e.target as HTMLInputElement).value })"
                  />
                  <span class="color-value">{{ settings.customColors.hlNumber }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- ========== 编辑器设置 ========== -->
        <div v-show="activeTab === 'editor'" class="settings-section">
          <h3 class="section-title">字体设置</h3>
          <div class="setting-row">
            <label>字体</label>
            <el-select
              :model-value="settings.editor.fontFamily"
              @update:model-value="(v: string) => updateEditorSettings({ fontFamily: v })"
              size="small"
              style="width: 280px"
            >
              <el-option
                v-for="font in fontOptions"
                :key="font"
                :label="font"
                :value="font"
              />
            </el-select>
          </div>
          <div class="setting-row">
            <label>字号</label>
            <el-select
              :model-value="settings.editor.fontSize"
              @update:model-value="(v: number) => updateEditorSettings({ fontSize: v })"
              size="small"
              style="width: 100px"
            >
              <el-option
                v-for="size in fontSizeOptions"
                :key="size"
                :label="`${size}px`"
                :value="size"
              />
            </el-select>
          </div>
          <div class="setting-row">
            <label>行高</label>
            <el-select
              :model-value="settings.editor.lineHeight"
              @update:model-value="(v: number) => updateEditorSettings({ lineHeight: v })"
              size="small"
              style="width: 100px"
            >
              <el-option
                v-for="lh in lineHeightOptions"
                :key="lh"
                :label="String(lh)"
                :value="lh"
              />
            </el-select>
          </div>

          <h3 class="section-title">缩进与格式</h3>
          <div class="setting-row">
            <label>Tab 大小</label>
            <el-select
              :model-value="settings.editor.tabSize"
              @update:model-value="(v: number) => updateEditorSettings({ tabSize: v })"
              size="small"
              style="width: 100px"
            >
              <el-option
                v-for="ts in tabSizeOptions"
                :key="ts"
                :label="String(ts)"
                :value="ts"
              />
            </el-select>
          </div>
          <div class="setting-row">
            <label>自动换行</label>
            <el-switch
              :model-value="settings.editor.wordWrap"
              @update:model-value="(v: boolean) => updateEditorSettings({ wordWrap: v })"
            />
          </div>

          <h3 class="section-title">显示</h3>
          <div class="setting-row">
            <label>显示行号</label>
            <el-switch
              :model-value="settings.editor.showLineNumbers"
              @update:model-value="(v: boolean) => updateEditorSettings({ showLineNumbers: v })"
            />
          </div>
          <div class="setting-row">
            <label>代码自动补全</label>
            <el-switch
              :model-value="settings.editor.autoComplete"
              @update:model-value="(v: boolean) => updateEditorSettings({ autoComplete: v })"
            />
          </div>
          <div class="setting-row">
            <label>迷你地图</label>
            <el-switch
              :model-value="settings.editor.minimap"
              @update:model-value="(v: boolean) => updateEditorSettings({ minimap: v })"
            />
          </div>

          <!-- 预览 -->
          <div class="editor-preview" :style="{
            fontFamily: settings.editor.fontFamily,
            fontSize: settings.editor.fontSize + 'px',
            lineHeight: settings.editor.lineHeight,
          }">
            <div class="preview-code">
              <span class="hl-elang">.子程序</span> <span class="hl-func">启动子程序</span>, <span class="hl-type">整数型</span>, <span class="hl-elang">公开</span><br/>
              <span class="hl-elang">.参数</span> <span class="hl-var">参数1</span>, <span class="hl-type">整数型</span>, <span class="hl-elang">参考</span><br/>
              <span class="hl-elang">.局部变量</span> <span class="hl-var">变量名</span>, <span class="hl-type">整数型</span>, <span class="hl-elang">静态</span><br/>
              <span class="hl-elang">.如果</span> (<span class="hl-var">变量名</span> ＝ <span class="hl-num">1</span>)<br/>
              &nbsp;&nbsp;<span class="hl-func">调试输出</span> (<span class="hl-str">"变量名是1"</span>)<br/>
              <span class="hl-elang">.如果结束</span><br/>
              <span class="hl-elang">.返回</span> (<span class="hl-num">0</span>)
            </div>
          </div>
        </div>

        <!-- ========== 通用设置 ========== -->
        <div v-show="activeTab === 'general'" class="settings-section">
          <h3 class="section-title">基本设置</h3>
          <div class="setting-row">
            <label>界面语言</label>
            <el-select
              :model-value="settings.general.language"
              @update:model-value="(v: string) => updateGeneralSettings({ language: v })"
              size="small"
              style="width: 160px"
            >
              <el-option label="简体中文" value="zh-CN" />
              <el-option label="English" value="en-US" />
            </el-select>
          </div>
          <div class="setting-row">
            <label>退出时确认</label>
            <el-switch
              :model-value="settings.general.confirmOnExit"
              @update:model-value="(v: boolean) => updateGeneralSettings({ confirmOnExit: v })"
            />
          </div>
          <div class="setting-row">
            <label>显示欢迎页</label>
            <el-switch
              :model-value="settings.general.showWelcomePage"
              @update:model-value="(v: boolean) => updateGeneralSettings({ showWelcomePage: v })"
            />
          </div>

          <h3 class="section-title">自动保存</h3>
          <div class="setting-row">
            <label>启用自动保存</label>
            <el-switch
              :model-value="settings.general.autoSave"
              @update:model-value="(v: boolean) => updateGeneralSettings({ autoSave: v })"
            />
          </div>
          <div class="setting-row" v-if="settings.general.autoSave">
            <label>保存间隔</label>
            <el-select
              :model-value="settings.general.autoSaveInterval"
              @update:model-value="(v: number) => updateGeneralSettings({ autoSaveInterval: v })"
              size="small"
              style="width: 140px"
            >
              <el-option
                v-for="interval in autoSaveOptions"
                :key="interval"
                :label="interval >= 60 ? `${interval / 60} 分钟` : `${interval} 秒`"
                :value="interval"
              />
            </el-select>
          </div>
        </div>

        <!-- ========== AI 设置 ========== -->
        <div v-show="activeTab === 'ai'" class="settings-section">
          <h3 class="section-title">🧠 智能体管理</h3>
          <div class="ai-manage-list">
            <div v-for="agent in aiAgents" :key="agent.id" class="ai-manage-item">
              <div class="ai-item-icon">
                <el-icon :size="18"><component :is="agent.icon" /></el-icon>
              </div>
              <div class="ai-item-info">
                <div class="ai-item-name">{{ agent.name }}</div>
                <div class="ai-item-desc">{{ agent.description }}</div>
              </div>
              <div class="ai-item-actions">
                <el-tag v-if="currentAgent.id === agent.id" type="success" size="small" effect="plain">当前</el-tag>
                <el-button v-else size="small" text @click="selectAgent(agent.id)">使用</el-button>
                <el-button size="small" text @click="handleEditAgent(agent)" :disabled="agent.isBuiltin">✏️</el-button>
                <el-button size="small" text @click="handleRemoveAgent(agent.id)" :disabled="agent.isBuiltin" style="color: #f56c6c">🗑️</el-button>
              </div>
            </div>
          </div>
          <el-button size="small" type="primary" plain @click="showAIAddAgent = true; aiEditingAgentId = null; aiNewAgent.name = ''; aiNewAgent.icon = 'Cpu'; aiNewAgent.description = ''; aiNewAgent.systemPrompt = ''" style="margin-top: 8px">+ 新增智能体</el-button>

          <h3 class="section-title">🤖 模型管理</h3>
          <div class="ai-tabs">
            <el-radio-group v-model="aiModelTab" size="small">
              <el-radio-button value="local">本地模型</el-radio-button>
              <el-radio-button value="ollama">Ollama</el-radio-button>
              <el-radio-button value="cloud">云端 API</el-radio-button>
            </el-radio-group>
          </div>

          <div class="ai-manage-list">
            <template v-if="aiModelTab === 'local'">
              <div v-for="model in localModels" :key="model.id" class="ai-manage-item">
                <div class="ai-item-icon">
                  <el-icon :size="18"><Cpu /></el-icon>
                </div>
                <div class="ai-item-info">
                  <div class="ai-item-name">{{ model.name }}</div>
                  <div class="ai-item-desc">{{ model.modelSize }} · GGUF · {{ downloadStatusLabel(model.downloadStatus) }}</div>
                  <el-progress
                    v-if="model.downloadStatus === 'downloading'"
                    :percentage="model.downloadProgress || 0"
                    :stroke-width="4"
                    style="width: 140px; margin-top: 4px"
                  />
                  <span v-if="model.downloadStatus === 'error'" style="color: #f56c6c; font-size: 12px">{{ downloadError }}</span>
                </div>
                <div class="ai-item-actions">
                  <el-tag v-if="currentModel.id === model.id" type="success" size="small" effect="plain">当前</el-tag>
                  <el-button v-else-if="model.downloadStatus === 'installed'" size="small" text @click="selectModel(model.id)">使用</el-button>
                  <el-button
                    v-if="model.downloadStatus === 'not_installed'"
                    size="small"
                    type="primary"
                    :loading="downloadingModelId === model.id"
                    @click="handleDownloadModel(model.id)"
                  >下载</el-button>
                  <el-button
                    v-if="model.downloadStatus === 'error'"
                    size="small"
                    type="warning"
                    @click="handleDownloadModel(model.id)"
                  >重试</el-button>
                </div>
              </div>
              <div v-if="localModels.length === 0" class="ai-empty">暂无本地模型</div>
            </template>

            <template v-if="aiModelTab === 'ollama'">
              <div v-for="model in ollamaModels" :key="model.id" class="ai-manage-item">
                <div class="ai-item-icon">
                  <el-icon :size="18"><Cpu /></el-icon>
                </div>
                <div class="ai-item-info">
                  <div class="ai-item-name">{{ model.name }}</div>
                  <div class="ai-item-desc">Ollama · {{ model.apiUrl }}</div>
                </div>
                <div class="ai-item-actions">
                  <el-tag v-if="currentModel.id === model.id" type="success" size="small" effect="plain">当前</el-tag>
                  <el-button v-else size="small" text @click="selectModel(model.id)">使用</el-button>
                </div>
              </div>
              <div v-if="ollamaModels.length === 0" class="ai-empty">暂无 Ollama 模型</div>
            </template>

            <template v-if="aiModelTab === 'cloud'">
              <div v-for="model in cloudModels" :key="model.id" class="ai-manage-item">
                <div class="ai-item-icon">
                  <el-icon :size="18"><Connection /></el-icon>
                </div>
                <div class="ai-item-info">
                  <div class="ai-item-name">{{ model.name }}</div>
                  <div class="ai-item-desc">云端 API · {{ model.apiUrl }} · {{ model.supportsImage ? '支持图片' : '纯文本' }}</div>
                </div>
                <div class="ai-item-actions">
                  <el-tag v-if="currentModel.id === model.id" type="success" size="small" effect="plain">当前</el-tag>
                  <el-button v-else size="small" text @click="selectModel(model.id)">使用</el-button>
                  <el-button size="small" text @click="handleRemoveModel(model.id)" style="color: #f56c6c">🗑️</el-button>
                </div>
              </div>
              <div v-if="cloudModels.length === 0" class="ai-empty">暂无云端模型</div>
            </template>
          </div>
          <el-button size="small" type="primary" plain @click="showAIAddModel = true" style="margin-top: 8px">+ 添加自定义模型</el-button>

          <h3 class="section-title">⚙️ 推理参数</h3>
          <div class="setting-row">
            <label>推理引擎</label>
            <el-select
              :model-value="settings.ai.provider"
              @update:model-value="(v: string) => updateAISettings({ provider: v as AISettings['provider'] })"
              size="small"
              style="width: 160px"
            >
              <el-option label="本地模型 (llama.cpp)" value="local" />
              <el-option label="Ollama" value="ollama" />
              <el-option label="云端 API" value="cloud" />
            </el-select>
          </div>

          <template v-if="settings.ai.provider === 'local'">
            <div class="setting-row">
              <label>模型文件路径</label>
              <el-input
                :model-value="settings.ai.localModelPath"
                @update:model-value="(v: string) => updateAISettings({ localModelPath: v })"
                size="small"
                placeholder="留空则使用内置模型"
                style="width: 260px"
              />
            </div>
            <div class="setting-row">
              <label>模型名称</label>
              <el-input
                :model-value="settings.ai.localModelName"
                @update:model-value="(v: string) => updateAISettings({ localModelName: v })"
                size="small"
                placeholder="qwen2.5-coder-1.5b-q4_k_m"
                style="width: 260px"
              />
            </div>
            <div class="setting-row">
              <label>启用 GPU 加速</label>
              <el-switch
                :model-value="settings.ai.enableGPU"
                @update:model-value="(v: boolean) => updateAISettings({ enableGPU: v })"
              />
            </div>
            <div class="setting-row" v-if="settings.ai.enableGPU">
              <label>GPU 推理层数</label>
              <el-input-number
                :model-value="settings.ai.gpuLayers"
                @update:model-value="(v: number | undefined) => updateAISettings({ gpuLayers: v ?? 28 })"
                size="small"
                :min="0"
                :max="99"
                controls-position="right"
                style="width: 120px"
              />
            </div>
          </template>

          <template v-if="settings.ai.provider === 'ollama'">
            <div class="setting-row">
              <label>服务地址</label>
              <el-input
                :model-value="settings.ai.ollamaUrl"
                @update:model-value="(v: string) => updateAISettings({ ollamaUrl: v })"
                size="small"
                placeholder="http://localhost:11434"
                style="width: 260px"
              />
            </div>
            <div class="setting-row">
              <label>模型名称</label>
              <el-input
                :model-value="settings.ai.ollamaModel"
                @update:model-value="(v: string) => updateAISettings({ ollamaModel: v })"
                size="small"
                placeholder="qwen2.5-coder:1.5b"
                style="width: 260px"
              />
            </div>
          </template>

          <template v-if="settings.ai.provider === 'cloud'">
            <div class="setting-row">
              <label>API 地址</label>
              <el-input
                :model-value="settings.ai.cloudApiUrl"
                @update:model-value="(v: string) => updateAISettings({ cloudApiUrl: v })"
                size="small"
                placeholder="https://api.openai.com/v1"
                style="width: 260px"
              />
            </div>
            <div class="setting-row">
              <label>API Key</label>
              <el-input
                :model-value="settings.ai.cloudApiKey"
                @update:model-value="(v: string) => updateAISettings({ cloudApiKey: v })"
                size="small"
                type="password"
                show-password
                placeholder="sk-..."
                style="width: 260px"
              />
            </div>
            <div class="setting-row">
              <label>模型名称</label>
              <el-input
                :model-value="settings.ai.cloudModel"
                @update:model-value="(v: string) => updateAISettings({ cloudModel: v })"
                size="small"
                placeholder="gpt-4o-mini"
                style="width: 260px"
              />
            </div>
          </template>

          <div class="setting-row">
            <label>最大 Token 数</label>
            <el-input-number
              :model-value="settings.ai.maxTokens"
              @update:model-value="(v: number | undefined) => updateAISettings({ maxTokens: v ?? 2048 })"
              size="small"
              :min="256"
              :max="32768"
              :step="256"
              controls-position="right"
              style="width: 140px"
            />
          </div>
          <div class="setting-row">
            <label>温度 (Temperature)</label>
            <el-slider
              :model-value="settings.ai.temperature"
              @update:model-value="(v: number) => updateAISettings({ temperature: v })"
              size="small"
              :min="0"
              :max="2"
              :step="0.1"
              style="width: 200px"
              :show-input="true"
              :show-input-controls="false"
            />
          </div>
          <div class="setting-row">
            <label>Top-P</label>
            <el-slider
              :model-value="settings.ai.topP"
              @update:model-value="(v: number) => updateAISettings({ topP: v })"
              size="small"
              :min="0"
              :max="1"
              :step="0.05"
              style="width: 200px"
              :show-input="true"
              :show-input-controls="false"
            />
          </div>
          <div class="setting-row">
            <label>上下文大小</label>
            <el-select
              :model-value="settings.ai.contextSize"
              @update:model-value="(v: number) => updateAISettings({ contextSize: v })"
              size="small"
              style="width: 140px"
            >
              <el-option label="4K" :value="4096" />
              <el-option label="8K" :value="8192" />
              <el-option label="16K" :value="16384" />
              <el-option label="32K" :value="32768" />
            </el-select>
          </div>

          <div class="setting-row">
            <label>启用代码补全</label>
            <el-switch
              :model-value="settings.ai.enableCodeCompletion"
              @update:model-value="(v: boolean) => updateAISettings({ enableCodeCompletion: v })"
            />
          </div>
          <div class="setting-row" v-if="settings.ai.enableCodeCompletion">
            <label>补全延迟 (ms)</label>
            <el-input-number
              :model-value="settings.ai.completionDelay"
              @update:model-value="(v: number | undefined) => updateAISettings({ completionDelay: v ?? 150 })"
              size="small"
              :min="50"
              :max="1000"
              :step="50"
              controls-position="right"
              style="width: 120px"
            />
          </div>
          <div class="setting-row">
            <label>保存对话历史</label>
            <el-switch
              :model-value="settings.ai.enableChatHistory"
              @update:model-value="(v: boolean) => updateAISettings({ enableChatHistory: v })"
            />
          </div>

          <div class="setting-row" style="flex-direction: column; align-items: flex-start; gap: 8px;">
            <label>System Prompt</label>
            <el-input
              :model-value="settings.ai.systemPrompt"
              @update:model-value="(v: string) => updateAISettings({ systemPrompt: v })"
              type="textarea"
              :rows="3"
              size="small"
              placeholder="设置 AI 助手的角色和行为..."
              style="width: 100%"
            />
          </div>

          <!-- 新增智能体弹窗 -->
          <el-dialog
            v-model="showAIAddAgent"
            :title="aiEditingAgentId ? '编辑智能体' : '新增智能体'"
            width="480px"
            :close-on-click-modal="false"
          >
            <div class="ai-form">
              <div class="ai-form-item">
                <label>名称</label>
                <el-input v-model="aiNewAgent.name" placeholder="智能体名称" size="small" />
              </div>
              <div class="ai-form-item">
                <label>图标</label>
                <el-select v-model="aiNewAgent.icon" size="small" style="width: 100%">
                  <el-option label="CPU" value="Cpu" />
                  <el-option label="文件" value="Files" />
                  <el-option label="Bug" value="Bug" />
                  <el-option label="切换" value="Switch" />
                  <el-option label="查看" value="View" />
                  <el-option label="编辑" value="Edit" />
                  <el-option label="聊天" value="ChatDotRound" />
                  <el-option label="搜索" value="Search" />
                </el-select>
              </div>
              <div class="ai-form-item">
                <label>描述</label>
                <el-input v-model="aiNewAgent.description" placeholder="智能体描述" size="small" />
              </div>
              <div class="ai-form-item">
                <label>系统提示词</label>
                <el-input v-model="aiNewAgent.systemPrompt" type="textarea" :rows="4" placeholder="定义智能体的角色和行为..." size="small" />
              </div>
            </div>
            <template #footer>
              <el-button size="small" @click="showAIAddAgent = false">取消</el-button>
              <el-button size="small" type="primary" @click="handleAddAgent">确定</el-button>
            </template>
          </el-dialog>

          <!-- 新增自定义模型弹窗 -->
          <el-dialog
            v-model="showAIAddModel"
            title="添加自定义模型"
            width="480px"
            :close-on-click-modal="false"
          >
            <div class="ai-form">
              <div class="ai-form-item">
                <label>模型 ID</label>
                <el-input v-model="aiNewModel.id" placeholder="模型的唯一标识" size="small" />
              </div>
              <div class="ai-form-item">
                <label>模型名称</label>
                <el-input v-model="aiNewModel.name" placeholder="显示名称" size="small" />
              </div>
              <div class="ai-form-item">
                <label>提供商</label>
                <el-select v-model="aiNewModel.provider" size="small" style="width: 100%">
                  <el-option label="本地 (llama.cpp)" value="local" />
                  <el-option label="Ollama" value="ollama" />
                  <el-option label="云端 API" value="cloud" />
                </el-select>
              </div>
              <div class="ai-form-item">
                <label>API 地址</label>
                <el-input v-model="aiNewModel.apiUrl" placeholder="API 端点地址" size="small" />
              </div>
              <div class="ai-form-item">
                <label>API Key</label>
                <el-input v-model="aiNewModel.apiKey" type="password" show-password placeholder="API 密钥" size="small" />
              </div>
              <div class="ai-form-item">
                <label>支持图片输入</label>
                <el-switch v-model="aiNewModel.supportsImage" />
              </div>
            </div>
            <template #footer>
              <el-button size="small" @click="showAIAddModel = false">取消</el-button>
              <el-button size="small" type="primary" @click="handleAddCustomModel">确定</el-button>
            </template>
          </el-dialog>
        </div>

        <!-- ========== 项目设置 ========== -->
        <div v-show="activeTab === 'project'" class="settings-section">
          <div v-if="!currentProjectId" class="project-empty">
            <el-icon :size="40"><FolderOpened /></el-icon>
            <p>请先打开一个项目以查看项目设置</p>
            <p class="hint">项目设置仅对当前项目生效，独立于全局设置</p>
          </div>
          <template v-else>
            <h3 class="section-title">项目信息</h3>
            <div class="setting-row">
              <label>项目名称</label>
              <span class="setting-value">{{ projectSettings.projectName || '-' }}</span>
            </div>
            <div class="setting-row">
              <label>项目ID</label>
              <span class="setting-value mono">{{ projectSettings.projectId }}</span>
            </div>

            <h3 class="section-title">构建设置</h3>
            <div class="setting-row">
              <label>启动窗口</label>
              <el-input
                :model-value="projectSettings.startupWindow"
                @update:model-value="(v: string) => { projectSettings.startupWindow = v; handleProjectSettingsChange() }"
                size="small"
                placeholder="主窗口名称"
                style="width: 200px"
              />
            </div>
            <div class="setting-row">
              <label>目标平台</label>
              <el-select
                :model-value="projectSettings.targetPlatform"
                @update:model-value="(v: string) => { projectSettings.targetPlatform = v; handleProjectSettingsChange() }"
                size="small"
                style="width: 200px"
              >
                <el-option label="Windows 64位" value="windows/amd64" />
                <el-option label="Windows 32位" value="windows/386" />
                <el-option label="Linux 64位" value="linux/amd64" />
                <el-option label="macOS 64位" value="darwin/amd64" />
              </el-select>
            </div>
            <div class="setting-row">
              <label>构建模式</label>
              <el-select
                :model-value="projectSettings.buildMode"
                @update:model-value="(v: string) => { projectSettings.buildMode = v; handleProjectSettingsChange() }"
                size="small"
                style="width: 200px"
              >
                <el-option label="发布版 (Release)" value="release" />
                <el-option label="调试版 (Debug)" value="debug" />
              </el-select>
            </div>
            <div class="setting-row">
              <label>输出目录</label>
              <el-input
                :model-value="projectSettings.outputDir"
                @update:model-value="(v: string) => { projectSettings.outputDir = v; handleProjectSettingsChange() }"
                size="small"
                placeholder="build"
                style="width: 200px"
              />
            </div>

            <h3 class="section-title">版本信息</h3>
            <div class="setting-row">
              <label>版本号</label>
              <el-input
                :model-value="projectSettings.version"
                @update:model-value="(v: string) => { projectSettings.version = v; handleProjectSettingsChange() }"
                size="small"
                placeholder="1.0.0"
                style="width: 160px"
              />
            </div>
            <div class="setting-row">
              <label>注入版本信息</label>
              <el-switch
                :model-value="projectSettings.injectVersion"
                @update:model-value="(v: boolean) => { projectSettings.injectVersion = v; handleProjectSettingsChange() }"
              />
            </div>

            <h3 class="section-title">高级选项</h3>
            <div class="setting-row">
              <label>压缩输出</label>
              <el-switch
                :model-value="projectSettings.compress"
                @update:model-value="(v: boolean) => { projectSettings.compress = v; handleProjectSettingsChange() }"
              />
            </div>
            <div class="setting-row">
              <label>代码混淆</label>
              <el-switch
                :model-value="projectSettings.obfuscate"
                @update:model-value="(v: boolean) => { projectSettings.obfuscate = v; handleProjectSettingsChange() }"
              />
            </div>
            <div class="setting-row">
              <label>调试模式</label>
              <el-switch
                :model-value="projectSettings.debug"
                @update:model-value="(v: boolean) => { projectSettings.debug = v; handleProjectSettingsChange() }"
              />
            </div>
            <div class="setting-row">
              <label>程序图标</label>
              <el-input
                :model-value="projectSettings.iconPath"
                @update:model-value="(v: string) => { projectSettings.iconPath = v; handleProjectSettingsChange() }"
                size="small"
                placeholder="图标路径 (.ico)"
                style="width: 240px"
              />
            </div>
          </template>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="settings-footer">
        <el-button @click="handleReset" type="warning" plain size="small">
          <el-icon><RefreshLeft /></el-icon>
          恢复默认设置
        </el-button>
        <div class="footer-right">
          <el-button @click="emit('close')" size="small">取消</el-button>
          <el-button type="primary" @click="emit('close')" size="small">确定</el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<style scoped>
.settings-dialog :deep(.el-dialog__header) {
  background: var(--dialog-header-bg);
  border-bottom: 1px solid var(--dialog-border);
  padding: 12px 20px;
  margin: 0;
}

.settings-dialog :deep(.el-dialog__title) {
  color: var(--text-primary);
  font-size: 15px;
  font-weight: 600;
}

.settings-dialog :deep(.el-dialog__body) {
  padding: 0;
  background: var(--dialog-bg);
}

.settings-dialog :deep(.el-dialog__footer) {
  padding: 0;
}

.settings-container {
  display: flex;
  height: 480px;
}

.settings-sidebar {
  width: 160px;
  background: var(--bg-tertiary);
  border-right: 1px solid var(--border-base);
  padding: 8px 0;
  flex-shrink: 0;
}

.settings-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  cursor: pointer;
  color: var(--text-regular);
  font-size: 13px;
  transition: all 0.15s ease;
  border-left: 3px solid transparent;
}

.settings-tab:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.settings-tab.active {
  background: var(--bg-selected);
  color: var(--color-primary);
  border-left-color: var(--color-primary);
  font-weight: 600;
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
}

.settings-section {
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 20px 0 12px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-lighter);
}

.section-title:first-child {
  margin-top: 0;
}

.sub-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 16px 0 10px 0;
}

/* 主题卡片 */
.theme-cards {
  display: flex;
  gap: 12px;
}

.theme-card {
  flex: 1;
  border: 2px solid var(--border-base);
  border-radius: 8px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--bg-primary);
}

.theme-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-sm);
}

.theme-card.selected {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
}

.theme-preview {
  height: 80px;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
  border: 1px solid var(--border-light);
}

.theme-preview-light {
  background: #f5f7fa;
}

.theme-preview-light .preview-menu {
  height: 10px;
  background: #ffffff;
  border-bottom: 1px solid #e4e7ed;
}

.theme-preview-light .preview-body {
  display: flex;
  height: 70px;
}

.theme-preview-light .preview-sidebar {
  width: 30px;
  background: #ffffff;
  border-right: 1px solid #e4e7ed;
}

.theme-preview-light .preview-editor {
  flex: 1;
  background: #ffffff;
  padding: 6px;
}

.theme-preview-light .preview-line {
  height: 4px;
  background: #e4e7ed;
  margin-bottom: 4px;
  border-radius: 2px;
}

.theme-preview-light .preview-line.short {
  width: 60%;
}

.theme-preview-dark {
  background: #1e1e1e;
}

.theme-preview-dark .preview-menu {
  height: 10px;
  background: #2c2c2c;
  border-bottom: 1px solid #3e3e42;
}

.theme-preview-dark .preview-body {
  display: flex;
  height: 70px;
}

.theme-preview-dark .preview-sidebar {
  width: 30px;
  background: #252526;
  border-right: 1px solid #3e3e42;
}

.theme-preview-dark .preview-editor {
  flex: 1;
  background: #1e1e1e;
  padding: 6px;
}

.theme-preview-dark .preview-line {
  height: 4px;
  background: #3e3e42;
  margin-bottom: 4px;
  border-radius: 2px;
}

.theme-preview-dark .preview-line.short {
  width: 60%;
}

.theme-preview-custom {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.theme-preview-custom .preview-menu {
  height: 10px;
  background: rgba(255,255,255,0.2);
}

.theme-preview-custom .preview-body {
  display: flex;
  height: 70px;
}

.theme-preview-custom .preview-sidebar {
  width: 30px;
  background: rgba(255,255,255,0.15);
}

.theme-preview-custom .preview-editor {
  flex: 1;
  background: rgba(255,255,255,0.1);
  padding: 6px;
}

.theme-preview-custom .preview-line {
  height: 4px;
  background: rgba(255,255,255,0.3);
  margin-bottom: 4px;
  border-radius: 2px;
}

.theme-preview-custom .preview-line.short {
  width: 60%;
}

.theme-info {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
}

.theme-icon {
  font-size: 16px;
}

.theme-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.theme-desc {
  font-size: 11px;
  color: var(--text-secondary);
  line-height: 1.4;
}

/* 自定义颜色 */
.custom-colors {
  margin-top: 8px;
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.color-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.color-item label {
  font-size: 12px;
  color: var(--text-secondary);
}

.color-input-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-input-row input[type="color"] {
  width: 32px;
  height: 28px;
  border: 1px solid var(--border-base);
  border-radius: 4px;
  cursor: pointer;
  padding: 2px;
  background: var(--input-bg);
}

.color-value {
  font-size: 11px;
  color: var(--text-secondary);
  font-family: 'Consolas', monospace;
}

/* 设置行 */
.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-lighter);
}

.setting-row label {
  font-size: 13px;
  color: var(--text-regular);
}

/* 编辑器预览 */
.editor-preview {
  margin-top: 16px;
  padding: 12px 16px;
  background: var(--editor-bg);
  border: 1px solid var(--editor-border);
  border-radius: 6px;
  color: var(--editor-text);
}

.preview-code {
  white-space: pre-wrap;
}

.hl-elang { color: var(--hl-elang); }
.hl-func { color: var(--hl-function); }
.hl-type { color: var(--hl-type); }
.hl-var { color: var(--hl-variable); }
.hl-num { color: var(--hl-number); }
.hl-str { color: var(--hl-string); }

/* 底部 */
.settings-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  border-top: 1px solid var(--dialog-border);
  background: var(--dialog-header-bg);
}

.footer-right {
  display: flex;
  gap: 8px;
}

.project-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: var(--text-secondary);
  gap: 8px;
}

.project-empty p {
  font-size: 14px;
  margin: 0;
}

.project-empty .hint {
  font-size: 12px;
  color: var(--text-placeholder);
}

.setting-value {
  font-size: 13px;
  color: var(--text-primary);
}

.setting-value.mono {
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 12px;
  color: var(--text-secondary);
}

.ai-manage-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 160px;
  overflow-y: auto;
  border: 1px solid var(--border-lighter);
  border-radius: 6px;
  padding: 4px;
}

.ai-manage-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 6px;
  transition: background 0.15s ease;
}

.ai-manage-item:hover {
  background: var(--bg-hover);
}

.ai-item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: var(--bg-tertiary);
  border-radius: 6px;
  color: var(--color-primary);
  flex-shrink: 0;
}

.ai-item-info {
  flex: 1;
  min-width: 0;
}

.ai-item-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.ai-item-desc {
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 2px;
}

.ai-item-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.ai-tabs {
  margin-bottom: 8px;
}

.ai-empty {
  text-align: center;
  padding: 16px;
  color: var(--text-secondary);
  font-size: 12px;
}

.ai-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ai-form-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ai-form-item label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 600;
}
</style>