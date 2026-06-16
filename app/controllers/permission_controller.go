package controllers

import (
	"fmt"
	"net/http"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// PermissionController handles CRUD operations for Permission.
type PermissionController struct {
	DB *orm.DB
}

// NewPermissionController creates a new controller instance.
func NewPermissionController(db *orm.DB) *PermissionController {
	return &PermissionController{DB: db}
}

// Index lists all permissions with pagination and search.
func (ctrl *PermissionController) Index(c *ghttp.Context) {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	search := c.Query("search")

	q := orm.Query[models.Permission](ctrl.DB)
	if search != "" {
		q = q.Where("name", "LIKE", "%"+search+"%")
	}
	q = q.OrderBy("created_at", "DESC")

	result, err := orm.Paginate[models.Permission](q, page, perPage)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.Success(result)
}

// Show returns a single permission by ID.
func (ctrl *PermissionController) Show(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Permission](
		orm.Query[models.Permission](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Permission not found")
		return
	}
	c.Success(item)
}

// Store creates a new permission.
func (ctrl *PermissionController) Store(c *ghttp.Context) {
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

	item := &models.Permission{}
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
	if v, ok := input["type"]; ok {
		if typed, ok := v.(string); ok {
			item.Type = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field type: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["parent_id"]; ok {
		if typed, ok := v.(int); ok {
			item.ParentId = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field parent_id: expected int, got %T", v))
			return
		}
	}
	if v, ok := input["sort_order"]; ok {
		if typed, ok := v.(int); ok {
			item.SortOrder = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field sort_order: expected int, got %T", v))
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

	result, err := orm.Create[models.Permission](ctrl.DB, item)
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

// Update modifies an existing permission.
func (ctrl *PermissionController) Update(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Permission](
		orm.Query[models.Permission](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Permission not found")
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
	if v, ok := input["type"]; ok {
		if typed, ok := v.(string); ok {
			item.Type = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field type: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["parent_id"]; ok {
		if typed, ok := v.(int); ok {
			item.ParentId = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field parent_id: expected int, got %T", v))
			return
		}
	}
	if v, ok := input["sort_order"]; ok {
		if typed, ok := v.(int); ok {
			item.SortOrder = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field sort_order: expected int, got %T", v))
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

	if err := orm.Update[models.Permission](ctrl.DB, item); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.Success(item)
}

// Destroy deletes a permission by ID.
func (ctrl *PermissionController) Destroy(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Permission](
		orm.Query[models.Permission](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Permission not found")
		return
	}

	if err := orm.Delete[models.Permission](ctrl.DB, item); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.NoContent()
}

// Ensure PermissionController satisfies the ResourceController interface.
var _ interface {
	Index(c *ghttp.Context)
	Show(c *ghttp.Context)
	Store(c *ghttp.Context)
	Update(c *ghttp.Context)
	Destroy(c *ghttp.Context)
} = (*PermissionController)(nil)
