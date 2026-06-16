package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// SettingController handles system settings.
type SettingController struct {
	DB *orm.DB
}

// NewSettingController creates a new setting controller.
func NewSettingController(db *orm.DB) *SettingController {
	return &SettingController{DB: db}
}

// Cache for settings
var (
	settingCache     = make(map[string]string)
	settingCacheMu   sync.RWMutex
	settingCacheInit = false
)

// GetSetting returns a setting value by key.
func GetSetting(db *orm.DB, key string, defaultValue string) string {
	// Check cache first
	settingCacheMu.RLock()
	if v, ok := settingCache[key]; ok {
		settingCacheMu.RUnlock()
		return v
	}
	settingCacheMu.RUnlock()

	// Query database
	var value string
	err := db.SQL.QueryRow(
		"SELECT value FROM settings WHERE key = ? AND deleted_at IS NULL LIMIT 1",
		key,
	).Scan(&value)

	if err != nil {
		return defaultValue
	}

	// Cache it
	settingCacheMu.Lock()
	settingCache[key] = value
	settingCacheMu.Unlock()

	return value
}

// GetSettingBool returns a setting as boolean.
func GetSettingBool(db *orm.DB, key string, defaultValue bool) bool {
	val := GetSetting(db, key, strconv.FormatBool(defaultValue))
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return b
}

// GetSettingInt returns a setting as integer.
func GetSettingInt(db *orm.DB, key string, defaultValue int) int {
	val := GetSetting(db, key, strconv.Itoa(defaultValue))
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

// ClearSettingCache clears the setting cache.
func ClearSettingCache() {
	settingCacheMu.Lock()
	settingCache = make(map[string]string)
	settingCacheMu.Unlock()
}

// List returns all settings by group.
func (ctrl *SettingController) List(c *ghttp.Context) {
	group := c.Query("group")

	q := orm.Query[models.Setting](ctrl.DB)
	if group != "" {
		q = q.Where("group_name", "=", group)
	}
	q = q.OrderBy("group_name", "ASC").OrderBy("key", "ASC")

	settings, err := orm.Get[models.Setting](q)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.Success(map[string]any{
		"items": settings,
		"total": len(settings),
	})
}

// Public returns public settings (no auth required).
func (ctrl *SettingController) Public(c *ghttp.Context) {
	settings, err := orm.Get[models.Setting](
		orm.Query[models.Setting](ctrl.DB).Where("is_public", "=", true),
	)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	// Return as key-value map
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	c.Success(result)
}

// Get returns a single setting.
func (ctrl *SettingController) Get(c *ghttp.Context) {
	key := c.Param("key")

	setting, err := orm.First[models.Setting](
		orm.Query[models.Setting](ctrl.DB).Where("key", "=", key),
	)
	if err != nil || setting == nil {
		c.Error(http.StatusNotFound, "Setting not found")
		return
	}

	c.Success(setting)
}

// Update updates a setting.
func (ctrl *SettingController) Update(c *ghttp.Context) {
	key := c.Param("key")

	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}

	setting, err := orm.First[models.Setting](
		orm.Query[models.Setting](ctrl.DB).Where("key", "=", key),
	)
	if err != nil || setting == nil {
		c.Error(http.StatusNotFound, "Setting not found")
		return
	}

	if v, ok := input["value"].(string); ok {
		setting.Value = v
	}
	if v, ok := input["label"].(string); ok {
		setting.Label = v
	}

	if err := orm.Update[models.Setting](ctrl.DB, setting); err != nil {
		c.Error(http.StatusInternalServerError, "Failed to update setting")
		return
	}

	// Clear cache
	ClearSettingCache()

	c.Success(setting)
}

// BatchUpdate updates multiple settings at once.
func (ctrl *SettingController) BatchUpdate(c *ghttp.Context) {
	var input map[string]string
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}

	for key, value := range input {
		ctrl.DB.SQL.Exec(
			"UPDATE settings SET value = ?, updated_at = datetime('now') WHERE key = ?",
			value, key,
		)
	}

	// Clear cache
	ClearSettingCache()

	c.Success(map[string]string{"message": "Settings updated"})
}

// Create creates a new setting.
func (ctrl *SettingController) Create(c *ghttp.Context) {
	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate required fields
	key, _ := input["key"].(string)
	value, _ := input["value"].(string)
	if key == "" {
		c.Error(http.StatusBadRequest, "Key is required")
		return
	}

	// Check if exists
	count := 0
	ctrl.DB.SQL.QueryRow("SELECT COUNT(*) FROM settings WHERE key = ?", key).Scan(&count)
	if count > 0 {
		c.Error(http.StatusConflict, "Setting key already exists")
		return
	}

	setting := &models.Setting{
		Key:       key,
		Value:     value,
		Type:      getStringOr(input, "type", "string"),
		GroupName: getStringOr(input, "group", "system"),
		Label:     getStringOr(input, "label", key),
		Options:   ptrString(getStringOr(input, "options", "")),
		IsPublic:  getBoolOr(input, "is_public", false),
	}

	result, err := orm.Create[models.Setting](ctrl.DB, setting)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to create setting")
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"code":    0,
		"message": "ok",
		"data":    result,
	})
}

// Delete deletes a setting.
func (ctrl *SettingController) Delete(c *ghttp.Context) {
	key := c.Param("key")

	setting, err := orm.First[models.Setting](
		orm.Query[models.Setting](ctrl.DB).Where("key", "=", key),
	)
	if err != nil || setting == nil {
		c.Error(http.StatusNotFound, "Setting not found")
		return
	}

	if err := orm.Delete[models.Setting](ctrl.DB, setting); err != nil {
		c.Error(http.StatusInternalServerError, "Failed to delete setting")
		return
	}

	ClearSettingCache()
	c.NoContent()
}

// Helper functions
func getStringOr(m map[string]any, key string, def string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return def
}

func getBoolOr(m map[string]any, key string, def bool) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	// Try string "true"
	if v, ok := m[key].(string); ok {
		b, _ := strconv.ParseBool(v)
		return b
	}
	return def
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// GetSettingsJSON returns settings as a nested map by group.
func GetSettingsJSON(db *orm.DB) map[string]map[string]string {
	settings, err := orm.Get[models.Setting](
		orm.Query[models.Setting](db).OrderBy("group_name", "ASC"),
	)
	if err != nil {
		return nil
	}

	result := make(map[string]map[string]string)
	for _, s := range settings {
		if result[s.GroupName] == nil {
			result[s.GroupName] = make(map[string]string)
		}
		result[s.GroupName][s.Key] = s.Value
	}
	return result
}

// MarshalJSON helps with JSON encoding
func init() {
	_ = json.Marshal
}
