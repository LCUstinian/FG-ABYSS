<template>
  <div class="command-panel">
    <div class="command-panel-container">
      <!-- 命令模板列表 -->
      <div class="command-templates">
        <div class="templates-header">
          <h3 class="templates-title">⚡ 快速命令</h3>
          <n-button size="small" quaternary @click="handleRefresh">
            <template #icon>
              <span>🔄</span>
            </template>
          </n-button>
        </div>
        <n-scrollbar style="max-height: calc(100vh - 300px)">
          <div class="template-list">
            <div 
              v-for="cmd in commandTemplates" 
              :key="cmd.id"
              class="template-item"
              :class="{ active: selectedTemplate?.id === cmd.id }"
              @click="selectTemplate(cmd)"
            >
              <div class="template-icon">{{ cmd.icon }}</div>
              <div class="template-info">
                <div class="template-name">{{ cmd.name }}</div>
                <div class="template-desc">{{ cmd.description }}</div>
              </div>
              <div class="template-action">
                <n-button 
                  size="tiny" 
                  type="primary"
                  @click.stop="executeCommand(cmd)"
                >
                  执行
                </n-button>
              </div>
            </div>
          </div>
        </n-scrollbar>
      </div>

      <!-- 命令执行区域 -->
      <div class="command-execution">
        <div class="execution-header">
          <h3 class="execution-title">命令执行</h3>
          <n-space>
            <n-tag size="small" :type="connected ? 'success' : 'default'">
              {{ connected ? '已连接' : '未连接' }}
            </n-tag>
          </n-space>
        </div>
        
        <div class="execution-body">
          <!-- 命令输入 -->
          <div class="command-input-area">
            <n-input
              v-model:value="customCommand"
              type="textarea"
              placeholder="输入自定义命令..."
              :rows="3"
              style="font-family: monospace;"
            />
            <n-space justify="end" style="margin-top: 12px;">
              <n-button @click="handleClear">清空</n-button>
              <n-button type="primary" @click="handleExecuteCustom" :loading="executing">
                <template #icon>
                  <span>▶️</span>
                </template>
                执行命令
              </n-button>
            </n-space>
          </div>

          <!-- 执行结果 -->
          <div class="command-output">
            <div class="output-header">
              <span class="output-title">执行结果</span>
              <n-space>
                <n-button 
                  size="tiny" 
                  quaternary
                  @click="handleCopyOutput"
                  :disabled="!output"
                >
                  <template #icon>
                    <span>📋</span>
                  </template>
                  复制
                </n-button>
                <n-button 
                  size="tiny" 
                  quaternary
                  @click="handleClearOutput"
                >
                  <template #icon>
                    <span>🗑️</span>
                  </template>
                  清空
                </n-button>
              </n-space>
            </div>
            <div class="output-content">
              <n-scrollbar ref="outputScrollbar">
                <pre v-if="output" v-text="output"></pre>
                <div v-else class="output-placeholder">
                  <span style="font-size: 48px; opacity: 0.3;">📝</span>
                  <p>暂无执行结果</p>
                </div>
              </n-scrollbar>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NSpace, NButton, NTag, NInput, NScrollbar, useMessage } from 'naive-ui'

interface CommandTemplate {
  id: string
  name: string
  description: string
  command: string
  icon: string
  category: string
}

defineProps<{
  webshellId?: string
}>()

const message = useMessage()
const connected = ref(false)
const executing = ref(false)
const customCommand = ref('')
const output = ref('')
const selectedTemplate = ref<CommandTemplate | null>(null)

// Mock 命令模板
const commandTemplates = ref<CommandTemplate[]>([
  {
    id: '1',
    name: '系统信息',
    description: '查看服务器基本信息',
    command: 'uname -a && whoami && pwd && date',
    icon: '💻',
    category: 'system'
  },
  {
    id: '2',
    name: '网络配置',
    description: '查看网络接口信息',
    command: 'ifconfig || ip addr',
    icon: '🌐',
    category: 'network'
  },
  {
    id: '3',
    name: '进程列表',
    description: '查看运行中的进程',
    command: 'ps aux | head -20',
    icon: '📊',
    category: 'system'
  },
  {
    id: '4',
    name: '磁盘空间',
    description: '查看磁盘使用情况',
    command: 'df -h',
    icon: '💾',
    category: 'system'
  },
  {
    id: '5',
    name: '内存使用',
    description: '查看内存使用情况',
    command: 'free -m',
    icon: '🧠',
    category: 'system'
  },
  {
    id: '6',
    name: 'PHP 版本',
    description: '查看 PHP 版本信息',
    command: 'php -v',
    icon: '🐘',
    category: 'php'
  },
  {
    id: '7',
    name: 'MySQL 状态',
    description: '检查 MySQL 服务状态',
    command: 'systemctl status mysqld || service mysql status',
    icon: '🐬',
    category: 'database'
  },
  {
    id: '8',
    name: 'Apache 状态',
    description: '检查 Apache 服务状态',
    command: 'systemctl status apache2 || service httpd status',
    icon: '🔷',
    category: 'web'
  },
  {
    id: '9',
    name: 'Nginx 状态',
    description: '检查 Nginx 服务状态',
    command: 'systemctl status nginx || service nginx status',
    icon: '⚡',
    category: 'web'
  },
  {
    id: '10',
    name: '查看日志',
    description: '查看系统日志',
    command: 'tail -f /var/log/syslog',
    icon: '📜',
    category: 'logs'
  }
])

const selectTemplate = (template: CommandTemplate) => {
  selectedTemplate.value = template
  customCommand.value = template.command
}

const executeCommand = async (template: CommandTemplate) => {
  executing.value = true
  output.value = `正在执行：${template.name}...\n命令：${template.command}\n\n`
  
  // 模拟命令执行
  setTimeout(() => {
    const mockOutputs: Record<string, string> = {
      '1': 'Linux webserver 5.4.0-42-generic #46-Ubuntu SMP Fri Jul 10 00:24:02 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux\nwww-data\n/var/www/html\nMon Jan 15 10:30:45 UTC 2024',
      '2': 'eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500\n        inet 192.168.1.100  netmask 255.255.255.0  broadcast 192.168.1.255\n        inet6 fe80::1  prefixlen 64  scopeid 0x20<link>',
      '3': 'USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND\nroot         1  0.0  0.1  18172  3456 ?        Ss   Jan14   0:02 /sbin/init\nroot         2  0.0  0.0      0     0 ?        S    Jan14   0:00 [kthreadd]',
      '4': 'Filesystem      Size  Used Avail Use% Mounted on\n/dev/sda1        50G   20G   28G  42% /\n/dev/sda2       100G   50G   45G  53% /home',
      '5': '              total        used        free      shared  buff/cache   available\nMem:           7982        2341        3456         123        2185        5318\nSwap:          2048           0        2048',
      '6': 'PHP 7.4.3 (cli) (built: Aug 17 2020 14:29:48) ( NTS )\nCopyright (c) The PHP Group\nZend Engine v3.4.0, Copyright (c) Zend Technologies\n    with Zend OPcache v7.4.3, Copyright (c), by Zend Technologies',
      '7': '● mysqld.service - MySQL Community Server\n   Loaded: loaded (/lib/systemd/system/mysqld.service; enabled)\n   Active: active (running) since Mon 2024-01-15 10:00:00 UTC',
      '8': '● apache2.service - The Apache HTTP Server\n   Loaded: loaded (/lib/systemd/system/apache2.service; enabled)\n   Active: active (running) since Mon 2024-01-15 09:00:00 UTC',
      '9': '● nginx.service - A high performance web server and a reverse proxy server\n   Loaded: loaded (/lib/systemd/system/nginx.service; enabled)\n   Active: active (running) since Mon 2024-01-15 08:00:00 UTC',
      '10': 'Jan 15 10:30:45 webserver systemd[1]: Started Session 1 of user www-data.\nJan 15 10:30:46 webserver apache2[1234]: AH00094: Command line: \'/usr/sbin/apache2\'\nJan 15 10:30:47 webserver kernel: [UFW BLOCK] IN=eth0 OUT= MAC=00:00:00:00:00:00'
    }
    
    output.value += '执行结果:\n' + (mockOutputs[template.id] || '命令执行成功')
    output.value += '\n\n[命令执行完成]'
    executing.value = false
    message.success(`${template.name} 执行成功`)
  }, 1000)
}

const handleExecuteCustom = async () => {
  if (!customCommand.value.trim()) {
    message.warning('请输入命令')
    return
  }
  
  executing.value = true
  output.value = `正在执行自定义命令...\n命令：${customCommand.value}\n\n`
  
  // 模拟执行
  setTimeout(() => {
    output.value += '执行结果:\n命令执行成功\n\n[命令执行完成]'
    executing.value = false
    message.success('命令执行成功')
  }, 1000)
}

const handleRefresh = () => {
  message.success('刷新命令模板')
}

const handleClear = () => {
  customCommand.value = ''
  selectedTemplate.value = null
}

const handleClearOutput = () => {
  output.value = ''
}

const handleCopyOutput = () => {
  if (output.value) {
    navigator.clipboard.writeText(output.value)
    message.success('复制成功')
  }
}

// 模拟连接
connected.value = true
</script>

<style scoped>
.command-panel {
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: var(--card-bg);
}

.command-panel-container {
  display: grid;
  grid-template-columns: 350px 1fr;
  height: 100%;
  gap: 1px;
  background: var(--border-color);
}

/* 命令模板列表 */
.command-templates {
  background: var(--card-bg);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.templates-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}

.templates-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color);
  margin: 0;
}

.template-list {
  padding: 4px 0;
}

.template-item {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  gap: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  border-left: 3px solid transparent;
}

.template-item:hover {
  background: var(--hover-color);
}

.template-item.active {
  background: var(--active-color-suppl);
  border-left-color: var(--active-color);
}

.template-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.template-info {
  flex: 1;
  min-width: 0;
}

.template-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 4px;
}

.template-desc {
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.template-action {
  flex-shrink: 0;
}

/* 命令执行区域 */
.command-execution {
  background: var(--card-bg);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.execution-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}

.execution-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color);
  margin: 0;
}

.execution-body {
  flex: 1;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
}

.command-input-area {
  flex-shrink: 0;
}

.command-output {
  flex: 1;
  background: var(--content-bg);
  border-radius: var(--border-radius-md);
  border: 1px solid var(--border-color);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.output-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-color);
}

.output-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color);
}

.output-content {
  flex: 1;
  overflow: hidden;
  background: #1e1e1e;
}

.output-content pre {
  margin: 0;
  padding: 12px 16px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #d4d4d4;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.output-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-secondary);
  text-align: center;
}

/* 滚动条样式 */
.output-content ::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

.output-content ::-webkit-scrollbar-track {
  background: #1e1e1e;
}

.output-content ::-webkit-scrollbar-thumb {
  background: #444;
  border-radius: 5px;
}

.output-content ::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
