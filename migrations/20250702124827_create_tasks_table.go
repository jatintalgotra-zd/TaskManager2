package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createTableTasks = `CREATE TABLE IF NOT EXISTS tasks (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(150),
    status TINYINT(1),
    user_id INT
);`

func createTasksTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// write your migrations here
			_, err := d.SQL.Exec(createTableTasks)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
