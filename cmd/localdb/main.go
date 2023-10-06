package main

import (
	"github.com/gustavoteixeira8/db-go/pkg/localdb/dbmgr"
	"github.com/gustavoteixeira8/db-go/pkg/localdb/repository"
)

type User struct {
	*repository.Base
	Name string
}

func main() {
	db := dbmgr.New(&dbmgr.DBManagerConfig{Path: "."})

	err := db.Migrate(&User{})
	if err != nil {
		panic(err)
	}

	r := repository.New[User](db)

	// model, err := r.Add(func(model User) *repository.AddResponse[User] {
	// 	return &repository.AddResponse[User]{
	// 		Value: User{Base: repository.NewBase(), Name: "GUSTAVO"},
	// 	}
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(model)

	err = r.Delete(func(model User) *repository.DeleteResponse[User] {
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
