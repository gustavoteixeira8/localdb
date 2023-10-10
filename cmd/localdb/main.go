package main

import (
	"github.com/gustavoteixeira8/localdb/pkg/localdb/dbmgr"
	"github.com/gustavoteixeira8/localdb/pkg/localdb/repository"
)

type Admin struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type User struct {
	*repository.Base
	Name     string `json:"name"`
	Username string `json:"username"`
	Admin    Admin  `json:"admin"`
}

func main() {
	db := dbmgr.New(&dbmgr.DBManagerConfig{FileType: dbmgr.FileTypeJSON})

	err := db.Migrate(&User{})
	if err != nil {
		panic(err)
	}

	r := repository.New[User](db)

	// model, err := r.Add(
	// 	User{
	// 		Base:     repository.NewBase(),
	// 		Name:     "GUSTAVO",
	// 		Username: "",
	// 		Admin: struct{ Name string }{
	// 			Name: "abc",
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(model)

	err = r.DeleteWithQuery(func(model User) *repository.DeleteResponse[User] {
		return &repository.DeleteResponse[User]{
			Query: model.Admin.Name == "ADMIN",
		}
	})
	if err != nil {
		panic(err)
	}
	// fmt.Println(m)
}

// model, err := r.Find(func(model User) *repository.FindResponse[User] {
// 	return &repository.FindResponse[User]{StopOnFirst: true, Query: model.Name == "GUSTAVOa"}
// })
// if err != nil {
// 	panic(err)
// }
