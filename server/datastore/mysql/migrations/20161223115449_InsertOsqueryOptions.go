package migration

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/datastore/internal/appstate"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up_20161223115449, Down_20161223115449)
}

func Up_20161223115449(tx *sql.Tx) error {
	sqlStatement := `
    INSERT INTO options (
      name,
      type,
      value,
      read_only
    ) VALUES( ?, ?, ?, ?)
  `

	for _, opt := range appstate.Options {
		_, err := tx.Exec(sqlStatement, opt.Name, opt.Type, opt.Value, opt.ReadOnly)
		if err != nil {
			return err
		}

	}
	return nil
}

func Down_20161223115449(tx *sql.Tx) error {
	sqlStatement := `
		DELETE FROM options
		WHERE name = ?
	`
	for _, opt := range appstate.Options {
		_, err := tx.Exec(sqlStatement, opt.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
