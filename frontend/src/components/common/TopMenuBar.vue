<script lang="ts" setup>
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'menu-action', action: string): void
}>()

const handleMenuClick = (action: string) => {
  emit('menu-action', action)
}
</script>

<template>
  <div class="top-menu-bar">
    <div class="menu-left">
      <div class="app-logo">
        <span class="logo-icon">🐕</span>
        <span class="logo-text">易狗 IDE</span>
      </div>
    </div>
    <div class="menu-center">
      <el-menu
        mode="horizontal"
        :ellipsis="false"
        class="main-menu"
        background-color="#2c3e50"
        text-color="#ecf0f1"
        active-text-color="#409eff"
      >
        <el-sub-menu index="file">
          <template #title>文件</template>
          <el-menu-item index="new-project" @click="handleMenuClick('new-project')">新建项目</el-menu-item>
          <el-menu-item index="open-project" @click="handleMenuClick('open-project')">打开项目</el-menu-item>
          <el-menu-item index="save-project" @click="handleMenuClick('save-project')">保存项目</el-menu-item>
          <el-menu-item index="save-as" @click="handleMenuClick('save-as')">另存为</el-menu-item>
          <el-divider style="margin: 4px 0" />
          <el-menu-item index="exit" @click="handleMenuClick('exit')">退出</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="edit">
          <template #title>编辑</template>
          <el-menu-item index="undo" @click="handleMenuClick('undo')">撤销</el-menu-item>
          <el-menu-item index="redo" @click="handleMenuClick('redo')">重做</el-menu-item>
          <el-divider style="margin: 4px 0" />
          <el-menu-item index="cut" @click="handleMenuClick('cut')">剪切</el-menu-item>
          <el-menu-item index="copy" @click="handleMenuClick('copy')">复制</el-menu-item>
          <el-menu-item index="paste" @click="handleMenuClick('paste')">粘贴</el-menu-item>
          <el-menu-item index="delete" @click="handleMenuClick('delete')">删除</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="run">
          <template #title>运行</template>
          <el-menu-item index="run-project" @click="handleMenuClick('run-project')">
            <el-icon><VideoPlay /></el-icon>
            运行项目
          </el-menu-item>
          <el-menu-item index="stop-project" @click="handleMenuClick('stop-project')">
            <el-icon><VideoPause /></el-icon>
            停止项目
          </el-menu-item>
          <el-divider style="margin: 4px 0" />
          <el-menu-item index="build" @click="handleMenuClick('build')">
            <el-icon><Box /></el-icon>
            编译打包
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="component">
          <template #title>组件</template>
          <el-menu-item index="import-component" @click="handleMenuClick('import-component')">导入组件</el-menu-item>
          <el-menu-item index="export-component" @click="handleMenuClick('export-component')">导出组件</el-menu-item>
          <el-menu-item index="component-manager" @click="handleMenuClick('component-manager')">组件管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="tools">
          <template #title>工具</template>
          <el-menu-item index="settings" @click="handleMenuClick('settings')">设置</el-menu-item>
          <el-menu-item index="plugin-manager" @click="handleMenuClick('plugin-manager')">插件管理</el-menu-item>
          <el-menu-item index="theme-toggle" @click="handleMenuClick('theme-toggle')">切换主题</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="help">
          <template #title>帮助</template>
          <el-menu-item index="docs" @click="handleMenuClick('docs')">使用文档</el-menu-item>
          <el-menu-item index="about" @click="handleMenuClick('about')">关于易狗 IDE</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </div>
    <div class="menu-right">
      <span class="version-text">v1.0.0</span>
    </div>
  </div>
</template>

<style scoped>
.top-menu-bar {
  display: flex;
  align-items: center;
  height: 40px;
  background-color: #2c3e50;
  border-bottom: 1px solid #1a252f;
  user-select: none;
  flex-shrink: 0;
}

.menu-left {
  display: flex;
  align-items: center;
  padding: 0 12px;
  min-width: 180px;
}

.app-logo {
  display: flex;
  align-items: center;
  gap: 6px;
}

.logo-icon {
  font-size: 20px;
}

.logo-text {
  color: #ecf0f1;
  font-size: 14px;
  font-weight: bold;
  white-space: nowrap;
}

.menu-center {
  flex: 1;
  display: flex;
}

.main-menu {
  border-bottom: none !important;
  height: 40px;
}

.main-menu :deep(.el-sub-menu__title) {
  height: 40px;
  line-height: 40px;
  padding: 0 16px;
  font-size: 13px;
}

.main-menu :deep(.el-menu-item) {
  height: 36px;
  line-height: 36px;
  font-size: 13px;
}

.menu-right {
  display: flex;
  align-items: center;
  padding: 0 12px;
}

.version-text {
  color: #7f8c8d;
  font-size: 11px;
}
</style>