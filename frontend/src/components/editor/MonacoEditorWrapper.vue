<script lang="ts" setup>
import { ref, onMounted, onUnmounted, watch, nextTick, computed } from 'vue'
import * as monaco from 'monaco-editor'
import { ElMessage } from 'element-plus'
import { useTheme } from '../../stores/theme'

const { currentTheme } = useTheme()

const getDefaultTheme = () => currentTheme.value === 'dark' ? 'vs-dark' : 'vs'

const props = defineProps<{
  modelValue: string
  language?: string
  theme?: string
  readOnly?: boolean
  minimap?: boolean
  fontSize?: number
  options?: Record<string, any>
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
  (e: 'cursor-change', position: { lineNumber: number; column: number }): void
}>()

const editorContainer = ref<HTMLElement | null>(null)
let editorInstance: monaco.editor.IStandaloneCodeEditor | null = null
let gotoLineDialog: HTMLDivElement | null = null
let findWidgetVisible = false

const cursorPosition = computed(() => {
  if (!editorInstance) return { lineNumber: 1, column: 1 }
  const pos = editorInstance.getPosition()
  return pos || { lineNumber: 1, column: 1 }
})

const lineCount = computed(() => {
  if (!editorInstance) return 0
  return editorInstance.getModel()?.getLineCount() || 0
})

const eLangCompletions = [
  { label: '.子程序', insertText: '.子程序 ${1:名称}, ${2:整数型}, ${3:公开}, ${4:备注}\n\t\n.返回 (${5:0})', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '定义子程序', documentation: '定义一个子程序，指定名称、返回类型、是否公开、备注' },
  { label: '.参数', insertText: '.参数 ${1:参数名}, ${2:整数型}, ${3:参考}, ${4:可空}, ${5:备注}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '定义参数', documentation: '定义子程序参数' },
  { label: '.局部变量', insertText: '.局部变量 ${1:变量名}, ${2:整数型}, ${3:静态}, "${4:默认值}", ${5:备注}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '定义局部变量', documentation: '定义局部变量' },
  { label: '.如果', insertText: '.如果 (${1:条件})\n\t${2}\n.如果结束', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '条件判断', documentation: '如果条件成立则执行代码块' },
  { label: '.否则', insertText: '.否则\n\t${1}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '否则分支', documentation: '条件不成立时执行的代码块' },
  { label: '.如果真', insertText: '.如果真 (${1:条件})\n\t${2}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '单分支判断', documentation: '如果条件成立则执行，无否则分支' },
  { label: '.判断', insertText: '.判断 (${1:条件})\n.默认\n\t${2}\n.判断结束', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '多分支判断', documentation: '多条件分支判断' },
  { label: '.计次循环', insertText: '.计次循环 (${1:变量}, ${2:起始值}, ${3:结束值})\n\t${4}\n.计次循环尾 ()', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '计次循环', documentation: '按次数循环执行' },
  { label: '.变量循环', insertText: '.变量循环 (${1:变量}, ${2:起始值}, ${3:结束值}, ${4:步长})\n\t${5}\n.变量循环尾 ()', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '变量循环', documentation: '按变量步长循环' },
  { label: '.循环判断', insertText: '.循环判断 ()\n\t${1}\n.循环判断尾 ()', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '循环判断', documentation: '先循环后判断' },
  { label: '.判断循环', insertText: '.判断循环 (${1:条件})\n\t${2}\n.判断循环尾 ()', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '判断循环', documentation: '先判断后循环' },
  { label: '.返回', insertText: '.返回 (${1:返回值})', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet, detail: '返回', documentation: '从子程序返回' },
]

const providerDisposables: monaco.IDisposable[] = []

function setupChineseContextMenu() {
  if (!editorInstance) return

  editorInstance.addAction({
    id: 'editor-cut',
    label: '剪切',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyX],
    contextMenuGroupId: '9_cutcopypaste',
    contextMenuOrder: 1,
    run: (ed) => {
      const selection = ed.getSelection()
      if (selection && !selection.isEmpty()) {
        const text = ed.getModel()?.getValueInRange(selection) || ''
        navigator.clipboard.writeText(text).catch(() => {})
        ed.executeEdits('cut', [{ range: selection, text: '' }])
      }
    }
  })

  editorInstance.addAction({
    id: 'editor-copy',
    label: '复制',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyC],
    contextMenuGroupId: '9_cutcopypaste',
    contextMenuOrder: 2,
    run: (ed) => {
      const selection = ed.getSelection()
      if (selection && !selection.isEmpty()) {
        const text = ed.getModel()?.getValueInRange(selection) || ''
        navigator.clipboard.writeText(text).catch(() => {})
      }
    }
  })

  editorInstance.addAction({
    id: 'editor-paste',
    label: '粘贴',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyV],
    contextMenuGroupId: '9_cutcopypaste',
    contextMenuOrder: 3,
    run: async (ed) => {
      try {
        const text = await navigator.clipboard.readText()
        if (text) {
          const selection = ed.getSelection()
          if (selection) {
            ed.executeEdits('paste', [{ range: selection, text }])
          }
        }
      } catch {}
    }
  })

  editorInstance.addAction({
    id: 'editor-select-all',
    label: '全选',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyA],
    contextMenuGroupId: '9_cutcopypaste',
    contextMenuOrder: 4,
    run: (ed) => {
      const model = ed.getModel()
      if (model) {
        ed.setSelection(new monaco.Range(1, 1, model.getLineCount(), model.getLineMaxColumn(model.getLineCount())))
      }
    }
  })

  editorInstance.addAction({
    id: 'editor-find',
    label: '查找',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyF],
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 1,
    run: (ed) => ed.getAction('actions.find')?.run()
  })

  editorInstance.addAction({
    id: 'editor-replace',
    label: '替换',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyH],
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2,
    run: (ed) => ed.getAction('editor.action.startFindReplaceAction')?.run()
  })

  editorInstance.addAction({
    id: 'editor-goto-line',
    label: '转到行...',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyG],
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 3,
    run: () => showGoToLineDialog()
  })

  editorInstance.addAction({
    id: 'editor-format',
    label: '格式化文档',
    keybindings: [monaco.KeyMod.Shift | monaco.KeyMod.Alt | monaco.KeyCode.KeyF],
    contextMenuGroupId: 'modification',
    contextMenuOrder: 1,
    run: (ed) => ed.getAction('editor.action.formatDocument')?.run()
  })

  editorInstance.addAction({
    id: 'editor-comment-line',
    label: '注释/取消注释',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.Slash],
    contextMenuGroupId: 'modification',
    contextMenuOrder: 2,
    run: (ed) => ed.getAction('editor.action.commentLine')?.run()
  })

  editorInstance.addAction({
    id: 'editor-undo',
    label: '撤销',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyZ],
    contextMenuGroupId: 'modification',
    contextMenuOrder: 3,
    run: (ed) => ed.trigger('keyboard', 'undo', null)
  })

  editorInstance.addAction({
    id: 'editor-redo',
    label: '重做',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyMod.Shift | monaco.KeyCode.KeyZ],
    contextMenuGroupId: 'modification',
    contextMenuOrder: 4,
    run: (ed) => ed.trigger('keyboard', 'redo', null)
  })
}

function showGoToLineDialog() {
  if (!editorInstance) return
  const maxLine = editorInstance.getModel()?.getLineCount() || 1

  ElMessage({
    message: `转到行 (1-${maxLine})`,
    type: 'info',
    duration: 0,
    showClose: true,
    customClass: 'goto-line-message',
    onClose: () => {
      gotoLineDialog = null
    }
  })
}

function setupELangCompletion() {
  const disposable = monaco.languages.registerCompletionItemProvider('go', {
    triggerCharacters: ['.'],
    provideCompletionItems: (model, position) => {
      const textUntilPosition = model.getValueInRange({
        startLineNumber: position.lineNumber,
        startColumn: 1,
        endLineNumber: position.lineNumber,
        endColumn: position.column,
      })

      if (textUntilPosition.startsWith('.')) {
        const word = model.getWordUntilPosition(position)
        const range = {
          startLineNumber: position.lineNumber,
          endLineNumber: position.lineNumber,
          startColumn: word.startColumn,
          endColumn: word.endColumn,
        }
        return {
          suggestions: eLangCompletions.map(item => ({ ...item, range })) as monaco.languages.CompletionItem[],
        }
      }

      return { suggestions: [] }
    },
  })
  providerDisposables.push(disposable)
}

function setupGoKeywordCompletion() {
  const goKeywords = [
    { label: 'package', insertText: 'package ${1:main}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Keyword },
    { label: 'import', insertText: 'import (\n\t"${1}"\n)', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Keyword },
    { label: 'func', insertText: 'func ${1:name}(${2}) ${3} {\n\t${4}\n}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Keyword },
    { label: 'type', insertText: 'type ${1:Name} struct {\n\t${2}\n}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Keyword },
    { label: 'struct', insertText: 'struct {\n\t${1}\n}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Keyword },
    { label: 'interface', insertText: 'interface {\n\t${1}\n}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Keyword },
    { label: 'if', insertText: 'if ${1:condition} {\n\t${2}\n}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet },
    { label: 'for', insertText: 'for ${1:i} := 0; ${1:i} < ${2:n}; ${1:i}++ {\n\t${3}\n}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet },
    { label: 'switch', insertText: 'switch ${1:expr} {\ncase ${2}:\n\t${3}\ndefault:\n\t${4}\n}', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Snippet },
    { label: 'go', insertText: 'go func() {\n\t${1}\n}()', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Keyword },
    { label: 'defer', insertText: 'defer ${1:func}()', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Keyword },
    { label: 'fmt.Println', insertText: 'fmt.Println("${1}")', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Function, detail: '打印输出' },
    { label: 'fmt.Printf', insertText: 'fmt.Printf("${1:%v}", ${2})', insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet, kind: monaco.languages.CompletionItemKind.Function, detail: '格式化输出' },
  ]

  const disposable = monaco.languages.registerCompletionItemProvider('go', {
    provideCompletionItems: () => {
      return { suggestions: goKeywords as monaco.languages.CompletionItem[] }
    },
  })
  providerDisposables.push(disposable)
}

function setupPinyinCompletion() {
  const pinyinMap: { pinyin: string; chinese: string; goFunc: string }[] = [
    { pinyin: 'shuchu', chinese: '输出', goFunc: 'Print' },
    { pinyin: 'dayin', chinese: '打印', goFunc: 'Print' },
    { pinyin: 'wenjian', chinese: '文件', goFunc: 'File' },
    { pinyin: 'duqu', chinese: '读取', goFunc: 'Read' },
    { pinyin: 'xieru', chinese: '写入', goFunc: 'Write' },
    { pinyin: 'shanchu', chinese: '删除', goFunc: 'Remove' },
    { pinyin: 'chuangjian', chinese: '创建', goFunc: 'Mkdir' },
    { pinyin: 'kaobei', chinese: '拷贝', goFunc: 'Copy' },
    { pinyin: 'yidong', chinese: '移动', goFunc: 'Move' },
    { pinyin: 'bianli', chinese: '遍历', goFunc: 'Walk' },
    { pinyin: 'guolv', chinese: '过滤', goFunc: 'Filter' },
    { pinyin: 'paixu', chinese: '排序', goFunc: 'Sort' },
    { pinyin: 'lianjie', chinese: '连接', goFunc: 'Join' },
    { pinyin: 'fenge', chinese: '分割', goFunc: 'Split' },
    { pinyin: 'tihuan', chinese: '替换', goFunc: 'Replace' },
    { pinyin: 'pipei', chinese: '匹配', goFunc: 'Match' },
    { pinyin: 'geshihua', chinese: '格式化', goFunc: 'Format' },
    { pinyin: 'bianma', chinese: '编码', goFunc: 'Encode' },
    { pinyin: 'jiema', chinese: '解码', goFunc: 'Decode' },
    { pinyin: 'yasuo', chinese: '压缩', goFunc: 'Compress' },
    { pinyin: 'jieya', chinese: '解压', goFunc: 'Decompress' },
    { pinyin: 'fasong', chinese: '发送', goFunc: 'Send' },
    { pinyin: 'jieshou', chinese: '接收', goFunc: 'Receive' },
    { pinyin: 'qingqiu', chinese: '请求', goFunc: 'Request' },
    { pinyin: 'xiangying', chinese: '响应', goFunc: 'Response' },
    { pinyin: 'yichang', chinese: '异常', goFunc: 'Error' },
    { pinyin: 'cuowu', chinese: '错误', goFunc: 'Error' },
    { pinyin: 'rizhi', chinese: '日志', goFunc: 'Log' },
    { pinyin: 'yanshi', chinese: '延时', goFunc: 'Sleep' },
    { pinyin: 'dingshi', chinese: '定时', goFunc: 'Timer' },
    { pinyin: 'xiancheng', chinese: '线程', goFunc: 'Goroutine' },
    { pinyin: 'xiecheng', chinese: '协程', goFunc: 'Goroutine' },
    { pinyin: 'tongdao', chinese: '通道', goFunc: 'Channel' },
    { pinyin: 'suoding', chinese: '锁定', goFunc: 'Lock' },
    { pinyin: 'jiesuo', chinese: '解锁', goFunc: 'Unlock' },
    { pinyin: 'duilie', chinese: '队列', goFunc: 'Queue' },
    { pinyin: 'zhan', chinese: '栈', goFunc: 'Stack' },
    { pinyin: 'haxibiao', chinese: '哈希表', goFunc: 'Map' },
    { pinyin: 'shuzu', chinese: '数组', goFunc: 'Array' },
    { pinyin: 'juzheng', chinese: '矩阵', goFunc: 'Matrix' },
    { pinyin: 'zifuchuan', chinese: '字符串', goFunc: 'String' },
    { pinyin: 'zhengshu', chinese: '整数', goFunc: 'Int' },
    { pinyin: 'fudianshu', chinese: '浮点数', goFunc: 'Float' },
    { pinyin: 'buer', chinese: '布尔', goFunc: 'Bool' },
    { pinyin: 'jiekou', chinese: '接口', goFunc: 'Interface' },
    { pinyin: 'jiegou', chinese: '结构', goFunc: 'Struct' },
    { pinyin: 'hanshu', chinese: '函数', goFunc: 'Func' },
  ]

  const disposable = monaco.languages.registerCompletionItemProvider('go', {
    triggerCharacters: [...'abcdefghijklmnopqrstuvwxyz'],
    provideCompletionItems: (model, position) => {
      const word = model.getWordUntilPosition(position)
      const prefix = word.word.toLowerCase()
      if (prefix.length < 2) return { suggestions: [] }

      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn,
      }

      const matches = pinyinMap.filter(item =>
        item.pinyin.startsWith(prefix) || item.chinese.startsWith(prefix)
      )

      if (matches.length === 0) return { suggestions: [] }

      const suggestions = matches.map(item => ({
        label: item.chinese,
        insertText: item.chinese,
        kind: monaco.languages.CompletionItemKind.Function,
        detail: `拼音: ${item.pinyin} → Go: ${item.goFunc}`,
        documentation: `中文关键词: ${item.chinese}\n对应 Go 概念: ${item.goFunc}\n拼音: ${item.pinyin}`,
        range,
        filterText: `${item.pinyin} ${item.chinese}`,
      }))

      return { suggestions: suggestions as monaco.languages.CompletionItem[] }
    },
  })
  providerDisposables.push(disposable)
}

function setupGoPackageCompletion() {
  const goPackages: { label: string; detail: string }[] = [
    { label: 'fmt', detail: '格式化输入输出' },
    { label: 'os', detail: '操作系统接口' },
    { label: 'io', detail: '基础 I/O 接口' },
    { label: 'bufio', detail: '带缓冲的 I/O' },
    { label: 'strings', detail: '字符串操作' },
    { label: 'strconv', detail: '字符串与其他类型转换' },
    { label: 'bytes', detail: '字节切片操作' },
    { label: 'math', detail: '数学函数' },
    { label: 'math/rand', detail: '随机数' },
    { label: 'time', detail: '时间与日期' },
    { label: 'sort', detail: '排序' },
    { label: 'sync', detail: '同步原语 (Mutex, WaitGroup)' },
    { label: 'sync/atomic', detail: '原子操作' },
    { label: 'context', detail: '上下文传播' },
    { label: 'encoding/json', detail: 'JSON 编解码' },
    { label: 'encoding/xml', detail: 'XML 编解码' },
    { label: 'encoding/csv', detail: 'CSV 编解码' },
    { label: 'encoding/base64', detail: 'Base64 编解码' },
    { label: 'encoding/hex', detail: '十六进制编解码' },
    { label: 'net', detail: '网络 I/O' },
    { label: 'net/http', detail: 'HTTP 客户端/服务端' },
    { label: 'net/url', detail: 'URL 解析' },
    { label: 'net/mail', detail: '邮件解析' },
    { label: 'crypto', detail: '加密' },
    { label: 'crypto/md5', detail: 'MD5 哈希' },
    { label: 'crypto/sha256', detail: 'SHA256 哈希' },
    { label: 'crypto/tls', detail: 'TLS 加密' },
    { label: 'database/sql', detail: 'SQL 数据库接口' },
    { label: 'html/template', detail: 'HTML 模板' },
    { label: 'text/template', detail: '文本模板' },
    { label: 'log', detail: '日志记录' },
    { label: 'log/slog', detail: '结构化日志' },
    { label: 'path', detail: '路径操作' },
    { label: 'path/filepath', detail: '文件路径操作' },
    { label: 'reflect', detail: '运行时反射' },
    { label: 'regexp', detail: '正则表达式' },
    { label: 'errors', detail: '错误处理' },
    { label: 'flag', detail: '命令行参数' },
    { label: 'json', detail: 'JSON 编解码' },
  ]

  const disposable = monaco.languages.registerCompletionItemProvider('go', {
    triggerCharacters: ['"'],
    provideCompletionItems: (model, position) => {
      const lineContent = model.getLineContent(position.lineNumber)
      const textBefore = lineContent.substring(0, position.column - 1)

      if (!textBefore.includes('import "') && !textBefore.includes('import (') && !textBefore.trimEnd().endsWith('"')) {
        const importBlockMatch = textBefore.includes('import (')
        if (!importBlockMatch) return { suggestions: [] }
      }

      const word = model.getWordUntilPosition(position)
      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn,
      }

      const suggestions = goPackages.map(pkg => ({
        label: pkg.label,
        insertText: pkg.label,
        kind: monaco.languages.CompletionItemKind.Module,
        detail: pkg.detail,
        documentation: `Go 标准库: ${pkg.label}\n${pkg.detail}`,
        range,
      }))

      return { suggestions: suggestions as monaco.languages.CompletionItem[] }
    },
  })
  providerDisposables.push(disposable)
}

function setupGlobalSymbolCompletion() {
  const disposable = monaco.languages.registerCompletionItemProvider('go', {
    provideCompletionItems: (model, _position) => {
      const suggestions: monaco.languages.CompletionItem[] = []
      const text = model.getValue()
      const seenLabels = new Set<string>()

      const varRegex = /(?:var|:=)\s+(\w+)/g
      let match: RegExpExecArray | null
      while ((match = varRegex.exec(text)) !== null) {
        const name = match[1]
        if (!seenLabels.has(name)) {
          seenLabels.add(name)
          suggestions.push({
            label: name,
            insertText: name,
            kind: monaco.languages.CompletionItemKind.Variable,
            detail: '变量',
            documentation: '当前文件中定义的变量',
          } as monaco.languages.CompletionItem)
        }
      }

      const funcRegex = /func\s+(?:\([^)]+\)\s+)?(\w+)/g
      while ((match = funcRegex.exec(text)) !== null) {
        const name = match[1]
        if (!seenLabels.has(name)) {
          seenLabels.add(name)
          suggestions.push({
            label: name,
            insertText: `${name}($1)`,
            insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            kind: monaco.languages.CompletionItemKind.Function,
            detail: '函数',
            documentation: '当前文件中定义的函数',
          } as monaco.languages.CompletionItem)
        }
      }

      const typeRegex = /type\s+(\w+)\s+struct/g
      while ((match = typeRegex.exec(text)) !== null) {
        const name = match[1]
        if (!seenLabels.has(name)) {
          seenLabels.add(name)
          suggestions.push({
            label: name,
            insertText: name,
            kind: monaco.languages.CompletionItemKind.Struct,
            detail: '结构体',
            documentation: '当前文件中定义的结构体',
          } as monaco.languages.CompletionItem)
        }
      }

      return { suggestions }
    },
  })
  providerDisposables.push(disposable)
}

function setupNxEGOFunctionCompletion() {
  const nxegoFunctions: { name: string; goFunc: string; description: string; category: string; params: string[]; returnType: string }[] = [
    { name: '输出.打印', goFunc: 'yigou.Print', description: '输出内容到控制台', category: '输出', params: ['内容'], returnType: '' },
    { name: '输出.打印行', goFunc: 'yigou.PrintLine', description: '输出一行内容到控制台', category: '输出', params: ['内容'], returnType: '' },
    { name: '输出.格式化', goFunc: 'yigou.FormatString', description: '格式化字符串', category: '输出', params: ['格式', '参数...'], returnType: '文本型' },
    { name: '窗口.弹出提示', goFunc: 'yigou.MsgBox', description: '弹出提示信息', category: '窗口', params: ['提示内容'], returnType: '' },
    { name: '窗口.弹出询问', goFunc: 'yigou.ConfirmBox', description: '弹出询问对话框，返回是/否', category: '窗口', params: ['询问内容'], returnType: '逻辑型' },
    { name: '窗口.弹出输入', goFunc: 'yigou.InputBox', description: '弹出输入对话框', category: '窗口', params: ['提示内容', '默认值'], returnType: '文本型' },
    { name: '窗口.关闭', goFunc: 'yigou.Exit', description: '关闭程序', category: '窗口', params: ['退出码'], returnType: '' },
    { name: '文件.读取文本', goFunc: 'yigou.ReadFile', description: '读取文本文件内容', category: '文件', params: ['文件路径'], returnType: '文本型' },
    { name: '文件.写入文本', goFunc: 'yigou.WriteFile', description: '写入文本内容到文件', category: '文件', params: ['文件路径', '文本内容'], returnType: '逻辑型' },
    { name: '文件.是否存在', goFunc: 'yigou.FileExists', description: '判断文件是否存在', category: '文件', params: ['文件路径'], returnType: '逻辑型' },
    { name: '文件.删除', goFunc: 'yigou.DeleteFile', description: '删除文件', category: '文件', params: ['文件路径'], returnType: '逻辑型' },
    { name: '文件.重命名', goFunc: 'yigou.RenameFile', description: '重命名文件', category: '文件', params: ['原路径', '新路径'], returnType: '逻辑型' },
    { name: '文件.复制', goFunc: 'yigou.CopyFile', description: '复制文件', category: '文件', params: ['源路径', '目标路径'], returnType: '逻辑型' },
    { name: '文件.获取大小', goFunc: 'yigou.FileSize', description: '获取文件大小（字节）', category: '文件', params: ['文件路径'], returnType: '整数型' },
    { name: '文件.获取修改时间', goFunc: 'yigou.FileModTime', description: '获取文件最后修改时间', category: '文件', params: ['文件路径'], returnType: '文本型' },
    { name: '文件.创建目录', goFunc: 'yigou.MkdirAll', description: '创建多级目录', category: '文件', params: ['目录路径'], returnType: '逻辑型' },
    { name: '文件.删除目录', goFunc: 'yigou.RemoveAll', description: '删除目录及所有内容', category: '文件', params: ['目录路径'], returnType: '逻辑型' },
    { name: '文件.列出文件', goFunc: 'yigou.ReadDir', description: '列出目录中的文件', category: '文件', params: ['目录路径'], returnType: '文本型数组' },
    { name: '字符串.截取', goFunc: 'yigou.Substring', description: '截取字符串的一部分', category: '字符串', params: ['原字符串', '起始位置', '截取长度'], returnType: '文本型' },
    { name: '字符串.长度', goFunc: 'yigou.StrLen', description: '获取字符串长度', category: '字符串', params: ['字符串'], returnType: '整数型' },
    { name: '字符串.查找', goFunc: 'yigou.StrIndex', description: '查找子串位置', category: '字符串', params: ['原字符串', '查找内容'], returnType: '整数型' },
    { name: '字符串.替换', goFunc: 'yigou.StrReplaceAll', description: '替换字符串中的内容', category: '字符串', params: ['原字符串', '旧内容', '新内容'], returnType: '文本型' },
    { name: '字符串.分割', goFunc: 'yigou.StrSplit', description: '按分隔符分割字符串', category: '字符串', params: ['原字符串', '分隔符'], returnType: '文本型数组' },
    { name: '字符串.合并', goFunc: 'yigou.StrJoin', description: '合并字符串数组', category: '字符串', params: ['字符串数组', '分隔符'], returnType: '文本型' },
    { name: '字符串.转大写', goFunc: 'yigou.StrToUpper', description: '转换为大写', category: '字符串', params: ['字符串'], returnType: '文本型' },
    { name: '字符串.转小写', goFunc: 'yigou.StrToLower', description: '转换为小写', category: '字符串', params: ['字符串'], returnType: '文本型' },
    { name: '字符串.去空格', goFunc: 'yigou.StrTrimSpace', description: '去除首尾空格', category: '字符串', params: ['字符串'], returnType: '文本型' },
    { name: '字符串.是否包含', goFunc: 'yigou.StrContains', description: '判断是否包含子串', category: '字符串', params: ['原字符串', '子串'], returnType: '逻辑型' },
    { name: '数值.转文本', goFunc: 'yigou.Itoa', description: '将整数转换为文本', category: '数值', params: ['整数'], returnType: '文本型' },
    { name: '数值.转整数', goFunc: 'yigou.Atoi', description: '将文本转换为整数', category: '数值', params: ['文本'], returnType: '整数型' },
    { name: '数值.转小数', goFunc: 'yigou.ParseFloat', description: '将文本转换为小数', category: '数值', params: ['文本'], returnType: '小数型' },
    { name: '数值.取整', goFunc: 'yigou.Floor', description: '向下取整', category: '数值', params: ['数值'], returnType: '小数型' },
    { name: '数值.四舍五入', goFunc: 'yigou.Round', description: '四舍五入取整', category: '数值', params: ['数值'], returnType: '小数型' },
    { name: '数值.取绝对值', goFunc: 'yigou.Abs', description: '取绝对值', category: '数值', params: ['数值'], returnType: '小数型' },
    { name: '数值.取最大值', goFunc: 'yigou.Max', description: '取两个数中的最大值', category: '数值', params: ['数值1', '数值2'], returnType: '小数型' },
    { name: '数值.取最小值', goFunc: 'yigou.Min', description: '取两个数中的最小值', category: '数值', params: ['数值1', '数值2'], returnType: '小数型' },
    { name: '数值.取随机数', goFunc: 'yigou.RandomInt', description: '取0到N-1的随机整数', category: '数值', params: ['最大值'], returnType: '整数型' },
    { name: '数值.求平方根', goFunc: 'yigou.Sqrt', description: '求平方根', category: '数值', params: ['数值'], returnType: '小数型' },
    { name: '数值.求幂', goFunc: 'yigou.Pow', description: '求x的y次方', category: '数值', params: ['底数', '指数'], returnType: '小数型' },
    { name: '延时', goFunc: 'yigou.Sleep', description: '程序延时指定毫秒数', category: '系统', params: ['毫秒数'], returnType: '' },
    { name: '系统.获取当前时间', goFunc: 'yigou.NowTime', description: '获取当前系统时间', category: '系统', params: [], returnType: '文本型' },
    { name: '系统.执行命令', goFunc: 'yigou.RunCommand', description: '执行系统命令', category: '系统', params: ['命令', '参数...'], returnType: '文本型' },
    { name: '系统.获取环境变量', goFunc: 'yigou.GetEnv', description: '获取环境变量值', category: '系统', params: ['变量名'], returnType: '文本型' },
    { name: '系统.设置环境变量', goFunc: 'yigou.SetEnv', description: '设置环境变量', category: '系统', params: ['变量名', '值'], returnType: '' },
    { name: '系统.获取命令行参数', goFunc: 'yigou.GetArgs', description: '获取命令行参数', category: '系统', params: [], returnType: '文本型数组' },
    { name: '系统.获取当前目录', goFunc: 'yigou.Getwd', description: '获取当前工作目录', category: '系统', params: [], returnType: '文本型' },
    { name: '系统.设置当前目录', goFunc: 'yigou.Chdir', description: '设置当前工作目录', category: '系统', params: ['目录路径'], returnType: '逻辑型' },
    { name: '系统.打开文件夹', goFunc: 'yigou.OpenDir', description: '打开指定文件夹', category: '系统', params: ['文件夹路径'], returnType: '' },
    { name: '网络.发送请求', goFunc: 'yigou.HTTPGet', description: '发送HTTP GET请求', category: '网络', params: ['网址'], returnType: '文本型' },
    { name: '网络.发送POST', goFunc: 'yigou.HTTPPost', description: '发送HTTP POST请求', category: '网络', params: ['网址', '数据类型', '数据'], returnType: '文本型' },
    { name: '网络.下载文件', goFunc: 'yigou.DownloadFile', description: '下载文件到本地', category: '网络', params: ['网址', '保存路径'], returnType: '逻辑型' },
    { name: '网络.编码URL', goFunc: 'yigou.URLEncode', description: 'URL编码', category: '网络', params: ['文本'], returnType: '文本型' },
    { name: '网络.解码URL', goFunc: 'yigou.URLDecode', description: 'URL解码', category: '网络', params: ['文本'], returnType: '文本型' },
    { name: 'JSON.编码', goFunc: 'yigou.JSONMarshal', description: '将数据编码为JSON文本', category: 'JSON', params: ['数据'], returnType: '文本型' },
    { name: 'JSON.解码', goFunc: 'yigou.JSONUnmarshal', description: '将JSON文本解码为数据', category: 'JSON', params: ['JSON文本', '目标变量'], returnType: '' },
    { name: '加密.MD5', goFunc: 'yigou.MD5Hash', description: '计算MD5哈希值', category: '加密', params: ['文本'], returnType: '文本型' },
    { name: '加密.SHA256', goFunc: 'yigou.SHA256Hash', description: '计算SHA256哈希值', category: '加密', params: ['文本'], returnType: '文本型' },
    { name: '加密.Base64编码', goFunc: 'yigou.Base64Encode', description: 'Base64编码', category: '加密', params: ['数据'], returnType: '文本型' },
    { name: '加密.Base64解码', goFunc: 'yigou.Base64Decode', description: 'Base64解码', category: '加密', params: ['文本'], returnType: '文本型' },
    { name: '调试输出', goFunc: 'yigou.DebugPrint', description: '输出调试信息到控制台', category: '调试', params: ['输出内容'], returnType: '' },
    { name: '读入文件', goFunc: 'yigou.ReadBinaryFile', description: '读取文件字节内容', category: '文件', params: ['文件路径'], returnType: '字节集' },
    { name: '写到文件', goFunc: 'yigou.WriteBinaryFile', description: '写入字节到文件', category: '文件', params: ['文件路径', '数据'], returnType: '逻辑型' },
    { name: '信息框', goFunc: 'yigou.MessageBox', description: '弹出信息提示框', category: '界面', params: ['提示信息', '标题', '按钮类型'], returnType: '整数型' },
    { name: '取文本长度', goFunc: 'yigou.Len', description: '取文本长度', category: '文本', params: ['文本'], returnType: '整数型' },
    { name: '到文本', goFunc: 'yigou.Str', description: '转换到文本', category: '转换', params: ['数据'], returnType: '文本型' },
    { name: '到整数', goFunc: 'yigou.Int', description: '转换到整数', category: '转换', params: ['文本'], returnType: '整数型' },
    { name: '取启动时间', goFunc: 'yigou.GetTickCount', description: '取系统启动时间', category: '系统', params: [], returnType: '整数型' },
    { name: '取现行时间', goFunc: 'yigou.Now', description: '获取当前时间', category: '时间', params: [], returnType: '日期时间型' },
    { name: '运行', goFunc: 'yigou.Exec', description: '执行外部命令', category: '系统', params: ['命令', '是否等待'], returnType: '整数型' },
    { name: '线程.启动', goFunc: 'yigou.StartThread', description: '启动子线程', category: '线程', params: ['子程序地址', '参数'], returnType: '线程句柄' },
    { name: '线程.等待', goFunc: 'yigou.WaitThread', description: '等待线程结束', category: '线程', params: ['线程句柄'], returnType: '' },
    { name: '时间.取现行时间戳', goFunc: 'yigou.Timestamp', description: '获取当前Unix时间戳', category: '时间', params: [], returnType: '整数型' },
    { name: '时间.格式化', goFunc: 'yigou.FormatTime', description: '格式化时间为文本', category: '时间', params: ['时间戳', '格式文本'], returnType: '文本型' },
    { name: '时间.时间戳转文本', goFunc: 'yigou.TimestampToStr', description: '时间戳转日期文本', category: '时间', params: ['时间戳'], returnType: '文本型' },
    { name: '数组.取成员数', goFunc: 'yigou.ArrayLen', description: '取数组的成员数量', category: '数组', params: ['数组'], returnType: '整数型' },
    { name: '数组.加入成员', goFunc: 'yigou.ArrayAppend', description: '向数组末尾添加成员', category: '数组', params: ['数组', '成员'], returnType: '' },
    { name: '数组.删除成员', goFunc: 'yigou.ArrayRemove', description: '删除数组指定位置成员', category: '数组', params: ['数组', '位置'], returnType: '' },
    { name: '数组.排序', goFunc: 'yigou.ArraySort', description: '对数组进行排序', category: '数组', params: ['数组', '是否升序'], returnType: '' },
    { name: '正则.是否匹配', goFunc: 'yigou.RegexMatch', description: '正则表达式是否匹配', category: '正则', params: ['文本', '正则表达式'], returnType: '逻辑型' },
    { name: '正则.查找', goFunc: 'yigou.RegexFind', description: '正则查找第一个匹配', category: '正则', params: ['文本', '正则表达式'], returnType: '文本型' },
    { name: '正则.替换', goFunc: 'yigou.RegexReplace', description: '正则替换全部匹配', category: '正则', params: ['文本', '正则表达式', '替换为'], returnType: '文本型' },
    { name: '数据.MD5校验', goFunc: 'yigou.MD5Check', description: '校验数据MD5', category: '加密', params: ['数据', '期望值'], returnType: '逻辑型' },
    { name: '数据.CRC32', goFunc: 'yigou.CRC32', description: '计算CRC32校验和', category: '加密', params: ['数据'], returnType: '整数型' },
    { name: '系统.取CPU核心数', goFunc: 'yigou.GetCPUCount', description: '获取CPU核心数量', category: '系统', params: [], returnType: '整数型' },
    { name: '系统.取内存信息', goFunc: 'yigou.GetMemoryInfo', description: '获取内存使用信息', category: '系统', params: [], returnType: '文本型' },
    { name: '系统.取磁盘空间', goFunc: 'yigou.GetDiskSpace', description: '获取磁盘剩余空间', category: '系统', params: ['盘符'], returnType: '整数型' },
    { name: '组件.取属性', goFunc: 'yigou.GetComponentProp', description: '获取组件属性值', category: '界面', params: ['组件名', '属性名'], returnType: '文本型' },
    { name: '组件.置属性', goFunc: 'yigou.SetComponentProp', description: '设置组件属性值', category: '界面', params: ['组件名', '属性名', '属性值'], returnType: '逻辑型' },
    { name: '程序.取参数', goFunc: 'yigou.GetCmdArg', description: '获取命令行参数', category: '系统', params: ['索引'], returnType: '文本型' },
    { name: '程序.退出', goFunc: 'yigou.ExitApp', description: '退出程序', category: '系统', params: [], returnType: '' },
    { name: '画板.画直线', goFunc: 'yigou.CanvasLine', description: '在画板上画直线', category: '绘图', params: ['起点X', '起点Y', '终点X', '终点Y'], returnType: '' },
    { name: '画板.画矩形', goFunc: 'yigou.CanvasRect', description: '在画板上画矩形', category: '绘图', params: ['X', 'Y', '宽', '高'], returnType: '' },
    { name: '画板.画文本', goFunc: 'yigou.CanvasText', description: '在画板上画文本', category: '绘图', params: ['X', 'Y', '文本'], returnType: '' },
    { name: '画板.清空', goFunc: 'yigou.CanvasClear', description: '清空画板', category: '绘图', params: [], returnType: '' },
    { name: '按钮.禁用', goFunc: 'yigou.ButtonDisable', description: '禁用按钮', category: '界面', params: ['按钮名'], returnType: '' },
    { name: '按钮.启用', goFunc: 'yigou.ButtonEnable', description: '启用按钮', category: '界面', params: ['按钮名'], returnType: '' },
    { name: '编辑框.取内容', goFunc: 'yigou.EditGetText', description: '获取编辑框内容', category: '界面', params: ['编辑框名'], returnType: '文本型' },
    { name: '编辑框.置内容', goFunc: 'yigou.EditSetText', description: '设置编辑框内容', category: '界面', params: ['编辑框名', '内容'], returnType: '' },
  ]

  const disposable = monaco.languages.registerCompletionItemProvider('go', {
    triggerCharacters: [...'abcdefghijklmnopqrstuvwxyz'.split(''), '输', '窗', '文', '字', '数', '系', '网', '延', '加', '调', '信', '读', '写', '取', '到', '线', '时', '数', '正', '组', '组', '按', '编', '画', '程', '设', '创', '删', '复', '移', '重', '命', '列'],
    provideCompletionItems: (model, position) => {
      const word = model.getWordUntilPosition(position)
      const prefix = word.word
      if (prefix.length < 1) return { suggestions: [] }

      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn,
      }

      const matches = nxegoFunctions.filter(fn =>
        fn.name.toLowerCase().includes(prefix.toLowerCase()) ||
        fn.name.startsWith(prefix) ||
        fn.goFunc.toLowerCase().includes(prefix.toLowerCase())
      )

      if (matches.length === 0) return { suggestions: [] }

      const suggestions = matches.map(fn => ({
        label: fn.name,
        insertText: fn.name + '(',
        kind: monaco.languages.CompletionItemKind.Function,
        detail: `${fn.goFunc} | ${fn.category}`,
        documentation: `**${fn.name}**\n\n${fn.description}\n\n- 对应Go函数: \`${fn.goFunc}\`\n- 参数: ${fn.params.length > 0 ? fn.params.join(', ') : '无'}\n- 返回值: ${fn.returnType || '无'}\n- 分类: ${fn.category}`,
        range,
      }))

      return { suggestions: suggestions as monaco.languages.CompletionItem[] }
    },
  })
  providerDisposables.push(disposable)
}

function setupNxEGOHoverProvider() {
  const nxegoHoverMap = new Map<string, string>()
  const functions = [
    { name: '输出.打印', goFunc: 'yigou.Print', description: '输出内容到控制台', params: ['内容'], returnType: '' },
    { name: '输出.打印行', goFunc: 'yigou.PrintLine', description: '输出一行内容到控制台', params: ['内容'], returnType: '' },
    { name: '输出.格式化', goFunc: 'yigou.FormatString', description: '格式化字符串', params: ['格式', '参数...'], returnType: '文本型' },
    { name: '窗口.弹出提示', goFunc: 'yigou.MsgBox', description: '弹出提示信息', params: ['提示内容'], returnType: '' },
    { name: '窗口.弹出询问', goFunc: 'yigou.ConfirmBox', description: '弹出询问对话框，返回是/否', params: ['询问内容'], returnType: '逻辑型' },
    { name: '窗口.弹出输入', goFunc: 'yigou.InputBox', description: '弹出输入对话框', params: ['提示内容', '默认值'], returnType: '文本型' },
    { name: '窗口.关闭', goFunc: 'yigou.Exit', description: '关闭程序', params: ['退出码'], returnType: '' },
    { name: '文件.读取文本', goFunc: 'yigou.ReadFile', description: '读取文本文件内容', params: ['文件路径'], returnType: '文本型' },
    { name: '文件.写入文本', goFunc: 'yigou.WriteFile', description: '写入文本内容到文件', params: ['文件路径', '文本内容'], returnType: '逻辑型' },
    { name: '文件.是否存在', goFunc: 'yigou.FileExists', description: '判断文件是否存在', params: ['文件路径'], returnType: '逻辑型' },
    { name: '文件.删除', goFunc: 'yigou.DeleteFile', description: '删除文件', params: ['文件路径'], returnType: '逻辑型' },
    { name: '文件.重命名', goFunc: 'yigou.RenameFile', description: '重命名文件', params: ['原路径', '新路径'], returnType: '逻辑型' },
    { name: '文件.复制', goFunc: 'yigou.CopyFile', description: '复制文件', params: ['源路径', '目标路径'], returnType: '逻辑型' },
    { name: '文件.获取大小', goFunc: 'yigou.FileSize', description: '获取文件大小（字节）', params: ['文件路径'], returnType: '整数型' },
    { name: '文件.获取修改时间', goFunc: 'yigou.FileModTime', description: '获取文件最后修改时间', params: ['文件路径'], returnType: '文本型' },
    { name: '文件.创建目录', goFunc: 'yigou.MkdirAll', description: '创建多级目录', params: ['目录路径'], returnType: '逻辑型' },
    { name: '文件.删除目录', goFunc: 'yigou.RemoveAll', description: '删除目录及所有内容', params: ['目录路径'], returnType: '逻辑型' },
    { name: '文件.列出文件', goFunc: 'yigou.ReadDir', description: '列出目录中的文件', params: ['目录路径'], returnType: '文本型数组' },
    { name: '字符串.截取', goFunc: 'yigou.Substring', description: '截取字符串的一部分', params: ['原字符串', '起始位置', '截取长度'], returnType: '文本型' },
    { name: '字符串.长度', goFunc: 'yigou.StrLen', description: '获取字符串长度', params: ['字符串'], returnType: '整数型' },
    { name: '字符串.查找', goFunc: 'yigou.StrIndex', description: '查找子串位置', params: ['原字符串', '查找内容'], returnType: '整数型' },
    { name: '字符串.替换', goFunc: 'yigou.StrReplaceAll', description: '替换字符串中的内容', params: ['原字符串', '旧内容', '新内容'], returnType: '文本型' },
    { name: '字符串.分割', goFunc: 'yigou.StrSplit', description: '按分隔符分割字符串', params: ['原字符串', '分隔符'], returnType: '文本型数组' },
    { name: '字符串.合并', goFunc: 'yigou.StrJoin', description: '合并字符串数组', params: ['字符串数组', '分隔符'], returnType: '文本型' },
    { name: '字符串.转大写', goFunc: 'yigou.StrToUpper', description: '转换为大写', params: ['字符串'], returnType: '文本型' },
    { name: '字符串.转小写', goFunc: 'yigou.StrToLower', description: '转换为小写', params: ['字符串'], returnType: '文本型' },
    { name: '字符串.去空格', goFunc: 'yigou.StrTrimSpace', description: '去除首尾空格', params: ['字符串'], returnType: '文本型' },
    { name: '字符串.是否包含', goFunc: 'yigou.StrContains', description: '判断是否包含子串', params: ['原字符串', '子串'], returnType: '逻辑型' },
    { name: '数值.转文本', goFunc: 'yigou.Itoa', description: '将整数转换为文本', params: ['整数'], returnType: '文本型' },
    { name: '数值.转整数', goFunc: 'yigou.Atoi', description: '将文本转换为整数', params: ['文本'], returnType: '整数型' },
    { name: '数值.转小数', goFunc: 'yigou.ParseFloat', description: '将文本转换为小数', params: ['文本'], returnType: '小数型' },
    { name: '数值.取整', goFunc: 'yigou.Floor', description: '向下取整', params: ['数值'], returnType: '小数型' },
    { name: '数值.四舍五入', goFunc: 'yigou.Round', description: '四舍五入取整', params: ['数值'], returnType: '小数型' },
    { name: '数值.取绝对值', goFunc: 'yigou.Abs', description: '取绝对值', params: ['数值'], returnType: '小数型' },
    { name: '数值.求平方根', goFunc: 'yigou.Sqrt', description: '求平方根', params: ['数值'], returnType: '小数型' },
    { name: '数值.求幂', goFunc: 'yigou.Pow', description: '求x的y次方', params: ['底数', '指数'], returnType: '小数型' },
    { name: '数值.取最大值', goFunc: 'yigou.Max', description: '取两个数中的最大值', params: ['数值1', '数值2'], returnType: '小数型' },
    { name: '数值.取最小值', goFunc: 'yigou.Min', description: '取两个数中的最小值', params: ['数值1', '数值2'], returnType: '小数型' },
    { name: '数值.取随机数', goFunc: 'yigou.RandomInt', description: '取0到N-1的随机整数', params: ['最大值'], returnType: '整数型' },
    { name: '延时', goFunc: 'yigou.Sleep', description: '程序延时指定毫秒数', params: ['毫秒数'], returnType: '' },
    { name: '系统.获取当前时间', goFunc: 'yigou.NowTime', description: '获取当前系统时间', params: [], returnType: '文本型' },
    { name: '系统.执行命令', goFunc: 'yigou.RunCommand', description: '执行系统命令', params: ['命令', '参数...'], returnType: '文本型' },
    { name: '系统.获取环境变量', goFunc: 'yigou.GetEnv', description: '获取环境变量值', params: ['变量名'], returnType: '文本型' },
    { name: '系统.获取命令行参数', goFunc: 'yigou.GetArgs', description: '获取命令行参数', params: [], returnType: '文本型数组' },
    { name: '系统.获取当前目录', goFunc: 'yigou.Getwd', description: '获取当前工作目录', params: [], returnType: '文本型' },
    { name: '网络.发送请求', goFunc: 'yigou.HTTPGet', description: '发送HTTP GET请求', params: ['网址'], returnType: '文本型' },
    { name: '网络.下载文件', goFunc: 'yigou.DownloadFile', description: '下载文件到本地', params: ['网址', '保存路径'], returnType: '逻辑型' },
    { name: 'JSON.编码', goFunc: 'yigou.JSONMarshal', description: '将数据编码为JSON文本', params: ['数据'], returnType: '文本型' },
    { name: 'JSON.解码', goFunc: 'yigou.JSONUnmarshal', description: '将JSON文本解码为数据', params: ['JSON文本', '目标变量'], returnType: '' },
    { name: '加密.MD5', goFunc: 'yigou.MD5Hash', description: '计算MD5哈希值', params: ['文本'], returnType: '文本型' },
    { name: '加密.Base64编码', goFunc: 'yigou.Base64Encode', description: 'Base64编码', params: ['数据'], returnType: '文本型' },
    { name: '加密.Base64解码', goFunc: 'yigou.Base64Decode', description: 'Base64解码', params: ['文本'], returnType: '文本型' },
    { name: '调试输出', goFunc: 'yigou.DebugPrint', description: '输出调试信息到控制台', params: ['输出内容'], returnType: '' },
    { name: '读入文件', goFunc: 'yigou.ReadBinaryFile', description: '读取文件字节内容', params: ['文件路径'], returnType: '字节集' },
    { name: '写到文件', goFunc: 'yigou.WriteBinaryFile', description: '写入字节到文件', params: ['文件路径', '数据'], returnType: '逻辑型' },
    { name: '信息框', goFunc: 'yigou.MessageBox', description: '弹出信息提示框', params: ['提示信息', '标题', '按钮类型'], returnType: '整数型' },
    { name: '取文本长度', goFunc: 'yigou.Len', description: '取文本长度', params: ['文本'], returnType: '整数型' },
    { name: '到文本', goFunc: 'yigou.Str', description: '转换到文本', params: ['数据'], returnType: '文本型' },
    { name: '到整数', goFunc: 'yigou.Int', description: '转换到整数', params: ['文本'], returnType: '整数型' },
    { name: '运行', goFunc: 'yigou.Exec', description: '执行外部命令', params: ['命令', '是否等待'], returnType: '整数型' },
    { name: '线程.启动', goFunc: 'yigou.StartThread', description: '启动子线程', params: ['子程序地址', '参数'], returnType: '线程句柄' },
    { name: '时间.格式化', goFunc: 'yigou.FormatTime', description: '格式化时间为文本', params: ['时间戳', '格式文本'], returnType: '文本型' },
    { name: '数组.取成员数', goFunc: 'yigou.ArrayLen', description: '取数组的成员数量', params: ['数组'], returnType: '整数型' },
    { name: '正则.是否匹配', goFunc: 'yigou.RegexMatch', description: '正则表达式是否匹配', params: ['文本', '正则表达式'], returnType: '逻辑型' },
    { name: '画板.画直线', goFunc: 'yigou.CanvasLine', description: '在画板上画直线', params: ['起点X', '起点Y', '终点X', '终点Y'], returnType: '' },
    { name: '编辑框.取内容', goFunc: 'yigou.EditGetText', description: '获取编辑框内容', params: ['编辑框名'], returnType: '文本型' },
    { name: '编辑框.置内容', goFunc: 'yigou.EditSetText', description: '设置编辑框内容', params: ['编辑框名', '内容'], returnType: '' },
  ]

  for (const fn of functions) {
    const paramsText = fn.params.length > 0 ? fn.params.join(', ') : '无'
    const returnText = fn.returnType || '无'
    nxegoHoverMap.set(fn.name, `**${fn.name}** - ${fn.description}\n\n参数: ${paramsText}\n返回值: ${returnText}\n对应Go: \`${fn.goFunc}\``)
  }

  const disposable = monaco.languages.registerHoverProvider('go', {
    provideHover: (model, position) => {
      const word = model.getWordAtPosition(position)
      if (!word) return null

      const wordText = word.word
      const lineContent = model.getLineContent(position.lineNumber)
      const lineBeforeCursor = lineContent.substring(0, position.column - 1)
      const lineAfterCursor = lineContent.substring(position.column - 1)

      for (const [fnName, doc] of nxegoHoverMap) {
        if (wordText === fnName || lineBeforeCursor.trimEnd().endsWith(fnName)) {
          return {
            contents: [{ value: doc }],
            range: {
              startLineNumber: position.lineNumber,
              endLineNumber: position.lineNumber,
              startColumn: word.startColumn,
              endColumn: word.endColumn,
            },
          }
        }
      }

      const yigouMatch = wordText.match(/^yigou\.(\w+)$/)
      if (yigouMatch) {
        const goFuncName = yigouMatch[1]
        for (const fn of functions) {
          const goFuncParts = fn.goFunc.split('.')
          if (goFuncParts.length > 1 && goFuncParts[1] === goFuncName) {
            const paramsText = fn.params.length > 0 ? fn.params.join(', ') : '无'
            const returnText = fn.returnType || '无'
            return {
              contents: [{ value: `**${fn.goFunc}** - ${fn.description}\n\n中文名: ${fn.name}\n参数: ${paramsText}\n返回值: ${returnText}` }],
              range: {
                startLineNumber: position.lineNumber,
                endLineNumber: position.lineNumber,
                startColumn: word.startColumn,
                endColumn: word.endColumn,
              },
            }
          }
        }
      }

      return null
    },
  })
  providerDisposables.push(disposable)
}

const defaultOptions: monaco.editor.IStandaloneEditorConstructionOptions = {
  value: props.modelValue,
  language: props.language || 'go',
  theme: props.theme || getDefaultTheme(),
  readOnly: props.readOnly || false,
  fontSize: props.fontSize || 14,
  fontFamily: "'Consolas', 'Courier New', 'Microsoft YaHei', monospace",
  lineNumbers: 'on',
  roundedSelection: true,
  scrollBeyondLastLine: false,
  automaticLayout: true,
  tabSize: 4,
  wordWrap: 'on',
  minimap: { enabled: props.minimap !== false },
  folding: true,
  foldingStrategy: 'indentation',
  showFoldingControls: 'always',
  bracketPairColorization: { enabled: true },
  guides: {
    indentation: true,
    bracketPairs: true
  },
  suggestOnTriggerCharacters: true,
  quickSuggestions: {
    other: true,
    comments: false,
    strings: false,
  },
  acceptSuggestionOnEnter: 'on',
  snippetSuggestions: 'top',
  wordBasedSuggestions: 'allDocuments' as any,
  parameterHints: { enabled: true, cycle: true },
  formatOnPaste: true,
  formatOnType: false,
  renderWhitespace: 'selection',
  cursorBlinking: 'smooth',
  cursorSmoothCaretAnimation: 'on',
  smoothScrolling: true,
  padding: { top: 12, bottom: 12 },
  matchBrackets: 'always',
  autoClosingBrackets: 'always',
  autoClosingQuotes: 'always',
  autoIndent: 'full',
  dragAndDrop: true,
  links: true,
  colorDecorators: true,
  lightbulb: { enabled: true },
  ...props.options
}

const initEditor = () => {
  if (!editorContainer.value) return

  editorInstance = monaco.editor.create(editorContainer.value, defaultOptions)

  editorInstance.onDidChangeModelContent(() => {
    const value = editorInstance?.getValue() || ''
    emit('update:modelValue', value)
    emit('change', value)
  })

  editorInstance.onDidChangeCursorPosition((e) => {
    emit('cursor-change', e.position)
  })

  editorInstance.addAction({
    id: 'save-file',
    label: '保存文件',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS],
    run: () => {
      ElMessage.success('文件已保存 (Ctrl+S)')
    }
  })

  editorInstance.addAction({
    id: 'goto-line-action',
    label: '转到行',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyG],
    run: () => showGoToLineDialog()
  })

  setupChineseContextMenu()
  setupELangCompletion()
  setupGoKeywordCompletion()
  setupPinyinCompletion()
  setupGoPackageCompletion()
  setupGlobalSymbolCompletion()
  setupNxEGOFunctionCompletion()
  setupNxEGOHoverProvider()
}

const getValue = (): string => {
  return editorInstance?.getValue() || ''
}

const setValue = (value: string) => {
  if (editorInstance) {
    const currentVal = editorInstance.getValue()
    if (currentVal !== value) {
      editorInstance.setValue(value)
    }
  }
}

const formatDocument = async () => {
  if (editorInstance) {
    await editorInstance.getAction('editor.action.formatDocument')?.run()
  }
}

const insertText = (text: string) => {
  if (editorInstance) {
    const position = editorInstance.getPosition()
    if (position) {
      editorInstance.executeEdits('insert-text', [{
        range: new monaco.Range(position.lineNumber, position.column, position.lineNumber, position.column),
        text
      }])
      editorInstance.focus()
    }
  }
}

const getSelectedText = (): string => {
  return editorInstance?.getModel()?.getValueInRange(editorInstance.getSelection()!) || ''
}

const replaceSelectedText = (text: string) => {
  if (editorInstance) {
    const selection = editorInstance.getSelection()
    if (selection) {
      editorInstance.executeEdits('replace-selection', [{
        range: selection,
        text
      }])
    }
  }
}

const focus = () => {
  editorInstance?.focus()
}

const dispose = () => {
  if (editorInstance) {
    editorInstance.dispose()
    editorInstance = null
  }
  providerDisposables.forEach(d => d.dispose())
  providerDisposables.length = 0
}

watch(() => props.modelValue, (newVal) => {
  if (editorInstance && editorInstance.getValue() !== newVal) {
    setValue(newVal)
  }
})

watch(() => props.language, (newLang) => {
  if (editorInstance) {
    const model = editorInstance.getModel()
    if (model) {
      monaco.editor.setModelLanguage(model, newLang || 'go')
    }
  }
})

watch(() => props.theme, (newTheme) => {
  if (editorInstance) {
    monaco.editor.setTheme(newTheme || getDefaultTheme())
  }
})

watch(currentTheme, () => {
  if (editorInstance && !props.theme) {
    monaco.editor.setTheme(getDefaultTheme())
  }
})

watch(() => props.readOnly, (readOnly) => {
  if (editorInstance) {
    editorInstance.updateOptions({ readOnly: readOnly || false })
  }
})

onMounted(() => {
  nextTick(() => {
    initEditor()
  })
})

onUnmounted(() => {
  dispose()
})

defineExpose({
  getInstance: () => editorInstance,
  getValue,
  setValue,
  formatDocument,
  insertText,
  getSelectedText,
  replaceSelectedText,
  focus,
  dispose,
  cursorPosition,
  lineCount,
})
</script>

<template>
  <div class="monaco-editor-wrapper">
    <div ref="editorContainer" class="monaco-editor-container"></div>
    <div class="editor-statusbar">
      <span class="status-item">行 {{ cursorPosition.lineNumber }}, 列 {{ cursorPosition.column }}</span>
      <span class="status-item">{{ lineCount }} 行</span>
      <span class="status-item">Go</span>
      <span class="status-item">UTF-8</span>
    </div>
  </div>
</template>

<style scoped>
.monaco-editor-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #1e1e1e;
}

.monaco-editor-container {
  flex: 1;
  width: 100%;
  height: 100%;
  min-height: 200px;
}

.editor-statusbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 2px 12px;
  background-color: #007acc;
  color: #ffffff;
  font-size: 12px;
  flex-shrink: 0;
  height: 22px;
  user-select: none;
}

.status-item {
  white-space: nowrap;
  opacity: 0.9;
}
</style>