<script lang="ts" setup>
import { ref, reactive, onMounted, onUnmounted, computed } from 'vue'

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
  zIndex: number
  enabled: boolean
  fontName: string
  fontBold: boolean
  fontItalic: boolean
  fontUnderline: boolean
  hint: string
  cursor: string
  align: string
  anchors: { left: boolean; top: boolean; right: boolean; bottom: boolean }
  padding: { left: number; top: number; right: number; bottom: number }
  margin: { left: number; top: number; right: number; bottom: number }
  autoSize: boolean
  tabStop: boolean
  readOnly: boolean
  passwordMode: boolean
  passwordChar: string
  maxLength: number
  multiLine: boolean
  numberMode: boolean
  checked: boolean
  items: string
  selectedIndex: number
  dropDownStyle: string
  min: number
  max: number
  position: number
  orientation: string
  rowCount: number
  colCount: number
  fixedRows: number
  fixedCols: number
  picturePath: string
  stretch: boolean
  center: boolean
  proportional: boolean
  interval: number
  url: string
  shapeType: string
  penColor: string
  penWidth: number
  brushColor: string
  events: { name: string; handler: string }[]
}

interface HistoryEntry {
  components: CanvasComponent[]
}

const canvasRef = ref<HTMLDivElement>()
const components = reactive<CanvasComponent[]>([])
const selectedIds = ref<Set<string>>(new Set())
const nextId = ref(1)
const showGrid = ref(true)
const snapToGrid = ref(true)
const gridSize = 20
const windowTitle = ref('我的窗口')
const windowWidth = ref(800)
const windowHeight = ref(600)
const windowBgColor = ref('#f0f0f0')
const showMenuBar = ref(false)
const showMinBtn = ref(true)
const showMaxBtn = ref(true)
const showCloseBtn = ref(true)
const titleBarHeight = ref(32)

const history = ref<HistoryEntry[]>([])
const historyIndex = ref(-1)
const maxHistory = 50

const clipboard = ref<CanvasComponent[]>([])

const isDragging = ref(false)
const isResizing = ref(false)
const isBoxSelecting = ref(false)
const boxSelection = ref({ x: 0, y: 0, width: 0, height: 0 })
const dragStart = ref({ x: 0, y: 0 })
const compStartPositions = ref<Map<string, { x: number; y: number; width: number; height: number }>>(new Map())
const resizeHandle = ref('')
const alignmentGuides = ref<{ x: number[]; y: number[] }>({ x: [], y: [] })

const emit = defineEmits<{
  (e: 'component-focus', component: CanvasComponent | null): void
  (e: 'components-update', components: CanvasComponent[]): void
}>()

const windowFrameStyle = computed(() => ({
  width: windowWidth.value + 'px',
  height: windowHeight.value + 'px',
  backgroundColor: windowBgColor.value,
}))

const getComponentStyle = (comp: CanvasComponent) => ({
  left: comp.x + 'px',
  top: comp.y + 'px',
  width: comp.width + 'px',
  height: comp.height + 'px',
  borderColor: comp.color,
  color: comp.color,
  fontSize: comp.fontSize + 'px',
  display: comp.visible ? 'flex' : 'none',
  zIndex: comp.zIndex || 0,
})

const onWindowSizeChange = () => {
  emit('components-update', [...components])
}

const onTitleFocus = (e: FocusEvent) => {
  const target = e.target as HTMLInputElement
  target?.select()
}

const snapToGridValue = (value: number): number => {
  if (!snapToGrid.value) return value
  return Math.round(value / gridSize) * gridSize
}

const saveHistory = () => {
  const entry: HistoryEntry = {
    components: components.map(c => JSON.parse(JSON.stringify(c))),
  }
  history.value = history.value.slice(0, historyIndex.value + 1)
  history.value.push(entry)
  if (history.value.length > maxHistory) {
    history.value.shift()
  }
  historyIndex.value = history.value.length - 1
}

const undo = () => {
  if (historyIndex.value <= 0) return
  historyIndex.value--
  restoreHistory()
}

const redo = () => {
  if (historyIndex.value >= history.value.length - 1) return
  historyIndex.value++
  restoreHistory()
}

const restoreHistory = () => {
  const entry = history.value[historyIndex.value]
  components.splice(0, components.length)
  entry.components.forEach(c => components.push({ ...c }))
  emit('components-update', [...components])
}

const onDragOver = (event: DragEvent) => {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'copy'
  }
}

const onDrop = (event: DragEvent) => {
  event.preventDefault()
  if (!event.dataTransfer || !canvasRef.value) return

  const data = event.dataTransfer.getData('application/json')
  if (!data) return

  const compData = JSON.parse(data)
  const rect = canvasRef.value.getBoundingClientRect()
  const scrollLeft = canvasRef.value.scrollLeft
  const scrollTop = canvasRef.value.scrollTop
  const x = snapToGridValue(event.clientX - rect.left + scrollLeft - 50)
  const y = snapToGridValue(event.clientY - rect.top + scrollTop - 16)

  saveHistory()

  const maxZ = components.reduce((max, c) => Math.max(max, c.zIndex || 0), 0)

  const newComp: CanvasComponent = {
    id: `comp_${nextId.value++}`,
    type: compData.type,
    name: compData.name,
    x: Math.max(0, x),
    y: Math.max(0, y),
    width: 100,
    height: 32,
    text: compData.name,
    visible: true,
    color: '#409eff',
    fontSize: 14,
    zIndex: maxZ + 1,
    enabled: true,
    fontName: '微软雅黑',
    fontBold: false,
    fontItalic: false,
    fontUnderline: false,
    hint: '',
    cursor: 'default',
    align: 'left',
    anchors: { left: true, top: true, right: false, bottom: false },
    padding: { left: 0, top: 0, right: 0, bottom: 0 },
    margin: { left: 0, top: 0, right: 0, bottom: 0 },
    autoSize: false,
    tabStop: true,
    readOnly: false,
    passwordMode: false,
    passwordChar: '*',
    maxLength: 0,
    multiLine: false,
    numberMode: false,
    checked: false,
    items: '',
    selectedIndex: -1,
    dropDownStyle: 'dropdown',
    min: 0,
    max: 100,
    position: 0,
    orientation: 'horizontal',
    rowCount: 2,
    colCount: 2,
    fixedRows: 0,
    fixedCols: 0,
    picturePath: '',
    stretch: false,
    center: false,
    proportional: false,
    interval: 1000,
    url: '',
    shapeType: 'rectangle',
    penColor: '#000000',
    penWidth: 1,
    brushColor: '#ffffff',
    events: [],
  }

  components.push(newComp)
  selectedIds.value = new Set([newComp.id])
  emit('component-focus', newComp)
  emit('components-update', [...components])
}

const selectComponent = (comp: CanvasComponent, event: MouseEvent) => {
  if (event.ctrlKey || event.metaKey) {
    const newSet = new Set(selectedIds.value)
    if (newSet.has(comp.id)) {
      newSet.delete(comp.id)
    } else {
      newSet.add(comp.id)
    }
    selectedIds.value = newSet
  } else {
    selectedIds.value = new Set([comp.id])
  }
  emit('component-focus', comp)
}

const clearSelection = () => {
  if (isBoxSelecting.value) return
  selectedIds.value = new Set()
  emit('component-focus', null)
}

const startBoxSelect = (event: MouseEvent) => {
  if (event.button !== 0) return
  if ((event.target as HTMLElement).closest('.canvas-component')) return

  const rect = canvasRef.value?.getBoundingClientRect()
  if (!rect) return

  isBoxSelecting.value = true
  dragStart.value = {
    x: event.clientX - rect.left,
    y: event.clientY - rect.top,
  }
  boxSelection.value = { x: dragStart.value.x, y: dragStart.value.y, width: 0, height: 0 }
}

const onBoxSelectMove = (event: MouseEvent) => {
  if (!isBoxSelecting.value) return
  const rect = canvasRef.value?.getBoundingClientRect()
  if (!rect) return

  const currentX = event.clientX - rect.left
  const currentY = event.clientY - rect.top

  boxSelection.value = {
    x: Math.min(dragStart.value.x, currentX),
    y: Math.min(dragStart.value.y, currentY),
    width: Math.abs(currentX - dragStart.value.x),
    height: Math.abs(currentY - dragStart.value.y),
  }
}

const endBoxSelect = () => {
  if (!isBoxSelecting.value) return

  const wasClick = boxSelection.value.width <= 3 && boxSelection.value.height <= 3

  if (wasClick) {
    selectedIds.value = new Set()
    emit('component-focus', null)
  } else if (boxSelection.value.width > 3 && boxSelection.value.height > 3) {
    const box = boxSelection.value
    const newSelected = new Set<string>()
    components.forEach(comp => {
      const compCenterX = comp.x + comp.width / 2
      const compCenterY = comp.y + comp.height / 2
      if (
        compCenterX >= box.x &&
        compCenterX <= box.x + box.width &&
        compCenterY >= box.y &&
        compCenterY <= box.y + box.height
      ) {
        newSelected.add(comp.id)
      }
    })
    selectedIds.value = newSelected
    if (newSelected.size > 0) {
      emit('component-focus', components.find(c => c.id === [...newSelected][0]) || null)
    } else {
      emit('component-focus', null)
    }
  }

  isBoxSelecting.value = false
  boxSelection.value = { x: 0, y: 0, width: 0, height: 0 }
}

const startDrag = (event: MouseEvent, comp: CanvasComponent) => {
  if (event.button !== 0) return
  if ((event.target as HTMLElement).classList.contains('resize-handle')) return

  if (!selectedIds.value.has(comp.id)) {
    selectComponent(comp, event)
  }

  isDragging.value = true
  dragStart.value = { x: event.clientX, y: event.clientY }

  const positions = new Map<string, { x: number; y: number; width: number; height: number }>()
  selectedIds.value.forEach(id => {
    const c = components.find(c => c.id === id)
    if (c) {
      positions.set(id, { x: c.x, y: c.y, width: c.width, height: c.height })
    }
  })
  compStartPositions.value = positions

  document.addEventListener('mousemove', onDragMove)
  document.addEventListener('mouseup', onDragEnd)
}

const onDragMove = (event: MouseEvent) => {
  if (!isDragging.value) return

  const dx = event.clientX - dragStart.value.x
  const dy = event.clientY - dragStart.value.y

  selectedIds.value.forEach(id => {
    const comp = components.find(c => c.id === id)
    const start = compStartPositions.value.get(id)
    if (comp && start) {
      comp.x = snapToGridValue(Math.max(0, start.x + dx))
      comp.y = snapToGridValue(Math.max(0, start.y + dy))
    }
  })

  calculateAlignmentGuides()
}

const onDragEnd = () => {
  if (isDragging.value) {
    saveHistory()
    emit('components-update', [...components])
  }
  isDragging.value = false
  alignmentGuides.value = { x: [], y: [] }
  document.removeEventListener('mousemove', onDragMove)
  document.removeEventListener('mouseup', onDragEnd)
}

const startResize = (event: MouseEvent, comp: CanvasComponent, handle: string) => {
  event.stopPropagation()
  event.preventDefault()

  isResizing.value = true
  resizeHandle.value = handle
  dragStart.value = { x: event.clientX, y: event.clientY }

  const positions = new Map<string, { x: number; y: number; width: number; height: number }>()
  positions.set(comp.id, { x: comp.x, y: comp.y, width: comp.width, height: comp.height })
  compStartPositions.value = positions

  document.addEventListener('mousemove', onResizeMove)
  document.addEventListener('mouseup', onResizeEnd)
}

const onResizeMove = (event: MouseEvent) => {
  if (!isResizing.value) return

  const dx = event.clientX - dragStart.value.x
  const dy = event.clientY - dragStart.value.y

  compStartPositions.value.forEach((start, id) => {
    const comp = components.find(c => c.id === id)
    if (!comp) return

    const minSize = 20

    switch (resizeHandle.value) {
      case 'se':
        comp.width = snapToGridValue(Math.max(minSize, start.width + dx))
        comp.height = snapToGridValue(Math.max(minSize, start.height + dy))
        break
      case 'sw':
        comp.width = snapToGridValue(Math.max(minSize, start.width - dx))
        comp.x = snapToGridValue(Math.max(0, start.x + dx))
        comp.height = snapToGridValue(Math.max(minSize, start.height + dy))
        break
      case 'ne':
        comp.width = snapToGridValue(Math.max(minSize, start.width + dx))
        comp.y = snapToGridValue(Math.max(0, start.y + dy))
        comp.height = snapToGridValue(Math.max(minSize, start.height - dy))
        break
      case 'nw':
        comp.width = snapToGridValue(Math.max(minSize, start.width - dx))
        comp.x = snapToGridValue(Math.max(0, start.x + dx))
        comp.y = snapToGridValue(Math.max(0, start.y + dy))
        comp.height = snapToGridValue(Math.max(minSize, start.height - dy))
        break
      case 'n':
        comp.y = snapToGridValue(Math.max(0, start.y + dy))
        comp.height = snapToGridValue(Math.max(minSize, start.height - dy))
        break
      case 's':
        comp.height = snapToGridValue(Math.max(minSize, start.height + dy))
        break
      case 'e':
        comp.width = snapToGridValue(Math.max(minSize, start.width + dx))
        break
      case 'w':
        comp.width = snapToGridValue(Math.max(minSize, start.width - dx))
        comp.x = snapToGridValue(Math.max(0, start.x + dx))
        break
    }
  })
}

const onResizeEnd = () => {
  if (isResizing.value) {
    saveHistory()
    emit('components-update', [...components])
  }
  isResizing.value = false
  resizeHandle.value = ''
  document.removeEventListener('mousemove', onResizeMove)
  document.removeEventListener('mouseup', onResizeEnd)
}

const calculateAlignmentGuides = () => {
  const guides = { x: [] as number[], y: [] as number[] }
  const selectedComps = components.filter(c => selectedIds.value.has(c.id))
  const otherComps = components.filter(c => !selectedIds.value.has(c.id))

  const threshold = 5

  selectedComps.forEach(sc => {
    const scCenterX = sc.x + sc.width / 2
    const scCenterY = sc.y + sc.height / 2
    const scRight = sc.x + sc.width
    const scBottom = sc.y + sc.height

    otherComps.forEach(oc => {
      const ocCenterX = oc.x + oc.width / 2
      const ocCenterY = oc.y + oc.height / 2
      const ocRight = oc.x + oc.width
      const ocBottom = oc.y + oc.height

      if (Math.abs(sc.x - oc.x) < threshold) guides.x.push(oc.x)
      if (Math.abs(scRight - ocRight) < threshold) guides.x.push(ocRight)
      if (Math.abs(scCenterX - ocCenterX) < threshold) guides.x.push(ocCenterX)
      if (Math.abs(sc.y - oc.y) < threshold) guides.y.push(oc.y)
      if (Math.abs(scBottom - ocBottom) < threshold) guides.y.push(ocBottom)
      if (Math.abs(scCenterY - ocCenterY) < threshold) guides.y.push(ocCenterY)
    })
  })

  alignmentGuides.value = guides
}

const deleteSelected = () => {
  if (selectedIds.value.size === 0) return
  saveHistory()
  selectedIds.value.forEach(id => {
    const idx = components.findIndex(c => c.id === id)
    if (idx !== -1) components.splice(idx, 1)
  })
  selectedIds.value = new Set()
  emit('component-focus', null)
  emit('components-update', [...components])
}

const copySelected = () => {
  clipboard.value = components
    .filter(c => selectedIds.value.has(c.id))
    .map(c => JSON.parse(JSON.stringify(c)))
}

const pasteComponents = () => {
  if (clipboard.value.length === 0) return
  saveHistory()

  const newIds: string[] = []
  const maxZ = components.reduce((max, c) => Math.max(max, c.zIndex || 0), 0)

  clipboard.value.forEach(c => {
    const newComp: CanvasComponent = {
      ...c,
      id: `comp_${nextId.value++}`,
      x: c.x + 20,
      y: c.y + 20,
      zIndex: maxZ + 1,
    }
    components.push(newComp)
    newIds.push(newComp.id)
  })

  selectedIds.value = new Set(newIds)
  emit('components-update', [...components])
}

const bringToFront = () => {
  const maxZ = components.reduce((max, c) => Math.max(max, c.zIndex || 0), 0)
  selectedIds.value.forEach(id => {
    const comp = components.find(c => c.id === id)
    if (comp) comp.zIndex = maxZ + 1
  })
  emit('components-update', [...components])
}

const sendToBack = () => {
  const minZ = components.reduce((min, c) => Math.min(min, c.zIndex || 0), 0)
  selectedIds.value.forEach(id => {
    const comp = components.find(c => c.id === id)
    if (comp) comp.zIndex = minZ - 1
  })
  emit('components-update', [...components])
}

const getSelectedComponents = () => {
  return components.filter(c => selectedIds.value.has(c.id))
}

const alignLeft = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const minX = Math.min(...selected.map(c => c.x))
  selected.forEach(c => { c.x = minX })
  emit('components-update', [...components])
}

const alignCenterH = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const bounds = getSelectionBounds(selected)
  const centerX = bounds.x + bounds.width / 2
  selected.forEach(c => { c.x = centerX - c.width / 2 })
  emit('components-update', [...components])
}

const alignRight = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const maxX = Math.max(...selected.map(c => c.x + c.width))
  selected.forEach(c => { c.x = maxX - c.width })
  emit('components-update', [...components])
}

const alignTop = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const minY = Math.min(...selected.map(c => c.y))
  selected.forEach(c => { c.y = minY })
  emit('components-update', [...components])
}

const alignCenterV = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const bounds = getSelectionBounds(selected)
  const centerY = bounds.y + bounds.height / 2
  selected.forEach(c => { c.y = centerY - c.height / 2 })
  emit('components-update', [...components])
}

const alignBottom = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const maxY = Math.max(...selected.map(c => c.y + c.height))
  selected.forEach(c => { c.y = maxY - c.height })
  emit('components-update', [...components])
}

const distributeH = () => {
  const selected = getSelectedComponents().sort((a, b) => a.x - b.x)
  if (selected.length < 3) return
  saveHistory()
  const first = selected[0]
  const last = selected[selected.length - 1]
  const totalWidth = last.x + last.width - first.x - first.width
  const gap = totalWidth / (selected.length - 1)
  let currentX = first.x + first.width + gap
  for (let i = 1; i < selected.length - 1; i++) {
    selected[i].x = currentX
    currentX += selected[i].width + gap
  }
  emit('components-update', [...components])
}

const distributeV = () => {
  const selected = getSelectedComponents().sort((a, b) => a.y - b.y)
  if (selected.length < 3) return
  saveHistory()
  const first = selected[0]
  const last = selected[selected.length - 1]
  const totalHeight = last.y + last.height - first.y - first.height
  const gap = totalHeight / (selected.length - 1)
  let currentY = first.y + first.height + gap
  for (let i = 1; i < selected.length - 1; i++) {
    selected[i].y = currentY
    currentY += selected[i].height + gap
  }
  emit('components-update', [...components])
}

const makeSameWidth = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const maxWidth = Math.max(...selected.map(c => c.width))
  selected.forEach(c => { c.width = maxWidth })
  emit('components-update', [...components])
}

const makeSameHeight = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const maxHeight = Math.max(...selected.map(c => c.height))
  selected.forEach(c => { c.height = maxHeight })
  emit('components-update', [...components])
}

const makeSameSize = () => {
  const selected = getSelectedComponents()
  if (selected.length < 2) return
  saveHistory()
  const maxWidth = Math.max(...selected.map(c => c.width))
  const maxHeight = Math.max(...selected.map(c => c.height))
  selected.forEach(c => { c.width = maxWidth; c.height = maxHeight })
  emit('components-update', [...components])
}

const getSelectionBounds = (selected: CanvasComponent[]) => {
  const xs = selected.map(c => c.x)
  const ys = selected.map(c => c.y)
  const right = selected.map(c => c.x + c.width)
  const bottom = selected.map(c => c.y + c.height)
  return {
    x: Math.min(...xs),
    y: Math.min(...ys),
    width: Math.max(...right) - Math.min(...xs),
    height: Math.max(...bottom) - Math.min(...ys),
  }
}

const updateComponentProperty = (id: string, key: string, value: any) => {
  const comp = components.find(c => c.id === id)
  if (comp) {
    (comp as any)[key] = value
    emit('components-update', [...components])
  }
}

const updateWindowProperty = (key: string, value: any) => {
  if (key === 'title') {
    windowTitle.value = value
  } else if (key === 'width') {
    windowWidth.value = value
  } else if (key === 'height') {
    windowHeight.value = value
  } else if (key === 'backgroundColor') {
    windowBgColor.value = value
  } else if (key === 'showMenuBar') {
    showMenuBar.value = value
  } else if (key === 'showMinBtn') {
    showMinBtn.value = value
  } else if (key === 'showMaxBtn') {
    showMaxBtn.value = value
  } else if (key === 'showCloseBtn') {
    showCloseBtn.value = value
  } else if (key === 'titleBarHeight') {
    titleBarHeight.value = value
  }
}

const getComponents = () => {
  return [...components]
}

const loadComponents = (comps: CanvasComponent[]) => {
  components.splice(0, components.length)
  comps.forEach(c => components.push({ ...c, zIndex: c.zIndex || 0 }))
  nextId.value = Math.max(...comps.map(c => {
    const num = parseInt(c.id.replace('comp_', ''))
    return isNaN(num) ? 0 : num
  }), 0) + 1
  saveHistory()
}

const clearCanvas = () => {
  saveHistory()
  components.splice(0, components.length)
  selectedIds.value = new Set()
  emit('component-focus', null)
  emit('components-update', [...components])
}

const addComponent = (compData: any) => {
  saveHistory()
  const maxZ = components.reduce((max, c) => Math.max(max, c.zIndex || 0), 0)
  const newComp: CanvasComponent = {
    id: `comp_${nextId.value++}`,
    type: compData.type,
    name: compData.name,
    x: 20,
    y: components.length * 40 + 20,
    width: 100,
    height: 32,
    text: compData.name,
    visible: true,
    color: '#409eff',
    fontSize: 14,
    zIndex: maxZ + 1,
    enabled: true,
    fontName: '微软雅黑',
    fontBold: false,
    fontItalic: false,
    fontUnderline: false,
    hint: '',
    cursor: '默认',
    align: '左对齐',
    anchors: { left: true, top: true, right: false, bottom: false },
    padding: { left: 0, top: 0, right: 0, bottom: 0 },
    margin: { left: 0, top: 0, right: 0, bottom: 0 },
    autoSize: false,
    tabStop: true,
    readOnly: false,
    passwordMode: false,
    passwordChar: '*',
    maxLength: 0,
    multiLine: false,
    numberMode: false,
    checked: false,
    items: '',
    selectedIndex: -1,
    dropDownStyle: '下拉列表',
    min: 0,
    max: 100,
    position: 0,
    orientation: '水平',
    rowCount: 5,
    colCount: 3,
    fixedRows: 0,
    fixedCols: 0,
    picturePath: '',
    stretch: false,
    center: false,
    proportional: false,
    interval: 1000,
    url: '',
    shapeType: '矩形',
    penColor: '#000000',
    penWidth: 1,
    brushColor: '#ffffff',
    events: [],
  }
  components.push(newComp)
  selectedIds.value = new Set([newComp.id])
  emit('component-focus', newComp)
  emit('components-update', [...components])
}

const onKeyDown = (event: KeyboardEvent) => {
  if (event.target !== document.body && (event.target as HTMLElement).tagName === 'INPUT') return
  if ((event.target as HTMLElement).isContentEditable) return
  const isTextArea = (event.target as HTMLElement).tagName === 'TEXTAREA'
  const isInput = (event.target as HTMLElement).tagName === 'INPUT'
  if (isTextArea || isInput) return

  if (event.ctrlKey || event.metaKey) {
    switch (event.key.toLowerCase()) {
      case 'z':
        event.preventDefault()
        if (event.shiftKey) redo()
        else undo()
        break
      case 'y':
        event.preventDefault()
        redo()
        break
      case 'c':
        if (selectedIds.value.size > 0 && canvasRef.value?.contains(event.target as Node)) {
          event.preventDefault()
          copySelected()
        }
        break
      case 'v':
        if (canvasRef.value?.contains(event.target as Node)) {
          event.preventDefault()
          pasteComponents()
        }
        break
      case 'a':
        if (canvasRef.value?.contains(event.target as Node)) {
          event.preventDefault()
          selectedIds.value = new Set(components.map(c => c.id))
        }
        break
    }
  } else if (event.key === 'Delete' || event.key === 'Backspace') {
    if (!isTextArea && !isInput) {
      deleteSelected()
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeyDown)
  document.addEventListener('mousemove', onBoxSelectMove)
  document.addEventListener('mouseup', endBoxSelect)
  saveHistory()
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeyDown)
  document.removeEventListener('mousemove', onDragMove)
  document.removeEventListener('mouseup', onDragEnd)
  document.removeEventListener('mousemove', onResizeMove)
  document.removeEventListener('mouseup', onResizeEnd)
  document.removeEventListener('mousemove', onBoxSelectMove)
  document.removeEventListener('mouseup', endBoxSelect)
})

defineExpose({
  components,
  selectedIds,
  updateComponentProperty,
  updateWindowProperty,
  deleteSelected,
  getComponents,
  loadComponents,
  clearCanvas,
  addComponent,
  undo,
  redo,
  getSelectedComponents: () => components.filter(c => selectedIds.value.has(c.id)),
})
</script>

<template>
  <div class="design-canvas-wrapper">
    <div class="canvas-toolbar">
      <el-button-group size="small">
        <el-button :type="showGrid ? 'primary' : 'default'" @click="showGrid = !showGrid" title="显示/隐藏网格">
          网格
        </el-button>
        <el-button :type="snapToGrid ? 'primary' : 'default'" @click="snapToGrid = !snapToGrid" title="吸附到网格">
          吸附
        </el-button>
      </el-button-group>
      <el-divider direction="vertical" />
      <el-button-group size="small">
        <el-button :disabled="historyIndex <= 0" @click="undo" title="撤销 (Ctrl+Z)">
          撤销
        </el-button>
        <el-button :disabled="historyIndex >= history.length - 1" @click="redo" title="重做 (Ctrl+Y)">
          重做
        </el-button>
      </el-button-group>
      <el-divider direction="vertical" />
      <el-button-group size="small">
        <el-button :disabled="selectedIds.size === 0" @click="bringToFront" title="置于顶层">
          置顶
        </el-button>
        <el-button :disabled="selectedIds.size === 0" @click="sendToBack" title="置于底层">
          置底
        </el-button>
      </el-button-group>

      <!-- 对齐工具组 -->
      <template v-if="selectedIds.size >= 2">
        <el-divider direction="vertical" />
        <div class="align-toolbar">
          <div class="align-group" title="水平对齐">
            <el-button size="small" :disabled="selectedIds.size < 2" @click="alignLeft" title="左对齐">⫷</el-button>
            <el-button size="small" :disabled="selectedIds.size < 2" @click="alignCenterH" title="水平居中">⋮</el-button>
            <el-button size="small" :disabled="selectedIds.size < 2" @click="alignRight" title="右对齐">⫸</el-button>
          </div>
          <div class="align-sep"></div>
          <div class="align-group" title="垂直对齐">
            <el-button size="small" :disabled="selectedIds.size < 2" @click="alignTop" title="顶对齐">⫷</el-button>
            <el-button size="small" :disabled="selectedIds.size < 2" @click="alignCenterV" title="垂直居中">⋮</el-button>
            <el-button size="small" :disabled="selectedIds.size < 2" @click="alignBottom" title="底对齐">⫸</el-button>
          </div>
          <div class="align-sep"></div>
          <div class="align-group" title="等间距分布">
            <el-button size="small" :disabled="selectedIds.size < 3" @click="distributeH" title="水平等间距">≡</el-button>
            <el-button size="small" :disabled="selectedIds.size < 3" @click="distributeV" title="垂直等间距">⋮⋮</el-button>
          </div>
          <div class="align-sep"></div>
          <div class="align-group" title="统一尺寸">
            <el-button size="small" :disabled="selectedIds.size < 2" @click="makeSameWidth" title="统一宽度">↔</el-button>
            <el-button size="small" :disabled="selectedIds.size < 2" @click="makeSameHeight" title="统一高度">↕</el-button>
            <el-button size="small" :disabled="selectedIds.size < 2" @click="makeSameSize" title="统一大小">⤢</el-button>
          </div>
        </div>
      </template>

      <el-divider direction="vertical" />
      <span class="selection-info" v-if="selectedIds.size > 0">已选中 {{ selectedIds.size }} 个组件</span>
      <el-button
        v-if="selectedIds.size > 0"
        size="small"
        type="danger"
        @click="deleteSelected"
        title="删除选中组件 (Delete)"
      >
        删除
      </el-button>
    </div>

    <div class="canvas-viewport">
      <div class="window-frame" :style="windowFrameStyle">
        <div class="window-titlebar" :style="{ height: titleBarHeight + 'px' }">
          <div class="window-titlebar-left">
            <span class="window-icon">🐕</span>
            <input
              class="window-title-input"
              v-model="windowTitle"
              @focus="onTitleFocus"
              spellcheck="false"
            />
          </div>
          <div class="window-titlebar-right">
            <span v-if="showMinBtn" class="window-btn minimize" title="最小化">─</span>
            <span v-if="showMaxBtn" class="window-btn maximize" title="最大化">□</span>
            <span v-if="showCloseBtn" class="window-btn close" title="关闭">✕</span>
          </div>
        </div>
        <div v-if="showMenuBar" class="window-menu-bar">
          <span class="menu-item">文件</span>
          <span class="menu-item">编辑</span>
          <span class="menu-item">视图</span>
          <span class="menu-item">帮助</span>
        </div>
        <div
          ref="canvasRef"
          class="design-canvas"
          @dragover="onDragOver"
          @drop="onDrop"
          @mousedown="startBoxSelect"
          @mousemove="onBoxSelectMove"
          @mouseup="endBoxSelect"
        >
          <div v-if="showGrid" class="canvas-grid"></div>

          <!-- 框选矩形 -->
          <div
            v-if="isBoxSelecting"
            class="selection-box"
            :style="{
              left: boxSelection.x + 'px',
              top: boxSelection.y + 'px',
              width: boxSelection.width + 'px',
              height: boxSelection.height + 'px',
            }"
          ></div>

          <div
            v-for="guide in alignmentGuides.x"
            :key="'gx-' + guide"
            class="alignment-guide vertical"
            :style="{ left: guide + 'px' }"
          />
          <div
            v-for="guide in alignmentGuides.y"
            :key="'gy-' + guide"
            class="alignment-guide horizontal"
            :style="{ top: guide + 'px' }"
          />

          <div
            v-for="comp in components"
            :key="comp.id"
            class="canvas-component"
            :class="{
              selected: selectedIds.has(comp.id),
              dragging: isDragging && selectedIds.has(comp.id),
            }"
            :style="getComponentStyle(comp)"
            @mousedown.stop="startDrag($event, comp)"
          >
            <template v-if="comp.type === 'button'">
              <button class="comp-render comp-btn" :style="{ backgroundColor: comp.color, fontSize: comp.fontSize + 'px' }">
                {{ comp.text }}
              </button>
            </template>
            <template v-else-if="comp.type === 'input'">
              <input class="comp-render comp-input" :value="comp.text" readonly :style="{ fontSize: comp.fontSize + 'px' }" />
            </template>
            <template v-else-if="comp.type === 'label'">
              <label class="comp-render comp-label-render" :style="{ color: comp.color, fontSize: comp.fontSize + 'px' }">{{ comp.text }}</label>
            </template>
            <template v-else-if="comp.type === 'checkbox'">
              <label class="comp-render comp-checkbox" :style="{ fontSize: comp.fontSize + 'px' }">
                <input type="checkbox" disabled /> {{ comp.text }}
              </label>
            </template>
            <template v-else-if="comp.type === 'combo'">
              <select class="comp-render comp-combo" :style="{ fontSize: comp.fontSize + 'px' }">
                <option>{{ comp.text }}</option>
              </select>
            </template>
            <template v-else-if="comp.type === 'radio'">
              <label class="comp-render comp-radio" :style="{ fontSize: comp.fontSize + 'px' }">
                <input type="radio" disabled /> {{ comp.text }}
              </label>
            </template>
            <template v-else-if="comp.type === 'textarea'">
              <textarea class="comp-render comp-textarea" readonly :style="{ fontSize: comp.fontSize + 'px' }">{{ comp.text }}</textarea>
            </template>
            <template v-else-if="comp.type === 'progress'">
              <div class="comp-render comp-progress">
                <div class="progress-fill" :style="{ width: '60%', backgroundColor: comp.color }"></div>
              </div>
            </template>
            <template v-else-if="comp.type === 'slider'">
              <input type="range" class="comp-render comp-slider" disabled />
            </template>
            <template v-else-if="comp.type === 'listbox'">
              <div class="comp-render comp-listbox" :style="{ borderColor: comp.color }">
                <div class="listbox-item">项目1</div>
                <div class="listbox-item">项目2</div>
                <div class="listbox-item">项目3</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'table'">
              <div class="comp-render comp-table" :style="{ borderColor: comp.color }">
                <div class="table-header">列1</div>
                <div class="table-header">列2</div>
                <div class="table-cell">-</div>
                <div class="table-cell">-</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'group'">
              <fieldset class="comp-render comp-group" :style="{ borderColor: comp.color }">
                <legend>{{ comp.text }}</legend>
              </fieldset>
            </template>
            <template v-else-if="comp.type === 'tab'">
              <div class="comp-render comp-tab">
                <div class="tab-header" :style="{ backgroundColor: comp.color }">选项卡1</div>
                <div class="tab-header">选项卡2</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'tree'">
              <div class="comp-render comp-tree">
                <div>📁 根节点</div>
                <div style="padding-left:12px">📄 子节点1</div>
                <div style="padding-left:12px">📄 子节点2</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'image'">
              <div class="comp-render comp-image" :style="{ borderColor: comp.color }">
                🖼 {{ comp.text }}
              </div>
            </template>
            <template v-else-if="comp.type === 'datepicker'">
              <input type="date" class="comp-render comp-datepicker" disabled />
            </template>
            <template v-else-if="comp.type === 'link'">
              <a class="comp-render comp-link" :style="{ color: comp.color, fontSize: comp.fontSize + 'px' }">{{ comp.text }}</a>
            </template>
            <template v-else-if="comp.type === 'divider'">
              <hr class="comp-render comp-divider" />
            </template>
            <template v-else-if="comp.type === 'statusbar'">
              <div class="comp-render comp-statusbar">{{ comp.text }}</div>
            </template>
            <template v-else-if="comp.type === 'menu'">
              <div class="comp-render comp-menu">
                <span class="menu-item-render">文件</span>
                <span class="menu-item-render">编辑</span>
                <span class="menu-item-render">帮助</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'toolbar'">
              <div class="comp-render comp-toolbar">
                <span class="toolbar-btn">🔧</span>
                <span class="toolbar-btn">💾</span>
                <span class="toolbar-btn">▶</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'container'">
              <div class="comp-render comp-container" :style="{ borderColor: comp.color }">
                {{ comp.text }}
              </div>
            </template>
            <template v-else-if="comp.type === 'text'">
              <span class="comp-render comp-text-render" :style="{ color: comp.color, fontSize: comp.fontSize + 'px' }">{{ comp.text }}</span>
            </template>
            <template v-else-if="comp.type === 'icon'">
              <span class="comp-render comp-icon-render" :style="{ fontSize: comp.fontSize + 'px' }">⭐</span>
            </template>
            <template v-else-if="comp.type === 'avatar'">
              <div class="comp-render comp-avatar" :style="{ backgroundColor: comp.color, fontSize: comp.fontSize + 'px' }">A</div>
            </template>
            <template v-else-if="comp.type === 'badge'">
              <span class="comp-render comp-badge" :style="{ backgroundColor: comp.color }">3</span>
            </template>
            <template v-else-if="comp.type === 'tag'">
              <span class="comp-render comp-tag-render" :style="{ backgroundColor: comp.color, fontSize: comp.fontSize + 'px' }">{{ comp.text }}</span>
            </template>
            <template v-else-if="comp.type === 'spinner'">
              <div class="comp-render comp-spinner"></div>
            </template>
            <template v-else-if="comp.type === 'skeleton'">
              <div class="comp-render comp-skeleton">
                <div class="skeleton-line"></div>
                <div class="skeleton-line short"></div>
              </div>
            </template>
            <template v-else-if="comp.type === 'tooltip'">
              <div class="comp-render comp-tooltip-render">💬 {{ comp.text }}</div>
            </template>
            <template v-else-if="comp.type === 'popover'">
              <div class="comp-render comp-popover-render">
                <span>📌</span>
                <div class="popover-bubble">{{ comp.text }}</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'switch'">
              <label class="comp-render comp-switch-render">
                <span class="switch-track" :style="{ backgroundColor: comp.color }">
                  <span class="switch-thumb"></span>
                </span>
              </label>
            </template>
            <template v-else-if="comp.type === 'number'">
              <input type="number" class="comp-render comp-number" value="0" readonly :style="{ fontSize: comp.fontSize + 'px' }" />
            </template>
            <template v-else-if="comp.type === 'select'">
              <select class="comp-render comp-select-render" :style="{ fontSize: comp.fontSize + 'px' }">
                <option>{{ comp.text }}</option>
              </select>
            </template>
            <template v-else-if="comp.type === 'upload'">
              <div class="comp-render comp-upload">📤 上传文件</div>
            </template>
            <template v-else-if="comp.type === 'color-picker'">
              <div class="comp-render comp-color-picker">
                <span class="color-swatch" :style="{ backgroundColor: comp.color }"></span>
                <span class="color-label">{{ comp.color }}</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'rate'">
              <div class="comp-render comp-rate">★★★★☆</div>
            </template>
            <template v-else-if="comp.type === 'transfer'">
              <div class="comp-render comp-transfer">
                <div class="transfer-panel">📋 源</div>
                <span>⇄</span>
                <div class="transfer-panel">📋 目标</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'cascader'">
              <div class="comp-render comp-cascader" :style="{ fontSize: comp.fontSize + 'px' }">
                <span>📍 {{ comp.text }}</span>
                <span>▾</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'time-picker'">
              <input type="time" class="comp-render comp-time-picker" disabled />
            </template>
            <template v-else-if="comp.type === 'form'">
              <div class="comp-render comp-form-render" :style="{ borderColor: comp.color }">
                <div class="form-row"><label>字段1</label><input disabled /></div>
                <div class="form-row"><label>字段2</label><input disabled /></div>
              </div>
            </template>
            <template v-else-if="comp.type === 'form-item'">
              <div class="comp-render comp-form-item">
                <label class="form-item-label">{{ comp.text }}</label>
                <input disabled class="form-item-input" />
              </div>
            </template>
            <template v-else-if="comp.type === 'calendar'">
              <div class="comp-render comp-calendar">
                <div class="calendar-header">📅 2026年5月</div>
                <div class="calendar-grid">
                  <span v-for="d in 7" :key="d" class="calendar-cell">{{ d }}</span>
                </div>
              </div>
            </template>
            <template v-else-if="comp.type === 'timeline'">
              <div class="comp-render comp-timeline">
                <div class="timeline-item"><span class="timeline-dot"></span>事件1</div>
                <div class="timeline-item"><span class="timeline-dot"></span>事件2</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'stat'">
              <div class="comp-render comp-stat">
                <span class="stat-value">1,024</span>
                <span class="stat-label">{{ comp.text }}</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'description'">
              <div class="comp-render comp-description">
                <div class="desc-row"><span class="desc-label">名称</span><span class="desc-value">值</span></div>
                <div class="desc-row"><span class="desc-label">状态</span><span class="desc-value">正常</span></div>
              </div>
            </template>
            <template v-else-if="comp.type === 'empty'">
              <div class="comp-render comp-empty">📭 暂无数据</div>
            </template>
            <template v-else-if="comp.type === 'result'">
              <div class="comp-render comp-result">
                <span class="result-icon">✅</span>
                <span>{{ comp.text }}</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'carousel'">
              <div class="comp-render comp-carousel">
                <span>◀</span>
                <span class="carousel-dot active"></span>
                <span class="carousel-dot"></span>
                <span class="carousel-dot"></span>
                <span>▶</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'collapse'">
              <div class="comp-render comp-collapse-render">
                <div class="collapse-header">▾ {{ comp.text }}</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'pagination'">
              <div class="comp-render comp-pagination">
                <span class="page-btn active">1</span>
                <span class="page-btn">2</span>
                <span class="page-btn">3</span>
                <span class="page-btn">⟩</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'steps'">
              <div class="comp-render comp-steps-render">
                <span class="step active">① 步骤1</span>
                <span class="step-arrow">→</span>
                <span class="step">② 步骤2</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'breadcrumb'">
              <div class="comp-render comp-breadcrumb">
                <span>首页</span><span>/</span><span>分类</span><span>/</span><span class="active">{{ comp.text }}</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'dropdown'">
              <div class="comp-render comp-dropdown-render">
                <span>{{ comp.text }}</span><span> ▾</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'anchor'">
              <div class="comp-render comp-anchor">
                <div class="anchor-item active">标题1</div>
                <div class="anchor-item">标题2</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'backtop'">
              <div class="comp-render comp-backtop">⬆</div>
            </template>
            <template v-else-if="comp.type === 'drawer'">
              <div class="comp-render comp-drawer-render" :style="{ borderColor: comp.color }">
                <div class="drawer-header">{{ comp.text }}</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'affix'">
              <div class="comp-render comp-affix">📌 {{ comp.text }}</div>
            </template>
            <template v-else-if="comp.type === 'alert'">
              <div class="comp-render comp-alert-render" :style="{ borderColor: comp.color }">
                <span>⚠</span><span>{{ comp.text }}</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'message'">
              <div class="comp-render comp-message-render" :style="{ backgroundColor: comp.color }">
                {{ comp.text }}
              </div>
            </template>
            <template v-else-if="comp.type === 'notification'">
              <div class="comp-render comp-notification">
                <span class="notif-icon">🔔</span>
                <div class="notif-body">
                  <div class="notif-title">{{ comp.text }}</div>
                  <div class="notif-desc">通知内容</div>
                </div>
              </div>
            </template>
            <template v-else-if="comp.type === 'dialog'">
              <div class="comp-render comp-dialog-render" :style="{ borderColor: comp.color }">
                <div class="dialog-title">{{ comp.text }}</div>
                <div class="dialog-btns">
                  <span class="dialog-btn">确定</span>
                  <span class="dialog-btn">取消</span>
                </div>
              </div>
            </template>
            <template v-else-if="comp.type === 'loading'">
              <div class="comp-render comp-loading">
                <span class="loading-spin"></span>
                <span>{{ comp.text }}</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'popconfirm'">
              <div class="comp-render comp-popconfirm">
                <span>❓ {{ comp.text }}</span>
                <span class="popconfirm-btns">确定/取消</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'watermark'">
              <div class="comp-render comp-watermark">{{ comp.text }}</div>
            </template>
            <template v-else-if="comp.type === 'tour'">
              <div class="comp-render comp-tour">
                <span>🚩 步骤 1/3</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'scrollbar'">
              <div class="comp-render comp-scrollbar">
                <div class="scrollbar-track">
                  <div class="scrollbar-thumb"></div>
                </div>
              </div>
            </template>
            <template v-else-if="comp.type === 'panel'">
              <div class="comp-render comp-panel-render" :style="{ borderColor: comp.color }">
                <div class="panel-title">{{ comp.text }}</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'splitter'">
              <div class="comp-render comp-splitter">
                <div class="splitter-pane">左</div>
                <div class="splitter-bar"></div>
                <div class="splitter-pane">右</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'scroll-view'">
              <div class="comp-render comp-scroll-view" :style="{ borderColor: comp.color }">
                <div class="scroll-content">{{ comp.text }}</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'dock-panel'">
              <div class="comp-render comp-dock-panel">
                <div class="dock-top">顶</div>
                <div class="dock-mid">
                  <div class="dock-left">左</div>
                  <div class="dock-center">中</div>
                  <div class="dock-right">右</div>
                </div>
                <div class="dock-bottom">底</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'card'">
              <div class="comp-render comp-card-render" :style="{ borderColor: comp.color }">
                <div class="card-header">{{ comp.text }}</div>
                <div class="card-body">卡片内容</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'line-chart'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 100 40">
                <polyline points="5,35 25,20 45,25 65,10 85,15 95,5" fill="none" :stroke="comp.color" stroke-width="2"/>
                <circle cx="5" cy="35" r="2" :fill="comp.color"/>
                <circle cx="25" cy="20" r="2" :fill="comp.color"/>
                <circle cx="45" cy="25" r="2" :fill="comp.color"/>
                <circle cx="65" cy="10" r="2" :fill="comp.color"/>
                <circle cx="85" cy="15" r="2" :fill="comp.color"/>
                <circle cx="95" cy="5" r="2" :fill="comp.color"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'bar-chart'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 100 40">
                <rect x="5" y="20" width="10" height="20" :fill="comp.color" rx="1"/>
                <rect x="20" y="10" width="10" height="30" :fill="comp.color" rx="1"/>
                <rect x="35" y="15" width="10" height="25" :fill="comp.color" rx="1"/>
                <rect x="50" y="5" width="10" height="35" :fill="comp.color" rx="1"/>
                <rect x="65" y="12" width="10" height="28" :fill="comp.color" rx="1"/>
                <rect x="80" y="8" width="10" height="32" :fill="comp.color" rx="1"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'pie-chart'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 40 40">
                <circle cx="20" cy="20" r="16" fill="none" :stroke="comp.color" stroke-width="16"
                  stroke-dasharray="30 70" stroke-dashoffset="25" transform="rotate(-90 20 20)"/>
                <circle cx="20" cy="20" r="16" fill="none" stroke="#e9ecef" stroke-width="16"
                  stroke-dasharray="70 30" stroke-dashoffset="5" transform="rotate(-90 20 20)"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'area-chart'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 100 40">
                <polygon points="5,40 5,30 25,15 45,20 65,8 85,12 95,5 95,40" :fill="comp.color" opacity="0.3"/>
                <polyline points="5,30 25,15 45,20 65,8 85,12 95,5" fill="none" :stroke="comp.color" stroke-width="2"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'scatter-chart'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 100 40">
                <circle cx="10" cy="30" r="3" :fill="comp.color"/>
                <circle cx="25" cy="20" r="3" :fill="comp.color"/>
                <circle cx="40" cy="25" r="3" :fill="comp.color"/>
                <circle cx="55" cy="10" r="3" :fill="comp.color"/>
                <circle cx="70" cy="15" r="3" :fill="comp.color"/>
                <circle cx="85" cy="8" r="3" :fill="comp.color"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'radar-chart'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 40 40">
                <polygon points="20,4 35,14 30,34 10,34 5,14" fill="none" :stroke="comp.color" stroke-width="1.5"/>
                <polygon points="20,10 30,17 26,30 14,30 10,17" :fill="comp.color" opacity="0.3"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'funnel-chart'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 60 40">
                <polygon points="5,2 55,2 45,14 15,14" :fill="comp.color" opacity="0.9"/>
                <polygon points="10,14 50,14 40,26 20,26" :fill="comp.color" opacity="0.6"/>
                <polygon points="15,26 45,26 35,38 25,38" :fill="comp.color" opacity="0.3"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'heatmap'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 60 40">
                <rect x="2" y="2" width="12" height="10" :fill="comp.color" opacity="0.3" rx="1"/>
                <rect x="16" y="2" width="12" height="10" :fill="comp.color" opacity="0.5" rx="1"/>
                <rect x="30" y="2" width="12" height="10" :fill="comp.color" opacity="0.7" rx="1"/>
                <rect x="44" y="2" width="12" height="10" :fill="comp.color" opacity="0.9" rx="1"/>
                <rect x="2" y="14" width="12" height="10" :fill="comp.color" opacity="0.4" rx="1"/>
                <rect x="16" y="14" width="12" height="10" :fill="comp.color" opacity="0.6" rx="1"/>
                <rect x="30" y="14" width="12" height="10" :fill="comp.color" opacity="0.8" rx="1"/>
                <rect x="44" y="14" width="12" height="10" :fill="comp.color" opacity="1" rx="1"/>
                <rect x="2" y="26" width="12" height="10" :fill="comp.color" opacity="0.2" rx="1"/>
                <rect x="16" y="26" width="12" height="10" :fill="comp.color" opacity="0.5" rx="1"/>
                <rect x="30" y="26" width="12" height="10" :fill="comp.color" opacity="0.7" rx="1"/>
                <rect x="44" y="26" width="12" height="10" :fill="comp.color" opacity="0.9" rx="1"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'gauge'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 40 30">
                <path d="M5,28 A18,18 0 0,1 35,28" fill="none" stroke="#e9ecef" stroke-width="6"/>
                <path d="M5,28 A18,18 0 0,1 25,6" fill="none" :stroke="comp.color" stroke-width="6" stroke-linecap="round"/>
                <circle cx="20" cy="22" r="2" :fill="comp.color"/>
              </svg>
            </template>
            <template v-else-if="comp.type === 'ring-progress'">
              <svg class="comp-render comp-chart-svg" viewBox="0 0 40 40">
                <circle cx="20" cy="20" r="14" fill="none" stroke="#e9ecef" stroke-width="5"/>
                <circle cx="20" cy="20" r="14" fill="none" :stroke="comp.color" stroke-width="5"
                  stroke-dasharray="55 33" stroke-linecap="round" transform="rotate(-90 20 20)"/>
                <text x="20" y="23" text-anchor="middle" font-size="8" fill="#333">62%</text>
              </svg>
            </template>
            <template v-else-if="comp.type === 'space'">
              <div class="comp-render comp-space-render"></div>
            </template>
            <template v-else-if="comp.type === 'heading'">
              <h3 class="comp-render comp-heading" :style="{ color: comp.color, fontSize: (comp.fontSize || 18) + 'px' }">{{ comp.text }}</h3>
            </template>
            <template v-else-if="comp.type === 'paragraph'">
              <p class="comp-render comp-paragraph" :style="{ fontSize: (comp.fontSize || 12) + 'px' }">{{ comp.text }}</p>
            </template>
            <template v-else-if="comp.type === 'blockquote'">
              <blockquote class="comp-render comp-blockquote" :style="{ borderColor: comp.color }">{{ comp.text }}</blockquote>
            </template>
            <template v-else-if="comp.type === 'code-block'">
              <pre class="comp-render comp-code-block"><code>{{ comp.text }}</code></pre>
            </template>
            <template v-else-if="comp.type === 'video'">
              <div class="comp-render comp-video-render">🎬 {{ comp.text }}</div>
            </template>
            <template v-else-if="comp.type === 'input-number'">
              <input type="number" class="comp-render comp-number" value="0" readonly :style="{ fontSize: comp.fontSize + 'px' }" />
            </template>
            <template v-else-if="comp.type === 'date-picker'">
              <input type="date" class="comp-render comp-datepicker" disabled />
            </template>
            <template v-else-if="comp.type === 'autocomplete'">
              <div class="comp-render comp-autocomplete">
                <input disabled :value="comp.text" :style="{ fontSize: comp.fontSize + 'px' }" />
                <div class="autocomplete-drop">🔍 匹配结果...</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'rich-text'">
              <div class="comp-render comp-rich-text" :style="{ borderColor: comp.color }">
                <span class="rt-bold">B</span>
                <span class="rt-italic">I</span>
                <span class="rt-underline">U</span>
                <div class="rt-content">{{ comp.text }}</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'statistic'">
              <div class="comp-render comp-stat">
                <span class="stat-value">1,024</span>
                <span class="stat-label">{{ comp.text }}</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'list'">
              <div class="comp-render comp-list-render">
                <div class="list-item-render">• 项目 1</div>
                <div class="list-item-render">• 项目 2</div>
                <div class="list-item-render">• 项目 3</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'image-viewer'">
              <div class="comp-render comp-image-viewer">
                <span>🖼</span>
                <span class="viewer-nav">◀ 1/3 ▶</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'virtual-list'">
              <div class="comp-render comp-virtual-list">
                <div class="vlist-item">行 1</div>
                <div class="vlist-item">行 2</div>
                <div class="vlist-item">行 3</div>
                <div class="vlist-more">...</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'descriptions'">
              <div class="comp-render comp-description">
                <div class="desc-row"><span class="desc-label">名称</span><span class="desc-value">值</span></div>
                <div class="desc-row"><span class="desc-label">状态</span><span class="desc-value">正常</span></div>
              </div>
            </template>
            <template v-else-if="comp.type === 'tabs'">
              <div class="comp-render comp-tab">
                <div class="tab-header" :style="{ backgroundColor: comp.color }">选项卡1</div>
                <div class="tab-header">选项卡2</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'navbar'">
              <div class="comp-render comp-navbar">
                <span class="navbar-brand">🏠</span>
                <span class="navbar-item">导航1</span>
                <span class="navbar-item">导航2</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'progress-bar'">
              <div class="comp-render comp-progress">
                <div class="progress-fill" :style="{ width: '60%', backgroundColor: comp.color }"></div>
              </div>
            </template>
            <template v-else-if="comp.type === 'confirm'">
              <div class="comp-render comp-popconfirm">
                <span>❓ {{ comp.text }}</span>
                <span class="popconfirm-btns">确定/取消</span>
              </div>
            </template>
            <template v-else-if="comp.type === 'layout'">
              <div class="comp-render comp-layout-render" :style="{ borderColor: comp.color }">
                <div class="layout-header">Header</div>
                <div class="layout-body">
                  <div class="layout-side">Side</div>
                  <div class="layout-main">Content</div>
                </div>
                <div class="layout-footer">Footer</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'row'">
              <div class="comp-render comp-row-render">
                <div class="row-col">Col</div>
                <div class="row-col">Col</div>
                <div class="row-col">Col</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'col'">
              <div class="comp-render comp-col-render">
                <div class="col-row">Row</div>
                <div class="col-row">Row</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'flex'">
              <div class="comp-render comp-flex-render">
                <div class="flex-item">A</div>
                <div class="flex-item">B</div>
                <div class="flex-item">C</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'split-pane'">
              <div class="comp-render comp-splitter">
                <div class="splitter-pane">左</div>
                <div class="splitter-bar"></div>
                <div class="splitter-pane">右</div>
              </div>
            </template>
            <template v-else-if="comp.type === 'tab-container'">
              <div class="comp-render comp-tab-container">
                <div class="tab-header-row">
                  <span class="tab-h active">Tab1</span>
                  <span class="tab-h">Tab2</span>
                </div>
                <div class="tab-body">内容区域</div>
              </div>
            </template>
            <template v-else>
              <span class="comp-type-badge">{{ comp.type }}</span>
              <span class="comp-label">{{ comp.text }}</span>
            </template>

            <div v-if="selectedIds.has(comp.id)" class="resize-handles">
              <div class="resize-handle nw" @mousedown.stop="startResize($event, comp, 'nw')" />
              <div class="resize-handle n" @mousedown.stop="startResize($event, comp, 'n')" />
              <div class="resize-handle ne" @mousedown.stop="startResize($event, comp, 'ne')" />
              <div class="resize-handle e" @mousedown.stop="startResize($event, comp, 'e')" />
              <div class="resize-handle se" @mousedown.stop="startResize($event, comp, 'se')" />
              <div class="resize-handle s" @mousedown.stop="startResize($event, comp, 's')" />
              <div class="resize-handle sw" @mousedown.stop="startResize($event, comp, 'sw')" />
              <div class="resize-handle w" @mousedown.stop="startResize($event, comp, 'w')" />
            </div>
          </div>

          <div v-if="components.length === 0" class="canvas-placeholder">
            <el-icon :size="48" color="#c0c4cc"><Plus /></el-icon>
            <p>从左侧组件库拖拽组件到此处</p>
            <p class="hint">或点击组件库中的组件添加</p>
          </div>
        </div>
        <div class="window-statusbar">
          <span>{{ windowTitle }} - {{ windowWidth }}×{{ windowHeight }}</span>
          <span>{{ components.length }} 个组件</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.design-canvas-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.canvas-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: var(--panel-header-bg);
  border-bottom: 1px solid var(--panel-border);
  flex-shrink: 0;
  flex-wrap: wrap;
}

.align-toolbar {
  display: flex;
  align-items: center;
  gap: 2px;
}

.align-group {
  display: flex;
  align-items: center;
  gap: 1px;
  background: var(--bg-hover);
  border-radius: 4px;
  padding: 2px;
}

.align-group .el-button {
  padding: 4px 6px !important;
  font-size: 13px !important;
  min-width: 26px;
}

.align-sep {
  width: 1px;
  height: 20px;
  background: var(--border-light);
  margin: 0 4px;
}

.selection-info {
  font-size: 12px;
  color: var(--text-secondary);
  margin-left: 8px;
}

.canvas-viewport {
  flex: 1;
  overflow: auto;
  background-color: var(--bg-canvas);
  background-image:
    linear-gradient(45deg, var(--border-light) 25%, transparent 25%),
    linear-gradient(-45deg, var(--border-light) 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, var(--border-light) 75%),
    linear-gradient(-45deg, transparent 75%, var(--border-light) 75%);
  background-size: 16px 16px;
  background-position: 0 0, 0 8px, 8px -8px, -8px 0;
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  padding: 24px;
}

.window-frame {
  background-color: #f0f0f0;
  border: 1px solid #0078d4;
  border-radius: 8px 8px 0 0;
  box-shadow:
    0 0 0 1px rgba(0, 0, 0, 0.1),
    0 8px 32px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow: hidden;
  min-width: 200px;
  min-height: 150px;
}

.window-titlebar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 32px;
  background: linear-gradient(180deg, #0078d4 0%, #005a9e 100%);
  padding: 0 8px;
  flex-shrink: 0;
}

.window-titlebar-left {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.window-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.window-title-input {
  background: transparent;
  border: 1px solid transparent;
  color: #fff;
  font-size: 12px;
  font-family: 'Microsoft YaHei', sans-serif;
  padding: 2px 4px;
  border-radius: 2px;
  outline: none;
  width: 100%;
  max-width: 200px;
}

.window-title-input:hover {
  border-color: rgba(255, 255, 255, 0.3);
}

.window-title-input:focus {
  border-color: rgba(255, 255, 255, 0.6);
  background-color: rgba(255, 255, 255, 0.1);
}

.window-titlebar-right {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
}

.window-btn {
  width: 28px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  color: #fff;
  border-radius: 2px;
  cursor: pointer;
  user-select: none;
}

.window-btn:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.window-btn.close:hover {
  background-color: #e81123;
}

.window-menu-bar {
  display: flex;
  align-items: center;
  height: 24px;
  background-color: #f0f0f0;
  border-bottom: 1px solid #d0d0d0;
  padding: 0 4px;
  gap: 2px;
  flex-shrink: 0;
}

.menu-item {
  font-size: 12px;
  color: #333;
  padding: 2px 8px;
  border-radius: 2px;
  cursor: default;
  user-select: none;
}

.menu-item:hover {
  background-color: #d0d0d0;
}

.design-canvas {
  flex: 1;
  background-color: #fff;
  position: relative;
  overflow: hidden;
  min-height: 100px;
}

.canvas-grid {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image:
    linear-gradient(rgba(0, 0, 0, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 0, 0, 0.04) 1px, transparent 1px);
  background-size: 20px 20px;
  pointer-events: none;
}

.alignment-guide {
  position: absolute;
  pointer-events: none;
  z-index: 1000;
}

.alignment-guide.vertical {
  top: 0;
  bottom: 0;
  width: 1px;
  background-color: #ff4d4f;
}

.alignment-guide.horizontal {
  left: 0;
  right: 0;
  height: 1px;
  background-color: #ff4d4f;
}

.canvas-component {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px dashed rgba(64, 158, 255, 0.4);
  border-radius: 3px;
  background-color: transparent;
  cursor: move;
  transition: border-color 0.15s, box-shadow 0.15s;
  user-select: none;
  min-width: 20px;
  min-height: 20px;
  overflow: visible;
}

.canvas-component:hover {
  border-color: rgba(64, 158, 255, 0.7);
}

.canvas-component.selected {
  border: 2px solid #409eff !important;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.25), 0 2px 8px rgba(64, 158, 255, 0.2);
  z-index: 100;
}

.canvas-component.dragging {
  opacity: 0.9;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.25);
}

.comp-type-badge {
  position: absolute;
  top: -10px;
  left: 4px;
  font-size: 10px;
  background-color: #409eff;
  color: #fff;
  padding: 0 4px;
  border-radius: 2px;
  line-height: 16px;
  white-space: nowrap;
  z-index: 1;
}

.comp-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 4px;
}

.comp-render {
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  pointer-events: none;
}

.comp-btn {
  background-color: #409eff;
  color: #fff;
  border-radius: 3px;
  cursor: default;
  font-family: 'Microsoft YaHei', sans-serif;
  padding: 0 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.comp-input {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 8px;
  background-color: #fff;
  font-family: 'Microsoft YaHei', sans-serif;
  color: #333;
}

.comp-label-render {
  display: flex;
  align-items: center;
  font-family: 'Microsoft YaHei', sans-serif;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.comp-checkbox {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'Microsoft YaHei', sans-serif;
  white-space: nowrap;
}

.comp-combo {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 4px;
  background-color: #fff;
  font-family: 'Microsoft YaHei', sans-serif;
}

.comp-radio {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'Microsoft YaHei', sans-serif;
  white-space: nowrap;
}

.comp-textarea {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 4px 8px;
  resize: none;
  font-family: 'Microsoft YaHei', sans-serif;
  color: #333;
  background-color: #fff;
}

.comp-progress {
  display: flex;
  align-items: center;
  background-color: #e9ecef;
  border-radius: 10px;
  overflow: hidden;
  height: 60%;
  width: 90%;
}

.progress-fill {
  height: 100%;
  border-radius: 10px;
  transition: width 0.3s;
}

.comp-slider {
  width: 90%;
}

.comp-listbox {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  background-color: #fff;
  width: 100%;
  height: 100%;
  overflow: hidden;
  font-size: 11px;
}

.listbox-item {
  padding: 2px 6px;
  border-bottom: 1px solid #eee;
  color: #333;
}

.comp-table {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  background-color: #fff;
  width: 100%;
  height: 100%;
  font-size: 11px;
}

.table-header {
  background-color: #f5f7fa;
  padding: 2px 6px;
  font-weight: bold;
  border-bottom: 1px solid #c0c4cc;
  color: #333;
}

.table-cell {
  padding: 2px 6px;
  border-bottom: 1px solid #eee;
  color: #999;
}

.comp-group {
  border: 2px solid #c0c4cc;
  border-radius: 4px;
  width: 100%;
  height: 100%;
  padding: 4px;
}

.comp-group legend {
  font-size: 12px;
  padding: 0 4px;
}

.comp-tab {
  display: flex;
  width: 100%;
  height: 100%;
  font-size: 11px;
}

.tab-header {
  padding: 4px 12px;
  background-color: #e9ecef;
  border: 1px solid #c0c4cc;
  border-bottom: none;
  border-radius: 4px 4px 0 0;
  color: #333;
  white-space: nowrap;
}

.comp-tree {
  font-size: 12px;
  color: #333;
  text-align: left;
  padding: 4px;
  width: 100%;
}

.comp-image {
  border: 1px dashed #c0c4cc;
  border-radius: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #999;
  background-color: #fafafa;
}

.comp-datepicker {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 4px;
  font-family: 'Microsoft YaHei', sans-serif;
}

.comp-link {
  color: #409eff;
  text-decoration: underline;
  cursor: default;
  font-family: 'Microsoft YaHei', sans-serif;
  white-space: nowrap;
}

.comp-divider {
  border: none;
  border-top: 1px solid #dcdfe6;
  width: 100%;
}

.comp-statusbar {
  background-color: #f5f5f5;
  border-top: 1px solid #d0d0d0;
  font-size: 11px;
  color: #666;
  padding: 2px 8px;
  width: 100%;
  text-align: left;
}

.comp-menu {
  display: flex;
  gap: 2px;
  background-color: #f0f0f0;
  width: 100%;
  height: 100%;
  align-items: center;
  padding: 0 4px;
}

.menu-item-render {
  font-size: 11px;
  color: #333;
  padding: 2px 6px;
}

.comp-toolbar {
  display: flex;
  gap: 4px;
  background-color: #f0f0f0;
  width: 100%;
  height: 100%;
  align-items: center;
  padding: 0 4px;
}

.toolbar-btn {
  font-size: 14px;
  cursor: default;
}

.comp-container {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  background-color: #fafafa;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #999;
}

.resize-handles {
  position: absolute;
  inset: -6px;
  pointer-events: none;
  z-index: 10;
}

.resize-handle {
  position: absolute;
  width: 10px;
  height: 10px;
  background-color: #fff;
  border: 2px solid #409eff;
  border-radius: 2px;
  pointer-events: auto;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  transition: transform 0.15s, background-color 0.15s;
}

.resize-handle:hover {
  transform: scale(1.2);
  background-color: #409eff;
}

.resize-handle.nw { top: -5px; left: -5px; cursor: nw-resize; }
.resize-handle.n { top: -5px; left: 50%; margin-left: -5px; cursor: n-resize; }
.resize-handle.ne { top: -5px; right: -5px; cursor: ne-resize; }
.resize-handle.e { top: 50%; margin-top: -5px; right: -5px; cursor: e-resize; }
.resize-handle.se { bottom: -5px; right: -5px; cursor: se-resize; }
.resize-handle.s { bottom: -5px; left: 50%; margin-left: -5px; cursor: s-resize; }
.resize-handle.sw { bottom: -5px; left: -5px; cursor: sw-resize; }
.resize-handle.w { top: 50%; margin-top: -5px; left: -5px; cursor: w-resize; }

.selection-box {
  position: absolute;
  border: 1px dashed #409eff;
  background-color: rgba(64, 158, 255, 0.1);
  pointer-events: none;
  z-index: 1000;
}

.canvas-placeholder {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--text-placeholder);
  pointer-events: none;
}

.canvas-placeholder p {
  font-size: 14px;
  margin: 0;
}

.canvas-placeholder .hint {
  font-size: 12px;
  color: var(--border-base);
}

.window-statusbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 22px;
  background-color: #f0f0f0;
  border-top: 1px solid #d0d0d0;
  padding: 0 8px;
  font-size: 11px;
  color: #666;
  flex-shrink: 0;
}

.comp-text-render {
  display: flex;
  align-items: center;
  font-family: 'Microsoft YaHei', sans-serif;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.comp-icon-render {
  display: flex;
  align-items: center;
  justify-content: center;
}

.comp-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-family: 'Microsoft YaHei', sans-serif;
}

.comp-badge {
  background-color: #f56c6c;
  color: #fff;
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 8px;
  line-height: 1.4;
}

.comp-tag-render {
  color: #fff;
  padding: 1px 6px;
  border-radius: 3px;
  font-size: 11px;
  white-space: nowrap;
}

.comp-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #e9ecef;
  border-top-color: #409eff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.comp-skeleton {
  width: 100%;
  padding: 4px;
}

.skeleton-line {
  height: 8px;
  background: linear-gradient(90deg, #e9ecef 25%, #f5f5f5 50%, #e9ecef 75%);
  border-radius: 2px;
  margin-bottom: 4px;
}

.skeleton-line.short {
  width: 60%;
}

.comp-tooltip-render {
  font-size: 11px;
  color: #666;
  white-space: nowrap;
}

.comp-popover-render {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
}

.popover-bubble {
  background-color: #fff;
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 10px;
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
}

.comp-switch-render {
  display: flex;
  align-items: center;
}

.switch-track {
  width: 36px;
  height: 18px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  padding: 0 2px;
  background-color: #409eff;
}

.switch-thumb {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background-color: #fff;
  margin-left: auto;
  box-shadow: 0 1px 3px rgba(0,0,0,0.2);
}

.comp-number {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 6px;
  width: 60px;
  background-color: #fff;
  font-family: 'Microsoft YaHei', sans-serif;
}

.comp-select-render {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 4px;
  background-color: #fff;
  font-family: 'Microsoft YaHei', sans-serif;
}

.comp-upload {
  border: 1px dashed #c0c4cc;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 11px;
  color: #666;
  background-color: #fafafa;
  white-space: nowrap;
}

.comp-color-picker {
  display: flex;
  align-items: center;
  gap: 6px;
}

.color-swatch {
  width: 20px;
  height: 20px;
  border-radius: 3px;
  border: 1px solid #d0d0d0;
}

.color-label {
  font-size: 10px;
  color: #666;
  font-family: monospace;
}

.comp-rate {
  color: #f7ba2a;
  font-size: 14px;
  letter-spacing: 2px;
}

.comp-transfer {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
}

.transfer-panel {
  border: 1px solid #d0d0d0;
  border-radius: 3px;
  padding: 2px 4px;
  background-color: #fafafa;
}

.comp-cascader {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 6px;
  background-color: #fff;
  font-family: 'Microsoft YaHei', sans-serif;
  white-space: nowrap;
}

.comp-time-picker {
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 4px;
  font-family: 'Microsoft YaHei', sans-serif;
}

.comp-form-render {
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  padding: 4px;
  width: 100%;
  background-color: #fff;
}

.form-row {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 2px;
  font-size: 10px;
}

.form-row label {
  width: 30px;
  color: #666;
  text-align: right;
}

.form-row input {
  flex: 1;
  border: 1px solid #e0e0e0;
  border-radius: 2px;
  padding: 1px 4px;
  font-size: 10px;
  height: 16px;
}

.comp-form-item {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;
}

.form-item-label {
  font-size: 11px;
  color: #666;
  white-space: nowrap;
}

.form-item-input {
  flex: 1;
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 6px;
  height: 22px;
  font-size: 11px;
  background-color: #fff;
}

.comp-calendar {
  width: 100%;
  font-size: 10px;
}

.calendar-header {
  text-align: center;
  font-weight: bold;
  margin-bottom: 2px;
  color: #333;
}

.calendar-grid {
  display: flex;
  gap: 2px;
  justify-content: center;
}

.calendar-cell {
  width: 16px;
  height: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 2px;
  color: #666;
}

.comp-timeline {
  width: 100%;
  font-size: 10px;
  padding: 2px 4px;
}

.timeline-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 1px 0;
  color: #666;
}

.timeline-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: #409eff;
  flex-shrink: 0;
}

.comp-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  font-family: 'Microsoft YaHei', sans-serif;
}

.stat-label {
  font-size: 10px;
  color: #999;
}

.comp-description {
  width: 100%;
  font-size: 10px;
}

.desc-row {
  display: flex;
  gap: 8px;
  padding: 1px 0;
  border-bottom: 1px solid #f0f0f0;
}

.desc-label {
  color: #999;
  width: 30px;
}

.desc-value {
  color: #333;
}

.comp-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-size: 11px;
  color: #999;
}

.comp-result {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  font-size: 11px;
  color: #333;
}

.result-icon {
  font-size: 20px;
}

.comp-carousel {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  color: #999;
}

.carousel-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: #d0d0d0;
}

.carousel-dot.active {
  background-color: #409eff;
}

.comp-collapse-render {
  width: 100%;
}

.collapse-header {
  font-size: 11px;
  color: #333;
  padding: 2px 4px;
  background-color: #f5f7fa;
  border: 1px solid #e0e0e0;
  border-radius: 3px;
}

.comp-pagination {
  display: flex;
  align-items: center;
  gap: 2px;
}

.page-btn {
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #d0d0d0;
  border-radius: 2px;
  font-size: 10px;
  color: #666;
  background-color: #fff;
}

.page-btn.active {
  background-color: #409eff;
  color: #fff;
  border-color: #409eff;
}

.comp-steps-render {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 10px;
}

.step {
  color: #999;
}

.step.active {
  color: #409eff;
  font-weight: bold;
}

.step-arrow {
  color: #d0d0d0;
}

.comp-breadcrumb {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 10px;
  color: #999;
}

.comp-breadcrumb .active {
  color: #333;
}

.comp-dropdown-render {
  display: flex;
  align-items: center;
  gap: 2px;
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 8px;
  font-size: 11px;
  background-color: #fff;
  color: #333;
  white-space: nowrap;
}

.comp-anchor {
  font-size: 10px;
  padding: 2px;
}

.anchor-item {
  padding: 1px 4px;
  color: #666;
  border-left: 2px solid transparent;
}

.anchor-item.active {
  color: #409eff;
  border-left-color: #409eff;
}

.comp-backtop {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background-color: #e9ecef;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: #666;
}

.comp-drawer-render {
  border: 1px solid #d0d0d0;
  border-radius: 4px 0 0 4px;
  background-color: #fff;
  width: 100%;
}

.drawer-header {
  padding: 4px 8px;
  font-size: 11px;
  font-weight: bold;
  color: #333;
  border-bottom: 1px solid #eee;
}

.comp-affix {
  font-size: 10px;
  color: #666;
  white-space: nowrap;
}

.comp-alert-render {
  display: flex;
  align-items: center;
  gap: 4px;
  border: 1px solid #e6a23c;
  border-radius: 3px;
  padding: 2px 6px;
  background-color: #fdf6ec;
  font-size: 10px;
  color: #e6a23c;
  white-space: nowrap;
}

.comp-message-render {
  color: #fff;
  padding: 2px 8px;
  border-radius: 3px;
  font-size: 11px;
  white-space: nowrap;
}

.comp-notification {
  display: flex;
  align-items: flex-start;
  gap: 4px;
  padding: 4px;
  background-color: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  width: 100%;
}

.notif-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.notif-body {
  flex: 1;
  min-width: 0;
}

.notif-title {
  font-size: 11px;
  font-weight: bold;
  color: #333;
}

.notif-desc {
  font-size: 9px;
  color: #999;
}

.comp-dialog-render {
  border: 1px solid #d0d0d0;
  border-radius: 6px;
  background-color: #fff;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  width: 100%;
  overflow: hidden;
}

.dialog-title {
  padding: 4px 8px;
  font-size: 11px;
  font-weight: bold;
  color: #333;
  border-bottom: 1px solid #eee;
}

.dialog-btns {
  display: flex;
  justify-content: flex-end;
  gap: 4px;
  padding: 4px 8px;
}

.dialog-btn {
  padding: 1px 8px;
  border: 1px solid #d0d0d0;
  border-radius: 3px;
  font-size: 10px;
  color: #666;
  background-color: #fff;
}

.comp-loading {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #666;
}

.loading-spin {
  width: 14px;
  height: 14px;
  border: 2px solid #e9ecef;
  border-top-color: #409eff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.comp-popconfirm {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  color: #666;
}

.popconfirm-btns {
  color: #409eff;
  font-size: 9px;
}

.comp-watermark {
  font-size: 14px;
  color: rgba(0,0,0,0.06);
  transform: rotate(-20deg);
  white-space: nowrap;
  font-weight: bold;
}

.comp-tour {
  font-size: 10px;
  color: #666;
  white-space: nowrap;
}

.comp-scrollbar {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
}

.scrollbar-track {
  width: 100%;
  height: 6px;
  background-color: #f0f0f0;
  border-radius: 3px;
  position: relative;
}

.scrollbar-thumb {
  width: 30%;
  height: 100%;
  background-color: #c0c4cc;
  border-radius: 3px;
}

.comp-panel-render {
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  background-color: #fff;
  width: 100%;
  overflow: hidden;
}

.panel-title {
  padding: 3px 8px;
  font-size: 11px;
  font-weight: bold;
  color: #333;
  background-color: #f5f7fa;
  border-bottom: 1px solid #e0e0e0;
}

.comp-splitter {
  display: flex;
  width: 100%;
  height: 100%;
  font-size: 9px;
}

.splitter-pane {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
  color: #999;
}

.splitter-bar {
  width: 4px;
  background-color: #d0d0d0;
  cursor: col-resize;
  flex-shrink: 0;
}

.comp-scroll-view {
  border: 1px solid #d0d0d0;
  border-radius: 3px;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background-color: #fff;
}

.scroll-content {
  padding: 4px;
  font-size: 10px;
  color: #666;
}

.comp-dock-panel {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  font-size: 8px;
  color: #999;
}

.dock-top, .dock-bottom {
  height: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f0f0;
  border: 1px solid #e0e0e0;
}

.dock-mid {
  flex: 1;
  display: flex;
}

.dock-left, .dock-right {
  width: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  border: 1px solid #e0e0e0;
}

.dock-center {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fff;
  border: 1px solid #e0e0e0;
}

.comp-card-render {
  border: 1px solid #d0d0d0;
  border-radius: 6px;
  background-color: #fff;
  width: 100%;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
}

.card-header {
  padding: 3px 8px;
  font-size: 11px;
  font-weight: bold;
  color: #333;
  border-bottom: 1px solid #eee;
}

.card-body {
  padding: 4px 8px;
  font-size: 10px;
  color: #999;
}

.comp-chart-svg {
  width: 100%;
  height: 100%;
}

.comp-space-render {
  width: 100%;
  height: 100%;
  background: repeating-linear-gradient(45deg, transparent, transparent 3px, rgba(64,158,255,0.08) 3px, rgba(64,158,255,0.08) 6px);
}

.comp-heading {
  font-weight: bold;
  font-family: 'Microsoft YaHei', sans-serif;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.comp-paragraph {
  font-family: 'Microsoft YaHei', sans-serif;
  color: #666;
  margin: 0;
  line-height: 1.5;
  overflow: hidden;
}

.comp-blockquote {
  border-left: 3px solid #409eff;
  padding: 2px 8px;
  margin: 0;
  color: #666;
  font-style: italic;
  background-color: #f8f9fa;
  font-size: 11px;
}

.comp-code-block {
  background-color: #282c34;
  color: #abb2bf;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 10px;
  margin: 0;
  overflow: hidden;
  white-space: nowrap;
}

.comp-video-render {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #000;
  color: #fff;
  border-radius: 4px;
  font-size: 12px;
}

.comp-autocomplete {
  width: 100%;
}

.comp-autocomplete input {
  width: 100%;
  border: 1px solid #c0c4cc;
  border-radius: 3px;
  padding: 0 6px;
  height: 22px;
  font-size: 11px;
  background-color: #fff;
}

.autocomplete-drop {
  font-size: 9px;
  color: #999;
  padding: 2px 4px;
  background-color: #fafafa;
  border: 1px solid #eee;
  border-top: none;
  border-radius: 0 0 3px 3px;
}

.comp-rich-text {
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  padding: 4px;
  width: 100%;
  background-color: #fff;
}

.rt-bold { font-weight: bold; margin-right: 4px; font-size: 10px; cursor: default; }
.rt-italic { font-style: italic; margin-right: 4px; font-size: 10px; cursor: default; }
.rt-underline { text-decoration: underline; margin-right: 4px; font-size: 10px; cursor: default; }

.rt-content {
  margin-top: 4px;
  font-size: 10px;
  color: #666;
  min-height: 16px;
  border-top: 1px solid #eee;
  padding-top: 2px;
}

.comp-list-render {
  width: 100%;
  font-size: 10px;
}

.list-item-render {
  padding: 1px 4px;
  color: #666;
  border-bottom: 1px solid #f5f5f5;
}

.comp-image-viewer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  background-color: #f0f0f0;
  border-radius: 4px;
  font-size: 12px;
}

.viewer-nav {
  font-size: 9px;
  color: #999;
}

.comp-virtual-list {
  width: 100%;
  font-size: 9px;
  border: 1px solid #e0e0e0;
  border-radius: 3px;
  overflow: hidden;
}

.vlist-item {
  padding: 1px 4px;
  border-bottom: 1px solid #f0f0f0;
  color: #666;
}

.vlist-more {
  text-align: center;
  color: #999;
  padding: 1px;
}

.comp-navbar {
  display: flex;
  align-items: center;
  gap: 6px;
  background-color: #f8f9fa;
  padding: 2px 8px;
  border-bottom: 1px solid #e0e0e0;
  width: 100%;
  font-size: 10px;
}

.navbar-brand {
  font-size: 14px;
}

.navbar-item {
  color: #666;
  cursor: default;
}

.comp-layout-render {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  overflow: hidden;
  font-size: 8px;
}

.layout-header, .layout-footer {
  height: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f0f0;
  color: #999;
}

.layout-body {
  flex: 1;
  display: flex;
}

.layout-side {
  width: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  color: #999;
  border-right: 1px solid #e0e0e0;
}

.layout-main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}

.comp-row-render {
  display: flex;
  gap: 2px;
  width: 100%;
}

.row-col {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(64,158,255,0.08);
  border: 1px dashed rgba(64,158,255,0.3);
  border-radius: 2px;
  font-size: 9px;
  color: #999;
  height: 100%;
}

.comp-col-render {
  display: flex;
  flex-direction: column;
  gap: 2px;
  width: 100%;
  height: 100%;
}

.col-row {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(64,158,255,0.08);
  border: 1px dashed rgba(64,158,255,0.3);
  border-radius: 2px;
  font-size: 9px;
  color: #999;
}

.comp-flex-render {
  display: flex;
  gap: 2px;
  width: 100%;
  height: 100%;
}

.flex-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(64,158,255,0.1);
  border-radius: 2px;
  font-size: 10px;
  color: #409eff;
  font-weight: bold;
}

.comp-tab-container {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  border: 1px solid #d0d0d0;
  border-radius: 4px;
  overflow: hidden;
}

.tab-header-row {
  display: flex;
  background-color: #f5f7fa;
  border-bottom: 1px solid #e0e0e0;
}

.tab-h {
  padding: 2px 8px;
  font-size: 9px;
  color: #999;
  border-right: 1px solid #e0e0e0;
  cursor: default;
}

.tab-h.active {
  color: #409eff;
  background-color: #fff;
  border-bottom: 2px solid #409eff;
}

.tab-body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  color: #999;
}
</style>