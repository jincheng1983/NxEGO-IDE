<script lang="ts" setup>
import { ref } from 'vue'
import { useTheme } from '../../stores/theme'

const emit = defineEmits<{
  (e: 'menu-action', action: string): void
}>()

const { currentTheme, toggleTheme } = useTheme()

const handleCommand = (command: string) => {
  emit('menu-action', command)
}
</script>

<template>
  <div class="top-menu-bar">
    <div class="menu-logo">
      <span class="logo-icon">🐕</span>
      <span class="logo-text">易狗 IDE</span>
    </div>

    <div class="menu-items">
      <el-menu
        mode="horizontal"
        :ellipsis="false"
        :background-color="'var(--menu-bg)'"
        :text-color="'var(--menu-text)'"
        :active-text-color="'var(--menu-active-text)'"
        @select="handleCommand"
      >
        <el-sub-menu index="file">
          <template #title>文件</template>
          <el-menu-item index="new-project">新建项目</el-menu-item>
          <el-menu-item index="open-project">打开项目</el-menu-item>
          <el-menu-item index="save-project">保存项目</el-menu-item>
          <el-menu-item index="save-as">另存为...</el-menu-item>
          <el-menu-item index="export-go">导出Go代码</el-menu-item>
          <el-menu-item index="exit">退出</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="edit">
          <template #title>编辑</template>
          <el-menu-item index="undo">撤销</el-menu-item>
          <el-menu-item index="redo">重做</el-menu-item>
          <el-menu-item index="delete">删除选中</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="build">
          <template #title>构建</template>
          <el-menu-item index="run-normal">运行 (F5)</el-menu-item>
          <el-menu-item index="stop-program">停止运行</el-menu-item>
          <el-menu-item divided index="build-exe">打包为EXE</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="view">
          <template #title>视图</template>
          <el-menu-item index="toggle-component-lib">组件库</el-menu-item>
          <el-menu-item index="toggle-property-panel">属性面板</el-menu-item>
          <el-menu-item index="toggle-bottom-panel">底部面板</el-menu-item>
          <el-menu-item divided index="toggle-theme">
            <span>{{ currentTheme === 'light' ? '🌙 切换到深色主题' : '☀️ 切换到浅色主题' }}</span>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="tools">
          <template #title>工具</template>
          <el-menu-item index="open-settings">
            <el-icon><Setting /></el-icon>
            软件设置...
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="help">
          <template #title>帮助</template>
          <el-menu-item index="about">关于易狗</el-menu-item>
          <el-menu-item index="docs">使用文档</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </div>

    <div class="menu-actions">
      <el-button size="small" @click="emit('menu-action', 'new-project')" title="新建项目">
        <el-icon><DocumentAdd /></el-icon>
      </el-button>
      <el-button size="small" @click="emit('menu-action', 'open-project')" title="打开项目">
        <el-icon><FolderOpened /></el-icon>
      </el-button>
      <el-button size="small" @click="emit('menu-action', 'save-project')" title="保存项目">
        <el-icon><FolderChecked /></el-icon>
      </el-button>
      <el-divider direction="vertical" style="height: 20px; margin: 0 4px;" />
      <el-button size="small" @click="emit('menu-action', 'undo')" title="撤销">
        <el-icon><RefreshLeft /></el-icon>
      </el-button>
      <el-button size="small" @click="emit('menu-action', 'redo')" title="重做">
        <el-icon><RefreshRight /></el-icon>
      </el-button>
      <el-divider direction="vertical" style="height: 20px; margin: 0 4px;" />
      <div class="run-buttons-group">
        <el-tooltip content="运行 (F5)" placement="bottom">
          <el-button size="small" type="primary" @click="emit('menu-action', 'run-normal')">
            <el-icon><VideoPlay /></el-icon>
            <span>运行</span>
          </el-button>
        </el-tooltip>
      </div>
      <el-button size="small" type="danger" @click="emit('menu-action', 'stop-program')" title="停止运行">
        <el-icon><VideoPause /></el-icon>
      </el-button>
      <el-divider direction="vertical" style="height: 20px; margin: 0 4px;" />
      <el-dropdown trigger="click" @command="(cmd: string) => emit('menu-action', cmd)">
        <el-button size="small" type="success" title="打包/构建选项">
          <el-icon><Box /></el-icon>
          <span>打包</span>
          <el-icon class="el-icon--right"><ArrowDown /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="build-exe">
              <el-icon><Box /></el-icon> 打包为EXE
            </el-dropdown-item>
            <el-dropdown-item command="build-config">
              <el-icon><Setting /></el-icon> 构建配置...
            </el-dropdown-item>
            <el-dropdown-item divided command="export-go">
              <el-icon><Document /></el-icon> 导出Go源码
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-divider direction="vertical" style="height: 20px; margin: 0 4px;" />
      <el-button size="small" @click="emit('menu-action', 'open-settings')" title="软件设置">
        <el-icon><Setting /></el-icon>
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.top-menu-bar {
  height: 36px;
  background-color: var(--menu-bg);
  display: flex;
  align-items: center;
  padding: 0 8px;
  border-bottom: 1px solid var(--border-base);
  flex-shrink: 0;
}

.menu-logo {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-right: 16px;
  flex-shrink: 0;
}

.logo-icon {
  font-size: 18px;
}

.logo-text {
  font-size: 14px;
  font-weight: bold;
  color: var(--color-primary);
}

.menu-items {
  flex: 1;
}

.menu-items :deep(.el-menu) {
  border-bottom: none !important;
  height: 36px;
}

.menu-items :deep(.el-menu--horizontal > .el-sub-menu .el-sub-menu__title) {
  height: 36px;
  line-height: 36px;
  border-bottom: none !important;
  font-size: 13px;
  padding: 0 20px 0 12px;
}

.menu-items :deep(.el-sub-menu__icon-arrow) {
  right: 6px;
  margin-top: -2px;
}

.menu-items :deep(.el-menu--horizontal > .el-menu-item) {
  height: 36px;
  line-height: 36px;
  font-size: 13px;
}

.menu-items :deep(.el-menu--horizontal .el-menu .el-menu-item) {
  font-size: 13px;
}

.menu-actions {
  flex-shrink: 0;
  margin-left: 8px;
  display: flex;
  align-items: center;
  gap: 2px;
}

.menu-actions :deep(.el-button) {
  padding: 4px 8px;
}

.menu-actions :deep(.el-button .el-icon) {
  font-size: 14px;
}

.menu-actions :deep(.el-divider--vertical) {
  border-left-color: var(--divider-color);
}

.run-buttons-group {
  display: flex;
  align-items: center;
  gap: 2px;
}

.debug-btn {
  background-color: var(--color-warning) !important;
  border-color: var(--color-warning) !important;
}

.debug-btn:hover {
  background-color: var(--color-warning) !important;
  opacity: 0.85;
}
</style>