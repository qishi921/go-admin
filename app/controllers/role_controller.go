package controllers

import (
	"fmt"
	"net/http"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// RoleController handles CRUD operations for Role.
type RoleController struct {
	DB *orm.DB
}

// NewRoleController creates a new controller instance.
func NewRoleController(db *orm.DB) *RoleController {
	return &RoleController{DB: db}
}

// Index lists all roles with pagination and search.
func (ctrl *RoleController) Index(c *ghttp.Context) {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	search := c.Query("search")

	q := orm.Query[models.Role](ctrl.DB)
	if search != "" {
		q = q.Where("name", "LIKE", "%"+search+"%")
	}
	q = q.OrderBy("created_at", "DESC")

	result, err := orm.Paginate[models.Role](q, page, perPage)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.Success(result)
}

// Show returns a single role by ID.
func (ctrl *RoleController) Show(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Role](
		orm.Query[models.Role](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Role not found")
		return
	}
	c.Success(item)
}

// Store creates a new role.
func (ctrl *RoleController) Store(c *ghttp.Context) {
	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}
	validator := ghttp.NewValidator(input, map[string]string{
		"name": "required|min:2|max:50",
		"code": "required|min:2|max:50",
	})
	if errs := validator.Validate(); errs != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]any{
			"code":    422,
			"message": "Validation failed",
			"errors":  errs,
		})
		return
	}

	item := &models.Role{}
	if v, ok := input["name"]; ok {
		if typed, ok := v.(string); ok {
			item.Name = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field name: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["code"]; ok {
		if typed, ok := v.(string); ok {
			item.Code = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field code: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["description"]; ok {
		if typed, ok := v.(string); ok {
			item.Description = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field description: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["status"]; ok {
		if typed, ok := v.(string); ok {
			item.Status = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field status: expected string, got %T", v))
			return
		}
	}

	result, err := orm.Create[models.Role](ctrl.DB, item)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]any{
		"code":    0,
		"message": "ok",
		"data":    result,
	})
}

// Update modifies an existing role.
func (ctrl *RoleController) Update(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Role](
		orm.Query[models.Role](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Role not found")
		return
	}

	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}
	if v, ok := input["name"]; ok {
		if typed, ok := v.(string); ok {
			item.Name = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field name: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["code"]; ok {
		if typed, ok := v.(string); ok {
			item.Code = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field code: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["description"]; ok {
		if typed, ok := v.(string); ok {
			item.Description = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field description: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["status"]; ok {
		if typed, ok := v.(string); ok {
			item.Status = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field status: expected string, got %T", v))
			return
		}
	}

	if err := orm.Update[models.Role](ctrl.DB, item); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.Success(item)
}

// Destroy deletes a role by ID.
func (ctrl *RoleController) Destroy(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Role](
		orm.Query[models.Role](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Role not found")
		return
	}

	// Check if role has associated users
	var userCount int
	err = ctrl.DB.SQL.QueryRow(
		"SELECT COUNT(*) FROM role_user WHERE role_id = ?",
		id,
	).Scan(&userCount)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to check role associations")
		return
	}
	if userCount > 0 {
		c.Error(http.StatusConflict, fmt.Sprintf("无法删除：该角色已分配给 %d 个用户，请先移除用户关联", userCount))
		return
	}

	// Check if role has associated permissions
	var permCount int
	err = ctrl.DB.SQL.QueryRow(
		"SELECT COUNT(*) FROM role_permission WHERE role_id = ?",
		id,
	).Scan(&permCount)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to check role associations")
		return
	}

	// Delete role-permission associations first
	if permCount > 0 {
		_, err = ctrl.DB.SQL.Exec("DELETE FROM role_permission WHERE role_id = ?", id)
		if err != nil {
			c.Error(http.StatusInternalServerError, "Failed to remove role permissions")
			return
		}
	}

	// Check if it's a system role (cannot be deleted)
	if item.Code == "super_admin" {
		c.Error(http.StatusForbidden, "无法删除系统超级管理员角色")
		return
	}

	if err := orm.Delete[models.Role](ctrl.DB, item); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.NoContent()
}

// Ensure RoleController satisfies the ResourceController interface.
var _ interface {
	Index(c *ghttp.Context)
	Show(c *ghttp.Context)
	Store(c *ghttp.Context)
	Update(c *ghttp.Context)
	Destroy(c *ghttp.Context)
} = (*RoleController)(nil)
