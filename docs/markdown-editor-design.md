# Markdown 编辑器设计文档

## 1. 概述
使用 Tiptap 开发一个功能完善的 Markdown 编辑器，支持实时预览和编辑功能。

## 2. 技术栈
- Vue 3 + TypeScript
- Tiptap 2.x
- TailwindCSS
- DaisyUI

## 3. 功能需求

### 3.1 基础功能
- Markdown 语法支持
  - 标题（h1-h6）
  - 粗体、斜体
  - 列表（有序、无序）
  - 代码块和行内代码（支持语法高亮）
  - 引用块（自定义样式）
  - 链接和图片（支持上传）
  - 表格
  - 任务列表
- 实时预览
- 双栏布局

### 3.2 高级功能
- 快捷键支持
- 自动保存
- 历史记录（撤销/重做）
- 图片上传
- 代码块语法高亮
- 导出功能（纯文本、HTML）

## 4. 组件设计

### 4.1 主要组件
```
MarkdownEditor/
├── index.vue                 # 主组件
├── components/
│   ├── Toolbar.vue          # 工具栏组件
│   ├── MenuBubble.vue       # 悬浮菜单组件
│   ├── CodeBlockComponent.vue # 代码块组件
│   └── ImageComponent.vue    # 图片组件
└── extensions/              # Tiptap 扩展
    ├── CustomImage.ts       # 自定义图片扩展
    └── CustomCode.ts        # 自定义代码块扩展
```

### 4.2 扩展实现

#### 4.2.1 CustomImage 扩展
- 功能：实现图片的插入、编辑和渲染
- 特性：
  - 支持图片上传
  - 支持图片标题和替代文本
  - 自定义图片组件界面
  - 支持拖拽调整大小

#### 4.2.2 CustomCode 扩展
- 功能：实现代码块的编辑和语法高亮
- 特性：
  - 支持多种编程语言
  - 实时语法高亮
  - 代码复制功能
  - 行号显示

### 4.3 工具栏功能
- 文本格式化按钮（粗体、斜体、删除线）
- 标题级别选择
- 列表类型（有序、无序、任务列表）
- 引用块
- 代码块
- 表格插入
- 链接插入
- 图片上传
- 撤销/重做

## 5. 数据流设计

### 5.1 状态管理
- 编辑器内容状态
  - 使用 v-model 双向绑定
  - 支持实时保存
- 历史记录状态
  - 使用 Tiptap 内置的历史功能
- 工具栏状态
  - 根据当前选区更新按钮状态

### 5.2 事件处理
- 内容更新事件
  - 触发 update:content 事件
  - 自动保存到本地存储
- 格式化命令事件
  - 通过 editor.commands 执行
- 文件上传事件
  - 支持拖拽和点击上传
  - 自动压缩和优化图片

## 6. 样式设计
- 使用 TailwindCSS 进行样式设计
- 支持亮色/暗色主题
- 响应式布局
- 自定义组件样式
  - 代码块样式
  - 引用块样式
  - 任务列表样式

## 7. 性能优化
- 编辑器内容懒加载
- 图片懒加载和压缩
- 代码高亮按需加载
- 防抖处理自动保存

## 8. 后续优化方向
- 协同编辑支持
- 更多 Markdown 扩展语法
- 自定义主题
- 插件系统
- 国际化支持

## 9. 使用示例

### 9.1 基本使用
```vue
<template>
  <MarkdownEditor
    v-model:content="content"
    placeholder="开始编写..."
    @save="handleSave"
  />
</template>
```

### 9.2 自定义配置
```vue
<template>
  <MarkdownEditor
    v-model:content="content"
    :autofocus="true"
    :readonly="false"
    placeholder="开始编写..."
    @save="handleSave"
  />
</template>
``` 