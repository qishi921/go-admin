package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hlgxz/gai/auth"
	ghttp "github.com/Hlgxz/gai/http"
	"github.com/user/admin-system/app/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthController_Login(t *testing.T) {
	db := testutil.SetupTestDB(t)

	testutil.CreateTable(t, db, `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			avatar TEXT,
			status TEXT DEFAULT 'active',
			last_login_at TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)

	// 创建测试用户
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testutil.InsertTestData(t, db, "users",
		[]string{"username", "password", "email", "status"},
		[][]any{{"admin", string(hashedPassword), "admin@example.com", "active"}},
	)

	// 创建 auth manager
	authMgr := auth.NewManager("jwt")
	authMgr.RegisterGuard(auth.NewJWTGuard("test-secret-key-for-testing", 3600))

	ctrl := NewAuthController(db, authMgr)

	tests := []struct {
		name         string
		username     string
		password     string
		expectStatus int
	}{
		{"正确登录", "admin", "password123", http.StatusOK},
		{"错误密码", "admin", "wrongpassword", http.StatusUnauthorized},
		{"不存在用户", "nonexistent", "password123", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]any{
				"username": tt.username,
				"password": tt.password,
			}
			bodyBytes, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c := ghttp.NewContext(w, req)

			ctrl.Login(c)

			var resp map[string]any
			json.NewDecoder(w.Body).Decode(&resp)

			if tt.expectStatus == http.StatusOK {
				if w.Code != http.StatusOK {
					t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
				}
				if _, ok := resp["data"].(map[string]any); !ok {
					t.Error("Expected data in response")
				}
				data := resp["data"].(map[string]any)
				if _, ok := data["token"].(string); !ok {
					t.Error("Expected token in response")
				}
			} else {
				if code, ok := resp["code"].(float64); ok && int(code) != tt.expectStatus {
					t.Errorf("Expected code %d, got %d", tt.expectStatus, int(code))
				}
			}
		})
	}
}

func TestAuthController_Register(t *testing.T) {
	db := testutil.SetupTestDB(t)

	testutil.CreateTable(t, db, `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			email TEXT UNIQUE,
			phone TEXT,
			avatar TEXT,
			status TEXT DEFAULT 'active',
			last_login_at TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)

	authMgr := auth.NewManager("jwt")
	authMgr.RegisterGuard(auth.NewJWTGuard("test-secret-key", 3600))

	ctrl := NewAuthController(db, authMgr)

	t.Run("成功注册", func(t *testing.T) {
		body := map[string]any{
			"username": "newuser",
			"password": "password123",
			"email":    "newuser@example.com",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c := ghttp.NewContext(w, req)

		ctrl.Register(c)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		count := testutil.RowCount(t, db, "users")
		if count != 1 {
			t.Errorf("Expected 1 user, got %d", count)
		}
	})

	t.Run("重复用户名", func(t *testing.T) {
		body := map[string]any{
			"username": "newuser", // 已存在
			"password": "password123",
			"email":    "another@example.com",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c := ghttp.NewContext(w, req)

		ctrl.Register(c)

		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)

		if code, ok := resp["code"].(float64); ok && int(code) != http.StatusConflict {
			t.Errorf("Expected code %d, got %d", http.StatusConflict, int(code))
		}
	})
}

func TestAuthController_Me(t *testing.T) {
	db := testutil.SetupTestDB(t)

	testutil.CreateTable(t, db, `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			avatar TEXT,
			real_name TEXT,
			status TEXT DEFAULT 'active',
			last_login_at TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)

	testutil.InsertTestData(t, db, "users",
		[]string{"username", "password", "email", "status"},
		[][]any{{"testuser", "hash", "test@example.com", "active"}},
	)

	authMgr := auth.NewManager("jwt")
	authMgr.RegisterGuard(auth.NewJWTGuard("test-secret-key", 3600))

	ctrl := NewAuthController(db, authMgr)

	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)
	c.Set("auth_user_id", uint64(1))

	ctrl.Me(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.Fatal("Expected data in response")
	}
	if data["username"] != "testuser" {
		t.Errorf("Expected username 'testuser', got %v", data["username"])
	}
}

func TestAuthController_UpdatePassword(t *testing.T) {
	db := testutil.SetupTestDB(t)

	testutil.CreateTable(t, db, `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			email TEXT,
			status TEXT DEFAULT 'active',
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT
		)
	`)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	testutil.InsertTestData(t, db, "users",
		[]string{"username", "password", "email", "status"},
		[][]any{{"testuser", string(hashedPassword), "test@example.com", "active"}},
	)

	authMgr := auth.NewManager("jwt")
	authMgr.RegisterGuard(auth.NewJWTGuard("test-secret-key", 3600))

	ctrl := NewAuthController(db, authMgr)

	body := map[string]any{
		"old_password": "oldpassword",
		"new_password": "newpassword123",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/v1/auth/password", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := ghttp.NewContext(w, req)
	c.Set("auth_user_id", uint64(1))

	ctrl.UpdatePassword(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// 验证密码已更新
	var newPassword string
	db.SQL.QueryRow("SELECT password FROM users WHERE id = 1").Scan(&newPassword)
	if bcrypt.CompareHashAndPassword([]byte(newPassword), []byte("newpassword123")) != nil {
		t.Error("Password was not updated correctly")
	}
}
