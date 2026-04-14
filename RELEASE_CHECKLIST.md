# FG-ABYSS v0.1.0 发布检查清单

## 发布前检查

### 代码质量
- [x] 所有 Rust 代码通过编译
- [x] 所有 TypeScript 代码通过类型检查
- [x] 运行所有单元测试并 passing
- [x] 代码格式化（`cargo fmt` / `npm run lint`）
- [x] 无严重警告

### 功能完整性
- [x] 项目管理功能正常
- [x] WebShell 管理功能正常
- [x] 载荷管理功能正常
- [x] 加密通信功能正常
- [x] 控制台窗口功能正常
- [x] 插件系统功能正常
- [x] 文件管理功能正常
- [x] 数据库管理功能正常
- [x] 终端管理功能正常
- [x] 设置页面功能正常

### 测试覆盖
- [x] Rust 单元测试（50+ 用例）
- [x] 前端组件测试
- [x] 集成测试
- [x] 手动测试关键路径

### 文档完整性
- [x] README.md 更新
- [x] BUILD.md 创建
- [x] RELEASE.md 创建
- [x] CHANGELOG.md 更新
- [x] API 文档（如适用）

### 安全性检查
- [x] 无硬编码密钥
- [x] 敏感信息脱敏
- [x] 依赖项安全扫描
- [x] 权限配置正确

### 性能检查
- [x] 二进制大小 < 30MB
- [x] 启动时间 < 3 秒
- [x] 内存使用 < 200MB
- [x] 无内存泄漏

## 构建步骤

### 1. 版本号确认

```bash
# 检查 package.json
cat package.json | grep version

# 检查 Cargo.toml
cat src-tauri/Cargo.toml | grep version

# 应该都是 0.1.0
```

### 2. 创建 Git 标签

```bash
git tag -a v0.1.0 -m "Release v0.1.0 - Initial public beta"
git push origin v0.1.0
```

### 3. 触发 CI/CD

推送标签后，GitHub Actions 自动执行：
- Windows 构建
- Linux 构建
- macOS 构建
- 创建 Release

### 4. 本地验证构建

```bash
# Windows
npm run tauri build

# 检查输出
ls src-tauri/target/release/bundle/
```

### 5. 测试安装包

- [ ] Windows MSI 安装测试
- [ ] Linux DEB 安装测试
- [ ] macOS DMG 安装测试
- [ ] 便携版测试

## 发布步骤

### GitHub Release

1. 访问 https://github.com/fg-abyss/fg-abyss/releases
2. 编辑自动创建的 draft release
3. 添加发布说明（从 RELEASE.md）
4. 上传所有构建产物
5. 标记为 latest release
6. 发布

### 文档更新

1. 更新官网下载链接
2. 更新 README 徽章
3. 更新文档站点
4. 发送发布通知

### 通知渠道

- [ ] GitHub Discussions
- [ ] Twitter/X
- [ ] Reddit (r/rust, r/webdev)
- [ ] Discord 社区
- [ ] 邮件列表（如有）

## 发布后检查

### 24 小时内
- [ ] 监控 GitHub Issues
- [ ] 回复用户反馈
- [ ] 收集崩溃报告
- [ ] 统计下载量

### 一周内
- [ ] 分析用户反馈
- [ ] 整理 bug 列表
- [ ] 规划 v0.2.0
- [ ] 发布修复版本（如需要）

## 回滚计划

如果出现严重问题：

1. 立即下架 release
2. 创建 hotfix 分支
3. 修复问题
4. 重新发布 v0.1.1
5. 通知用户升级

## 成功标准

- [ ] 所有平台构建成功
- [ ] 安装包测试通过
- [ ] 无严重 bug 报告
- [ ] 文档可访问
- [ ] 社区反馈积极

## 联系人

- **发布负责人**: @TODO
- **技术支持**: @TODO
- **文档维护**: @TODO

---

**发布状态**: ✅ 准备就绪
**最后更新**: 2024-01-01
