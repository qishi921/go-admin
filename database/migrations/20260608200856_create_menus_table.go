package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260608200856_create_menus_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("menus", drv)
			b.ID()
			b.String("name", 50)
			b.String("path", 100)
			b.String("icon", 50).SetNullable()
			b.String("component", 100).SetNullable()
			b.Integer("sort_order").SetDefault("0")
			b.Integer("parent_id").SetNullable()
			b.String("status", 50).SetDefault("'active'")
			b.Timestamps()
			b.SoftDeletes()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("menus", drv)
			return b.ToDropSQL()
		},
	})
}
