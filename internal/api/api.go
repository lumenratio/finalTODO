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
	// auth обработчик
	mux.HandleFunc("POST /api/signin", authGenTokenHandler)
	// Обработчики задачи
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", auth(manageTaskHandler(db)))
	mux.HandleFunc("PUT /api/task", auth(manageTaskHandler(db)))
	mux.HandleFunc("GET /api/task", auth(getTaskHandler(db)))
	mux.HandleFunc("DELETE /api/task", auth(deleteTaskHandler(db)))
	mux.HandleFunc("POST /api/task/done", auth(doneTaskHandler(db)))
	// Получение списка задач
	mux.HandleFunc("GET /api/tasks", auth(getTaskListHandler(db)))
}
