package controllers

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/middleware"
	"github.com/user/admin-system/app/models"
)

// UploadController handles file upload operations.
type UploadController struct {
	DB         *orm.DB
	UploadDir  string
	MaxSize    int64 // in bytes
	AllowTypes []string
}

// NewUploadController creates a new upload controller.
func NewUploadController(db *orm.DB, uploadDir string) *UploadController {
	// Ensure upload directory exists
	os.MkdirAll(uploadDir, 0755)

	return &UploadController{
		DB:        db,
		UploadDir: uploadDir,
		MaxSize:   10 * 1024 * 1024, // 10MB default
		AllowTypes: []string{
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"text/plain", "text/csv",
		},
	}
}

// Upload handles single file upload.
func (ctrl *UploadController) Upload(c *ghttp.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > ctrl.MaxSize {
		c.Error(http.StatusBadRequest, fmt.Sprintf("File too large, max %dMB", ctrl.MaxSize/1024/1024))
		return
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !ctrl.isAllowedType(contentType) {
		c.Error(http.StatusBadRequest, "File type not allowed")
		return
	}

	// Sanitize filename
	filename := middleware.SanitizeFilename(header.Filename)
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".bin"
	}

	// Generate unique filename
	timestamp := time.Now().Format("20060102150405")
	randomStr := randomString(6)
	newFilename := fmt.Sprintf("%s_%s%s", timestamp, randomStr, ext)

	// Get module from form (optional)
	module := middleware.SanitizeInput(c.Request.FormValue("module"))
	if module == "" {
		module = "general"
	}

	// Create module directory
	moduleDir := filepath.Join(ctrl.UploadDir, module)
	os.MkdirAll(moduleDir, 0755)

	// Create dated subdirectory
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(moduleDir, dateDir)
	os.MkdirAll(fullDir, 0755)

	// Full file path
	filePath := filepath.Join(fullDir, newFilename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		c.Error(http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Get user ID from context
	var userId *int
	if userIDVal, ok := c.Get("auth_user_id"); ok {
		if id, ok := userIDVal.(uint64); ok {
			uid := int(id)
			userId = &uid
		}
	}

	// Save upload record
	upload := &models.Upload{
		FileName:      newFilename,
		OriginalName:  filename,
		FilePath:      filepath.Join(module, dateDir, newFilename),
		FileSize:      header.Size,
		MimeType:      contentType,
		Extension:     ext,
		UserId:        userId,
		Module:        module,
		Status:        "active",
	}

	result, err := orm.Create[models.Upload](ctrl.DB, upload)
	if err != nil {
		os.Remove(filePath)
		c.Error(http.StatusInternalServerError, "Failed to save upload record")
		return
	}

	c.Success(map[string]any{
		"id":            result.ID,
		"file_name":     result.FileName,
		"original_name": result.OriginalName,
		"file_path":     result.FilePath,
		"file_size":     result.FileSize,
		"mime_type":     result.MimeType,
		"url":           "/uploads/" + result.FilePath,
	})
}

// Delete soft deletes an upload.
func (ctrl *UploadController) Delete(c *ghttp.Context) {
	id := c.ParamInt("id")

	upload, err := orm.First[models.Upload](
		orm.Query[models.Upload](ctrl.DB).Where("id", "=", id),
	)
	if err != nil || upload == nil {
		c.Error(http.StatusNotFound, "Upload not found")
		return
	}

	if err := orm.Delete[models.Upload](ctrl.DB, upload); err != nil {
		c.Error(http.StatusInternalServerError, "Failed to delete upload")
		return
	}

	c.NoContent()
}

// List returns uploads list.
func (ctrl *UploadController) List(c *ghttp.Context) {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	module := c.Query("module")

	q := orm.Query[models.Upload](ctrl.DB)
	if module != "" {
		q = q.Where("module", "=", module)
	}
	q = q.OrderBy("created_at", "DESC")

	result, err := orm.Paginate[models.Upload](q, page, perPage)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.Success(result)
}

// Serve serves uploaded files.
func (ctrl *UploadController) Serve(c *ghttp.Context) {
	path := c.Param("path")
	if path == "" {
		c.Error(http.StatusNotFound, "File not found")
		return
	}

	// Security: prevent directory traversal
	if strings.Contains(path, "..") {
		c.Error(http.StatusForbidden, "Access denied")
		return
	}

	filePath := filepath.Join(ctrl.UploadDir, path)

	// Check file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.Error(http.StatusNotFound, "File not found")
		return
	}

	// Set content type based on extension
	ext := strings.ToLower(filepath.Ext(filePath))
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	case ".pdf":
		contentType = "application/pdf"
	case ".txt":
		contentType = "text/plain"
	case ".csv":
		contentType = "text/csv"
	}

	c.SetHeader("Content-Type", contentType)
	c.SetHeader("Cache-Control", "public, max-age=31536000")
	c.File(filePath)
}

// isAllowedType checks if the content type is allowed.
func (ctrl *UploadController) isAllowedType(contentType string) bool {
	for _, t := range ctrl.AllowTypes {
		if t == contentType {
			return true
		}
	}
	return false
}

// randomString generates a cryptographically secure random string of given length.
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	randBytes := make([]byte, n)
	_, err := rand.Read(randBytes)
	if err != nil {
		// 回退到时间戳作为后备方案
		for i := range b {
			b[i] = letters[time.Now().Nanosecond()%len(letters)]
		}
		return string(b)
	}
	for i := range b {
		b[i] = letters[int(randBytes[i])%len(letters)]
	}
	return string(b)
}
