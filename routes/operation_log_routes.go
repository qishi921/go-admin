package routes

import (
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	"github.com/Hlgxz/gai/auth"
	"github.com/user/admin-system/app/controllers"
)

// RegisterOperationLogRoutes sets up the OperationLog resource routes.
func RegisterOperationLogRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	ctrl := controllers.NewOperationLogController(db)

	r.Group("/api/v1/logs", func(g *router.Group) {
		g.Use(authMgr.Middleware("jwt"))
		g.Get("", ctrl.Index)
		g.Get("/:id", ctrl.Show)
	})
}
