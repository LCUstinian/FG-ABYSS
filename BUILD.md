# FG-ABYSS 构建指南

## 系统要求

### Windows
- Windows 10/11 (64-bit)
- Visual Studio 2022 with C++ tools
- Rust 1.75+ (通过 rustup 安装)
- Node.js 18+ (LTS 版本)

### Linux
- Ubuntu 20.04+ / Debian 11+ / Fedora 36+
- GCC, make, pkg-config, libssl-dev
- Rust 1.75+
- Node.js 18+

### macOS
- macOS 11+ (Big Sur 或更高版本)
- Xcode Command Line Tools
- Rust 1.75+
- Node.js 18+

## 安装依赖

### 1. 安装 Rust

```bash
# Windows (PowerShell)
winget install Rustlang.Rustup
# 或访问 https://rustup.rs 下载安装

# Linux/macOS
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

### 2. 安装 Node.js

```bash
# Windows (PowerShell)
winget install OpenJS.NodeJS.LTS

# Linux (使用 nvm)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18

# macOS (使用 Homebrew)
brew install node@18
```

### 3. 安装系统依赖

#### Windows
```powershell
# 安装 Visual Studio Build Tools
winget install Microsoft.VisualStudio.2022.BuildTools

# 安装 WebView2 (通常已预装)
winget install Microsoft.WebView2
```

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install -y build-essential pkg-config libssl-dev libgtk-3-dev libayatana-appindicator3-dev librsvg2-dev
```

#### Fedora
```bash
sudo dnf install -y gcc make pkg-config openssl-devel gtk3-devel libappindicator-gtk3-devel librsvg2-devel
```

#### macOS
```bash
xcode-select --install
brew install pkg-config openssl gtk3
```

## 构建步骤

### 1. 克隆项目

```bash
git clone https://github.com/fg-abyss/fg-abyss.git
cd fg-abyss
```

### 2. 安装前端依赖

```bash
npm install
```

### 3. 开发模式

```bash
# 启动开发服务器（热重载）
npm run tauri:dev

# 或分别启动
npm run dev        # Vite 开发服务器
npm run tauri dev  # Tauri 开发模式
```

### 4. 生产构建

```bash
# 构建所有平台
npm run tauri:build

# 构建特定平台
# Windows
npm run tauri build -- --target x86_64-pc-windows-msvc

# Linux
npm run tauri build -- --target x86_64-unknown-linux-gnu

# macOS
npm run tauri build -- --target x86_64-apple-darwin
```

## 构建输出

构建完成后，产物位于：

```
src-tauri/target/release/
├── bundle/
│   ├── msis/           # Windows MSI 安装包
│   ├── nsis/           # Windows NSIS 安装包
│   ├── deb/            # Linux DEB 包
│   ├── appimage/       # Linux AppImage
│   ├── dmg/            # macOS DMG
│   └── macos/          # macOS .app
└── fg-abyss(.exe)      # 可执行文件
```

## 构建优化

### 减小二进制文件大小

项目已配置以下优化（见 `src-tauri/Cargo.toml`）：

```toml
[profile.release]
panic = "abort"           # 更小的二进制
codegen-units = 1         # 更好的优化
lto = true                # 链接时优化
opt-level = "s"           # 优化大小
strip = true              # 移除调试信息
```

### 构建时间优化

```bash
# 使用 sccache 加速编译
cargo install sccache

# 在 .cargo/config.toml 中配置
echo '[build]
rustc-wrapper = "sccache"' >> src-tauri/.cargo/config.toml
```

## 常见问题

### 1. Rust 版本不兼容

```bash
# 切换到指定版本
rustup install 1.75.0
rustup default 1.75.0
```

### 2. Node.js 版本问题

```bash
# 使用 nvm 切换版本
nvm install 18
nvm use 18
```

### 3. Windows 构建错误

```powershell
# 以管理员身份运行 PowerShell
# 安装必要的证书
certutil -generateSSTFromWU root.sst
Import-Certificate -FilePath root.sst -CertStoreLocation Cert:\LocalMachine\Root
```

### 4. Linux 依赖缺失

```bash
# Ubuntu/Debian
sudo apt install libwebkit2gtk-4.0-dev libgtk-3-dev libayatana-appindicator3-dev

# Fedora
sudo dnf install webkit2gtk3-devel gtk3-devel libappindicator-gtk3-devel
```

### 5. macOS 签名问题

```bash
# 开发模式可以跳过签名
# 生产环境需要配置证书
export APPLE_CERTIFICATE_NAME="Developer ID Application: Your Name"
export APPLE_TEAM_ID="YOUR_TEAM_ID"
```

## 代码签名

### Windows

1. 获取代码签名证书（如 DigiCert, Sectigo）
2. 在 `tauri.conf.json` 中配置：

```json
{
  "tauri": {
    "bundle": {
      "windows": {
        "certificateThumbprint": "YOUR_CERT_THUMBPRINT",
        "timestampUrl": "http://timestamp.digicert.com"
      }
    }
  }
}
```

### macOS

1. 在 Apple Developer 创建证书
2. 配置环境变量：

```bash
export APPLE_CERTIFICATE_NAME="Developer ID Application: Your Name"
export APPLE_TEAM_ID="YOUR_TEAM_ID"
export APPLE_SIGNING_IDENTITY="YOUR_SIGNING_IDENTITY"
```

## 性能基准

在典型配置下（Intel i7, 16GB RAM）：

- **开发构建**: ~2-3 分钟
- **生产构建**: ~5-8 分钟
- **二进制大小**: ~15-25 MB (取决于平台)

## 持续集成

项目包含 GitHub Actions 配置（`.github/workflows/build.yml`），支持：

- 自动构建 Windows/Linux/macOS
- 自动运行测试
- 自动创建 Release
- 自动上传产物

## 下一步

构建完成后，请查看：
- `README.md` - 项目介绍和使用指南
- `docs/` - 详细文档
- `CHANGELOG.md` - 版本更新日志
