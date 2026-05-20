<script lang="ts" setup>
import { ref, computed, nextTick, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  ConvertCodeRowsToGo,
  GetFunctionSuggestions,
  GetAllFunctions,
  ValidateCode,
  GetCategories,
  GetTemplates,
  GetFunctionsByCategory
} from '../../../wailsjs/go/main/App'

interface SubParam {
  name: string
  type: string
  isRef: boolean
  isNullable: boolean
  isArray: boolean
  note: string
}

interface SubLocal {
  name: string
  type: string
  isStatic: boolean
  defaultValue: string
  note: string
}

interface SubProgram {
  name: string
  returnType: string
  isPublic: boolean
  note: string
  params: SubParam[]
  locals: SubLocal[]
  codeBody: string
  collapsed: boolean
}

interface ChineseFunction {
  chineseName: string
  goFunction: string
  description: string
  category: string
  params: string[]
  returnType: string
}

interface CodeTemplate {
  name: string
  description: string
  category: string
  code: string
}

const emit = defineEmits<{
  (e: 'code-rows-change', rows: any[]): void
}>()

const editMode = ref<'visual' | 'code'>('code')
const isSyncing = ref(false)

const defaultCode = `.子程序 _启动子程序, 整数型, 公开, 程序入口
.局部变量 返回值, 整数型, , "0", 返回值变量
.局部变量 计数, 整数型, 静态, "0", 循环计数

.如果 (计数 ＝ 0)
    调试输出 ("程序启动")
.否则
    调试输出 ("计数不为0")
.如果结束

返回 (0)
`

const subPrograms = ref<SubProgram[]>([
  {
    name: '窗口程序集_启动窗口',
    returnType: '',
    isPublic: false,
    note: '',
    params: [],
    locals: [],
    codeBody: '',
    collapsed: false,
  },
])

const allFunctions = ref<ChineseFunction[]>([])
const categories = ref<string[]>([])
const templates = ref<CodeTemplate[]>([])
const selectedCategory = ref('')
const categoryFunctions = ref<ChineseFunction[]>([])
const showFunctionPanel = ref(false)
const showOutline = ref(true)
const codeErrors = ref<string[]>([])

const codeContent = ref('')
const highlightedCode = ref('')

const lineCount = computed(() => {
  return Math.max(codeContent.value.split('\n').length, 1)
})

const dataTypes = ['整数型', '小数型', '逻辑型', '文本型', '字节集', '日期时间型', '子程序指针', '通用型', '空']

const loadAllFunctions = async () => {
  try {
    const funcs = await GetAllFunctions()
    allFunctions.value = funcs || []
    const cats = await GetCategories()
    categories.value = cats || []
  } catch (e) {
    console.error('加载函数失败:', e)
    loadDefaultFunctions()
  }
}

const loadDefaultFunctions = () => {
  allFunctions.value = [
    { chineseName: '信息框', goFunction: 'MessageBox', description: '弹出信息提示框', category: '界面', params: ['提示信息', '标题', '按钮'], returnType: '整数型' },
    { chineseName: '延迟', goFunction: 'Sleep', description: '延迟指定毫秒数', category: '系统', params: ['毫秒数'], returnType: '' },
    { chineseName: '取文本长度', goFunction: 'Len', description: '取文本长度', category: '文本', params: ['文本'], returnType: '整数型' },
    { chineseName: '到文本', goFunction: 'Str', description: '转换到文本', category: '转换', params: ['数据'], returnType: '文本型' },
    { chineseName: '到整数', goFunction: 'Int', description: '转换到整数', category: '转换', params: ['文本'], returnType: '整数型' },
    { chineseName: '取启动时间', goFunction: 'GetTickCount', description: '取系统启动时间', category: '系统', params: [], returnType: '整数型' },
    { chineseName: '取现行时间', goFunction: 'Now', description: '获取当前时间', category: '时间', params: [], returnType: '日期时间型' },
    { chineseName: '运行', goFunction: 'Exec', description: '执行外部命令', category: '系统', params: ['命令', '是否等待'], returnType: '整数型' },
    { chineseName: '读入文件', goFunction: 'ReadFile', description: '读取文件内容', category: '文件', params: ['文件路径'], returnType: '字节集' },
    { chineseName: '写到文件', goFunction: 'WriteFile', description: '写入文件内容', category: '文件', params: ['文件路径', '数据'], returnType: '逻辑型' },
    { chineseName: '调试输出', goFunction: 'DebugPrint', description: '输出调试信息', category: '调试', params: ['内容'], returnType: '' },
  ]
  categories.value = ['全部', '界面', '系统', '文本', '转换', '时间', '文件', '调试']
}

const selectCategory = async (cat: string) => {
  selectedCategory.value = cat
  if (cat === '全部') {
    categoryFunctions.value = allFunctions.value
  } else {
    try {
      const funcs = await GetFunctionsByCategory(cat)
      categoryFunctions.value = funcs || allFunctions.value.filter(f => f.category === cat)
    } catch (e) {
      categoryFunctions.value = allFunctions.value.filter(f => f.category === cat)
    }
  }
}

const addSubProgram = () => {
  subPrograms.value.push({
    name: `_子程序${subPrograms.value.length}`,
    returnType: '',
    isPublic: false,
    note: '',
    params: [],
    locals: [],
    codeBody: '',
    collapsed: false,
  })
  syncCodeFromVisual()
}

const removeSubProgram = (index: number) => {
  if (subPrograms.value.length > 1) {
    subPrograms.value.splice(index, 1)
    syncCodeFromVisual()
  }
}

const toggleCollapse = (index: number) => {
  subPrograms.value[index].collapsed = !subPrograms.value[index].collapsed
}

const addParam = (progIndex: number) => {
  subPrograms.value[progIndex].params.push({
    name: '', type: '整数型', isRef: false, isNullable: false, isArray: false, note: ''
  })
  syncCodeFromVisual()
}

const removeParam = (progIndex: number, paramIndex: number) => {
  subPrograms.value[progIndex].params.splice(paramIndex, 1)
  syncCodeFromVisual()
}

const addLocal = (progIndex: number) => {
  subPrograms.value[progIndex].locals.push({
    name: '', type: '整数型', isStatic: false, defaultValue: '', note: ''
  })
  syncCodeFromVisual()
}

const removeLocal = (progIndex: number, localIndex: number) => {
  subPrograms.value[progIndex].locals.splice(localIndex, 1)
  syncCodeFromVisual()
}

const insertFunctionToCode = (fn: ChineseFunction) => {
  const textarea = document.querySelector('.code-textarea') as HTMLTextAreaElement
  if (textarea) {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = textarea.value
    const insertText = fn.chineseName + ' ()'
    textarea.value = value.substring(0, start) + insertText + value.substring(end)
    textarea.selectionStart = textarea.selectionEnd = start + insertText.length
    codeContent.value = textarea.value
    highlightCode(textarea.value)
    syncCodeToVisual(textarea.value)
    textarea.focus()
  }
}

const validateCurrentCode = async () => {
  try {
    const errors = await ValidateCode(codeContent.value)
    codeErrors.value = errors || []
  } catch (e) {
    codeErrors.value = ['验证失败: ' + String(e)]
  }
}

const convertToGo = async () => {
  try {
    const codeRows = getAllCodeRows()
    const hasValidCode = codeRows.some((row: any) => row.code && row.code.trim().length > 0)

    if (!hasValidCode) {
      ElMessage.warning('请先编写代码内容再转换')
      return
    }

    ElMessage.info('正在转换为Go代码...')
    const result = await ConvertCodeRowsToGo(codeRows)
    console.log('生成的Go代码:', result)
    ElMessage.success('Go代码生成成功！请查看控制台输出')
  } catch (e: any) {
    console.error('转换失败详情:', e)
    ElMessage.error('转换失败: ' + (e?.message || String(e)))
  }
}

const getAllCodeRows = () => {
  const allRows: any[] = []
  subPrograms.value.forEach(prog => {
    allRows.push({
      event: prog.name || '未命名子程序',
      code: prog.codeBody || '',
    })
  })
  return allRows
}

defineExpose({
  subPrograms,
  getAllCodeRows,
})

function parseELangCode(fullCode: string): SubProgram[] {
  const programs: SubProgram[] = []

  const subRegex = /^\.子程序\s+(.+)$/gm
  const matches: { index: number; line: string }[] = []
  let match: RegExpExecArray | null
  while ((match = subRegex.exec(fullCode)) !== null) {
    matches.push({ index: match.index, line: match[1] })
  }

  if (matches.length === 0) {
    programs.push({
      name: '_启动子程序',
      returnType: '',
      isPublic: false,
      note: '',
      params: [],
      locals: [],
      codeBody: fullCode.trim(),
      collapsed: false,
    })
    return programs
  }

  for (let i = 0; i < matches.length; i++) {
    const startIdx = matches[i].index
    const endIdx = i + 1 < matches.length ? matches[i + 1].index : fullCode.length
    const block = fullCode.substring(startIdx, endIdx).trim()

    const prog = parseSubProgramBlock(block)
    programs.push(prog)
  }

  return programs
}

function parseSubProgramBlock(block: string): SubProgram {
  const lines = block.split('\n')

  const prog: SubProgram = {
    name: '_子程序',
    returnType: '',
    isPublic: false,
    note: '',
    params: [],
    locals: [],
    codeBody: '',
    collapsed: false,
  }

  let codeBodyLines: string[] = []
  let i = 0

  for (i = 0; i < lines.length; i++) {
    const trimmed = lines[i].trim()

    if (trimmed.startsWith('.子程序')) {
      const content = trimmed.substring(4).trim()
      const parts = parseCommaSeparated(content)
      prog.name = parts[0] || '_子程序'
      prog.returnType = parts[1] || ''

      for (let j = 2; j < parts.length; j++) {
        const p = parts[j].trim()
        if (p === '公开') {
          prog.isPublic = true
        } else if (p) {
          prog.note = p
        }
      }
      continue
    }

    if (trimmed.startsWith('.参数')) {
      const content = trimmed.substring(3).trim()
      if (content) {
        const param = parseParamLine(content)
        prog.params.push(param)
      }
      continue
    }

    if (trimmed.startsWith('.局部变量')) {
      const content = trimmed.substring(5).trim()
      if (content) {
        const local = parseLocalLine(content)
        prog.locals.push(local)
      }
      continue
    }

    codeBodyLines.push(lines[i])
  }

  prog.codeBody = codeBodyLines.join('\n').trim()
  return prog
}

function parseCommaSeparated(str: string): string[] {
  const result: string[] = []
  let current = ''
  let inQuote = false
  let quoteChar = ''

  for (let i = 0; i < str.length; i++) {
    const ch = str[i]
    if (inQuote) {
      current += ch
      if (ch === quoteChar) {
        inQuote = false
      }
    } else if (ch === '"' || ch === '\'') {
      inQuote = true
      quoteChar = ch
      current += ch
    } else if (ch === ',') {
      result.push(current.trim())
      current = ''
    } else {
      current += ch
    }
  }
  if (current.trim()) {
    result.push(current.trim())
  }
  return result
}

function parseParamLine(content: string): SubParam {
  const parts = parseCommaSeparated(content)
  const param: SubParam = {
    name: parts[0] || '',
    type: parts[1] || '整数型',
    isRef: false,
    isNullable: false,
    isArray: false,
    note: '',
  }

  for (let i = 2; i < parts.length; i++) {
    const p = parts[i].trim()
    if (p === '参考') param.isRef = true
    else if (p === '可空') param.isNullable = true
    else if (p === '数组') param.isArray = true
    else if (p) param.note = p
  }

  return param
}

function parseLocalLine(content: string): SubLocal {
  const parts = parseCommaSeparated(content)
  const local: SubLocal = {
    name: parts[0] || '',
    type: parts[1] || '整数型',
    isStatic: false,
    defaultValue: '',
    note: '',
  }

  for (let i = 2; i < parts.length; i++) {
    const p = parts[i].trim()
    if (p === '静态') local.isStatic = true
    else if (p.startsWith('"') || p.startsWith('\'')) local.defaultValue = p.replace(/^["']|["']$/g, '')
    else if (p && !local.defaultValue) local.defaultValue = p
    else if (p) local.note = p
  }

  return local
}

function generateELangCode(prog: SubProgram): string {
  let code = `.子程序 ${prog.name}`
  if (prog.returnType) code += `, ${prog.returnType}`
  if (prog.isPublic) code += `, 公开`
  if (prog.note) code += `, ${prog.note}`
  code += '\n'

  if (prog.params.length > 0) {
    for (const p of prog.params) {
      let line = `.参数 ${p.name}, ${p.type}`
      if (p.isRef) line += ', 参考'
      if (p.isNullable) line += ', 可空'
      if (p.isArray) line += ', 数组'
      if (p.note) line += `, ${p.note}`
      code += line + '\n'
    }
  }

  if (prog.locals.length > 0) {
    for (const l of prog.locals) {
      let line = `.局部变量 ${l.name}, ${l.type}`
      if (l.isStatic) line += ', 静态'
      if (l.defaultValue) line += `, "${l.defaultValue}"`
      if (l.note) line += `, ${l.note}`
      code += line + '\n'
    }
  }

  if (prog.codeBody.trim()) {
    code += '\n' + prog.codeBody.trim() + '\n'
  }

  return code
}

function generateFullCode(): string {
  return subPrograms.value.map(p => generateELangCode(p)).join('\n\n')
}

const syncCodeFromVisual = () => {
  if (isSyncing.value) return
  isSyncing.value = true

  const code = generateFullCode()
  codeContent.value = code
  highlightCode(code)

  nextTick(() => {
    isSyncing.value = false
  })
}

const syncCodeToVisual = (code: string) => {
  if (isSyncing.value) return
  isSyncing.value = true

  const parsed = parseELangCode(code)
  if (parsed.length > 0) {
    subPrograms.value = parsed
  }

  emit('code-rows-change', getAllCodeRows())

  nextTick(() => {
    isSyncing.value = false
  })
}

const onVisualFieldChange = () => {
  syncCodeFromVisual()
}

const onVisualCodeBodyInput = (progIndex: number, event: Event) => {
  const target = event.target as HTMLTextAreaElement
  const prog = subPrograms.value[progIndex]
  if (!prog) return
  prog.codeBody = target.value
  syncCodeFromVisual()
}

const scrollToSubProgram = (progName: string) => {
  nextTick(() => {
    const textarea = document.querySelector('.code-textarea') as HTMLTextAreaElement
    if (!textarea) return

    const lines = textarea.value.split('\n')
    let targetLine = 0
    for (let i = 0; i < lines.length; i++) {
      if (lines[i].includes(`.子程序 ${progName}`)) {
        targetLine = i
        break
      }
    }

    let pos = 0
    for (let i = 0; i < targetLine; i++) {
      pos += lines[i].length + 1
    }

    textarea.focus()
    textarea.setSelectionRange(pos, pos)

    const lineHeight = 20
    textarea.scrollTop = Math.max(0, (targetLine - 3) * lineHeight)
  })
}

const goKeywords = ['package', 'import', 'func', 'var', 'const', 'if', 'else', 'for', 'range', 'switch', 'case', 'default', 'break', 'continue', 'return', 'go', 'defer', 'select', 'chan', 'struct', 'interface', 'map', 'type', 'nil', 'true', 'false', 'make', 'new', 'append', 'len', 'cap', 'copy', 'delete', 'print', 'println', 'fmt', 'string', 'int', 'int8', 'int16', 'int32', 'int64', 'uint', 'uint8', 'uint16', 'uint32', 'uint64', 'float32', 'float64', 'complex64', 'complex128', 'bool', 'byte', 'rune', 'uintptr', 'error', 'iota']

const elangKeywords = ['.子程序', '.参数', '.局部变量', '.判断', '.如果', '.否则', '.否则如果', '.如果结束', '.判断结束', '.循环判断首', '.循环判断尾', '.计次循环首', '.计次循环尾', '.变量循环首', '.变量循环尾', '.循环', '.循环尾', '.返回', '.公开', '.私有', '.参考', '.数组', '.静态', '.真', '.假', '.空', '.默认', '.判断开始']

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function highlightCode(code: string): void {
  const lines = code.split('\n')
  const highlightedLines: string[] = []

  lines.forEach(line => {
    let escaped = escapeHtml(line)

    escaped = escaped.replace(/(\/\/.*)/g, '<span class="hl-comment">$1</span>')
    escaped = escaped.replace(/("[^"]*")/g, '<span class="hl-string">$1</span>')
    escaped = escaped.replace(/(`[^`]*`)/g, '<span class="hl-string">$1</span>')
    escaped = escaped.replace(/('.*?')/g, '<span class="hl-string">$1</span>')
    escaped = escaped.replace(/\b(\d+\.?\d*)\b/g, '<span class="hl-number">$1</span>')

    elangKeywords.forEach(kw => {
      const regex = new RegExp(`${escapeRegex(kw)}`, 'g')
      escaped = escaped.replace(regex, `<span class="hl-elang">${kw}</span>`)
    })

    goKeywords.forEach(kw => {
      const regex = new RegExp(`\\b${kw}\\b`, 'g')
      escaped = escaped.replace(regex, `<span class="hl-keyword">${kw}</span>`)
    })

    allFunctions.value.forEach(fn => {
      const regex = new RegExp(`${escapeRegex(fn.chineseName)}`, 'g')
      escaped = escaped.replace(regex, `<span class="hl-function" title="${fn.description}">${fn.chineseName}</span>`)
    })

    escaped = escaped.replace(/(＝|≠|≥|≤|＜|＞)/g, '<span class="hl-operator">$1</span>')

    highlightedLines.push(escaped || ' ')
  })

  highlightedCode.value = highlightedLines.join('\n')
}

function escapeRegex(string: string): string {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

const handleCodeInput = (event: Event) => {
  const target = event.target as HTMLTextAreaElement
  codeContent.value = target.value
  highlightCode(target.value)
  syncCodeToVisual(target.value)
}

const handleCodeKeydown = (event: KeyboardEvent) => {
  const textarea = event.target as HTMLTextAreaElement

  if (event.key === 'Tab') {
    event.preventDefault()
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = textarea.value
    textarea.value = value.substring(0, start) + '    ' + value.substring(end)
    textarea.selectionStart = textarea.selectionEnd = start + 4
    handleCodeInput({ target: textarea } as unknown as Event)
  }

  if (event.key === '{') {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = textarea.value
    textarea.value = value.substring(0, start) + '{}' + value.substring(end)
    textarea.selectionStart = textarea.selectionEnd = start + 1
    handleCodeInput({ target: textarea } as unknown as Event)
    event.preventDefault()
  }

  if (event.key === '(') {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = textarea.value
    textarea.value = value.substring(0, start) + '()' + value.substring(end)
    textarea.selectionStart = textarea.selectionEnd = start + 1
    handleCodeInput({ target: textarea } as unknown as Event)
    event.preventDefault()
  }

  if ((event.ctrlKey || event.metaKey) && event.key === 's') {
    event.preventDefault()
    validateCurrentCode()
  }
}

const scrollToLine = (lineNum: number) => {
  nextTick(() => {
    const textarea = document.querySelector('.code-textarea') as HTMLTextAreaElement
    if (!textarea) return

    const lines = textarea.value.split('\n')
    let pos = 0
    for (let i = 0; i < Math.min(lineNum - 1, lines.length); i++) {
      pos += lines[i].length + 1
    }

    textarea.focus()
    textarea.setSelectionRange(pos, pos)

    const lineHeight = 20
    textarea.scrollTop = (lineNum - 5) * lineHeight
  })
}

onMounted(() => {
  loadAllFunctions()
  if (subPrograms.value.length === 1 && !subPrograms.value[0].codeBody) {
    const parsed = parseELangCode(defaultCode)
    if (parsed.length > 0) {
      subPrograms.value = parsed
    }
  }
  syncCodeFromVisual()
})
</script>

<template>
  <div class="ecode-editor">
    <div class="editor-toolbar">
      <div class="toolbar-left">
        <el-radio-group v-model="editMode" size="small">
          <el-radio-button value="visual">
            <el-icon><Grid /></el-icon> 可视化
          </el-radio-button>
          <el-radio-button value="code">
            <el-icon><Document /></el-icon> 代码模式
          </el-radio-button>
        </el-radio-group>

        <el-button size="small" type="primary" @click="addSubProgram">
          <el-icon><Plus /></el-icon> 新建子程序
        </el-button>
      </div>

      <div class="toolbar-right">
        <el-button size="small" @click="showOutline = !showOutline" v-if="editMode === 'code'">
          <el-icon><List /></el-icon> 大纲
        </el-button>
        <el-button size="small" @click="showFunctionPanel = !showFunctionPanel">
          <el-icon><Collection /></el-icon> 函数库
        </el-button>
        <el-button size="small" @click="validateCurrentCode" title="Ctrl+S 保存验证">
          <el-icon><Check /></el-icon> 验证
        </el-button>
        <el-tooltip content="转换为Go代码">
          <el-button size="small" @click="convertToGo">
            <el-icon><Cpu /></el-icon>
          </el-button>
        </el-tooltip>
      </div>
    </div>

    <div class="editor-body" :class="{ 'with-panel': showFunctionPanel }">
      <div class="editor-main-area">

      <!-- 左侧大纲导航 -->
      <div v-if="showOutline && editMode === 'code'" class="outline-sidebar">
        <div class="outline-header">
          <span>📑 大纲</span>
          <el-icon class="outline-toggle" @click="showOutline = false"><Close /></el-icon>
        </div>
        <div class="outline-list">
          <div
            v-for="(prog, index) in subPrograms"
            :key="index"
            class="outline-item"
            @click="scrollToSubProgram(prog.name)"
            :title="prog.name"
          >
            <span class="outline-icon">{{ prog.isPublic ? '🌐' : '📦' }}</span>
            <span class="outline-name">{{ prog.name }}</span>
            <span v-if="prog.returnType" class="outline-return">{{ prog.returnType }}</span>
          </div>
        </div>
      </div>

      <!-- ==================== 代码模式 ==================== -->
      <div v-if="editMode === 'code'" class="code-mode-container">
        <div class="code-editor-wrapper">
          <div class="line-numbers">
            <div
              v-for="(line, idx) in lineCount"
              :key="idx"
              class="line-number"
              @click="scrollToLine(idx + 1)"
            >
              {{ idx + 1 }}
            </div>
          </div>

          <div class="code-editor-main">
            <pre class="code-highlight" v-html="highlightedCode"></pre>
            <textarea
              class="code-textarea"
              v-model="codeContent"
              spellcheck="false"
              @input="handleCodeInput"
              @keydown="handleCodeKeydown"
              placeholder="在此编写易语言代码...&#10;&#10;.子程序 函数名, 返回类型, 公开, 备注&#10;.参数 参数名, 整数型, 参考, 备注&#10;.局部变量 变量名, 整数型, 静态, &quot;0&quot;, 备注&#10;&#10;.如果 (条件)&#10;    调试输出 (&quot;内容&quot;)&#10;.如果结束&#10;&#10;按 Tab 缩进，Ctrl+S 验证代码"
            ></textarea>
          </div>
        </div>

        <div class="code-status-bar">
          <span class="status-left">
            <span class="status-item">行 {{ lineCount }}</span>
            <span class="status-item">UTF-8</span>
            <span class="status-item">{{ subPrograms.length }} 个子程序</span>
          </span>
          <span class="status-right">
            <span class="status-item mode-badge">E语言</span>
            <span class="status-item">空格: 4</span>
          </span>
        </div>
      </div>

      <!-- ==================== 可视化模式 ==================== -->
      <div v-else class="visual-mode-container">
        <div class="programs-area">
          <div
            v-for="(prog, progIndex) in subPrograms"
            :key="progIndex"
            class="sub-program-panel"
            :class="{ collapsed: prog.collapsed }"
          >
            <div class="panel-header-bar" @click="toggleCollapse(progIndex)">
              <span class="collapse-btn">{{ prog.collapsed ? '▶' : '▼' }}</span>
              <span class="panel-title">{{ prog.name || '未命名' }}</span>
              <span v-if="prog.returnType" class="return-badge">{{ prog.returnType }}</span>
              <span v-if="prog.isPublic" class="public-badge">公开</span>
              <el-button
                size="small"
                type="danger"
                circle
                class="del-btn"
                @click.stop="removeSubProgram(progIndex)"
                :disabled="subPrograms.length <= 1"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>

            <div v-show="!prog.collapsed" class="panel-body">

              <!-- 子程序表头 -->
              <table class="def-table">
                <thead>
                  <tr>
                    <th class="col-label">.子程序</th>
                    <th class="col-name">名称</th>
                    <th class="col-type">返回类型</th>
                    <th class="col-flags">属性</th>
                    <th class="col-note">备注</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td class="col-label"><span class="dot-cmd">.子程序</span></td>
                    <td class="col-name">
                      <el-input
                        v-model="prog.name"
                        size="small"
                        placeholder="子程序名"
                        @input="onVisualFieldChange()"
                      />
                    </td>
                    <td class="col-type">
                      <el-select
                        v-model="prog.returnType"
                        size="small"
                        clearable
                        placeholder="无"
                        @change="onVisualFieldChange()"
                      >
                        <el-option v-for="t in dataTypes" :key="t" :label="t" :value="t" />
                      </el-select>
                    </td>
                    <td class="col-flags">
                      <el-checkbox
                        v-model="prog.isPublic"
                        size="small"
                        @change="onVisualFieldChange()"
                      >公开</el-checkbox>
                    </td>
                    <td class="col-note">
                      <el-input
                        v-model="prog.note"
                        size="small"
                        placeholder="备注"
                        @input="onVisualFieldChange()"
                      />
                    </td>
                  </tr>
                </tbody>
              </table>

              <!-- 参数表 -->
              <div class="section-block">
                <div class="section-label">
                  <span class="dot-cmd">.参数</span>
                  <el-button size="small" text type="primary" @click="addParam(progIndex)">+ 添加参数</el-button>
                </div>
                <table class="def-table" v-if="prog.params.length > 0">
                  <thead>
                    <tr>
                      <th class="col-name">名称</th>
                      <th class="col-type">类型</th>
                      <th class="col-flags">参考</th>
                      <th class="col-flags">可空</th>
                      <th class="col-flags">数组</th>
                      <th class="col-note">备注</th>
                      <th class="col-action"></th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(param, pi) in prog.params" :key="pi">
                      <td class="col-name">
                        <el-input v-model="param.name" size="small" placeholder="参数名" @input="onVisualFieldChange()" />
                      </td>
                      <td class="col-type">
                        <el-select v-model="param.type" size="small" @change="onVisualFieldChange()">
                          <el-option v-for="t in dataTypes" :key="t" :label="t" :value="t" />
                        </el-select>
                      </td>
                      <td class="col-flags">
                        <el-checkbox v-model="param.isRef" size="small" @change="onVisualFieldChange()" />
                      </td>
                      <td class="col-flags">
                        <el-checkbox v-model="param.isNullable" size="small" @change="onVisualFieldChange()" />
                      </td>
                      <td class="col-flags">
                        <el-checkbox v-model="param.isArray" size="small" @change="onVisualFieldChange()" />
                      </td>
                      <td class="col-note">
                        <el-input v-model="param.note" size="small" placeholder="备注" @input="onVisualFieldChange()" />
                      </td>
                      <td class="col-action">
                        <el-button size="small" type="danger" circle @click="removeParam(progIndex, pi)">
                          <el-icon><Delete /></el-icon>
                        </el-button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- 局部变量表 -->
              <div class="section-block">
                <div class="section-label">
                  <span class="dot-cmd">.局部变量</span>
                  <el-button size="small" text type="primary" @click="addLocal(progIndex)">+ 添加变量</el-button>
                </div>
                <table class="def-table" v-if="prog.locals.length > 0">
                  <thead>
                    <tr>
                      <th class="col-name">名称</th>
                      <th class="col-type">类型</th>
                      <th class="col-flags">静态</th>
                      <th class="col-default">默认值</th>
                      <th class="col-note">备注</th>
                      <th class="col-action"></th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(local, li) in prog.locals" :key="li">
                      <td class="col-name">
                        <el-input v-model="local.name" size="small" placeholder="变量名" @input="onVisualFieldChange()" />
                      </td>
                      <td class="col-type">
                        <el-select v-model="local.type" size="small" @change="onVisualFieldChange()">
                          <el-option v-for="t in dataTypes" :key="t" :label="t" :value="t" />
                        </el-select>
                      </td>
                      <td class="col-flags">
                        <el-checkbox v-model="local.isStatic" size="small" @change="onVisualFieldChange()" />
                      </td>
                      <td class="col-default">
                        <el-input v-model="local.defaultValue" size="small" placeholder="默认值" @input="onVisualFieldChange()" />
                      </td>
                      <td class="col-note">
                        <el-input v-model="local.note" size="small" placeholder="备注" @input="onVisualFieldChange()" />
                      </td>
                      <td class="col-action">
                        <el-button size="small" type="danger" circle @click="removeLocal(progIndex, li)">
                          <el-icon><Delete /></el-icon>
                        </el-button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- 代码体 - 纯文本编辑区 -->
              <div class="section-block">
                <div class="section-label">
                  <span>📝 代码体</span>
                  <span class="hint-text">支持 .如果 .否则 .循环 等流程控制及函数调用</span>
                </div>
                <div class="visual-code-editor">
                  <textarea
                    class="visual-code-textarea"
                    :value="prog.codeBody"
                    spellcheck="false"
                    @input="(e: Event) => onVisualCodeBodyInput(progIndex, e)"
                    placeholder="在此编写代码体...&#10;&#10;.如果 (条件)&#10;    调试输出 (&quot;内容&quot;)&#10;.否则&#10;    返回 (0)&#10;.如果结束"
                  ></textarea>
                </div>
              </div>

            </div>
          </div>
        </div>
      </div>
      </div>

      <!-- 函数面板 -->
      <div v-if="showFunctionPanel" class="function-side-panel">
        <div class="panel-header">
          <h4>📚 函数库</h4>
          <el-icon class="close-panel" @click="showFunctionPanel = false"><Close /></el-icon>
        </div>

        <div class="side-panel-tabs">
          <span
            v-for="cat in categories"
            :key="cat"
            :class="{ active: selectedCategory === cat }"
            @click="selectCategory(cat)"
          >{{ cat }}</span>
        </div>

        <div class="side-panel-list">
          <div
            v-for="fn in (selectedCategory ? categoryFunctions : allFunctions)"
            :key="fn.chineseName"
            class="fn-item"
            @click="insertFunctionToCode(fn)"
          >
            <div>
              <div class="fn-name">{{ fn.chineseName }}</div>
              <div class="fn-desc">{{ fn.description }}</div>
            </div>
            <span class="fn-return">{{ fn.returnType || 'void' }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="codeErrors.length > 0" class="error-panel">
      <div class="error-header">
        <el-icon><WarningFilled /></el-icon> 代码错误 ({{ codeErrors.length }})
      </div>
      <div v-for="(err, i) in codeErrors" :key="i" class="error-item">{{ err }}</div>
    </div>
  </div>
</template>

<style scoped>
.ecode-editor {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--editor-bg);
  overflow: hidden;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 12px;
  background-color: var(--editor-toolbar-bg);
  border-bottom: 1px solid var(--editor-toolbar-border);
  flex-shrink: 0;
}

.toolbar-left, .toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.editor-toolbar :deep(.el-radio-group) {
  background-color: var(--editor-input-bg);
}

.editor-toolbar :deep(.el-radio-button__inner) {
  background-color: var(--editor-input-bg);
  border-color: var(--editor-input-border);
  color: var(--editor-input-text);
}

.editor-toolbar :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background-color: var(--color-primary);
  border-color: var(--color-primary-hover);
  box-shadow: none;
}

.editor-toolbar :deep(.el-button) {
  background-color: var(--editor-input-bg) !important;
  border-color: var(--editor-input-border) !important;
  color: var(--editor-input-text) !important;
}

.editor-toolbar :deep(.el-button--primary) {
  background-color: var(--color-primary) !important;
  border-color: var(--color-primary-hover) !important;
}

.editor-toolbar :deep(.el-button:hover) {
  background-color: var(--editor-toolbar-btn-hover) !important;
  border-color: var(--border-dark) !important;
}

.editor-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-main-area {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-body.with-panel .editor-main-area {
  flex: 1;
}

/* ========== 左侧大纲导航 ========== */
.outline-sidebar {
  width: 180px;
  background-color: var(--bg-secondary);
  border-right: 1px solid var(--editor-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.outline-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 1px solid var(--editor-border);
  flex-shrink: 0;
}

.outline-toggle {
  cursor: pointer;
  font-size: 14px;
  color: var(--text-secondary);
}

.outline-toggle:hover { color: var(--text-primary); }

.outline-list {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.outline-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 10px;
  cursor: pointer;
  font-size: 12px;
  color: var(--text-regular);
  transition: background-color 0.12s;
  border-left: 2px solid transparent;
}

.outline-item:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.outline-icon { font-size: 11px; flex-shrink: 0; }

.outline-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.outline-return {
  font-size: 10px;
  color: var(--hl-type);
  background-color: rgba(78, 201, 176, 0.12);
  padding: 0 4px;
  border-radius: 2px;
  flex-shrink: 0;
}

/* ========== 代码模式 ========== */
.code-mode-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.code-editor-wrapper {
  flex: 1;
  display: flex;
  overflow: hidden;
  background-color: var(--editor-bg);
}

.line-numbers {
  width: 48px;
  background-color: var(--editor-gutter-bg);
  border-right: 1px solid var(--editor-border);
  padding: 8px 0;
  overflow-y: auto;
  user-select: none;
  text-align: right;
  flex-shrink: 0;
}

.line-number {
  height: 20px;
  line-height: 20px;
  padding-right: 12px;
  font-size: 12px;
  color: var(--editor-line-number);
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  cursor: pointer;
}

.line-number:hover { color: var(--text-primary); }

.code-editor-main {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.code-highlight {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  margin: 0;
  padding: 8px 12px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 20px;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: var(--editor-text);
  background: transparent;
  pointer-events: none;
  overflow-y: auto;
  tab-size: 4;
}

.code-textarea {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 8px 12px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 20px;
  color: transparent;
  background: transparent;
  caret-color: var(--editor-cursor);
  border: none;
  outline: none;
  resize: none;
  white-space: pre-wrap;
  word-wrap: break-word;
  tab-size: 4;
  overflow-y: auto;
}

.code-textarea::selection {
  background-color: var(--editor-selection-bg);
  color: white;
}

.hl-comment { color: var(--hl-comment); font-style: italic; }
.hl-string { color: var(--hl-string); }
.hl-number { color: var(--hl-number); }
.hl-keyword { color: var(--hl-keyword); }
.hl-elang { color: var(--hl-elang); font-weight: 600; }
.hl-function { color: var(--hl-function); cursor: pointer; border-bottom: 1px dotted var(--hl-function); }
.hl-operator { color: var(--hl-operator); }

.code-status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2px 12px;
  background-color: var(--color-primary);
  font-size: 11px;
  color: var(--text-inverse);
  flex-shrink: 0;
}

.status-left, .status-right {
  display: flex;
  gap: 16px;
}

.status-item { opacity: 0.9; }

.mode-badge {
  background-color: rgba(255,255,255,0.2);
  padding: 1px 6px;
  border-radius: 3px;
  font-weight: 500;
}

/* ========== 可视化模式 ========== */
.visual-mode-container {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.programs-area {
  max-width: 100%;
}

.sub-program-panel {
  margin-bottom: 8px;
  border: 1px solid var(--editor-border);
  border-radius: 6px;
  background-color: var(--bg-secondary);
  overflow: hidden;
}

.panel-header-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background-color: var(--editor-toolbar-bg);
  cursor: pointer;
  transition: background-color 0.15s;
  border-bottom: 1px solid var(--editor-border);
}

.panel-header-bar:hover {
  background-color: var(--bg-hover);
}

.collapse-btn {
  font-size: 10px;
  color: var(--editor-gutter-text);
  width: 14px;
  user-select: none;
}

.panel-title {
  font-size: 13px;
  color: var(--hl-function);
  font-weight: 600;
}

.return-badge {
  font-size: 10px;
  color: var(--hl-type);
  background-color: rgba(78,201,176,0.15);
  padding: 1px 6px;
  border-radius: 3px;
}

.public-badge {
  font-size: 10px;
  color: var(--hl-keyword);
  background-color: rgba(86,156,214,0.15);
  padding: 1px 6px;
  border-radius: 3px;
  border: 1px solid var(--hl-keyword);
}

.del-btn {
  margin-left: auto !important;
  opacity: 0;
  transition: opacity 0.15s;
}

.panel-header-bar:hover .del-btn { opacity: 1; }

.panel-body {
  padding: 8px;
}

/* 定义表格 */
.def-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  margin-bottom: 4px;
}

.def-table th {
  background-color: var(--table-header-bg);
  padding: 4px 8px;
  text-align: left;
  font-weight: 600;
  color: var(--text-secondary);
  border: 1px solid var(--table-border);
  font-size: 11px;
  white-space: nowrap;
}

.def-table td {
  padding: 3px 6px;
  border: 1px solid var(--table-border);
  background-color: var(--table-bg);
  vertical-align: middle;
}

.col-label { width: 80px; }
.col-name { min-width: 120px; }
.col-type { width: 110px; }
.col-flags { width: 55px; text-align: center; }
.col-note { min-width: 100px; }
.col-default { width: 100px; }
.col-action { width: 36px; text-align: center; }

.dot-cmd {
  color: var(--hl-elang);
  font-weight: 700;
  font-size: 12px;
}

.def-table :deep(.el-input__wrapper) {
  background-color: var(--input-bg) !important;
  box-shadow: none !important;
  border: 1px solid var(--input-border) !important;
  padding: 1px 8px !important;
}

.def-table :deep(.el-input__inner) {
  color: var(--input-text) !important;
  font-size: 12px !important;
  height: 26px !important;
}

.def-table :deep(.el-select .el-input__wrapper) {
  background-color: var(--input-bg) !important;
  box-shadow: none !important;
  border: 1px solid var(--input-border) !important;
}

.def-table :deep(.el-checkbox) {
  height: 20px;
  --el-checkbox-text-color: var(--text-primary);
  --el-checkbox-checked-bg-color: var(--color-primary);
  --el-checkbox-checked-border-color: var(--color-primary);
}

.def-table :deep(.el-checkbox__label) {
  font-size: 11px;
  padding-left: 4px;
}

/* 区块 */
.section-block {
  margin-top: 8px;
}

.section-label {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 12px;
  color: var(--text-primary);
}

.section-label .el-button {
  font-size: 11px !important;
  padding: 2px 8px !important;
  height: 24px !important;
}

.hint-text {
  font-size: 10px;
  color: var(--text-secondary);
  font-style: italic;
}

/* 可视化模式代码编辑区 */
.visual-code-editor {
  border: 1px solid var(--editor-border);
  border-radius: 4px;
  overflow: hidden;
  background-color: var(--editor-bg);
}

.visual-code-textarea {
  width: 100%;
  min-height: 200px;
  padding: 8px 12px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 20px;
  color: var(--editor-text);
  background-color: var(--editor-bg);
  border: none;
  outline: none;
  resize: vertical;
  white-space: pre-wrap;
  word-wrap: break-word;
  tab-size: 4;
}

.visual-code-textarea::placeholder {
  color: var(--text-secondary);
}

.visual-code-textarea:focus {
  background-color: var(--bg-secondary);
}

/* 函数侧边面板 */
.function-side-panel {
  width: 260px;
  border-left: 1px solid var(--editor-border);
  display: flex;
  flex-direction: column;
  background-color: var(--bg-secondary);
  flex-shrink: 0;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid var(--editor-border);
}

.panel-header h4 {
  margin: 0;
  font-size: 13px;
  color: var(--text-primary);
}

.close-panel {
  cursor: pointer;
  color: var(--editor-gutter-text);
  font-size: 16px;
}

.close-panel:hover { color: var(--text-primary); }

.side-panel-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  padding: 6px;
  border-bottom: 1px solid var(--editor-border);
}

.side-panel-tabs span {
  font-size: 11px;
  padding: 3px 8px;
  border-radius: 3px;
  cursor: pointer;
  color: var(--text-primary);
  white-space: nowrap;
}

.side-panel-tabs span:hover { background-color: var(--bg-hover); }

.side-panel-tabs span.active {
  background-color: var(--color-primary);
  color: var(--text-inverse);
}

.side-panel-list {
  flex: 1;
  overflow-y: auto;
  padding: 4px;
}

.fn-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  cursor: pointer;
  border-radius: 3px;
  margin-bottom: 2px;
}

.fn-item:hover { background-color: var(--bg-selected); }

.fn-name {
  color: var(--hl-function);
  font-weight: 500;
  font-size: 12px;
}

.fn-desc {
  color: var(--text-secondary);
  font-size: 10px;
  margin-top: 2px;
}

.fn-return {
  font-size: 10px;
  color: var(--hl-type);
  background-color: rgba(78, 201, 176, 0.15);
  padding: 1px 6px;
  border-radius: 2px;
  flex-shrink: 0;
}

.error-panel {
  background-color: var(--color-danger-light);
  border-top: 1px solid var(--color-danger);
  padding: 8px 12px;
  max-height: 150px;
  overflow-y: auto;
  flex-shrink: 0;
}

.error-header {
  color: var(--color-danger);
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.error-item {
  color: var(--color-danger);
  font-size: 11px;
  padding: 3px 0;
  padding-left: 16px;
  opacity: 0.85;
}
</style>