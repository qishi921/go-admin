package controllers

import (
	"database/sql"
	"strconv"
	"time"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
)

// NotificationController 通知控制器
type NotificationController struct {
	DB *orm.DB
}

// NotificationListResponse 通知列表响应
type NotificationListResponse struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Priority  int    `json:"priority"`
	IsRead    bool   `json:"is_read"`
	ReadAt    string `json:"read_at,omitempty"`
	Channel   string `json:"channel"`
	SentAt    string `json:"sent_at,omitempty"`
	SendStatus string `json:"send_status"`
	CreatedAt string `json:"created_at"`
}

// List 获取通知列表
func (nc *NotificationController) List(c *ghttp.Context) {
	userID := uint64(0)
	if claims, ok := c.Get("auth_claims"); ok {
		if jwtClaims, ok := claims.(*auth.Claims); ok {
			userID = jwtClaims.UserID
		}
	}
	if userID == 0 {
		c.Error(401, "请先登录")
		return
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 使用原始 SQL 查询
	rows, err := nc.DB.SQL.Query(`
		SELECT id, user_id, title, content, type, priority, is_read,
		       COALESCE(read_at, ''), channel, COALESCE(sent_at, ''),
		       send_status, created_at
		FROM notifications
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userID, pageSize, offset)
	if err != nil {
		c.Error(500, "查询失败: "+err.Error())
		return
	}
	defer rows.Close()

	var items []NotificationListResponse
	for rows.Next() {
		var item NotificationListResponse
		var readAt, sentAt sql.NullString
		var isRead int
		err := rows.Scan(
			&item.ID, &item.UserID, &item.Title, &item.Content, &item.Type,
			&item.Priority, &isRead, &readAt, &item.Channel, &sentAt,
			&item.SendStatus, &item.CreatedAt,
		)
		if err != nil {
			continue
		}
		item.IsRead = isRead == 1
		if readAt.Valid {
			item.ReadAt = readAt.String
		}
		if sentAt.Valid {
			item.SentAt = sentAt.String
		}
		items = append(items, item)
	}

	// 获取总数
	var total int
	nc.DB.SQL.QueryRow(`
		SELECT COUNT(*) FROM notifications
		WHERE user_id = ? AND deleted_at IS NULL
	`, userID).Scan(&total)

	// 获取未读数量
	var unreadCount int
	nc.DB.SQL.QueryRow(`
		SELECT COUNT(*) FROM notifications
		WHERE user_id = ? AND is_read = 0 AND deleted_at IS NULL
	`, userID).Scan(&unreadCount)

	c.Success(map[string]any{
		"items":        items,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
		"unread_count": unreadCount,
	})
}

// GetUnreadCount 获取未读通知数量
func (nc *NotificationController) GetUnreadCount(c *ghttp.Context) {
	userID := uint64(0)
	if claims, ok := c.Get("auth_claims"); ok {
		if jwtClaims, ok := claims.(*auth.Claims); ok {
			userID = jwtClaims.UserID
		}
	}
	if userID == 0 {
		c.Error(401, "请先登录")
		return
	}

	var count int
	nc.DB.SQL.QueryRow(`
		SELECT COUNT(*) FROM notifications
		WHERE user_id = ? AND is_read = 0 AND deleted_at IS NULL
	`, userID).Scan(&count)

	c.Success(map[string]int{"count": count})
}

// MarkRead 标记通知为已读
func (nc *NotificationController) MarkRead(c *ghttp.Context) {
	userID := uint64(0)
	if claims, ok := c.Get("auth_claims"); ok {
		if jwtClaims, ok := claims.(*auth.Claims); ok {
			userID = jwtClaims.UserID
		}
	}
	if userID == 0 {
		c.Error(401, "请先登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(400, "无效的通知ID")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := nc.DB.SQL.Exec(`
		UPDATE notifications
		SET is_read = 1, read_at = ?, updated_at = ?
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`, now, now, id, userID)
	if err != nil {
		c.Error(500, "操作失败")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.Error(404, "通知不存在")
		return
	}

	c.Success(map[string]string{"message": "已标记为已读"})
}

// MarkAllRead 标记所有通知为已读
func (nc *NotificationController) MarkAllRead(c *ghttp.Context) {
	userID := uint64(0)
	if claims, ok := c.Get("auth_claims"); ok {
		if jwtClaims, ok := claims.(*auth.Claims); ok {
			userID = jwtClaims.UserID
		}
	}
	if userID == 0 {
		c.Error(401, "请先登录")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	nc.DB.SQL.Exec(`
		UPDATE notifications
		SET is_read = 1, read_at = ?, updated_at = ?
		WHERE user_id = ? AND is_read = 0 AND deleted_at IS NULL
	`, now, now, userID)

	c.Success(map[string]string{"message": "已全部标记为已读"})
}

// Delete 删除通知
func (nc *NotificationController) Delete(c *ghttp.Context) {
	userID := uint64(0)
	if claims, ok := c.Get("auth_claims"); ok {
		if jwtClaims, ok := claims.(*auth.Claims); ok {
			userID = jwtClaims.UserID
		}
	}
	if userID == 0 {
		c.Error(401, "请先登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(400, "无效的通知ID")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := nc.DB.SQL.Exec(`
		UPDATE notifications
		SET deleted_at = ?
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`, now, id, userID)
	if err != nil {
		c.Error(500, "删除失败")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.Error(404, "通知不存在")
		return
	}

	c.Success(map[string]string{"message": "删除成功"})
}

// Create 创建通知（管理员）
func (nc *NotificationController) Create(c *ghttp.Context) {
	var input struct {
		UserID   uint64 `json:"user_id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Type     string `json:"type"`
		Priority int    `json:"priority"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.Error(400, "无效的请求数据")
		return
	}

	if input.UserID == 0 || input.Title == "" {
		c.Error(400, "用户ID和标题不能为空")
		return
	}

	if input.Type == "" {
		input.Type = "system"
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := nc.DB.SQL.Exec(`
		INSERT INTO notifications (user_id, title, content, type, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, input.UserID, input.Title, input.Content, input.Type, input.Priority, now, now)
	if err != nil {
		c.Error(500, "创建失败")
		return
	}

	id, _ := result.LastInsertId()

	c.Success(map[string]any{
		"id":      id,
		"message": "创建成功",
	})
}