package mysql

import (
	"database/sql"
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) SaveUserAdmin(id uint, isAdmin bool) error {
	sqlStatement := `
		UPDATE users SET
			admin = ?
		WHERE id = ?
	`
	result, err := d.db.Exec(sqlStatement, isAdmin, id)
	if err != nil {
		return errors.Wrap(err, "updating user admin value")
	}
	count, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "fetching rows effected updating user admin")
	}
	if count == 0 {
		return notFound("User")
	}
	return nil
}

func (d *Datastore) SaveUserEnabled(id uint, isEnabled bool) error {
	sqlStatement := `
			UPDATE users SET
				enabled = ?
			WHERE id = ?
		`
	result, err := d.db.Exec(sqlStatement, isEnabled, id)
	if err != nil {
		return errors.Wrap(err, "updating user enabled value")
	}
	count, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "fetching rows effected updating user enabled")
	}
	if count == 0 {
		return notFound("User")
	}
	return nil
}

// NewUser creates a new user
func (d *Datastore) NewUser(user *kolide.User) (*kolide.User, error) {
	sqlStatement := `
		INSERT INTO users (
			password,
			salt,
			name,
			username,
			email,
			admin,
			enabled,
			admin_forced_password_reset,
			gravatar_url,
			position
		) VALUES (?,?,?,?,?,?,?,?,?,?)
	`
	result, err := d.db.Exec(sqlStatement, user.Password, user.Salt, user.Name,
		user.Username, user.Email, user.Admin, user.Enabled,
		user.AdminForcedPasswordReset, user.GravatarURL, user.Position)
	if err != nil {
		return nil, errors.Wrap(err, "create new user")
	}

	id, _ := result.LastInsertId()
	user.ID = uint(id)
	return user, nil
}

func (d *Datastore) findUser(searchCol string, searchVal interface{}) (*kolide.User, error) {
	sqlStatement := fmt.Sprintf(
		"SELECT * FROM users "+
			"WHERE %s = ? AND NOT deleted LIMIT 1",
		searchCol,
	)

	user := &kolide.User{}

	err := d.db.Get(user, sqlStatement, searchVal)
	if err != nil && err == sql.ErrNoRows {
		return nil, notFound("User").
			WithMessage(fmt.Sprintf("with %s=%v", searchCol, searchVal))
	} else if err != nil {
		return nil, errors.Wrap(err, "find user")
	}

	return user, nil
}

// User retrieves a user by name
func (d *Datastore) User(username string) (*kolide.User, error) {
	return d.findUser("username", username)
}

// ListUsers lists all users with limit, sort and offset passed in with
// kolide.ListOptions
func (d *Datastore) ListUsers(opt kolide.ListOptions) ([]*kolide.User, error) {
	sqlStatement := `
		SELECT * FROM users WHERE NOT deleted
	`
	sqlStatement = appendListOptionsToSQL(sqlStatement, opt)
	users := []*kolide.User{}

	if err := d.db.Select(&users, sqlStatement); err != nil {
		return nil, errors.Wrap(err, "list users")
	}

	return users, nil

}

func (d *Datastore) UserByEmail(email string) (*kolide.User, error) {
	return d.findUser("email", email)
}

func (d *Datastore) UserByID(id uint) (*kolide.User, error) {
	return d.findUser("id", id)
}

func (d *Datastore) SaveUser(user *kolide.User) error {
	sqlStatement := `
		UPDATE users SET
			username = ?,
			password = ?,
			salt = ?,
			name = ?,
			email = ?,
			admin = ?,
			enabled = ?,
			admin_forced_password_reset = ?,
			gravatar_url = ?,
			position = ?
		WHERE id = ?
	`
	_, err := d.db.Exec(sqlStatement, user.Username, user.Password,
		user.Salt, user.Name, user.Email, user.Admin, user.Enabled,
		user.AdminForcedPasswordReset, user.GravatarURL, user.Position, user.ID)
	if err != nil {
		return errors.Wrap(err, "save user")
	}

	return nil
}
