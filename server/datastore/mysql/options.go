package mysql

import (
	"database/sql"

	"github.com/kolide/kolide/server/datastore/internal/appstate"
	"github.com/kolide/kolide/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) ResetOptions() ([]kolide.Option, error) {
	// Atomically remove all existing options, reset auto increment so id's will be the
	// same as original defaults, and re-insert defaults in option table.
	txn, err := d.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "reset options begin transaction")
	}
	var success bool
	defer func() {
		if success {
			if err = txn.Commit(); err == nil {
				return
			}
			d.logger.Log(
				"method", "ResetOptions",
				"activity", "db commit",
				"err", err,
			)
		}
		if err = txn.Rollback(); err == nil {
			return
		}
		d.logger.Log(
			"method", "ResetOptions",
			"activity", "db rollback",
			"err", err,
		)
	}()

	_, err = d.db.Exec("DELETE FROM options")
	if err != nil {
		return nil, errors.Wrap(err, "deleting options in reset options")
	}
	// Reset auto increment
	_, err = d.db.Exec("ALTER TABLE `options` AUTO_INCREMENT = 1")
	if err != nil {
		return nil, errors.Wrap(err, "resetting auto increment counter in reset options")
	}
	sqlStatement := `
		INSERT INTO options (
			name,
			type,
			value,
			read_only
		) VALUES (?, ?, ?, ?)
	`
	var options []kolide.Option
	for _, defaultOpt := range appstate.Options() {
		opt := kolide.Option{
			Name:     defaultOpt.Name,
			ReadOnly: defaultOpt.ReadOnly,
			Type:     defaultOpt.Type,
			Value: kolide.OptionValue{
				Val: defaultOpt.Value,
			},
		}
		dbResponse, err := d.db.Exec(
			sqlStatement,
			opt.Name,
			opt.Type,
			opt.Value,
			opt.ReadOnly,
		)
		if err != nil {
			return nil, errors.Wrap(err, "inserting default option in reset options")
		}
		id, err := dbResponse.LastInsertId()
		if err != nil {
			return nil, errors.Wrap(err, "fetching id in reset options")
		}
		opt.ID = uint(id)
		options = append(options, opt)
	}
	// We've removed all old options and restored defaults, indicate success
	// so our work will be committed in our defer func
	success = true
	return options, nil
}

func (d *Datastore) OptionByName(name string) (*kolide.Option, error) {
	sqlStatement := `
			SELECT *
			FROM options
			WHERE name = ?
		`
	var option kolide.Option
	if err := d.db.Get(&option, sqlStatement, name); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Option")
		}
		return nil, errors.Wrap(err, sqlStatement)
	}
	return &option, nil
}

func (d *Datastore) SaveOptions(opts []kolide.Option) (err error) {
	sqlStatement := `
		UPDATE options
		SET value = ?
		WHERE id = ? AND type = ? AND NOT read_only
	`
	txn, err := d.db.Begin()
	if err != nil {
		return errors.Wrap(err, "update options begin transaction")
	}
	var success bool
	defer func() {
		if success {
			if err = txn.Commit(); err == nil {
				return
			}
		}
		txn.Rollback()
	}()

	for _, opt := range opts {
		result, err := txn.Exec(sqlStatement, opt.Value, opt.ID, opt.Type)
		if err != nil {
			return errors.Wrap(err, "update options")
		}
		rowsChanged, err := result.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "option rows affected")
		}
		if rowsChanged != 1 {
			return notFound("Option").WithID(opt.ID)
		}
	}
	// If all the updates succeed, set the success flag, this will cause the
	// function we defined in defer to commit the transaction. Otherwise, all
	// changes will be rolled back
	success = true
	return err
}

func (d *Datastore) Option(id uint) (*kolide.Option, error) {
	sqlStatement := `
		SELECT *
		FROM options
		WHERE id = ?
	`
	var opt kolide.Option
	if err := d.db.Get(&opt, sqlStatement, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Option").WithID(id)
		}
		return nil, errors.Wrap(err, "select option by ID")
	}
	return &opt, nil
}

func (d *Datastore) ListOptions() ([]kolide.Option, error) {
	sqlStatement := `
    SELECT *
    FROM options
    ORDER BY name ASC
  `
	var opts []kolide.Option
	if err := d.db.Select(&opts, sqlStatement); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Option")
		}
		return nil, errors.Wrap(err, "select from options")
	}
	return opts, nil
}

func (d *Datastore) GetOsqueryConfigOptions() (map[string]interface{}, error) {
	// Retrieve all the options that are set. The value field is JSON formatted so
	// to retrieve options that are set, we check JSON null keyword
	sqlStatement := `
		SELECT *
		FROM options
		WHERE value != "null"
	`
	var opts []kolide.Option
	if err := d.db.Select(&opts, sqlStatement); err != nil {
		if err == sql.ErrNoRows {
			return nil, notFound("Option")
		}
		return nil, errors.Wrap(err, "select from options")
	}
	optConfig := map[string]interface{}{}
	for _, opt := range opts {
		optConfig[opt.Name] = opt.GetValue()
	}
	return optConfig, nil
}
