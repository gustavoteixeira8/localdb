package main

import (
	"github.com/gustavoteixeira8/db-go/pkg/localdb/dbmgr"
)

type User struct {
	Name string
}

func main() {
	db := dbmgr.New(&dbmgr.DBManagerConfig{Path: "."})
	// fmt.Println(db.Migrate(&User{}))
	// fmt.Println(db.GetTableNames())
	db.Backup()
}
