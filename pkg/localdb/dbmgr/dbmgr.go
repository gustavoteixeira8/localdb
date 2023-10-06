package dbmgr

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

type FileType string

const (
	FileTypeJSON FileType = "json"
	FileTypeYAML FileType = "yaml"
)

const _DB_PREFFIX = "localdb"
const _DB_BACKUP_PREFFIX = "backup"
const _BACKUP_DEFAULT_PATH = "./backup"

type DBManagerConfig struct {
	Path       string   `json:"path"`
	BackupPath string   `json:"backupPath"`
	FileType   FileType `json:"fileType"`
}

type DBManager struct {
	config *DBManagerConfig
}

func (db *DBManager) GetConfig() *DBManagerConfig {
	return db.config
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

		fullpath := fmt.Sprintf("%s.%s", tablename, db.config.FileType)
		fullpath = filepath.Join(db.config.Path, fullpath)

		_, err := os.Stat(fullpath)
		if err == nil {
			continue
		}

		err = os.WriteFile(fullpath, []byte(""), 0777)
		if err != nil {
			return fmt.Errorf("error creating %s file to %s", db.config.FileType, tablename)
		}
	}

	return nil
}

func (db DBManager) GetTableName(v any) string {
	typeofEntity := reflect.TypeOf(v)

	if typeofEntity.Kind() == reflect.Ptr {
		typeofEntity = typeofEntity.Elem()
	}

	if typeofEntity.Kind() != reflect.Struct && typeofEntity.Kind() != reflect.Map {
		return ""
	}

	tablename := db.formatTableName(typeofEntity.Name())

	return fmt.Sprintf("%s.%s", tablename, db.config.FileType)
}

func (db *DBManager) GetTableNames() ([]string, error) {
	tablenames := []string{}
	dirFiles, err := os.ReadDir(db.config.Path)

	if err != nil {
		return nil, err
	}

	for _, file := range dirFiles {
		name := file.Name()
		filetype := fmt.Sprintf(".%s", db.config.FileType)
		if strings.HasPrefix(name, _DB_PREFFIX) && strings.HasSuffix(name, filetype) {
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

func (db *DBManager) formatTableName(s string) string {
	sf := strings.ToLower(s)
	sf = strings.ReplaceAll(sf, " ", "_")
	sf = strings.ReplaceAll(sf, "-", "_")
	return fmt.Sprintf("%s_%s", _DB_PREFFIX, sf)
}

func New(config *DBManagerConfig) *DBManager {
	if config == nil {
		wd, _ := os.Getwd()
		config = &DBManagerConfig{
			Path:       wd,
			BackupPath: wd,
			FileType:   "json",
		}
	} else if config.FileType == "" {
		config.FileType = "json"
	}

	mgr := &DBManager{config: config}

	return mgr
}
