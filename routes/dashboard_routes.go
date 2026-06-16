package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
)

// RegisterDashboardRoutes 注册仪表盘路由
func RegisterDashboardRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	dc := &controllers.DashboardController{DB: db}

	r.Group("/api/v1/dashboard", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Get("/stats", dc.Stats)
		g.Get("/recent-logs", dc.RecentLogs)
		g.Get("/system-info", dc.SystemInfo)
	})
}