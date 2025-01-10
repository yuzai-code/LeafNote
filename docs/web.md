# LeafNote 前端开发文档

## 技术栈
- Vue 3 + TypeScript
- Vite + pnpm
- DaisyUI + TailwindCSS
- Vue Router + Pinia
- Tauri

## 项目结构
```
web/                  # 前端子仓库
├── src/             # Vue3 源代码
│   ├── assets/      # 静态资源
│   ├── components/  # 公共组件
│   ├── composables/ # 组合式函数
│   ├── layouts/     # 布局组件
│   ├── router/     # 路由配置
│   ├── stores/     # 状态管理
│   ├── styles/     # 全局样式
│   ├── types/      # 类型定义
│   ├── utils/      # 工具函数
│   └── views/      # 页面组件
├── src-tauri/      # Tauri 相关代码
└── public/         # 静态资源
```

## 开发环境要求
- Node.js 18+
- pnpm 8+
- Rust（Tauri 要求）
- 系统依赖：
  - Linux: `libwebkit2gtk-4.0-dev build-essential curl wget file libssl-dev libgtk-3-dev libayatana-appindicator3-dev librsvg2-dev`

## 开发命令
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

## UI 组件库使用

### DaisyUI 主题配置
在 `tailwind.config.js` 中配置：
```js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light", "dark"], // 启用亮色和暗色主题
  }
}
```

### 主题切换
在 `composables/useTheme.ts` 中实现：
```typescript
export function useTheme() {
  const isDark = ref(document.documentElement.getAttribute('data-theme') === 'dark')

  const toggleDark = () => {
    isDark.value = !isDark.value
    document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
  }

  return {
    isDark,
    toggleDark
  }
}
```

### 布局组件
主布局使用 DaisyUI 的 Drawer 组件实现响应式侧边栏：
```vue
<template>
  <div class="drawer lg:drawer-open">
    <input id="main-drawer" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">
      <!-- 内容区域 -->
    </div>
    <div class="drawer-side">
      <!-- 侧边栏内容 -->
    </div>
  </div>
</template>
```

### 常用组件示例

1. 按钮：
```vue
<button class="btn btn-primary">主要按钮</button>
<button class="btn btn-secondary">次要按钮</button>
<button class="btn btn-accent">强调按钮</button>
```

2. 输入框：
```vue
<input type="text" class="input input-bordered w-full" />
<div class="form-control">
  <label class="label">
    <span class="label-text">标签</span>
  </label>
  <input type="text" class="input input-bordered" />
</div>
```

3. 卡片：
```vue
<div class="card bg-base-100 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">标题</h2>
    <p>内容</p>
    <div class="card-actions justify-end">
      <button class="btn btn-primary">确认</button>
    </div>
  </div>
</div>
```

4. 标签页：
```vue
<div class="tabs tabs-boxed">
  <a class="tab" :class="{ 'tab-active': activeTab === 'tab1' }">标签1</a>
  <a class="tab" :class="{ 'tab-active': activeTab === 'tab2' }">标签2</a>
</div>
```

5. 徽章：
```vue
<div class="badge">默认</div>
<div class="badge badge-primary">主要</div>
<div class="badge badge-secondary">次要</div>
```

## 状态管理

### Pinia Store 示例
```typescript
// stores/notes.ts
import { defineStore } from 'pinia'
import type { Note } from '../types/note'

export const useNotesStore = defineStore('notes', {
  state: () => ({
    notes: [] as Note[],
    currentNote: null as Note | null,
  }),
  
  actions: {
    async fetchNotes() {
      // 获取笔记列表
    },
    async saveNote(note: Note) {
      // 保存笔记
    }
  }
})
```

## 路由配置
```typescript
// router/index.ts
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/notes',
      name: 'notes',
      component: () => import('../views/Notes.vue')
    }
  ]
})

export default router
```

## API 调用
使用 Fetch API 进行后端通信：
```typescript
// 获取数据
const fetchData = async () => {
  try {
    const res = await fetch('/api/v1/endpoint')
    const data = await res.json()
    return data
  } catch (error) {
    console.error('请求失败:', error)
  }
}

// 提交数据
const submitData = async (data: any) => {
  try {
    const res = await fetch('/api/v1/endpoint', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    })
    return await res.json()
  } catch (error) {
    console.error('提交失败:', error)
  }
}
```

## 开发规范

### 1. 组件命名
- 使用 PascalCase 命名组件文件和组件名
- 页面组件放在 `views` 目录
- 通用组件放在 `components` 目录

### 2. TypeScript 使用
- 为所有的 props 和响应式数据定义类型
- 使用 interface 而不是 type 定义对象类型
- 导出类型定义到 `types` 目录

### 3. 样式规范
- 使用 TailwindCSS 的工具类
- 组件特定样式使用 scoped style
- 全局样式定义在 `styles` 目录

### 4. Git 提交规范
- feat: 新功能
- fix: 修复问题
- docs: 文档修改
- style: 代码格式修改
- refactor: 代码重构
- test: 测试用例修改
- chore: 其他修改

## 部署

### 构建
```bash
# 构建前端
pnpm build

# 构建 Tauri 应用
pnpm tauri build
```

### 输出目录
- Web 构建输出：`dist/`
- Tauri 构建输出：`src-tauri/target/release/`

## 性能优化

### 1. 代码分割
- 使用动态导入进行路由懒加载
- 大型组件库按需导入

### 2. 资源优化
- 图片使用适当的格式和大小
- 使用 vite 的资源导入优化

### 3. 缓存策略
- 合理使用浏览器缓存
- 实现数据本地缓存

## 测试

### 单元测试
使用 Vitest 进行单元测试：
```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MyComponent from './MyComponent.vue'

describe('MyComponent', () => {
  it('renders properly', () => {
    const wrapper = mount(MyComponent)
    expect(wrapper.text()).toContain('Hello')
  })
})
```

### E2E 测试
使用 Cypress 进行端到端测试：
```typescript
describe('My App', () => {
  it('visits the app root url', () => {
    cy.visit('/')
    cy.contains('h1', 'Welcome')
  })
})
```
