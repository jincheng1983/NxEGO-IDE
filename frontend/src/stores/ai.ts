import { ref, reactive, computed, watch } from 'vue'
import { DownloadModelFile, GetModelDownloadStatus } from '../../wailsjs/go/main/App'

export interface AIModel {
  id: string
  name: string
  provider: 'local' | 'ollama' | 'cloud'
  apiUrl?: string
  apiKey?: string
  supportsImage: boolean
  modelSize?: string
  modelFile?: string
  downloadUrl?: string
  downloadProgress?: number
  downloadStatus?: 'not_installed' | 'downloading' | 'installed' | 'error'
  isActive?: boolean
}

export interface AIAgent {
  id: string
  name: string
  icon: string
  description: string
  systemPrompt: string
  isBuiltin?: boolean
}

export interface AIChatMessage {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: number
  streaming?: boolean
  images?: string[]
}

export interface AISettings {
  provider: 'local' | 'ollama' | 'cloud'
  localModelPath: string
  localModelName: string
  ollamaUrl: string
  ollamaModel: string
  cloudApiKey: string
  cloudApiUrl: string
  cloudModel: string
  enableGPU: boolean
  gpuLayers: number
  maxTokens: number
  temperature: number
  topP: number
  contextSize: number
  enableCodeCompletion: boolean
  completionDelay: number
  enableChatHistory: boolean
  systemPrompt: string
}

const STORAGE_KEY = 'yigou-ide-ai-store'

const HF_MIRROR = 'https://hf-mirror.com'
const HF_FALLBACK_MIRRORS = [
  'https://hf-mirror.com',
  'https://huggingface.co',
]

const defaultBuiltinModels: AIModel[] = [
  {
    id: 'qwen2.5-coder-1.5b',
    name: 'Qwen2.5-Coder 1.5B (内置)',
    provider: 'local',
    supportsImage: false,
    modelSize: '0.9 GB',
    modelFile: 'qwen2.5-coder-1.5b-q4_k_m.gguf',
    downloadUrl: 'https://hf-mirror.com/Qwen/Qwen2.5-Coder-1.5B-GGUF/resolve/main/qwen2.5-coder-1.5b-q4_k_m.gguf',
    downloadStatus: 'not_installed',
    downloadProgress: 0,
  },
  {
    id: 'qwen2.5-coder-7b',
    name: 'Qwen2.5-Coder 7B',
    provider: 'local',
    supportsImage: false,
    modelSize: '4.7 GB',
    modelFile: 'qwen2.5-coder-7b-q4_k_m.gguf',
    downloadUrl: 'https://hf-mirror.com/Qwen/Qwen2.5-Coder-7B-GGUF/resolve/main/qwen2.5-coder-7b-q4_k_m.gguf',
    downloadStatus: 'not_installed',
    downloadProgress: 0,
  },
  {
    id: 'codellama-7b',
    name: 'CodeLlama 7B',
    provider: 'local',
    supportsImage: false,
    modelSize: '4.4 GB',
    modelFile: 'codellama-7b-q4_k_m.gguf',
    downloadUrl: 'https://hf-mirror.com/TheBloke/CodeLlama-7B-GGUF/resolve/main/codellama-7b.Q4_K_M.gguf',
    downloadStatus: 'not_installed',
    downloadProgress: 0,
  },
  {
    id: 'deepseek-coder-1.3b',
    name: 'DeepSeek-Coder 1.3B',
    provider: 'local',
    supportsImage: false,
    modelSize: '0.8 GB',
    modelFile: 'deepseek-coder-1.3b-q4_k_m.gguf',
    downloadUrl: 'https://hf-mirror.com/TheBloke/deepseek-coder-1.3b-GGUF/resolve/main/deepseek-coder-1.3b.Q4_K_M.gguf',
    downloadStatus: 'not_installed',
    downloadProgress: 0,
  },
  {
    id: 'starcoder2-3b',
    name: 'StarCoder2 3B',
    provider: 'local',
    supportsImage: false,
    modelSize: '2.1 GB',
    modelFile: 'starcoder2-3b-q4_k_m.gguf',
    downloadUrl: 'https://hf-mirror.com/bigcode/starcoder2-3b-GGUF/resolve/main/starcoder2-3b-Q4_K_M.gguf',
    downloadStatus: 'not_installed',
    downloadProgress: 0,
  },
]

const defaultCloudModels: AIModel[] = [
  {
    id: 'gpt-4o-mini',
    name: 'GPT-4o Mini',
    provider: 'cloud',
    apiUrl: 'https://api.openai.com/v1',
    supportsImage: true,
    isActive: false,
  },
  {
    id: 'gpt-4o',
    name: 'GPT-4o',
    provider: 'cloud',
    apiUrl: 'https://api.openai.com/v1',
    supportsImage: true,
    isActive: false,
  },
  {
    id: 'claude-3.5-sonnet',
    name: 'Claude 3.5 Sonnet',
    provider: 'cloud',
    apiUrl: 'https://api.anthropic.com/v1',
    supportsImage: true,
    isActive: false,
  },
  {
    id: 'claude-3-opus',
    name: 'Claude 3 Opus',
    provider: 'cloud',
    apiUrl: 'https://api.anthropic.com/v1',
    supportsImage: true,
    isActive: false,
  },
  {
    id: 'deepseek-v3',
    name: 'DeepSeek-V3',
    provider: 'cloud',
    apiUrl: 'https://api.deepseek.com/v1',
    supportsImage: false,
    isActive: false,
  },
  {
    id: 'glm-4-plus',
    name: 'GLM-4 Plus',
    provider: 'cloud',
    apiUrl: 'https://open.bigmodel.cn/api/paas/v4',
    supportsImage: true,
    isActive: false,
  },
  {
    id: 'moonshot-v1',
    name: 'Moonshot v1',
    provider: 'cloud',
    apiUrl: 'https://api.moonshot.cn/v1',
    supportsImage: false,
    isActive: false,
  },
]

const defaultOllamaModels: AIModel[] = [
  { id: 'qwen2.5-coder:1.5b', name: 'Qwen2.5-Coder 1.5B (Ollama)', provider: 'ollama', apiUrl: 'http://localhost:11434', supportsImage: false },
  { id: 'qwen2.5-coder:7b', name: 'Qwen2.5-Coder 7B (Ollama)', provider: 'ollama', apiUrl: 'http://localhost:11434', supportsImage: false },
  { id: 'codellama:7b', name: 'CodeLlama 7B (Ollama)', provider: 'ollama', apiUrl: 'http://localhost:11434', supportsImage: false },
  { id: 'deepseek-coder:6.7b', name: 'DeepSeek-Coder 6.7B (Ollama)', provider: 'ollama', apiUrl: 'http://localhost:11434', supportsImage: false },
]

const defaultAgents: AIAgent[] = [
  {
    id: 'developer',
    name: '编程专家',
    icon: 'Cpu',
    description: '精通 Go 语言和中文函数体系，帮助你编写高质量代码',
    systemPrompt: '你是一位资深的 Go 语言编程专家，精通原生 Go 语言和中文函数映射体系。你可以帮助用户编写代码、优化逻辑、排查问题，并支持中英文函数的双向转换。请提供完整可运行的代码示例，并附上详细解释。',
    isBuiltin: true,
  },
  {
    id: 'architect',
    name: '架构师',
    icon: 'Files',
    description: '帮助你进行系统架构设计、模块划分和技术选型',
    systemPrompt: '你是一位经验丰富的系统架构师，专注于 Go 语言项目。你可以帮助用户进行系统架构设计、模块划分、接口定义和技术选型，确保系统的可扩展性、可维护性和高性能。',
    isBuiltin: true,
  },
  {
    id: 'debugger',
    name: '排错专家',
    icon: 'Bug',
    description: '快速定位 BUG，分析错误原因，提供修复方案',
    systemPrompt: '你是一位经验丰富的调试专家，擅长分析和解决代码问题。你可以帮助用户快速定位BUG、分析堆栈跟踪、检查内存泄漏和竞态条件，并提供完整的修复方案。',
    isBuiltin: true,
  },
  {
    id: 'translator',
    name: '翻译官',
    icon: 'Switch',
    description: '中英文函数互转、代码注释翻译、文档生成',
    systemPrompt: '你是一位专业的技术翻译，精通中英文技术文档翻译和代码转换。你可以帮助用户将中文函数名转为英文，或英文转为中文，生成代码注释和API文档。',
    isBuiltin: true,
  },
  {
    id: 'reviewer',
    name: '代码审查',
    icon: 'View',
    description: '审查代码质量、安全性和性能问题，提供改进建议',
    systemPrompt: '你是一位资深的代码审查专家，专注于 Go 语言。你可以帮助用户审查代码质量、安全漏洞、性能瓶颈和代码风格问题，并提供具体的改进建议和最佳实践指导。',
    isBuiltin: true,
  },
]

export interface ContextCompressionState {
  isCompressed: boolean
  originalMessageCount: number
  compressedMessageCount: number
  summaryText: string
  lastCompressionTime: number
  estimatedTokensSaved: number
}

interface AIStoreState {
  models: AIModel[]
  agents: AIAgent[]
  currentModelId: string
  currentAgentId: string
  chatHistory: Record<string, AIChatMessage[]>
  contextCompression: Record<string, ContextCompressionState>
}

function loadState(): AIStoreState {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const parsed = JSON.parse(saved)
      // Merge saved models with defaults, preserving download status
      const mergedModels = [...defaultBuiltinModels]
      const mergedOllamaModels = [...defaultOllamaModels]
      const mergedCloudModels = [...defaultCloudModels]
      if (parsed.models) {
        parsed.models.forEach((savedModel: AIModel) => {
          const idx = mergedModels.findIndex(m => m.id === savedModel.id)
          if (idx >= 0) {
            mergedModels[idx] = { ...mergedModels[idx], ...savedModel }
          }
          const oIdx = mergedOllamaModels.findIndex(m => m.id === savedModel.id)
          if (oIdx >= 0) {
            mergedOllamaModels[oIdx] = { ...mergedOllamaModels[oIdx], ...savedModel }
          }
          const cIdx = mergedCloudModels.findIndex(m => m.id === savedModel.id)
          if (cIdx >= 0) {
            mergedCloudModels[cIdx] = { ...mergedCloudModels[cIdx], ...savedModel }
          }
        })
      }
      return {
        models: [...mergedModels, ...mergedOllamaModels, ...mergedCloudModels],
        agents: parsed.agents || [...defaultAgents],
        currentModelId: parsed.currentModelId || mergedModels[0].id,
        currentAgentId: parsed.currentAgentId || defaultAgents[0].id,
        chatHistory: parsed.chatHistory || {},
        contextCompression: parsed.contextCompression || {},
      }
    }
  } catch (e) {
    console.warn('加载 AI 设置失败', e)
  }
  return {
    models: [...defaultBuiltinModels, ...defaultOllamaModels, ...defaultCloudModels],
    agents: [...defaultAgents],
    currentModelId: defaultBuiltinModels[0].id,
    currentAgentId: defaultAgents[0].id,
    chatHistory: {},
    contextCompression: {},
  }
}

function saveState(state: AIStoreState) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
  } catch (e) {
    console.warn('保存 AI 设置失败', e)
  }
}

export function useAIStore() {
  const state = reactive<AIStoreState>(loadState())

  const allModels = computed(() => state.models)
  const availableModels = computed(() => state.models.filter(m => {
    if (m.provider === 'local') return m.downloadStatus === 'installed'
    if (m.provider === 'cloud') return m.isActive === true
    return true
  }))
  const localModels = computed(() => state.models.filter(m => m.provider === 'local'))
  const ollamaModels = computed(() => state.models.filter(m => m.provider === 'ollama'))
  const cloudModels = computed(() => state.models.filter(m => m.provider === 'cloud'))
  const installedLocalModels = computed(() => state.models.filter(m => m.provider === 'local' && m.downloadStatus === 'installed'))
  const activeCloudModels = computed(() => state.models.filter(m => m.provider === 'cloud' && m.isActive))

  const currentModel = computed(() => state.models.find(m => m.id === state.currentModelId) || state.models[0])
  const currentAgent = computed(() => state.agents.find(a => a.id === state.currentAgentId) || state.agents[0])

  const currentChatHistory = computed(() => state.chatHistory[state.currentAgentId] || [])

  function selectModel(modelId: string) {
    state.currentModelId = modelId
    saveState({ ...state })
  }

  function selectAgent(agentId: string) {
    state.currentAgentId = agentId
    saveState({ ...state })
  }

  function addAgent(agent: AIAgent) {
    state.agents.push(agent)
    saveState({ ...state })
  }

  function updateAgent(agentId: string, updates: Partial<AIAgent>) {
    const idx = state.agents.findIndex(a => a.id === agentId)
    if (idx >= 0) {
      state.agents[idx] = { ...state.agents[idx], ...updates }
      saveState({ ...state })
    }
  }

  function removeAgent(agentId: string) {
    state.agents = state.agents.filter(a => a.id !== agentId)
    if (state.currentAgentId === agentId) {
      state.currentAgentId = state.agents[0]?.id || ''
    }
    saveState({ ...state })
  }

  function addCustomModel(model: AIModel) {
    state.models.push({ ...model, downloadStatus: model.provider === 'local' ? 'not_installed' : undefined })
    saveState({ ...state })
  }

  function updateCustomModel(modelId: string, updates: Partial<AIModel>) {
    const idx = state.models.findIndex(m => m.id === modelId)
    if (idx >= 0) {
      state.models[idx] = { ...state.models[idx], ...updates }
      saveState({ ...state })
    }
  }

  function removeCustomModel(modelId: string) {
    state.models = state.models.filter(m => m.id !== modelId)
    if (state.currentModelId === modelId) {
      state.currentModelId = state.models[0]?.id || ''
    }
    saveState({ ...state })
  }

  function getFallbackUrl(url: string): string {
    if (url.startsWith(HF_MIRROR)) {
      return url.replace(HF_MIRROR, 'https://huggingface.co')
    }
    return url
  }

  async function downloadModel(modelId: string): Promise<void> {
    const model = state.models.find(m => m.id === modelId)
    if (!model || !model.downloadUrl) return

    model.downloadStatus = 'downloading'
    model.downloadProgress = 0
    saveState({ ...state })

    try {
      try {
        await DownloadModelFile(modelId, model.downloadUrl)
      } catch {
        const fallbackUrl = getFallbackUrl(model.downloadUrl)
        if (fallbackUrl !== model.downloadUrl) {
          console.log('HF mirror failed, trying original source:', fallbackUrl)
          await DownloadModelFile(modelId, fallbackUrl)
        } else {
          throw new Error('下载失败')
        }
      }
      model.downloadStatus = 'installed'
      model.downloadProgress = 100
      saveState({ ...state })
    } catch (e: any) {
      model.downloadStatus = 'error'
      model.downloadProgress = 0
      saveState({ ...state })
      throw e
    }
  }

  function setModelInstalled(modelId: string, installed: boolean) {
    const model = state.models.find(m => m.id === modelId)
    if (model) {
      model.downloadStatus = installed ? 'installed' : 'not_installed'
      model.downloadProgress = installed ? 100 : 0
      saveState({ ...state })
    }
  }

  async function checkInstalledModels() {
    for (const model of state.models) {
      if (model.provider !== 'local') continue
      try {
        const status = await GetModelDownloadStatus(model.id)
        if (status && status.status === 'installed') {
          model.downloadStatus = 'installed'
          model.downloadProgress = 100
          saveState({ ...state })
        }
      } catch {
        // offline or error, keep current status
      }
    }
  }

  function updateModelProgress(modelId: string, progress: number, status: string) {
    const model = state.models.find(m => m.id === modelId)
    if (model) {
      model.downloadProgress = progress
      if (status === 'installed' || status === 'error') {
        model.downloadStatus = status as 'installed' | 'error'
      }
      saveState({ ...state })
    }
  }

  function addChatMessage(message: AIChatMessage) {
    if (!state.chatHistory[state.currentAgentId]) {
      state.chatHistory[state.currentAgentId] = []
    }
    state.chatHistory[state.currentAgentId].push(message)
    saveState({ ...state })
  }

  function clearChatHistory() {
    state.chatHistory[state.currentAgentId] = []
    saveState({ ...state })
  }

  function deleteChatMessage(messageId: string) {
    const history = state.chatHistory[state.currentAgentId]
    if (history) {
      state.chatHistory[state.currentAgentId] = history.filter(m => m.id !== messageId)
      saveState({ ...state })
    }
  }

  function estimateTokens(text: string): number {
    let count = 0
    for (let i = 0; i < text.length; i++) {
      const ch = text.charCodeAt(i)
      if (ch >= 0x4e00 && ch <= 0x9fff) {
        count += 2
      } else if (ch >= 0x3000 && ch <= 0x303f) {
        count += 1
      } else if (ch > 127) {
        count += 1
      } else {
        count += text[i] === ' ' ? 0.25 : 0.25
      }
    }
    return Math.ceil(count)
  }

  function estimateChatTokens(agentId?: string): number {
    const id = agentId || state.currentAgentId
    const history = state.chatHistory[id] || []
    let total = 0
    for (const msg of history) {
      total += estimateTokens(msg.content)
    }
    const compression = state.contextCompression[id]
    if (compression?.summaryText) {
      total += estimateTokens(compression.summaryText)
    }
    return total
  }

  const COMPRESSION_THRESHOLD = 6000
  const KEEP_RECENT_COUNT = 8

  function compressContext(agentId?: string): boolean {
    const id = agentId || state.currentAgentId
    const history = state.chatHistory[id]
    if (!history || history.length <= KEEP_RECENT_COUNT + 4) return false

    const totalTokens = estimateChatTokens(id)
    if (totalTokens < COMPRESSION_THRESHOLD) return false

    const toCompress = history.slice(0, history.length - KEEP_RECENT_COUNT)
    const keepRecent = history.slice(-KEEP_RECENT_COUNT)

    const summaryParts: string[] = ['[对话历史摘要]']
    for (const msg of toCompress) {
      const roleLabel = msg.role === 'user' ? '用户' : '助手'
      const truncated = msg.content.length > 200 ? msg.content.substring(0, 200) + '...' : msg.content
      summaryParts.push(`${roleLabel}: ${truncated}`)
    }
    const summaryText = summaryParts.join('\n')

    const originalTokens = toCompress.reduce((sum, m) => sum + estimateTokens(m.content), 0)
    const summaryTokens = estimateTokens(summaryText)
    const tokensSaved = originalTokens - summaryTokens

    state.chatHistory[id] = keepRecent
    state.contextCompression[id] = {
      isCompressed: true,
      originalMessageCount: history.length,
      compressedMessageCount: toCompress.length,
      summaryText,
      lastCompressionTime: Date.now(),
      estimatedTokensSaved: Math.max(0, tokensSaved),
    }
    saveState({ ...state })
    return true
  }

  function getCompressionState(agentId?: string): ContextCompressionState | null {
    const id = agentId || state.currentAgentId
    return state.contextCompression[id] || null
  }

  function clearCompression(agentId?: string) {
    const id = agentId || state.currentAgentId
    delete state.contextCompression[id]
    saveState({ ...state })
  }

  const contextCompressionState = computed(() => state.contextCompression[state.currentAgentId] || null)
  const estimatedContextTokens = computed(() => estimateChatTokens())

  return {
    state,
    agents: computed(() => state.agents),
    allModels,
    availableModels,
    localModels,
    ollamaModels,
    cloudModels,
    installedLocalModels,
    activeCloudModels,
    currentModel,
    currentAgent,
    currentChatHistory,
    selectModel,
    selectAgent,
    addAgent,
    updateAgent,
    removeAgent,
    addCustomModel,
    updateCustomModel,
    removeCustomModel,
    downloadModel,
    setModelInstalled,
    checkInstalledModels,
    updateModelProgress,
    addChatMessage,
    clearChatHistory,
    deleteChatMessage,
    estimateTokens,
    estimateChatTokens,
    compressContext,
    getCompressionState,
    clearCompression,
    contextCompressionState,
    estimatedContextTokens,
  }
}