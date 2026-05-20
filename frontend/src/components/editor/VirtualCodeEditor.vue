<script lang="ts" setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useGoParser } from '../../composables/useGoParser'
import type { RenderBlock, TokenSpan, FlowSegment } from '../../workers/goParserShared'

const props = defineProps<{
  modelValue: string
  readOnly?: boolean
  fontSize?: number
  lineHeight?: number
  showFlowGuides?: boolean
  showLineNumbers?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'cursor-change', pos: { line: number; column: number }): void
}>()

const { blocks, flowMaxDepth, lineCount, isParsing, parse } = useGoParser()

const containerRef = ref<HTMLDivElement | null>(null)
const scrollContainerRef = ref<HTMLDivElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)
const inputRef = ref<HTMLTextAreaElement | null>(null)

const fontSize = computed(() => props.fontSize || 14)
const lineHeight = computed(() => props.lineHeight || 22)
const showFlowGuides = computed(() => props.showFlowGuides !== false)
const showLineNumbers = computed(() => props.showLineNumbers !== false)

const scrollTop = ref(0)
const viewportHeight = ref(600)
const viewportWidth = ref(800)
const cursorLine = ref(0)
const cursorColumn = ref(0)
const selectionStart = ref<{ line: number; column: number } | null>(null)
const selectionEnd = ref<{ line: number; column: number } | null>(null)

const gutterWidth = computed(() => {
  let w = showLineNumbers.value ? 48 : 0
  if (showFlowGuides.value && flowMaxDepth.value > 0) {
    w += flowMaxDepth.value * 14 + 8
  }
  return w
})

const visibleLineCount = computed(() => Math.ceil(viewportHeight.value / lineHeight.value) + 2)
const startLine = computed(() => Math.max(0, Math.floor(scrollTop.value / lineHeight.value) - 1))
const endLine = computed(() => Math.min(blocks.value.length, startLine.value + visibleLineCount.value))
const totalHeight = computed(() => Math.max(blocks.value.length * lineHeight.value, viewportHeight.value))

const visibleBlocks = computed(() => {
  return blocks.value.slice(startLine.value, endLine.value)
})

const codeText = computed({
  get: () => props.modelValue,
  set: (val: string) => emit('update:modelValue', val),
})

watch(() => props.modelValue, (newVal) => {
  if (newVal !== undefined) {
    parse(newVal)
  }
}, { immediate: true })

const handleScroll = () => {
  if (!scrollContainerRef.value) return
  scrollTop.value = scrollContainerRef.value.scrollTop
}

const handleResize = () => {
  if (!containerRef.value) return
  viewportHeight.value = containerRef.value.clientHeight
  viewportWidth.value = containerRef.value.clientWidth
}

const handleInput = (e: Event) => {
  const target = e.target as HTMLTextAreaElement
  codeText.value = target.value
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Tab') {
    e.preventDefault()
    const textarea = inputRef.value
    if (!textarea) return
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const val = textarea.value
    textarea.value = val.slice(0, start) + '\t' + val.slice(end)
    textarea.selectionStart = textarea.selectionEnd = start + 1
    codeText.value = textarea.value
  }
}

const handleTextareaScroll = () => {
  if (!inputRef.value || !scrollContainerRef.value) return
  scrollContainerRef.value.scrollTop = inputRef.value.scrollTop
}

const updateCursorFromTextarea = () => {
  if (!inputRef.value) return
  const pos = inputRef.value.selectionStart
  const text = inputRef.value.value
  let line = 0
  let col = 0
  for (let i = 0; i < pos && i < text.length; i++) {
    if (text[i] === '\n') {
      line++
      col = 0
    } else {
      col++
    }
  }
  cursorLine.value = line
  cursorColumn.value = col
  emit('cursor-change', { line, column: col })
}

const scrollToLine = (line: number) => {
  if (!scrollContainerRef.value) return
  scrollContainerRef.value.scrollTop = line * lineHeight.value
}

const getTokenColor = (type: TokenSpan['type']): string => {
  switch (type) {
    case 'keyword': return 'var(--hl-keyword, #569cd6)'
    case 'string': return 'var(--hl-string, #ce9178)'
    case 'comment': return 'var(--hl-comment, #6a9955)'
    case 'number': return 'var(--hl-number, #b5cea8)'
    case 'type': return 'var(--hl-type, #4ec9b0)'
    case 'function': return 'var(--hl-function, #dcdcaa)'
    case 'operator': return 'var(--hl-operator, #d4d4d4)'
    case 'elang': return 'var(--hl-elang, #c586c0)'
    case 'package': return 'var(--hl-keyword, #569cd6)'
    default: return 'var(--editor-text, #d4d4d4)'
  }
}

const getFlowColor = (depth: number): string => {
  const colors = [
    '#569cd6', '#6a9955', '#ce9178', '#c586c0',
    '#4ec9b0', '#dcdcaa', '#b5cea8', '#d16969',
  ]
  return colors[depth % colors.length]
}

const renderCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const dpr = window.devicePixelRatio || 1
  const w = viewportWidth.value
  const h = viewportHeight.value
  canvas.width = w * dpr
  canvas.height = h * dpr
  canvas.style.width = w + 'px'
  canvas.style.height = h + 'px'
  ctx.scale(dpr, dpr)

  ctx.clearRect(0, 0, w, h)

  const fh = lineHeight.value
  const fs = fontSize.value
  const gutterW = gutterWidth.value
  const codeX = gutterW + 8

  ctx.font = `${fs}px "Consolas", "Courier New", "Microsoft YaHei", monospace`
  ctx.textBaseline = 'middle'

  const bgColor = getComputedStyle(document.documentElement).getPropertyValue('--editor-bg').trim() || '#1e1e1e'
  const textColor = getComputedStyle(document.documentElement).getPropertyValue('--editor-text').trim() || '#d4d4d4'
  const gutterBg = getComputedStyle(document.documentElement).getPropertyValue('--editor-gutter-bg').trim() || '#1e1e1e'
  const gutterText = getComputedStyle(document.documentElement).getPropertyValue('--editor-gutter-text').trim() || '#858585'
  const lineActiveBg = getComputedStyle(document.documentElement).getPropertyValue('--editor-line-active-bg').trim() || '#2a2d2e'
  const selectionBg = getComputedStyle(document.documentElement).getPropertyValue('--editor-selection-bg').trim() || '#264f78'

  ctx.fillStyle = bgColor
  ctx.fillRect(0, 0, w, h)

  ctx.fillStyle = gutterBg
  ctx.fillRect(0, 0, gutterW, h)

  const blocks_ = visibleBlocks.value
  const offsetY = -scrollTop.value + startLine.value * fh

  for (let i = 0; i < blocks_.length; i++) {
    const block = blocks_[i]
    const y = offsetY + i * fh
    const lineIdx = block.lineIndex

    if (lineIdx === cursorLine.value) {
      ctx.fillStyle = lineActiveBg
      ctx.fillRect(0, y, w, fh)
    }

    if (showFlowGuides.value && block.flowSegments.length > 0) {
      const guideStartX = showLineNumbers.value ? 48 : 0
      for (const seg of block.flowSegments) {
        const gx = guideStartX + seg.depth * 14 + 7
        ctx.strokeStyle = getFlowColor(seg.depth)
        ctx.lineWidth = 1
        ctx.globalAlpha = 0.4

        if (seg.kind === 'start') {
          ctx.beginPath()
          ctx.moveTo(gx, y + fh / 2)
          ctx.lineTo(gx + 6, y + fh / 2)
          ctx.lineTo(gx + 6, y + fh)
          ctx.stroke()
        } else if (seg.kind === 'line') {
          ctx.beginPath()
          ctx.moveTo(gx, y)
          ctx.lineTo(gx, y + fh)
          ctx.stroke()
        } else if (seg.kind === 'end') {
          ctx.beginPath()
          ctx.moveTo(gx, y)
          ctx.lineTo(gx, y + fh / 2)
          ctx.lineTo(gx + 6, y + fh / 2)
          ctx.stroke()
        }
        ctx.globalAlpha = 1
      }
    }

    if (showLineNumbers.value) {
      ctx.fillStyle = gutterText
      ctx.textAlign = 'right'
      ctx.fillText(String(lineIdx + 1), gutterW - 8, y + fh / 2)
    }

    if (block.kind === 'blank') {
      continue
    }

    let cx = codeX
    ctx.textAlign = 'left'

    for (const token of block.tokens) {
      ctx.fillStyle = getTokenColor(token.type)
      const tokenText = token.text
      const metrics = ctx.measureText(tokenText)
      ctx.fillText(tokenText, cx, y + fh / 2)
      cx += metrics.width
    }
  }
}

let renderFrameId = 0
const scheduleRender = () => {
  if (renderFrameId) return
  renderFrameId = requestAnimationFrame(() => {
    renderFrameId = 0
    renderCanvas()
  })
}

watch([visibleBlocks, scrollTop, viewportWidth, viewportHeight, cursorLine, flowMaxDepth], () => {
  nextTick(scheduleRender)
}, { deep: false })

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)

  if (props.modelValue) {
    parse(props.modelValue)
  }

  nextTick(() => {
    scheduleRender()
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (renderFrameId) {
    cancelAnimationFrame(renderFrameId)
  }
})

defineExpose({
  scrollToLine,
  getCodeText: () => codeText.value,
})
</script>

<template>
  <div ref="containerRef" class="virtual-code-editor">
    <div
      ref="scrollContainerRef"
      class="virtual-code-editor__scroll"
      @scroll="handleScroll"
    >
      <div
        class="virtual-code-editor__spacer"
        :style="{ height: totalHeight + 'px' }"
      ></div>
    </div>

    <canvas
      ref="canvasRef"
      class="virtual-code-editor__canvas"
    ></canvas>

    <textarea
      v-if="!readOnly"
      ref="inputRef"
      class="virtual-code-editor__input"
      :value="modelValue"
      :style="{
        fontSize: fontSize + 'px',
        lineHeight: lineHeight + 'px',
        paddingLeft: gutterWidth + 8 + 'px',
        fontFamily: '\'Consolas\', \'Courier New\', \'Microsoft YaHei\', monospace',
      }"
      spellcheck="false"
      @input="handleInput"
      @keydown="handleKeydown"
      @scroll="handleTextareaScroll"
      @click="updateCursorFromTextarea"
      @keyup="updateCursorFromTextarea"
    ></textarea>

    <div
      v-if="isParsing"
      class="virtual-code-editor__parsing"
    >
      解析中...
    </div>
  </div>
</template>

<style scoped>
.virtual-code-editor {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background-color: var(--editor-bg, #1e1e1e);
}

.virtual-code-editor__scroll {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: auto;
  overflow-x: hidden;
  z-index: 1;
}

.virtual-code-editor__spacer {
  width: 100%;
  pointer-events: none;
}

.virtual-code-editor__canvas {
  position: absolute;
  top: 0;
  left: 0;
  pointer-events: none;
  z-index: 2;
}

.virtual-code-editor__input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  resize: none;
  background: transparent;
  color: transparent;
  caret-color: var(--editor-cursor, #aeafad);
  white-space: pre;
  overflow-wrap: normal;
  overflow-x: auto;
  overflow-y: auto;
  z-index: 3;
  tab-size: 4;
}

.virtual-code-editor__input::selection {
  background-color: var(--editor-selection-bg, #264f78);
}

.virtual-code-editor__parsing {
  position: absolute;
  top: 8px;
  right: 12px;
  padding: 2px 8px;
  background-color: var(--color-primary, #0078d4);
  color: #fff;
  font-size: 11px;
  border-radius: 3px;
  z-index: 10;
  opacity: 0.9;
}
</style>