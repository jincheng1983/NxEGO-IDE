<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { BuildWithConfig, GetSupportedPlatforms, SelectSavePath } from '../../../wailsjs/go/main/App'
import { ElMessage } from 'element-plus'

interface BuildConfig {
  appName: string
  version: string
  description: string
  company: string
  copyright: string
  iconPath: string
  outputPath: string
  targetOS: string
  targetArch: string
  buildMode: string
  optimize: boolean
  stripDebug: boolean
  upxCompress: boolean
  embedAssets: boolean
  staticLink: boolean
  cgoEnabled: boolean
  raceDetector: boolean
  customLdflags: string
  customBuildTags: string
  outputType: string
}

interface Platform {
  os: string
  arch: string
  label: string
}

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'build-result', result: any): void
}>()

const props = defineProps<{
  goCode?: string
}>()

const config = ref<BuildConfig>({
  appName: '我的应用',
  version: '1.0.0',
  description: '',
  company: '',
  copyright: '',
  iconPath: '',
  outputPath: '',
  targetOS: 'windows',
  targetArch: 'amd64',
  buildMode: 'release',
  optimize: true,
  stripDebug: true,
  upxCompress: false,
  embedAssets: false,
  staticLink: false,
  cgoEnabled: false,
  raceDetector: false,
  customLdflags: '',
  customBuildTags: '',
  outputType: 'exe',
})

const platforms = ref<Record<string, string>[]>([])
const isBuilding = ref(false)
const buildProgress = ref(0)
const buildStage = ref('')
const buildMessage = ref('')

const loadPlatforms = async () => {
  try {
    platforms.value = await GetSupportedPlatforms()
  } catch (e) {
    console.error('加载平台列表失败:', e)
  }
}

const selectOutputPath = async () => {
  try {
    const path = await SelectSavePath()
    if (path) {
      config.value.outputPath = path
    }
  } catch (e) {
    console.error('选择输出路径失败:', e)
  }
}

const startBuild = async () => {
  if (!config.value.appName.trim()) {
    ElMessage.warning('请输入应用名称')
    return
  }

  if (!props.goCode) {
    ElMessage.warning('正在准备代码，请稍后重试')
    return
  }

  isBuilding.value = true
  buildProgress.value = 0
  buildStage.value = '准备中...'
  buildMessage.value = ''

  try {
    const result = await BuildWithConfig(props.goCode || '', config.value)
    isBuilding.value = false

    if (result.success) {
      buildProgress.value = 100
      buildStage.value = '完成'
      buildMessage.value = result.message
      ElMessage.success('构建成功！')
      emit('build-result', result)
    } else {
      buildProgress.value = 0
      buildStage.value = '失败'
      buildMessage.value = result.message
      ElMessage.error(result.message)
    }
  } catch (e: any) {
    isBuilding.value = false
    buildProgress.value = 0
    buildStage.value = '错误'
    buildMessage.value = String(e)
    ElMessage.error('构建失败: ' + e)
  }
}

const getPlatformLabel = (): string => {
  const p = platforms.value.find(pl => pl.os === config.value.targetOS && pl.arch === config.value.targetArch)
  return p ? p.label : `${config.value.targetOS}/${config.value.targetArch}`
}

onMounted(() => {
  loadPlatforms()
})
</script>

<template>
  <div class="build-dialog-overlay" @click.self="emit('close')">
    <div class="build-dialog">
      <div class="bd-header">
        <h3>构建配置</h3>
        <el-button circle size="small" @click="emit('close')">
          <el-icon><Close /></el-icon>
        </el-button>
      </div>

      <div class="bd-body">
        <div class="bd-section">
          <div class="bd-section-title">基本信息</div>
          <el-form label-width="80px" label-position="left" size="small">
            <el-form-item label="应用名称">
              <el-input v-model="config.appName" placeholder="输入应用名称" />
            </el-form-item>
            <el-form-item label="版本号">
              <el-input v-model="config.version" placeholder="1.0.0" />
            </el-form-item>
            <el-form-item label="描述">
              <el-input v-model="config.description" placeholder="应用描述（可选）" />
            </el-form-item>
            <el-form-item label="公司">
              <el-input v-model="config.company" placeholder="公司名称（可选）" />
            </el-form-item>
            <el-form-item label="版权">
              <el-input v-model="config.copyright" placeholder="版权信息（可选）" />
            </el-form-item>
          </el-form>
        </div>

        <div class="bd-section">
          <div class="bd-section-title">目标平台</div>
          <div class="platform-grid">
            <div
              v-for="p in platforms"
              :key="p.os + p.arch"
              class="platform-card"
              :class="{ selected: config.targetOS === p.os && config.targetArch === p.arch }"
              @click="config.targetOS = p.os; config.targetArch = p.arch"
            >
              <div class="platform-icon">
                <el-icon v-if="p.os === 'windows'" size="20"><Monitor /></el-icon>
                <el-icon v-else-if="p.os === 'linux'" size="20"><Platform /></el-icon>
                <el-icon v-else-if="p.os === 'darwin'" size="20"><Apple /></el-icon>
              </div>
              <div class="platform-label">{{ p.label }}</div>
            </div>
          </div>
        </div>

        <div class="bd-section">
          <div class="bd-section-title">输出设置</div>
          <el-form label-width="80px" label-position="left" size="small">
            <el-form-item label="输出路径">
              <div style="display: flex; gap: 8px; width: 100%">
                <el-input v-model="config.outputPath" placeholder="选择输出路径" readonly />
                <el-button @click="selectOutputPath">选择</el-button>
              </div>
            </el-form-item>
          </el-form>
        </div>

        <div class="bd-section">
          <div class="bd-section-title">编译模式</div>
          <el-form label-width="80px" label-position="left" size="small">
            <el-form-item label="构建模式">
              <el-radio-group v-model="config.buildMode" size="small">
                <el-radio-button value="debug">Debug (调试)</el-radio-button>
                <el-radio-button value="release">Release (发布)</el-radio-button>
              </el-radio-group>
              <div class="form-hint">Debug: 包含调试信息 | Release: 优化编译</div>
            </el-form-item>
            <el-form-item label="输出类型">
              <el-select v-model="config.outputType" size="small" style="width: 100%">
                <el-option label="可执行文件 (.exe)" value="exe" />
                <el-option label="动态链接库 (.dll)" value="dll" />
                <el-option label="静态库 (.lib/.a)" value="static-lib" />
                <el-option label="共享对象 (.so)" value="shared-lib" />
              </el-select>
            </el-form-item>
          </el-form>
        </div>

        <div class="bd-section">
          <div class="bd-section-title">高级编译选项</div>
          <div class="build-options">
            <div class="option-row">
              <el-checkbox v-model="config.optimize" size="small">
                优化编译 (-ldflags "-s -w")
              </el-checkbox>
              <span class="option-hint">减小体积约30%</span>
            </div>
            <div class="option-row">
              <el-checkbox v-model="config.stripDebug" size="small">
                去除调试符号
              </el-checkbox>
              <span class="option-hint">进一步减小体积</span>
            </div>
            <div class="option-row">
              <el-checkbox v-model="config.staticLink" size="small">
                静态编译 (CGO_ENABLED=0 -ldflags "-extldflags -static")
              </el-checkbox>
              <span class="option-hint">无依赖运行</span>
            </div>
            <div class="option-row">
              <el-checkbox v-model="config.cgoEnabled" size="small">
                启用 CGO
              </el-checkbox>
              <span class="option-hint">支持C代码调用</span>
            </div>
            <div class="option-row">
              <el-checkbox v-model="config.raceDetector" size="small">
                竞态检测器 (-race)
              </el-checkbox>
              <span class="option-hint">检测并发问题</span>
            </div>
            <div class="option-row">
              <el-checkbox v-model="config.upxCompress" size="small">
                UPX压缩
              </el-checkbox>
              <span class="option-hint">需安装UPX工具</span>
            </div>
            <div class="option-row">
              <el-checkbox v-model="config.embedAssets" size="small">
                嵌入资源文件
              </el-checkbox>
              <span class="option-hint">打包资源到可执行文件</span>
            </div>
          </div>
        </div>

        <div class="bd-section">
          <div class="bd-section-title">自定义参数</div>
          <el-form label-width="100px" label-position="left" size="small">
            <el-form-item label="自定义 LDFLAGS">
              <el-input
                v-model="config.customLdflags"
                placeholder="-X main.version=1.0.0"
                type="textarea"
                :rows="2"
              />
              <div class="form-hint">用于注入版本信息等变量</div>
            </el-form-item>
            <el-form-item label="Build Tags">
              <el-input
                v-model="config.customBuildTags"
                placeholder="linux,darwin,!windows"
              />
              <div class="form-hint">条件编译标签，逗号分隔</div>
            </el-form-item>
          </el-form>
        </div>

        <div class="bd-section">
          <div class="bd-section-title">版本信息注入</div>
          <div class="version-inject-info">
            <el-alert type="info" :closable="false" show-icon>
              <template #title>自动注入以下变量到程序中</template>
              <div class="inject-vars">
                <code>main.AppName = "{{ config.appName }}"</code><br>
                <code>main.Version = "{{ config.version }}"</code><br>
                <code>main.Company = "{{ config.company || 'N/A' }}"</code><br>
                <code>main.BuildTime = "当前时间"</code><br>
                <code>main.GoVersion = "Go版本号"</code>
              </div>
            </el-alert>
          </div>
        </div>

        <div v-if="isBuilding || buildMessage" class="bd-section">
          <div class="bd-section-title">构建状态</div>
          <div class="build-status">
            <el-progress
              :percentage="buildProgress"
              :status="buildProgress === 100 ? 'success' : buildProgress > 0 ? '' : 'exception'"
              :stroke-width="8"
            />
            <div class="build-stage">{{ buildStage }}</div>
            <div v-if="buildMessage" class="build-message">{{ buildMessage }}</div>
          </div>
        </div>
      </div>

      <div class="bd-footer">
        <el-button @click="emit('close')">取消</el-button>
        <el-button type="primary" @click="startBuild" :loading="isBuilding">
          <el-icon><VideoPlay /></el-icon>
          开始构建
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.build-dialog-overlay {
  position: fixed;
  inset: 0;
  background: var(--dialog-overlay);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.build-dialog {
  width: 600px;
  max-height: 85vh;
  background: var(--dialog-bg);
  border-radius: 8px;
  box-shadow: var(--shadow-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.bd-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border-lighter);
}

.bd-header h3 {
  margin: 0;
  font-size: 16px;
  color: var(--text-primary);
}

.bd-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

.bd-section {
  margin-bottom: 16px;
}

.bd-section-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-regular);
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--border-lighter);
}

.platform-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.platform-card {
  padding: 10px 8px;
  border: 1px solid var(--border-lighter);
  border-radius: 6px;
  cursor: pointer;
  text-align: center;
  transition: all 0.2s;
}

.platform-card:hover {
  border-color: var(--color-primary);
}

.platform-card.selected {
  border-color: var(--color-primary);
  background-color: var(--color-primary-light);
}

.platform-icon {
  margin-bottom: 4px;
  color: var(--text-regular);
}

.platform-card.selected .platform-icon {
  color: var(--color-primary);
}

.platform-label {
  font-size: 11px;
  color: var(--text-regular);
}

.build-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.option-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
}

.option-hint {
  font-size: 11px;
  color: var(--text-secondary);
}

.form-hint {
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 2px;
}

.version-inject-info {
  margin-top: 4px;
}

.inject-vars {
  margin-top: 8px;
  font-size: 12px;
  line-height: 1.6;
}

.inject-vars code {
  background-color: var(--bg-secondary);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Consolas', monospace;
  color: var(--color-primary);
}

.build-status {
  padding: 12px;
  background: var(--bg-secondary);
  border-radius: 6px;
}

.build-stage {
  font-size: 13px;
  color: var(--text-regular);
  margin-top: 8px;
}

.build-message {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 4px;
  white-space: pre-wrap;
  word-break: break-all;
}

.bd-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 20px;
  border-top: 1px solid var(--border-lighter);
}
</style>