<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { GetAllFunctions, GetCategories, GetFunctionsByCategory } from '../../../wailsjs/go/main/App'

interface UIComponent {
  id: string
  name: string
  icon: string
  category: string
  type: string
}

interface ChineseFunction {
  chineseName: string
  goFunction: string
  description: string
  category: string
  params: string[]
  returnType: string
}

interface CodeFile {
  id: string
  name: string
  code: string
}

interface ProjectWindow {
  id: string
  name: string
  type: string
  icon: string
  isActive: boolean
  codeFiles: CodeFile[]
}

const props = defineProps<{
  forceActiveTab?: string
  hideTabs?: boolean
}>()

const activeTab = ref('components')

const effectiveTab = computed(() => props.forceActiveTab || activeTab.value)

const emit = defineEmits<{
  (e: 'select-component', component: UIComponent): void
  (e: 'select-function', fn: ChineseFunction): void
  (e: 'select-window-type', type: string): void
  (e: 'open-window-design', windowId: string): void
  (e: 'open-window-code', windowId: string, codeFileId: string): void
  (e: 'add-new-window', windowData: ProjectWindow): void
  (e: 'delete-window', windowId: string): void
  (e: 'rename-window', windowId: string, newName: string): void
  (e: 'add-code-file', windowId: string, codeFile: CodeFile): void
  (e: 'delete-code-file', windowId: string, codeFileId: string): void
  (e: 'rename-code-file', windowId: string, codeFileId: string, newName: string): void
  (e: 'add-independent-code-file', codeFile: CodeFile): void
  (e: 'delete-independent-code-file', fileId: string): void
  (e: 'rename-independent-code-file', fileId: string, newName: string): void
  (e: 'open-independent-code-file', fileId: string): void
}>()

function defaultCodeFiles(): CodeFile[] {
  return [{ id: 'cf_1', name: 'main.ego', code: '' }]
}

const projectWindows = ref<ProjectWindow[]>([
  { id: 'win_1', name: '启动窗口', type: 'normal', icon: 'Monitor', isActive: true, codeFiles: defaultCodeFiles() },
])

const nextWindowId = ref(2)
const nextCodeFileId = ref(2)
const nextIndependentFileId = ref(1)
const showAddMenu = ref(false)
const independentCodeFiles = ref<CodeFile[]>([])
const editingWindowId = ref<string | null>(null)
const editingCodeFileId = ref<{ windowId: string; fileId: string } | null>(null)
const editName = ref('')
const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const contextMenuTarget = ref<{ type: 'window'; data: any } | null>(null)

const categories = ref([
  {
    name: '🧱 基础',
    components: [
      { id: 'button', name: '按钮', icon: 'Button', category: '基础', type: 'button' },
      { id: 'label', name: '标签', icon: 'Document', category: '基础', type: 'label' },
      { id: 'text', name: '文本', icon: 'ChatLineSquare', category: '基础', type: 'text' },
      { id: 'link', name: '链接', icon: 'Link', category: '基础', type: 'link' },
      { id: 'icon', name: '图标', icon: 'Picture', category: '基础', type: 'icon' },
      { id: 'divider', name: '分割线', icon: 'Minus', category: '基础', type: 'divider' },
      { id: 'space', name: '间距', icon: 'Rank', category: '基础', type: 'space' },
      { id: 'heading', name: '标题', icon: 'EditPen', category: '基础', type: 'heading' },
      { id: 'paragraph', name: '段落', icon: 'Notebook', category: '基础', type: 'paragraph' },
      { id: 'blockquote', name: '引用', icon: 'ChatLineRound', category: '基础', type: 'blockquote' },
      { id: 'code-block', name: '代码块', icon: 'Code', category: '基础', type: 'code-block' },
      { id: 'image', name: '图片', icon: 'PictureFilled', category: '基础', type: 'image' },
      { id: 'video', name: '视频', icon: 'VideoCamera', category: '基础', type: 'video' },
    ] as UIComponent[],
  },
  {
    name: '📝 表单',
    components: [
      { id: 'input', name: '输入框', icon: 'Edit', category: '表单', type: 'input' },
      { id: 'textarea', name: '文本域', icon: 'Notebook', category: '表单', type: 'textarea' },
      { id: 'input-number', name: '数字输入', icon: 'Sort', category: '表单', type: 'input-number' },
      { id: 'checkbox', name: '复选框', icon: 'Select', category: '表单', type: 'checkbox' },
      { id: 'radio', name: '单选框', icon: 'CircleCheck', category: '表单', type: 'radio' },
      { id: 'switch', name: '开关', icon: 'Switch', category: '表单', type: 'switch' },
      { id: 'slider', name: '滑块', icon: 'Operation', category: '表单', type: 'slider' },
      { id: 'select', name: '选择器', icon: 'ArrowDown', category: '表单', type: 'select' },
      { id: 'rate', name: '评分', icon: 'Star', category: '表单', type: 'rate' },
      { id: 'color-picker', name: '颜色选择', icon: 'Brush', category: '表单', type: 'color-picker' },
      { id: 'date-picker', name: '日期选择', icon: 'Calendar', category: '表单', type: 'date-picker' },
      { id: 'time-picker', name: '时间选择', icon: 'Clock', category: '表单', type: 'time-picker' },
      { id: 'upload', name: '上传', icon: 'Upload', category: '表单', type: 'upload' },
      { id: 'cascader', name: '级联选择', icon: 'DArrowRight', category: '表单', type: 'cascader' },
      { id: 'transfer', name: '穿梭框', icon: 'Sort', category: '表单', type: 'transfer' },
      { id: 'autocomplete', name: '自动完成', icon: 'Search', category: '表单', type: 'autocomplete' },
      { id: 'form', name: '表单容器', icon: 'Collection', category: '表单', type: 'form' },
      { id: 'rich-text', name: '富文本', icon: 'DocumentCopy', category: '表单', type: 'rich-text' },
    ] as UIComponent[],
  },
  {
    name: '📊 数据展示',
    components: [
      { id: 'table', name: '表格', icon: 'List', category: '数据展示', type: 'table' },
      { id: 'tag', name: '标签', icon: 'PriceTag', category: '数据展示', type: 'tag' },
      { id: 'badge', name: '徽标', icon: 'Notification', category: '数据展示', type: 'badge' },
      { id: 'progress', name: '进度条', icon: 'Loading', category: '数据展示', type: 'progress' },
      { id: 'avatar', name: '头像', icon: 'User', category: '数据展示', type: 'avatar' },
      { id: 'card', name: '卡片', icon: 'Postcard', category: '数据展示', type: 'card' },
      { id: 'collapse', name: '折叠面板', icon: 'Fold', category: '数据展示', type: 'collapse' },
      { id: 'timeline', name: '时间线', icon: 'Timer', category: '数据展示', type: 'timeline' },
      { id: 'statistic', name: '统计数值', icon: 'DataAnalysis', category: '数据展示', type: 'statistic' },
      { id: 'empty', name: '空状态', icon: 'FolderOpened', category: '数据展示', type: 'empty' },
      { id: 'skeleton', name: '骨架屏', icon: 'SemiSelection', category: '数据展示', type: 'skeleton' },
      { id: 'tree', name: '树形控件', icon: 'Guide', category: '数据展示', type: 'tree' },
      { id: 'descriptions', name: '描述列表', icon: 'Tickets', category: '数据展示', type: 'descriptions' },
      { id: 'calendar', name: '日历', icon: 'Calendar', category: '数据展示', type: 'calendar' },
      { id: 'carousel', name: '走马灯', icon: 'PictureRounded', category: '数据展示', type: 'carousel' },
      { id: 'list', name: '列表', icon: 'List', category: '数据展示', type: 'list' },
      { id: 'image-viewer', name: '图片查看', icon: 'ZoomIn', category: '数据展示', type: 'image-viewer' },
      { id: 'virtual-list', name: '虚拟列表', icon: 'Grid', category: '数据展示', type: 'virtual-list' },
    ] as UIComponent[],
  },
  {
    name: '🧭 导航',
    components: [
      { id: 'menu', name: '菜单', icon: 'Menu', category: '导航', type: 'menu' },
      { id: 'tabs', name: '标签页', icon: 'Tickets', category: '导航', type: 'tabs' },
      { id: 'breadcrumb', name: '面包屑', icon: 'DArrowRight', category: '导航', type: 'breadcrumb' },
      { id: 'dropdown', name: '下拉菜单', icon: 'CaretBottom', category: '导航', type: 'dropdown' },
      { id: 'pagination', name: '分页', icon: 'More', category: '导航', type: 'pagination' },
      { id: 'steps', name: '步骤条', icon: 'SetUp', category: '导航', type: 'steps' },
      { id: 'backtop', name: '回到顶部', icon: 'UploadFilled', category: '导航', type: 'backtop' },
      { id: 'anchor', name: '锚点', icon: 'Aim', category: '导航', type: 'anchor' },
      { id: 'affix', name: '固钉', icon: 'Pushpin', category: '导航', type: 'affix' },
      { id: 'navbar', name: '导航栏', icon: 'HomeFilled', category: '导航', type: 'navbar' },
    ] as UIComponent[],
  },
  {
    name: '💬 反馈',
    components: [
      { id: 'dialog', name: '对话框', icon: 'ChatDotSquare', category: '反馈', type: 'dialog' },
      { id: 'drawer', name: '抽屉', icon: 'DArrowLeft', category: '反馈', type: 'drawer' },
      { id: 'tooltip', name: '文字提示', icon: 'QuestionFilled', category: '反馈', type: 'tooltip' },
      { id: 'popover', name: '气泡卡片', icon: 'ChatRound', category: '反馈', type: 'popover' },
      { id: 'alert', name: '警告', icon: 'Warning', category: '反馈', type: 'alert' },
      { id: 'loading', name: '加载', icon: 'Loading', category: '反馈', type: 'loading' },
      { id: 'result', name: '结果页', icon: 'CircleCheck', category: '反馈', type: 'result' },
      { id: 'message', name: '消息提示', icon: 'ChatLineSquare', category: '反馈', type: 'message' },
      { id: 'notification', name: '通知', icon: 'Bell', category: '反馈', type: 'notification' },
      { id: 'watermark', name: '水印', icon: 'Stamp', category: '反馈', type: 'watermark' },
      { id: 'progress-bar', name: '进度条', icon: 'Loading', category: '反馈', type: 'progress-bar' },
      { id: 'confirm', name: '确认框', icon: 'WarningFilled', category: '反馈', type: 'confirm' },
    ] as UIComponent[],
  },
  {
    name: '📦 容器',
    components: [
      { id: 'panel', name: '面板', icon: 'Grid', category: '容器', type: 'panel' },
      { id: 'group', name: '分组框', icon: 'Collection', category: '容器', type: 'group' },
      { id: 'container', name: '容器', icon: 'Box', category: '容器', type: 'container' },
      { id: 'scrollbar', name: '滚动条', icon: 'Rank', category: '容器', type: 'scrollbar' },
      { id: 'layout', name: '布局', icon: 'Grid', category: '容器', type: 'layout' },
      { id: 'row', name: '行', icon: 'Rank', category: '容器', type: 'row' },
      { id: 'col', name: '列', icon: 'Sort', category: '容器', type: 'col' },
      { id: 'flex', name: '弹性布局', icon: 'Operation', category: '容器', type: 'flex' },
      { id: 'split-pane', name: '分割面板', icon: 'DArrowRight', category: '容器', type: 'split-pane' },
      { id: 'tab-container', name: '选项卡容器', icon: 'Tickets', category: '容器', type: 'tab-container' },
    ] as UIComponent[],
  },
  {
    name: '📈 图表',
    components: [
      { id: 'line-chart', name: '折线图', icon: 'TrendCharts', category: '图表', type: 'line-chart' },
      { id: 'bar-chart', name: '柱状图', icon: 'Histogram', category: '图表', type: 'bar-chart' },
      { id: 'gauge', name: '仪表盘', icon: 'Odometer', category: '图表', type: 'gauge' },
      { id: 'ring-progress', name: '环形进度', icon: 'PieChart', category: '图表', type: 'ring-progress' },
      { id: 'pie-chart', name: '饼图', icon: 'PieChart', category: '图表', type: 'pie-chart' },
      { id: 'scatter-chart', name: '散点图', icon: 'ScatterPlot', category: '图表', type: 'scatter-chart' },
      { id: 'radar-chart', name: '雷达图', icon: 'Radar', category: '图表', type: 'radar-chart' },
      { id: 'area-chart', name: '面积图', icon: 'TrendCharts', category: '图表', type: 'area-chart' },
      { id: 'funnel-chart', name: '漏斗图', icon: 'DataAnalysis', category: '图表', type: 'funnel-chart' },
      { id: 'heatmap', name: '热力图', icon: 'Grid', category: '图表', type: 'heatmap' },
    ] as UIComponent[],
  },
])

const fnCategories = ref<string[]>([])
const allFunctions = ref<ChineseFunction[]>([])
const selectedFnCategory = ref('')
const categoryFunctions = ref<ChineseFunction[]>([])
const expandedFnCategories = ref<Set<string>>(new Set())
const expandedCodeFiles = ref<Set<string>>(new Set())

const solutionExpanded = ref(true)
const groupExpands = ref({
  windows: true,
  shared: true,
  resources: false,
})
const folderExpands = ref({
  src: true,
  build: false,
  int: false,
})
const selectedTreeItem = ref<string | null>(null)
const projectName = ref('未命名项目')

const toggleSolutionExpand = () => {
  solutionExpanded.value = !solutionExpanded.value
}

const toggleGroupExpand = (group: 'windows' | 'shared' | 'resources') => {
  groupExpands.value[group] = !groupExpands.value[group]
}

const toggleFolderExpand = (folder: 'src' | 'build' | 'int') => {
  folderExpands.value[folder] = !folderExpands.value[folder]
}

const handleToolbarCommand = (cmd: string) => {
  if (cmd === 'new-window') {
    addNewWindow()
  } else if (cmd === 'new-code') {
    addIndependentCodeFile()
  }
}

const refreshExplorer = () => {
  solutionExpanded.value = true
  folderExpands.value.src = true
}

const collapseAll = () => {
  solutionExpanded.value = false
  folderExpands.value.src = false
  folderExpands.value.build = false
  folderExpands.value.int = false
}

const windowTypes = [
  { type: 'normal', name: '普通窗口', desc: '标准可调大小窗口', icon: 'Monitor' },
  { type: 'dialog', name: '对话框', desc: '模态对话框', icon: 'ChatDotSquare' },
  { type: 'popup', name: '弹出窗口', desc: '无边框弹出窗口', icon: 'Promotion' },
  { type: 'tool', name: '工具窗口', desc: '小标题栏工具窗口', icon: 'Tools' },
  { type: 'splash', name: '启动窗口', desc: '无边框启动画面', icon: 'PictureFilled' },
]

const expandedCompCategories = ref<Set<string>>(new Set())

const toggleCompCategory = (catName: string) => {
  if (expandedCompCategories.value.has(catName)) {
    expandedCompCategories.value.delete(catName)
  } else {
    expandedCompCategories.value.add(catName)
  }
}

const loadDefaultFunctions = () => {
  const defaultFns: ChineseFunction[] = [
    { chineseName: '信息框', goFunction: 'MessageBox', description: '弹出信息提示框', category: '界面', params: ['提示信息', '标题', '按钮类型'], returnType: '整数型' },
    { chineseName: '调试输出', goFunction: 'DebugPrint', description: '输出调试信息到控制台', category: '调试', params: ['输出内容'], returnType: '' },
    { chineseName: '延迟', goFunction: 'Sleep', description: '延迟指定毫秒数', category: '系统', params: ['毫秒数'], returnType: '' },
    { chineseName: '取文本长度', goFunction: 'Len', description: '取文本的字符长度', category: '文本', params: ['文本'], returnType: '整数型' },
    { chineseName: '到文本', goFunction: 'Str', description: '将数据转换为文本', category: '转换', params: ['数据'], returnType: '文本型' },
    { chineseName: '到整数', goFunction: 'Int', description: '将文本转换为整数', category: '转换', params: ['文本'], returnType: '整数型' },
    { chineseName: '到数值', goFunction: 'Float', description: '将文本转换为小数', category: '转换', params: ['文本'], returnType: '小数型' },
    { chineseName: '到字节集', goFunction: 'Bytes', description: '将文本转换为字节集', category: '转换', params: ['文本'], returnType: '字节集' },
    { chineseName: '取启动时间', goFunction: 'GetTickCount', description: '取系统启动以来的毫秒数', category: '系统', params: [], returnType: '整数型' },
    { chineseName: '取现行时间', goFunction: 'Now', description: '获取当前日期时间', category: '时间', params: [], returnType: '日期时间型' },
    { chineseName: '取年份', goFunction: 'Year', description: '取日期中的年份', category: '时间', params: ['日期时间'], returnType: '整数型' },
    { chineseName: '取月份', goFunction: 'Month', description: '取日期中的月份', category: '时间', params: ['日期时间'], returnType: '整数型' },
    { chineseName: '取日', goFunction: 'Day', description: '取日期中的日', category: '时间', params: ['日期时间'], returnType: '整数型' },
    { chineseName: '取小时', goFunction: 'Hour', description: '取时间中的小时', category: '时间', params: ['日期时间'], returnType: '整数型' },
    { chineseName: '取分钟', goFunction: 'Minute', description: '取时间中的分钟', category: '时间', params: ['日期时间'], returnType: '整数型' },
    { chineseName: '取秒', goFunction: 'Second', description: '取时间中的秒', category: '时间', params: ['日期时间'], returnType: '整数型' },
    { chineseName: '格式化时间', goFunction: 'FormatTime', description: '格式化日期时间为文本', category: '时间', params: ['日期时间', '格式'], returnType: '文本型' },
    { chineseName: '运行', goFunction: 'Exec', description: '执行外部命令或程序', category: '系统', params: ['命令行', '是否等待'], returnType: '整数型' },
    { chineseName: '取命令行', goFunction: 'GetCmdArgs', description: '取程序启动命令行参数', category: '系统', params: [], returnType: '文本型数组' },
    { chineseName: '取环境变量', goFunction: 'GetEnv', description: '取系统环境变量值', category: '系统', params: ['变量名'], returnType: '文本型' },
    { chineseName: '置环境变量', goFunction: 'SetEnv', description: '设置系统环境变量', category: '系统', params: ['变量名', '值'], returnType: '逻辑型' },
    { chineseName: '取当前目录', goFunction: 'GetCurrentDir', description: '获取程序当前工作目录', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '置当前目录', goFunction: 'SetCurrentDir', description: '设置程序当前工作目录', category: '系统', params: ['目录路径'], returnType: '逻辑型' },
    { chineseName: '取临时目录', goFunction: 'GetTempDir', description: '获取系统临时目录路径', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取计算机名', goFunction: 'GetComputerName', description: '获取本机计算机名', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取用户名', goFunction: 'GetUserName', description: '获取当前登录用户名', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '读入文件', goFunction: 'ReadFile', description: '读取文件全部内容', category: '文件', params: ['文件路径'], returnType: '字节集' },
    { chineseName: '写到文件', goFunction: 'WriteFile', description: '写入数据到文件', category: '文件', params: ['文件路径', '数据'], returnType: '逻辑型' },
    { chineseName: '创建目录', goFunction: 'Mkdir', description: '创建目录（支持多级）', category: '文件', params: ['目录路径'], returnType: '逻辑型' },
    { chineseName: '删除文件', goFunction: 'RemoveFile', description: '删除指定文件', category: '文件', params: ['文件路径'], returnType: '逻辑型' },
    { chineseName: '删除目录', goFunction: 'RemoveDir', description: '删除指定目录', category: '文件', params: ['目录路径'], returnType: '逻辑型' },
    { chineseName: '文件是否存在', goFunction: 'FileExists', description: '判断文件是否存在', category: '文件', params: ['文件路径'], returnType: '逻辑型' },
    { chineseName: '目录是否存在', goFunction: 'DirExists', description: '判断目录是否存在', category: '文件', params: ['目录路径'], returnType: '逻辑型' },
    { chineseName: '取文件大小', goFunction: 'GetFileSize', description: '获取文件大小（字节）', category: '文件', params: ['文件路径'], returnType: '整数型' },
    { chineseName: '取文件时间', goFunction: 'GetFileTime', description: '获取文件修改时间', category: '文件', params: ['文件路径'], returnType: '日期时间型' },
    { chineseName: '复制文件', goFunction: 'CopyFile', description: '复制文件到目标位置', category: '文件', params: ['源文件', '目标文件'], returnType: '逻辑型' },
    { chineseName: '移动文件', goFunction: 'MoveFile', description: '移动文件到目标位置', category: '文件', params: ['源文件', '目标文件'], returnType: '逻辑型' },
    { chineseName: '取文件扩展名', goFunction: 'GetFileExt', description: '取文件的扩展名', category: '文件', params: ['文件路径'], returnType: '文本型' },
    { chineseName: '取文件名', goFunction: 'GetFileName', description: '取路径中的文件名', category: '文件', params: ['文件路径'], returnType: '文本型' },
    { chineseName: '取目录名', goFunction: 'GetDirName', description: '取路径中的目录名', category: '文件', params: ['文件路径'], returnType: '文本型' },
    { chineseName: '文本替换', goFunction: 'Replace', description: '替换文本中的指定内容', category: '文本', params: ['原文本', '被替换文本', '替换为'], returnType: '文本型' },
    { chineseName: '取文本中间', goFunction: 'Mid', description: '取文本中间指定部分', category: '文本', params: ['文本', '起始位置', '长度'], returnType: '文本型' },
    { chineseName: '取文本左边', goFunction: 'Left', description: '取文本左边指定长度', category: '文本', params: ['文本', '长度'], returnType: '文本型' },
    { chineseName: '取文本右边', goFunction: 'Right', description: '取文本右边指定长度', category: '文本', params: ['文本', '长度'], returnType: '文本型' },
    { chineseName: '分割文本', goFunction: 'Split', description: '按分隔符分割文本为数组', category: '文本', params: ['文本', '分隔符'], returnType: '文本型数组' },
    { chineseName: '到大写', goFunction: 'ToUpper', description: '将文本转换为大写', category: '文本', params: ['文本'], returnType: '文本型' },
    { chineseName: '到小写', goFunction: 'ToLower', description: '将文本转换为小写', category: '文本', params: ['文本'], returnType: '文本型' },
    { chineseName: '删首尾空', goFunction: 'Trim', description: '删除文本首尾空格', category: '文本', params: ['文本'], returnType: '文本型' },
    { chineseName: '删首空', goFunction: 'TrimLeft', description: '删除文本首部空格', category: '文本', params: ['文本'], returnType: '文本型' },
    { chineseName: '删尾空', goFunction: 'TrimRight', description: '删除文本尾部空格', category: '文本', params: ['文本'], returnType: '文本型' },
    { chineseName: '寻找文本', goFunction: 'FindStr', description: '在文本中寻找指定内容的位置', category: '文本', params: ['文本', '寻找内容', '起始位置'], returnType: '整数型' },
    { chineseName: '倒找文本', goFunction: 'FindStrRev', description: '从尾部寻找指定内容的位置', category: '文本', params: ['文本', '寻找内容'], returnType: '整数型' },
    { chineseName: '取文本出现次数', goFunction: 'CountStr', description: '统计子文本出现次数', category: '文本', params: ['文本', '子文本'], returnType: '整数型' },
    { chineseName: '取随机数', goFunction: 'Random', description: '取指定范围内的随机整数', category: '数学', params: ['最小值', '最大值'], returnType: '整数型' },
    { chineseName: '四舍五入', goFunction: 'Round', description: '对数值进行四舍五入', category: '数学', params: ['数值', '小数位数'], returnType: '小数型' },
    { chineseName: '取绝对值', goFunction: 'Abs', description: '取数值的绝对值', category: '数学', params: ['数值'], returnType: '小数型' },
    { chineseName: '取整', goFunction: 'Floor', description: '向下取整', category: '数学', params: ['数值'], returnType: '整数型' },
    { chineseName: '求次方', goFunction: 'Pow', description: '求数值的次方', category: '数学', params: ['底数', '指数'], returnType: '小数型' },
    { chineseName: '求平方根', goFunction: 'Sqrt', description: '求数值的平方根', category: '数学', params: ['数值'], returnType: '小数型' },
    { chineseName: '求正弦', goFunction: 'Sin', description: '求角度的正弦值', category: '数学', params: ['角度'], returnType: '小数型' },
    { chineseName: '求余弦', goFunction: 'Cos', description: '求角度的余弦值', category: '数学', params: ['角度'], returnType: '小数型' },
    { chineseName: '求正切', goFunction: 'Tan', description: '求角度的正切值', category: '数学', params: ['角度'], returnType: '小数型' },
    { chineseName: '取最大值', goFunction: 'Max', description: '取两个数值中的最大值', category: '数学', params: ['数值1', '数值2'], returnType: '小数型' },
    { chineseName: '取最小值', goFunction: 'Min', description: '取两个数值中的最小值', category: '数学', params: ['数值1', '数值2'], returnType: '小数型' },
    { chineseName: 'HTTP请求', goFunction: 'HTTPRequest', description: '发送HTTP请求', category: '网络', params: ['URL', '方法', '数据'], returnType: '文本型' },
    { chineseName: 'HTTP下载', goFunction: 'HTTPDownload', description: '下载网络文件到本地', category: '网络', params: ['URL', '保存路径'], returnType: '逻辑型' },
    { chineseName: 'JSON解析', goFunction: 'JSONParse', description: '解析JSON文本', category: '网络', params: ['JSON文本'], returnType: '对象' },
    { chineseName: 'JSON生成', goFunction: 'JSONStringify', description: '生成JSON文本', category: '网络', params: ['对象'], returnType: '文本型' },
    { chineseName: 'Base64编码', goFunction: 'Base64Encode', description: '将数据编码为Base64', category: '编码', params: ['数据'], returnType: '文本型' },
    { chineseName: 'Base64解码', goFunction: 'Base64Decode', description: '将Base64文本解码', category: '编码', params: ['Base64文本'], returnType: '字节集' },
    { chineseName: 'MD5加密', goFunction: 'MD5', description: '计算文本的MD5值', category: '编码', params: ['文本'], returnType: '文本型' },
    { chineseName: 'SHA256加密', goFunction: 'SHA256', description: '计算文本的SHA256值', category: '编码', params: ['文本'], returnType: '文本型' },
    { chineseName: 'URL编码', goFunction: 'URLEncode', description: '对文本进行URL编码', category: '编码', params: ['文本'], returnType: '文本型' },
    { chineseName: 'URL解码', goFunction: 'URLDecode', description: '对URL编码文本进行解码', category: '编码', params: ['文本'], returnType: '文本型' },
    { chineseName: '创建窗口', goFunction: 'CreateForm', description: '创建并显示窗口', category: '窗口', params: ['窗口模板'], returnType: '' },
    { chineseName: '载入窗口', goFunction: 'LoadForm', description: '载入窗口（模态/非模态）', category: '窗口', params: ['窗口变量', '是否模态'], returnType: '' },
    { chineseName: '销毁窗口', goFunction: 'DestroyForm', description: '销毁当前窗口', category: '窗口', params: [], returnType: '' },
    { chineseName: '创建按钮', goFunction: 'NewButton', description: '创建按钮组件', category: '窗口', params: ['父窗口'], returnType: '按钮' },
    { chineseName: '创建标签', goFunction: 'NewLabel', description: '创建标签组件', category: '窗口', params: ['父窗口'], returnType: '标签' },
    { chineseName: '创建编辑框', goFunction: 'NewEdit', description: '创建编辑框组件', category: '窗口', params: ['父窗口'], returnType: '编辑框' },
    { chineseName: '创建选择框', goFunction: 'NewCheckBox', description: '创建选择框组件', category: '窗口', params: ['父窗口'], returnType: '选择框' },
    { chineseName: '创建单选框', goFunction: 'NewRadioButton', description: '创建单选框组件', category: '窗口', params: ['父窗口'], returnType: '单选框' },
    { chineseName: '创建组合框', goFunction: 'NewComboBox', description: '创建组合框组件', category: '窗口', params: ['父窗口'], returnType: '组合框' },
    { chineseName: '创建列表框', goFunction: 'NewListBox', description: '创建列表框组件', category: '窗口', params: ['父窗口'], returnType: '列表框' },
    { chineseName: '创建图片框', goFunction: 'NewImage', description: '创建图片框组件', category: '窗口', params: ['父窗口'], returnType: '图片框' },
    { chineseName: '创建面板', goFunction: 'NewPanel', description: '创建面板组件', category: '窗口', params: ['父窗口'], returnType: '面板' },
    { chineseName: '创建进度条', goFunction: 'NewProgressBar', description: '创建进度条组件', category: '窗口', params: ['父窗口'], returnType: '进度条' },
    { chineseName: '创建滑块', goFunction: 'NewTrackBar', description: '创建滑块组件', category: '窗口', params: ['父窗口'], returnType: '滑块' },
    { chineseName: '创建表格', goFunction: 'NewStringGrid', description: '创建表格组件', category: '窗口', params: ['父窗口'], returnType: '表格' },
    { chineseName: '创建树形框', goFunction: 'NewTreeView', description: '创建树形框组件', category: '窗口', params: ['父窗口'], returnType: '树形框' },
    { chineseName: '创建选项卡', goFunction: 'NewPageControl', description: '创建选项卡组件', category: '窗口', params: ['父窗口'], returnType: '选项卡' },
    { chineseName: '创建菜单', goFunction: 'NewMainMenu', description: '创建菜单组件', category: '窗口', params: ['父窗口'], returnType: '菜单' },
    { chineseName: '创建工具栏', goFunction: 'NewToolBar', description: '创建工具栏组件', category: '窗口', params: ['父窗口'], returnType: '工具栏' },
    { chineseName: '创建状态栏', goFunction: 'NewStatusBar', description: '创建状态栏组件', category: '窗口', params: ['父窗口'], returnType: '状态栏' },
    { chineseName: '创建定时器', goFunction: 'NewTimer', description: '创建定时器组件', category: '窗口', params: ['父窗口', '间隔毫秒'], returnType: '定时器' },
    { chineseName: '创建浏览器', goFunction: 'NewMiniBlink', description: '创建浏览器组件', category: '窗口', params: ['父窗口'], returnType: '浏览器' },
    { chineseName: '创建分组框', goFunction: 'NewGroupBox', description: '创建分组框组件', category: '窗口', params: ['父窗口'], returnType: '分组框' },
    { chineseName: '创建滚动框', goFunction: 'NewScrollBox', description: '创建滚动框组件', category: '窗口', params: ['父窗口'], returnType: '滚动框' },
    { chineseName: '创建日期框', goFunction: 'NewDateTimePicker', description: '创建日期选择框', category: '窗口', params: ['父窗口'], returnType: '日期框' },
    { chineseName: '创建颜色框', goFunction: 'NewColorBox', description: '创建颜色选择框', category: '窗口', params: ['父窗口'], returnType: '颜色框' },
    { chineseName: '创建月历', goFunction: 'NewMonthCalendar', description: '创建月历组件', category: '窗口', params: ['父窗口'], returnType: '月历' },
    { chineseName: '创建外形框', goFunction: 'NewShape', description: '创建外形框组件', category: '窗口', params: ['父窗口'], returnType: '外形框' },
    { chineseName: '创建分隔条', goFunction: 'NewSplitter', description: '创建分隔条组件', category: '窗口', params: ['父窗口'], returnType: '分隔条' },
    { chineseName: '创建滚动条', goFunction: 'NewScrollBar', description: '创建滚动条组件', category: '窗口', params: ['父窗口'], returnType: '滚动条' },
    { chineseName: '创建调节器', goFunction: 'NewSpinEdit', description: '创建调节器组件', category: '窗口', params: ['父窗口'], returnType: '调节器' },
    { chineseName: '创建多行编辑框', goFunction: 'NewMemo', description: '创建多行编辑框', category: '窗口', params: ['父窗口'], returnType: '多行编辑框' },
    { chineseName: '创建斜角框', goFunction: 'NewBevel', description: '创建斜角框组件', category: '窗口', params: ['父窗口'], returnType: '斜角框' },
    { chineseName: '创建打开对话框', goFunction: 'NewOpenDialog', description: '创建打开文件对话框', category: '窗口', params: ['父窗口'], returnType: '对话框' },
    { chineseName: '创建保存对话框', goFunction: 'NewSaveDialog', description: '创建保存文件对话框', category: '窗口', params: ['父窗口'], returnType: '对话框' },
    { chineseName: '置组件标题', goFunction: 'SetCaption', description: '设置组件的标题/文本', category: '窗口', params: ['组件', '标题'], returnType: '' },
    { chineseName: '取组件标题', goFunction: 'GetCaption', description: '获取组件的标题/文本', category: '窗口', params: ['组件'], returnType: '文本型' },
    { chineseName: '置组件边界', goFunction: 'SetBounds', description: '设置组件位置和大小', category: '窗口', params: ['组件', '左', '顶', '宽', '高'], returnType: '' },
    { chineseName: '置组件可视', goFunction: 'SetVisible', description: '设置组件是否可见', category: '窗口', params: ['组件', '是否可见'], returnType: '' },
    { chineseName: '置组件可用', goFunction: 'SetEnabled', description: '设置组件是否可用', category: '窗口', params: ['组件', '是否可用'], returnType: '' },
    { chineseName: '置组件颜色', goFunction: 'SetColor', description: '设置组件背景颜色', category: '窗口', params: ['组件', '颜色值'], returnType: '' },
    { chineseName: '置组件字体', goFunction: 'SetFont', description: '设置组件字体', category: '窗口', params: ['组件', '字体名', '大小', '是否加粗'], returnType: '' },
    { chineseName: '置组件提示', goFunction: 'SetHint', description: '设置组件鼠标悬停提示', category: '窗口', params: ['组件', '提示文本'], returnType: '' },
    { chineseName: '置组件光标', goFunction: 'SetCursor', description: '设置组件鼠标光标样式', category: '窗口', params: ['组件', '光标类型'], returnType: '' },
    { chineseName: '置组件对齐', goFunction: 'SetAlign', description: '设置组件对齐方式', category: '窗口', params: ['组件', '对齐方式'], returnType: '' },
    { chineseName: '置组件锚点', goFunction: 'SetAnchors', description: '设置组件锚点', category: '窗口', params: ['组件', '左锚', '顶锚', '右锚', '底锚'], returnType: '' },
    { chineseName: '置窗口标题', goFunction: 'SetFormCaption', description: '设置窗口标题', category: '窗口', params: ['窗口', '标题'], returnType: '' },
    { chineseName: '置窗口大小', goFunction: 'SetFormSize', description: '设置窗口大小', category: '窗口', params: ['窗口', '宽度', '高度'], returnType: '' },
    { chineseName: '置窗口图标', goFunction: 'SetFormIcon', description: '设置窗口图标', category: '窗口', params: ['窗口', '图标路径'], returnType: '' },
    { chineseName: '置窗口背景', goFunction: 'SetFormColor', description: '设置窗口背景颜色', category: '窗口', params: ['窗口', '颜色值'], returnType: '' },
    { chineseName: '置窗口居中', goFunction: 'SetFormCenter', description: '将窗口置于屏幕中央', category: '窗口', params: ['窗口'], returnType: '' },
    { chineseName: '置窗口置顶', goFunction: 'SetFormTopMost', description: '设置窗口是否置顶', category: '窗口', params: ['窗口', '是否置顶'], returnType: '' },
    { chineseName: '置窗口边框', goFunction: 'SetFormBorder', description: '设置窗口边框样式', category: '窗口', params: ['窗口', '边框样式'], returnType: '' },
    { chineseName: '置窗口透明度', goFunction: 'SetFormAlpha', description: '设置窗口透明度(0-255)', category: '窗口', params: ['窗口', '透明度'], returnType: '' },
    { chineseName: '取屏幕宽度', goFunction: 'GetScreenWidth', description: '获取屏幕宽度', category: '窗口', params: [], returnType: '整数型' },
    { chineseName: '取屏幕高度', goFunction: 'GetScreenHeight', description: '获取屏幕高度', category: '窗口', params: [], returnType: '整数型' },
    { chineseName: '取剪贴板文本', goFunction: 'GetClipboard', description: '获取剪贴板文本内容', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '置剪贴板文本', goFunction: 'SetClipboard', description: '设置剪贴板文本内容', category: '系统', params: ['文本'], returnType: '' },
    { chineseName: '播放声音', goFunction: 'PlaySound', description: '播放系统声音或音频文件', category: '系统', params: ['声音文件', '是否循环'], returnType: '逻辑型' },
    { chineseName: '停止播放', goFunction: 'StopSound', description: '停止播放声音', category: '系统', params: [], returnType: '' },
    { chineseName: '取进程列表', goFunction: 'GetProcessList', description: '获取系统进程列表', category: '系统', params: [], returnType: '文本型数组' },
    { chineseName: '结束进程', goFunction: 'KillProcess', description: '结束指定进程', category: '系统', params: ['进程名或PID'], returnType: '逻辑型' },
    { chineseName: '注册表读', goFunction: 'RegRead', description: '读取注册表值', category: '系统', params: ['路径', '键名'], returnType: '文本型' },
    { chineseName: '注册表写', goFunction: 'RegWrite', description: '写入注册表值', category: '系统', params: ['路径', '键名', '值'], returnType: '逻辑型' },
    { chineseName: '注册表删除', goFunction: 'RegDelete', description: '删除注册表键', category: '系统', params: ['路径', '键名'], returnType: '逻辑型' },
    { chineseName: '取CPU使用率', goFunction: 'GetCPUUsage', description: '获取CPU使用率百分比', category: '系统', params: [], returnType: '整数型' },
    { chineseName: '取内存使用率', goFunction: 'GetMemUsage', description: '获取内存使用率百分比', category: '系统', params: [], returnType: '整数型' },
    { chineseName: '取磁盘剩余空间', goFunction: 'GetDiskFree', description: '获取磁盘剩余空间(MB)', category: '系统', params: ['盘符'], returnType: '整数型' },
    { chineseName: '取磁盘总空间', goFunction: 'GetDiskTotal', description: '获取磁盘总空间(MB)', category: '系统', params: ['盘符'], returnType: '整数型' },
    { chineseName: '打开网址', goFunction: 'OpenURL', description: '用默认浏览器打开网址', category: '系统', params: ['网址'], returnType: '' },
    { chineseName: '打开文件', goFunction: 'OpenFile', description: '用默认程序打开文件', category: '系统', params: ['文件路径'], returnType: '' },
    { chineseName: '打开目录', goFunction: 'OpenDir', description: '用资源管理器打开目录', category: '系统', params: ['目录路径'], returnType: '' },
    { chineseName: '取配置项', goFunction: 'GetConfig', description: '读取INI配置项', category: '文件', params: ['文件名', '节名', '键名'], returnType: '文本型' },
    { chineseName: '写配置项', goFunction: 'SetConfig', description: '写入INI配置项', category: '文件', params: ['文件名', '节名', '键名', '值'], returnType: '逻辑型' },
    { chineseName: '取配置节名', goFunction: 'GetConfigSections', description: '读取INI所有节名', category: '文件', params: ['文件名'], returnType: '文本型数组' },
    { chineseName: '取配置键名', goFunction: 'GetConfigKeys', description: '读取INI节下所有键名', category: '文件', params: ['文件名', '节名'], returnType: '文本型数组' },
    { chineseName: 'SQLite打开', goFunction: 'SQLiteOpen', description: '打开SQLite数据库', category: '数据库', params: ['数据库路径'], returnType: '数据库句柄' },
    { chineseName: 'SQLite关闭', goFunction: 'SQLiteClose', description: '关闭SQLite数据库', category: '数据库', params: ['数据库句柄'], returnType: '' },
    { chineseName: 'SQLite执行', goFunction: 'SQLiteExec', description: '执行SQL语句', category: '数据库', params: ['数据库句柄', 'SQL语句'], returnType: '逻辑型' },
    { chineseName: 'SQLite查询', goFunction: 'SQLiteQuery', description: '查询SQLite数据', category: '数据库', params: ['数据库句柄', 'SQL语句'], returnType: '查询结果' },
    { chineseName: '取数组成员数', goFunction: 'ArrayLen', description: '取数组的成员数量', category: '数组', params: ['数组'], returnType: '整数型' },
    { chineseName: '加入成员', goFunction: 'ArrayAppend', description: '向数组末尾添加成员', category: '数组', params: ['数组', '成员'], returnType: '' },
    { chineseName: '删除成员', goFunction: 'ArrayRemove', description: '删除数组指定位置成员', category: '数组', params: ['数组', '位置'], returnType: '' },
    { chineseName: '插入成员', goFunction: 'ArrayInsert', description: '在数组指定位置插入成员', category: '数组', params: ['数组', '位置', '成员'], returnType: '' },
    { chineseName: '清除数组', goFunction: 'ArrayClear', description: '清除数组所有成员', category: '数组', params: ['数组'], returnType: '' },
    { chineseName: '数组排序', goFunction: 'ArraySort', description: '对数组进行排序', category: '数组', params: ['数组', '是否升序'], returnType: '' },
    { chineseName: '数组查找', goFunction: 'ArrayFind', description: '在数组中查找成员位置', category: '数组', params: ['数组', '查找值'], returnType: '整数型' },
    { chineseName: '加密数据', goFunction: 'EncryptData', description: '使用AES加密数据', category: '编码', params: ['数据', '密钥'], returnType: '字节集' },
    { chineseName: '解密数据', goFunction: 'DecryptData', description: '使用AES解密数据', category: '编码', params: ['数据', '密钥'], returnType: '字节集' },
    { chineseName: '取随机字符串', goFunction: 'RandomStr', description: '生成随机字符串', category: '文本', params: ['长度'], returnType: '文本型' },
    { chineseName: '取文本相似度', goFunction: 'StrSimilarity', description: '计算两段文本的相似度', category: '文本', params: ['文本1', '文本2'], returnType: '小数型' },
    { chineseName: '正则匹配', goFunction: 'RegexMatch', description: '使用正则表达式匹配文本', category: '文本', params: ['文本', '正则表达式'], returnType: '文本型数组' },
    { chineseName: '正则替换', goFunction: 'RegexReplace', description: '使用正则表达式替换文本', category: '文本', params: ['文本', '正则表达式', '替换为'], returnType: '文本型' },
    { chineseName: '取IP地址', goFunction: 'GetLocalIP', description: '获取本机IP地址', category: '网络', params: [], returnType: '文本型' },
    { chineseName: 'Ping主机', goFunction: 'PingHost', description: 'Ping远程主机', category: '网络', params: ['主机地址', '超时毫秒'], returnType: '逻辑型' },
    { chineseName: 'TCP连接', goFunction: 'TCPConnect', description: '建立TCP客户端连接', category: '网络', params: ['主机', '端口'], returnType: '连接句柄' },
    { chineseName: 'TCP发送', goFunction: 'TCPSend', description: '通过TCP发送数据', category: '网络', params: ['连接句柄', '数据'], returnType: '逻辑型' },
    { chineseName: 'TCP接收', goFunction: 'TCPRecv', description: '通过TCP接收数据', category: '网络', params: ['连接句柄', '超时毫秒'], returnType: '字节集' },
    { chineseName: 'TCP关闭', goFunction: 'TCPClose', description: '关闭TCP连接', category: '网络', params: ['连接句柄'], returnType: '' },
    { chineseName: '创建线程', goFunction: 'CreateThread', description: '创建并启动新线程', category: '系统', params: ['线程函数', '参数'], returnType: '线程句柄' },
    { chineseName: '等待线程', goFunction: 'WaitThread', description: '等待线程执行完毕', category: '系统', params: ['线程句柄', '超时毫秒'], returnType: '逻辑型' },
    { chineseName: '创建互斥体', goFunction: 'CreateMutex', description: '创建线程互斥体', category: '系统', params: ['名称'], returnType: '互斥体句柄' },
    { chineseName: '锁定互斥体', goFunction: 'LockMutex', description: '锁定互斥体', category: '系统', params: ['互斥体句柄'], returnType: '' },
    { chineseName: '解锁互斥体', goFunction: 'UnlockMutex', description: '解锁互斥体', category: '系统', params: ['互斥体句柄'], returnType: '' },
    { chineseName: '取错误信息', goFunction: 'GetLastError', description: '获取最后一次错误信息', category: '调试', params: [], returnType: '文本型' },
    { chineseName: '置随机数种子', goFunction: 'SetRandomSeed', description: '设置随机数种子', category: '数学', params: ['种子值'], returnType: '' },
    { chineseName: '取时间戳', goFunction: 'GetTimestamp', description: '获取Unix时间戳(秒)', category: '时间', params: [], returnType: '整数型' },
    { chineseName: '时间戳转时间', goFunction: 'TimestampToTime', description: '将Unix时间戳转为日期时间', category: '时间', params: ['时间戳'], returnType: '日期时间型' },
    { chineseName: '取时间间隔', goFunction: 'TimeDiff', description: '计算两个时间的间隔秒数', category: '时间', params: ['时间1', '时间2'], returnType: '整数型' },
    { chineseName: '置窗口位置', goFunction: 'SetFormPosition', description: '设置窗口位置', category: '窗口', params: ['窗口', '左边', '顶边'], returnType: '' },
    { chineseName: '取窗口宽度', goFunction: 'GetFormWidth', description: '获取窗口宽度', category: '窗口', params: ['窗口'], returnType: '整数型' },
    { chineseName: '取窗口高度', goFunction: 'GetFormHeight', description: '获取窗口高度', category: '窗口', params: ['窗口'], returnType: '整数型' },
    { chineseName: '最小化窗口', goFunction: 'MinimizeForm', description: '最小化窗口', category: '窗口', params: ['窗口'], returnType: '' },
    { chineseName: '最大化窗口', goFunction: 'MaximizeForm', description: '最大化窗口', category: '窗口', params: ['窗口'], returnType: '' },
    { chineseName: '还原窗口', goFunction: 'RestoreForm', description: '还原窗口大小', category: '窗口', params: ['窗口'], returnType: '' },
    { chineseName: '关闭窗口', goFunction: 'CloseForm', description: '关闭窗口', category: '窗口', params: ['窗口'], returnType: '' },
    { chineseName: '刷新窗口', goFunction: 'RefreshForm', description: '刷新窗口显示', category: '窗口', params: ['窗口'], returnType: '' },
    { chineseName: '置组件事件', goFunction: 'SetEvent', description: '设置组件事件回调', category: '窗口', params: ['组件', '事件类型', '回调函数'], returnType: '' },
    { chineseName: '取组件宽度', goFunction: 'GetWidth', description: '获取组件宽度', category: '窗口', params: ['组件'], returnType: '整数型' },
    { chineseName: '取组件高度', goFunction: 'GetHeight', description: '获取组件高度', category: '窗口', params: ['组件'], returnType: '整数型' },
    { chineseName: '取组件左边', goFunction: 'GetLeft', description: '获取组件左边位置', category: '窗口', params: ['组件'], returnType: '整数型' },
    { chineseName: '取组件顶边', goFunction: 'GetTop', description: '获取组件顶边位置', category: '窗口', params: ['组件'], returnType: '整数型' },
    { chineseName: '置组件Tab顺序', goFunction: 'SetTabOrder', description: '设置组件Tab键顺序', category: '窗口', params: ['组件', '顺序'], returnType: '' },
    { chineseName: '组件获得焦点', goFunction: 'SetFocus', description: '让组件获得输入焦点', category: '窗口', params: ['组件'], returnType: '' },
    { chineseName: '取焦点组件', goFunction: 'GetFocused', description: '获取当前焦点所在组件', category: '窗口', params: ['窗口'], returnType: '组件' },
    { chineseName: '置组件文本颜色', goFunction: 'SetFontColor', description: '设置组件文本颜色', category: '窗口', params: ['组件', '颜色值'], returnType: '' },
    { chineseName: '置组件背景颜色', goFunction: 'SetBgColor', description: '设置组件背景颜色', category: '窗口', params: ['组件', '颜色值'], returnType: '' },
    { chineseName: '置组件透明', goFunction: 'SetTransparent', description: '设置组件是否透明', category: '窗口', params: ['组件', '是否透明'], returnType: '' },
    { chineseName: '置组件边框', goFunction: 'SetBorder', description: '设置组件边框样式', category: '窗口', params: ['组件', '边框样式'], returnType: '' },
    { chineseName: '置组件圆角', goFunction: 'SetRoundRect', description: '设置组件圆角', category: '窗口', params: ['组件', '圆角半径'], returnType: '' },
    { chineseName: '置组件图片', goFunction: 'SetPicture', description: '设置组件显示的图片', category: '窗口', params: ['组件', '图片路径'], returnType: '' },
    { chineseName: '置组件禁用', goFunction: 'DisableComponent', description: '禁用组件', category: '窗口', params: ['组件'], returnType: '' },
    { chineseName: '置组件启用', goFunction: 'EnableComponent', description: '启用组件', category: '窗口', params: ['组件'], returnType: '' },
    { chineseName: '置组件隐藏', goFunction: 'HideComponent', description: '隐藏组件', category: '窗口', params: ['组件'], returnType: '' },
    { chineseName: '置组件显示', goFunction: 'ShowComponent', description: '显示组件', category: '窗口', params: ['组件'], returnType: '' },
    { chineseName: '置组件刷新', goFunction: 'RefreshComponent', description: '刷新组件显示', category: '窗口', params: ['组件'], returnType: '' },
    { chineseName: '置组件置顶', goFunction: 'BringToFront', description: '将组件置于顶层', category: '窗口', params: ['组件'], returnType: '' },
    { chineseName: '置组件置底', goFunction: 'SendToBack', description: '将组件置于底层', category: '窗口', params: ['组件'], returnType: '' },
    { chineseName: '取组件数量', goFunction: 'GetComponentCount', description: '获取窗口中的组件数量', category: '窗口', params: ['窗口'], returnType: '整数型' },
    { chineseName: '取所有组件', goFunction: 'GetAllComponents', description: '获取窗口中所有组件', category: '窗口', params: ['窗口'], returnType: '组件数组' },
    { chineseName: '取组件类型', goFunction: 'GetComponentType', description: '获取组件的类型名称', category: '窗口', params: ['组件'], returnType: '文本型' },
    { chineseName: '取组件名称', goFunction: 'GetComponentName', description: '获取组件的名称', category: '窗口', params: ['组件'], returnType: '文本型' },
    { chineseName: '置组件名称', goFunction: 'SetComponentName', description: '设置组件的名称', category: '窗口', params: ['组件', '名称'], returnType: '' },
    { chineseName: '取组件标记', goFunction: 'GetTag', description: '获取组件的标记值', category: '窗口', params: ['组件'], returnType: '整数型' },
    { chineseName: '置组件标记', goFunction: 'SetTag', description: '设置组件的标记值', category: '窗口', params: ['组件', '标记值'], returnType: '' },
    { chineseName: '取鼠标X', goFunction: 'GetMouseX', description: '获取鼠标X坐标', category: '窗口', params: [], returnType: '整数型' },
    { chineseName: '取鼠标Y', goFunction: 'GetMouseY', description: '获取鼠标Y坐标', category: '窗口', params: [], returnType: '整数型' },
    { chineseName: '模拟按键', goFunction: 'SendKeys', description: '模拟键盘按键', category: '系统', params: ['按键码', '是否按下'], returnType: '' },
    { chineseName: '模拟鼠标', goFunction: 'SendMouse', description: '模拟鼠标操作', category: '系统', params: ['X坐标', 'Y坐标', '操作类型'], returnType: '' },
    { chineseName: '取按键状态', goFunction: 'GetKeyState', description: '获取按键是否按下', category: '系统', params: ['按键码'], returnType: '逻辑型' },
    { chineseName: '取系统语言', goFunction: 'GetSystemLang', description: '获取系统语言', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取系统版本', goFunction: 'GetOSVersion', description: '获取操作系统版本', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '是否64位系统', goFunction: 'Is64BitOS', description: '判断是否为64位系统', category: '系统', params: [], returnType: '逻辑型' },
    { chineseName: '取程序路径', goFunction: 'GetAppPath', description: '获取当前程序完整路径', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取程序目录', goFunction: 'GetAppDir', description: '获取当前程序所在目录', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取桌面目录', goFunction: 'GetDesktopDir', description: '获取桌面目录路径', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取我的文档目录', goFunction: 'GetMyDocDir', description: '获取我的文档目录路径', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取系统目录', goFunction: 'GetSystemDir', description: '获取Windows系统目录', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取Windows目录', goFunction: 'GetWindowsDir', description: '获取Windows安装目录', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取特殊目录', goFunction: 'GetSpecialDir', description: '获取指定特殊目录', category: '系统', params: ['目录类型'], returnType: '文本型' },
    { chineseName: '文件枚举', goFunction: 'EnumFiles', description: '枚举目录中的文件', category: '文件', params: ['目录', '通配符', '是否包含子目录'], returnType: '文本型数组' },
    { chineseName: '目录枚举', goFunction: 'EnumDirs', description: '枚举目录中的子目录', category: '文件', params: ['目录'], returnType: '文本型数组' },
    { chineseName: '驱动器枚举', goFunction: 'EnumDrives', description: '枚举所有驱动器', category: '文件', params: [], returnType: '文本型数组' },
    { chineseName: '取驱动器类型', goFunction: 'GetDriveType', description: '获取驱动器类型', category: '文件', params: ['盘符'], returnType: '文本型' },
    { chineseName: '取文件版本', goFunction: 'GetFileVersion', description: '获取可执行文件版本信息', category: '文件', params: ['文件路径'], returnType: '文本型' },
    { chineseName: '取文件MD5', goFunction: 'GetFileMD5', description: '计算文件的MD5值', category: '文件', params: ['文件路径'], returnType: '文本型' },
    { chineseName: '取文件SHA1', goFunction: 'GetFileSHA1', description: '计算文件的SHA1值', category: '文件', params: ['文件路径'], returnType: '文本型' },
    { chineseName: '压缩文件', goFunction: 'ZipCompress', description: '压缩文件或目录', category: '文件', params: ['源路径', '目标ZIP路径'], returnType: '逻辑型' },
    { chineseName: '解压文件', goFunction: 'ZipExtract', description: '解压ZIP文件', category: '文件', params: ['ZIP路径', '目标目录'], returnType: '逻辑型' },
    { chineseName: '取文本编码', goFunction: 'DetectEncoding', description: '检测文本的编码格式', category: '文本', params: ['字节集'], returnType: '文本型' },
    { chineseName: '编码转换', goFunction: 'ConvertEncoding', description: '转换文本编码', category: '文本', params: ['文本', '源编码', '目标编码'], returnType: '文本型' },
    { chineseName: '取拼音首字母', goFunction: 'GetPinyin', description: '取中文文本的拼音首字母', category: '文本', params: ['中文文本'], returnType: '文本型' },
    { chineseName: '是否为汉字', goFunction: 'IsChinese', description: '判断字符是否为汉字', category: '文本', params: ['字符'], returnType: '逻辑型' },
    { chineseName: '是否为数字', goFunction: 'IsDigit', description: '判断字符是否为数字', category: '文本', params: ['字符'], returnType: '逻辑型' },
    { chineseName: '是否为字母', goFunction: 'IsLetter', description: '判断字符是否为字母', category: '文本', params: ['字符'], returnType: '逻辑型' },
    { chineseName: '取文本哈希', goFunction: 'HashStr', description: '计算文本的哈希值', category: '文本', params: ['文本', '算法'], returnType: '文本型' },
    { chineseName: '取UUID', goFunction: 'GetUUID', description: '生成唯一标识UUID', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '取GUID', goFunction: 'GetGUID', description: '生成全局唯一标识GUID', category: '系统', params: [], returnType: '文本型' },
    { chineseName: '置窗口最小尺寸', goFunction: 'SetFormMinSize', description: '设置窗口最小尺寸', category: '窗口', params: ['窗口', '最小宽度', '最小高度'], returnType: '' },
    { chineseName: '置窗口最大尺寸', goFunction: 'SetFormMaxSize', description: '设置窗口最大尺寸', category: '窗口', params: ['窗口', '最大宽度', '最大高度'], returnType: '' },
    { chineseName: '置窗口可否最大化', goFunction: 'SetFormMaximizable', description: '设置窗口是否可最大化', category: '窗口', params: ['窗口', '是否可最大化'], returnType: '' },
    { chineseName: '置窗口可否最小化', goFunction: 'SetFormMinimizable', description: '设置窗口是否可最小化', category: '窗口', params: ['窗口', '是否可最小化'], returnType: '' },
    { chineseName: '置窗口可否改变大小', goFunction: 'SetFormResizable', description: '设置窗口是否可改变大小', category: '窗口', params: ['窗口', '是否可改变'], returnType: '' },
    { chineseName: '置窗口Esc关闭', goFunction: 'SetFormEscClose', description: '设置按Esc是否关闭窗口', category: '窗口', params: ['窗口', '是否关闭'], returnType: '' },
    { chineseName: '置窗口Enter提交', goFunction: 'SetFormEnterSubmit', description: '设置按Enter是否提交', category: '窗口', params: ['窗口', '是否提交'], returnType: '' },
    { chineseName: '置窗口全屏', goFunction: 'SetFormFullScreen', description: '设置窗口全屏显示', category: '窗口', params: ['窗口', '是否全屏'], returnType: '' },
    { chineseName: '取窗口状态', goFunction: 'GetFormState', description: '获取窗口状态(正常/最小化/最大化)', category: '窗口', params: ['窗口'], returnType: '文本型' },
    { chineseName: '置窗口鼠标形状', goFunction: 'SetFormCursor', description: '设置窗口鼠标形状', category: '窗口', params: ['窗口', '光标类型'], returnType: '' },
    { chineseName: '置窗口任务栏显示', goFunction: 'SetFormTaskBar', description: '设置窗口任务栏显示', category: '窗口', params: ['窗口', '是否显示'], returnType: '' },
    { chineseName: '置窗口总在最前', goFunction: 'SetFormStayOnTop', description: '设置窗口总在最前', category: '窗口', params: ['窗口', '是否最前'], returnType: '' },
    { chineseName: '取窗口句柄', goFunction: 'GetFormHandle', description: '获取窗口句柄', category: '窗口', params: ['窗口'], returnType: '整数型' },
    { chineseName: '取组件句柄', goFunction: 'GetComponentHandle', description: '获取组件句柄', category: '窗口', params: ['组件'], returnType: '整数型' },
    { chineseName: '置组件拖放模式', goFunction: 'SetDragMode', description: '设置组件拖放模式', category: '窗口', params: ['组件', '拖放模式'], returnType: '' },
    { chineseName: '置组件接受拖放', goFunction: 'SetAcceptDrop', description: '设置组件是否接受拖放', category: '窗口', params: ['组件', '是否接受'], returnType: '' },
    { chineseName: '置组件双击间隔', goFunction: 'SetDblClickInterval', description: '设置双击间隔时间', category: '窗口', params: ['组件', '间隔毫秒'], returnType: '' },
    { chineseName: '置组件右键菜单', goFunction: 'SetPopupMenu', description: '设置组件右键弹出菜单', category: '窗口', params: ['组件', '菜单'], returnType: '' },
    { chineseName: '置组件内边距', goFunction: 'SetPadding', description: '设置组件内边距', category: '窗口', params: ['组件', '左', '顶', '右', '底'], returnType: '' },
    { chineseName: '置组件外边距', goFunction: 'SetMargin', description: '设置组件外边距', category: '窗口', params: ['组件', '左', '顶', '右', '底'], returnType: '' },
    { chineseName: '置组件自动大小', goFunction: 'SetAutoSize', description: '设置组件自动调整大小', category: '窗口', params: ['组件', '是否自动'], returnType: '' },
    { chineseName: '置组件Tab停靠', goFunction: 'SetTabStop', description: '设置组件是否Tab停靠', category: '窗口', params: ['组件', '是否停靠'], returnType: '' },
    { chineseName: '置组件只读', goFunction: 'SetReadOnly', description: '设置编辑框只读模式', category: '窗口', params: ['编辑框', '是否只读'], returnType: '' },
    { chineseName: '置组件密码模式', goFunction: 'SetPasswordMode', description: '设置编辑框密码模式', category: '窗口', params: ['编辑框', '是否密码模式', '密码字符'], returnType: '' },
    { chineseName: '置组件最大长度', goFunction: 'SetMaxLength', description: '设置编辑框最大输入长度', category: '窗口', params: ['编辑框', '最大长度'], returnType: '' },
    { chineseName: '置组件多行模式', goFunction: 'SetMultiLine', description: '设置编辑框多行模式', category: '窗口', params: ['编辑框', '是否多行'], returnType: '' },
    { chineseName: '置组件数字模式', goFunction: 'SetNumberMode', description: '设置编辑框仅允许数字', category: '窗口', params: ['编辑框', '是否数字模式'], returnType: '' },
    { chineseName: '置组件选择文本', goFunction: 'SetSelText', description: '设置编辑框选中文本', category: '窗口', params: ['编辑框', '起始位置', '长度'], returnType: '' },
    { chineseName: '取组件选择文本', goFunction: 'GetSelText', description: '获取编辑框选中文本', category: '窗口', params: ['编辑框'], returnType: '文本型' },
    { chineseName: '置组件滚动位置', goFunction: 'SetScrollPos', description: '设置组件滚动位置', category: '窗口', params: ['组件', '横向位置', '纵向位置'], returnType: '' },
    { chineseName: '取组件滚动位置', goFunction: 'GetScrollPos', description: '获取组件滚动位置', category: '窗口', params: ['组件'], returnType: '对象' },
    { chineseName: '置进度条位置', goFunction: 'SetProgressPos', description: '设置进度条当前位置', category: '窗口', params: ['进度条', '位置值'], returnType: '' },
    { chineseName: '取进度条位置', goFunction: 'GetProgressPos', description: '获取进度条当前位置', category: '窗口', params: ['进度条'], returnType: '整数型' },
    { chineseName: '置进度条范围', goFunction: 'SetProgressRange', description: '设置进度条范围', category: '窗口', params: ['进度条', '最小值', '最大值'], returnType: '' },
    { chineseName: '置滑块位置', goFunction: 'SetSliderPos', description: '设置滑块当前位置', category: '窗口', params: ['滑块', '位置值'], returnType: '' },
    { chineseName: '取滑块位置', goFunction: 'GetSliderPos', description: '获取滑块当前位置', category: '窗口', params: ['滑块'], returnType: '整数型' },
    { chineseName: '置滑块范围', goFunction: 'SetSliderRange', description: '设置滑块范围', category: '窗口', params: ['滑块', '最小值', '最大值'], returnType: '' },
    { chineseName: '置表格行列数', goFunction: 'SetGridSize', description: '设置表格行列数', category: '窗口', params: ['表格', '行数', '列数'], returnType: '' },
    { chineseName: '置表格单元格', goFunction: 'SetGridCell', description: '设置表格单元格内容', category: '窗口', params: ['表格', '行', '列', '内容'], returnType: '' },
    { chineseName: '取表格单元格', goFunction: 'GetGridCell', description: '获取表格单元格内容', category: '窗口', params: ['表格', '行', '列'], returnType: '文本型' },
    { chineseName: '置表格列宽', goFunction: 'SetGridColWidth', description: '设置表格列宽', category: '窗口', params: ['表格', '列', '宽度'], returnType: '' },
    { chineseName: '置表格行高', goFunction: 'SetGridRowHeight', description: '设置表格行高', category: '窗口', params: ['表格', '行', '高度'], returnType: '' },
    { chineseName: '置表格固定行列', goFunction: 'SetGridFixed', description: '设置表格固定行列数', category: '窗口', params: ['表格', '固定行', '固定列'], returnType: '' },
    { chineseName: '置树形框节点', goFunction: 'SetTreeNode', description: '设置树形框节点', category: '窗口', params: ['树形框', '父节点', '节点文本'], returnType: '节点句柄' },
    { chineseName: '取树形框选中', goFunction: 'GetTreeSelected', description: '获取树形框选中节点', category: '窗口', params: ['树形框'], returnType: '节点句柄' },
    { chineseName: '置树形框展开', goFunction: 'SetTreeExpand', description: '展开树形框节点', category: '窗口', params: ['树形框', '节点句柄', '是否展开'], returnType: '' },
    { chineseName: '置组合框项目', goFunction: 'SetComboItems', description: '设置组合框项目列表', category: '窗口', params: ['组合框', '项目数组'], returnType: '' },
    { chineseName: '取组合框选中', goFunction: 'GetComboSelected', description: '获取组合框选中项索引', category: '窗口', params: ['组合框'], returnType: '整数型' },
    { chineseName: '取组合框文本', goFunction: 'GetComboText', description: '获取组合框当前文本', category: '窗口', params: ['组合框'], returnType: '文本型' },
    { chineseName: '置列表框项目', goFunction: 'SetListItems', description: '设置列表框项目列表', category: '窗口', params: ['列表框', '项目数组'], returnType: '' },
    { chineseName: '取列表框选中', goFunction: 'GetListSelected', description: '获取列表框选中项索引', category: '窗口', params: ['列表框'], returnType: '整数型' },
    { chineseName: '置选项卡页面', goFunction: 'SetTabPage', description: '添加选项卡页面', category: '窗口', params: ['选项卡', '页面标题'], returnType: '页面句柄' },
    { chineseName: '取选项卡当前页', goFunction: 'GetTabActivePage', description: '获取选项卡当前页面索引', category: '窗口', params: ['选项卡'], returnType: '整数型' },
    { chineseName: '置菜单项', goFunction: 'SetMenuItem', description: '添加菜单项', category: '窗口', params: ['菜单', '父菜单', '菜单文本', '快捷键'], returnType: '菜单项句柄' },
    { chineseName: '置菜单分隔线', goFunction: 'SetMenuSeparator', description: '添加菜单分隔线', category: '窗口', params: ['菜单', '父菜单'], returnType: '' },
    { chineseName: '置工具栏按钮', goFunction: 'SetToolButton', description: '添加工具栏按钮', category: '窗口', params: ['工具栏', '按钮文本', '图标路径', '提示文本'], returnType: '按钮句柄' },
    { chineseName: '置工具栏分隔', goFunction: 'SetToolSeparator', description: '添加工具栏分隔符', category: '窗口', params: ['工具栏'], returnType: '' },
    { chineseName: '置状态栏文本', goFunction: 'SetStatusText', description: '设置状态栏文本', category: '窗口', params: ['状态栏', '文本', '面板索引'], returnType: '' },
    { chineseName: '置状态栏面板', goFunction: 'SetStatusPanel', description: '添加状态栏面板', category: '窗口', params: ['状态栏', '面板宽度'], returnType: '' },
    { chineseName: '置定时器启动', goFunction: 'StartTimer', description: '启动定时器', category: '窗口', params: ['定时器'], returnType: '' },
    { chineseName: '置定时器停止', goFunction: 'StopTimer', description: '停止定时器', category: '窗口', params: ['定时器'], returnType: '' },
    { chineseName: '置定时器间隔', goFunction: 'SetTimerInterval', description: '设置定时器间隔', category: '窗口', params: ['定时器', '间隔毫秒'], returnType: '' },
    { chineseName: '置浏览器URL', goFunction: 'SetBrowserURL', description: '设置浏览器导航URL', category: '窗口', params: ['浏览器', 'URL地址'], returnType: '' },
    { chineseName: '置浏览器HTML', goFunction: 'SetBrowserHTML', description: '设置浏览器HTML内容', category: '窗口', params: ['浏览器', 'HTML内容'], returnType: '' },
    { chineseName: '浏览器后退', goFunction: 'BrowserGoBack', description: '浏览器后退', category: '窗口', params: ['浏览器'], returnType: '' },
    { chineseName: '浏览器前进', goFunction: 'BrowserGoForward', description: '浏览器前进', category: '窗口', params: ['浏览器'], returnType: '' },
    { chineseName: '浏览器刷新', goFunction: 'BrowserRefresh', description: '浏览器刷新', category: '窗口', params: ['浏览器'], returnType: '' },
    { chineseName: '浏览器执行JS', goFunction: 'BrowserExecJS', description: '浏览器执行JavaScript', category: '窗口', params: ['浏览器', 'JS代码'], returnType: '' },
    { chineseName: '置图片框图片', goFunction: 'SetImagePicture', description: '设置图片框显示的图片', category: '窗口', params: ['图片框', '图片路径'], returnType: '' },
    { chineseName: '置图片框缩放', goFunction: 'SetImageStretch', description: '设置图片框缩放模式', category: '窗口', params: ['图片框', '是否缩放'], returnType: '' },
    { chineseName: '置图片框居中', goFunction: 'SetImageCenter', description: '设置图片框居中显示', category: '窗口', params: ['图片框', '是否居中'], returnType: '' },
    { chineseName: '置图片框等比', goFunction: 'SetImageProportional', description: '设置图片框等比缩放', category: '窗口', params: ['图片框', '是否等比'], returnType: '' },
    { chineseName: '置日期框日期', goFunction: 'SetDatePickerDate', description: '设置日期框日期', category: '窗口', params: ['日期框', '日期'], returnType: '' },
    { chineseName: '取日期框日期', goFunction: 'GetDatePickerDate', description: '获取日期框日期', category: '窗口', params: ['日期框'], returnType: '日期时间型' },
    { chineseName: '置日期框格式', goFunction: 'SetDatePickerFormat', description: '设置日期框显示格式', category: '窗口', params: ['日期框', '格式文本'], returnType: '' },
    { chineseName: '置颜色框颜色', goFunction: 'SetColorBoxColor', description: '设置颜色框颜色', category: '窗口', params: ['颜色框', '颜色值'], returnType: '' },
    { chineseName: '取颜色框颜色', goFunction: 'GetColorBoxColor', description: '获取颜色框颜色', category: '窗口', params: ['颜色框'], returnType: '文本型' },
    { chineseName: '置月历日期', goFunction: 'SetCalendarDate', description: '设置月历选中日期', category: '窗口', params: ['月历', '日期'], returnType: '' },
    { chineseName: '取月历日期', goFunction: 'GetCalendarDate', description: '获取月历选中日期', category: '窗口', params: ['月历'], returnType: '日期时间型' },
    { chineseName: '置外形框形状', goFunction: 'SetShapeType', description: '设置外形框形状', category: '窗口', params: ['外形框', '形状类型'], returnType: '' },
    { chineseName: '置外形框画笔', goFunction: 'SetShapePen', description: '设置外形框画笔', category: '窗口', params: ['外形框', '颜色', '宽度', '样式'], returnType: '' },
    { chineseName: '置外形框画刷', goFunction: 'SetShapeBrush', description: '设置外形框画刷', category: '窗口', params: ['外形框', '颜色', '样式'], returnType: '' },
    { chineseName: '置分隔条方向', goFunction: 'SetSplitterAlign', description: '设置分隔条方向', category: '窗口', params: ['分隔条', '方向'], returnType: '' },
    { chineseName: '置滚动条范围', goFunction: 'SetScrollBarRange', description: '设置滚动条范围', category: '窗口', params: ['滚动条', '最小值', '最大值'], returnType: '' },
    { chineseName: '置滚动条位置', goFunction: 'SetScrollBarPos', description: '设置滚动条位置', category: '窗口', params: ['滚动条', '位置值'], returnType: '' },
    { chineseName: '取滚动条位置', goFunction: 'GetScrollBarPos', description: '获取滚动条位置', category: '窗口', params: ['滚动条'], returnType: '整数型' },
    { chineseName: '置调节器范围', goFunction: 'SetSpinEditRange', description: '设置调节器范围', category: '窗口', params: ['调节器', '最小值', '最大值'], returnType: '' },
    { chineseName: '置调节器值', goFunction: 'SetSpinEditValue', description: '设置调节器值', category: '窗口', params: ['调节器', '值'], returnType: '' },
  ]

  allFunctions.value = defaultFns
  const cats = new Set(defaultFns.map(f => f.category))
  fnCategories.value = Array.from(cats).sort()
}

const selectFnCategory = async (cat: string) => {
  if (expandedFnCategories.value.has(cat)) {
    expandedFnCategories.value.delete(cat)
  } else {
    expandedFnCategories.value.add(cat)
  }
  selectedFnCategory.value = cat
  try {
    categoryFunctions.value = await GetFunctionsByCategory(cat)
  } catch (e) {
    categoryFunctions.value = allFunctions.value.filter(f => f.category === cat)
  }
}

const getCategoryColor = (cat: string): string => {
  const colors: Record<string, string> = {
    '窗口': '#409eff', '文件': '#67c23a', '字符串': '#e6a23c',
    '数值': '#f56c6c', '数组': '#909399', '组件': '#409eff',
    '系统': '#67c23a', '网络': '#e6a23c', '数据库': '#f56c6c',
    'JSON': '#909399', '加密': '#409eff', '多媒体': '#67c23a',
    '定时器': '#e6a23c', '线程': '#f56c6c',
  }
  return colors[cat] || '#909399'
}

const onDragStart = (event: DragEvent, component: UIComponent) => {
  if (event.dataTransfer) {
    event.dataTransfer.setData('application/json', JSON.stringify(component))
    event.dataTransfer.effectAllowed = 'copy'
  }
}

const onSelectWindowType = (type: string) => {
  emit('select-window-type', type)
}

const onSelectFunction = (fn: ChineseFunction) => {
  emit('select-function', fn)
}

const addNewWindow = () => {
  const newWindow: ProjectWindow = {
    id: `win_${nextWindowId.value}`,
    name: `窗口${nextWindowId.value}`,
    type: 'normal',
    icon: 'Monitor',
    isActive: false,
    codeFiles: defaultCodeFiles(),
  }
  nextWindowId.value++
  projectWindows.value.push(newWindow)
  emit('add-new-window', newWindow)
}

const selectWindow = (win: ProjectWindow) => {
  projectWindows.value.forEach(w => w.isActive = false)
  win.isActive = true
  emit('open-window-design', win.id)
}

const openWindowCode = (win: ProjectWindow, codeFileId?: string) => {
  if (codeFileId) {
    emit('open-window-code', win.id, codeFileId)
  } else if (win.codeFiles.length > 0) {
    emit('open-window-code', win.id, win.codeFiles[0].id)
  }
}

const deleteWindow = (win: ProjectWindow) => {
  if (projectWindows.value.length <= 1) return
  const index = projectWindows.value.findIndex(w => w.id === win.id)
  if (index > -1) {
    if (win.isActive && index > 0) {
      projectWindows.value[index - 1].isActive = true
    } else if (win.isActive && projectWindows.value.length > 1) {
      projectWindows.value[Math.min(index + 1, projectWindows.value.length - 1)].isActive = true
    }
    projectWindows.value.splice(index, 1)
    emit('delete-window', win.id)
  }
}

const startRenameWindow = (win: ProjectWindow) => {
  editingWindowId.value = win.id
  editName.value = win.name
}

const confirmRenameWindow = (win: ProjectWindow) => {
  if (editName.value.trim()) {
    win.name = editName.value.trim()
    emit('rename-window', win.id, win.name)
  }
  editingWindowId.value = null
}

const cancelRename = () => {
  editingWindowId.value = null
  editingCodeFileId.value = null
}

const showContextMenu = (e: MouseEvent, target: { type: 'window'; data: any }) => {
  e.preventDefault()
  contextMenuTarget.value = target
  contextMenuPosition.value = { x: e.clientX, y: e.clientY }
  contextMenuVisible.value = true
}

const hideContextMenu = () => {
  contextMenuVisible.value = false
  contextMenuTarget.value = null
}

const handleContextMenuAction = (action: string) => {
  if (!contextMenuTarget.value) return
  const { type, data } = contextMenuTarget.value
  if (type === 'window') {
    handleWindowContextAction(action, data as ProjectWindow)
  }
  hideContextMenu()
}

const handleWindowContextAction = (action: string, win: ProjectWindow) => {
  switch (action) {
    case 'open-design':
      selectWindow(win)
      break
    case 'open-code':
      openWindowCode(win)
      break
    case 'rename':
      startRenameWindow(win)
      break
    case 'copy':
      const newWin: ProjectWindow = {
        id: `win_${nextWindowId.value}`,
        name: win.name + '_副本',
        type: win.type,
        icon: win.icon,
        isActive: false,
        codeFiles: [{ id: `cf_${nextCodeFileId.value}`, name: 'main.ego', code: '' }],
      }
      nextWindowId.value++
      nextCodeFileId.value++
      projectWindows.value.push(newWin)
      emit('add-new-window', newWin)
      break
    case 'delete':
      deleteWindow(win)
      break
    case 'set-startup':
      projectWindows.value.forEach(w => { w.isActive = false })
      win.isActive = true
      emit('open-window-design', win.id)
      break
  }
}

const toggleCodeFilesExpand = (winId: string) => {
  if (expandedCodeFiles.value.has(winId)) {
    expandedCodeFiles.value.delete(winId)
  } else {
    expandedCodeFiles.value.add(winId)
  }
}

const openCodeFile = (win: ProjectWindow, file: CodeFile) => {
  emit('open-window-code', win.id, file.id)
}

const addCodeFile = (win: ProjectWindow) => {
  const newFile: CodeFile = {
    id: `cf_${nextCodeFileId.value}`,
    name: `code_${nextCodeFileId.value}.ego`,
    code: '',
  }
  nextCodeFileId.value++
  win.codeFiles.push(newFile)
  if (!expandedCodeFiles.value.has(win.id)) {
    expandedCodeFiles.value.add(win.id)
  }
  emit('add-code-file', win.id, newFile)
}

const deleteCodeFile = (win: ProjectWindow, file: CodeFile) => {
  if (win.codeFiles.length <= 1) return
  const index = win.codeFiles.findIndex(f => f.id === file.id)
  if (index > -1) {
    win.codeFiles.splice(index, 1)
    emit('delete-code-file', win.id, file.id)
  }
}

const startRenameCodeFile = (win: ProjectWindow, file: CodeFile) => {
  editingCodeFileId.value = { windowId: win.id, fileId: file.id }
  editName.value = file.name
}

const confirmRenameCodeFile = (win: ProjectWindow, file: CodeFile) => {
  if (editName.value.trim()) {
    file.name = editName.value.trim()
    emit('rename-code-file', win.id, file.id, file.name)
  }
  editingCodeFileId.value = null
}

const toggleAddMenu = () => {
  showAddMenu.value = !showAddMenu.value
}

const closeAddMenu = () => {
  showAddMenu.value = false
}

const addIndependentCodeFile = () => {
  showAddMenu.value = false
  const newFile: CodeFile = {
    id: `icf_${nextIndependentFileId.value}`,
    name: `module_${nextIndependentFileId.value}.ego`,
    code: '',
  }
  nextIndependentFileId.value++
  independentCodeFiles.value.push(newFile)
  emit('add-independent-code-file', newFile)
}

const openIndependentCodeFile = (file: CodeFile) => {
  emit('open-independent-code-file', file.id)
}

const deleteIndependentCodeFile = (file: CodeFile) => {
  const index = independentCodeFiles.value.findIndex(f => f.id === file.id)
  if (index > -1) {
    independentCodeFiles.value.splice(index, 1)
    emit('delete-independent-code-file', file.id)
  }
}

const startRenameIndependentCodeFile = (file: CodeFile) => {
  editingCodeFileId.value = { windowId: '', fileId: file.id }
  editName.value = file.name
}

const confirmRenameIndependentCodeFile = (file: CodeFile) => {
  if (editName.value.trim()) {
    file.name = editName.value.trim()
    emit('rename-independent-code-file', file.id, file.name)
  }
  editingCodeFileId.value = null
}

const loadFunctions = async () => {
  try {
    allFunctions.value = await GetAllFunctions()
    fnCategories.value = await GetCategories()
  } catch (e) {
    console.error('加载函数库失败，使用默认函数库:', e)
    loadDefaultFunctions()
  }
}

defineExpose({
  projectWindows,
  independentCodeFiles,
  addNewWindow,
  addCodeFile,
  addIndependentCodeFile,
})

onMounted(() => {
  loadFunctions()
})
</script>

<template>
  <div class="left-panel">
    <div v-if="!hideTabs" class="panel-tabs">
      <div
        class="panel-tab"
        :class="{ active: effectiveTab === 'functions' }"
        @click="activeTab = 'functions'"
      >
        <el-icon><FolderOpened /></el-icon>
        <span>函数库</span>
      </div>
      <div
        class="panel-tab"
        :class="{ active: effectiveTab === 'components' }"
        @click="activeTab = 'components'"
      >
        <el-icon><Grid /></el-icon>
        <span>组件库</span>
      </div>
      <div
        class="panel-tab"
        :class="{ active: effectiveTab === 'construct' }"
        @click="activeTab = 'construct'"
      >
        <el-icon><Setting /></el-icon>
        <span>构造</span>
      </div>
    </div>

    <div class="panel-body">
      <!-- 函数库 -->
      <div v-show="effectiveTab === 'functions'" class="functions-view">
        <div class="fn-search">
          <el-input size="small" placeholder="搜索函数..." clearable />
        </div>
        <div class="fn-tree">
          <div v-for="cat in fnCategories" :key="cat" class="fn-category">
            <div class="fn-cat-header" @click="selectFnCategory(cat)">
              <span class="cat-arrow">{{ expandedFnCategories.has(cat) ? '▼' : '▶' }}</span>
              <span class="cat-dot" :style="{ backgroundColor: getCategoryColor(cat) }"></span>
              <span class="cat-name">{{ cat }}</span>
              <span class="cat-count">{{ allFunctions.filter(f => f.category === cat).length }}</span>
            </div>
            <div v-show="expandedFnCategories.has(cat)" class="fn-list">
              <div
                v-for="fn in (selectedFnCategory === cat ? categoryFunctions : allFunctions.filter(f => f.category === cat))"
                :key="fn.chineseName"
                class="fn-item"
                @click="onSelectFunction(fn)"
                :title="fn.description"
              >
                <span class="fn-name">{{ fn.chineseName }}</span>
                <span v-if="fn.returnType" class="fn-return">{{ fn.returnType }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 组件库 - 树列表模式，默认收缩 -->
      <div v-show="effectiveTab === 'components'" class="components-view">
        <div v-for="category in categories" :key="category.name" class="comp-category-tree">
          <div class="comp-cat-header" @click="toggleCompCategory(category.name)">
            <span class="cat-arrow">{{ expandedCompCategories.has(category.name) ? '▼' : '▶' }}</span>
            <span class="cat-name">{{ category.name }}</span>
            <span class="cat-count">{{ category.components.length }}</span>
          </div>
          <div v-show="expandedCompCategories.has(category.name)" class="comp-tree-list">
            <div
              v-for="comp in category.components"
              :key="comp.id"
              class="comp-tree-item"
              draggable="true"
              @dragstart="onDragStart($event, comp)"
            >
              <el-icon class="comp-icon"><component :is="comp.icon" /></el-icon>
              <span class="comp-name">{{ comp.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 构造面板 - 专业IDE风格 -->
      <div v-show="effectiveTab === 'construct'" class="construct-view" @click="closeAddMenu(); hideContextMenu()">

        <!-- ========== 项目结构树 - VSCode风格 ========== -->
        <div class="construct-section">
          <div class="section-header">
            <span class="section-title">
              <el-icon><FolderOpened /></el-icon>
              项目资源管理器
            </span>
            <div class="toolbar-actions">
              <el-dropdown trigger="click" @command="handleToolbarCommand" size="small">
                <el-button size="small" text title="新建">
                  <el-icon><Plus /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="new-window">
                      <el-icon><Monitor /></el-icon> 新建窗体
                    </el-dropdown-item>
                    <el-dropdown-item command="new-code">
                      <el-icon><DocumentAdd /></el-icon> 新建代码文件
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <el-button size="small" text title="刷新" @click="refreshExplorer">
                <el-icon><Refresh /></el-icon>
              </el-button>
              <el-button size="small" text title="折叠全部" @click="collapseAll">
                <el-icon><DArrowLeft /></el-icon>
              </el-button>
            </div>
          </div>

          <!-- 解决方案/项目根节点 -->
          <div class="explorer-tree">
            <div class="tree-item solution-root" @click="toggleSolutionExpand">
              <span class="tree-indent"></span>
              <span class="tree-chevron" :class="{ expanded: solutionExpanded }">▶</span>
              <span class="tree-icon solution-icon">�</span>
              <span class="tree-label solution-label">{{ projectName || '未命名项目' }}</span>
              <span class="tree-badge">{{ projectWindows.length + independentCodeFiles.length }} 项</span>
            </div>

                        <template v-if="solutionExpanded">
              <!-- 源码目录 src/ -->
              <div class="tree-group-container">
                <div class="tree-item folder-item" @click="toggleFolderExpand('src')">
                  <span class="tree-indent"></span>
                  <span class="tree-chevron" :class="{ expanded: folderExpands.src }">&#9654;</span>
                  <span class="tree-icon folder-icon-open" v-if="folderExpands.src">&#128194;</span>
                  <span class="tree-icon folder-icon" v-else>&#128193;</span>
                  <span class="tree-label folder-label">src</span>
                  <span class="tree-badge-folder">源码</span>
                </div>

                <template v-if="folderExpands.src">
                  <!-- .ego 布局文件 -->
                  <div class="tree-item section-label">
                    <span class="tree-indent"></span>
                    <span class="tree-indent"></span>
                    <span class="tree-label sub-label">布局文件</span>
                  </div>
                  <div
                    v-for="win in projectWindows"
                    :key="'explorer_win_' + win.id"
                    class="tree-item file-item"
                    :class="{ active: win.isActive, selected: selectedTreeItem === win.id }"
                    @click="selectWindow(win); selectedTreeItem = win.id"
                    @contextmenu.prevent="showContextMenu($event, { type: 'window', data: win })"
                  >
                    <span class="tree-indent"></span>
                    <span class="tree-indent"></span>
                    <span class="tree-indent"></span>
                    <span class="tree-chevron invisible">&#9654;</span>
                    <span class="tree-icon ego-icon">&#128196;</span>
                    <span class="tree-label file-label">{{ win.name }}.ego</span>
                    <span v-if="win.isActive" class="file-badge startup-badge">启动</span>
                  </div>

                  <!-- .go 源码文件 -->
                  <div class="tree-item section-label" v-if="independentCodeFiles.length > 0">
                    <span class="tree-indent"></span>
                    <span class="tree-indent"></span>
                    <span class="tree-label sub-label">源码文件</span>
                  </div>
                  <div
                    v-for="file in independentCodeFiles"
                    :key="'explorer_code_' + file.id"
                    class="tree-item file-item"
                    :class="{ selected: selectedTreeItem === file.id }"
                    @click="openIndependentCodeFile(file); selectedTreeItem = file.id"
                  >
                    <span class="tree-indent"></span>
                    <span class="tree-indent"></span>
                    <span class="tree-indent"></span>
                    <span class="tree-chevron invisible">&#9654;</span>
                    <span class="tree-icon go-icon">&#9889;</span>
                    <span class="tree-label file-label">{{ file.name }}</span>
                  </div>
                </template>
              </div>

              <!-- 构建输出 build/ -->
              <div class="tree-group-container">
                <div class="tree-item folder-item" @click="toggleFolderExpand('build')">
                  <span class="tree-indent"></span>
                  <span class="tree-chevron" :class="{ expanded: folderExpands.build }">&#9654;</span>
                  <span class="tree-icon folder-icon-open" v-if="folderExpands.build">&#128194;</span>
                  <span class="tree-icon folder-icon" v-else>&#128193;</span>
                  <span class="tree-label folder-label">build</span>
                  <span class="tree-badge-folder">输出</span>
                </div>

                <template v-if="folderExpands.build">
                  <div class="tree-item file-item empty-folder">
                    <span class="tree-indent"></span>
                    <span class="tree-indent"></span>
                    <span class="tree-chevron invisible">&#9654;</span>
                    <span class="tree-icon">&#128160;</span>
                    <span class="tree-label file-label">{{ projectName || '项目' }}.exe</span>
                    <span class="file-badge exe-badge">编译</span>
                  </div>
                </template>
              </div>

              <!-- 内部资源 _int/ -->
              <div class="tree-group-container">
                <div class="tree-item folder-item" @click="toggleFolderExpand('int')">
                  <span class="tree-indent"></span>
                  <span class="tree-chevron" :class="{ expanded: folderExpands.int }">&#9654;</span>
                  <span class="tree-icon folder-icon-open" v-if="folderExpands.int">&#128194;</span>
                  <span class="tree-icon folder-icon" v-else>&#128193;</span>
                  <span class="tree-label folder-label">_int</span>
                  <span class="tree-badge-folder">内部</span>
                </div>

                <template v-if="folderExpands.int">
                  <div class="tree-item file-item empty-folder">
                    <span class="tree-indent"></span>
                    <span class="tree-indent"></span>
                    <span class="tree-chevron invisible">&#9654;</span>
                    <span class="tree-icon resource-icon">&#9881;&#65039;</span>
                    <span class="tree-label file-label empty-label">运行时库</span>
                  </div>
                </template>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <div
      v-if="contextMenuVisible && contextMenuTarget"
      class="context-menu"
      :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
      @click.stop
    >
      <template v-if="contextMenuTarget.type === 'window'">
        <div class="ctx-title">{{ contextMenuTarget.data.name }}</div>
        <div class="ctx-item" @click="handleContextMenuAction('open-design')">
          <el-icon><Monitor /></el-icon> 打开设计视图
        </div>
        <div class="ctx-item" @click="handleContextMenuAction('open-code')">
          <el-icon><Document /></el-icon> 打开代码编辑器
        </div>
        <div class="ctx-divider"></div>
        <div class="ctx-item" @click="handleContextMenuAction('rename')">
          <el-icon><EditPen /></el-icon> 重命名
        </div>
        <div class="ctx-item" @click="handleContextMenuAction('copy')">
          <el-icon><CopyDocument /></el-icon> 复制窗体
        </div>
        <div class="ctx-divider"></div>
        <div
          class="ctx-item"
          :class="{ disabled: contextMenuTarget.data?.isActive }"
          @click="!contextMenuTarget.data?.isActive && handleContextMenuAction('set-startup')"
        >
          <el-icon><Star /></el-icon> 设为启动窗体
        </div>
        <div class="ctx-divider"></div>
        <div
          class="ctx-item danger"
          :class="{ disabled: projectWindows.length <= 1 }"
          @click="projectWindows.length > 1 && handleContextMenuAction('delete')"
        >
          <el-icon><Delete /></el-icon> 删除窗体
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.left-panel {
  width: 240px;
  height: 100%;
  background-color: var(--panel-bg);
  border-right: 1px solid var(--panel-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.panel-tabs {
  display: flex;
  border-bottom: 1px solid var(--panel-border);
  background-color: var(--panel-header-bg);
  flex-shrink: 0;
}

.panel-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 8px 4px;
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
  user-select: none;
}

.panel-tab:hover {
  color: var(--color-primary);
  background-color: var(--color-primary-light);
}

.panel-tab.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
  font-weight: 600;
}

.panel-body {
  flex: 1;
  overflow-y: auto;
}

/* 函数库样式 */
.functions-view {
  padding: 4px;
}

.fn-search {
  padding: 6px;
}

.fn-category {
  margin-bottom: 2px;
}

.fn-cat-header {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  cursor: pointer;
  border-radius: 4px;
  font-size: 12px;
  user-select: none;
}

.fn-cat-header:hover {
  background-color: var(--bg-secondary);
}

.cat-arrow {
  font-size: 10px;
  color: var(--text-secondary);
  width: 12px;
}

.cat-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.cat-name {
  flex: 1;
  font-weight: 500;
  color: var(--text-primary);
}

.cat-count {
  font-size: 10px;
  color: var(--text-placeholder);
  background-color: var(--bg-secondary);
  padding: 0 6px;
  border-radius: 8px;
}

.fn-list {
  padding-left: 20px;
}

.fn-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px;
  cursor: pointer;
  border-radius: 3px;
  font-size: 12px;
}

.fn-item:hover {
  background-color: var(--color-primary-light);
}

.fn-name {
  color: var(--text-regular);
}

.fn-return {
  font-size: 10px;
  color: var(--color-success);
  background-color: var(--color-success-light);
  padding: 0 4px;
  border-radius: 2px;
}

/* 组件库样式 - 树列表模式 */
.components-view {
  padding: 4px;
}

.comp-category-tree {
  margin-bottom: 2px;
}

.comp-cat-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
  user-select: none;
}

.comp-cat-header:hover {
  background-color: var(--bg-hover);
}

.comp-cat-header .cat-arrow {
  font-size: 10px;
  color: var(--text-secondary);
  width: 12px;
  display: inline-block;
}

.comp-cat-header .cat-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  flex: 1;
}

.comp-cat-header .cat-count {
  font-size: 11px;
  color: var(--text-secondary);
  background: var(--bg-hover);
  padding: 1px 6px;
  border-radius: 10px;
}

.comp-tree-list {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding-left: 18px;
}

.comp-tree-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: grab;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.comp-tree-item:hover {
  background-color: var(--color-primary-light);
  border-color: var(--color-primary-light);
}

.comp-tree-item:active {
  cursor: grabbing;
  background-color: var(--bg-active);
}

.comp-icon {
  font-size: 15px;
  color: var(--color-primary);
}

.comp-name {
  font-size: 13px;
  color: var(--text-regular);
}

/* 构造面板样式 */
.construct-view {
  padding: 8px;
  position: relative;
}

.construct-section {
  margin-bottom: 12px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  margin-bottom: 4px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-primary);
}

.item-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 8px;
  border-radius: 4px;
  transition: all 0.15s;
  border: 1px solid transparent;
}

.list-item:hover {
  background-color: var(--bg-secondary);
}

.list-item.active {
  background-color: var(--color-primary-light);
  border-color: var(--color-primary);
}

.item-main {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.item-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.window-icon {
  color: #409eff;
}

.module-icon {
  color: #67c23a;
}

.item-name {
  font-size: 12px;
  color: var(--text-regular);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
}

.item-name:hover {
  color: var(--color-primary);
}

.item-name-edit {
  flex: 1;
  max-width: 120px;
}

.item-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.15s;
}

.list-item:hover .item-actions {
  opacity: 1;
}

.badge {
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 8px;
  flex-shrink: 0;
}

.badge-startup {
  background-color: #e6f7ff;
  color: #1890ff;
  border: 1px solid #91d5ff;
}

.class-count {
  font-size: 10px;
  color: var(--text-placeholder);
  background-color: var(--bg-secondary);
  padding: 1px 6px;
  border-radius: 8px;
  flex-shrink: 0;
}

/* 代码集样式 */
.module-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.module-item {
  border: 1px solid var(--panel-border);
  border-radius: 4px;
  overflow: hidden;
  transition: all 0.15s;
}

.module-item.active {
  border-color: #67c23a;
  background-color: #f0f9eb;
}

.module-header {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  cursor: pointer;
}

.module-header:hover {
  background-color: var(--bg-secondary);
}

.expand-btn {
  font-size: 10px;
  color: var(--text-secondary);
  width: 14px;
  cursor: pointer;
  user-select: none;
}

.expand-btn:hover {
  color: var(--color-primary);
}

.class-list {
  border-top: 1px solid var(--panel-border);
  background-color: var(--bg-secondary);
  padding: 4px;
}

.class-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 11px;
}

.class-item:hover {
  background-color: var(--bg-active);
}

.class-icon {
  font-size: 12px;
  color: #e6a23c;
}

.class-name {
  flex: 1;
  color: var(--text-regular);
}

.add-class-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  font-size: 11px;
  color: var(--color-success);
  cursor: pointer;
  border-radius: 3px;
  border: 1px dashed var(--color-success);
  margin-top: 2px;
  justify-content: center;
}

.add-class-btn:hover {
  background-color: var(--color-success-light);
}

.empty-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px;
  color: var(--text-placeholder);
  font-size: 12px;
}

.empty-hint .el-icon {
  font-size: 24px;
}

/* 窗口类型网格 */
.type-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 4px;
}

.type-card {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 4px;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.15s;
}

.type-card:hover {
  background-color: var(--color-primary-light);
  border-color: var(--color-primary);
}

.type-icon {
  font-size: 18px;
  color: var(--color-primary);
}

.type-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.type-name {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
}

.type-desc {
  font-size: 10px;
  color: var(--text-placeholder);
}

/* 项目结构树 - VSCode风格 */
.explorer-tree {
  font-size: 12px;
  user-select: none;
}

.tree-item {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 3px 8px;
  cursor: pointer;
  border-radius: 3px;
  transition: all 0.1s ease;
  position: relative;
}

.tree-item:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.solution-root {
  font-weight: 600;
  font-size: 12px;
  padding: 6px 8px;
}

.solution-root:hover {
  background-color: rgba(0, 0, 0, 0.06);
}

.solution-icon {
  font-size: 14px;
  margin-right: 4px;
}

.solution-label {
  flex: 1;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tree-badge {
  font-size: 10px;
  color: var(--text-placeholder);
  background-color: var(--bg-hover);
  padding: 1px 6px;
  border-radius: 10px;
  font-weight: normal;
}

.tree-chevron {
  font-size: 10px;
  color: var(--text-secondary);
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s ease;
  flex-shrink: 0;
}

.tree-chevron.expanded {
  transform: rotate(90deg);
}

.tree-chevron.invisible {
  visibility: hidden;
}

.tree-indent {
  width: 16px;
  flex-shrink: 0;
}

.tree-icon {
  font-size: 13px;
  margin-right: 4px;
  flex-shrink: 0;
  width: 16px;
  text-align: center;
}

.folder-icon,
.folder-icon-open {
  color: #dcb67a;
}

.form-icon {
  color: #519aba;
}

.code-icon {
  color: #6ec1a7;
}

.resource-icon {
  color: #dcb67a;
}

.tree-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
}

.group-header {
  font-weight: 500;
  font-size: 11px;
}

.folder-item {
  font-weight: 500;
  font-size: 12px;
  padding: 4px 8px;
}

.folder-item:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.folder-label {
  color: var(--text-primary);
  font-weight: 500;
}

.tree-badge-folder {
  font-size: 9px;
  color: var(--text-placeholder);
  background-color: var(--bg-hover);
  padding: 1px 5px;
  border-radius: 8px;
  font-weight: 400;
  flex-shrink: 0;
  margin-left: 4px;
}

.section-label {
  padding: 2px 8px 2px 32px;
  cursor: default;
}

.section-label:hover {
  background-color: transparent;
}

.sub-label {
  font-size: 11px;
  color: var(--text-secondary);
  font-weight: 400;
  letter-spacing: 0.3px;
}

.ego-icon {
  color: #519aba;
}

.go-icon {
  color: #00add8;
}

.exe-badge {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
  font-size: 9px;
  padding: 1px 5px;
  border-radius: 8px;
  font-weight: 500;
  flex-shrink: 0;
  margin-left: 4px;
}

.group-label {
  color: var(--text-primary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.tree-count {
  font-size: 10px;
  color: var(--text-placeholder);
  background-color: transparent;
  padding: 0 4px;
  font-weight: normal;
}

.tree-group-container {
  margin-left: 0;
}

.file-item {
  padding-left: 20px;
}

.file-item.active,
.file-item.selected {
  background-color: rgba(0, 122, 204, 0.1);
}

.file-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: var(--color-primary);
}

.file-label {
  color: var(--text-regular);
}

.file-item:hover .file-label {
  color: var(--text-primary);
}

.file-badge {
  font-size: 9px;
  padding: 1px 5px;
  border-radius: 8px;
  font-weight: 500;
  flex-shrink: 0;
  margin-left: 4px;
}

.startup-badge {
  background-color: #e6f7ff;
  color: #1890ff;
  border: 1px solid #91d5ff;
}

.empty-folder .empty-label {
  color: var(--text-disabled);
  font-style: italic;
}

.toolbar-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.15s;
}

.section-header:hover .toolbar-actions {
  opacity: 1;
}

/* 右键菜单 */
.context-menu {
  position: fixed;
  z-index: 9999;
  background: var(--bg-primary);
  border: 1px solid var(--panel-border);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 4px 0;
  min-width: 160px;
}

.ctx-title {
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 1px solid var(--panel-border);
  margin-bottom: 2px;
}

.ctx-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  font-size: 12px;
  color: var(--text-regular);
  cursor: pointer;
  transition: all 0.1s;
}

.ctx-item:hover {
  background-color: var(--color-primary-light);
  color: var(--color-primary);
}

.ctx-item.danger:hover {
  background-color: #fef0f0;
  color: #f56c6c;
}

.ctx-item.disabled {
  color: var(--text-disabled);
  cursor: not-allowed;
}

.ctx-item.disabled:hover {
  background-color: transparent;
  color: var(--text-disabled);
}

.ctx-divider {
  height: 1px;
  background-color: var(--panel-border);
  margin: 4px 8px;
}

.add-dropdown-wrapper {
  position: relative;
}

.add-dropdown-menu {
  position: absolute;
  right: 0;
  top: 100%;
  margin-top: 4px;
  background: var(--bg-primary);
  border: 1px solid var(--panel-border);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 4px 0;
  min-width: 150px;
  z-index: 1000;
}

.add-dropdown-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 14px;
  font-size: 12px;
  color: var(--text-regular);
  cursor: pointer;
  transition: all 0.1s;
}

.add-dropdown-item:hover {
  background-color: var(--color-primary-light);
  color: var(--color-primary);
}

.codefile-icon {
  color: #67c23a;
}

.independent-file-item {
  cursor: pointer;
}

.independent-file-item:hover {
  background-color: var(--bg-secondary);
}

.empty-hint {
  padding: 12px 8px;
  font-size: 11px;
  color: var(--text-placeholder);
  text-align: center;
}
</style>
