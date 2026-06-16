package testutil

import (
	"database/sql"
	"os"
	"testing"

	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/orm"
)

// SetupTestDB 创建内存 SQLite 数据库用于测试
func SetupTestDB(t *testing.T) *orm.DB {
	t.Helper()

	drv, err := driver.Get("sqlite")
	if err != nil {
		t.Fatalf("Failed to get sqlite driver: %v", err)
	}

	// 使用内存数据库
	sqlDB, err := drv.Open(":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	db := &orm.DB{
		SQL:        sqlDB,
		DriverName: "sqlite",
		QuoteIdent: drv.QuoteIdent,
	}

	t.Cleanup(func() {
		sqlDB.Close()
	})

	return db
}

// CreateTable 创建测试表
func CreateTable(t *testing.T, db *orm.DB, schema string) {
	t.Helper()
	_, err := db.SQL.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
}

// InsertTestData 插入测试数据
func InsertTestData(t *testing.T, db *orm.DB, table string, columns []string, rows [][]any) {
	t.Helper()
	for _, row := range rows {
	 placeholders := make([]string, len(row))
	 args := make([]any, len(row))
	 for i, v := range row {
		 placeholders[i] = "?"
		 args[i] = v
	 }
	 query := "INSERT INTO " + table + " (" + columns[0]
	 for i := 1; i < len(columns); i++ {
		 query += ", " + columns[i]
	 }
	 query += ") VALUES (" + placeholders[0]
	 for i := 1; i < len(placeholders); i++ {
		 query += ", " + placeholders[i]
	 }
	 query += ")"
	 _, err := db.SQL.Exec(query, args...)
	 if err != nil {
		 t.Fatalf("Failed to insert test data: %v", err)
	 }
	}
}

// Setenv 设置环境变量并在测试后恢复
func Setenv(t *testing.T, key, value string) {
	t.Helper()
	old := os.Getenv(key)
	os.Setenv(key, value)
	t.Cleanup(func() {
		if old == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, old)
		}
	})
}

// RowCount 返回表中的行数
func RowCount(t *testing.T, db *orm.DB, table string) int {
	t.Helper()
	var count int
	err := db.SQL.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count rows: %v", err)
	}
	return count
}

// RowExists 检查行是否存在
func RowExists(t *testing.T, db *orm.DB, query string, args ...any) bool {
	t.Helper()
	var exists int
	err := db.SQL.QueryRow("SELECT 1 WHERE EXISTS ("+query+")", args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		t.Fatalf("Failed to check row existence: %v", err)
	}
	return true
}
