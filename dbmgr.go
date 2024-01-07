package localdb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type StorageType string

const (
	StorageTypeJSON   StorageType = "json"
	StorageTypeYAML   StorageType = "yaml"
	StorageTypeMemory StorageType = "memory"
)

const _DB_PREFFIX = "localdb"
const _DB_BACKUP_PREFFIX = "backup"
const _BACKUP_DEFAULT_PATH = "./backup"

type DBManagerConfig struct {
	Path        string      `json:"path" yaml:"path"`
	BackupPath  string      `json:"backupPath" yaml:"backupPath"`
	StorageType StorageType `json:"storageType" yaml:"storageType"`
}

type DBManager[T any] struct {
	config  *DBManagerConfig
	storage StorageMgr[[]T]
}

func (db *DBManager[T]) GetConfig() *DBManagerConfig {
	return db.config
}

func (db *DBManager[T]) GetStorage() StorageMgr[[]T] {
	return db.storage
}

func (db *DBManager[T]) Start() error {
	if db.config == nil {
		return errors.New("db config cannot be nil")
	}

	var err error

	if db.config.StorageType == StorageTypeMemory {
		return nil
	}

	db.config.Path, err = filepath.Abs(db.config.Path)

	if err != nil {
		return err
	}

	db.config.BackupPath, err = filepath.Abs(db.config.BackupPath)

	if err != nil {
		return err
	}

	err = os.MkdirAll(db.config.Path, 0777)

	if err != nil {
		return err
	}

	return nil
}

func (db *DBManager[T]) Migrate(v ...any) error {
	tablename := ""

	for _, value := range v {
		typeofEntity := reflect.TypeOf(value)

		if typeofEntity.Kind() == reflect.Ptr {
			typeofEntity = typeofEntity.Elem()
		}

		if typeofEntity.Kind() != reflect.Struct && typeofEntity.Kind() != reflect.Map {
			return errors.New("only structs and maps are allowed to be migrated")
		}

		tablename = db.formatTableName(typeofEntity.Name())

		if strings.Trim(tablename, " ") == "" {
			return errors.New("cannot use a struct/map without name")
		}

		if db.config.StorageType == StorageTypeMemory {
			// do something
			continue
		}

		fullpath := fmt.Sprintf("%s.%s", tablename, db.config.StorageType)
		fullpath = filepath.Join(db.config.Path, fullpath)

		_, err := os.Stat(fullpath)
		if err == nil {
			continue
		}

		err = os.WriteFile(fullpath, []byte("[]"), 0777)
		if err != nil {
			return fmt.Errorf("error creating %s file to %s", db.config.StorageType, tablename)
		}
	}

	return nil
}

func (db DBManager[T]) GetTableName(v any) string {
	typeofEntity := reflect.TypeOf(v)

	if typeofEntity.Kind() == reflect.Ptr {
		typeofEntity = typeofEntity.Elem()
	}

	if typeofEntity.Kind() != reflect.Struct && typeofEntity.Kind() != reflect.Map {
		return ""
	}

	tablename := db.formatTableName(typeofEntity.Name())

	return fmt.Sprintf("%s.%s", tablename, db.config.StorageType)
}

func (db *DBManager[T]) GetTableNames() ([]string, error) {
	tablenames := []string{}
	dirFiles, err := os.ReadDir(db.config.Path)

	if err != nil {
		return nil, err
	}

	for _, file := range dirFiles {
		name := file.Name()
		StorageType := fmt.Sprintf(".%s", db.config.StorageType)
		if strings.HasPrefix(name, _DB_PREFFIX) && strings.HasSuffix(name, StorageType) {
			tablenames = append(tablenames, name)
		}
	}

	return tablenames, nil
}

func (db *DBManager[T]) Backup() error {
	var err error
	if db.config.BackupPath == "" {
		db.config.BackupPath, err = filepath.Abs(_BACKUP_DEFAULT_PATH)
		if err != nil {
			return err
		}
	}

	tablenames, err := db.GetTableNames()
	if err != nil {
		return err
	}

	backuppath := fmt.Sprintf("%s_%d", db.config.BackupPath, time.Now().UnixMilli())

	err = os.MkdirAll(backuppath, 0777)

	if err != nil {
		return err
	}

	for _, name := range tablenames {
		fullpath := filepath.Join(db.config.Path, name)
		fileBytes, err := os.ReadFile(fullpath)
		if err != nil {
			log.Printf("error reading file to backup (%v)\n", err)
			continue
		}
		backupName := fmt.Sprintf("%s_%s", _DB_BACKUP_PREFFIX, name)
		backupName = filepath.Join(backuppath, backupName)
		err = os.WriteFile(backupName, fileBytes, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DBManager[T]) formatTableName(s string) string {
	sf := strings.ToLower(s)
	sf = strings.ReplaceAll(sf, " ", "_")
	sf = strings.ReplaceAll(sf, "-", "_")
	return fmt.Sprintf("%s_%s", _DB_PREFFIX, sf)
}

func newDBMgr[T any](config *DBManagerConfig) *DBManager[T] {
	wd, _ := os.Getwd()
	if config == nil {
		fmt.Println(wd)
		config = &DBManagerConfig{
			Path:        wd,
			BackupPath:  wd,
			StorageType: "json",
		}
	}
	if config.StorageType == "" {
		config.StorageType = "json"
	}
	if config.Path == "" {
		config.Path = wd
	}
	if config.BackupPath == "" {
		config.BackupPath = config.Path
	}

	mgr := &DBManager[T]{config: config}

	if config.StorageType == StorageTypeJSON {
		mgr.storage = NewJSONStorage[[]T]()
	} else if config.StorageType == StorageTypeYAML {
		mgr.storage = NewYAMLStorage[[]T]()
	} else if config.StorageType == StorageTypeMemory {
		mgr.storage = NewMemoryStorage[[]T]()
	}

	return mgr
}
