import { ref, shallowRef, onUnmounted } from 'vue'
import type { ParseRequest, ParseResponse, RenderBlock, ParseLanguage } from '../workers/goParserShared'

interface PendingRequest {
  id: number
  resolve: (response: ParseResponse) => void
  reject: (error: Error) => void
}

let workerInstance: Worker | null = null
let requestIdCounter = 0
const pendingRequests = new Map<number, PendingRequest>()

function getWorker(): Worker {
  if (!workerInstance) {
    workerInstance = new Worker(
      new URL('../workers/goParser.worker.ts', import.meta.url),
      { type: 'module' }
    )

    workerInstance.onmessage = (event: MessageEvent<ParseResponse>) => {
      const { id } = event.data
      const pending = pendingRequests.get(id)
      if (pending) {
        pending.resolve(event.data)
        pendingRequests.delete(id)
      }
    }

    workerInstance.onerror = (err) => {
      console.error('Go Parser Worker error:', err)
      pendingRequests.forEach((pending) => {
        pending.reject(new Error('Worker error'))
      })
      pendingRequests.clear()
    }
  }
  return workerInstance
}

function parseInWorker(text: string, language: ParseLanguage = 'auto'): Promise<ParseResponse> {
  return new Promise((resolve, reject) => {
    const id = ++requestIdCounter
    const worker = getWorker()

    pendingRequests.set(id, { id, resolve, reject })

    const request: ParseRequest = {
      id,
      text,
      isClassModule: false,
      language,
    }

    worker.postMessage(request)

    setTimeout(() => {
      if (pendingRequests.has(id)) {
        pendingRequests.delete(id)
        reject(new Error('Parse timeout'))
      }
    }, 10000)
  })
}

export function useGoParser() {
  const blocks = shallowRef<RenderBlock[]>([])
  const flowMaxDepth = ref(0)
  const lineCount = ref(0)
  const isParsing = ref(false)
  const parseError = ref<string | null>(null)
  let latestRequestId = 0

  const parse = async (text: string, language: ParseLanguage = 'auto') => {
    const requestId = ++latestRequestId
    isParsing.value = true
    parseError.value = null

    try {
      const response = await parseInWorker(text, language)

      if (requestId !== latestRequestId) return

      blocks.value = response.blocks
      flowMaxDepth.value = response.flowMaxDepth
      lineCount.value = response.lineCount
    } catch (err) {
      if (requestId !== latestRequestId) return
      parseError.value = err instanceof Error ? err.message : 'Parse failed'
      console.error('Go parser error:', err)
    } finally {
      if (requestId === latestRequestId) {
        isParsing.value = false
      }
    }
  }

  const terminate = () => {
    if (workerInstance) {
      workerInstance.terminate()
      workerInstance = null
      pendingRequests.clear()
    }
  }

  onUnmounted(() => {
    terminate()
  })

  return {
    blocks,
    flowMaxDepth,
    lineCount,
    isParsing,
    parseError,
    parse,
    terminate,
  }
}

export function terminateGoParserWorker() {
  if (workerInstance) {
    workerInstance.terminate()
    workerInstance = null
    pendingRequests.clear()
  }
}