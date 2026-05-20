<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import {
  GetAllProjects, GetRecentProjects, GetFavoriteProjects,
  CreateProject, OpenProject, DeleteProject, ToggleFavorite,
  SearchProjects, GetProjectTemplates, SelectFolder
} from '../../../wailsjs/go/main/App'
import { ElMessage, ElMessageBox } from 'element-plus'

interface Project {
  id: number
  name: string
  description: string
  path: string
  type: string
  author: string
  version: string
  icon: string
  createdAt: string
  updatedAt: string
  lastOpened: string
  isFavorite: boolean
  tags: string
}

interface ProjectTemplate {
  id: number
  name: string
  description: string
  category: string
  icon: string
  content: string
}

const emit = defineEmits<{
  (e: 'open-project', project: Project): void
  (e: 'close'): void
}>()

const props = defineProps<{
  mustSelect?: boolean
}>()

const activeTab = ref('recent')
const projects = ref<Project[]>([])
const templates = ref<ProjectTemplate[]>([])
const searchQuery = ref('')
const showCreateDialog = ref(false)
const createForm = ref({
  name: '',
  description: '',
  type: 'window',
  author: '',
  path: '',
  templateId: 0,
})

const projectTypes = [
  { value: 'window', label: '窗口程序' },
  { value: 'console', label: '控制台程序' },
  { value: 'service', label: '服务程序' },
  { value: 'library', label: '库项目' },
]

const loadProjects = async () => {
  try {
    switch (activeTab.value) {
      case 'recent':
        projects.value = await GetRecentProjects(20)
        break
      case 'all':
        projects.value = await GetAllProjects()
        break
      case 'favorites':
        projects.value = await GetFavoriteProjects()
        break
    }
  } catch (e) {
    console.error('加载项目列表失败:', e)
  }
}

const loadTemplates = async () => {
  try {
    templates.value = await GetProjectTemplates()
  } catch (e) {
    console.error('加载模板失败:', e)
  }
}

const onSearch = async () => {
  if (searchQuery.value.trim()) {
    try {
      projects.value = await SearchProjects(searchQuery.value)
    } catch (e) {
      projects.value = []
    }
  } else {
    loadProjects()
  }
}

const selectFolder = async () => {
  try {
    const folder = await SelectFolder()
    if (folder) {
      createForm.value.path = folder
    }
  } catch (e) {
    console.error('选择文件夹失败:', e)
  }
}

const createProject = async () => {
  if (!createForm.value.name.trim()) {
    ElMessage.warning('请输入项目名称')
    return
  }
  if (!createForm.value.path.trim()) {
    ElMessage.warning('请选择项目保存路径')
    return
  }

  try {
    const project = await CreateProject(
      createForm.value.name,
      createForm.value.description,
      createForm.value.type,
      createForm.value.author,
      createForm.value.path
    )
    if (project) {
      ElMessage.success('项目创建成功！')
      showCreateDialog.value = false
      createForm.value = { name: '', description: '', type: 'window', author: '', path: '', templateId: 0 }
      loadProjects()
      emit('open-project', project)
    }
  } catch (e) {
    ElMessage.error('创建项目失败: ' + e)
  }
}

const openProject = async (project: Project) => {
  try {
    const p = await OpenProject(project.id)
    if (p) {
      emit('open-project', p)
    }
  } catch (e) {
    ElMessage.error('打开项目失败: ' + e)
  }
}

const removeProject = async (project: Project) => {
  try {
    await ElMessageBox.confirm(`确定要删除项目"${project.name}"吗？此操作不可撤销。`, '确认删除', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await DeleteProject(project.id)
    ElMessage.success('项目已删除')
    loadProjects()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败: ' + e)
    }
  }
}

const toggleFav = async (project: Project) => {
  try {
    const result = await ToggleFavorite(project.id)
    project.isFavorite = result
    if (activeTab.value === 'favorites' && !result) {
      loadProjects()
    }
  } catch (e) {
    console.error('切换收藏失败:', e)
  }
}

const formatDate = (dateStr: string): string => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const getTypeLabel = (type: string): string => {
  const t = projectTypes.find(p => p.value === type)
  return t ? t.label : type
}

onMounted(() => {
  loadProjects()
  loadTemplates()
})
</script>

<template>
  <div class="project-manager-overlay" @click.self="!props.mustSelect && emit('close')">
    <div class="project-manager">
      <div class="pm-header">
        <h2>项目管理</h2>
        <el-button v-if="!props.mustSelect" circle size="small" @click="emit('close')">
          <el-icon><Close /></el-icon>
        </el-button>
      </div>

      <div class="pm-toolbar">
        <div class="pm-tabs">
          <span :class="{ active: activeTab === 'recent' }" @click="activeTab = 'recent'; loadProjects()">最近项目</span>
          <span :class="{ active: activeTab === 'all' }" @click="activeTab = 'all'; loadProjects()">全部项目</span>
          <span :class="{ active: activeTab === 'favorites' }" @click="activeTab = 'favorites'; loadProjects()">收藏项目</span>
        </div>
        <div class="pm-actions">
          <el-input
            v-model="searchQuery"
            size="small"
            placeholder="搜索项目..."
            clearable
            @input="onSearch"
            style="width: 200px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" size="small" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            新建项目
          </el-button>
        </div>
      </div>

      <div class="pm-content">
        <div v-if="projects.length === 0" class="pm-empty">
          <el-empty description="暂无项目">
            <el-button type="primary" @click="showCreateDialog = true">创建第一个项目</el-button>
          </el-empty>
        </div>

        <div v-else class="pm-list">
          <div
            v-for="project in projects"
            :key="project.id"
            class="pm-item"
            @dblclick="openProject(project)"
          >
            <div class="pm-item-icon">
              <el-icon size="24" color="#409eff"><FolderOpened /></el-icon>
            </div>
            <div class="pm-item-info">
              <div class="pm-item-name">
                {{ project.name }}
                <el-tag size="small" type="info">{{ getTypeLabel(project.type) }}</el-tag>
              </div>
              <div class="pm-item-desc">{{ project.description || '无描述' }}</div>
              <div class="pm-item-meta">
                <span>版本 {{ project.version }}</span>
                <span v-if="project.author">作者: {{ project.author }}</span>
                <span>最后打开: {{ formatDate(project.lastOpened) }}</span>
              </div>
            </div>
            <div class="pm-item-actions">
              <el-button
                circle
                size="small"
                :type="project.isFavorite ? 'warning' : 'default'"
                @click.stop="toggleFav(project)"
              >
                <el-icon><StarFilled v-if="project.isFavorite" /><Star v-else /></el-icon>
              </el-button>
              <el-button circle size="small" type="primary" @click.stop="openProject(project)">
                <el-icon><VideoPlay /></el-icon>
              </el-button>
              <el-button circle size="small" type="danger" @click.stop="removeProject(project)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <el-dialog
        v-model="showCreateDialog"
        title="新建项目"
        width="560px"
        :close-on-click-modal="false"
      >
        <el-form :model="createForm" label-width="80px" label-position="left">
          <el-form-item label="项目名称" required>
            <el-input v-model="createForm.name" placeholder="输入项目名称" />
          </el-form-item>
          <el-form-item label="项目描述">
            <el-input
              v-model="createForm.description"
              type="textarea"
              :rows="2"
              placeholder="输入项目描述（可选）"
            />
          </el-form-item>
          <el-form-item label="项目类型">
            <el-select v-model="createForm.type" style="width: 100%">
              <el-option
                v-for="t in projectTypes"
                :key="t.value"
                :label="t.label"
                :value="t.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="作者">
            <el-input v-model="createForm.author" placeholder="输入作者名称（可选）" />
          </el-form-item>
          <el-form-item label="保存路径" required>
            <div style="display: flex; gap: 8px; width: 100%">
              <el-input v-model="createForm.path" placeholder="选择项目保存目录" readonly />
              <el-button @click="selectFolder">选择</el-button>
            </div>
          </el-form-item>
          <el-form-item label="项目模板">
            <div class="template-grid">
              <div
                v-for="tpl in templates"
                :key="tpl.id"
                class="template-card"
                :class="{ selected: createForm.templateId === tpl.id }"
                @click="createForm.templateId = tpl.id"
              >
                <div class="template-icon">
                  <el-icon size="28"><FolderOpened /></el-icon>
                </div>
                <div class="template-name">{{ tpl.name }}</div>
                <div class="template-desc">{{ tpl.description }}</div>
              </div>
            </div>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="createProject">创建项目</el-button>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<style scoped>
.project-manager-overlay {
  position: fixed;
  inset: 0;
  background: var(--dialog-overlay);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.project-manager {
  width: 800px;
  max-height: 80vh;
  background: var(--dialog-bg);
  border-radius: 8px;
  box-shadow: var(--shadow-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.pm-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-lighter);
}

.pm-header h2 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.pm-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  border-bottom: 1px solid var(--border-lighter);
}

.pm-tabs {
  display: flex;
  gap: 4px;
}

.pm-tabs span {
  padding: 6px 14px;
  font-size: 13px;
  color: var(--text-regular);
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
}

.pm-tabs span:hover {
  background-color: var(--color-primary-light);
  color: var(--color-primary);
}

.pm-tabs span.active {
  background-color: var(--color-primary);
  color: var(--text-inverse);
}

.pm-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.pm-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px 20px;
}

.pm-empty {
  display: flex;
  justify-content: center;
  padding: 40px 0;
}

.pm-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.pm-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border: 1px solid var(--border-lighter);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.pm-item:hover {
  border-color: var(--color-primary);
  background-color: var(--color-primary-light);
}

.pm-item-icon {
  flex-shrink: 0;
}

.pm-item-info {
  flex: 1;
  min-width: 0;
}

.pm-item-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.pm-item-desc {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pm-item-meta {
  display: flex;
  gap: 16px;
  font-size: 11px;
  color: var(--text-placeholder);
}

.pm-item-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  width: 100%;
}

.template-card {
  padding: 12px;
  border: 1px solid var(--border-lighter);
  border-radius: 6px;
  cursor: pointer;
  text-align: center;
  transition: all 0.2s;
}

.template-card:hover {
  border-color: var(--color-primary);
}

.template-card.selected {
  border-color: var(--color-primary);
  background-color: var(--color-primary-light);
}

.template-icon {
  margin-bottom: 6px;
}

.template-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.template-desc {
  font-size: 11px;
  color: var(--text-secondary);
}
</style>