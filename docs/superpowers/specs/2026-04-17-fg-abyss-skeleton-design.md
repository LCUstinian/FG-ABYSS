# FG-ABYSS 项目骨架设计文档

> **日期**: 2026-04-17
> **状态**: 设计完成，待实现
> **范围**: 项目骨架——架构、分层、基础设施、协议、约定
> **版本**: v4.0（含全部 36 轮 brainstorm 改进，目标：超越哥斯拉+冰蝎，成为下一代 WebShell 管理工具）

---

## 背景

上一版代码存在模块耦合严重、错误处理混乱、分层边界模糊等问题，决定从零重建。本文档是骨架设计的唯一权威来源，所有功能开发以此为准。

---

## 1. 整体架构

采用**按功能分层架构（Feature-based Layered Architecture）**：

- `commands/`：薄接口层，只做参数解构和调用转发，不含业务逻辑
- `features/`：业务逻辑层，按功能模块组织，模块间不互相依赖
  - `features/*/repo.rs`：Repo trait（测试边界）+ `DbXxxRepo`（infra 适配器），知道业务 model，不暴露给其他 feature
- `infra/`：基础设施层（DB / HTTP / 加密 / 审计 / C2 Profile），不了解业务 model
- `AppState`：依赖注入中枢，在此实例化所有 `DbXxxRepo` → `XxxService`，注册到 Tauri

**核心约束**：
- Service 层不依赖 `AppHandle`，确保可单元测试
- `AppHandle` 只在 `commands/` 层使用（开窗口、发 Tauri 事件）
- 所有基础设施依赖通过 trait 定义，测试时可 mock（见第 28 节）
- 目标能力：超越哥斯拉的流量伪装、编码链、多语言混淆、审计日志

---

## 2. 项目结构

### Rust 后端

```
src-tauri/
├── build.rs                     # 预编译 Java/ASPX 内置插件
├── src/
│   ├── main.rs                  # 入口：bootstrap → AppState → 注册 commands
│   ├── lib.rs                   # Tauri 库入口
│   ├── error.rs                 # 全局 AppError + Result<T>
│   ├── state.rs                 # AppState 定义与初始化
│   │
│   ├── commands/
│   │   ├── mod.rs
│   │   ├── webshell.rs
│   │   ├── project.rs
│   │   ├── payload.rs
│   │   ├── console.rs
│   │   ├── plugin.rs
│   │   └── settings.rs
│   │
│   ├── features/
│   │   ├── webshell/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs       # WebshellService（依赖 WebshellRepo trait）
│   │   │   ├── repo.rs          # WebshellRepo trait + DbWebshellRepo（包裹 Database::call）
│   │   │   ├── models.rs        # Webshell, WebshellDto, CreateWebshellInput
│   │   │   ├── session.rs       # WebshellSession（两阶段协议状态）
│   │   │   └── queue.rs         # per-WebShell 请求队列
│   │   ├── project/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs
│   │   │   ├── repo.rs          # ProjectRepo trait + DbProjectRepo
│   │   │   └── models.rs
│   │   ├── payload/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs       # PayloadService
│   │   │   ├── repo.rs          # PayloadRepo trait + DbPayloadRepo
│   │   │   ├── models.rs
│   │   │   ├── generator/       # PHP/JSP/ASP/ASPX 代码生成
│   │   │   │   ├── mod.rs
│   │   │   │   ├── php.rs
│   │   │   │   ├── jsp.rs
│   │   │   │   ├── asp.rs
│   │   │   │   └── aspx.rs
│   │   │   └── obfuscator/      # 混淆引擎
│   │   │       ├── mod.rs       # ObfuscateEngine + Obfuscator trait
│   │   │       ├── php.rs       # tree-sitter AST 变换
│   │   │       ├── jsp.rs       # 模板替换 + Java字符串编码
│   │   │       ├── asp.rs       # VBScript 技巧
│   │   │       └── aspx.rs      # C# 反射化 + 字符串分片
│   │   ├── console/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs       # ConsoleService
│   │   │   ├── file/
│   │   │   ├── database/
│   │   │   └── terminal/
│   │   ├── plugin/
│   │   │   ├── mod.rs
│   │   │   ├── service.rs
│   │   │   └── repo.rs          # PluginRepo trait + DbPluginRepo
│   │   └── settings/
│   │       ├── mod.rs
│   │       └── service.rs
│   │
│   └── infra/
│       ├── mod.rs
│       ├── paths.rs             # 应用数据目录
│       ├── db.rs                # SQLite async 连接池（tokio-rusqlite）+ 迁移
│       ├── http.rs              # reqwest 客户端 + 熔断器
│       ├── crypto.rs            # AES-256-GCM / Argon2id / Sensitive<T>
│       ├── config.rs            # TOML 配置文件读写 + schema 迁移
│       ├── logger.rs            # tracing 初始化，返回 WorkerGuard
│       ├── audit.rs             # 审计日志写入
│       └── c2_profile.rs        # C2 流量伪装 Profile
│
├── migrations/
│   ├── 001_init.sql       # projects + webshells
│   ├── 002_plugins.sql    # plugins
│   ├── 003_payload.sql    # payloads + payload_history（含 webshell_id 关联）
│   └── 004_audit.sql      # audit_events
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
│   ├── AppLayout.vue
│   └── components/
│       ├── CustomTitlebar.vue
│       ├── NavSidebar.vue
│       └── StatusBar.vue
│
├── router/index.ts
│
├── features/
│   ├── webshell/
│   │   ├── components/
│   │   ├── views/
│   │   ├── store.ts
│   │   ├── api.ts
│   │   └── types.ts
│   ├── project/
│   ├── payload/
│   ├── console/
│   │   ├── file/
│   │   ├── database/
│   │   └── terminal/
│   ├── plugin/
│   └── settings/
│
├── shared/
│   ├── components/
│   ├── types/
│   │   ├── error.ts             # AppError 类型
│   │   └── loading.ts           # LoadingMap 类型
│   └── utils/
│       ├── invoke.ts            # 统一 invoke 封装
│       └── clipboard.ts         # 剪贴板 + 敏感数据自动清除
│
└── i18n/index.ts
```

---

## 3. 路由与页面职责

```
/              → HomeView         首页（统计概览）
/project       → ProjectView      项目管理 + WebShell 管理
/payload       → PayloadView      载荷生成 + 历史记录
/plugin        → PluginView       插件管理
/settings      → SettingsView     设置（外观/连接/安全/日志）
/console       → ConsoleView      控制台（独立 Tauri 窗口）
/unlock        → UnlockView       主密码解锁（主密码功能启用时）
```

**ConsoleView 获取 webshell_id**：控制台窗口 URL 为 `tauri://localhost/#/console?id={webshell_id}`（Hash 路由模式），必须用 `useRoute().query.id` 读取，不能用 `useRoute().params.id`（params 只用于 `/console/:id` 动态路由模式）：

```typescript
// features/console/ConsoleView.vue
import { useRoute } from 'vue-router'
const route = useRoute()
const webshellId = route.query.id as string  // "?id=xxx" → query.id
if (!webshellId) throw new Error('ConsoleView: missing webshell id in URL query')
```

`ProjectView` 组合 `projectStore` + `webshellStore`，View 层可组合多 store，feature 之间不互相 import。

**WebShell 列表展示模式**：仅使用**表格模式**（`n-data-table`），不提供树形视图。树形模式对渗透场景不实用（WebShell 无层级关系，project 分组通过列过滤即可），且维护两套视图成本高——此为有意的设计取舍。

---

## 4. 错误处理

### Rust 侧

```rust
// src/error.rs
#[derive(Debug, thiserror::Error)]
pub enum AppError {
    #[error("数据库错误: {0}")]
    Database(#[from] rusqlite::Error),
    #[error("数据库连接错误: {0}")]
    DbConnect(String),
    #[error("HTTP 错误: {0}")]
    Http(#[from] reqwest::Error),
    #[error("加密错误: {0}")]
    Crypto(String),
    #[error("WebShell 连接失败: {0}")]
    Connection(String),
    #[error("WebShell 响应验证失败: {0}")]
    InvalidResponse(String),
    #[error("熔断器开启: {0}")]
    CircuitOpen(String),
    #[error("未找到: {0}")]
    NotFound(String),
    #[error("参数无效: {0}")]
    InvalidInput(String),
    #[error("IO 错误: {0}")]
    Io(#[from] std::io::Error),
    #[error("序列化错误: {0}")]
    Serialize(#[from] serde_json::Error),
    #[error("应用已锁定")]
    Locked,
    #[error("插件错误: {0}")]
    Plugin(String),
    #[error("WebShell 需要重新部署: {0}")]
    NeedsRedeploy(String),
    #[error("内存马会话已失效（服务端重启导致）")]
    MemShellExpired,
}

impl AppError {
    pub fn kind(&self) -> &'static str {
        match self {
            Self::Database(_) | Self::DbConnect(_) => "Database",
            Self::Http(_)                     => "Http",
            Self::Crypto(_)                   => "Crypto",
            Self::Connection(_)               => "Connection",
            Self::InvalidResponse(_)          => "InvalidResponse",
            Self::CircuitOpen(_)              => "CircuitOpen",
            Self::NotFound(_)                 => "NotFound",
            Self::InvalidInput(_)             => "InvalidInput",
            Self::Io(_)                       => "Io",
            Self::Serialize(_)                => "Serialize",
            Self::Locked                      => "Locked",
            Self::Plugin(_)                   => "Plugin",
            Self::NeedsRedeploy(_)            => "NeedsRedeploy",
            Self::MemShellExpired             => "MemShellExpired",
        }
    }
}

impl serde::Serialize for AppError {
    fn serialize<S>(&self, s: S) -> Result<S::Ok, S::Error>
    where S: serde::Serializer {
        use serde::ser::SerializeStruct;
        let mut state = s.serialize_struct("AppError", 2)?;
        state.serialize_field("kind", self.kind())?;
        state.serialize_field("message", &self.to_string())?;
        state.finish()
    }
}

pub type Result<T> = std::result::Result<T, AppError>;
```

**Tauri V2 兼容性说明**：Tauri V2 commands 完全支持 `-> Result<T, AppError>`，只要 `AppError: serde::Serialize`。我们的自定义 `impl serde::Serialize` 输出 `{ kind, message }` 格式，前端 `catch` 收到的就是这个对象。**无需改成 `Result<T, String>`**，当前设计正确。

### TypeScript 侧

```typescript
// shared/types/error.ts
export interface AppError {
  kind: 'Database' | 'Http' | 'Crypto' | 'Connection' | 'InvalidResponse'
      | 'CircuitOpen' | 'NotFound' | 'InvalidInput' | 'Io' | 'Serialize'
      | 'Locked' | 'Plugin' | 'NeedsRedeploy' | 'MemShellExpired'
  message: string
}

// shared/utils/invoke.ts
import { invoke as tauriInvoke } from '@tauri-apps/api/core'

export async function invoke<T>(cmd: string, args?: Record<string, unknown>): Promise<T> {
  try {
    return await tauriInvoke<T>(cmd, args)
  } catch (e) {
    throw e as AppError
  }
}
```

---

## 5. AppState 与数据层

### AppState

```rust
// src/state.rs
pub struct AppState {
    pub paths:             AppPaths,               // 数据目录路径集合（get_app_info 读取）
    pub webshell_service:  Arc<WebshellService>,   // Arc：BatchService 也持有引用
    pub project_service:   ProjectService,
    pub payload_service:   PayloadService,
    pub console_service:   ConsoleService,
    pub batch_service:     BatchService,
    pub plugin_service:    PluginService,
    pub settings_service:  SettingsService,
    pub audit_log:         AuditLog,      // infra 层组件，非 feature service
    pub is_locked:         AtomicBool,    // 主密码锁定状态（可选功能）
    _log_guard:            WorkerGuard,   // 持有以确保日志刷盘
}
```

`main.rs` 初始化一次，通过 `tauri::Builder::manage()` 注册，所有 command 通过 `State<'_, AppState>` 注入。

### 数据库连接

使用 `tokio-rusqlite`——rusqlite 的 async 原生包装，与 tokio 完全兼容，无需 `spawn_blocking`：

```rust
// infra/db.rs
#[derive(Clone)]
pub struct Database(Arc<tokio_rusqlite::Connection>);

impl Database {
    pub async fn open(path: &Path) -> Result<Self> {
        let conn = tokio_rusqlite::Connection::open(path).await
            .map_err(|e| AppError::DbConnect(e.to_string()))?;
        conn.call(|c| {
            c.execute_batch("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;")?;
            Ok(())
        }).await.map_err(|e| AppError::DbConnect(e.to_string()))?;
        Ok(Self(Arc::new(conn)))
    }

    // 所有 DB 操作通过 call() 在专用线程执行，不阻塞 tokio 线程池
    pub async fn call<F, T>(&self, f: F) -> Result<T>
    where
        F: FnOnce(&rusqlite::Connection) -> rusqlite::Result<T> + Send + 'static,
        T: Send + 'static,
    {
        self.0.call(f).await.map_err(|e| AppError::DbConnect(e.to_string()))
    }

    /// 运行所有未应用的迁移（按文件名排序）
    pub async fn migrate(&self) -> Result<()> {
        // 迁移文件用 include_str! 在编译期内嵌，无需运行时文件 IO
        const MIGRATIONS: &[(&str, &str)] = &[
            ("001_init",    include_str!("../../migrations/001_init.sql")),
            ("002_plugins", include_str!("../../migrations/002_plugins.sql")),
            ("003_payload", include_str!("../../migrations/003_payload.sql")),
            ("004_audit",   include_str!("../../migrations/004_audit.sql")),
        ];

        self.call(|conn| {
            // 确保追踪表存在
            conn.execute_batch(
                "CREATE TABLE IF NOT EXISTS _migrations (
                    name TEXT PRIMARY KEY, applied_at INTEGER NOT NULL
                );"
            )?;
            for (name, sql) in MIGRATIONS {
                let exists: bool = conn.query_row(
                    "SELECT COUNT(*) FROM _migrations WHERE name = ?1",
                    [name], |r| r.get::<_, i64>(0),
                ).map(|n| n > 0).unwrap_or(false);

                if !exists {
                    conn.execute_batch(sql)?;
                    conn.execute(
                        "INSERT INTO _migrations (name, applied_at) VALUES (?1, strftime('%s','now'))",
                        [name],
                    )?;
                }
            }
            Ok(())
        }).await
    }
}
```

Service 中调用示例：
```rust
pub async fn find_by_id(&self, id: &str) -> Result<Webshell> {
    let id = id.to_string();
    self.db.call(move |conn| {
        conn.query_row("SELECT * FROM webshells WHERE id=?1", [&id], |r| { ... })
    }).await
}
```

测试时使用 `:memory:` 数据库：
```rust
#[cfg(test)]
async fn test_db() -> Database {
    Database::open(Path::new(":memory:")).await.unwrap()
}
```

### per-WebShell 请求队列

每个 WebShell 有独立的 `mpsc` 通道，确保请求串行执行，避免并发乱序。响应通过 **oneshot channel** 回传给调用方：

```rust
// features/webshell/queue.rs

pub struct QueuedRequest {
    pub config:     WebshellHttpConfig,    // 该请求的 HTTP 配置
    pub body:       Vec<u8>,               // 已加密的请求体
    pub respond_to: oneshot::Sender<Result<reqwest::Response>>,
}

/// key = webshell_id，value = (发送端, worker JoinHandle)
pub struct WebshellQueue {
    workers:   DashMap<String, (mpsc::Sender<QueuedRequest>, JoinHandle<()>)>,
    http_pool: Arc<HttpClientPool>,
}

impl WebshellQueue {
    /// 发送请求并等待响应。队列不存在时**懒惰创建** worker。
    pub async fn send(
        &self,
        webshell_id: &str,
        config: WebshellHttpConfig,
        body: Vec<u8>,
    ) -> Result<reqwest::Response> {
        let (resp_tx, resp_rx) = oneshot::channel();

        let tx = if let Some(entry) = self.workers.get(webshell_id) {
            entry.0.clone()
        } else {
            let (tx, rx) = mpsc::channel::<QueuedRequest>(32);
            let pool    = self.http_pool.clone();
            let id      = webshell_id.to_string();
            let handle  = tokio::spawn(queue_worker(rx, pool, id));
            self.workers.insert(webshell_id.to_string(), (tx.clone(), handle));
            tx
        };

        tx.send(QueuedRequest { config, body, respond_to: resp_tx }).await
            .map_err(|_| AppError::Connection("queue dropped".into()))?;
        resp_rx.await
            .map_err(|_| AppError::Connection("response channel closed".into()))?
    }

    /// 窗口关闭或断开连接时调用，终止后台 worker
    pub fn cleanup(&self, webshell_id: &str) {
        if let Some((_, (_, handle))) = self.workers.remove(webshell_id) {
            handle.abort();
        }
    }
}

/// 每个 WebShell 对应一个后台任务，串行执行队列中的 HTTP 请求
async fn queue_worker(
    mut rx: mpsc::Receiver<QueuedRequest>,
    pool:   Arc<HttpClientPool>,
    id:     String,
) {
    while let Some(req) = rx.recv().await {
        // C2 Profile jitter：在发送前随机延迟，模拟人工操作节奏，降低被检测概率
        if let Some((min, max)) = req.config.jitter_ms {
            let delay = rand::thread_rng().gen_range(min..=max);
            tokio::time::sleep(tokio::time::Duration::from_millis(delay)).await;
        }

        // body → RequestBuilder：worker 负责组装请求，HttpClientPool 只做熔断 + 发送
        let client = pool.get_client(&req.config);
        let method = req.config.method.to_uppercase();
        let mut builder = match method.as_str() {
            "GET"  => client.get(&req.config.url),
            _      => client.post(&req.config.url),
        };
        for (k, v) in &req.config.custom_headers {
            builder = builder.header(k, v);
        }
        // Cookie 拼接为 "k=v; k2=v2"
        if !req.config.cookies.is_empty() {
            let cookie_str = req.config.cookies.iter()
                .map(|(k, v)| format!("{k}={v}")).collect::<Vec<_>>().join("; ");
            builder = builder.header("Cookie", cookie_str);
        }
        // body 以 form 参数形式发送：param_name=base64(encrypted_body)
        // 服务端 PHP 通过 $_POST['pass']（或自定义参数名）接收
        let b64_body = base64::Engine::encode(&base64::engine::general_purpose::STANDARD, &req.body);
        builder = builder
            .form(&[(&req.config.request_param, b64_body)])
            .timeout(std::time::Duration::from_secs(req.config.timeout_secs));
        let result = pool.send_raw(&id, builder).await;
        let _      = req.respond_to.send(result);  // 忽略 SendError（调用方已取消等待）
    }
}
```

队列 **懒惰创建**：WebShell 第一次发送请求时才 spawn worker，无需显式 connect/disconnect 生命周期管理。

### WebShell Session 管理

```rust
// features/webshell/session.rs
pub struct WebshellSession {
    pub session_id:    String,      // Init 阶段服务端返回
    pub init_time:     Instant,
    pub response_mark: (String, String), // 响应真实性标记，Init 时协商
    pub timeout_secs:  u64,
}

impl WebshellSession {
    pub fn is_expired(&self) -> bool {
        self.init_time.elapsed().as_secs() > self.timeout_secs
    }
}
```

Session 存储在 `WebshellService` 内部的 `DashMap<String, WebshellSession>`（key = webshell_id），App 启动时为空，Init 成功后写入，WebShell 断开或 App 退出时清除。不持久化到 DB。

Session 过期时，下一次 Exec 请求**透明地**重新执行 Init，对调用方无感知。

**`response_mark` 全局唯一性保证**：同一进程内两个 Session 若碰巧生成相同 mark，会互相污染响应解析。在写入 sessions DashMap 前检查并重试：

```rust
fn generate_unique_response_mark(sessions: &DashMap<String, WebshellSession>) -> (String, String) {
    loop {
        let start = random_ascii(8);
        let end   = random_ascii(8);
        // 检查当前所有 session 是否有冲突
        let conflict = sessions.iter().any(|s| {
            s.response_mark.0 == start || s.response_mark.1 == end
        });
        if !conflict { return (start, end); }
        // 碰撞极罕见（概率 ~1/62^8），最多 3 次重试
    }
}
```

### 超越冰蝎：客户端生成会话密钥

冰蝎的密钥协商缺陷：session_key 由服务端生成后**明文返回**，MITM 可截获。FG-ABYSS 改为由**客户端生成 session_key**，在 Init 时包裹在 master_key 加密的 payload 内传给服务端，服务端仅能通过解密拿到它：

```
Init 请求体（master_key 加密）：
  { action:"init", session_key:"<base64 random 32B>",
    response_mark_start:"<8char>", response_mark_end:"<8char>",
    plugin_code:"<base64 bytecode>" }

服务端响应（session_key 加密）：
  <mark_start><session_key 加密的 { session_id, version }><mark_end>
```

- session_key 从不明文出现在网络上（即使抓包也无法获得）
- 每个 session 使用独立密钥，session 过期后密钥丢弃（`Zeroizing` drop）
- 后续所有 Exec 使用 session_key，不使用 master_key（限制 master_key 暴露面）

`WebshellSession` 扩展 `session_key` 字段：

```rust
pub struct WebshellSession {
    pub session_id:    String,
    pub session_key:   zeroize::Zeroizing<[u8; 32]>,  // Drop 时自动清零
    pub init_time:     Instant,
    pub response_mark: (String, String),
    pub timeout_secs:  u64,
}
```

### 共享模型与输入类型

```rust
// features/webshell/models.rs

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "lowercase", tag = "type")]
pub enum PayloadType {
    Php, Jsp, Asp, Aspx,
    JavaMemShell(JavaMemShellKind),   // JVM 内存马（v1.1.0）
    AspNetMemShell,                    // ASP.NET 内存马（v1.1.0）
}

/// 内存马类型（JVM/ASP.NET，不落盘，重启后失效）
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub enum JavaMemShellKind {
    TomcatFilter,      // 注入 Filter chain，拦截所有请求
    TomcatServlet,     // 注入 Servlet
    SpringController,  // 注入 Spring MVC Controller
    JavaAgent,         // JVMTI attach（需要 tools.jar）
}

impl PayloadType {
    pub fn extension(&self) -> &'static str {
        match self {
            Self::Php => "php", Self::Jsp => "jsp",
            Self::Asp => "asp", Self::Aspx => "aspx",
            Self::JavaMemShell(_) | Self::AspNetMemShell => "mem",
        }
    }
    /// 内存马不落盘，Session 失效后无法自动重新初始化
    pub fn is_memshell(&self) -> bool {
        matches!(self, Self::JavaMemShell(_) | Self::AspNetMemShell)
    }
}

/// 数据库行 → Service 层对象（password 已解密）
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Webshell {
    pub id: String, pub name: String, pub url: String,
    pub password: String,            // 明文（Service 层解密，commands 层直接传前端）
    pub payload_type: PayloadType,
    pub project_id: Option<String>,
    pub status: String,              // inactive / active / needs_redeploy
    pub tags: Vec<String>,
    pub custom_headers: HashMap<String, String>,
    pub cookies: HashMap<String, String>,
    pub proxy_override: Option<String>,
    pub http_method: String,
    pub c2_profile: String,
    pub crypto_chain: CryptoChain,
    pub notes: Option<String>,
    pub last_connected_at: Option<i64>,
    pub created_at: i64,
    pub updated_at: i64,
}

/// 来自前端的创建请求，Tauri 自动将 camelCase JSON → snake_case Rust 字段
#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CreateWebshellInput {
    pub name: String, pub url: String, pub password: String,
    pub payload_type: PayloadType,
    pub project_id: Option<String>,
    #[serde(default)] pub tags: Vec<String>,
    #[serde(default)] pub custom_headers: HashMap<String, String>,
    #[serde(default)] pub cookies: HashMap<String, String>,
    pub proxy_override: Option<String>,
    #[serde(default = "default_post")] pub http_method: String,
    #[serde(default = "default_c2")] pub c2_profile: String,
    #[serde(default)] pub crypto_chain: CryptoChain,
    pub notes: Option<String>,
}
fn default_post() -> String { "post".to_string() }
fn default_c2()   -> String { "default".to_string() }

/// 部分更新——所有字段可选，None 表示不修改
#[derive(Debug, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct UpdateWebshellInput {
    pub name: Option<String>, pub url: Option<String>, pub password: Option<String>,
    pub project_id: Option<String>,
    pub tags: Option<Vec<String>>,
    pub custom_headers: Option<HashMap<String, String>>,
    pub cookies: Option<HashMap<String, String>>,
    pub proxy_override: Option<String>,
    pub http_method: Option<String>,
    pub c2_profile: Option<String>,
    pub crypto_chain: Option<CryptoChain>,  // 变化时触发 needs_redeploy
    pub notes: Option<String>,
    pub status: Option<String>,             // 仅允许用户手动重置 active→inactive
}
```

```rust
/// 探测远程服务器环境，辅助选择 Payload 版本和混淆策略
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct WebshellFingerprint {
    pub os:            Option<String>,   // "Linux", "Windows Server 2019"
    pub runtime:       Option<String>,   // "PHP", "Java", "ASP.NET"
    pub runtime_ver:   Option<String>,   // "8.1.12", "11.0.14"
    pub username:      Option<String>,   // "www-data", "IUSR", "nobody"
    pub cwd:           Option<String>,   // "/var/www/html"
    pub disabled_fns:  Vec<String>,      // PHP disable_functions，影响混淆策略选择
    pub writable:      Option<bool>,     // 当前目录是否可写（判断是否能上传文件）
    pub last_probed_at: Option<i64>,
}

/// 连接测试结果
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ConnectionResult {
    pub success:     bool,
    pub latency_ms:  Option<u64>,
    pub error:       Option<String>,
    pub fingerprint: Option<WebshellFingerprint>,  // 连接成功时顺带探测
}
```

**`status` 状态机转换规则**：

```
inactive ──── test_connection() 成功 ────► active
  ▲                                           │
  │     circuit_breaker 触发（失败 ≥3）       │
  └──────────────────────────────────────────┘
  
active ──── crypto_chain 被修改 ────► needs_redeploy
                                           │
                          用户确认重部署后  │
                    reset_redeploy_status()│
                                           ▼
                                        inactive
```

| 转换 | 触发条件 | 执行者 |
|------|---------|--------|
| inactive → active | `test_connection()` / `exec()` 成功 | WebshellService |
| active → inactive | 连续失败 ≥ `FAILURE_THRESHOLD` 次（熔断） | HttpClientPool 回调 |
| any → needs_redeploy | `update()` 检测到 `crypto_chain` 变化 | WebshellService::update() |
| needs_redeploy → inactive | 用户调用 `reset_deploy_status(id)` 确认已重部署 | commands 层显式调用 |

前端根据 `status` 控制操作按钮：`needs_redeploy` 状态下"连接/执行"按钮禁用，显示"重部署警告"横幅。

TypeScript 对应类型在 `features/webshell/types.ts` 中镜像定义（手动维护，骨架阶段先写，后续可考虑 tauri-specta 自动生成）。

### WebshellService 字段定义

```rust
// features/webshell/service.rs
pub struct WebshellService {
    repo:     Arc<dyn WebshellRepo>,  // 通过 Repo trait 访问 DB，测试时可 mock
    crypto:   Arc<CryptoContext>,
    queue:    Arc<WebshellQueue>,
    sessions: Arc<DashMap<String, WebshellSession>>,
}

impl WebshellService {
    pub fn new(
        repo:   Arc<dyn WebshellRepo>,
        crypto: Arc<CryptoContext>,
        queue:  Arc<WebshellQueue>,
    ) -> Self {
        Self { repo, crypto, queue, sessions: Arc::new(DashMap::new()) }
    }

    /// 创建时加密密码，URL 规范化后写入 Repo
    pub async fn create(&self, input: CreateWebshellInput) -> Result<Webshell> {
        let url          = normalize_url(&input.url)?;
        let enc_password = self.crypto.encrypt(input.password.as_bytes())?;
        let plain_pw     = input.password.clone();
        let mut w        = Webshell::from_input(input, url, enc_password);  // 见 models.rs
        self.repo.insert(&w).await?;
        w.password = plain_pw;  // 返回给调用方的是明文
        Ok(w)
    }

    pub async fn get(&self, id: &str) -> Result<Webshell> {
        let mut w = self.repo.find_by_id(id).await?;
        w.password = self.crypto.decrypt_str(&w.password)?;  // 解密覆盖加密存储的值
        Ok(w)
    }

    /// 连接测试：Phase 1 Init → 成功则写回 fingerprint + status=active
    pub async fn test_connection(&self, id: &str) -> Result<ConnectionResult> {
        let shell = self.get(id).await?;
        let start = std::time::Instant::now();
        match self.phase1_init(&shell).await {
            Ok(session) => {
                let latency     = start.elapsed().as_millis() as u64;
                let fingerprint = self.probe_fingerprint(&shell, &session).await.ok().flatten();
                self.sessions.insert(id.to_string(), session);
                self.repo.update_status(id, "active", fingerprint.as_ref()).await?;
                Ok(ConnectionResult { success: true, latency_ms: Some(latency),
                    error: None, fingerprint })
            }
            Err(e) => {
                self.repo.update_status(id, "inactive", None).await?;
                Ok(ConnectionResult { success: false, latency_ms: None,
                    error: Some(e.to_string()), fingerprint: None })
            }
        }
    }

    // 唯一性检查（同 URL）：insert 前调用 repo.find_by_url()，重复则返回 InvalidInput

    /// Phase 1 Init（私有）：生成 32 字节 session_key，用 master_key 加密后随 plugin_code 发出
    async fn phase1_init(&self, shell: &Webshell) -> Result<WebshellSession> {
        todo!("协议初始化实现")
    }

    /// Phase 2 Exec（私有）：用 session_key 加密调用，从响应中提取 response_mark 包裹的数据
    async fn phase2_exec(
        &self,
        shell:   &Webshell,
        session: &WebshellSession,
        method:  &str,
        args:    serde_json::Value,
    ) -> Result<serde_json::Value> {
        todo!("协议执行实现")
    }

    /// 指纹探测（私有）：Init 成功后立即发一次 fingerprint 查询（OS/版本/disabled_fns 等）
    async fn probe_fingerprint(
        &self,
        shell:   &Webshell,
        session: &WebshellSession,
    ) -> Result<Option<WebshellFingerprint>> {
        todo!("指纹探测实现")
    }
}

/// URL 规范化：scheme 小写、hostname 小写、保留端口号、去尾随 /、trim 空白
fn normalize_url(raw: &str) -> Result<String> {
    let trimmed = raw.trim();
    let url = url::Url::parse(trimmed)
        .map_err(|e| AppError::InvalidInput(format!("invalid URL: {e}")))?;
    let scheme = url.scheme();
    let host   = url.host_str().unwrap_or("");
    // 只有非标准端口才拼接（标准端口 http=80/https=443 omit，与浏览器行为一致）
    let port_str = match url.port() {
        Some(p) => format!(":{p}"),
        None    => String::new(),
    };
    let path = url.path().trim_end_matches('/');
    let path = if path.is_empty() { "/" } else { path };
    Ok(format!("{scheme}://{host}{port_str}{path}"))
}

    /// Exec 入口：needs_redeploy 状态下禁止执行，避免解密乱码或蜜罐响应
    /// 透明 Init-or-reuse 逻辑（伪代码）：
    ///
    /// ```
    /// 1. shell = get(id)；status == needs_redeploy → Err(NeedsRedeploy)
    /// 2. session = sessions.get(id)
    /// 3. 若 session 为 None 或 session.is_expired() → phase1_init(shell) → 写入 sessions
    /// 4. result = phase2_exec(shell, session, method, args)
    /// 5. 若 result == InvalidResponse:
    ///      若 shell.payload_type.is_memshell() → Err(MemShellExpired)  // 不重试
    ///      否则 → sessions.remove(id)；重试一次 init+exec（session 可能已过期）
    ///      二次重试仍失败 → 透传 InvalidResponse
    /// 6. Ok(decrypted_data)
    /// ```
    pub async fn exec(&self, id: &str, method: &str, args: serde_json::Value) -> Result<serde_json::Value> {
        let shell = self.get(id).await?;
        if shell.status == "needs_redeploy" {
            return Err(AppError::NeedsRedeploy(
                format!("WebShell {id} payload has changed, redeploy required before exec")
            ));
        }
        todo!("实现透明 Init-or-reuse + Exec")
    }
}
```

---

## 6. 数据库 Schema

```sql
-- 迁移追踪表（所有迁移文件运行前已存在）
-- Database::migrate() 第一步创建此表（若不存在）
CREATE TABLE IF NOT EXISTS _migrations (
    name       TEXT PRIMARY KEY,   -- 文件名，如 "001_init"
    applied_at INTEGER NOT NULL    -- Unix 时间戳
);

-- migrations/001_init.sql
CREATE TABLE projects (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  INTEGER NOT NULL,
    updated_at  INTEGER NOT NULL,
    deleted_at  INTEGER
);

CREATE TABLE webshells (
    id                TEXT PRIMARY KEY,   -- UUID v4，Service 层自动生成
    name              TEXT NOT NULL,
    url               TEXT NOT NULL,
    password          TEXT NOT NULL,      -- AES-256-GCM 加密存储
    payload_type      TEXT NOT NULL,      -- php / jsp / asp / aspx
    project_id        TEXT REFERENCES projects(id),
    status            TEXT NOT NULL DEFAULT 'inactive'
                          CHECK (status IN ('inactive', 'active', 'needs_redeploy')),
    tags              TEXT NOT NULL DEFAULT '[]',     -- JSON 数组；当前用 LIKE 过滤，高频过滤场景升级为 webshell_tags 关联表（P1）
    custom_headers    TEXT NOT NULL DEFAULT '{}',     -- JSON 对象
    cookies           TEXT NOT NULL DEFAULT '{}',     -- JSON 对象
    proxy_override    TEXT,               -- NULL = 使用全局代理；格式："http://host:port" 或 "socks5://host:port"
    http_method       TEXT NOT NULL DEFAULT 'post',
    c2_profile        TEXT NOT NULL DEFAULT 'default',
    crypto_chain      TEXT NOT NULL DEFAULT '[]',     -- JSON CryptoChain
    fingerprint       TEXT,               -- JSON WebshellFingerprint
    notes             TEXT,               -- 用户备注
    last_connected_at INTEGER,            -- 最近一次成功连接时间
    created_at        INTEGER NOT NULL,
    updated_at        INTEGER NOT NULL,
    deleted_at        INTEGER
);
CREATE INDEX idx_webshells_active ON webshells(deleted_at, project_id);
CREATE INDEX idx_webshells_status ON webshells(status, deleted_at);

-- migrations/002_plugins.sql
CREATE TABLE plugins (
    id          TEXT PRIMARY KEY,   -- UUID v4
    name        TEXT NOT NULL,
    version     TEXT NOT NULL,
    enabled     INTEGER NOT NULL DEFAULT 1,  -- 0/1 boolean
    config      TEXT NOT NULL DEFAULT '{}',  -- JSON 配置
    source      TEXT NOT NULL DEFAULT 'builtin',  -- builtin / user
    created_at  INTEGER NOT NULL,
    updated_at  INTEGER NOT NULL
);

-- migrations/003_payload.sql
CREATE TABLE payloads (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    payload_type TEXT NOT NULL,
    config       TEXT NOT NULL,      -- JSON 序列化的配置
    created_at   INTEGER NOT NULL,
    deleted_at   INTEGER
);

CREATE TABLE payload_history (
    id               TEXT PRIMARY KEY,
    payload_id       TEXT REFERENCES payloads(id),
    webshell_id      TEXT,               -- 从 WebShell 页面触发生成时填入，便于 needs_redeploy 场景快速找回
    code             TEXT NOT NULL,
    template_version TEXT NOT NULL DEFAULT 'v1',  -- 记录生成时的模板版本
    created_at       INTEGER NOT NULL
);
CREATE INDEX idx_payload_history_webshell ON payload_history(webshell_id, created_at);

-- migrations/004_audit.sql
CREATE TABLE audit_events (
    id           TEXT PRIMARY KEY,
    webshell_id  TEXT,
    action       TEXT NOT NULL,   -- FileRead/FileWrite/CommandExec/DBQuery
    detail       TEXT NOT NULL,   -- 操作摘要（非完整数据）
    created_at   INTEGER NOT NULL
);
CREATE INDEX idx_audit_webshell ON audit_events(webshell_id, created_at);
```

所有业务查询默认附加 `WHERE deleted_at IS NULL`。

**软删除级联策略**：Project 软删除必须与其下所有 WebShell 的软删除在**同一个 SQLite 事务**内完成，防止进程崩溃造成孤儿数据。`ProjectRepo` 提供原子方法：

```rust
// features/project/repo.rs
async fn soft_delete_with_webshells(&self, project_id: &str, ts: i64) -> Result<()> {
    self.db.call(move |conn| {
        conn.execute_batch("BEGIN;")?;
        conn.execute(
            "UPDATE webshells SET deleted_at=?1 WHERE project_id=?2 AND deleted_at IS NULL",
            rusqlite::params![ts, project_id],
        )?;
        conn.execute(
            "UPDATE projects SET deleted_at=?1 WHERE id=?2 AND deleted_at IS NULL",
            rusqlite::params![ts, project_id],
        )?;
        conn.execute_batch("COMMIT;")?;
        Ok(())
    }).await
}
```

恢复 Project 时对应 `restore_with_webshells(project_id, ts)` 同样在事务内执行（`WHERE deleted_at = ts`，只恢复本次删除的 WebShell）。

**`updated_at` 维护**：SQLite 无自动更新时间戳，所有 UPDATE 语句必须显式 `SET updated_at = strftime('%s','now')`。

**JSON 字段反序列化**：`tags`、`custom_headers`、`cookies`、`crypto_chain`、`fingerprint` 等字段存储为 TEXT（JSON 字符串），`DbWebshellRepo` 读取时必须显式 `serde_json::from_str(&row.tags)?` 反序列化，写入时 `serde_json::to_string(&val)?`。不能直接 `rusqlite::Row::get` 到 Rust 结构体类型，需要先 get 为 `String` 再 deserialize。

---

## 7. 配置文件结构

```toml
# config.toml

[meta]
config_version = 1   # schema 版本，升级时触发迁移

[appearance]
theme = "dark"
language = "zh-CN"
accent_color = "#2080f0"

[window]
x = 100
y = 100
width = 1400
height = 900
maximized = false

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
salt = ""                    # 首次启动随机生成，Argon2id 每次派生主密钥
master_password_enabled = false   # 可选：启用后需要主密码解锁
master_password_hash = ""         # Argon2id hash，空=未设置（明文不落盘）
idle_lock_minutes = 30            # 空闲锁定时间（master_password_enabled=true 时生效）

[logging]
level = "info"
max_file_size_mb = 10
max_files = 5
```

**Config Rust 类型**（与 TOML 结构一一对应）：

```rust
// infra/config.rs
#[derive(Debug, Serialize, Deserialize, Default)]
pub struct Config {
    pub meta:       MetaConfig,
    pub appearance: AppearanceConfig,
    pub window:     WindowConfig,
    pub connection: ConnectionConfig,
    pub security:   SecurityConfig,
    pub logging:    LoggingConfig,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct MetaConfig {
    pub config_version: u32,
}
impl Default for MetaConfig { fn default() -> Self { Self { config_version: 1 } } }

#[derive(Debug, Serialize, Deserialize)]
pub struct AppearanceConfig {
    pub theme: String, pub language: String, pub accent_color: String,
}
impl Default for AppearanceConfig {
    fn default() -> Self { Self { theme: "dark".into(), language: "zh-CN".into(), accent_color: "#2080f0".into() } }
}

#[derive(Debug, Serialize, Deserialize, Default)]
pub struct WindowConfig {
    pub x: i32, pub y: i32, pub width: u32, pub height: u32, pub maximized: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ConnectionConfig {
    pub timeout_secs: u64, pub retry_count: u32,
    pub proxy_enabled: bool, pub proxy_type: String,
    pub proxy_host: String, pub proxy_port: u16,
}
impl Default for ConnectionConfig {
    fn default() -> Self { Self { timeout_secs: 30, retry_count: 3, proxy_enabled: false,
        proxy_type: "http".into(), proxy_host: "".into(), proxy_port: 7890 } }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SecurityConfig {
    pub encryption: String, pub key_rotation_days: u32,
    pub salt: String,
    pub master_password_enabled: bool,
    pub master_password_hash: String,  // Argon2id hash；空 = 未启用
    pub idle_lock_minutes: u32,
}
impl Default for SecurityConfig {
    fn default() -> Self { Self { encryption: "aes-256-gcm".into(), key_rotation_days: 30,
        salt: String::new(), master_password_enabled: false,
        master_password_hash: String::new(), idle_lock_minutes: 30 } }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LoggingConfig {
    pub level: String, pub max_file_size_mb: u64, pub max_files: usize,
}
impl Default for LoggingConfig {
    fn default() -> Self { Self { level: "info".into(), max_file_size_mb: 10, max_files: 5 } }
}

pub fn load_or_default(path: &Path) -> Result<Config> {
    if path.exists() {
        let text = std::fs::read_to_string(path)?;
        let mut cfg: Config = toml::from_str(&text)
            .map_err(|e| AppError::InvalidInput(e.to_string()))?;
        migrate_config(&mut cfg);    // 补全缺失字段
        Ok(cfg)
    } else {
        Ok(Config::default())
    }
}

pub fn save(path: &Path, cfg: &Config) -> Result<()> {
    let text = toml::to_string_pretty(cfg)
        .map_err(|e| AppError::InvalidInput(e.to_string()))?;
    std::fs::write(path, text)?;
    Ok(())
}
```

`toml` crate 需加入 Cargo.toml：`toml = "0.8"`。

**SettingsService — 运行时配置共享**：

配置在运行期可被用户修改（修改超时、切换主题等），其他 Service 需要读取最新值（如 `timeout_secs`）。采用 `Arc<RwLock<Config>>` 共享，避免 Service 持有过期的构造时快照：

```rust
// features/settings/service.rs
pub struct SettingsService {
    config_path: PathBuf,
    config:      Arc<RwLock<Config>>,   // 与 AppState 共享同一实例
}

impl SettingsService {
    pub fn new(config_path: PathBuf, config: Config) -> (Self, Arc<RwLock<Config>>) {
        let shared = Arc::new(RwLock::new(config));
        (Self { config_path, config: shared.clone() }, shared)
    }

    pub async fn get(&self) -> Config {
        self.config.read().await.clone()
    }

    pub async fn update_appearance(&self, theme: String, language: String, accent: String) -> Result<()> {
        let mut cfg = self.config.write().await;
        cfg.appearance = AppearanceConfig { theme, language, accent_color: accent };
        config::save(&self.config_path, &cfg)
    }

    pub async fn update_connection(&self, input: ConnectionConfig) -> Result<()> {
        let mut cfg = self.config.write().await;
        cfg.connection = input;
        config::save(&self.config_path, &cfg)
    }

    /// 设置/更换主密码。首次设置 old_password 为 None；已设置后必须提供旧密码验证，
    /// 防止在用户短暂离开期间被他人悄悄替换密码。
    pub async fn set_master_password(
        &self,
        crypto: &CryptoContext,
        old_password: Option<&str>,
        new_password: &str,
    ) -> Result<()> {
        let mut cfg = self.config.write().await;
        // 已有密码时验证旧密码
        if cfg.security.master_password_enabled {
            match old_password {
                Some(old) if crypto.verify_password(old, &cfg.security.master_password_hash) => {}
                Some(_)  => return Err(AppError::InvalidInput("old password incorrect".into())),
                None     => return Err(AppError::InvalidInput("old password required".into())),
            }
        }
        cfg.security.master_password_hash    = CryptoContext::hash_password(new_password)?;
        cfg.security.master_password_enabled = true;
        config::save(&self.config_path, &cfg)
    }
}
```

`AppState::new()` 中：
```rust
let (settings_service, shared_config) = SettingsService::new(paths.config.clone(), config);
// shared_config 可在构建其他需要运行时配置的组件时传入
// 例如 HttpClientPool 可持有 Arc<RwLock<Config>> 以动态读取 timeout
```

**config_version 迁移策略**：启动时读取版本号，若低于当前版本则逐步应用迁移函数（为缺失字段填充默认值），迁移完成后更新版本号写回文件。

---

## 8. 应用数据目录

```rust
// infra/paths.rs
pub struct AppPaths {
    pub data_dir:  PathBuf,
    pub db_path:   PathBuf,
    pub config:    PathBuf,
    pub logs_dir:  PathBuf,
    pub plugins_dir: PathBuf,
    pub exports_dir: PathBuf,
}

impl AppPaths {
    pub fn resolve(app: &AppHandle) -> Result<Self> {
        let base = app.path().app_data_dir()?;
        // Windows: %APPDATA%\fg-abyss\
        // macOS:   ~/Library/Application Support/fg-abyss/
        // Linux:   ~/.local/share/fg-abyss/
        Ok(Self {
            db_path:     base.join("data.db"),
            config:      base.join("config.toml"),
            logs_dir:    base.join("logs"),
            plugins_dir: base.join("plugins"),
            exports_dir: base.join("exports"),
            data_dir:    base,
        })
    }
}
```

所有模块从 `AppPaths` 取路径，不硬编码任何路径。

---

## 9. 首次运行 Bootstrap 流程

`main.rs` 调用 `bootstrap(app)` 完成所有初始化，顺序严格固定：

```rust
// main.rs
async fn bootstrap(app: &AppHandle) -> Result<AppState> {
    let paths = AppPaths::resolve(app)?;
    std::fs::create_dir_all(&paths.logs_dir)?;
    std::fs::create_dir_all(&paths.plugins_dir)?;
    std::fs::create_dir_all(&paths.exports_dir)?;

    let log_guard = logger::init(&paths.logs_dir)?;   // 1. 最先初始化日志
    let mut config = config::load_or_default(&paths.config)?; // 2. 配置（含迁移）

    if config.security.salt.is_empty() {              // 3. 首次运行生成 salt
        config.security.salt = crypto::generate_salt();
        config::save(&paths.config, &config)?;
    }

    let master_key = crypto::derive_key(&config.security.salt)?; // 4. 派生主密钥
    let db = Database::open(&paths.db_path).await?;   // 5. 数据库连接
    db.migrate().await?;                              // 6. 运行迁移

    // 7. VACUUM 在后台异步执行，不阻塞启动
    let db_bg = db.clone();
    tokio::spawn(async move {
        if let Err(e) = db_bg.vacuum_if_needed().await {
            tracing::warn!("VACUUM failed: {}", e);
        }
    });

    // 8. 初始化各 Service
    AppState::new(paths, config, master_key, db, log_guard)
}
```

### 优雅关闭（Graceful Shutdown）

Tauri 应用退出时需要有序释放资源，防止文件传输中断、终端进程残留：

```rust
// state.rs
impl AppState {
    /// 在主窗口 CloseRequested 事件中同步调用
    pub async fn shutdown(&self) {
        tracing::info!("graceful shutdown started");

        // 1. 中止所有 WebShell 队列 worker（停止接受新请求）
        self.webshell_service.queue.shutdown_all();

        // 2. 清理所有活跃控制台会话（终止进程、取消传输）
        //    DashMap 遍历快照，逐一 cleanup
        let ids: Vec<String> = self.console_service
            .active_webshell_ids()
            .collect();
        for id in ids {
            self.console_service.cleanup(&id).await;
        }

        // 3. _log_guard 在 AppState drop 时自动 flush 日志（无需显式调用）
        tracing::info!("graceful shutdown complete");
    }
}

// WebshellQueue 增加 shutdown_all()
impl WebshellQueue {
    pub fn shutdown_all(&self) {
        for entry in self.workers.iter() {
            entry.value().1.abort();
        }
        self.workers.clear();
    }
}
```

注册关闭钩子（在 `setup` 完成后，主窗口创建时）：

```rust
let main_win = app.get_webview_window("main").unwrap();
let handle   = app.handle().clone();
main_win.on_window_event(move |event| {
    if let tauri::WindowEvent::CloseRequested { .. } = event {
        let state: tauri::State<AppState> = handle.state();
        tauri::async_runtime::block_on(state.shutdown());
    }
});
```

### tauri::Builder 注册骨架

**tauri-specta 2.x 集成**：不使用 `tauri::generate_handler![]`，改用 `tauri_specta::Builder` 同时完成类型收集和 handler 注册，两者在一处维护，避免遗漏：

```rust
// main.rs — fn main()
fn main() {
    // tauri-specta Builder：收集 commands + types，debug 构建时写出 bindings.ts
    // tauri-specta 2.x API：collect_commands! / collect_types! 宏名称（非 1.x 的旧名）
    let specta_builder = {
        use tauri_specta::{Builder, collect_commands, collect_types};
        Builder::<tauri::Wry>::new()
            .commands(collect_commands![
                // webshell
                commands::webshell::list_webshells,
                commands::webshell::get_webshell,
                commands::webshell::create_webshell,
                commands::webshell::update_webshell,
                commands::webshell::delete_webshell,
                commands::webshell::test_connection,
                commands::webshell::reset_redeploy_status,  // 用户确认已重部署后调用，清除 needs_redeploy
                commands::webshell::inject_memshell,    // 内存马注入（v1.1.0）
                // project
                commands::project::list_projects,
                commands::project::create_project,
                commands::project::update_project,
                commands::project::delete_project,
                // payload
                commands::payload::list_payloads,
                commands::payload::create_payload,
                commands::payload::generate_payload,
                commands::payload::list_payload_history,
                // console
                commands::console::open_console,
                commands::console::exec_command,
                commands::console::list_files,
                commands::console::download_file,
                commands::console::upload_file,
                commands::console::connect_database,
                commands::console::execute_query,
                // plugin
                commands::plugin::list_plugins,
                commands::plugin::enable_plugin,
                commands::plugin::disable_plugin,
                // settings
                commands::settings::get_settings,
                commands::settings::update_settings,
                commands::settings::unlock,
                commands::settings::set_master_password,
                // system
                commands::system::get_app_info,
                commands::system::get_audit_log,
                // batch
                commands::batch::test_connections,
            ])
            .types(collect_types![
                Webshell, CreateWebshellInput, UpdateWebshellInput,
                ConnectionResult, WebshellFingerprint, PayloadType,
                PayloadConfig, CryptoChain, BatchTestResult, AppErrorDto,
            ].unwrap())
    };

    #[cfg(debug_assertions)]
    specta_builder
        .export(
            specta_typescript::Typescript::default()
                .formatter(specta_typescript::formatter::prettier),
            "../src/bindings.ts",
        )
        .expect("failed to export TypeScript bindings");

    tauri::Builder::default()
        .setup(|app| {
            // setup 在事件循环启动前同步调用，block_on 是安全的
            // 保证 State 在第一个 command 可能被调用之前就已注册
            let state = tauri::async_runtime::block_on(bootstrap(app.handle()))
                .map_err(|e| {
                    eprintln!("Fatal: bootstrap failed: {e}");
                    e
                })?;
            app.manage(state);

            // 后端全局窗口事件：控制台窗口关闭时触发 cleanup
            app.on_window_event(move |window, event| {
                if let tauri::WindowEvent::Destroyed = event {
                    if let Some(id) = window.label().strip_prefix("console-") {
                        let webshell_id = id.to_string();
                        let handle = window.app_handle().clone();
                        tauri::async_runtime::spawn(async move {
                            let state: tauri::State<AppState> = handle.state();
                            state.console_service.cleanup(&webshell_id).await;
                        });
                    }
                }
            });

            Ok(())
        })
        .invoke_handler(specta_builder.invoke_handler())  // 替代 tauri::generate_handler![]
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

### AppState::new()

```rust
// state.rs
impl AppState {
    pub fn new(
        paths:      AppPaths,
        config:     Config,
        master_key: [u8; 32],
        db:         Database,
        log_guard:  WorkerGuard,
    ) -> Result<Self> {
        let crypto       = CryptoContext::new(master_key);  // Arc<CryptoContext>，各 Service 共享
        let http_pool    = Arc::new(HttpClientPool::new(&config.connection));
        let audit_log    = AuditLog::new(db.clone());
        let queue        = Arc::new(WebshellQueue::new(http_pool.clone()));
        let is_locked    = config.security.master_password_enabled;  // 先读，config 后续被移动

        // Repo 层：DbXxxRepo 包裹 Database，实现对应的 XxxRepo trait
        let webshell_repo = Arc::new(DbWebshellRepo::new(db.clone()));
        let project_repo  = Arc::new(DbProjectRepo::new(db.clone()));
        let payload_repo  = Arc::new(DbPayloadRepo::new(db.clone()));
        let plugin_repo   = Arc::new(DbPluginRepo::new(db.clone()));

        let webshell_service = Arc::new(WebshellService::new(webshell_repo, crypto.clone(), queue.clone()));
        // ConsoleService 持有 webshell_service（查 fingerprint）+ audit_log（记录操作）
        // audit_log.clone() 是 Arc<Database> 内部克隆，成本极低
        let console_service = ConsoleService::new(
            webshell_service.clone(), http_pool.clone(), audit_log.clone()
        );
        // paths.plugins_dir 在 paths 被 move 进 Self 前提取
        let plugins_dir = paths.plugins_dir.clone();
        // SettingsService::new 返回 (service, shared_config Arc)
        // shared_config 将来可传给 HttpClientPool 以便运行时动态读取 timeout
        let (settings_service, _shared_config) = SettingsService::new(paths.config.clone(), config);
        Ok(Self {
            paths,  // move 后无法再访问 paths 字段，故前面提前克隆
            batch_service:    BatchService::new(webshell_service.clone()),
            webshell_service,
            project_service:  ProjectService::new(project_repo),
            payload_service:  PayloadService::new(payload_repo, crypto.clone()),
            console_service,
            plugin_service:   PluginService::new(plugin_repo, plugins_dir),
            settings_service,
            audit_log,
            is_locked:        AtomicBool::new(is_locked),
            _log_guard:       log_guard,
        })
    }
}
```

---

## 10. 日志基础设施

```rust
// infra/logger.rs
pub fn init(logs_dir: &Path) -> Result<tracing_appender::non_blocking::WorkerGuard> {
    let file_appender = tracing_appender::rolling::daily(logs_dir, "fg-abyss.log");
    let (non_blocking, guard) = tracing_appender::non_blocking(file_appender);

    tracing_subscriber::registry()
        .with(tracing_subscriber::fmt::layer().with_writer(non_blocking))
        .with(tracing_subscriber::fmt::layer().with_writer(std::io::stderr))
        .with(tracing_subscriber::EnvFilter::new("info"))
        .init();

    Ok(guard)  // 由 AppState 持有，应用退出时自动 flush
}
```

---

## 11. 崩溃处理

在 `bootstrap()` 最开头（日志初始化后）注册 panic hook：

```rust
let logs_dir_c = paths.logs_dir.clone();  // move 捕获需要 owned PathBuf
std::panic::set_hook(Box::new(move |info| {
    let msg = info.to_string();
    tracing::error!("PANIC: {}", msg);
    // 同步写入 logs/crash-{timestamp}.txt，确保不依赖异步日志
    let path = logs_dir_c.join(format!("crash-{}.txt", chrono::Utc::now().timestamp()));
    let _ = std::fs::write(path, &msg);
}));
```

`Cargo.toml` release profile 保留行号信息：
```toml
[profile.release]
debug = 1
strip = false
```

---

## 12. WebShell 通信协议（两阶段）

超越哥斯拉的核心协议设计：

### Phase 1 — Init（会话建立）

```
客户端 → POST {url}
  参数名: 由 C2 Profile 定义（默认 "pass"）
  内容:  加密({ action: "init", plugin_code: "<base64 bytecode>" })

服务端 → 响应
  内容:  {response_mark.0}<加密({ session_id: "xxx", version: "php8.2" })>{response_mark.1}
```

- `plugin_code`：JSP/PHP/ASPX 的插件字节码或代码，服务端动态加载。来源是 `build.rs` 预编译的**内置通用执行插件**（`FILE_MANAGER_JAR`、`FILE_MANAGER_PHP` 等静态字节数组），在 `WebshellService` 内部通过 `include_bytes!` 直接引用，**不依赖 `PayloadService`**。`PayloadService` 只负责生成用户落盘的 WebShell 文件，与运行时通信插件是两个独立概念。
- `response_mark`：Init 时客户端随机生成并注入 payload 的起止标记，防蜜罐/WAF 伪造响应
- JSP 用 `ClassLoader.defineClass()` 加载字节码；ASPX 用 `Assembly.Load(bytes)`；PHP 用 `eval()`

### Phase 2 — Exec（方法调用）

```
客户端 → POST {url}
  内容:  加密({ session_id: "xxx", method: "file.list", args: {...} })

服务端 → 响应
  内容:  {response_mark.0}<加密({ ok: true, data: {...} })>{response_mark.1}
```

- 同一 session 下所有请求复用已加载的插件实例
- session 过期（服务端重启等）：客户端检测到错误响应后透明重新 Init，对上层无感知
- 无状态降级：若 Init 失败，自动降级为 per-request 模式（每次请求都携带 plugin_code）

### 响应真实性验证

解析响应时，若 `response_mark` 不存在或位置异常，抛出 `AppError::InvalidResponse`，不解密，不展示数据。

```rust
/// 从响应体中提取 response_mark 包裹的内容。
/// 用 find() 而非 starts_with()，允许前缀有若干字节噪声（页面头/BOM/HTTP 响应残留）。
/// 搜索范围限制在前 1KB，防止超大响应体引起性能问题。
pub fn extract_response(body: &[u8], mark: &(String, String)) -> Result<Vec<u8>> {
    let search_limit  = body.len().min(1024);
    let start_bytes   = mark.0.as_bytes();
    let end_bytes     = mark.1.as_bytes();

    let start_pos = body[..search_limit]
        .windows(start_bytes.len())
        .position(|w| w == start_bytes)
        .ok_or_else(|| AppError::InvalidResponse("missing response_mark start".into()))?;

    let content_start = start_pos + start_bytes.len();
    let end_pos = body[content_start..]
        .windows(end_bytes.len())
        .position(|w| w == end_bytes)
        .ok_or_else(|| AppError::InvalidResponse("missing response_mark end".into()))?;

    Ok(body[content_start..content_start + end_pos].to_vec())
}
```

**内存马 Session 失效处理**：内存马不落盘，服务端重启后内存马消失，此时 Exec 请求会因响应不含 `response_mark` 而抛出 `InvalidResponse`。`WebshellService::exec()` 检测到 `payload_type.is_memshell()` 为 true 时，将 `InvalidResponse` 转换为 `AppError::MemShellExpired`，前端显示"内存马已失效，需重新注入"提示而非通用错误。

---

## 13. C2 Profile 系统

定义通信外观，使流量看起来像合法业务流量。每个 WebShell 可绑定独立 Profile。

```rust
// infra/c2_profile.rs
pub struct C2Profile {
    pub name:             String,
    pub request_param:    String,             // POST 参数名，默认 "pass"
    pub extra_headers:    HashMap<String, String>,
    pub user_agent:       String,
    pub request_wrapper:  Option<String>,     // 模板：将加密数据包裹进 JSON/XML 等
    pub response_prefix:  String,             // 解析响应时跳过的前缀格式
    pub jitter_ms:        (u64, u64),         // 请求间隔随机抖动 [min, max]
    pub padding_range:    Option<(usize, usize)>, // 加密前追加随机填充字节 [min, max]，消除流量大小特征
}
```

内置 Profile 预设：

| Profile | 外观 |
|---------|------|
| `default` | 原始模式，无伪装 |
| `cdn-callback` | 伪装成 CDN 回源请求，Content-Type: application/octet-stream |
| `api-json` | 请求/响应包裹在 `{"code":0,"data":"..."}` 格式 |
| `form-submit` | multipart/form-data，参数名改为 `_token` |

**流量大小随机化**：`padding_range` 在加密前向 payload 末尾追加 `[min, max]` 字节的随机填充，服务端 payload 通过 JSON 结构解析自然忽略多余字节（填充在加密内，不影响解密后的 JSON key 解析）。每次请求大小随机化，IDS 无法靠固定包长指纹识别。

**`request_wrapper` 服务端对称解包**：客户端将加密数据包裹进指定格式（如 JSON）后发送，服务端 Payload 必须先解包再解密。Payload Generator 根据绑定的 C2 Profile 生成对应解包代码：

| Profile | 客户端包裹格式 | 服务端解包代码（PHP 示例） |
|---------|-------------|------------------------|
| `default` | 裸加密字节 | `$data = $_POST['pass'];` |
| `api-json` | `{"code":0,"data":"<b64>"}` | `$data = json_decode(file_get_contents('php://input'),true)['data'];` |
| `form-submit` | multipart `_token=<b64>` | `$data = $_POST['_token'];` |

Generator 在生成 PHP/JSP/ASPX payload 时，根据 `c2_profile` 字段选择对应的 `UNWRAP_TEMPLATE` 插入到 payload 头部。自定义 Profile 的 `request_wrapper` 字段作为模板字符串，`{DATA}` 占位符替换为实际解包表达式。

自定义 Profile 存储在 `config.toml` 的 `[c2_profiles]` 节（数组）。

---

## 14. 自定义编码链（CryptoChain）

每个 WebShell 可配置独立的编码链，超越哥斯拉的固定 AES/XOR：

```rust
// infra/crypto.rs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum CodecStep {
    Aes256Gcm,
    XorKey(String),
    Base64,
    UrlEncode,
    GzipCompress,
    HexEncode,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CryptoChain {
    pub steps: Vec<CodecStep>,   // 发送：顺序执行；接收：逆序执行
}

impl Default for CryptoChain {
    fn default() -> Self { Self { steps: vec![CodecStep::Aes256Gcm, CodecStep::Base64] } }
}

impl CryptoChain {
    /// key 为 32 字节——Init 阶段传 master_key，Exec 阶段传 session_key（两者 shape 相同）
    pub fn encode(&self, data: &[u8], key: &[u8; 32]) -> Result<Vec<u8>> { ... }
    pub fn decode(&self, data: &[u8], key: &[u8; 32]) -> Result<Vec<u8>> { ... }
}
```

Payload 生成时，Generator 同步输出与 CryptoChain 对应的服务端解码代码（PHP/JSP/ASP/ASPX 各一份）。

**⚠️ 修改约束**：`crypto_chain` 与已部署的 WebShell payload 内嵌解码逻辑**强绑定**。`WebshellService::update()` 检测到 `crypto_chain` 变化时，自动将 `status` 改为 `'needs_redeploy'`，前端展示"需要重新部署 Payload"警告，用户确认后才能继续操作。

**update() 变更检测逻辑**：

```rust
pub async fn update(&self, id: &str, input: UpdateWebshellInput) -> Result<Webshell> {
    let current = self.get(id).await?;

    // crypto_chain 变化 → 强制 needs_redeploy，覆盖调用方传入的 status
    let new_status = match &input.crypto_chain {
        Some(chain) if *chain != current.crypto_chain => "needs_redeploy",
        _ => input.status.as_deref().unwrap_or(&current.status),
    };

    // 处于 needs_redeploy 状态时禁止 Exec 操作（commands 层检查）
    let enc_password = match &input.password {
        Some(p) => Some(self.crypto.encrypt(p.as_bytes())?),
        None    => None,
    };

    self.db.call(move |conn| {
        conn.execute(
            "UPDATE webshells SET name=COALESCE(?1,name), url=COALESCE(?2,url),
             password=COALESCE(?3,password), status=?4,
             crypto_chain=COALESCE(?5,crypto_chain),
             updated_at=strftime('%s','now')
             WHERE id=?6 AND deleted_at IS NULL",
            rusqlite::params![input.name, input.url, enc_password, new_status,
                              input.crypto_chain.map(|c| serde_json::to_string(&c).unwrap()),
                              id],
        )
    }).await?;

    self.get(id).await
}
```

---

## 15. 加密基础设施

```rust
// infra/crypto.rs

// AES-256-GCM，随机 12 字节 IV，存储格式：base64(IV[12] || TAG[16] || CIPHERTEXT)
pub fn encrypt(plain: &[u8], key: &[u8; 32]) -> Result<String> { ... }
pub fn decrypt(cipher: &str,  key: &[u8; 32]) -> Result<Vec<u8>> { ... }

// Argon2id 从 salt 派生 32 字节主密钥（每次启动重新派生，不落盘）
pub fn derive_key(salt: &str) -> Result<[u8; 32]> { ... }

// 随机生成 base64 编码的 32 字节 salt
pub fn generate_salt() -> String { ... }

/// 加解密上下文——持有主密钥，通过 Arc 在多个 Service 间共享，避免 master_key 复制
/// Drop 时自动 zeroize（master_key 字段标注 #[zeroize(drop)]）
pub struct CryptoContext {
    #[zeroize(drop)]
    master_key: [u8; 32],
}

impl CryptoContext {
    pub fn new(key: [u8; 32]) -> Arc<Self> { Arc::new(Self { master_key: key }) }

    pub fn encrypt(&self, plain: &[u8]) -> Result<String> {
        encrypt(plain, &self.master_key)
    }
    pub fn decrypt(&self, cipher: &str) -> Result<Vec<u8>> {
        decrypt(cipher, &self.master_key)
    }
    pub fn decrypt_str(&self, cipher: &str) -> Result<String> {
        String::from_utf8(self.decrypt(cipher)?)
            .map_err(|e| AppError::Crypto(e.to_string()))
    }
    /// 验证主密码（Argon2id hash 比对）
    pub fn verify_password(&self, input: &str, hash: &str) -> bool {
        argon2::Argon2::default()
            .verify_password(input.as_bytes(), &argon2::PasswordHash::new(hash).unwrap_or_default())
            .is_ok()
    }
    /// 生成主密码 Argon2id hash（设置主密码时调用）
    pub fn hash_password(password: &str) -> Result<String> {
        use argon2::password_hash::{rand_core::OsRng, SaltString, PasswordHasher};
        let salt = SaltString::generate(&mut OsRng);
        argon2::Argon2::default()
            .hash_password(password.as_bytes(), &salt)
            .map(|h| h.to_string())
            .map_err(|e| AppError::Crypto(e.to_string()))
    }
}

// 敏感值包装器，防止 Debug/Display 泄漏
pub struct Sensitive<T>(T);

impl<T> Sensitive<T> {
    pub fn new(val: T) -> Self { Self(val) }
    pub fn inner(&self) -> &T { &self.0 }
    pub fn into_inner(self) -> T { self.0 }
}

impl<T> std::fmt::Debug for Sensitive<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        f.write_str("[REDACTED]")
    }
}
```

密码内存安全：存储密码的 `String` 离开作用域前用 `zeroize` crate 清零。

**`Arc<CryptoContext>` 与 `zeroize` 的局限性**：`#[zeroize(drop)]` 在 `CryptoContext` 的最后一个 `Arc` 引用释放时才触发 drop 并清零 `master_key`。只要有任何一个 Arc 克隆存活，内存不会清零。设计约束：`CryptoContext` 的 Arc 克隆**只允许存储在 `AppState` 持有的各 Service 字段中**，禁止将 Arc 克隆到 thread-local、静态变量、或跨越 `AppState` 生命周期的异步任务中，否则 `master_key` 将在 App 退出后仍驻留内存。`shutdown()` 确保所有 Service 不再使用 crypto 后，`AppState` drop → 所有 Arc 引用计数归零 → `master_key` 被 zeroize 清零。

**PHP 版本兼容加密矩阵**：

| PHP 版本 | 方案 |
|---------|------|
| 7.1+ | AES-256-GCM（`openssl_encrypt` with AEAD） |
| 5.6–7.0 | AES-256-CBC + HMAC-SHA256 |
| 5.3–5.5 | XOR + Base64 |

Payload 生成时根据目标 PHP 版本选择对应加密方案，并在 payload 内嵌入版本检测代码。

---

## 16. 混淆器设计

### 统一抽象

```rust
// features/payload/obfuscator/mod.rs
pub struct ObfuscateConfig {
    pub level: u8,           // 0=无混淆 1=变量重命名 2=控流平坦 3=字符串加密 4=垃圾代码 5=全组合
    pub seed: Option<u64>,   // None=随机（每次输出不同）Some=固定（测试用）
    pub target_version: Option<String>, // PHP版本感知混淆
}

impl Default for ObfuscateConfig {
    fn default() -> Self { Self { level: 1, seed: None, target_version: None } }
}

pub trait Obfuscator: Send + Sync {
    fn obfuscate(&self, code: &str, config: &ObfuscateConfig) -> Result<String>;
    fn max_level(&self) -> u8;
}

pub struct ObfuscateEngine {
    php:  PhpObfuscator,
    jsp:  JspObfuscator,
    asp:  AspObfuscator,
    aspx: AspxObfuscator,
}
```

### PHP 混淆（tree-sitter AST）

- 使用 `tree-sitter` + `tree-sitter-php` 进行 AST 变换
- **多态生成**：引入随机 seed，每次输出结果不同但语义等价，变量名从字典随机采样
- **函数别名替换**：`count()→sizeof()`, `exit()→die()`, `is_int()→is_integer()` 等
- **多策略字符串编码**（随机混合）：
  - `"\x65\x76\x61\x6c"` 十六进制转义
  - `"\145\166\141\154"` 八进制转义
  - `chr(101).chr(118).chr(97).chr(108)` ASCII 拼接
  - 字符串随机分片：`"ev"."al"`
- **PHP 版本感知**：根据目标版本选择可用 transform（PHP 7.2+ 禁用 `create_function()`）
- **注释伪装**：插入看起来像正常业务代码的注释和死代码（`// TODO: optimize`）
- **语义验证（测试环境）**：混淆后重新解析 AST，验证语义结构等价

### JSP 混淆（模板 + Java 字符串编码）

- **ClassLoader 反射化**：`"defineClass"` 用字符数组构造 `new char[]{'\u0064',...}`
- **Java Unicode 转义**：字符串中的敏感词用 `\uXXXX` 转义
- **多模板结构**：随机选择 `<% scriptlet %>`、`<%! declaration %>`、`<script runat="server">` 等不同语法形式
- **内嵌字节码多层编码**：XOR + Base64，不直接 base64 暴露 class 文件特征

### ASP（VBScript）混淆

- **Execute 别名链**：`Execute` → `ExecuteGlobal` → `Eval`，可以动态构造函数名
- **大小写随机混合**：`eXeCuTe`, `ExEcUtE`（VBScript 不区分大小写）
- **Chr() 字符串构造**：所有敏感词（execute、eval、request、response）用 `Chr(101)&Chr(120)...` 拼接
- **注释混淆**：插入 `Rem Legacy code` 等误导性注释

### ASPX（C#）混淆

- **Assembly.Load 反射化**：通过 `typeof(Assembly).GetMethod(new string(...))` 调用，避免直接出现 `"Load"`
- **C# 字符串分片**：`"Sys" + "tem.Re" + "flec" + "tion"`
- **ASPX 指令变体**：随机选择 `<% %>` 或 `<script language="C#" runat="server">` 形式

---

## 17. 插件系统

### 内置插件（随 App 分发）

文件管理、数据库管理、终端管理插件在构建期预编译，打包进 App：

```rust
// build.rs
fn main() {
    // 编译 Java 插件（--release 8，兼容 JVM 8+）
    // 若 javac 不可用则跳过：PHP-only 开发者无需 JDK 也能正常编译
    if std::process::Command::new("javac").arg("-version").status().is_ok() {
        std::process::Command::new("javac")
            .args(["--release", "8", "-d", "out/", "plugins/java/FileManager.java"])
            .status().expect("javac compilation failed");
        println!("cargo:rerun-if-changed=plugins/java/FileManager.java");
    } else {
        println!("cargo:warning=javac not found, Java plugin skipped");
    }

    // 编译 .NET 插件（net20 框架）
    if std::process::Command::new("dotnet").arg("--version").status().is_ok() {
        std::process::Command::new("dotnet")
            .args(["build", "--framework", "net20", "-o", "out/"])
            .current_dir("plugins/dotnet")
            .status().expect("dotnet build failed");
        println!("cargo:rerun-if-changed=plugins/dotnet/");
    } else {
        println!("cargo:warning=dotnet not found, ASPX plugin skipped");
    }
}

// 在代码中以字节数组形式引用
// ⚠️ include_bytes! 是编译期宏，路径必须存在否则编译报错。
// build.rs 若因缺少 javac/dotnet 跳过编译，需要保证 out/ 下有空占位文件：
//   echo "" > out/FileManager.jar.placeholder
// 使用 cfg 条件编译区分有/无 JDK 的构建环境：
#[cfg(feature = "java-plugin")]
static FILE_MANAGER_JAR: &[u8] = include_bytes!("../out/FileManager.jar");
#[cfg(not(feature = "java-plugin"))]
static FILE_MANAGER_JAR: &[u8] = &[];   // 无 javac 时为空，运行时跳过 Java plugin 注入

#[cfg(feature = "dotnet-plugin")]
static FILE_MANAGER_DLL: &[u8] = include_bytes!("../out/FileManager.dll");
#[cfg(not(feature = "dotnet-plugin"))]
static FILE_MANAGER_DLL: &[u8] = &[];

static FILE_MANAGER_PHP: &[u8] = include_bytes!("../plugins/php/FileManager.php"); // PHP 无需编译
```

### PHP 服务端插件基类

```php
abstract class FGPlugin {
    protected $session_id;
    abstract public function execute(string $method, array $args): array;
    protected function response(bool $ok, $data = null): array {
        return ['ok' => $ok, 'data' => $data];
    }
}
```

### 插件文件 Serve 策略

用户安装的插件（HTML + JS 包）存储在 `plugins_dir/{plugin_id}/index.html`。Tauri V2 通过 `protocol-asset` feature 提供 `asset://` 协议访问本地文件：

```typescript
// ConsoleView.vue — 加载插件 iframe
import { convertFileSrc } from '@tauri-apps/api/core'

// asset:// 协议将本地路径转换为可供 webview 访问的 URL
const pluginSrc = convertFileSrc(`${pluginsDir}/${pluginId}/index.html`)
// → "asset://localhost/C:/Users/.../plugins/plugin-id/index.html"（Windows）
```

```vue
<iframe
  :src="pluginSrc"
  sandbox="allow-scripts"
  style="border:none; width:100%; height:100%"
/>
```

CSP 需在 `tauri.conf.json` 中允许 `asset:` 来源：
```json
"csp": "default-src 'self' asset: https://asset.localhost; style-src 'self' 'unsafe-inline'"
```

内置插件（文件/终端/DB）不用 iframe，直接作为 Vue 组件集成，性能更好，无需 postMessage 开销。**iframe 沙箱仅用于用户安装的第三方插件**。

### 用户自定义插件（沙箱）

用户插件在独立 `<iframe sandbox="allow-scripts">` 中运行，**无法访问 `window.__TAURI__`**：

```typescript
// 插件 SDK（通过 provide/inject 注入，不暴露原始 IPC）
export interface PluginSDK {
  shell: {
    exec(cmd: string): Promise<string>
    readFile(path: string): Promise<Uint8Array>
    writeFile(path: string, data: Uint8Array): Promise<void>
  }
  ui: {
    showMessage(msg: string, type: 'info' | 'success' | 'error'): void
  }
  http: {
    get(url: string): Promise<string>
    post(url: string, body: string): Promise<string>
  }
}
```

主窗口通过 `postMessage` 接收插件调用，**先验证权限**后代为执行，将结果回传。

**插件 HTTP 白名单**：插件 SDK 的 `http.get/post` 只能访问已注册 WebShell 的 URL（以及用户显式配置的扩展域名），防止恶意插件将数据外泄到任意地址。

```rust
// features/plugin/service.rs
pub struct PluginPermissions {
    /// 自动注入：当前控制台关联 WebShell 的 URL（基础域名，含路径前缀）
    pub allowed_webshell_urls: Vec<String>,
    /// 用户显式追加的可信域名（插件配置页可设置）
    pub extra_allowed_domains: Vec<String>,
}

impl PluginPermissions {
    pub fn is_allowed(&self, url: &str) -> bool {
        self.allowed_webshell_urls.iter().any(|u| url.starts_with(u))
            || self.extra_allowed_domains.iter().any(|d| url.contains(d))
    }
}
```

```typescript
// plugin_host.ts — 主窗口 postMessage 中间人
window.addEventListener('message', async (e: MessageEvent) => {
  if (e.origin !== 'null') return  // iframe sandbox origin 固定为 "null"
  const { type, payload, requestId } = e.data

  if (type === 'http.get' || type === 'http.post') {
    if (!pluginPermissions.isAllowed(payload.url)) {
      e.source?.postMessage({ requestId, error: 'URL not in allowlist' }, '*')
      return
    }
    // 代为执行，隔离插件与 Tauri IPC
    const resp = await fetch(payload.url, {
      method: type === 'http.post' ? 'POST' : 'GET',
      body: type === 'http.post' ? payload.body : undefined,
    })
    e.source?.postMessage({ requestId, data: await resp.text() }, '*')
  }
  // shell.exec / ui.showMessage 等其他类型类似处理
})
```

---

## 18. HTTP 客户端

每个 WebShell 可有独立的 `proxy_override`，`reqwest::Client` 创建后代理不可更改，因此使用**按代理配置懒加载的 Client 缓存**：

```rust
// infra/http.rs

/// HttpClientPool 的 per-request 配置，从 webshell 字段组装
pub struct WebshellHttpConfig {
    pub url:            String,
    pub method:         String,             // "post" / "get"
    pub request_param:  String,             // POST 参数名，来自 C2 Profile（默认 "pass"）
    pub proxy_override: Option<String>,     // "http://host:port" or "socks5://host:port"
    pub custom_headers: HashMap<String, String>,
    pub cookies:        HashMap<String, String>,
    pub timeout_secs:   u64,
    pub jitter_ms:      Option<(u64, u64)>, // C2 Profile 抖动 [min, max]，None = 无延迟
}

impl WebshellHttpConfig {
    /// 用于区分 reqwest::Client 实例（相同代理共享同一 Client）
    pub fn proxy_hash(&self) -> u64 {
        use std::hash::{Hash, Hasher};
        use std::collections::hash_map::DefaultHasher;
        let mut h = DefaultHasher::new();
        self.proxy_override.hash(&mut h);
        h.finish()
    }
}

pub struct HttpClientPool {
    // key = proxy 配置的 hash，相同代理共享同一 Client（连接池复用）
    clients: DashMap<u64, reqwest::Client>,
    circuit:       DashMap<String, CircuitState>,  // key = webshell_id
    failure_count: DashMap<String, u32>,
}

impl HttpClientPool {
    pub fn get_client(&self, config: &WebshellHttpConfig) -> reqwest::Client {
        let key = config.proxy_hash();
        self.clients
            .entry(key)
            .or_insert_with(|| build_client(config))
            .clone()
    }
}

// 熔断器状态
pub enum CircuitState { Closed, Open(Instant), HalfOpen }
```

impl HttpClientPool {
    /// queue_worker 调用：接收已组装好的 RequestBuilder（含 body/headers/timeout），
    /// 只负责熔断检查 + 发送。RequestBuilder 已绑定 reqwest::Client，直接 .send()。
    /// get_client() 在 queue_worker 里调用，send_raw 不再需要知道 Client 实例。
    pub async fn send_raw(&self, webshell_id: &str, req: reqwest::RequestBuilder) -> Result<reqwest::Response> {
        // 检查熔断状态
        match self.circuit.get(webshell_id).as_deref() {
            Some(CircuitState::Open(opened_at)) => {
                if opened_at.elapsed().as_secs() < CIRCUIT_RESET_SECS {
                    return Err(AppError::CircuitOpen(webshell_id.to_string()));
                }
                // 进入 HalfOpen，允许一次探测
            }
            _ => {}
        }

        match req.send().await {
            Ok(resp) => {
                self.failure_count.remove(webshell_id);
                self.circuit.insert(webshell_id.to_string(), CircuitState::Closed);
                Ok(resp)
            }
            Err(e) => {
                let count = self.failure_count.entry(webshell_id.to_string())
                    .and_modify(|c| *c += 1).or_insert(1);
                if *count >= FAILURE_THRESHOLD {
                    self.circuit.insert(webshell_id.to_string(),
                        CircuitState::Open(Instant::now()));
                }
                Err(AppError::Http(e))
            }
        }
    }
}
```

**Charset 处理**：响应体解码前，用 `encoding_rs` 检测并转换 GBK/GB2312 编码（常见于旧版 PHP 服务器），统一转为 UTF-8 后再解析。

**HTTP/2**：reqwest 默认启用 HTTP/2，无需额外配置。

---

## 19. ConsoleService 状态与生命周期

ConsoleService 管理每个 WebShell 控制台窗口的活跃资源，窗口关闭时统一释放：

```rust
// features/console/service.rs
pub struct ConsoleService {
    webshell_service: Arc<WebshellService>,  // 查询 fingerprint（select_exec_method 路径）
    http_pool:        Arc<HttpClientPool>,
    audit_log:        AuditLog,              // 记录文件/命令/DB 操作，不影响主业务

    // 进行中的文件传输：webshell_id -> transfer_id -> AbortHandle
    file_transfers: DashMap<String, DashMap<String, AbortHandle>>,

    // 终端会话状态（命令历史、CWD）：webshell_id -> TerminalState
    // 注：WebShell 终端是 HTTP 请求执行，不产生本地子进程；TerminalState 存服务端状态快照
    terminal_sessions: DashMap<String, TerminalState>,

    // 远程 DB 会话（连接参数 + 查询状态）：webshell_id -> RemoteDbState
    remote_db_sessions: DashMap<String, RemoteDbState>,
}

/// 终端会话在客户端侧的状态快照（服务端 CWD 随命令变化，每次 exec_command 返回后更新）
pub struct TerminalState {
    pub cwd:     String,              // 最后已知的远程工作目录
    pub history: VecDeque<String>,    // 命令历史（最近 200 条，VecDeque 自动弹出旧记录）
}

pub struct RemoteDbState {
    pub db_type: String,       // mysql / postgresql / sqlite / mssql
    pub conn_str: String,      // 加密存储的连接字符串
    pub last_query: Option<String>,
}

impl ConsoleService {
    /// 控制台窗口关闭时调用（监听 window-destroyed 事件触发）
    pub async fn cleanup(&self, webshell_id: &str) {
        // 1. 中止所有文件传输
        if let Some((_, transfers)) = self.file_transfers.remove(webshell_id) {
            for entry in transfers.iter() {
                entry.value().abort();
            }
        }

        // 2. 清理终端会话状态（纯内存，无需异步操作）
        self.terminal_sessions.remove(webshell_id);

        // 3. 丢弃远程 DB 会话
        self.remote_db_sessions.remove(webshell_id);

        tracing::info!("console cleanup done: {}", webshell_id);
    }

    /// 返回所有有活跃资源的 webshell_id（用于 shutdown() 遍历清理）
    pub fn active_webshell_ids(&self) -> impl Iterator<Item = String> + '_ {
        // 取三个 DashMap key 的并集
        let from_transfers = self.file_transfers.iter().map(|e| e.key().clone());
        let from_terminals = self.terminal_sessions.iter().map(|e| e.key().clone());
        let from_db        = self.remote_db_sessions.iter().map(|e| e.key().clone());
        let mut seen = std::collections::HashSet::new();
        from_transfers.chain(from_terminals).chain(from_db)
            .filter(move |id| seen.insert(id.clone()))
    }
}
```

**文件传输注册**：`download_file` command 启动传输前，将 `tokio::spawn` 返回的 `JoinHandle` 的 `abort_handle()` 注册到 `file_transfers[webshell_id][transfer_id]`；传输完成后主动移除。

---

## 20. 流式响应架构

文件传输和终端命令输出使用流式推送，不一次性加载到内存：

```rust
// ConsoleService 方法签名区分一次性响应和流式响应
pub async fn exec_command(&self, id: &str, cmd: &str) -> Result<String> { ... }

// Service 层：返回字节流，不持有 AppHandle
pub async fn download_file_stream(
    &self,
    id: &str,
    path: &str,
) -> Result<impl Stream<Item = Result<Bytes>>> { ... }

// commands/console.rs：持有 AppHandle，负责 emit 进度事件
#[tauri::command]
pub async fn download_file(
    app: AppHandle,
    state: State<'_, AppState>,
    webshell_id: String,
    path: String,
    transfer_id: String,
) -> Result<()> {
    let mut stream = state.console_service.download_file_stream(&webshell_id, &path).await?;
    while let Some(chunk) = stream.next().await {
        let bytes = chunk?;
        app.emit_to(&format!("console-{}", webshell_id), "file-progress",
            FileProgressEvent { transfer_id: transfer_id.clone(), chunk: bytes })?;
    }
    Ok(())
}
```

注意：流式方法需要 `AppHandle`，因此在 `commands/` 层调用，不在 Service 内部 emit。Service 只返回数据流（`impl Stream<Item = Result<Bytes>>`），command 层负责 emit。

---

## 20. PHP disable_functions 自动绕过

指纹探测后，`ConsoleService` 根据 `WebshellFingerprint.disabled_fns` 自动选择可用的命令执行方式，对用户完全透明——这是哥斯拉/冰蝎都没有做到的。

```rust
// features/console/terminal/php_exec.rs

/// 按优先级排列的 PHP 命令执行备用链
pub const PHP_EXEC_CHAIN: &[PhpExecMethod] = &[
    PhpExecMethod::System,          // system()
    PhpExecMethod::Exec,            // exec()
    PhpExecMethod::Passthru,        // passthru()
    PhpExecMethod::ShellExec,       // shell_exec()
    PhpExecMethod::ProcOpen,        // proc_open()
    PhpExecMethod::Popen,           // popen()
    PhpExecMethod::PcntlExec,       // pcntl_exec()（Linux，返回值需重定向）
    PhpExecMethod::MailInject,      // mail() 第5参数注入（需 sendmail）
    PhpExecMethod::PutenvLdPreload, // putenv()+mail() LD_PRELOAD（Linux root）
    PhpExecMethod::ImageMagick,     // convert 命令注入（需 ImageMagick）
    PhpExecMethod::Imap,            // imap_open() RSHELL（需 IMAP 扩展）
];

#[derive(Debug, Clone)]
pub enum PhpExecMethod {
    System, Exec, Passthru, ShellExec, ProcOpen, Popen, PcntlExec,
    MailInject, PutenvLdPreload, ImageMagick, Imap,
}

impl PhpExecMethod {
    /// 生成对应的 PHP 代码片段（command = $cmd 变量）
    pub fn php_snippet(&self) -> &'static str { ... }

    /// 所依赖的函数名（用于与 disabled_fns 比对）
    pub fn required_fns(&self) -> &'static [&'static str] {
        match self {
            Self::System           => &["system"],
            Self::Exec             => &["exec"],
            Self::Passthru         => &["passthru"],
            Self::ShellExec        => &["shell_exec"],
            Self::ProcOpen         => &["proc_open", "proc_get_status", "stream_get_contents"],
            Self::Popen            => &["popen", "fread", "pclose"],
            Self::PcntlExec        => &["pcntl_exec"],
            Self::MailInject       => &["mail", "putenv"],
            Self::PutenvLdPreload  => &["putenv", "mail"],   // 需要额外写 .so 文件，仅 root
            Self::ImageMagick      => &[],                   // 通过 CLI 注入，无 PHP 函数依赖
            Self::Imap             => &["imap_open"],
        }
    }
}

pub fn select_exec_method(disabled: &[String]) -> Option<&'static PhpExecMethod> {
    let disabled_set: HashSet<&str> = disabled.iter().map(|s| s.as_str()).collect();
    PHP_EXEC_CHAIN.iter().find(|m| {
        m.required_fns().iter().all(|f| !disabled_set.contains(f))
    })
}
```

`ConsoleService::exec_command()` 在执行前通过 `WebshellService::get()` 获取 `Webshell.fingerprint`（存储在 DB，非 Session），调用 `select_exec_method()` 选择策略，将选定的 PHP 代码片段注入终端插件请求。fingerprint 未探测（`None`）时默认使用 `System`。

`ConsoleService` 需持有 `Arc<WebshellService>` 引用（在 `AppState::new()` 中传入），以便在 exec 路径中查询 fingerprint，避免数据冗余存储：

（ConsoleService 完整字段定义见第 19 节，此处仅说明 `webshell_service` 字段的用途。）

---

## 21. 多窗口策略

每个 WebShell 对应一个独立控制台窗口：

```rust
// commands/console.rs
#[tauri_specta::command]   // 所有 commands 用 tauri_specta::command，而非 tauri::command
pub async fn open_console(
    app:   AppHandle,
    state: State<'_, AppState>,
    webshell_id: String,
) -> Result<()> {
    let label = format!("console-{}", webshell_id);

    if let Some(win) = app.get_webview_window(&label) {
        win.set_focus()?;
        return Ok(());
    }

    // 必须先 fetch shell 才能拿到 name 和 url 用于窗口标题
    let shell = state.webshell_service.get(&webshell_id).await?;

    tauri::WebviewWindowBuilder::new(
        &app, &label,
        tauri::WebviewUrl::App(format!("/console?id={}", webshell_id).into()),
    )
    .title(format!("{} — {}", shell.name, shell.url))
    .inner_size(1200.0, 800.0)
    .min_inner_size(900.0, 600.0)
    .decorations(false)
    .build()?;

    Ok(())
}
```

**窗口关闭时清理资源**：控制台窗口关闭后 JS 运行时已停止，无法在前端监听。应在 `main.rs` 的 `setup` 回调里注册**后端全局窗口事件监听**，对 `console-` 前缀窗口触发 cleanup：

```rust
// main.rs — setup 回调内，bootstrap 完成后
app.on_window_event(move |window, event| {
    if let tauri::WindowEvent::Destroyed = event {
        let label = window.label();
        if let Some(id) = label.strip_prefix("console-") {
            let webshell_id = id.to_string();
            let handle = window.app_handle().clone();
            tauri::async_runtime::spawn(async move {
                let state: tauri::State<AppState> = handle.state();
                state.console_service.cleanup(&webshell_id).await;
            });
        }
    }
});
```

不使用前端 `listen('window-destroyed', ...)`——控制台窗口销毁时其内部 JS 已停止，事件无法接收。

**Windows 边缘拖拽调整大小**：`decorations: false` 在 Windows 下禁用系统边框后，边缘拖拽失效。需在前端各窗口边缘添加透明 resize 热区，调用 `win.startResizeDragging(ResizeDirection)` API 处理拖拽事件。

---

## 21. Tauri 配置

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
        "resizable": true,
        "center": true
      }
    ],
    "security": {
      "csp": "default-src 'self'; style-src 'self' 'unsafe-inline'; frame-src 'self'"
    }
  }
}
```

### 权限配置

```json
// capabilities/default.json
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
    "shell:allow-open",
    "clipboard-manager:allow-write-text"
  ]
}

// capabilities/console.json
{
  "identifier": "console",
  "windows": ["console-*"],
  "permissions": ["core:default", "http:default"]
}
```

---

## 22. 审计日志

命名为 `AuditLog`（infra 层技术组件，区别于 features 层 Service）：

```rust
// infra/audit.rs
#[derive(Clone)]  // Database 内部是 Arc，clone 成本极低；便于 ConsoleService 持有副本
pub struct AuditLog {
    db: Database,
}

#[derive(Debug, Serialize, Deserialize)]
pub enum AuditAction {
    FileRead, FileWrite, FileDelete,
    CommandExec,
    DBQuery, DBExport,
    WebshellConnect, WebshellDisconnect,
}

impl AuditLog {
    pub async fn record(&self, webshell_id: &str, action: AuditAction, detail: &str) -> Result<()> {
        let id = Uuid::new_v4().to_string();
        let action_str = serde_json::to_string(&action)?;
        let detail = detail.to_string();
        let ts = chrono::Utc::now().timestamp();
        self.db.call(move |conn| {
            conn.execute(
                "INSERT INTO audit_events (id, webshell_id, action, detail, created_at)
                 VALUES (?1, ?2, ?3, ?4, ?5)",
                rusqlite::params![id, webshell_id, action_str, detail, ts],
            )
        }).await?;
        Ok(())
    }
}
```

审计日志独立于 tracing 日志，存储在 SQLite `audit_events` 表，可在设置页导出为 CSV/JSON。

**审计失败不阻断主业务**：`record()` 写入失败（磁盘满、DB 锁）不应导致命令执行失败。调用方统一使用 fire-and-forget 模式而非 `?` 传播：

```rust
// ConsoleService::exec_command 内
if let Err(e) = self.audit_log.record(id, AuditAction::CommandExec, &cmd).await {
    tracing::error!("audit write failed (non-fatal): {}", e);
    // 不 return Err，继续执行
}
```

---

## 23. 应用级身份验证（可选）

主密码功能**完全可选**，不影响未启用用户的使用体验。

```rust
// AppState 中
pub is_locked: AtomicBool,

// 全局 Command 守卫（仅在主密码启用时生效）
fn check_locked(state: &AppState) -> Result<()> {
    if state.is_locked.load(Ordering::Relaxed) {
        Err(AppError::Locked)
    } else {
        Ok(())
    }
}
```

- 首次启动时可选设置主密码（跳过则永不锁定）
- 主密码通过 Argon2id 验证（不存储明文，验证 hash）
- 空闲 `idle_lock_minutes` 分钟后自动锁定（仅主密码启用时）
- 前端检测到 `AppError.kind === 'Locked'` 时弹出解锁界面

**`check_locked` 应用策略**：只在"有副作用"的 command 入口处调用，只读/系统信息类命令不需要：

| 调用 check_locked | 不需要 |
|---------|--------|
| `create_webshell`, `update_webshell`, `delete_webshell` | `get_app_info` |
| `test_connection`, `open_console`, `exec_command` | `get_settings`（读取） |
| `list_files`, `download_file`, `upload_file` | `get_audit_log` |
| `connect_database`, `execute_query` | `list_webshells`, `list_projects`（列表） |
| `create_project`, `update_project`, `delete_project` | `unlock`（解锁本身） |
| `create_payload`, `generate_payload` | |

```rust
// commands/webshell.rs 示例
#[tauri_specta::command]
pub async fn create_webshell(
    state: State<'_, AppState>,
    input: CreateWebshellInput,
) -> Result<Webshell> {
    check_locked(&state)?;  // 第一行，主密码启用时阻断
    state.webshell_service.create(input).await
}
```

---

## 24. SQLite 维护

**VACUUM 策略**（bootstrap 完成后后台异步执行，不阻塞启动）：

```rust
impl Database {
    pub async fn vacuum_if_needed(&self) -> Result<()> {
        let deleted_count: i64 = self.call(|conn| {
            conn.query_row(
                "SELECT COUNT(*) FROM webshells WHERE deleted_at IS NOT NULL",
                [], |r| r.get(0),
            )
        }).await?;
        if deleted_count > 100 {
            self.call(|conn| conn.execute_batch("VACUUM;")).await?;
            tracing::info!("VACUUM completed, {} soft-deleted rows reclaimed", deleted_count);
        }
        Ok(())
    }
}
```

设置页也提供手动触发"清理数据库"按钮。

---

## 25. Payload 模板版本化

```rust
/// Payload 生成配置——包含生成一份完整 WebShell 所需的全部参数
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PayloadConfig {
    pub payload_type:   PayloadType,
    pub password:       String,           // WebShell 密码（明文，生成时内嵌到代码）
    pub crypto_chain:   CryptoChain,      // 编码链（生成对应的服务端解码逻辑）
    pub c2_profile:     String,           // Profile 名称（生成对应的参数名/格式）
    pub obfuscate:      ObfuscateConfig,  // 混淆配置
    pub target_version: Option<String>,   // 目标运行时版本，如 "php7.4"、"java8"
    pub extra:          HashMap<String, String>,  // 自定义模板变量（扩展点）
}

pub struct PayloadTemplate {
    pub id:           String,
    pub version:      &'static str,  // "v1", "v2"，随 App 版本固定
    pub payload_type: PayloadType,
    pub description:  String,
    pub generate:     fn(&PayloadConfig) -> Result<String>,
}
```

`payload_history` 表记录生成时的 `template_version`，保证历史记录可重现（即使当前模板已升级）。旧版本模板以只读形式保留在代码中。

---

## 26. 前端初始化顺序与主题绑定

### main.ts

```typescript
// i18n/index.ts
import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN.json'
import en   from './locales/en.json'

export const i18n = createI18n({
  legacy: false,           // Composition API 模式（useI18n()）
  locale: 'zh-CN',
  fallbackLocale: 'en',
  messages: { 'zh-CN': zhCN, en },
})

// settings store 切换语言时调用
export function setLocale(locale: 'zh-CN' | 'en') {
  i18n.global.locale.value = locale
}
```

`i18n/locales/zh-CN.json` 顶层结构（按 feature 分组，便于维护）：
```json
{
  "nav":      { "home": "首页", "project": "项目", "payload": "载荷", "plugin": "插件", "settings": "设置" },
  "webshell": { "status": { "active": "已连接", "inactive": "未连接", "needs_redeploy": "需重部署" } },
  "error":    { "locked": "请输入主密码以继续", "network": "网络连接失败", "circuit_open": "目标暂时不可达" },
  "action":   { "connect": "连接", "delete": "删除", "copy": "复制", "confirm": "确认", "cancel": "取消" }
}
```

```typescript
// main.ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import { i18n } from './i18n'
import { globalMessage } from './shared/utils/message'  // 初始化 discrete API
import App from './App.vue'

const app = createApp(App)
app.use(i18n)      // 1. i18n 最先（store action 里 t() 依赖它）
app.use(createPinia())
app.use(router)
app.mount('#app')

// globalMessage 在 import 时已经初始化（Pinia 之后），无需额外操作
void globalMessage  // 消除 unused import 警告

// 初始语言从持久化配置同步：App.vue onMounted 里调用
// settingsStore.init() → invoke('get_settings') → setLocale(cfg.appearance.language)
// i18n 默认 zh-CN，首帧可能短暂显示默认语言，但对用户无感（≈1 frame）
```

### App.vue 主题绑定

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

## 27. 前端错误展示策略

| 场景 | 方式 | 说明 |
|------|------|------|
| 操作成功提示、轻量失败（超时重试）| `n-message` toast | 不打断流程 |
| 不可恢复错误（DB崩溃、加密失败）| `n-dialog` 模态框 | 需要用户确认知晓 |
| 风险操作确认（删除 WebShell）| `n-dialog` 确认框 | 防误操作 |
| 表单字段校验失败 | 行内 error 提示 | 精确定位 |
| `AppError.kind === 'Locked'` | 全屏解锁遮罩 | 主密码模式 |

所有错误展示由 store action 处理，组件不直接处理错误。

---

## 28. 前端 Loading 状态规范

单一 `loading: boolean` 无法区分多个并行操作，统一使用 `LoadingMap`：

**注意**：Naive UI 的 `useMessage()` 是 Composition API hook，只能在组件 `setup()` 内调用，Pinia store action 中无法使用。在 `main.ts` 创建全局 discrete API 实例：

```typescript
// shared/utils/message.ts  — 在 main.ts 中 Pinia 初始化后立即调用
import { createDiscreteApi } from 'naive-ui'
export const { message: globalMessage } = createDiscreteApi(['message'])
```

Store action 一律用 `globalMessage.error(...)` 替代 `useMessage().error(...)`。

```typescript
// shared/types/loading.ts
export type LoadingMap = Record<string, boolean>

// store 使用模板
export const useWebshellStore = defineStore('webshell', {
  state: () => ({
    items: [] as Webshell[],
    loading: {} as LoadingMap,
    // loading.list, loading.create, loading['delete-uuid'] 各自独立
    connectingStatus: {} as Record<string, 'connecting' | 'slow' | 'success' | 'failed'>,
  }),
  actions: {
    async fetchAll() {
      this.loading.list = true
      try {
        this.items = await webshellApi.list()
      } catch (e) {
        globalMessage.error((e as AppError).message)
      } finally {
        this.loading.list = false
      }
    },

    // 渐进式超时 UX：5 秒后仍在连接中则提示"连接缓慢"，消除用户疑惑
    async testConnection(id: string) {
      this.loading[`test-${id}`] = true
      this.connectingStatus[id]  = 'connecting'

      const slowTimer = setTimeout(() => {
        if (this.loading[`test-${id}`]) {
          this.connectingStatus[id] = 'slow'
          globalMessage.info(t('webshell.connecting_slow'))
        }
      }, 5000)

      try {
        const result = await webshellApi.testConnection(id)
        this.connectingStatus[id] = result.success ? 'success' : 'failed'
        const item = this.items.find(w => w.id === id)
        if (item) item.status = result.success ? 'active' : 'inactive'
        return result
      } catch (e) {
        this.connectingStatus[id] = 'failed'
        globalMessage.error((e as AppError).message)
        throw e
      } finally {
        clearTimeout(slowTimer)
        this.loading[`test-${id}`] = false
      }
    },
  },
})
```

---

## 29. WebShell 列表右键菜单

最高频操作入口，组件职责划分：

```
ProjectView
└── WebshellTable.vue        # n-data-table，行右键触发菜单
    └── WebshellContextMenu.vue  # n-dropdown，声明式 options 数组
```

菜单项与 store action 的映射（按使用频率排序）：

| 菜单项 | store action | 条件 |
|-------|-------------|------|
| 连接控制台 | `openConsole(id)` | status=active |
| 测试连接   | `testConnection(id)` | 任意 status |
| 编辑       | 路由跳转 `/webshell/edit?id=` | 任意 |
| 复制 URL   | clipboard | 任意 |
| 复制密码   | clipboard（5秒后自动清除） | 任意 |
| 重新生成 Payload | 路由跳转 `/payload?webshellId=` | needs_redeploy 时高亮 |
| 移动到项目 | `moveToProject(id, projectId)` | 任意 |
| 删除       | `deleteWebshell(id)`（confirm dialog） | 任意 |

```typescript
// features/webshell/components/WebshellContextMenu.vue
// 通过 v-model:show + x,y 坐标定位，Naive UI n-dropdown 实现
// options 数组根据当前行 status 动态计算 disabled 状态
const menuOptions = computed(() => buildMenuOptions(props.webshell, webshellStore))
```

---

## 30. 前端状态管理约定

- 每个 feature 有自己的 Pinia store，store 之间不互相 import
- View 层可组合多个 store
- 组件不直接调用 `api.ts`，所有数据操作通过 store action
- 错误在 store action 内捕获，按第 27 节策略展示
- loading 状态用 `LoadingMap`，组件只读取对应 key

---

## 30. 前端路由守卫与 Locked 状态

`Locked` 错误有两个触发来源，分别在不同层拦截：

**层 1 — invoke 拦截器**（在任意命令返回 Locked 时立即跳转）：

```typescript
// shared/utils/invoke.ts
import { invoke as tauriInvoke } from '@tauri-apps/api/core'
import { useRouter } from 'vue-router'

export async function invoke<T>(cmd: string, args?: Record<string, unknown>): Promise<T> {
  try {
    return await tauriInvoke<T>(cmd, args)
  } catch (e) {
    const err = e as AppError
    if (err.kind === 'Locked') {
      // 不能在模块顶层 useRouter()，通过 import router 实例
      router.push('/unlock')
    }
    throw err
  }
}
```

**层 2 — 路由守卫**（防止直接导航到功能页时绕过锁定）：

```typescript
// router/index.ts
import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/',        component: () => import('@/features/home/HomeView.vue') },
    { path: '/project', component: () => import('@/features/webshell/ProjectView.vue') },
    { path: '/payload', component: () => import('@/features/payload/PayloadView.vue') },
    { path: '/plugin',  component: () => import('@/features/plugin/PluginView.vue') },
    { path: '/settings',component: () => import('@/features/settings/SettingsView.vue') },
    { path: '/console', component: () => import('@/features/console/ConsoleView.vue') },
    { path: '/unlock',  component: () => import('@/features/settings/UnlockView.vue') },
  ],
})

// 路由守卫：主密码启用时，非 /unlock 路由需确认解锁状态
router.beforeEach((to) => {
  if (to.path === '/unlock') return true
  const settingsStore = useSettingsStore()
  if (settingsStore.config.security.masterPasswordEnabled && settingsStore.isLocked) {
    return '/unlock'
  }
  return true
})

export default router
export { router }  // 供 invoke.ts 导入
```

`UnlockView` 调用 `commands::settings::unlock(password)` → 后端验证 Argon2id hash → 成功后 `is_locked.store(false)` → 前端 `router.back()`。

---

## 31. 前端 api.ts 调用模式

每个 feature 的 `api.ts` 是对 `invoke` 的薄封装，仅做类型标注，不含逻辑：

```typescript
// features/webshell/api.ts
import { invoke } from '@/shared/utils/invoke'
import type { Webshell, CreateWebshellInput, UpdateWebshellInput, ConnectionResult } from './types'

// Tauri V2 命令名保持 snake_case（与 Rust fn 名一致）
// 参数对象的字段使用 camelCase（Tauri 自动序列化 Rust snake_case → JS camelCase）
export const webshellApi = {
  list:           ()                            => invoke<Webshell[]>('list_webshells'),
  get:            (id: string)                  => invoke<Webshell>('get_webshell', { id }),
  create:         (input: CreateWebshellInput)  => invoke<Webshell>('create_webshell', { input }),
  update:         (id: string, input: UpdateWebshellInput) => invoke<Webshell>('update_webshell', { id, input }),
  delete:         (id: string)                  => invoke<void>('delete_webshell', { id }),
  testConnection: (id: string)                  => invoke<ConnectionResult>('test_connection', { id }),
}
```

**命名约定**：
- Rust command fn：`snake_case`（`list_webshells`）
- TypeScript `invoke` 调用名：与 Rust 保持一致（`'list_webshells'`）
- 参数字段：Tauri 自动将 Rust `webshell_id` 映射为 JS `webshellId`，api.ts 使用 camelCase

---

## 32. Tauri 事件规范

后端通过 `app.emit_to(window_label, event, payload)` 向特定窗口推送事件。**命名格式：`{domain}:{verb}`**（kebab domain，kebab verb）。

### 事件清单

| 事件名 | 方向 | payload 类型 | 说明 |
|--------|------|-------------|------|
| `file:progress` | backend→console | `{ transferId, bytesTotal, bytesDone }` | 文件传输进度 |
| `file:complete` | backend→console | `{ transferId, savedPath? }` | 传输完成 |
| `file:error`    | backend→console | `{ transferId, error: AppError }` | 传输失败 |
| `terminal:output` | backend→console | `{ sessionId, data: number[] }` | 终端字节流 |
| `terminal:exit`   | backend→console | `{ sessionId, code: number }` | 进程退出 |
| `webshell:status-changed` | backend→main | `{ id, status }` | 连接状态更新 |
| `console:cleanup-done`    | backend→console | `{ webshellId }` | cleanup 完成确认 |

### TypeScript 监听模式

```typescript
// 在 ConsoleView 的 onMounted 中
import { listen } from '@tauri-apps/api/event'
import type { UnlistenFn } from '@tauri-apps/api/event'

const unlisten: UnlistenFn[] = []

onMounted(async () => {
  unlisten.push(await listen<FileProgressPayload>('file:progress', (e) => {
    fileStore.updateProgress(e.payload)
  }))
})

onUnmounted(() => {
  unlisten.forEach(fn => fn())
})
```

**规则**：每个 `listen` 必须在组件卸载时调用返回的 `unlisten()` 函数，防止内存泄漏。

---

## 33. 自定义标题栏

```vue
<!-- layouts/components/CustomTitlebar.vue -->
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

## 34. 窗口状态持久化

主窗口关闭时将大小/位置写入 `config.toml` 的 `[window]` 节，下次启动时恢复。在 `main.rs` 监听 `window-close-requested` 事件：

```rust
main_window.on_window_event(|event| {
    if let WindowEvent::CloseRequested { .. } = event {
        if let Ok(pos)  = window.outer_position() { /* 写入 config */ }
        if let Ok(size) = window.outer_size()     { /* 写入 config */ }
    }
});
```

---

## 35. WebShell 导出 / 导入（预留接口）

在 `WebshellService` 预留接口，实现延后：

```rust
pub async fn export_encrypted(&self, ids: &[String], password: &str) -> Result<Vec<u8>> {
    // 导出为 AES 加密的 JSON 包，可选包含密码
    todo!("P1 功能，骨架期预留")
}
pub async fn import_encrypted(&self, data: &[u8], password: &str) -> Result<usize> {
    todo!("P1 功能，骨架期预留")
}
```

---

## 36. 测试架构

**Repo trait 策略**：`Database` 只提供通用 `call()` 方法，具体 SQL 在 Service 内部写，无法为 `Database` 统一 mock。正确做法是为每个 feature Service 定义**专属 Repo trait**，Service 依赖 Repo trait 而非具体 `Database`。测试 mock Repo，不 mock Database 内部实现。

```rust
// features/webshell/service.rs

#[async_trait::async_trait]
#[cfg_attr(test, mockall::automock)]
pub trait WebshellRepo: Send + Sync {
    async fn find_all(&self)                     -> Result<Vec<Webshell>>;
    async fn find_by_id(&self, id: &str)         -> Result<Webshell>;
    async fn find_by_url(&self, url: &str)       -> Result<Option<Webshell>>;  // 去重检查
    async fn insert(&self, w: &Webshell)         -> Result<()>;
    async fn update(&self, w: &Webshell)         -> Result<()>;
    async fn soft_delete(&self, id: &str)        -> Result<()>;
    /// 更新 status 和 fingerprint（test_connection / circuit breaker 回调路径）
    async fn update_status(
        &self,
        id:          &str,
        status:      &str,
        fingerprint: Option<&WebshellFingerprint>,
    ) -> Result<()>;
}

// 真实实现：DbWebshellRepo 包装 Database
pub struct DbWebshellRepo(Database);

#[async_trait::async_trait]
impl WebshellRepo for DbWebshellRepo {
    async fn find_all(&self) -> Result<Vec<Webshell>> {
        self.0.call(|conn| {
            // 具体 SQL ...
        }).await
    }
    // ...
}

// Service 依赖 trait，测试时可注入 Mock
pub struct WebshellService {
    repo:   Arc<dyn WebshellRepo>,
    crypto: Arc<CryptoContext>,
    queue:  Arc<WebshellQueue>,
    sessions: Arc<DashMap<String, WebshellSession>>,
}
```

**其他 feature 的 Repo trait 定义**（模式相同，均标注 `#[cfg_attr(test, mockall::automock)]`）：

```rust
// features/project/repo.rs
#[async_trait::async_trait]
#[cfg_attr(test, mockall::automock)]
pub trait ProjectRepo: Send + Sync {
    async fn find_all(&self)                      -> Result<Vec<Project>>;
    async fn find_by_id(&self, id: &str)          -> Result<Project>;
    async fn insert(&self, p: &Project)           -> Result<()>;
    async fn update(&self, p: &Project)           -> Result<()>;
    async fn soft_delete(&self, id: &str)         -> Result<()>;
    async fn soft_delete_with_webshells(&self, project_id: &str, ts: i64) -> Result<()>;
    async fn restore_with_webshells(&self, project_id: &str, ts: i64)     -> Result<()>;
}

// features/payload/repo.rs
#[async_trait::async_trait]
#[cfg_attr(test, mockall::automock)]
pub trait PayloadRepo: Send + Sync {
    async fn find_all(&self)                                      -> Result<Vec<Payload>>;
    async fn find_by_id(&self, id: &str)                          -> Result<Payload>;
    async fn insert(&self, p: &Payload)                           -> Result<()>;
    async fn update(&self, p: &Payload)                           -> Result<()>;
    async fn soft_delete(&self, id: &str)                         -> Result<()>;
    async fn insert_history(&self, h: &PayloadHistory)            -> Result<()>;
    async fn find_history_by_payload(&self, payload_id: &str)     -> Result<Vec<PayloadHistory>>;
    async fn find_history_by_webshell(&self, webshell_id: &str)   -> Result<Vec<PayloadHistory>>;
}

// features/plugin/repo.rs
#[async_trait::async_trait]
#[cfg_attr(test, mockall::automock)]
pub trait PluginRepo: Send + Sync {
    async fn find_all(&self)                           -> Result<Vec<Plugin>>;
    async fn find_by_id(&self, id: &str)               -> Result<Plugin>;
    async fn insert(&self, p: &Plugin)                 -> Result<()>;
    async fn update_enabled(&self, id: &str, v: bool)  -> Result<()>;
    async fn update_config(&self, id: &str, cfg: &serde_json::Value) -> Result<()>;
    async fn delete(&self, id: &str)                   -> Result<()>;  // 插件卸载为硬删除
}
```

```rust
// 每个 feature 都有对应的 Repo trait：
// WebshellRepo, ProjectRepo, PayloadRepo, PluginRepo
// 测试时：
#[cfg(test)]
mod tests {
    use super::*;
    use mockall::predicate::*;

    #[tokio::test]
    async fn test_create_webshell_encrypts_password() {
        let mut mock_repo = MockWebshellRepo::new();
        mock_repo.expect_insert()
            .withf(|w| w.password != "plaintext")  // 验证密码被加密
            .returning(|_| Ok(()));

        let service = WebshellService::new(
            Arc::new(mock_repo),
            CryptoContext::new([0u8; 32]),
            Arc::new(WebshellQueue::new_noop()),  // noop queue：不发网络请求，直接返回 Ok
        );
        // WebshellQueue::new_noop() 实现：内部 http_pool 为 MockHttpPool，send_raw 直接返回 Ok(dummy_response)
        service.create(CreateWebshellInput { password: "plaintext".into(), .. }).await.unwrap();
    }
}
```

集成测试使用 `:memory:` SQLite（`DbWebshellRepo(Database::open(":memory:").await?)`），无需 mock，验证真实 SQL 逻辑。

**两层测试策略**：
- 单元测试 → Mock Repo → 测试 Service 业务逻辑（加密、状态机、规范化）
- 集成测试 → :memory: DB → 测试 SQL 正确性（查询、索引、迁移）

---

## 37. 批量连接测试

批量命令执行用的少（P2），但**批量连接测试**是日常必需：确认哪些 WebShell 还存活，无需逐个手动测试。

```rust
// features/batch/service.rs
pub struct BatchService {
    webshell_service: Arc<WebshellService>,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct BatchTestResult {
    pub id:     String,
    pub result: ConnectionResult,
}

impl BatchService {
    /// 并发测试多个 WebShell，信号量限制并发防止批量请求被检测
    pub async fn test_connections(
        &self,
        ids: Vec<String>,
        concurrency: usize,  // 建议默认 5，可在设置中调整
    ) -> Vec<BatchTestResult> {
        let sem     = Arc::new(tokio::sync::Semaphore::new(concurrency));
        let service = Arc::clone(&self.webshell_service);

        let tasks: Vec<_> = ids.into_iter().map(|id| {
            let sem     = sem.clone();
            let service = service.clone();
            tokio::spawn(async move {
                let _permit = sem.acquire_owned().await.unwrap();
                BatchTestResult { result: service.test_connection(&id).await.unwrap_or_else(|e| {
                    ConnectionResult { success: false, latency_ms: None, error: Some(e.to_string()), fingerprint: None }
                }), id }
            })
        }).collect();

        futures::future::join_all(tasks).await
            .into_iter().filter_map(|r| r.ok()).collect()
    }
}
```

`BatchService` 加入 `AppState`：`batch_service: BatchService`。前端通过 `commands::batch::test_connections` 触发，结果通过正常 command 返回值返回（数量不大，无需流式推送）。

---

## 38. Cargo.toml 关键依赖

```toml
[dependencies]
tauri           = { version = "2", features = ["protocol-asset"] }
tokio           = { version = "1", features = ["full"] }
serde           = { version = "1", features = ["derive"] }
serde_json      = "1"
thiserror       = "1"
uuid            = { version = "1", features = ["v4"] }
chrono          = { version = "0.4", features = ["serde"] }
tracing         = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter"] }
tracing-appender   = "0.2"
tokio-rusqlite  = "0.5"                # async SQLite，内含 rusqlite bundled
reqwest         = { version = "0.12", features = ["json", "stream"] }
aes-gcm         = "0.10"   # 纯 Rust 实现，非 ring
argon2          = "0.5"
rand            = "0.8"    # 随机 salt、response_mark、混淆 seed
toml            = "0.8"    # config.toml 序列化/反序列化
zeroize         = { version = "1", features = ["derive"] }
dashmap         = "6"
encoding_rs     = "0.8"
url             = "2"      # URL 规范化（WebshellService::create 入口）
base64          = "0.22"   # queue_worker 中 body → base64 form 参数
futures         = "0.3"    # Stream trait（文件/终端流式响应）
bytes           = "1"      # Bytes 类型（流式数据块）
async-trait     = "0.1"    # async fn in traits（XxxRepo trait + mockall 兼容）
tree-sitter     = "0.22"
tree-sitter-php = "0.22"
mockall         = { version = "0.12", optional = true }
specta          = { version = "2", features = ["derive"] }           # 类型导出
tauri-specta    = { version = "2", features = ["derive", "typescript"] }  # TS 绑定生成
specta-typescript = "0.0.7"  # tauri_specta::Builder::export() 使用的 Typescript 格式化器
tauri-plugin-clipboard-manager = "2"   # 剪贴板写入（copyToClipboard 工具函数使用）

[features]
test-utils     = ["mockall"]
java-plugin    = []   # 需要 javac，build.rs 检测后自动激活（cargo:rustc-cfg）
dotnet-plugin  = []   # 需要 dotnet，build.rs 检测后自动激活

[profile.release]
debug = 1        # 保留行号信息，便于 panic 定位
strip = false
```

### 前端 package.json 关键依赖

```json
{
  "dependencies": {
    "vue":              "^3.4",
    "naive-ui":         "^2.38",
    "pinia":            "^2.1",
    "vue-router":       "^4.3",
    "vue-i18n":         "^9.13",
    "@tauri-apps/api":  "^2",
    "@tauri-apps/plugin-opener": "^2",
    "@tauri-apps/plugin-clipboard-manager": "^2"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5",
    "vite":               "^5",
    "typescript":         "^5",
    "vue-tsc":            "^2",
    "@tauri-apps/cli":    "^2"
  },
  "scripts": {
    "dev":   "vite",
    "build": "vue-tsc && vite build",
    "tauri": "tauri"
  }
}
```

`pnpm` 作为包管理器（与 Tauri 官方模板一致）。`naive-ui` 按需引入（`unplugin-auto-import` + `unplugin-vue-components`）以减小打包体积。

---

## 39. tauri-specta 自动类型生成

手动维护 Rust/TypeScript 双份类型是最常见的 bug 来源。`tauri-specta` 从 Rust 类型自动生成 `src/bindings.ts`。

**完整集成代码见第 9 节**（`tauri::Builder` 注册骨架）。tauri-specta 2.x 的正确用法是通过 `tauri_specta::Builder` 同时完成 command 收集和 invoke_handler 注册，**不使用** `tauri::generate_handler![]`，不使用旧版 `ts::export_with_cfg()`。

**所有需要导出的类型标注 `#[derive(specta::Type)]`**，命令函数标注 `#[tauri_specta::command]`。

**`AppError` 导出问题**：`AppError` 实现了自定义 `serde::Serialize`（输出 `{kind, message}`），但 specta 的 `#[derive(specta::Type)]` 依赖 derive 宏静态分析结构——自定义 Serialize 会导致 specta 推断的 TypeScript 类型与实际序列化 shape 不符。解决方案：定义一个专用 DTO 用于 specta 导出，Commands 签名不变：

```rust
// error.rs — 仅用于 specta 类型生成，不用于实际序列化
#[derive(specta::Type, serde::Serialize)]
pub struct AppErrorDto {
    pub kind:    String,
    pub message: String,
}

// collect_types! 中使用 AppErrorDto 替代 AppError
// Commands 仍然返回 Result<T, AppError>，Tauri 调用自定义 Serialize
// 前端 bindings.ts 里的 error 形状来自 AppErrorDto，与实际 JSON 吻合
```

`src/bindings.ts` 不手动编辑，加入 `.gitignore` 中的 "不追踪但不忽略"（实际用 gitattributes 标记为生成文件）。前端 import 从此改为：

```typescript
// 替代手动写的 features/webshell/types.ts
import type { Webshell, CreateWebshellInput } from '@/bindings'
```

**首次克隆冷启动问题**：`pnpm tauri dev` 会先启动 Vite dev server，Vite 编译时若 `src/bindings.ts` 不存在（首次 clone），会报 `Cannot find module '@/bindings'`。解决方案：在 git 中追踪一个**空占位文件** `src/bindings.ts`，内容为所有类型的 `export type X = any`（宽松，仅保证 Vite 能启动）。首次 `cargo build` 完成后自动覆盖为真实绑定。`package.json` 的 dev 脚本无需修改，占位文件兜底即可。

```typescript
// src/bindings.ts — git 追踪的占位，cargo build 后自动覆盖
// DO NOT EDIT — auto-generated by tauri-specta, committed as bootstrap placeholder
export type Webshell = any
export type CreateWebshellInput = any
export type UpdateWebshellInput = any
export type ConnectionResult = any
export type AppErrorDto = any
// ... 其余类型同样 any
```

---

## 40. system 命令模块

`commands/system.rs` 提供应用级只读查询，无需 `check_locked`：

```rust
// commands/system.rs

#[derive(Debug, Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct AppInfo {
    pub version:    String,   // tauri::VERSION 或 package.json version
    pub data_dir:   String,   // 数据目录路径（供用户打开文件夹）
    pub db_size_kb: u64,      // SQLite 文件大小（settings 页显示）
}

#[tauri_specta::command]
pub async fn get_app_info(
    app:   AppHandle,
    state: State<'_, AppState>,
) -> Result<AppInfo> {
    let meta = app.package_info();
    let db_size_kb = std::fs::metadata(&state.paths.db_path)
        .map(|m| m.len() / 1024).unwrap_or(0);
    Ok(AppInfo {
        version:    meta.version.to_string(),
        data_dir:   state.paths.data_dir.to_string_lossy().into_owned(),
        db_size_kb,
    })
}

/// 返回审计日志（按时间倒序，支持分页）
#[tauri_specta::command]
pub async fn get_audit_log(
    state:  State<'_, AppState>,
    limit:  u32,
    offset: u32,
) -> Result<Vec<AuditEvent>> {
    state.audit_log.query(limit, offset).await
}
```

`AppState` 需暴露 `pub paths: AppPaths`（目前 paths 仅在 `new()` 内部使用），或通过 `settings_service.data_dir()` 转发。

`AuditLog` 增加 `query` 方法：

```rust
// infra/audit.rs
#[derive(Debug, Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct AuditEvent {
    pub id:          String,
    pub webshell_id: Option<String>,
    pub action:      String,
    pub detail:      String,
    pub created_at:  i64,
}

impl AuditLog {
    pub async fn query(&self, limit: u32, offset: u32) -> Result<Vec<AuditEvent>> {
        self.db.call(move |conn| {
            let mut stmt = conn.prepare(
                "SELECT id, webshell_id, action, detail, created_at
                 FROM audit_events ORDER BY created_at DESC LIMIT ?1 OFFSET ?2"
            )?;
            // ... map rows to AuditEvent
        }).await
    }
}
```

**`AppState` 路径暴露**：将 `paths: AppPaths` 设为 `pub` 字段，`get_app_info` command 直接读取，无需额外 getter。

---

## 41. ConsoleView 控制台窗口布局

控制台窗口（`/console?id=xxx`）展示三个内置功能 Tab：

```
┌────────────────────────────────────────────────────────┐
│  [文件管理] [终端] [数据库]           连接状态 ●活跃    │
├────────────────────────────────────────────────────────┤
│                                                        │
│   Tab 内容区（占满剩余空间）                            │
│                                                        │
└────────────────────────────────────────────────────────┘
```

```vue
<!-- features/console/ConsoleView.vue -->
<template>
  <n-tabs v-model:value="activeTab" type="card" animated>
    <n-tab-pane name="file"     tab="文件管理"><FileManagerPane /></n-tab-pane>
    <n-tab-pane name="terminal" tab="终端">    <TerminalPane /></n-tab-pane>
    <n-tab-pane name="database" tab="数据库">  <DatabasePane /></n-tab-pane>
    <!-- 用户安装插件的 Tab（动态渲染） -->
    <n-tab-pane v-for="p in activePanes" :key="p.id" :name="p.id" :tab="p.name">
      <PluginPane :plugin-id="p.id" :webshell-id="webshellId" />
    </n-tab-pane>
  </n-tabs>
</template>
```

- 内置 Tab（文件/终端/DB）是 Vue 组件，直接调用 `consoleApi`，**不使用 iframe**
- 用户插件 Tab（`PluginPane`）使用 iframe 沙箱（见第 17 节）
- 控制台 store（`features/console/store.ts`）独立于主窗口 store，各 Tab 通过 `inject` 共享 `webshellId`
- 三个内置 Pane 组件路径：`features/console/file/FileManagerPane.vue`、`features/console/terminal/TerminalPane.vue`、`features/console/database/DatabasePane.vue`

---

## 42. HomeView 统计概览 + StatusBar

### HomeView

```vue
<!-- features/home/HomeView.vue -->
```

首页展示三类快速统计（数据来自各 feature store，App 启动时 `fetchAll` 已加载）：

| 卡片 | 数据来源 |
|------|---------|
| WebShell 总数 / 存活数 | `webshellStore.items.length` / `.filter(active)` |
| 项目数 | `projectStore.items.length` |
| 最近审计操作（最新 5 条） | `settingsStore.recentAudit`（首页 onMounted 单独 fetch） |

布局：三个统计卡片横排，下方审计列表。首页**不做额外后端 command**，复用已有 store 数据。

### StatusBar

```vue
<!-- layouts/components/StatusBar.vue -->
<template>
  <div class="statusbar">
    <span :class="statusDot(activeCount)">●</span>
    <span>{{ t('status.active', { n: activeCount }) }}</span>
    <n-divider vertical />
    <span>{{ t('status.projects', { n: projectCount }) }}</span>
    <n-divider vertical />
    <span>{{ t('status.payloads', { n: payloadCount }) }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useWebshellStore } from '@/features/webshell/store'
import { useProjectStore }  from '@/features/project/store'
import { usePayloadStore }  from '@/features/payload/store'

const ws = useWebshellStore()
const ps = useProjectStore()
const pl = usePayloadStore()

const activeCount  = computed(() => ws.items.filter(w => w.status === 'active').length)
const projectCount = computed(() => ps.items.length)
const payloadCount = computed(() => pl.items.length)
</script>
```

StatusBar 是纯展示组件，读取 store 计算属性，**不直接 invoke**。

---

## 43. 剪贴板工具

```typescript
// shared/utils/clipboard.ts
import { writeText } from '@tauri-apps/plugin-clipboard-manager'

/**
 * 复制文本到剪贴板。
 * sensitive=true：复制成功后 5 秒自动清除，防止密码残留在系统剪贴板。
 */
export async function copyToClipboard(text: string, sensitive = false): Promise<void> {
  await writeText(text)
  if (sensitive) {
    setTimeout(() => writeText(''), 5000)
  }
}
```

用法（右键菜单"复制密码"）：
```typescript
// WebshellContextMenu.vue
import { copyToClipboard } from '@/shared/utils/clipboard'

async function onCopyPassword() {
  await copyToClipboard(props.webshell.password, true)  // sensitive=true → 5s 后清除
  globalMessage.success(t('action.copied'))
}
```

`@tauri-apps/plugin-clipboard-manager` 需加入 `package.json`（`"^2"`）并在 `capabilities/default.json` 中添加 `"clipboard-manager:allow-write-text"` 权限。

---

## 45. 实现优先级

骨架阶段严格按以下顺序，每步完成后再进入下一步：

1. **AppError + 全局 Result** — 所有模块依赖，最先定义
2. **Cargo.toml + package.json** — 声明全部依赖
3. **tauri.conf.json + capabilities** — 权限和窗口配置
4. **infra 层**：paths → logger → panic hook → crypto → db（含迁移）→ config（含迁移）→ http（HttpClientPool + send_raw）→ audit_log → c2_profile
5. **Repo trait 层**：`WebshellRepo`、`ProjectRepo`、`PayloadRepo`、`PluginRepo` trait 定义 + `DbXxxRepo` 真实实现（包裹 `Database::call()`）。位于各 feature 目录的 `repo.rs`，trait 使用 `async-trait`，标注 `#[cfg_attr(test, mockall::automock)]`
6. **bootstrap() 函数** — 串联所有基础设施，首次运行初始化
7. **AppState** — 整合所有 Service（使用 `Arc<WebshellService>`），注册到 Tauri
8. **features 目录骨架** — 各模块 Service 空实现（struct + 方法签名 + todo!()），依赖注入 Repo trait
9. **commands 层注册** — tauri_specta::Builder 收集所有 command，验证编译通过
10. **前端 main.ts + App.vue** — i18n、Pinia、Router、主题绑定
11. **AppLayout + CustomTitlebar + NavSidebar + StatusBar** — 整体布局
12. **各 feature 的 store + api.ts 空实现** — 前后端通信链路验证
13. **WebShell 通信协议** — Init/Exec 两阶段，含响应验证、透明重试、MemShellExpired 分支
14. **混淆器骨架** — PhpObfuscator level 1（变量重命名），其他语言模板占位
15. **C2 Profile default** — 默认 Profile + request_wrapper 通路跑通

---

*设计文档完成，实现阶段参考本文档，不在代码中重复解释设计决策。*
