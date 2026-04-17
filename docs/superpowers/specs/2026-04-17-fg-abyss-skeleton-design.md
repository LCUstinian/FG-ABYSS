# FG-ABYSS 项目骨架设计文档

> **日期**: 2026-04-17
> **状态**: 待实现
> **范围**: 项目骨架——架构、分层、基础设施、约定

---

## 背景

上一版代码存在模块耦合严重、错误处理混乱、分层边界模糊等问题，决定从零重建。本文档定义项目骨架的全部关键决策，后续所有功能开发以此为准，不重复返工。

---

## 1. 整体架构

采用**按功能分层架构（Feature-based Layered Architecture）**：

- `commands/`：薄接口层，只做参数解构和调用转发，不含业务逻辑
- `features/`：业务逻辑层，按功能模块组织，模块间不互相依赖
- `infra/`：基础设施层（DB / HTTP / 加密），可独立替换
- `AppState`：依赖注入中枢，持有所有 Service 实例

**核心约束**：Service 层不依赖 `AppHandle`，确保可单元测试。`AppHandle` 只在 `commands/` 层使用（如开窗口）。

---

## 2. 项目结构

### Rust 后端

```
src-tauri/
├── src/
│   ├── main.rs                  # 入口：注册 commands，初始化 AppState
│   ├── lib.rs                   # Tauri 库入口
│   ├── error.rs                 # 全局 AppError + Result<T>
│   ├── state.rs                 # AppState 定义与初始化
│   │
│   ├── commands/                # 薄接口层
│   │   ├── mod.rs
│   │   ├── payload.rs
│   │   ├── webshell.rs
│   │   ├── project.rs
│   │   └── console.rs
│   │
│   ├── features/                # 业务逻辑层
│   │   ├── payload/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs       # PayloadService
│   │   │   └── models.rs        # 领域模型 + DTO
│   │   ├── webshell/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs       # WebshellService
│   │   │   └── models.rs
│   │   ├── project/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs
│   │   │   └── models.rs
│   │   ├── console/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs
│   │   │   ├── file/
│   │   │   ├── database/
│   │   │   └── terminal/
│   │   ├── plugin/
│   │   │   ├── mod.rs
│   │   │   └── service.rs
│   │   └── settings/
│   │       ├── mod.rs
│   │       └── service.rs
│   │
│   └── infra/                   # 基础设施层
│       ├── mod.rs
│       ├── paths.rs             # 应用数据目录
│       ├── db.rs                # SQLite 连接 + 迁移
│       ├── http.rs              # reqwest 客户端封装
│       ├── crypto.rs            # AES-256-GCM / XOR 加密
│       ├── config.rs            # TOML 配置文件读写
│       └── logger.rs            # tracing 初始化
│
├── migrations/
│   ├── 001_init.sql             # projects, webshells 基础表
│   └── 002_payload.sql          # payloads, payload_history 表
│
└── Cargo.toml
```

### Vue 前端

```
src/
├── main.ts                      # 初始化顺序：i18n → Pinia → Router
├── App.vue                      # n-config-provider 主题绑定
│
├── layouts/
│   └── AppLayout.vue            # 标题栏 + 导航 + 内容区 + 状态栏
│
├── router/
│   └── index.ts
│
├── features/                    # 按功能组织，镜像后端
│   ├── payload/
│   │   ├── components/          # PayloadForm, PayloadPreview, ...
│   │   ├── views/PayloadView.vue
│   │   ├── store.ts             # Pinia store（功能级）
│   │   ├── api.ts               # invoke() 类型化封装
│   │   └── types.ts             # TypeScript 类型定义
│   ├── webshell/
│   ├── project/
│   ├── console/
│   │   ├── file/
│   │   ├── database/
│   │   └── terminal/
│   ├── plugin/
│   └── settings/
│
├── shared/
│   ├── components/              # 通用 UI 组件
│   ├── types/
│   │   └── error.ts             # AppError 类型
│   └── utils/
│       └── invoke.ts            # 统一 invoke 封装
│
└── i18n/
    └── index.ts
```

---

## 3. 路由与页面职责

```
/              → HomeView         首页（统计概览）
/project       → ProjectView      项目管理 + WebShell 管理
/payload       → PayloadView      载荷生成 + 历史记录
/plugin        → PluginView       插件管理
/settings      → SettingsView     设置（子路由）
/console       → ConsoleView      控制台（独立 Tauri 窗口，不在主窗口路由）
```

| 页面 | 是否远程操作 | 使用的 features |
|------|------------|----------------|
| HomeView | 否 | webshell, project, payload（只读统计）|
| ProjectView | 是 | project, webshell |
| PayloadView | 否 | payload |
| PluginView | 否 | plugin |
| SettingsView | 否 | settings |
| ConsoleView | 是 | console（file / database / terminal）|

`ProjectView` 同时使用 `project` 和 `webshell` 两个 feature——View 层可以组合多个 feature，feature 之间不互相依赖。

---

## 4. 错误处理

### Rust 侧

```rust
// src/error.rs
#[derive(Debug, thiserror::Error)]
pub enum AppError {
    #[error("数据库错误: {0}")]
    Database(#[from] rusqlite::Error),
    #[error("HTTP 错误: {0}")]
    Http(#[from] reqwest::Error),
    #[error("加密错误: {0}")]
    Crypto(String),
    #[error("WebShell 连接失败: {0}")]
    Connection(String),
    #[error("未找到: {0}")]
    NotFound(String),
    #[error("参数无效: {0}")]
    InvalidInput(String),
    #[error("IO 错误: {0}")]
    Io(#[from] std::io::Error),
}

// 序列化为前端可消费的结构
impl serde::Serialize for AppError {
    fn serialize<S>(&self, s: S) -> Result<S::Ok, S::Error>
    where S: serde::Serializer {
        use serde::ser::SerializeStruct;
        let mut state = s.serialize_struct("AppError", 2)?;
        state.serialize_field("kind", &self.kind())?;
        state.serialize_field("message", &self.to_string())?;
        state.finish()
    }
}

pub type Result<T> = std::result::Result<T, AppError>;
```

Command 统一返回 `std::result::Result<T, AppError>`，Tauri 自动序列化错误。

### TypeScript 侧

```typescript
// src/shared/types/error.ts
export interface AppError {
  kind: 'Database' | 'Http' | 'Crypto' | 'Connection' | 'NotFound' | 'InvalidInput' | 'Io'
  message: string
}

// src/shared/utils/invoke.ts
import { invoke as tauriInvoke } from '@tauri-apps/api/core'

export async function invoke<T>(cmd: string, args?: Record<string, unknown>): Promise<T> {
  try {
    return await tauriInvoke<T>(cmd, args)
  } catch (e) {
    throw e as AppError
  }
}
```

各 feature 的 `api.ts` 通过 `invoke` 封装调用，store action 捕获错误并通过 `useMessage()` 展示——组件不直接处理错误。

---

## 5. 数据层与 AppState

### AppState

```rust
// src/state.rs
pub struct AppState {
    pub payload_service:  PayloadService,
    pub webshell_service: WebshellService,
    pub project_service:  ProjectService,
    pub console_service:  ConsoleService,
    pub plugin_service:   PluginService,
    pub settings_service: SettingsService,
}
```

`main.rs` 初始化一次，通过 `tauri::Builder::manage()` 注册，所有 command 通过 `State<'_, AppState>` 注入。

**Service 不依赖 AppHandle**，只依赖 `Database`、`HttpClient`、`AppConfig` 等基础设施，保证可单元测试。

### 异步数据库策略

使用 `rusqlite`（同步）+ `tokio::task::spawn_blocking` 模式：

```rust
// src/infra/db.rs
#[derive(Clone)]
pub struct Database(Arc<Mutex<rusqlite::Connection>>);

impl Database {
    pub fn open(path: &Path) -> Result<Self> { ... }
    pub fn migrate(&self) -> Result<()> { ... }
}

// Service 中的 DB 调用方式
pub async fn find_by_id(&self, id: &str) -> Result<Webshell> {
    let db = self.db.clone();
    let id = id.to_string();
    tokio::task::spawn_blocking(move || -> Result<Webshell> {
        let conn = db.lock().map_err(|_| AppError::Database(...))?;
        // 执行同步查询
    }).await.map_err(|e| AppError::InvalidInput(e.to_string()))?
}
```

### 数据库 Schema

所有业务表从第一天加入软删除字段：

```sql
-- migrations/001_init.sql
CREATE TABLE projects (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  INTEGER NOT NULL,
    updated_at  INTEGER NOT NULL,
    deleted_at  INTEGER          -- NULL = 未删除
);

CREATE TABLE webshells (
    id           TEXT PRIMARY KEY,   -- UUID v4，Service 层自动生成
    name         TEXT NOT NULL,
    url          TEXT NOT NULL,
    password     TEXT NOT NULL,      -- AES-256-GCM 加密存储
    payload_type TEXT NOT NULL,      -- php / jsp / asp / aspx
    project_id   TEXT REFERENCES projects(id),
    status       TEXT NOT NULL DEFAULT 'inactive',
    created_at   INTEGER NOT NULL,
    updated_at   INTEGER NOT NULL,
    deleted_at   INTEGER
);

-- migrations/002_payload.sql
CREATE TABLE payloads (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    payload_type TEXT NOT NULL,
    config       TEXT NOT NULL,      -- JSON 序列化的配置
    created_at   INTEGER NOT NULL,
    deleted_at   INTEGER
);

CREATE TABLE payload_history (
    id         TEXT PRIMARY KEY,
    payload_id TEXT REFERENCES payloads(id),
    code       TEXT NOT NULL,        -- 生成的代码
    created_at INTEGER NOT NULL
);
```

所有查询默认附加 `WHERE deleted_at IS NULL`，回收站查询则 `WHERE deleted_at IS NOT NULL`。

### WebShell ID 生成

UUID 在 Service 层创建时自动生成，不由外部传入：

```rust
use uuid::Uuid;

pub async fn create(&self, input: CreateWebshellInput) -> Result<WebshellDto> {
    let id = Uuid::new_v4().to_string();
    // ...
}
```

### 密码加密存取

WebShell 密码加密后存入 DB，取出后解密使用：

```rust
// infra/crypto.rs
pub fn encrypt(plain: &str, key: &[u8; 32]) -> Result<String> {
    // AES-256-GCM，随机 IV，输出 base64(iv + ciphertext + tag)
}
pub fn decrypt(cipher: &str, key: &[u8; 32]) -> Result<String> { ... }
```

主密钥（master key）由 Argon2id 从用户配置中的随机 salt 派生，存储在 `config.toml` 的 `[security]` 节（salt 明文存储，key 每次启动重新派生，不落盘）。

---

## 6. 配置文件 vs 数据库分工

| 存 SQLite（业务数据） | 存 config.toml（用户偏好）|
|---------------------|------------------------|
| WebShell、Project、Payload、历史记录 | 主题、语言、强调色 |
| 插件信息 | 代理设置、超时、重试次数 |
| | 加密算法、密钥参数、Salt |
| | 日志级别、轮转配置 |

```toml
# config.toml 结构
[appearance]
theme = "dark"
language = "zh-CN"
accent_color = "#2080f0"

[connection]
timeout_secs = 30
retry_count = 3
proxy_enabled = false
proxy_type = "http"
proxy_host = ""
proxy_port = 7890

[security]
encryption = "aes-256-gcm"
key_rotation_days = 30
salt = ""    # 首次启动随机生成并写入

[logging]
level = "info"
max_file_size_mb = 10
max_files = 5
```

---

## 7. 应用数据目录

```rust
// src/infra/paths.rs
pub struct AppPaths {
    pub data_dir: PathBuf,
    pub db_path:  PathBuf,
    pub config:   PathBuf,
    pub logs_dir: PathBuf,
}

impl AppPaths {
    pub fn resolve(app: &AppHandle) -> Result<Self> {
        let base = app.path().app_data_dir()?;
        // Windows: %APPDATA%\fg-abyss\
        // macOS:   ~/Library/Application Support/fg-abyss/
        // Linux:   ~/.local/share/fg-abyss/
        Ok(Self {
            db_path:  base.join("data.db"),
            config:   base.join("config.toml"),
            logs_dir: base.join("logs"),
            data_dir: base,
        })
    }
}
```

所有模块从 `AppPaths` 取路径，不硬编码任何路径。

---

## 8. 日志基础设施

```rust
// src/infra/logger.rs
pub fn init(logs_dir: &Path) -> Result<tracing_appender::non_blocking::WorkerGuard> {
    std::fs::create_dir_all(logs_dir)?;
    let file_appender = tracing_appender::rolling::daily(logs_dir, "fg-abyss.log");
    let (non_blocking, guard) = tracing_appender::non_blocking(file_appender);

    tracing_subscriber::registry()
        .with(tracing_subscriber::fmt::layer().with_writer(non_blocking))
        .with(tracing_subscriber::fmt::layer().with_writer(std::io::stderr))
        .with(tracing_subscriber::EnvFilter::new("info"))
        .init();

    Ok(guard)  // 返回 guard，由 AppState 持有，确保日志在应用退出前刷盘
}
```

`AppState` 持有 `_log_guard: WorkerGuard` 字段，应用退出时自动 flush 日志缓冲区。

`main.rs` 第一行调用 `logger::init()`，之后所有模块直接用 `tracing::info!()`。

---

## 9. 多窗口策略

每个 WebShell 对应一个独立的控制台 Tauri 窗口：

- **窗口 label**：`console-{webshell_id}`（唯一，用于查找已开窗口）
- **窗口 URL**：`/console?id={webshell_id}`
- **复用逻辑**：已开启则聚焦，未开启则新建

```rust
// src/commands/console.rs
#[tauri::command]
pub async fn open_console(app: AppHandle, webshell_id: String) -> Result<()> {
    let label = format!("console-{}", webshell_id);

    if let Some(win) = app.get_webview_window(&label) {
        win.set_focus()?;
        return Ok(());
    }

    tauri::WebviewWindowBuilder::new(
        &app,
        &label,
        tauri::WebviewUrl::App(format!("/console?id={}", webshell_id).into()),
    )
    .title("控制台")
    .inner_size(1200.0, 800.0)
    .min_inner_size(900.0, 600.0)
    .decorations(false)
    .build()?;

    Ok(())
}
```

前端控制台页面通过路由 query 获取 WebShell ID：

```typescript
// src/features/console/views/ConsoleView.vue
const webshellId = useRoute().query.id as string
```

主窗口与控制台窗口通过 Tauri 事件通信，按 label 定向发送：

```rust
app.emit_to(format!("console-{}", id), "event-name", payload)?;
```

---

## 10. Tauri 配置

### tauri.conf.json（核心部分）

```json
{
  "productName": "FG-ABYSS",
  "identifier": "com.fg-abyss.app",
  "app": {
    "windows": [
      {
        "label": "main",
        "title": "FG-ABYSS",
        "width": 1400,
        "height": 900,
        "minWidth": 1024,
        "minHeight": 768,
        "decorations": false,
        "center": true,
        "visible": true
      }
    ],
    "security": {
      "csp": "default-src 'self'; style-src 'self' 'unsafe-inline'"
    }
  }
}
```

### 权限配置（capabilities）

```json
// src-tauri/capabilities/default.json
{
  "identifier": "default",
  "windows": ["main"],
  "permissions": [
    "core:default",
    "core:window:allow-create",
    "core:window:allow-set-focus",
    "fs:allow-app-data-read-recursive",
    "fs:allow-app-data-write-recursive",
    "http:default",
    "shell:allow-open"
  ]
}

// src-tauri/capabilities/console.json
{
  "identifier": "console",
  "windows": ["console-*"],
  "permissions": [
    "core:default",
    "http:default"
  ]
}
```

---

## 11. 前端初始化顺序与主题绑定

### main.ts 初始化顺序

```typescript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter } from './router'
import { i18n } from './i18n'
import App from './App.vue'

const app = createApp(App)
app.use(i18n)           // 1. i18n 最先（store 可能依赖翻译）
app.use(createPinia())  // 2. Pinia
app.use(createRouter()) // 3. Router
app.mount('#app')
```

### App.vue 主题响应式绑定

```vue
<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <router-view />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { darkTheme } from 'naive-ui'
import { useSettingsStore } from '@/features/settings/store'

const settingsStore = useSettingsStore()
const theme = computed(() =>
  settingsStore.appearance.theme === 'dark' ? darkTheme : null
)
const themeOverrides = computed(() => ({
  common: { primaryColor: settingsStore.appearance.accentColor }
}))
</script>
```

---

## 12. 前端状态管理约定

- 每个 feature 有自己的 Pinia store，store 之间不互相 import
- View 层可组合多个 store（`ProjectView` 同时用 `projectStore` + `webshellStore`）
- 组件不直接调用 `api.ts`，所有数据操作通过 store action
- 错误在 store action 内捕获，通过 `useMessage()` 通知用户
- `loading` 状态统一放 store，组件只读 `store.loading`

### Store 模板

```typescript
// src/features/{name}/store.ts
export const use{Name}Store = defineStore('{name}', {
  state: () => ({
    items: [] as {Model}[],
    loading: false,
  }),
  actions: {
    async fetchAll() {
      this.loading = true
      try {
        this.items = await {name}Api.list()
      } catch (e) {
        useMessage().error((e as AppError).message)
      } finally {
        this.loading = false
      }
    },
  },
})
```

---

## 13. 自定义标题栏

前端标题栏组件需要拖拽区域标记和窗口控制按钮：

```vue
<!-- src/layouts/components/CustomTitlebar.vue -->
<template>
  <div class="titlebar" data-tauri-drag-region>
    <div class="left" data-tauri-drag-region>
      <img src="/logo.svg" />
      <span>FG-ABYSS</span>
    </div>
    <div class="controls">
      <button @click="minimize">－</button>
      <button @click="maximize">□</button>
      <button @click="close">×</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getCurrentWindow } from '@tauri-apps/api/window'
const win = getCurrentWindow()
const minimize = () => win.minimize()
const maximize = () => win.toggleMaximize()
const close    = () => win.close()
</script>
```

---

## 14. 实现优先级

骨架阶段按以下顺序搭建，每步完成后再进入下一步：

1. **Cargo.toml + package.json** — 依赖声明
2. **AppError + 全局 Result** — 所有后续模块都依赖它，最先定义
3. **tauri.conf.json + capabilities** — 权限和窗口配置
4. **infra 层** — paths, logger, db（含迁移）, config, crypto, http
5. **AppState + main.rs 初始化** — 串联所有基础设施
6. **features 目录骨架** — 各模块空实现（struct + 方法签名）
7. **commands 层注册** — 所有 command 注册到 Tauri
8. **前端 main.ts + App.vue** — i18n、Pinia、Router、主题绑定
9. **AppLayout + CustomTitlebar** — 整体布局
10. **各 feature 的 store + api.ts 空实现** — 前后端连通验证

---

*设计文档完成，实现阶段参考本文档，不在代码中重复解释设计决策。*
