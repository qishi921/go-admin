package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "storage/database.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	fmt.Printf("检查数据库: %s\n\n", dbPath)

	// 检查表
	tables := []string{
		"users", "roles", "menus", "permissions",
		"notifications", "notification_templates",
		"scheduled_tasks", "task_executions",
		"settings", "operation_logs",
		"role_user", "role_permission",
	}

	fmt.Println("=== 表结构检查 ===")
	for _, table := range tables {
		checkTable(db, table)
	}

	// 检查数据
	fmt.Println("\n=== 数据统计 ===")
	for _, table := range tables {
		countRows(db, table)
	}

	// 检查用户
	fmt.Println("\n=== 用户列表 ===")
	listUsers(db)
}

func checkTable(db *sql.DB, table string) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		fmt.Printf("❌ 表 %s 不存在\n", table)
		return
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue sql.NullString
		rows.Scan(&cid, &name, &ctype, &notnull, &pk, &dfltValue)
		columns = append(columns, name)
	}

	fmt.Printf("✓ 表 %s 存在，列: %s\n", table, strings.Join(columns, ", "))
}

func countRows(db *sql.DB, table string) {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted_at IS NULL", table)).Scan(&count)
	if err != nil {
		fmt.Printf("❌ 无法查询 %s\n", table)
		return
	}
	fmt.Printf("  %s: %d 条记录\n", table, count)
}

func listUsers(db *sql.DB) {
	rows, err := db.Query("SELECT id, username, email, status FROM users WHERE deleted_at IS NULL")
	if err != nil {
		fmt.Println("❌ 无法查询用户")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, email, status string
		rows.Scan(&id, &username, &email, &status)
		fmt.Printf("  [%d] %s (%s) - %s\n", id, username, email, status)
	}
}