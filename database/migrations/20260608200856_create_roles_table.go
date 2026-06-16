package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260608200856_create_roles_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("roles", drv)
			b.ID()
			b.String("name", 50)
			b.String("code", 50).SetUnique()
			b.String("description", 255).SetNullable()
			b.String("status", 50).SetDefault("'active'")
			b.Timestamps()
			b.SoftDeletes()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("roles", drv)
			return b.ToDropSQL()
		},
	})
}
