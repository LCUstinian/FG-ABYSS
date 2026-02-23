<script setup lang="ts">
import TitleBar from './components/TitleBar.vue'
import { ref, onMounted, computed } from 'vue'
import { Home, Folder, Package, Plug, Settings, Server, Cpu, Database } from 'lucide-vue-next'

const isDarkTheme = ref(localStorage.getItem('theme') === 'dark' || window.matchMedia('(prefers-color-scheme: dark)').matches)
const currentNavItem = ref('home')

// 导航项配置
const navItems = [
  { id: 'home', label: '首页', icon: Home },
  { id: 'projects', label: '项目', icon: Folder },
  { id: 'payloads', label: '载荷', icon: Package },
  { id: 'plugins', label: '插件', icon: Plug },
  { id: 'settings', label: '设置', icon: Settings }
]

// 系统状态数据（模拟）
const systemStatus = ref({
  memoryUsage: '1.2 GB / 8 GB',
  processId: '12345',
  cpuUsage: '15%',
  uptime: '2 hours'
})

// 计算当前内容组件
const currentContent = computed(() => {
  return currentNavItem.value
})

// 切换导航项
const switchNavItem = (itemId: string) => {
  currentNavItem.value = itemId
}

onMounted(() => {
  document.documentElement.classList.toggle('dark', isDarkTheme.value)
})
</script>

<template>
  <div class="app-container">
    <TitleBar />
    <div class="main-content">
      <!-- 左边导航区 -->
      <div class="sidebar">
        <div class="nav-item" 
             v-for="item in navItems" 
             :key="item.id"
             :class="{ active: currentNavItem === item.id }"
             @click="switchNavItem(item.id)"
        >
          <component :is="item.icon" :size="24" />
          <span>{{ item.label }}</span>
        </div>
      </div>
      
      <!-- 右边内容区 -->
      <div class="content-area">
        <!-- 首页内容 -->
        <div v-if="currentContent === 'home'" class="content-section">
          <h1>欢迎使用 FG-ABYSS</h1>
          <div class="home-content">
            <div class="home-card">
              <Server :size="48" />
              <h3>项目介绍</h3>
              <p>FG-ABYSS 是一个功能强大的安全工具，用于管理项目、生成载荷和使用插件。</p>
            </div>
            <div class="home-card">
              <Cpu :size="48" />
              <h3>系统状态</h3>
              <p>内存: {{ systemStatus.memoryUsage }}</p>
              <p>CPU: {{ systemStatus.cpuUsage }}</p>
              <p>运行时间: {{ systemStatus.uptime }}</p>
            </div>
            <div class="home-card">
              <Database :size="48" />
              <h3>许可协议</h3>
              <p>本软件采用 MIT 许可证开源。</p>
            </div>
          </div>
        </div>
        
        <!-- 项目内容 -->
        <div v-else-if="currentContent === 'projects'" class="content-section">
          <h1>项目管理</h1>
          <div class="projects-content">
            <div class="projects-sidebar">
              <h3>项目目录</h3>
              <div class="directory-tree">
                <div class="tree-item">项目 1</div>
                <div class="tree-item">项目 2</div>
                <div class="tree-item">项目 3</div>
              </div>
            </div>
            <div class="projects-main">
              <h3>WebShell 列表</h3>
              <div class="webshell-list">
                <div class="webshell-item">
                  <span>webshell1.php</span>
                  <span>192.168.1.100</span>
                  <span>2024-01-01</span>
                  <button class="action-button">进入</button>
                </div>
                <div class="webshell-item">
                  <span>webshell2.asp</span>
                  <span>192.168.1.101</span>
                  <span>2024-01-02</span>
                  <button class="action-button">进入</button>
                </div>
                <div class="webshell-item">
                  <span>webshell3.aspx</span>
                  <span>192.168.1.102</span>
                  <span>2024-01-03</span>
                  <button class="action-button">进入</button>
                </div>
              </div>
              <div class="pagination">
                <button class="page-button">上一页</button>
                <button class="page-button active">1</button>
                <button class="page-button">2</button>
                <button class="page-button">3</button>
                <button class="page-button">下一页</button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 载荷内容 -->
        <div v-else-if="currentContent === 'payloads'" class="content-section">
          <h1>载荷生成</h1>
          <div class="payloads-content">
            <div class="payload-form">
              <div class="form-group">
                <label>脚本类型</label>
                <select>
                  <option>PHP</option>
                  <option>ASP</option>
                  <option>ASPX</option>
                  <option>JSP</option>
                </select>
              </div>
              <div class="form-group">
                <label>功能类型</label>
                <select>
                  <option>基础版</option>
                  <option>文件管理</option>
                  <option>命令执行</option>
                  <option>数据库管理</option>
                  <option>完整版</option>
                </select>
              </div>
              <div class="form-group">
                <label>连接密码</label>
                <input type="password" placeholder="请输入连接密码">
              </div>
              <div class="form-group">
                <label>混淆级别</label>
                <select>
                  <option>无混淆</option>
                  <option>轻度混淆</option>
                  <option>中度混淆</option>
                  <option>高度混淆</option>
                </select>
              </div>
              <div class="form-group">
                <label>文件名</label>
                <input type="text" placeholder="shell.php">
              </div>
              <button class="generate-button">生成载荷</button>
            </div>
          </div>
        </div>
        
        <!-- 插件内容 -->
        <div v-else-if="currentContent === 'plugins'" class="content-section">
          <h1>插件管理</h1>
          <div class="plugins-content">
            <div class="plugins-tabs">
              <button class="tab-button active">本地插件</button>
              <button class="tab-button">插件商店</button>
            </div>
            <div class="plugins-list">
              <div class="plugin-item">
                <h3>插件 1</h3>
                <p>这是一个本地插件，用于执行特定任务。</p>
                <button class="plugin-button">启用</button>
              </div>
              <div class="plugin-item">
                <h3>插件 2</h3>
                <p>这是另一个本地插件，提供额外功能。</p>
                <button class="plugin-button">启用</button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 设置内容 -->
        <div v-else-if="currentContent === 'settings'" class="content-section">
          <h1>设置</h1>
          <div class="settings-content">
            <div class="settings-group">
              <h3>外观设置</h3>
              <div class="setting-item">
                <label>主题</label>
                <select>
                  <option>浅色</option>
                  <option>深色</option>
                  <option>跟随系统</option>
                </select>
              </div>
            </div>
            <div class="settings-group">
              <h3>系统设置</h3>
              <div class="setting-item">
                <label>自动更新</label>
                <input type="checkbox">
              </div>
              <div class="setting-item">
                <label>检查更新</label>
                <button>立即检查</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 底部状态栏 -->
    <div class="status-bar">
      <div class="status-item">
        <span>内存: {{ systemStatus.memoryUsage }}</span>
      </div>
      <div class="status-item">
        <span>进程 ID: {{ systemStatus.processId }}</span>
      </div>
      <div class="status-item">
        <span>CPU: {{ systemStatus.cpuUsage }}</span>
      </div>
      <div class="status-item">
        <span>运行时间: {{ systemStatus.uptime }}</span>
      </div>
    </div>
  </div>
</template>

<style>
/* 全局样式 - 液态玻璃风格 */
:root {
  --bg-color: #ffffff;
  --text-color: #1e293b;
  --border-color: rgba(0, 0, 0, 0.1);
  --hover-color: rgba(0, 0, 0, 0.06);
  --active-color: #667eea;
  --sidebar-bg: rgba(255, 255, 255, 0.95);
  --content-bg: rgba(255, 255, 255, 0.98);
  --status-bar-bg: rgba(255, 255, 255, 0.98);
  --card-bg: rgba(255, 255, 255, 0.98);
  --glass-blur: blur(20px);
  --shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.08);
  --shadow-md: 0 4px 16px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 8px 32px rgba(0, 0, 0, 0.15);
  --glass-border: 1px solid rgba(255, 255, 255, 0.4);
  --glass-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.6), inset 0 -1px 0 rgba(255, 255, 255, 0.2);
  --glass-gradient: linear-gradient(135deg, rgba(255, 255, 255, 0.4) 0%, rgba(255, 255, 255, 0.1) 100%);
  --glass-highlight: rgba(255, 255, 255, 0.8);
}

.dark {
  --bg-color: #0f172a;
  --text-color: #e2e8f0;
  --border-color: rgba(255, 255, 255, 0.1);
  --hover-color: rgba(255, 255, 255, 0.06);
  --active-color: #818cf8;
  --sidebar-bg: rgba(30, 41, 59, 0.85);
  --content-bg: rgba(15, 23, 42, 0.9);
  --status-bar-bg: rgba(30, 41, 59, 0.95);
  --card-bg: rgba(30, 41, 59, 0.9);
  --shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.4);
  --shadow-md: 0 4px 16px rgba(0, 0, 0, 0.5);
  --shadow-lg: 0 8px 32px rgba(0, 0, 0, 0.6);
  --glass-border: 1px solid rgba(255, 255, 255, 0.12);
  --glass-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.15), inset 0 -1px 0 rgba(255, 255, 255, 0.05);
  --glass-gradient: linear-gradient(135deg, rgba(255, 255, 255, 0.08) 0%, rgba(255, 255, 255, 0.02) 100%);
  --glass-highlight: rgba(255, 255, 255, 0.15);
}

body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark body {
  background: linear-gradient(135deg, #1e3a8a 0%, #0f172a 100%);
}

/* 应用容器样式 */
.app-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--bg-color);
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  gap: 12px;
  padding: 12px;
  width: 100%;
  box-sizing: border-box;
  max-width: 100%;
}

/* 侧边栏样式 - 液态玻璃风格 */
.sidebar {
  width: 90px;
  background: var(--sidebar-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
  gap: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  position: relative;
  overflow: hidden;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.nav-item {
  width: 75px;
  height: 75px;
  min-width: 75px;
  min-height: 75px;
  max-width: 75px;
  max-height: 75px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  gap: 8px;
  position: relative;
  overflow: hidden;
  padding: 10px 6px;
  z-index: 1;
  box-sizing: border-box;
}

.nav-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255, 0.1) 0%, rgba(255, 255, 255, 0) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 0;
}

.nav-item:hover::before {
  opacity: 1;
}

.nav-item:hover {
  background: var(--hover-color);
  transform: translateY(-3px);
  box-shadow: var(--shadow-sm);
}

.nav-item.active {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.4), 0 1px 3px rgba(99, 102, 241, 0.2);
  transform: scale(1.05);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.nav-item.active::before {
  opacity: 1;
}

.nav-item span {
  font-size: 12px;
  text-align: center;
  font-weight: 500;
  letter-spacing: 0.3px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  z-index: 1;
}

.nav-item svg {
  z-index: 1;
  flex-shrink: 0;
}

/* 内容区域样式 - 液态玻璃风格 */
.content-area {
  flex: 1;
  background: var(--content-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 12px;
  overflow-y: auto;
  padding: 24px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  box-sizing: border-box;
  min-width: 0;
  position: relative;
}

.content-area::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.content-section {
  position: relative;
  z-index: 1;
}

.content-section h1 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 24px;
  color: var(--text-color);
  transition: color 0.4s ease;
}

/* 首页卡片样式 */
.home-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 32px;
}

.home-card {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  border-radius: 12px;
  padding: 24px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  z-index: 1;
}

.home-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.home-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
}

.home-card svg {
  margin-bottom: 16px;
  color: var(--active-color);
}

.home-card h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.home-card p {
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-color);
  opacity: 0.8;
  position: relative;
  z-index: 1;
}

/* 项目管理样式 */
.projects-content {
  display: flex;
  gap: 20px;
  height: calc(100% - 40px);
}

.projects-sidebar {
  width: 250px;
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  border-radius: 12px;
  padding: 20px;
  overflow-y: auto;
  flex-shrink: 0;
  position: relative;
  overflow: hidden;
}

.projects-sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.projects-sidebar h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.directory-tree {
  position: relative;
  z-index: 1;
}

.tree-item {
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 4px;
  color: var(--text-color);
}

.tree-item:hover {
  background: var(--hover-color);
}

.projects-main {
  flex: 1;
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  border-radius: 12px;
  padding: 20px;
  overflow-y: auto;
  position: relative;
  overflow: hidden;
}

.projects-main::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.projects-main h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.webshell-list {
  position: relative;
  z-index: 1;
}

.webshell-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.2s ease;
  border: 1px solid var(--border-color);
}

.dark .webshell-item {
  background: rgba(255, 255, 255, 0.05);
}

.webshell-item:hover {
  background: var(--hover-color);
  transform: translateX(4px);
}

.webshell-item span {
  flex: 1;
  margin-right: 16px;
  font-size: 14px;
  color: var(--text-color);
}

.action-button {
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-button:hover {
  background: #5a67d8;
  transform: translateY(-1px);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 32px;
  gap: 8px;
  position: relative;
  z-index: 1;
}

.page-button {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 8px 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-color);
}

.page-button:hover {
  background: var(--hover-color);
  transform: translateY(-1px);
}

.page-button.active {
  background: var(--active-color);
  color: white;
  border-color: var(--active-color);
}

/* 载荷生成样式 */
.payloads-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  position: relative;
  z-index: 1;
}

.payload-form {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 12px;
  padding: 32px;
  width: 100%;
  max-width: 500px;
  position: relative;
  overflow: hidden;
}

.payload-form::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.form-group {
  margin-bottom: 24px;
  position: relative;
  z-index: 1;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
  color: var(--text-color);
}

.form-group select,
.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.8);
  color: var(--text-color);
  font-size: 14px;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.dark .form-group select,
.dark .form-group input {
  background: rgba(255, 255, 255, 0.1);
}

.form-group select:hover,
.form-group input:hover {
  border-color: var(--active-color);
}

.form-group select:focus,
.form-group input:focus {
  outline: none;
  border-color: var(--active-color);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.generate-button {
  width: 100%;
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 14px 20px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  z-index: 1;
}

.generate-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.generate-button:active {
  transform: translateY(0);
}

/* 插件管理样式 */
.plugins-content {
  position: relative;
  z-index: 1;
}

.plugins-tabs {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 12px;
}

.tab-button {
  background: none;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-color);
  font-size: 14px;
  font-weight: 500;
}

.tab-button:hover {
  background: var(--hover-color);
}

.tab-button.active {
  background: var(--active-color);
  color: white;
}

.plugins-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.plugin-item {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.plugin-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.plugin-item:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
}

.plugin-item h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.plugin-item p {
  font-size: 14px;
  line-height: 1.4;
  margin-bottom: 16px;
  color: var(--text-color);
  opacity: 0.8;
  position: relative;
  z-index: 1;
}

.plugin-button {
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  z-index: 1;
}

.plugin-button:hover {
  background: #5a67d8;
  transform: translateY(-1px);
}

/* 设置样式 */
.settings-content {
  max-width: 600px;
  position: relative;
  z-index: 1;
}

.settings-group {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  position: relative;
  overflow: hidden;
}

.settings-group::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.settings-group h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
  position: relative;
  z-index: 1;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item label {
  font-size: 14px;
  color: var(--text-color);
}

.setting-item select {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 6px 12px;
  color: var(--text-color);
  font-size: 14px;
}

.setting-item input[type="checkbox"] {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.setting-item button {
  background: var(--active-color);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.setting-item button:hover {
  background: #5a67d8;
  transform: translateY(-1px);
}

/* 底部状态栏样式 */
.status-bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 20px;
  padding: 8px 20px;
  background: var(--status-bar-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border-top: 1px solid var(--border-color);
  font-size: 12px;
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  z-index: 1;
}

.status-bar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.status-item {
  display: flex;
  align-items: center;
  position: relative;
  z-index: 1;
}

.status-item span {
  opacity: 0.8;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

.dark ::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
  transition: background 0.3s ease;
}

.dark ::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 60px;
  }
  
  .nav-item {
    width: 50px;
    height: 50px;
    min-width: 50px;
    min-height: 50px;
    max-width: 50px;
    max-height: 50px;
    padding: 8px 4px;
  }
  
  .nav-item span {
    font-size: 10px;
  }
  
  .nav-item svg {
    width: 16px;
    height: 16px;
  }
  
  .main-content {
    padding: 8px;
    gap: 8px;
  }
  
  .content-area {
    padding: 16px;
  }
  
  .home-content {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .projects-content {
    flex-direction: column;
  }
  
  .projects-sidebar {
    width: 100%;
    max-height: 200px;
  }
  
  .status-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
    padding: 12px 20px;
  }
}
</style>

<style scoped>
.app-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--bg-color);
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  gap: 12px;
  padding: 12px;
  width: 100%;
  box-sizing: border-box;
  max-width: 100%;
}

/* 侧边栏样式 - 液态玻璃风格 */
.sidebar {
  width: 90px;
  background: var(--sidebar-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
  gap: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  position: relative;
  overflow: hidden;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
}

.nav-item {
  width: 75px;
  height: 75px;
  min-width: 75px;
  min-height: 75px;
  max-width: 75px;
  max-height: 75px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  gap: 8px;
  position: relative;
  overflow: hidden;
  padding: 10px 6px;
  z-index: 1;
  box-sizing: border-box;
}

.nav-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255, 0.1) 0%, rgba(255, 255, 255, 0) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 0;
}

.nav-item:hover::before {
  opacity: 1;
}

.nav-item:hover {
  background: var(--hover-color);
  transform: translateY(-3px);
  box-shadow: var(--shadow-sm);
}

.nav-item.active {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.4), 0 1px 3px rgba(99, 102, 241, 0.2);
  transform: scale(1.05);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.nav-item.active::before {
  opacity: 1;
}

.nav-item span {
  font-size: 12px;
  text-align: center;
  font-weight: 500;
  letter-spacing: 0.3px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  z-index: 1;
}

.nav-item svg {
  z-index: 1;
  flex-shrink: 0;
}

/* 内容区域样式 - 液态玻璃风格 */
.content-area {
  flex: 1;
  background: var(--content-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 12px;
  overflow-y: auto;
  padding: 24px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  box-sizing: border-box;
  min-width: 0;
  position: relative;
}

.content-area::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
  border-radius: 12px;
}

.content-section {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.content-section h1 {
  margin: 0 0 24px 0;
  font-size: 28px;
  font-weight: 600;
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-align: center;
  width: 100%;
  box-sizing: border-box;
}

/* 首页内容样式 - 液态玻璃风格 */
.home-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.home-card {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 16px;
  padding: 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  width: 100%;
  box-sizing: border-box;
}

.home-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
  opacity: 0;
  transition: opacity 0.4s ease;
}

.home-card:hover::before {
  opacity: 1;
}

.home-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-lg);
}

.home-card svg {
  z-index: 1;
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.home-card h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
  z-index: 1;
  position: relative;
}

.home-card p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color);
  opacity: 0.8;
  z-index: 1;
  position: relative;
}

/* 项目内容样式 - 液态玻璃风格 */
.projects-content {
  display: flex;
  gap: 20px;
  height: calc(100% - 60px);
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.projects-sidebar {
  width: 200px;
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.projects-sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
  border-radius: 12px;
}

.directory-tree {
  margin-top: 16px;
  position: relative;
  z-index: 1;
}

.tree-item {
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 8px;
  margin-bottom: 6px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  position: relative;
  z-index: 1;
}

.tree-item:hover {
  background: var(--hover-color);
  transform: translateX(4px);
  box-shadow: var(--shadow-sm);
}

.tree-item.active {
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  color: white;
  box-shadow: var(--shadow-md);
}

.projects-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.webshell-list {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 12px;
  margin-bottom: 16px;
  flex: 1;
  overflow: hidden;
  position: relative;
}

.webshell-list::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
  border-radius: 12px;
}

.webshell-item {
  display: flex;
  align-items: center;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border-color);
  gap: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 1;
}

.webshell-item:hover {
  background: var(--hover-color);
  transform: translateX(4px);
  box-shadow: var(--shadow-sm);
}

.webshell-item:last-child {
  border-bottom: none;
}

.webshell-item span {
  flex: 1;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: flex;
  align-items: center;
  position: relative;
  z-index: 1;
}

.action-button {
  padding: 10px 20px;
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 80px;
}

.action-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.pagination {
  display: flex;
  gap: 8px;
  justify-content: center;
  align-items: center;
  width: 100%;
  flex-wrap: wrap;
}

.page-button {
  padding: 10px 18px;
  border: var(--glass-border);
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 14px;
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 50px;
}

.page-button:hover {
  background: var(--hover-color);
  transform: translateY(-2px);
}

.page-button.active {
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
  box-shadow: var(--shadow-md);
}

/* 载荷内容样式 - 液态玻璃风格 */
.payloads-content {
  height: calc(100% - 60px);
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-y: auto;
  overflow-x: hidden;
}

.payload-form {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 16px;
  padding: 28px;
  max-width: 560px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
  overflow: hidden;
  position: relative;
}

.payload-form::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
  border-radius: 16px;
}

.form-group {
  margin-bottom: 20px;
  width: 100%;
  box-sizing: border-box;
  position: relative;
  z-index: 1;
}

.form-group label {
  display: block;
  margin-bottom: 10px;
  font-weight: 600;
  font-size: 14px;
  color: var(--text-color);
  width: 100%;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
  z-index: 1;
}

.form-group label::before {
  content: '';
  width: 4px;
  height: 16px;
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  border-radius: 2px;
  flex-shrink: 0;
  box-shadow: 0 0 8px rgba(102, 126, 234, 0.4);
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: var(--glass-border);
  border-radius: 12px;
  background: var(--content-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  color: var(--text-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-sm);
  font-size: 14px;
  box-sizing: border-box;
  position: relative;
  z-index: 1;
}

.form-group input::placeholder {
  color: var(--text-color);
  opacity: 0.5;
}

.form-group select {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23666' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  background-size: 16px;
  padding-right: 40px;
}

.form-group select:focus {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23667eea' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
}

.generate-button {
  padding: 16px 36px;
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  font-size: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-md);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 140px;
  max-width: 200px;
  margin: 0 auto;
  box-sizing: border-box;
  letter-spacing: 0.5px;
}

.generate-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.generate-button:active {
  transform: translateY(0);
}

/* 插件内容样式 - 液态玻璃风格 */
.plugins-content {
  height: calc(100% - 60px);
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.plugins-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  justify-content: center;
}

.tab-button {
  padding: 12px 24px;
  border: var(--glass-border);
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;
  font-size: 14px;
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 90px;
}

.tab-button:hover {
  background: var(--hover-color);
  transform: translateY(-2px);
}

.tab-button.active {
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
  box-shadow: var(--shadow-md);
}

.plugins-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  width: 100%;
  box-sizing: border-box;
}

.plugin-item {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 16px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.plugin-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
  border-radius: 16px;
}

.plugin-item:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.plugin-item h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.plugin-item p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color);
  opacity: 0.8;
  position: relative;
  z-index: 1;
}

.plugin-button {
  margin-top: 12px;
  padding: 12px 24px;
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-sm);
  width: fit-content;
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 90px;
}

.plugin-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

/* 设置内容样式 - 液态玻璃风格 */
.settings-content {
  height: calc(100% - 60px);
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.settings-group {
  background: var(--card-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  box-shadow: var(--shadow-md), var(--glass-shadow);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 20px;
  max-width: 540px;
  margin-left: auto;
  margin-right: auto;
  width: 100%;
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
}

.settings-group::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--glass-gradient);
  pointer-events: none;
  z-index: 0;
  border-radius: 16px;
}

.settings-group h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);
  position: relative;
  z-index: 1;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 1;
}

.setting-item:hover {
  background: var(--hover-color);
  margin: 0 -24px;
  padding: 16px 24px;
  box-shadow: var(--shadow-sm);
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item label {
  font-weight: 500;
  font-size: 14px;
  color: var(--text-color);
}

.setting-item select {
  padding: 8px 14px;
  border: var(--glass-border);
  border-radius: 8px;
  background: var(--content-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  color: var(--text-color);
  font-size: 14px;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.setting-item button {
  padding: 10px 20px;
  background: linear-gradient(135deg, var(--active-color) 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 90px;
}

.setting-item button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

/* 状态栏样式 - 液态玻璃风格 */
.status-bar {
  height: 28px;
  background: var(--status-bar-bg);
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  border-top: var(--glass-border);
  box-shadow: var(--shadow-sm), var(--glass-shadow);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 20px;
  gap: 24px;
  font-size: 12px;
  color: var(--text-color);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  width: 100%;
  box-sizing: border-box;
  flex-wrap: wrap;
}

.status-item {
  display: flex;
  align-items: center;
  font-weight: 500;
}

/* 深色主题样式 - 液态玻璃风格 */
.app-container.dark-theme {
  background-color: #0f172a;
  color: #e2e8f0;
}

.app-container.dark-theme .sidebar {
  background: rgba(30, 41, 59, 0.7);
}

.app-container.dark-theme .content-area {
  background: rgba(15, 23, 42, 0.5);
}

.app-container.dark-theme .status-bar {
  background: rgba(30, 41, 59, 0.8);
}

.app-container.dark-theme .home-card,
.app-container.dark-theme .projects-sidebar,
.app-container.dark-theme .webshell-list,
.app-container.dark-theme .payload-form,
.app-container.dark-theme .plugin-item,
.app-container.dark-theme .settings-group {
  background: rgba(30, 41, 59, 0.6);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .home-content {
    max-width: 100%;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  }
  
  .projects-content {
    max-width: 100%;
  }
  
  .plugins-content {
    max-width: 100%;
  }
  
  .payloads-content {
    max-width: 100%;
  }
  
  .settings-content {
    max-width: 100%;
  }
}

@media (max-width: 768px) {
  .main-content {
    padding: 8px;
    gap: 8px;
  }
  
  .sidebar {
    width: 70px;
    padding: 16px 0;
    gap: 8px;
  }
  
  .nav-item {
    width: 58px;
    height: 62px;
    padding: 8px 4px;
    gap: 6px;
  }
  
  .nav-item span {
    font-size: 10px;
  }
  
  .content-area {
    padding: 16px;
  }
  
  .status-bar {
    gap: 12px;
    padding: 0 12px;
    font-size: 11px;
  }
  
  .status-item {
    white-space: nowrap;
  }
  
  .home-content {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .projects-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .projects-sidebar {
    width: 100%;
    height: auto;
    max-height: 150px;
  }
  
  .webshell-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .webshell-item span {
    width: 100%;
  }
  
  .pagination {
    flex-wrap: wrap;
    gap: 6px;
  }
  
  .page-button {
    padding: 8px 12px;
    min-width: 40px;
  }
  
  .plugins-tabs {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .plugins-list {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .settings-group {
    max-width: 100%;
  }
  
  .payload-form {
    max-width: 100%;
  }
  
  .form-group input,
  .form-group select {
    padding: 12px 16px;
    font-size: 14px;
  }
  
  .form-group label {
    font-size: 13px;
  }
  
  .generate-button {
    padding: 14px 28px;
    font-size: 15px;
  }
}

@media (max-width: 480px) {
  .sidebar {
    width: 60px;
  }
  
  .nav-item {
    width: 50px;
    height: 55px;
    padding: 6px 3px;
    gap: 4px;
  }
  
  .nav-item span {
    font-size: 9px;
  }
  
  .content-section h1 {
    font-size: 22px;
  }
  
  .home-card {
    padding: 20px;
  }
  
  .action-button {
    padding: 8px 16px;
    min-width: 70px;
  }
  
  .generate-button {
    padding: 12px 24px;
    min-width: 100px;
  }
  
  .tab-button {
    padding: 10px 18px;
    min-width: 70px;
  }
  
  .plugin-button {
    padding: 10px 18px;
    min-width: 70px;
  }
  
  .payload-form {
    padding: 24px;
    max-width: 100%;
  }
  
  .form-group {
    margin-bottom: 20px;
  }
  
  .form-group input,
  .form-group select {
    padding: 12px 16px;
    font-size: 14px;
  }
  
  .form-group label {
    font-size: 13px;
    margin-bottom: 10px;
  }
  
  .generate-button {
    padding: 14px 28px;
    font-size: 15px;
    min-width: 120px;
  }
}
</style>
