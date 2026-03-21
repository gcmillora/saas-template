package provider

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DbProvider struct {
	*sql.DB
}

// this is a transaction wrapper method.
// if tx is given, callbackFn would be executed under the given transaction. rollback/commit would not be handled in this call.
// if tx is not given, callbackFn would be executed under a new transaction created. rollback/commit would be handled in this call.
func (db *DbProvider) WithTransaction(tx *sql.Tx, callbackFn func(rootTx *sql.Tx) error) error {
	var rootTx = tx
	var err error

	// start new transaction if tx not given
	if rootTx == nil {
		tx, err = db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			err = tx.Rollback()
		}()
	}

	// execute callback function
	err = callbackFn(tx)
	if err != nil {
		return err
	}

	// commit transaction if tx not given
	if rootTx == nil {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewDbProvider(env *EnvProvider) *DbProvider {
	db, err := sql.Open("postgres", env.databaseUrl)
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}

	return &DbProvider{db}
}
