package services

import (
	"database/sql"
	"fmt"
	"time"

	"fg-abyss/internal/domain/entity"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
)

// DatabaseService 数据库管理服务
type DatabaseService struct {
	connections map[int64]*sql.DB
}

// NewDatabaseService 创建数据库服务
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{
		connections: make(map[int64]*sql.DB),
	}
}

// Connect 连接到数据库
func (s *DatabaseService) Connect(conn *entity.DatabaseConnection) error {
	dsn := s.buildDSN(conn)
	
	db, err := sql.Open(string(conn.Type), dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	
	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	
	s.connections[conn.ID] = db
	return nil
}

// Disconnect 断开数据库连接
func (s *DatabaseService) Disconnect(connID int64) error {
	if db, ok := s.connections[connID]; ok {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		delete(s.connections, connID)
		return nil
	}
	return fmt.Errorf("connection not found: %d", connID)
}

// ExecuteQuery 执行 SQL 查询
func (s *DatabaseService) ExecuteQuery(connID int64, query entity.DatabaseQuery) (*entity.DatabaseQueryResult, error) {
	db, ok := s.connections[connID]
	if !ok {
		return nil, fmt.Errorf("connection not found: %d", connID)
	}
	
	startTime := time.Now()
	
	// 设置超时
	if query.Timeout > 0 {
		db.SetConnMaxLifetime(time.Duration(query.Timeout) * time.Second)
	}
	
	// 执行查询
	rows, err := db.Query(query.SQL)
	if err != nil {
		return &entity.DatabaseQueryResult{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	defer rows.Close()
	
	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	
	// 读取数据
	var results []map[string]interface{}
	count := 0
	
	for rows.Next() {
		if query.Limit > 0 && count >= query.Limit {
			break
		}
		
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(values))
		
		for i := range values {
			scanArgs[i] = &values[i]
		}
		
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			var value interface{}
			if values[i] == nil {
				value = nil
			} else {
				switch v := values[i].(type) {
				case []byte:
					value = string(v)
				default:
					value = v
				}
			}
			rowMap[col] = value
		}
		
		results = append(results, rowMap)
		count++
	}
	
	duration := time.Since(startTime)
	
	return &entity.DatabaseQueryResult{
		Columns:  columns,
		Rows:     results,
		Affected: int64(len(results)),
		Duration: duration.String(),
		Success:  true,
		Message:  "Query executed successfully",
	}, nil
}

// ExecuteUpdate 执行更新操作
func (s *DatabaseService) ExecuteUpdate(connID int64, sql string) (*entity.DatabaseQueryResult, error) {
	db, ok := s.connections[connID]
	if !ok {
		return nil, fmt.Errorf("connection not found: %d", connID)
	}
	
	startTime := time.Now()
	
	result, err := db.Exec(sql)
	if err != nil {
		return &entity.DatabaseQueryResult{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	
	affected, _ := result.RowsAffected()
	duration := time.Since(startTime)
	
	return &entity.DatabaseQueryResult{
		Affected: affected,
		Duration: duration.String(),
		Success:  true,
		Message:  "Update executed successfully",
	}, nil
}

// GetTables 获取所有表
func (s *DatabaseService) GetTables(connID int64, connType entity.DatabaseType) ([]entity.DatabaseTable, error) {
	db, ok := s.connections[connID]
	if !ok {
		return nil, fmt.Errorf("connection not found: %d", connID)
	}
	
	var query string
	switch connType {
	case entity.DatabaseMySQL:
		query = `SELECT 
			TABLE_NAME as name,
			TABLE_SCHEMA as schema,
			TABLE_ROWS as rows,
			DATA_LENGTH as size,
			ENGINE as engine,
			TABLE_COLLATION as collation,
			TABLE_COMMENT as comment,
			CREATE_TIME as created_at,
			UPDATE_TIME as updated_at
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = DATABASE()`
	case entity.DatabasePostgreSQL:
		query = `SELECT 
			tablename as name,
			schemaname as schema,
			0 as rows,
			0 as size,
			'' as engine,
			'' as collation,
			'' as comment,
			'' as created_at,
			'' as updated_at
		FROM pg_tables
		WHERE schemaname NOT IN ('information_schema', 'pg_catalog')`
	case entity.DatabaseSQLite:
		query = `SELECT 
			name,
			'main' as schema,
			0 as rows,
			0 as size,
			'' as engine,
			'' as collation,
			'' as comment,
			'' as created_at,
			'' as updated_at
		FROM sqlite_master
		WHERE type='table'`
	default:
		return nil, fmt.Errorf("unsupported database type: %s", connType)
	}
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()
	
	var tables []entity.DatabaseTable
	for rows.Next() {
		var table entity.DatabaseTable
		if err := rows.Scan(
			&table.Name,
			&table.Schema,
			&table.Rows,
			&table.Size,
			&table.Engine,
			&table.Collation,
			&table.Comment,
			&table.CreatedAt,
			&table.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		tables = append(tables, table)
	}
	
	return tables, nil
}

// GetTableColumns 获取表列信息
func (s *DatabaseService) GetTableColumns(connID int64, connType entity.DatabaseType, tableName string) ([]entity.TableColumn, error) {
	db, ok := s.connections[connID]
	if !ok {
		return nil, fmt.Errorf("connection not found: %d", connID)
	}
	
	var query string
	switch connType {
	case entity.DatabaseMySQL:
		query = fmt.Sprintf(`SHOW COLUMNS FROM %s`, tableName)
	case entity.DatabasePostgreSQL:
		query = fmt.Sprintf(`SELECT 
			column_name,
			data_type,
			character_maximum_length,
			is_nullable,
			column_default,
			FALSE as auto_inc,
			FALSE as primary_key,
			'' as comment
		FROM information_schema.COLUMNS
		WHERE table_name = '%s'`, tableName)
	case entity.DatabaseSQLite:
		query = fmt.Sprintf(`PRAGMA table_info(%s)`, tableName)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", connType)
	}
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()
	
	var columns []entity.TableColumn
	for rows.Next() {
		var col entity.TableColumn
		switch connType {
		case entity.DatabaseMySQL:
			var extra sql.NullString
			err := rows.Scan(&col.Name, &col.Type, &col.Nullable, &col.PrimaryKey, &col.Default, &extra)
			if err != nil {
				return nil, err
			}
			col.AutoInc = extra.Valid && extra.String == "auto_increment"
		case entity.DatabasePostgreSQL:
			var length sql.NullInt64
			err := rows.Scan(&col.Name, &col.Type, &length, &col.Nullable, &col.Default, &col.AutoInc, &col.PrimaryKey, &col.Comment)
			if err != nil {
				return nil, err
			}
			if length.Valid {
				col.Length = fmt.Sprintf("%d", length.Int64)
			}
		case entity.DatabaseSQLite:
			var pk int
			err := rows.Scan(&col.Name, &col.Type, &col.Nullable, &col.Default, &pk)
			if err != nil {
				return nil, err
			}
			col.PrimaryKey = pk == 1
		}
		columns = append(columns, col)
	}
	
	return columns, nil
}

// buildDSN 构建数据库连接字符串
func (s *DatabaseService) buildDSN(conn *entity.DatabaseConnection) string {
	switch conn.Type {
	case entity.DatabaseMySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			conn.Username, conn.Password, conn.Host, conn.Port, conn.Database, conn.Charset)
	case entity.DatabasePostgreSQL:
		sslMode := "disable"
		if conn.SSLMode {
			sslMode = "require"
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			conn.Host, conn.Port, conn.Username, conn.Password, conn.Database, sslMode)
	case entity.DatabaseMSSQL:
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			conn.Username, conn.Password, conn.Host, conn.Port, conn.Database)
	case entity.DatabaseSQLite:
		return conn.Database
	default:
		return ""
	}
}

// TestConnection 测试数据库连接
func (s *DatabaseService) TestConnection(conn *entity.DatabaseConnection) (bool, string) {
	dsn := s.buildDSN(conn)
	
	db, err := sql.Open(string(conn.Type), dsn)
	if err != nil {
		return false, fmt.Sprintf("Failed to open: %v", err)
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		return false, fmt.Sprintf("Connection failed: %v", err)
	}
	
	return true, "Connection successful"
}

// GetConnection 获取数据库连接
func (s *DatabaseService) GetConnection(connID int64) (*sql.DB, bool) {
	db, ok := s.connections[connID]
	return db, ok
}

// GetAllConnections 获取所有连接
func (s *DatabaseService) GetAllConnections() map[int64]*sql.DB {
	return s.connections
}
