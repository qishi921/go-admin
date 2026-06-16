package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	ghttp "github.com/Hlgxz/gai/http"
)

// cacheEntry holds cached data with expiration time.
type cacheEntry[T any] struct {
	value      T
	expiredAt  time.Time
}

// PermissionCache caches user permissions to avoid repeated database queries.
type PermissionCache struct {
	mu         sync.RWMutex
	userPerms  map[int]cacheEntry[map[string]bool] // userID -> set of permission codes
	superAdmin map[int]cacheEntry[bool]            // userID -> is super admin
	ttl        time.Duration                        // cache TTL
	stopCleanup chan struct{}
}

var permCache = &PermissionCache{
	userPerms:   make(map[int]cacheEntry[map[string]bool]),
	superAdmin:  make(map[int]cacheEntry[bool]),
	ttl:         5 * time.Minute, // Cache expires after 5 minutes
	stopCleanup: make(chan struct{}),
}

// init 启动定期清理过期缓存
func init() {
	go permCache.cleanupExpiredCache()
}

// cleanupExpiredCache 定期清理过期的缓存条目
func (pc *PermissionCache) cleanupExpiredCache() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			pc.mu.Lock()
			now := time.Now()
			for uid, entry := range pc.userPerms {
				if now.After(entry.expiredAt) {
					delete(pc.userPerms, uid)
				}
			}
			for uid, entry := range pc.superAdmin {
				if now.After(entry.expiredAt) {
					delete(pc.superAdmin, uid)
				}
			}
			pc.mu.Unlock()
		case <-pc.stopCleanup:
			return
		}
	}
}

// ClearUserCache clears cached permissions for a specific user.
// Call this when user's roles or permissions are changed.
func ClearUserCache(userID int) {
	permCache.mu.Lock()
	defer permCache.mu.Unlock()
	delete(permCache.userPerms, userID)
	delete(permCache.superAdmin, userID)
}

// ClearAllCache clears all cached permissions.
// Call this when roles or permissions are modified globally.
func ClearAllCache() {
	permCache.mu.Lock()
	defer permCache.mu.Unlock()
	permCache.userPerms = make(map[int]cacheEntry[map[string]bool])
	permCache.superAdmin = make(map[int]cacheEntry[bool])
}

// isCacheExpired checks if a cache entry has expired.
func isCacheExpired[T any](entry cacheEntry[T]) bool {
	return time.Now().After(entry.expiredAt)
}

// SetCacheTTL sets the cache TTL duration.
func SetCacheTTL(ttl time.Duration) {
	permCache.mu.Lock()
	defer permCache.mu.Unlock()
	permCache.ttl = ttl
}

// RBACMiddleware returns middleware that checks user permissions.
// It verifies the user has the required permission code to access the endpoint.
// Super admins (role code "super_admin") bypass all permission checks.
func RBACMiddleware(db *orm.DB) ghttp.HandlerFunc {
	return func(c *ghttp.Context) {
		// Get user ID from JWT claims
		claims, ok := c.Get("auth_claims")
		if !ok {
			c.Error(http.StatusUnauthorized, "Unauthorized")
			return
		}

		jwtClaims, ok := claims.(*auth.Claims)
		if !ok {
			c.Error(http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID := int(jwtClaims.UserID)

		// Check if user is super admin (bypass all checks)
		if isSuperAdmin(db, userID) {
			c.Set("is_super_admin", true)
			c.Next()
			return
		}

		// Get required permission for this endpoint
		requiredPerm := getRequiredPermission(c)
		if requiredPerm == "" {
			// No permission required for this endpoint
			c.Next()
			return
		}

		// Check if user has the required permission
		if !hasPermission(db, userID, requiredPerm) {
			c.Error(http.StatusForbidden, "Permission denied: "+requiredPerm)
			return
		}

		c.Next()
	}
}

// isSuperAdmin checks if the user has the super_admin role.
func isSuperAdmin(db *orm.DB, userID int) bool {
	permCache.mu.RLock()
	if entry, ok := permCache.superAdmin[userID]; ok && !isCacheExpired(entry) {
		permCache.mu.RUnlock()
		return entry.value
	}
	permCache.mu.RUnlock()

	// Query database
	query := `
		SELECT r.code FROM roles r
		INNER JOIN role_user ru ON r.id = ru.role_id
		WHERE ru.user_id = ? AND r.code = 'super_admin' AND r.deleted_at IS NULL
	`
	var code string
	err := db.SQL.QueryRow(query, userID).Scan(&code)
	isSuper := err == nil && code == "super_admin"

	// Cache result with TTL
	permCache.mu.Lock()
	permCache.superAdmin[userID] = cacheEntry[bool]{
		value:     isSuper,
		expiredAt: time.Now().Add(permCache.ttl),
	}
	permCache.mu.Unlock()

	return isSuper
}

// hasPermission checks if the user has a specific permission code.
func hasPermission(db *orm.DB, userID int, permCode string) bool {
	permCache.mu.RLock()
	if entry, ok := permCache.userPerms[userID]; ok && !isCacheExpired(entry) {
		permCache.mu.RUnlock()
		return entry.value[permCode]
	}
	permCache.mu.RUnlock()

	// Query all permissions for this user
	query := `
		SELECT p.code FROM permissions p
		INNER JOIN role_permission rp ON p.id = rp.permission_id
		INNER JOIN role_user ru ON rp.role_id = ru.role_id
		WHERE ru.user_id = ? AND p.status = 'active' AND p.deleted_at IS NULL
	`
	rows, err := db.SQL.Query(query, userID)
	if err != nil {
		return false
	}
	defer rows.Close()

	perms := make(map[string]bool)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil {
			perms[code] = true
		}
	}

	// Cache result with TTL
	permCache.mu.Lock()
	permCache.userPerms[userID] = cacheEntry[map[string]bool]{
		value:     perms,
		expiredAt: time.Now().Add(permCache.ttl),
	}
	permCache.mu.Unlock()

	return perms[permCode]
}

// getRequiredPermission determines the required permission code for the current request.
// Mapping: /api/v1/users -> user:manage, user:create, user:edit, user:delete based on method
func getRequiredPermission(c *ghttp.Context) string {
	path := c.Request.URL.Path
	method := c.Request.Method

	// Skip permission check for certain paths
	skipPaths := []string{
		"/api/v1/auth/me",
		"/api/v1/auth/password",
		"/api/v1/logs",
	}
	for _, skip := range skipPaths {
		if strings.HasPrefix(path, skip) {
			return ""
		}
	}

	// Map path patterns to permission prefixes
	pathPerms := map[string]string{
		"/api/v1/users":       "user",
		"/api/v1/roles":       "role",
		"/api/v1/permissions": "permission",
		"/api/v1/menus":       "menu",
	}

	for prefix, permPrefix := range pathPerms {
		if strings.HasPrefix(path, prefix) {
			// Determine specific permission based on method
			switch method {
			case "GET":
				// Read operations need :manage permission
				return permPrefix + ":manage"
			case "POST":
				return permPrefix + ":create"
			case "PUT":
				return permPrefix + ":edit"
			case "DELETE":
				return permPrefix + ":delete"
			}
		}
	}

	// Role-Permission and Role-User management
	if strings.Contains(path, "/permissions") && strings.Contains(path, "/roles/") {
		return "role:manage"
	}
	if strings.Contains(path, "/users") && strings.Contains(path, "/roles/") {
		return "role:manage"
	}
	if strings.Contains(path, "/roles") && strings.Contains(path, "/users/") {
		return "user:manage"
	}

	// Default: require permission based on path
	return ""
}

// GetCurrentUserPermissions returns all permission codes for the current user.
// Useful for frontend to determine which buttons to show.
func GetCurrentUserPermissions(c *ghttp.Context, db *orm.DB) ([]string, error) {
	claims, ok := c.Get("auth_claims")
	if !ok {
		return nil, fmt.Errorf("no auth claims")
	}

	jwtClaims, ok := claims.(*auth.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	userID := int(jwtClaims.UserID)

	// Check cache first
	permCache.mu.RLock()
	if entry, ok := permCache.userPerms[userID]; ok && !isCacheExpired(entry) {
		permCache.mu.RUnlock()
		result := make([]string, 0, len(entry.value))
		for code := range entry.value {
			result = append(result, code)
		}
		return result, nil
	}
	permCache.mu.RUnlock()

	// Query database
	query := `
		SELECT p.code FROM permissions p
		INNER JOIN role_permission rp ON p.id = rp.permission_id
		INNER JOIN role_user ru ON rp.role_id = ru.role_id
		WHERE ru.user_id = ? AND p.status = 'active' AND p.deleted_at IS NULL
	`
	rows, err := db.SQL.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	permsMap := make(map[string]bool)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil {
			permissions = append(permissions, code)
			permsMap[code] = true
		}
	}

	// Cache result with TTL
	permCache.mu.Lock()
	permCache.userPerms[userID] = cacheEntry[map[string]bool]{
		value:     permsMap,
		expiredAt: time.Now().Add(permCache.ttl),
	}
	permCache.mu.Unlock()

	return permissions, nil
}