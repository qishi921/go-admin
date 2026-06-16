package controllers

import (
	"net/http"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// RoleUserController handles role-user association operations.
type RoleUserController struct {
	DB *orm.DB
}

// NewRoleUserController creates a new controller instance.
func NewRoleUserController(db *orm.DB) *RoleUserController {
	return &RoleUserController{DB: db}
}

// GetRoleUsers returns the list of users assigned to a role.
func (ctrl *RoleUserController) GetRoleUsers(c *ghttp.Context) {
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

	// Find user IDs from role_user table
	var roleUsers []models.RoleUser
	roleUsers, err = orm.Get[models.RoleUser](
		orm.Query[models.RoleUser](ctrl.DB).Where("role_id", "=", roleID),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	if len(roleUsers) == 0 {
		c.Success([]models.User{})
		return
	}

	userIDs := make([]any, len(roleUsers))
	for i, ru := range roleUsers {
		userIDs[i] = ru.UserID
	}

	users, err := orm.Get[models.User](
		orm.Query[models.User](ctrl.DB).WhereIn("id", userIDs),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.Success(users)
}

// AssignUserToRole assigns a user to a role.
func (ctrl *RoleUserController) AssignUserToRole(c *ghttp.Context) {
	roleID := c.ParamInt("id")

	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}

	userIDVal, ok := input["user_id"]
	if !ok {
		c.Error(http.StatusBadRequest, "user_id is required")
		return
	}

	userID, ok := userIDVal.(float64)
	if !ok {
		c.Error(http.StatusBadRequest, "user_id must be a number")
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

	// Verify user exists
	user, err := orm.First[models.User](
		orm.Query[models.User](ctrl.DB).Where("id", "=", int(userID)),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		c.Error(http.StatusNotFound, "User not found")
		return
	}

	// Check if already assigned
	existing, err := orm.First[models.RoleUser](
		orm.Query[models.RoleUser](ctrl.DB).
			Where("role_id", "=", roleID).
			Where("user_id", "=", int(userID)),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if existing != nil {
		c.Error(http.StatusConflict, "User is already assigned to this role")
		return
	}

	// Create association
	roleUser := &models.RoleUser{
		RoleID: roleID,
		UserID: int(userID),
	}

	_, err = orm.Create[models.RoleUser](ctrl.DB, roleUser)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"code":    0,
		"message": "User assigned to role successfully",
		"data":    roleUser,
	})
}

// RemoveUserFromRole removes a user from a role.
func (ctrl *RoleUserController) RemoveUserFromRole(c *ghttp.Context) {
	roleID := c.ParamInt("id")
	userID := c.ParamInt("userId")

	// Find the association
	roleUser, err := orm.First[models.RoleUser](
		orm.Query[models.RoleUser](ctrl.DB).
			Where("role_id", "=", roleID).
			Where("user_id", "=", userID),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if roleUser == nil {
		c.Error(http.StatusNotFound, "Association not found")
		return
	}

	// Delete the association
	if err := orm.Delete[models.RoleUser](ctrl.DB, roleUser); err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.NoContent()
}

// GetUserRoles returns the list of roles assigned to a user.
func (ctrl *RoleUserController) GetUserRoles(c *ghttp.Context) {
	userID := c.ParamInt("id")

	// Verify user exists
	user, err := orm.First[models.User](
		orm.Query[models.User](ctrl.DB).Where("id", "=", userID),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		c.Error(http.StatusNotFound, "User not found")
		return
	}

	// Find role IDs from role_user table
	var roleUsers []models.RoleUser
	roleUsers, err = orm.Get[models.RoleUser](
		orm.Query[models.RoleUser](ctrl.DB).Where("user_id", "=", userID),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	if len(roleUsers) == 0 {
		c.Success([]models.Role{})
		return
	}

	roleIDs := make([]any, len(roleUsers))
	for i, ru := range roleUsers {
		roleIDs[i] = ru.RoleID
	}

	roles, err := orm.Get[models.Role](
		orm.Query[models.Role](ctrl.DB).WhereIn("id", roleIDs),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.Success(roles)
}
