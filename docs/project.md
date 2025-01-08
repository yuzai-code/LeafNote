
# LeafNote 项目开发文档

## 项目概述
LeafNote 是一个基于 Go 语言开发的跨平台桌面笔记应用，采用 Gin + Vue 3 + Tauri 技术栈，提供一个简洁、高效的个人笔记管理解决方案。

## 技术栈
- 后端（主仓库）：
  - Go + Gin 框架
  - GORM + SQLite
  - Zap 日志系统
  - Viper 配置管理
  - Air 热更新
- 前端（子仓库）：
  - Vue 3 + TypeScript
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

## 开发路线图

### 第一阶段：基础架构 [已完成]
- [x] 项目目录结构设置
- [x] 基础配置系统
  - [x] 使用 Viper 实现配置加载
  - [x] 配置文件结构定义
- [x] 日志系统
  - [x] 使用 Zap 实现日志记录
  - [x] 日志轮转功能
- [x] 数据库设计与实现
  - [x] GORM + SQLite 设置
  - [ ] 数据模型定义
  - [x] 自动迁移功能
- [x] 前端子仓库初始化
  - [x] Vue 3 + TypeScript 设置
  - [x] Vite 配置优化
  - [x] UI 框架集成
  - [x] 自动导入配置
  - [x] 路由系统搭建
  - [x] 状态管理配置
- [ ] 后端路由和中间件设置

### 第二阶段：核心功能实现 [计划中]
- [ ] 文件系统集成
  - [ ] 文件监控系统
  - [ ] 双向同步机制
  - [ ] 冲突处理
- [ ] 加密系统
  - [ ] 密钥管理
  - [ ] 端到端加密实现
  - [ ] 搜索索引加密
- [ ] 数据模型完善
  - [ ] YAML 解析器
  - [ ] 标签系统
  - [ ] 搜索索引

### 第三阶段：用户界面和体验 [计划中]
- [ ] 编辑器实现
  - [ ] Markdown 实时预览
  - [ ] YAML 编辑器
  - [ ] 快捷键支持
- [ ] 搜索功能
  - [ ] 全文搜索
  - [ ] 标签搜索
  - [ ] 高级过滤
- [ ] UI 优化
  - [ ] 响应式设计
  - [ ] 主题支持
  - [ ] 性能优化

### 第四阶段：高级功能 [计划中]
- [ ] 插件系统
- [ ] 自动备份
- [ ] 版本控制
- [ ] 导入导出
- [ ] 统计分析 
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
- Node. Js 18+
- Pnpm 8+
- Rust（Tauri 要求）
- 系统依赖：
  - Linux: `libwebkit2gtk-4.0-dev build-essential curl wget file libssl-dev libgtk-3-dev libayatana-appindicator3-dev librsvg2-dev`


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
- 2024-03-19: 更新技术栈，调整为 Gin + Vue 3 + Tauri 架构
- 2024-03-19: 完成基础架构搭建
- 2024-03-19: 更新数据库实现，迁移到 GORM
- 2024-03-19: 调整项目结构，改为后端主仓库模式
- 2024-03-19: 初始化前端子仓库（Vue 3 + Tauri）
- 2024-03-19: 优化前端项目配置和结构
- 2024-03-19: 更新子仓库管理流程 

## 数据库设计

### 核心表结构

1. Notes（笔记表）
```sql
CREATE TABLE notes (
    id          VARCHAR(36) PRIMARY KEY,    -- UUID
    title       TEXT NOT NULL,              -- 笔记标题
    content     TEXT,                       -- 加密后的笔记内容
    yaml_meta   TEXT,  -- 加密后的YAML元数据
    file_path   TEXT NOT NULL,              -- 文件路径（相对于笔记根目录）
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    deleted_at  TIMESTAMP,                  -- 软删除
    version     INTEGER NOT NULL DEFAULT 1,  -- 版本号，用于同步
    checksum    TEXT NOT NULL               -- 内容校验和，用于同步
);
```

2. Tags（标签表）
```sql
CREATE TABLE tags (
    id          VARCHAR(36) PRIMARY KEY,    -- UUID
    name        TEXT NOT NULL,              -- 标签名称
    parent_id   VARCHAR(36),                -- 父标签ID，支持多层标签
    created_at  TIMESTAMP NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES tags(id)
);
```

3. Note_Tags（笔记-标签关联表）
```sql
CREATE TABLE note_tags (
    note_id     VARCHAR(36),
    tag_id      VARCHAR(36),
    PRIMARY KEY (note_id, tag_id),
    FOREIGN KEY (note_id) REFERENCES notes(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
```

4. Categories（目录表）
```sql
CREATE TABLE categories (
    id          VARCHAR(36) PRIMARY KEY,    -- UUID
    name        TEXT NOT NULL,              -- 目录名称
    parent_id   VARCHAR(36),                -- 父目录ID
    path        TEXT NOT NULL,              -- 完整路径
    created_at  TIMESTAMP NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES categories(id)
);
```

5. Search_Index（搜索索引表）
```sql
CREATE TABLE search_index (
    id          VARCHAR(36) PRIMARY KEY,    -- UUID
    note_id     VARCHAR(36) NOT NULL,       -- 关联的笔记ID
    content     TEXT NOT NULL,              -- 加密的搜索索引内容
    type        VARCHAR(10) NOT NULL,       -- 索引类型：title/content/tag
    FOREIGN KEY (note_id) REFERENCES notes(id)
);
```

### 数据加密方案

1. 端到端加密实现：
   - 使用用户主密码生成主密钥
   - 使用 AES-256-GCM 进行内容加密
   - 每个笔记使用唯一的 IV（Initialization Vector）
   - 密钥派生使用 Argon 2 id 算法
   - 支持离线工作模式

2. 搜索索引加密：
   - 采用可搜索加密（Searchable Encryption）技术
   - 使用确定性加密保证搜索功能
   - 实现前缀搜索和模糊匹配

## 核心功能设计

### 1. 文件系统集成
- 兼容 Obsidian 的文件组织方式
- 支持实时文件系统监控和同步
- 维护文件系统与数据库的双向同步
- 支持外部编辑器修改笔记

### 2. YAML 前置元数据处理
- 自动解析和保存 YAML 前置元数据
- 支持多层级标签系统
- 标签自动补全和建议
- 元数据版本控制

### 3. 高性能搜索系统
- 实现全文模糊搜索
- 标签和元数据索引
- 支持正则表达式搜索
- 搜索结果预览
- 搜索历史记录

### 4. 数据安全
- 端到端加密
- 本地主密钥管理
- 加密搜索索引
- 安全的密钥派生和存储
- 自动备份机制

### 5. 用户界面
- 双栏布局（目录树 + 编辑器）
- 支持实时预览
- 标签云和标签树视图
- 快捷键支持
- 黑暗模式
