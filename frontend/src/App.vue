<script lang="ts" setup>
import { ref, onMounted, nextTick } from 'vue'
import TopMenuBar from './components/layout/TopMenuBar.vue'
import LeftPanelWrapper from './components/designer/LeftPanelWrapper.vue'
import CenterTabs from './components/editor/CenterTabs.vue'
import RightPanel from './components/designer/RightPanel.vue'
import BottomPanel from './components/editor/BottomPanel.vue'
import ProjectManager from './components/project/ProjectManager.vue'
import BuildDialog from './components/build/BuildDialog.vue'
import StatusBar from './components/common/StatusBar.vue'
import SettingsPanel from './components/settings/SettingsPanel.vue'
import OperationConfirmDialog from './components/ai/OperationConfirmDialog.vue'
import { SaveFile, LoadFile, OpenDialog, SaveDialog, ExportGoCode, ConvertCodeRowsToGo, CompileAndRun, StopProgram, IsProgramRunning, BuildExe, SelectSavePath } from '../wailsjs/go/main/App'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useTheme } from './stores/theme'
import './styles/theme.css'

const { currentTheme, toggleTheme, initTheme } = useTheme()

const centerTabsRef = ref()
const leftPanelRef = ref()
const rightPanelRef = ref()
const bottomPanelRef = ref()

const getActiveCanvasRef = () => {
  const tabs = centerTabsRef.value
  if (!tabs) return null
  return tabs.canvasRefs?.[tabs.activeWindowId] || null
}

const getActiveCodeEditorRef = () => {
  const tabs = centerTabsRef.value
  if (!tabs) return null
  return tabs.codeEditorRefs?.[tabs.activeWindowId] || null
}
const currentFilePath = ref('')
const projectName = ref('未命名项目')
const showPropertyPanel = ref(false)
const showBottomPanel = ref(true)
const showProjectManager = ref(false)
const showBuildDialog = ref(false)
const buildDialogGoCode = ref('')
const showSettings = ref(false)
const selectedComponent = ref<any>(null)
const selectedComponents = ref<any[]>([])
const currentProject = ref<any>(null)
const hasProject = ref(false)
const isProgramRunning = ref(false)
const goVersion = ref('-')
const codeLineCount = ref(0)
const windowProps = ref({
  title: '我的窗口',
  width: 800,
  height: 600,
  windowType: 'normal',
  backgroundColor: '#f0f0f0',
  resizable: true,
  showInTaskbar: true,
  showIcon: true,
  iconPath: '',
  minWidth: 200,
  minHeight: 150,
  maxWidth: 0,
  maxHeight: 0,
  opacity: 100,
  alwaysOnTop: false,
  borderStyle: 'sizable',
  showMenuBar: false,
  showMinBtn: true,
  showMaxBtn: true,
  showCloseBtn: true,
  titleBarHeight: 32,
})

const handleMenuAction = async (action: string) => {
  switch (action) {
    case 'new-project':
      showProjectManager.value = true
      break

    case 'open-project':
      showProjectManager.value = true
      break

    case 'save-project':
      await saveProject()
      break

    case 'save-as':
      await saveProjectAs()
      break

    case 'export-go':
      await exportGoCode()
      break

    case 'compile-run':
      await runNormal()
      break

    case 'run-normal':
      await runNormal()
      break

    case 'debug-run':
      await runNormal()
      break

    case 'step-over':
      bottomPanelRef.value?.addConsoleLog('单步跳过 (F10)')
      ElMessage.info('调试功能：单步跳过')
      break

    case 'step-into':
      bottomPanelRef.value?.addConsoleLog('单步进入 (F11)')
      ElMessage.info('调试功能：单步进入')
      break

    case 'stop-program':
      await stopProgram()
      break

    case 'build-exe':
      await buildExe()
      break

    case 'build-config':
      await buildExe()
      break

    case 'toggle-component-lib':
      break

    case 'toggle-property-panel':
      if (rightPanelRef.value) {
        rightPanelRef.value.setExpanded(!rightPanelRef.value.isExpanded())
      }
      break

    case 'toggle-bottom-panel':
      showBottomPanel.value = !showBottomPanel.value
      break

    case 'toggle-theme':
      toggleTheme()
      ElMessage.info(`已切换到${currentTheme.value === 'dark' ? '深色' : '浅色'}主题`)
      break

    case 'open-settings':
      showSettings.value = true
      break

    case 'about':
      ElMessageBox.alert(
        '易狗 IDE v1.0.0\n\n一款全中文可视化Go语言桌面应用开发工具。\n拖拽设计UI，中文编写逻辑，一键编译打包。',
        '关于易狗',
        { confirmButtonText: '确定' }
      )
      break

    case 'exit':
      try {
        await ElMessageBox.confirm('确定要退出易狗 IDE 吗？', '退出', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        })
        window.close()
      } catch {
        // 用户取消
      }
      break
  }
}

const getProjectDir = () => {
  if (currentProject.value && currentProject.value.path && currentProject.value.name) {
    return currentProject.value.path + '\\' + currentProject.value.name
  }
  return ''
}

const getAllCodeRows = () => {
  const editor = getActiveCodeEditorRef()
  if (!editor) return []
  if (editor.getAllCodeRows && typeof editor.getAllCodeRows === 'function') {
    return editor.getAllCodeRows()
  }
  return []
}

const getCodeContent = () => {
  const editor = getActiveCodeEditorRef()
  if (!editor) return ''
  if (editor.getValue && typeof editor.getValue === 'function') {
    return editor.getValue()
  }
  return ''
}

const setCodeContent = (code: string) => {
  const editor = getActiveCodeEditorRef()
  if (!editor) return
  if (editor.setValue && typeof editor.setValue === 'function') {
    editor.setValue(code)
  }
}

const saveProject = async () => {
  if (!currentFilePath.value) {
    await saveProjectAs()
    return
  }
  try {
    const components = getActiveCanvasRef()?.getComponents() || []
    const codeRows = getAllCodeRows()
    const codeContent = getCodeContent()
    const data = JSON.stringify({
      appName: projectName.value,
      version: '1.0.0',
      components,
      codeRows,
      codeContent,
      windowProps: windowProps.value,
    })
    await SaveFile(currentFilePath.value, data)
    bottomPanelRef.value?.addConsoleLog(`项目已保存: ${currentFilePath.value}`)
    ElMessage.success('项目已保存')
  } catch (e: any) {
    ElMessage.error(`保存失败: ${e}`)
  }
}

const saveProjectAs = async () => {
  try {
    const path = await SaveDialog()
    if (path) {
      currentFilePath.value = path
      await saveProject()
    }
  } catch (e: any) {
    ElMessage.error(`保存失败: ${e}`)
  }
}

const exportGoCode = async () => {
  try {
    const goCode = await prepareGoCode()
    bottomPanelRef.value?.addConsoleLog('Go代码已导出')
    bottomPanelRef.value?.addConsoleLog(goCode)
    ElMessage.success('Go代码已导出到控制台')
  } catch (e: any) {
    ElMessage.error(`导出失败: ${e}`)
  }
}

const prepareGoCode = async () => {
  const components = getActiveCanvasRef()?.getComponents() || []
  const codeContent = getCodeContent()
  const codeRows = getAllCodeRows()

  if (codeContent.trim().length > 0) {
    return codeContent
  }

  const data = JSON.stringify({
    appName: projectName.value,
    version: '1.0.0',
    components,
    codeRows,
    codeContent: '',
    windowProps: windowProps.value,
  })
  return await ExportGoCode(data)
}

const runNormal = async () => {
  try {
    isProgramRunning.value = true
    bottomPanelRef.value?.addConsoleLog('▶ 正常模式运行...')
    ElMessage.info('正在编译运行（正常模式）...')

    const goCode = await prepareGoCode()
    codeLineCount.value = goCode.split('\n').length
    const projectDir = getProjectDir()

    const result = await CompileAndRun(goCode, projectName.value, projectDir, false)

    if (result.success) {
      bottomPanelRef.value?.addConsoleLog('✅ 程序启动成功！（正常模式）')
      if (result.exePath) {
        bottomPanelRef.value?.addConsoleLog(`📁 编译输出: ${result.exePath}`)
      }
      if (result.output) {
        bottomPanelRef.value?.addConsoleLog(`程序输出: ${result.output}`)
      }
      ElMessage.success('程序运行成功！')
    } else {
      isProgramRunning.value = false
      bottomPanelRef.value?.addConsoleLog(`❌ 运行失败: ${result.message}`)
      if (result.output) {
        bottomPanelRef.value?.addConsoleLog(`错误详情: ${result.output}`)
      }
      ElMessage.error(result.message)
    }
  } catch (e: any) {
    isProgramRunning.value = false
    bottomPanelRef.value?.addConsoleLog(`运行异常: ${e}`)
    ElMessage.error('运行失败')
  }
}

const runDebug = async () => {
  try {
    isProgramRunning.value = true
    bottomPanelRef.value?.addConsoleLog('🐛 调试模式运行...')
    ElMessage.info('正在调试运行...')

    const goCode = await prepareGoCode()
    codeLineCount.value = goCode.split('\n').length
    const projectDir = getProjectDir()

    const result = await CompileAndRun(goCode, projectName.value, projectDir, true)

    if (result.success) {
      bottomPanelRef.value?.addConsoleLog('✅ 调试程序启动成功！')
      if (result.exePath) {
        bottomPanelRef.value?.addConsoleLog(`📁 编译输出: ${result.exePath}`)
      }
      bottomPanelRef.value?.addConsoleLog('💡 提示：调试模式已禁用优化，便于排查问题')
      if (result.output) {
        bottomPanelRef.value?.addConsoleLog(`调试输出: ${result.output}`)
      }
      ElMessage.success('调试模式已启动！')
    } else {
      isProgramRunning.value = false
      bottomPanelRef.value?.addConsoleLog(`❌ 调试失败: ${result.message}`)
      if (result.output) {
        bottomPanelRef.value?.addConsoleLog(`错误详情: ${result.output}`)
      }
      ElMessage.error(result.message)
    }
  } catch (e: any) {
    isProgramRunning.value = false
    bottomPanelRef.value?.addConsoleLog(`调试异常: ${e}`)
    ElMessage.error('调试失败')
  }
}

const compileAndRun = async () => {
  try {
    bottomPanelRef.value?.addConsoleLog('正在编译运行...')
    ElMessage.info('正在编译运行，请稍候...')

    const goCode = await prepareGoCode()
    const projectDir = getProjectDir()

    const result = await CompileAndRun(goCode, projectName.value, projectDir, false)

    if (result.success) {
      bottomPanelRef.value?.addConsoleLog('编译运行成功！')
      if (result.output) {
        bottomPanelRef.value?.addConsoleLog(`程序输出: ${result.output}`)
      }
      ElMessage.success('编译运行成功！')
    } else {
      bottomPanelRef.value?.addConsoleLog(`编译失败: ${result.message}`)
      if (result.output) {
        bottomPanelRef.value?.addConsoleLog(`错误详情: ${result.output}`)
      }
      ElMessage.error(result.message)
    }
  } catch (e: any) {
    bottomPanelRef.value?.addConsoleLog(`编译运行异常: ${e}`)
    ElMessage.error(`编译运行失败: ${e}`)
  }
}

const stopProgram = async () => {
  try {
    await StopProgram()
    isProgramRunning.value = false
    bottomPanelRef.value?.addConsoleLog('程序已停止')
    ElMessage.info('程序已停止')
  } catch (e: any) {
    ElMessage.error(`停止失败: ${e}`)
  }
}

const buildExe = async () => {
  try {
    const goCode = await prepareGoCode()
    buildDialogGoCode.value = goCode
    showBuildDialog.value = true
  } catch (e: any) {
    bottomPanelRef.value?.addConsoleLog(`准备打包失败: ${e}`)
    ElMessage.error(`准备打包失败: ${e}`)
  }
}

const onComponentSelect = (component: any) => {
  getActiveCanvasRef()?.addComponent(component)
  bottomPanelRef.value?.addConsoleLog(`已添加组件: ${component.name}`)
}

const onComponentFocus = (component: any) => {
  selectedComponent.value = component
  const canvasRef = getActiveCanvasRef()
  if (canvasRef) {
    selectedComponents.value = canvasRef.getSelectedComponents() || []
  } else {
    selectedComponents.value = component ? [component] : []
  }
}

const onPropertyChange = (id: string, key: string, value: any) => {
  getActiveCanvasRef()?.updateComponentProperty(id, key, value)
}

const onPropertyBatchChange = (ids: string[], key: string, value: any) => {
  ids.forEach(id => {
    getActiveCanvasRef()?.updateComponentProperty(id, key, value)
  })
}

const onWindowPropertyChange = (key: string, value: any) => {
  windowProps.value = { ...windowProps.value, [key]: value }
  getActiveCanvasRef()?.updateWindowProperty(key, value)
}

const onSelectFunction = (fn: any) => {
  bottomPanelRef.value?.addConsoleLog(`已选择函数: ${fn.chineseName} - ${fn.description}`)
}

const onSelectWindowType = (type: string) => {
  windowProps.value = { ...windowProps.value, windowType: type }
  getActiveCanvasRef()?.updateWindowProperty('windowType', type)
  bottomPanelRef.value?.addConsoleLog(`窗口类型已切换为: ${type}`)
}

const onOpenWindowDesign = (windowId: string) => {
  centerTabsRef.value?.switchWindow(windowId)
  const win = centerTabsRef.value?.windows?.find((w: any) => w.id === windowId)
  if (win) win.subTab = 'design'
  bottomPanelRef.value?.addConsoleLog(`已切换到窗体设计视图: ${windowId}`)
}

const onOpenWindowCode = (windowId: string, codeFileId?: string) => {
  centerTabsRef.value?.switchWindow(windowId)
  const win = centerTabsRef.value?.windows?.find((w: any) => w.id === windowId)
  if (win) {
    win.subTab = 'code'
    if (codeFileId) win.activeCodeFileId = codeFileId
  }
  bottomPanelRef.value?.addConsoleLog(`已打开窗体代码编辑器: ${windowId}`)
}

const onAddNewWindow = (windowData: any) => {
  centerTabsRef.value?.addWindow()
  bottomPanelRef.value?.addConsoleLog(`已创建新窗体: ${windowData.name}`)
}

const onDeleteWindow = (windowId: string) => {
  centerTabsRef.value?.removeWindow(windowId)
  bottomPanelRef.value?.addConsoleLog(`已删除窗体: ${windowId}`)
}

const onRenameWindow = (windowId: string, newName: string) => {
  const win = centerTabsRef.value?.windows?.find((w: any) => w.id === windowId)
  if (win) win.name = newName + '.ego'
  bottomPanelRef.value?.addConsoleLog(`窗体已重命名: ${newName}`)
}

const onAddIndependentCodeFile = (codeFile: any) => {
  centerTabsRef.value?.addIndependentCodeFile(codeFile.id, codeFile.name)
  bottomPanelRef.value?.addConsoleLog(`已创建共享代码文件: ${codeFile.name}`)
}

const onDeleteIndependentCodeFile = (fileId: string) => {
  centerTabsRef.value?.removeIndependentCodeFile(fileId)
  bottomPanelRef.value?.addConsoleLog(`已删除共享代码文件: ${fileId}`)
}

const onRenameIndependentCodeFile = (fileId: string, newName: string) => {
  centerTabsRef.value?.renameIndependentCodeFile(fileId, newName + '.ego')
  bottomPanelRef.value?.addConsoleLog(`共享代码文件已重命名: ${newName}`)
}

const onOpenIndependentCodeFile = (fileId: string) => {
  centerTabsRef.value?.switchIndependentFile(fileId)
  bottomPanelRef.value?.addConsoleLog(`已打开共享代码文件: ${fileId}`)
}

const onProjectOpen = async (project: any) => {
  showProjectManager.value = false
  currentProject.value = project
  projectName.value = project.name
  hasProject.value = true

  try {
    const baseDir = project.path + '\\' + project.name
    // 新结构: src/项目名.ego，旧结构: 项目名.ego
    let egoPath = baseDir + '\\src\\' + project.name + '.ego'
    let data = ''
    try {
      data = await LoadFile(egoPath)
    } catch {
      egoPath = baseDir + '\\' + project.name + '.ego'
      data = await LoadFile(egoPath)
    }
    const egoData = JSON.parse(data)
    if (egoData.components) {
      getActiveCanvasRef()?.loadComponents(egoData.components)
    }
    if (egoData.codeContent && egoData.codeContent.trim()) {
      const editor = getActiveCodeEditorRef()
      if (editor && editor.setValue && typeof editor.setValue === 'function') {
        editor.setValue(egoData.codeContent)
      }
    } else if (egoData.codeRows && getActiveCodeEditorRef()) {
      const editor = getActiveCodeEditorRef()
      if (editor && editor.setValue && typeof editor.setValue === 'function') {
        if (typeof egoData.codeRows === 'string') {
          editor.setValue(egoData.codeRows)
        } else if (Array.isArray(egoData.codeRows)) {
          const codeText = egoData.codeRows.map((row: any) => row.code || '').join('\n')
          editor.setValue(codeText)
        }
      }
    }
    if (egoData.windowProps) {
      Object.assign(windowProps.value, egoData.windowProps)
    }
    await nextTick()
    if (egoData.components && !egoData.codeRows) {
      centerTabsRef.value?.syncComponentsToCodeEditor(egoData.components)
    }
    currentFilePath.value = egoPath
    bottomPanelRef.value?.addConsoleLog(`已打开项目: ${project.name}`)
    ElMessage.success(`已打开项目: ${project.name}`)
  } catch (e: any) {
    bottomPanelRef.value?.addConsoleLog(`已创建新项目: ${project.name}`)
    ElMessage.success(`已创建新项目: ${project.name}`)
  }
}

const onProjectManagerClose = () => {
  if (!hasProject.value) {
    ElMessage.warning('请先打开或创建一个项目')
    return
  }
  showProjectManager.value = false
}

const onBuildResult = (result: any) => {
  if (result.success) {
    bottomPanelRef.value?.addConsoleLog(`构建成功: ${result.message}`)
    if (result.exePath) {
      bottomPanelRef.value?.addConsoleLog(`输出文件: ${result.exePath}`)
    }
  } else {
    bottomPanelRef.value?.addConsoleLog(`构建失败: ${result.message}`)
  }
}

const onToggleLeftPanel = () => {
  if (leftPanelRef.value) {
    leftPanelRef.value.toggleSidebar()
  }
}

const onToggleRightPanel = () => {
  if (rightPanelRef.value) {
    rightPanelRef.value.setExpanded(!rightPanelRef.value.isExpanded())
  }
}

onMounted(() => {
  initTheme()
  showProjectManager.value = true
})
</script>

<template>
  <div class="app-container">
    <TopMenuBar @menu-action="handleMenuAction" />

    <div class="main-content">
      <LeftPanelWrapper
        ref="leftPanelRef"
        :window-props="windowProps"
        :selected-component="selectedComponent"
        :selected-components="selectedComponents"
        @select-component="onComponentSelect"
        @select-function="onSelectFunction"
        @select-window-type="onSelectWindowType"
        @open-window-design="onOpenWindowDesign"
        @open-window-code="onOpenWindowCode"
        @add-new-window="onAddNewWindow"
        @delete-window="onDeleteWindow"
        @rename-window="onRenameWindow"
        @add-independent-code-file="onAddIndependentCodeFile"
        @delete-independent-code-file="onDeleteIndependentCodeFile"
        @rename-independent-code-file="onRenameIndependentCodeFile"
        @open-independent-code-file="onOpenIndependentCodeFile"
        @window-property-change="onWindowPropertyChange"
        @property-change="onPropertyChange"
        @property-batch-change="onPropertyBatchChange"
      />

      <div class="center-wrapper">
        <CenterTabs
          ref="centerTabsRef"
          @component-focus="onComponentFocus"
        />
        <BottomPanel
          v-if="showBottomPanel"
          ref="bottomPanelRef"
        />
      </div>

      <RightPanel
        ref="rightPanelRef"
      />
    </div>

    <StatusBar
      :project-name="projectName"
      :is-running="isProgramRunning"
      :go-version="goVersion"
      :line-count="codeLineCount"
      @toggle-left-panel="onToggleLeftPanel"
      @toggle-right-panel="onToggleRightPanel"
      @toggle-bottom-panel="showBottomPanel = !showBottomPanel"
    />

    <ProjectManager
      v-if="showProjectManager"
      :must-select="!hasProject"
      @open-project="onProjectOpen"
      @close="onProjectManagerClose"
    />

    <BuildDialog
      v-if="showBuildDialog"
      :go-code="buildDialogGoCode"
      @close="showBuildDialog = false"
      @build-result="onBuildResult"
    />

    <SettingsPanel
      v-if="showSettings"
      @close="showSettings = false"
    />
    <OperationConfirmDialog />
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  width: 100%;
  overflow: hidden;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
  background-color: var(--bg-secondary);
  color: var(--text-primary);
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--bg-secondary);
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
  background-color: var(--bg-secondary);
}

.center-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}
</style>