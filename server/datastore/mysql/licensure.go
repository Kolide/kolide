package mysql

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

// validateToken will insure that a jwt token is signed properly and that we
// have the public key we need to validate it.  The public key string is returned
// on success
func (ds *Datastore) validateToken(token string) (string, error) {
	var (
		publicKeyHash string
		publicKey     string
	)
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		h, ok := token.Header["kid"]
		if !ok {
			return "", errors.New("missing kid header")
		}
		publicKeyHash, ok = h.(string)
		if !ok {
			return "", errors.New("kid is not expected type")
		}

		sql := `SELECT pk.key
						FROM public_keys pk
						WHERE hash = ?`

		err := ds.db.Get(&publicKey, sql, publicKeyHash)
		if err != nil {
			return "", errors.Wrap(err, "could not find public key matching hash")
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	})
	return publicKey, err
}

func (ds *Datastore) SaveLicense(token string) (*kolide.License, error) {
	publicKeyString, err := ds.validateToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "token validation failed")
	}

	sqlStatement := "UPDATE licensure SET " +
		"  token = ?, " +
		"	`key` = ? " +
		"WHERE id = 1"

	_, err = ds.db.Exec(sqlStatement, token, publicKeyString)
	if err != nil {
		return nil, errors.Wrap(err, "saving license")
	}
	result, err := ds.License()
	if err != nil {
		return nil, errors.Wrap(err, "fetching license")
	}
	return result, nil
}

func (ds *Datastore) License() (*kolide.License, error) {
	query := `
  SELECT * FROM licensure
    WHERE id = 1
  `
	var license kolide.License
	err := ds.db.Get(&license, query)
	if err != nil {
		return nil, errors.Wrap(err, "fetching license information")
	}
	query = `
    SELECT count(*)
      FROM hosts
      WHERE NOT deleted
  `
	err = ds.db.Get(&license.HostCount, query)
	if err != nil {
		return nil, errors.Wrap(err, "fetching host count for license")
	}
	return &license, nil
}
