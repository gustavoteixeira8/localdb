package dbmanager

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const _DB_PREFFIX = "db_go"
const _DB_BACKUP_PREFFIX = "backup"
const _BACKUP_DEFAULT_PATH = "./backup"

type DBManagerConfig struct {
	Path       string `json:"path"`
	BackupPath string `json:"backupPath"`
}

type DBManager struct {
	config *DBManagerConfig
}

func (db *DBManager) Start() error {
	if db.config == nil {
		return errors.New("db config cannot be nil")
	}

	var err error

	db.config.Path, err = filepath.Abs(db.config.Path)

	if err != nil {
		return err
	}

	db.config.Path, err = filepath.Abs(db.config.BackupPath)

	if err != nil {
		return err
	}

	err = os.MkdirAll(db.config.Path, 0777)

	if err != nil {
		return err
	}

	return nil
}

func (db *DBManager) Migrate(v ...any) error {
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

		fullpath := fmt.Sprintf("%s.json", tablename)
		fullpath = filepath.Join(db.config.Path, fullpath)

		_, err := os.Stat(fullpath)
		if err == nil {
			continue
		}

		err = os.WriteFile(fullpath, []byte("[]"), 0777)
		if err != nil {
			return fmt.Errorf("error creating json file to %s", tablename)
		}
	}

	return nil
}

func (db *DBManager) GetTableNames() ([]string, error) {
	tablenames := []string{}
	dirFiles, err := os.ReadDir(db.config.Path)

	if err != nil {
		return nil, err
	}

	for _, file := range dirFiles {
		name := file.Name()
		if strings.HasPrefix(name, _DB_PREFFIX) && strings.HasSuffix(name, ".json") {
			tablenames = append(tablenames, name)
		}
	}

	return tablenames, nil
}

func (db *DBManager) Backup() error {
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

	err = os.MkdirAll(db.config.BackupPath, 0777)

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
		backupName = filepath.Join(db.config.BackupPath, backupName)
		err = os.WriteFile(backupName, fileBytes, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DBManager) formatTableName(s string) string {
	sf := strings.ToLower(s)
	sf = strings.ReplaceAll(sf, " ", "_")
	sf = strings.ReplaceAll(sf, "-", "_")
	return fmt.Sprintf("%s_%s", _DB_PREFFIX, sf)
}

func New(config *DBManagerConfig) *DBManager {
	return &DBManager{config}
}
