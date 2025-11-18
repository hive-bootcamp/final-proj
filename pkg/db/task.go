package db

import (
	"database/sql"
	"errors"
)

type Task struct {
	ID      int64  `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

// DB должен быть инициализирован в db.go (Init() или что-то подобное)
// тип: var DB *sqlx.DB

func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, errors.New("db is not initialized")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat)
	          VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}
