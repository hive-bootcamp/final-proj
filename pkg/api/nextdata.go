package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return y1 > y2 || (y1 == y2 && (m1 > m2 || (m1 == m2 && d1 > d2)))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("bad date format")
	}

	parts := strings.Split(strings.TrimSpace(repeat), " ")
	rule := parts[0]

	// d N
	if rule == "d" {
		if len(parts) != 2 {
			return "", errors.New("bad d rule")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", errors.New("bad d interval")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}
	}

	// y
	if rule == "y" {
		if len(parts) != 1 {
			return "", errors.New("bad y rule")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}
	}

	return "", errors.New("unsupported rule")
}

// HTTP handler
func nextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")
	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "bad now", 400)
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	res, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprint(w, res)
}
