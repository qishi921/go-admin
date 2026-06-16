package controllers

import (
	"net/http"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// RolePermissionController handles role-permission association operations.
type RolePermissionController struct {
	DB *orm.DB
}

// NewRolePermissionController creates a new controller instance.
func NewRolePermissionController(db *orm.DB) *RolePermissionController {
	return &RolePermissionController{DB: db}
}

// GetRolePermissions returns the list of permissions assigned to a role.
func (ctrl *RolePermissionController) GetRolePermissions(c *ghttp.Context) {
	roleID := c.ParamInt("id")

	// Verify role exists
	role, err := orm.First[models.Role](
		orm.Query[models.Role](ctrl.DB).Where("id", "=", roleID),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		c.Error(http.StatusNotFound, "Role not found")
		return
	}

	// Find permission IDs from role_permission table
	var rolePermissions []models.RolePermission
	rolePermissions, err = orm.Get[models.RolePermission](
		orm.Query[models.RolePermission](ctrl.DB).Where("role_id", "=", roleID),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	if len(rolePermissions) == 0 {
		c.Success([]models.Permission{})
		return
	}

	permissionIDs := make([]any, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissionIDs[i] = rp.PermissionID
	}

	permissions, err := orm.Get[models.Permission](
		orm.Query[models.Permission](ctrl.DB).WhereIn("id", permissionIDs),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.Success(permissions)
}

// AssignPermissionToRole assigns a permission to a role.
func (ctrl *RolePermissionController) AssignPermissionToRole(c *ghttp.Context) {
	roleID := c.ParamInt("id")

	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}

	permissionIDVal, ok := input["permission_id"]
	if !ok {
		c.Error(http.StatusBadRequest, "permission_id is required")
		return
	}

	permissionID, ok := permissionIDVal.(float64)
	if !ok {
		c.Error(http.StatusBadRequest, "permission_id must be a number")
		return
	}

	// Verify role exists
	role, err := orm.First[models.Role](
		orm.Query[models.Role](ctrl.DB).Where("id", "=", roleID),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		c.Error(http.StatusNotFound, "Role not found")
		return
	}

	// Verify permission exists
	permission, err := orm.First[models.Permission](
		orm.Query[models.Permission](ctrl.DB).Where("id", "=", int(permissionID)),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if permission == nil {
		c.Error(http.StatusNotFound, "Permission not found")
		return
	}

	// Check if already assigned
	existing, err := orm.First[models.RolePermission](
		orm.Query[models.RolePermission](ctrl.DB).
			Where("role_id", "=", roleID).
			Where("permission_id", "=", int(permissionID)),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if existing != nil {
		c.Error(http.StatusConflict, "Permission is already assigned to this role")
		return
	}

	// Create association
	rolePermission := &models.RolePermission{
		RoleID:       roleID,
		PermissionID: int(permissionID),
	}

	_, err = orm.Create[models.RolePermission](ctrl.DB, rolePermission)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"code":    0,
		"message": "Permission assigned to role successfully",
		"data":    rolePermission,
	})
}

// RemovePermissionFromRole removes a permission from a role.
func (ctrl *RolePermissionController) RemovePermissionFromRole(c *ghttp.Context) {
	roleID := c.ParamInt("id")
	permissionID := c.ParamInt("permissionId")

	// Find the association
	rolePermission, err := orm.First[models.RolePermission](
		orm.Query[models.RolePermission](ctrl.DB).
			Where("role_id", "=", roleID).
			Where("permission_id", "=", permissionID),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if rolePermission == nil {
		c.Error(http.StatusNotFound, "Association not found")
		return
	}

	// Delete the association
	if err := orm.Delete[models.RolePermission](ctrl.DB, rolePermission); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.NoContent()
}
