package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	// 通知表
	Migrations = append(Migrations, migration.Migration{
		Name: "20260612000003_create_notifications_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("notifications", drv)
			b.ID()
			b.Integer("user_id")
			b.String("title", 200)
			b.Text("content")
			b.String("type", 50).SetDefault("'system'")
			b.Integer("priority").SetDefault("0")
			b.Boolean("is_read").SetDefault("false")
			b.DateTime("read_at").SetNullable()
			b.String("channel", 50).SetDefault("'in_app'")
			b.DateTime("sent_at").SetNullable()
			b.String("send_status", 20).SetDefault("'pending'")
			b.Text("error_msg").SetNullable()
			b.Text("metadata").SetNullable()
			b.Timestamps()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("notifications", drv)
			return b.ToDropSQL()
		},
	})

	// 通知模板表
	Migrations = append(Migrations, migration.Migration{
		Name: "20260612000004_create_notification_templates_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("notification_templates", drv)
			b.ID()
			b.String("code", 100).SetUnique()
			b.String("name", 200)
			b.String("title", 200)
			b.Text("content")
			b.String("type", 50).SetDefault("'system'")
			b.String("channels", 200)
			b.Text("variables").SetNullable()
			b.Boolean("is_active").SetDefault("true")
			b.Timestamps()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("notification_templates", drv)
			return b.ToDropSQL()
		},
	})
}
