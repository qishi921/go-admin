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

## 许可证

MIT
