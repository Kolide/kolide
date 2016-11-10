// Package datastore implements Kolide's interactions with the database backend
package datastore

import (
	"errors"
	"fmt"

	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/kolide"
)

var (
	// ErrNotFound is returned when the datastore resource cannot be found
	ErrNotFound = errors.New("resource not found")

	// ErrExists is returned when creating a datastore resource that already exists
	ErrExists = errors.New("resource already created")
)

// New creates a kolide.Datastore with a database connection
// Use DBOption to pass optional arguments
func New(driver, conn string, opts ...DBOption) (kolide.Datastore, error) {
	opt := &dbOptions{
		maxAttempts: defaultMaxAttempts,
	}
	for _, option := range opts {
		if err := option(opt); err != nil {
			return nil, err
		}
	}

	// check if datastore is already present
	if opt.db != nil {
		return opt.db, nil
	}
	switch driver {
	case "gorm-mysql":
		db, err := openGORM("mysql", conn, opt.maxAttempts)
		if err != nil {
			return nil, err
		}
		ds := gormDB{
			DB:     db,
			Driver: "mysql",
		}
		// configure logger
		if opt.logger != nil {
			db.SetLogger(opt.logger)
			db.LogMode(opt.debug)
		}
		if err := ds.Migrate(); err != nil {
			return nil, err
		}
		return ds, nil
	case "inmem":
		ds := &inmem{
			Driver: "inmem",
		}

		err := ds.Migrate()
		if err != nil {
			return nil, err
		}

		return ds, nil
	default:
		return nil, fmt.Errorf("unsupported datastore driver %s", driver)
	}
}

// GetMysqlConnectionString returns a MySQL connection string using the
// provided configuration.
func GetMysqlConnectionString(conf config.MysqlConfig) string {
	return fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Address,
		conf.Database,
	)
}
