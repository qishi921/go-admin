package controllers

import (
	"fmt"
	"net/http"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/middleware"
	"github.com/user/admin-system/app/models"
	"github.com/user/admin-system/app/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles CRUD operations for User.
type UserController struct {
	DB *orm.DB
}

// NewUserController creates a new controller instance.
func NewUserController(db *orm.DB) *UserController {
	return &UserController{DB: db}
}

// Index lists all users with pagination and search.
func (ctrl *UserController) Index(c *ghttp.Context) {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	search := c.Query("search")
	status := c.Query("status")

	// Build query with multiple search conditions
	query := "SELECT id, username, email, phone, avatar, real_name, status, created_at, updated_at, last_login_at FROM users WHERE deleted_at IS NULL"
	countQuery := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"
	args := make([]any, 0, 3)

	if search != "" {
		query += " AND (username LIKE ? OR email LIKE ? OR real_name LIKE ?)"
		countQuery += " AND (username LIKE ? OR email LIKE ? OR real_name LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	if status != "" {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, status)
	}

	// Get total count
	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	err := ctrl.DB.SQL.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Add pagination
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, perPage, (page-1)*perPage)

	// Execute query
	rows, err := ctrl.DB.SQL.Query(query, args...)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	// Parse results
	type UserListItem struct {
		ID         uint64  `json:"id"`
		Username   string  `json:"username"`
		Email      string  `json:"email"`
		Phone      *string `json:"phone"`
		Avatar     *string `json:"avatar"`
		RealName   *string `json:"real_name"`
		Status     string  `json:"status"`
		CreatedAt  string  `json:"created_at"`
		UpdatedAt  string  `json:"updated_at"`
		LastLoginAt *string `json:"last_login_at"`
	}

	users := make([]UserListItem, 0)
	for rows.Next() {
		var u UserListItem
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Phone, &u.Avatar, &u.RealName, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt)
		if err != nil {
			continue
		}
		users = append(users, u)
	}

	// Calculate total pages
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	c.Success(map[string]any{
		"items":       users,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// Show returns a single user by ID.
func (ctrl *UserController) Show(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.User](
		orm.Query[models.User](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "User not found")
		return
	}
	c.Success(item)
}

// Store creates a new user.
func (ctrl *UserController) Store(c *ghttp.Context) {
	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Sanitize input
	input = middleware.SanitizeMapInput(input)

	validator := ghttp.NewValidator(input, map[string]string{
		"username": "required|min:3|max:50",
		"password": "required|min:8",
		"email":    "required|email",
		"phone":    "phone",
	})
	if errs := validator.Validate(); errs != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]any{
			"code":    422,
			"message": "Validation failed",
			"errors":  errs,
		})
		return
	}

	// Validate password strength
	if password, ok := input["password"].(string); ok {
		if err := utils.ValidatePasswordSimple(password); err != nil {
			c.Error(http.StatusBadRequest, err.Error())
			return
		}
		if utils.IsWeakPassword(password) {
			c.Error(http.StatusBadRequest, "密码过于简单，请使用更复杂的密码")
			return
		}
	}

	item := &models.User{}
	if v, ok := input["username"]; ok {
		if typed, ok := v.(string); ok {
			item.Username = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field username: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["password"]; ok {
		if typed, ok := v.(string); ok {
			hashed, err := bcrypt.GenerateFromPassword([]byte(typed), bcrypt.DefaultCost)
			if err != nil {
				c.Error(http.StatusInternalServerError, "Failed to hash password")
				return
			}
			item.Password = string(hashed)
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field password: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["email"]; ok {
		if typed, ok := v.(string); ok {
			item.Email = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field email: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["phone"]; ok {
		if typed, ok := v.(string); ok {
			item.Phone = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field phone: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["avatar"]; ok {
		if typed, ok := v.(string); ok {
			item.Avatar = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field avatar: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["real_name"]; ok {
		if typed, ok := v.(string); ok {
			item.RealName = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field real_name: expected string, got %T", v))
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
	if v, ok := input["last_login_at"]; ok {
		if typed, ok := v.(string); ok {
			item.LastLoginAt = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field last_login_at: expected time.Time, got %T", v))
			return
		}
	}
	if v, ok := input["role_id"]; ok {
		if typed, ok := v.(int); ok {
			typedVal := typed; item.RoleId = &typedVal
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field role_id: expected int, got %T", v))
			return
		}
	}

	result, err := orm.Create[models.User](ctrl.DB, item)
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

// Update modifies an existing user.
func (ctrl *UserController) Update(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.User](
		orm.Query[models.User](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "User not found")
		return
	}

	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}
	if v, ok := input["username"]; ok {
		if typed, ok := v.(string); ok {
			item.Username = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field username: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["password"]; ok {
		if typed, ok := v.(string); ok {
			hashed, err := bcrypt.GenerateFromPassword([]byte(typed), bcrypt.DefaultCost)
			if err != nil {
				c.Error(http.StatusInternalServerError, "Failed to hash password")
				return
			}
			item.Password = string(hashed)
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field password: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["email"]; ok {
		if typed, ok := v.(string); ok {
			item.Email = typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field email: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["phone"]; ok {
		if typed, ok := v.(string); ok {
			item.Phone = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field phone: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["avatar"]; ok {
		if typed, ok := v.(string); ok {
			item.Avatar = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field avatar: expected string, got %T", v))
			return
		}
	}
	if v, ok := input["real_name"]; ok {
		if typed, ok := v.(string); ok {
			item.RealName = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field real_name: expected string, got %T", v))
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
	if v, ok := input["last_login_at"]; ok {
		if typed, ok := v.(string); ok {
			item.LastLoginAt = &typed
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field last_login_at: expected time.Time, got %T", v))
			return
		}
	}
	if v, ok := input["role_id"]; ok {
		if typed, ok := v.(int); ok {
			typedVal := typed; item.RoleId = &typedVal
		} else {
			c.Error(http.StatusBadRequest, fmt.Sprintf("invalid type for field role_id: expected int, got %T", v))
			return
		}
	}

	if err := orm.Update[models.User](ctrl.DB, item); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.Success(item)
}

// Destroy deletes a user by ID.
func (ctrl *UserController) Destroy(c *ghttp.Context) {
	id := c.ParamInt("id")
	item, err := orm.First[models.User](
		orm.Query[models.User](ctrl.DB).Where("id", "=", id),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		c.Error(http.StatusNotFound, "User not found")
		return
	}

	if err := orm.Delete[models.User](ctrl.DB, item); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.NoContent()
}

// Ensure UserController satisfies the ResourceController interface.
var _ interface {
	Index(c *ghttp.Context)
	Show(c *ghttp.Context)
	Store(c *ghttp.Context)
	Update(c *ghttp.Context)
	Destroy(c *ghttp.Context)
} = (*UserController)(nil)
