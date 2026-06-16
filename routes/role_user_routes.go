package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
	"github.com/user/admin-system/app/middleware"
)

// RegisterRoleUserRoutes sets up the role-user association routes.
func RegisterRoleUserRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	ctrl := controllers.NewRoleUserController(db)

	// Role-centric routes
	r.Group("/api/v1/roles/:id/users", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))
		g.Use(middleware.RBACMiddleware(db))
		g.Use(middleware.OperationLogMiddleware(db))
		g.Get("", ctrl.GetRoleUsers)
		g.Post("", ctrl.AssignUserToRole)
		g.Delete("/:userId", ctrl.RemoveUserFromRole)
	})

	// User-centric routes
	r.Group("/api/v1/users/:id/roles", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))
		g.Use(middleware.RBACMiddleware(db))
		g.Use(middleware.OperationLogMiddleware(db))
		g.Get("", ctrl.GetUserRoles)
	})
}
