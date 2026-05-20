<script lang="ts" setup>
import { ref, computed, watch } from 'vue'

const props = defineProps<{
  visible: boolean
  filePath?: string
  fileName?: string
  fileType?: string
  fileData?: Uint8Array | string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const previewType = ref<'image' | 'audio' | 'video' | 'text' | 'binary' | 'unknown'>('unknown')
const previewUrl = ref('')
const textContent = ref('')
const binaryInfo = ref('')
const zoom = ref(1)
const isPlaying = ref(false)

const imageExtensions = new Set(['png', 'jpg', 'jpeg', 'gif', 'bmp', 'ico', 'webp', 'svg', 'tiff'])
const audioExtensions = new Set(['mp3', 'wav', 'ogg', 'flac', 'aac', 'm4a', 'wma'])
const videoExtensions = new Set(['mp4', 'avi', 'mkv', 'mov', 'wmv', 'flv', 'webm'])
const textExtensions = new Set(['txt', 'md', 'json', 'xml', 'yaml', 'yml', 'toml', 'ini', 'cfg', 'log', 'csv'])

const detectType = () => {
  if (!props.fileName && !props.fileType) {
    previewType.value = 'unknown'
    return
  }

  const ext = (props.fileName || props.fileType || '').split('.').pop()?.toLowerCase() || ''

  if (imageExtensions.has(ext)) {
    previewType.value = 'image'
  } else if (audioExtensions.has(ext)) {
    previewType.value = 'audio'
  } else if (videoExtensions.has(ext)) {
    previewType.value = 'video'
  } else if (textExtensions.has(ext)) {
    previewType.value = 'text'
  } else {
    previewType.value = 'binary'
  }
}

const loadPreview = () => {
  detectType()

  if (props.fileData instanceof Uint8Array) {
    const blob = new Blob([props.fileData])
    previewUrl.value = URL.createObjectURL(blob)

    if (previewType.value === 'text') {
      const decoder = new TextDecoder('utf-8')
      textContent.value = decoder.decode(props.fileData)
    } else if (previewType.value === 'binary') {
      binaryInfo.value = `二进制文件 (${formatBytes(props.fileData.length)})`
    }
  } else if (typeof props.fileData === 'string') {
    if (previewType.value === 'text') {
      textContent.value = props.fileData
    } else if (previewType.value === 'image' || previewType.value === 'audio' || previewType.value === 'video') {
      previewUrl.value = props.fileData
    }
  }
}

const formatBytes = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

const zoomIn = () => {
  zoom.value = Math.min(5, zoom.value + 0.25)
}

const zoomOut = () => {
  zoom.value = Math.max(0.25, zoom.value - 0.25)
}

const resetZoom = () => {
  zoom.value = 1
}

const close = () => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = ''
  }
  emit('close')
}

watch(() => props.visible, (val) => {
  if (val) {
    loadPreview()
  } else {
    if (previewUrl.value) {
      URL.revokeObjectURL(previewUrl.value)
      previewUrl.value = ''
    }
  }
})
</script>

<template>
  <div v-if="visible" class="resource-preview">
    <div class="resource-preview__header">
      <span class="resource-preview__title">
        📁 资源预览: {{ fileName || '未知文件' }}
      </span>
      <div class="resource-preview__actions">
        <template v-if="previewType === 'image'">
          <button class="rp-btn" @click="zoomOut" title="缩小">🔍-</button>
          <span class="rp-zoom">{{ Math.round(zoom * 100) }}%</span>
          <button class="rp-btn" @click="zoomIn" title="放大">🔍+</button>
          <button class="rp-btn" @click="resetZoom" title="重置">↺</button>
        </template>
        <button class="rp-btn rp-btn-close" @click="close" title="关闭">✕</button>
      </div>
    </div>

    <div class="resource-preview__body">
      <div v-if="previewType === 'image'" class="rp-image-container">
        <img
          v-if="previewUrl"
          :src="previewUrl"
          :style="{ transform: `scale(${zoom})` }"
          class="rp-image"
          alt="预览"
        />
        <div v-else class="rp-empty">暂无图片数据</div>
      </div>

      <div v-else-if="previewType === 'audio'" class="rp-audio-container">
        <div class="rp-audio-icon">🎵</div>
        <audio
          v-if="previewUrl"
          :src="previewUrl"
          controls
          class="rp-audio"
        ></audio>
        <div v-else class="rp-empty">暂无音频数据</div>
      </div>

      <div v-else-if="previewType === 'video'" class="rp-video-container">
        <video
          v-if="previewUrl"
          :src="previewUrl"
          controls
          class="rp-video"
        ></video>
        <div v-else class="rp-empty">暂无视频数据</div>
      </div>

      <div v-else-if="previewType === 'text'" class="rp-text-container">
        <pre class="rp-text">{{ textContent }}</pre>
      </div>

      <div v-else-if="previewType === 'binary'" class="rp-binary-container">
        <div class="rp-binary-icon">📦</div>
        <div class="rp-binary-info">{{ binaryInfo }}</div>
        <div class="rp-binary-hint">二进制文件无法预览</div>
      </div>

      <div v-else class="rp-unknown-container">
        <div class="rp-unknown-icon">❓</div>
        <div class="rp-unknown-text">无法识别文件类型</div>
      </div>
    </div>

    <div class="resource-preview__footer">
      <span>类型: {{ previewType }}</span>
      <span v-if="fileName">文件: {{ fileName }}</span>
    </div>
  </div>
</template>

<style scoped>
.resource-preview {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80vw;
  max-width: 900px;
  height: 70vh;
  max-height: 700px;
  background-color: var(--bg-primary, #fff);
  border: 1px solid var(--border-base, #dcdfe6);
  border-radius: 8px;
  box-shadow: var(--shadow-lg, 0 8px 32px rgba(0, 0, 0, 0.2));
  display: flex;
  flex-direction: column;
  z-index: 1000;
}

.resource-preview__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background-color: var(--panel-header-bg, #f5f7fa);
  border-bottom: 1px solid var(--border-light, #e4e7ed);
  border-radius: 8px 8px 0 0;
  flex-shrink: 0;
}

.resource-preview__title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary, #303133);
}

.resource-preview__actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.rp-btn {
  width: 26px;
  height: 26px;
  border: 1px solid var(--border-base, #dcdfe6);
  border-radius: 4px;
  background-color: var(--bg-secondary, #f5f7fa);
  color: var(--text-regular, #606266);
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.rp-btn:hover {
  background-color: var(--bg-hover, #ecf5ff);
  border-color: var(--color-primary, #409eff);
}

.rp-btn-close:hover {
  background-color: var(--color-danger-light, #fef0f0);
  border-color: var(--color-danger, #f56c6c);
  color: var(--color-danger, #f56c6c);
}

.rp-zoom {
  font-size: 11px;
  color: var(--text-secondary, #909399);
  min-width: 40px;
  text-align: center;
}

.resource-preview__body {
  flex: 1;
  overflow: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-canvas, #e8e8e8);
}

.rp-image-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;
}

.rp-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  transition: transform 0.2s ease;
}

.rp-audio-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.rp-audio-icon {
  font-size: 48px;
}

.rp-audio {
  width: 400px;
}

.rp-video-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.rp-video {
  max-width: 100%;
  max-height: 100%;
}

.rp-text-container {
  width: 100%;
  height: 100%;
  overflow: auto;
  padding: 12px;
}

.rp-text {
  margin: 0;
  font-family: 'Consolas', 'Courier New', 'Microsoft YaHei', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-primary, #303133);
  white-space: pre-wrap;
  word-break: break-all;
}

.rp-binary-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.rp-binary-icon {
  font-size: 64px;
}

.rp-binary-info {
  font-size: 16px;
  color: var(--text-primary, #303133);
  font-weight: 600;
}

.rp-binary-hint {
  font-size: 13px;
  color: var(--text-secondary, #909399);
}

.rp-unknown-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.rp-unknown-icon {
  font-size: 64px;
}

.rp-unknown-text {
  font-size: 16px;
  color: var(--text-secondary, #909399);
}

.rp-empty {
  font-size: 14px;
  color: var(--text-secondary, #909399);
}

.resource-preview__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  background-color: var(--statusbar-bg, #f5f7fa);
  border-top: 1px solid var(--border-light, #e4e7ed);
  border-radius: 0 0 8px 8px;
  font-size: 11px;
  color: var(--text-secondary, #909399);
  flex-shrink: 0;
}
</style>