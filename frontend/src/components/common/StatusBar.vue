<script lang="ts" setup>
import { ref, computed } from 'vue'

interface StatusBarProps {
  projectName?: string
  isRunning?: boolean
  platform?: string
  goVersion?: string
  lineCount?: number
  charCount?: number
}

const props = withDefaults(defineProps<StatusBarProps>(), {
  projectName: '未命名项目',
  isRunning: false,
  platform: 'windows/amd64',
  goVersion: '-',
  lineCount: 0,
  charCount: 0,
})

const emit = defineEmits<{
  (e: 'toggle-left-panel'): void
  (e: 'toggle-right-panel'): void
  (e: 'toggle-bottom-panel'): void
}>()

const currentTime = ref('')
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
}
setInterval(updateTime, 1000)
updateTime()

const statusText = computed(() => props.isRunning ? '● 运行中' : '就绪')
const statusClass = computed(() => props.isRunning ? 'running' : 'ready')
</script>

<template>
  <div class="status-bar">
    <div class="status-left">
      <span class="status-item status-text" :class="statusClass">
        {{ statusText }}
      </span>
      <span class="status-separator">|</span>

      <!-- 面板开关按钮（带文字标签，更明显） -->
      <div class="panel-toggles">
        <button class="panel-toggle-btn" title="切换左边栏" @click="emit('toggle-left-panel')">
          <el-icon :size="12"><DArrowRight /></el-icon>
          <span>左</span>
        </button>
        <button class="panel-toggle-btn" title="切换右边AI栏" @click="emit('toggle-right-panel')">
          <el-icon :size="12"><ChatDotRound /></el-icon>
          <span>AI</span>
        </button>
        <button class="panel-toggle-btn" title="切换底部栏" @click="emit('toggle-bottom-panel')">
          <el-icon :size="12"><ArrowUp /></el-icon>
          <span>底</span>
        </button>
      </div>

      <span class="status-separator">|</span>
      <span class="status-item" title="项目名称">
        <el-icon><Folder /></el-icon> {{ projectName }}
      </span>
      <span class="status-separator">|</span>
      <span class="status-item" title="目标平台">
        <el-icon><Monitor /></el-icon> {{ platform }}
      </span>
    </div>

    <div class="status-center">
      <span class="status-item shortcut-hint">F5 运行 · Ctrl+S 保存</span>
    </div>

    <div class="status-right">
      <span v-if="lineCount > 0" class="status-item" title="代码统计">
        行: {{ lineCount }} | 字符: {{ charCount }}
      </span>
      <span class="status-separator" v-if="lineCount > 0">|</span>
      <span class="status-item">UTF-8</span>
      <span class="status-separator">|</span>
      <span class="status-item" title="当前时间">
        <el-icon><Clock /></el-icon> {{ currentTime }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.status-bar {
  height: 26px;
  background-color: var(--statusbar-bg);
  color: var(--statusbar-text);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 10px;
  font-size: 12px;
  flex-shrink: 0;
  user-select: none;
}

.status-left,
.status-center,
.status-right {
  display: flex;
  align-items: center;
  gap: 2px;
}

.status-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.status-item {
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 0 6px;
  white-space: nowrap;
}

.status-item .el-icon {
  font-size: 13px;
}

.status-text.ready { color: var(--color-success); }
.status-text.running {
  color: var(--color-warning);
  animation: pulse 1.5s infinite;
}
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }

.status-separator {
  color: var(--text-secondary);
  margin: 0 2px;
}

.shortcut-hint {
  color: var(--text-regular);
  font-size: 11px;
  opacity: 0.7;
}
.shortcut-hint:hover { opacity: 1; }

/* ====== 面板开关按钮组 ====== */
.panel-toggles {
  display: flex;
  align-items: center;
  gap: 2px;
  background: rgba(255,255,255,0.08);
  border-radius: 4px;
  padding: 1px;
}

.panel-toggle-btn {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 2px 6px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 11px;
  color: var(--statusbar-text);
  background: transparent;
  transition: all 0.15s ease;
  white-space: nowrap;
  line-height: 1.4;
}

.panel-toggle-btn:hover {
  background: rgba(255,255,255,0.15);
  color: #fff;
}

.panel-toggle-btn:active {
  background: rgba(255,255,255,0.22);
}

.panel-toggle-btn .el-icon {
  font-size: 12px;
}
</style>