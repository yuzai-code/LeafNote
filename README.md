# LeafNote

LeafNote 是一个跨平台的桌面笔记应用，使用 Go + Gin + Vue3 + Tauri 技术栈开发。

## 项目结构

本项目采用主仓库（后端）+ 子仓库（前端）的方式管理：
- 主仓库：Go + Gin 后端服务
- 子仓库 `web/`：Vue3 + Tauri 前端应用

## 开发环境要求

- Go 1.21+
- Node.js 18+
- Rust (用于 Tauri)
- Git

## 快速开始

1. 克隆项目及其子模块
```bash
git clone --recursive https://github.com/yourusername/leafNote.git
cd leafNote
```

2. 初始化后端
```bash
go mod tidy
```

3. 初始化前端
```bash
cd web
npm install
```

## 开发指南

详细的开发文档请查看 `docs/project.md`。

## 许可证

MIT License 