<script lang="ts" setup>
import { ref, reactive, computed, watch, nextTick, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import MonacoEditorWrapper from './MonacoEditorWrapper.vue'
import { useTheme } from '../../stores/theme'

const { currentTheme } = useTheme()
const editorTheme = computed(() => currentTheme.value === 'dark' ? 'vs-dark' : 'vs')

interface ClassMember {
  id: string; name: string; type: string; public: boolean; static: boolean;
  ref: boolean; value: string; comment: string
  x?: number; y?: number; width?: number; height?: number
}

interface MethodParam {
  name: string; type: string; comment: string
}

interface LocalVar {
  name: string; type: string; isStatic: boolean; isRef: boolean; value: string; comment: string
}

interface ClassMethod {
  id: string; name: string; type: '通常' | '接收事件' | '定义事件'
  public: boolean; isStatic: boolean; returnType: string; returnComment: string
  comment: string; params: MethodParam[]; localVars: LocalVar[]
  code: string; collapsed: boolean; activeTab: 'code' | 'params' | 'locals'
}

interface ClassProp {
  name: string; value: string; comment: string
}

interface ClassInfo {
  id: string; name: string; baseClass: string; isPublic: boolean; comment: string
  properties: ClassProp[]; members: ClassMember[]; methods: ClassMethod[]
  collapsed: boolean
}

interface OutlineItem {
  type: 'class' | 'method' | 'member' | 'property'
  name: string
  icon: string
  line: number
  children?: OutlineItem[]
}

const props = defineProps<{
  modelValue?: string
  canvasComponents?: any[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', text: string): void
  (e: 'component-add', component: any): void
  (e: 'component-update', components: any[]): void
  (e: 'code-change', code: string): void
}>()

let idCounter = 100
const newId = () => `_${idCounter++}`

const programName = ref('package main')
const programAttr = ref('')
const programAttrVal = ref('')

const classes = ref<ClassInfo[]>([
  {
    id: `cls${newId()}`, name: '启动类', baseClass: '',
    isPublic: true, comment: '',
    properties: [],
    members: [
      { id: `m${newId()}`, name: 'mainWindow', type: 'MainWindow', public: false, static: false, ref: false, value: '', comment: '' }
    ],
    methods: [{
      id: `mt${newId()}`, name: 'main', type: '通常',
      public: true, isStatic: false, returnType: '', returnComment: '', comment: '程序入口',
      params: [], localVars: [],
      code: `func main() {\n\t// 在此编写启动代码\n\tmainWindow.Create()\n}`,
      collapsed: false, activeTab: 'code'
    }],
    collapsed: false
  }
])

const currentCode = ref('')
const editorMode = ref<'visual' | 'code'>('code')
const showOutline = ref(true)
const showFunctionPanel = ref(false)
const activeTabIndex = ref(0)
const cursorPosition = ref({ line: 1, column: 1 })

interface EditCell {
  type: 'member-name' | 'member-type' | 'method-name' | 'class-name' | 'prop-name' | 'prop-value'
  clsIdx: number
  itemIdx: number
  field: string
}

const editingCell = ref<EditCell | null>(null)
const editValue = ref('')
const editInputRef = ref<HTMLInputElement | null>(null)

const startCellEdit = (cell: EditCell, currentValue: string) => {
  editingCell.value = cell
  editValue.value = currentValue
  nextTick(() => {
    editInputRef.value?.focus()
    editInputRef.value?.select()
  })
}

const commitCellEdit = () => {
  if (!editingCell.value) return
  const { type, clsIdx, itemIdx, field } = editingCell.value
  const val = editValue.value

  if (type === 'member-name' || type === 'member-type') {
    ;(classes.value[clsIdx].members[itemIdx] as any)[field] = val
  } else if (type === 'method-name') {
    ;(classes.value[clsIdx].methods[itemIdx] as any)[field] = val
  } else if (type === 'class-name') {
    ;(classes.value[clsIdx] as any)[field] = val
  } else if (type === 'prop-name' || type === 'prop-value') {
    ;(classes.value[clsIdx].properties[itemIdx] as any)[field] = val
  }

  editingCell.value = null
  syncVisualToCode()
}

const cancelCellEdit = () => {
  editingCell.value = null
}

const handleEditKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter') {
    e.preventDefault()
    commitCellEdit()
  } else if (e.key === 'Escape') {
    e.preventDefault()
    cancelCellEdit()
  }
}

const outlineItems = computed<OutlineItem[]>(() => {
  const items: OutlineItem[] = []
  let lineNum = 1

  classes.value.forEach(cls => {
    const classItem: OutlineItem = {
      type: 'class', name: cls.name, icon: '📁', line: lineNum++, children: []
    }

    cls.properties.forEach(p => {
      classItem.children!.push({ type: 'property', name: p.name, icon: '⚙️', line: lineNum++ })
    })

    cls.members.forEach(m => {
      classItem.children!.push({ type: 'member', name: `${m.name}: ${m.type}`, icon: '📊', line: lineNum++ })
    })

    cls.methods.forEach(mtd => {
      classItem.children!.push({ type: 'method', name: mtd.name, icon: '⚡', line: lineNum++ })
    })

    items.push(classItem)
  })

  return items
})

const generateGoCode = (): string => {
  let lines: string[] = []

  lines.push(`// ${programName.value}`)
  lines.push('')

  classes.value.forEach(cls => {
    if (cls.isPublic) {
      lines.push(`// ${cls.comment || cls.name}`)
    }

    cls.members.forEach(m => {
      let decl = `var ${m.name} `
      decl += m.type === 'MainWindow' ? '*' + m.name : m.type
      if (m.value) decl += ` = ${m.value}`
      lines.push(decl)
    })

    if (cls.members.length > 0) lines.push('')

    cls.methods.forEach(mtd => {
      let sig = 'func '
      sig += `${mtd.name} (`

      mtd.params.forEach((p, i) => {
        if (i > 0) sig += ', '
        sig += `${p.name} ${p.type}`
      })

      sig += ')'

      if (mtd.returnType && mtd.returnType !== '') {
        sig += ` ${mtd.returnType}`
      }

      lines.push(sig + ' {')
      mtd.code.split('\n').filter(l => l.trim()).forEach(cl => {
        lines.push('\t' + cl.trim())
      })
      lines.push('}')
      lines.push('')
    })
  })

  const code = lines.join('\n')
  currentCode.value = code
  emit('update:modelValue', code)
  emit('code-change', code)

  return code
}

const syncVisualToCode = () => {
  generateGoCode()
}

const syncCodeToVisual = (code: string) => {
  currentCode.value = code
  ElMessage.info('代码已更新')
}

const handleCodeChange = (value: string) => {
  currentCode.value = value
  emit('update:modelValue', value)
  emit('code-change', value)
}

const handleCursorChange = (position: { lineNumber: number; column: number }) => {
  cursorPosition.value = { line: position.lineNumber, column: position.column }
}

const addClass = () => {
  classes.value.push({
    id: `cls${newId()}`, name: `NewClass${classes.value.length + 1}`, baseClass: '',
    isPublic: false, comment: '',
    properties: [],
    members: [], methods: [], collapsed: false
  })
  syncVisualToCode()
}

const removeClass = (idx: number) => {
  if (classes.value.length <= 1) { ElMessage.warning('至少保留一个类'); return }
  classes.value.splice(idx, 1)
  syncVisualToCode()
}

const addMember = (clsIdx: number) => {
  const cls = classes.value[clsIdx]
  if (!cls) return
  const m: ClassMember = {
    id: `m${newId()}`, name: `member${cls.members.length + 1}`, type: 'int',
    public: false, static: false, ref: false, value: '', comment: ''
  }
  cls.members.push(m)
  syncMemberToCanvas(clsIdx, m)
  syncVisualToCode()
}

const removeMember = (clsIdx: number, memIdx: number) => {
  classes.value[clsIdx].members.splice(memIdx, 1)
  syncVisualToCode()
}

const updateMember = (clsIdx: number, memIdx: number, field: string, val: any) => {
  ;(classes.value[clsIdx].members[memIdx] as any)[field] = val
  syncVisualToCode()
}

const addMethod = (clsIdx: number) => {
  const cls = classes.value[clsIdx]
  if (!cls) return
  cls.methods.push({
    id: `mt${newId()}`, name: `newMethod${cls.methods.length + 1}`, type: '通常',
    public: true, isStatic: false, returnType: '', returnComment: '', comment: '',
    params: [], localVars: [],
    code: '// TODO: implement\n',
    collapsed: false, activeTab: 'code'
  })
  syncVisualToCode()
}

const removeMethod = (clsIdx: number, mtdIdx: number) => {
  classes.value[clsIdx].methods.splice(mtdIdx, 1)
  syncVisualToCode()
}

const updateMethod = (clsIdx: number, mtdIdx: number, field: string, val: any) => {
  ;(classes.value[clsIdx].methods[mtdIdx] as any)[field] = val
  syncVisualToCode()
}

const updateMethodCode = (clsIdx: number, mtdIdx: number, code: string) => {
  classes.value[clsIdx].methods[mtdIdx].code = code
  syncVisualToCode()
}

const addClassProp = (clsIdx: number) => {
  classes.value[clsIdx].properties.push({ name: '', value: '', comment: '' })
  syncVisualToCode()
}

const removeClassProp = (clsIdx: number, pIdx: number) => {
  classes.value[clsIdx].properties.splice(pIdx, 1)
  syncVisualToCode()
}

const updateClassProp = (clsIdx: number, pIdx: number, field: string, val: string) => {
  ;(classes.value[clsIdx].properties[pIdx] as any)[field] = val
  syncVisualToCode()
}

const syncMemberToCanvas = (clsIdx: number, member: ClassMember) => {
  const typeMap: Record<string, string> = {
    'Button': 'button', 'Input': 'input', 'Label': 'label', 'Edit': 'edit',
    'CheckBox': 'checkbox', 'Radio': 'radio', 'ListBox': 'listbox',
    'ComboBox': 'combobox', 'Image': 'image', 'Container': 'container',
    'WebView': 'webview', 'PDFView': 'pdfview'
  }

  if (typeMap[member.type]) {
    emit('component-add', {
      id: member.id, type: typeMap[member.type], name: member.name,
      x: 50 + clsIdx * 20, y: 50 + classes.value[clsIdx].members.indexOf(member) * 30,
      width: member.type === 'Button' ? 80 : 120,
      height: member.type === 'Button' ? 32 : 32,
      text: member.name, visible: true
    })
  }
}

const getClassComponents = (): any[] => {
  const result: any[] = []
  const typeMap: Record<string, string> = {
    'Button': 'button', 'Input': 'input', 'Label': 'label',
    'Edit': 'edit', 'WebView': 'webview'
  }

  classes.value.forEach(cls => {
    cls.members.forEach(m => {
      if (typeMap[m.type]) {
        result.push({
          id: m.id, type: typeMap[m.type], name: m.name,
          x: 0, y: 0, width: 100, height: 32,
          text: m.name, color: '#409eff', fontSize: 14, visible: true
        })
      }
    })
  })

  return result
}

const syncFromCanvas = (components: any[]) => {
  components.forEach(comp => {
    classes.value.forEach(cls => {
      const mem = cls.members.find(m => m.id === comp.id)
      if (mem) Object.assign(mem, { x: comp.x, y: comp.y, width: comp.width, height: comp.height })
    })
  })
  syncVisualToCode()
}

const goToLine = (line: number) => {
  ElMessage.info(`跳转到第 ${line} 行`)
}

const formatCode = () => {
  ElMessage.success('代码格式化完成')
}

const toggleMode = () => {
  editorMode.value = editorMode.value === 'visual' ? 'code' : 'visual'
  if (editorMode.value === 'code') {
    nextTick(() => {
      generateGoCode()
    })
  }
}

watch(() => props.modelValue, (newVal) => {
  if (newVal && newVal !== currentCode.value) {
    currentCode.value = newVal
  }
}, { immediate: true })

watch(() => props.canvasComponents, (comps) => {
  if (comps?.length) syncFromCanvas(comps)
}, { deep: true })

onMounted(() => {
  generateGoCode()
})

defineExpose({
  getProgramData: () => ({ programName: programName.value, classes: JSON.parse(JSON.stringify(classes.value)) }),
  setProgramData: (d: any) => {},
  generateCode: generateGoCode,
  syncFromCanvas,
  getClassComponents,
  addClass,
  getValue: () => currentCode.value,
  setValue: (v: string) => { currentCode.value = v }
})
</script>

<template>
  <div class="ide-code-editor">
    <!-- 工具栏 -->
    <div class="editor-toolbar">
      <div class="toolbar-left">
        <el-button-group size="small">
          <el-button :type="editorMode === 'code' ? 'primary' : ''" @click="toggleMode">
            📝 代码
          </el-button>
          <el-button :type="editorMode === 'visual' ? 'primary' : ''" @click="toggleMode">
            👁️ 视图
          </el-button>
        </el-button-group>

        <span class="toolbar-separator"></span>

        <el-tooltip content="保存 (Ctrl+S)" placement="bottom">
          <el-button size="small" @click="emit('update:modelValue', currentCode)">
            💾
          </el-button>
        </el-tooltip>

        <el-tooltip content="格式化代码" placement="bottom">
          <el-button size="small" @click="formatCode">✨</el-button>
        </el-tooltip>

        <el-tooltip content="文件大纲" placement="bottom">
          <el-button size="small" :type="showOutline ? 'primary' : ''" @click="showOutline = !showOutline">
            📋
          </el-button>
        </el-tooltip>
      </div>

      <div class="toolbar-center">
        <span class="file-name">{{ programName }}</span>
      </div>

      <div class="toolbar-right">
        <span class="cursor-position">Ln {{ cursorPosition.line }}, Col {{ cursorPosition.column }}</span>
        <span class="lang-badge">Go</span>
      </div>
    </div>

    <!-- 主编辑区 -->
    <div class="editor-main">
      <!-- 左侧大纲面板 -->
      <div v-if="showOutline" class="outline-panel">
        <div class="panel-header">
          <span>📋 大纲</span>
        </div>
        <div class="outline-content">
          <div v-for="(item, idx) in outlineItems" :key="idx" class="outline-item outline-class">
            <span class="item-icon">{{ item.icon }}</span>
            <span class="item-name" @click="goToLine(item.line)">{{ item.name }}</span>
            <div v-if="item.children" class="outline-children">
              <div v-for="(child, cIdx) in item.children" :key="cIdx"
                   class="outline-item outline-child"
                   @click="goToLine(child.line)">
                <span class="item-icon">{{ child.icon }}</span>
                <span class="item-name">{{ child.name }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 中间编辑器 -->
      <div class="editor-container">
        <!-- 代码模式 -->
        <div v-if="editorMode === 'code'" class="code-mode">
          <MonacoEditorWrapper
            v-model="currentCode"
            language="go"
            :theme="editorTheme"
            :minimap="true"
            :font-size="14"
            @change="handleCodeChange"
            @cursor-change="handleCursorChange"
          />
        </div>

        <!-- 视觉模式 -->
        <div v-else class="visual-mode">
          <!-- 包名行 -->
          <div class="visual-row package-row">
            <span class="row-icon">📦</span>
            <span class="row-label">Package:</span>
            <el-input size="small" v-model="programName" @change="syncVisualToCode" style="width:200px" />
          </div>

          <!-- 类列表 -->
          <div v-for="(cls, clsIdx) in classes" :key="cls.id" class="class-section">
            <!-- 类头 -->
            <div class="visual-row class-header" @click="cls.collapsed = !cls.collapsed">
              <span class="collapse-icon">{{ cls.collapsed ? '▶' : '▼' }}</span>
              <span class="row-icon">📁</span>
              <template v-if="editingCell?.type === 'class-name' && editingCell?.clsIdx === clsIdx">
                <input
                  ref="editInputRef"
                  v-model="editValue"
                  class="inline-edit-input"
                  @keydown="handleEditKeydown"
                  @blur="commitCellEdit"
                  @click.stop
                />
              </template>
              <template v-else>
                <el-input size="small" v-model="cls.name" @change="syncVisualToCode" style="width:140px" @click.stop
                  @dblclick.stop="startCellEdit({ type: 'class-name', clsIdx, itemIdx: 0, field: 'name' }, cls.name)" />
              </template>
              <el-input size="small" v-model="cls.baseClass" placeholder="base class" @change="syncVisualToCode" style="width:120px" @click.stop />
              <el-checkbox v-model="cls.isPublic" @change="syncVisualToCode" @click.stop>public</el-checkbox>
              <el-button v-if="classes.length > 1" size="small" type="danger" link @click.stop="removeClass(clsIdx)">✕</el-button>
            </div>

            <template v-if="!cls.collapsed">
              <!-- 属性 -->
              <div v-if="cls.properties.length > 0" class="sub-section">
                <div class="sub-header">Properties</div>
                <div v-for="(p, pIdx) in cls.properties" :key="pIdx" class="prop-row">
                  <el-input size="small" v-model="p.name" placeholder="name" @change="syncVisualToCode" style="width:120px" />
                  <el-input size="small" v-model="p.value" placeholder="value" @change="syncVisualToCode" style="flex:1" />
                  <el-button size="small" type="danger" link @click="removeClassProp(clsIdx, pIdx)">✕</el-button>
                </div>
              </div>

              <!-- 成员变量 -->
              <div class="sub-section">
                <div class="sub-header">
                  <span>Members</span>
                  <el-button size="small" link type="primary" @click="addMember(clsIdx)">+</el-button>
                </div>
                <div class="members-grid">
                  <div v-for="(m, mIdx) in cls.members" :key="m.id" class="member-chip">
                    <template v-if="editingCell?.type === 'member-name' && editingCell?.clsIdx === clsIdx && editingCell?.itemIdx === mIdx">
                      <input
                        ref="editInputRef"
                        v-model="editValue"
                        class="inline-edit-input"
                        style="width:90px"
                        @keydown="handleEditKeydown"
                        @blur="commitCellEdit"
                      />
                    </template>
                    <template v-else>
                      <el-input size="small" v-model="m.name" @change="syncVisualToCode" style="width:90px"
                        @dblclick.stop="startCellEdit({ type: 'member-name', clsIdx, itemIdx: mIdx, field: 'name' }, m.name)" />
                    </template>
                    <span>:</span>
                    <template v-if="editingCell?.type === 'member-type' && editingCell?.clsIdx === clsIdx && editingCell?.itemIdx === mIdx">
                      <input
                        ref="editInputRef"
                        v-model="editValue"
                        class="inline-edit-input"
                        style="width:100px"
                        @keydown="handleEditKeydown"
                        @blur="commitCellEdit"
                      />
                    </template>
                    <template v-else>
                      <el-select size="small" v-model="m.type" @change="syncVisualToCode" style="width:100px"
                        @dblclick.stop="startCellEdit({ type: 'member-type', clsIdx, itemIdx: mIdx, field: 'type' }, m.type)">
                        <el-option label="int" value="int" />
                        <el-option label="string" value="string" />
                        <el-option label="bool" value="bool" />
                        <el-option label="float64" value="float64" />
                        <el-option label="Button" value="Button" />
                        <el-option label="Label" value="Label" />
                        <el-option label="Edit" value="Edit" />
                        <el-option label="WebView" value="WebView" />
                      </el-select>
                    </template>
                    <el-button size="small" type="danger" link @click="removeMember(clsIdx, mIdx)">✕</el-button>
                  </div>
                </div>
              </div>

              <!-- 方法 -->
              <div v-for="(mtd, mtdIdx) in cls.methods" :key="mtd.id" class="method-section">
                <div class="method-header" @click="mtd.collapsed = !mtd.collapsed">
                  <span class="collapse-icon">{{ mtd.collapsed ? '▶' : '▼' }}</span>
                  <span class="row-icon">⚡</span>
                  <template v-if="editingCell?.type === 'method-name' && editingCell?.clsIdx === clsIdx && editingCell?.itemIdx === mtdIdx">
                    <input
                      ref="editInputRef"
                      v-model="editValue"
                      class="inline-edit-input"
                      style="width:150px"
                      @keydown="handleEditKeydown"
                      @blur="commitCellEdit"
                      @click.stop
                    />
                  </template>
                  <template v-else>
                    <el-input size="small" v-model="mtd.name" @change="syncVisualToCode" style="width:150px" @click.stop
                      @dblclick.stop="startCellEdit({ type: 'method-name', clsIdx, itemIdx: mtdIdx, field: 'name' }, mtd.name)" />
                  </template>
                  <el-select size="small" v-model="mtd.type" @change="syncVisualToCode" style="width:85px" @click.stop>
                    <el-option label="func" value="通常" />
                    <el-option label="event" value="接收事件" />
                  </el-select>
                  <el-checkbox v-model="mtd.public" @change="syncVisualToCode" @click.stop>public</el-checkbox>
                  <el-button size="small" type="danger" link @click.stop="removeMethod(clsIdx, mtdIdx)">✕</el-button>
                </div>

                <template v-if="!mtd.collapsed">
                  <div class="method-body">
                    <div class="tabs">
                      <div class="tab" :class="{ active: mtd.activeTab === 'code' }" @click="mtd.activeTab = 'code'">Code</div>
                      <div class="tab" :class="{ active: mtd.activeTab === 'params' }" @click="mtd.activeTab = 'params'">Params ({{ mtd.params.length }})</div>
                      <div class="tab" :class="{ active: mtd.activeTab === 'locals' }" @click="mtd.activeTab = 'locals'">Locals ({{ mtd.localVars.length }})</div>
                    </div>

                    <div v-show="mtd.activeTab === 'code'" class="code-area">
                      <textarea
                        v-model="mtd.code"
                        @input="updateMethodCode(clsIdx, mtdIdx, ($event.target as HTMLTextAreaElement).value)"
                        class="method-code-textarea"
                        spellcheck="false"
                        placeholder="// Write your code here..."
                      ></textarea>
                    </div>

                    <div v-show="mtd.activeTab === 'params'" class="meta-list">
                      <div v-for="(p, pIdx) in mtd.params" :key="pIdx" class="meta-row">
                        <el-input size="small" v-model="p.name" placeholder="name" style="width:100px" />
                        <el-input size="small" v-model="p.type" placeholder="type" style="width:90px" />
                        <el-button size="small" type="danger" link @click="mtd.params.splice(pIdx, 1)">✕</el-button>
                      </div>
                      <el-button size="small" link type="primary" @click="mtd.params.push({name:'', type:'', comment:''})">+ Param</el-button>
                    </div>

                    <div v-show="mtd.activeTab === 'locals'" class="meta-list">
                      <div v-for="(lv, lvIdx) in mtd.localVars" :key="lvIdx" class="meta-row">
                        <el-input size="small" v-model="lv.name" placeholder="name" style="width:100px" />
                        <el-input size="small" v-model="lv.type" placeholder="type" style="width:85px" />
                        <el-button size="small" type="danger" link @click="mtd.localVars.splice(lvIdx, 1)">✕</el-button>
                      </div>
                      <el-button size="small" link type="primary" @click="mtd.localVars.push({name:'var', type:'int', isStatic:false, isRef:false, value:'', comment:''})">+ Local</el-button>
                    </div>
                  </div>
                </template>
              </div>

              <el-button size="small" link type="primary" class="add-method-btn" @click="addMethod(clsIdx)">+ Method</el-button>
            </template>
          </div>

          <div class="add-class-bar">
            <el-button size="small" type="primary" @click="addClass">+ Add Class</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ide-code-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #1e1e1e;
  color: #d4d4d4;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  overflow: hidden;
}

/* ========== 工具栏 ========== */
.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 12px;
  background-color: #252526;
  border-bottom: 1px solid #3c3c3c;
  min-height: 36px;
}

.toolbar-left, .toolbar-right {
  display: flex;
  align-items: center;
  gap: 6px;
}

.toolbar-separator {
  width: 1px;
  height: 20px;
  background-color: #3c3c3c;
  margin: 0 6px;
}

.toolbar-center {
  display: flex;
  align-items: center;
}

.file-name {
  font-size: 13px;
  color: #9cdcfe;
  font-weight: 500;
}

.cursor-position {
  font-size: 11px;
  color: #858585;
  font-family: 'Consolas', monospace;
}

.lang-badge {
  background-color: #007acc;
  color: #fff;
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 3px;
  font-weight: 600;
}

/* ========== 主编辑区 ========== */
.editor-main {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* ========== 大纲面板 ========== */
.outline-panel {
  width: 220px;
  min-width: 180px;
  background-color: #252526;
  border-right: 1px solid #3c3c3c;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  padding: 10px 14px;
  font-size: 12px;
  font-weight: 600;
  color: #bbbbbb;
  border-bottom: 1px solid #3c3c3c;
  background-color: #2d2d30;
}

.outline-content {
  flex: 1;
  overflow-y: auto;
  padding: 6px 0;
}

.outline-content::-webkit-scrollbar { width: 8px; }
.outline-content::-webkit-scrollbar-track { background: #252526; }
.outline-content::-webkit-scrollbar-thumb { background: #424242; border-radius: 0; }

.outline-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 14px;
  cursor: pointer;
  transition: background-color 0.15s;
  font-size: 12px;
}

.outline-item:hover { background-color: #2a2d2e; }

.outline-class {
  font-weight: 500;
  color: #d4d4d4;
}

.outline-child {
  padding-left: 28px;
  color: #9cdcfe;
}

.item-icon {
  font-size: 12px;
  flex-shrink: 0;
}

.item-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.outline-children {
  margin-top: 2px;
}

/* ========== 编辑器容器 ========== */
.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #1e1e1e;
}

.code-mode {
  flex: 1;
  overflow: hidden;
}

/* ========== 视觉模式 ========== */
.visual-mode {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
}

.visual-mode::-webkit-scrollbar { width: 10px; }
.visual-mode::-webkit-scrollbar-track { background: #1e1e1e; }
.visual-mode::-webkit-scrollbar-thumb { background: #424242; border-radius: 0; }

.visual-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  margin-bottom: 4px;
  background-color: #252526;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  transition: all 0.15s;
}
.visual-row:hover { border-color: #007acc; }

.package-row {
  background: linear-gradient(135deg, #1a3a5c, #1e4060);
  border-color: #0e4c7f;
}

.row-icon { font-size: 14px; flex-shrink: 0; }
.row-label { color: #9cdcfe; font-size: 12px; font-weight: 500; }

.collapse-icon {
  font-size: 10px;
  color: #858585;
  width: 16px;
  text-align: center;
  cursor: pointer;
}

.class-header {
  cursor: pointer;
  user-select: none;
  background: linear-gradient(135deg, #2d2d30, #28282a);
}
.class-header:hover { background: linear-gradient(135deg, #35383b, #303234); }

.class-section {
  margin-bottom: 12px;
}

.sub-section {
  margin: 8px 0 8px 20px;
  background-color: #1e1e1e;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  overflow: hidden;
}

.sub-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 12px;
  font-size: 11px;
  font-weight: 600;
  color: #4ec9b0;
  background-color: #252526;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.prop-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-top: 1px solid #2d2d30;
}

.members-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 10px 12px;
}

.member-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background-color: #252526;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  font-size: 12px;
  transition: all 0.15s;
}
.member-chip:hover { border-color: #007acc; }

.inline-edit-input {
  padding: 4px 8px;
  font-size: 12px;
  color: #d4d4d4;
  background-color: #1e1e1e;
  border: 1px solid #007acc;
  border-radius: 3px;
  outline: none;
  font-family: 'Consolas', 'Courier New', monospace;
  box-shadow: 0 0 0 1px #007acc;
}
.inline-edit-input:focus {
  border-color: #4fc1ff;
  box-shadow: 0 0 0 2px rgba(0, 122, 204, 0.3);
}

.method-section {
  margin: 8px 0 8px 16px;
  border: 1px solid #3c3c3c;
  border-radius: 4px;
  background-color: #1e1e1e;
  overflow: hidden;
}

.method-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  cursor: pointer;
  user-select: none;
  background: linear-gradient(to right, #2d2d30, #28282a);
  border-bottom: 1px solid #3c3c3c;
}
.method-header:hover { background: linear-gradient(to right, #35383b, #303234); }

.method-body {
  background-color: #1e1e1e;
}

.tabs {
  display: flex;
  background-color: #252526;
  border-bottom: 1px solid #3c3c3c;
}

.tab {
  padding: 7px 16px;
  font-size: 12px;
  cursor: pointer;
  color: #9cdcfe;
  border-right: 1px solid #3c3c3c;
  transition: all 0.15s;
}
.tab:hover { background-color: #2d2d30; color: #fff; }
.tab.active {
  background-color: #1e1e1e;
  color: #4fc1ff;
  font-weight: 600;
  border-bottom: 2px solid #007acc;
}

.code-area {
  padding: 0;
}

.method-code-textarea {
  width: 100%;
  min-height: 150px;
  max-height: 400px;
  padding: 14px 18px;
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #d4d4d4;
  background-color: #1e1e1e;
  border: none;
  outline: none;
  resize: vertical;
  tab-size: 4;
}
.method-code-textarea:focus { box-shadow: inset 0 0 0 1px #007acc; }
.method-code-textarea::placeholder { color: #5a5a5a; }

.meta-list {
  padding: 12px 14px;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.add-method-btn {
  margin: 8px 0 8px 24px;
}

.add-class-bar {
  padding: 12px;
  text-align: center;
  border-top: 1px solid #3c3c3c;
  margin-top: 8px;
}

/* ========== Element Plus 深色适配 ========== */
:deep(.el-input__wrapper) {
  box-shadow: 0 0 0 1px #3c3c3c inset;
  background-color: #1e1e1e;
}
:deep(.el-input__wrapper:hover) { box-shadow: 0 0 0 1px #007acc inset; }
:deep(.el-input__inner) { color: #d4d4d4; font-size: 12px; }
:deep(.el-checkbox__label) { color: #9cdcfe; font-size: 11px; }
:deep(.el-checkbox__inner) { background-color: #1e1e1e; border-color: #3c3c3c; }
:deep(.el-checkbox__input.is-checked .el-checkbox__inner) { background-color: #007acc; }
:deep(.el-button--primary) { --el-button-bg-color: #007acc; --el-button-border-color: #007acc; }
:deep(.el-button--primary.is-link) { --el-button-text-color: #4fc1ff; }
:deep(.el-button--danger.is-link) { --el-button-text-color: #f48771; }
:deep(.el-select .el-input__wrapper) { box-shadow: 0 0 0 1px #3c3c3c inset; background-color: #1e1e1e; }
</style>