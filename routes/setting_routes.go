package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
)

// RegisterSettingRoutes sets up setting routes.
func RegisterSettingRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	ctrl := controllers.NewSettingController(db)

	// Public settings (no auth required)
	r.Get("/api/v1/settings/public", ctrl.Public)

	// Admin settings (requires auth, no special permission needed for basic settings)
	r.Group("/api/v1/settings", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))
		g.Get("", ctrl.List)
		g.Post("", ctrl.Create)
		g.Get("/:key", ctrl.Get)
		g.Put("/:key", ctrl.Update)
		g.Delete("/:key", ctrl.Delete)
		g.Put("/batch", ctrl.BatchUpdate)
	})
}