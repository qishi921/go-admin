package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/user/admin-system/app/testutil"
)

func TestDashboardController_Stats(t *testing.T) {
	db := testutil.SetupTestDB(t)

	// 创建必要的表
	testutil.CreateTable(t, db, `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			status TEXT DEFAULT 'active',
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)
	testutil.CreateTable(t, db, `
		CREATE TABLE roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			slug TEXT NOT NULL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)
	testutil.CreateTable(t, db, `
		CREATE TABLE menus (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			path TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)
	testutil.CreateTable(t, db, `
		CREATE TABLE permissions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			slug TEXT NOT NULL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)
	testutil.CreateTable(t, db, `
		CREATE TABLE operation_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			action TEXT,
			method TEXT,
			path TEXT,
			ip TEXT,
			status TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`)
	testutil.CreateTable(t, db, `
		CREATE TABLE notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT,
			content TEXT,
			is_read INTEGER DEFAULT 0,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)
	testutil.CreateTable(t, db, `
		CREATE TABLE scheduled_tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			status TEXT DEFAULT 'enabled',
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)

	// 插入测试数据
	testutil.InsertTestData(t, db, "users",
		[]string{"username", "password", "status"},
		[][]any{
			{"admin", "hash1", "active"},
			{"user1", "hash2", "active"},
			{"user2", "hash3", "inactive"},
		},
	)

	dc := &DashboardController{DB: db}

	req := httptest.NewRequest("GET", "/api/v1/dashboard/stats", nil)
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)

	dc.Stats(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDashboardController_RecentLogs(t *testing.T) {
	db := testutil.SetupTestDB(t)

	testutil.CreateTable(t, db, `
		CREATE TABLE operation_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			action TEXT,
			method TEXT,
			path TEXT,
			ip TEXT,
			status TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`)

	testutil.InsertTestData(t, db, "operation_logs",
		[]string{"username", "action", "method", "path", "ip", "status"},
		[][]any{
			{"admin", "登录", "POST", "/api/v1/auth/login", "127.0.0.1", "success"},
			{"admin", "查看用户", "GET", "/api/v1/users", "127.0.0.1", "success"},
		},
	)

	dc := &DashboardController{DB: db}

	req := httptest.NewRequest("GET", "/api/v1/dashboard/recent-logs", nil)
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)

	dc.RecentLogs(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDashboardController_SystemInfo(t *testing.T) {
	db := testutil.SetupTestDB(t)
	dc := &DashboardController{DB: db}

	req := httptest.NewRequest("GET", "/api/v1/dashboard/system-info", nil)
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)

	dc.SystemInfo(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
