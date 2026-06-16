package routes

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hlgxz/gai"
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/orm"
	"github.com/Hlgxz/gai/router"
	ghttp "github.com/Hlgxz/gai/http"
	"github.com/user/admin-system/app/controllers"
	"github.com/user/admin-system/app/middleware"
)

// Register sets up all application routes.
func Register(app *gai.Application) {
	r := app.Router()

	// Apply security headers to all routes
	r.Use(middleware.SecurityHeadersMiddleware())

	// Serve static files (frontend)
	webDir := filepath.Join(app.BasePath(), "web")

	// Serve login page at /login
	r.Get("/login", func(c *ghttp.Context) {
		fullPath := filepath.Join(webDir, "login.html")
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.Error(http.StatusNotFound, "Not Found")
			return
		}
		c.SetHeader("Content-Type", "text/html; charset=utf-8")
		content, _ := os.ReadFile(fullPath)
		c.Writer.Write(content)
	})

	// 首页重定向到多页面版本
	r.Get("/", func(c *ghttp.Context) {
		c.Redirect(302, "/pages/dashboard/index.html")
	})

	// Serve pages 目录下的静态 HTML 文件
	// 使用 :dir/:file 格式匹配两段路径
	r.Get("/pages/:dir/:file", func(c *ghttp.Context) {
		dir := c.Param("dir")
		file := c.Param("file")
		pathParam := dir + "/" + file
		fullPath := filepath.Join(webDir, "pages", pathParam)
		if !strings.HasPrefix(fullPath, webDir) {
			c.Error(http.StatusForbidden, "Forbidden")
			return
		}
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.Error(http.StatusNotFound, "Not Found")
			return
		}
		c.SetHeader("Content-Type", "text/html; charset=utf-8")
		content, _ := os.ReadFile(fullPath)
		c.Writer.Write(content)
	})

	// Serve Layui static files
	r.Get("/layui/:file", func(c *ghttp.Context) {
		filePath := c.Param("file")
		fullPath := filepath.Join(webDir, "layui", filePath)
		if !strings.HasPrefix(fullPath, webDir) {
			c.Error(http.StatusForbidden, "Forbidden")
			return
		}
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.Error(http.StatusNotFound, "Not Found")
			return
		}
		ext := strings.ToLower(filepath.Ext(fullPath))
		switch ext {
		case ".css":
			c.SetHeader("Content-Type", "text/css; charset=utf-8")
		case ".js":
			c.SetHeader("Content-Type", "application/javascript; charset=utf-8")
		case ".ttf", ".woff", ".woff2", ".eot":
			c.SetHeader("Content-Type", "application/font-"+ext[1:])
		default:
			c.SetHeader("Content-Type", mime.TypeByExtension(ext))
		}
		c.Writer.WriteHeader(http.StatusOK)
		content, _ := os.ReadFile(fullPath)
		c.Writer.Write(content)
	})

	// Serve Layui subdirectory files (css/*.css, font/*.woff2, etc.)
	r.Get("/layui/:dir/:file", func(c *ghttp.Context) {
		dir := c.Param("dir")
		filePath := c.Param("file")
		fullPath := filepath.Join(webDir, "layui", dir, filePath)
		if !strings.HasPrefix(fullPath, webDir) {
			c.Error(http.StatusForbidden, "Forbidden")
			return
		}
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.Error(http.StatusNotFound, "Not Found")
			return
		}
		ext := strings.ToLower(filepath.Ext(fullPath))
		switch ext {
		case ".css":
			c.SetHeader("Content-Type", "text/css; charset=utf-8")
		case ".js":
			c.SetHeader("Content-Type", "application/javascript; charset=utf-8")
		case ".ttf", ".woff", ".woff2", ".eot":
			c.SetHeader("Content-Type", "application/font-"+ext[1:])
		default:
			c.SetHeader("Content-Type", mime.TypeByExtension(ext))
		}
		c.Writer.WriteHeader(http.StatusOK)
		content, _ := os.ReadFile(fullPath)
		c.Writer.Write(content)
	})

	// Serve static files (css, js) with proper MIME types
	r.Get("/css/:file", func(c *ghttp.Context) {
		filePath := c.Param("file")
		fullPath := filepath.Join(webDir, "css", filePath)
		if !strings.HasPrefix(fullPath, webDir) {
			c.Error(http.StatusForbidden, "Forbidden")
			return
		}
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.Error(http.StatusNotFound, "Not Found")
			return
		}
		c.SetHeader("Content-Type", "text/css; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		content, _ := os.ReadFile(fullPath)
		c.Writer.Write(content)
	})

	r.Get("/js/:file", func(c *ghttp.Context) {
		filePath := c.Param("file")
		fullPath := filepath.Join(webDir, "js", filePath)
		if !strings.HasPrefix(fullPath, webDir) {
			c.Error(http.StatusForbidden, "Forbidden")
			return
		}
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.Error(http.StatusNotFound, "Not Found")
			return
		}
		c.SetHeader("Content-Type", "application/javascript; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		content, _ := os.ReadFile(fullPath)
		c.Writer.Write(content)
	})

	r.Get("/web/:file", func(c *ghttp.Context) {
		filePath := c.Param("file")
		if filePath == "" {
			c.File(filepath.Join(webDir, "index.html"))
			return
		}
		fullPath := filepath.Join(webDir, filePath)
		if !strings.HasPrefix(fullPath, webDir) {
			c.Error(http.StatusForbidden, "Forbidden")
			return
		}
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.File(filepath.Join(webDir, "index.html"))
			return
		}
		ext := strings.ToLower(filepath.Ext(fullPath))
		switch ext {
		case ".css":
			c.SetHeader("Content-Type", "text/css; charset=utf-8")
		case ".js":
			c.SetHeader("Content-Type", "application/javascript; charset=utf-8")
		case ".html", ".htm":
			c.SetHeader("Content-Type", "text/html; charset=utf-8")
		case ".json":
			c.SetHeader("Content-Type", "application/json; charset=utf-8")
		case ".png":
			c.SetHeader("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			c.SetHeader("Content-Type", "image/jpeg")
		case ".gif":
			c.SetHeader("Content-Type", "image/gif")
		case ".svg":
			c.SetHeader("Content-Type", "image/svg+xml")
		case ".ico":
			c.SetHeader("Content-Type", "image/x-icon")
		default:
			c.SetHeader("Content-Type", mime.TypeByExtension(ext))
		}
		c.Writer.WriteHeader(http.StatusOK)
		content, _ := os.ReadFile(fullPath)
		c.Writer.Write(content)
	})

	// API Documentation
	r.Get("/docs", func(c *ghttp.Context) {
		c.File(filepath.Join(app.BasePath(), "docs", "API.md"))
	})
}

// RegisterAuthRoutes sets up authentication routes.
func RegisterAuthRoutes(r *router.Router, db *orm.DB, authMgr *auth.Manager) {
	authController := controllers.NewAuthController(db, authMgr)
	permCheckController := controllers.NewPermissionCheckController(db)

	// Login endpoint with rate limiting
	loginWithRateLimit := func(c *ghttp.Context) {
		ip := c.ClientIP()
		if middleware.IsLoginBlocked(ip) {
			c.Error(http.StatusTooManyRequests, "登录尝试次数过多，请稍后再试")
			return
		}
		authController.Login(c)
	}

	r.Post("/api/v1/auth/login", loginWithRateLimit)

	// 注册接口使用严格的速率限制
	registerWithRateLimit := func(c *ghttp.Context) {
		ip := c.ClientIP()
		if middleware.IsLoginBlocked(ip) {
			c.Error(http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			return
		}
		authController.Register(c)
	}

	r.Group("/api/v1/auth", func(g *router.Group) {
		g.Post("/register", registerWithRateLimit)
		g.Post("/logout", authController.Logout)
		g.Use(authMgr.Middleware("jwt"))
		g.Get("/me", authController.Me)
		g.Put("/me", authController.UpdateProfile)
		g.Put("/password", authController.UpdatePassword)
		g.Get("/permissions", permCheckController.MyPermissions)
	})
}
