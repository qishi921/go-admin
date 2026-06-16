package controllers

import (
	"fmt"
	"net/http"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// MenuController handles CRUD operations for Menu.
type MenuController struct {
	DB *orm.DB
}

// NewMenuController creates a new controller instance.
func NewMenuController(db *orm.DB) *MenuController {
	return &MenuController{DB: db}
}

// MenuTreeNode represents a menu node in tree structure.
type MenuTreeNode struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Icon      string         `json:"icon"`
	SortOrder int            `json:"sort_order"`
	Status    string         `json:"status"`
	ParentId  *int           `json:"parent_id"`
	Children  []*MenuTreeNode `json:"children,omitempty"`
}

// Tree returns all menus in tree structure.
func (ctrl *MenuController) Tree(c *ghttp.Context) {
	menus, err := orm.Get[models.Menu](
		orm.Query[models.Menu](ctrl.DB).OrderBy("sort_order", "ASC"),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Build tree
	tree := buildMenuTree(menus, nil)
	c.Success(tree)
}

// buildMenuTree converts flat menu list to tree structure.
func buildMenuTree(menus []models.Menu, parentId *int) []*MenuTreeNode {
	var nodes []*MenuTreeNode
	for _, menu := range menus {
		if menu.ParentId == nil && parentId == nil ||
		   (menu.ParentId != nil && parentId != nil && *menu.ParentId == *parentId) {
			node := &MenuTreeNode{
				ID:        int(menu.ID),
				Name:      menu.Name,
				Path:      menu.Path,
				Icon:      menu.Icon,
				SortOrder: menu.SortOrder,
				Status:    menu.Status,
				ParentId:  menu.ParentId,
			}
			// Find children
			node.Children = buildMenuTree(menus, &node.ID)
			if len(node.Children) == 0 {
				node.Children = nil
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// Index lists all menus with pagination and search.
func (ctrl *MenuController) Index(c *ghttp.Context) {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	search := c.Query("search")

	q := orm.Query[models.Menu](ctrl.DB)
	if search != "" {
		q = q.Where("name", "LIKE", "%"+search+"%")
	}
	q = q.OrderBy("sort_order", "ASC")

	result, err := orm.Paginate[models.Menu](q, page, perPage)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.Success(result)
}

// Show returns a single menu by ID.
func (ctrl *MenuController) Show(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Menu](
		orm.Query[models.Menu](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Menu not found")
		return
	}
	c.Success(item)
}

// Store creates a new menu.
func (ctrl *MenuController) Store(c *ghttp.Context) {
	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}
	validator := ghttp.NewValidator(input, map[string]string{
		"name": "required|min:2|max:50",
		"path": "required",
	})
	if errs := validator.Validate(); errs != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]any{
			"code":    422,
			"message": "Validation failed",
			"errors":  errs,
		})
		return
	}

	item := &models.Menu{}
	if v, ok := input["name"]; ok {
		if typed, ok := v.(string); ok {
			item.Name = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field name: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["path"]; ok {
		if typed, ok := v.(string); ok {
			item.Path = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field path: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["icon"]; ok {
		if typed, ok := v.(string); ok {
			item.Icon = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field icon: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["component"]; ok {
		if typed, ok := v.(string); ok {
			item.Component = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field component: expected string, got %T", v))
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
	if v, ok := input["parent_id"]; ok {
		if typed, ok := v.(int); ok {
			item.ParentId = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field parent_id: expected int, got %T", v))
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

	result, err := orm.Create[models.Menu](ctrl.DB, item)
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

// Update modifies an existing menu.
func (ctrl *MenuController) Update(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Menu](
		orm.Query[models.Menu](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Menu not found")
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
	if v, ok := input["path"]; ok {
		if typed, ok := v.(string); ok {
			item.Path = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field path: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["icon"]; ok {
		if typed, ok := v.(string); ok {
			item.Icon = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field icon: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["component"]; ok {
		if typed, ok := v.(string); ok {
			item.Component = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field component: expected string, got %T", v))
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
	if v, ok := input["parent_id"]; ok {
		if typed, ok := v.(int); ok {
			item.ParentId = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field parent_id: expected int, got %T", v))
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

	if err := orm.Update[models.Menu](ctrl.DB, item); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.Success(item)
}

// Destroy deletes a menu by ID.
func (ctrl *MenuController) Destroy(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.Menu](
		orm.Query[models.Menu](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "Menu not found")
		return
	}

	if err := orm.Delete[models.Menu](ctrl.DB, item); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.NoContent()
}

// Ensure MenuController satisfies the ResourceController interface.
var _ interface {
	Index(c *ghttp.Context)
	Show(c *ghttp.Context)
	Store(c *ghttp.Context)
	Update(c *ghttp.Context)
	Destroy(c *ghttp.Context)
} = (*MenuController)(nil)
