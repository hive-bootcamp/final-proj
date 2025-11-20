package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

func Tasks(limit int, search string) ([]*Task, error) {
	if DB == nil {
		return nil, errors.New("db is not initialized")
	}

	var rows *sql.Rows
	var err error

	//  поиск не указан → список задач
	if search == "" {
		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat 
             FROM scheduler 
             ORDER BY date 
             LIMIT ?`, limit)
		if err != nil {
			return nil, err
		}
		return scanTasks(rows)
	}

	// поиск — дата?
	if t, errDate := time.Parse("02.01.2006", search); errDate == nil {
		date := t.Format("20060102")
		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat
             FROM scheduler 
             WHERE date = ?
             ORDER BY date 
             LIMIT ?`, date, limit)
		if err != nil {
			return nil, err
		}
		return scanTasks(rows)
	}

	//поиск по подстроке в title/comment
	like := "%" + strings.ToLower(search) + "%"

	rows, err = DB.Query(
		`SELECT id, date, title, comment, repeat
         FROM scheduler
         WHERE lower(title) LIKE ? OR lower(comment) LIKE ?
         ORDER BY date 
         LIMIT ?`,
		like, like, limit)

	if err != nil {
		return nil, err
	}

	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		t := &Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, errors.New("db is not initialized")
	}

	query := `
	INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?, ?, ?, ?)
	`

	res, err := DB.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
