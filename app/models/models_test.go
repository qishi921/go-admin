package models

import (
	"testing"
)

func TestBaseModel_GetID(t *testing.T) {
	m := &BaseModel{ID: 123}
	if m.GetID() != 123 {
		t.Errorf("Expected ID 123, got %d", m.GetID())
	}
}

func TestBaseModel_SetID(t *testing.T) {
	m := &BaseModel{}
	m.SetID(456)
	if m.ID != 456 {
		t.Errorf("Expected ID 456, got %d", m.ID)
	}
}

func TestNowStr(t *testing.T) {
	s := NowStr()
	if len(s) != 19 {
		t.Errorf("Expected 19 chars, got %d: %s", len(s), s)
	}
	// 验证格式 YYYY-MM-DD HH:MM:SS
	if s[4] != '-' || s[7] != '-' || s[10] != ' ' || s[13] != ':' || s[16] != ':' {
		t.Errorf("Invalid format: %s", s)
	}
}

func TestUser_TableName(t *testing.T) {
	u := &User{}
	if u.TableName() != "users" {
		t.Errorf("Expected 'users', got %s", u.TableName())
	}
}

func TestRole_TableName(t *testing.T) {
	r := &Role{}
	if r.TableName() != "roles" {
		t.Errorf("Expected 'roles', got %s", r.TableName())
	}
}

func TestUser_Fields(t *testing.T) {
	phone := "13800138000"
	avatar := "/uploads/avatar.png"
	realName := "测试用户"
	lastLogin := "2026-01-01 12:00:00"
	roleId := 1

	u := &User{
		BaseModel: BaseModel{ID: 1},
		Username:  "testuser",
		Password:  "hashedpassword",
		Email:     "test@example.com",
		Phone:     &phone,
		Avatar:    &avatar,
		RealName:  &realName,
		Status:    "active",
		LastLoginAt: &lastLogin,
		RoleId:    &roleId,
	}

	if u.ID != 1 {
		t.Errorf("Expected ID 1, got %d", u.ID)
	}
	if u.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got %s", u.Username)
	}
	if *u.Phone != phone {
		t.Errorf("Expected Phone %s, got %s", phone, *u.Phone)
	}
	if u.Status != "active" {
		t.Errorf("Expected Status 'active', got %s", u.Status)
	}
}

func TestRole_Fields(t *testing.T) {
	r := &Role{
		BaseModel:   BaseModel{ID: 1},
		Name:        "管理员",
		Code:        "admin",
		Description: "系统管理员角色",
		Status:      "active",
	}

	if r.ID != 1 {
		t.Errorf("Expected ID 1, got %d", r.ID)
	}
	if r.Name != "管理员" {
		t.Errorf("Expected Name '管理员', got %s", r.Name)
	}
	if r.Code != "admin" {
		t.Errorf("Expected Code 'admin', got %s", r.Code)
	}
}

func TestBaseModel_SoftDelete(t *testing.T) {
	deletedAt := "2026-01-01 00:00:00"
	m := &BaseModel{DeletedAt: &deletedAt}

	if m.DeletedAt == nil {
		t.Error("Expected DeletedAt to be set")
	}
	if *m.DeletedAt != deletedAt {
		t.Errorf("Expected DeletedAt %s, got %s", deletedAt, *m.DeletedAt)
	}

	// 未删除的情况
	m2 := &BaseModel{}
	if m2.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil for non-deleted record")
	}
}
