package server

import (
	"final-proj/pkg/api"
	"final-proj/pkg/db"
	"fmt"
	"net/http"
	"os"
)

func Run() error {
	port := 7540

	if env := os.Getenv("TODO_PORT"); env != "" {
		fmt.Sscanf(env, "%d", &port)
	}

	// Инициализация базы (ВАЖНО!)
	if err := db.Init("scheduler.db"); err != nil {
		return err
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}
