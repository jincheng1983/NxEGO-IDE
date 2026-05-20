<script lang="ts" setup>
import { ref, computed } from 'vue'

interface UIComponent {
  id: string
  name: string
  icon: string
  category: string
  type: string
  description?: string
}

const searchKeyword = ref('')
const collapsedCategories = ref<Record<string, boolean>>({})

const toggleCategory = (name: string) => {
  collapsedCategories.value[name] = !collapsedCategories.value[name]
}

const categories = ref([
  {
    name: '🧱 基础',
    components: [
      { id: 'button', name: '按钮', icon: 'Button', category: '基础', type: 'button', description: '触发操作' },
      { id: 'label', name: '标签', icon: 'Document', category: '基础', type: 'label', description: '文本展示' },
      { id: 'text', name: '文本', icon: 'ChatLineSquare', category: '基础', type: 'text', description: '多行文本' },
      { id: 'link', name: '链接', icon: 'Link', category: '基础', type: 'link', description: '超链接' },
      { id: 'icon', name: '图标', icon: 'Picture', category: '基础', type: 'icon', description: '图标展示' },
      { id: 'divider', name: '分割线', icon: 'Minus', category: '基础', type: 'divider', description: '内容分隔' },
      { id: 'space', name: '间距', icon: 'Rank', category: '基础', type: 'space', description: '布局间距' },
      { id: 'heading', name: '标题', icon: 'EditPen', category: '基础', type: 'heading', description: '章节标题' },
      { id: 'paragraph', name: '段落', icon: 'Notebook', category: '基础', type: 'paragraph', description: '文本段落' },
      { id: 'blockquote', name: '引用', icon: 'ChatLineRound', category: '基础', type: 'blockquote', description: '引用文本' },
      { id: 'code-block', name: '代码块', icon: 'Code', category: '基础', type: 'code-block', description: '代码展示' },
      { id: 'image', name: '图片', icon: 'PictureFilled', category: '基础', type: 'image', description: '图片展示' },
      { id: 'video', name: '视频', icon: 'VideoCamera', category: '基础', type: 'video', description: '视频播放' },
    ] as UIComponent[],
  },
  {
    name: '📝 表单',
    components: [
      { id: 'input', name: '输入框', icon: 'Edit', category: '表单', type: 'input', description: '文本输入' },
      { id: 'textarea', name: '文本域', icon: 'Notebook', category: '表单', type: 'textarea', description: '多行输入' },
      { id: 'input-number', name: '数字输入', icon: 'Sort', category: '表单', type: 'input-number', description: '数值输入' },
      { id: 'checkbox', name: '复选框', icon: 'Select', category: '表单', type: 'checkbox', description: '多选' },
      { id: 'radio', name: '单选框', icon: 'CircleCheck', category: '表单', type: 'radio', description: '单选' },
      { id: 'switch', name: '开关', icon: 'Switch', category: '表单', type: 'switch', description: '布尔切换' },
      { id: 'slider', name: '滑块', icon: 'Operation', category: '表单', type: 'slider', description: '范围选择' },
      { id: 'select', name: '选择器', icon: 'ArrowDown', category: '表单', type: 'select', description: '下拉选择' },
      { id: 'rate', name: '评分', icon: 'Star', category: '表单', type: 'rate', description: '星级评分' },
      { id: 'color-picker', name: '颜色选择', icon: 'Brush', category: '表单', type: 'color-picker', description: '颜色选取' },
      { id: 'date-picker', name: '日期选择', icon: 'Calendar', category: '表单', type: 'date-picker', description: '日期选取' },
      { id: 'time-picker', name: '时间选择', icon: 'Clock', category: '表单', type: 'time-picker', description: '时间选取' },
      { id: 'upload', name: '上传', icon: 'Upload', category: '表单', type: 'upload', description: '文件上传' },
      { id: 'cascader', name: '级联选择', icon: 'DArrowRight', category: '表单', type: 'cascader', description: '层级选择' },
      { id: 'transfer', name: '穿梭框', icon: 'Sort', category: '表单', type: 'transfer', description: '数据穿梭' },
      { id: 'autocomplete', name: '自动完成', icon: 'Search', category: '表单', type: 'autocomplete', description: '输入建议' },
      { id: 'form', name: '表单容器', icon: 'Collection', category: '表单', type: 'form', description: '表单布局' },
      { id: 'rich-text', name: '富文本', icon: 'DocumentCopy', category: '表单', type: 'rich-text', description: '富文本编辑' },
    ] as UIComponent[],
  },
  {
    name: '📊 数据展示',
    components: [
      { id: 'table', name: '表格', icon: 'List', category: '数据展示', type: 'table', description: '数据表格' },
      { id: 'tag', name: '标签', icon: 'PriceTag', category: '数据展示', type: 'tag', description: '标记分类' },
      { id: 'badge', name: '徽标', icon: 'Notification', category: '数据展示', type: 'badge', description: '角标提示' },
      { id: 'progress', name: '进度条', icon: 'Loading', category: '数据展示', type: 'progress', description: '进度展示' },
      { id: 'avatar', name: '头像', icon: 'User', category: '数据展示', type: 'avatar', description: '用户头像' },
      { id: 'card', name: '卡片', icon: 'Postcard', category: '数据展示', type: 'card', description: '信息卡片' },
      { id: 'collapse', name: '折叠面板', icon: 'Fold', category: '数据展示', type: 'collapse', description: '可折叠区域' },
      { id: 'timeline', name: '时间线', icon: 'Timer', category: '数据展示', type: 'timeline', description: '时间轴' },
      { id: 'statistic', name: '统计数值', icon: 'DataAnalysis', category: '数据展示', type: 'statistic', description: '数值统计' },
      { id: 'empty', name: '空状态', icon: 'FolderOpened', category: '数据展示', type: 'empty', description: '空数据占位' },
      { id: 'skeleton', name: '骨架屏', icon: 'SemiSelection', category: '数据展示', type: 'skeleton', description: '加载占位' },
      { id: 'tree', name: '树形控件', icon: 'Guide', category: '数据展示', type: 'tree', description: '层级数据' },
      { id: 'descriptions', name: '描述列表', icon: 'Tickets', category: '数据展示', type: 'descriptions', description: '键值对展示' },
      { id: 'calendar', name: '日历', icon: 'Calendar', category: '数据展示', type: 'calendar', description: '日期面板' },
      { id: 'carousel', name: '走马灯', icon: 'PictureRounded', category: '数据展示', type: 'carousel', description: '轮播展示' },
      { id: 'list', name: '列表', icon: 'List', category: '数据展示', type: 'list', description: '数据列表' },
      { id: 'image-viewer', name: '图片查看', icon: 'ZoomIn', category: '数据展示', type: 'image-viewer', description: '图片预览' },
      { id: 'virtual-list', name: '虚拟列表', icon: 'Grid', category: '数据展示', type: 'virtual-list', description: '大数据列表' },
    ] as UIComponent[],
  },
  {
    name: '🧭 导航',
    components: [
      { id: 'menu', name: '菜单', icon: 'Menu', category: '导航', type: 'menu', description: '导航菜单' },
      { id: 'tabs', name: '标签页', icon: 'Tickets', category: '导航', type: 'tabs', description: '选项卡切换' },
      { id: 'breadcrumb', name: '面包屑', icon: 'DArrowRight', category: '导航', type: 'breadcrumb', description: '路径导航' },
      { id: 'dropdown', name: '下拉菜单', icon: 'CaretBottom', category: '导航', type: 'dropdown', description: '下拉操作' },
      { id: 'pagination', name: '分页', icon: 'More', category: '导航', type: 'pagination', description: '分页导航' },
      { id: 'steps', name: '步骤条', icon: 'SetUp', category: '导航', type: 'steps', description: '流程步骤' },
      { id: 'backtop', name: '回到顶部', icon: 'UploadFilled', category: '导航', type: 'backtop', description: '快速回顶' },
      { id: 'anchor', name: '锚点', icon: 'Aim', category: '导航', type: 'anchor', description: '页面锚点' },
      { id: 'affix', name: '固钉', icon: 'Pushpin', category: '导航', type: 'affix', description: '固定定位' },
      { id: 'navbar', name: '导航栏', icon: 'HomeFilled', category: '导航', type: 'navbar', description: '顶部导航' },
    ] as UIComponent[],
  },
  {
    name: '💬 反馈',
    components: [
      { id: 'dialog', name: '对话框', icon: 'ChatDotSquare', category: '反馈', type: 'dialog', description: '模态弹窗' },
      { id: 'drawer', name: '抽屉', icon: 'DArrowLeft', category: '反馈', type: 'drawer', description: '侧边抽屉' },
      { id: 'tooltip', name: '文字提示', icon: 'QuestionFilled', category: '反馈', type: 'tooltip', description: '悬浮提示' },
      { id: 'popover', name: '气泡卡片', icon: 'ChatRound', category: '反馈', type: 'popover', description: '气泡弹窗' },
      { id: 'alert', name: '警告', icon: 'Warning', category: '反馈', type: 'alert', description: '提示信息' },
      { id: 'loading', name: '加载', icon: 'Loading', category: '反馈', type: 'loading', description: '加载状态' },
      { id: 'result', name: '结果页', icon: 'CircleCheck', category: '反馈', type: 'result', description: '结果反馈' },
      { id: 'message', name: '消息提示', icon: 'ChatLineSquare', category: '反馈', type: 'message', description: '全局消息' },
      { id: 'notification', name: '通知', icon: 'Bell', category: '反馈', type: 'notification', description: '通知提醒' },
      { id: 'watermark', name: '水印', icon: 'Stamp', category: '反馈', type: 'watermark', description: '页面水印' },
      { id: 'progress-bar', name: '进度条', icon: 'Loading', category: '反馈', type: 'progress-bar', description: '操作进度' },
      { id: 'confirm', name: '确认框', icon: 'WarningFilled', category: '反馈', type: 'confirm', description: '确认操作' },
    ] as UIComponent[],
  },
  {
    name: '📦 容器',
    components: [
      { id: 'panel', name: '面板', icon: 'Grid', category: '容器', type: 'panel', description: '通用面板' },
      { id: 'group', name: '分组框', icon: 'Collection', category: '容器', type: 'group', description: '分组容器' },
      { id: 'container', name: '容器', icon: 'Box', category: '容器', type: 'container', description: '布局容器' },
      { id: 'scrollbar', name: '滚动条', icon: 'Rank', category: '容器', type: 'scrollbar', description: '滚动容器' },
      { id: 'layout', name: '布局', icon: 'Grid', category: '容器', type: 'layout', description: '页面布局' },
      { id: 'row', name: '行', icon: 'Rank', category: '容器', type: 'row', description: '栅格行' },
      { id: 'col', name: '列', icon: 'Sort', category: '容器', type: 'col', description: '栅格列' },
      { id: 'flex', name: '弹性布局', icon: 'Operation', category: '容器', type: 'flex', description: 'Flex布局' },
      { id: 'split-pane', name: '分割面板', icon: 'DArrowRight', category: '容器', type: 'split-pane', description: '可拖拽分割' },
      { id: 'tab-container', name: '选项卡容器', icon: 'Tickets', category: '容器', type: 'tab-container', description: '选项卡布局' },
    ] as UIComponent[],
  },
  {
    name: '📈 图表',
    components: [
      { id: 'line-chart', name: '折线图', icon: 'TrendCharts', category: '图表', type: 'line-chart', description: '趋势图表' },
      { id: 'bar-chart', name: '柱状图', icon: 'Histogram', category: '图表', type: 'bar-chart', description: '对比图表' },
      { id: 'gauge', name: '仪表盘', icon: 'Odometer', category: '图表', type: 'gauge', description: '仪表展示' },
      { id: 'ring-progress', name: '环形进度', icon: 'PieChart', category: '图表', type: 'ring-progress', description: '环形指标' },
      { id: 'pie-chart', name: '饼图', icon: 'PieChart', category: '图表', type: 'pie-chart', description: '占比图表' },
      { id: 'scatter-chart', name: '散点图', icon: 'ScatterPlot', category: '图表', type: 'scatter-chart', description: '分布图表' },
      { id: 'radar-chart', name: '雷达图', icon: 'Radar', category: '图表', type: 'radar-chart', description: '多维对比' },
      { id: 'area-chart', name: '面积图', icon: 'TrendCharts', category: '图表', type: 'area-chart', description: '面积趋势' },
      { id: 'funnel-chart', name: '漏斗图', icon: 'DataAnalysis', category: '图表', type: 'funnel-chart', description: '转化分析' },
      { id: 'heatmap', name: '热力图', icon: 'Grid', category: '图表', type: 'heatmap', description: '密度分布' },
    ] as UIComponent[],
  },
])

const filteredCategories = computed(() => {
  if (!searchKeyword.value.trim()) return categories.value
  const kw = searchKeyword.value.toLowerCase()
  return categories.value
    .map(cat => ({
      ...cat,
      components: cat.components.filter(c =>
        c.name.includes(kw) || c.type.includes(kw) || (c.description && c.description.includes(kw))
      ),
    }))
    .filter(cat => cat.components.length > 0)
})

const emit = defineEmits<{
  (e: 'drag-start', component: UIComponent): void
}>()

const onDragStart = (event: DragEvent, component: UIComponent) => {
  if (event.dataTransfer) {
    event.dataTransfer.setData('application/json', JSON.stringify(component))
    event.dataTransfer.effectAllowed = 'copy'
  }
  emit('drag-start', component)
}
</script>

<template>
  <div class="component-library">
    <div class="library-header">
      <span>组件库</span>
      <span class="component-count">{{ categories.reduce((s, c) => s + c.components.length, 0) }}</span>
    </div>
    <div class="library-search">
      <el-input
        v-model="searchKeyword"
        size="small"
        placeholder="搜索组件..."
        clearable
        :prefix-icon="'Search'"
      />
    </div>
    <div class="library-content">
      <div v-for="category in filteredCategories" :key="category.name" class="category-group">
        <div class="category-title" @click="toggleCategory(category.name)">
          <span class="collapse-arrow">{{ collapsedCategories[category.name] ? '▶' : '▼' }}</span>
          <span>{{ category.name }}</span>
          <span class="category-count">{{ category.components.length }}</span>
        </div>
        <div v-show="!collapsedCategories[category.name]" class="component-list">
          <div
            v-for="comp in category.components"
            :key="comp.id"
            class="component-item"
            draggable="true"
            @dragstart="onDragStart($event, comp)"
          >
            <el-icon class="comp-icon"><component :is="comp.icon" /></el-icon>
            <div class="comp-info">
              <span class="comp-name">{{ comp.name }}</span>
              <span v-if="comp.description" class="comp-desc">{{ comp.description }}</span>
            </div>
          </div>
        </div>
      </div>
      <div v-if="filteredCategories.length === 0" class="empty-hint">
        未找到匹配的组件
      </div>
    </div>
  </div>
</template>

<style scoped>
.component-library {
  width: 220px;
  height: 100%;
  background-color: var(--panel-bg);
  border-right: 1px solid var(--panel-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.library-header {
  padding: 10px 12px;
  font-size: 14px;
  font-weight: bold;
  color: var(--panel-header-text);
  border-bottom: 1px solid var(--panel-border);
  background-color: var(--panel-header-bg);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.component-count {
  font-size: 11px;
  color: var(--text-secondary);
  background-color: var(--bg-active);
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: normal;
}

.library-search {
  padding: 8px;
  border-bottom: 1px solid var(--panel-border);
}

.library-content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.category-group {
  margin-bottom: 2px;
}

.category-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary);
  padding: 6px 12px;
  font-weight: 500;
  cursor: pointer;
  user-select: none;
  transition: background-color 0.15s;
}

.category-title:hover {
  background-color: var(--bg-hover);
}

.collapse-arrow {
  font-size: 10px;
  width: 12px;
  color: var(--text-secondary);
}

.category-count {
  margin-left: auto;
  font-size: 10px;
  color: var(--text-secondary);
  background-color: var(--bg-active);
  padding: 1px 6px;
  border-radius: 8px;
}

.component-list {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: 2px 0;
}

.component-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 12px 7px 24px;
  cursor: grab;
  transition: all 0.15s;
  border-left: 2px solid transparent;
}

.component-item:hover {
  background-color: var(--color-primary-light);
  border-left-color: var(--color-primary);
}

.component-item:active {
  cursor: grabbing;
  background-color: var(--bg-active);
}

.comp-icon {
  font-size: 15px;
  color: var(--color-primary);
  flex-shrink: 0;
}

.comp-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  overflow: hidden;
}

.comp-name {
  font-size: 13px;
  color: var(--text-regular);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.comp-desc {
  font-size: 10px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-hint {
  text-align: center;
  padding: 24px 12px;
  font-size: 12px;
  color: var(--text-secondary);
}
</style>