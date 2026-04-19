# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

FG-ABYSS (非攻 - 渊渟) is a modern WebShell management tool built with **Tauri V2 + Rust + Vue 3 + TypeScript**. The goal is to surpass Godzilla (哥斯拉) in features and design quality.

**Status**: Skeleton scaffolding phase. Authoritative design spec:
`docs/superpowers/specs/2026-04-17-fg-abyss-skeleton-design.md`

Read this spec before making any architectural decisions. It supersedes `docs/DEV.md` and `docs/PRD.md` where they conflict.

---

## Build Commands

```bash
# Development (hot reload)
pnpm tauri dev

# Production build
pnpm tauri build

# Frontend only
pnpm dev

# Rust tests
cargo test --manifest-path src-tauri/Cargo.toml

# Rust tests with mock infra
cargo test --manifest-path src-tauri/Cargo.toml --features test-utils
```

---

## Architecture

Feature-based Layered Architecture — `commands/` → `features/` → `infra/`.

**Hard rules:**
- `commands/` layer: thin, no business logic, the only place `AppHandle` is used
- `features/` services: no `AppHandle` dependency, fully unit-testable via mock traits
- `features/` modules: no cross-imports between feature modules (View layer composes them)
- All infra dependencies injected via traits (`#[cfg_attr(test, mockall::automock)]`)

**AppState** is the DI root — initialized once in `bootstrap()`, registered via `tauri::Builder::manage()`.

---

## Key Infrastructure Decisions

| Concern | Choice | Reason |
|---------|--------|--------|
| DB access | `r2d2` + `r2d2_sqlite` connection pool | WAL mode multi-read concurrency |
| Crypto | `aes-gcm` crate (pure Rust) | Not `ring` |
| Async DB calls | Direct pool `.get()` | No `spawn_blocking` needed with pool |
| Logging | `tracing` + `tracing_appender` daily rolling | `WorkerGuard` held in `AppState._log_guard` |
| Password memory safety | `zeroize` on drop | Prevent heap scraping |
| Secrets in logs | `Sensitive<T>` wrapper | `Debug` impl prints `[REDACTED]` |
| Charset detection | `encoding_rs` | GBK/GB2312 servers |
| PHP obfuscation | `tree-sitter` + `tree-sitter-php` | AST-level transforms |

---

## WebShell Communication Protocol

Two-phase protocol (Init → Exec):

1. **Init**: client sends plugin bytecode/code encrypted in request body; server dynamically loads it and returns `session_id` + `response_mark`
2. **Exec**: subsequent calls reference `session_id`; responses wrapped in `response_mark` for authenticity verification

If `response_mark` is absent in a response → throw `AppError::InvalidResponse`, do not decrypt.

Session expiry → transparent re-Init, invisible to the caller.

---

## Database Schema Conventions

- All business tables have `deleted_at INTEGER` (soft delete). Queries always append `WHERE deleted_at IS NULL`.
- UUIDs generated in the Service layer (`Uuid::new_v4()`), never passed in from frontend.
- Migrations in `src-tauri/migrations/` numbered `001_`, `002_`, etc.
- Run `db.vacuum_if_needed()` at startup when soft-deleted count > 100.

---

## Frontend Conventions

- **LoadingMap** (`Record<string, boolean>`) instead of single `loading: boolean` — allows `loading.list`, `loading.create`, `loading['delete-uuid']` to be independent.
- **Error display policy**: toast (`n-message`) for transient errors; modal (`n-dialog`) for unrecoverable errors or destructive confirmations; inline for form validation.
- Components never call `api.ts` directly — all data ops go through store actions.
- Store actions catch errors and display them; components only read store state.
- Stores never import other feature stores.

---

## Streaming Responses

For file transfer and terminal output, Service layer returns a `Stream`, never takes `AppHandle`. The `commands/` layer owns the stream → Tauri event loop:

```rust
// Service: returns stream
pub async fn download_file_stream(&self, ...) -> Result<impl Stream<Item = Result<Bytes>>>

// Command: drives stream, emits progress events
app.emit_to(&label, "file-progress", event)?;
```

---

## C2 Profile & Encoding Chain

Each WebShell can have:
- A **C2 Profile** (`infra/c2_profile.rs`): controls request param name, headers, User-Agent, response wrapper format, timing jitter
- A **CryptoChain** (`infra/crypto.rs`): ordered `Vec<CodecStep>` (AES-GCM, XOR, Base64, Gzip, etc.) — encode on send, reverse-decode on receive

Payload generator must emit matching server-side decode code for each chain config.

---

## Multi-Window Strategy

Console windows: label = `console-{webshell_id}`, URL = `/console?id={webshell_id}`.

On window close → listen `window-destroyed` → call `console_service.cleanup(webshell_id)`.

On window already open → `set_focus()` instead of creating new.
