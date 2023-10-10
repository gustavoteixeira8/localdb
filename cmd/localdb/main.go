package main

import (
	"fmt"

	"github.com/gustavoteixeira8/localdb/pkg/localdb/dbmgr"
	"github.com/gustavoteixeira8/localdb/pkg/localdb/repository"
)

type User struct {
	*repository.Base
	Name     string
	Username string
}

func main() {
	db := dbmgr.New(&dbmgr.DBManagerConfig{FileType: dbmgr.FileTypeYAML})

	err := db.Migrate(&User{})
	if err != nil {
		panic(err)
	}

	r := repository.New[User](db)

	model, err := r.Add(User{Base: repository.NewBase(), Name: "GUSTAVO", Username: ""})
	if err != nil {
		panic(err)
	}

	fmt.Println(model)

	err = r.DeleteWithQuery(func(model User) *repository.DeleteResponse[User] {
		return &repository.DeleteResponse[User]{
			Query:       model.Name == "GUSTAVO",
			StopOnFirst: false,
		}
	})
	if err != nil {
		panic(err)
	}
}

// model, err := r.Find(func(model User) *repository.FindResponse[User] {
// 	return &repository.FindResponse[User]{StopOnFirst: true, Query: model.Name == "GUSTAVOa"}
// })
// if err != nil {
// 	panic(err)
// }
