package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hlgxz/gai"
	"github.com/Hlgxz/gai/auth"
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
	"github.com/Hlgxz/gai/database/orm"
	ghttp "github.com/Hlgxz/gai/http"
	"github.com/user/admin-system/app/config"
	"github.com/user/admin-system/app/middleware"
	"github.com/user/admin-system/app/scheduler"
	"github.com/user/admin-system/database/migrations"
	"github.com/user/admin-system/routes"
)

func main() {
	app := gai.New()
	app.LoadConfig("config")
	app.UseDefaults()

	// Build and validate configuration
	appCfg := &config.AppConfig{
		Port:      app.Config().GetInt("port", 8080),
		Env:       app.Config().GetString("env", "development"),
		Debug:     app.Config().GetBool("debug", false),
		DBDriver:  app.Config().GetString("database.driver", "sqlite"),
		DBDsn:     app.Config().GetString("database.dsn", "storage/database.db"),
		JWTSecret: app.Config().GetString("auth.guards.jwt.secret", ""),
		JWTTL:     app.Config().GetInt("auth.guards.jwt.ttl", 7200),
		LogLevel:  app.Config().GetString("log.level", "info"),
		LogFormat: app.Config().GetString("log.format", "json"),
		LogFile:   app.Config().GetString("log.file", "storage/logs/app.log"),
	}

	// Validate configuration
	if err := config.ValidateConfig(appCfg); err != nil {
		log.Fatal("Configuration error:\n", err)
	}

	// Initialize structured logger
	logConfig := &middleware.LoggerConfig{
		Level:       appCfg.LogLevel,
		Format:      appCfg.LogFormat,
		OutputFile:  appCfg.LogFile,
		LogRequests: app.Config().GetBool("log.requests", true),
	}
	if err := middleware.InitLogger(logConfig); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer middleware.CloseLogger()
	middleware.Info("Application starting", "env", appCfg.Env)

	// 验证 JWT 密钥安全性
	weakSecrets := []string{"change-me", "change-me-to-a-random-string", "secret", "password", "admin", "test"}
	for _, weak := range weakSecrets {
		if appCfg.JWTSecret == weak || len(appCfg.JWTSecret) < 16 {
			middleware.Error("JWT secret is too weak or too short", "secret_length", len(appCfg.JWTSecret))
			log.Fatal("JWT secret must be at least 16 characters and not a common weak value. Set JWT_SECRET environment variable.")
		}
	}

	// Initialize auth manager with validated secret
	authMgr := auth.NewManager("jwt")
	authMgr.RegisterGuard(auth.NewJWTGuard(appCfg.JWTSecret, appCfg.JWTTL))
	app.Instance("auth", authMgr)

	// Initialize database
	drv, err := driver.Get(appCfg.DBDriver)
	if err != nil {
		middleware.Error("Failed to get database driver", "error", err)
		log.Fatal("Failed to get database driver:", err)
	}
	sqlDB, err := drv.Open(appCfg.DBDsn)
	if err != nil {
		middleware.Error("Failed to connect to database", "error", err)
		log.Fatal("Failed to connect to database:", err)
	}
	db := &orm.DB{
		SQL:        sqlDB,
		DriverName: drv.Name(),
		QuoteIdent: drv.QuoteIdent,
	}
	app.Instance("db", db)

	// Run migrations
	migrator := migration.NewMigrator(sqlDB, drv)
	for _, m := range migrations.Migrations {
		migrator.Add(m)
	}
	if err := migrator.Migrate(); err != nil {
		middleware.Error("Failed to run migrations", "error", err)
		log.Fatal("Failed to run migrations:", err)
	}
	middleware.Info("Migrations completed successfully")

	// Initialize scheduler
	sched := scheduler.NewScheduler(db)
	// Register built-in tasks
	sched.RegisterTask("cleanup_logs", func(params map[string]any) error {
		// 清理超过30天的操作日志
		retentionDays := 30
		if days, ok := params["retention_days"].(int); ok {
			retentionDays = days
		}
		cutoffDate := time.Now().AddDate(0, 0, -retentionDays).Format("2006-01-02 15:04:05")
		result, err := db.SQL.Exec("DELETE FROM operation_logs WHERE created_at < ?", cutoffDate)
		if err != nil {
			middleware.Error("Failed to cleanup logs", "error", err)
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		middleware.Info("Logs cleanup completed", "deleted_rows", rowsAffected)
		return nil
	})
	sched.RegisterTask("cleanup_uploads", func(params map[string]any) error {
		// 清理超过7天的临时上传文件
		retentionDays := 7
		if days, ok := params["retention_days"].(int); ok {
			retentionDays = days
		}
		cutoffDate := time.Now().AddDate(0, 0, -retentionDays).Format("2006-01-02 15:04:05")
		result, err := db.SQL.Exec("DELETE FROM uploads WHERE status = 'temp' AND created_at < ?", cutoffDate)
		if err != nil {
			middleware.Error("Failed to cleanup uploads", "error", err)
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		middleware.Info("Uploads cleanup completed", "deleted_files", rowsAffected)
		return nil
	})
	sched.RegisterTask("cleanup_operation_logs", func(params map[string]any) error {
		// 清理超过90天的操作日志
		retentionDays := 90
		if days, ok := params["retention_days"].(int); ok {
			retentionDays = days
		}
		cutoffDate := time.Now().AddDate(0, 0, -retentionDays).Format("2006-01-02 15:04:05")
		result, err := db.SQL.Exec("DELETE FROM operation_logs WHERE created_at < ?", cutoffDate)
		if err != nil {
			middleware.Error("Failed to cleanup operation logs", "error", err)
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		middleware.Info("Operation logs cleanup completed", "deleted_rows", rowsAffected)
		return nil
	})
	sched.RegisterTask("cleanup_task_executions", func(params map[string]any) error {
		// 清理超过30天的任务执行记录
		retentionDays := 30
		if days, ok := params["retention_days"].(int); ok {
			retentionDays = days
		}
		cutoffDate := time.Now().AddDate(0, 0, -retentionDays).Format("2006-01-02 15:04:05")
		result, err := db.SQL.Exec("DELETE FROM task_executions WHERE started_at < ?", cutoffDate)
		if err != nil {
			middleware.Error("Failed to cleanup task executions", "error", err)
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		middleware.Info("Task executions cleanup completed", "deleted_rows", rowsAffected)
		return nil
	})
	// Start scheduler (runs in background)
	sched.Start()
	defer sched.Stop()
	middleware.Info("Scheduler started")

	// 注册全局中间件
	app.Router().Use(middleware.RequestIDMiddleware())
	app.Router().Use(middleware.RecoveryMiddleware())

	// 健康检查端点（带数据库连接验证）
	app.Router().Get("/health", func(c *ghttp.Context) {
		ctx := context.Background()

		// 检查数据库连接
		dbErr := db.SQL.PingContext(ctx)
		if dbErr != nil {
			middleware.Error("Health check failed: database connection error", "error", dbErr)
			c.Error(http.StatusServiceUnavailable, "database unavailable")
			return
		}

		c.Success(map[string]any{
			"status":   "ok",
			"database": "connected",
		})
	})

	// Register routes (includes static files)
	routes.Register(app)
	routes.RegisterAuthRoutes(app.Router(), db, authMgr)
	routes.RegisterUserRoutes(app.Router(), db, authMgr)
	routes.RegisterRoleRoutes(app.Router(), db, authMgr)
	routes.RegisterPermissionRoutes(app.Router(), db, authMgr)
	routes.RegisterMenuRoutes(app.Router(), db, authMgr)
	routes.RegisterOperationLogRoutes(app.Router(), db, authMgr)
	routes.RegisterRoleUserRoutes(app.Router(), db, authMgr)
	routes.RegisterRolePermissionRoutes(app.Router(), db, authMgr)
	routes.RegisterUploadRoutes(app.Router(), db, authMgr, "storage/uploads")
	routes.RegisterSettingRoutes(app.Router(), db, authMgr)
		routes.RegisterDashboardRoutes(app.Router(), db, authMgr)
	// P2 features
	routes.RegisterNotificationRoutes(app.Router(), db, authMgr)
	routes.RegisterScheduledTaskRoutes(app.Router(), db, authMgr, sched)
	routes.RegisterExportRoutes(app.Router(), db, authMgr)
	routes.RegisterAuditRoutes(app.Router(), db, authMgr)

	addr := fmt.Sprintf(":%d", appCfg.Port)
	middleware.Info("Server starting", "addr", addr)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Router(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器（非阻塞）
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed: ", err)
		}
	}()

	middleware.Info("Server started successfully", "addr", addr)

	// 等待中断信号进行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	middleware.Info("Shutting down server...")

	// 给正在处理的请求最多30秒完成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 停止定时任务调度器
	sched.Stop()

	if err := srv.Shutdown(ctx); err != nil {
		middleware.Error("Server forced to shutdown", "error", err)
	}

	middleware.Info("Server exited")
}
