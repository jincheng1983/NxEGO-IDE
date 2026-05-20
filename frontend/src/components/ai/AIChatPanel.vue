<script lang="ts" setup>
import { ref, nextTick, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useAIStore } from '../../stores/ai'
import { ReadEgoMemory, WriteEgoMemory, GetEgoMemorySummary } from '../../../wailsjs/go/main/App'

interface ChatMessage {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  images?: string[]
  timestamp: number
  streaming?: boolean
}

const {
  agents: storeAgents,
  availableModels: storeModels,
  currentAgent: storeCurrentAgent,
  currentModel: storeCurrentModel,
  selectAgent: storeSelectAgent,
  selectModel: storeSelectModel,
  compressContext,
  contextCompressionState: storeCompressionState,
  estimatedContextTokens,
  getCompressionState,
  clearCompression,
} = useAIStore()

const messages = ref<ChatMessage[]>([])
const inputText = ref('')
const isStreaming = ref(false)
const chatContainerRef = ref<HTMLElement | null>(null)

const showAgentPicker = ref(false)
const showModelPicker = ref(false)
const showCompressionSummary = ref(false)

const selectedImages = ref<string[]>([])

const showMemoryPanel = ref(false)
const memoryContent = ref('')
const memorySummary = ref('')
const isMemoryLoading = ref(false)
const isMemorySaving = ref(false)
const currentProjectPath = ref('')

const props = defineProps<{
  projectPath?: string
  projectName?: string
}>()

watch(() => props.projectPath, (newPath) => {
  if (newPath && newPath !== currentProjectPath.value) {
    currentProjectPath.value = newPath
    loadMemory()
  }
}, { immediate: true })

async function loadMemory() {
  isMemoryLoading.value = true
  try {
    const content = await ReadEgoMemory(currentProjectPath.value)
    memoryContent.value = content || ''
    const summary = await GetEgoMemorySummary(currentProjectPath.value)
    memorySummary.value = summary || '暂无项目记忆'
  } catch {
    memoryContent.value = ''
    memorySummary.value = '暂无项目记忆'
  } finally {
    isMemoryLoading.value = false
  }
}

async function saveMemory() {
  isMemorySaving.value = true
  try {
    await WriteEgoMemory(currentProjectPath.value, memoryContent.value)
    ElMessage.success('项目记忆已保存')
    showMemoryPanel.value = false
    await loadMemory()
  } catch (e: any) {
    ElMessage.error('保存失败: ' + (e?.message || e))
  } finally {
    isMemorySaving.value = false
  }
}

function toggleMemoryPanel() {
  showMemoryPanel.value = !showMemoryPanel.value
  if (showMemoryPanel.value && currentProjectPath.value) {
    loadMemory()
  }
}

onMounted(() => {
  messages.value.push({
    id: 'welcome',
    role: 'assistant',
    content: `您正在与 ${storeCurrentAgent.value.name} 聊天`,
    timestamp: Date.now(),
  })
})

function scrollToBottom() {
  nextTick(() => {
    if (chatContainerRef.value) {
      chatContainerRef.value.scrollTop = chatContainerRef.value.scrollHeight
    }
  })
}

function generateId(): string {
  return `msg_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`
}

function onSelectAgent(agent: typeof storeAgents.value[number]) {
  storeSelectAgent(agent.id)
  showAgentPicker.value = false
}

function onSelectModel(model: typeof storeModels.value[number]) {
  storeSelectModel(model.id)
  showModelPicker.value = false
}

async function sendMessage() {
  const text = inputText.value.trim()
  if ((!text && selectedImages.value.length === 0) || isStreaming.value) return

  messages.value.push({
    id: generateId(), role: 'user', content: text,
    images: selectedImages.value.length > 0 ? [...selectedImages.value] : undefined,
    timestamp: Date.now(),
  })
  inputText.value = ''
  selectedImages.value = []
  autoResize()
  scrollToBottom()

  const assistantMsg: ChatMessage = { id: generateId(), role: 'assistant', content: '', timestamp: Date.now(), streaming: true }
  messages.value.push(assistantMsg)
  isStreaming.value = true
  scrollToBottom()

  try {
    await simulateAIResponse(assistantMsg, text)
  } catch (e: any) {
    assistantMsg.content = `抱歉，发生了错误：${e?.message || e}`
    assistantMsg.streaming = false
    isStreaming.value = false
  }
}

async function simulateAIResponse(msg: ChatMessage, userInput: string) {
  let fullResponse = ''
  const lowerInput = userInput.toLowerCase()
  const prefix = `[${storeCurrentAgent.value.name}]`

  if (lowerInput.includes('你好') || lowerInput.includes('hello') || lowerInput.includes('hi')) {
    fullResponse = `${prefix} 你好！我是 ${storeCurrentAgent.value.name}，有什么可以帮助你的吗？`
  } else if (lowerInput.includes('帮助') || lowerInput.includes('help')) {
    fullResponse = `${prefix} 作为${storeCurrentAgent.value.name}：\n• ${storeCurrentAgent.value.description}\n\n直接告诉我你的需求即可！`
  } else if (lowerInput.includes('代码') || lowerInput.includes('code')) {
    fullResponse = `${prefix} 好的，我来帮你处理代码问题。\n\n需要什么功能？中文函数还是原生 Go？有具体代码片段吗？`
  } else if (lowerInput.includes('函数') || lowerInput.includes('function')) {
    fullResponse = `${prefix} 易狗 NxEGO 中文函数示例：\n• \`输出.打印("Hello")\` → \`yigou.Print("Hello")\`\n• \`文件.读取文本("path")\` → \`yigou.ReadFile("path")\``
  } else if (lowerInput.includes('bug') || lowerInput.includes('错误') || lowerInput.includes('修复')) {
    fullResponse = `${prefix} 我来帮你排查。请提供：\n1. 错误信息或现象\n2. 相关代码片段\n3. 期望的正确行为`
  } else {
    fullResponse = `${prefix} 收到：「${userInput}」\n\n作为易狗 NxEGO AI（${storeCurrentModel.value.name}），我可以帮你处理 Go 代码编写、中文函数转换等问题。\n\n请告诉我更多细节！`
  }

  for (let i = 0; i < fullResponse.length; i++) {
    await new Promise(resolve => setTimeout(resolve, 10 + Math.random() * 15))
    msg.content = fullResponse.substring(0, i + 1)
    scrollToBottom()
  }
  msg.streaming = false
  isStreaming.value = false
  compressContext()
}

function clearMessages() {
  messages.value = [{ id: 'welcome', role: 'assistant', content: `您正在与 ${storeCurrentAgent.value.name} 聊天`, timestamp: Date.now() }]
  clearCompression()
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

const textareaRef = ref<HTMLTextAreaElement | null>(null)
const MIN_HEIGHT = 36
const MAX_HEIGHT = 200

function autoResize() {
  nextTick(() => {
    const el = textareaRef.value
    if (!el) return
    el.style.height = 'auto'
    const scrollHeight = el.scrollHeight
    el.style.height = Math.min(Math.max(scrollHeight, MIN_HEIGHT), MAX_HEIGHT) + 'px'
  })
}

function formatTime(ts: number): string {
  const d = new Date(ts)
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

function handleImageUpload(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    if (!storeCurrentModel.value.supportsImage) {
      ElMessage.warning(`当前模型不支持图片`)
      return
    }
    const reader = new FileReader()
    reader.onload = () => { selectedImages.value.push(reader.result as string) }
    reader.readAsDataURL(target.files[0])
    target.value = ''
  }
}

function removeImage(index: number) {
  selectedImages.value.splice(index, 1)
}

function providerLabel(provider: string): string {
  if (provider === 'local') return '本地'
  if (provider === 'ollama') return 'Ollama'
  return '云端'
}
</script>

<template>
  <div class="trae-chat">
    <div class="chat-topbar">
      <div class="topbar-agent" @click.stop="showAgentPicker = !showAgentPicker">
        <span class="topbar-at">@</span>
        <span class="topbar-name">{{ storeCurrentAgent.name }}</span>
        <el-icon :size="10"><ArrowDown /></el-icon>

        <div v-if="showAgentPicker" class="picker-dropdown" @click.stop>
          <div class="picker-header">
            <span>选择智能体</span>
          </div>
          <div
            v-for="agent in storeAgents"
            :key="agent.id"
            class="picker-item"
            :class="{ active: storeCurrentAgent.id === agent.id }"
            @click="onSelectAgent(agent)"
          >
            <div class="picker-item-icon">
              <el-icon :size="16"><component v-if="agent.icon" :is="agent.icon" /></el-icon>
            </div>
            <div class="picker-info">
              <span class="picker-name">{{ agent.name }}</span>
              <span class="picker-desc">{{ agent.description }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showMemoryPanel" class="memory-panel">
      <div class="memory-header">
        <el-icon :size="16"><Document /></el-icon>
        <span class="memory-title">.ego.md 项目记忆</span>
        <span class="memory-hint">AI 会参考此文件理解你的项目</span>
        <button
          class="memory-save-btn"
          :disabled="isMemorySaving"
          @click="saveMemory"
        >
          {{ isMemorySaving ? '保存中...' : '保存' }}
        </button>
        <button class="memory-close-btn" @click="showMemoryPanel = false">
          <el-icon :size="14"><Close /></el-icon>
        </button>
      </div>
      <textarea
        v-model="memoryContent"
        class="memory-editor"
        placeholder="在此编辑项目记忆，AI 会自动读取...&#10;&#10;你可记录：&#10;- 项目技术栈和依赖&#10;- 重要的代码约定&#10;- 当前开发阶段&#10;- 给 AI 的特定指令"
        :disabled="isMemoryLoading"
      />
    </div>

    <div ref="chatContainerRef" class="chat-messages">
      <div v-for="msg in messages" :key="msg.id" :class="['chat-msg', `msg-${msg.role}`]">
        <div v-if="msg.images && msg.images.length > 0" class="msg-images">
          <img v-for="(img, idx) in msg.images" :key="idx" :src="img" class="msg-img" />
        </div>
        <pre class="msg-text">{{ msg.content }}<span v-if="msg.streaming" class="cursor">|</span></pre>
        <span v-if="msg.role !== 'system'" class="msg-time">{{ formatTime(msg.timestamp) }}</span>
      </div>
    </div>

    <div v-if="storeCompressionState?.isCompressed" class="compression-bar">
      <div class="compression-info">
        <el-icon :size="14"><InfoFilled /></el-icon>
        <span>已压缩 {{ storeCompressionState.compressedMessageCount }} 条历史消息，节省约 {{ storeCompressionState.estimatedTokensSaved }} tokens</span>
        <span class="compression-tokens">当前 ~{{ estimatedContextTokens }} tokens</span>
      </div>
      <div class="compression-actions">
        <button class="compression-btn" @click="showCompressionSummary = !showCompressionSummary">
          {{ showCompressionSummary ? '收起' : '查看摘要' }}
        </button>
        <button class="compression-btn" @click="clearCompression(); messages = [{ id: 'welcome', role: 'assistant', content: `您正在与 ${storeCurrentAgent.name} 聊天`, timestamp: Date.now() }]">
          清除
        </button>
      </div>
      <div v-if="showCompressionSummary" class="compression-summary">
        <pre>{{ storeCompressionState.summaryText }}</pre>
      </div>
    </div>

    <div class="chat-toolbar">
      <div class="toolbar-row">
        <div class="toolbar-left">
          <button class="toolbar-icon-btn" title="项目记忆 (.ego.md)" @click="toggleMemoryPanel">
            <el-icon :size="16"><Notebook /></el-icon>
          </button>
          <button class="toolbar-icon-btn" title="清空对话" @click="clearMessages">
            <el-icon :size="16"><Delete /></el-icon>
          </button>
          <label v-if="storeCurrentModel.supportsImage" class="toolbar-icon-btn" title="上传图片">
            <el-icon :size="16"><Picture /></el-icon>
            <input type="file" accept="image/*" hidden @change="handleImageUpload" />
          </label>
        </div>

        <div class="toolbar-center">
          <div class="model-selector" @click.stop="showModelPicker = !showModelPicker">
            <span class="model-name">{{ storeCurrentModel.name }}</span>
            <el-icon :size="10"><ArrowDown /></el-icon>

            <div v-if="showModelPicker" class="picker-dropdown model-picker" @click.stop>
              <div class="picker-header">
                <span>选择模型</span>
              </div>
              <div
                v-for="model in storeModels"
                :key="model.id"
                class="picker-item"
                :class="{ active: storeCurrentModel.id === model.id }"
                @click="onSelectModel(model)"
              >
                <el-icon :size="16"><Cpu /></el-icon>
                <div class="picker-info">
                  <span class="picker-name">{{ model.name }}</span>
                  <span class="picker-meta">
                    {{ providerLabel(model.provider) }}
                    <el-tag v-if="model.supportsImage" size="small" type="success" effect="plain">图片</el-tag>
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="toolbar-right">
          <button
            class="send-btn"
            :disabled="(!inputText.trim() && selectedImages.length === 0) || isStreaming"
            @click="sendMessage"
          >
            <el-icon v-if="!isStreaming" :size="16"><Top /></el-icon>
            <el-icon v-else :size="16" class="spin"><Loading /></el-icon>
          </button>
        </div>
      </div>

      <div v-if="selectedImages.length > 0" class="image-strip">
        <div v-for="(img, idx) in selectedImages" :key="idx" class="thumb">
          <img :src="img" />
          <span class="thumb-x" @click="removeImage(idx)">×</span>
        </div>
      </div>

      <textarea
        ref="textareaRef"
        v-model="inputText"
        placeholder="输入你的问题，Enter 发送，Shift+Enter 换行..."
        :disabled="isStreaming"
        @keydown="handleKeydown"
        @input="autoResize"
        class="chat-textarea"
        rows="1"
      />
    </div>
  </div>
</template>

<style scoped>
.trae-chat {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--panel-bg);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

/* ====== Top Bar ====== */
.chat-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-lighter);
  flex-shrink: 0;
}

.topbar-agent {
  display: flex;
  align-items: center;
  gap: 2px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background 0.15s;
  position: relative;
}

.topbar-agent:hover { background: rgba(255,255,255,0.08); }

.topbar-at {
  color: var(--text-placeholder);
  font-size: 14px;
  font-weight: 600;
}

.topbar-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

/* ====== Messages ====== */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-msg {
  animation: fadeUp 0.2s ease;
}
@keyframes fadeUp { from { opacity: 0; transform: translateY(6px); } to { opacity: 1; transform: translateY(0); } }

.msg-user .msg-text {
  background: var(--color-primary);
  color: #fff;
  padding: 10px 14px;
  border-radius: 12px;
  border-bottom-right-radius: 4px;
  max-width: 85%;
  margin-left: auto;
  word-break: break-word;
}

.msg-assistant .msg-text {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  padding: 10px 14px;
  border-radius: 12px;
  border-bottom-left-radius: 4px;
  max-width: 90%;
  word-break: break-word;
}

.msg-text {
  margin: 0;
  white-space: pre-wrap;
  font-family: inherit;
  font-size: 13.5px;
  line-height: 1.65;
}

.msg-images { margin-bottom: 6px; }
.msg-img {
  max-width: 180px;
  max-height: 120px;
  border-radius: 8px;
  object-fit: cover;
  border: 1px solid var(--border-lighter);
}

.msg-time {
  display: block;
  font-size: 10px;
  color: var(--text-placeholder);
  margin-top: 4px;
  padding: 0 4px;
}

.msg-user .msg-time { text-align: right; }

.cursor { animation: blink 1s step-end infinite; color: var(--color-primary); font-weight: bold; }
@keyframes blink { 50% { opacity: 0; } }

/* ====== Bottom Toolbar ====== */
.chat-toolbar {
  border-top: 1px solid var(--border-lighter);
  padding: 8px 14px;
  background: var(--bg-secondary);
  flex-shrink: 0;
}

.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

.toolbar-left { display: flex; gap: 2px; }
.toolbar-center { flex: 1; display: flex; justify-content: center; }
.toolbar-right { display: flex; gap: 2px; }

.toolbar-icon-btn {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: var(--text-secondary);
  background: transparent;
  transition: all 0.15s;
}

.toolbar-icon-btn:hover { color: var(--text-primary); background: rgba(255,255,255,0.08); }

.model-selector {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: background 0.15s;
  position: relative;
}

.model-selector:hover { background: rgba(255,255,255,0.06); color: var(--text-primary); }

.model-name { max-width: 120px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.send-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  background: var(--color-primary);
  color: #fff;
  transition: all 0.15s;
}

.send-btn:hover:not(:disabled) { transform: scale(1.05); }
.send-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.chat-textarea {
  width: 100%;
  border: 1px solid var(--border-lighter);
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 14px;
  font-family: inherit;
  background: var(--bg-primary);
  color: var(--text-primary);
  resize: none;
  outline: none;
  min-height: 36px;
  max-height: 200px;
  overflow-y: auto;
  line-height: 1.5;
  transition: height 0.1s ease;
}

.chat-textarea:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.15);
}

.chat-textarea:disabled { opacity: 0.5; }

/* ====== Image Strip ====== */
.image-strip { display: flex; gap: 6px; margin-bottom: 6px; flex-wrap: wrap; }

.thumb {
  position: relative;
  width: 48px;
  height: 48px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid var(--border-lighter);
}

.thumb img { width: 100%; height: 100%; object-fit: cover; }

.thumb-x {
  position: absolute;
  top: 0;
  right: 0;
  width: 16px;
  height: 16px;
  background: rgba(0,0,0,0.6);
  color: #fff;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 0 6px 0 4px;
}

/* ====== Picker Dropdown ====== */
.picker-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 4px;
  min-width: 240px;
  max-width: 300px;
  max-height: 320px;
  overflow-y: auto;
  background: var(--bg-secondary);
  border: 1px solid var(--border-lighter);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.25);
  z-index: 100;
}

.model-picker {
  bottom: 100%;
  top: auto;
  margin-bottom: 4px;
  margin-top: 0;
}

.picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-lighter);
}

.picker-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  cursor: pointer;
  transition: background 0.1s;
  border-bottom: 1px solid rgba(255,255,255,0.03);
}

.picker-item:last-child { border-bottom: none; }
.picker-item:hover { background: rgba(255,255,255,0.05); }
.picker-item.active { background: rgba(64, 158, 255, 0.1); }

.picker-item-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: rgba(255,255,255,0.06);
  flex-shrink: 0;
}

.picker-info { display: flex; flex-direction: column; gap: 2px; overflow: hidden; }

.picker-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.picker-desc {
  font-size: 11px;
  color: var(--text-placeholder);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.picker-meta {
  font-size: 11px;
  color: var(--text-placeholder);
  display: flex;
  align-items: center;
  gap: 4px;
}

/* ====== Compression Bar ====== */
.compression-bar {
  border-top: 1px solid var(--border-lighter);
  background: rgba(245, 158, 11, 0.08);
  padding: 6px 14px;
  flex-shrink: 0;
}

.compression-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.compression-tokens {
  color: var(--text-placeholder);
  font-size: 10px;
  flex-shrink: 0;
}

.compression-actions {
  display: flex;
  gap: 6px;
}

.compression-btn {
  font-size: 11px;
  padding: 2px 8px;
  border: 1px solid var(--border-lighter);
  border-radius: 4px;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}

.compression-btn:hover {
  background: rgba(255,255,255,0.06);
  color: var(--text-primary);
}

.compression-summary {
  margin-top: 6px;
  padding: 8px 10px;
  background: var(--bg-primary);
  border: 1px solid var(--border-lighter);
  border-radius: 6px;
  max-height: 160px;
  overflow-y: auto;
}

.compression-summary pre {
  margin: 0;
  font-size: 11px;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  line-height: 1.5;
}

.memory-panel {
  border-top: 1px solid var(--border-lighter);
  border-bottom: 1px solid var(--border-lighter);
  background: var(--bg-secondary);
  flex-shrink: 0;
}

.memory-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-bottom: 1px solid var(--border-lighter);
}

.memory-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-primary);
}

.memory-hint {
  font-size: 10px;
  color: var(--text-placeholder);
  flex: 1;
}

.memory-save-btn {
  font-size: 11px;
  padding: 3px 10px;
  border: 1px solid var(--color-primary);
  border-radius: 4px;
  background: var(--color-primary);
  color: #fff;
  cursor: pointer;
  transition: all 0.15s;
}

.memory-save-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.memory-save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.memory-close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: var(--text-placeholder);
  cursor: pointer;
}

.memory-close-btn:hover {
  background: rgba(255,255,255,0.06);
  color: var(--text-primary);
}

.memory-editor {
  width: 100%;
  min-height: 160px;
  max-height: 260px;
  padding: 10px 12px;
  border: none;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 12px;
  font-family: 'Consolas', 'Courier New', monospace;
  line-height: 1.6;
  resize: vertical;
  outline: none;
}

.memory-editor::placeholder {
  color: var(--text-placeholder);
}

.memory-editor:focus {
  background: var(--bg-primary);
}
</style>