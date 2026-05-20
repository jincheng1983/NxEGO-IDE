export interface TokenSpan {
  type: 'keyword' | 'string' | 'comment' | 'number' | 'type' | 'function' | 'operator' | 'plain' | 'elang' | 'package'
  text: string
  start: number
  end: number
}

export interface FlowSegment {
  kind: 'start' | 'line' | 'end' | 'branch'
  depth: number
}

export interface RenderBlock {
  kind: 'codeline' | 'blank' | 'header'
  lineIndex: number
  codeLine: string
  tokens: TokenSpan[]
  flowSegments: FlowSegment[]
  indentLevel: number
  isVirtual: boolean
}

export type ParseLanguage = 'go' | 'elang' | 'auto'

export interface ParseRequest {
  id: number
  text: string
  isClassModule: boolean
  language?: ParseLanguage
}

export interface ParseResponse {
  id: number
  blocks: RenderBlock[]
  flowMaxDepth: number
  lineCount: number
}

export const GO_KEYWORDS = new Set([
  'break', 'case', 'chan', 'const', 'continue', 'default', 'defer',
  'else', 'fallthrough', 'for', 'func', 'go', 'goto', 'if',
  'import', 'interface', 'map', 'package', 'range', 'return',
  'select', 'struct', 'switch', 'type', 'var',
])

export const GO_TYPES = new Set([
  'bool', 'byte', 'complex64', 'complex128', 'error', 'float32', 'float64',
  'int', 'int8', 'int16', 'int32', 'int64',
  'rune', 'string', 'uint', 'uint8', 'uint16', 'uint32', 'uint64', 'uintptr',
])

export const GO_CONSTANTS = new Set([
  'true', 'false', 'nil', 'iota',
])

export const FLOW_KEYWORDS = new Set([
  'if', 'else', 'for', 'switch', 'case', 'default', 'select',
  'func', 'return', 'break', 'continue', 'goto', 'fallthrough',
  'defer', 'go', 'range',
])

export const ELANG_KEYWORDS = new Set([
  '.子程序', '.参数', '.局部变量',
  '.如果', '.否则', '.否则如果', '.如果结束',
  '.判断', '.判断开始', '.默认', '.判断结束',
  '.循环', '.循环尾', '.循环判断首', '.循环判断尾',
  '.计次循环首', '.计次循环尾', '.变量循环首', '.变量循环尾',
  '.返回', '.公开', '.私有',
  '.参考', '.数组', '.静态',
  '.真', '.假', '.空',
  '.程序集', '.类', '.类方法',
  '.常量', '.全局变量', '.图片', '.声音', '.资源',
  '.版本', '.数据类型',
])

export const ELANG_FLOW_KEYWORDS = new Set([
  '.如果', '.否则', '.否则如果', '.如果结束',
  '.判断', '.判断开始', '.默认', '.判断结束',
  '.循环', '.循环尾', '.循环判断首', '.循环判断尾',
  '.计次循环首', '.计次循环尾', '.变量循环首', '.变量循环尾',
  '.子程序', '.返回',
])

export const ELANG_TYPES = new Set([
  '整数型', '小数型', '逻辑型', '文本型', '字节集',
  '日期时间型', '子程序指针', '通用型', '空',
  '短整数型', '长整数型', '双精度小数型',
])

export const ELANG_OPERATORS = new Set([
  '＝', '≠', '≥', '≤', '＜', '＞',
  '＋', '－', '×', '÷',
])