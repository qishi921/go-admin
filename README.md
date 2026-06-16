# 后台管理系统 (Admin System)

基于 Gai 框架构建的 Go 后台管理系统脚手架。

## 功能模块

### 核心功能
- [x] 用户管理 (User Management) - 用户 CRUD、密码修改、状态管理
- [x] 角色管理 (Role Management) - 角色 CRUD、状态管理
- [x] 权限管理 (Permission Management) - 权限 CRUD、权限分配
- [x] 菜单管理 (Menu Management) - 菜单树形结构、权限绑定
- [x] 操作日志 (Operation Logs) - 自动记录操作、日志查询
- [x] 数据看板 (Dashboard) - 统计数据、最近日志、系统信息

### P2 功能
- [x] 通知系统 (Notifications) - 通知创建、已读标记、未读统计
- [x] 定时任务 (Scheduled Tasks) - 任务调度、执行记录、手动触发
- [x] 数据导出 (Export) - CSV/JSON 导出、导入模板
- [x] 审计日志 (Audit Logs) - 操作审计、变更追踪
- [x] 系统设置 (Settings) - 分组配置、动态更新
- [x] 文件上传 (Upload) - 文件上传、图片处理

### 安全特性
- [x] JWT 认证 - Token 登录、自动续期
- [x] 登录限流 - 防暴力破解
- [x] API 限流 - 防滥用保护
- [x] XSS 防护 - 输入过滤
- [x] CSRF 防护 - 表单保护
- [x] RBAC 权限 - 角色权限控制

## 技术栈

- **后端**: Go 1.22+ + Gai Framework
- **数据库**: SQLite (开发) / MySQL / PostgreSQL (生产)
- **认证**: JWT
- **前端**: Layui (内置管理界面)

## 快速开始

```bash
# 安装依赖
go mod tidy

# 构建项目
go build -o bin/admin-server main.go

# 启动服务
./bin/admin-server

# 或直接运行
go run main.go
```

服务启动后访问：
- 后台界面：http://localhost:8080/
- API 文档：http://localhost:8080/api/
- 默认账号：admin / admin123

## 开发命令

```bash
# 运行测试
go test ./... -v

# 测试覆盖率
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 代码格式化
go fmt ./...

# 代码检查
go vet ./...

# 热重载开发 (需要安装 air)
air
```

## API 端点

### 认证模块
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/login | 用户登录 |
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/logout | 用户登出 |
| GET | /api/v1/auth/me | 获取当前用户 |
| PUT | /api/v1/auth/password | 修改密码 |
| PUT | /api/v1/auth/profile | 更新资料 |

### 用户管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/users | 用户列表 |
| GET | /api/v1/users/:id | 用户详情 |
| POST | /api/v1/users | 创建用户 |
| PUT | /api/v1/users/:id | 更新用户 |
| DELETE | /api/v1/users/:id | 删除用户 |

### 角色管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/roles | 角色列表 |
| GET | /api/v1/roles/:id | 角色详情 |
| POST | /api/v1/roles | 创建角色 |
| PUT | /api/v1/roles/:id | 更新角色 |
| DELETE | /api/v1/roles/:id | 删除角色 |

### 权限管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/permissions | 权限列表 |
| POST | /api/v1/permissions | 创建权限 |
| PUT | /api/v1/permissions/:id | 更新权限 |
| DELETE | /api/v1/permissions/:id | 删除权限 |

### 菜单管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/menus | 菜单树 |
| GET | /api/v1/menus/:id | 菜单详情 |
| POST | /api/v1/menus | 创建菜单 |
| PUT | /api/v1/menus/:id | 更新菜单 |
| DELETE | /api/v1/menus/:id | 删除菜单 |

### 仪表盘
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/dashboard/stats | 统计数据 |
| GET | /api/v1/dashboard/recent-logs | 最近日志 |
| GET | /api/v1/dashboard/system-info | 系统信息 |

### 通知系统
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/notifications | 通知列表 |
| GET | /api/v1/notifications/unread-count | 未读数量 |
| POST | /api/v1/notifications | 创建通知 |
| PUT | /api/v1/notifications/:id/read | 标记已读 |
| PUT | /api/v1/notifications/read-all | 全部已读 |
| DELETE | /api/v1/notifications/:id | 删除通知 |

### 定时任务
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/scheduled-tasks | 任务列表 |
| GET | /api/v1/scheduled-tasks/:id | 任务详情 |
| POST | /api/v1/scheduled-tasks | 创建任务 |
| PUT | /api/v1/scheduled-tasks/:id | 更新任务 |
| DELETE | /api/v1/scheduled-tasks/:id | 删除任务 |
| POST | /api/v1/scheduled-tasks/:id/toggle | 启用/禁用 |
| POST | /api/v1/scheduled-tasks/:id/run | 立即执行 |

### 数据导出
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/export/csv | 导出 CSV |
| POST | /api/v1/export/json | 导出 JSON |
| GET | /api/v1/export/template | 下载模板 |

### 系统设置
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/settings | 设置列表 |
| GET | /api/v1/settings/:group | 分组设置 |
| PUT | /api/v1/settings/:key | 更新设置 |

## 项目结构

```
admin-system/
├── app/
│   ├── controllers/     # HTTP 控制器
│   ├── models/          # ORM 模型
│   ├── middleware/      # 自定义中间件
│   ├── config/          # 配置管理
│   ├── scheduler/       # 定时任务调度器
│   └── testutil/        # 测试工具
├── cmd/
│   ├── admin-cli/       # CLI 工具
│   └── db-check/        # 数据库检查工具
├── config/
│   └── app.yaml         # 应用配置
├── database/
│   └── migrations/      # 数据库迁移
├── routes/              # 路由注册
├── scripts/
│   └── test-api.sh      # API 测试脚本
├── storage/
│   ├── database.db      # SQLite 数据库
│   ├── logs/            # 日志文件
│   └── uploads/         # 上传文件
├── web/                 # 前端界面 (Layui)
├── .air.toml            # 热重载配置
├── Makefile             # 构建命令
└── main.go              # 入口文件
```

## 配置说明

### 环境变量 (.env)

```bash
# 服务配置
PORT=8080
ENV=development
DEBUG=true

# 数据库
DB_DRIVER=sqlite
DB_DSN=storage/database.db

# JWT
JWT_SECRET=your-secret-key
JWT_TTL=7200

# 日志
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=storage/logs/app.log
```

## 测试覆盖

- 控制器测试：认证、用户、角色、仪表盘
- 模型测试：BaseModel、User、Role
- 中间件测试：CSRF、限流、XSS 防护
- 配置测试：校验规则

运行 `go test ./... -v` 查看所有测试结果。

## 开发规范

1. **导入别名**: `import ghttp "github.com/Hlgxz/gai/http"`
2. **Handler 签名**: `func Handler(c *ghttp.Context)`
3. **ORM 查询**: 使用泛型函数 `orm.Query[T]()`, `orm.Get[T]()`, `orm.Create[T]()`
4. **模型定义**: 嵌入 `BaseModel`，使用 `gai:"..."` 标签
5. **错误处理**: 使用 `c.Error(code, message)` 返回错误
6. **成功响应**: 使用 `c.Success(data)` 返回数据


# Admin-System 项目完善记录

> 本文档记录了项目从初始搭建到生产可用级别的所有改进工作。

---

## 一、项目概述

### 基本信息
- **框架**: Gai Framework (github.com/Hlgxz/gai) - AI原生 Go Web 框架
- **数据库**: SQLite (开发) / MySQL/PostgreSQL (生产推荐)
- **认证**: JWT Token
- **权限**: RBAC (基于角色的访问控制)
- **前端**: Layui + 原生 JavaScript

### 目录结构
```
admin-system/
├── app/
│   ├── controllers/      # HTTP 控制器
│   ├── models/           # ORM 模型
│   ├── middleware/       # 自定义中间件
│   ├── utils/            # 工具函数
│   ├── config/           # 配置解析
│   └── scheduler/        # 定时任务调度器
├── config/
│   ├── app.yaml          # 开发环境配置
│   └── app.production.yaml  # 生产环境配置示例
├── database/
│   └── migrations/       # 数据库迁移
├── routes/
│   └── routes.go         # 路由注册
├── web/
│   ├── layui/            # Layui 框架
│   ├── css/              # 自定义样式
│   ├── js/               # JavaScript 文件
│   └── pages/            # 多页面 HTML
│       ├── dashboard/    # 数据看板
│       ├── system/       # 系统管理
│       ├── data/         # 数据管理
│       └── logs/         # 日志管理
├── storage/
│   ├── logs/             # 日志文件
│   ├── uploads/          # 上传文件
│   └── database/         # SQLite 数据库
├── Dockerfile            # Docker 部署
├── .dockerignore         # Docker 忽略文件
├── .env                  # 环境变量 (开发)
├── .env.example          # 环境变量模板
└── main.go               # 应用入口
```

---

## 二、功能模块

### 已实现功能

| 模块 | 功能 | 状态 |
|------|------|------|
| **认证** | 用户登录/注册 | ✅ |
| | JWT Token 认证 | ✅ |
| | 密码强度验证 | ✅ |
| | 登录速率限制 | ✅ |
| | 注册速率限制 | ✅ |
| | Token 过期自动跳转 | ✅ |
| **用户管理** | 用户 CRUD | ✅ |
| | 用户状态管理 | ✅ |
| | 密码修改 | ✅ |
| **角色管理** | 角色 CRUD | ✅ |
| | 角色状态管理 | ✅ |
| | 角色-用户关联 | ✅ |
| | 角色-权限关联 | ✅ |
| | 删除前关联检查 | ✅ |
| **权限管理** | 权限 CRUD | ✅ |
| | 权限分组展示 | ✅ |
| | 权限缓存机制 | ✅ |
| **菜单管理** | 菌单 CRUD | ✅ |
| | 父级菜单选择 | ✅ |
| | 菜单排序 | ✅ |
| **操作日志** | 自动记录操作 | ✅ |
| | 日志列表/详情 | ✅ |
| | 日志导出 | ✅ |
| | 敏感字段脱敏 | ✅ |
| | 定期清理机制 | ✅ |
| **通知管理** | 通知 CRUD | ✅ |
| | 标记已读 | ✅ |
| | 批量标记已读 | ✅ |
| **定时任务** | 任务 CRUD | ✅ |
| | 手动执行 | ✅ |
| | 执行记录 | ✅ |
| | 定期清理机制 | ✅ |
| **文件上传** | 文件上传 | ✅ |
| | 文件大小限制 (10MB) | ✅ |
| | MIME 类型验证 | ✅ |
| | 安全随机文件名 | ✅ |
| **数据导出** | CSV/JSON 导出 | ✅ |
| | 表名白名单验证 | ✅ |
| | 列名安全验证 | ✅ |
| | 参数化查询防注入 | ✅ |
| **数据导入** | CSV/JSON 导入 | ✅ |
| | 模板下载 | ✅ |
| | 导入结果反馈 | ✅ |
| **系统设置** | 设置 CRUD | ✅ |
| | 批量更新 | ✅ |
| **数据看板** | 统计数据展示 | ✅ |
| | 系统信息 | ✅ |

---

## 三、安全改进

### 3.1 认证安全

#### JWT 密钥验证
**位置**: `main.go` 启动时

```go
// 验证 JWT 密钥安全性
weakSecrets := []string{"change-me", "change-me-to-a-random-string", "secret", "password", "admin", "test"}
for _, weak := range weakSecrets {
    if appCfg.JWTSecret == weak || len(appCfg.JWTSecret) < 16 {
        log.Fatal("JWT secret must be at least 16 characters and not a common weak value.")
    }
}
```

**效果**: 生产环境启动时强制要求强密钥，拒绝弱密钥。

#### 登录速率限制
**位置**: `app/middleware/login_rate_limit.go`

- 最大尝试次数: 5 次
- 锁定时长: 15 分钟
- 窗口时长: 5 分钟
- 自动清理过期记录

#### 注册速率限制
**位置**: `routes/routes.go`

```go
registerWithRateLimit := func(c *ghttp.Context) {
    ip := c.ClientIP()
    if middleware.IsLoginBlocked(ip) {
        c.Error(http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
        return
    }
    authController.Register(c)
}
```

#### 密码强度验证
**位置**: `app/utils/password.go`

- 最小长度: 8 字符
- 必须包含: 字母 + 数字
- 弱密码黑名单检测

### 3.2 SQL 注入防护

#### 表名白名单
**位置**: `app/controllers/export_controller.go`

```go
var allowedTables = map[string]bool{
    "users":           true,
    "roles":           true,
    "menus":           true,
    "permissions":     true,
    "operation_logs":  true,
    "notifications":   true,
    "scheduled_tasks": true,
    "uploads":         true,
    "settings":        true,
}
```

#### 列名安全验证
```go
func validateColumnName(col string) error {
    for _, c := range col {
        if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || 
             (c >= '0' && c <= '9') || c == '_') {
            return fmt.Errorf("无效的列名: %s", col)
        }
    }
    return nil
}
```

#### 参数化查询
所有动态查询使用参数化方式，避免 SQL 拼接。

### 3.3 敏感信息保护

#### 操作日志脱敏
**位置**: `app/middleware/operation_log.go`

```go
var sensitiveFields = map[string]bool{
    "password":        true,
    "new_password":    true,
    "token":           true,
    "secret":          true,
    "authorization":   true,
    // ...
}

func maskSensitiveFields(data any) any {
    // 递归脱敏，将敏感字段值替换为 "***"
}
```

#### 错误信息处理
**位置**: `app/middleware/error_handler.go`

```go
func SafeError(err error, userMessage string) string {
    if err != nil {
        slog.Error("Internal error", "error", err.Error())
    }
    if IsProduction() {
        return userMessage  // 生产环境返回通用消息
    }
    // 开发环境可返回详细信息
}
```

### 3.4 CSRF/CORS 安全

#### CSRF 动态配置
**位置**: `app/middleware/csrf.go`

```go
func isSecureEnv() bool {
    return os.Getenv("APP_ENV") == "production" || os.Getenv("HTTPS") == "true"
}
```

- 生产环境自动启用 `Secure: true`
- 开发环境使用 `Secure: false`

#### CORS 环境变量配置
**位置**: `app/middleware/security_headers.go`

```go
// 通过 CORS_ALLOW_ORIGINS 环境变量配置允许的来源
if envOrigins := os.Getenv("CORS_ALLOW_ORIGINS"); envOrigins != "" {
    allowOrigins = strings.Split(envOrigins, ",")
}
```

---

## 四、稳定性改进

### 4.1 优雅关闭机制
**位置**: `main.go`

```go
// 等待中断信号进行优雅关闭
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 给正在处理的请求最多30秒完成
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// 停止定时任务调度器
sched.Stop()

// 关闭 HTTP 服务器
srv.Shutdown(ctx)
```

**效果**: 服务终止时等待请求完成，避免数据丢失。

### 4.2 Panic Recovery 中间件
**位置**: `app/middleware/recovery.go`

```go
func RecoveryMiddleware() ghttp.HandlerFunc {
    return func(c *ghttp.Context) {
        defer func() {
            if err := recover() {
                // 记录堆栈信息到日志
                stack := string(debug.Stack())
                Error("Panic recovered", "error", err, "stack", stack)
                
                // 返回通用 500 错误
                c.Error(http.StatusInternalServerError, "服务器内部错误")
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

**效果**: 单个请求 panic 不会导致整个服务崩溃。

### 4.3 健康检查完善
**位置**: `main.go`

```go
app.Router().Get("/health", func(c *ghttp.Context) {
    // 检查数据库连接
    dbErr := db.SQL.PingContext(ctx)
    if dbErr != nil {
        c.Error(http.StatusServiceUnavailable, "database unavailable")
        return
    }
    c.Success(map[string]any{
        "status":   "ok",
        "database": "connected",
    })
})
```

**效果**: 容器编排可检测服务健康状态。

### 4.4 内存泄漏防护

#### 限流器清理
**位置**: `app/middleware/rate_limit.go`

```go
// 定期清理过期条目 (每 5 分钟)
func (rl *RateLimiter) cleanupExpiredEntries() {
    ticker := time.NewTicker(5 * time.Minute)
    for {
        select {
        case <-ticker.C:
            // 清理过期的窗口或解封的条目
        case <-rl.stopCleanup:
            return
        }
    }
}
```

#### 登录限流清理
**位置**: `app/middleware/login_rate_limit.go`

```go
// 定期清理过期登录尝试记录 (每 10 分钟)
func (l *LoginRateLimiter) cleanupExpiredAttempts() {
    ticker := time.NewTicker(10 * time.Minute)
    // ...
}
```

#### 权限缓存清理
**位置**: `app/middleware/rbac.go`

```go
// 定期清理过期缓存条目 (每 1 分钟)
func (pc *PermissionCache) cleanupExpiredCache() {
    ticker := time.NewTicker(time.Minute)
    // ...
}
```

---

## 五、可观测性改进

### 5.1 请求 ID 中间件
**位置**: `app/middleware/request_id.go`

```go
func RequestIDMiddleware() ghttp.HandlerFunc {
    return func(c *ghttp.Context) {
        requestID := c.Header("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.SetHeader("X-Request-ID", requestID)
        c.Set("request_id", requestID)
        c.Next()
    }
}
```

**效果**: 每个请求有唯一 ID，便于日志追踪和问题排查。

### 5.2 定时清理任务
**位置**: `main.go`

| 任务名 | 清理对象 | 保留天数 |
|--------|----------|----------|
| cleanup_logs | 日志文件 | 30 天 |
| cleanup_uploads | 上传文件 | 90 天 |
| cleanup_operation_logs | 操作日志 | 90 天 |
| cleanup_task_executions | 任务执行记录 | 30 天 |

---

## 六、部署配置

### 6.1 Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o admin-system .

# 运行阶段
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/admin-system .
RUN mkdir -p /app/storage/logs /app/storage/uploads /app/storage/database

ENV TZ=Asia/Shanghai
ENV APP_ENV=production
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --spider http://localhost:8080/health || exit 1

CMD ["./admin-system"]
```

### 6.2 环境变量配置 (.env.example)

```bash
# 应用配置
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# 数据库配置
# SQLite (开发环境)
DB_DRIVER=sqlite
DB_DATABASE=storage/database.db

# MySQL (生产环境推荐)
# DB_DRIVER=mysql
# DB_DATABASE=user:password@tcp(localhost:3306)/admin_system?charset=utf8mb4&parseTime=True&loc=Local

# JWT 配置 (生产环境必须设置强密钥)
JWT_SECRET=your-strong-jwt-secret-at-least-16-characters
JWT_TTL=3600

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=storage/logs/app.log
LOG_REQUESTS=true

# CORS 配置
# CORS_ALLOW_ORIGINS=https://example.com,https://api.example.com
```

### 6.3 生产环境配置 (config/app.production.yaml)

```yaml
name: admin-system
port: 8080
env: production
debug: false

database:
  driver: mysql
  dsn: ${DB_DSN}

auth:
  default: jwt
  guards:
    jwt:
      driver: jwt
      secret: ${JWT_SECRET}
      ttl: 3600

log:
  level: warn
  format: json
  file: storage/logs/app.log
  requests: false
```

---

## 七、前端改进

### 7.1 多页面架构
从单页面重构为多页面架构:

| 页面 | 路径 | 功能 |
|------|------|------|
| 登录页 | `/login` | 用户认证 |
| 数据看板 | `/pages/dashboard/index.html` | 统计展示 |
| 用户管理 | `/pages/system/users.html` | 用户 CRUD |
| 角色管理 | `/pages/system/roles.html` | 角色 CRUD + 关联 |
| 权限管理 | `/pages/system/permissions.html` | 权限 CRUD |
| 菌单管理 | `/pages/system/menus.html` | 菜单 CRUD |
| 定时任务 | `/pages/system/tasks.html` | 任务管理 |
| 通知管理 | `/pages/dashboard/notifications.html` | 通知管理 |
| 个人中心 | `/pages/dashboard/profile.html` | 个人设置 |
| 操作日志 | `/pages/logs/operation.html` | 日志查看 |
| 数据导出 | `/pages/data/export.html` | 数据导出 |
| 数据导入 | `/pages/data/import.html` | 数据导入 |

### 7.2 字段名修复
修复前端与后端 API 字段名不匹配问题:

| 页面 | 前端字段 | 后端字段 | 状态 |
|------|----------|----------|------|
| menus.html | title | name | ✅ 已修复 |
| menus.html | sort | sort_order | ✅ 已修复 |
| permissions.html | slug | code | ✅ 已修复 |
| tasks.html | cron | cron_expr | ✅ 已修复 |
| tasks.html | handler | command | ✅ 已修复 |
| notifications.html | /read-all | /mark-all-read | ✅ 已修复 |
| profile.html | /auth/profile | /auth/me | ✅ 已修复 |

### 7.3 功能增强
- 角色-权限关联对话框
- 菜单父级选择下拉框
- 前端密码强度验证
- Token 过期自动跳转登录
- 表格数据加载提示
- 操作成功/失败提示

---

## 八、代码清理

### 已删除文件
- `web/test-menu.html` - 测试文件
- `web/test-index.html` - 测试文件
- `web/layouts/main.html` - 未使用模板
- `web/js/page.js` - 未使用脚本

### 已修复问题
- `main.go` 中 `audit_logs` 表名改为 `operation_logs`
- `upload_controller.go` 中 `randomString` 使用 `crypto/rand`

---

## 九、部署指南

### 本地开发
```bash
# 1. 安装依赖
go mod tidy

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 设置 JWT_SECRET

# 3. 启动服务
go run main.go

# 4. 访问系统
http://localhost:8080/login
默认账号: admin / admin123
```

### Docker 部署
```bash
# 1. 构建镜像
docker build -t admin-system:latest .

# 2. 运行容器
docker run -d \
  --name admin-system \
  -p 8080:8080 \
  -e JWT_SECRET=your-strong-secret \
  -e DB_DRIVER=mysql \
  -e DB_DSN="user:pass@tcp(host:3306)/db" \
  admin-system:latest

# 3. 健康检查
curl http://localhost:8080/health
```

### 生产环境清单

| 检查项 | 要求 |
|--------|------|
| APP_ENV | production |
| APP_DEBUG | false |
| JWT_SECRET | >= 16 字符，非弱密钥 |
| DB_DRIVER | mysql 或 postgres |
| CORS_ALLOW_ORIGINS | 限制为指定域名 |
| HTTPS | 强制 HTTPS |
| 数据备份 | 定期备份数据库 |

---

## 十、后续优化建议

### 可继续完善项

| 优先级 | 项目 | 说明 |
|--------|------|------|
| Medium | 其他控制器错误处理 | 统一使用 SafeError |
| Medium | 分页最大值限制 | 限制 perPage 最大 100 |
| Low | 数据库连接池 | 配置 MaxOpenConns 等 |
| Low | Prometheus metrics | 添加监控端点 |
| Low | 乐观锁 | 添加 version 字段防并发 |

---

## 十一、版本历史

| 版本 | 日期 | 主要改进 |
|------|------|----------|
| v1.0 | 2026-06-12 | 初始功能实现 |
| v1.1 | 2026-06-13 | 安全加固、稳定性改进 |
| v1.2 | 2026-06-16 | 生产级完善、部署配置 |

---

**文档维护**: 本文档记录项目改进历程，后续更新请在此文档追加内容。

## 许可证

MIT
