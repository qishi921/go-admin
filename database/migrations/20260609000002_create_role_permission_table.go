package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260609000002_create_role_permission_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("role_permission", drv)
			b.ID()
			b.Integer("role_id")
			b.Integer("permission_id")
			b.Timestamps()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("role_permission", drv)
			return b.ToDropSQL()
		},
	})
}
