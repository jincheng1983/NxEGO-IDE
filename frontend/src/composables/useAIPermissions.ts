import { ref, reactive, computed } from 'vue'

export type OperationType =
  | 'read_file'
  | 'write_file'
  | 'delete_file'
  | 'create_file'
  | 'move_file'
  | 'exec_command'
  | 'network_request'
  | 'git_operation'
  | 'install_package'
  | 'modify_config'

export type RiskLevel = 'safe' | 'moderate' | 'dangerous'

export interface OperationRequest {
  id: string
  type: OperationType
  title: string
  description: string
  targetPath?: string
  targetContent?: string
  command?: string
  risk: RiskLevel
  requiresConfirmation: boolean
  autoApproveSimilar: boolean
  timestamp: number
}

export interface PermissionRule {
  id: string
  type: OperationType
  pathPattern?: string
  allow: boolean
  askEveryTime: boolean
  createdAt: number
}

interface ResolveReject {
  resolve: (approved: boolean) => void
  reject: (err: Error) => void
}

const STORAGE_KEY = 'yigou-ai-permissions'

const pendingRequests = ref<OperationRequest[]>([])
const currentRequest = ref<OperationRequest | null>(null)
const resolveRejectMap = new Map<string, ResolveReject>()

const permissionRules = reactive<PermissionRule[]>(loadRules())
const autoApproveCount = ref(0)
const autoDenyCount = ref(0)
const manualApproveCount = ref(0)

function loadRules(): PermissionRule[] {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    return saved ? JSON.parse(saved) : []
  } catch {
    return []
  }
}

function saveRules() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(permissionRules))
}

function generateId(): string {
  return `op_${Date.now()}_${Math.random().toString(36).substring(2, 7)}`
}

function getRiskLevel(type: OperationType, targetPath?: string): RiskLevel {
  const highRisk: OperationType[] = ['delete_file', 'exec_command', 'git_operation', 'install_package']
  const moderateRisk: OperationType[] = ['write_file', 'move_file', 'modify_config', 'network_request']
  const lowRisk: OperationType[] = ['read_file', 'create_file']

  if (highRisk.includes(type)) return 'dangerous'
  if (moderateRisk.includes(type)) return 'moderate'

  if (type === 'read_file' || type === 'create_file') {
    if (targetPath && (targetPath.includes('/etc/') || targetPath.includes('/system/') || targetPath.includes('C:\\Windows\\'))) {
      return 'moderate'
    }
  }

  return 'safe'
}

function getOperationLabel(type: OperationType): string {
  const labels: Record<OperationType, string> = {
    read_file: '读取文件',
    write_file: '写入文件',
    delete_file: '删除文件',
    create_file: '创建文件',
    move_file: '移动文件',
    exec_command: '执行命令',
    network_request: '网络请求',
    git_operation: 'Git 操作',
    install_package: '安装包',
    modify_config: '修改配置',
  }
  return labels[type] || type
}

function getOperationIcon(type: OperationType): string {
  const icons: Record<OperationType, string> = {
    read_file: '📖',
    write_file: '✏️',
    delete_file: '🗑️',
    create_file: '📄',
    move_file: '📦',
    exec_command: '⚡',
    network_request: '🌐',
    git_operation: '🔀',
    install_package: '📥',
    modify_config: '⚙️',
  }
  return icons[type] || '🔧'
}

function matchRule(type: OperationType, targetPath?: string): PermissionRule | null {
  for (const rule of permissionRules) {
    if (rule.type !== type) continue
    if (rule.pathPattern && targetPath) {
      try {
        const regex = new RegExp(rule.pathPattern.replace(/\*/g, '.*').replace(/\?/g, '.'), 'i')
        if (regex.test(targetPath)) return rule
      } catch {
        continue
      }
    } else if (rule.pathPattern && !targetPath) {
      continue
    } else {
      return rule
    }
  }
  return null
}

async function requestPermission(opts: {
  type: OperationType
  title?: string
  description?: string
  targetPath?: string
  targetContent?: string
  command?: string
  timeout?: number
}): Promise<boolean> {
  const risk = getRiskLevel(opts.type, opts.targetPath)

  const matchedRule = matchRule(opts.type, opts.targetPath)
  if (matchedRule) {
    if (matchedRule.askEveryTime) {
      // fall through to dialog
    } else if (matchedRule.allow) {
      autoApproveCount.value++
      return true
    } else {
      autoDenyCount.value++
      return false
    }
  }

  const request: OperationRequest = {
    id: generateId(),
    type: opts.type,
    title: opts.title || getOperationLabel(opts.type),
    description: opts.description || `AI 请求执行：${getOperationLabel(opts.type)}`,
    targetPath: opts.targetPath,
    targetContent: opts.targetContent,
    command: opts.command,
    risk,
    requiresConfirmation: risk !== 'safe',
    autoApproveSimilar: false,
    timestamp: Date.now(),
  }

  if (risk === 'safe' && permissionRules.some(r => r.type === opts.type && r.allow && !r.askEveryTime)) {
    autoApproveCount.value++
    return true
  }

  pendingRequests.value.push(request)
  currentRequest.value = request

  return new Promise<boolean>((resolve, reject) => {
    const timeout = opts.timeout || 60000
    const timeoutId = setTimeout(() => {
      resolveRejectMap.delete(request.id)
      removeRequest(request.id)
      reject(new Error('操作确认超时'))
    }, timeout)

    resolveRejectMap.set(request.id, {
      resolve: (approved: boolean) => {
        clearTimeout(timeoutId)
        resolve(approved)
      },
      reject: (err: Error) => {
        clearTimeout(timeoutId)
        reject(err)
      },
    })
  })
}

function approveRequest(requestId: string, saveRule = false, askEveryTime = false) {
  const rr = resolveRejectMap.get(requestId)
  if (rr) {
    resolveRejectMap.delete(requestId)

    const request = currentRequest.value
    if (request && saveRule) {
      permissionRules.push({
        id: generateId(),
        type: request.type,
        pathPattern: request.autoApproveSimilar ? request.targetPath || '' : '',
        allow: true,
        askEveryTime,
        createdAt: Date.now(),
      })
      saveRules()
    }

    manualApproveCount.value++
    rr.resolve(true)
  }
  removeRequest(requestId)
}

function denyRequest(requestId: string) {
  const rr = resolveRejectMap.get(requestId)
  if (rr) {
    resolveRejectMap.delete(requestId)
    rr.resolve(false)
  }
  removeRequest(requestId)
}

function removeRequest(requestId: string) {
  pendingRequests.value = pendingRequests.value.filter(r => r.id !== requestId)
  if (currentRequest.value?.id === requestId) {
    currentRequest.value = pendingRequests.value[0] || null
  }
}

function cleanupTimeout() {
  const now = Date.now()
  for (const request of pendingRequests.value) {
    if (now - request.timestamp > 120000) {
      const rr = resolveRejectMap.get(request.id)
      if (rr) {
        resolveRejectMap.delete(request.id)
        rr.reject(new Error('请求超时已自动取消'))
      }
    }
  }
  pendingRequests.value = pendingRequests.value.filter(r => {
    const rr = resolveRejectMap.get(r.id)
    return rr !== undefined
  })
}

function addPermissionRule(rule: PermissionRule) {
  permissionRules.push(rule)
  saveRules()
}

function removePermissionRule(ruleId: string) {
  const idx = permissionRules.findIndex(r => r.id === ruleId)
  if (idx >= 0) {
    permissionRules.splice(idx, 1)
    saveRules()
  }
}

function clearAllRules() {
  permissionRules.length = 0
  saveRules()
}

export function useAIPermissions() {
  cleanupTimeout()

  return {
    pendingRequests: computed(() => [...pendingRequests.value]),
    currentRequest: computed(() => currentRequest.value),
    permissionRules: computed(() => [...permissionRules]),
    hasPending: computed(() => pendingRequests.value.length > 0),
    autoApproveCount: computed(() => autoApproveCount.value),
    autoDenyCount: computed(() => autoDenyCount.value),
    manualApproveCount: computed(() => manualApproveCount.value),
    requestPermission,
    approveRequest,
    denyRequest,
    addPermissionRule,
    removePermissionRule,
    clearAllRules,
    getOperationLabel,
    getOperationIcon,
    getRiskLevel,
  }
}