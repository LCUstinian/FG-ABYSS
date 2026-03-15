# 时间字段标准规范

## 📋 概述

本文档定义了 FG-ABYSS 项目中所有时间字段的标准规范，确保在整个项目中时间字段的命名、存储、传输和显示保持一致性。

## 🎯 适用范围

本规范适用于：
- 后端 Go 模型定义
- 数据库表结构设计
- 前端 TypeScript 接口定义
- API 数据传输
- UI 时间显示

## 📐 时间字段标准

### 1. 命名规范

#### 后端 Go 模型
```go
// ✅ 正确：使用 CreatedAt 和 UpdatedAt
type Model struct {
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-"`
}

// ❌ 错误：避免使用以下命名
type Model struct {
    CreateTime string  // 不使用
    UpdateTime string  // 不使用
    CreatedAt  string  // 类型错误
}
```

#### 前端 TypeScript 接口
```typescript
// ✅ 正确：使用 createdAt 和 updatedAt
interface Entity {
  id: string
  createdAt: string
  updatedAt: string
}

// ❌ 错误：避免使用以下命名
interface Entity {
  createTime: string   // 不使用
  updateTime: string   // 不使用
  created_at: string   // 不使用 snake_case
}
```

### 2. 数据类型

#### 后端（Go）
- **类型**: `time.Time`
- **GORM 标签**: `gorm:"type:text"`
- **JSON 标签**: `json:"createdAt"` / `json:"updatedAt"`
- **自动管理**: GORM 自动处理时间戳的创建和更新

```go
type WebShell struct {
    ID        string    `gorm:"type:text;primaryKey" json:"id"`
    Name      string    `gorm:"type:text" json:"name"`
    CreatedAt time.Time `gorm:"type:text" json:"createdAt"`
    UpdatedAt time.Time `gorm:"type:text" json:"updatedAt"`
}
```

#### 前端（TypeScript）
- **类型**: `string`（ISO 8601 格式）
- **格式**: `YYYY-MM-DDTHH:mm:ss.sssZ` 或 `YYYY-MM-DD HH:mm:ss`
- **处理**: 使用工具函数统一格式化

```typescript
interface WebShell {
  id: string
  name: string
  createdAt: string  // ISO 8601 或标准格式
  updatedAt: string
}
```

### 3. 存储格式

#### 数据库存储
- **SQLite**: TEXT 类型
- **格式**: RFC3339 / ISO 8601
- **示例**: `2024-01-15T10:30:00.000Z`

#### JSON 序列化
- **Go 后端**: 自动序列化为 RFC3339 格式
- **示例**: `"2024-01-15T10:30:00Z"`

### 4. 时区处理

#### 后端处理
```go
// ✅ 正确：使用 UTC 时间
now := time.Now().UTC()

// ✅ 正确：转换为特定时区
loc, _ := time.LoadLocation("Asia/Shanghai")
localTime := time.Now().In(loc)

// ❌ 错误：避免使用本地时间而不标注时区
```

#### 前端处理
```typescript
// ✅ 正确：使用 Date 对象自动处理时区
const date = new Date(isoString)
const localTime = date.toLocaleString()

// ✅ 正确：使用工具函数
import { formatTime } from '@/utils/formatTime'
const formatted = formatTime(isoString)
```

### 5. 显示格式

#### 标准时间格式
```typescript
import { formatTime } from '@/utils/formatTime'

// 格式：YYYY-MM-DD HH:mm:ss
// 示例：2024-01-15 10:30:00
formatTime(date)
```

#### 日期格式
```typescript
import { formatDate } from '@/utils/formatTime'

// 格式：YYYY-MM-DD
// 示例：2024-01-15
formatDate(date)
```

#### 相对时间
```typescript
import { formatRelativeTime } from '@/utils/formatTime'

// 示例：刚刚、5 分钟前、1 小时前、3 天前
formatRelativeTime(date)
```

#### 详细时间（含毫秒）
```typescript
import { formatTimeWithMs } from '@/utils/formatTime'

// 格式：YYYY-MM-DD HH:mm:ss.SSS
// 示例：2024-01-15 10:30:00.123
formatTimeWithMs(date)
```

## 🔧 工具函数

### formatTime.ts

位置：`frontend/src/utils/formatTime.ts`

#### 可用函数

| 函数名 | 说明 | 格式 | 示例 |
|--------|------|------|------|
| `formatTime()` | 标准时间格式化 | YYYY-MM-DD HH:mm:ss | `2024-01-15 10:30:00` |
| `formatDate()` | 日期格式化 | YYYY-MM-DD | `2024-01-15` |
| `formatRelativeTime()` | 相对时间 | 刚刚/分钟前/小时前 | `5 分钟前` |
| `formatTimeWithMs()` | 详细时间（含毫秒） | YYYY-MM-DD HH:mm:ss.SSS | `2024-01-15 10:30:00.123` |
| `formatTimeRange()` | 时间范围 | 开始时间 至 结束时间 | `2024-01-15 至 2024-01-20` |

#### 使用示例

```typescript
import { 
  formatTime, 
  formatDate, 
  formatRelativeTime,
  formatTimeWithMs 
} from '@/utils/formatTime'

const now = new Date()

// 标准格式
console.log(formatTime(now)) // "2024-01-15 10:30:00"

// 日期格式
console.log(formatDate(now)) // "2024-01-15"

// 相对时间
console.log(formatRelativeTime(now)) // "刚刚"

// 详细时间
console.log(formatTimeWithMs(now)) // "2024-01-15 10:30:00.123"
```

## 📊 模型示例

### WebShell 模型

#### 后端（Go）
```go
type WebShell struct {
    ID        string    `gorm:"type:text;primaryKey" json:"id"`
    ProjectID string    `gorm:"type:text;index" json:"projectId"`
    Url       string    `gorm:"type:text;not null" json:"url"`
    Payload   string    `gorm:"type:text" json:"payload"`
    Cryption  string    `gorm:"type:text" json:"cryption"`
    Encoding  string    `gorm:"type:text" json:"encoding"`
    ProxyType string    `gorm:"type:text" json:"proxyType"`
    Remark    string    `gorm:"type:text" json:"remark"`
    Status    string    `gorm:"type:text" json:"status"`
    CreatedAt time.Time `gorm:"type:text" json:"createdAt"`
    UpdatedAt time.Time `gorm:"type:text" json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

#### 前端（TypeScript）
```typescript
interface WebShell {
  id: string
  projectId: string
  url: string
  payload: string
  cryption: string
  encoding: string
  proxyType: string
  remark: string
  status: string
  createdAt: string
  updatedAt: string
}
```

### Project 模型

#### 后端（Go）
```go
type Project struct {
    ID          string    `gorm:"type:text;uniqueIndex;not null" json:"id"`
    Name        string    `gorm:"type:text;uniqueIndex;not null" json:"name"`
    Description string    `gorm:"type:text" json:"description"`
    Status      int       `gorm:"type:integer;default:0" json:"status"`
    CreatedAt   time.Time `gorm:"type:text" json:"createdAt"`
    UpdatedAt   time.Time `gorm:"type:text" json:"updatedAt"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

#### 前端（TypeScript）
```typescript
interface Project {
  id: string
  name: string
  description: string
  status: number
  createdAt: string
  updatedAt: string
}
```

## 🎨 UI 展示规范

### 表格中的时间列
```vue
<template>
  <table>
    <thead>
      <tr>
        <th>创建时间</th>
        <th>更新时间</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>{{ formatTime(item.createdAt) }}</td>
        <td>{{ formatTime(item.updatedAt) }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script setup lang="ts">
import { formatTime } from '@/utils/formatTime'
</script>
```

### 卡片中的时间显示
```vue
<template>
  <div class="card">
    <div class="time-info">
      <span class="label">创建时间：</span>
      <span class="value">{{ formatTime(item.createdAt) }}</span>
    </div>
    <div class="time-info">
      <span class="label">更新时间：</span>
      <span class="value">{{ formatTime(item.updatedAt) }}</span>
    </div>
  </div>
</template>
```

## ✅ 检查清单

在提交代码前，请检查以下内容：

- [ ] 所有时间字段使用 `createdAt` / `updatedAt` 命名
- [ ] 后端使用 `time.Time` 类型
- [ ] 前端使用 `string` 类型（ISO 格式）
- [ ] 时间显示使用 `formatTime()` 等工具函数
- [ ] 不使用硬编码的时间格式化逻辑
- [ ] 时区处理正确（后端使用 UTC，前端自动转换）
- [ ] 数据库迁移正确执行

## 📝 迁移指南

### 从旧格式迁移

如果项目中存在旧的时间字段命名（如 `createTime` / `updateTime`），请按以下步骤迁移：

1. **更新后端模型**
   ```go
   // 旧代码
   CreateTime string `json:"createTime"`
   UpdateTime string `json:"updateTime"`
   
   // 新代码
   CreatedAt time.Time `json:"createdAt"`
   UpdatedAt time.Time `json:"updatedAt"`
   ```

2. **更新前端接口**
   ```typescript
   // 旧代码
   interface Model {
     createTime: string
     updateTime: string
   }
   
   // 新代码
   interface Model {
     createdAt: string
     updatedAt: string
   }
   ```

3. **更新数据库**
   ```sql
   -- 如果需要，执行数据库迁移
   ALTER TABLE webshells ADD COLUMN createdAt TEXT;
   ALTER TABLE webshells ADD COLUMN updatedAt TEXT;
   -- 迁移数据后删除旧列
   ```

4. **更新时间显示**
   ```vue
   <!-- 旧代码 -->
   <td>{{ item.createTime }}</td>
   
   <!-- 新代码 -->
   <td>{{ formatTime(item.createdAt) }}</td>
   ```

## 🔗 相关资源

- [Go time 包文档](https://pkg.go.dev/time)
- [TypeScript Date 对象](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/Date)
- [ISO 8601 标准](https://zh.wikipedia.org/wiki/ISO_8601)
- [GORM 时间戳处理](https://gorm.io/docs/models.html#Fields-Tags)

## 📅 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0 | 2024-01-15 | 初始版本，统一时间字段标准 |

---

**制定**: FG-ABYSS Team  
**审核**: FG-ABYSS Team  
**生效日期**: 2024-01-15
