package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
	"github.com/user/admin-system/app/middleware"
)

// RegisterRolePermissionRoutes sets up the role-permission association routes.
func RegisterRolePermissionRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	ctrl := controllers.NewRolePermissionController(db)

	r.Group("/api/v1/roles/:id/permissions", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))
		g.Use(middleware.RBACMiddleware(db))
		g.Use(middleware.OperationLogMiddleware(db))
		g.Get("", ctrl.GetRolePermissions)
		g.Post("", ctrl.AssignPermissionToRole)
		g.Delete("/:permissionId", ctrl.RemovePermissionFromRole)
	})
}
