package datastore

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestEnrollHostMySQLGORM(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}
	db := setupMySQLGORM(t)
	defer teardownMySQLGORM(t, db)

	testEnrollHost(t, db)
}

func TestAuthenticateHostMySQLGORM(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}
	db := setup(t)
	defer teardown(t, db)

	testAuthenticateHost(t, db)
}

func TestUserByIDMySQLGORM(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}
	db := setup(t)
	defer teardown(t, db)

	testUserByID(t, db)
}

// TestCreateUser tests the UserStore interface
// this test uses the MySQL GORM backend
func TestCreateUserMySQLGORM(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	db := setupMySQLGORM(t)
	defer teardownMySQLGORM(t, db)

	testCreateUser(t, db)
}

func TestSaveUserMySQLGORM(t *testing.T) {
	address := os.Getenv("MYSQL_ADDR")
	if address == "" {
		t.SkipNow()
	}

	db := setupMySQLGORM(t)
	defer teardownMySQLGORM(t, db)

	testSaveUser(t, db)
}

func setupMySQLGORM(t *testing.T) Datastore {
	user := "kolide"
	password := "kolide"
	dbName := "kolide"

	// try container first
	host := os.Getenv("MYSQL_PORT_3306_TCP_ADDR")
	if host == "" {
		host = "127.0.0.1:3306"
	}

	conn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err := New("gorm-mysql", conn, LimitAttempts(1))
	if err != nil {
		t.Fatal(err)
	}

	backend := db.(gormDB)
	if err := backend.Migrate(); err != nil {
		t.Fatal(err)
	}

	return db
}

func teardownMySQLGORM(t *testing.T, db Datastore) {
	if err := db.Drop(); err != nil {
		t.Fatal(err)
	}
}
