package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260608200856_create_users_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("users", drv)
			b.ID()
			b.String("username", 50)
			b.String("password", 255)
			b.String("email", 100).SetUnique()
			b.String("phone", 20).SetNullable()
			b.String("avatar", 255).SetNullable()
			b.String("real_name", 50).SetNullable()
			b.String("status", 50).SetDefault("'active'")
			b.DateTime("last_login_at").SetNullable()
			b.Integer("role_id").SetNullable()
			b.Timestamps()
			b.SoftDeletes()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("users", drv)
			return b.ToDropSQL()
		},
	})
}
