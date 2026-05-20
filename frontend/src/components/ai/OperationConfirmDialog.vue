<script lang="ts" setup>
import { ref, watch } from 'vue'
import { ElDialog, ElButton, ElCheckbox, ElTag } from 'element-plus'
import { useAIPermissions, type OperationType } from '../../composables/useAIPermissions'

const {
  currentRequest,
  requestPermission,
  approveRequest,
  denyRequest,
  getOperationLabel,
  getOperationIcon,
  getRiskLevel,
} = useAIPermissions()

const showDialog = ref(false)
const rememberChoice = ref(false)
const alwaysAsk = ref(false)
const isAnimating = ref(false)

watch(currentRequest, (req) => {
  if (req) {
    showDialog.value = true
    rememberChoice.value = false
    alwaysAsk.value = false
  }
}, { immediate: true })

function onApprove() {
  if (!currentRequest.value) return
  isAnimating.value = true
  approveRequest(currentRequest.value.id, rememberChoice.value, alwaysAsk.value)
  showDialog.value = false
  setTimeout(() => { isAnimating.value = false }, 300)
}

function onDeny() {
  if (!currentRequest.value) return
  isAnimating.value = true
  denyRequest(currentRequest.value.id)
  showDialog.value = false
  setTimeout(() => { isAnimating.value = false }, 300)
}

function onDialogClosed() {
  if (currentRequest.value && showDialog.value) {
    onDeny()
  }
}

function formatTimestamp(ts: number): string {
  const d = new Date(ts)
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:${String(d.getSeconds()).padStart(2, '0')}`
}

function truncateContent(content: string, maxLen = 500): string {
  if (content.length <= maxLen) return content
  return content.substring(0, maxLen) + '\n... (内容已截断)'
}
</script>

<template>
  <ElDialog
    v-if="currentRequest"
    :model-value="showDialog"
    :title="`🔒 AI 操作确认`"
    width="520px"
    :close-on-click-modal="false"
    :close-on-press-escape="true"
    @close="onDialogClosed"
  >
    <div class="perm-dialog-body">
      <div class="perm-risk-banner" :class="`risk-${currentRequest.risk}`">
        <span class="perm-risk-icon">
          {{ currentRequest.risk === 'dangerous' ? '⚠️' : currentRequest.risk === 'moderate' ? '🔶' : '✅' }}
        </span>
        <span class="perm-risk-text">
          {{ currentRequest.risk === 'dangerous' ? '高风险操作' : currentRequest.risk === 'moderate' ? '中等风险' : '安全操作' }}
        </span>
      </div>

      <div class="perm-section">
        <div class="perm-row">
          <span class="perm-label">操作类型</span>
          <span class="perm-value">
            {{ getOperationIcon(currentRequest.type) }} {{ getOperationLabel(currentRequest.type) }}
          </span>
        </div>

        <div class="perm-row">
          <span class="perm-label">操作说明</span>
          <span class="perm-value">{{ currentRequest.description }}</span>
        </div>

        <div v-if="currentRequest.targetPath" class="perm-row">
          <span class="perm-label">目标路径</span>
          <code class="perm-path">{{ currentRequest.targetPath }}</code>
        </div>

        <div v-if="currentRequest.command" class="perm-row">
          <span class="perm-label">执行命令</span>
          <code class="perm-path">{{ currentRequest.command }}</code>
        </div>

        <div v-if="currentRequest.targetContent" class="perm-row">
          <span class="perm-label">操作内容</span>
          <pre class="perm-content">{{ truncateContent(currentRequest.targetContent) }}</pre>
        </div>

        <div class="perm-row">
          <span class="perm-label">时间</span>
          <span class="perm-value perm-time">{{ formatTimestamp(currentRequest.timestamp) }}</span>
        </div>
      </div>

      <div class="perm-options">
        <ElCheckbox v-model="rememberChoice">
          记住此选择（同类型操作自动处理）
        </ElCheckbox>
        <ElCheckbox v-if="rememberChoice" v-model="alwaysAsk" class="perm-sub-check">
          每次都询问
        </ElCheckbox>
      </div>
    </div>

    <template #footer>
      <div class="perm-footer">
        <ElButton @click="onDeny" :disabled="isAnimating">
          {{ currentRequest.risk === 'dangerous' ? '🚫 拒绝' : '拒绝' }}
        </ElButton>
        <ElButton type="primary" @click="onApprove" :disabled="isAnimating">
          {{ currentRequest.risk === 'dangerous' ? '⚠️ 确认允许' : '允许' }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<style scoped>
.perm-dialog-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.perm-risk-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
}

.risk-safe {
  background: rgba(76, 175, 80, 0.1);
  border: 1px solid rgba(76, 175, 80, 0.3);
}

.risk-moderate {
  background: rgba(255, 152, 0, 0.1);
  border: 1px solid rgba(255, 152, 0, 0.3);
}

.risk-dangerous {
  background: rgba(244, 67, 54, 0.1);
  border: 1px solid rgba(244, 67, 54, 0.3);
}

.perm-risk-icon {
  font-size: 18px;
}

.perm-risk-text {
  color: var(--text-primary);
}

.perm-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.perm-row {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.perm-label {
  font-size: 12px;
  color: var(--text-placeholder);
  font-weight: 500;
  text-transform: uppercase;
}

.perm-value {
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.5;
}

.perm-path {
  font-size: 13px;
  font-family: 'Consolas', 'Courier New', monospace;
  background: var(--bg-tertiary);
  padding: 6px 10px;
  border-radius: 4px;
  color: var(--color-primary);
  word-break: break-all;
}

.perm-content {
  font-size: 12px;
  font-family: 'Consolas', 'Courier New', monospace;
  background: var(--bg-tertiary);
  padding: 8px 10px;
  border-radius: 4px;
  color: var(--text-secondary);
  max-height: 160px;
  overflow-y: auto;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.perm-time {
  font-size: 12px;
  color: var(--text-placeholder);
}

.perm-options {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px 0;
  border-top: 1px solid var(--border-lighter);
}

.perm-sub-check {
  margin-left: 24px;
}

.perm-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>