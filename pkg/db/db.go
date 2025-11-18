package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

// Инициализация базы данных
func Init(dbfile string) error {

	// создаём файл, если его нет
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

	// создаём таблицу scheduler
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
