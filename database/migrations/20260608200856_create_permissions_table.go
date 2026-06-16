package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260608200856_create_permissions_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("permissions", drv)
			b.ID()
			b.String("name", 50)
			b.String("code", 50).SetUnique()
			b.String("description", 255).SetNullable()
			b.String("type", 50).SetDefault("'menu'")
			b.Integer("parent_id").SetNullable()
			b.Integer("sort_order").SetDefault("0")
			b.String("status", 50).SetDefault("'active'")
			b.Timestamps()
			b.SoftDeletes()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("permissions", drv)
			return b.ToDropSQL()
		},
	})
}
