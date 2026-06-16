package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
)

// RegisterScheduledTaskRoutes 注册定时任务路由
func RegisterScheduledTaskRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager, sched interface{}) {
	tc := &controllers.ScheduledTaskController{
		DB: db,
	}

	// 定时任务管理（需要登录）
	r.Group("/api/v1/scheduled-tasks", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Get("", tc.List)
		g.Get("/stats", tc.Stats)
		g.Post("", tc.Create)
		g.Get("/:id", tc.Get)
		g.Put("/:id", tc.Update)
		g.Delete("/:id", tc.Delete)
		g.Post("/:id/toggle", tc.Toggle)
		g.Post("/:id/run", tc.RunNow)
	})

	// 执行记录
	r.Group("/api/v1/task-executions", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))

		g.Get("", tc.Executions)
	})
}