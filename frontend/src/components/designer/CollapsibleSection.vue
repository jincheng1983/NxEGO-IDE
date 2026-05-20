<script lang="ts" setup>
import { ref } from 'vue'

const props = defineProps<{
  title: string
  defaultExpanded?: boolean
}>()

const isExpanded = ref(props.defaultExpanded !== false)

function toggle() {
  isExpanded.value = !isExpanded.value
}
</script>

<template>
  <div class="collapsible-section">
    <div class="section-header" @click="toggle">
      <el-icon :size="12" class="section-arrow" :class="{ expanded: isExpanded }">
        <ArrowRight />
      </el-icon>
      <span class="section-title-text">{{ title }}</span>
    </div>
    <div v-show="isExpanded" class="section-body">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.collapsible-section {
  border-bottom: 1px solid var(--panel-border);
}

.collapsible-section:last-child {
  border-bottom: none;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  cursor: pointer;
  user-select: none;
  background: var(--panel-header-bg);
  transition: background 0.12s ease;
}

.section-header:hover {
  background: var(--bg-hover);
}

.section-arrow {
  color: var(--text-disabled);
  transition: transform 0.15s ease;
  flex-shrink: 0;
}

.section-arrow.expanded {
  transform: rotate(90deg);
}

.section-title-text {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.section-body {
  padding: 4px 8px 6px;
}
</style>