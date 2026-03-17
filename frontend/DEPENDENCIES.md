# 前端依赖更新说明

## 新增依赖

为了支持 WebShell 终端功能，需要安装以下依赖：

### 核心依赖

```bash
npm install @xterm/xterm @xterm/addon-fit @xterm/addon-web-links
```

**说明：**
- `@xterm/xterm`: 终端模拟器核心库
- `@xterm/addon-fit`: 终端自适应大小插件
- `@xterm/addon-web-links`: Web 链接检测插件

### 替代方案（可选）

如果遇到问题，可以使用以下替代方案：

**方案 1：使用稳定版本**
```bash
npm install xterm xterm-addon-fit xterm-addon-web-links
```

然后修改导入：
```typescript
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { WebLinksAddon } from 'xterm-addon-web-links'
import 'xterm/css/xterm.css'
```

**方案 2：简化版本（不使用 xterm）**
如果安装依赖有问题，可以使用简单的文本区域替代：
- 移除 Terminal 相关代码
- 使用 `<textarea>` 或 `<pre>` 显示输出
- 功能会简化，但基本命令执行可用

## 安装步骤

1. 进入 frontend 目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install @xterm/xterm @xterm/addon-fit @xterm/addon-web-links
```

3. 验证安装
```bash
npm list @xterm/xterm
```

4. 重新启动开发服务器
```bash
npm run dev
```

## 版本要求

- Node.js: >= 16.0.0
- npm: >= 7.0.0

## 注意事项

1. **依赖冲突**：如果已有旧版本 xterm，建议先卸载
   ```bash
   npm uninstall xterm xterm-addon-fit xterm-addon-web-links
   ```

2. **构建问题**：如果遇到构建问题，清理缓存
   ```bash
   rm -rf node_modules package-lock.json
   npm install
   ```

3. **TypeScript 类型**：新版本已包含类型定义，无需额外安装 @types/xterm

## 测试

安装完成后，访问 WebShell 页面测试终端功能：
- 终端应该正常显示
- 支持 ANSI 颜色
- 支持复制粘贴
- 支持滚动

## 故障排除

如果终端无法显示，检查：
1. 浏览器控制台是否有错误
2. 终端容器是否有正确的高度
3. 是否正确调用了 `terminal.open()` 和 `fitAddon.fit()`
