package db

import "github.com/gustavoteixeira8/db-go/src/dbmanager"

type DBEntity struct {
	ID string `json:"id"`
}

type DB struct {
	DBManager *dbmanager.DBManager
}

func (db *DB) Start() error {
	err := db.DBManager.Start()
	return err
}

func New(config *dbmanager.DBManagerConfig) *DB {
	return &DB{DBManager: dbmanager.New(config)}
}
