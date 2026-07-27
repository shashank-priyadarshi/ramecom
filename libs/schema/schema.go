package schema

import "database/sql"

func CreateTables(db *sql.DB) error {

	if err := createUsersTable(db); err != nil {
		return err
	}

	return nil
}

func createUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,

		username VARCHAR(100) NOT NULL UNIQUE,

		email VARCHAR(255) NOT NULL UNIQUE,

		password VARCHAR(255) NOT NULL,

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			ON UPDATE CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)

	return err
}
