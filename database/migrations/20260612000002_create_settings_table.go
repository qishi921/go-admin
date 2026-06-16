package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260612000002_create_settings_table",
		Up: func(drv driver.Driver) string {
			return `CREATE TABLE IF NOT EXISTS settings (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				key VARCHAR(100) NOT NULL UNIQUE,
				value VARCHAR(500) NOT NULL,
				type VARCHAR(20) DEFAULT 'string',
				group_name VARCHAR(50) DEFAULT 'system',
				label VARCHAR(100) NOT NULL,
				options VARCHAR(500),
				is_public INTEGER DEFAULT 0,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				deleted_at DATETIME
			);
			CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key);
			CREATE INDEX IF NOT EXISTS idx_settings_group_name ON settings(group_name);

			INSERT OR IGNORE INTO settings (key, value, type, group_name, label, is_public) VALUES
			('site_name', 'Gai Admin', 'string', 'system', '站点名称', 1),
			('site_logo', '/layui/logo.png', 'string', 'system', '站点Logo', 1),
			('site_footer', 'Powered by Gai Framework', 'string', 'system', '页脚文字', 1),
			('login_timeout', '7200', 'number', 'system', '登录有效期(秒)', 0),
			('password_min_length', '6', 'number', 'system', '密码最小长度', 0),
			('upload_max_size', '10', 'number', 'system', '上传文件最大(MB)', 0),
			('allow_register', 'false', 'boolean', 'system', '允许注册', 0),
			('email_enabled', 'false', 'boolean', 'email', '启用邮件通知', 0),
			('email_host', '', 'string', 'email', 'SMTP服务器', 0),
			('email_port', '25', 'number', 'email', 'SMTP端口', 0),
			('email_user', '', 'string', 'email', 'SMTP用户名', 0),
			('email_pass', '', 'string', 'email', 'SMTP密码', 0),
			('email_from', '', 'string', 'email', '发件人地址', 0);`
		},
		Down: func(drv driver.Driver) string {
			return `DROP TABLE IF EXISTS settings;`
		},
	})
}
