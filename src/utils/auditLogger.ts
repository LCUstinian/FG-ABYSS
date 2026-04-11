import { invoke } from '@/utils/tauri-mock-adapter';

/**
 * 审计日志操作类型
 */
export enum AuditActionType {
  PayloadGenerate = 'PayloadGenerate',
  WebShellConnect = 'WebShellConnect',
  WebShellCommand = 'WebShellCommand',
  PluginLoad = 'PluginLoad',
  ProjectCreate = 'ProjectCreate',
  ProjectDelete = 'ProjectDelete',
  WebShellCreate = 'WebShellCreate',
  WebShellDelete = 'WebShellDelete',
  SettingsChange = 'SettingsChange',
  Other = 'Other'
}

/**
 * 审计日志记录器
 */
export class AuditLogger {
  /**
   * 记录审计日志
   * @param actionType 操作类型
   * @param content 操作内容
   * @param payloadHash 载荷哈希（可选）
   */
  static async log(
    actionType: AuditActionType,
    content: string,
    payloadHash?: string
  ): Promise<void> {
    try {
      await invoke('add_audit_log', {
        action_type: actionType,
        payload_hash: payloadHash,
        content
      });
    } catch (error) {
      console.error('Failed to log audit event:', error);
    }
  }

  /**
   * 记录载荷生成操作
   * @param payloadType 载荷类型
   * @param scriptType 脚本类型
   * @param payloadHash 载荷哈希
   */
  static async logPayloadGenerate(
    payloadType: string,
    scriptType: string,
    payloadHash: string
  ): Promise<void> {
    await this.log(
      AuditActionType.PayloadGenerate,
      `生成${scriptType}类型的${payloadType}载荷`,
      payloadHash
    );
  }

  /**
   * 记录WebShell连接操作
   * @param url WebShell URL
   * @param status 连接状态
   */
  static async logWebShellConnect(
    url: string,
    status: string
  ): Promise<void> {
    await this.log(
      AuditActionType.WebShellConnect,
      `连接WebShell: ${url}, 状态: ${status}`
    );
  }

  /**
   * 记录WebShell命令执行操作
   * @param url WebShell URL
   * @param command 执行的命令
   */
  static async logWebShellCommand(
    url: string,
    command: string
  ): Promise<void> {
    await this.log(
      AuditActionType.WebShellCommand,
      `在WebShell ${url} 执行命令: ${command}`
    );
  }

  /**
   * 记录插件加载操作
   * @param pluginName 插件名称
   * @param version 插件版本
   */
  static async logPluginLoad(
    pluginName: string,
    version: string
  ): Promise<void> {
    await this.log(
      AuditActionType.PluginLoad,
      `加载插件: ${pluginName} v${version}`
    );
  }

  /**
   * 记录项目创建操作
   * @param projectName 项目名称
   */
  static async logProjectCreate(
    projectName: string
  ): Promise<void> {
    await this.log(
      AuditActionType.ProjectCreate,
      `创建项目: ${projectName}`
    );
  }

  /**
   * 记录项目删除操作
   * @param projectName 项目名称
   */
  static async logProjectDelete(
    projectName: string
  ): Promise<void> {
    await this.log(
      AuditActionType.ProjectDelete,
      `删除项目: ${projectName}`
    );
  }

  /**
   * 记录WebShell创建操作
   * @param projectName 项目名称
   * @param url WebShell URL
   */
  static async logWebShellCreate(
    projectName: string,
    url: string
  ): Promise<void> {
    await this.log(
      AuditActionType.WebShellCreate,
      `在项目 ${projectName} 中创建WebShell: ${url}`
    );
  }

  /**
   * 记录WebShell删除操作
   * @param projectName 项目名称
   * @param url WebShell URL
   */
  static async logWebShellDelete(
    projectName: string,
    url: string
  ): Promise<void> {
    await this.log(
      AuditActionType.WebShellDelete,
      `在项目 ${projectName} 中删除WebShell: ${url}`
    );
  }

  /**
   * 记录设置修改操作
   * @param settingName 设置名称
   * @param oldValue 旧值
   * @param newValue 新值
   */
  static async logSettingsChange(
    settingName: string,
    oldValue: string,
    newValue: string
  ): Promise<void> {
    await this.log(
      AuditActionType.SettingsChange,
      `修改设置 ${settingName}: ${oldValue} → ${newValue}`
    );
  }

  /**
   * 记录其他操作
   * @param content 操作内容
   */
  static async logOther(
    content: string
  ): Promise<void> {
    await this.log(
      AuditActionType.Other,
      content
    );
  }
}
