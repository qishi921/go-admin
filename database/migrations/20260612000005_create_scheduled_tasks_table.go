package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	// 定时任务表
	Migrations = append(Migrations, migration.Migration{
		Name: "20260612000005_create_scheduled_tasks_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("scheduled_tasks", drv)
			b.ID()
			b.String("name", 200)
			b.String("code", 100).SetUnique()
			b.Text("description").SetNullable()
			b.String("cron_expr", 100)
			b.Text("command").SetNullable()
			b.Text("params").SetNullable()
			b.String("status", 20).SetDefault("'enabled'")
			b.Integer("priority").SetDefault("0")
			b.Integer("max_retries").SetDefault("3")
			b.Integer("timeout").SetDefault("3600")
			b.DateTime("last_run_at").SetNullable()
			b.DateTime("next_run_at").SetNullable()
			b.Integer("run_count").SetDefault("0")
			b.Integer("fail_count").SetDefault("0")
			b.Text("last_error").SetNullable()
			b.Timestamps()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("scheduled_tasks", drv)
			return b.ToDropSQL()
		},
	})

	// 任务执行记录表
	Migrations = append(Migrations, migration.Migration{
		Name: "20260612000006_create_task_executions_table",
		Up: func(drv driver.Driver) string {
			b := migration.NewBlueprint("task_executions", drv)
			b.ID()
			b.Integer("task_id")
			b.String("task_name", 200)
			b.String("status", 20)
			b.DateTime("started_at")
			b.DateTime("finished_at").SetNullable()
			b.Integer("duration").SetDefault("0")
			b.Text("output").SetNullable()
			b.Text("error_msg").SetNullable()
			b.Integer("retry_count").SetDefault("0")
			b.String("triggered_by", 50)
			b.Timestamps()
			return b.ToCreateSQL()
		},
		Down: func(drv driver.Driver) string {
			b := migration.NewBlueprint("task_executions", drv)
			return b.ToDropSQL()
		},
	})
}
