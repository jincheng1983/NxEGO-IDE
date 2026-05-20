# 易狗 IDE (YiGou IDE)

一款全中文可视化 Go 语言桌面应用开发工具。

## 特性

- **可视化设计器** — 拖拽组件，所见即所得的 UI 设计
- **中文编程** — 使用中文函数名和语法编写 Go 代码
- **一键编译** — Wails v2 + WebView2，编译为原生 Windows 应用
- **AI 辅助** — 内置 AI 对话面板，智能辅助编码
- **VSCode 风格界面** — 现代化 IDE 布局，项目资源管理器、属性面板、控制台

## 技术栈

| 层 | 技术 |
|---|------|
| 框架 | [Wails v2](https://wails.io) (Go + WebView2) |
| 后端 | Go 1.22+ |
| 前端 | Vue 3 + TypeScript + Element Plus |
| 编辑器 | Monaco Editor (VSCode 同款) |
| 构建 | wails build |

## 项目结构

```
yigou-ide/
├── app.go              # Wails 应用入口
├── main.go             # 主程序入口
├── backend/            # Go 后端
│   ├── core/           # 核心逻辑（代码转换）
│   ├── models/         # 数据模型
│   ├── runtime/        # 运行时嵌入
│   └── services/       # 业务服务（构建/AI/文件/项目）
├── frontend/           # Vue 前端
│   ├── src/
│   │   ├── components/ # UI 组件（设计器/编辑器/面板等）
│   │   ├── stores/     # Pinia 状态管理
│   │   ├── composables/# 组合式函数
│   │   └── styles/     # 主题样式
│   └── package.json
├── build/              # 构建配置与资源
├── tools/              # 工具脚本
├── wails.json          # Wails 项目配置
├── go.mod / go.sum     # Go 模块依赖
└── .gitignore
```

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Windows 10+ (WebView2)

### 开发模式

```bash
cd yigou-ide
wails dev
```

### 生产构建

```bash
wails build
```

输出: `build/bin/易狗IDE.exe`

## 许可证

MIT License — 详见 [LICENSE](LICENSE) 文件