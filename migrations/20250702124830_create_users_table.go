package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createTableUsers = `CREATE TABLE IF NOT EXISTS users (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(50)
);`

func createUsersTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// write your migrations here
			_, err := d.SQL.Exec(createTableUsers)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
