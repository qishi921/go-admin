package controllers

import (
	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
)

// DashboardController 仪表盘控制器
type DashboardController struct {
	DB *orm.DB
}

// Stats 返回仪表盘统计数据
func (dc *DashboardController) Stats(c *ghttp.Context) {
	stats := make(map[string]any)

	// 用户统计
	var userTotal, userActive int
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&userTotal)
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM users WHERE status = 'active' AND deleted_at IS NULL`).Scan(&userActive)
	stats["users"] = map[string]int{
		"total":  userTotal,
		"active": userActive,
	}

	// 角色统计
	var roleTotal int
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL`).Scan(&roleTotal)
	stats["roles"] = map[string]int{
		"total": roleTotal,
	}

	// 菜单统计
	var menuTotal int
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM menus WHERE deleted_at IS NULL`).Scan(&menuTotal)
	stats["menus"] = map[string]int{
		"total": menuTotal,
	}

	// 权限统计
	var permissionTotal int
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM permissions WHERE deleted_at IS NULL`).Scan(&permissionTotal)
	stats["permissions"] = map[string]int{
		"total": permissionTotal,
	}

	// 操作日志统计（今日）
	var logToday int
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM operation_logs WHERE DATE(created_at) = DATE('now')`).Scan(&logToday)
	stats["logs_today"] = logToday

	// 通知统计
	var notificationTotal, notificationUnread int
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM notifications WHERE deleted_at IS NULL`).Scan(&notificationTotal)
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM notifications WHERE is_read = 0 AND deleted_at IS NULL`).Scan(&notificationUnread)
	stats["notifications"] = map[string]int{
		"total":  notificationTotal,
		"unread": notificationUnread,
	}

	// 定时任务统计
	var taskTotal, taskEnabled int
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM scheduled_tasks WHERE deleted_at IS NULL`).Scan(&taskTotal)
	dc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM scheduled_tasks WHERE status = 'enabled' AND deleted_at IS NULL`).Scan(&taskEnabled)
	stats["tasks"] = map[string]int{
		"total":   taskTotal,
		"enabled": taskEnabled,
	}

	c.Success(stats)
}

// RecentLogs 返回最近操作日志
func (dc *DashboardController) RecentLogs(c *ghttp.Context) {
	rows, err := dc.DB.SQL.Query(`
		SELECT id, username, action, method, path, ip, status, created_at
		FROM operation_logs
		ORDER BY created_at DESC
		LIMIT 10
	`)
	if err != nil {
		c.Error(500, "查询失败")
		return
	}
	defer rows.Close()

	var items []map[string]any
	for rows.Next() {
		var id int
		var username, action, method, path, ip, status, createdAt string
		rows.Scan(&id, &username, &action, &method, &path, &ip, &status, &createdAt)
		items = append(items, map[string]any{
			"id":         id,
			"username":   username,
			"action":     action,
			"method":     method,
			"path":       path,
			"ip":         ip,
			"status":     status,
			"created_at": createdAt,
		})
	}

	c.Success(map[string]any{
		"items": items,
	})
}

// SystemInfo 返回系统信息
func (dc *DashboardController) SystemInfo(c *ghttp.Context) {
	info := map[string]any{
		"version":     "1.0.0",
		"framework":   "Gai Framework",
		"go_version":  "1.22+",
		"database":    "SQLite",
		"environment": "development",
	}

	c.Success(info)
}