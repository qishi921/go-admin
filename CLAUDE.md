# CLAUDE.md — Gai Framework Project

## Project

- **Module**: github.com/user/admin-system
- **Framework**: Gai (github.com/Hlgxz/gai) — AI-native Go web framework
- **Docs**: https://github.com/Hlgxz/gai

## Setup

```bash
go mod tidy && go build ./...
```

## Architecture

```
admin-system/
├── app/controllers/   # HTTP 控制器
├── app/models/        # ORM 模型
├── app/middleware/     # 自定义中间件
├── config/app.yaml    # 配置 (支持 ${ENV_VAR:default})
├── database/migrations/
├── routes/routes.go   # 路由注册入口
├── schemas/           # Schema 定义 (YAML → 自动生成代码)
├── storage/           # 日志、SQLite、上传文件
├── .env               # 环境变量
└── main.go            # 入口
```

## Gai Framework Conventions

### 1. Import 别名 (强制)
```go
import ghttp "github.com/Hlgxz/gai/http"  // 必须用 ghttp 别名
```

### 2. Handler 签名
```go
func Handler(c *ghttp.Context) {
    c.Success(data)        // → {"code":0,"message":"ok","data":...}
    c.Error(400, "msg")    // → {"code":400,"message":"msg"}
}
```

### 3. ORM (泛型函数)
```go
import "github.com/Hlgxz/gai/database/orm"

// 查询
query := orm.Query[User](db).Where("status", "=", "active").OrderBy("id", "DESC")
users, _ := orm.Get[User](query)       // []User
user, _ := orm.First[User](query)      // *User
page, _ := orm.Paginate[User](query, 1, 20)

// 增删改
created, _ := orm.Create[User](db, &User{Name: "test"})
orm.Update[User](db, user)
orm.Delete[User](db, user)             // 软删除
```

### 4. Model 定义
```go
import "github.com/Hlgxz/gai/database/orm"

type User struct {
    orm.Model                                        // ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string `json:"name"  gai:"column:name;size:100"`
    Email string `json:"email" gai:"column:email;unique"`
    Posts []Post `json:"-"     gai:"hasMany"`
}
```

### 5. 路由
```go
import "github.com/Hlgxz/gai/router"

r.Get("/path/:id", handler)
r.Group("/api/v1", func(g *router.Group) {
    g.Use(authManager.Middleware("jwt"))
    g.Resource("/users", userController) // 自动 CRUD 五条路由
})
```

### 6. 中间件
```go
func MyMiddleware() ghttp.HandlerFunc {
    return func(c *ghttp.Context) {
        // 前置逻辑
        c.Next()  // 必须调用
        // 后置逻辑
    }
}
```

### 7. 校验
```go
v := ghttp.NewValidator(data, map[string]string{
    "email": "required|email",
    "name":  "required|min:2|max:50",
    "phone": "phone",
})
if errs := v.Validate(); errs != nil { /* 422 */ }
```

### 8. Schema 驱动开发
在 schemas/ 中创建 YAML，运行 `gai generate --schema schemas/` 自动生成 Model + Controller + Migration + Routes。

## Package 速查

| 导入路径 | 别名 | 用途 |
|---------|------|------|
| github.com/Hlgxz/gai | gai | Application, Container, Make[T] |
| github.com/Hlgxz/gai/http | ghttp | Context, HandlerFunc, Validator |
| github.com/Hlgxz/gai/router | router | Router, Group, ResourceController |
| github.com/Hlgxz/gai/database/orm | orm | DB, Model, Query[T], Get[T], Create[T] |
| github.com/Hlgxz/gai/auth | auth | Manager, Guard, JWTGuard |
| github.com/Hlgxz/gai/middleware | middleware | CORS, Logger, Recovery, RateLimit |
| github.com/Hlgxz/gai/miniapp/wechat | wechat | Client, Auth, Pay, Message |
| github.com/Hlgxz/gai/miniapp/alipay | alipay | Client, Auth |
| github.com/Hlgxz/gai/support | support | Snake, Camel, Hash, Env |

## 禁止

- 导入 gai/http 不加 ghttp 别名
- 在 handler 中使用 net/http 的 ResponseWriter/Request
- 跳过中间件的 c.Next()
- 直接写 SQL（用 query builder 或 migration blueprint）
- 添加 gin/echo/chi 等第三方路由库
