package main

import (
	"fmt"

	"github.com/gustavoteixeira8/db-go/src/dbmanager"
)

func main() {
	db := dbmanager.New(&dbmanager.DBManagerConfig{Path: "."})
	fmt.Println(db.Migrate(&dbmanager.DBManagerConfig{}))
	fmt.Println(db.GetTableNames())
	fmt.Println(db.Backup())
}
