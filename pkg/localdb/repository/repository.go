package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gustavoteixeira8/db-go/pkg/localdb/dbmgr"
)

type RepositoryFindCallback[T any] func(model T) bool

type Repository[T any] struct {
	DBManager *dbmgr.DBManager
}

func (r *Repository[T]) Find(cb RepositoryFindCallback[T]) ([]T, error) {
	dataPath := r.DBManager.GetConfig().Path
	tablenames, err := r.DBManager.GetTableNames()
	if err != nil {
		return nil, err
	}

	rightTablename := ""

	for _, tablename := range tablenames {
		n := r.DBManager.GetTableName(new(T))
		if tablename == n {
			rightTablename = tablename
			break
		}
	}

	fullpath := fmt.Sprintf("%s/%s", dataPath, rightTablename)

	databs, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	tval := []T{}

	err = json.Unmarshal(databs, &tval)
	if err != nil {
		return nil, err
	}

	tvalFound := []T{}

	for _, v := range tval {
		if cb == nil {
			tvalFound = tval
			break
		}

		isToAppend := cb(v)
		if isToAppend {
			tvalFound = append(tvalFound, v)
		}
	}

	return tvalFound, nil
}
