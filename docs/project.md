# LeafNote 项目开发文档

## 项目概述
LeafNote 是一个基于 Go 语言开发的跨平台桌面笔记应用，采用 Gin + Vue3 + Tauri 技术栈，提供一个简洁、高效的个人笔记管理解决方案。

## 技术栈
- 后端（主仓库）：
  - Go + Gin 框架
  - GORM + SQLite
  - Zap 日志系统
  - Viper 配置管理
- 前端（子仓库）：
  - Vue3 + TypeScript
  - Vite + pnpm
  - Naive UI
  - Vue Router + Pinia
  - Tauri

## 项目结构
```
leafNote/                 # 主仓库（后端）
├── cmd/                  # 主要的应用程序入口
│   └── app/             # 主应用入口
├── internal/            # 私有应用程序和库代码
│   ├── config/         # 配置相关代码
│   ├── handler/        # HTTP 处理器
│   ├── model/          # 数据模型
│   └── service/        # 业务逻辑
├── pkg/                # 可以被外部应用程序使用的库代码
├── api/               # API 协议定义
├── configs/           # 配置文件
├── docs/             # 项目文档
├── test/             # 测试文件
├── web/              # 前端子仓库（Vue3 + Tauri）
│   ├── src/          # Vue3 源代码
│   │   ├── assets/   # 静态资源
│   │   ├── components/# 公共组件
│   │   ├── composables/# 组合式函数
│   │   ├── layouts/   # 布局组件
│   │   ├── router/   # 路由配置
│   │   ├── stores/   # 状态管理
│   │   ├── styles/   # 全局样式
│   │   ├── types/    # 类型定义
│   │   ├── utils/    # 工具函数
│   │   └── views/    # 页面组件
│   ├── src-tauri/    # Tauri 相关代码
│   └── public/       # 静态资源
└── README.md         # 项目说明
```

## 子仓库管理

### 初始化子仓库
```bash
# 在主仓库中初始化前端子仓库
git submodule add <frontend-repo-url> web

# 克隆项目（包含子模块）
git clone --recursive <main-repo-url>

# 更新子模块
git submodule update --remote
```

### 开发工作流
1. 主仓库开发
   ```bash
   # 后端开发
   go run cmd/app/main.go
   ```

2. 子仓库开发
   ```bash
   # 进入前端目录
   cd web

   # 安装依赖
   pnpm install

   # 开发模式
   pnpm dev

   # 构建
   pnpm build
   ```

3. 提交更改
   ```bash
   # 提交子仓库更改
   cd web
   git add .
   git commit -m "feat: update frontend"
   git push

   # 提交主仓库更改
   cd ..
   git add .
   git commit -m "feat: update submodule"
   git push
   ```

## 开发环境要求

### 后端要求
- Go 1.21+
- SQLite 3

### 前端要求
- Node.js 18+
- pnpm 8+
- Rust（Tauri 要求）
- 系统依赖：
  - Linux: `libwebkit2gtk-4.0-dev build-essential curl wget file libssl-dev libgtk-3-dev libayatana-appindicator3-dev librsvg2-dev`

## 开发进度

### 第一阶段：基础架构搭建 [进行中]
- [x] 项目目录结构设置
- [x] 基础配置系统
  - [x] 使用 Viper 实现配置加载
  - [x] 配置文件结构定义
- [x] 日志系统
  - [x] 使用 Zap 实现日志记录
  - [x] 日志轮转功能
- [x] 数据库设计与实现
  - [x] GORM + SQLite 设置
  - [x] 数据模型定义
  - [x] 自动迁移功能
- [x] 前端子仓库初始化
  - [x] Vue3 + TypeScript 设置
  - [x] Vite 配置优化
  - [x] UI 框架集成
  - [x] 自动导入配置
  - [x] 路由系统搭建
  - [x] 状态管理配置
- [ ] 后端路由和中间件设置

### 前端开发说明

1. 项目配置
   - Vite 配置：自动导入、组件解析、路径别名
   - TypeScript 配置：严格模式、路径映射
   - ESLint + Prettier：代码规范和格式化
   - 环境变量：区分开发和生产环境

2. 开发命令
   ```bash
   # 安装依赖
   pnpm install

   # 开发模式
   pnpm dev

   # 代码格式化
   pnpm format

   # 代码检查
   pnpm lint

   # 构建
   pnpm build

   # Tauri 开发模式
   pnpm tauri dev
   ```

3. 目录说明
   - `src/components/`: 可复用组件
   - `src/composables/`: 组合式函数
   - `src/layouts/`: 布局组件
   - `src/views/`: 页面组件
   - `src/stores/`: Pinia 状态管理
   - `src/types/`: TypeScript 类型定义

4. 开发规范
   - 使用 TypeScript 进行开发
   - 使用 Composition API
   - 组件命名采用 PascalCase
   - 使用 ESLint + Prettier 进行代码格式化

## 下一步计划
1. 实现后端路由系统和中间件
2. 完善前端页面组件
3. 实现前后端通信
4. 添加 Markdown 编辑器

## 更新日志
- 2024-03-19: 创建项目文档
- 2024-03-19: 更新技术栈，调整为 Gin + Vue3 + Tauri 架构
- 2024-03-19: 完成基础架构搭建
- 2024-03-19: 更新数据库实现，迁移到 GORM
- 2024-03-19: 调整项目结构，改为后端主仓库模式
- 2024-03-19: 初始化前端子仓库（Vue3 + Tauri）
- 2024-03-19: 优化前端项目配置和结构
- 2024-03-19: 更新子仓库管理流程 