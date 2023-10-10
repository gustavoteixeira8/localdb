package main

import (
	"github.com/gustavoteixeira8/localdb/pkg/localdb/dbmgr"
	"github.com/gustavoteixeira8/localdb/pkg/localdb/repository"
)

type User struct {
	*repository.Base
	Name     string
	Username string
}

func main() {
	db := dbmgr.New(&dbmgr.DBManagerConfig{})

	err := db.Migrate(&User{})
	if err != nil {
		panic(err)
	}

	r := repository.New[User](db)

	// model, err := r.Add(User{Base: repository.NewBase(), Name: "Juliana", Username: "JUJU"})
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(model)

	err = r.UpdateWithQuery(func(model User) *repository.UpdateResponse[User] {
		return &repository.UpdateResponse[User]{
			Query: model.Name == "GUSTAVO",
			Value: User{Username: "GUGU"},
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
