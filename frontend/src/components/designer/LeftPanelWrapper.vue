<script lang="ts" setup>
import { ref, computed, watch, onUnmounted } from 'vue'
import LeftPanel from './LeftPanel.vue'
import PropertyPanel from './PropertyPanel.vue'

interface WindowProps {
  title: string
  width: number
  height: number
  [key: string]: any
}

const props = defineProps<{
  windowProps?: WindowProps
  selectedComponent?: any
  selectedComponents?: any[]
}>()

const emit = defineEmits<{
  (e: 'select-component', component: any): void
  (e: 'select-function', fn: any): void
  (e: 'select-window-type', type: string): void
  (e: 'open-window-design', windowId: string): void
  (e: 'open-window-code', windowId: string): void
  (e: 'add-new-window', windowData: any): void
  (e: 'delete-window', windowId: string): void
  (e: 'rename-window', windowId: string, newName: string): void
  (e: 'add-new-module', moduleData: any): void
  (e: 'delete-module', moduleId: string): void
  (e: 'rename-module', moduleId: string, newName: string): void
  (e: 'open-module-code', moduleId: string): void
  (e: 'add-independent-code-file', codeFile: any): void
  (e: 'delete-independent-code-file', fileId: string): void
  (e: 'rename-independent-code-file', fileId: string, newName: string): void
  (e: 'open-independent-code-file', fileId: string): void
  (e: 'window-property-change', key: string, value: any): void
  (e: 'property-change', id: string, key: string, value: any): void
  (e: 'property-batch-change', ids: string[], key: string, value: any): void
}>()

type ActivityItem = 'components' | 'functions' | 'construct' | 'properties'

const activityItems = [
  { id: 'components' as ActivityItem, icon: 'Grid', title: '组件库' },
  { id: 'functions' as ActivityItem, icon: 'Collection', title: '函数库' },
  { id: 'construct' as ActivityItem, icon: 'Monitor', title: '构造' },
  { id: 'properties' as ActivityItem, icon: 'Setting', title: '窗口属性' },
]

const activeItem = ref<ActivityItem>('components')
const sidebarVisible = ref(true)
const userWidth = ref(240)
const minWidth = 200
const maxWidth = 450

const totalWidth = computed(() => {
  const barWidth = 48
  return sidebarVisible.value ? barWidth + userWidth.value : barWidth
})

function onActivityClick(item: ActivityItem) {
  if (activeItem.value === item && sidebarVisible.value) {
    sidebarVisible.value = false
  } else {
    activeItem.value = item
    if (!sidebarVisible.value) {
      sidebarVisible.value = true
    }
  }
}

function toggleSidebar() {
  sidebarVisible.value = !sidebarVisible.value
}

const leftPanelRef = ref<InstanceType<typeof LeftPanel> | null>(null)

const isResizing = ref(false)
const resizeStartX = ref(0)
const resizeStartWidth = ref(0)

function onResizeMouseDown(e: MouseEvent) {
  e.preventDefault()
  isResizing.value = true
  resizeStartX.value = e.clientX
  resizeStartWidth.value = userWidth.value
  document.addEventListener('mousemove', onResizeMouseMove)
  document.addEventListener('mouseup', onResizeMouseUp)
  document.body.style.cursor = 'ew-resize'
  document.body.style.userSelect = 'none'
}

function onResizeMouseMove(e: MouseEvent) {
  if (!isResizing.value) return
  const delta = e.clientX - resizeStartX.value
  const newWidth = Math.min(maxWidth, Math.max(minWidth, resizeStartWidth.value + delta))
  userWidth.value = newWidth
}

function onResizeMouseUp() {
  isResizing.value = false
  document.removeEventListener('mousemove', onResizeMouseMove)
  document.removeEventListener('mouseup', onResizeMouseUp)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
}

watch(() => props.selectedComponent, (comp) => {
  if (comp) {
    activeItem.value = 'properties'
    if (!sidebarVisible.value) {
      sidebarVisible.value = true
    }
  }
})

onUnmounted(() => {
  document.removeEventListener('mousemove', onResizeMouseMove)
  document.removeEventListener('mouseup', onResizeMouseUp)
})

defineExpose({
  toggleSidebar,
})
</script>

<template>
  <div class="vscode-left-panel" :style="{ width: totalWidth + 'px' }">
    <!-- Activity Bar -->
    <div class="activity-bar">
      <div class="activity-header">
        <div class="activity-toggle" @click="toggleSidebar" :title="sidebarVisible ? '收起侧边栏' : '展开侧边栏'">
          <el-icon :size="16"><component :is="sidebarVisible ? 'DArrowLeft' : 'DArrowRight'" /></el-icon>
        </div>
      </div>
      <div class="activity-items">
        <div
          v-for="item in activityItems"
          :key="item.id"
          :class="['activity-item', { active: activeItem === item.id && sidebarVisible }]"
          :title="item.title"
          @click="onActivityClick(item.id)"
        >
          <el-icon :size="20"><component :is="item.icon" /></el-icon>
          <div v-if="activeItem === item.id && sidebarVisible" class="activity-indicator" />
        </div>
      </div>
    </div>

    <!-- Sidebar -->
    <div
      v-show="sidebarVisible"
      class="sidebar-area"
      :style="{ width: userWidth + 'px' }"
    >
      <div
        class="sidebar-resize-handle"
        @mousedown="onResizeMouseDown"
        title="拖拽调整宽度"
      >
        <div class="resize-grip" />
      </div>
      <div class="sidebar-content">
        <LeftPanel
          v-show="activeItem !== 'properties'"
          ref="leftPanelRef"
          :force-active-tab="activeItem"
          hide-tabs
          @select-component="(c: any) => emit('select-component', c)"
          @select-function="(f: any) => emit('select-function', f)"
          @select-window-type="(t: string) => emit('select-window-type', t)"
          @open-window-design="(id: string) => emit('open-window-design', id)"
          @open-window-code="(id: string) => emit('open-window-code', id)"
          @add-new-window="(d: any) => emit('add-new-window', d)"
          @delete-window="(id: string) => emit('delete-window', id)"
          @rename-window="(id: string, n: string) => emit('rename-window', id, n)"
          @add-new-module="(d: any) => emit('add-new-module', d)"
          @delete-module="(id: string) => emit('delete-module', id)"
          @rename-module="(id: string, n: string) => emit('rename-module', id, n)"
          @open-module-code="(id: string) => emit('open-module-code', id)"
          @add-independent-code-file="(cf: any) => emit('add-independent-code-file', cf)"
          @delete-independent-code-file="(fid: string) => emit('delete-independent-code-file', fid)"
          @rename-independent-code-file="(fid: string, n: string) => emit('rename-independent-code-file', fid, n)"
          @open-independent-code-file="(fid: string) => emit('open-independent-code-file', fid)"
        />
        <PropertyPanel
          v-show="activeItem === 'properties'"
          :selected-component="props.selectedComponent || null"
          :selected-components="props.selectedComponents || undefined"
          :window-props="windowProps"
          @window-property-change="(key: string, value: any) => emit('window-property-change', key, value)"
          @property-change="(id: string, key: string, value: any) => emit('property-change', id, key, value)"
          @property-batch-change="(ids: string[], key: string, value: any) => emit('property-batch-change', ids, key, value)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.vscode-left-panel {
  display: flex;
  flex-shrink: 0;
  background: var(--panel-bg);
  border-right: 1px solid var(--panel-border);
  overflow: hidden;
  transition: width 0.2s ease;
}

/* ====== Activity Bar ====== */
.activity-bar {
  display: flex;
  flex-direction: column;
  width: 48px;
  height: 100%;
  background: var(--activity-bar-bg);
  border-right: 1px solid var(--activity-bar-border);
  flex-shrink: 0;
  z-index: 10;
}

.activity-header {
  display: flex;
  justify-content: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--activity-bar-border);
}

.activity-toggle {
  width: 32px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 4px;
  color: var(--activity-toggle-color);
  transition: all 0.15s ease;
}

.activity-toggle:hover {
  color: var(--activity-toggle-hover);
  background: var(--bg-hover);
}

.activity-items {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 0;
  flex: 1;
}

.activity-item {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 6px;
  color: var(--activity-icon-color);
  transition: all 0.12s ease;
}

.activity-item:hover {
  color: var(--activity-icon-hover);
  background: var(--bg-hover);
}

.activity-item.active {
  color: var(--activity-icon-active);
}

.activity-indicator {
  position: absolute;
  left: -1px;
  top: 6px;
  bottom: 6px;
  width: 2px;
  border-radius: 0 2px 2px 0;
  background: var(--activity-indicator);
}

/* ====== Sidebar Area ====== */
.sidebar-area {
  position: relative;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.sidebar-resize-handle {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  cursor: ew-resize;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sidebar-resize-handle:hover {
  background: var(--color-primary-light);
}

.resize-grip {
  width: 2px;
  height: 32px;
  border-radius: 1px;
  background: var(--text-disabled);
  opacity: 0;
  transition: opacity 0.15s;
}

.sidebar-resize-handle:hover .resize-grip {
  opacity: 1;
  background: var(--color-primary);
}

.sidebar-content {
  flex: 1;
  overflow: hidden;
}
</style>