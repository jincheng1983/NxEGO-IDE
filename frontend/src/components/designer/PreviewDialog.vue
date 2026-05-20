<script lang="ts" setup>
import { computed } from 'vue'

interface CanvasComponent {
  id: string
  type: string
  name: string
  x: number
  y: number
  width: number
  height: number
  text: string
  visible: boolean
  color: string
  fontSize: number
}

const props = defineProps<{
  visible: boolean
  components: CanvasComponent[]
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
}>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
})

const previewStyle = computed(() => {
  if (props.components.length === 0) return {}
  const maxX = Math.max(...props.components.map(c => c.x + c.width)) + 40
  const maxY = Math.max(...props.components.map(c => c.y + c.height)) + 40
  return {
    width: Math.max(400, maxX) + 'px',
    height: Math.max(300, maxY) + 'px',
  }
})
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    title="预览窗口"
    width="80%"
    :close-on-click-modal="false"
    draggable
  >
    <div class="preview-container" :style="previewStyle">
      <div
        v-for="comp in components"
        :key="comp.id"
        class="preview-component"
        :class="'preview-' + comp.type"
        :style="{
          left: comp.x + 'px',
          top: comp.y + 'px',
          width: comp.width + 'px',
          height: comp.height + 'px',
          borderColor: comp.color,
          color: comp.color,
          fontSize: comp.fontSize + 'px',
          display: comp.visible ? 'flex' : 'none',
        }"
      >
        <span v-if="comp.type === 'button'" class="preview-btn">{{ comp.text }}</span>
        <span v-else-if="comp.type === 'input'" class="preview-input">{{ comp.text }}</span>
        <span v-else-if="comp.type === 'checkbox'" class="preview-checkbox">
          <input type="checkbox" disabled /> {{ comp.text }}
        </span>
        <span v-else>{{ comp.text }}</span>
      </div>
      <div v-if="components.length === 0" class="preview-empty">
        暂无组件，请先在画布中添加组件
      </div>
    </div>
  </el-dialog>
</template>

<style scoped>
.preview-container {
  position: relative;
  background-color: var(--bg-primary);
  border: 1px solid var(--border-base);
  border-radius: 4px;
  min-height: 300px;
  overflow: auto;
}

.preview-component {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-base);
  border-radius: 4px;
  background-color: var(--bg-primary);
}

.preview-btn {
  background-color: var(--color-primary);
  color: var(--text-inverse);
  padding: 4px 12px;
  border-radius: 4px;
  font-size: inherit;
}

.preview-input {
  border: 1px solid var(--border-base);
  padding: 4px 8px;
  border-radius: 4px;
  width: 100%;
}

.preview-checkbox {
  display: flex;
  align-items: center;
  gap: 4px;
}

.preview-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: var(--text-placeholder);
  font-size: 14px;
}
</style>