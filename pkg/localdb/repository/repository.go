package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoteixeira8/db-go/pkg/localdb/dbmgr"
	"github.com/gustavoteixeira8/db-go/pkg/localdb/file"
)

type Repository[T Model] struct {
	DBManager *dbmgr.DBManager
	file      file.File[[]T]
}

func (r *Repository[T]) getTablePath() (string, error) {
	dataPath := r.DBManager.GetConfig().Path
	tablenames, err := r.DBManager.GetTableNames()
	if err != nil {
		return "", err
	}

	rightTablename := ""

	for _, tablename := range tablenames {
		n := r.DBManager.GetTableName(new(T))
		if tablename == n {
			rightTablename = tablename
			break
		}
	}

	return fmt.Sprintf("%s/%s", dataPath, rightTablename), nil
}

func (r *Repository[T]) Find(cb RepositoryFindCallback[T]) ([]T, error) {
	fullpath, err := r.getTablePath()
	if err != nil {
		return nil, err
	}

	val, err := r.file.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	tval := *val

	tvalFound := []T{}

	for _, v := range tval {
		if cb == nil {
			tvalFound = tval
			break
		}

		resp := cb(v)
		if resp.Query {
			tvalFound = append(tvalFound, v)
			if resp.StopOnFirst {
				break
			}
		}
	}

	return tvalFound, nil
}

func (r *Repository[T]) Add(cb RepositoryAddCallback[T]) (*T, error) {
	fullpath, err := r.getTablePath()
	if err != nil {
		return nil, err
	}

	val, err := r.file.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	tval := *val

	if len(tval) == 0 {
		resp := cb(*new(T))
		if resp == nil {
			return nil, errors.New("callback from Add cannot be nil")
		}
		resp.Value.SetID(uuid.NewString())
		resp.Value.SetCreatedAt(time.Now())
		resp.Value.SetUpdatedAt(time.Now())
		tval = append(tval, resp.Value)
	} else {
		for _, v := range tval {
			if cb == nil {
				break
			}

			resp := cb(v)

			resp.Value.SetID(uuid.NewString())
			resp.Value.SetCreatedAt(time.Now())
			resp.Value.SetUpdatedAt(time.Now())

			if resp.Query && resp.UseQuery {
				tval = append(tval, resp.Value)
				break
			}

			if !resp.UseQuery {
				tval = append(tval, resp.Value)
				break
			}
		}
	}

	err = r.file.WriteFile(fullpath, tval)

	return nil, err
}

func New[T Model](mgr *dbmgr.DBManager) *Repository[T] {
	r := &Repository[T]{DBManager: mgr}

	if mgr.GetConfig().FileType == dbmgr.FileTypeJSON {
		r.file = file.NewJSONFile[[]T]()
	} else if mgr.GetConfig().FileType == dbmgr.FileTypeYAML {
		r.file = file.NewYAMLFile[[]T]()
	}

	return r
}
