package controllers

import (
	"net/http"

	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	ghttp "github.com/Hlgxz/gai/http"
	"github.com/user/admin-system/app/middleware"
)

// PermissionCheckController handles permission check endpoints.
type PermissionCheckController struct {
	DB *orm.DB
}

// NewPermissionCheckController creates a new controller instance.
func NewPermissionCheckController(db *orm.DB) *PermissionCheckController {
	return &PermissionCheckController{DB: db}
}

// MyPermissions returns all permission codes for the current user.
// Frontend uses this to determine which UI elements to show.
func (ctrl *PermissionCheckController) MyPermissions(c *ghttp.Context) {
	permissions, err := middleware.GetCurrentUserPermissions(c, ctrl.DB)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Check if user is super admin
	isSuperAdmin := false
	if claims, ok := c.Get("auth_claims"); ok {
		if jwtClaims, ok := claims.(*auth.Claims); ok {
			userID := int(jwtClaims.UserID)
			// Query directly for super admin check
			var count int
			err := ctrl.DB.SQL.QueryRow(`
				SELECT COUNT(*) FROM roles r
				INNER JOIN role_user ru ON r.id = ru.role_id
				WHERE ru.user_id = ? AND r.code = 'super_admin' AND r.deleted_at IS NULL
			`, userID).Scan(&count)
			if err == nil && count > 0 {
				isSuperAdmin = true
			}
		}
	}

	c.Success(map[string]any{
		"permissions":    permissions,
		"is_super_admin": isSuperAdmin,
	})
}