package migrations

import (
	"github.com/Hlgxz/gai/database/driver"
	"github.com/Hlgxz/gai/database/migration"
)

func init() {
	Migrations = append(Migrations, migration.Migration{
		Name: "20260612000001_create_uploads_table",
		Up: func(drv driver.Driver) string {
			return `CREATE TABLE IF NOT EXISTS uploads (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				file_name VARCHAR(255) NOT NULL,
				original_name VARCHAR(255) NOT NULL,
				file_path VARCHAR(500) NOT NULL,
				file_size INTEGER NOT NULL,
				mime_type VARCHAR(100) NOT NULL,
				extension VARCHAR(20) NOT NULL,
				user_id INTEGER,
				module VARCHAR(50) DEFAULT 'general',
				status VARCHAR(20) DEFAULT 'active',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				deleted_at DATETIME
			);
			CREATE INDEX IF NOT EXISTS idx_uploads_user_id ON uploads(user_id);
			CREATE INDEX IF NOT EXISTS idx_uploads_module ON uploads(module);
			CREATE INDEX IF NOT EXISTS idx_uploads_status ON uploads(status);`
		},
		Down: func(drv driver.Driver) string {
			return `DROP TABLE IF EXISTS uploads;`
		},
	})
}
