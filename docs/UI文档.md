📘 FG-ABYSS (非攻 - 渊渟) UI 设计规格说明书
版本号: v1.0.0 (Final)
适用技术栈: Rust + Tauri + Vue 3 + TypeScript + Naive UI
设计风格: 沉浸式暗黑 · 专业红队工具 · 紫色极光强调色
文档状态: ✅ 定稿 (基于用户确认的界面原型)

1. 设计系统规范 (Design System)
本系统采用 Dark Mode First 策略，所有颜色定义均基于深色环境。
1.1 色彩体系 (Color Palette)
表格
语义	变量名	Hex 值	用途说明
Primary (主色)	--primary-gradient	#7B61FF
 → #A78BFA	核心特征：左侧导航选中态、关键按钮、高亮边框的线性渐变。
Secondary (次色)	--color-secondary	#3B82F6	次要操作按钮（如“新建 WebShell”）、链接。
Success (成功)	--color-success	#52C41A	Shell 在线状态、连接成功、签名验证通过。
Warning (警告)	--color-warning	#FAAD14	混淆等级 L2、代理延迟高、部分活跃状态。
Danger (危险)	--color-danger	#FF4D4F	删除操作、WAF 拦截、Shell 离线/死亡、混淆 L3。
Info (信息)	--color-info	#00E5FF	提示性文字、辅助图标。
Bg Base (基底)	--bg-base	#0B0C10	应用最底层背景，近乎黑色的深蓝灰。
Bg Surface 1	--bg-surface-1	#1F2833	左侧导航栏、底部状态栏、卡片背景。
Bg Surface 2	--bg-surface-2	#2C3542	输入框背景、表格行悬停、模态框内容区。
Text Primary	--text-main	#E0E6ED	标题、正文主要信息。
Text Secondary	--text-sub	#8B9BB4	辅助说明、占位符、禁用状态文字。
Border	--border-color	#333333	分割线、卡片边框。
1.2 字体排印 (Typography)
● 中文字体: Inter, PingFang SC, Microsoft YaHei, sans-serif
● 代码/数字字体: JetBrains Mono, Roboto Mono, monospace (用于 Shell 列表、代码预览、状态栏)
● 字号规范:
  ○ H1 (页面标题): 20px / Bold (font-weight: 600)
  ○ H2 (模块标题): 16px / Semi-Bold
  ○ Body (正文): 14px / Regular
  ○ Caption (辅助): 12px / Regular (颜色 --text-sub)
  ○ Code: 13px / Mono (行高 1.6)
1.3 形状与质感 (Shape & Texture)
● 圆角 (Radius):
  ○ 按钮/输入框: 8px
  ○ 卡片/容器: 10px
  ○ 模态框: 12px
● 阴影 (Shadow):
  ○ 默认无阴影，依靠背景色差区分层级。
  ○ 悬浮层/模态框: 0 8px 24px rgba(0, 0, 0, 0.5)
  ○ 光晕效果: 选中态元素增加 box-shadow: 0 0 8px rgba(123, 97, 255, 0.3)
1.4 图标 (Iconography)
● 库: Lucide Vue Next (Outline 风格)
● 线宽: 1.5px
● 尺寸: 16px (列表内), 20px (导航), 24px (空状态)

2. 全局布局架构 (Global Layout)
2.1 框架结构
采用经典的 Tauri 桌面应用布局：
文本
编辑
[ 自定义标题栏 (32px) ]
-------------------------------------------
[ 左侧导航 (240px) ] [ 右侧内容区 (自适应) ]
-------------------------------------------
[ 底部状态栏 (28px) ]
2.2 自定义标题栏 (Custom TitleBar)
● 高度: 32px
● 背景: --bg-surface-1
● 左侧:
  ○ Logo 图标 + 文本 "FG-ABYSS" (作为拖拽区域 app-region: drag)
● 右侧:
  ○ 主题切换图标 (☀️/🌙)
  ○ 语言切换 (CN/US)
  ○ 窗口控制组 (最小化 _ , 最大化 □, 关闭 ×) - 关闭按钮悬停变红
2.3 左侧导航栏 (Sidebar)
● 宽度: 240px (固定)
● 背景: --bg-surface-1
● 菜单项:
  ○ 垂直排列：首页、项目、载荷、插件、设置
  ○ 选中态: 背景为 紫色线性渐变 (linear-gradient(90deg, #7B61FF, #A78BFA))，文字白色，图标白色。
  ○ 默认态: 文字 --text-sub，悬停时文字变 --text-main，背景微亮。
2.4 底部状态栏 (Status Bar)
● 高度: 28px
● 背景: --bg-surface-1，顶部有 1px 分割线 --border-color
● 布局 (Flex Space-Between):
  ○ Left: 内存监控 8.5 GB / 16 GB
  ○ Center: 进程信息 ID: 12345 | CPU 25% | 延迟 2s
  ○ Right: 运行时长 01:34:58

3. 核心页面详细设计 (Page Specifications)
3.1 项目模块 (Project Module)
对应截图: 项目管理与 WebShell 列表
A. 布局结构
● 左侧子栏 (项目树): 宽度约 240px (含在左侧导航旁或独立分栏，视具体实现而定，截图中为独立深灰区域)。
  ○ 顶部: 紫色圆角矩形按钮 [+] (新建项目)。
  ○ 列表: 文件夹图标 + 项目名称 + 删除图标 (垃圾桶)。
    ■ 选中态: 紫色渐变背景。
  ○ 底部: “回收站”分割线 + 橙色圆角按钮 [↺ 恢复项目]。
● 右侧主区 (WebShell 列表):
  ○ 顶部工具栏:
    ■ 左：搜索框 [🔍 搜索 WebShell...]
    ■ 右：统计文本 (总数:0 活跃:0 非活跃:0) + [🗑️ 回收站] 按钮 + [+] 蓝色新建按钮 + [5 条/页 ▼] 分页器。
  ○ 内容区:
    ■ 空状态: 中央显示纸箱图标 📦 + 文本“暂无数据” + “点击'+'按钮创建新的 WebShell”。
    ■ 数据态: 表格展示 (状态点、名称、URL、脚本类型、最后活跃、操作)。
  ○ 底部分页: < 1 > 样式，当前页码蓝色高亮。
B. 交互逻辑
● 右键菜单: 在 WebShell 行上右键，弹出深色菜单：[连接控制台] [清理缓存] [编辑] [移入回收站] [导出]。
● 软删除: 点击删除仅标记 is_deleted=true，需在回收站视图恢复或彻底删除。

3.2 载荷工厂模块 (Payload Factory)
对应截图: WebShell 生成器
A. 布局结构
● 顶部 Tabs: 载荷生成 (激活) | 载荷历史 | 载荷模板
● 双栏布局:
  ○ 左侧 (配置区):
    ■ 生成模式: 单选按钮组 (文件通用/代理/混合...)。
    ■ 脚本语言: 标签式单选 (PHP/JSP/ASPX/ASP)，选中为蓝色背景。
    ■ 加密算法: 标签式单选 (无加密/AES-CBC/AES-GCM✅)。
    ■ 凭证配置:
      ● 密码输入框 + [随机生成] 按钮。
      ● 文件名输入框 + [智能推荐] 按钮。
    ■ 混淆强度:
      ● 滑块控件: 从 L1 到 L3。
      ● 刻度说明: L1(变量重命名) - L2(垃圾代码) - L3(控制流平坦化)。
      ● 实时反馈: 下方显示“安全等级预估”卡片 (如：高安全性，推荐实战使用)。
  ○ 右侧 (预览区):
    ■ 头部: 标题“代码预览” + [📋 复制] [⬇️ 下载] 图标按钮。
    ■ 内容: 深色代码编辑器区域。
      ● 未生成时: 显示图标 + “点击生成按钮预览代码”。
      ● 生成后: 显示高亮代码。
    ■ 底部警示: 黄色背景卡片 ⚠️ “本工具仅供授权渗透测试...严禁非法用途”。
B. 交互逻辑
● 联动: 修改脚本语言自动过滤支持的加密算法；拖动混淆滑块实时更新安全评估文案。
● 生成: 点击底部“生成载荷”大按钮（截图中未完全显示，推测在左下角），右侧立即渲染代码。

3.3 系统设置模块 (Settings)
对应截图: 系统配置和外观设置
A. 布局结构
● 左侧子菜单: 外观 (激活) | 连接 | 关于
● 右侧配置面板 (卡片式):
  ○ 主题卡片:
    ■ 三列布局：浅色模式 | 深色模式 (选中，带紫色边框和对勾) | 跟随系统。
    ■ 每个选项带图标 (太阳/月亮/屏幕)。
  ○ 语言卡片:
    ■ 两列布局：CN 中文 (选中) | US English。
  ○ 全局强调色:
    ■ 一行彩色圆点选择器 (蓝/紫/粉/橙/绿...)。
    ■ 点击圆点即时预览界面变色。
B. 扩展设计 (连接/关于 - 基于需求推导)
● 连接设置:
  ○ 表单形式：超时时间 (Input)、重试次数 (Input)、User-Agent 池 (Textarea + 导入按钮)。
  ○ 代理配置：开关 + 类型下拉 (HTTP/SOCKS5) + 主机/端口输入。
● 关于页面:
  ○ 版本信息卡片、作者介绍、开源协议链接、[检查更新] 按钮。

4. 交互与状态定义 (Interaction & States)
4.1 加载状态 (Loading)
● 骨架屏: 列表加载时，显示灰色条纹流动的骨架块。
● 按钮加载: 按钮文字消失，变为旋转 Spinner，禁用点击。
● 全屏加载: 启动时显示 Logo 呼吸灯动画。
4.2 空状态 (Empty States)
● 通用样式: 中央图标 (线性风格，颜色 --text-sub) + 主标题 + 副标题引导 + 行动按钮。
● 文案示例: “暂无数据，点击 '+' 按钮创建新的 WebShell”。
4.3 异常反馈 (Error Handling)
● Toast 通知:
  ○ 位置：右上角堆叠。
  ○ 样式：深色背景，左侧彩色竖条 (红/黄/绿)，自动 3s 消失。
  ○ 内容：[图标] 标题：详细描述 (错误码)。
● 模态框: 危险操作 (如彻底删除) 需二次确认，按钮为红色。
4.4 微交互 (Micro-interactions)
● Hover: 卡片/行悬停时，背景变为 --bg-surface-2，轻微上浮 translateY(-2px)。
● Active: 按钮点击时有缩放 scale(0.98) 效果。
● Transition: 所有颜色变化、宽度变化需添加 transition: all 0.2s ease。

5. 给开发者的实施建议 (Implementation Notes)
5.1 技术栈组件映射 (Naive UI)
表格
界面元素	推荐组件	备注
整体布局	n-layout
, n-layout-sider	开启 native-scrollbar=false
 自定义滚动条
导航菜单	n-menu	自定义 render-label
 实现渐变背景
数据表格	n-data-table	必须开启 virtual-scroll
 优化性能
表单输入	n-input
, n-select
, n-radio-group	统一设置 size="medium"
滑块	n-slider	自定义 mark
 显示 L1/L2/L3
代码编辑	vue-monaco-editor	主题设为 vs-dark
，只读模式
模态框	n-modal	preset="card"
，去除默认 padding
提示	n-message
, n-notification	全局挂载
5.2 性能优化策略
1. 虚拟滚动: 项目列表和 WebShell 列表若超过 50 条，必须启用虚拟滚动。
2. 状态节流: 底部状态栏的 CPU/内存数据更新频率限制在 1Hz (每秒 1 次)，避免频繁重绘。
3. 按需加载: 载荷工厂的代码高亮组件仅在切换到该 Tab 时加载。
4. Rust 通信: 大量日志或列表数据通过 Tauri Event 批量发送，避免逐条 IPC 调用。
5.3 CSS 变量定义 (Tailwind/SCSS 参考)
css
编辑
:root {
  /* Colors */
  --abyss-primary-start: #7B61FF;
  --abyss-primary-end: #A78BFA;
  --abyss-bg-base: #0B0C10;
  --abyss-bg-surface-1: #1F2833;
  --abyss-bg-surface-2: #2C3542;
  --abyss-text-main: #E0E6ED;
  --abyss-text-sub: #8B9BB4;
  --abyss-border: #333333;
  
  /* Gradients */
  --gradient-primary: linear-gradient(90deg, var(--abyss-primary-start), var(--abyss-primary-end));
  
  /* Radius */
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
}

/* 全局暗黑重置 */
body {
  background-color: var(--abyss-bg-base);
  color: var(--abyss-text-main);
  font-family: 'Inter', sans-serif;
}
5.4 安全与合规细节
● 剪贴板: 复制敏感信息 (密码/Key) 后，调用 Rust 后端在 60s 后清空剪贴板。
● 输入保护: 所有密码输入框设置 autocomplete="new-password" 防止浏览器缓存。
● 日志脱敏: 前端渲染日志时，正则替换 URL 中的 password= 参数为 ***。

✅ 交付清单
1. 设计规范: 包含完整的色值、字体、圆角定义。
2. 页面原型描述: 针对项目、载荷、设置三大核心模块的详细布局与交互说明。
3. 开发指南: 具体的 Naive UI 组件选型与性能优化建议。