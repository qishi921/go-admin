package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
	"github.com/user/admin-system/app/middleware"
)

// RegisterUploadRoutes sets up upload routes.
func RegisterUploadRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager, uploadDir string) {
	ctrl := controllers.NewUploadController(db, uploadDir)

	// Upload API (requires auth)
	r.Group("/api/v1/uploads", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))
		g.Use(middleware.RBACMiddleware(db))
		g.Use(middleware.OperationLogMiddleware(db))
		g.Post("", ctrl.Upload)
		g.Delete("/:id", ctrl.Delete)
		g.Get("", ctrl.List)
	})

	// File serving (public)
	r.Get("/uploads/*", ctrl.Serve)
}