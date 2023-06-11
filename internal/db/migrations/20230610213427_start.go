package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	
)

func init() {
	goose.AddMigration(upStart, downStart)
}

func upStart(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	query := `CREATE TABLE links (
		id bigserial PRIMARY KEY,
		short_link varchar(10),
		original_link varchar,
		created_at timestamp
	);`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("creating table: %s", err.Error())
	}
	return nil
}

func downStart(tx *sql.Tx) error {
	query := `DROP TABLE links;`
	_, err := tx.Exec(query)
	if err != nil {
		return fmt.Errorf("dropping table: %s", err.Error())
	}
	// This code is executed when the migration is rolled back.
	return nil
}
