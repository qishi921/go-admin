package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
)

// RegisterNotificationRoutes 注册通知相关路由
func RegisterNotificationRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	nc := &controllers.NotificationController{DB: db}

	// 用户通知接口（需要登录，无需特殊权限）
	r.Group("/api/v1/notifications", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Get("", nc.List)
		g.Get("/unread-count", nc.GetUnreadCount)
		g.Post("/mark-read/:id", nc.MarkRead)
		g.Post("/mark-all-read", nc.MarkAllRead)
		g.Delete("/:id", nc.Delete)
	})

	// 管理员通知接口
	r.Group("/api/v1/admin/notifications", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Post("", nc.Create)
	})
}