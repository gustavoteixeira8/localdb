package main

import (
	"fmt"

	"github.com/gustavoteixeira8/db-go/pkg/localdb/dbmgr"
	"github.com/gustavoteixeira8/db-go/pkg/localdb/repository"
)

type User struct {
	Name string
}

func main() {
	db := dbmgr.New(&dbmgr.DBManagerConfig{Path: "."})
	fmt.Println(db.Migrate(&User{}))
	// fmt.Println(db.GetTableNames())
	r := repository.Repository[User]{DBManager: db}
	fmt.Println(r.Find(func(model User) bool { return model.Name == "GUSTAVO" }))
}
