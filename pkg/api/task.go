package api

import (
	"encoding/json"
	"final-proj/pkg/db"
	"net/http"
	"strconv"
	"time"
)

type errorResponse struct {
	Error string `json:"error"`
}

type idResponse struct {
	ID string `json:"id"`
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(v)
}

// общий обработчик /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		// для других методов пока просто ошибка в JSON
		writeJSON(w, errorResponse{Error: "unsupported method"})
	}
}

// проверка и корректировка даты
func checkDate(task *db.Task) error {
	now := time.Now()

	// если дата не указана или пустая → ставим сегодня
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	// парсим дату в формате 20060102
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	var next string
	// если есть правило повторения → проверяем его и заодно считаем следующую дату
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// если сегодня больше даты задачи
	if afterNow(now, t) {
		if task.Repeat == "" {
			// без повторения — ставим сегодня
			task.Date = now.Format(DateFormat)
		} else {
			// с повторением — ставим ближайшую следующую дату
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// читаем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, errorResponse{Error: "bad json"})
		return
	}

	// обязательный title
	if task.Title == "" {
		writeJSON(w, errorResponse{Error: "empty title"})
		return
	}

	// проверяем/нормализуем дату + правило повторения
	if err := checkDate(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// добавляем в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// успех — возвращаем id как строку
	writeJSON(w, idResponse{ID: strconv.FormatInt(id, 10)})
}
