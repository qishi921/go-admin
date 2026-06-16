package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	ghttp "github.com/Hlgxz/gai/http"
	"github.com/user/admin-system/app/models"
)

// AuditLogger 审计日志记录器
type AuditLogger struct {
	DB        *orm.DB
	ExcludePaths []string // 排除的路径
	SensitiveFields []string // 敏感字段（需要脱敏）
}

// NewAuditLogger 创建审计日志记录器
func NewAuditLogger(db *orm.DB) *AuditLogger {
	return &AuditLogger{
		DB: db,
		ExcludePaths: []string{
			"/api/operation-logs",
			"/api/notifications/unread-count",
			"/api/notifications",
			"/static/",
			"/favicon.ico",
		},
		SensitiveFields: []string{
			"password",
			"password_confirmation",
			"old_password",
			"new_password",
			"token",
			"secret",
			"api_key",
		},
	}
}

// AuditMiddleware 审计日志中间件
func (al *AuditLogger) AuditMiddleware() ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		// 检查是否排除的路径
		path := c.Request.URL.Path
		for _, excludePath := range al.ExcludePaths {
			if strings.HasPrefix(path, excludePath) {
				c.Next()
				return
			}
		}

		// 只记录写操作（POST, PUT, DELETE, PATCH）
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" {
			c.Next()
			return
		}

		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 使用 response writer wrapper 捕获响应
		wrapper := &auditResponseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = wrapper

		c.Next()

		// 计算耗时
		duration := time.Since(start).Milliseconds()

		// 获取用户信息
		var userID *int
		var username string
		if claims, ok := c.Get("auth_claims"); ok {
			if jwtClaims, ok := claims.(*auth.Claims); ok {
				uid := int(jwtClaims.UserID)
				userID = &uid
				if uname, ok := jwtClaims.Extra["username"].(string); ok {
					username = uname
				}
			}
		}

		// 确定操作类型
		action := al.determineAction(method, path)

		// 脱敏请求参数
		params := al.sanitizeData(string(requestBody))

		// 截取响应结果（限制大小）
		result := wrapper.body.String()
		if len(result) > 2000 {
			result = result[:2000] + "..."
		}

		// 确定状态
		status := "success"
		if wrapper.status >= 400 {
			status = "failed"
		}

		// 创建日志记录
		log := &models.OperationLog{
			UserId:     userID,
			Username:   username,
			Action:     action,
			Method:     method,
			Path:       path,
			Ip:         al.getClientIP(c),
			UserAgent:  c.Request.UserAgent(),
			Params:     params,
			Result:     result,
			Duration:   int(duration),
			Status:     status,
			CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		}

		// 异步保存日志
		go al.saveLog(log)
	}
}

// DataChangeMiddleware 数据变更快照中间件
// 用于记录数据修改前后的快照
func (al *AuditLogger) DataChangeMiddleware(tableName string, getRecordFunc func(id uint64) (map[string]any, error)) ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		// 只处理 PUT 和 DELETE（更新和删除需要记录旧数据）
		var oldData map[string]any
		var recordID uint64

		if method == "PUT" || method == "DELETE" {
			// 从路径中提取 ID
			idStr := c.Param("id")
			if idStr != "" {
				// 解析 ID
				for _, ch := range idStr {
					if ch >= '0' && ch <= '9' {
						recordID = recordID*10 + uint64(ch-'0')
					}
				}
				if recordID > 0 {
					oldData, _ = getRecordFunc(recordID)
				}
			}
		}

		c.Next()

		// 只处理成功的请求
		if c.Writer.(*auditResponseWriter) == nil {
			return
		}
		wrapper := c.Writer.(*auditResponseWriter)
		if wrapper.status >= 400 {
			return
		}

		// 确定变更类型
		changeType := ""
		switch method {
		case "POST":
			changeType = "create"
		case "PUT", "PATCH":
			changeType = "update"
		case "DELETE":
			changeType = "delete"
		}

		if changeType == "" {
			return
		}

		// 获取新数据（对于创建和更新）
		var newData map[string]any
		if changeType == "create" || changeType == "update" {
			// 从响应中提取数据
			var response struct {
				Data map[string]any `json:"data"`
			}
			if err := json.Unmarshal(wrapper.body.Bytes(), &response); err == nil {
				newData = response.Data
			}
		}

		// 序列化数据
		var oldDataStr, newDataStr string
		if oldData != nil {
			b, _ := json.Marshal(oldData)
			oldDataStr = string(b)
		}
		if newData != nil {
			b, _ := json.Marshal(newData)
			newDataStr = string(b)
		}

		// 获取用户信息
		var userID *int
		var username string
		if claims, ok := c.Get("auth_claims"); ok {
			if jwtClaims, ok := claims.(*auth.Claims); ok {
				uid := int(jwtClaims.UserID)
				userID = &uid
				if uname, ok := jwtClaims.Extra["username"].(string); ok {
					username = uname
				}
			}
		}

		// 创建带快照的日志
		log := &models.OperationLog{
			UserId:        userID,
			Username:      username,
			Action:        al.determineAction(method, path),
			Method:        method,
			Path:          path,
			Ip:            al.getClientIP(c),
			UserAgent:     c.Request.UserAgent(),
			Duration:      0,
			Status:        "success",
			ResourceTable: tableName,
			RecordId:      &recordID,
			OldData:       oldDataStr,
			NewData:       newDataStr,
			ChangeType:    changeType,
			CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		}

		go al.saveLog(log)
	}
}

// determineAction 根据方法和路径确定操作名称
func (al *AuditLogger) determineAction(method, path string) string {
	// 从路径推断资源名称
	parts := strings.Split(strings.Trim(path, "/"), "/")
	resource := ""
	action := ""

	if len(parts) >= 2 {
		// 获取资源名称（如 users, roles 等）
		for i, part := range parts {
			if part == "api" && i+1 < len(parts) {
				resource = parts[i+1]
				break
			}
		}
	}

	// 根据方法确定动作
	switch method {
	case "GET":
		action = "查看"
	case "POST":
		action = "创建"
	case "PUT", "PATCH":
		action = "更新"
	case "DELETE":
		action = "删除"
	}

	if resource != "" {
		return action + resource
	}
	return action
}

// getClientIP 获取客户端 IP
func (al *AuditLogger) getClientIP(c *ghttp.Context) string {
	// 尝试从代理头获取
	if ip := c.Header("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	if ip := c.Header("X-Real-IP"); ip != "" {
		return ip
	}
	// 使用 RemoteAddr
	addr := c.Request.RemoteAddr
	if colonIdx := strings.LastIndex(addr, ":"); colonIdx != -1 {
		return addr[:colonIdx]
	}
	return addr
}

// sanitizeData 脱敏敏感数据
func (al *AuditLogger) sanitizeData(data string) string {
	if data == "" {
		return ""
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return data
	}

	for _, field := range al.SensitiveFields {
		if _, exists := obj[field]; exists {
			obj[field] = "******"
		}
	}

	b, _ := json.Marshal(obj)
	return string(b)
}

// saveLog 保存日志
func (al *AuditLogger) saveLog(log *models.OperationLog) {
	// 使用原生 SQL 插入，避免 ORM 的软删除等问题
	query := `INSERT INTO operation_logs
		(user_id, username, action, method, path, ip, user_agent, params, result, duration, status,
		 resource_table, record_id, old_data, new_data, change_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	al.DB.SQL.Exec(query,
		log.UserId, log.Username, log.Action, log.Method, log.Path, log.Ip, log.UserAgent,
		log.Params, log.Result, log.Duration, log.Status,
		log.ResourceTable, log.RecordId, log.OldData, log.NewData, log.ChangeType,
		log.CreatedAt, log.CreatedAt,
	)
}

// auditResponseWriter 响应包装器
type auditResponseWriter struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *auditResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *auditResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
