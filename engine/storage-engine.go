package engine

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	// needed for SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

const createStatement string = `
		CREATE TABLE IF NOT EXISTS %v (
		key TEXT NOT NULL PRIMARY KEY,
		value TEXT
		);
`
const file string = "KeyValueStorage.db"

type StorageEngine interface {
	Create(storageName string, key string, value string) error
	Retrieve(storageName string, key string) (string, error)
	Update(storageName string, key string, value string) error
	Delete(storageName string, key string) error
}

type StorageEngineImpl struct {
	db *sql.DB
}

func (c *StorageEngineImpl) NewStorageEngine() (*StorageEngineImpl, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	return &StorageEngineImpl{
		db: db,
	}, nil
}
func (c *StorageEngineImpl) Create(storageName string, key string, value string) error {
	res, err := c.db.Exec(
		fmt.Sprintf("INSERT INTO %v VALUES(?,?);", strings.ToUpper(storageName)), key, value)
	if err != nil {
		//TODO check how to recognize that table does not exist yet and if table does not exist, create it and try again
		if err != nil {
			log.Printf("Storage %v does not yet exist, creating...", storageName)
			if _, err := c.db.Exec(fmt.Sprintf(createStatement, storageName)); err != nil {
				return err
			}
			err = nil
			res, err = c.db.Exec(
				fmt.Sprintf("INSERT INTO %v VALUES(?,?);", strings.ToUpper(storageName)), key, value)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if _, err = res.LastInsertId(); err != nil {
		return err
	}
	log.Printf("Created on storage %v with key %v value %v", storageName, key, value)
	return nil
}

func (c *StorageEngineImpl) Update(storageName string, key string, value string) error {
	_, err := c.db.Exec(
		fmt.Sprintf("UPDATE %v SET value = ? where key = ?;", storageName), value, key)
	if err != nil {
		return err
	}

	return nil
}

func (c *StorageEngineImpl) Delete(storageName string, key string) error {
	_, err := c.db.Exec(
		fmt.Sprintf("DELETE FROM %v where key = ?;", storageName), key)
	if err != nil {
		return err
	}

	return nil
}

var ErrIDNotFound = errors.New("Key not found")

func (c *StorageEngineImpl) Retrieve(storageName string, key string) (string, error) {
	log.Printf("Getting %v", key)

	// Query DB row based on ID
	row := c.db.QueryRow(fmt.Sprintf("SELECT value FROM %v WHERE key=?", storageName), key)

	var err error
	var value *string

	if err = row.Scan(&value); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return "", ErrIDNotFound
	}

	return *value, err
}

//TODO use as example when implementing queries
/**func (c *StorageEngineImpl) List(offset int) ([]*api.Activity, error) {
	log.Printf("Getting list from offset %d\n", offset)

	// Query DB row based on ID
	rows, err := c.db.Query("SELECT * FROM StorageEngineImpl WHERE ID > ? ORDER BY id DESC LIMIT 100", offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*api.Activity{}
	for rows.Next() {
		i := api.Activity{}
		var time time.Time
		err = rows.Scan(&i.Id, &time, &i.Description)
		if err != nil {
			return nil, err
		}
		i.Time = timestamppb.New(time)
		data = append(data, &i)
	}
	return data, nil
}**/
