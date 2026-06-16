package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260611000001_seed_initial_data",
		Up: func(drv driver.Driver) string {
			return `INSERT OR IGNORE INTO roles (name, code, description, status, created_at, updated_at) VALUES
('超级管理员', 'super_admin', '拥有系统所有权限', 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('管理员', 'admin', '拥有大部分管理权限', 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('编辑者', 'editor', '可编辑内容，无法管理用户', 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('查看者', 'viewer', '只读权限', 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00');

INSERT OR IGNORE INTO permissions (name, code, type, description, status, sort_order, created_at, updated_at) VALUES
('用户管理', 'user:manage', 'menu', '用户增删改查', 'active', 1, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('用户创建', 'user:create', 'action', '创建新用户', 'active', 2, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('用户编辑', 'user:edit', 'action', '编辑用户信息', 'active', 3, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('用户删除', 'user:delete', 'action', '删除用户', 'active', 4, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('角色管理', 'role:manage', 'menu', '角色增删改查', 'active', 5, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('角色创建', 'role:create', 'action', '创建新角色', 'active', 6, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('角色编辑', 'role:edit', 'action', '编辑角色信息', 'active', 7, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('角色删除', 'role:delete', 'action', '删除角色', 'active', 8, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('菜单管理', 'menu:manage', 'menu', '菜单增删改查', 'active', 9, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('权限管理', 'permission:manage', 'menu', '权限增删改查', 'active', 10, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('日志查看', 'log:view', 'menu', '查看操作日志', 'active', 11, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('数据看板', 'dashboard:view', 'menu', '查看数据看板', 'active', 12, '2026-06-11 07:30:00', '2026-06-11 07:30:00');

INSERT OR IGNORE INTO menus (name, path, icon, component, sort_order, status, created_at, updated_at) VALUES
('数据看板', '/dashboard', 'layui-icon-chart', '', 1, 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('系统管理', '/system', 'layui-icon-set', '', 2, 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('用户管理', '/system/users', 'layui-icon-username', '', 3, 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('角色管理', '/system/roles', 'layui-icon-group', '', 4, 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('菜单管理', '/system/menus', 'layui-icon-menu-fill', '', 5, 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('权限管理', '/system/permissions', 'layui-icon-auz', '', 6, 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('日志管理', '/logs', 'layui-icon-log', '', 7, 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
('操作日志', '/logs/operations', 'layui-icon-log', '', 8, 'active', '2026-06-11 07:30:00', '2026-06-11 07:30:00');

INSERT OR IGNORE INTO role_user (role_id, user_id, created_at, updated_at) VALUES (1, 1, '2026-06-11 07:30:00', '2026-06-11 07:30:00');

INSERT OR IGNORE INTO role_permission (role_id, permission_id, created_at, updated_at) VALUES
(1, 1, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 2, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 3, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 4, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 5, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 6, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 7, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 8, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 9, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 10, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 11, '2026-06-11 07:30:00', '2026-06-11 07:30:00'),
(1, 12, '2026-06-11 07:30:00', '2026-06-11 07:30:00');`
		},
		Down: func(drv driver.Driver) string {
			return `DELETE FROM role_permission WHERE role_id = 1;
DELETE FROM role_user WHERE role_id = 1;
DELETE FROM menus;
DELETE FROM permissions;
DELETE FROM roles WHERE code IN ('super_admin', 'admin', 'editor', 'viewer');`
		},
	})
}