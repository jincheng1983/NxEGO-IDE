<script lang="ts" setup>
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { GetAllFunctions, GetCategories, GetFunctionsByCategory } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime/runtime'

interface ChineseFunction {
  chineseName: string
  goFunction: string
  description: string
  category: string
  params: string[]
  returnType: string
}

const activeTab = ref('console')
const consoleLogs = ref<string[]>([
  '[系统] 易狗 IDE v1.0.0 已启动',
  '[系统] 等待用户操作...',
])
const allFunctions = ref<ChineseFunction[]>([])
const categories = ref<string[]>([])
const selectedCategory = ref('')
const categoryFunctions = ref<ChineseFunction[]>([])
const selectedFunction = ref<ChineseFunction | null>(null)

const loadFunctions = async () => {
  try {
    allFunctions.value = await GetAllFunctions()
    categories.value = await GetCategories()
  } catch (e) {
    console.error('加载函数库失败:', e)
  }
}

loadFunctions()

const selectCategory = async (cat: string) => {
  selectedCategory.value = cat
  try {
    categoryFunctions.value = await GetFunctionsByCategory(cat)
  } catch (e) {
    categoryFunctions.value = allFunctions.value.filter(f => f.category === cat)
  }
}

const showFunctionDetail = (fn: ChineseFunction) => {
  selectedFunction.value = fn
}

const getCategoryColor = (cat: string): string => {
  const colors: Record<string, string> = {
    '窗口': '#409eff', '文件': '#67c23a', '字符串': '#e6a23c',
    '数值': '#f56c6c', '数组': '#909399', '组件': '#409eff',
    '系统': '#67c23a', '网络': '#e6a23c', '数据库': '#f56c6c',
    'JSON': '#909399', '加密': '#409eff', '多媒体': '#67c23a',
    '定时器': '#e6a23c', '线程': '#f56c6c',
  }
  return colors[cat] || '#909399'
}

const consoleOutputRef = ref<HTMLElement | null>(null)

const addConsoleLog = (msg: string) => {
  const time = new Date().toLocaleTimeString()
  consoleLogs.value.push(`[${time}] ${msg}`)
  nextTick(() => {
    if (consoleOutputRef.value) {
      consoleOutputRef.value.scrollTop = consoleOutputRef.value.scrollHeight
    }
  })
}

const getLevelLabel = (level: string): string => {
  const labels: Record<string, string> = {
    'info': '信息',
    'warn': '警告',
    'error': '错误',
    'output': '输出',
  }
  return labels[level] || level
}

const getLevelClass = (level: string): string => {
  const classes: Record<string, string> = {
    'info': 'is-system',
    'warn': 'is-warn',
    'error': 'is-error',
    'output': 'is-output',
  }
  return classes[level] || ''
}

let offBuildLog: (() => void) | null = null

onMounted(() => {
  EventsOn('buildLog', (data: any) => {
    const level = data?.level || 'info'
    const message = data?.message || ''
    const label = getLevelLabel(level)
    addConsoleLog(`[${label}] ${message}`)
  })
})

onUnmounted(() => {
  if (offBuildLog) {
    offBuildLog()
    offBuildLog = null
  }
})

const clearConsole = () => {
  consoleLogs.value = []
}

const copyConsole = () => {
  const text = consoleLogs.value.join('\n')
  navigator.clipboard.writeText(text).then(() => {
    addConsoleLog('[系统] 控制台内容已复制到剪贴板')
  }).catch(() => {
    addConsoleLog('[系统] 复制失败，请手动选择复制')
  })
}

defineExpose({ addConsoleLog, consoleLogs })
</script>

<template>
  <div class="bottom-panel">
    <div class="panel-tabs">
      <div
        class="tab-item"
        :class="{ active: activeTab === 'console' }"
        @click="activeTab = 'console'"
      >
        <el-icon><Monitor /></el-icon>
        <span>控制台</span>
      </div>
      <div
        class="tab-item"
        :class="{ active: activeTab === 'functions' }"
        @click="activeTab = 'functions'"
      >
        <el-icon><FolderOpened /></el-icon>
        <span>函数说明</span>
      </div>
    </div>

    <div class="panel-body">
      <!-- 控制台 -->
      <div v-show="activeTab === 'console'" class="console-view">
        <div class="console-toolbar">
          <el-button size="small" @click="clearConsole">清空</el-button>
          <el-button size="small" @click="copyConsole">复制全部</el-button>
        </div>
        <div class="console-output" ref="consoleOutputRef">
          <div v-for="(log, index) in consoleLogs" :key="index" class="console-line" :class="{ 'is-system': log.includes('[信息]') || log.includes('[系统]'), 'is-error': log.includes('[错误]'), 'is-warn': log.includes('[警告]'), 'is-output': log.includes('[输出]') }">
            {{ log }}
          </div>
          <div v-if="consoleLogs.length === 0" class="console-empty">控制台为空</div>
        </div>
      </div>

      <!-- 函数说明 -->
      <div v-show="activeTab === 'functions'" class="functions-view">
        <div class="fn-layout">
          <div class="fn-sidebar">
            <div
              v-for="cat in categories"
              :key="cat"
              class="fn-cat-item"
              :class="{ active: selectedCategory === cat }"
              @click="selectCategory(cat)"
            >
              <span class="cat-dot" :style="{ backgroundColor: getCategoryColor(cat) }"></span>
              <span class="cat-name">{{ cat }}</span>
              <span class="cat-count">{{ allFunctions.filter(f => f.category === cat).length }}</span>
            </div>
          </div>
          <div class="fn-list">
            <div class="fn-list-header">
              {{ selectedCategory || '请选择分类' }}
              <span v-if="selectedCategory" class="fn-list-count">{{ categoryFunctions.length }} 个函数</span>
            </div>
            <div
              v-for="fn in categoryFunctions"
              :key="fn.chineseName"
              class="fn-list-item"
              :class="{ active: selectedFunction?.chineseName === fn.chineseName }"
              @click="showFunctionDetail(fn)"
            >
              <span class="fn-list-name">{{ fn.chineseName }}</span>
              <span v-if="fn.returnType" class="fn-list-return">→ {{ fn.returnType }}</span>
            </div>
            <div v-if="!selectedCategory" class="fn-list-empty">请在左侧选择一个分类</div>
          </div>
          <div class="fn-detail">
            <template v-if="selectedFunction">
              <div class="fn-detail-header">
                <span class="fn-detail-name">{{ selectedFunction.chineseName }}</span>
                <span
                  class="fn-detail-cat"
                  :style="{ backgroundColor: getCategoryColor(selectedFunction.category) + '20', color: getCategoryColor(selectedFunction.category) }"
                >{{ selectedFunction.category }}</span>
              </div>
              <div class="fn-detail-section">
                <div class="fn-detail-label">功能描述</div>
                <div class="fn-detail-value">{{ selectedFunction.description }}</div>
              </div>
              <div class="fn-detail-section">
                <div class="fn-detail-label">Go函数</div>
                <div class="fn-detail-value code">{{ selectedFunction.goFunction }}</div>
              </div>
              <div class="fn-detail-section">
                <div class="fn-detail-label">参数列表</div>
                <div class="fn-detail-value">
                  <span v-if="selectedFunction.params.length === 0">无参数</span>
                  <el-tag v-for="(p, i) in selectedFunction.params" :key="i" size="small" class="param-tag">{{ p }}</el-tag>
                </div>
              </div>
              <div class="fn-detail-section">
                <div class="fn-detail-label">返回值</div>
                <div class="fn-detail-value">{{ selectedFunction.returnType || '无返回值' }}</div>
              </div>
            </template>
            <div v-else class="fn-detail-empty">请选择一个函数查看详情</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bottom-panel {
  height: 220px;
  background-color: var(--panel-bg);
  border-top: 1px solid var(--panel-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.panel-tabs {
  display: flex;
  background-color: var(--panel-header-bg);
  border-bottom: 1px solid var(--panel-border);
  padding: 0 8px;
  flex-shrink: 0;
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  font-size: 12px;
  color: var(--text-regular);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
  user-select: none;
}

.tab-item:hover {
  color: var(--color-primary);
}

.tab-item.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
}

.panel-body {
  flex: 1;
  overflow: hidden;
}

.console-view {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.console-toolbar {
  display: flex;
  gap: 6px;
  padding: 4px 8px;
  border-bottom: 1px solid var(--border-lighter);
  flex-shrink: 0;
}

.console-output {
  flex: 1;
  overflow-y: auto;
  padding: 6px 10px;
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 12px;
  background-color: var(--bg-console);
  color: var(--text-console);
  user-select: text;
}

.console-line {
  padding: 1px 0;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
}

.console-line.is-system {
  color: var(--hl-comment);
}

.console-line.is-error {
  color: var(--color-danger);
}

.console-line.is-warn {
  color: #e6a23c;
}

.console-line.is-output {
  color: #67c23a;
}

.console-empty {
  color: var(--text-secondary);
  text-align: center;
  padding: 20px;
}

.functions-view {
  height: 100%;
}

.fn-layout {
  display: flex;
  height: 100%;
}

.fn-sidebar {
  width: 120px;
  border-right: 1px solid var(--border-lighter);
  overflow-y: auto;
  flex-shrink: 0;
  padding: 4px 0;
}

.fn-cat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 10px;
  font-size: 12px;
  color: var(--text-regular);
  cursor: pointer;
  transition: background-color 0.2s;
}

.fn-cat-item:hover {
  background-color: var(--bg-secondary);
}

.fn-cat-item.active {
  background-color: var(--color-primary-light);
  color: var(--color-primary);
}

.cat-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.cat-name {
  flex: 1;
}

.cat-count {
  font-size: 10px;
  color: var(--text-placeholder);
}

.fn-list {
  width: 180px;
  border-right: 1px solid var(--border-lighter);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.fn-list-header {
  padding: 6px 10px;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-lighter);
  flex-shrink: 0;
}

.fn-list-count {
  font-size: 10px;
  color: var(--text-secondary);
  font-weight: normal;
  margin-left: 4px;
}

.fn-list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-lighter);
}

.fn-list-item:hover {
  background-color: var(--bg-secondary);
}

.fn-list-item.active {
  background-color: var(--color-primary-light);
  color: var(--color-primary);
}

.fn-list-name {
  color: var(--text-primary);
}

.fn-list-return {
  font-size: 10px;
  color: var(--color-success);
}

.fn-list-empty {
  padding: 20px;
  text-align: center;
  color: var(--text-placeholder);
  font-size: 12px;
}

.fn-detail {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.fn-detail-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.fn-detail-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.fn-detail-cat {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 3px;
}

.fn-detail-section {
  margin-bottom: 10px;
}

.fn-detail-label {
  font-size: 11px;
  color: var(--text-secondary);
  margin-bottom: 3px;
}

.fn-detail-value {
  font-size: 13px;
  color: var(--text-primary);
}

.fn-detail-value.code {
  font-family: 'Consolas', monospace;
  font-size: 12px;
  background-color: #f5f7fa;
  padding: 4px 8px;
  border-radius: 3px;
}

.param-tag {
  margin-right: 4px;
  margin-bottom: 2px;
}

.fn-detail-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #c0c4cc;
  font-size: 13px;
}
</style>