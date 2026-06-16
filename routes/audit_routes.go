package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
)

// RegisterAuditRoutes 注册审计日志路由
func RegisterAuditRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	ac := &controllers.AuditController{DB: db}

	r.Group("/api/v1/audit-logs", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Get("", ac.List)
		g.Get("/stats", ac.Stats)
		g.Get("/export", ac.Export)
		g.Get("/:id", ac.Detail)
	})

	// 数据变更历史
	r.Group("/api/v1/data-changes", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Get("/history", ac.DataChangeHistory)
	})
}
