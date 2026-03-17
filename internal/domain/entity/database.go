package entity

import "time"

// DatabaseType 数据库类型
type DatabaseType string

const (
	DatabaseMySQL    DatabaseType = "mysql"
	DatabasePostgreSQL DatabaseType = "postgresql"
	DatabaseMSSQL    DatabaseType = "mssql"
	DatabaseOracle   DatabaseType = "oracle"
	DatabaseSQLite   DatabaseType = "sqlite"
)

// DatabaseConnection 数据库连接配置
type DatabaseConnection struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Type      DatabaseType   `json:"type" gorm:"size:20;not null"`
	Host      string         `json:"host" gorm:"size:255;not null"`
	Port      int            `json:"port" gorm:"not null"`
	Username  string         `json:"username" gorm:"size:100;not null"`
	Password  string         `json:"password" gorm:"size:255;not null"`
	Database  string         `json:"database" gorm:"size:255;not null"`
	Charset   string         `json:"charset" gorm:"size:50;default:utf8"`
	SSLMode   bool           `json:"ssl_mode" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// DatabaseTable 数据库表信息
type DatabaseTable struct {
	Name        string `json:"name"`
	Schema      string `json:"schema"`
	Rows        int64  `json:"rows"`
	Size        int64  `json:"size"`
	Engine      string `json:"engine"`
	Collation   string `json:"collation"`
	Comment     string `json:"comment"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// TableColumn 数据库列信息
type TableColumn struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Length     string `json:"length"`
	Nullable   bool   `json:"nullable"`
	Default    string `json:"default"`
	AutoInc    bool   `json:"auto_inc"`
	PrimaryKey bool   `json:"primary_key"`
	Comment    string `json:"comment"`
}

// DatabaseQuery 数据库查询请求
type DatabaseQuery struct {
	SQL     string `json:"sql"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Timeout int    `json:"timeout"`
}

// DatabaseQueryResult 数据库查询结果
type DatabaseQueryResult struct {
	Columns    []string          `json:"columns"`
	Rows       []map[string]interface{} `json:"rows"`
	Affected   int64             `json:"affected"`
	Duration   string            `json:"duration"`
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
}

// DatabaseOperation 数据库操作类型
type DatabaseOperation string

const (
	OperationSelect  DatabaseOperation = "select"
	OperationInsert  DatabaseOperation = "insert"
	OperationUpdate  DatabaseOperation = "update"
	OperationDelete  DatabaseOperation = "delete"
	OperationCreate  DatabaseOperation = "create"
	OperationDrop    DatabaseOperation = "drop"
	OperationAlter   DatabaseOperation = "alter"
)

// DatabaseRecord 数据库记录
type DatabaseRecord struct {
	ID        int64             `json:"id"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// DatabaseBackup 数据库备份信息
type DatabaseBackup struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	ConnID    int64     `json:"conn_id" gorm:"not null"`
	Filename  string    `json:"filename" gorm:"size:255;not null"`
	Size      int64     `json:"size" gorm:"not null"`
	Status    string    `json:"status" gorm:"size:50;default:pending"`
	Path      string    `json:"path" gorm:"size:500"`
	CreatedAt time.Time `json:"created_at"`
}
