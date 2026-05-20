import { ref, watch, reactive } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'custom'

export interface EditorSettings {
  fontSize: number
  tabSize: number
  lineHeight: number
  wordWrap: boolean
  showLineNumbers: boolean
  fontFamily: string
  autoComplete: boolean
  minimap: boolean
}

export interface GeneralSettings {
  language: string
  autoSave: boolean
  autoSaveInterval: number
  confirmOnExit: boolean
  showWelcomePage: boolean
}

export interface AISettings {
  provider: 'local' | 'ollama' | 'cloud'
  localModelPath: string
  localModelName: string
  ollamaUrl: string
  ollamaModel: string
  cloudApiKey: string
  cloudApiUrl: string
  cloudModel: string
  enableGPU: boolean
  gpuLayers: number
  maxTokens: number
  temperature: number
  topP: number
  contextSize: number
  enableCodeCompletion: boolean
  completionDelay: number
  enableChatHistory: boolean
  systemPrompt: string
}

export interface CustomThemeColors {
  bgPrimary: string
  bgSecondary: string
  bgTertiary: string
  textPrimary: string
  textRegular: string
  colorPrimary: string
  editorBg: string
  editorText: string
  hlKeyword: string
  hlString: string
  hlComment: string
  hlElang: string
  hlFunction: string
  hlNumber: string
}

export interface AppSettings {
  theme: ThemeMode
  customColors: CustomThemeColors
  editor: EditorSettings
  general: GeneralSettings
  ai: AISettings
}

export interface ProjectSettings {
  projectId: string
  projectName: string
  startupWindow: string
  targetPlatform: string
  buildMode: string
  outputDir: string
  iconPath: string
  version: string
  injectVersion: boolean
  compress: boolean
  obfuscate: boolean
  debug: boolean
}

const defaultEditorSettings: EditorSettings = {
  fontSize: 14,
  tabSize: 4,
  lineHeight: 1.6,
  wordWrap: false,
  showLineNumbers: true,
  fontFamily: "'Consolas', 'Courier New', 'Microsoft YaHei', monospace",
  autoComplete: true,
  minimap: false,
}

const defaultGeneralSettings: GeneralSettings = {
  language: 'zh-CN',
  autoSave: true,
  autoSaveInterval: 60,
  confirmOnExit: true,
  showWelcomePage: true,
}

const defaultAISettings: AISettings = {
  provider: 'local',
  localModelPath: '',
  localModelName: 'qwen2.5-coder-1.5b-q4_k_m',
  ollamaUrl: 'http://localhost:11434',
  ollamaModel: 'qwen2.5-coder:1.5b',
  cloudApiKey: '',
  cloudApiUrl: 'https://api.openai.com/v1',
  cloudModel: 'gpt-4o-mini',
  enableGPU: true,
  gpuLayers: 28,
  maxTokens: 2048,
  temperature: 0.7,
  topP: 0.9,
  contextSize: 32768,
  enableCodeCompletion: true,
  completionDelay: 150,
  enableChatHistory: true,
  systemPrompt: '你是易狗 NxEGO IDE 的 AI 编程助手，精通 Go 语言和中文函数映射体系。你可以帮助用户编写代码、排查问题、优化逻辑，并支持中文函数和原生 Go 代码的双向转换。',
}

const defaultCustomColors: CustomThemeColors = {
  bgPrimary: '#1e1e1e',
  bgSecondary: '#252526',
  bgTertiary: '#2d2d30',
  textPrimary: '#cccccc',
  textRegular: '#999999',
  colorPrimary: '#0078d4',
  editorBg: '#1e1e1e',
  editorText: '#d4d4d4',
  hlKeyword: '#569cd6',
  hlString: '#ce9178',
  hlComment: '#6a9955',
  hlElang: '#c586c0',
  hlFunction: '#dcdcaa',
  hlNumber: '#b5cea8',
}

const STORAGE_KEY = 'yigou-ide-settings'

function loadSettings(): AppSettings {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const parsed = JSON.parse(saved)
      return {
        theme: parsed.theme || 'light',
        customColors: { ...defaultCustomColors, ...parsed.customColors },
        editor: { ...defaultEditorSettings, ...parsed.editor },
        general: { ...defaultGeneralSettings, ...parsed.general },
        ai: { ...defaultAISettings, ...parsed.ai },
      }
    }
  } catch (e) {
    console.warn('加载设置失败，使用默认设置', e)
  }
  return {
    theme: 'light',
    customColors: { ...defaultCustomColors },
    editor: { ...defaultEditorSettings },
    general: { ...defaultGeneralSettings },
    ai: { ...defaultAISettings },
  }
}

function saveSettings(settings: AppSettings) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings))
  } catch (e) {
    console.warn('保存设置失败', e)
  }
}

const PROJECT_SETTINGS_KEY = 'yigou-ide-project-settings'

const defaultProjectSettings: ProjectSettings = {
  projectId: '',
  projectName: '',
  startupWindow: '',
  targetPlatform: 'windows/amd64',
  buildMode: 'release',
  outputDir: 'build',
  iconPath: '',
  version: '1.0.0',
  injectVersion: true,
  compress: false,
  obfuscate: false,
  debug: false,
}

function loadAllProjectSettings(): Record<string, ProjectSettings> {
  try {
    const saved = localStorage.getItem(PROJECT_SETTINGS_KEY)
    if (saved) {
      return JSON.parse(saved)
    }
  } catch (e) {
    console.warn('加载项目设置失败', e)
  }
  return {}
}

function saveAllProjectSettings(allSettings: Record<string, ProjectSettings>) {
  try {
    localStorage.setItem(PROJECT_SETTINGS_KEY, JSON.stringify(allSettings))
  } catch (e) {
    console.warn('保存项目设置失败', e)
  }
}

function loadProjectSettings(projectId: string): ProjectSettings {
  const all = loadAllProjectSettings()
  if (all[projectId]) {
    return { ...defaultProjectSettings, ...all[projectId] }
  }
  return { ...defaultProjectSettings, projectId }
}

function saveProjectSettings(ps: ProjectSettings) {
  const all = loadAllProjectSettings()
  all[ps.projectId] = ps
  saveAllProjectSettings(all)
}

function deleteProjectSettings(projectId: string) {
  const all = loadAllProjectSettings()
  delete all[projectId]
  saveAllProjectSettings(all)
}

const settings = reactive<AppSettings>(loadSettings())

const currentTheme = ref<ThemeMode>(settings.theme)

function applyTheme(theme: ThemeMode) {
  document.documentElement.setAttribute('data-theme', theme)

  if (theme === 'custom' && settings.customColors) {
    const c = settings.customColors
    const root = document.documentElement.style
    root.setProperty('--bg-primary', c.bgPrimary)
    root.setProperty('--bg-secondary', c.bgSecondary)
    root.setProperty('--bg-tertiary', c.bgTertiary)
    root.setProperty('--text-primary', c.textPrimary)
    root.setProperty('--text-regular', c.textRegular)
    root.setProperty('--color-primary', c.colorPrimary)
    root.setProperty('--editor-bg', c.editorBg)
    root.setProperty('--editor-text', c.editorText)
    root.setProperty('--hl-keyword', c.hlKeyword)
    root.setProperty('--hl-string', c.hlString)
    root.setProperty('--hl-comment', c.hlComment)
    root.setProperty('--hl-elang', c.hlElang)
    root.setProperty('--hl-function', c.hlFunction)
    root.setProperty('--hl-number', c.hlNumber)
  }
}

function setTheme(theme: ThemeMode) {
  currentTheme.value = theme
  settings.theme = theme
  applyTheme(theme)
  saveSettings({ ...settings })
}

function toggleTheme() {
  const next = currentTheme.value === 'light' ? 'dark' : 'light'
  setTheme(next)
}

function updateEditorSettings(partial: Partial<EditorSettings>) {
  Object.assign(settings.editor, partial)
  saveSettings({ ...settings })
}

function updateGeneralSettings(partial: Partial<GeneralSettings>) {
  Object.assign(settings.general, partial)
  saveSettings({ ...settings })
}

function updateAISettings(partial: Partial<AISettings>) {
  Object.assign(settings.ai, partial)
  saveSettings({ ...settings })
}

function updateCustomColors(partial: Partial<CustomThemeColors>) {
  Object.assign(settings.customColors, partial)
  if (currentTheme.value === 'custom') {
    applyTheme('custom')
  }
  saveSettings({ ...settings })
}

function resetSettings() {
  settings.theme = 'light'
  settings.customColors = { ...defaultCustomColors }
  settings.editor = { ...defaultEditorSettings }
  settings.general = { ...defaultGeneralSettings }
  settings.ai = { ...defaultAISettings }
  currentTheme.value = 'light'
  applyTheme('light')
  saveSettings({ ...settings })
}

function initTheme() {
  applyTheme(currentTheme.value)
}

watch(currentTheme, (newTheme) => {
  applyTheme(newTheme)
})

export function useTheme() {
  return {
    currentTheme,
    settings,
    toggleTheme,
    setTheme,
    initTheme,
    updateEditorSettings,
    updateGeneralSettings,
    updateAISettings,
    updateCustomColors,
    resetSettings,
    defaultEditorSettings,
    defaultGeneralSettings,
    defaultAISettings,
    defaultCustomColors,
    defaultProjectSettings,
    loadProjectSettings,
    saveProjectSettings,
    deleteProjectSettings,
    loadAllProjectSettings,
  }
}