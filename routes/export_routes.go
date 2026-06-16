package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
)

// RegisterExportRoutes 注册导入导出路由
func RegisterExportRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	ec := &controllers.ExportController{DB: db}

	r.Group("/api/v1/export", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Post("/csv", ec.ExportCSV)
		g.Post("/json", ec.ExportJSON)
		g.Get("/template", ec.ExportTemplate)
	})

	r.Group("/api/v1/import", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Post("/csv", ec.ImportCSV)
		g.Post("/json", ec.ImportJSON)
	})
}
