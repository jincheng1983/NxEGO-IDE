<script lang="ts" setup>
import { ref, watch, computed } from 'vue'
import CollapsibleSection from './CollapsibleSection.vue'

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
  zIndex?: number
  enabled?: boolean
  fontName?: string
  fontBold?: boolean
  fontItalic?: boolean
  fontUnderline?: boolean
  hint?: string
  cursor?: string
  align?: string
  anchors?: { left: boolean; top: boolean; right: boolean; bottom: boolean }
  padding?: { left: number; top: number; right: number; bottom: number }
  margin?: { left: number; top: number; right: number; bottom: number }
  autoSize?: boolean
  tabStop?: boolean
  readOnly?: boolean
  passwordMode?: boolean
  passwordChar?: string
  maxLength?: number
  multiLine?: boolean
  numberMode?: boolean
  checked?: boolean
  items?: string
  selectedIndex?: number
  dropDownStyle?: string
  min?: number
  max?: number
  position?: number
  orientation?: string
  rowCount?: number
  colCount?: number
  fixedRows?: number
  fixedCols?: number
  picturePath?: string
  stretch?: boolean
  center?: boolean
  proportional?: boolean
  interval?: number
  url?: string
  shapeType?: string
  penColor?: string
  penWidth?: number
  brushColor?: string
  events?: { name: string; handler: string }[]
  value?: number
  step?: number
  multiSelect?: boolean
  columns?: string
  rows?: string
  groupName?: string
  tabs?: string
  nodes?: string
  src?: string
  stretchMode?: string
  format?: string
  buttons?: string
  backgroundColor?: string
  borderStyle?: string
  textAlign?: string
  fontWeight?: string
  fontStyle?: string
  showText?: boolean
  showValue?: boolean
  showHeader?: boolean
  showLines?: boolean
  showRoot?: boolean
  collapsible?: boolean
  editable?: boolean
  threeState?: boolean
  defaultButton?: boolean
  openInNewWindow?: boolean
}

interface WindowProps {
  title: string
  width: number
  height: number
  windowType?: string
  backgroundColor?: string
  resizable?: boolean
  showInTaskbar?: boolean
  showIcon?: boolean
  iconPath?: string
  minWidth?: number
  minHeight?: number
  maxWidth?: number
  maxHeight?: number
  opacity?: number
  alwaysOnTop?: boolean
  borderStyle?: string
  showMenuBar?: boolean
  showMinBtn?: boolean
  showMaxBtn?: boolean
  showCloseBtn?: boolean
  titleBarHeight?: number
}

const props = defineProps<{
  selectedComponent: CanvasComponent | null
  selectedComponents?: CanvasComponent[]
  windowProps?: WindowProps
}>()

const emit = defineEmits<{
  (e: 'property-change', id: string, key: string, value: any): void
  (e: 'property-batch-change', ids: string[], key: string, value: any): void
  (e: 'window-property-change', key: string, value: any): void
}>()

const localName = ref('')
const localText = ref('')
const localX = ref(0)
const localY = ref(0)
const localWidth = ref(100)
const localHeight = ref(32)
const localColor = ref('#409eff')
const localFontSize = ref(14)
const localVisible = ref(true)
const localZIndex = ref(0)
const localEnabled = ref(true)
const localChecked = ref(false)
const localValue = ref(0)
const localMin = ref(0)
const localMax = ref(100)
const localStep = ref(1)
const localItems = ref('')
const localSelectedIndex = ref(0)
const localPasswordMode = ref(false)
const localReadOnly = ref(false)
const localMultiSelect = ref(false)
const localColumns = ref('')
const localRows = ref('')
const localGroupName = ref('')
const localTabs = ref('')
const localNodes = ref('')
const localSrc = ref('')
const localStretchMode = ref('fill')
const localFormat = ref('yyyy-MM-dd')
const localUrl = ref('')
const localOrientation = ref('horizontal')
const localButtons = ref('')
const localBackgroundColor = ref('#ffffff')
const localBorderStyle = ref('solid')
const localTextAlign = ref('left')
const localFontWeight = ref('normal')
const localFontStyle = ref('normal')
const localShowText = ref(false)
const localShowValue = ref(false)
const localShowHeader = ref(true)
const localShowLines = ref(true)
const localShowRoot = ref(true)
const localCollapsible = ref(false)
const localEditable = ref(false)
const localThreeState = ref(false)
const localDefaultButton = ref(false)
const localOpenInNewWindow = ref(false)
const localMaxLength = ref(0)
const localMultiLine = ref(false)
const localNumberMode = ref(false)
const localPasswordChar = ref('*')
const localDropDownStyle = ref('dropdown')
const localRowCount = ref(2)
const localColCount = ref(2)
const localFixedRows = ref(0)
const localFixedCols = ref(0)
const localPicturePath = ref('')
const localStretch = ref(false)
const localCenter = ref(false)
const localProportional = ref(false)
const localInterval = ref(1000)
const localShapeType = ref('rectangle')
const localPenColor = ref('#000000')
const localPenWidth = ref(1)
const localBrushColor = ref('#ffffff')
const localAutoSize = ref(false)
const localTabStop = ref(true)
const localHint = ref('')
const localCursor = ref('default')
const localAlign = ref('left')
const localFontName = ref('微软雅黑')
const localFontBold = ref(false)
const localFontItalic = ref(false)
const localFontUnderline = ref(false)
const localAnchors = ref({ left: true, top: true, right: false, bottom: false })
const localPadding = ref({ left: 0, top: 0, right: 0, bottom: 0 })
const localMargin = ref({ left: 0, top: 0, right: 0, bottom: 0 })
const localEvents = ref<{ name: string; handler: string }[]>([])

const localWinTitle = ref('我的窗口')
const localWinWidth = ref(800)
const localWinHeight = ref(600)
const localWinType = ref('normal')
const localWinBgColor = ref('#f0f0f0')
const localWinResizable = ref(true)
const localWinShowInTaskbar = ref(true)
const localWinShowIcon = ref(true)
const localWinIconPath = ref('')
const localWinMinWidth = ref(200)
const localWinMinHeight = ref(150)
const localWinMaxWidth = ref(0)
const localWinMaxHeight = ref(0)
const localWinOpacity = ref(100)
const localWinAlwaysOnTop = ref(false)
const localWinBorderStyle = ref('sizable')
const localWinShowMenuBar = ref(false)
const localWinShowMinBtn = ref(true)
const localWinShowMaxBtn = ref(true)
const localWinShowCloseBtn = ref(true)
const localWinTitleBarHeight = ref(32)

watch(() => props.selectedComponent, (comp) => {
  if (comp) {
    localName.value = comp.name
    localText.value = comp.text
    localX.value = comp.x
    localY.value = comp.y
    localWidth.value = comp.width
    localHeight.value = comp.height
    localColor.value = comp.color
    localFontSize.value = comp.fontSize
    localVisible.value = comp.visible
    localZIndex.value = comp.zIndex || 0
    localEnabled.value = comp.enabled !== false
    localChecked.value = comp.checked || false
    localValue.value = comp.value || 0
    localMin.value = comp.min || 0
    localMax.value = comp.max || 100
    localStep.value = comp.step || 1
    localItems.value = comp.items || ''
    localSelectedIndex.value = comp.selectedIndex || 0
    localPasswordMode.value = comp.passwordMode || false
    localReadOnly.value = comp.readOnly || false
    localMultiSelect.value = comp.multiSelect || false
    localColumns.value = comp.columns || ''
    localRows.value = comp.rows || ''
    localGroupName.value = comp.groupName || ''
    localTabs.value = comp.tabs || ''
    localNodes.value = comp.nodes || ''
    localSrc.value = comp.src || ''
    localStretchMode.value = comp.stretchMode || 'fill'
    localFormat.value = comp.format || 'yyyy-MM-dd'
    localUrl.value = comp.url || ''
    localOrientation.value = comp.orientation || 'horizontal'
    localButtons.value = comp.buttons || ''
    localBackgroundColor.value = comp.backgroundColor || '#ffffff'
    localBorderStyle.value = comp.borderStyle || 'solid'
    localTextAlign.value = comp.textAlign || 'left'
    localFontWeight.value = comp.fontWeight || 'normal'
    localFontStyle.value = comp.fontStyle || 'normal'
    localShowText.value = comp.showText || false
    localShowValue.value = comp.showValue || false
    localShowHeader.value = comp.showHeader !== false
    localShowLines.value = comp.showLines !== false
    localShowRoot.value = comp.showRoot !== false
    localCollapsible.value = comp.collapsible || false
    localEditable.value = comp.editable || false
    localThreeState.value = comp.threeState || false
    localDefaultButton.value = comp.defaultButton || false
    localOpenInNewWindow.value = comp.openInNewWindow || false
    localMaxLength.value = comp.maxLength || 0
    localMultiLine.value = comp.multiLine || false
    localNumberMode.value = comp.numberMode || false
    localPasswordChar.value = comp.passwordChar || '*'
    localDropDownStyle.value = comp.dropDownStyle || 'dropdown'
    localRowCount.value = comp.rowCount || 2
    localColCount.value = comp.colCount || 2
    localFixedRows.value = comp.fixedRows || 0
    localFixedCols.value = comp.fixedCols || 0
    localPicturePath.value = comp.picturePath || ''
    localStretch.value = comp.stretch || false
    localCenter.value = comp.center || false
    localProportional.value = comp.proportional || false
    localInterval.value = comp.interval || 1000
    localShapeType.value = comp.shapeType || 'rectangle'
    localPenColor.value = comp.penColor || '#000000'
    localPenWidth.value = comp.penWidth || 1
    localBrushColor.value = comp.brushColor || '#ffffff'
    localAutoSize.value = comp.autoSize || false
    localTabStop.value = comp.tabStop !== false
    localHint.value = comp.hint || ''
    localCursor.value = comp.cursor || 'default'
    localAlign.value = comp.align || 'left'
    localFontName.value = comp.fontName || '微软雅黑'
    localFontBold.value = comp.fontBold || false
    localFontItalic.value = comp.fontItalic || false
    localFontUnderline.value = comp.fontUnderline || false
    localAnchors.value = comp.anchors || { left: true, top: true, right: false, bottom: false }
    localPadding.value = comp.padding || { left: 0, top: 0, right: 0, bottom: 0 }
    localMargin.value = comp.margin || { left: 0, top: 0, right: 0, bottom: 0 }
    localEvents.value = comp.events || []
  }
}, { immediate: true })

watch(() => props.windowProps, (wp) => {
  if (wp) {
    localWinTitle.value = wp.title || '我的窗口'
    localWinWidth.value = wp.width || 800
    localWinHeight.value = wp.height || 600
    localWinType.value = wp.windowType || 'normal'
    localWinBgColor.value = wp.backgroundColor || '#f0f0f0'
    localWinResizable.value = wp.resizable !== false
    localWinShowInTaskbar.value = wp.showInTaskbar !== false
    localWinShowIcon.value = wp.showIcon !== false
    localWinIconPath.value = wp.iconPath || ''
    localWinMinWidth.value = wp.minWidth || 200
    localWinMinHeight.value = wp.minHeight || 150
    localWinMaxWidth.value = wp.maxWidth || 0
    localWinMaxHeight.value = wp.maxHeight || 0
    localWinOpacity.value = wp.opacity || 100
    localWinAlwaysOnTop.value = wp.alwaysOnTop || false
    localWinBorderStyle.value = wp.borderStyle || 'sizable'
    localWinShowMenuBar.value = wp.showMenuBar || false
    localWinShowMinBtn.value = wp.showMinBtn !== false
    localWinShowMaxBtn.value = wp.showMaxBtn !== false
    localWinShowCloseBtn.value = wp.showCloseBtn !== false
    localWinTitleBarHeight.value = wp.titleBarHeight || 32
  }
}, { immediate: true })

const emitChange = (key: string, value: any) => {
  if (props.selectedComponents && props.selectedComponents.length > 1) {
    const ids = props.selectedComponents.map(c => c.id)
    emit('property-batch-change', ids, key, value)
  } else if (props.selectedComponent) {
    emit('property-change', props.selectedComponent.id, key, value)
  }
}

const emitWinChange = (key: string, value: any) => {
  emit('window-property-change', key, value)
}

const compType = () => props.selectedComponent?.type || ''

const isType = (...types: string[]) => types.includes(compType())
</script>

<template>
  <div class="property-panel">
    <div class="panel-header">
      <span v-if="selectedComponents && selectedComponents.length > 1">已选中 {{ selectedComponents.length }} 个组件</span>
      <span v-else-if="selectedComponent">{{ selectedComponent.name }}</span>
      <span v-else>窗口属性</span>
    </div>

    <div class="panel-content">
      <!-- ========== 窗口属性 ========== -->
      <div v-if="!selectedComponent" class="props-section">
        <CollapsibleSection title="基本属性" :default-expanded="true">
          <div class="prop-row">
            <span class="prop-label">标题</span>
            <el-input v-model="localWinTitle" size="small" @change="emitWinChange('title', localWinTitle)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">宽度</span>
            <el-input-number v-model="localWinWidth" size="small" :min="200" :max="3000" :step="10" controls-position="right" @change="emitWinChange('width', localWinWidth)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">高度</span>
            <el-input-number v-model="localWinHeight" size="small" :min="150" :max="3000" :step="10" controls-position="right" @change="emitWinChange('height', localWinHeight)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">窗口类型</span>
            <el-select v-model="localWinType" size="small" @change="emitWinChange('windowType', localWinType)">
              <el-option label="普通窗口" value="normal" />
              <el-option label="对话框" value="dialog" />
              <el-option label="弹出窗口" value="popup" />
              <el-option label="工具窗口" value="tool" />
              <el-option label="启动窗口" value="splash" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">背景颜色</span>
            <el-color-picker v-model="localWinBgColor" size="small" @change="emitWinChange('backgroundColor', localWinBgColor)" />
          </div>
        </CollapsibleSection>

        <CollapsibleSection title="外观">
          <div class="prop-row">
            <span class="prop-label">边框样式</span>
            <el-select v-model="localWinBorderStyle" size="small" @change="emitWinChange('borderStyle', localWinBorderStyle)">
              <el-option label="可调大小" value="sizable" />
              <el-option label="固定大小" value="fixed" />
              <el-option label="无边框" value="none" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">可调整大小</span>
            <el-switch v-model="localWinResizable" size="small" @change="emitWinChange('resizable', localWinResizable)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示在任务栏</span>
            <el-switch v-model="localWinShowInTaskbar" size="small" @change="emitWinChange('showInTaskbar', localWinShowInTaskbar)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示图标</span>
            <el-switch v-model="localWinShowIcon" size="small" @change="emitWinChange('showIcon', localWinShowIcon)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">透明度(%)</span>
            <el-slider v-model="localWinOpacity" size="small" :min="10" :max="100" @change="emitWinChange('opacity', localWinOpacity)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">总在最前</span>
            <el-switch v-model="localWinAlwaysOnTop" size="small" @change="emitWinChange('alwaysOnTop', localWinAlwaysOnTop)" />
          </div>
        </CollapsibleSection>

        <CollapsibleSection title="标题栏">
          <div class="prop-row">
            <span class="prop-label">显示菜单栏</span>
            <el-switch v-model="localWinShowMenuBar" size="small" @change="emitWinChange('showMenuBar', localWinShowMenuBar)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">标题栏高度</span>
            <el-input-number v-model="localWinTitleBarHeight" size="small" :min="20" :max="60" controls-position="right" @change="emitWinChange('titleBarHeight', localWinTitleBarHeight)" />
          </div>
        </CollapsibleSection>

        <CollapsibleSection title="控制按钮">
          <div class="prop-row">
            <span class="prop-label">显示最小化</span>
            <el-switch v-model="localWinShowMinBtn" size="small" @change="emitWinChange('showMinBtn', localWinShowMinBtn)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示最大化</span>
            <el-switch v-model="localWinShowMaxBtn" size="small" @change="emitWinChange('showMaxBtn', localWinShowMaxBtn)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示关闭</span>
            <el-switch v-model="localWinShowCloseBtn" size="small" @change="emitWinChange('showCloseBtn', localWinShowCloseBtn)" />
          </div>
        </CollapsibleSection>

        <CollapsibleSection title="尺寸限制">
          <div class="prop-row">
            <span class="prop-label">最小宽度</span>
            <el-input-number v-model="localWinMinWidth" size="small" :min="0" :max="3000" controls-position="right" @change="emitWinChange('minWidth', localWinMinWidth)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最小高度</span>
            <el-input-number v-model="localWinMinHeight" size="small" :min="0" :max="3000" controls-position="right" @change="emitWinChange('minHeight', localWinMinHeight)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最大宽度</span>
            <el-input-number v-model="localWinMaxWidth" size="small" :min="0" :max="3000" controls-position="right" @change="emitWinChange('maxWidth', localWinMaxWidth)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最大高度</span>
            <el-input-number v-model="localWinMaxHeight" size="small" :min="0" :max="3000" controls-position="right" @change="emitWinChange('maxHeight', localWinMaxHeight)" />
          </div>
        </CollapsibleSection>
      </div>

      <!-- ========== 组件属性 ========== -->
      <div v-if="selectedComponent">
        <CollapsibleSection title="通用属性" :default-expanded="true">
          <div class="prop-row">
            <span class="prop-label">名称</span>
            <el-input v-model="localName" size="small" @change="emitChange('name', localName)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">X坐标</span>
            <el-input-number v-model="localX" size="small" :min="0" :max="3000" controls-position="right" @change="emitChange('x', localX)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">Y坐标</span>
            <el-input-number v-model="localY" size="small" :min="0" :max="3000" controls-position="right" @change="emitChange('y', localY)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">宽度</span>
            <el-input-number v-model="localWidth" size="small" :min="10" :max="3000" controls-position="right" @change="emitChange('width', localWidth)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">高度</span>
            <el-input-number v-model="localHeight" size="small" :min="10" :max="3000" controls-position="right" @change="emitChange('height', localHeight)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">层级</span>
            <el-input-number v-model="localZIndex" size="small" :min="0" :max="999" controls-position="right" @change="emitChange('zIndex', localZIndex)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">可见</span>
            <el-switch v-model="localVisible" size="small" @change="emitChange('visible', localVisible)" />
          </div>
        </CollapsibleSection>

        <!-- 按钮 -->
        <template v-if="isType('button')">
          <CollapsibleSection title="按钮属性">
            <div class="prop-row">
              <span class="prop-label">文本</span>
              <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
            </div>
            <div class="prop-row">
              <span class="prop-label">背景颜色</span>
              <el-color-picker v-model="localColor" size="small" @change="emitChange('color', localColor)" />
            </div>
            <div class="prop-row">
              <span class="prop-label">字体大小</span>
              <el-input-number v-model="localFontSize" size="small" :min="8" :max="72" controls-position="right" @change="emitChange('fontSize', localFontSize)" />
            </div>
            <div class="prop-row">
              <span class="prop-label">启用</span>
              <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
            </div>
            <div class="prop-row">
              <span class="prop-label">默认按钮</span>
              <el-switch v-model="localDefaultButton" size="small" @change="emitChange('defaultButton', localDefaultButton)" />
            </div>
          </CollapsibleSection>
        </template>

        <!-- 标签 -->
        <template v-if="isType('label')">
          <CollapsibleSection title="标签属性">
          <div class="prop-row">
            <span class="prop-label">文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">文字颜色</span>
            <el-color-picker v-model="localColor" size="small" @change="emitChange('color', localColor)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">字体大小</span>
            <el-input-number v-model="localFontSize" size="small" :min="8" :max="72" controls-position="right" @change="emitChange('fontSize', localFontSize)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">对齐</span>
            <el-select v-model="localTextAlign" size="small" @change="emitChange('textAlign', localTextAlign)">
              <el-option label="左" value="left" />
              <el-option label="中" value="center" />
              <el-option label="右" value="right" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">粗体</span>
            <el-switch v-model="localFontWeight" size="small" true-value="bold" false-value="normal" @change="emitChange('fontWeight', localFontWeight)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">斜体</span>
            <el-switch v-model="localFontStyle" size="small" true-value="italic" false-value="normal" @change="emitChange('fontStyle', localFontStyle)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 输入框 -->
        <template v-if="isType('input')">
          <CollapsibleSection title="输入框属性">
          <div class="prop-row">
            <span class="prop-label">占位文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">字体大小</span>
            <el-input-number v-model="localFontSize" size="small" :min="8" :max="72" controls-position="right" @change="emitChange('fontSize', localFontSize)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最大长度</span>
            <el-input-number v-model="localMaxLength" size="small" :min="0" :max="99999" controls-position="right" @change="emitChange('maxLength', localMaxLength)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">密码模式</span>
            <el-switch v-model="localPasswordMode" size="small" @change="emitChange('passwordMode', localPasswordMode)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">只读</span>
            <el-switch v-model="localReadOnly" size="small" @change="emitChange('readOnly', localReadOnly)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">启用</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 文本域 -->
        <template v-if="isType('textarea')">
          <CollapsibleSection title="文本域属性">
          <div class="prop-row">
            <span class="prop-label">占位文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">字体大小</span>
            <el-input-number v-model="localFontSize" size="small" :min="8" :max="72" controls-position="right" @change="emitChange('fontSize', localFontSize)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">只读</span>
            <el-switch v-model="localReadOnly" size="small" @change="emitChange('readOnly', localReadOnly)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">启用</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 复选框 -->
        <template v-if="isType('checkbox')">
          <CollapsibleSection title="复选框属性">
          <div class="prop-row">
            <span class="prop-label">文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">选中</span>
            <el-switch v-model="localChecked" size="small" @change="emitChange('checked', localChecked)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">三态</span>
            <el-switch v-model="localThreeState" size="small" @change="emitChange('threeState', localThreeState)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">启用</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 单选框 -->
        <template v-if="isType('radio')">
          <CollapsibleSection title="单选框属性">
          <div class="prop-row">
            <span class="prop-label">文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">分组名</span>
            <el-input v-model="localGroupName" size="small" @change="emitChange('groupName', localGroupName)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">选中</span>
            <el-switch v-model="localChecked" size="small" @change="emitChange('checked', localChecked)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">启用</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 下拉框 -->
        <template v-if="isType('combo')">
          <CollapsibleSection title="下拉框属性">
          <div class="prop-row prop-row-multi">
            <span class="prop-label">选项列表</span>
            <el-input v-model="localItems" size="small" type="textarea" :rows="3" @change="emitChange('items', localItems)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">默认索引</span>
            <el-input-number v-model="localSelectedIndex" size="small" :min="0" :max="99" controls-position="right" @change="emitChange('selectedIndex', localSelectedIndex)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">可编辑</span>
            <el-switch v-model="localEditable" size="small" @change="emitChange('editable', localEditable)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">启用</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 进度条 -->
        <template v-if="isType('progress')">
          <CollapsibleSection title="进度条属性">
          <div class="prop-row">
            <span class="prop-label">当前值</span>
            <el-input-number v-model="localValue" size="small" :min="0" :max="localMax" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最小值</span>
            <el-input-number v-model="localMin" size="small" :min="0" :max="99999" controls-position="right" @change="emitChange('min', localMin)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最大值</span>
            <el-input-number v-model="localMax" size="small" :min="1" :max="99999" controls-position="right" @change="emitChange('max', localMax)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">颜色</span>
            <el-color-picker v-model="localColor" size="small" @change="emitChange('color', localColor)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示百分比</span>
            <el-switch v-model="localShowText" size="small" @change="emitChange('showText', localShowText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 滑块 -->
        <template v-if="isType('slider')">
          <CollapsibleSection title="滑块属性">
          <div class="prop-row">
            <span class="prop-label">当前值</span>
            <el-input-number v-model="localValue" size="small" :min="localMin" :max="localMax" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最小值</span>
            <el-input-number v-model="localMin" size="small" :min="0" :max="99999" controls-position="right" @change="emitChange('min', localMin)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最大值</span>
            <el-input-number v-model="localMax" size="small" :min="1" :max="99999" controls-position="right" @change="emitChange('max', localMax)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">步长</span>
            <el-input-number v-model="localStep" size="small" :min="1" :max="1000" controls-position="right" @change="emitChange('step', localStep)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示数值</span>
            <el-switch v-model="localShowValue" size="small" @change="emitChange('showValue', localShowValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 列表框 -->
        <template v-if="isType('listbox')">
          <CollapsibleSection title="列表框属性">
          <div class="prop-row prop-row-multi">
            <span class="prop-label">列表项</span>
            <el-input v-model="localItems" size="small" type="textarea" :rows="3" @change="emitChange('items', localItems)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">多选</span>
            <el-switch v-model="localMultiSelect" size="small" @change="emitChange('multiSelect', localMultiSelect)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 表格 -->
        <template v-if="isType('table')">
          <CollapsibleSection title="表格属性">
          <div class="prop-row">
            <span class="prop-label">列定义</span>
            <el-input v-model="localColumns" size="small" @change="emitChange('columns', localColumns)" />
          </div>
          <div class="prop-row prop-row-multi">
            <span class="prop-label">行数据</span>
            <el-input v-model="localRows" size="small" type="textarea" :rows="3" @change="emitChange('rows', localRows)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示表头</span>
            <el-switch v-model="localShowHeader" size="small" @change="emitChange('showHeader', localShowHeader)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 分组框 -->
        <template v-if="isType('group')">
          <CollapsibleSection title="分组框属性">
          <div class="prop-row">
            <span class="prop-label">标题</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">可折叠</span>
            <el-switch v-model="localCollapsible" size="small" @change="emitChange('collapsible', localCollapsible)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 选项卡 -->
        <template v-if="isType('tab')">
          <CollapsibleSection title="选项卡属性">
          <div class="prop-row prop-row-multi">
            <span class="prop-label">标签页</span>
            <el-input v-model="localTabs" size="small" type="textarea" :rows="3" @change="emitChange('tabs', localTabs)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">默认索引</span>
            <el-input-number v-model="localSelectedIndex" size="small" :min="0" :max="99" controls-position="right" @change="emitChange('selectedIndex', localSelectedIndex)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 树形框 -->
        <template v-if="isType('tree')">
          <CollapsibleSection title="树形框属性">
          <div class="prop-row prop-row-multi">
            <span class="prop-label">节点数据</span>
            <el-input v-model="localNodes" size="small" type="textarea" :rows="3" @change="emitChange('nodes', localNodes)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">连接线</span>
            <el-switch v-model="localShowLines" size="small" @change="emitChange('showLines', localShowLines)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">根节点</span>
            <el-switch v-model="localShowRoot" size="small" @change="emitChange('showRoot', localShowRoot)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 图片框 -->
        <template v-if="isType('image')">
          <CollapsibleSection title="图片框属性">
          <div class="prop-row">
            <span class="prop-label">图片路径</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">缩放模式</span>
            <el-select v-model="localStretchMode" size="small" @change="emitChange('stretchMode', localStretchMode)">
              <el-option label="填充" value="fill" />
              <el-option label="居中" value="center" />
              <el-option label="拉伸" value="stretch" />
              <el-option label="缩放" value="zoom" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 日期选择 -->
        <template v-if="isType('datepicker')">
          <CollapsibleSection title="日期选择属性">
          <div class="prop-row">
            <span class="prop-label">日期格式</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="yyyy-MM-dd" value="yyyy-MM-dd" />
              <el-option label="yyyy/MM/dd" value="yyyy/MM/dd" />
              <el-option label="yyyy-MM-dd HH:mm:ss" value="yyyy-MM-dd HH:mm:ss" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 超链接 -->
        <template v-if="isType('link')">
          <CollapsibleSection title="超链接属性">
          <div class="prop-row">
            <span class="prop-label">文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">链接地址</span>
            <el-input v-model="localUrl" size="small" @change="emitChange('url', localUrl)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">新窗口打开</span>
            <el-switch v-model="localOpenInNewWindow" size="small" @change="emitChange('openInNewWindow', localOpenInNewWindow)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 分隔线 -->
        <template v-if="isType('divider')">
          <CollapsibleSection title="分隔线属性">
          <div class="prop-row">
            <span class="prop-label">方向</span>
            <el-select v-model="localOrientation" size="small" @change="emitChange('orientation', localOrientation)">
              <el-option label="水平" value="horizontal" />
              <el-option label="垂直" value="vertical" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">颜色</span>
            <el-color-picker v-model="localColor" size="small" @change="emitChange('color', localColor)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 状态栏 -->
        <template v-if="isType('statusbar')">
          <CollapsibleSection title="状态栏属性">
          <div class="prop-row">
            <span class="prop-label">文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 菜单栏 -->
        <template v-if="isType('menu')">
          <CollapsibleSection title="菜单栏属性">
          <div class="prop-row prop-row-multi">
            <span class="prop-label">菜单项</span>
            <el-input v-model="localItems" size="small" type="textarea" :rows="3" @change="emitChange('items', localItems)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 工具栏 -->
        <template v-if="isType('toolbar')">
          <CollapsibleSection title="工具栏属性">
          <div class="prop-row">
            <span class="prop-label">按钮列表</span>
            <el-input v-model="localButtons" size="small" @change="emitChange('buttons', localButtons)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 容器 -->
        <template v-if="isType('container')">
          <CollapsibleSection title="容器属性">
          <div class="prop-row">
            <span class="prop-label">背景颜色</span>
            <el-color-picker v-model="localBackgroundColor" size="small" @change="emitChange('backgroundColor', localBackgroundColor)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">边框样式</span>
            <el-select v-model="localBorderStyle" size="small" @change="emitChange('borderStyle', localBorderStyle)">
              <el-option label="实线" value="solid" />
              <el-option label="虚线" value="dashed" />
              <el-option label="点线" value="dotted" />
              <el-option label="无" value="none" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 面板 -->
        <template v-if="isType('panel')">
          <CollapsibleSection title="面板属性">
          <div class="prop-row">
            <span class="prop-label">背景颜色</span>
            <el-color-picker v-model="localBackgroundColor" size="small" @change="emitChange('backgroundColor', localBackgroundColor)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">边框样式</span>
            <el-select v-model="localBorderStyle" size="small" @change="emitChange('borderStyle', localBorderStyle)">
              <el-option label="实线" value="solid" />
              <el-option label="虚线" value="dashed" />
              <el-option label="点线" value="dotted" />
              <el-option label="无" value="none" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 滚动框 -->
        <template v-if="isType('scrollbox')">
          <CollapsibleSection title="滚动框属性">
          <div class="prop-row">
            <span class="prop-label">背景颜色</span>
            <el-color-picker v-model="localBackgroundColor" size="small" @change="emitChange('backgroundColor', localBackgroundColor)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">边框样式</span>
            <el-select v-model="localBorderStyle" size="small" @change="emitChange('borderStyle', localBorderStyle)">
              <el-option label="实线" value="solid" />
              <el-option label="虚线" value="dashed" />
              <el-option label="点线" value="dotted" />
              <el-option label="无" value="none" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 分隔条 -->
        <template v-if="isType('splitter')">
          <CollapsibleSection title="分隔条属性">
          <div class="prop-row">
            <span class="prop-label">方向</span>
            <el-select v-model="localOrientation" size="small" @change="emitChange('orientation', localOrientation)">
              <el-option label="水平" value="horizontal" />
              <el-option label="垂直" value="vertical" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">颜色</span>
            <el-color-picker v-model="localColor" size="small" @change="emitChange('color', localColor)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 外形框 -->
        <template v-if="isType('shape')">
          <CollapsibleSection title="外形框属性">
          <div class="prop-row">
            <span class="prop-label">形状类型</span>
            <el-select v-model="localShapeType" size="small" @change="emitChange('shapeType', localShapeType)">
              <el-option label="矩形" value="rectangle" />
              <el-option label="圆形" value="circle" />
              <el-option label="椭圆" value="ellipse" />
              <el-option label="圆角矩形" value="roundRect" />
              <el-option label="三角形" value="triangle" />
              <el-option label="菱形" value="diamond" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">画笔颜色</span>
            <el-color-picker v-model="localPenColor" size="small" @change="emitChange('penColor', localPenColor)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">画笔宽度</span>
            <el-input-number v-model="localPenWidth" size="small" :min="1" :max="20" controls-position="right" @change="emitChange('penWidth', localPenWidth)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">画刷颜色</span>
            <el-color-picker v-model="localBrushColor" size="small" @change="emitChange('brushColor', localBrushColor)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 斜角框 -->
        <template v-if="isType('bevel')">
          <CollapsibleSection title="斜角框属性">
          <div class="prop-row">
            <span class="prop-label">样式</span>
            <el-select v-model="localBorderStyle" size="small" @change="emitChange('borderStyle', localBorderStyle)">
              <el-option label="凸起" value="raised" />
              <el-option label="凹陷" value="lowered" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 数字输入/调节器 -->
        <template v-if="isType('number', 'spin')">
          <CollapsibleSection title="数字输入属性">
          <div class="prop-row">
            <span class="prop-label">当前值</span>
            <el-input-number v-model="localValue" size="small" :min="localMin" :max="localMax" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最小值</span>
            <el-input-number v-model="localMin" size="small" :min="0" :max="99999" controls-position="right" @change="emitChange('min', localMin)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最大值</span>
            <el-input-number v-model="localMax" size="small" :min="1" :max="99999" controls-position="right" @change="emitChange('max', localMax)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">步长</span>
            <el-input-number v-model="localStep" size="small" :min="1" :max="1000" controls-position="right" @change="emitChange('step', localStep)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">只读</span>
            <el-switch v-model="localReadOnly" size="small" @change="emitChange('readOnly', localReadOnly)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 时间选择 -->
        <template v-if="isType('timepicker')">
          <CollapsibleSection title="时间选择属性">
          <div class="prop-row">
            <span class="prop-label">时间格式</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="HH:mm:ss" value="HH:mm:ss" />
              <el-option label="HH:mm" value="HH:mm" />
              <el-option label="hh:mm:ss tt" value="hh:mm:ss tt" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 颜色选择 -->
        <template v-if="isType('colorpicker')">
          <CollapsibleSection title="颜色选择属性">
          <div class="prop-row">
            <span class="prop-label">默认颜色</span>
            <el-color-picker v-model="localColor" size="small" @change="emitChange('color', localColor)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 文件选择 -->
        <template v-if="isType('filepicker')">
          <CollapsibleSection title="文件选择属性">
          <div class="prop-row">
            <span class="prop-label">占位文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">只读</span>
            <el-switch v-model="localReadOnly" size="small" @change="emitChange('readOnly', localReadOnly)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 热键框 -->
        <template v-if="isType('hotkey')">
          <CollapsibleSection title="热键框属性">
          <div class="prop-row">
            <span class="prop-label">默认热键</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- IP地址框 -->
        <template v-if="isType('ipedit')">
          <CollapsibleSection title="IP地址框属性">
          <div class="prop-row">
            <span class="prop-label">默认地址</span>
            <el-input v-model="localText" size="small" placeholder="0.0.0.0" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 遮罩框 -->
        <template v-if="isType('maskedit')">
          <CollapsibleSection title="遮罩框属性">
          <div class="prop-row">
            <span class="prop-label">掩码</span>
            <el-input v-model="localText" size="small" placeholder="(999) 000-0000" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 多行编辑框 -->
        <template v-if="isType('memo')">
          <CollapsibleSection title="多行编辑框属性">
          <div class="prop-row">
            <span class="prop-label">占位文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">字体大小</span>
            <el-input-number v-model="localFontSize" size="small" :min="8" :max="72" controls-position="right" @change="emitChange('fontSize', localFontSize)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">只读</span>
            <el-switch v-model="localReadOnly" size="small" @change="emitChange('readOnly', localReadOnly)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">启用</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 选择列表框 -->
        <template v-if="isType('checklistbox')">
          <CollapsibleSection title="选择列表框属性">
          <div class="prop-row prop-row-multi">
            <span class="prop-label">列表项</span>
            <el-input v-model="localItems" size="small" type="textarea" :rows="3" @change="emitChange('items', localItems)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 超级列表框 -->
        <template v-if="isType('listview')">
          <CollapsibleSection title="超级列表框属性">
          <div class="prop-row">
            <span class="prop-label">列定义</span>
            <el-input v-model="localColumns" size="small" @change="emitChange('columns', localColumns)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示表头</span>
            <el-switch v-model="localShowHeader" size="small" @change="emitChange('showHeader', localShowHeader)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">显示网格线</span>
            <el-switch v-model="localShowLines" size="small" @change="emitChange('showLines', localShowLines)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 图表 -->
        <template v-if="isType('chart')">
          <CollapsibleSection title="图表属性">
          <div class="prop-row">
            <span class="prop-label">图表类型</span>
            <el-select v-model="localOrientation" size="small" @change="emitChange('orientation', localOrientation)">
              <el-option label="柱状图" value="bar" />
              <el-option label="折线图" value="line" />
              <el-option label="饼图" value="pie" />
              <el-option label="散点图" value="scatter" />
              <el-option label="面积图" value="area" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">标题</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 动画框 -->
        <template v-if="isType('animate')">
          <CollapsibleSection title="动画框属性">
          <div class="prop-row">
            <span class="prop-label">动画文件</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">自动播放</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 媒体播放器 -->
        <template v-if="isType('mediaplayer')">
          <CollapsibleSection title="媒体播放器属性">
          <div class="prop-row">
            <span class="prop-label">媒体文件</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">自动播放</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">循环播放</span>
            <el-switch v-model="localChecked" size="small" @change="emitChange('checked', localChecked)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 浏览器 -->
        <template v-if="isType('webview')">
          <CollapsibleSection title="浏览器属性">
          <div class="prop-row">
            <span class="prop-label">URL</span>
            <el-input v-model="localUrl" size="small" @change="emitChange('url', localUrl)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 条形码 -->
        <template v-if="isType('barcode')">
          <CollapsibleSection title="条形码属性">
          <div class="prop-row">
            <span class="prop-label">编码内容</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">编码类型</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="Code128" value="Code128" />
              <el-option label="Code39" value="Code39" />
              <el-option label="EAN13" value="EAN13" />
              <el-option label="UPC" value="UPC" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 二维码 -->
        <template v-if="isType('qrcode')">
          <CollapsibleSection title="二维码属性">
          <div class="prop-row">
            <span class="prop-label">编码内容</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">容错级别</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="L (7%)" value="L" />
              <el-option label="M (15%)" value="M" />
              <el-option label="Q (25%)" value="Q" />
              <el-option label="H (30%)" value="H" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 滚动条 -->
        <template v-if="isType('scrollbar')">
          <CollapsibleSection title="滚动条属性">
          <div class="prop-row">
            <span class="prop-label">当前值</span>
            <el-input-number v-model="localValue" size="small" :min="localMin" :max="localMax" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最小值</span>
            <el-input-number v-model="localMin" size="small" :min="0" :max="99999" controls-position="right" @change="emitChange('min', localMin)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最大值</span>
            <el-input-number v-model="localMax" size="small" :min="1" :max="99999" controls-position="right" @change="emitChange('max', localMax)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">方向</span>
            <el-select v-model="localOrientation" size="small" @change="emitChange('orientation', localOrientation)">
              <el-option label="水平" value="horizontal" />
              <el-option label="垂直" value="vertical" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 定时器 -->
        <template v-if="isType('timer')">
          <CollapsibleSection title="定时器属性">
          <div class="prop-row">
            <span class="prop-label">间隔(ms)</span>
            <el-input-number v-model="localInterval" size="small" :min="1" :max="999999" controls-position="right" @change="emitChange('interval', localInterval)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">启用</span>
            <el-switch v-model="localEnabled" size="small" @change="emitChange('enabled', localEnabled)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 托盘图标 -->
        <template v-if="isType('trayicon')">
          <CollapsibleSection title="托盘图标属性">
          <div class="prop-row">
            <span class="prop-label">提示文本</span>
            <el-input v-model="localHint" size="small" @change="emitChange('hint', localHint)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">图标路径</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">可见</span>
            <el-switch v-model="localVisible" size="small" @change="emitChange('visible', localVisible)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 打印机 -->
        <template v-if="isType('printer')">
          <CollapsibleSection title="打印机属性">
          <div class="prop-row">
            <span class="prop-label">纸张方向</span>
            <el-select v-model="localOrientation" size="small" @change="emitChange('orientation', localOrientation)">
              <el-option label="纵向" value="portrait" />
              <el-option label="横向" value="landscape" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 数据源 -->
        <template v-if="isType('datasource')">
          <CollapsibleSection title="数据源属性">
          <div class="prop-row">
            <span class="prop-label">连接字符串</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 客户端Socket -->
        <template v-if="isType('clientsocket')">
          <CollapsibleSection title="客户端Socket属性">
          <div class="prop-row">
            <span class="prop-label">远程地址</span>
            <el-input v-model="localText" size="small" placeholder="127.0.0.1" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">远程端口</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="65535" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 服务器Socket -->
        <template v-if="isType('serversocket')">
          <CollapsibleSection title="服务器Socket属性">
          <div class="prop-row">
            <span class="prop-label">监听端口</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="65535" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- UDP数据报 -->
        <template v-if="isType('udpsocket')">
          <CollapsibleSection title="UDP数据报属性">
          <div class="prop-row">
            <span class="prop-label">本地端口</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="65535" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 串口 -->
        <template v-if="isType('serialport')">
          <CollapsibleSection title="串口属性">
          <div class="prop-row">
            <span class="prop-label">端口号</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="256" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">波特率</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="9600" value="9600" />
              <el-option label="19200" value="19200" />
              <el-option label="38400" value="38400" />
              <el-option label="57600" value="57600" />
              <el-option label="115200" value="115200" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- FTP -->
        <template v-if="isType('ftp')">
          <CollapsibleSection title="FTP属性">
          <div class="prop-row">
            <span class="prop-label">服务器地址</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">端口</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="65535" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- HTTP -->
        <template v-if="isType('http')">
          <CollapsibleSection title="HTTP属性">
          <div class="prop-row">
            <span class="prop-label">基础URL</span>
            <el-input v-model="localUrl" size="small" @change="emitChange('url', localUrl)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">超时(秒)</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="300" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 邮件 -->
        <template v-if="isType('mail')">
          <CollapsibleSection title="邮件属性">
          <div class="prop-row">
            <span class="prop-label">SMTP服务器</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">端口</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="65535" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- XML解析器 -->
        <template v-if="isType('xmlparser')">
          <CollapsibleSection title="XML解析器属性">
          <div class="prop-row prop-row-multi">
            <span class="prop-label">XML文本</span>
            <el-input v-model="localText" size="small" type="textarea" :rows="3" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- JSON解析器 -->
        <template v-if="isType('jsonparser')">
          <CollapsibleSection title="JSON解析器属性">
          <div class="prop-row prop-row-multi">
            <span class="prop-label">JSON文本</span>
            <el-input v-model="localText" size="small" type="textarea" :rows="3" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 正则表达式 -->
        <template v-if="isType('regex')">
          <CollapsibleSection title="正则表达式属性">
          <div class="prop-row">
            <span class="prop-label">表达式</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 数据库连接 -->
        <template v-if="isType('dbconnection')">
          <CollapsibleSection title="数据库连接属性">
          <div class="prop-row">
            <span class="prop-label">连接字符串</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 记录集 -->
        <template v-if="isType('recordset')">
          <CollapsibleSection title="记录集属性">
          <div class="prop-row">
            <span class="prop-label">SQL语句</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- Excel对象 -->
        <template v-if="isType('excel')">
          <CollapsibleSection title="Excel对象属性">
          <div class="prop-row">
            <span class="prop-label">文件路径</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">可见</span>
            <el-switch v-model="localVisible" size="small" @change="emitChange('visible', localVisible)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- Word对象 -->
        <template v-if="isType('word')">
          <CollapsibleSection title="Word对象属性">
          <div class="prop-row">
            <span class="prop-label">文件路径</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">可见</span>
            <el-switch v-model="localVisible" size="small" @change="emitChange('visible', localVisible)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- PDF对象 -->
        <template v-if="isType('pdf')">
          <CollapsibleSection title="PDF对象属性">
          <div class="prop-row">
            <span class="prop-label">文件路径</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 压缩对象 -->
        <template v-if="isType('zip')">
          <CollapsibleSection title="压缩对象属性">
          <div class="prop-row">
            <span class="prop-label">压缩级别</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="9" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 加密对象 -->
        <template v-if="isType('crypto')">
          <CollapsibleSection title="加密对象属性">
          <div class="prop-row">
            <span class="prop-label">加密算法</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="MD5" value="MD5" />
              <el-option label="SHA1" value="SHA1" />
              <el-option label="SHA256" value="SHA256" />
              <el-option label="SHA512" value="SHA512" />
              <el-option label="AES" value="AES" />
              <el-option label="DES" value="DES" />
              <el-option label="RSA" value="RSA" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">密钥</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 线程池 -->
        <template v-if="isType('threadpool')">
          <CollapsibleSection title="线程池属性">
          <div class="prop-row">
            <span class="prop-label">最大线程数</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="1024" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 临界区 -->
        <template v-if="isType('criticalsection')">
          <CollapsibleSection title="临界区属性">
          <div class="prop-row">
            <span class="prop-label">名称</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 互斥体 -->
        <template v-if="isType('mutex')">
          <CollapsibleSection title="互斥体属性">
          <div class="prop-row">
            <span class="prop-label">名称</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 信号量 -->
        <template v-if="isType('semaphore')">
          <CollapsibleSection title="信号量属性">
          <div class="prop-row">
            <span class="prop-label">初始计数</span>
            <el-input-number v-model="localValue" size="small" :min="0" :max="99999" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">最大计数</span>
            <el-input-number v-model="localMax" size="small" :min="1" :max="99999" controls-position="right" @change="emitChange('max', localMax)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 事件对象 -->
        <template v-if="isType('event')">
          <CollapsibleSection title="事件对象属性">
          <div class="prop-row">
            <span class="prop-label">名称</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">手动重置</span>
            <el-switch v-model="localChecked" size="small" @change="emitChange('checked', localChecked)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 内存映射 -->
        <template v-if="isType('memorymap')">
          <CollapsibleSection title="内存映射属性">
          <div class="prop-row">
            <span class="prop-label">映射名称</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">大小(字节)</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="999999999" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 管道 -->
        <template v-if="isType('pipe')">
          <CollapsibleSection title="管道属性">
          <div class="prop-row">
            <span class="prop-label">管道名称</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 注册表 -->
        <template v-if="isType('registry')">
          <CollapsibleSection title="注册表属性">
          <div class="prop-row">
            <span class="prop-label">根键</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="HKEY_CLASSES_ROOT" value="HKCR" />
              <el-option label="HKEY_CURRENT_USER" value="HKCU" />
              <el-option label="HKEY_LOCAL_MACHINE" value="HKLM" />
              <el-option label="HKEY_USERS" value="HKU" />
              <el-option label="HKEY_CURRENT_CONFIG" value="HKCC" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- INI文件 -->
        <template v-if="isType('inifile')">
          <CollapsibleSection title="INI文件属性">
          <div class="prop-row">
            <span class="prop-label">文件路径</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 剪贴板 -->
        <template v-if="isType('clipboard')">
          <CollapsibleSection title="剪贴板属性">
          <div class="prop-row">
            <span class="prop-label">数据格式</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="文本" value="text" />
              <el-option label="图片" value="image" />
              <el-option label="文件" value="file" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 屏幕截图 -->
        <template v-if="isType('screencapture')">
          <CollapsibleSection title="屏幕截图属性">
          <div class="prop-row">
            <span class="prop-label">保存路径</span>
            <el-input v-model="localSrc" size="small" @change="emitChange('src', localSrc)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 语音对象 -->
        <template v-if="isType('speech')">
          <CollapsibleSection title="语音对象属性">
          <div class="prop-row">
            <span class="prop-label">语速</span>
            <el-input-number v-model="localValue" size="small" :min="-10" :max="10" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">音量</span>
            <el-input-number v-model="localMax" size="small" :min="0" :max="100" controls-position="right" @change="emitChange('max', localMax)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 摄像头 -->
        <template v-if="isType('camera')">
          <CollapsibleSection title="摄像头属性">
          <div class="prop-row">
            <span class="prop-label">设备索引</span>
            <el-input-number v-model="localValue" size="small" :min="0" :max="9" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">分辨率</span>
            <el-select v-model="localFormat" size="small" @change="emitChange('format', localFormat)">
              <el-option label="640x480" value="640x480" />
              <el-option label="1280x720" value="1280x720" />
              <el-option label="1920x1080" value="1920x1080" />
            </el-select>
          </div>
        </CollapsibleSection>
        </template>

        <!-- 对话框 -->
        <template v-if="isType('dialog')">
          <CollapsibleSection title="对话框属性">
          <div class="prop-row">
            <span class="prop-label">标题</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">过滤器</span>
            <el-input v-model="localItems" size="small" @change="emitChange('items', localItems)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 字体对话框 -->
        <template v-if="isType('fontdialog')">
          <CollapsibleSection title="字体对话框属性">
          <div class="prop-row">
            <span class="prop-label">默认字体</span>
            <el-input v-model="localFontName" size="small" @change="emitChange('fontName', localFontName)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">默认大小</span>
            <el-input-number v-model="localFontSize" size="small" :min="8" :max="72" controls-position="right" @change="emitChange('fontSize', localFontSize)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 颜色对话框 -->
        <template v-if="isType('colordialog')">
          <CollapsibleSection title="颜色对话框属性">
          <div class="prop-row">
            <span class="prop-label">默认颜色</span>
            <el-color-picker v-model="localColor" size="small" @change="emitChange('color', localColor)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 打印对话框 -->
        <template v-if="isType('printdialog')">
          <CollapsibleSection title="打印对话框属性">
          <div class="prop-row">
            <span class="prop-label">份数</span>
            <el-input-number v-model="localValue" size="small" :min="1" :max="99" controls-position="right" @change="emitChange('value', localValue)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 查找/替换对话框 -->
        <template v-if="isType('finddialog', 'replacedialog')">
          <CollapsibleSection title="查找对话框属性">
          <div class="prop-row">
            <span class="prop-label">查找文本</span>
            <el-input v-model="localText" size="small" @change="emitChange('text', localText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 月历 -->
        <template v-if="isType('monthcalendar')">
          <CollapsibleSection title="月历属性">
          <div class="prop-row">
            <span class="prop-label">显示周数</span>
            <el-switch v-model="localShowText" size="small" @change="emitChange('showText', localShowText)" />
          </div>
        </CollapsibleSection>
        </template>

        <!-- 通用外观属性 -->
        <template v-if="isType('button', 'label', 'input', 'textarea', 'checkbox', 'radio', 'combo', 'listbox', 'number', 'spin', 'datepicker', 'timepicker', 'colorpicker', 'filepicker', 'hotkey', 'ipedit', 'maskedit', 'memo', 'checklistbox', 'listview', 'progress', 'slider', 'scrollbar', 'group', 'panel', 'container', 'scrollbox', 'image', 'link', 'chart', 'barcode', 'qrcode', 'animate', 'mediaplayer', 'webview', 'shape', 'bevel', 'divider', 'splitter', 'tab', 'table', 'tree', 'menu', 'toolbar', 'statusbar', 'monthcalendar')">
          <CollapsibleSection title="外观属性">
          <div class="prop-row">
            <span class="prop-label">字体</span>
            <el-input v-model="localFontName" size="small" @change="emitChange('fontName', localFontName)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">粗体</span>
            <el-switch v-model="localFontBold" size="small" @change="emitChange('fontBold', localFontBold)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">斜体</span>
            <el-switch v-model="localFontItalic" size="small" @change="emitChange('fontItalic', localFontItalic)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">下划线</span>
            <el-switch v-model="localFontUnderline" size="small" @change="emitChange('fontUnderline', localFontUnderline)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">提示文本</span>
            <el-input v-model="localHint" size="small" @change="emitChange('hint', localHint)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">光标样式</span>
            <el-select v-model="localCursor" size="small" @change="emitChange('cursor', localCursor)">
              <el-option label="默认" value="default" />
              <el-option label="手型" value="hand" />
              <el-option label="十字" value="cross" />
              <el-option label="等待" value="wait" />
              <el-option label="文本" value="text" />
              <el-option label="移动" value="move" />
            </el-select>
          </div>
          <div class="prop-row">
            <span class="prop-label">自动大小</span>
            <el-switch v-model="localAutoSize" size="small" @change="emitChange('autoSize', localAutoSize)" />
          </div>
          <div class="prop-row">
            <span class="prop-label">Tab键停止</span>
            <el-switch v-model="localTabStop" size="small" @change="emitChange('tabStop', localTabStop)" />
          </div>
        </CollapsibleSection>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.property-panel {
  width: 240px;
  height: 100%;
  background-color: var(--panel-bg);
  border-left: 1px solid var(--panel-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.panel-header {
  padding: 8px 10px;
  font-size: 13px;
  font-weight: 600;
  color: var(--panel-header-text);
  border-bottom: 1px solid var(--panel-border);
  background-color: var(--panel-header-bg);
  flex-shrink: 0;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 6px;
}

.prop-row {
  display: flex;
  align-items: center;
  padding: 3px 4px;
  gap: 6px;
}

.prop-row-multi {
  align-items: flex-start;
}

.prop-label {
  width: 95px;
  font-size: 12px;
  color: var(--text-regular);
  flex-shrink: 0;
  text-align: right;
  line-height: 28px;
}

.prop-row-multi .prop-label {
  line-height: 20px;
  padding-top: 4px;
}

.prop-row :deep(.el-input),
.prop-row :deep(.el-select),
.prop-row :deep(.el-input-number) {
  flex: 1;
}

.prop-row :deep(.el-color-picker) {
  flex: 1;
}

.prop-row :deep(.el-switch) {
  margin-left: 0;
}

.prop-row :deep(.el-slider) {
  flex: 1;
  margin-right: 8px;
}
</style>