package main

import (
	"fmt"

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
	dbconfig := &dbmgr.DBManagerConfig{StorageType: dbmgr.StorageTypeJSON}

	r := repository.New[User](dbconfig)

	err := r.Migrate(User{})
	if err != nil {
		panic(err)
	}

	model, err := r.Add(
		User{
			Base:     repository.NewBase(),
			Name:     "GUSTAVO",
			Username: "",
			Admin: Admin{
				Name: "abc",
			},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(model)

	err = r.DeleteWithQuery(func(model User) *repository.DeleteResponse[User] {
		return &repository.DeleteResponse[User]{
			Query: model.Name == "GUSTAVO1",
		}
	})
	if err != nil {
		panic(err)
	}

	model1, err := r.Find(func(model User) *repository.FindResponse[User] {
		return &repository.FindResponse[User]{
			Query: model.Name == "GUSTAVO",
		}
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(model1)

	model1, err = r.Find(func(model User) *repository.FindResponse[User] {
		return &repository.FindResponse[User]{StopOnFirst: true, Query: model.Name == "GUSTAVOa"}
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(model1)

	err = r.DeleteWithQuery(func(model User) *repository.DeleteResponse[User] {
		return &repository.DeleteResponse[User]{
			Query: model.Admin.Name == "ADMIN",
		}
	})
	if err != nil {
		panic(err)
	}

	err = r.DeleteWithQuery(func(model User) *repository.DeleteResponse[User] {
		return &repository.DeleteResponse[User]{
			Query: true,
		}
	})
	if err != nil {
		panic(err)
	}
}
