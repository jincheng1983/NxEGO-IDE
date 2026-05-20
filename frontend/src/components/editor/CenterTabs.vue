<script lang="ts" setup>
import { ref, reactive, computed, watch, nextTick } from 'vue'
import DesignCanvas from '../designer/DesignCanvas.vue'
import MonacoEditorWrapper from './MonacoEditorWrapper.vue'
import { useTheme } from '../../stores/theme'

const { currentTheme } = useTheme()

const editorTheme = computed(() => currentTheme.value === 'dark' ? 'vs-dark' : 'vs')

interface WindowTab {
  id: string
  name: string
  subTab: 'design' | 'code'
  activeCodeFileId?: string
}

interface ChineseFunction {
  chineseName: string
  goFunction: string
  description: string
  category: string
  params: string[]
  returnType: string
}

const windows = reactive<WindowTab[]>([
  { id: 'win_1', name: '启动窗口.ego', subTab: 'design' },
])

interface IndependentFile {
  id: string
  name: string
  code: string
}

const independentFiles = reactive<IndependentFile[]>([])
const activeIndependentFileId = ref<string | null>(null)

const activeWindowId = ref('win_1')
const viewMode = ref<'tab' | 'split'>('tab')
let nextWindowId = 2

const activeWindow = computed(() =>
  windows.find(w => w.id === activeWindowId.value)
)

const activeIndependentFile = computed(() =>
  independentFiles.find(f => f.id === activeIndependentFileId.value)
)

const emit = defineEmits<{
  (e: 'component-focus', component: any): void
  (e: 'code-rows-change', rows: any[]): void
  (e: 'windows-change', windows: WindowTab[]): void
}>()

const canvasRefs = ref<Record<string, any>>({})
const codeEditorRefs = ref<Record<string, any>>({})

const splitRatio = ref(50)
const isDragging = ref(false)
const generatedCode = ref('')
const focusedComponent = ref<any>(null)
const showFnLibrary = ref(true)
const fnSearchText = ref('')
const expandedFnCats = ref<Set<string>>(new Set())

const chineseFunctions = ref<ChineseFunction[]>([
  { chineseName: '信息框', goFunction: 'MessageBox', description: '弹出信息提示框', category: '界面', params: ['提示信息', '标题', '按钮类型'], returnType: '整数型' },
  { chineseName: '调试输出', goFunction: 'DebugPrint', description: '输出调试信息到控制台', category: '调试', params: ['输出内容'], returnType: '' },
  { chineseName: '延迟', goFunction: 'Sleep', description: '延迟指定毫秒数', category: '系统', params: ['毫秒数'], returnType: '' },
  { chineseName: '取文本长度', goFunction: 'Len', description: '取文本长度', category: '文本', params: ['文本'], returnType: '整数型' },
  { chineseName: '到文本', goFunction: 'Str', description: '转换到文本', category: '转换', params: ['数据'], returnType: '文本型' },
  { chineseName: '到整数', goFunction: 'Int', description: '转换到整数', category: '转换', params: ['文本'], returnType: '整数型' },
  { chineseName: '到数值', goFunction: 'Float', description: '转换到小数', category: '转换', params: ['文本'], returnType: '小数型' },
  { chineseName: '取启动时间', goFunction: 'GetTickCount', description: '取系统启动时间', category: '系统', params: [], returnType: '整数型' },
  { chineseName: '取现行时间', goFunction: 'Now', description: '获取当前时间', category: '时间', params: [], returnType: '日期时间型' },
  { chineseName: '运行', goFunction: 'Exec', description: '执行外部命令', category: '系统', params: ['命令', '是否等待'], returnType: '整数型' },
  { chineseName: '读入文件', goFunction: 'ReadFile', description: '读取文件内容', category: '文件', params: ['文件路径'], returnType: '字节集' },
  { chineseName: '写到文件', goFunction: 'WriteFile', description: '写入文件内容', category: '文件', params: ['文件路径', '数据'], returnType: '逻辑型' },
  { chineseName: '创建目录', goFunction: 'Mkdir', description: '创建目录', category: '文件', params: ['目录路径'], returnType: '逻辑型' },
  { chineseName: '删除文件', goFunction: 'RemoveFile', description: '删除文件', category: '文件', params: ['文件路径'], returnType: '逻辑型' },
  { chineseName: '文件是否存在', goFunction: 'FileExists', description: '判断文件是否存在', category: '文件', params: ['文件路径'], returnType: '逻辑型' },
  { chineseName: '取当前目录', goFunction: 'GetCurrentDir', description: '获取当前目录', category: '系统', params: [], returnType: '文本型' },
  { chineseName: '取随机数', goFunction: 'Random', description: '取随机数', category: '数学', params: ['最小值', '最大值'], returnType: '整数型' },
  { chineseName: '四舍五入', goFunction: 'Round', description: '四舍五入', category: '数学', params: ['数值', '小数位数'], returnType: '小数型' },
  { chineseName: '取绝对值', goFunction: 'Abs', description: '取绝对值', category: '数学', params: ['数值'], returnType: '小数型' },
  { chineseName: '求次方', goFunction: 'Pow', description: '求次方', category: '数学', params: ['底数', '指数'], returnType: '小数型' },
  { chineseName: '求平方根', goFunction: 'Sqrt', description: '求平方根', category: '数学', params: ['数值'], returnType: '小数型' },
  { chineseName: '求正弦', goFunction: 'Sin', description: '求正弦值', category: '数学', params: ['弧度'], returnType: '小数型' },
  { chineseName: '求余弦', goFunction: 'Cos', description: '求余弦值', category: '数学', params: ['弧度'], returnType: '小数型' },
  { chineseName: '求正切', goFunction: 'Tan', description: '求正切值', category: '数学', params: ['弧度'], returnType: '小数型' },
  { chineseName: '文本替换', goFunction: 'Replace', description: '替换文本内容', category: '文本', params: ['原文本', '被替换', '替换为'], returnType: '文本型' },
  { chineseName: '取文本中间', goFunction: 'Mid', description: '取文本中间部分', category: '文本', params: ['文本', '起始位置', '长度'], returnType: '文本型' },
  { chineseName: '取文本左边', goFunction: 'Left', description: '取文本左边部分', category: '文本', params: ['文本', '长度'], returnType: '文本型' },
  { chineseName: '取文本右边', goFunction: 'Right', description: '取文本右边部分', category: '文本', params: ['文本', '长度'], returnType: '文本型' },
  { chineseName: '分割文本', goFunction: 'Split', description: '分割文本为数组', category: '文本', params: ['文本', '分隔符'], returnType: '文本型数组' },
  { chineseName: '寻找文本', goFunction: 'FindStr', description: '寻找文本位置', category: '文本', params: ['被搜寻文本', '欲寻找文本', '起始位置'], returnType: '整数型' },
  { chineseName: '倒找文本', goFunction: 'RFindStr', description: '反向寻找文本位置', category: '文本', params: ['被搜寻文本', '欲寻找文本', '起始位置'], returnType: '整数型' },
  { chineseName: '到大写', goFunction: 'ToUpper', description: '转换到大写', category: '文本', params: ['文本'], returnType: '文本型' },
  { chineseName: '到小写', goFunction: 'ToLower', description: '转换到小写', category: '文本', params: ['文本'], returnType: '文本型' },
  { chineseName: '删首尾空', goFunction: 'Trim', description: '删除首尾空白', category: '文本', params: ['文本'], returnType: '文本型' },
  { chineseName: '删首空', goFunction: 'TrimLeft', description: '删除首部空白', category: '文本', params: ['文本'], returnType: '文本型' },
  { chineseName: '删尾空', goFunction: 'TrimRight', description: '删除尾部空白', category: '文本', params: ['文本'], returnType: '文本型' },
  { chineseName: '取重复文本', goFunction: 'Repeat', description: '取重复文本', category: '文本', params: ['重复次数', '待重复文本'], returnType: '文本型' },
  { chineseName: '创建窗口', goFunction: 'CreateForm', description: '创建窗口', category: '窗口', params: ['窗口模板'], returnType: '' },
  { chineseName: '载入窗口', goFunction: 'LoadForm', description: '载入窗口', category: '窗口', params: ['窗口变量', '是否模态'], returnType: '' },
  { chineseName: '销毁窗口', goFunction: 'DestroyForm', description: '销毁窗口', category: '窗口', params: [], returnType: '' },
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
  { chineseName: '创建定时器', goFunction: 'NewTimer', description: '创建定时器组件', category: '窗口', params: ['父窗口'], returnType: '定时器' },
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
  { chineseName: '创建字体对话框', goFunction: 'NewFontDialog', description: '创建字体选择对话框', category: '窗口', params: ['父窗口'], returnType: '对话框' },
  { chineseName: '创建颜色对话框', goFunction: 'NewColorDialog', description: '创建颜色选择对话框', category: '窗口', params: ['父窗口'], returnType: '对话框' },
  { chineseName: '创建超级列表框', goFunction: 'NewListView', description: '创建超级列表框', category: '窗口', params: ['父窗口'], returnType: '超级列表框' },
  { chineseName: '创建选择列表框', goFunction: 'NewCheckListBox', description: '创建选择列表框', category: '窗口', params: ['父窗口'], returnType: '选择列表框' },
  { chineseName: '创建超链接', goFunction: 'NewLinkLabel', description: '创建超链接标签', category: '窗口', params: ['父窗口'], returnType: '超链接' },
  { chineseName: '创建动画框', goFunction: 'NewAnimate', description: '创建动画框', category: '窗口', params: ['父窗口'], returnType: '动画框' },
  { chineseName: '创建图表', goFunction: 'NewChart', description: '创建图表组件', category: '窗口', params: ['父窗口'], returnType: '图表' },
  { chineseName: '创建热键框', goFunction: 'NewHotKey', description: '创建热键框', category: '窗口', params: ['父窗口'], returnType: '热键框' },
  { chineseName: '创建IP框', goFunction: 'NewIPEdit', description: '创建IP地址输入框', category: '窗口', params: ['父窗口'], returnType: 'IP框' },
  { chineseName: '创建遮罩框', goFunction: 'NewMaskEdit', description: '创建遮罩编辑框', category: '窗口', params: ['父窗口'], returnType: '遮罩框' },
  { chineseName: '创建媒体播放器', goFunction: 'NewMediaPlayer', description: '创建媒体播放器', category: '窗口', params: ['父窗口'], returnType: '媒体播放器' },
  { chineseName: '创建打印机', goFunction: 'NewPrinter', description: '创建打印机对象', category: '窗口', params: [], returnType: '打印机' },
  { chineseName: '创建打印对话框', goFunction: 'NewPrintDialog', description: '创建打印对话框', category: '窗口', params: ['父窗口'], returnType: '对话框' },
  { chineseName: '创建页面设置对话框', goFunction: 'NewPageSetupDialog', description: '创建页面设置对话框', category: '窗口', params: ['父窗口'], returnType: '对话框' },
  { chineseName: '创建查找对话框', goFunction: 'NewFindDialog', description: '创建查找对话框', category: '窗口', params: ['父窗口'], returnType: '对话框' },
  { chineseName: '创建替换对话框', goFunction: 'NewReplaceDialog', description: '创建替换对话框', category: '窗口', params: ['父窗口'], returnType: '对话框' },
  { chineseName: '创建托盘图标', goFunction: 'NewTrayIcon', description: '创建系统托盘图标', category: '窗口', params: ['父窗口'], returnType: '托盘图标' },
  { chineseName: '创建数据源', goFunction: 'NewDataSource', description: '创建数据源', category: '窗口', params: ['父窗口'], returnType: '数据源' },
  { chineseName: '创建客户端', goFunction: 'NewClientSocket', description: '创建客户端套接字', category: '窗口', params: ['父窗口'], returnType: '客户端' },
  { chineseName: '创建服务器', goFunction: 'NewServerSocket', description: '创建服务器套接字', category: '窗口', params: ['父窗口'], returnType: '服务器' },
  { chineseName: '创建数据报', goFunction: 'NewUDPSocket', description: '创建UDP套接字', category: '窗口', params: ['父窗口'], returnType: '数据报' },
  { chineseName: '创建串口', goFunction: 'NewSerialPort', description: '创建串口对象', category: '窗口', params: ['父窗口'], returnType: '串口' },
  { chineseName: '创建FTP', goFunction: 'NewFTP', description: '创建FTP对象', category: '窗口', params: ['父窗口'], returnType: 'FTP' },
  { chineseName: '创建HTTP', goFunction: 'NewHTTP', description: '创建HTTP对象', category: '窗口', params: ['父窗口'], returnType: 'HTTP' },
  { chineseName: '创建邮件', goFunction: 'NewMail', description: '创建邮件对象', category: '窗口', params: ['父窗口'], returnType: '邮件' },
  { chineseName: '创建XML解析器', goFunction: 'NewXMLParser', description: '创建XML解析器', category: '窗口', params: ['父窗口'], returnType: 'XML解析器' },
  { chineseName: '创建JSON解析器', goFunction: 'NewJSONParser', description: '创建JSON解析器', category: '窗口', params: ['父窗口'], returnType: 'JSON解析器' },
  { chineseName: '创建正则表达式', goFunction: 'NewRegex', description: '创建正则表达式对象', category: '窗口', params: ['父窗口'], returnType: '正则表达式' },
  { chineseName: '创建数据库连接', goFunction: 'NewDBConnection', description: '创建数据库连接', category: '窗口', params: ['父窗口'], returnType: '数据库连接' },
  { chineseName: '创建记录集', goFunction: 'NewRecordSet', description: '创建记录集', category: '窗口', params: ['父窗口'], returnType: '记录集' },
  { chineseName: '创建Excel对象', goFunction: 'NewExcel', description: '创建Excel操作对象', category: '窗口', params: ['父窗口'], returnType: 'Excel对象' },
  { chineseName: '创建Word对象', goFunction: 'NewWord', description: '创建Word操作对象', category: '窗口', params: ['父窗口'], returnType: 'Word对象' },
  { chineseName: '创建PDF对象', goFunction: 'NewPDF', description: '创建PDF操作对象', category: '窗口', params: ['父窗口'], returnType: 'PDF对象' },
  { chineseName: '创建压缩对象', goFunction: 'NewZip', description: '创建压缩解压对象', category: '窗口', params: ['父窗口'], returnType: '压缩对象' },
  { chineseName: '创建加密对象', goFunction: 'NewCrypto', description: '创建加密解密对象', category: '窗口', params: ['父窗口'], returnType: '加密对象' },
  { chineseName: '创建线程池', goFunction: 'NewThreadPool', description: '创建线程池', category: '窗口', params: ['父窗口'], returnType: '线程池' },
  { chineseName: '创建临界区', goFunction: 'NewCriticalSection', description: '创建临界区', category: '窗口', params: ['父窗口'], returnType: '临界区' },
  { chineseName: '创建互斥体', goFunction: 'NewMutex', description: '创建互斥体', category: '窗口', params: ['父窗口'], returnType: '互斥体' },
  { chineseName: '创建信号量', goFunction: 'NewSemaphore', description: '创建信号量', category: '窗口', params: ['父窗口'], returnType: '信号量' },
  { chineseName: '创建事件对象', goFunction: 'NewEvent', description: '创建事件对象', category: '窗口', params: ['父窗口'], returnType: '事件对象' },
  { chineseName: '创建内存映射', goFunction: 'NewMemoryMap', description: '创建内存映射文件', category: '窗口', params: ['父窗口'], returnType: '内存映射' },
  { chineseName: '创建管道', goFunction: 'NewPipe', description: '创建命名管道', category: '窗口', params: ['父窗口'], returnType: '管道' },
  { chineseName: '创建注册表', goFunction: 'NewRegistry', description: '创建注册表操作对象', category: '窗口', params: ['父窗口'], returnType: '注册表' },
  { chineseName: '创建INI文件', goFunction: 'NewINIFile', description: '创建INI文件操作对象', category: '窗口', params: ['父窗口'], returnType: 'INI文件' },
  { chineseName: '创建剪贴板', goFunction: 'NewClipboard', description: '创建剪贴板对象', category: '窗口', params: ['父窗口'], returnType: '剪贴板' },
  { chineseName: '创建屏幕截图', goFunction: 'NewScreenCapture', description: '创建屏幕截图对象', category: '窗口', params: ['父窗口'], returnType: '屏幕截图' },
  { chineseName: '创建语音对象', goFunction: 'NewSpeech', description: '创建语音合成对象', category: '窗口', params: ['父窗口'], returnType: '语音对象' },
  { chineseName: '创建摄像头', goFunction: 'NewCamera', description: '创建摄像头对象', category: '窗口', params: ['父窗口'], returnType: '摄像头' },
  { chineseName: '创建条形码', goFunction: 'NewBarcode', description: '创建条形码对象', category: '窗口', params: ['父窗口'], returnType: '条形码' },
  { chineseName: '创建二维码', goFunction: 'NewQRCode', description: '创建二维码对象', category: '窗口', params: ['父窗口'], returnType: '二维码' },
  { chineseName: '置窗口标题', goFunction: 'SetFormTitle', description: '设置窗口标题', category: '窗口', params: ['窗口', '标题'], returnType: '' },
  { chineseName: '置窗口大小', goFunction: 'SetFormSize', description: '设置窗口大小', category: '窗口', params: ['窗口', '宽度', '高度'], returnType: '' },
  { chineseName: '置窗口位置', goFunction: 'SetFormPosition', description: '设置窗口位置', category: '窗口', params: ['窗口', '左边', '顶边'], returnType: '' },
  { chineseName: '置组件文本', goFunction: 'SetComponentText', description: '设置组件文本', category: '窗口', params: ['组件', '文本'], returnType: '' },
  { chineseName: '置组件位置', goFunction: 'SetComponentPosition', description: '设置组件位置', category: '窗口', params: ['组件', '左边', '顶边'], returnType: '' },
  { chineseName: '置组件大小', goFunction: 'SetComponentSize', description: '设置组件大小', category: '窗口', params: ['组件', '宽度', '高度'], returnType: '' },
  { chineseName: '置组件可视', goFunction: 'SetComponentVisible', description: '设置组件是否可见', category: '窗口', params: ['组件', '是否可见'], returnType: '' },
  { chineseName: '置组件可用', goFunction: 'SetComponentEnabled', description: '设置组件是否可用', category: '窗口', params: ['组件', '是否可用'], returnType: '' },
  { chineseName: '置组件颜色', goFunction: 'SetComponentColor', description: '设置组件颜色', category: '窗口', params: ['组件', '颜色值'], returnType: '' },
  { chineseName: '置组件字体', goFunction: 'SetComponentFont', description: '设置组件字体', category: '窗口', params: ['组件', '字体名', '大小', '是否粗体'], returnType: '' },
  { chineseName: '取组件文本', goFunction: 'GetComponentText', description: '获取组件文本', category: '窗口', params: ['组件'], returnType: '文本型' },
  { chineseName: '取组件宽度', goFunction: 'GetComponentWidth', description: '获取组件宽度', category: '窗口', params: ['组件'], returnType: '整数型' },
  { chineseName: '取组件高度', goFunction: 'GetComponentHeight', description: '获取组件高度', category: '窗口', params: ['组件'], returnType: '整数型' },
  { chineseName: '置进度条位置', goFunction: 'SetProgressBarPosition', description: '设置进度条位置', category: '窗口', params: ['进度条', '位置'], returnType: '' },
  { chineseName: '置滑块位置', goFunction: 'SetTrackBarPosition', description: '设置滑块位置', category: '窗口', params: ['滑块', '位置'], returnType: '' },
  { chineseName: '置调节器值', goFunction: 'SetSpinEditValue', description: '设置调节器值', category: '窗口', params: ['调节器', '值'], returnType: '' },
  { chineseName: '取调节器值', goFunction: 'GetSpinEditValue', description: '获取调节器值', category: '窗口', params: ['调节器'], returnType: '整数型' },
  { chineseName: '置组合框项目', goFunction: 'SetComboBoxItems', description: '设置组合框项目列表', category: '窗口', params: ['组合框', '项目列表'], returnType: '' },
  { chineseName: '置列表框项目', goFunction: 'SetListBoxItems', description: '设置列表框项目列表', category: '窗口', params: ['列表框', '项目列表'], returnType: '' },
  { chineseName: '取组合框选中项', goFunction: 'GetComboBoxSelected', description: '获取组合框选中项', category: '窗口', params: ['组合框'], returnType: '文本型' },
  { chineseName: '取列表框选中项', goFunction: 'GetListBoxSelected', description: '获取列表框选中项', category: '窗口', params: ['列表框'], returnType: '文本型' },
  { chineseName: '置表格单元格', goFunction: 'SetGridCell', description: '设置表格单元格文本', category: '窗口', params: ['表格', '行', '列', '文本'], returnType: '' },
  { chineseName: '取表格单元格', goFunction: 'GetGridCell', description: '获取表格单元格文本', category: '窗口', params: ['表格', '行', '列'], returnType: '文本型' },
  { chineseName: '置表格行数', goFunction: 'SetGridRowCount', description: '设置表格行数', category: '窗口', params: ['表格', '行数'], returnType: '' },
  { chineseName: '置表格列数', goFunction: 'SetGridColCount', description: '设置表格列数', category: '窗口', params: ['表格', '列数'], returnType: '' },
  { chineseName: '置树形框节点', goFunction: 'SetTreeNode', description: '设置树形框节点', category: '窗口', params: ['树形框', '节点文本', '父节点'], returnType: '' },
  { chineseName: '取树形框选中项', goFunction: 'GetTreeSelected', description: '获取树形框选中项', category: '窗口', params: ['树形框'], returnType: '文本型' },
  { chineseName: '置选项卡标题', goFunction: 'SetPageControlTitle', description: '设置选项卡页面标题', category: '窗口', params: ['选项卡', '页面索引', '标题'], returnType: '' },
  { chineseName: '置图片框图片', goFunction: 'SetImagePicture', description: '设置图片框图片', category: '窗口', params: ['图片框', '图片路径'], returnType: '' },
  { chineseName: '置浏览器URL', goFunction: 'SetBrowserURL', description: '设置浏览器URL', category: '窗口', params: ['浏览器', 'URL'], returnType: '' },
  { chineseName: '置定时器周期', goFunction: 'SetTimerInterval', description: '设置定时器周期', category: '窗口', params: ['定时器', '毫秒'], returnType: '' },
  { chineseName: '启动定时器', goFunction: 'StartTimer', description: '启动定时器', category: '窗口', params: ['定时器'], returnType: '' },
  { chineseName: '停止定时器', goFunction: 'StopTimer', description: '停止定时器', category: '窗口', params: ['定时器'], returnType: '' },
  { chineseName: '置菜单项', goFunction: 'SetMenuItem', description: '设置菜单项', category: '窗口', params: ['菜单', '标题', '事件名'], returnType: '' },
  { chineseName: '置工具栏按钮', goFunction: 'SetToolBarButton', description: '设置工具栏按钮', category: '窗口', params: ['工具栏', '标题', '图标', '事件名'], returnType: '' },
  { chineseName: '置状态栏文本', goFunction: 'SetStatusBarText', description: '设置状态栏文本', category: '窗口', params: ['状态栏', '文本'], returnType: '' },
  { chineseName: '置选择框状态', goFunction: 'SetCheckBoxState', description: '设置选择框选中状态', category: '窗口', params: ['选择框', '是否选中'], returnType: '' },
  { chineseName: '取选择框状态', goFunction: 'GetCheckBoxState', description: '获取选择框选中状态', category: '窗口', params: ['选择框'], returnType: '逻辑型' },
  { chineseName: '置单选框状态', goFunction: 'SetRadioButtonState', description: '设置单选框选中状态', category: '窗口', params: ['单选框', '是否选中'], returnType: '' },
  { chineseName: '取单选框状态', goFunction: 'GetRadioButtonState', description: '获取单选框选中状态', category: '窗口', params: ['单选框'], returnType: '逻辑型' },
  { chineseName: '置日期框日期', goFunction: 'SetDateTimePickerDate', description: '设置日期框日期', category: '窗口', params: ['日期框', '日期'], returnType: '' },
  { chineseName: '取日期框日期', goFunction: 'GetDateTimePickerDate', description: '获取日期框日期', category: '窗口', params: ['日期框'], returnType: '日期时间型' },
  { chineseName: '置颜色框颜色', goFunction: 'SetColorBoxColor', description: '设置颜色框颜色', category: '窗口', params: ['颜色框', '颜色值'], returnType: '' },
  { chineseName: '取颜色框颜色', goFunction: 'GetColorBoxColor', description: '获取颜色框颜色', category: '窗口', params: ['颜色框'], returnType: '整数型' },
  { chineseName: '置外形框形状', goFunction: 'SetShapeType', description: '设置外形框形状', category: '窗口', params: ['外形框', '形状类型'], returnType: '' },
  { chineseName: '置外形框画笔', goFunction: 'SetShapePen', description: '设置外形框画笔颜色和宽度', category: '窗口', params: ['外形框', '颜色', '宽度'], returnType: '' },
  { chineseName: '置外形框画刷', goFunction: 'SetShapeBrush', description: '设置外形框画刷颜色', category: '窗口', params: ['外形框', '颜色'], returnType: '' },
  { chineseName: '置滚动条范围', goFunction: 'SetScrollBarRange', description: '设置滚动条范围', category: '窗口', params: ['滚动条', '最小值', '最大值'], returnType: '' },
  { chineseName: '置滚动条位置', goFunction: 'SetScrollBarPosition', description: '设置滚动条位置', category: '窗口', params: ['滚动条', '位置'], returnType: '' },
  { chineseName: '取滚动条位置', goFunction: 'GetScrollBarPosition', description: '获取滚动条位置', category: '窗口', params: ['滚动条'], returnType: '整数型' },
  { chineseName: '置超级列表框列', goFunction: 'SetListViewColumn', description: '设置超级列表框列标题', category: '窗口', params: ['超级列表框', '列索引', '标题', '宽度'], returnType: '' },
  { chineseName: '置超级列表框项目', goFunction: 'SetListViewItem', description: '设置超级列表框项目', category: '窗口', params: ['超级列表框', '行', '列', '文本'], returnType: '' },
  { chineseName: '取超级列表框选中项', goFunction: 'GetListViewSelected', description: '获取超级列表框选中项索引', category: '窗口', params: ['超级列表框'], returnType: '整数型' },
  { chineseName: '置超链接URL', goFunction: 'SetLinkLabelURL', description: '设置超链接URL', category: '窗口', params: ['超链接', 'URL'], returnType: '' },
  { chineseName: '置动画框文件', goFunction: 'SetAnimateFile', description: '设置动画框文件', category: '窗口', params: ['动画框', '文件路径'], returnType: '' },
  { chineseName: '播放动画', goFunction: 'PlayAnimate', description: '播放动画', category: '窗口', params: ['动画框'], returnType: '' },
  { chineseName: '停止动画', goFunction: 'StopAnimate', description: '停止动画', category: '窗口', params: ['动画框'], returnType: '' },
  { chineseName: '置图表数据', goFunction: 'SetChartData', description: '设置图表数据', category: '窗口', params: ['图表', '数据'], returnType: '' },
  { chineseName: '置图表类型', goFunction: 'SetChartType', description: '设置图表类型', category: '窗口', params: ['图表', '类型'], returnType: '' },
  { chineseName: '置热键框热键', goFunction: 'SetHotKeyValue', description: '设置热键框热键值', category: '窗口', params: ['热键框', '键值'], returnType: '' },
  { chineseName: '取热键框热键', goFunction: 'GetHotKeyValue', description: '获取热键框热键值', category: '窗口', params: ['热键框'], returnType: '整数型' },
  { chineseName: '置IP框地址', goFunction: 'SetIPEditAddress', description: '设置IP框地址', category: '窗口', params: ['IP框', 'IP地址'], returnType: '' },
  { chineseName: '取IP框地址', goFunction: 'GetIPEditAddress', description: '获取IP框地址', category: '窗口', params: ['IP框'], returnType: '文本型' },
  { chineseName: '置遮罩框掩码', goFunction: 'SetMaskEditMask', description: '设置遮罩框掩码', category: '窗口', params: ['遮罩框', '掩码'], returnType: '' },
  { chineseName: '置媒体播放器文件', goFunction: 'SetMediaPlayerFile', description: '设置媒体播放器文件', category: '窗口', params: ['媒体播放器', '文件路径'], returnType: '' },
  { chineseName: '播放媒体', goFunction: 'PlayMedia', description: '播放媒体', category: '窗口', params: ['媒体播放器'], returnType: '' },
  { chineseName: '暂停媒体', goFunction: 'PauseMedia', description: '暂停媒体', category: '窗口', params: ['媒体播放器'], returnType: '' },
  { chineseName: '停止媒体', goFunction: 'StopMedia', description: '停止媒体', category: '窗口', params: ['媒体播放器'], returnType: '' },
  { chineseName: '置托盘图标提示', goFunction: 'SetTrayIconHint', description: '设置托盘图标提示文本', category: '窗口', params: ['托盘图标', '提示文本'], returnType: '' },
  { chineseName: '置托盘图标图片', goFunction: 'SetTrayIconPicture', description: '设置托盘图标图片', category: '窗口', params: ['托盘图标', '图片路径'], returnType: '' },
  { chineseName: '显示托盘气泡', goFunction: 'ShowTrayBalloon', description: '显示托盘气泡提示', category: '窗口', params: ['托盘图标', '标题', '内容', '图标类型'], returnType: '' },
  { chineseName: '连接数据库', goFunction: 'ConnectDB', description: '连接数据库', category: '窗口', params: ['数据库连接', '连接字符串'], returnType: '逻辑型' },
  { chineseName: '断开数据库', goFunction: 'DisconnectDB', description: '断开数据库连接', category: '窗口', params: ['数据库连接'], returnType: '' },
  { chineseName: '执行SQL', goFunction: 'ExecuteSQL', description: '执行SQL语句', category: '窗口', params: ['数据库连接', 'SQL语句'], returnType: '逻辑型' },
  { chineseName: '取记录集行数', goFunction: 'GetRecordSetRowCount', description: '获取记录集行数', category: '窗口', params: ['记录集'], returnType: '整数型' },
  { chineseName: '取记录集字段', goFunction: 'GetRecordSetField', description: '获取记录集字段值', category: '窗口', params: ['记录集', '行', '字段名'], returnType: '文本型' },
  { chineseName: '打开Excel', goFunction: 'OpenExcel', description: '打开Excel文件', category: '窗口', params: ['Excel对象', '文件路径'], returnType: '逻辑型' },
  { chineseName: '置Excel单元格', goFunction: 'SetExcelCell', description: '设置Excel单元格值', category: '窗口', params: ['Excel对象', '行', '列', '值'], returnType: '' },
  { chineseName: '取Excel单元格', goFunction: 'GetExcelCell', description: '获取Excel单元格值', category: '窗口', params: ['Excel对象', '行', '列'], returnType: '文本型' },
  { chineseName: '保存Excel', goFunction: 'SaveExcel', description: '保存Excel文件', category: '窗口', params: ['Excel对象'], returnType: '逻辑型' },
  { chineseName: '关闭Excel', goFunction: 'CloseExcel', description: '关闭Excel对象', category: '窗口', params: ['Excel对象'], returnType: '' },
  { chineseName: '压缩文件', goFunction: 'ZipFiles', description: '压缩文件', category: '窗口', params: ['压缩对象', '文件列表', '目标路径'], returnType: '逻辑型' },
  { chineseName: '解压文件', goFunction: 'UnzipFiles', description: '解压文件', category: '窗口', params: ['压缩对象', '压缩包路径', '目标目录'], returnType: '逻辑型' },
  { chineseName: 'MD5加密', goFunction: 'MD5Encrypt', description: 'MD5加密', category: '窗口', params: ['加密对象', '文本'], returnType: '文本型' },
  { chineseName: 'SHA256加密', goFunction: 'SHA256Encrypt', description: 'SHA256加密', category: '窗口', params: ['加密对象', '文本'], returnType: '文本型' },
  { chineseName: 'Base64编码', goFunction: 'Base64Encode', description: 'Base64编码', category: '窗口', params: ['加密对象', '数据'], returnType: '文本型' },
  { chineseName: 'Base64解码', goFunction: 'Base64Decode', description: 'Base64解码', category: '窗口', params: ['加密对象', '文本'], returnType: '文本型' },
  { chineseName: 'AES加密', goFunction: 'AESEncrypt', description: 'AES加密', category: '窗口', params: ['加密对象', '文本', '密钥'], returnType: '文本型' },
  { chineseName: 'AES解密', goFunction: 'AESDecrypt', description: 'AES解密', category: '窗口', params: ['加密对象', '文本', '密钥'], returnType: '文本型' },
  { chineseName: '提交线程任务', goFunction: 'SubmitThreadTask', description: '提交线程任务', category: '窗口', params: ['线程池', '任务函数'], returnType: '' },
  { chineseName: '等待线程任务', goFunction: 'WaitThreadTasks', description: '等待所有线程任务完成', category: '窗口', params: ['线程池'], returnType: '' },
  { chineseName: '进入临界区', goFunction: 'EnterCriticalSection', description: '进入临界区', category: '窗口', params: ['临界区'], returnType: '' },
  { chineseName: '离开临界区', goFunction: 'LeaveCriticalSection', description: '离开临界区', category: '窗口', params: ['临界区'], returnType: '' },
  { chineseName: '锁定互斥体', goFunction: 'LockMutex', description: '锁定互斥体', category: '窗口', params: ['互斥体'], returnType: '' },
  { chineseName: '解锁互斥体', goFunction: 'UnlockMutex', description: '解锁互斥体', category: '窗口', params: ['互斥体'], returnType: '' },
  { chineseName: '等待信号量', goFunction: 'WaitSemaphore', description: '等待信号量', category: '窗口', params: ['信号量'], returnType: '' },
  { chineseName: '释放信号量', goFunction: 'ReleaseSemaphore', description: '释放信号量', category: '窗口', params: ['信号量'], returnType: '' },
  { chineseName: '置事件信号', goFunction: 'SetEventSignal', description: '设置事件信号', category: '窗口', params: ['事件对象'], returnType: '' },
  { chineseName: '等待事件信号', goFunction: 'WaitEventSignal', description: '等待事件信号', category: '窗口', params: ['事件对象', '超时毫秒'], returnType: '逻辑型' },
  { chineseName: '重置事件信号', goFunction: 'ResetEventSignal', description: '重置事件信号', category: '窗口', params: ['事件对象'], returnType: '' },
  { chineseName: '写注册表', goFunction: 'WriteRegistry', description: '写入注册表值', category: '窗口', params: ['注册表', '路径', '键名', '值'], returnType: '逻辑型' },
  { chineseName: '读注册表', goFunction: 'ReadRegistry', description: '读取注册表值', category: '窗口', params: ['注册表', '路径', '键名'], returnType: '文本型' },
  { chineseName: '删除注册表', goFunction: 'DeleteRegistry', description: '删除注册表键', category: '窗口', params: ['注册表', '路径', '键名'], returnType: '逻辑型' },
  { chineseName: '写INI值', goFunction: 'WriteINIValue', description: '写入INI配置值', category: '窗口', params: ['INI文件', '节名', '键名', '值'], returnType: '逻辑型' },
  { chineseName: '读INI值', goFunction: 'ReadINIValue', description: '读取INI配置值', category: '窗口', params: ['INI文件', '节名', '键名'], returnType: '文本型' },
  { chineseName: '置剪贴板文本', goFunction: 'SetClipboardText', description: '设置剪贴板文本', category: '窗口', params: ['剪贴板', '文本'], returnType: '' },
  { chineseName: '取剪贴板文本', goFunction: 'GetClipboardText', description: '获取剪贴板文本', category: '窗口', params: ['剪贴板'], returnType: '文本型' },
  { chineseName: '屏幕截图', goFunction: 'CaptureScreen', description: '截取屏幕图像', category: '窗口', params: ['屏幕截图', '左边', '顶边', '宽度', '高度'], returnType: '' },
  { chineseName: '保存截图', goFunction: 'SaveScreenCapture', description: '保存截图到文件', category: '窗口', params: ['屏幕截图', '文件路径'], returnType: '逻辑型' },
  { chineseName: '语音合成', goFunction: 'SpeechSynthesize', description: '语音合成文本', category: '窗口', params: ['语音对象', '文本'], returnType: '' },
  { chineseName: '打开摄像头', goFunction: 'OpenCamera', description: '打开摄像头', category: '窗口', params: ['摄像头'], returnType: '逻辑型' },
  { chineseName: '拍照', goFunction: 'CapturePhoto', description: '拍摄照片', category: '窗口', params: ['摄像头', '保存路径'], returnType: '逻辑型' },
  { chineseName: '关闭摄像头', goFunction: 'CloseCamera', description: '关闭摄像头', category: '窗口', params: ['摄像头'], returnType: '' },
  { chineseName: '生成条形码', goFunction: 'GenerateBarcode', description: '生成条形码图片', category: '窗口', params: ['条形码', '内容', '保存路径'], returnType: '逻辑型' },
  { chineseName: '生成二维码', goFunction: 'GenerateQRCode', description: '生成二维码图片', category: '窗口', params: ['二维码', '内容', '保存路径'], returnType: '逻辑型' },
  { chineseName: '取文件大小', goFunction: 'GetFileSize', description: '获取文件大小', category: '文件', params: ['文件路径'], returnType: '整数型' },
  { chineseName: '取文件时间', goFunction: 'GetFileTime', description: '获取文件修改时间', category: '文件', params: ['文件路径'], returnType: '日期时间型' },
  { chineseName: '复制文件', goFunction: 'CopyFile', description: '复制文件', category: '文件', params: ['源文件', '目标文件'], returnType: '逻辑型' },
  { chineseName: '移动文件', goFunction: 'MoveFile', description: '移动文件', category: '文件', params: ['源文件', '目标文件'], returnType: '逻辑型' },
  { chineseName: '取文件扩展名', goFunction: 'GetFileExt', description: '获取文件扩展名', category: '文件', params: ['文件路径'], returnType: '文本型' },
  { chineseName: '取文件名', goFunction: 'GetFileName', description: '获取文件名（不含路径）', category: '文件', params: ['文件路径'], returnType: '文本型' },
  { chineseName: '取文件目录', goFunction: 'GetFileDir', description: '获取文件所在目录', category: '文件', params: ['文件路径'], returnType: '文本型' },
  { chineseName: '枚举文件', goFunction: 'EnumFiles', description: '枚举目录下文件', category: '文件', params: ['目录路径', '通配符'], returnType: '文本型数组' },
  { chineseName: '枚举子目录', goFunction: 'EnumSubDirs', description: '枚举子目录', category: '文件', params: ['目录路径'], returnType: '文本型数组' },
  { chineseName: '取临时目录', goFunction: 'GetTempDir', description: '获取临时目录路径', category: '文件', params: [], returnType: '文本型' },
  { chineseName: '取特殊目录', goFunction: 'GetSpecialDir', description: '获取特殊目录路径', category: '文件', params: ['目录类型'], returnType: '文本型' },
  { chineseName: '取环境变量', goFunction: 'GetEnv', description: '获取环境变量值', category: '系统', params: ['变量名'], returnType: '文本型' },
  { chineseName: '置环境变量', goFunction: 'SetEnv', description: '设置环境变量值', category: '系统', params: ['变量名', '值'], returnType: '逻辑型' },
  { chineseName: '取命令行', goFunction: 'GetCmdLine', description: '获取命令行参数', category: '系统', params: [], returnType: '文本型数组' },
  { chineseName: '取CPU核心数', goFunction: 'GetCPUCoreCount', description: '获取CPU核心数', category: '系统', params: [], returnType: '整数型' },
  { chineseName: '取内存信息', goFunction: 'GetMemoryInfo', description: '获取内存信息', category: '系统', params: [], returnType: '文本型' },
  { chineseName: '取磁盘信息', goFunction: 'GetDiskInfo', description: '获取磁盘信息', category: '系统', params: ['盘符'], returnType: '文本型' },
  { chineseName: '取系统版本', goFunction: 'GetOSVersion', description: '获取操作系统版本', category: '系统', params: [], returnType: '文本型' },
  { chineseName: '取计算机名', goFunction: 'GetComputerName', description: '获取计算机名称', category: '系统', params: [], returnType: '文本型' },
  { chineseName: '取用户名', goFunction: 'GetUserName', description: '获取当前用户名', category: '系统', params: [], returnType: '文本型' },
  { chineseName: '取屏幕宽度', goFunction: 'GetScreenWidth', description: '获取屏幕宽度', category: '系统', params: [], returnType: '整数型' },
  { chineseName: '取屏幕高度', goFunction: 'GetScreenHeight', description: '获取屏幕高度', category: '系统', params: [], returnType: '整数型' },
  { chineseName: '播放声音', goFunction: 'PlaySound', description: '播放系统声音', category: '系统', params: ['声音类型'], returnType: '' },
  { chineseName: '关闭系统', goFunction: 'ShutdownSystem', description: '关闭/重启系统', category: '系统', params: ['操作类型'], returnType: '逻辑型' },
  { chineseName: '取时间戳', goFunction: 'GetTimestamp', description: '获取Unix时间戳', category: '时间', params: [], returnType: '整数型' },
  { chineseName: '时间戳到时间', goFunction: 'TimestampToTime', description: '时间戳转日期时间', category: '时间', params: ['时间戳'], returnType: '日期时间型' },
  { chineseName: '时间到时间戳', goFunction: 'TimeToTimestamp', description: '日期时间转时间戳', category: '时间', params: ['日期时间'], returnType: '整数型' },
  { chineseName: '格式化时间', goFunction: 'FormatTime', description: '格式化日期时间', category: '时间', params: ['日期时间', '格式'], returnType: '文本型' },
  { chineseName: '取年份', goFunction: 'GetYear', description: '取年份', category: '时间', params: ['日期时间'], returnType: '整数型' },
  { chineseName: '取月份', goFunction: 'GetMonth', description: '取月份', category: '时间', params: ['日期时间'], returnType: '整数型' },
  { chineseName: '取日', goFunction: 'GetDay', description: '取日', category: '时间', params: ['日期时间'], returnType: '整数型' },
  { chineseName: '取小时', goFunction: 'GetHour', description: '取小时', category: '时间', params: ['日期时间'], returnType: '整数型' },
  { chineseName: '取分钟', goFunction: 'GetMinute', description: '取分钟', category: '时间', params: ['日期时间'], returnType: '整数型' },
  { chineseName: '取秒', goFunction: 'GetSecond', description: '取秒', category: '时间', params: ['日期时间'], returnType: '整数型' },
  { chineseName: '取星期几', goFunction: 'GetWeekday', description: '取星期几', category: '时间', params: ['日期时间'], returnType: '整数型' },
  { chineseName: '增减时间', goFunction: 'AddTime', description: '增减时间', category: '时间', params: ['日期时间', '单位', '数量'], returnType: '日期时间型' },
  { chineseName: '取时间间隔', goFunction: 'GetTimeInterval', description: '取两个时间间隔', category: '时间', params: ['时间1', '时间2', '单位'], returnType: '整数型' },
  { chineseName: 'HTTP请求', goFunction: 'HTTPRequest', description: '发送HTTP请求', category: '网络', params: ['URL', '方法', '数据', '超时'], returnType: '文本型' },
  { chineseName: 'HTTP下载', goFunction: 'HTTPDownload', description: 'HTTP下载文件', category: '网络', params: ['URL', '保存路径', '进度回调'], returnType: '逻辑型' },
  { chineseName: 'TCP连接', goFunction: 'TCPConnect', description: 'TCP连接服务器', category: '网络', params: ['客户端', '地址', '端口'], returnType: '逻辑型' },
  { chineseName: 'TCP发送', goFunction: 'TCPSend', description: 'TCP发送数据', category: '网络', params: ['客户端', '数据'], returnType: '逻辑型' },
  { chineseName: 'TCP接收', goFunction: 'TCPReceive', description: 'TCP接收数据', category: '网络', params: ['客户端', '超时'], returnType: '字节集' },
  { chineseName: 'TCP断开', goFunction: 'TCPDisconnect', description: 'TCP断开连接', category: '网络', params: ['客户端'], returnType: '' },
  { chineseName: 'TCP监听', goFunction: 'TCPListen', description: 'TCP开始监听', category: '网络', params: ['服务器', '端口'], returnType: '逻辑型' },
  { chineseName: 'TCP停止监听', goFunction: 'TCPStopListen', description: 'TCP停止监听', category: '网络', params: ['服务器'], returnType: '' },
  { chineseName: 'UDP发送', goFunction: 'UDPSend', description: 'UDP发送数据', category: '网络', params: ['数据报', '地址', '端口', '数据'], returnType: '逻辑型' },
  { chineseName: 'UDP接收', goFunction: 'UDPReceive', description: 'UDP接收数据', category: '网络', params: ['数据报', '超时'], returnType: '字节集' },
  { chineseName: '串口打开', goFunction: 'SerialOpen', description: '打开串口', category: '网络', params: ['串口', '端口号', '波特率'], returnType: '逻辑型' },
  { chineseName: '串口发送', goFunction: 'SerialSend', description: '串口发送数据', category: '网络', params: ['串口', '数据'], returnType: '逻辑型' },
  { chineseName: '串口接收', goFunction: 'SerialReceive', description: '串口接收数据', category: '网络', params: ['串口', '超时'], returnType: '字节集' },
  { chineseName: '串口关闭', goFunction: 'SerialClose', description: '关闭串口', category: '网络', params: ['串口'], returnType: '' },
  { chineseName: 'FTP连接', goFunction: 'FTPConnect', description: 'FTP连接服务器', category: '网络', params: ['FTP', '地址', '端口', '用户名', '密码'], returnType: '逻辑型' },
  { chineseName: 'FTP上传', goFunction: 'FTPUpload', description: 'FTP上传文件', category: '网络', params: ['FTP', '本地文件', '远程路径'], returnType: '逻辑型' },
  { chineseName: 'FTP下载', goFunction: 'FTPDownload', description: 'FTP下载文件', category: '网络', params: ['FTP', '远程文件', '本地路径'], returnType: '逻辑型' },
  { chineseName: 'FTP断开', goFunction: 'FTPDisconnect', description: 'FTP断开连接', category: '网络', params: ['FTP'], returnType: '' },
  { chineseName: '发送邮件', goFunction: 'SendMail', description: '发送邮件', category: '网络', params: ['邮件', '收件人', '主题', '内容'], returnType: '逻辑型' },
  { chineseName: 'JSON解析', goFunction: 'JSONParse', description: '解析JSON文本', category: '网络', params: ['JSON解析器', 'JSON文本'], returnType: '逻辑型' },
  { chineseName: 'JSON取值', goFunction: 'JSONGetValue', description: '取JSON节点值', category: '网络', params: ['JSON解析器', '路径'], returnType: '文本型' },
  { chineseName: 'JSON置值', goFunction: 'JSONSetValue', description: '设置JSON节点值', category: '网络', params: ['JSON解析器', '路径', '值'], returnType: '' },
  { chineseName: 'JSON到文本', goFunction: 'JSONToString', description: 'JSON对象转文本', category: '网络', params: ['JSON解析器'], returnType: '文本型' },
  { chineseName: 'XML解析', goFunction: 'XMLParse', description: '解析XML文本', category: '网络', params: ['XML解析器', 'XML文本'], returnType: '逻辑型' },
  { chineseName: 'XML取值', goFunction: 'XMLGetValue', description: '取XML节点值', category: '网络', params: ['XML解析器', '路径'], returnType: '文本型' },
  { chineseName: '正则匹配', goFunction: 'RegexMatch', description: '正则表达式匹配', category: '网络', params: ['正则表达式', '文本', '模式'], returnType: '文本型数组' },
  { chineseName: '正则替换', goFunction: 'RegexReplace', description: '正则表达式替换', category: '网络', params: ['正则表达式', '文本', '替换为'], returnType: '文本型' },
  { chineseName: '取反', goFunction: 'Not', description: '逻辑取反', category: '逻辑', params: ['逻辑值'], returnType: '逻辑型' },
  { chineseName: '且', goFunction: 'And', description: '逻辑与运算', category: '逻辑', params: ['值1', '值2'], returnType: '逻辑型' },
  { chineseName: '或', goFunction: 'Or', description: '逻辑或运算', category: '逻辑', params: ['值1', '值2'], returnType: '逻辑型' },
  { chineseName: '位与', goFunction: 'BitAnd', description: '位与运算', category: '逻辑', params: ['值1', '值2'], returnType: '整数型' },
  { chineseName: '位或', goFunction: 'BitOr', description: '位或运算', category: '逻辑', params: ['值1', '值2'], returnType: '整数型' },
  { chineseName: '位取反', goFunction: 'BitNot', description: '位取反运算', category: '逻辑', params: ['值'], returnType: '整数型' },
  { chineseName: '左移', goFunction: 'BitLeftShift', description: '位左移运算', category: '逻辑', params: ['值', '位数'], returnType: '整数型' },
  { chineseName: '右移', goFunction: 'BitRightShift', description: '位右移运算', category: '逻辑', params: ['值', '位数'], returnType: '整数型' },
  { chineseName: '取数组成员数', goFunction: 'ArrayLen', description: '取数组长度', category: '数组', params: ['数组'], returnType: '整数型' },
  { chineseName: '加入成员', goFunction: 'ArrayAppend', description: '数组加入成员', category: '数组', params: ['数组', '成员'], returnType: '' },
  { chineseName: '删除成员', goFunction: 'ArrayRemove', description: '数组删除成员', category: '数组', params: ['数组', '索引'], returnType: '' },
  { chineseName: '插入成员', goFunction: 'ArrayInsert', description: '数组插入成员', category: '数组', params: ['数组', '索引', '成员'], returnType: '' },
  { chineseName: '清除数组', goFunction: 'ArrayClear', description: '清除数组所有成员', category: '数组', params: ['数组'], returnType: '' },
  { chineseName: '数组排序', goFunction: 'ArraySort', description: '数组排序', category: '数组', params: ['数组', '是否升序'], returnType: '' },
  { chineseName: '数组查找', goFunction: 'ArrayFind', description: '数组查找成员', category: '数组', params: ['数组', '成员'], returnType: '整数型' },
  { chineseName: '数组反转', goFunction: 'ArrayReverse', description: '数组反转', category: '数组', params: ['数组'], returnType: '' },
  { chineseName: '取颜色值', goFunction: 'RGB', description: '取RGB颜色值', category: '图形', params: ['红', '绿', '蓝'], returnType: '整数型' },
  { chineseName: '取红色值', goFunction: 'GetRValue', description: '取颜色红色分量', category: '图形', params: ['颜色值'], returnType: '整数型' },
  { chineseName: '取绿色值', goFunction: 'GetGValue', description: '取颜色绿色分量', category: '图形', params: ['颜色值'], returnType: '整数型' },
  { chineseName: '取蓝色值', goFunction: 'GetBValue', description: '取颜色蓝色分量', category: '图形', params: ['颜色值'], returnType: '整数型' },
  { chineseName: '取图片宽度', goFunction: 'GetImageWidth', description: '获取图片宽度', category: '图形', params: ['图片路径'], returnType: '整数型' },
  { chineseName: '取图片高度', goFunction: 'GetImageHeight', description: '获取图片高度', category: '图形', params: ['图片路径'], returnType: '整数型' },
  { chineseName: '缩放图片', goFunction: 'ResizeImage', description: '缩放图片', category: '图形', params: ['源文件', '目标文件', '宽度', '高度'], returnType: '逻辑型' },
  { chineseName: '旋转图片', goFunction: 'RotateImage', description: '旋转图片', category: '图形', params: ['源文件', '目标文件', '角度'], returnType: '逻辑型' },
  { chineseName: '裁剪图片', goFunction: 'CropImage', description: '裁剪图片', category: '图形', params: ['源文件', '目标文件', '左', '顶', '宽', '高'], returnType: '逻辑型' },
  { chineseName: '转换图片格式', goFunction: 'ConvertImageFormat', description: '转换图片格式', category: '图形', params: ['源文件', '目标文件', '格式'], returnType: '逻辑型' },
  { chineseName: '取文本宽度', goFunction: 'GetTextWidth', description: '获取文本渲染宽度', category: '图形', params: ['文本', '字体名', '字体大小'], returnType: '整数型' },
  { chineseName: '取文本高度', goFunction: 'GetTextHeight', description: '获取文本渲染高度', category: '图形', params: ['文本', '字体名', '字体大小'], returnType: '整数型' },
  { chineseName: '置随机数种子', goFunction: 'SetRandomSeed', description: '设置随机数种子', category: '数学', params: ['种子'], returnType: '' },
  { chineseName: '取整', goFunction: 'Floor', description: '向下取整', category: '数学', params: ['数值'], returnType: '整数型' },
  { chineseName: '取上限整', goFunction: 'Ceil', description: '向上取整', category: '数学', params: ['数值'], returnType: '整数型' },
  { chineseName: '取最大值', goFunction: 'Max', description: '取最大值', category: '数学', params: ['值1', '值2'], returnType: '小数型' },
  { chineseName: '取最小值', goFunction: 'Min', description: '取最小值', category: '数学', params: ['值1', '值2'], returnType: '小数型' },
  { chineseName: '取符号', goFunction: 'Sign', description: '取数值符号', category: '数学', params: ['数值'], returnType: '整数型' },
  { chineseName: '求自然对数', goFunction: 'Log', description: '求自然对数', category: '数学', params: ['数值'], returnType: '小数型' },
  { chineseName: '求以10为底对数', goFunction: 'Log10', description: '求以10为底对数', category: '数学', params: ['数值'], returnType: '小数型' },
  { chineseName: '求反正弦', goFunction: 'Asin', description: '求反正弦值', category: '数学', params: ['数值'], returnType: '小数型' },
  { chineseName: '求反余弦', goFunction: 'Acos', description: '求反余弦值', category: '数学', params: ['数值'], returnType: '小数型' },
  { chineseName: '求反正切', goFunction: 'Atan', description: '求反正切值', category: '数学', params: ['数值'], returnType: '小数型' },
  { chineseName: '求反正切2', goFunction: 'Atan2', description: '求反正切值（双参数）', category: '数学', params: ['Y值', 'X值'], returnType: '小数型' },
  { chineseName: '角度转弧度', goFunction: 'DegToRad', description: '角度转弧度', category: '数学', params: ['角度'], returnType: '小数型' },
  { chineseName: '弧度转角度', goFunction: 'RadToDeg', description: '弧度转角度', category: '数学', params: ['弧度'], returnType: '小数型' },
])
const fnLibraryCategories = computed(() => {
  const cats = new Set(chineseFunctions.value.map(f => f.category))
  return Array.from(cats).sort()
})

const filteredFnLibrary = computed(() => {
  if (!fnSearchText.value.trim()) return chineseFunctions.value
  const q = fnSearchText.value.toLowerCase()
  return chineseFunctions.value.filter(f =>
    f.chineseName.toLowerCase().includes(q) ||
    f.goFunction.toLowerCase().includes(q) ||
    f.description.toLowerCase().includes(q)
  )
})

const getFnCategoryColor = (cat: string): string => {
  const colors: Record<string, string> = {
    '窗口': '#409eff', '文件': '#67c23a', '文本': '#e6a23c',
    '数学': '#f56c6c', '数组': '#909399', '系统': '#409eff',
    '网络': '#e6a23c', '数据库': '#f56c6c', '界面': '#909399',
    '编码': '#409eff', '时间': '#67c23a', '转换': '#e6a23c',
    '调试': '#f56c6c',
  }
  return colors[cat] || '#909399'
}

const toggleFnCategory = (cat: string) => {
  if (expandedFnCats.value.has(cat)) {
    expandedFnCats.value.delete(cat)
  } else {
    expandedFnCats.value.add(cat)
  }
}

const insertFunctionAtCursor = (fn: ChineseFunction) => {
  const editorRef = codeEditorRefs.value[activeWindowId.value]
  if (!editorRef) return
  const instance = editorRef.getInstance?.()
  if (!instance) return

  const params = fn.params.map((p, i) => `\${${i + 1}:${p}}`).join(', ')
  const snippet = fn.returnType
    ? `${fn.chineseName}(${params})  // → ${fn.returnType}`
    : `${fn.chineseName}(${params})`

  instance.executeEdits('insert-function', [{
    range: instance.getSelection() || {
      startLineNumber: instance.getPosition().lineNumber,
      startColumn: instance.getPosition().column,
      endLineNumber: instance.getPosition().lineNumber,
      endColumn: instance.getPosition().column,
    },
    text: snippet,
  }])
  instance.focus()
}

const onComponentFocus = (component: any) => {
  focusedComponent.value = component
  emit('component-focus', component)

  if (component && viewMode.value === 'split') {
    nextTick(() => {
      navigateToComponentCode(component)
    })
  }
}

const navigateToComponentCode = (component: any) => {
  const editorRef = codeEditorRefs.value[activeWindowId.value]
  if (!editorRef) return

  const instance = editorRef.getInstance?.()
  if (!instance) return

  const code = generatedCode.value
  const lines = code.split('\n')
  const fieldName = toGoFieldName(component.name || component.type)

  for (let i = 0; i < lines.length; i++) {
    if (lines[i].includes(fieldName) || lines[i].includes(component.id)) {
      instance.revealLineInCenter(i + 1)
      instance.setPosition({ lineNumber: i + 1, column: 1 })
      break
    }
  }
}

const componentToChineseCreate = (type: string): string => {
  const map: Record<string, string> = {
    'button': '创建按钮', 'input': '创建编辑框', 'label': '创建标签',
    'edit': '创建编辑框', 'checkbox': '创建选择框', 'radio': '创建单选框',
    'listbox': '创建列表框', 'combobox': '创建组合框', 'image': '创建图片框',
    'container': '创建面板', 'webview': '创建浏览器', 'pdfview': '创建面板',
    'progress': '创建进度条', 'slider': '创建滑块', 'switch': '创建选择框',
    'table': '创建表格', 'tree': '创建树形框', 'tabs': '创建选项卡',
    'menu': '创建菜单', 'toolbar': '创建工具栏', 'statusbar': '创建状态栏',
    'dialog': '创建打开对话框', 'timer': '创建定时器',
    'chart': '创建面板', 'map': '创建面板', 'video': '创建面板', 'audio': '创建面板',
    'richtext': '创建多行编辑框', 'datepicker': '创建日期框',
    'timepicker': '创建日期框', 'colorpicker': '创建颜色框',
    'number': '创建编辑框', 'password': '创建编辑框', 'textarea': '创建多行编辑框',
    'dropdown': '创建组合框', 'breadcrumb': '创建面板', 'pagination': '创建面板',
    'steps': '创建面板', 'card': '创建面板', 'collapse': '创建面板',
    'divider': '创建斜角框', 'badge': '创建标签', 'tag': '创建标签',
    'avatar': '创建图片框', 'spin': '创建调节器', 'transfer': '创建列表框',
    'upload': '创建按钮', 'rate': '创建面板', 'slider-range': '创建滑块',
    'input-number': '创建调节器', 'select': '创建组合框', 'cascader': '创建组合框',
    'form': '创建面板', 'descriptions': '创建面板', 'result': '创建面板',
    'skeleton': '创建面板', 'empty': '创建面板', 'popover': '创建面板',
    'tooltip': '创建面板', 'popconfirm': '创建面板', 'alert': '创建面板',
    'message': '创建面板', 'notification': '创建面板', 'modal': '创建面板',
    'drawer': '创建面板', 'loading': '创建面板', 'affix': '创建面板',
    'anchor': '创建面板', 'backtop': '创建按钮', 'watermark': '创建面板',
    'carousel': '创建面板', 'timeline': '创建面板', 'calendar': '创建月历',
    'statistic': '创建标签', 'qrcode': '创建图片框', 'barcode': '创建图片框',
    'icon': '创建图片框', 'text': '创建标签', 'link': '创建标签',
    'title': '创建标签', 'paragraph': '创建标签', 'group': '创建分组框',
    'splitter': '创建分隔条', 'scrollbar': '创建滚动条', 'scrollbox': '创建滚动框',
    'shape': '创建外形框', 'gauge': '创建进度条', 'led': '创建外形框',
    'knob': '创建滑块', 'meter': '创建进度条', 'thermometer': '创建进度条',
    'compass': '创建面板', 'odometer': '创建标签', 'trafficlight': '创建面板',
    'battery': '创建进度条', 'tank': '创建进度条', 'pipe': '创建面板',
    'valve': '创建面板', 'motor': '创建面板', 'conveyor': '创建面板',
    'robotarm': '创建面板', 'sensor': '创建面板', 'camera': '创建面板',
    'plc': '创建面板', 'hmi': '创建面板', 'scada': '创建面板',
    'alarm': '创建面板', 'trend': '创建面板', 'pid': '创建面板',
    'heading': '创建标签', 'blockquote': '创建面板',
    'code-block': '创建多行编辑框',
    'color-picker': '创建颜色框', 'date-picker': '创建日期框',
    'time-picker': '创建日期框',
    'autocomplete': '创建编辑框', 'rich-text': '创建多行编辑框',
    'list': '创建列表框', 'image-viewer': '创建图片框',
    'virtual-list': '创建列表框',
    'navbar': '创建面板', 'progress-bar': '创建进度条',
    'confirm': '创建打开对话框', 'panel': '创建面板',
    'layout': '创建面板', 'row': '创建面板', 'col': '创建面板',
    'flex': '创建面板',
    'split-pane': '创建分隔条', 'tab-container': '创建选项卡',
    'line-chart': '创建面板', 'bar-chart': '创建面板',
    'ring-progress': '创建进度条', 'pie-chart': '创建面板',
    'scatter-chart': '创建面板', 'radar-chart': '创建面板',
    'area-chart': '创建面板', 'funnel-chart': '创建面板',
    'heatmap': '创建面板', 'space': '创建面板',
  }
  return map[type] || '创建面板'
}

const generateGoCodeFromComponents = (components: any[]): string => {
  const lines: string[] = []

  lines.push('// ==========================================')
  lines.push('//  易狗 IDE 自动生成 - 中文函数库风格')
  lines.push('//  设计画布组件 → Go 代码实时联动')
  lines.push('// ==========================================')
  lines.push('')
  lines.push('package main')
  lines.push('')
  lines.push('import (')
  lines.push('\t"yigou"')
  lines.push(')')
  lines.push('')

  if (!components || components.length === 0) {
    lines.push('// 主窗口结构')
    lines.push('type 主窗口 struct {')
    lines.push('\t*yigou.窗口')
    lines.push('}')
    lines.push('')
    lines.push('// 窗口创建完毕事件')
    lines.push('func (f *主窗口) 创建完毕() {')
    lines.push('\tyigou.置窗口标题(f, "启动窗口")')
    lines.push('\tyigou.置窗口大小(f, 800, 600)')
    lines.push('}')
    lines.push('')
    lines.push('func main() {')
    lines.push('\tyigou.应用初始化()')
    lines.push('\tyigou.应用创建窗口(&主窗口{})')
    lines.push('\tyigou.应用运行()')
    lines.push('}')
    return lines.join('\n')
  }

  lines.push('// 主窗口结构 - 包含所有画布组件')
  lines.push('type 主窗口 struct {')
  lines.push('\t*yigou.窗口')
  lines.push('')

  const usedNames = new Set<string>()
  components.forEach(comp => {
    let fieldName = toGoFieldName(comp.name || comp.type)
    if (usedNames.has(fieldName)) {
      fieldName = fieldName + '_' + comp.id.replace(/[^a-zA-Z0-9]/g, '')
    }
    usedNames.add(fieldName)
    const chineseType = componentToChineseType(comp.type)
    lines.push(`\t${fieldName} *yigou.${chineseType}`)
  })

  lines.push('}')
  lines.push('')

  lines.push('// 窗口创建完毕 - 初始化所有组件')
  lines.push('func (f *主窗口) 创建完毕() {')
  lines.push('\tyigou.置窗口标题(f, "启动窗口")')
  lines.push('\tyigou.置窗口大小(f, 800, 600)')
  lines.push('')

  const fieldNames: string[] = []
  usedNames.clear()
  components.forEach(comp => {
    let fieldName = toGoFieldName(comp.name || comp.type)
    if (usedNames.has(fieldName)) {
      fieldName = fieldName + '_' + comp.id.replace(/[^a-zA-Z0-9]/g, '')
    }
    usedNames.add(fieldName)
    fieldNames.push(fieldName)

    const createFunc = componentToChineseCreate(comp.type)
    lines.push(`\t// ${comp.name || comp.type} (ID: ${comp.id})`)
    lines.push(`\tf.${fieldName} = yigou.${createFunc}(f)`)
    lines.push(`\tf.${fieldName}.置边界(${comp.x}, ${comp.y}, ${comp.width}, ${comp.height})`)

    if (comp.text) {
      lines.push(`\tf.${fieldName}.置标题("${comp.text}")`)
    }
    if (comp.visible === false) {
      lines.push(`\tf.${fieldName}.置可视(假)`)
    }
    if (comp.enabled === false) {
      lines.push(`\tf.${fieldName}.置可用(假)`)
    }
    if (comp.color && comp.color !== '#409eff') {
      lines.push(`\tf.${fieldName}.置背景颜色("${comp.color}")`)
    }
    if (comp.fontName && comp.fontName !== '微软雅黑') {
      lines.push(`\tf.${fieldName}.置字体名称("${comp.fontName}")`)
    }
    if (comp.fontSize && comp.fontSize !== 14) {
      lines.push(`\tf.${fieldName}.置字体大小(${comp.fontSize})`)
    }
    if (comp.fontBold) {
      lines.push(`\tf.${fieldName}.置字体加粗(真)`)
    }
    if (comp.fontItalic) {
      lines.push(`\tf.${fieldName}.置字体倾斜(真)`)
    }
    if (comp.fontUnderline) {
      lines.push(`\tf.${fieldName}.置字体下划线(真)`)
    }
    if (comp.hint) {
      lines.push(`\tf.${fieldName}.置提示("${comp.hint}")`)
    }
    if (comp.cursor && comp.cursor !== '默认') {
      lines.push(`\tf.${fieldName}.置光标("${comp.cursor}")`)
    }
    if (comp.align && comp.align !== '左对齐') {
      const alignMap: Record<string, string> = { '左对齐': '左', '居中': '中', '右对齐': '右', '顶对齐': '顶', '底对齐': '底' }
      lines.push(`\tf.${fieldName}.置对齐(yigou.对齐_${alignMap[comp.align] || '左'})`)
    }
    if (comp.anchors) {
      const a = comp.anchors
      if (!a.left || !a.top || a.right || a.bottom) {
        lines.push(`\tf.${fieldName}.置锚点(${a.left ? '真' : '假'}, ${a.top ? '真' : '假'}, ${a.right ? '真' : '假'}, ${a.bottom ? '真' : '假'})`)
      }
    }
    if (comp.padding && (comp.padding.left || comp.padding.top || comp.padding.right || comp.padding.bottom)) {
      lines.push(`\tf.${fieldName}.置内边距(${comp.padding.left}, ${comp.padding.top}, ${comp.padding.right}, ${comp.padding.bottom})`)
    }
    if (comp.margin && (comp.margin.left || comp.margin.top || comp.margin.right || comp.margin.bottom)) {
      lines.push(`\tf.${fieldName}.置外边距(${comp.margin.left}, ${comp.margin.top}, ${comp.margin.right}, ${comp.margin.bottom})`)
    }
    if (comp.autoSize) {
      lines.push(`\tf.${fieldName}.置自动大小(真)`)
    }
    if (comp.tabStop === false) {
      lines.push(`\tf.${fieldName}.置Tab停靠(假)`)
    }

    if (comp.readOnly) {
      lines.push(`\tf.${fieldName}.置只读(真)`)
    }
    if (comp.passwordMode) {
      lines.push(`\tf.${fieldName}.置密码模式(真, "${comp.passwordChar || '*'}")`)
    }
    if (comp.maxLength && comp.maxLength > 0) {
      lines.push(`\tf.${fieldName}.置最大长度(${comp.maxLength})`)
    }
    if (comp.multiLine) {
      lines.push(`\tf.${fieldName}.置多行模式(真)`)
    }
    if (comp.numberMode) {
      lines.push(`\tf.${fieldName}.置数字模式(真)`)
    }

    if (comp.checked) {
      lines.push(`\tf.${fieldName}.置选中(真)`)
    }

    if (comp.items && comp.items.length > 0) {
      const itemsStr = comp.items.map((item: string) => `"${item}"`).join(', ')
      lines.push(`\tf.${fieldName}.置项目列表([${itemsStr}])`)
    }
    if (comp.selectedIndex >= 0) {
      lines.push(`\tf.${fieldName}.置选中索引(${comp.selectedIndex})`)
    }
    if (comp.dropDownStyle && comp.dropDownStyle !== '下拉列表') {
      lines.push(`\tf.${fieldName}.置下拉样式("${comp.dropDownStyle}")`)
    }

    if (comp.min !== 0 || comp.max !== 100) {
      lines.push(`\tf.${fieldName}.置范围(${comp.min}, ${comp.max})`)
    }
    if (comp.position !== 0) {
      lines.push(`\tf.${fieldName}.置位置(${comp.position})`)
    }
    if (comp.orientation && comp.orientation !== '水平') {
      lines.push(`\tf.${fieldName}.置方向("${comp.orientation}")`)
    }

    if (comp.rowCount !== 5 || comp.colCount !== 3) {
      lines.push(`\tf.${fieldName}.置行列数(${comp.rowCount}, ${comp.colCount})`)
    }
    if (comp.fixedRows > 0 || comp.fixedCols > 0) {
      lines.push(`\tf.${fieldName}.置固定行列(${comp.fixedRows}, ${comp.fixedCols})`)
    }

    if (comp.picturePath) {
      lines.push(`\tf.${fieldName}.置图片("${comp.picturePath}")`)
    }
    if (comp.stretch) {
      lines.push(`\tf.${fieldName}.置缩放(真)`)
    }
    if (comp.center) {
      lines.push(`\tf.${fieldName}.置居中(真)`)
    }
    if (comp.proportional) {
      lines.push(`\tf.${fieldName}.置等比缩放(真)`)
    }

    if (comp.interval && comp.interval !== 1000) {
      lines.push(`\tf.${fieldName}.置间隔(${comp.interval})`)
    }

    if (comp.url) {
      lines.push(`\tf.${fieldName}.置URL("${comp.url}")`)
    }

    if (comp.shapeType && comp.shapeType !== '矩形') {
      lines.push(`\tf.${fieldName}.置形状("${comp.shapeType}")`)
    }
    if (comp.penColor && comp.penColor !== '#000000') {
      lines.push(`\tf.${fieldName}.置画笔颜色("${comp.penColor}")`)
    }
    if (comp.penWidth && comp.penWidth !== 1) {
      lines.push(`\tf.${fieldName}.置画笔宽度(${comp.penWidth})`)
    }
    if (comp.brushColor && comp.brushColor !== '#ffffff') {
      lines.push(`\tf.${fieldName}.置画刷颜色("${comp.brushColor}")`)
    }

    lines.push('')
  })

  lines.push('\t// 注册组件事件')
  components.forEach((comp, i) => {
    const fn = fieldNames[i]
    if (comp.events && comp.events.length > 0) {
      comp.events.forEach((evt: { name: string; handler: string }) => {
        const eventConst = eventNameToConst(evt.name)
        lines.push(`\tf.${fn}.置事件(yigou.${eventConst}, f.${fn}_${evt.handler})`)
      })
    } else {
      const eventName = `${fn}_被单击`
      lines.push(`\t// f.${fn}.置事件(yigou.事件_被单击, f.${eventName})`)
    }
  })

  lines.push('}')
  lines.push('')

  components.forEach((comp, i) => {
    const fn = fieldNames[i]
    if (comp.events && comp.events.length > 0) {
      comp.events.forEach((evt: { name: string; handler: string }) => {
        lines.push(`// ${comp.name || comp.type} ${evt.name}事件`)
        lines.push(`func (f *主窗口) ${fn}_${evt.handler}() {`)
        lines.push(`\tyigou.调试输出("${comp.name || comp.type} ${evt.name}")`)
        lines.push('}')
        lines.push('')
      })
    } else {
      const eventName = `${fn}_被单击`
      lines.push(`// ${comp.name || comp.type} 被单击事件`)
      lines.push(`func (f *主窗口) ${eventName}() {`)
      lines.push(`\tyigou.调试输出("${comp.name || comp.type} 被单击")`)
      lines.push('}')
      lines.push('')
    }
  })

  lines.push('func main() {')
  lines.push('\tyigou.应用初始化()')
  lines.push('\tyigou.应用创建窗口(&主窗口{})')
  lines.push('\tyigou.应用运行()')
  lines.push('}')

  return lines.join('\n')
}

const toGoFieldName = (name: string): string => {
  if (!name) return '未命名组件'
  const cleaned = name.replace(/[^a-zA-Z0-9_\u4e00-\u9fa5]/g, '')
  if (!cleaned) return '未命名组件'
  if (/^[a-zA-Z_]/.test(cleaned)) {
    return cleaned.charAt(0).toUpperCase() + cleaned.slice(1)
  }
  return cleaned
}

const componentToChineseType = (type: string): string => {
  const map: Record<string, string> = {
    'button': '按钮', 'input': '编辑框', 'label': '标签',
    'edit': '编辑框', 'checkbox': '选择框', 'radio': '单选框',
    'listbox': '列表框', 'combobox': '组合框', 'image': '图片框',
    'container': '面板', 'webview': '浏览器', 'pdfview': '面板',
    'progress': '进度条', 'slider': '滑块', 'switch': '选择框',
    'table': '表格', 'tree': '树形框', 'tabs': '选项卡',
    'menu': '菜单', 'toolbar': '工具栏', 'statusbar': '状态栏',
    'dialog': '对话框', 'timer': '定时器', 'chart': '面板',
    'map': '面板', 'video': '面板', 'audio': '面板',
    'richtext': '多行编辑框', 'datepicker': '日期框',
    'timepicker': '日期框', 'colorpicker': '颜色框',
    'number': '编辑框', 'password': '编辑框', 'textarea': '多行编辑框',
    'dropdown': '组合框', 'breadcrumb': '面板', 'pagination': '面板',
    'steps': '面板', 'card': '面板', 'collapse': '面板',
    'divider': '斜角框', 'badge': '标签', 'tag': '标签',
    'avatar': '图片框', 'spin': '调节器', 'transfer': '列表框',
    'upload': '按钮', 'rate': '面板', 'slider-range': '滑块',
    'input-number': '调节器', 'select': '组合框', 'cascader': '组合框',
    'form': '面板', 'descriptions': '面板', 'result': '面板',
    'skeleton': '面板', 'empty': '面板', 'popover': '面板',
    'tooltip': '面板', 'popconfirm': '面板', 'alert': '面板',
    'message': '面板', 'notification': '面板', 'modal': '面板',
    'drawer': '面板', 'loading': '面板', 'affix': '面板',
    'anchor': '面板', 'backtop': '按钮', 'watermark': '面板',
    'carousel': '面板', 'timeline': '面板', 'calendar': '月历',
    'statistic': '标签', 'qrcode': '图片框', 'barcode': '图片框',
    'icon': '图片框', 'text': '标签', 'link': '标签',
    'title': '标签', 'paragraph': '标签', 'group': '分组框',
    'splitter': '分隔条', 'scrollbar': '滚动条', 'scrollbox': '滚动框',
    'shape': '外形框', 'gauge': '进度条', 'led': '外形框',
    'knob': '滑块', 'meter': '进度条', 'thermometer': '进度条',
    'compass': '面板', 'odometer': '标签', 'trafficlight': '面板',
    'battery': '进度条', 'tank': '进度条', 'pipe': '面板',
    'valve': '面板', 'motor': '面板', 'conveyor': '面板',
    'robotarm': '面板', 'sensor': '面板', 'camera': '面板',
    'plc': '面板', 'hmi': '面板', 'scada': '面板',
    'alarm': '面板', 'trend': '面板', 'pid': '面板',
  }
  return map[type] || '面板'
}

const eventNameToConst = (name: string): string => {
  const map: Record<string, string> = {
    '被单击': '事件_被单击',
    '被双击': '事件_被双击',
    '内容改变': '事件_内容改变',
    '选择改变': '事件_选择改变',
    '获得焦点': '事件_获得焦点',
    '失去焦点': '事件_失去焦点',
    '鼠标按下': '事件_鼠标按下',
    '鼠标释放': '事件_鼠标释放',
    '鼠标移动': '事件_鼠标移动',
    '鼠标滚轮': '事件_鼠标滚轮',
    '按键按下': '事件_按键按下',
    '按键释放': '事件_按键释放',
    '大小改变': '事件_大小改变',
    '位置改变': '事件_位置改变',
    '定时器': '事件_定时器',
    '关闭': '事件_关闭',
    '创建完毕': '事件_创建完毕',
    '即将销毁': '事件_即将销毁',
  }
  return map[name] || '事件_被单击'
}

const syncComponentsToCodeEditor = (components: any[]) => {
  generatedCode.value = generateGoCodeFromComponents(components)
  const editorRef = codeEditorRefs.value[activeWindowId.value]
  if (editorRef) {
    editorRef.setValue?.(generatedCode.value)
  }
}

watch(activeWindowId, (newId) => {
  nextTick(() => {
    const canvasRef = canvasRefs.value[newId]
    if (canvasRef) {
      const comps = canvasRef.getComponents?.()
      if (comps && comps.length > 0) {
        syncComponentsToCodeEditor(comps)
      }
    }
  })
})

const onCodeRowsChange = (rows: any[]) => {
  emit('code-rows-change', rows)
}

const switchWindow = (id: string) => {
  activeWindowId.value = id
}

const addWindow = () => {
  const newWin: WindowTab = {
    id: `win_${nextWindowId++}`,
    name: `窗口${nextWindowId - 1}.ego`,
    subTab: 'design',
  }
  windows.push(newWin)
  activeWindowId.value = newWin.id
  emit('windows-change', [...windows])
}

const removeWindow = (id: string) => {
  const idx = windows.findIndex(w => w.id === id)
  if (idx === -1 || windows.length <= 1) return
  windows.splice(idx, 1)
  delete canvasRefs.value[id]
  delete codeEditorRefs.value[id]
  if (activeWindowId.value === id) {
    activeWindowId.value = windows[Math.min(idx, windows.length - 1)].id
  }
  emit('windows-change', [...windows])
}

const switchSubTab = (winId: string, tab: 'design' | 'code') => {
  const win = windows.find(w => w.id === winId)
  if (win) {
    win.subTab = tab
  }
}

const toggleViewMode = () => {
  viewMode.value = viewMode.value === 'tab' ? 'split' : 'tab'
}

const getCanvasRef = (winId: string, el: any) => {
  if (el) {
    canvasRefs.value[winId] = el
  }
}

const getCodeEditorRef = (winId: string, el: any) => {
  if (el) {
    codeEditorRefs.value[winId] = el
  }
}

const startDrag = (e: MouseEvent) => {
  isDragging.value = true
  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
  e.preventDefault()
}

const onDrag = (e: MouseEvent) => {
  if (!isDragging.value) return
  const container = (e.target as HTMLElement)?.closest('.tab-content') as HTMLElement
  if (!container) return
  const rect = container.getBoundingClientRect()
  const newRatio = ((e.clientX - rect.left) / rect.width) * 100
  splitRatio.value = Math.max(25, Math.min(75, newRatio))
}

const stopDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
}

const addIndependentCodeFile = (id: string, name: string) => {
  const newFile: IndependentFile = {
    id,
    name,
    code: '',
  }
  independentFiles.push(newFile)
  activeIndependentFileId.value = newFile.id
}

const removeIndependentCodeFile = (id: string) => {
  const idx = independentFiles.findIndex(f => f.id === id)
  if (idx === -1) return
  independentFiles.splice(idx, 1)
  delete codeEditorRefs.value[id]
  if (activeIndependentFileId.value === id) {
    activeIndependentFileId.value = independentFiles.length > 0
      ? independentFiles[Math.min(idx, independentFiles.length - 1)].id
      : null
  }
}

const renameIndependentCodeFile = (id: string, newName: string) => {
  const file = independentFiles.find(f => f.id === id)
  if (file) file.name = newName
}

const switchIndependentFile = (id: string) => {
  activeIndependentFileId.value = id
}

defineExpose({
  windows,
  independentFiles,
  activeWindowId,
  activeIndependentFileId,
  activeWindow,
  activeIndependentFile,
  canvasRefs,
  codeEditorRefs,
  addWindow,
  removeWindow,
  switchWindow,
  addIndependentCodeFile,
  removeIndependentCodeFile,
  renameIndependentCodeFile,
  switchIndependentFile,
  syncComponentsToCodeEditor,
  chineseFunctions,
})
</script>

<template>
  <div class="center-tabs">
    <div class="main-tab-bar">
      <div class="tab-group window-group">
        <div class="tabs-container">
          <div
            v-for="win in windows"
            :key="win.id"
            class="main-tab"
            :class="{ active: activeWindowId === win.id }"
            @click="switchWindow(win.id)"
          >
            <span class="tab-icon">🪟</span>
            <span class="tab-name">{{ win.name }}</span>
            <el-icon
              v-if="windows.length > 1"
              class="tab-close"
              @click.stop="removeWindow(win.id)"
            >
              <Close />
            </el-icon>
          </div>
          <el-button class="add-tab-btn" size="small" text @click="addWindow" title="新建窗体">
            <el-icon><Plus /></el-icon>
          </el-button>
        </div>
      </div>
      <div v-if="independentFiles.length > 0" class="tab-group file-group">
        <div class="tabs-container">
          <div
            v-for="file in independentFiles"
            :key="file.id"
            class="main-tab file-tab"
            :class="{ active: activeIndependentFileId === file.id }"
            @click="switchIndependentFile(file.id)"
          >
            <span class="tab-icon">📄</span>
            <span class="tab-name">{{ file.name }}</span>
            <el-icon
              class="tab-close"
              @click.stop="removeIndependentCodeFile(file.id)"
            >
              <Close />
            </el-icon>
          </div>
        </div>
      </div>
    </div>

    <div v-if="activeWindow" class="sub-tab-bar">
      <div
        class="sub-tab-item"
        :class="{ active: activeWindow.subTab === 'design' }"
        @click="switchSubTab(activeWindow.id, 'design')"
      >
        <el-icon><Grid /></el-icon>
        <span>窗体设计</span>
      </div>
      <div
        class="sub-tab-item"
        :class="{ active: activeWindow.subTab === 'code' }"
        @click="switchSubTab(activeWindow.id, 'code')"
      >
        <el-icon><Edit /></el-icon>
        <span>Go 代码</span>
        <span class="code-lang-badge">中文库</span>
      </div>

      <div class="sub-tab-actions">
        <el-tooltip
          :content="viewMode === 'tab' ? '切换为左右分屏' : '切换为标签模式'"
          placement="bottom"
        >
          <el-button
            size="small"
            :type="viewMode === 'split' ? 'primary' : ''"
            text
            @click="toggleViewMode"
          >
            <el-icon><component :is="viewMode === 'tab' ? 'Right' : 'Grid'" /></el-icon>
            <span class="mode-text">{{ viewMode === 'tab' ? '分屏' : '标签' }}</span>
          </el-button>
        </el-tooltip>
      </div>
    </div>

    <div class="tab-content">
      <template v-for="win in windows" :key="win.id">
        <div v-show="win.id === activeWindowId" class="window-pane">
            <template v-if="viewMode === 'tab'">
              <div v-show="win.subTab === 'design'" class="sub-pane">
                <DesignCanvas
                  :ref="(el: any) => getCanvasRef(win.id, el)"
                  @component-focus="onComponentFocus"
                  @components-update="syncComponentsToCodeEditor"
                />
              </div>
              <div v-show="win.subTab === 'code'" class="sub-pane code-pane">
                <div class="code-pane-header">
                  <span class="code-pane-title">
                    <el-icon><Edit /></el-icon>
                    Go 代码 — 中文函数库风格
                  </span>
                  <span v-if="focusedComponent" class="focus-badge">
                    已选中: {{ focusedComponent.name || focusedComponent.type }}
                  </span>
                  <el-button
                    size="small"
                    text
                    class="fn-lib-toggle"
                    :type="showFnLibrary ? 'primary' : ''"
                    @click="showFnLibrary = !showFnLibrary"
                    title="切换函数库面板"
                  >
                    <el-icon><Collection /></el-icon>
                    <span class="mode-text">函数库</span>
                  </el-button>
                </div>
                <div class="code-pane-body code-with-fnlib">
                  <div class="code-editor-area" :style="{ width: showFnLibrary ? '70%' : '100%' }">
                    <MonacoEditorWrapper
                      :ref="(el: any) => getCodeEditorRef(win.id, el)"
                      :model-value="generatedCode"
                      language="go"
                      :theme="editorTheme"
                      :font-size="13"
                      :read-only="false"
                    />
                  </div>
                  <div v-if="showFnLibrary" class="fn-lib-panel">
                    <div class="fn-lib-search">
                      <el-input
                        v-model="fnSearchText"
                        size="small"
                        placeholder="搜索函数..."
                        clearable
                      />
                    </div>
                    <div class="fn-lib-list">
                      <div v-for="cat in fnLibraryCategories" :key="cat" class="fn-lib-cat">
                        <div class="fn-lib-cat-header" @click="toggleFnCategory(cat)">
                          <span class="cat-arrow">{{ expandedFnCats.has(cat) ? '▼' : '▶' }}</span>
                          <span class="cat-dot" :style="{ backgroundColor: getFnCategoryColor(cat) }"></span>
                          <span class="cat-name">{{ cat }}</span>
                          <span class="cat-count">{{ filteredFnLibrary.filter(f => f.category === cat).length }}</span>
                        </div>
                        <div v-show="expandedFnCats.has(cat) || fnSearchText" class="fn-lib-items">
                          <div
                            v-for="fn in filteredFnLibrary.filter(f => f.category === cat)"
                            :key="fn.chineseName"
                            class="fn-lib-item"
                            @click="insertFunctionAtCursor(fn)"
                            :title="fn.description"
                          >
                            <span class="fn-lib-name">{{ fn.chineseName }}</span>
                            <span v-if="fn.returnType" class="fn-lib-return">{{ fn.returnType }}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </template>

            <template v-else>
              <div class="split-pane">
                <div class="split-left" :style="{ width: splitRatio + '%' }">
                  <div class="pane-header">
                    <el-icon><Grid /></el-icon>
                    <span>窗体设计</span>
                    <span v-if="focusedComponent" class="focus-badge">
                      已选中: {{ focusedComponent.name || focusedComponent.type }}
                    </span>
                  </div>
                  <div class="pane-body">
                    <DesignCanvas
                      :ref="(el: any) => getCanvasRef(win.id, el)"
                      @component-focus="onComponentFocus"
                      @components-update="syncComponentsToCodeEditor"
                    />
                  </div>
                </div>

                <div
                  class="split-handle"
                  :class="{ active: isDragging }"
                  @mousedown="startDrag"
                >
                  <div class="handle-line"></div>
                </div>

                <div class="split-right" :style="{ width: (100 - splitRatio) + '%' }">
                  <div class="pane-header code-header">
                    <el-icon><Edit /></el-icon>
                    <span>Go 代码 — 中文函数库风格</span>
                    <el-button
                      size="small"
                      text
                      class="fn-lib-toggle"
                      :type="showFnLibrary ? 'primary' : ''"
                      @click="showFnLibrary = !showFnLibrary"
                      title="切换函数库面板"
                    >
                      <el-icon><Collection /></el-icon>
                      <span class="mode-text">函数库</span>
                    </el-button>
                  </div>
                  <div class="pane-body code-body code-with-fnlib">
                    <div class="code-editor-area" :style="{ width: showFnLibrary ? '65%' : '100%' }">
                      <MonacoEditorWrapper
                        :ref="(el: any) => getCodeEditorRef(win.id, el)"
                        :model-value="generatedCode"
                        language="go"
                        :theme="editorTheme"
                        :font-size="13"
                        :read-only="false"
                      />
                    </div>
                    <div v-if="showFnLibrary" class="fn-lib-panel">
                      <div class="fn-lib-search">
                        <el-input
                          v-model="fnSearchText"
                          size="small"
                          placeholder="搜索函数..."
                          clearable
                        />
                      </div>
                      <div class="fn-lib-list">
                        <div v-for="cat in fnLibraryCategories" :key="cat" class="fn-lib-cat">
                          <div class="fn-lib-cat-header" @click="toggleFnCategory(cat)">
                            <span class="cat-arrow">{{ expandedFnCats.has(cat) ? '▼' : '▶' }}</span>
                            <span class="cat-dot" :style="{ backgroundColor: getFnCategoryColor(cat) }"></span>
                            <span class="cat-name">{{ cat }}</span>
                            <span class="cat-count">{{ filteredFnLibrary.filter(f => f.category === cat).length }}</span>
                          </div>
                          <div v-show="expandedFnCats.has(cat) || fnSearchText" class="fn-lib-items">
                            <div
                              v-for="fn in filteredFnLibrary.filter(f => f.category === cat)"
                              :key="fn.chineseName"
                              class="fn-lib-item"
                              @click="insertFunctionAtCursor(fn)"
                              :title="fn.description"
                            >
                              <span class="fn-lib-name">{{ fn.chineseName }}</span>
                              <span v-if="fn.returnType" class="fn-lib-return">{{ fn.returnType }}</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </template>

      <template v-for="file in independentFiles" :key="file.id">
        <div v-show="file.id === activeIndependentFileId" class="window-pane">
          <div class="code-pane independent-code-pane">
            <div class="code-pane-header">
              <span class="code-pane-title">
                <el-icon><Document /></el-icon>
                共享代码 — {{ file.name }}
              </span>
            </div>
            <div class="code-pane-body code-with-fnlib">
              <div class="code-editor-area" :style="{ width: showFnLibrary ? '70%' : '100%' }">
                <MonacoEditorWrapper
                  :ref="(el: any) => getCodeEditorRef(file.id, el)"
                  :model-value="file.code"
                  language="go"
                  :theme="editorTheme"
                  :font-size="13"
                  :read-only="false"
                  @update:model-value="(val: string) => file.code = val"
                />
              </div>
              <div v-if="showFnLibrary" class="fn-lib-panel">
                <div class="fn-lib-search">
                  <el-input
                    v-model="fnSearchText"
                    size="small"
                    placeholder="搜索函数..."
                    clearable
                  />
                </div>
                <div class="fn-lib-list">
                  <div v-for="cat in fnLibraryCategories" :key="cat" class="fn-lib-cat">
                    <div class="fn-lib-cat-header" @click="toggleFnCategory(cat)">
                      <span class="cat-arrow">{{ expandedFnCats.has(cat) ? '▼' : '▶' }}</span>
                      <span class="cat-dot" :style="{ backgroundColor: getFnCategoryColor(cat) }"></span>
                      <span class="cat-name">{{ cat }}</span>
                      <span class="cat-count">{{ filteredFnLibrary.filter(f => f.category === cat).length }}</span>
                    </div>
                    <div v-show="expandedFnCats.has(cat) || fnSearchText" class="fn-lib-items">
                      <div
                        v-for="fn in filteredFnLibrary.filter(f => f.category === cat)"
                        :key="fn.chineseName"
                        class="fn-lib-item"
                        @click="insertFunctionAtCursor(fn)"
                        :title="fn.description"
                      >
                        <span class="fn-lib-name">{{ fn.chineseName }}</span>
                        <span v-if="fn.returnType" class="fn-lib-return">{{ fn.returnType }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <div v-if="!activeWindow && !activeIndependentFile" class="empty-state">
        <el-icon :size="48"><FolderOpened /></el-icon>
        <p>请从左侧构造面板选择或创建项目元素</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.center-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100%;
}

.main-tab-bar {
  display: flex;
  flex-direction: column;
  background-color: var(--bg-tertiary);
  border-bottom: 1px solid var(--border-dark);
  flex-shrink: 0;
}

.tab-group {
  border-bottom: 1px solid var(--border-light);
}

.tab-group:last-child {
  border-bottom: none;
}

.group-label {
  padding: 4px 10px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-secondary);
  background-color: var(--bg-secondary);
  letter-spacing: 0.5px;
}

.window-group .group-label {
  color: #409eff;
}

.module-group .group-label {
  color: #67c23a;
}

.tabs-container {
  display: flex;
  align-items: center;
  overflow-x: auto;
  scrollbar-width: thin;
  min-height: 32px;
}

.main-tab {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
  border-right: 1px solid var(--border-light);
  white-space: nowrap;
  transition: all 0.15s;
  user-select: none;
  min-width: 0;
  max-width: 160px;
}

.main-tab:hover {
  color: var(--text-primary);
  background-color: var(--bg-hover);
}

.main-tab.active {
  color: var(--text-primary);
  background-color: var(--bg-primary);
  border-bottom: 2px solid var(--color-primary);
}

.module-tab.active {
  border-bottom-color: #67c23a;
}

.file-tab.active {
  border-bottom-color: #e6a23c;
}

.tab-icon {
  font-size: 12px;
  flex-shrink: 0;
}

.tab-name {
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.tab-close {
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.15s;
  flex-shrink: 0;
}

.main-tab:hover .tab-close {
  opacity: 1;
}

.tab-close:hover {
  color: var(--color-danger);
}

.add-tab-btn {
  flex-shrink: 0;
  padding: 4px 8px !important;
  min-height: auto !important;
}

.sub-tab-bar {
  display: flex;
  align-items: center;
  background-color: var(--bg-canvas);
  border-bottom: 1px solid var(--border-dark);
  flex-shrink: 0;
  padding: 0 4px;
}

.module-sub-bar {
  background-color: #f0f9eb;
}

.sub-tab-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 14px;
  font-size: 12px;
  color: var(--text-regular);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
  user-select: none;
}

.sub-tab-item:hover {
  color: var(--color-primary);
  background-color: var(--color-primary-light);
}

.sub-tab-item.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
  background-color: var(--bg-primary);
}

.module-sub-bar .sub-tab-item.active {
  color: #67c23a;
  border-bottom-color: #67c23a;
  background-color: #f0f9eb;
}

.module-hint {
  font-size: 10px;
  color: var(--text-placeholder);
  margin-left: 8px;
}

.code-lang-badge {
  font-size: 9px;
  background: linear-gradient(135deg, #e74c3c, #f39c12);
  color: #fff;
  padding: 1px 5px;
  border-radius: 6px;
  margin-left: 4px;
  font-weight: 600;
}

.sub-tab-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 4px;
  padding-right: 8px;
}

.mode-text {
  font-size: 11px;
  margin-left: 2px;
}

.tab-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.window-pane {
  height: 100%;
}

.module-pane {
  height: 100%;
  background-color: #fafafa;
}

.sub-pane {
  height: 100%;
}

.code-pane {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #1e1e1e;
}

.code-pane-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  font-size: 11px;
  font-weight: 600;
  color: var(--panel-header-text);
  background-color: var(--panel-header-bg);
  border-bottom: 1px solid var(--panel-border);
  flex-shrink: 0;
  min-height: 28px;
}

.code-pane-title {
  display: flex;
  align-items: center;
  gap: 4px;
}

.code-pane-body {
  flex: 1;
  overflow: hidden;
}

.split-pane {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.split-left,
.split-right {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.split-left {
  background-color: var(--bg-canvas);
}

.split-right {
  background-color: var(--editor-bg);
}

.split-handle {
  width: 6px;
  cursor: col-resize;
  background-color: var(--border-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background-color 0.15s;
  z-index: 10;
}

.split-handle:hover,
.split-handle.active {
  background-color: var(--color-primary);
}

.handle-line {
  width: 2px;
  height: 32px;
  background-color: var(--text-placeholder);
  border-radius: 1px;
  transition: background-color 0.15s;
}

.split-handle:hover .handle-line,
.split-handle.active .handle-line {
  background-color: #fff;
}

.pane-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-secondary);
  background-color: var(--bg-tertiary);
  border-bottom: 1px solid var(--border-dark);
  flex-shrink: 0;
  min-height: 28px;
}

.code-header {
  background-color: var(--panel-header-bg);
  color: var(--panel-header-text);
  border-bottom-color: var(--panel-border);
}

.focus-badge {
  margin-left: auto;
  padding: 1px 8px;
  background-color: var(--color-primary-light);
  color: var(--color-primary);
  border-radius: 10px;
  font-size: 10px;
  font-weight: 500;
}

.pane-body {
  flex: 1;
  overflow: hidden;
}

.code-body {
  flex: 1;
  overflow: hidden;
}

.empty-state {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: var(--text-placeholder);
}

.empty-state p {
  font-size: 13px;
  margin: 0;
}

.fn-lib-toggle {
  margin-left: auto;
  color: var(--text-primary);
}

.fn-lib-toggle:hover {
  color: var(--color-primary);
}

.code-with-fnlib {
  display: flex;
  flex-direction: row;
}

.code-editor-area {
  flex: 1;
  overflow: hidden;
  transition: width 0.2s;
}

.fn-lib-panel {
  width: 200px;
  min-width: 160px;
  background-color: var(--panel-bg);
  border-left: 1px solid var(--panel-border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;
}

.fn-lib-search {
  padding: 6px;
  flex-shrink: 0;
}

.fn-lib-search :deep(.el-input__wrapper) {
  background-color: var(--input-bg);
  box-shadow: none;
  border: 1px solid var(--input-border);
}

.fn-lib-search :deep(.el-input__inner) {
  color: var(--input-text);
}

.fn-lib-list {
  flex: 1;
  overflow-y: auto;
  padding: 2px 4px;
}

.fn-lib-cat {
  margin-bottom: 2px;
}

.fn-lib-cat-header {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 6px;
  cursor: pointer;
  border-radius: 3px;
  font-size: 11px;
  user-select: none;
}

.fn-lib-cat-header:hover {
  background-color: var(--bg-hover);
}

.fn-lib-cat-header .cat-arrow {
  font-size: 9px;
  color: var(--text-secondary);
  width: 10px;
}

.fn-lib-cat-header .cat-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.fn-lib-cat-header .cat-name {
  flex: 1;
  font-weight: 500;
  color: var(--text-primary);
}

.fn-lib-cat-header .cat-count {
  font-size: 9px;
  color: var(--text-secondary);
  background-color: var(--bg-tertiary);
  padding: 0 5px;
  border-radius: 6px;
}

.fn-lib-items {
  padding-left: 16px;
}

.fn-lib-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 3px 6px;
  cursor: pointer;
  border-radius: 3px;
  font-size: 11px;
  transition: background-color 0.1s;
}

.fn-lib-item:hover {
  background-color: var(--bg-selected);
}

.fn-lib-name {
  color: var(--text-primary);
}

.fn-lib-return {
  font-size: 9px;
  color: var(--hl-type);
  background-color: var(--color-success-light);
  padding: 0 4px;
  border-radius: 2px;
  flex-shrink: 0;
}
</style>