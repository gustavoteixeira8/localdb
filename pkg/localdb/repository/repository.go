package repository

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoteixeira8/localdb/pkg/localdb/dbmgr"
	"github.com/gustavoteixeira8/localdb/pkg/localdb/filemgr"
)

type Repository[T Model] struct {
	DBManager *dbmgr.DBManager
	file      filemgr.FileMgr[[]T]
}

func (r *Repository[T]) getTablePath() (string, error) {
	dataPath := r.DBManager.GetConfig().Path
	tablenames, err := r.DBManager.GetTableNames()
	if err != nil {
		return "", err
	}

	rightTablename := ""

	tval := new(T)

	for _, tablename := range tablenames {
		n := r.DBManager.GetTableName(tval)
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

	tval, err := r.file.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

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

func (r *Repository[T]) AddWithQuery(cb RepositoryAddCallback[T]) (*T, error) {
	if cb == nil {
		return nil, errors.New("callback cannot be nil")
	}

	fullpath, err := r.getTablePath()
	if err != nil {
		return nil, err
	}

	tval, err := r.file.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	var tvalAdded *T

	if len(tval) == 0 {
		resp := cb(*new(T))
		if resp == nil {
			return nil, errors.New("response cannot be nil")
		}

		resp.Value.SetID(uuid.NewString())
		resp.Value.SetCreatedAt(time.Now())
		resp.Value.SetUpdatedAt(time.Now())
		tval = append(tval, resp.Value)
		tvalAdded = &resp.Value
	} else {
		for _, v := range tval {
			resp := cb(v)

			resp.Value.SetID(uuid.NewString())
			resp.Value.SetCreatedAt(time.Now())
			resp.Value.SetUpdatedAt(time.Now())

			tvalAdded = &resp.Value

			if resp.Query {
				tval = append(tval, resp.Value)
				break
			}
		}
	}

	err = r.file.WriteFile(fullpath, tval)

	return tvalAdded, err
}

func (r *Repository[T]) Add(data T) (*T, error) {
	fullpath, err := r.getTablePath()
	if err != nil {
		return nil, err
	}

	tval, err := r.file.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	iszero := reflect.ValueOf(data).FieldByName("Base").IsZero()

	if iszero {
		return nil, errors.New("data should have base model")
	}

	data.SetID(uuid.NewString())
	data.SetCreatedAt(time.Now())
	data.SetUpdatedAt(time.Now())
	tval = append(tval, data)

	err = r.file.WriteFile(fullpath, tval)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Repository[T]) DeleteWithQuery(cb RepositoryDeleteCallback[T]) error {
	if cb == nil {
		return errors.New("callback cannot be nil")
	}

	fullpath, err := r.getTablePath()
	if err != nil {
		return err
	}

	data, err := r.file.ReadFile(fullpath)
	if err != nil {
		return err
	}

	sort.Slice(data, func(i, j int) bool {
		return i > j
	})

	for i := len(data) - 1; i >= 0; i-- {
		val := data[i]

		resp := cb(val)
		if resp.Query {
			data = append(data[:i], data[i+1:]...)
			if resp.StopOnFirst {
				break
			}
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return i > j
	})

	err = r.file.WriteFile(fullpath, data)

	return err
}

func (r *Repository[T]) Delete(id string) error {
	fullpath, err := r.getTablePath()
	if err != nil {
		return err
	}

	alldata, err := r.file.ReadFile(fullpath)
	if err != nil {
		return err
	}

	for i, data := range alldata {
		if data.GetID() == id {
			alldata = append(alldata[:i], alldata[i+1:]...)
			break
		}
	}

	err = r.file.WriteFile(fullpath, alldata)

	return err
}

func (r *Repository[T]) Update(id string, newdata T) error {
	fullpath, err := r.getTablePath()
	if err != nil {
		return err
	}

	alldata, err := r.file.ReadFile(fullpath)
	if err != nil {
		return err
	}

	for i, data := range alldata {
		if data.GetID() == id {
			dataValue := reflect.Indirect(reflect.ValueOf(&data))
			newdataType := reflect.TypeOf(newdata)
			dataType := reflect.TypeOf(data)
			newdataValue := reflect.Indirect(reflect.ValueOf(&newdata))

			for _, newField := range reflect.VisibleFields(newdataType) {
				verifyBaseFields := newField.Name == "Base" || newField.Name == "ID"
				verifyBaseFields = verifyBaseFields && newField.Name == "CreatedAt" || newField.Name == "UpdatedAt"

				if verifyBaseFields {
					continue
				}

				for _, oldField := range reflect.VisibleFields(dataType) {
					verifyBaseFields := oldField.Name == "Base" || oldField.Name == "ID"
					verifyBaseFields = verifyBaseFields && oldField.Name == "CreatedAt" || oldField.Name == "UpdatedAt"

					if verifyBaseFields {
						continue
					}

					if newField.Name == oldField.Name {
						newvalue := newdataValue.FieldByName(newField.Name)
						oldvalue := dataValue.FieldByName(oldField.Name)

						if newvalue.IsValid() && !newvalue.IsZero() {
							oldvalue.Set(newvalue)
						}

						break
					}
				}
			}

			finalnewdata := dataValue.Interface().(T)
			finalnewdata.SetUpdatedAt(time.Now())
			alldata[i] = finalnewdata

			break
		}
	}

	err = r.file.WriteFile(fullpath, alldata)

	return err
}

func (r *Repository[T]) UpdateWithQuery(cb RepositoryUpdateCallback[T]) error {
	fullpath, err := r.getTablePath()
	if err != nil {
		return err
	}

	alldata, err := r.file.ReadFile(fullpath)
	if err != nil {
		return err
	}

	for _, v := range alldata {
		resp := cb(v)
		if resp.Query {
			err = r.Update(v.GetID(), resp.Value)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func New[T Model](mgr *dbmgr.DBManager) *Repository[T] {
	r := &Repository[T]{DBManager: mgr}

	if mgr.GetConfig().FileType == dbmgr.FileTypeJSON {
		r.file = filemgr.NewJSONFile[[]T]()
	} else if mgr.GetConfig().FileType == dbmgr.FileTypeYAML {
		r.file = filemgr.NewYAMLFile[[]T]()
	}

	return r
}
