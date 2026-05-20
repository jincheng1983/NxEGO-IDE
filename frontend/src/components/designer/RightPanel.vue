<script lang="ts" setup>
import { ref, computed, onUnmounted } from 'vue'
import AIChatPanel from '../ai/AIChatPanel.vue'

const sidebarVisible = ref(true)
const userWidth = ref(380)
const minWidth = 280
const maxWidth = 700

const totalWidth = computed(() => sidebarVisible.value ? userWidth.value : 0)

function toggleSidebar() {
  sidebarVisible.value = !sidebarVisible.value
}

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
  const delta = resizeStartX.value - e.clientX
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

onUnmounted(() => {
  document.removeEventListener('mousemove', onResizeMouseMove)
  document.removeEventListener('mouseup', onResizeMouseUp)
})

defineExpose({
  isExpanded: () => sidebarVisible.value,
  setExpanded: (v: boolean) => { sidebarVisible.value = v },
})
</script>

<template>
  <div
    v-show="sidebarVisible"
    class="right-panel"
    :style="{ width: totalWidth + 'px' }"
  >
    <div
      class="resize-handle"
      @mousedown="onResizeMouseDown"
      title="拖拽调整宽度"
    >
      <div class="resize-grip" />
    </div>

    <AIChatPanel />
  </div>
</template>

<style scoped>
.right-panel {
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  background: var(--panel-bg, #ffffff);
  border-left: 1px solid var(--panel-border, #e0e0e0);
  overflow: hidden;
  transition: width 0.2s ease;
  position: relative;
}

.resize-handle {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  cursor: ew-resize;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
}

.resize-handle:hover {
  background: var(--color-primary-light, rgba(64, 158, 255, 0.15));
}

.resize-grip {
  width: 2px;
  height: 32px;
  border-radius: 1px;
  background: var(--text-disabled, #c0c4cc);
  opacity: 0;
  transition: opacity 0.15s;
}

.resize-handle:hover .resize-grip {
  opacity: 1;
  background: var(--color-primary, #409eff);
}
</style>