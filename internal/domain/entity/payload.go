package entity

// PayloadType Payload 类型
type PayloadType string

const (
	// PayloadTypePHP PHP Payload
	PayloadTypePHP PayloadType = "php"
	// PayloadTypeASP ASP Payload
	PayloadTypeASP PayloadType = "asp"
	// PayloadTypeASPX ASPX Payload
	PayloadTypeASPX PayloadType = "aspx"
	// PayloadTypeJSP JSP Payload
	PayloadTypeJSP PayloadType = "jsp"
)

// PayloadFunction Payload 功能类型
type PayloadFunction string

const (
	// PayloadFunctionBasic 基础版本
	PayloadFunctionBasic PayloadFunction = "basic"
	// PayloadFunctionFile 文件管理
	PayloadFunctionFile PayloadFunction = "file"
	// PayloadFunctionCommand 命令执行
	PayloadFunctionCommand PayloadFunction = "command"
	// PayloadFunctionDatabase 数据库管理
	PayloadFunctionDatabase PayloadFunction = "database"
	// PayloadFunctionFull 完整版本
	PayloadFunctionFull PayloadFunction = "full"
)

// ObfuscationLevel 混淆级别
type ObfuscationLevel string

const (
	// ObfuscationNone 无混淆
	ObfuscationNone ObfuscationLevel = "none"
	// ObfuscationLow 轻度混淆
	ObfuscationLow ObfuscationLevel = "low"
	// ObfuscationMedium 中度混淆
	ObfuscationMedium ObfuscationLevel = "medium"
	// ObfuscationHigh 高度混淆
	ObfuscationHigh ObfuscationLevel = "high"
)

// PayloadConfig Payload 生成配置
type PayloadConfig struct {
	// Type Payload 类型（php/asp/aspx/jsp）
	Type PayloadType `json:"type"`
	// Function 功能类型
	Function PayloadFunction `json:"function"`
	// Password 连接密码
	Password string `json:"password"`
	// Encoder 编码器类型
	Encoder string `json:"encoder"`
	// EncryptionKey 加密密钥
	EncryptionKey string `json:"encryption_key"`
	// ObfuscationLevel 混淆级别
	ObfuscationLevel ObfuscationLevel `json:"obfuscation_level"`
	// CustomTemplate 自定义模板路径
	CustomTemplate string `json:"custom_template"`
	// OutputFilename 输出文件名
	OutputFilename string `json:"output_filename"`
	// Options 额外选项
	Options map[string]interface{} `json:"options,omitempty"`
}

// PayloadResult Payload 生成结果
type PayloadResult struct {
	// Success 是否成功
	Success bool `json:"success"`
	// Content Payload 内容
	Content string `json:"content"`
	// Filename 文件名
	Filename string `json:"filename"`
	// Size 文件大小（字节）
	Size int `json:"size"`
	// Message 消息
	Message string `json:"message"`
	// Warnings 警告信息
	Warnings []string `json:"warnings,omitempty"`
}

// PayloadTemplate Payload 模板
type PayloadTemplate struct {
	// Name 模板名称
	Name string `json:"name"`
	// Description 模板描述
	Description string `json:"description"`
	// Type Payload 类型
	Type PayloadType `json:"type"`
	// Function 功能类型
	Function PayloadFunction `json:"function"`
	// Content 模板内容
	Content string `json:"content"`
	// Parameters 模板参数
	Parameters []TemplateParameter `json:"parameters"`
	// IsBuiltin 是否为内置模板
	IsBuiltin bool `json:"is_builtin"`
}

// TemplateParameter 模板参数
type TemplateParameter struct {
	// Name 参数名称
	Name string `json:"name"`
	// Description 参数描述
	Description string `json:"description"`
	// Type 参数类型
	Type string `json:"type"`
	// Required 是否必需
	Required bool `json:"required"`
	// Default 默认值
	Default string `json:"default"`
}
