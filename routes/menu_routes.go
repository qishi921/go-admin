package routes

import (
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/user/admin-system/app/controllers"
	"github.com/user/admin-system/app/middleware"
)

// RegisterMenuRoutes sets up the Menu resource routes.
func RegisterMenuRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	ctrl := controllers.NewMenuController(db)

	r.Group("/api/v1/menus", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))
		g.Use(middleware.RBACMiddleware(db))
		g.Use(middleware.OperationLogMiddleware(db))
		g.Get("", ctrl.Index)
		g.Get("/tree", ctrl.Tree)
		g.Post("", ctrl.Store)
		g.Get("/:id", ctrl.Show)
		g.Put("/:id", ctrl.Update)
		g.Delete("/:id", ctrl.Destroy)
	})
}
