package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

var ErrDBNotInitialized = fmt.Errorf("db is not initialized")


func Init(dbfile string) error {
	if _, err := os.Stat(dbfile); os.IsNotExist(err) {
		f, err := os.Create(dbfile)
		if err != nil {
			return err
		}
		f.Close()
	}

	var err error
	DB, err = sqlx.Open("sqlite", dbfile)
	if err != nil {
		return err
	}
	schema := `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT,
	title TEXT,
	comment TEXT,
	repeat TEXT
);`
	_, err = DB.Exec(schema)
	return err
}
