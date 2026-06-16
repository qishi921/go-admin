package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260612000007_seed_settings_data",
		Up: func(drv driver.Driver) string {
			return `INSERT OR IGNORE INTO settings (key, value, type, group_name, label, is_public) VALUES
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
			return `DELETE FROM settings WHERE key IN ('site_name', 'site_logo', 'site_footer', 'login_timeout', 'password_min_length', 'upload_max_size', 'allow_register', 'email_enabled', 'email_host', 'email_port', 'email_user', 'email_pass', 'email_from');`
		},
	})
}
