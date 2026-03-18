/**
 * 载荷生成模块类型定义
 */

// 生成模式
export type PayloadMode = 'simple' | 'advanced'

// 脚本类型
export type ScriptType = 'php' | 'jsp' | 'aspx' | 'asp'

// 功能类型
export type FunctionType = 'basic' | 'file_manager' | 'process_manager' | 'registry' | 'network'

// 编码器类型
export type EncodeType = 
  | 'none' 
  | 'base64' 
  | 'xor' 
  | 'gzinflate' 
  | 'hex' 
  | 'urlencode' 
  | 'rot13'

// 加密算法
export type EncryptAlgo = 'aes128_cbc' | 'aes256_cbc' | 'xor'

// 混淆级别
export type ObfuscationLevel = 'low' | 'medium' | 'high'

/**
 * 载荷生成配置
 */
export interface PayloadConfig {
  mode: PayloadMode
  script_type: ScriptType
  function_type: FunctionType
  password: string
  encode_type?: EncodeType
  encrypt_algo?: EncryptAlgo
  obfuscation_level: ObfuscationLevel
  output_filename?: string
  template_name?: string
}

/**
 * 客户端配置 (仅 Advanced 模式)
 */
export interface ClientConfig {
  key: string
  iv: string
  algorithm: string
  options: Record<string, any>
}

/**
 * 载荷生成结果
 */
export interface PayloadResult {
  code: string
  client_config: ClientConfig | null
  filename: string
  size: number
  success: boolean
  message?: string
}

/**
 * 编码器选项
 */
export const ENCODE_OPTIONS: Array<{ label: string; value: EncodeType }> = [
  { label: '无编码', value: 'none' },
  { label: 'Base64', value: 'base64' },
  { label: 'XOR', value: 'xor' },
  { label: 'GZInflate', value: 'gzinflate' },
  { label: 'Hex', value: 'hex' },
  { label: 'URL Encode', value: 'urlencode' },
  { label: 'ROT13', value: 'rot13' },
]

/**
 * 加密算法选项
 */
export const ENCRYPT_OPTIONS: Array<{ label: string; value: EncryptAlgo }> = [
  { label: 'AES-128-CBC', value: 'aes128_cbc' },
  { label: 'AES-256-CBC', value: 'aes256_cbc' },
  { label: 'XOR', value: 'xor' },
]

/**
 * 脚本类型选项
 */
export const SCRIPT_TYPE_OPTIONS: Array<{ label: string; value: ScriptType }> = [
  { label: 'PHP', value: 'php' },
  { label: 'JSP', value: 'jsp' },
  { label: 'ASPX', value: 'aspx' },
  { label: 'ASP', value: 'asp' },
]

/**
 * 功能类型选项
 */
export const FUNCTION_TYPE_OPTIONS: Array<{ label: string; value: FunctionType }> = [
  { label: '基础连接', value: 'basic' },
  { label: '文件管理', value: 'file_manager' },
  { label: '进程管理', value: 'process_manager' },
  { label: '注册表操作', value: 'registry' },
  { label: '网络操作', value: 'network' },
]

/**
 * 混淆级别选项
 */
export const OBFUSCATION_OPTIONS: Array<{ label: string; value: ObfuscationLevel }> = [
  { label: '低', value: 'low' },
  { label: '中', value: 'medium' },
  { label: '高', value: 'high' },
]
