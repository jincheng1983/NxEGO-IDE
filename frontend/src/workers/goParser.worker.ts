import {
  type ParseRequest,
  type ParseResponse,
  type RenderBlock,
  type TokenSpan,
  type FlowSegment,
  type ParseLanguage,
  GO_KEYWORDS,
  GO_TYPES,
  GO_CONSTANTS,
  FLOW_KEYWORDS,
  ELANG_KEYWORDS,
  ELANG_FLOW_KEYWORDS,
  ELANG_TYPES,
  ELANG_OPERATORS,
} from './goParserShared'

function detectLanguage(text: string): 'go' | 'elang' {
  const sample = text.slice(0, 5000)
  let elangScore = 0
  let goScore = 0

  for (const kw of ELANG_KEYWORDS) {
    if (sample.includes(kw)) elangScore++
  }

  for (const kw of GO_KEYWORDS) {
    if (new RegExp(`\\b${kw}\\b`).test(sample)) goScore++
  }

  if (sample.includes('package ') && sample.includes('func ')) goScore += 5
  if (sample.includes('.子程序') || sample.includes('.如果')) elangScore += 5

  return elangScore > goScore ? 'elang' : 'go'
}

function tokenizeGoLine(line: string): TokenSpan[] {
  const tokens: TokenSpan[] = []
  let i = 0
  const len = line.length

  while (i < len) {
    if (line[i] === ' ' || line[i] === '\t' || line[i] === '\r') {
      i++
      continue
    }

    if (line[i] === '/' && i + 1 < len) {
      if (line[i + 1] === '/') {
        tokens.push({ type: 'comment', text: line.slice(i), start: i, end: len })
        break
      }
      if (line[i + 1] === '*') {
        const end = line.indexOf('*/', i + 2)
        if (end !== -1) {
          tokens.push({ type: 'comment', text: line.slice(i, end + 2), start: i, end: end + 2 })
          i = end + 2
          continue
        }
        tokens.push({ type: 'comment', text: line.slice(i), start: i, end: len })
        break
      }
    }

    if (line[i] === '"' || line[i] === '`') {
      const quote = line[i]
      let j = i + 1
      while (j < len) {
        if (line[j] === '\\' && quote === '"') { j += 2; continue }
        if (line[j] === quote) { j++; break }
        j++
      }
      tokens.push({ type: 'string', text: line.slice(i, j), start: i, end: j })
      i = j
      continue
    }

    if (line[i] === '\'') {
      let j = i + 1
      while (j < len && line[j] !== '\'') {
        if (line[j] === '\\') j++
        j++
      }
      if (j < len) j++
      tokens.push({ type: 'string', text: line.slice(i, j), start: i, end: j })
      i = j
      continue
    }

    if (/[0-9]/.test(line[i]) || (line[i] === '.' && i + 1 < len && /[0-9]/.test(line[i + 1]))) {
      let j = i
      if (line[j] === '0' && j + 1 < len && (line[j + 1] === 'x' || line[j + 1] === 'X')) {
        j += 2
        while (j < len && /[0-9a-fA-F]/.test(line[j])) j++
      } else if (line[j] === '0' && j + 1 < len && (line[j + 1] === 'b' || line[j + 1] === 'B')) {
        j += 2
        while (j < len && /[01]/.test(line[j])) j++
      } else {
        while (j < len && /[0-9]/.test(line[j])) j++
        if (j < len && line[j] === '.') { j++; while (j < len && /[0-9]/.test(line[j])) j++ }
        if (j < len && (line[j] === 'e' || line[j] === 'E')) {
          j++
          if (j < len && (line[j] === '+' || line[j] === '-')) j++
          while (j < len && /[0-9]/.test(line[j])) j++
        }
      }
      tokens.push({ type: 'number', text: line.slice(i, j), start: i, end: j })
      i = j
      continue
    }

    if (/[a-zA-Z_\u4e00-\u9fff]/.test(line[i])) {
      let j = i
      while (j < len && /[a-zA-Z0-9_\u4e00-\u9fff]/.test(line[j])) j++
      const word = line.slice(i, j)

      if (GO_KEYWORDS.has(word)) {
        tokens.push({ type: 'keyword', text: word, start: i, end: j })
      } else if (GO_TYPES.has(word)) {
        tokens.push({ type: 'type', text: word, start: i, end: j })
      } else if (GO_CONSTANTS.has(word)) {
        tokens.push({ type: 'keyword', text: word, start: i, end: j })
      } else if (/^[\u4e00-\u9fff]/.test(word)) {
        tokens.push({ type: 'elang', text: word, start: i, end: j })
      } else {
        const nextNonSpace = line.slice(j).trimStart()
        if (nextNonSpace.startsWith('(')) {
          tokens.push({ type: 'function', text: word, start: i, end: j })
        } else {
          tokens.push({ type: 'plain', text: word, start: i, end: j })
        }
      }
      i = j
      continue
    }

    const opChars = '+-*/%=<>!&|^~.:;,(){}[]'
    if (opChars.includes(line[i])) {
      let j = i + 1
      const twoCharOps = ['==', '!=', '<=', '>=', ':=', '+=', '-=', '*=', '/=', '%=', '&=', '|=', '^=', '<<', '>>', '++', '--', '&&', '||', '<-', '->', '::']
      if (j < len) {
        const two = line[i] + line[j]
        if (twoCharOps.includes(two)) j++
      }
      tokens.push({ type: 'operator', text: line.slice(i, j), start: i, end: j })
      i = j
      continue
    }

    tokens.push({ type: 'plain', text: line[i], start: i, end: i + 1 })
    i++
  }

  return tokens
}

function tokenizeELangLine(line: string): TokenSpan[] {
  const tokens: TokenSpan[] = []
  let i = 0
  const len = line.length

  while (i < len) {
    if (line[i] === ' ' || line[i] === '\t' || line[i] === '\r') {
      i++
      continue
    }

    if (line[i] === '\'' || line[i] === '‘') {
      tokens.push({ type: 'comment', text: line.slice(i), start: i, end: len })
      break
    }

    if (line[i] === '"' || line[i] === '\u201C') {
      const startQuote = line[i]
      const endQuote = startQuote === '"' ? '"' : '\u201D'
      let j = i + 1
      while (j < len && line[j] !== endQuote) j++
      if (j < len) j++
      tokens.push({ type: 'string', text: line.slice(i, j), start: i, end: j })
      i = j
      continue
    }

    if (line[i] === '#' && i + 1 < len && line[i + 1] === '"') {
      let j = i + 2
      while (j < len && line[j] !== '"') j++
      if (j < len) j++
      tokens.push({ type: 'string', text: line.slice(i, j), start: i, end: j })
      i = j
      continue
    }

    if (/[0-9]/.test(line[i])) {
      let j = i
      while (j < len && /[0-9]/.test(line[j])) j++
      if (j < len && line[j] === '.') { j++; while (j < len && /[0-9]/.test(line[j])) j++ }
      tokens.push({ type: 'number', text: line.slice(i, j), start: i, end: j })
      i = j
      continue
    }

    if (line[i] === '.') {
      let j = i + 1
      while (j < len && /[a-zA-Z0-9_\u4e00-\u9fff]/.test(line[j])) j++
      const word = line.slice(i, j)

      if (ELANG_KEYWORDS.has(word)) {
        tokens.push({ type: 'keyword', text: word, start: i, end: j })
      } else {
        tokens.push({ type: 'elang', text: word, start: i, end: j })
      }
      i = j
      continue
    }

    if (/[a-zA-Z_\u4e00-\u9fff]/.test(line[i])) {
      let j = i
      while (j < len && /[a-zA-Z0-9_\u4e00-\u9fff]/.test(line[j])) j++
      const word = line.slice(i, j)

      if (ELANG_TYPES.has(word)) {
        tokens.push({ type: 'type', text: word, start: i, end: j })
      } else if (ELANG_KEYWORDS.has(word)) {
        tokens.push({ type: 'keyword', text: word, start: i, end: j })
      } else if (/^[\u4e00-\u9fff]/.test(word)) {
        const nextNonSpace = line.slice(j).trimStart()
        if (nextNonSpace.startsWith('（') || nextNonSpace.startsWith('(')) {
          tokens.push({ type: 'function', text: word, start: i, end: j })
        } else {
          tokens.push({ type: 'elang', text: word, start: i, end: j })
        }
      } else {
        tokens.push({ type: 'plain', text: word, start: i, end: j })
      }
      i = j
      continue
    }

    if (ELANG_OPERATORS.has(line[i])) {
      tokens.push({ type: 'operator', text: line[i], start: i, end: i + 1 })
      i++
      continue
    }

    const opChars = '+-*/%=<>!&|^~.:;,(){}[]（）【】《》'
    if (opChars.includes(line[i])) {
      tokens.push({ type: 'operator', text: line[i], start: i, end: i + 1 })
      i++
      continue
    }

    tokens.push({ type: 'plain', text: line[i], start: i, end: i + 1 })
    i++
  }

  return tokens
}

function tokenizeLine(line: string, language: 'go' | 'elang'): TokenSpan[] {
  return language === 'elang' ? tokenizeELangLine(line) : tokenizeGoLine(line)
}

function extractFirstWord(line: string, language: 'go' | 'elang'): string {
  if (language === 'elang') {
    const match = line.match(/^\.[\u4e00-\u9fff]*[a-zA-Z]*[\u4e00-\u9fff]*/)
    return match ? match[0] : ''
  }
  const match = line.match(/^[a-zA-Z_]+/)
  return match ? match[0] : ''
}

function computeFlowSegments(blocks: RenderBlock[], language: 'go' | 'elang'): { maxDepth: number } {
  const stack: Array<{ depth: number; kind: string }> = []
  let maxDepth = 0

  const flowSet = language === 'elang' ? ELANG_FLOW_KEYWORDS : FLOW_KEYWORDS

  const isBlockEnd = (trimmed: string): boolean => {
    if (language === 'go') {
      return trimmed === '}' || trimmed === '})' || trimmed === '} else {' || trimmed.startsWith('} else')
    }
    return trimmed === '.如果结束' || trimmed === '.判断结束' ||
           trimmed === '.循环尾' || trimmed === '.循环判断尾' ||
           trimmed === '.计次循环尾' || trimmed === '.变量循环尾'
  }

  const isBlockStart = (trimmed: string): boolean => {
    const word = extractFirstWord(trimmed, language)
    return flowSet.has(word) && !isBlockEnd(trimmed)
  }

  const isBranch = (trimmed: string): boolean => {
    if (language === 'go') return trimmed.startsWith('case ') || trimmed.startsWith('default:') || trimmed === 'default:'
    return trimmed === '.否则' || trimmed.startsWith('.否则如果') || trimmed === '.默认'
  }

  for (const block of blocks) {
    if (block.kind !== 'codeline') continue

    const trimmed = block.codeLine.trimStart()
    const indent = block.codeLine.length - trimmed.length

    if (isBlockEnd(trimmed)) {
      if (stack.length > 0) stack.pop()
    }

    const segs: FlowSegment[] = []
    for (let d = 0; d < stack.length; d++) {
      segs.push({ kind: 'line', depth: d })
    }

    if (isBranch(trimmed)) {
      if (stack.length > 0) {
        segs.push({ kind: 'branch', depth: stack.length - 1 })
      }
    } else if (isBlockStart(trimmed)) {
      const depth = stack.length
      segs.push({ kind: 'start', depth })
      stack.push({ depth, kind: extractFirstWord(trimmed, language) })
      if (depth + 1 > maxDepth) maxDepth = depth + 1
    }

    block.flowSegments = segs
    block.indentLevel = Math.floor(indent / 2)
  }

  return { maxDepth }
}

function buildBlocks(text: string, language: 'go' | 'elang'): RenderBlock[] {
  const lines = text.split('\n')
  const blocks: RenderBlock[] = []

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const trimmed = line.trim()

    if (trimmed === '') {
      blocks.push({
        kind: 'blank',
        lineIndex: i,
        codeLine: '',
        tokens: [],
        flowSegments: [],
        indentLevel: 0,
        isVirtual: false,
      })
    } else {
      const tokens = tokenizeLine(line, language)
      blocks.push({
        kind: 'codeline',
        lineIndex: i,
        codeLine: line,
        tokens,
        flowSegments: [],
        indentLevel: 0,
        isVirtual: false,
      })
    }
  }

  computeFlowSegments(blocks, language)
  return blocks
}

function buildHeadLineBlocks(text: string, language: 'go' | 'elang', maxLines: number = 220): RenderBlock[] {
  const limit = Math.max(1, maxLines)
  const lines: string[] = []
  let start = 0

  for (let i = 0; i < text.length; i++) {
    if (text.charCodeAt(i) !== 10) continue
    lines.push(text.slice(start, i))
    if (lines.length >= limit) break
    start = i + 1
  }

  if (lines.length < limit) {
    lines.push(text.slice(start))
  }

  return lines.map((line, idx): RenderBlock => ({
    kind: 'codeline',
    lineIndex: idx,
    codeLine: line,
    tokens: tokenizeLine(line, language),
    flowSegments: [],
    indentLevel: 0,
    isVirtual: false,
  }))
}

self.onmessage = (event: MessageEvent<ParseRequest>) => {
  const { id, text, language: reqLang } = event.data

  try {
    const language: 'go' | 'elang' = (reqLang === 'auto' || !reqLang)
      ? detectLanguage(text)
      : reqLang as 'go' | 'elang'

    const blocks = buildBlocks(text, language)
    const maxDepth = Math.max(...blocks.map(b => b.flowSegments.length), 0)

    const response: ParseResponse = {
      id,
      blocks,
      flowMaxDepth: maxDepth,
      lineCount: blocks.length,
    }

    self.postMessage(response)
  } catch (err) {
    const fallbackBlocks = buildHeadLineBlocks(text, 'go', 500)
    const response: ParseResponse = {
      id,
      blocks: fallbackBlocks,
      flowMaxDepth: 0,
      lineCount: fallbackBlocks.length,
    }
    self.postMessage(response)
  }
}

const _self: Worker = self as any
export {}