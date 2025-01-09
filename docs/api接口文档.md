## API 接口文档

### 笔记管理接口

#### 获取笔记列表
```http
GET /api/v1/notes
```
返回所有笔记列表。

**响应示例：**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "笔记标题",
      "content": "笔记内容",
      "yaml_meta": "yaml元数据",
      "file_path": "/path/to/note.md",
      "category_id": "分类ID",
      "created_at": "2024-01-09T12:00:00Z",
      "updated_at": "2024-01-09T12:00:00Z",
      "version": 1,
      "checksum": "内容校验和",
      "category": {
        "id": "uuid",
        "name": "分类名称"
      },
      "tags": [
        {
          "id": "uuid",
          "name": "标签名称"
        }
      ]
    }
  ],
  "status": "success"
}
```

#### 创建笔记
```http
POST /api/v1/notes
```
创建新笔记。

**请求体：**
```json
{
  "title": "笔记标题",
  "content": "笔记内容",
  "yaml_meta": "yaml元数据",
  "file_path": "/path/to/note.md",
  "category_id": "分类ID",
  "tag_ids": ["标签ID1", "标签ID2"]
}
```

**响应示例：**
```json
{
  "data": {
    "id": "uuid",
    "title": "笔记标题",
    "content": "笔记内容",
    "yaml_meta": "yaml元数据",
    "file_path": "/path/to/note.md",
    "category_id": "分类ID",
    "created_at": "2024-01-09T12:00:00Z",
    "updated_at": "2024-01-09T12:00:00Z",
    "version": 1,
    "checksum": "内容校验和"
  },
  "status": "success"
}
```

#### 获取笔记详情
```http
GET /api/v1/notes/:id
```
获取指定笔记的详细信息。

**响应示例：**
```json
{
  "data": {
    "id": "uuid",
    "title": "笔记标题",
    "content": "笔记内容",
    "yaml_meta": "yaml元数据",
    "file_path": "/path/to/note.md",
    "category_id": "分类ID",
    "created_at": "2024-01-09T12:00:00Z",
    "updated_at": "2024-01-09T12:00:00Z",
    "version": 1,
    "checksum": "内容校验和",
    "category": {
      "id": "uuid",
      "name": "分类名称"
    },
    "tags": [
      {
        "id": "uuid",
        "name": "标签名称"
      }
    ]
  },
  "status": "success"
}
```

#### 更新笔记
```http
PUT /api/v1/notes/:id
```
更新指定笔记的信息。

**请求体：**
```json
{
  "title": "新笔记标题",
  "content": "新笔记内容",
  "yaml_meta": "新yaml元数据",
  "category_id": "新分类ID",
  "tag_ids": ["新标签ID1", "新标签ID2"]
}
```

**响应示例：**
```json
{
  "data": {
    "id": "uuid",
    "title": "新笔记标题",
    "content": "新笔记内容",
    "yaml_meta": "新yaml元数据",
    "file_path": "/path/to/note.md",
    "category_id": "新分类ID",
    "created_at": "2024-01-09T12:00:00Z",
    "updated_at": "2024-01-09T12:00:00Z",
    "version": 2,
    "checksum": "新内容校验和"
  },
  "status": "success"
}
```

#### 删除笔记
```http
DELETE /api/v1/notes/:id
```
删除指定笔记。

**响应示例：**
```json
{
  "message": "删除成功",
  "status": "success"
}
```

**错误响应：**
```json
{
  "error": "笔记不存在",
  "status": "error"
}
```

### 标签管理接口

#### 获取标签列表
```http
GET /api/v1/tags
```
返回所有顶级标签及其子标签。

**响应示例：**
```json
[
  {
    "id": "uuid",
    "name": "标签名称",
    "parent_id": null,
    "children": [
      {
        "id": "child-uuid",
        "name": "子标签名称",
        "parent_id": "uuid"
      }
    ]
  }
]
```

#### 创建标签
```http
POST /api/v1/tags
```
创建新标签。

**请求体：**
```json
{
  "name": "标签名称",
  "parent_id": "父标签ID"  // 可选
}
```

#### 获取标签详情
```http
GET /api/v1/tags/:id
```
获取指定标签的详细信息。

**响应示例：**
```json
{
  "id": "uuid",
  "name": "标签名称",
  "parent_id": "父标签ID",
  "children": []
}
```

#### 更新标签
```http
PUT /api/v1/tags/:id
```
更新指定标签的信息。

**请求体：**
```json
{
  "name": "新标签名称",
  "parent_id": "新父标签ID"  // 可选
}
```

#### 删除标签
```http
DELETE /api/v1/tags/:id
```
删除指定标签。如果标签有子标签，需要先删除子标签。

**响应示例：**
```json
{
  "message": "删除成功"
}
```

### 目录管理接口

#### 获取目录列表
```http
GET /api/v1/categories
```
返回所有顶级目录及其子目录。

**响应示例：**
```json
[
  {
    "id": "uuid",
    "name": "目录名称",
    "parent_id": null,
    "path": "/目录名称",
    "children": [
      {
        "id": "child-uuid",
        "name": "子目录名称",
        "parent_id": "uuid",
        "path": "/目录名称/子目录名称"
      }
    ]
  }
]
```

#### 创建目录
```http
POST /api/v1/categories
```
创建新目录。

**请求体：**
```json
{
  "name": "目录名称",
  "parent_id": "父目录ID"  // 可选
}
```

**响应示例：**
```json
{
  "id": "uuid",
  "name": "目录名称",
  "parent_id": "父目录ID",
  "path": "/父目录名称/目录名称",
  "children": []
}
```

#### 获取目录详情
```http
GET /api/v1/categories/:id
```
获取指定目录的详细信息。

**响应示例：**
```json
{
  "id": "uuid",
  "name": "目录名称",
  "parent_id": "父目录ID",
  "path": "/父目录名称/目录名称",
  "children": []
}
```

#### 更新目录
```http
PUT /api/v1/categories/:id
```
更新指定目录的信息。

**请求体：**
```json
{
  "name": "新目录名称",
  "parent_id": "新父目录ID"  // 可选
}
```

**响应示例：**
```json
{
  "id": "uuid",
  "name": "新目录名称",
  "parent_id": "新父目录ID",
  "path": "/新父目录名称/新目录名称",
  "children": []
}
```

#### 删除目录
```http
DELETE /api/v1/categories/:id
```
删除指定目录。如果目录有子目录或笔记，需要先删除子目录和笔记。

**响应示例：**
```json
{
  "message": "删除成功"
}
```

**错误响应：**
```json
{
  "error": "请先删除子目录"
}
```
或
```json
{
  "error": "请先删除目录下的笔记"
}
```
