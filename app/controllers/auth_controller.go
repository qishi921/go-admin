package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	ghttp "github.com/Hlgxz/gai/http"
	"github.com/user/admin-system/app/middleware"
	"github.com/user/admin-system/app/models"
	"github.com/user/admin-system/app/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthController handles authentication operations.
type AuthController struct {
	DB      *orm.DB
	AuthMgr *auth.Manager
}

// NewAuthController creates a new auth controller.
func NewAuthController(db *orm.DB, authMgr *auth.Manager) *AuthController {
	return &AuthController{DB: db, AuthMgr: authMgr}
}

// LoginRequest represents login credentials.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents registration data.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// UserResponse represents a user without sensitive fields.
type UserResponse struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	RealName  string `json:"real_name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// Login handles user login.
func (ctrl *AuthController) Login(c *ghttp.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		c.Error(http.StatusBadRequest, "Username and password are required")
		return
	}

	// Find user by username using raw SQL
	var user models.User
	var createdAtStr string
	var phone, avatar sql.NullString
	err := ctrl.DB.SQL.QueryRow(
		"SELECT id, username, password, email, phone, avatar, status, created_at FROM users WHERE username = ? AND deleted_at IS NULL LIMIT 1",
		req.Username,
	).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &phone,
		&avatar, &user.Status, &createdAtStr,
	)

	if err != nil {
		middleware.RecordLoginFailure(c.ClientIP())
		c.Error(http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if phone.Valid {
		user.Phone = &phone.String
	}
	if avatar.Valid {
		user.Avatar = &avatar.String
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		middleware.RecordLoginFailure(c.ClientIP())
		c.Error(http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// Check user status
	if user.Status != "active" {
		middleware.RecordLoginFailure(c.ClientIP())
		c.Error(http.StatusForbidden, "Account is disabled")
		return
	}

	// Clear login attempts on successful login
	middleware.RecordLoginSuccess(c.ClientIP())

	// Generate JWT token
	guard, err := ctrl.AuthMgr.Guard("jwt")
	if err != nil {
		c.Error(http.StatusInternalServerError, "Auth configuration error")
		return
	}

	jwtGuard := guard.(*auth.JWTGuard)
	token, err := jwtGuard.IssueToken(user.ID, map[string]any{
		"username": user.Username,
		"email":    user.Email,
	})
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Update last login time
	now := time.Now().Format("2006-01-02 15:04:05")
	ctrl.DB.SQL.Exec("UPDATE users SET last_login_at = ? WHERE id = ?", now, user.ID)

	// Build response with pointer handling
	respUser := map[string]any{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
	}
	if user.Phone != nil {
		respUser["phone"] = *user.Phone
	}
	if user.Avatar != nil {
		respUser["avatar"] = *user.Avatar
	}

	c.Success(map[string]any{
		"token": token,
		"user":  respUser,
	})
}

// Register handles user registration.
func (ctrl *AuthController) Register(c *ghttp.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" || req.Email == "" {
		c.Error(http.StatusBadRequest, "Username, password and email are required")
		return
	}

	// Validate password strength
	if err := utils.ValidatePasswordSimple(req.Password); err != nil {
		c.Error(http.StatusBadRequest, err.Error())
		return
	}

	// Check for weak passwords
	if utils.IsWeakPassword(req.Password) {
		c.Error(http.StatusBadRequest, "密码过于简单，请使用更复杂的密码")
		return
	}

	// Check if username already exists
	var count int
	err := ctrl.DB.SQL.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Database error")
		return
	}
	if count > 0 {
		c.Error(http.StatusConflict, "Username already exists")
		return
	}

	// Check if email already exists
	err = ctrl.DB.SQL.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Database error")
		return
	}
	if count > 0 {
		c.Error(http.StatusConflict, "Email already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user using raw SQL
	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := ctrl.DB.SQL.Exec(
		"INSERT INTO users (username, password, email, phone, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		req.Username, string(hashedPassword), req.Email, req.Phone, "active", now, now,
	)
	if err != nil {
		middleware.LogError("Failed to create user", err, "username", req.Username)
		c.Error(http.StatusInternalServerError, "用户创建失败")
		return
	}

	lastID, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, map[string]any{
		"code":    0,
		"message": "User registered successfully",
		"data": map[string]any{
			"id":       lastID,
			"username": req.Username,
			"email":    req.Email,
		},
	})
}

// Logout handles user logout.
func (ctrl *AuthController) Logout(c *ghttp.Context) {
	c.Success(map[string]string{
		"message": "Logged out successfully",
	})
}

// UpdatePasswordRequest represents password update data.
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UpdatePassword handles user password change.
func (ctrl *AuthController) UpdatePassword(c *ghttp.Context) {
	// Get user ID from context (set by auth middleware)
	userIDVal, ok := c.Get("auth_user_id")
	if !ok {
		c.Error(http.StatusUnauthorized, "Unauthorized")
		return
	}

	var userID uint64
	if id, ok := userIDVal.(uint64); ok {
		userID = id
	} else {
		c.Error(http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req UpdatePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		c.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		c.Error(http.StatusBadRequest, "Old password and new password are required")
		return
	}

	// Validate new password strength
	if err := utils.ValidatePasswordSimple(req.NewPassword); err != nil {
		c.Error(http.StatusBadRequest, err.Error())
		return
	}

	// Check for weak passwords
	if utils.IsWeakPassword(req.NewPassword) {
		c.Error(http.StatusBadRequest, "新密码过于简单，请使用更复杂的密码")
		return
	}

	// Find user
	var password string
	err := ctrl.DB.SQL.QueryRow(
		"SELECT password FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1",
		userID,
	).Scan(&password)

	if err != nil {
		c.Error(http.StatusNotFound, "User not found")
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.OldPassword)); err != nil {
		c.Error(http.StatusUnauthorized, "Old password is incorrect")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update password
	_, err = ctrl.DB.SQL.Exec("UPDATE users SET password = ? WHERE id = ?", string(hashedPassword), userID)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to update password")
		return
	}

	c.Success(map[string]string{
		"message": "Password updated successfully",
	})
}

// Me returns current user info.
func (ctrl *AuthController) Me(c *ghttp.Context) {
	// Get user ID from context (set by auth middleware)
	userIDVal, ok := c.Get("auth_user_id")
	if !ok {
		c.Error(http.StatusUnauthorized, "Unauthorized")
		return
	}

	var userID uint64
	if id, ok := userIDVal.(uint64); ok {
		userID = id
	} else {
		c.Error(http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var user models.User
	var createdAtStr string
	var phone, avatar, realName sql.NullString
	err := ctrl.DB.SQL.QueryRow(
		"SELECT id, username, email, phone, avatar, real_name, status, created_at FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1",
		userID,
	).Scan(
		&user.ID, &user.Username, &user.Email, &phone,
		&avatar, &realName, &user.Status, &createdAtStr,
	)

	if err != nil {
		c.Error(http.StatusNotFound, "User not found")
		return
	}

	resp := map[string]any{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
	}
	if phone.Valid {
		resp["phone"] = phone.String
	}
	if avatar.Valid {
		resp["avatar"] = avatar.String
	}
	if realName.Valid {
		resp["real_name"] = realName.String
	}

	c.Success(resp)
}

// UpdateProfile handles user profile update.
func (ctrl *AuthController) UpdateProfile(c *ghttp.Context) {
	userIDVal, ok := c.Get("auth_user_id")
	if !ok {
		c.Error(http.StatusUnauthorized, "Unauthorized")
		return
	}

	var userID uint64
	if id, ok := userIDVal.(uint64); ok {
		userID = id
	} else {
		c.Error(http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	// Build update query dynamically
	updates := make(map[string]any)
	if v, ok := input["email"].(string); ok && v != "" {
		updates["email"] = v
	}
	if v, ok := input["phone"].(string); ok {
		updates["phone"] = v
	}
	if v, ok := input["real_name"].(string); ok {
		updates["real_name"] = v
	}

	if len(updates) == 0 {
		c.Error(http.StatusBadRequest, "No fields to update")
		return
	}

	// Build SQL
	sql := "UPDATE users SET "
	args := make([]any, 0, len(updates)+2)
	first := true
	for k, v := range updates {
		if !first {
			sql += ", "
		}
		sql += k + " = ?"
		args = append(args, v)
		first = false
	}
	sql += ", updated_at = ? WHERE id = ?"
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))
	args = append(args, userID)

	_, err := ctrl.DB.SQL.Exec(sql, args...)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to update profile")
		return
	}

	c.Success(map[string]string{"message": "Profile updated"})
}