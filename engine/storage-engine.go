package engine

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	api "igbodb/grpc"

	// needed for SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

const create string = `
		CREATE TABLE IF NOT EXISTS StorageEngine (
		id TEXT NOT NULL PRIMARY KEY,
		time DATETIME NOT NULL,
		description TEXT
		);
`
const file string = "StorageEngine.db"

type StorageEngine struct {
	db *sql.DB
}

func (c *StorageEngine) NewStorageEngine() (*StorageEngine, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &StorageEngine{
		db: db,
	}, nil
}
func (c *StorageEngine) Insert(activity *api.Activity) (int, error) {
	res, err := c.db.Exec("INSERT INTO StorageEngine VALUES(?,?,?);", activity.Id, activity.Time.AsTime(), activity.Description)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	log.Printf("Added %v as %d", activity, id)
	return int(id), nil
}

func (c *StorageEngine) Update(activity *api.Activity) (int, error) {
	res, err := c.db.Exec("UPDATE StorageEngine SET description = ? where id = ?;", activity.Description, activity.Id)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.RowsAffected(); err != nil {
		return 0, err
	}
	log.Printf("Added %v as %d", activity, id)
	return int(id), nil
}

func (c *StorageEngine) Delete(id string) (int, error) {
	res, err := c.db.Exec("DELETE FROM StorageEngine where id = ?;", id)
	if err != nil {
		return 0, err
	}

	var retId int64
	if retId, err = res.RowsAffected(); err != nil {
		return 0, err
	}
	log.Printf("Added %v", id)
	return int(retId), nil
}

var ErrIDNotFound = errors.New("Id not found")

func (c *StorageEngine) Retrieve(id string) (*api.Activity, error) {
	log.Printf("Getting %v", id)

	// Query DB row based on ID
	row := c.db.QueryRow("SELECT id, time, description FROM StorageEngine WHERE id=?", id)

	// Parse row into Interval struct
	activity := api.Activity{}
	var err error
	var time time.Time
	if err = row.Scan(&activity.Id, &time, &activity.Description); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return &api.Activity{}, ErrIDNotFound
	}
	activity.Time = timestamppb.New(time)
	return &activity, err
}

func (c *StorageEngine) List(offset int) ([]*api.Activity, error) {
	log.Printf("Getting list from offset %d\n", offset)

	// Query DB row based on ID
	rows, err := c.db.Query("SELECT * FROM StorageEngine WHERE ID > ? ORDER BY id DESC LIMIT 100", offset)
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
}
