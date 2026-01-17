package api

import (
	"database/sql"
	"net/http"
)

const webRoot string = "./web"

func Init(mux *http.ServeMux, db *sql.DB) {
	// root folder
	//logger.Info.Println(os.Getwd())
	mux.Handle("GET /", http.FileServer(http.Dir(webRoot)))
	// API handlers
	// Обработчики задачи
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", manageTaskHandler(db))
	mux.HandleFunc("PUT /api/task", manageTaskHandler(db))
	mux.HandleFunc("GET /api/task", getTaskHandler(db))
	mux.HandleFunc("DELETE /api/task", deleteTaskHandler(db))
	mux.HandleFunc("POST /api/task/done", doneTaskHandler(db))
	// Получение списка задач
	mux.HandleFunc("GET /api/tasks", getTaskListHandler(db))
}
