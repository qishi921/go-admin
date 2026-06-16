package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260608200856_create_operation_logs_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("operation_logs", drv)
			b.ID()
			b.Integer("user_id").SetNullable()
			b.String("username", 50).SetNullable()
			b.String("action", 50)
			b.String("method", 10).SetNullable()
			b.String("path", 255).SetNullable()
			b.String("ip", 45).SetNullable()
			b.String("user_agent", 500).SetNullable()
			b.Text("params").SetNullable()
			b.Text("result").SetNullable()
			b.Integer("duration").SetDefault("0")
			b.String("status", 50).SetDefault("'success'")
			// 数据变更快照字段
			b.String("resource_table", 100).SetNullable()
			b.Integer("record_id").SetNullable()
			b.Text("old_data").SetNullable()
			b.Text("new_data").SetNullable()
			b.String("change_type", 20).SetNullable()
			b.Timestamps()
			b.SoftDeletes()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("operation_logs", drv)
			return b.ToDropSQL()
		},
	})
}
